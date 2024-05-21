pragma solidity ^0.8.0;

import {MaliciousChainlink} from "./MaliciousChainlink.sol";
import {MaliciousChainlinked, Chainlink} from "./MaliciousChainlinked.sol";
import {ChainlinkRequestInterface} from "../../../../interfaces/ChainlinkRequestInterface.sol";

contract MaliciousRequester is MaliciousChainlinked {
  uint256 private constant ORACLE_PAYMENT = 1 ether;
  uint256 private s_expiration;

  constructor(address _link, address _oracle) {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function maliciousWithdraw() public {
    MaliciousChainlink.WithdrawRequest memory req = newWithdrawRequest(
      "specId",
      address(this),
      this.doesNothing.selector
    );
    chainlinkWithdrawRequest(req, ORACLE_PAYMENT);
  }

  function request(bytes32 _id, address _target, bytes memory _callbackFunc) public returns (bytes32 requestId) {
    Chainlink.Request memory req = newRequest(_id, _target, bytes4(keccak256(_callbackFunc)));
    s_expiration = block.timestamp + 5 minutes; // solhint-disable-line not-rely-on-time
    return chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function maliciousPrice(bytes32 _id) public returns (bytes32 requestId) {
    Chainlink.Request memory req = newRequest(_id, address(this), this.doesNothing.selector);
    return chainlinkPriceRequest(req, ORACLE_PAYMENT);
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    Chainlink.Request memory req = newRequest("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    return chainlinkTargetRequest(_target, req, ORACLE_PAYMENT);
  }

  function maliciousRequestCancel(bytes32 _id, bytes memory _callbackFunc) public {
    ChainlinkRequestInterface oracle = ChainlinkRequestInterface(oracleAddress());
    oracle.cancelOracleRequest(
      request(_id, address(this), _callbackFunc),
      ORACLE_PAYMENT,
      this.maliciousRequestCancel.selector,
      s_expiration
    );
  }

  function doesNothing(bytes32, bytes32) public pure {} // solhint-disable-line no-empty-blocks
}
