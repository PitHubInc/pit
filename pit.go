package main

import (
	"fmt"
	"log"
	"os"

	"io/ioutil"
)

func check(err error) {
	if err != nil {
		fmt.Printf("Fatal Error: %s\n", err)
		failureExitCode := 1
		os.Exit(failureExitCode)
	}
}

func main() {
	verifyUserAccountAndLogFile()

	// At least two arguments are needed for any command.
	if len(os.Args) < 2 {
		help()
		os.Exit(0)
	}

	var secondArg string = os.Args[1]
	if (secondArg == "help") || (secondArg == "-help") || (secondArg == "-h") {
		help()
	} else if (secondArg == "version") || (secondArg == "-version") || (secondArg == "-v") {
		version()
	} else if (secondArg == "status") || (secondArg == "-status") || (secondArg == "-s") {
		status()
	} else if (secondArg == "init") || (secondArg == "-init") || (secondArg == "-i") {
		initialize()
	} else if (secondArg == "add") || (secondArg == "-add") || (secondArg == "-a") {
		if len(os.Args) < 3 {
			log.Println("Error: '-add' must include a [[document-name]] argument")
			os.Exit(0)
		} else {
			add(os.Args[2])
		}
	} else if (secondArg == "push") || (secondArg == "-push") || (secondArg == "-p") {
		push()
	} else if (secondArg == "test") || (secondArg == "-test") || (secondArg == "-t") {
		test()
	} else if (secondArg == "setproduction") || (secondArg == "-setproduction") || (secondArg == "-sp") {
		setProductionEnv()
	} else if (secondArg == "clone") {
		collectionClone()
	} else if (secondArg == "settest") || (secondArg == "-settest") || (secondArg == "-st") {
		setTestEnv()
	} else {
		introduction()
	}

	log.Println("Exiting Successfully")
}

func verifyUserAccountAndLogFile() {
	account := new(accountProperties)
	err := account.verify()
	check(err)

	// If the log file exists, append log messages to the file.
	logfilePathAndName := account.userAppPath() + string(os.PathSeparator) + userAppLogFileName
	logFile, err := os.OpenFile(logfilePathAndName, os.O_APPEND|os.O_RDWR, 0666)
	if err == nil {
		log.SetOutput(logFile) // For stderr logging utilize log.SetOutput(os.Stderr).

	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func introduction() {
	fmt.Println(
		`
The Pit application provides simple functionality that allows 
you to share your videos online. 

Common Pit commands init:
    init      Create a Pit collection
    add       Add or update a Pit collection document
    push      Copy all new or updated documents so they can be viewed online
    status    View the current status of the collection including document URLs
    help      View more detailed information about Pit functionality`)
}

func help() {
	fmt.Println(
		`
Example Usage:
    pit init
    pit add [[document-name]]
    pit push
    pit status
    pit help
    pit version`)
}

func version() {
	v := fmt.Sprintf("version/build: %s", productVersion)
	fmt.Printf(v+"\n")
	log.Println(v+" reported")
}

func status() {
	collectionStatus()
}

func add(documentName string) {
	fmt.Printf("Add \"%s\"\n", documentName)
	err := collectionAdd(documentName)
	check(err)
}

func initialize() {
	fmt.Printf("Initialize:\n")
	collectionInitialize()
}

func push() {
	fmt.Printf("Push:\n")
	collectionPush()
}

func test() {
	pitTest()
}

func setProductionEnv() {
	account := new(accountProperties)
	account.defaultAccountProperties("production")
}

func setTestEnv() {
	account := new(accountProperties)
	account.defaultAccountProperties("test")
}

// Todo: Consider adding header:
//     fmt.Println("\nBestOfTheBest Document Sharing")
//     fmt.Println("\nBestOfTheBest Video Management & Sharing")
//     fmt.Println("\nBestOfTheBest Document Management & Sharing")
