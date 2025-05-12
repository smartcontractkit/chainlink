##
# Takes Chainlink core as a base image and layers in private plugins.
##
ARG BASE_IMAGE=public.ecr.aws/chainlink/chainlink:v2.23.0-plugins

##
# Final image
##
FROM ${BASE_IMAGE} AS final
# This directory should contain a bin/ subdir with the plugins and an optional lib/ subdir with shared libraries.
ARG PKG_PATH=./build

# Copy/override any (optional) additional shared libraries.
COPY ${PKG_PATH}/li[b] /usr/lib/
# Copy/override plugins.
COPY ${PKG_PATH}/bin /usr/local/bin
