pragma solidity ^0.4.24;

import "chainlink/contracts/Chainlinked.sol";

contract RunLog is Chainlinked {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  event Fulfillment(bytes32 data);

  constructor(address _link, address _oracle) public {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function request(bytes32 _jobId) public {
    Chainlink.Request memory req = newRequest(_jobId, this, this.fulfill.selector);
    req.add("msg", "hello_chainlink");
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function fulfill(bytes32 _externalId, bytes32 _data)
    public
    recordChainlinkFulfillment(_externalId)
  {
      emit Fulfillment(_data);
  }
}
