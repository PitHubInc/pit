import os
import subprocess

# Removing automatic update of build number. We may want to put this back in later.
# subprocess.run(["python3", "scripts/update-build-number.py"])

print("go build -o pit")
os.system("go build -o pit")

print("./pit version")
os.system("./pit version")