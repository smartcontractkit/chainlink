# Chainlink Ingester

This directory contains the ingester project, which is responsible for taking 
data off the chain and aggregating it for reporting or analysis purposes.


## Configuration

The ingester application takes the following environment variables

```bash
# Ethereum Chain ID for the contracts you want to listen to
ETH_CHAIN_ID
# Websocket endpoint the monitor uses to watch the aggregator contracts
ETH_URL
# Postgres database host
DB_HOST
# Postgres database name
DB_NAME
# Postgres database port
DB_PORT
# Postgres database username
DB_USERNAME
# Postgres database password
DB_PASSWORD
```
