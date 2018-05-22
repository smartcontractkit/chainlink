pragma solidity ^0.4.23;

import "../../../solidity/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  bytes32 private jobId;

  constructor(address _link, address _oracle, bytes32 _jobId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function request() public {
    ChainlinkLib.Run memory run = newRun(jobId, this, "fulfill(bytes32,bytes32)");
    run.add("msg", "hello_chainlink");
    chainlinkRequest(run, LINK(1));
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    checkChainlinkFulfillment(_externalId)
  {
  }
}
