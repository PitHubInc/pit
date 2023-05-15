package main

// BugBug: Why is the following line included?
// const pitProductionContainerName = "nvm4zqwm"

// Warning: Be VERY careful about changing the constants, types, or even the names of the variable 
// below as any changes will likely be a breaking change for existing versions. Note that the
// variable names are utilized by Go in the creation of JSON files. 

const pitFileName = "pit.json"
const pitMD5tag = "pitmd5"
const pitSeparator = "-"
const userAppFolderName = ".pit"
const userAppLogFileName = "log.txt"
const userAppAccountFileName = "account.json"

// Note that the Go json.Marshal() function only exports fields that start with an upper case name. 
type collectionProperties struct {
	NameLocal      string
	NameRemote     string
	Created        string
	Updated        string
	ETag           string
	URL            string
	Documents      []documentProperties
}

type documentProperties struct {
	NameLocal      string
	ETag		   string
	MD5            string
	PreviousMD5s   []string
}