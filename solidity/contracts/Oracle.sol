pragma solidity 0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./interfaces/OracleInterface.sol";
import "./interfaces/LinkTokenInterface.sol";

contract Oracle is OracleInterface, Ownable {
  using SafeMath for uint256;

  LinkTokenInterface internal LINK;

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
  mapping(address => bool) private authorizedNodes;

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
    LINK = LinkTokenInterface(_link);
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
    external
    onlyLINK
    checkCallbackAddress(_callbackAddress)
  {
    uint256 internalId = uint256(keccak256(abi.encodePacked(_sender, _externalId)));
    require(callbacks[internalId].externalId != _externalId, "Must use a unique ID");
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
    external
    onlyAuthorizedNode
    hasInternalId(_internalId)
    returns (bool)
  {
    Callback memory callback = callbacks[_internalId];
    withdrawableWei = withdrawableWei.add(callback.amount);
    delete callbacks[_internalId];
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    return callback.addr.call(callback.functionId, callback.externalId, _data); // solium-disable-line security/no-low-level-calls
  }

  function getAuthorizationStatus(address _node) external view returns (bool) {
    return authorizedNodes[_node];
  }

  function setFulfillmentPermission(address _node, bool _allowed) external onlyOwner {
    authorizedNodes[_node] = _allowed;
  }

  function withdraw(address _recipient, uint256 _amount)
    external
    onlyOwner
    hasAvailableFunds(_amount)
  {
    withdrawableWei = withdrawableWei.sub(_amount);
    require(LINK.transfer(_recipient, _amount), "Failed to transfer LINK");
  }

  function withdrawable() external view onlyOwner returns (uint256) {
    return withdrawableWei.sub(oneForConsistentGasCost);
  }

  function cancel(bytes32 _externalId)
    external
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

  modifier onlyAuthorizedNode() {
    require(authorizedNodes[msg.sender] == true || msg.sender == owner, "Not an authorized node to fulfill requests");
    _;
  }

  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }

  modifier permittedFunctionsForLINK() {
    bytes4[1] memory funcSelector;
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      calldatacopy(funcSelector, 132, 4) // grab function selector from calldata
    }
    require(funcSelector[0] == this.requestData.selector, "Must use whitelisted functions");
    _;
  }

  modifier checkCallbackAddress(address _to) {
    require(_to != address(LINK), "Cannot callback to LINK");
    _;
  }

}
