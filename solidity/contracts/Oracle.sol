pragma solidity 0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./interfaces/ChainlinkRequestInterface.sol";
import "./interfaces/OracleInterface.sol";
import "./interfaces/LinkTokenInterface.sol";

contract Oracle is ChainlinkRequestInterface, OracleInterface, Ownable {
  using SafeMath for uint256;

  uint256 constant public EXPIRY_TIME = 5 minutes;
  uint256 constant private MINIMUM_CONSUMER_GAS_LIMIT = 400000;
  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 constant private ONE_FOR_CONSISTENT_GAS_COST = 1;
  uint256 constant private SELECTOR_LENGTH = 4;
  uint256 constant private EXPECTED_REQUEST_WORDS = 9;
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant private MINIMUM_REQUEST_LENGTH = SELECTOR_LENGTH + (32 * EXPECTED_REQUEST_WORDS);

  struct Callback {
    uint256 amount;
    address addr;
    bytes4 functionId;
    uint64 cancelExpiration;
  }

  LinkTokenInterface internal LINK;
  mapping(bytes32 => Callback) private callbacks;
  mapping(address => bool) private authorizedNodes;
  uint256 private withdrawableTokens = ONE_FOR_CONSISTENT_GAS_COST;

  event RunRequest(
    bytes32 indexed specId,
    address indexed requester,
    uint256 indexed amount,
    uint256 requestId,
    uint256 version,
    bytes data
  );

  event CancelRequest(
    bytes32 requestId
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
    validRequestLength(_data)
    permittedFunctionsForLINK(_data)
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
    uint256 _nonce,
    bytes _data
  )
    external
    onlyLINK
    checkCallbackAddress(_callbackAddress)
  {
    bytes32 requestId = keccak256(abi.encodePacked(_sender, _nonce));
    require(callbacks[requestId].cancelExpiration == 0, "Must use a unique ID");
    callbacks[requestId] = Callback(
      _amount,
      _callbackAddress,
      _callbackFunctionId,
      uint64(now.add(EXPIRY_TIME)));
    emit RunRequest(
      _specId,
      _sender,
      _amount,
      uint256(requestId),
      _version,
      _data);
  }

  function fulfillData(
    uint256 _requestId,
    bytes32 _data
  )
    external
    onlyAuthorizedNode
    isValidRequest(_requestId)
    returns (bool)
  {
    bytes32 requestId = bytes32(_requestId);
    Callback memory callback = callbacks[requestId];
    withdrawableTokens = withdrawableTokens.add(callback.amount);
    delete callbacks[requestId];
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    return callback.addr.call(callback.functionId, requestId, _data); // solium-disable-line security/no-low-level-calls
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
    withdrawableTokens = withdrawableTokens.sub(_amount);
    require(LINK.transfer(_recipient, _amount), "Failed to transfer LINK");
  }

  function withdrawable() external view onlyOwner returns (uint256) {
    return withdrawableTokens.sub(ONE_FOR_CONSISTENT_GAS_COST);
  }

  function cancel(bytes32 _requestId)
    external
  {
    require(msg.sender == callbacks[_requestId].addr, "Must be called from requester");
    require(callbacks[_requestId].cancelExpiration <= now, "Request is not expired");
    Callback memory cb = callbacks[_requestId];
    require(LINK.transfer(cb.addr, cb.amount), "Unable to transfer");
    delete callbacks[_requestId];
    emit CancelRequest(_requestId);
  }

  // MODIFIERS

  modifier hasAvailableFunds(uint256 _amount) {
    require(withdrawableTokens >= _amount.add(ONE_FOR_CONSISTENT_GAS_COST), "Amount requested is greater than withdrawable balance");
    _;
  }

  modifier isValidRequest(uint256 _requestId) {
    require(callbacks[bytes32(_requestId)].addr != address(0), "Must have a valid requestId");
    _;
  }

  modifier onlyAuthorizedNode() {
    require(authorizedNodes[msg.sender] || msg.sender == owner, "Not an authorized node to fulfill requests");
    _;
  }

  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }

  modifier permittedFunctionsForLINK(bytes _data) {
    bytes4 funcSelector;
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      funcSelector := mload(add(_data, 32))
    }
    require(funcSelector == this.requestData.selector, "Must use whitelisted functions");
    _;
  }

  modifier checkCallbackAddress(address _to) {
    require(_to != address(LINK), "Cannot callback to LINK");
    _;
  }

  modifier validRequestLength(bytes _data) {
    require(_data.length >= MINIMUM_REQUEST_LENGTH, "Cannot callback to LINK");
    _;
  }

}
