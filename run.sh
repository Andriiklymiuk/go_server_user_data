#!/bin/bash

# Check if fswatch is installed
if ! command -v fswatch &> /dev/null
then
    echo "fswatch could not be found. Please install it before running this script."
    exit
fi

# Function to check if a process is running
is_running() {
  if ps -p $1 > /dev/null; then
    return 0
  else
    return 1
  fi
}

# Function to stop running server
stop_server() {
  if [ ! -z "$PID" ] && is_running $PID; then
    kill -SIGTERM $PID
    wait $PID
  fi
}

# Function to start server
start_server() {
  make build
  ./build &
  PID=$!
}

# Trap SIGINT (Ctrl+C) and stop the server
trap stop_server INT

# Initial start
start_server

# Watch for changes in .go files and .env in the current directory and its subdirectories
fswatch -o -r -e ".*" -i "\\.go$" -i "\\.env$" . | while read
do
  echo "Changes detected. Restarting server..."
  stop_server
  start_server
done

printf ""
