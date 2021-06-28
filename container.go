/*
	Container represents a place where documents can be stored. It is initially modeled after a MS Azure Blob Storage
	Container.

	Additional references: https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blobs-introduction
	https://www.eventslooped.com/posts/use-golang-to-upload-files-to-azure-blob-storage/
	https://github.com/inemtsev/go_blob_uploader/blob/master/main.go
	https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go
	https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#example-Metadata--Blobs
	https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-properties
	https://github.com/Azure/azure-storage-blob-go/blob/master/azblob/zt_examples_test.go
	https://techcommunity.microsoft.com/t5/azure-developer-community-blog/build-and-deploy-your-first-app-with-the-azure-sdk-for-go-on/ba-p/1172973
	https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-go?tabs=linux
	https://github.com/inemtsev/go_blob_uploader/blob/master/main.go
	https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#example-BlockBlobURL

	Azure Go Authentication:
	https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization 

	Azure Go Resourced:
	https://docs.microsoft.com/en-us/azure/developer/go/
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func getContainerName() string {
	containerName := os.Getenv("AZURE_CONTAINER_NAME")
	if len(containerName) != 0 {
		return containerName
	}

	account := new(accountProperties)
	container, err := account.defaultContainer()
	if err != nil {
		log.Println("Fatal Error: Unable to obtain Container name")
		return ""
	}	
	
	return container.Name
}

func getDocumentNames(containerName string) []string {
	containerURL, _ := getContainerURL(containerName)
	ctx := context.Background() // This example uses a never-expiring context

	var blobNames []string
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})

		// Check if ContainerNotFound
		if err != nil {
			if serr, ok := err.(azblob.StorageError); ok { // This error is a Service-specific
				if serr.ServiceCode() == "ContainerNotFound" {
					fmt.Printf("Container \"%s\" not found\n", containerName)
					return blobNames
				}
			}
		}
	
		handleAzureErrors(err)
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			blobNames = append(blobNames, blobInfo.Name)
		}
	}

	return blobNames
}

func uploadDocument(containerName string, localName string, remoteName string) {
	accountName, _, accountCredential, err := getAccountNameKeyAndCredential()
	p := azblob.NewPipeline(accountCredential, azblob.PipelineOptions{
		Retry: azblob.RetryOptions{
			TryTimeout: 120 * time.Minute,
		},
	})

	// Documentation notes relating to "TryTimeout: value:
	// https://github.com/Azure/azure-storage-blob-go/issues/60
	//
	// TryTimeout indicates the maximum time allowed for any single try of an HTTP request.
	// A value of zero means that you accept our default timeout. NOTE: When transferring large amounts
	// of data, the default TryTimeout will probably not be sufficient. You should override this value
	// based on the bandwidth available to the host machine and proximity to the Storage service. A good
	// starting point may be something like (60 seconds per MB of anticipated-payload-size).

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background() // This example uses a never-expiring context

	// Attempt to create a new container.
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	handleAzureErrors(err)

	if localName != pitFileName {
		fmt.Println(fmt.Sprintf("Uploading: %s...", localName))
	} else {
		log.Println(fmt.Sprintf("Uploading: %s...", localName))
	}

	// Todo: Consider accessing the local file directly instead of creating a copy.
	// Create a copy of the local file to upload.
	copyFile(localName, remoteName)

	// Upload blob.
	blobURL := containerURL.NewBlockBlobURL(remoteName)
	file, err := os.Open(remoteName)
	handleAzureErrors(err)

	// Set content types for known file extensions.
	ext := filepath.Ext(remoteName)
	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{},
	}

	if strings.EqualFold(ext, ".json") {
		o = azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: "application/json",
			},
		}
	} else if strings.EqualFold(ext, ".mp4") {
		o = azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: "video/mp4",
			},
		}
	} else if strings.EqualFold(ext, ".docx") {
		o = azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
		}
	} else if strings.EqualFold(ext, ".pdf") {
		o = azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: "application/pdf",
			},
		}
	} else if strings.EqualFold(ext, ".html") {
		o = azblob.UploadToBlockBlobOptions{
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: "text/html",
			},
		}
	}
	
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, o)
	handleAzureErrors(err)

	// Todo: Analyze if we should delete ctx by including the following statement?
	//     containerURL.Delete(ctx, azblob.ContainerAccessConditions{})

	// Delete the local file that we created above.
	deleteFile(remoteName)
}

/*
func downloadDocument(containerName string, localName string, remoteName string) {
	accountName, _, accountCredential, err:= getAccountNameKeyAndCredential()
	fmt.Println("1")
	fmt.Println(containerName)
	fmt.Println(localName)
	fmt.Println(remoteName)

	p := azblob.NewPipeline(accountCredential, azblob.PipelineOptions{})
	fmt.Println("2")
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	fmt.Println("3")
	containerURL := azblob.NewContainerURL(*URL, p)
	ctx := context.Background() // This example uses a never-expiring context
	fmt.Println("4")
	blobURL := containerURL.NewBlockBlobURL(remoteName)
	fmt.Println("5")
	// Here's how to download the blob
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	fmt.Println("6")
	handleAzureErrors(err)

	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	fmt.Printf("7")
	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	handleAzureErrors(err)
	fmt.Printf("8")
	// The downloaded blob data is in downloadData's buffer. :Let's print it
	fmt.Printf("Downloaded the blob: " + downloadedData.String())

	fmt.Printf("Cleaning up.\n")
	// Careful: containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
}
*/

func getFileBlobURL(containerName string, fileName string) azblob.BlobURL {
	accountName, _, accountCredential, err := getAccountNameKeyAndCredential()
	check(err)
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	p := azblob.NewPipeline(accountCredential, azblob.PipelineOptions{})
	containerURL := azblob.NewContainerURL(*URL, p)
	blobURL := containerURL.NewBlobURL(fileName)
	return blobURL
}

func setDocumentMetadataMD5(containerName string, remoteName string, md5 string) {
	blobURL := getFileBlobURL(containerName, remoteName)
	ctx := context.Background() 

	blobProps, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	check(err)

	// Add or update md5 property 
	metadata := blobProps.NewMetadata()
	metadata[pitMD5tag] = md5 
	_, err = blobURL.SetMetadata(ctx, metadata, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	check(err)
}

func getRemoteFileMD5AndETag(remoteFileName string) (string, string, error) {
	MD5 := ""
	ETag := ""

	metaData, err := getDocumentMetadata(getContainerName(),remoteFileName)
	if err != nil {
		return MD5, ETag, err
	}

	for _, element := range metaData {
		// To view print all of the remote file's metadate, uncomment the following line: 
		//     fmt.Println("metaDataElement: "+element)
		if strings.Contains(element, pitMD5tag) {
			strParts := strings.Split(element, "=")
			MD5 = strParts[1]
		}

		if strings.Contains(element, "ETag=") {
			strParts := strings.Split(element, "=")

			// The ETag is returned as a string that includes double quotes so we need to 
			//     remove the leading and trailing double quotes.
			ETag = strings.TrimLeft(strParts[1], "\"")
			ETag = strings.TrimRight(ETag, "\"")
		}
	}

	return MD5, ETag, nil
}

func getDocumentMetadata(containerName string, remoteName string) ([]string, error) {
	blobURL := getFileBlobURL(containerName, remoteName)
	ctx := context.Background() 

	// Query the blob's properties and metadata.
	blobProps, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	var metadataElements []string

	// Add some of the Azure standard metadata.
	metadataElements = append(metadataElements, string("BlogType="+blobProps.BlobType()))
	metadataElements = append(metadataElements, string("ETag="+blobProps.ETag()))
	metadataElements = append(metadataElements, string("LastModified="+blobProps.LastModified().Format(time.RFC3339)))
	
	// Consider utilizing Azures MD5 property:
	//     fmt.Println("MD5=", get.ContentMD5())

	// Consider utilizing Azure ETag property: 
	//     References:
	//     https://cann0nf0dder.wordpress.com/2015/09/07/azure-blob-storage-and-managing-concurrency/

	// Add application custom metadata.
	metadata := blobProps.NewMetadata()
	for k, v := range metadata {
		metadataElements = append(metadataElements, k+"="+v)
	}
	return metadataElements, nil
}

func listBlobs(containerName string) {
	documentNames := getDocumentNames(containerName)
	if len(documentNames) == 0 {
		fmt.Printf("No files currently in container \"%sd\"\n", containerName)		
	} else {
		fmt.Printf("Files in container \"%sd\":\n", containerName)
		for _, documentName := range documentNames {
			fmt.Println("    "+documentName)
		}
		fmt.Printf("Number of documents %d\n", len(documentNames))
	}
}