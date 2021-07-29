// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/AggregatorValidatorInterface.sol";
import "../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/AccessControllerInterface.sol";
import "../SimpleWriteAccessController.sol";

/* dev dependencies - to be re/moved after audit */
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/bridge/interfaces/IInbox.sol";
import "./interfaces/FlagsInterface.sol";

/**
 * @title ArbitrumValidator
 * @notice Allows to raise and lower Flags on the Arbitrum network through its Layer 1 contracts
 *  - The internal AccessController controls the access of the validate method
 *  - Gas configuration is controlled by a configurable external SimpleWriteAccessController
 *  - Funds on the contract are managed by the owner
 */
contract ArbitrumValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  // Config for L1 -> L2 `createRetryableTicket` call
  struct GasConfiguration {
    uint256 maxSubmissionCost;
    uint256 maxGasPrice;
    uint256 gasCostL2;
    uint256 gasLimitL2;
    address refundableAddress;
  }

  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address constant private FLAG_ARBITRUM_SEQ_OFFLINE = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));
  bytes constant private CALL_RAISE_FLAG = abi.encodeWithSelector(FlagsInterface.raiseFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);
  bytes constant private CALL_LOWER_FLAG = abi.encodeWithSelector(FlagsInterface.lowerFlag.selector, FLAG_ARBITRUM_SEQ_OFFLINE);

  address private s_l2FlagsAddress;
  IInbox private s_inbox;
  AccessControllerInterface private s_gasConfigAccessController;
  GasConfiguration private s_gasConfig;

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  event GasConfigurationSet(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address indexed refundableAddress
  );

  /**
   * @notice emitted when a new gas access-control contract is set
   * @param previous the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event GasAccessControllerSet(
    address indexed previous,
    address indexed current
  );

  /**
   * @param inboxAddress address of the Arbitrum Inbox L1 contract
   * @param l2FlagsAddress address of the Chainlink L2 Flags contract
   * @param gasConfigAccessControllerAddress address of the access controller for managing gas price on Arbitrum
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param gasLimitL2 gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  constructor(
    address inboxAddress,
    address l2FlagsAddress,
    address gasConfigAccessControllerAddress,
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddress
  ) {
    require(inboxAddress != address(0), "Invalid Inbox contract address");
    require(l2FlagsAddress != address(0), "Invalid Flags contract address");
    s_inbox = IInbox(inboxAddress);
    s_gasConfigAccessController = AccessControllerInterface(gasConfigAccessControllerAddress);
    s_l2FlagsAddress = l2FlagsAddress;
    _setGasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddress);
  }

  /**
   * @notice versions:
   *
   * - ArbitrumValidator 0.1.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion()
    external
    pure
    virtual
    override
    returns (
      string memory
    )
  {
    return "ArbitrumValidator 0.1.0";
  }

  /// @return L2 Flags contract address
  function l2Flags()
    external
    view
    virtual
    returns (address)
  {
    return s_l2FlagsAddress;
  }

  /// @return Arbitrum Inbox contract address
  function inbox()
    external
    view
    virtual
    returns (address)
  {
    return address(s_inbox);
  }

  /// @return gas config AccessControllerInterface contract address
  function gasConfigAccessController()
    external
    view
    virtual
    returns (address)
  {
    return address(s_gasConfigAccessController);
  }

  /// @return stored GasConfiguration
  function gasConfig()
    external
    view
    virtual
    returns (GasConfiguration memory)
  {
    return s_gasConfig;
  }

  /// @notice makes this contract payable as it need funds to pay for L2 transactions fees on L1.
  receive() external payable {}

  /**
   * @notice withdraws all funds availbale in this contract to the msg.sender
   * @dev only owner can call this
   */
  function withdrawFunds()
    external
    onlyOwner()
  {
    address payable to = payable(msg.sender);
    to.transfer(address(this).balance);
  }

  /**
   * @notice withdraws all funds availbale in this contract to the address specified
   * @dev only owner can call this
   * @param to address where to send the funds
   */
  function withdrawFundsTo(
    address payable to
  ) 
    external
    onlyOwner()
  {
    to.transfer(address(this).balance);
  }

  /**
   * @notice sets gas config AccessControllerInterface contract
   * @dev only owner can call this
   * @param accessController new AccessControllerInterface contract address
   */
  function setGasAccessController(
    address accessController
  )
    external
    onlyOwner
  {
    _setGasAccessController(accessController);
  }

  /**
   * @notice sets Arbitrum gas configuration
   * @dev access control provided by s_gasConfigAccessController
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maxGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param gasLimitL2 gas limit for immediate L2 execution attempt. A value around 1M should be sufficient
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  function setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddress
  )
    external
  {
    require(s_gasConfigAccessController.hasAccess(msg.sender, msg.data), "Access required to set config");
    _setGasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddress);
  }

  /**
   * @notice validate method updates the state of an L2 Flag in case of change on the Arbitrum Sequencer.
   * A one answer considers the service as offline.
   * In case the previous answer is the same as the current it does not trigger any tx on L2. In other case,
   * a retryable ticket is created on the Arbitrum L1 Inbox contract. The tx gas fee can be paid from this
   * contract providing a value, or the same address on L2.
   * @dev access control provided internally by SimpleWriteAccessController
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer
   */
  function validate(
    uint256 /* previousRoundId */,
    int256 previousAnswer,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) 
    external
    override
    checkAccess()
    returns (bool)
  {
    // Avoids resending to L2 the same tx on every call
    if (previousAnswer == currentAnswer) {
      return true;
    }

    int isServiceOffline = 1;
    // NOTICE: if gasCostL2 is zero the payment is processed on L2 so the L2 address needs to be funded, as it will
    // paying the fee. We also ignore the returned msg number, that can be queried via the InboxMessageDelivered event.
    s_inbox.createRetryableTicket{value: s_gasConfig.gasCostL2}(
      s_l2FlagsAddress,
      0, // L2 call value
      // NOTICE: maxSubmissionCost info will possibly become available on L1 after the London fork. At that time this
      // contract could start querying/calculating it directly so we wouldn't need to configure it statically. On L2 this
      // info is available via `ArbRetryableTx.getSubmissionPrice`.
      s_gasConfig.maxSubmissionCost, // Max submission cost of sending data length
      s_gasConfig.refundableAddress, // excessFeeRefundAddress
      s_gasConfig.refundableAddress, // callValueRefundAddress
      s_gasConfig.gasLimitL2,
      s_gasConfig.maxGasPrice,
      currentAnswer == isServiceOffline ? CALL_RAISE_FLAG : CALL_LOWER_FLAG
    );
    return true;
  }

  function _setGasConfiguration(
    uint256 maxSubmissionCost,
    uint256 maxGasPrice,
    uint256 gasCostL2,
    uint256 gasLimitL2,
    address refundableAddress
  )
    internal
  {
    // L2 will pay the fee if gasCostL2 is zero
    if (gasCostL2 > 0) {
      uint256 minGasCostValue = maxSubmissionCost + gasLimitL2 * maxGasPrice;
      require(gasCostL2 >= minGasCostValue, "Gas cost provided is too low");
    }
    s_gasConfig = GasConfiguration(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddress);
    emit GasConfigurationSet(maxSubmissionCost, maxGasPrice, gasCostL2, gasLimitL2, refundableAddress);
  }

  function _setGasAccessController(
    address accessController
  )
    internal
  {
    address previousAccessController = address(s_gasConfigAccessController);
    if (accessController != previousAccessController) {
      s_gasConfigAccessController = AccessControllerInterface(accessController);
      emit GasAccessControllerSet(previousAccessController, accessController);
    }
  }
}
