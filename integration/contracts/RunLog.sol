pragma solidity ^0.4.23;

import "Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  event FulfilledEvent(bytes32 data);

  constructor(address _link, address _oracle) public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function request(bytes32 _jobId) public {
    ChainlinkLib.Run memory run = newRun(_jobId, this, this.fulfill.selector);
    run.add("msg", "hello_chainlink");
    chainlinkRequest(run, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    checkChainlinkFulfillment(_externalId) {
    emit FulfilledEvent(_data);
  }
}
