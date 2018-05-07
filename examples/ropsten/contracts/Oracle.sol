pragma solidity ^0.4.23;

import "github.com/OpenZeppelin/openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "github.com/smartcontractkit/LinkToken/contracts/LinkToken.sol";

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
  uint256 private currentInternalId = oneForConsistentGasCost;
  uint256 private currentAmount = oneForConsistentGasCost;
  uint256 private withdrawableWei = oneForConsistentGasCost;

  mapping(uint256 => Callback) private callbacks;

  event RunRequest(
    uint256 indexed id,
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
    currentInternalId += 1;
    callbacks[currentInternalId] = Callback(
      _externalId,
      currentAmount,
      _callbackAddress,
      _callbackFunctionId);
    emit RunRequest(currentInternalId, _jobId, currentAmount, _version, _data);
  }

  function fulfillData(
    uint256 _internalId,
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

  function cancel(uint256 _internalId) public fromCallback(_internalId) {
    Callback memory cb = callbacks[_internalId];
    LINK.transfer(cb.addr, cb.amount);
    delete callbacks[_internalId];
  }

  // MODIFIERS

  modifier fromCallback(uint256 _internalId) {
    require(msg.sender == callbacks[_internalId].addr);
    _;
  }

  modifier hasInternalId(uint256 _internalId) {
    require(callbacks[_internalId].addr != address(0));
    _;
  }

  modifier onlyLINK() {
    require(msg.sender == address(LINK));
    _;
  }

}
