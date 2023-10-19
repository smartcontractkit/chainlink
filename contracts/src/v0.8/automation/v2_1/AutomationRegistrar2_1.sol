// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {IKeeperRegistryMaster} from "../interfaces/v2_1/IKeeperRegistryMaster.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IERC677Receiver} from "../../shared/interfaces/IERC677Receiver.sol";

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
contract AutomationRegistrar2_1 is TypeAndVersionInterface, ConfirmedOwner, IERC677Receiver {
  /**
   * DISABLED: No auto approvals, all new upkeeps should be approved manually.
   * ENABLED_SENDER_ALLOWLIST: Auto approvals for allowed senders subject to max allowed. Manual for rest.
   * ENABLED_ALL: Auto approvals for all new upkeeps subject to max allowed.
   */
  enum AutoApproveType {
    DISABLED,
    ENABLED_SENDER_ALLOWLIST,
    ENABLED_ALL
  }

  bytes4 private constant REGISTER_REQUEST_SELECTOR = this.register.selector;

  mapping(bytes32 => PendingRequest) private s_pendingRequests;
  mapping(uint8 => TriggerRegistrationStorage) private s_triggerRegistrations;

  LinkTokenInterface public immutable LINK;

  /**
   * @notice versions:
   * - KeeperRegistrar 2.1.0: Update for compatability with registry 2.1.0
   *                          Add auto approval levels by type
   * - KeeperRegistrar 2.0.0: Remove source from register
   *                          Breaks our example of "Register an Upkeep using your own deployed contract"
   * - KeeperRegistrar 1.1.0: Add functionality for sender allowlist in auto approve
   *                        : Remove rate limit and add max allowed for auto approve
   * - KeeperRegistrar 1.0.0: initial release
   */
  string public constant override typeAndVersion = "AutomationRegistrar 2.1.0";

  /**
   * @notice TriggerRegistrationStorage stores the auto-approval levels for upkeeps by type
   * @member autoApproveType the auto approval setting (see enum)
   * @member autoApproveMaxAllowed the max number of upkeeps that can be auto approved of this type
   * @member approvedCount the count of upkeeps auto approved of this type
   */
  struct TriggerRegistrationStorage {
    AutoApproveType autoApproveType;
    uint32 autoApproveMaxAllowed;
    uint32 approvedCount;
  }

  /**
   * @notice InitialTriggerConfig configures the auto-approval levels for upkeeps by trigger type
   * @dev this struct is only used in the constructor to set the initial values for various trigger configs
   * @member triggerType the upkeep type to configure
   * @member autoApproveType the auto approval setting (see enum)
   * @member autoApproveMaxAllowed the max number of upkeeps that can be auto approved of this type
   */
  struct InitialTriggerConfig {
    uint8 triggerType;
    AutoApproveType autoApproveType;
    uint32 autoApproveMaxAllowed;
  }

  struct RegistrarConfig {
    IKeeperRegistryMaster keeperRegistry;
    uint96 minLINKJuels;
  }

  struct PendingRequest {
    address admin;
    uint96 balance;
  }

  struct RegistrationParams {
    string name;
    bytes encryptedEmail;
    address upkeepContract;
    uint32 gasLimit;
    address adminAddress;
    uint8 triggerType;
    bytes checkData;
    bytes triggerConfig;
    bytes offchainConfig;
    uint96 amount;
  }

  RegistrarConfig private s_config;
  // Only applicable if s_config.configType is ENABLED_SENDER_ALLOWLIST
  mapping(address => bool) private s_autoApproveAllowedSenders;

  event RegistrationRequested(
    bytes32 indexed hash,
    string name,
    bytes encryptedEmail,
    address indexed upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    uint8 triggerType,
    bytes triggerConfig,
    bytes offchainConfig,
    bytes checkData,
    uint96 amount
  );

  event RegistrationApproved(bytes32 indexed hash, string displayName, uint256 indexed upkeepId);

  event RegistrationRejected(bytes32 indexed hash);

  event AutoApproveAllowedSenderSet(address indexed senderAddress, bool allowed);

  event ConfigChanged(address keeperRegistry, uint96 minLINKJuels);

  event TriggerConfigSet(uint8 triggerType, AutoApproveType autoApproveType, uint32 autoApproveMaxAllowed);

  error InvalidAdminAddress();
  error RequestNotFound();
  error HashMismatch();
  error OnlyAdminOrOwner();
  error InsufficientPayment();
  error RegistrationRequestFailed();
  error OnlyLink();
  error AmountMismatch();
  error SenderMismatch();
  error FunctionNotPermitted();
  error LinkTransferFailed(address to);
  error InvalidDataLength();

  /**
   * @param LINKAddress Address of Link token
   * @param keeperRegistry keeper registry address
   * @param minLINKJuels minimum LINK that new registrations should fund their upkeep with
   * @param triggerConfigs the initial config for individual triggers
   */
  constructor(
    address LINKAddress,
    address keeperRegistry,
    uint96 minLINKJuels,
    InitialTriggerConfig[] memory triggerConfigs
  ) ConfirmedOwner(msg.sender) {
    LINK = LinkTokenInterface(LINKAddress);
    setConfig(keeperRegistry, minLINKJuels);
    for (uint256 idx = 0; idx < triggerConfigs.length; idx++) {
      setTriggerConfig(
        triggerConfigs[idx].triggerType,
        triggerConfigs[idx].autoApproveType,
        triggerConfigs[idx].autoApproveMaxAllowed
      );
    }
  }

  //EXTERNAL

  /**
   * @notice register can only be called through transferAndCall on LINK contract
   * @param name string of the upkeep to be registered
   * @param encryptedEmail email address of upkeep contact
   * @param upkeepContract address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when performing upkeep
   * @param adminAddress address to cancel upkeep and withdraw remaining funds
   * @param triggerType the type of trigger for the upkeep
   * @param checkData data passed to the contract when checking for upkeep
   * @param triggerConfig the config for the trigger
   * @param offchainConfig offchainConfig for upkeep in bytes
   * @param amount quantity of LINK upkeep is funded with (specified in Juels)
   * @param sender address of the sender making the request
   */
  function register(
    string memory name,
    bytes calldata encryptedEmail,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    uint8 triggerType,
    bytes memory checkData,
    bytes memory triggerConfig,
    bytes memory offchainConfig,
    uint96 amount,
    address sender
  ) external onlyLINK {
    _register(
      RegistrationParams({
        name: name,
        encryptedEmail: encryptedEmail,
        upkeepContract: upkeepContract,
        gasLimit: gasLimit,
        adminAddress: adminAddress,
        triggerType: triggerType,
        checkData: checkData,
        triggerConfig: triggerConfig,
        offchainConfig: offchainConfig,
        amount: amount
      }),
      sender
    );
  }

  /**
   * @notice Allows external users to register upkeeps; assumes amount is approved for transfer by the contract
   * @param requestParams struct of all possible registration parameters
   */
  function registerUpkeep(RegistrationParams calldata requestParams) external returns (uint256) {
    if (requestParams.amount < s_config.minLINKJuels) {
      revert InsufficientPayment();
    }

    LINK.transferFrom(msg.sender, address(this), requestParams.amount);

    return _register(requestParams, msg.sender);
  }

  /**
   * @dev register upkeep on KeeperRegistry contract and emit RegistrationApproved event
   */
  function approve(
    string memory name,
    address upkeepContract,
    uint32 gasLimit,
    address adminAddress,
    uint8 triggerType,
    bytes calldata checkData,
    bytes memory triggerConfig,
    bytes calldata offchainConfig,
    bytes32 hash
  ) external onlyOwner {
    PendingRequest memory request = s_pendingRequests[hash];
    if (request.admin == address(0)) {
      revert RequestNotFound();
    }
    bytes32 expectedHash = keccak256(
      abi.encode(upkeepContract, gasLimit, adminAddress, triggerType, checkData, triggerConfig, offchainConfig)
    );
    if (hash != expectedHash) {
      revert HashMismatch();
    }
    delete s_pendingRequests[hash];
    _approve(
      RegistrationParams({
        name: name,
        encryptedEmail: "",
        upkeepContract: upkeepContract,
        gasLimit: gasLimit,
        adminAddress: adminAddress,
        triggerType: triggerType,
        checkData: checkData,
        triggerConfig: triggerConfig,
        offchainConfig: offchainConfig,
        amount: request.balance
      }),
      expectedHash
    );
  }

  /**
   * @notice cancel will remove a registration request and return the refunds to the request.admin
   * @param hash the request hash
   */
  function cancel(bytes32 hash) external {
    PendingRequest memory request = s_pendingRequests[hash];
    if (!(msg.sender == request.admin || msg.sender == owner())) {
      revert OnlyAdminOrOwner();
    }
    if (request.admin == address(0)) {
      revert RequestNotFound();
    }
    delete s_pendingRequests[hash];
    bool success = LINK.transfer(request.admin, request.balance);
    if (!success) {
      revert LinkTransferFailed(request.admin);
    }
    emit RegistrationRejected(hash);
  }

  /**
   * @notice owner calls this function to set contract config
   * @param keeperRegistry new keeper registry address
   * @param minLINKJuels minimum LINK that new registrations should fund their upkeep with
   */
  function setConfig(address keeperRegistry, uint96 minLINKJuels) public onlyOwner {
    s_config = RegistrarConfig({minLINKJuels: minLINKJuels, keeperRegistry: IKeeperRegistryMaster(keeperRegistry)});
    emit ConfigChanged(keeperRegistry, minLINKJuels);
  }

  /**
   * @notice owner calls to set the config for this upkeep type
   * @param triggerType the upkeep type to configure
   * @param autoApproveType the auto approval setting (see enum)
   * @param autoApproveMaxAllowed the max number of upkeeps that can be auto approved of this type
   */
  function setTriggerConfig(
    uint8 triggerType,
    AutoApproveType autoApproveType,
    uint32 autoApproveMaxAllowed
  ) public onlyOwner {
    s_triggerRegistrations[triggerType].autoApproveType = autoApproveType;
    s_triggerRegistrations[triggerType].autoApproveMaxAllowed = autoApproveMaxAllowed;
    emit TriggerConfigSet(triggerType, autoApproveType, autoApproveMaxAllowed);
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
  function getConfig() external view returns (address keeperRegistry, uint256 minLINKJuels) {
    RegistrarConfig memory config = s_config;
    return (address(config.keeperRegistry), config.minLINKJuels);
  }

  /**
   * @notice read the config for this upkeep type
   * @param triggerType upkeep type to read config for
   */
  function getTriggerRegistrationDetails(uint8 triggerType) external view returns (TriggerRegistrationStorage memory) {
    return s_triggerRegistrations[triggerType];
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
  )
    external
    override
    onlyLINK
    permittedFunctionsForLINK(data)
    isActualAmount(amount, data)
    isActualSender(sender, data)
  {
    if (amount < s_config.minLINKJuels) {
      revert InsufficientPayment();
    }
    (bool success, ) = address(this).delegatecall(data);
    // calls register
    if (!success) {
      revert RegistrationRequestFailed();
    }
  }

  // ================================================================
  // |                           PRIVATE                            |
  // ================================================================

  /**
   * @dev verify registration request and emit RegistrationRequested event
   */
  function _register(RegistrationParams memory params, address sender) private returns (uint256) {
    if (params.adminAddress == address(0)) {
      revert InvalidAdminAddress();
    }
    bytes32 hash = keccak256(
      abi.encode(
        params.upkeepContract,
        params.gasLimit,
        params.adminAddress,
        params.triggerType,
        params.checkData,
        params.triggerConfig,
        params.offchainConfig
      )
    );

    emit RegistrationRequested(
      hash,
      params.name,
      params.encryptedEmail,
      params.upkeepContract,
      params.gasLimit,
      params.adminAddress,
      params.triggerType,
      params.triggerConfig,
      params.offchainConfig,
      params.checkData,
      params.amount
    );

    uint256 upkeepId;
    if (_shouldAutoApprove(s_triggerRegistrations[params.triggerType], sender)) {
      s_triggerRegistrations[params.triggerType].approvedCount++;
      upkeepId = _approve(params, hash);
    } else {
      uint96 newBalance = s_pendingRequests[hash].balance + params.amount;
      s_pendingRequests[hash] = PendingRequest({admin: params.adminAddress, balance: newBalance});
    }

    return upkeepId;
  }

  /**
   * @dev register upkeep on KeeperRegistry contract and emit RegistrationApproved event
   */
  function _approve(RegistrationParams memory params, bytes32 hash) private returns (uint256) {
    IKeeperRegistryMaster keeperRegistry = s_config.keeperRegistry;
    uint256 upkeepId = keeperRegistry.registerUpkeep(
      params.upkeepContract,
      params.gasLimit,
      params.adminAddress,
      params.triggerType,
      params.checkData,
      params.triggerConfig,
      params.offchainConfig
    );
    bool success = LINK.transferAndCall(address(keeperRegistry), params.amount, abi.encode(upkeepId));
    if (!success) {
      revert LinkTransferFailed(address(keeperRegistry));
    }
    emit RegistrationApproved(hash, params.name, upkeepId);
    return upkeepId;
  }

  /**
   * @dev verify sender allowlist if needed and check max limit
   */
  function _shouldAutoApprove(TriggerRegistrationStorage memory config, address sender) private view returns (bool) {
    if (config.autoApproveType == AutoApproveType.DISABLED) {
      return false;
    }
    if (config.autoApproveType == AutoApproveType.ENABLED_SENDER_ALLOWLIST && (!s_autoApproveAllowedSenders[sender])) {
      return false;
    }
    if (config.approvedCount < config.autoApproveMaxAllowed) {
      return true;
    }
    return false;
  }

  // ================================================================
  // |                          MODIFIERS                           |
  // ================================================================

  /**
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    if (msg.sender != address(LINK)) {
      revert OnlyLink();
    }
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
      funcSelector := mload(add(_data, 32)) // First 32 bytes contain length of data
    }
    if (funcSelector != REGISTER_REQUEST_SELECTOR) {
      revert FunctionNotPermitted();
    }
    _;
  }

  /**
   * @dev Reverts if the actual amount passed does not match the expected amount
   * @param expected amount that should match the actual amount
   * @param data bytes
   */
  modifier isActualAmount(uint256 expected, bytes calldata data) {
    // decode register function arguments to get actual amount
    (, , , , , , , , , uint96 amount, ) = abi.decode(
      data[4:],
      (string, bytes, address, uint32, address, uint8, bytes, bytes, bytes, uint96, address)
    );
    if (expected != amount) {
      revert AmountMismatch();
    }
    _;
  }

  /**
   * @dev Reverts if the actual sender address does not match the expected sender address
   * @param expected address that should match the actual sender address
   * @param data bytes
   */
  modifier isActualSender(address expected, bytes calldata data) {
    // decode register function arguments to get actual sender
    (, , , , , , , , , , address sender) = abi.decode(
      data[4:],
      (string, bytes, address, uint32, address, uint8, bytes, bytes, bytes, uint96, address)
    );
    if (expected != sender) {
      revert SenderMismatch();
    }
    _;
  }
}
