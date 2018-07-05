#!/bin/bash -e

export GOOGLE_APPLICATION_CREDENTIALS=/keyconfig.json

codeship_google authenticate
gcloud config set compute/zone us-central1-a
gcloud container clusters get-credentials acceptance

# set chainlink pod image to one recently pushed by build
kubectl set image \
  deployment.apps/chainlink-deploy \
  chainlink=gcr.io/chainlink-209321/chainlink:0.2.0-$CI_TIMESTAMP
