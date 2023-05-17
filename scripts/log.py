import sys
import os

# ChatGPT generated initial function with: Write a python function that creates a new file that contains a single line
# that gives the current date and time"
from datetime import datetime
def create_file_with_datetime(file_path):
    try:
        current_datetime = datetime.now()
        formatted_time = current_datetime.strftime("%Y/%m/%d %H:%M:%S")
        
        with open(file_path, "w") as file:
            file.write(formatted_time + " New log file created\n")
    except IOError:
        print(f"An error occurred while creating the file '{file_path}'.")

# ChatGPT generated initial function with: Write python code that removes a file if it exists but provides no user 
# messages.
def remove_file(file_path):
    if os.path.exists(file_path):
        try:
            os.remove(file_path)
        except OSError:
            pass
# End ChatGPT generated


# Todo: Read log file name from Go code so that it remains consistent.
logFilePathAndName = '/Users/epogue/.pit/log.txt'

# Stop logging by removing log file and exit.
if len(sys.argv) == 2:
	if sys.argv[1] == '-s' or sys.argv[1] == "-stop":
		remove_file(logFilePathAndName)
		sys.exit(0) 

# Truncate log file if requested.
if len(sys.argv) == 2:
	if sys.argv[1] == '-t' or sys.argv[1] == "-truncate":
		remove_file(logFilePathAndName)
		create_file_with_datetime(logFilePathAndName)

os.system(f'tail -f {logFilePathAndName}')


