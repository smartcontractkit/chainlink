pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "link_token/contracts/LinkToken.sol";

contract Oracle is Ownable {
  using SafeMath for uint256;

  LinkToken internal LINK;

  struct Callback {
    bytes32 externalId;
    uint256 amount;
    address addr;
    bytes4 functionId;
    uint64 cancelExpiration;
  }

  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 constant private oneForConsistentGasCost = 1;
  uint256 private withdrawableWei = oneForConsistentGasCost;

  mapping(uint256 => Callback) private callbacks;

  event RunRequest(
    bytes32 indexed specId,
    address indexed requester,
    uint256 indexed amount,
    uint256 internalId,
    uint256 version,
    bytes data
  );

  event CancelRequest(
    uint256 internalId
  );

  constructor(address _link) Ownable() public {
    LINK = LinkToken(_link);
  }

  function onTokenTransfer(
    address _sender,
    uint256 _amount,
    bytes _data
  )
    public
    onlyLINK
    permittedFunctionsForLINK
  {
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      mstore(add(_data, 36), _sender) // ensure correct sender is passed
      // solium-disable-next-line security/no-low-level-calls
      mstore(add(_data, 68), _amount)    // ensure correct amount is passed
    }
    // solium-disable-next-line security/no-low-level-calls
    require(address(this).delegatecall(_data), "Unable to create request"); // calls requestData
  }

  function requestData(
    address _sender,
    uint256 _amount,
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
    uint256 internalId = uint256(keccak256(abi.encodePacked(_sender, _externalId)));
    callbacks[internalId] = Callback(
      _externalId,
      _amount,
      _callbackAddress,
      _callbackFunctionId,
      uint64(now.add(5 minutes)));
    emit RunRequest(
      _specId,
      _sender,
      _amount,
      internalId,
      _version,
      _data);
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
    require(callbacks[internalId].cancelExpiration <= now, "Request is not expired");
    Callback memory cb = callbacks[internalId];
    require(LINK.transfer(cb.addr, cb.amount), "Unable to transfer");
    delete callbacks[internalId];
    emit CancelRequest(internalId);
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

  bytes4 constant private permittedFunc = bytes4(keccak256("requestData(address,uint256,uint256,bytes32,address,bytes4,bytes32,bytes)"));

  modifier permittedFunctionsForLINK() {
    bytes4[1] memory funcSelector;
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      calldatacopy(funcSelector, 132, 4) // grab function selector from calldata
    }
    require(funcSelector[0] == permittedFunc, "Must use whitelisted functions");
    _;
  }

}
