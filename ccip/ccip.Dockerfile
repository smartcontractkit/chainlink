FROM chainlink:latest

# Add the config directory
COPY ./config /ccip/config

# Expose the config directory as an environment variable
ENV CL_CHAIN_DEFAULTS=/ccip/config