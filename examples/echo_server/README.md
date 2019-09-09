# Echo Server

Using Chainlink (CL), this application simply echos incoming ethereum logs
as JSON, listened to by a CL job. It is intended to demonstrate the the first
step to bridging on chain to off chain activity.

![Log Echo Server](screenshot.jpg?raw=true 'Log Echo Server')

## Configure and run [Chainlink development environment](../README.md)

## Run EthLog (Raw Ethereum Logs)

Uses an `ethlog` initiator to echo all log events. An `ethlog` initiator starts
a job anytime a log event occurs. It can optionally be filtered by an `address`.

1. Complete the [Run Chainlink Development Environment](../README.md#run-chainlink-development-environment) steps.
2. `./create_ethlog_job` to create the Chainlink (CL) job
3. `yarn install`
4. `node echo.js`
5. `yarn truffle migrate` in another window
6. `node send_ethlog_transaction.js`
7. Wait for log to show up in echo server

## Run RunLog (Chainlink Specific Ethereum Logs)

Uses a `runlog` initiator to echo Chainlink log events with the matching job id.

1. Complete the [Run Chainlink Development Environment](../README.md#run-chainlink-development-environment) steps.
2. `yarn install`
3. `node echo.js`
4. `yarn truffle migrate` in another window
5. `node send_runlog_transaction.js`
6. Wait for log to show up in echo server
7. Investigate migrations/5_run_log.js for insight

## Further Reading

Please see the other examples in the repo, and take care to
identify the difference between an `HttpPost` task and a bridge to an external
adapter. The latter allows an asynchronous response from a service to
then potentially write back to the chain.

## Development

To run the tests, call `./node_modules/.bin/truffle --network=test test`
