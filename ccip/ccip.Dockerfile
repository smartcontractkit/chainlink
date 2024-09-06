FROM chainlink-ccip:latest

COPY ccip/config /chainlink/ccip-config

# Expose the config directory as an environment variable
ENV CL_CHAIN_DEFAULTS=/chainlink/ccip-config