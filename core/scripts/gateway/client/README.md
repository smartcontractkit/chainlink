# A gateway client script

This script is used to connect to a gateway server and send commands to it.

## Usage

All requests have to be signed on behalf of a user, you need to provide your private key in .env file, e.g.

```
PRIVATE_KEY=1a2b3c...
```

The script will automatically sign the message using the provided private key.
Run the script without arguments to get the list of available commands.

## Example

```
go run . -gateway_url https://01.functions-gateway.chain.link -don_id fun-avalanche-mainnet-2 -method secrets_list -message_id 123
```
