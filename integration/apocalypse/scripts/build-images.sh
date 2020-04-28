#!/bin/bash

pushd gethnet && \
docker build -t smartcontract/gethnet:apocalypse . && \
popd && \
pushd paritynet && \
docker build -t smartcontract/paritynet:apocalypse . && \
popd
