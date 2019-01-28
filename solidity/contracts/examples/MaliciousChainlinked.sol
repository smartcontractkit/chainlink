pragma solidity 0.4.24;

import "./MaliciousChainlink.sol";
import "../Chainlinked.sol";

contract MaliciousChainlinked is Chainlinked {
  using MaliciousChainlink for MaliciousChainlink.Run;
  using MaliciousChainlink for MaliciousChainlink.WithdrawRun;
  using Chainlink for Chainlink.Run;
  using SafeMath for uint256;

  uint256 private maliciousRequests = 1;
  mapping(bytes32 => address) private maliciousUnfulfilledRequests;

  function newWithdrawRun(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunction
  ) internal pure returns (MaliciousChainlink.WithdrawRun memory) {
    MaliciousChainlink.WithdrawRun memory run;
    return run.initializeWithdraw(_specId, _callbackAddress, _callbackFunction);
  }

  function chainlinkTargetRequest(address _target, Chainlink.Run memory _run, uint256 _amount)
    internal
    returns(bytes32 requestId)
  {
    requestId = keccak256(abi.encodePacked(_target, maliciousRequests));
    _run.nonce = maliciousRequests;
    maliciousUnfulfilledRequests[requestId] = oracleAddress();
    emit ChainlinkRequested(requestId);
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transferAndCall(oracleAddress(), _amount, encodeTargetRequest(_run)), "Unable to transferAndCall to oracle");
    maliciousRequests += 1;

    return requestId;
  }

  function chainlinkPriceRequest(Chainlink.Run memory _run, uint256 _amount)
    internal
    returns(bytes32 requestId)
  {
    requestId = keccak256(abi.encodePacked(this, maliciousRequests));
    _run.nonce = maliciousRequests;
    maliciousUnfulfilledRequests[requestId] = oracleAddress();
    emit ChainlinkRequested(requestId);
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transferAndCall(oracleAddress(), _amount, encodePriceRequest(_run)), "Unable to transferAndCall to oracle");
    maliciousRequests += 1;

    return requestId;
  }

  function chainlinkWithdrawRequest(MaliciousChainlink.WithdrawRun memory _run, uint256 _wei)
    internal
    returns(bytes32 requestId)
  {
    requestId = keccak256(abi.encodePacked(this, maliciousRequests));
    _run.nonce = maliciousRequests;
    maliciousUnfulfilledRequests[requestId] = oracleAddress();
    emit ChainlinkRequested(requestId);
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transferAndCall(oracleAddress(), _wei, encodeWithdrawRequest(_run)), "Unable to transferAndCall to oracle");
    maliciousRequests += 1;
    return requestId;
  }

  function encodeWithdrawRequest(MaliciousChainlink.WithdrawRun memory _run)
    internal pure returns (bytes memory)
  {
    return abi.encodeWithSelector(
      bytes4(keccak256("withdraw(address,uint256)")),
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.nonce,
      _run.buf.buf);
  }

  function encodeTargetRequest(Chainlink.Run memory _run)
    internal pure returns (bytes memory)
  {
    return abi.encodeWithSelector(
      bytes4(keccak256("requestData(address,uint256,uint256,bytes32,address,bytes4,uint256,bytes)")),
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      1,
      _run.id,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.nonce,
      _run.buf.buf);
  }

  function encodePriceRequest(Chainlink.Run memory _run)
    internal pure returns (bytes memory)
  {
    return abi.encodeWithSelector(
      bytes4(keccak256("requestData(address,uint256,uint256,bytes32,address,bytes4,uint256,bytes)")),
      0, // overridden by onTokenTransfer
      2000000000000000000, // overridden by onTokenTransfer
      1,
      _run.id,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.nonce,
      _run.buf.buf);
  }
}
