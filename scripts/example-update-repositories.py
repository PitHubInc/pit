import os

STUDENT_REPOSITORIES = [
	"https://github.com/KurtDankovich/360-kurt-dankovich.git",
	"https://github.com/BryanGabe00/360-bryan-gabe",
	"https://github.com/RichGol/360-richard-goluszka",
	"https://github.com/mcclint50/360-colin-mcclintic.git",
	"https://github.com/JuanMoncada23/360-juan-moncada.git",
	"https://github.com/BrennanP01/360-brennan-price.git",
	"https://github.com/mrodriguezdelcorral/360-Maria-Rodriguez.git"
]

SPRINT_3_TEAM_REPOSITORIES = [
	"https://github.com/mcclint50/360-Mongooses.git",
	"https://github.com/BrennanP01/360-gloriousKenobis.git",
	"https://github.com/JuanMoncada23/360-RedDragons.git"
]


SPRINT_4_TEAM_REPOSITORIES = [
	"https://github.com/mcclint50/360-Mongooses.git", 		# Mongooses (1)
	"https://github.com/BryanGabe00/QuizMaster.git", 		# RedDragons
	"https://github.com/BrennanP01/360-gloriousKenobis.git" # GloriousKenobis (3)
]

def printAndSystemExecute(executeString):
	print('Executing: ' + executeString)
	os.system(executeString)

def clone(repositoryLink):
	printAndSystemExecute("git clone %s" % repositoryLink)

def update(directoryName):
	os.chdir(directoryName)
	printAndSystemExecute("git pull")
	os.chdir("..")

def cloneOrUpdate(repositoryLink):
	# Split the path name so that we have the local directory name.
	#     https://github.com/KurtDankovich/360-kurt-dankovich.git -> 360-kurt-dankovich
	substrings = repositoryLink.split("/")
	substrings = substrings[len(substrings)-1].split(".")
	directoryName = substrings[0]

	isDir = os.path.isdir(directoryName)
	if isDir:
		print("Update %s" % repositoryLink)
		update(directoryName)
	else:
		print("Clone %s" % repositoryLink)
		clone(repositoryLink)

	print("")

os.system("clear")
print("Cloning or Updating Student Repositories:\n")
for repositoryLink in STUDENT_REPOSITORIES:
	cloneOrUpdate(repositoryLink)

print("Cloning or Updating Sprint 3 Team Repositories:\n")
for repositoryLink in SPRINT_3_TEAM_REPOSITORIES:
	cloneOrUpdate(repositoryLink)

print("Cloning or Updating Sprint 4 Team Repositories:\n")
for repositoryLink in SPRINT_4_TEAM_REPOSITORIES:
	cloneOrUpdate(repositoryLink)
