// User encapsulates user account information including Container and basic Collection information. User account
// information is NOT the same as storage account information. The associated JSON file is storied in the
// %HOME%/.pit/account.json.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

// Warning: Be VERY careful about changing the constants, types, or even the names of the variables directly below as
// any changes will likely be a breaking change for existing versions. Note that the variable names are utilized by Go
// in the creation of JSON files and that the Go json.Marshal() function only exports fields that start with an upper
// case name.
type accountProperties struct {
	Description string // Value:   accountFileDescription
	Email       string // Example: "epogue@epogue.com"
	Containers  []containerProperties
	Collections []basicCollectionProperties
}

type containerProperties struct {
	Type    string // Example: "azure"
	Account string // Example: "pithub"
	Key     string // Example: "dGs3xAXFgM7bwJN9GtU0HaahFzqa77rWU/TWl8Oryqon93V28sexQ80V8V6PNedgEMhVu3C2eEBGcWUtFDcEUA=="
	Name    string // Example: "nvm4zqwmtesttest"
	URL     string // Example: "https://pithub.blob.core.windows.net/nvm4zqwmtesttest/"
	Default string // Values:  "yes" or "no"
}

type basicCollectionProperties struct {
	NameLocal  string
	NameRemote string
}

// Todo: Consider removing this field from accountProperties
// const accountFileDescription = "User Account Information including Container and basic Collection information"
const accountFileDescription = "Account and Container Properties"

// End Warning

// User account methods.
func (ap *accountProperties) userAppPath() string {
	user, _ := user.Current()
	path := user.HomeDir + string(os.PathSeparator) + userAppFolderName
	return path
}

func (ap *accountProperties) filePathAndName() string {
	return ap.userAppPath() + string(os.PathSeparator) + userAppAccountFileName
}

func (ap *accountProperties) exists() bool {
	ap.filePathAndName()

	info, err := os.Stat(ap.filePathAndName())
	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		log.Println(fmt.Sprintf("Fatal Error: A directory exists at %s", ap.filePathAndName()))
		return false
	}

	return true
}

func (ap *accountProperties) verify() error {
	var account = new(accountProperties)
	if !account.exists() {
		return account.initialize()
	}

	return account.read()
}

func (ap *accountProperties) write() error {
	accountJSON, err := json.MarshalIndent(ap, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ap.filePathAndName(), accountJSON, 0644) // The "0644" parameter is filemode.
	return err
}

func (ap *accountProperties) initialize() error {
	if ap.exists() {
		// The exists() method should have been called previously.
		return errors.New(fmt.Sprintf("Account already exists at %s", ap.filePathAndName()))
	} else {
		_, err := os.Stat(ap.userAppPath())
		if os.IsNotExist(err) {
			os.Mkdir(ap.userAppPath(), os.ModePerm)
		}

		ap.Description = accountFileDescription
		err = ap.write()
		if err != nil {
			return err
		}
	}

	return nil
}

func (ap *accountProperties) read() error {
	if !ap.exists() {
		// The initialize() method should have already been called previously.
		return errors.New("Fatal Error: Account file does not exist.")
	}

	accountFileData, err := ioutil.ReadFile(ap.filePathAndName())
	if err != nil {
		return err
	}

	err = json.Unmarshal(accountFileData, &ap)
	if err != nil {
		return err
	}

	if ap.Description != accountFileDescription {
		return errors.New("Fatal Error: Account file not valid.")
	}

	return nil
}

func (ap *accountProperties) defaultContainer() (containerProperties, error) {
	// Assumes verify() called previously.
	err := ap.read()
	if err != nil {
		return ap.Containers[0], err
	}

	if ap.Containers == nil {
		var container containerProperties
		return container, err
	}

	return ap.Containers[0], nil
}

func (ap *accountProperties) defaultAccountProperties(environment string) {
	// BugBug: Remove default account and key and email!
	var container containerProperties
	container.Type = "azure"
	container.Account = "pithub"
	container.Key = "2q62fVoYfT6ZOudTALzXBSz6eKOXh4CRgpMfuMWpyRFlUh/QB+K3IpaAm/hAUjrbMoZN9t0Cbl4lYHMU3lV89A=="
	container.Name = "nvm4zqwmtesttest"
	container.URL = "https://pithub.blob.core.windows.net/nvm4zqwmtesttest/"
	container.Default = "yes"

	if environment == "production" {
		fmt.Println("Warning: **PRODUCTION** environment enabled!")
		container.Name = "nvm4zqwm"
	} else {
		fmt.Println("**TEST** environment enabled.")
	}

	ap.Description = accountFileDescription
	ap.Email = "epogue@epogue.com"
	ap.Containers = append(ap.Containers, container)

	err := ap.write()
	if err != nil {
		log.Println("Fatal Error: Unable to write account properties")
	}
}

func (ap *accountProperties) addBasicCollectionInfo(nameLocal string, nameRemote string) {
	var basicCollection basicCollectionProperties
	basicCollection.NameLocal = nameLocal
	basicCollection.NameRemote = nameRemote

	err := ap.read()
	check(err)

	ap.Collections = append(ap.Collections, basicCollection)
	ap.write()
}
