// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AggregatorValidatorInterface} from "../../../shared/interfaces/AggregatorValidatorInterface.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {ISequencerUptimeFeed} from "./../interfaces/ISequencerUptimeFeed.sol";

import {SimpleWriteAccessController} from "../../../shared/access/SimpleWriteAccessController.sol";

import {IBridgehub, L2TransactionRequestDirect} from "@zksync/contracts/l1-contracts/contracts/bridgehub/IBridgehub.sol";

/// @title ZKSyncValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
contract ZKSyncValidator is ITypeAndVersion, AggregatorValidatorInterface, SimpleWriteAccessController {
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L2_UPTIME_FEED_ADDR;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  uint32 private immutable CHAIN_ID;

  /// Contract state variables
  string public constant override typeAndVersion = "ZKSyncValidator 1.0.0";
  uint32 private constant TEST_NET_CHAIN_ID = 300;
  uint32 private constant MAIN_NET_CHAIN_ID = 324;
  int256 private constant ANSWER_SEQ_OFFLINE = 1;
  uint32 private s_l2GasPerPubdataByteLimit;
  uint32 private s_gasLimit;

  /// @notice emitted when gas cost to spend on L2 is updated
  /// @param gasLimit updated gas cost
  event GasLimitUpdated(uint32 gasLimit);

  /// @notice emitted when the gas per pubdata byte limit is updated
  /// @param l2GasPerPubdataByteLimit updated gas per pubdata byte limit
  event GasPerPubdataByteLimitUpdated(uint32 l2GasPerPubdataByteLimit);

  /// @notice ChainID is not a valid value
  error InvalidChainID();
  /// @notice Address cannot be zero
  error ZeroAddressNotAllowed(string message);

  /// @param l1CrossDomainMessengerAddress address the Bridgehub contract address
  /// @param l2UptimeFeedAddr the address of the SequencerUptimeFeedInterface contract address
  /// @param gasLimit the gasLimit to use for sending a message from L1 to L2
  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit,
    uint32 chainId,
    uint32 l2GasPerPubdataByteLimit
  ) {
    if (chainId != TEST_NET_CHAIN_ID && chainId != MAIN_NET_CHAIN_ID) {
      revert InvalidChainID();
    }

    if (l1CrossDomainMessengerAddress == address(0)) {
      revert ZeroAddressNotAllowed("Invalid xDomain Messenger address");
    }

    if (l2UptimeFeedAddr == address(0)) {
      revert ZeroAddressNotAllowed("Invalid SequencerUptimeFeedInterface contract address");
    }

    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
    s_l2GasPerPubdataByteLimit = l2GasPerPubdataByteLimit;
    s_gasLimit = gasLimit;
    CHAIN_ID = chainId;
  }

  /// @notice sets the new gas cost to spend when sending cross chain message
  /// @param gasLimit the updated gas cost
  function setGasLimit(uint32 gasLimit) external onlyOwner {
    if (s_gasLimit != gasLimit) {
      s_gasLimit = gasLimit;
      emit GasLimitUpdated(gasLimit);
    }
  }

  /// @notice fetches the gas cost of sending a cross chain message
  function getGasLimit() external view returns (uint32) {
    return s_gasLimit;
  }

  /// @notice sets the l2GasPerPubdataByteLimit for the L2 transaction request
  /// @param l2GasPerPubdataByteLimit the updated l2GasPerPubdataByteLimit
  function setL2GasPerPubdataByteLimit(uint32 l2GasPerPubdataByteLimit) external onlyOwner {
    if (s_l2GasPerPubdataByteLimit != l2GasPerPubdataByteLimit) {
      s_l2GasPerPubdataByteLimit = l2GasPerPubdataByteLimit;
      emit GasPerPubdataByteLimitUpdated(l2GasPerPubdataByteLimit);
    }
  }

  /// @notice fetches the l2GasPerPubdataByteLimit for the L2 transaction request
  function getL2GasPerPubdataByteLimit() external view returns (uint32) {
    return s_l2GasPerPubdataByteLimit;
  }

  /// @notice fetches the chain id
  function getChainId() external view returns (uint32) {
    return CHAIN_ID;
  }

  /// @notice makes this contract payable
  /// @dev receives funds:
  ///  - to use them (if configured) to pay for L2 execution on L1
  ///  - when withdrawing funds from L2 xDomain alias address (pay for L2 execution on L2)
  receive() external payable {}

  /// @notice validate method sends an xDomain L2 tx to update Uptime Feed contract on L2.
  /// @dev A message is sent using the Bridgehub. This method is accessed controlled.
  /// @param currentAnswer new aggregator answer - value of 1 considers the sequencer offline.
  function validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Encode the SequencerUptimeFeedInterface call with `status` and `timestamp`
    bytes memory message = abi.encodeWithSelector(
      ISequencerUptimeFeed.updateStatus.selector,
      currentAnswer == ANSWER_SEQ_OFFLINE,
      uint64(block.timestamp)
    );

    // Empty bytes for factoryDeps
    bytes[] memory emptyBytes;

    // Get a reference to the bridgehub contract
    IBridgehub bridgeHub = IBridgehub(L1_CROSS_DOMAIN_MESSENGER_ADDRESS);

    // Get the L2 transaction base cost
    uint256 transactionBaseCostEstimate = bridgeHub.l2TransactionBaseCost(
      CHAIN_ID,
      tx.gasprice,
      s_gasLimit,
      s_l2GasPerPubdataByteLimit
    );

    // Create the L2 transaction request
    L2TransactionRequestDirect memory l2TransactionRequestDirect = L2TransactionRequestDirect({
      chainId: CHAIN_ID,
      mintValue: transactionBaseCostEstimate,
      l2Contract: L2_UPTIME_FEED_ADDR,
      l2Value: 0,
      l2Calldata: message,
      l2GasLimit: s_gasLimit,
      l2GasPerPubdataByteLimit: s_l2GasPerPubdataByteLimit,
      factoryDeps: emptyBytes,
      refundRecipient: msg.sender
    });

    // Make the xDomain call
    bridgeHub.requestL2TransactionDirect{value: transactionBaseCostEstimate}(l2TransactionRequestDirect);

    // Return a success flag
    return true;
  }
}
