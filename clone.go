package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func collectionClone(localName string) {
	userAccount := new(accountProperties)
	remoteName, _ := userAccount.getRemoteName(localName)
	collectionJSONFileName := remoteName+".json"

	err := os.Mkdir(localName, os.ModePerm)
	if err == nil {
		err = os.Chdir(localName)
	}
	
	if err != nil {
		panic(err)
	}
	

	// Todo: Update so that path is not hard coded.
	remoteCollectionURL := "https://pithub.blob.core.windows.net/nvm4zqwm/"+collectionJSONFileName

	// Todo: Make download file name temp.
	err = DownloadFile(collectionJSONFileName, remoteCollectionURL)
	// Todo: Implement better error handling. 
	if err != nil {
		panic(err)
	}

	fmt.Println("Collection Description Downloaded")
	copyFile(collectionJSONFileName, pitFileName)

	// Todo: Implement collectionRead that takes the pathName of the pit file.
	props, err := collectionRead()
	check(err)
	
	// Download all files in the collection.
	for _, doc := range props.Documents {
		_, remoteFileURL := getRemoteFileNameAndURL(props, doc.NameLocal)

		fmt.Printf(remoteFileURL+"\n")
		err := DownloadFile(doc.NameLocal, remoteFileURL)
		// Todo: Implement better error handling. 
		if err != nil {
			panic(err)
		}
	}
}

// Code initially taken from "https://golangcode.com/download-a-file-with-progress/"

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	// fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	fmt.Printf("\rDownloading... %d complete", wc.Total)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}