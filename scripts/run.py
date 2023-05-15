import os
import sys

command = "./pit"
for param in range(1, len(sys.argv)):
    command = command + ' '+sys.argv[param]

print(command)
os.system(command)