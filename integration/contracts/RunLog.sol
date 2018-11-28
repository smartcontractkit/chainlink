pragma solidity ^0.4.23;

import "Chainlinked.sol";

contract RunLog is Chainlinked {

  constructor(address _link, address _oracle) public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function request(bytes32 _jobId) public {
    ChainlinkLib.Run memory run = newRun(_jobId, this, this.fulfill.selector);
    run.add("msg", "hello_chainlink");
    chainlinkRequest(run, LINK(1));
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    checkChainlinkFulfillment(_externalId)
  {
  }
}
