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
	check(err)
	
	err = os.Chdir(localName)
	check(err)
	
	// Todo: Update so that path is not hard coded.
	remoteCollectionURL := "https://pithub.blob.core.windows.net/nvm4zqwm/"+collectionJSONFileName

	err = DownloadFile(pitFileName, remoteCollectionURL)
	check(err)

	props, err := collectionRead()
	check(err)
	
	// Download all files in the collection.
	for _, doc := range props.Documents {
		_, remoteFileURL := getRemoteFileNameAndURL(props, doc.NameLocal)

		err := DownloadFile(doc.NameLocal, remoteFileURL)
		check(err)
	}
}

// The code below was initially taken from "https://golangcode.com/download-a-file-with-progress/"

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

var downloadFileName string = ""
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Return again and print current status of download
	fmt.Printf("\r%s...%d downloaded (kb)", downloadFileName, wc.Total/1000)
}

func DownloadFile(filepath string, url string) error {
	downloadFileName = filepath

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
	out.Close()

	// Extra spaces are to make sure that everything written previously to the line is deleted.
	fmt.Printf("\r%s...complete                   \n", downloadFileName)

	err = os.Rename(filepath+".tmp", filepath)
	check(err)

	return nil
}