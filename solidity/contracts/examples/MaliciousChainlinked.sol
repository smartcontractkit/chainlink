pragma solidity 0.4.24;

import "./MaliciousChainlinkLib.sol";
import "../interfaces/OracleInterface.sol";
import "../interfaces/LinkTokenInterface.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract MaliciousChainlinked {
  using MaliciousChainlinkLib for MaliciousChainlinkLib.Run;
  using MaliciousChainlinkLib for MaliciousChainlinkLib.WithdrawRun;
  using SafeMath for uint256;

  uint256 constant clArgsVersion = 1;

  LinkTokenInterface internal link;
  OracleInterface internal oracle;
  uint256 internal requests = 1;
  mapping(bytes32 => bool) internal unfulfilledRequests;

  event ChainlinkRequested(bytes32 id);
  event ChainlinkFulfilled(bytes32 id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunction
  ) internal pure returns (MaliciousChainlinkLib.Run memory) {
    MaliciousChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunction);
  }

  function newWithdrawRun(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunction
  ) internal pure returns (MaliciousChainlinkLib.WithdrawRun memory) {
    MaliciousChainlinkLib.WithdrawRun memory run;
    return run.initializeWithdraw(_specId, _callbackAddress, _callbackFunction);
  }

  function chainlinkRequest(MaliciousChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.close();
    require(link.transferAndCall(oracle, _wei, encodeForOracle(_run)), "Unable to transferAndCall to oracle");
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
    require(link.transferAndCall(oracle, _wei, encodeWithdrawForOracle(_run)), "Unable to transferAndCall to oracle");
    emit ChainlinkRequested(_run.requestId);
    unfulfilledRequests[_run.requestId] = true;
    return _run.requestId;
  }

  function LINK(uint256 _amount) internal pure returns (uint256) {
    return _amount.mul(10**18);
  }

  function setOracle(address _oracle) internal {
    oracle = OracleInterface(_oracle);
  }

  function setLinkToken(address _link) internal {
    link = LinkTokenInterface(_link);
  }

  function encodeForOracle(MaliciousChainlinkLib.Run memory _run)
    internal view returns (bytes memory)
  {
    return abi.encodeWithSelector(
      oracle.requestData.selector,
      address(this), // overridden by onTokenTransfer
      100 ether,     // overridden by onTokenTransfer
      clArgsVersion,
      _run.specId,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.requestId,
      _run.buf.buf);
  }

  function encodeWithdrawForOracle(MaliciousChainlinkLib.WithdrawRun memory _run)
    internal view returns (bytes memory)
  {
    return abi.encodeWithSelector(
      oracle.withdraw.selector,
      _run.callbackAddress,
      _run.amount,
      _run.buf.buf);
  }

  modifier checkChainlinkFulfillment(bytes32 _requestId) {
    require(msg.sender == address(oracle) && unfulfilledRequests[_requestId], "Source must be oracle with a valid requestId");
    _;
    unfulfilledRequests[_requestId] = false;
    emit ChainlinkFulfilled(_requestId);
  }
}