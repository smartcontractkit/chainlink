pragma solidity ^0.4.18;

import "../../../Chainlinked.sol";

contract UptimeSLA is Chainlinked {
  uint256 constant uptimeThreshold = 9999;
  bytes32 private jobId;
  uint256 private endAt;
  address private client;
  address private serviceProvider;
  uint256 public requestId;
  uint256 public uptime;

  function UptimeSLA(
    address _client,
    address _serviceProvider,
    address _oracle,
    bytes32 _jobId
  ) public payable {
    client = _client;
    serviceProvider = _serviceProvider;
    endAt = block.timestamp + 30 days;
    oracle = Oracle(_oracle);
    jobId = _jobId;
  }

  function updateUptime(string _when) public {
    Chainlink.Run memory run = newRun(jobId, this, "report(uint256,uint256)");
    run.add("url", "https://status.heroku.com/api/ui/availabilities");
    string[] memory path = new string[](4);
    path[0] = "data";
    path[1] = _when;
    path[2] = "attributes";
    path[3] = "calculation";
    run.add("path", path);
    requestId = chainlinkRequest(run);
  }

  function report(uint256 _requestId, uint256 _uptime)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    uptime = _uptime;
    if (_uptime < uptimeThreshold) {
      client.send(this.balance);
    } else if (block.timestamp >= endAt) {
      serviceProvider.send(this.balance);
    }
  }

  modifier checkRequestId(uint256 _requestId) {
    require(requestId == _requestId);
    _;
  }

}
