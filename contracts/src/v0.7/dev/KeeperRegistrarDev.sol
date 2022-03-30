// SPDX-License-Identifier: MIT
/*
 * This is a development version of UpkeepRegistrationRequests (soon to be renamed to KeeperRegistrar).
 * Once this is audited and finalised it will be copied to KeeperRegistrar
 */
pragma solidity ^0.7.0;
pragma abicoder v2;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/KeeperRegistryInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../vendor/SafeMath96.sol";
import "../ConfirmedOwner.sol";

/**
 * @notice Contract to accept requests for upkeep registrations
 * @dev There are 2 registration workflows in this contract
 * Flow 1. auto approve OFF / manual registration - UI calls `register` function on this contract, this contract owner at a later time then manually
 *  calls `approve` to register upkeep and emit events to inform UI and others interested.
 * Flow 2. auto approve ON / real time registration - UI calls `register` function as before, which calls the `registerUpkeep` function directly on
 *  keeper registry and then emits approved event to finish the flow automatically without manual intervention.
 * The idea is to have same interface(functions,events) for UI or anyone using this contract irrespective of auto approve being enabled or not.
 * they can just listen to `RegistrationRequested` & `RegistrationApproved` events and know the status on registrations.
 */
contract KeeperRegistrarDev is TypeAndVersionInterface, ConfirmedOwner {
  /**
   * DISABLED: No auto approvals, all new upkeeps should be approved manually.
   * ENABLED_SENDER_ALLOWLIST: Auto approvals for allowed senders subject to rate limits. Manual for rest.
   * ENABLED_ALL: Auto approvals for all new upkeeps subject to rate limits.
   */
  enum AutoApproveType {
    DISABLED,
    ENABLED_SENDER_ALLOWLIST,
    ENABLED_ALL
  }

  using SafeMath96 for uint96;

  bytes4 private constant REGISTER_REQUEST_SELECTOR = this.register.selector;

  uint256 private s_minLINKJuels;
  mapping(bytes32 => PendingRequest) private s_pendingRequests;

  LinkTokenInterface public immutable LINK;

  /**
   * @notice versions:
   * - KeeperRegistrar 1.1.0: Add functionality for sender allowlist in auto approve
   * - KeeperRegistrar 1.0.0: initial release
   */
  string public constant override typeAndVersion = "KeeperRegistrar 1.1.0";

  struct AutoApprovedConfig {
    AutoApproveType configType;
    uint16 allowedPerWindow;
    uint32 windowSizeInBlocks;
    uint64 windowStart;
    uint16 approvedInCurrentWindow;
  }

  struct PendingRequest {
    address admin;
    uint96 balance;
  }

  AutoApprovedConfig private s_config;
  // Only applicable if s_config.configType is ENABLED_SENDER_ALLOWLIST
  mapping(address => bool) private s_autoApproveAllowedSenders;
  KeeperRegistryBaseInterface private s_keeperRegistry;

  event RegistrationRequested(
    bytes32 indexed hash,
    string name,
    bytes encryptedEmail,
    address indexed upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes checkData,
    uint96 amount,
    uint8 indexed source
  );

  event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId);

  event RegistrationRejected(bytes32 indexed hash);

  event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed);

  event ConfigChanged(
    AutoApproveType configType,
    uint32 windowSizeInBlocks,
    uint16 allowedPerWindow,
    address keeperRegistry,
    uint256 minLINKJuels
  );

  constructor(address LINKAddress, uint256 minimumLINKJuels) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(LINKAddress);
    s_minLINKJuels = minimumLINKJuels;
  }

  //EXTERNAL

  /**
   * @notice register can only be called through transferAndCall on LINK contract
   * @param name string of the upkeep to be registered
   * @param encryptedEmail email address of upkeep contact
   * @param upkeepContract address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when performing upkeep
   * @param adminAddress address to cancel upkeep and withdraw remaining funds
   * @param checkData data passed to the contract when checking for upkeep
   * @param amount quantity of LINK upkeep is funded with (specified in Juels)
   * @param source application sending this request
   * @param sender address of the sender making the request
   */
  function register(
    string memory name,
    bytes calldata encryptedEmail,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes calldata checkData,
    uint96 amount,
    uint8 source,
    address sender
  ) external onlyLINK {
    require(adminAddress != address(0), "invalid admin address");
    bytes32 hash = keccak256(abi.encode(upkeepContract, gasLimit, adminAddress, checkData));

    emit RegistrationRequested(
      hash,
      name,
      encryptedEmail,
      upkeepContract,
      gasLimit,
      adminAddress,
      checkData,
      amount,
      source
    );

    AutoApprovedConfig memory config = s_config;
    if (_shouldAutoApprove(config, sender)) {
      _incrementApprovedCount(config);

      _approve(name, upkeepContract, gasLimit, adminAddress, checkData, amount, hash);
    } else {
      uint96 newBalance = s_pendingRequests[hash].balance.add(amount);
      s_pendingRequests[hash] = PendingRequest({admin: adminAddress, balance: newBalance});
    }
  }

  /**
   * @dev register upkeep on KeeperRegistry contract and emit RegistrationApproved event
   */
  function approve(
    string memory name,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes calldata checkData,
    bytes32 hash
  ) external onlyOwner {
    PendingRequest memory request = s_pendingRequests[hash];
    require(request.admin != address(0), "request not found");
    bytes32 expectedHash = keccak256(abi.encode(upkeepContract, gasLimit, adminAddress, checkData));
    require(hash == expectedHash, "hash and payload do not match");
    delete s_pendingRequests[hash];
    _approve(name, upkeepContract, gasLimit, adminAddress, checkData, request.balance, hash);
  }

  /**
   * @notice cancel will remove a registration request and return the refunds to the msg.sender
   * @param hash the request hash
   */
  function cancel(bytes32 hash) external {
    PendingRequest memory request = s_pendingRequests[hash];
    require(msg.sender == request.admin || msg.sender == owner(), "only admin / owner can cancel");
    require(request.admin != address(0), "request not found");
    delete s_pendingRequests[hash];
    require(LINK.transfer(msg.sender, request.balance), "LINK token transfer failed");
    emit RegistrationRejected(hash);
  }

  /**
   * @notice owner calls this function to set if registration requests should be sent directly to the Keeper Registry
   * @param configType setting for auto-approve registrations
   *                   note: autoApproveAllowedSenders list persists across config changes irrespective of type
   * @param windowSizeInBlocks window size defined in number of blocks
   * @param allowedPerWindow number of registrations that can be auto approved in above window
   * @param keeperRegistry new keeper registry address
   */
  function setRegistrationConfig(
    AutoApproveType configType,
    uint32 windowSizeInBlocks,
    uint16 allowedPerWindow,
    address keeperRegistry,
    uint256 minLINKJuels
  ) external onlyOwner {
    s_config = AutoApprovedConfig({
      configType: configType,
      allowedPerWindow: allowedPerWindow,
      windowSizeInBlocks: windowSizeInBlocks,
      windowStart: 0,
      approvedInCurrentWindow: 0
    });
    s_minLINKJuels = minLINKJuels;
    s_keeperRegistry = KeeperRegistryBaseInterface(keeperRegistry);

    emit ConfigChanged(configType, windowSizeInBlocks, allowedPerWindow, keeperRegistry, minLINKJuels);
  }

  /**
   * @notice owner calls this function to set allowlist status for senderAddress
   * @param senderAddress senderAddress to set the allowlist status for
   * @param allowed true if senderAddress needs to be added to allowlist, false if needs to be removed
   */
  function setAutoApproveAllowedSender(address senderAddress, bool allowed) external onlyOwner {
    s_autoApproveAllowedSenders[senderAddress] = allowed;

    emit AutoApproveAllowedSenderSet(senderAddress, allowed);
  }

  /**
   * @notice read the allowlist status of senderAddress
   * @param senderAddress address to read the allowlist status for
   */
  function getAutoApproveAllowedSender(address senderAddress) external view returns (bool) {
    return s_autoApproveAllowedSenders[senderAddress];
  }

  /**
   * @notice read the current registration configuration
   */
  function getRegistrationConfig()
    external
    view
    returns (
      AutoApproveType configType,
      uint32 windowSizeInBlocks,
      uint16 allowedPerWindow,
      address keeperRegistry,
      uint256 minLINKJuels,
      uint64 windowStart,
      uint16 approvedInCurrentWindow
    )
  {
    AutoApprovedConfig memory config = s_config;
    return (
      config.configType,
      config.windowSizeInBlocks,
      config.allowedPerWindow,
      address(s_keeperRegistry),
      s_minLINKJuels,
      config.windowStart,
      config.approvedInCurrentWindow
    );
  }

  /**
   * @notice gets the admin address and the current balance of a registration request
   */
  function getPendingRequest(bytes32 hash) external view returns (address, uint96) {
    PendingRequest memory request = s_pendingRequests[hash];
    return (request.admin, request.balance);
  }

  /**
   * @notice Called when LINK is sent to the contract via `transferAndCall`
   * @param sender Address of the sender transfering LINK
   * @param amount Amount of LINK sent (specified in Juels)
   * @param data Payload of the transaction
   */
  function onTokenTransfer(
    address sender,
    uint256 amount,
    bytes calldata data
  ) external onlyLINK permittedFunctionsForLINK(data) isActualAmount(amount, data) isActualSender(sender, data) {
    require(amount >= s_minLINKJuels, "Insufficient payment");
    (bool success, ) = address(this).delegatecall(data);
    // calls register
    require(success, "Unable to create request");
  }

  //PRIVATE

  /**
   * @dev reset auto approve window if passed end of current window
   */
  function _resetWindowIfRequired(AutoApprovedConfig memory config) private {
    uint64 blocksPassed = uint64(block.number - config.windowStart);
    if (blocksPassed >= config.windowSizeInBlocks) {
      config.windowStart = uint64(block.number);
      config.approvedInCurrentWindow = 0;
      s_config = config;
    }
  }

  /**
   * @dev register upkeep on KeeperRegistry contract and emit RegistrationApproved event
   */
  function _approve(
    string memory name,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    bytes calldata checkData,
    uint96 amount,
    bytes32 hash
  ) private {
    KeeperRegistryBaseInterface keeperRegistry = s_keeperRegistry;

    // register upkeep
    uint256 upkeepId = keeperRegistry.registerUpkeep(upkeepContract, gasLimit, adminAddress, checkData);
    // fund upkeep
    bool success = LINK.transferAndCall(address(keeperRegistry), amount, abi.encode(upkeepId));
    require(success, "failed to fund upkeep");

    emit RegistrationApproved(hash, name, upkeepId);
  }

  /**
   * @dev verify sender allowlist if needed and check rate limits
   */
  function _shouldAutoApprove(AutoApprovedConfig memory config, address sender) private returns (bool) {
    if (config.configType == AutoApproveType.DISABLED) {
      return false;
    }
    if (config.configType == AutoApproveType.ENABLED_SENDER_ALLOWLIST && (!s_autoApproveAllowedSenders[sender])) {
      return false;
    }
    _resetWindowIfRequired(config);
    if (config.approvedInCurrentWindow < config.allowedPerWindow) {
      return true;
    }
    return false;
  }

  /**
   * @dev record new latest approved count
   */
  function _incrementApprovedCount(AutoApprovedConfig memory config) private {
    config.approvedInCurrentWindow++;
    s_config = config;
  }

  //MODIFIERS

  /**
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }

  /**
   * @dev Reverts if the given data does not begin with the `register` function selector
   * @param _data The data payload of the request
   */
  modifier permittedFunctionsForLINK(bytes memory _data) {
    bytes4 funcSelector;
    assembly {
      // solhint-disable-next-line avoid-low-level-calls
      funcSelector := mload(add(_data, 32))
    }
    require(funcSelector == REGISTER_REQUEST_SELECTOR, "Must use whitelisted functions");
    _;
  }

  /**
   * @dev Reverts if the actual amount passed does not match the expected amount
   * @param expected amount that should match the actual amount
   * @param data bytes
   */
  modifier isActualAmount(uint256 expected, bytes memory data) {
    uint256 actual;
    assembly {
      actual := mload(add(data, 228))
    }
    require(expected == actual, "Amount mismatch");
    _;
  }

  /**
   * @dev Reverts if the actual sender address does not match the expected sender address
   * @param expected address that should match the actual sender address
   * @param data bytes
   */
  modifier isActualSender(address expected, bytes memory data) {
    address actual;
    assembly {
      actual := mload(add(data, 292))
    }
    require(expected == actual, "Sender address mismatch");
    _;
  }
}
