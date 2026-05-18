#!/bin/bash

go run main.go \
 deploy-jobspecs \
 --chainid=11155111 \
 --p2pport=6690 \
 --onlyreplay=false
