// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {Chainable} from "../Chainable.sol";
import {IERC20Metadata as IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {SafeCast} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";

contract AutomationRegistryLogicB2_3 is AutomationRegistryBase2_3, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;
  using SafeERC20 for IERC20;

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
  // |                      PIPELINE FUNCTIONS                      |
  // ================================================================

  /**
   * @notice called by the automation DON to check if work is needed
   * @param id the upkeep ID to check for work needed
   * @param triggerData extra contextual data about the trigger (not used in all code paths)
   * @dev this one of the core functions called in the hot path
   * @dev there is a 2nd checkUpkeep function (below) that is being maintained for backwards compatibility
   * @dev there is an incongruency on what gets returned during failure modes
   * ex sometimes we include price data, sometimes we omit it depending on the failure
   */
  function checkUpkeep(
    uint256 id,
    bytes memory triggerData
  )
    public
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkUSD
    )
  {
    _preventExecution();

    Trigger triggerType = _getTriggerType(id);
    HotVars memory hotVars = s_hotVars;
    Upkeep memory upkeep = s_upkeep[id];

    {
      uint256 nativeUSD;
      uint96 maxPayment;
      if (hotVars.paused) return (false, bytes(""), UpkeepFailureReason.REGISTRY_PAUSED, 0, upkeep.performGas, 0, 0);
      if (upkeep.maxValidBlocknumber != UINT32_MAX)
        return (false, bytes(""), UpkeepFailureReason.UPKEEP_CANCELLED, 0, upkeep.performGas, 0, 0);
      if (upkeep.paused) return (false, bytes(""), UpkeepFailureReason.UPKEEP_PAUSED, 0, upkeep.performGas, 0, 0);
      (fastGasWei, linkUSD, nativeUSD) = _getFeedData(hotVars);
      maxPayment = _getMaxPayment(
        id,
        hotVars,
        triggerType,
        upkeep.performGas,
        fastGasWei,
        linkUSD,
        nativeUSD,
        upkeep.billingToken
      );
      if (upkeep.balance < maxPayment) {
        return (false, bytes(""), UpkeepFailureReason.INSUFFICIENT_BALANCE, 0, upkeep.performGas, 0, 0);
      }
    }

    bytes memory callData = _checkPayload(id, triggerType, triggerData);

    gasUsed = gasleft();
    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory result) = upkeep.forwarder.getTarget().call{gas: s_storage.checkGasLimit}(callData);
    gasUsed = gasUsed - gasleft();

    if (!success) {
      // User's target check reverted. We capture the revert data here and pass it within performData
      if (result.length > s_storage.maxRevertDataSize) {
        return (
          false,
          bytes(""),
          UpkeepFailureReason.REVERT_DATA_EXCEEDS_LIMIT,
          gasUsed,
          upkeep.performGas,
          fastGasWei,
          linkUSD
        );
      }
      return (
        upkeepNeeded,
        result,
        UpkeepFailureReason.TARGET_CHECK_REVERTED,
        gasUsed,
        upkeep.performGas,
        fastGasWei,
        linkUSD
      );
    }

    (upkeepNeeded, performData) = abi.decode(result, (bool, bytes));
    if (!upkeepNeeded)
      return (false, bytes(""), UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed, upkeep.performGas, fastGasWei, linkUSD);

    if (performData.length > s_storage.maxPerformDataSize)
      return (
        false,
        bytes(""),
        UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
        gasUsed,
        upkeep.performGas,
        fastGasWei,
        linkUSD
      );

    return (upkeepNeeded, performData, upkeepFailureReason, gasUsed, upkeep.performGas, fastGasWei, linkUSD);
  }

  /**
   * @notice see other checkUpkeep function for description
   * @dev this function may be deprecated in a future version of chainlink automation
   */
  function checkUpkeep(
    uint256 id
  )
    external
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkUSD
    )
  {
    return checkUpkeep(id, bytes(""));
  }

  /**
   * @dev checkCallback is used specifically for automation data streams lookups (see StreamsLookupCompatibleInterface.sol)
   * @param id the upkeepID to execute a callback for
   * @param values the values returned from the data streams lookup
   * @param extraData the user-provided extra context data
   */
  function checkCallback(
    uint256 id,
    bytes[] memory values,
    bytes calldata extraData
  )
    external
    returns (bool upkeepNeeded, bytes memory performData, UpkeepFailureReason upkeepFailureReason, uint256 gasUsed)
  {
    bytes memory payload = abi.encodeWithSelector(CHECK_CALLBACK_SELECTOR, values, extraData);
    return executeCallback(id, payload);
  }

  /**
   * @notice this is a generic callback executor that forwards a call to a user's contract with the configured
   * gas limit
   * @param id the upkeepID to execute a callback for
   * @param payload the data (including function selector) to call on the upkeep target contract
   */
  function executeCallback(
    uint256 id,
    bytes memory payload
  )
    public
    returns (bool upkeepNeeded, bytes memory performData, UpkeepFailureReason upkeepFailureReason, uint256 gasUsed)
  {
    _preventExecution();

    Upkeep memory upkeep = s_upkeep[id];
    gasUsed = gasleft();
    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory result) = upkeep.forwarder.getTarget().call{gas: s_storage.checkGasLimit}(payload);
    gasUsed = gasUsed - gasleft();
    if (!success) {
      return (false, bytes(""), UpkeepFailureReason.CALLBACK_REVERTED, gasUsed);
    }
    (upkeepNeeded, performData) = abi.decode(result, (bool, bytes));
    if (!upkeepNeeded) {
      return (false, bytes(""), UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed);
    }
    if (performData.length > s_storage.maxPerformDataSize) {
      return (false, bytes(""), UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT, gasUsed);
    }
    return (upkeepNeeded, performData, upkeepFailureReason, gasUsed);
  }

  /**
   * @notice simulates the upkeep with the perform data returned from checkUpkeep
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   * @return success whether the call reverted or not
   * @return gasUsed the amount of gas the target contract consumed
   */
  function simulatePerformUpkeep(
    uint256 id,
    bytes calldata performData
  ) external returns (bool success, uint256 gasUsed) {
    _preventExecution();

    if (s_hotVars.paused) revert RegistryPaused();
    Upkeep memory upkeep = s_upkeep[id];
    (success, gasUsed) = _performUpkeep(upkeep.forwarder, upkeep.performGas, performData);
    return (success, gasUsed);
  }

  // ================================================================
  // |                      UPKEEP MANAGEMENT                       |
  // ================================================================

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
        revert InvalidToken();
      }
      amount = SafeCast.toUint96(msg.value);
    }

    s_upkeep[id].balance = upkeep.balance + amount;
    s_reserveAmounts[upkeep.billingToken] = s_reserveAmounts[upkeep.billingToken] + amount;

    if (msg.value == 0) {
      // ERC20 payment
      upkeep.billingToken.safeTransferFrom(msg.sender, address(this), amount);
    } else {
      // native payment
      i_wrappedNativeToken.deposit{value: amount}();
    }

    emit FundsAdded(id, msg.sender, amount);
  }

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
    emit BillingConfigOverridden(id, billingOverrides);
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
    upkeep.billingToken.safeTransfer(to, amountToWithdraw);
    emit FundsWithdrawn(id, amountToWithdraw, to);
  }

  // ================================================================
  // |                       FINANCE ACTIONS                        |
  // ================================================================

  /**
   * @notice withdraws excess LINK from the liquidity pool
   * @param to the address to send the fees to
   * @param amount the amount to withdraw
   */
  function withdrawLink(address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();

    int256 available = _linkAvailableForPayment();
    if (available < 0) {
      revert InsufficientBalance(0, amount);
    } else if (amount > uint256(available)) {
      revert InsufficientBalance(uint256(available), amount);
    }

    bool transferStatus = i_link.transfer(to, amount);
    if (!transferStatus) {
      revert TransferFailed();
    }
    emit FeesWithdrawn(address(i_link), to, amount);
  }

  /**
   * @notice withdraws non-LINK fees earned by the contract
   * @param asset the asset to withdraw
   * @param to the address to send the fees to
   * @param amount the amount to withdraw
   * @dev in ON_CHAIN mode, we prevent withdrawing non-LINK fees unless there is sufficient LINK liquidity
   * to cover all outstanding debts on the registry
   */
  function withdrawERC20Fees(IERC20 asset, address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    if (address(asset) == address(i_link)) revert InvalidToken();
    if (_linkAvailableForPayment() < 0 && s_payoutMode == PayoutMode.ON_CHAIN) revert InsufficientLinkLiquidity();
    uint256 available = asset.balanceOf(address(this)) - s_reserveAmounts[asset];
    if (amount > available) revert InsufficientBalance(available, amount);

    asset.safeTransfer(to, amount);
    emit FeesWithdrawn(address(asset), to, amount);
  }
}
