import os

# Read Pit version.go file.
versionFileName = "version.go"

versionGoFileRead = open(versionFileName, "r")
versionGoLines = versionGoFileRead.readlines()
versionGoFileRead.close()

versionGoLastLine = versionGoLines[len(versionGoLines)-1]
locationOfDot = versionGoLastLine.rfind(".")
locationOfDoubleQuote = versionGoLastLine.rfind("\"")

numberOfLastBuild = int(versionGoLastLine[locationOfDot+1:locationOfDoubleQuote])
numberOfCurrentBuild = numberOfLastBuild + 1

print("update build number: " + str(numberOfLastBuild) + " ⇒ " + str(numberOfCurrentBuild))

lineWithUpdatedNumber = versionGoLastLine[:locationOfDot+1]
lineWithUpdatedNumber += str(numberOfCurrentBuild) + "\""

versionGoLines[len(versionGoLines)-1] = lineWithUpdatedNumber

# Write Pit version.go file.
versionGoFileWrite = open(versionFileName, "w")
for line in versionGoLines:
	versionGoFileWrite.write(line)

versionGoFileRead.close()