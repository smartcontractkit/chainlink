pragma solidity ^0.4.23;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "./LinkToken.sol";

contract Oracle is Ownable {
  using SafeMath for uint256;

  LinkToken internal LINK;

  struct Callback {
    bytes32 externalId;
    uint256 amount;
    address addr;
    bytes4 functionId;
  }

  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 constant private oneForConsistentGasCost = 1;
  uint256 private withdrawableWei = oneForConsistentGasCost;
  uint256 private currentAmount = oneForConsistentGasCost;
  address private currentSender;

  mapping(bytes32 => Callback) private callbacks;

  event RunRequest(
    bytes32 indexed id,
    bytes32 indexed jobId,
    uint256 indexed amount,
    uint256 version,
    bytes data
  );

  constructor(address _link) Ownable() public {
    LINK = LinkToken(_link);
  }

  function onTokenTransfer(
    address _sender,
    uint256 _wei,
    bytes _data
  )
    public
    onlyLINK
  {
    currentAmount = _wei;
    currentSender = _sender;
    require(address(this).delegatecall(_data)); // calls requestData
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
    onlyLINK
  {
    bytes32 internalId = keccak256(currentSender, _externalId);
    callbacks[internalId] = Callback(
      _externalId,
      currentAmount,
      _callbackAddress,
      _callbackFunctionId);
    emit RunRequest(internalId, _jobId, currentAmount, _version, _data);
  }

  function fulfillData(
    bytes32 _internalId,
    bytes32 _data
  )
    public
    onlyOwner
    hasInternalId(_internalId)
  {
    Callback memory callback = callbacks[_internalId];
    require(callback.addr.call(callback.functionId, callback.externalId, _data));
    withdrawableWei = withdrawableWei.add(callback.amount);
    delete callbacks[_internalId];
  }

  function withdraw() public onlyOwner {
    LINK.transfer(owner, withdrawableWei.sub(oneForConsistentGasCost));
    withdrawableWei = oneForConsistentGasCost;
  }

  function cancel(bytes32 _externalId) public {
    bytes32 internalId = keccak256(msg.sender, _externalId);
    require(msg.sender == callbacks[internalId].addr);
    Callback memory cb = callbacks[internalId];
    LINK.transfer(cb.addr, cb.amount);
    delete callbacks[internalId];
  }

  // MODIFIERS

  modifier hasInternalId(bytes32 _internalId) {
    require(callbacks[_internalId].addr != address(0));
    _;
  }

  modifier onlyLINK() {
    require(msg.sender == address(LINK));
    _;
  }

}
