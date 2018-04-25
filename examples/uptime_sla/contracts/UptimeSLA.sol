pragma solidity ^0.4.23;

import "../../../solidity/contracts/Chainlinked.sol";
import "../../../solidity/contracts/ChainlinkLib.sol";

contract UptimeSLA is Chainlinked {
  uint256 constant uptimeThreshold = 9999;
  bytes32 private jobId;
  uint256 private endAt;
  address private client;
  address private serviceProvider;
  bytes32 public externalId;
  uint256 public uptime;

  function UptimeSLA(
    address _client,
    address _serviceProvider,
    address _link,
    address _oracle,
    bytes32 _jobId
  ) public payable {
    client = _client;
    serviceProvider = _serviceProvider;
    endAt = block.timestamp + 30 days;
    jobId = _jobId;

    setLinkToken(_link);
    setOracle(_oracle);
  }

  function updateUptime(string _when) public {
    ChainlinkLib.Run memory run = newRun(jobId, this, "report(bytes32,uint256)");
    run.add("url", "https://status.heroku.com/api/ui/availabilities");
    string[] memory path = new string[](4);
    path[0] = "data";
    path[1] = _when;
    path[2] = "attributes";
    path[3] = "calculation";
    run.addStringArray("path", path);
    externalId = chainlinkRequest(run, 1 szabo);
  }

  function report(bytes32 _externalId, uint256 _uptime)
    public
    onlyOracle
    checkRequestId(_externalId)
  {
    uptime = _uptime;
    if (_uptime < uptimeThreshold) {
      client.send(this.balance);
    } else if (block.timestamp >= endAt) {
      serviceProvider.send(this.balance);
    }
  }

  modifier checkRequestId(bytes32 _externalId) {
    require(externalId == _externalId);
    _;
  }

}
