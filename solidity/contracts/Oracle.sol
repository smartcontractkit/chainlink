pragma solidity ^0.4.23;

import "zeppelin-solidity/contracts/ownership/Ownable.sol";

contract Oracle is Ownable {

  struct Callback {
    bytes32 externalId;
    address addr;
    bytes4 functionId;
  }

  uint256 private currentInternalId = 1;
  mapping(uint256 => Callback) private callbacks;

  event RunRequest(
    uint256 indexed id,
    bytes32 indexed jobId,
    uint256 version,
    bytes data
  );

  function onTokenTransfer(address _sender, uint _amount, bytes _data)
    public
  {
    if (_data.length > 0) {
      require(address(this).delegatecall(_data)); // calls requestData
    }
  }

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
    currentInternalId += 1;
    callbacks[currentInternalId] = Callback(
      _externalId,
      _callbackAddress,
      _callbackFunctionId);
    emit RunRequest(currentInternalId, _jobId, _version, _data);
  }

  function fulfillData(uint256 _internalId, bytes32 _data)
    public
    onlyOwner
    hasInternalId(_internalId)
  {
    Callback memory callback = callbacks[_internalId];
    require(callback.addr.call(callback.functionId, callback.externalId, _data));
    delete callbacks[_internalId];
  }


  modifier hasInternalId(uint256 _internalId) {
    require(callbacks[_internalId].addr != address(0));
    _;
  }
}
