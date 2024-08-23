// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {AutomationRegistryLogicB2_3} from "./AutomationRegistryLogicB2_3.sol";
import {Chainable} from "../Chainable.sol";
import {AutomationForwarder} from "../AutomationForwarder.sol";
import {IAutomationForwarder} from "../interfaces/IAutomationForwarder.sol";
import {UpkeepTranscoderInterfaceV2} from "../interfaces/UpkeepTranscoderInterfaceV2.sol";
import {MigratableKeeperRegistryInterfaceV2} from "../interfaces/MigratableKeeperRegistryInterfaceV2.sol";
import {IERC20Metadata as IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC677Receiver} from "../../shared/interfaces/IERC677Receiver.sol";

/**
 * @notice Logic contract, works in tandem with AutomationRegistry as a proxy
 */
contract AutomationRegistryLogicA2_3 is AutomationRegistryBase2_3, Chainable, IERC677Receiver {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;
  using SafeERC20 for IERC20;

  /**
   * @param logicB the address of the second logic contract
   * @dev we cast the contract to logicC in order to call logicC functions (via fallback)
   */
  constructor(
    AutomationRegistryLogicB2_3 logicB
  )
    AutomationRegistryBase2_3(
      AutomationRegistryLogicC2_3(address(logicB)).getLinkAddress(),
      AutomationRegistryLogicC2_3(address(logicB)).getLinkUSDFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicB)).getNativeUSDFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicB)).getFastGasFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicB)).getAutomationForwarderLogic(),
      AutomationRegistryLogicC2_3(address(logicB)).getAllowedReadOnlyAddress(),
      AutomationRegistryLogicC2_3(address(logicB)).getPayoutMode(),
      AutomationRegistryLogicC2_3(address(logicB)).getWrappedNativeTokenAddress()
    )
    Chainable(address(logicB))
  {}

  /**
   * @notice uses LINK's transferAndCall to LINK and add funding to an upkeep
   * @dev safe to cast uint256 to uint96 as total LINK supply is under UINT96MAX
   * @param sender the account which transferred the funds
   * @param amount number of LINK transfer
   */
  function onTokenTransfer(address sender, uint256 amount, bytes calldata data) external override {
    if (msg.sender != address(i_link)) revert OnlyCallableByLINKToken();
    if (data.length != 32) revert InvalidDataLength();
    uint256 id = abi.decode(data, (uint256));
    if (s_upkeep[id].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (address(s_upkeep[id].billingToken) != address(i_link)) revert InvalidToken();
    s_upkeep[id].balance = s_upkeep[id].balance + uint96(amount);
    s_reserveAmounts[IERC20(address(i_link))] = s_reserveAmounts[IERC20(address(i_link))] + amount;
    emit FundsAdded(id, sender, uint96(amount));
  }

  // ================================================================
  // |                      UPKEEP MANAGEMENT                       |
  // ================================================================

  /**
   * @notice adds a new upkeep
   * @param target address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when
   * performing upkeep
   * @param admin address to cancel upkeep and withdraw remaining funds
   * @param triggerType the trigger for the upkeep
   * @param billingToken the billing token for the upkeep
   * @param checkData data passed to the contract when checking for upkeep
   * @param triggerConfig the config for the trigger
   * @param offchainConfig arbitrary offchain config for the upkeep
   */
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    Trigger triggerType,
    IERC20 billingToken,
    bytes calldata checkData,
    bytes memory triggerConfig,
    bytes memory offchainConfig
  ) public returns (uint256 id) {
    if (msg.sender != owner() && !s_registrars.contains(msg.sender)) revert OnlyCallableByOwnerOrRegistrar();
    if (!target.isContract()) revert NotAContract();
    id = _createID(triggerType);
    IAutomationForwarder forwarder = IAutomationForwarder(
      address(new AutomationForwarder(target, address(this), i_automationForwarderLogic))
    );
    _createUpkeep(
      id,
      Upkeep({
        overridesEnabled: false,
        performGas: gasLimit,
        balance: 0,
        maxValidBlocknumber: UINT32_MAX,
        lastPerformedBlockNumber: 0,
        amountSpent: 0,
        paused: false,
        forwarder: forwarder,
        billingToken: billingToken
      }),
      admin,
      checkData,
      triggerConfig,
      offchainConfig
    );
    s_storage.nonce++;
    emit UpkeepRegistered(id, gasLimit, admin);
    emit UpkeepCheckDataSet(id, checkData);
    emit UpkeepTriggerConfigSet(id, triggerConfig);
    emit UpkeepOffchainConfigSet(id, offchainConfig);
    return (id);
  }

  /**
   * @notice cancels an upkeep
   * @param id the upkeepID to cancel
   * @dev if a user cancels an upkeep, their funds are locked for CANCELLATION_DELAY blocks to
   * allow any pending performUpkeep txs time to get confirmed
   */
  function cancelUpkeep(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    bool isOwner = msg.sender == owner();
    uint96 minSpend = s_billingConfigs[upkeep.billingToken].minSpend;

    uint256 height = s_hotVars.chainModule.blockNumber();
    if (upkeep.maxValidBlocknumber == 0) revert CannotCancel();
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (!isOwner && msg.sender != s_upkeepAdmin[id]) revert OnlyCallableByOwnerOrAdmin();

    if (!isOwner) {
      height = height + CANCELLATION_DELAY;
    }
    s_upkeep[id].maxValidBlocknumber = uint32(height);
    s_upkeepIDs.remove(id);

    // charge the cancellation fee if the minSpend is not met
    uint96 cancellationFee = 0;
    // cancellationFee is min(max(minSpend - amountSpent, 0), amountLeft)
    if (upkeep.amountSpent < minSpend) {
      cancellationFee = minSpend - uint96(upkeep.amountSpent);
      if (cancellationFee > upkeep.balance) {
        cancellationFee = upkeep.balance;
      }
    }
    s_upkeep[id].balance = upkeep.balance - cancellationFee;
    s_reserveAmounts[upkeep.billingToken] = s_reserveAmounts[upkeep.billingToken] - cancellationFee;

    emit UpkeepCanceled(id, uint64(height));
  }

  /**
   * @notice migrates upkeeps from one registry to another.
   * @param ids the upkeepIDs to migrate
   * @param destination the destination registry address
   * @dev a transcoder must be set in order to enable migration
   * @dev migration permissions must be set on *both* sending and receiving registries
   * @dev only an upkeep admin can migrate their upkeeps
   * @dev this function is most gas-efficient if upkeepIDs are sorted by billing token
   * @dev s_billingOverrides and s_upkeepPrivilegeConfig are not migrated in this function
   */
  function migrateUpkeeps(uint256[] calldata ids, address destination) external {
    if (
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.OUTGOING &&
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    if (s_storage.transcoder == ZERO_ADDRESS) revert TranscoderNotSet();
    if (ids.length == 0) revert ArrayHasNoEntries();

    IERC20 billingToken;
    uint256 balanceToTransfer;
    uint256 id;
    Upkeep memory upkeep;
    address[] memory admins = new address[](ids.length);
    Upkeep[] memory upkeeps = new Upkeep[](ids.length);
    bytes[] memory checkDatas = new bytes[](ids.length);
    bytes[] memory triggerConfigs = new bytes[](ids.length);
    bytes[] memory offchainConfigs = new bytes[](ids.length);

    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];

      if (idx == 0) {
        billingToken = upkeep.billingToken;
        balanceToTransfer = upkeep.balance;
      }

      // if we encounter a new billing token, send the sum from the last billing token to the destination registry
      if (upkeep.billingToken != billingToken) {
        s_reserveAmounts[billingToken] = s_reserveAmounts[billingToken] - balanceToTransfer;
        billingToken.safeTransfer(destination, balanceToTransfer);
        billingToken = upkeep.billingToken;
        balanceToTransfer = upkeep.balance;
      } else if (idx != 0) {
        balanceToTransfer += upkeep.balance;
      }

      _requireAdminAndNotCancelled(id);
      upkeep.forwarder.updateRegistry(destination);

      upkeeps[idx] = upkeep;
      admins[idx] = s_upkeepAdmin[id];
      checkDatas[idx] = s_checkData[id];
      triggerConfigs[idx] = s_upkeepTriggerConfig[id];
      offchainConfigs[idx] = s_upkeepOffchainConfig[id];
      delete s_upkeep[id];
      delete s_checkData[id];
      delete s_upkeepTriggerConfig[id];
      delete s_upkeepOffchainConfig[id];
      // nullify existing proposed admin change if an upkeep is being migrated
      delete s_proposedAdmin[id];
      delete s_upkeepAdmin[id];
      s_upkeepIDs.remove(id);
      emit UpkeepMigrated(id, upkeep.balance, destination);
    }
    // always transfer the rolling sum in the end
    s_reserveAmounts[billingToken] = s_reserveAmounts[billingToken] - balanceToTransfer;
    billingToken.safeTransfer(destination, balanceToTransfer);

    bytes memory encodedUpkeeps = abi.encode(
      ids,
      upkeeps,
      new address[](ids.length),
      admins,
      checkDatas,
      triggerConfigs,
      offchainConfigs
    );
    MigratableKeeperRegistryInterfaceV2(destination).receiveUpkeeps(
      UpkeepTranscoderInterfaceV2(s_storage.transcoder).transcodeUpkeeps(
        UPKEEP_VERSION_BASE,
        MigratableKeeperRegistryInterfaceV2(destination).upkeepVersion(),
        encodedUpkeeps
      )
    );
  }

  /**
   * @notice received upkeeps migrated from another registry
   * @param encodedUpkeeps the raw upkeep data to import
   * @dev this function is never called directly, it is only called by another registry's migrate function
   * @dev s_billingOverrides and s_upkeepPrivilegeConfig are not handled in this function
   */
  function receiveUpkeeps(bytes calldata encodedUpkeeps) external {
    if (
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.INCOMING &&
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    (
      uint256[] memory ids,
      Upkeep[] memory upkeeps,
      address[] memory targets,
      address[] memory upkeepAdmins,
      bytes[] memory checkDatas,
      bytes[] memory triggerConfigs,
      bytes[] memory offchainConfigs
    ) = abi.decode(encodedUpkeeps, (uint256[], Upkeep[], address[], address[], bytes[], bytes[], bytes[]));
    for (uint256 idx = 0; idx < ids.length; idx++) {
      if (address(upkeeps[idx].forwarder) == ZERO_ADDRESS) {
        upkeeps[idx].forwarder = IAutomationForwarder(
          address(new AutomationForwarder(targets[idx], address(this), i_automationForwarderLogic))
        );
      }
      _createUpkeep(
        ids[idx],
        upkeeps[idx],
        upkeepAdmins[idx],
        checkDatas[idx],
        triggerConfigs[idx],
        offchainConfigs[idx]
      );
      emit UpkeepReceived(ids[idx], upkeeps[idx].balance, msg.sender);
    }
  }
}
