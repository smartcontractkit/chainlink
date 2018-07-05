#!/bin/bash -e

export GOOGLE_APPLICATION_CREDENTIALS=/keyconfig.json

codeship_google authenticate
gcloud config set compute/zone us-central1-a
gcloud container clusters get-credentials staging

# set chainlink pod image to one recently pushed by build
kubectl set image \
  deployment.apps/chainlink-deploy \
  chainlink=smartcontract/chainlink:0.2.0
