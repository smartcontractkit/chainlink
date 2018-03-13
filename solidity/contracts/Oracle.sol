pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

contract Oracle is Ownable {

  struct Callback {
    address addr;
    bytes4 functionId;
  }

  uint256 private requestId;
  mapping(uint256 => Callback) private callbacks;

  event Request(
    uint256 indexed id,
    bytes32 indexed jobId,
    string data
  );

  function requestData(
    bytes32 _jobId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    string _data
  )
    public
    returns (uint256)
  {
    requestId += 1;
    Callback memory callback = Callback(_callbackAddress, _callbackFunctionId);
    callbacks[requestId] = callback;
    Request(requestId, _jobId, _data);
    return requestId;
  }

  function fulfillData(uint256 _requestId, bytes32 _data)
    public
    onlyOwner
    hasRequestId(_requestId)
  {
    Callback memory callback = callbacks[_requestId];
    require(callback.addr.call(callback.functionId, _requestId, _data));
    delete callbacks[_requestId];
  }

  modifier hasRequestId(uint256 _requestId) {
    require(callbacks[_requestId].addr != address(0));
    _;
  }
}
