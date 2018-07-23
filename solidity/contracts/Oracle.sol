pragma solidity ^0.4.24;

import "./lib/Ownable.sol";
import "./lib/LinkToken.sol";

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

  mapping(uint256 => Callback) private callbacks;

  event RunRequest(
    uint256 indexed internalId,
    bytes32 indexed specId,
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
    // solium-disable-next-line security/no-low-level-calls
    require(address(this).delegatecall(_data), "Unable to create request"); // calls requestData
  }

  function requestData(
    uint256 _version,
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    bytes32 _externalId,
    bytes _data
  )
    public
    onlyLINK
  {
    uint256 internalId = uint256(keccak256(abi.encodePacked(currentSender, _externalId)));
    callbacks[internalId] = Callback(
      _externalId,
      currentAmount,
      _callbackAddress,
      _callbackFunctionId);
    emit RunRequest(internalId, _specId, currentAmount, _version, _data);
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
    withdrawableWei = withdrawableWei.add(callback.amount);
    delete callbacks[_internalId];
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    callback.addr.call(callback.functionId, callback.externalId, _data); // solium-disable-line security/no-low-level-calls
  }

  function withdraw(address _recipient, uint256 _amount)
    public
    onlyOwner
    hasAvailableFunds(_amount)
  {
    withdrawableWei = withdrawableWei.sub(_amount);
    LINK.transfer(_recipient, _amount);
  }

  function cancel(bytes32 _externalId)
    public
  {
    uint256 internalId = uint256(keccak256(abi.encodePacked(msg.sender, _externalId)));
    require(msg.sender == callbacks[internalId].addr, "Must be called from requester");
    Callback memory cb = callbacks[internalId];
    require(LINK.transfer(cb.addr, cb.amount), "Unable to transfer");
    delete callbacks[internalId];
  }

  // MODIFIERS

  modifier hasAvailableFunds(uint256 _amount) {
    require(withdrawableWei >= _amount.add(oneForConsistentGasCost), "Amount requested is greater than withdrawable balance");
    _;
  }

  modifier hasInternalId(uint256 _internalId) {
    require(callbacks[_internalId].addr != address(0), "Must have a valid internalId");
    _;
  }

  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }

}
