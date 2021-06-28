/*
	Account resents ownership and security. An Account also handles the creating and deleting of
	Containers. It is modeled after an Azure Storage Account. 

	References: https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blobs-introduction
*/
package main
import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func handleAzureErrors(err error) {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
			log.Println(serr.ServiceCode())
			switch serr.ServiceCode() { // Compare serviceCode to ServiceCodeXxx constants
			case azblob.ServiceCodeContainerAlreadyExists:
				log.Println("Verified container exists")
				return
			case "ContainerNotFound":
				log.Println("Container not found")
			}
		}
		log.Fatal(err)
	}
}

func getAccountNameAndKey() (string, string, error) {
	// From the Azure portal, get storage account name and key and set environment variables.
	//     export AZURE_STORAGE_ACCOUNT="pithub"
	//     export AZURE_STORAGE_ACCESS_KEY="2q62fVoYfT6ZOudTALzXBSz6eKOXh4CRgpMfuMWpyRFlUh/QB+K3IpaAm/hAUjrbMoZN9t0Cbl4lYHMU3lV89A=="

	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		account := new(accountProperties)
		container, err := account.defaultContainer()
		if err != nil {
			return accountName, accountKey, err
		}
		
		accountName = container.Account
		accountKey  = container.Key
	}
	
	if len(accountName) == 0 || len(accountKey) == 0 {
		err := errors.New("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
		return accountName, accountKey, err
	}
	
	return accountName, accountKey, nil
}

func getAccountNameKeyAndCredential() (string, string, azblob.Credential, error) {
	var err error = nil
	var accountCredential azblob.Credential = nil

	accountName, accountKey, err := getAccountNameAndKey()
	if err != nil {
		return accountName, accountKey, accountCredential, err
	}

	accountCredential, err = azblob.NewSharedKeyCredential(accountName, accountKey)
	return accountName, accountKey, accountCredential, err
}

func getContainerURL(containerName string) (azblob.ContainerURL, error) {
	accountName, _, accountCredential, err := getAccountNameKeyAndCredential()
	p := azblob.NewPipeline(accountCredential, azblob.PipelineOptions{})

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)) 

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)
	return containerURL, err
}

func createPublicContainer(containerName string) {
	log.Println(fmt.Sprintf("createPublicContainer(containerName=%s)", containerName))
	containerURL, err := getContainerURL(containerName)
	handleAzureErrors(err)
	
	ctx := context.Background()
	_, err2 := containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessBlob)
	// Consider other options for private or fully public containers including:
	//     azblob.PublicAccessContainer 
	//     azblob.PublicAccessNone 
	handleAzureErrors(err2)
}

func deleteContainer(containerName string) error {
	log.Println(fmt.Sprintf("deleteContainer(containerName=%s)", containerName))
	if !strings.Contains(containerName, "testtest") {
		log.Println(fmt.Sprintf("Fatal Error: Attempting to delete a non-test Containter named %s", containerName))
		return errors.New("The deleteContainer function only deletes containers with \"testtest\" in the name")
	}

	containerURL, _ := getContainerURL(containerName)
	ctx := context.Background() 
	containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	return nil
}

