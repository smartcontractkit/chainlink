# Uptime Service Level Agreement

An example SLA that uses ChainLink to determine the release of payment.

When the contract is deployed a client, service provider, and start time are specified. Additionally a deposit is made. The end of the contract is set to 30 days after the start time.

After the contract is created anyone can request updates from the oracle for the contract. If the oracle reports that the uptime is below 99.99% then the deposit is released to the client. If the rate is still above 99.99% after the contract ends, and the deposit has not been released, the deposit is sent to the service provider.

```solidity
function report(uint256 _requestId, uint256 _rate)
    public
    recordChainlinkFulfillment(_requestId)
  {
    if (_rate < uptimeThreshold) {
      client.send(this.balance);
    } else if (block.timestamp >= endAt) {
      serviceProvider.send(this.balance);
    }
  }
```

# ChainLink

Initiator: `runLog`

Job Pipeline: `httpGet` => `jsonParse` => `multiply` => `ethUint256` => `ethTx`

This contract displays ChainLinks ability to pull in data from outside data feeds and format it to be used by Ethereum contracts.

A float value is pulled out of a nested JSON object and multiplied to a precision level that is useful for the contract.

The ChainLink Job is configured to not take any specific URL or JSON path, so that this oracle and job can be reused for other APIs. Both `url` and `path` are passed into the oracle by the SLA contract, specifically which data point to use is passed into the contract:

```solidity
function updateUptime(string _when) public {
   Chainlink.Request memory req = newRequest(jobId, this, "report(uint256,uint256)");
   req.add("get", "https://status.heroku.com/api/ui/availabilities");
   string[] memory path = new string[](4);
   path[0] = "data";
   path[1] = _when;           //pick which data point in the array you want to examine
   path[2] = "attributes";
   path[3] = "calculation";
   req.add("path", path);
   chainlinkRequest(req, LINK(1));
}
```

The API returns the percentage as a float, for example the current value is `0.999999178716033`. The `multiply` adapter takes that result and multiplies it by 1000, which [a parameter specified in the `times` field](https://github.com/smartcontractkit/hello_chainlink/blob/4b42f127ddeca6541ac2aba1803f458d0a3bf460/uptime_sla/http_json_x10000_job.json). The result is `9999`, allowing the contract to check for "four nines" of uptime.

## Configure and run [Chainlink development environment](../README.md#run-chainlink-development-environment)

## Run and update the Uptime SLA contract.

1. `yarn install`
2. `./deploy` in another window
3. `./send_sla_transaction.js` to trigger an update to the SLA
4. `./get_uptime.js` get the latest uptime

## Development

To run the tests, call `./node_modules/.bin/truffle --network=test test`
