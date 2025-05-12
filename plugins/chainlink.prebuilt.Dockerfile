##
# Takes Chainlink core as a base image and layers in private plugins.
##
ARG BASE_IMAGE=public.ecr.aws/chainlink/chainlink:v2.23.0-plugins

##
# Final image
##
FROM ${BASE_IMAGE} AS final
# This directory should only contain the plugin binaries.
ARG LOCAL_PLUGIN_DIR=./plugins
ARG LOCAL_LIB_DIR=./lib

# Copy/override any additional shared libraries.
COPY ${LOCAL_LIB_DIR} /usr/lib/
# Copy/override plugins.
COPY ${LOCAL_PLUGIN_DIR} /usr/local/bin
