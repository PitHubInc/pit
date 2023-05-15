import os
import subprocess

subprocess.run(["python3", "scripts/update-build-number.py"])

print("go build -o pit")
os.system("go build -o pit")

print("./pit version")
os.system("./pit version")