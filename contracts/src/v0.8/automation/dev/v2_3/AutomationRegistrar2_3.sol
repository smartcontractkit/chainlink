// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {IERC677Receiver} from "../../../shared/interfaces/IERC677Receiver.sol";
import {IERC20Metadata as IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IWrappedNative} from "../interfaces/v2_3/IWrappedNative.sol";
import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

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
contract AutomationRegistrar2_3 is TypeAndVersionInterface, ConfirmedOwner, IERC677Receiver {
  using SafeERC20 for IERC20;

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

  /**
   * @notice versions:
   * - KeeperRegistrar 2.3.0: Update for compatability with registry 2.3.0
   *                          Add native billing and ERC20 billing support
   * - KeeperRegistrar 2.1.0: Update for compatability with registry 2.1.0
   *                          Add auto approval levels by type
   * - KeeperRegistrar 2.0.0: Remove source from register
   *                          Breaks our example of "Register an Upkeep using your own deployed contract"
   * - KeeperRegistrar 1.1.0: Add functionality for sender allowlist in auto approve
   *                        : Remove rate limit and add max allowed for auto approve
   * - KeeperRegistrar 1.0.0: initial release
   */
  string public constant override typeAndVersion = "AutomationRegistrar 2.3.0";

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
  // solhint-disable-next-line gas-struct-packing
  struct InitialTriggerConfig {
    uint8 triggerType;
    AutoApproveType autoApproveType;
    uint32 autoApproveMaxAllowed;
  }

  struct PendingRequest {
    address admin;
    uint96 balance;
    IERC20 billingToken;
  }
  /**
   * @member upkeepContract address to perform upkeep on
   * @member amount quantity of billing token upkeep is funded with (specified in the billing token's decimals)
   * @member adminAddress address to cancel upkeep and withdraw remaining funds
   * @member gasLimit amount of gas to provide the target contract when performing upkeep
   * @member triggerType the type of trigger for the upkeep
   * @member billingToken the token to pay with
   * @member name string of the upkeep to be registered
   * @member encryptedEmail email address of upkeep contact
   * @member checkData data passed to the contract when checking for upkeep
   * @member triggerConfig the config for the trigger
   * @member offchainConfig offchainConfig for upkeep in bytes
   */
  struct RegistrationParams {
    address upkeepContract;
    uint96 amount;
    // 1 word full
    address adminAddress;
    uint32 gasLimit;
    uint8 triggerType;
    // 7 bytes left in 2nd word
    IERC20 billingToken;
    // 12 bytes left in 3rd word
    string name;
    bytes encryptedEmail;
    bytes checkData;
    bytes triggerConfig;
    bytes offchainConfig;
  }

  LinkTokenInterface public immutable i_LINK;
  IWrappedNative public immutable i_WRAPPED_NATIVE_TOKEN;
  IAutomationRegistryMaster2_3 private s_registry;

  // Only applicable if trigger config is set to ENABLED_SENDER_ALLOWLIST
  mapping(address => bool) private s_autoApproveAllowedSenders;
  mapping(IERC20 => uint256) private s_minRegistrationAmounts;
  mapping(bytes32 => PendingRequest) private s_pendingRequests;
  mapping(uint8 => TriggerRegistrationStorage) private s_triggerRegistrations;

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

  event ConfigChanged();

  event TriggerConfigSet(uint8 triggerType, AutoApproveType autoApproveType, uint32 autoApproveMaxAllowed);

  error HashMismatch();
  error InsufficientPayment();
  error InvalidAdminAddress();
  error InvalidBillingToken();
  error InvalidDataLength();
  error TransferFailed(address to);
  error DuplicateEntry();
  error OnlyAdminOrOwner();
  error OnlyLink();
  error RequestNotFound();

  /**
   * @param LINKAddress Address of Link token
   * @param registry keeper registry address
   * @param triggerConfigs the initial config for individual triggers
   * @param billingTokens the tokens allowed for billing
   * @param minRegistrationFees the minimum amount for registering with each billing token
   * @param wrappedNativeToken wrapped native token
   */
  constructor(
    address LINKAddress,
    IAutomationRegistryMaster2_3 registry,
    InitialTriggerConfig[] memory triggerConfigs,
    IERC20[] memory billingTokens,
    uint256[] memory minRegistrationFees,
    IWrappedNative wrappedNativeToken
  ) ConfirmedOwner(msg.sender) {
    i_LINK = LinkTokenInterface(LINKAddress);
    i_WRAPPED_NATIVE_TOKEN = wrappedNativeToken;
    setConfig(registry, billingTokens, minRegistrationFees);
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
   * @notice Allows external users to register upkeeps; assumes amount is approved for transfer by the contract
   * @param requestParams struct of all possible registration parameters
   */
  function registerUpkeep(RegistrationParams memory requestParams) external payable returns (uint256) {
    if (requestParams.billingToken == IERC20(i_WRAPPED_NATIVE_TOKEN) && msg.value != 0) {
      requestParams.amount = SafeCast.toUint96(msg.value);
      // wrap and send native payment
      i_WRAPPED_NATIVE_TOKEN.deposit{value: msg.value}();
    } else {
      // send ERC20 payment, including wrapped native token
      requestParams.billingToken.safeTransferFrom(msg.sender, address(this), requestParams.amount);
    }

    return _register(requestParams, msg.sender);
  }

  /**
   * @dev register upkeep on AutomationRegistry contract and emit RegistrationApproved event
   * @param requestParams struct of all possible registration parameters
   */
  function approve(RegistrationParams calldata requestParams) external onlyOwner {
    bytes32 hash = keccak256(abi.encode(requestParams));

    PendingRequest memory request = s_pendingRequests[hash];
    if (request.admin == address(0)) {
      revert RequestNotFound();
    }

    delete s_pendingRequests[hash];
    _approve(requestParams, hash);
  }

  /**
   * @notice cancel will remove a registration request from the pending request queue and return the refunds to the request.admin
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

    request.billingToken.safeTransfer(request.admin, request.balance);

    emit RegistrationRejected(hash);
  }

  /**
   * @notice owner calls this function to set contract config
   * @param registry new keeper registry address
   * @param billingTokens the billing tokens that this registrar supports (registy must also support these)
   * @param minBalances minimum balances that users must supply to register with the corresponding billing token
   */
  function setConfig(
    IAutomationRegistryMaster2_3 registry,
    IERC20[] memory billingTokens,
    uint256[] memory minBalances
  ) public onlyOwner {
    if (billingTokens.length != minBalances.length) revert InvalidDataLength();
    s_registry = registry;
    for (uint256 i = 0; i < billingTokens.length; i++) {
      s_minRegistrationAmounts[billingTokens[i]] = minBalances[i];
    }
    emit ConfigChanged();
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
   * @notice gets the registry that this registrar is pointed to
   */
  function getRegistry() external view returns (IAutomationRegistryMaster2_3) {
    return s_registry;
  }

  /**
   * @notice get the minimum registration fee for a particular billing token
   */
  function getMinimumRegistrationAmount(IERC20 billingToken) external view returns (uint256) {
    return s_minRegistrationAmounts[billingToken];
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
  function onTokenTransfer(address sender, uint256 amount, bytes calldata data) external override {
    if (msg.sender != address(i_LINK)) revert OnlyLink();
    RegistrationParams memory params = abi.decode(data, (RegistrationParams));
    if (address(params.billingToken) != address(i_LINK)) revert OnlyLink();
    params.amount = uint96(amount); // ignore whatever is sent in registration params, use actual value; casting safe because max supply LINK < 2^96
    _register(params, sender);
  }

  // ================================================================
  // |                           PRIVATE                            |
  // ================================================================

  /**
   * @dev verify registration request and emit RegistrationRequested event
   * @dev we don't allow multiple duplicate registrations by adding to the original registration's balance
   * users can cancel and re-register if they want to update the registration
   */
  function _register(RegistrationParams memory params, address sender) private returns (uint256) {
    if (params.amount < s_minRegistrationAmounts[params.billingToken]) {
      revert InsufficientPayment();
    }
    if (params.adminAddress == address(0)) {
      revert InvalidAdminAddress();
    }
    if (!s_registry.supportsBillingToken(address(params.billingToken))) {
      revert InvalidBillingToken();
    }
    bytes32 hash = keccak256(abi.encode(params));

    if (s_pendingRequests[hash].admin != address(0)) {
      revert DuplicateEntry();
    }

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
      s_pendingRequests[hash] = PendingRequest({
        admin: params.adminAddress,
        balance: params.amount,
        billingToken: params.billingToken
      });
    }

    return upkeepId;
  }

  /**
   * @dev register upkeep on AutomationRegistry contract and emit RegistrationApproved event
   * @dev safeApprove is deprecated and removed from the latest (v5) OZ release, Use safeIncreaseAllowance when we upgrade OZ (we are on v4.8)
   * @dev we stick to the safeApprove because of the older version (v4.8) of safeIncreaseAllowance can't handle USDT correctly, but newer version can
   */
  function _approve(RegistrationParams memory params, bytes32 hash) private returns (uint256) {
    IAutomationRegistryMaster2_3 registry = s_registry;
    uint256 upkeepId = registry.registerUpkeep(
      params.upkeepContract,
      params.gasLimit,
      params.adminAddress,
      params.triggerType,
      address(params.billingToken), // have to cast as address because master interface doesn't use contract types
      params.checkData,
      params.triggerConfig,
      params.offchainConfig
    );

    if (address(params.billingToken) == address(i_LINK)) {
      bool success = i_LINK.transferAndCall(address(registry), params.amount, abi.encode(upkeepId));
      if (!success) {
        revert TransferFailed(address(registry));
      }
    } else {
      params.billingToken.safeApprove(address(registry), params.amount);
      registry.addFunds(upkeepId, params.amount);
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
}
