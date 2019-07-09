pragma solidity 0.4.24;

import "chainlink/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  bytes32 private jobId;

  constructor(address _link, address _oracle, bytes32 _jobId) public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function request() public {
    Chainlink.Request memory req = newRequest(jobId, this, this.fulfill.selector);
    req.add("msg", "hello_chainlink");
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    recordChainlinkFulfillment(_externalId)
  {} // solium-disable-line no-empty-blocks

}
