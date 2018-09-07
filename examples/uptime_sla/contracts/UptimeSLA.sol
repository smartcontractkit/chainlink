pragma solidity ^0.4.24;

import "../../../solidity/contracts/Chainlinked.sol";

contract UptimeSLA is Chainlinked {
  uint256 constant uptimeThreshold = 9999;
  bytes32 private jobId;
  uint256 private endAt;
  address private client;
  address private serviceProvider;
  uint256 public uptime;

  constructor(
    address _client,
    address _serviceProvider,
    address _link,
    address _oracle,
    bytes32 _jobId
  ) public payable {
    client = _client;
    serviceProvider = _serviceProvider;
    endAt = block.timestamp.add(30 days);
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
    chainlinkRequest(run, LINK(1));
  }

  function report(bytes32 _externalId, uint256 _uptime)
    public
    checkChainlinkFulfillment(_externalId)
  {
    uptime = _uptime;
    if (_uptime < uptimeThreshold) {
      client.transfer(address(this).balance);
    } else if (block.timestamp >= endAt) {
      serviceProvider.transfer(address(this).balance);
    }
  }

}
