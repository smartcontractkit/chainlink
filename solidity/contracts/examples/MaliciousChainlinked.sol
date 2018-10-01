pragma solidity ^0.4.24;

import "./MaliciousChainlinkLib.sol";
import "../Oracle.sol";
import "link_token/contracts/LinkToken.sol";

contract MaliciousChainlinked {
  using MaliciousChainlinkLib for MaliciousChainlinkLib.Run;
  using MaliciousChainlinkLib for MaliciousChainlinkLib.WithdrawRun;
  using SafeMath for uint256;

  uint256 constant clArgsVersion = 1;

  LinkToken internal link;
  Oracle internal oracle;
  uint256 internal requests = 1;
  mapping(bytes32 => bool) internal unfulfilledRequests;

  event ChainlinkRequested(bytes32 id);
  event ChainlinkFulfilled(bytes32 id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (MaliciousChainlinkLib.Run memory) {
    MaliciousChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function newWithdrawRun(
    bytes32 _specId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (MaliciousChainlinkLib.WithdrawRun memory) {
    MaliciousChainlinkLib.WithdrawRun memory run;
    return run.initializeWithdraw(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(MaliciousChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.close();
    require(link.transferAndCall(oracle, _wei, _run.encodeForOracle(clArgsVersion)), "Unable to transferAndCall to oracle");
    emit ChainlinkRequested(_run.requestId);
    unfulfilledRequests[_run.requestId] = true;
    return _run.requestId;
  }

  function chainlinkWithdrawRequest(MaliciousChainlinkLib.WithdrawRun memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.closeWithdraw();
    require(link.transferAndCall(oracle, _wei, _run.encodeWithdrawForOracle(clArgsVersion)), "Unable to transferAndCall to oracle");
    emit ChainlinkRequested(_run.requestId);
    unfulfilledRequests[_run.requestId] = true;
    return _run.requestId;
  }

  function LINK(uint256 _amount) internal view returns (uint256) {
    return _amount.mul(10**18);
  }

  function setOracle(address _oracle) internal {
    oracle = Oracle(_oracle);
  }

  function setLinkToken(address _link) internal {
    link = LinkToken(_link);
  }

  modifier checkChainlinkFulfillment(bytes32 _requestId) {
    require(msg.sender == address(oracle) && unfulfilledRequests[_requestId], "Source must be oracle with a valid requestId");
    _;
    unfulfilledRequests[_requestId] = false;
    emit ChainlinkFulfilled(_requestId);
  }
}