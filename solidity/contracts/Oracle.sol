pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

contract Oracle is Ownable {

  uint256 private requestId = 1;

  struct Callback {
    address addr;
    bytes32 externalId;
    bytes4 functionId;
  }

  mapping(uint256 => Callback) private callbacks;

  event RunRequest(
    uint256 indexed id,
    bytes32 indexed jobId,
    uint256 version,
    bytes data
  );

  function requestData(
    uint256 _version,
    bytes32 _jobId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    bytes32 _externalId,
    bytes _data
  )
    public
  {
    requestId += 1;
    Callback memory callback = Callback(
      _callbackAddress,
      _externalId,
      _callbackFunctionId);
    callbacks[requestId] = callback;
    emit RunRequest(requestId, _jobId, _version, _data);
  }

  function fulfillData(uint256 _requestId, bytes32 _data)
    public
    onlyOwner
    hasRequestId(_requestId)
  {
    Callback memory callback = callbacks[_requestId];
    require(callback.addr.call(callback.functionId, callback.externalId, _data));
    delete callbacks[_requestId];
  }

  function onTokenTransfer(address _sender, uint _amount, bytes _data)
    public
  {
    require(address(this).delegatecall(_data));
  }

  modifier hasRequestId(uint256 _requestId) {
    require(callbacks[_requestId].addr != address(0));
    _;
  }
}
