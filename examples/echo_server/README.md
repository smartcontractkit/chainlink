# Echo Server

Using Chainlink (CL), this application simply echos incoming ethereum logs
as JSON, listened to by a CL job. It is intended to demonstrate the the first
step to bridging on chain to off chain activity.

![Log Echo Server](screenshot.jpg?raw=true "Log Echo Server")

## Configure and run [Chainlink development environment](../README.md)

## Run Echo Server EthLog (Raw Ethereum Logs)

Uses an `ethlog` initiator to echo all log events. An `ethlog` initiator starts
a job anytime a log event occurs. It can optionally be filtered by an `address`.

1. Complete the [Run Chainlink Development Environment](../README.md#run-chainlink-development-environment) steps.
2. `./create_ethlog_job` to create the Chainlink (CL) job
3. `yarn install`
4. `node echo.js`
5. `./node_modules/.bin/truffle migrate` in another window
6. `node send_ethlog_transaction.js`
7. Wait for log to show up in echo server

## Run Echo Server RunLog (Chainlink Specific Ethereum Logs)

Uses a `runlog` initiator to echo Chainlink log events with the matching job id.

1. Complete the [Run Chainlink Development Environment](../README.md#run-chainlink-development-environment) steps.
2. `./create_runlog_job` to create CL job. Keep track of the returned job id.
3. `yarn install`
4. `node echo.js`
5. Add job id to `contracts/RunLog.sol` where it says `MY_JOB_ID`
5. `./node_modules/.bin/truffle migrate --reset` in another window
6. `node send_runlog_transaction.js`
7. Wait for log to show up in echo server


## Further Reading

Please see the other examples in the repo, and take care to
identify the difference between an `HttpPost` task and a bridge to an external
adapter. The latter allows an asynchronous response from a service to
then potentially write back to the chain.
