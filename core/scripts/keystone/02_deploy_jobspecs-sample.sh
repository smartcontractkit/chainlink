#!/bin/bash

go run main.go \
 deploy-jobspecs \
 --chainid=1337 \
 --p2pport=6690 \
 --onlyreplay=false
