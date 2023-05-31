# Pit Video Sharing

# Build
## go build -o pit
## OR
## python ./scripts/build.py

# Execute
## ./pit
## OR
## python ./scripts/run.py
##
## Note that if you have pit installed on your system already and you type "pit" you will execute the globally 
## installed version of pit and not the local copy that was presumably just compiled.

# View log:
## python ./scripts/log.py

# Update build number:
## python ./script/update-build-number.py

# Deploy pit to local computer
## python ./scripts/deploy.py

# Dependency
## go get -u github.com/Azure/azure-storage-blob-go/azblob

# Additional information:
## More Go build information: https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs
## More Go install information: https://golang.org/doc/code.html 

# If we need to reinstall MacOS developer tools
sudo rm -rf /Library/Developer/CommandLineTools
sudo xcode-select --install