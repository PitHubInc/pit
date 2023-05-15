/*
	Collection represents a logical grouping of Documents. Many collections can be stored in a Container.
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func collectionExists() bool {
	// Check if collection file exists.
	info, err := os.Stat(pitFileName)
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		// It would be very strange that someone created a directory with the same name as our collection file.
		fmt.Println(fmt.Sprintf("Fatal: A directory exists at %s", pitFileName))
		return false
	}

	return true
}

func collectionInitialize() {
	if collectionExists() {
		fmt.Println("Collection already initialized")
	} else {
		var props collectionProperties

		props.NameLocal = collectionGetInitialLocalName()
		props.NameRemote = collectionNewRemoteName()
		fmt.Printf("Initializing \"%s\" as \"%s\"\n", props.NameLocal, props.NameRemote)

		t := time.Now()
		props.Created = t.Format(time.RFC3339)
		props.Updated = t.Format(time.RFC3339)

		err := collectionWrite(props)
		check(err)
	}
}

func getRemoteNameAndURL(props collectionProperties, localFileName string) (string, string) {
	remoteFileName := props.NameRemote + pitSeparator + strings.ToLower(localFileName)
	remoteFileURL := fmt.Sprintf("https://pithub.blob.core.windows.net/%s/%s", getContainerName(), remoteFileName)
	return remoteFileName, remoteFileURL
}

func verifyCollectionDocument(props collectionProperties, doc documentProperties) error {
	// Verify local file exists.
	if !fileExists(doc.NameLocal) {
		return errors.New(fmt.Sprintf("%s Document has been deleted\n", doc.NameLocal))
	}

	// Verify local document has not been modified since is was last added.
	localFileMD5 := md5File(doc.NameLocal)
	if doc.MD5 != localFileMD5 {
		// Document has been updated locally but updated version has not been added. 
		fmt.Printf("%s updated but version has not been added with \"pit add %s\"\n", 
			padRight(doc.NameLocal, " ", 20), doc.NameLocal)
		return nil
	}

	// Verify remote document.
	remoteFileName, remoteFileURL := getRemoteNameAndURL(props, doc.NameLocal)
	remoteFileMD5, _, err := getRemoteFileMD5AndETag(remoteFileName)
	if err != nil {
		// Remote file not found likely because it has not been pushed. 
		fmt.Printf("%s added but not published with \"pit push\"\n", padRight(doc.NameLocal, " ", 20))
		return nil
	}
	if doc.MD5 != remoteFileMD5 {
		// Local and remote files are not the same like because the local file has been updated
		return nil
	}

	remotefileMD5, remoteFileETag, err := getRemoteFileMD5AndETag(remoteFileName)
	if doc.MD5 != remotefileMD5 {
		errorString := fmt.Sprintf("Fatal Error: %s MD5 has unexpected value\n", remoteFileURL)
		errorString += fmt.Sprintf("  Current cloud MD5:  %s\n", remoteFileMD5)
		errorString += fmt.Sprintf("  Local MD5: %s", doc.ETag)
		return errors.New(errorString)
	}

	if doc.ETag != remoteFileETag {
		errorString := fmt.Sprintf("Fatal Error: %s has been updated by another computer\n", remoteFileURL)
		errorString += fmt.Sprintf("  Local ETag: %s\n", doc.ETag)
		errorString += fmt.Sprintf("  Cloud ETag: %s", remoteFileETag)
		return errors.New(errorString)
	}

	fmt.Printf("%s verified and shared as %s\n", padRight(doc.NameLocal, " ", 20), remoteFileURL)
	return nil
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func verifyCollectionDocuments(props collectionProperties) {
	// Verify each document in collection.
	for _, doc := range props.Documents {
		err := verifyCollectionDocument(props, doc)
		if err != nil {
			fmt.Printf("%s\n", err)
		} 
	}
}

func collectionStatus() {
	if !collectionExists() {
		fmt.Printf("No collection initialzed\n")
	} else {
		props, err := collectionRead()
		check(err)

		fmt.Printf("Collection \"%s\"\n", props.NameLocal)

		verifyCollectionDocuments(props)
	}
}

func collectionRead() (collectionProperties, error) {
	var props collectionProperties
	if !collectionExists() {
		return props, errors.New("No collection initialized")
	}

	collectionFileData, err := ioutil.ReadFile(pitFileName)
	check(err)

	err = json.Unmarshal(collectionFileData, &props)
	check(err)

	return props, nil
}

func collectionWrite(props collectionProperties) error {
	collectionJSON, err := json.MarshalIndent(props, "", "    ")
	// Non-formatted JSON would be generated with the following statement:
	//     collectionJSON, err := json.Marshal(collectionProps)
	check(err)

	err = ioutil.WriteFile(pitFileName, collectionJSON, 0644)
	// The magic "0644" parameter above is filemode. Additional information can be found at:
	// https://golang.org/pkg/os/#FileMode

	return err
}


func collectionAddOrUpdate(filePathAndName string) error {
	if !fileExists(filePathAndName) {
		return errors.New(fmt.Sprintf("File \"%s\"does not exist", filePathAndName))
	}

	props, err := collectionRead()
	check(err)

	var doc documentProperties
	doc.NameLocal = filepath.Base(filePathAndName)
	doc.MD5 = md5File(filePathAndName)

	for _, element := range props.Documents {
		if doc.NameLocal == element.NameLocal {
			if doc.MD5 == element.MD5 {
				fmt.Printf("%s is already in the Collection and up to date\n", padRight(doc.NameLocal, " ", 20))
				return nil
			} else if doc.MD5 != element.MD5 {
				err = collectionUpdate(filePathAndName)
				return err
			}
		} 
	}

	// File was not in Collection so it needs to be added.
	props.Documents = append(props.Documents, doc)
	return collectionWrite(props)
}

func collectionUpdate(filePathAndName string) error {
	props, err := collectionRead()
	check(err)

	for i := 0; i < len(props.Documents); i++ {
		fileName := filepath.Base(filePathAndName)
		if props.Documents[i].NameLocal == fileName {
			currentMD5 := md5File(filePathAndName)
			originalMD5 := props.Documents[i].MD5

			if (currentMD5 != originalMD5) {
				props.Documents[i].PreviousMD5s = append(props.Documents[i].PreviousMD5s, originalMD5)
				props.Documents[i].MD5 = currentMD5
			}
		}
	}

	return collectionWrite(props)
}

func collectionAdd(filePathAndName string) error {
	if !fileExists(filePathAndName) {
		return errors.New(fmt.Sprintf("File \"%s\"does not exist", filePathAndName))
	}

	props, err := collectionRead()
	check(err)

	var doc documentProperties
	doc.NameLocal = filepath.Base(filePathAndName)
	doc.MD5 = md5File(filePathAndName)

	for _, element := range props.Documents {
		if doc.NameLocal == element.NameLocal {
			if doc.MD5 == element.MD5 {
				fmt.Printf("%s is already in the Collection and up to date\n", padRight(doc.NameLocal, " ", 20))
				return nil			
			} else {
				return collectionUpdate(filePathAndName)
			}
		}
	}

	props.Documents = append(props.Documents, doc)
	err = collectionWrite(props)
	check(err)
	return nil
}

func collectionPropsPrintJSON(props collectionProperties) {
	// Optionally store non-formatted json by utilizing:
	//     collectionJSON, err := json.Marshal(collectionProps)
	json, err := json.MarshalIndent(props, "", "    ")
	check(err)

	jsonString := string(json)
	fmt.Println(jsonString)
}

func collectionGetInitialLocalName() string {
	// Set the initial CollectionNameLocal to the current directory name.
	path, err := os.Getwd()
	check(err)
	return filepath.Base(path)
}

func collectionNewRemoteName() string {
	return randomFileName(8)
}

func collectionDeleteLocalIfExist() {
	if collectionExists() {
		_ = os.Remove(pitFileName)
	}
}

func collectionPush() {
	props, err := collectionRead()

	if (err != nil) { 
		fmt.Printf("Error: unable to read collection\n%s", err)
		return
	}

	// If the local Collection json file  is updated, we will need to upload it at the end of the function. 
	collectFileModified := false

	// Check each Document in the Collection to see if it need to be uploaded. 
	for index, doc := range props.Documents {
		newFile := false
		updatedFile := false
		uploadFile := false

		remoteFileName := props.NameRemote + pitSeparator + strings.ToLower(doc.NameLocal)
		remoteMDd5, remoteETag, err := getRemoteFileMD5AndETag(remoteFileName)
		if err != nil {
			if strings.Contains(err.Error(), "ServiceCode=BlobNotFound") {
				// Document exists locally, but not remotely.
				// Todo: Implement a better way to handle this error so that we are not comparing strings.
				newFile = true
				uploadFile = true
			} else {
				// Document exist remotely, be we are not able to get the remote file MD5 and ETag.
				fmt.Printf("%s Error: %s\n", padRight(doc.NameLocal, " ", 20), err.Error())
				break // Stop processing this document.
			}
		}

		if !newFile {
			if doc.ETag != remoteETag || doc.MD5 != remoteMDd5  {
				// Document exists locally and remotely, but has been updated locally.
				updatedFile = true
			}
		}

		_, remoteFileURL := getRemoteNameAndURL(props, doc.NameLocal)
		if !newFile && !updatedFile {
			// The Document exists locally and remotely and the files are the same (matching eTags and MD5s). 
			fmt.Printf("%s verified and shared as %s\n", padRight(doc.NameLocal, " ", 20), remoteFileURL)
		} 

		if updatedFile {
			// Check if the remote file was a previously uploaded file from this computer. 
			for _, previousLocalMD5 := range props.Documents[index].PreviousMD5s {
				if remoteMDd5 == previousLocalMD5  {
					// One of the previous MD5s match so it should be safe to upload the updated file.
					uploadFile = true
					break
				}
			}

			if !uploadFile {
				// The remote Document has an MD5 that is not recognized in the local Collection. Therefore, the Document 
				// was likely updated from a different computer and we risk overwriting changes. 
				fmt.Printf("%s Error: aborting upload due to version conflict\n", padRight(doc.NameLocal, " ", 20))
				break
			}
		}

		if uploadFile {
			collectFileModified = true

			uploadDocument(getContainerName(), doc.NameLocal, remoteFileName)
			setDocumentMetadataMD5(getContainerName(), remoteFileName, doc.MD5)
			
			newRemoteMD5, newRemoteETag, err := getRemoteFileMD5AndETag(remoteFileName)
			if err != nil {
				fmt.Printf("%s Error: unable to abtain MD5 or ETag for %s\n", padRight(doc.NameLocal, " ", 20), remoteFileURL)
				break
			}

			if doc.MD5 != newRemoteMD5 {
				fmt.Printf("%s Error: MD5 not updated correctly for for %s\n", padRight(doc.NameLocal, " ", 20), remoteFileURL)
			}

			props.Documents[index].ETag = newRemoteETag
		}
	}

	if collectFileModified {
		// Write the local Collection if was modified (e.g. ETag). 
		err = collectionWrite(props)
		if err != nil {
			fmt.Printf("Error: unable to updated Collection\n")
		}

		// Upload (overwrite if necessary) the modified Collection json file.
		collectionLocalFileName := pitFileName
		collectionRemoteFileName := props.NameRemote + ".json"
		uploadDocument(getContainerName(), collectionLocalFileName, collectionRemoteFileName)
	}
}

func collectionPushBK1() {
	props, err := collectionRead()
	check(err)

	// Todo: Refactor code so that we do not get ALL of the remote files to see if a given file exists and is current.
	documentNames := getDocumentNames(getContainerName())
	for docIndex, doc := range props.Documents {
		remoteFileName := props.NameRemote + pitSeparator + strings.ToLower(doc.NameLocal)
		remoteFileExists := false
		for _, documentName := range documentNames {
			// Todo: consider optimizing by breaking out of loop when file is found.
			if strings.EqualFold(remoteFileName, documentName) {
				remoteFileExists = true
				log.Println(fmt.Sprintf("Verified: %s", doc.NameLocal))

				// BugBug: This does not look right as it references 'ovun'.
				log.Println(fmt.Sprintf("    Is: %s%s", "https://pithub.blob.core.windows.net/ovun/", remoteFileName))
			}
		}

		if !remoteFileExists {
			uploadDocument(getContainerName(), doc.NameLocal, remoteFileName)
			setDocumentMetadataMD5(getContainerName(), remoteFileName, doc.MD5)
			// Todo: Verify MD5 and set local ETag.

			_, remoteFileETag, err := getRemoteFileMD5AndETag(remoteFileName)
			check(err)

			props.Documents[docIndex].ETag = remoteFileETag
		}
	}

	// Write the local collection as it may have been modified above (e.g. ETag). 
	err = collectionWrite(props)
	check(err)

	// Upload (overwrite if necessary) the collection json file.
	collectionLocalFileName := pitFileName
	collectionRemoteFileName := props.NameRemote + ".json"
	uploadDocument(getContainerName(), collectionLocalFileName, collectionRemoteFileName)
}

/* Todo: Consider if we need to check local file that are not in the Collection as part of the status feature.
func fileInCollection(props collectionProperties, fileName string) bool {
	for _, doc := range props.Documents {
		if doc.NameLocal == fileName {
			return true
		}
	}
	return false
}

func verifyLocalFileInCollection(fileName string) error {
	// Todo: Implement function.
	return nil
}

func collectionCheckDocuments(props collectionProperties) {
	// Verify all local files are in collection.
	var files []string

	root := "./"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	check(err)

	for _, file := range files {
		info, err := os.Stat(file)
		check(err)
		if !info.IsDir() {
			if !fileInCollection(props, file) {
				if (file != pitFileName) && (file != ".DS_Store") {
					fmt.Printf("\"%s\" not in Collection\n", file)
				}
			}
		}
	}

	// Verify remote copy of collection
	remoteCollectionFileName := props.NameRemote + ".json"
	fmt.Printf("ContainerName=\"%s\" remoteCollectionName=\"%s\n", getContainerName(), remoteCollectionFileName)

	// downloadDocument(containerName, pitFileName, remoteCollectionFileName)
}
*/