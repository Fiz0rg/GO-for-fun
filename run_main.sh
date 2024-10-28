#!/bin/bash

# load envs

cd "$(dirname "$0")"

# RUN SCRIPT
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Run script" >> "logfile.log"
/usr/local/go/bin/go run main.go >> "logfile.log" 2>&1

# Check that script executed
if [ $? -eq 0 ]; then
    echo "Script successfully executed" >> "logfile.log"
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Script error"  >> "logfile.log"
fi