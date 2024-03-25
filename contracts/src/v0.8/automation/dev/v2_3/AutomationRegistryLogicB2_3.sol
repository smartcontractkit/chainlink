// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {Chainable} from "../../Chainable.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeCast} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";

contract AutomationRegistryLogicB2_3 is AutomationRegistryBase2_3, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @param logicC the address of the third logic contract
   */
  constructor(
    AutomationRegistryLogicC2_3 logicC
  )
    AutomationRegistryBase2_3(
      logicC.getLinkAddress(),
      logicC.getLinkUSDFeedAddress(),
      logicC.getNativeUSDFeedAddress(),
      logicC.getFastGasFeedAddress(),
      logicC.getAutomationForwarderLogic(),
      logicC.getAllowedReadOnlyAddress(),
      logicC.getPayoutMode(),
      logicC.getWrappedNativeTokenAddress()
    )
    Chainable(address(logicC))
  {}

  // ================================================================
  // |                      UPKEEP MANAGEMENT                       |
  // ================================================================

  /**
   * @notice overrides the billing config for an upkeep
   * @param id the upkeepID
   * @param billingOverrides the override-able billing config
   */
  function setBillingOverrides(uint256 id, BillingOverrides calldata billingOverrides) external {
    _onlyPrivilegeManagerAllowed();
    if (s_upkeep[id].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();

    s_upkeep[id].overridesEnabled = true;
    s_billingOverrides[id] = billingOverrides;
    emit BillingConfigOverridden(id);
  }

  /**
   * @notice remove the overridden billing config for an upkeep
   * @param id the upkeepID
   */
  function removeBillingOverrides(uint256 id) external {
    _onlyPrivilegeManagerAllowed();

    s_upkeep[id].overridesEnabled = false;
    delete s_billingOverrides[id];
    emit BillingConfigOverrideRemoved(id);
  }

  /**
   * @notice adds fund to an upkeep
   * @param id the upkeepID
   * @param amount the amount of funds to add, in the upkeep's billing token
   */
  function addFunds(uint256 id, uint96 amount) external payable {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();

    if (msg.value != 0) {
      if (upkeep.billingToken != IERC20(i_wrappedNativeToken)) {
        revert InvalidBillingToken();
      }
      amount = SafeCast.toUint96(msg.value);
    }

    s_upkeep[id].balance = upkeep.balance + amount;
    s_reserveAmounts[upkeep.billingToken] = s_reserveAmounts[upkeep.billingToken] + amount;

    if (msg.value == 0) {
      // ERC20 payment
      bool success = upkeep.billingToken.transferFrom(msg.sender, address(this), amount);
      if (!success) revert TransferFailed();
    } else {
      // native payment
      i_wrappedNativeToken.deposit{value: amount}();
    }

    emit FundsAdded(id, msg.sender, amount);
  }

  /**
   * @notice transfers the address of an admin for an upkeep
   */
  function transferUpkeepAdmin(uint256 id, address proposed) external {
    _requireAdminAndNotCancelled(id);
    if (proposed == msg.sender) revert ValueNotChanged();

    if (s_proposedAdmin[id] != proposed) {
      s_proposedAdmin[id] = proposed;
      emit UpkeepAdminTransferRequested(id, msg.sender, proposed);
    }
  }

  /**
   * @notice accepts the transfer of an upkeep admin
   */
  function acceptUpkeepAdmin(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (s_proposedAdmin[id] != msg.sender) revert OnlyCallableByProposedAdmin();
    address past = s_upkeepAdmin[id];
    s_upkeepAdmin[id] = msg.sender;
    s_proposedAdmin[id] = ZERO_ADDRESS;

    emit UpkeepAdminTransferred(id, past, msg.sender);
  }

  /**
   * @notice pauses an upkeep - an upkeep will be neither checked nor performed while paused
   */
  function pauseUpkeep(uint256 id) external {
    _requireAdminAndNotCancelled(id);
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.paused) revert OnlyUnpausedUpkeep();
    s_upkeep[id].paused = true;
    s_upkeepIDs.remove(id);
    emit UpkeepPaused(id);
  }

  /**
   * @notice unpauses an upkeep
   */
  function unpauseUpkeep(uint256 id) external {
    _requireAdminAndNotCancelled(id);
    Upkeep memory upkeep = s_upkeep[id];
    if (!upkeep.paused) revert OnlyPausedUpkeep();
    s_upkeep[id].paused = false;
    s_upkeepIDs.add(id);
    emit UpkeepUnpaused(id);
  }

  /**
   * @notice updates the checkData for an upkeep
   */
  function setUpkeepCheckData(uint256 id, bytes calldata newCheckData) external {
    _requireAdminAndNotCancelled(id);
    if (newCheckData.length > s_storage.maxCheckDataSize) revert CheckDataExceedsLimit();
    s_checkData[id] = newCheckData;
    emit UpkeepCheckDataSet(id, newCheckData);
  }

  /**
   * @notice updates the gas limit for an upkeep
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external {
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_storage.maxPerformGas) revert GasLimitOutsideRange();
    _requireAdminAndNotCancelled(id);
    s_upkeep[id].performGas = gasLimit;

    emit UpkeepGasLimitSet(id, gasLimit);
  }

  /**
   * @notice updates the offchain config for an upkeep
   */
  function setUpkeepOffchainConfig(uint256 id, bytes calldata config) external {
    _requireAdminAndNotCancelled(id);
    s_upkeepOffchainConfig[id] = config;
    emit UpkeepOffchainConfigSet(id, config);
  }

  /**
   * @notice sets the upkeep trigger config
   * @param id the upkeepID to change the trigger for
   * @param triggerConfig the new trigger config
   */
  function setUpkeepTriggerConfig(uint256 id, bytes calldata triggerConfig) external {
    _requireAdminAndNotCancelled(id);
    s_upkeepTriggerConfig[id] = triggerConfig;
    emit UpkeepTriggerConfigSet(id, triggerConfig);
  }

  /**
   * @notice withdraws an upkeep's funds from an upkeep
   * @dev note that an upkeep must be cancelled first!!
   */
  function withdrawFunds(uint256 id, address to) external nonReentrant {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    Upkeep memory upkeep = s_upkeep[id];
    if (s_upkeepAdmin[id] != msg.sender) revert OnlyCallableByAdmin();
    if (upkeep.maxValidBlocknumber > s_hotVars.chainModule.blockNumber()) revert UpkeepNotCanceled();
    uint96 amountToWithdraw = s_upkeep[id].balance;
    s_reserveAmounts[upkeep.billingToken] = s_reserveAmounts[upkeep.billingToken] - amountToWithdraw;
    s_upkeep[id].balance = 0;
    bool success = upkeep.billingToken.transfer(to, amountToWithdraw);
    if (!success) revert TransferFailed();
    emit FundsWithdrawn(id, amountToWithdraw, to);
  }

  /**
   * @notice LINK available to withdraw by the finance team
   # @dev LINK max supply < 2^96, so casting to int256 is safe
   */
  function linkAvailableForPayment() public view returns (int256) {
    return int256(i_link.balanceOf(address(this))) - int256(s_reserveAmounts[IERC20(address(i_link))]);
  }

  function withdrawLinkFees(address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();

    int256 available = linkAvailableForPayment();
    if (available < 0 || amount > uint256(available)) revert InsufficientBalance(available, amount);

    bool transferStatus = i_link.transfer(to, amount);
    if (!transferStatus) {
      revert TransferFailed();
    }
    emit FeesWithdrawn(to, address(i_link), amount);
  }

  function withdrawERC20Fees(address assetAddress, address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();

    bool transferStatus = IERC20(assetAddress).transfer(to, amount);
    if (!transferStatus) {
      revert TransferFailed();
    }

    emit FeesWithdrawn(to, assetAddress, amount);
  }
}
