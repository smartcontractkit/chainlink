#!/bin/bash

echo Running start script... &&
chainlink local import /keys/$KEY_NAME &&
chainlink local node -d -p /password.txt -a /apicredentials