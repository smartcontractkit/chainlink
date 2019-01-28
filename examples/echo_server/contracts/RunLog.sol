pragma solidity 0.4.24;

import "../../../solidity/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  bytes32 private jobId;

  constructor(address _link, address _oracle, bytes32 _jobId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function request() public {
    Chainlink.Run memory run = newRun(jobId, this, this.fulfill.selector);
    run.add("msg", "hello_chainlink");
    chainlinkRequest(run, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    checkChainlinkFulfillment(_externalId)
  {
  }
}
