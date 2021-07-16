// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/ArbitrumInboxInterface.sol";
import "../interfaces/AggregatorValidatorInterface.sol";
import "./interfaces/FlagsInterface.sol";
import "../SimpleWriteAccessController.sol";

/**
 * @title ArbitrumValidator Contract
 * @notice Allows to raise and lower Flags on the Arbitrum network through its Layer 1 contracts
 * The internal AccessController controls the access of the validate method
 * Gas configuration is controlled by a configurable external SimpleWriteAccessController
 * Funds on the contract are managed by the owner
 */
contract ArbitrumValidator is SimpleWriteAccessController, AggregatorValidatorInterface {
  // Follows: https://eips.ethereum.org/EIPS/eip-1967
  address constant private s_arbitrumFlag = address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-offline")) - 1)));
  bytes constant private s_raiseFlagData = abi.encodeWithSelector(FlagsInterface.raiseFlag.selector, s_arbitrumFlag);
  bytes constant private s_lowerFlagData = abi.encodeWithSelector(FlagsInterface.lowerFlag.selector, s_arbitrumFlag);
  uint32 constant private s_L2GasLimit = 30000000;

  address private s_flagsAddress;
  ArbitrumInboxInterface private s_arbitrumInbox;
  SimpleWriteAccessController private s_gasConfigAccessController;

  struct GasConfiguration {
    uint256 maxSubmissionCost;
    uint32 maximumGasPrice;
    uint256 gasCostL2;
    address refundableAddress;
  }
  GasConfiguration private s_gasConfig;

  /**
   * @notice emitted when a new gas configuration is set
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maximumGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  event GasConfigurationSet(
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  );

  /**
   * @notice emitted when a new gas access-control contract is set
   * @param old the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event GasAccessControllerSet(
    address old, 
    address current
  );

  /**
   * @param inboxAddress address of the Arbitrum Inbox L1 contract
   * @param flagAddress address of the Chainlink L2 Flags contract
   * @param gasConfigAccessController address of the access controller for managing gas price on Arbitrum
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maximumGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   */
  constructor(
    address inboxAddress,
    address flagAddress,
    address gasConfigAccessController,
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  ) {
    require(flagAddress != address(0), "Invalid Flags contract address");
    s_arbitrumInbox = ArbitrumInboxInterface(inboxAddress);
    s_gasConfigAccessController = SimpleWriteAccessController(gasConfigAccessController);
    s_flagsAddress = flagAddress;
    setGasConfigurationInternal(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);
  }
  
  /**
   * @notice makes this contract payable. It need funds in order to pay for L2 transactions fees
   */
  fallback() external payable {}

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
   * @param to address where to send the funds
   * @dev only owner can call this
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
   * @notice sets gasAccessController
   * @param gasAccessController new gasAccessController contract address
   * @dev only owner can call this
   */
  function setGasAccessController(
    address gasAccessController
  )
    external
    onlyOwner
  {
    setGasAccessControllerInternal(gasAccessController);
  }
  
  /**
   * @notice sets Arbitrum gas configuration
   * @param maxSubmissionCost maximum cost willing to pay on L2
   * @param maximumGasPrice maximum gas price to pay on L2
   * @param gasCostL2 value to send to L2 to cover gas fee
   * @param refundableAddress address where gas excess on L2 will be sent
   * @dev access control provided by s_gasConfigAccessController
   */
  function setGasConfiguration(
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  )
    external
  {
    require(s_gasConfigAccessController.hasAccess(msg.sender, msg.data), "Only gas configuration admin can call");
    setGasConfigurationInternal(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);
  }

  
  /**
   * @notice validate method updates the state of an L2 Flag in case of change on the Arbitrum Sequencer. A zero answer considers the service as offline
   * In case the previous answer is the same as the current it does not trigger any tx on L2
   * In other case, a retryable ticket is created on the Arbitrum L1 Inbox contract. The tx gas fee can be paid from this contract providing a value, or the same address on L2
   * @param previousAnswer previous aggregator answer
   * @param currentAnswer new aggregator answer
   * @dev access control provided internally by SimpleWriteAccessController
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

    s_arbitrumInbox.createRetryableTicket{value: s_gasConfig.gasCostL2}(
      s_flagsAddress, 
      0, // L2 call value
      s_gasConfig.maxSubmissionCost, // Max submission cost of sending data length
      s_gasConfig.refundableAddress, // excessFeeRefundAddress
      s_gasConfig.refundableAddress, // callValueRefundAddress
      s_L2GasLimit,
      s_gasConfig.maximumGasPrice, 
      currentAnswer == 0 ? s_raiseFlagData : s_lowerFlagData
    );
    return true;
  }
  
  function setGasConfigurationInternal(
    uint256 maxSubmissionCost,
    uint32 maximumGasPrice,
    uint256 gasCostL2,
    address refundableAddress
  ) internal {
    s_gasConfig = GasConfiguration(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);
    emit GasConfigurationSet(maxSubmissionCost, maximumGasPrice, gasCostL2, refundableAddress);
  }

  function setGasAccessControllerInternal(
    address gasAccessController
  )
    internal
  {
    address oldController = address(s_gasConfigAccessController);
    if (gasAccessController != oldController) {
      s_gasConfigAccessController = SimpleWriteAccessController(gasAccessController);
      emit GasAccessControllerSet(
        oldController,
        gasAccessController
      );
    }
  }
}
