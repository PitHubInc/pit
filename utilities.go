package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
	"os/user"
)

var lowerCaseLettersAndNumbers = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func generateChecksum(s string) string {
	// Add a last character that is a checksum of the other characters
	var total int = 0
	var i int = 0
	for i < (len(s) - 1) {
		total = total + int(s[i])
		i++
	}

	var modTotal = total % len(lowerCaseLettersAndNumbers)
	return string(lowerCaseLettersAndNumbers[modTotal])
}

func validateRandomFileName(fileName string) bool {
	fmt.Println(fileName)
	r := []rune(fileName)

	var s string = string(r[:len(r)-1])
	fmt.Println(s)
	s = s + generateChecksum(s)
	fmt.Println(s)

	if fileName == s {
		fmt.Println("Valid")
	} else {
		fmt.Println("Not valid")
	}
	return false
}

func randomFileName(nameLength int) string {
	var lowerCaseLetters = []rune("abcdefghijklmnopqrstuvwxyz")

	rand.Seed(time.Now().UnixNano())

	// Create a random file name
	randFileName := make([]rune, (nameLength - 1))
	for i := range randFileName {
		if i == 0 {
			randFileName[i] = lowerCaseLetters[rand.Intn(len(lowerCaseLetters))]
		} else {
			randFileName[i] = lowerCaseLettersAndNumbers[rand.Intn(len(lowerCaseLettersAndNumbers))]
		}
	}

	return string(randFileName) + generateChecksum(string(randFileName))
}

func ensureDirectory(dirName string) bool {
	err := os.Mkdir(dirName, os.ModePerm)
	if err == nil || os.IsExist(err) {
		return true
	}

	return false
}

func fileOrDirectoryExists(fileOrDirectoryName string) bool {
	info, err := os.Stat(fileOrDirectoryName)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func md5File(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func copyFile(sourceFileName string, destinationFileName string) bool {
	var returnValue bool = true
	sourceFile, err := os.Open(sourceFileName)
	if err != nil {
		returnValue = false
		log.Fatal(err)
	}
	defer sourceFile.Close()

	// Create new file
	newFile, err := os.Create(destinationFileName)
	if err != nil {
		returnValue = false
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		returnValue = false
		log.Fatal(err)
	}

	return returnValue
}

func deleteFile(fileName string) {
	_ = os.Remove(fileName)
}

func getHomeFolderName() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	return usr.HomeDir
}

