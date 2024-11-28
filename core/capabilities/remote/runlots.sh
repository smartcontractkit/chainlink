#!/bin/bash

# Number of times to run the command
x=50

# Command to run
command="go test ./... --race"

# Loop to run the command x times
for ((i=1; i<=x; i++))
do
  $command
done
