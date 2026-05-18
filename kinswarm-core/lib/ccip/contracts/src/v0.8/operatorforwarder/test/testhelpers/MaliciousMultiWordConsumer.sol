// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ChainlinkClient} from "../../../ChainlinkClient.sol";
import {Chainlink} from "../../../Chainlink.sol";

contract MaliciousMultiWordConsumer is ChainlinkClient {
  uint256 private constant ORACLE_PAYMENT = 1 ether;
  uint256 private s_expiration;

  constructor(address _link, address _oracle) payable {
    _setChainlinkToken(_link);
    _setChainlinkOracle(_oracle);
  }

  receive() external payable {} // solhint-disable-line no-empty-blocks

  function requestData(bytes32 _id, bytes memory _callbackFunc) public {
    Chainlink.Request memory req = _buildChainlinkRequest(_id, address(this), bytes4(keccak256(_callbackFunc)));
    s_expiration = block.timestamp + 5 minutes; // solhint-disable-line not-rely-on-time
    _sendChainlinkRequest(req, ORACLE_PAYMENT);
  }

  function assertFail(bytes32, bytes memory) public pure {
    assert(1 == 2);
  }

  function cancelRequestOnFulfill(bytes32 _requestId, bytes memory) public {
    _cancelChainlinkRequest(_requestId, ORACLE_PAYMENT, this.cancelRequestOnFulfill.selector, s_expiration);
  }

  function remove() public {
    selfdestruct(payable(address(0)));
  }

  function stealEthCall(bytes32 _requestId, bytes memory) public recordChainlinkFulfillment(_requestId) {
    (bool success, ) = address(this).call{value: 100}(""); // solhint-disable-line avoid-call-value
    // solhint-disable-next-line gas-custom-errors
    require(success, "Call failed");
  }

  function stealEthSend(bytes32 _requestId, bytes memory) public recordChainlinkFulfillment(_requestId) {
    // solhint-disable-next-line check-send-result
    bool success = payable(address(this)).send(100); // solhint-disable-line multiple-sends
    // solhint-disable-next-line gas-custom-errors
    require(success, "Send failed");
  }

  function stealEthTransfer(bytes32 _requestId, bytes memory) public recordChainlinkFulfillment(_requestId) {
    payable(address(this)).transfer(100);
  }

  function doesNothing(bytes32, bytes memory) public pure {} // solhint-disable-line no-empty-blocks
}
