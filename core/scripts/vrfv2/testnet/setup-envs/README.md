1. go to `core/scripts/vrfv2/testnet/setup-envs` folder and `go build`
2. go to `core/scripts/vrfv2/testnet/docker` folder and start containers - `docker compose up`
3. execute from `core/scripts/vrfv2/testnet/setup-envs` folder
    ```
    ./setup-envs -remote-node-urls=http://localhost:6610 -creds-file ../docker/secrets/apicredentials
    ```


    ```
    ./setup-envs -vrf-primary-node-url=http://localhost:6610 -creds-file ../docker/secrets/apicredentials
    ```



   ```
   source eth-sepolia-staging.env
   
   go run . --vrf-primary-node-url=http://localhost:6610 --creds-file ../docker/secrets/apicredentials
   ```