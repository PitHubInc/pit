package main

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Edit "launch.json" to set execution folder using "cwd" field.

func pitTest() {
	log.Println("test()")
	log.Printf("version: %s\n", productVersion)
	log.Println(fmt.Sprintf("os.Args=%s", os.Args))
	log.Println(fmt.Sprintf("containerName=%s", getContainerName()))

	// fmt.Println("Test Executing: setTestEnv()")
	// setTestEnv()

	// testCase00()

	version()
	collectionClone()
}

func testCase00() {
	// Initialize container.
	log.Println("testCase00()")
	fmt.Println(fmt.Sprintf("Test Executing: createPublicContainer(%s)", getContainerName()))
	createPublicContainer(getContainerName())
}

func testCase01() {
	// Add documents and push documents.
	log.Println("testCase01()")

	os.Chdir("/Users/eric/Collections/test-case-01")
	path, _ := os.Getwd()
	fmt.Println(fmt.Sprintf("Test Current Folder: %s", path))

	fmt.Println("Test Executing: collectionDeleteLocalIfExist()")
	collectionDeleteLocalIfExist()

	// No collection initialize.
	fmt.Println("Test Executing: pit status")
	collectionStatus()

	fmt.Println("Test Executing: pit init")
	collectionInitialize()

	fmt.Println("Test Executing: pit add syllabus.docx")
	collectionAdd("syllabus.docx")

	fmt.Println("Test Executing: pit status")
	collectionStatus()

	fmt.Println("Test Executing: pit add syllabus.pdf")
	collectionAdd("syllabus.pdf")

	fmt.Println("Test Executing: pit status")
	collectionStatus()

	fmt.Println("Test Executing: pit push")
	collectionPush()

	fmt.Println("Test Executing: pit status")
	collectionStatus()
}

func testCase03() {
	testCase99()
	testCase00()
	testCase01()
}

func testCase99() {
	fmt.Println(fmt.Sprintf("Test Executing: deleteContainer(%s)", getContainerName()))
	err := deleteContainer(getContainerName())
	check(err)

	fmt.Println("Test Executing: Sleep for 5 minutes to give Azure time to delete the container...")
	time.Sleep(5 * 60 * time.Second)
	fmt.Println("Text Executing: Done sleeping")
}

func validate(arg1 string, arg2 string) {
	log.Println(md5File(arg2))

	log.Println("New test code ejp")
	log.Println(len(md5File(arg2)))

	// Todo: log.Println("MD5 Length="+strconv.Itoa(int(len(md5File(arg2)))))

	var fileExtension string = filepath.Ext(arg2)

	log.Println("Random file name with md5: ")
	var destinationFileName = randomFileName(8) + md5File(arg2) + fileExtension
	log.Println(destinationFileName)

	log.Println(randomFileName(128)) 
	log.Println(randomFileName(64))
	log.Println(randomFileName(32))
	log.Println(randomFileName(16))
	log.Println(randomFileName(8))
	log.Println(randomFileName(4))

	validateRandomFileName(randomFileName(8))

	log.Println("The following test should (generally) fail.")
	validateRandomFileName(randomFileName(7) + "j")

	var destinationFolderName string = getHomeFolderName()
	if ensureDirectory(destinationFolderName) {
		log.Println("Pit Directory OK: " + destinationFolderName)
	}

	var destination string = destinationFolderName + destinationFileName
	log.Println("destination: " + destination)

	copyFile(arg2, destination)
}



