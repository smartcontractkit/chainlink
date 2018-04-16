pragma solidity ^0.4.18;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

contract Oracle is Ownable {

  struct Callback {
    address addr;
    bytes4 functionId;
  }

  mapping(bytes32 => Callback) private callbacks;

  event RunRequest(
    bytes32 indexed id,
    bytes32 indexed jobId,
    uint256 version,
    bytes data
  );

  function requestData(
    uint256 _version,
    bytes32 _jobId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    bytes32 _requestId,
    bytes _data
  )
    public
  {
    Callback memory callback = Callback(_callbackAddress, _callbackFunctionId);
    callbacks[_requestId] = callback;
    emit RunRequest(_requestId, _jobId, _version, _data);
  }

  function fulfillData(bytes32 _requestId, bytes32 _data)
    public
    onlyOwner
    hasRequestId(_requestId)
  {
    Callback memory callback = callbacks[_requestId];
    require(callback.addr.call(callback.functionId, _requestId, _data));
    delete callbacks[_requestId];
  }

  function onTokenTransfer(address _sender, uint _amount, bytes _data)
    public
  {
    if (_data.length > 0) {
      require(address(this).delegatecall(_data));
    }
  }

  modifier hasRequestId(bytes32 _requestId) {
    require(callbacks[_requestId].addr != address(0));
    _;
  }
}
