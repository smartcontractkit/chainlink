#!/bin/bash

set -e

# authenticate to google cloud
codeship_google authenticate

# set compute zone
gcloud config set compute/zone us-central1-a

# set kubernetes cluste
gcloud container clusters get-credentials cluster-1

# update kubernetes Deployment
GOOGLE_APPLICATION_CREDENTIALS=/keyconfig.json kubectl set image deployment.apps/chainlink-deploy chainlink=smartcontract/chainlink:0.2.0-$CI_TIMESTAMP
