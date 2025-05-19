// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";

/// @dev test contract to test CCIPReader functionality, never deployed to real chains.
contract CCIPReaderTester {
  mapping(uint64 sourceChainSelector => OffRamp.SourceChainConfig sourceChainConfig) internal s_sourceChainConfigs;
  mapping(uint64 destChainSelector => uint64 sequenceNumber) internal s_destChainSeqNrs;
  mapping(uint64 sourceChainSelector => mapping(bytes sender => uint64 nonce)) internal s_senderNonce;
  uint64 private s_latestPriceSequenceNumber;

  /// @notice Gets the next sequence number to be used in the onRamp
  /// @param destChainSelector The destination chain selector
  /// @return nextSequenceNumber The next sequence number to be used
  function getExpectedNextSequenceNumber(
    uint64 destChainSelector
  ) external view returns (uint64) {
    return s_destChainSeqNrs[destChainSelector] + 1;
  }

  /// @notice Sets the sequence number in the onRamp
  /// @param destChainSelector The destination chain selector
  /// @param sequenceNumber The sequence number
  function setDestChainSeqNr(uint64 destChainSelector, uint64 sequenceNumber) external {
    s_destChainSeqNrs[destChainSelector] = sequenceNumber;
  }

  /// @notice Returns the inbound nonce for a given sender on a given source chain.
  /// @param sourceChainSelector The source chain selector.
  /// @param sender The encoded sender address.
  /// @return inboundNonce The inbound nonce.
  function getInboundNonce(uint64 sourceChainSelector, bytes calldata sender) external view returns (uint64) {
    return s_senderNonce[sourceChainSelector][sender];
  }

  function setInboundNonce(uint64 sourceChainSelector, uint64 testNonce, bytes calldata sender) external {
    s_senderNonce[sourceChainSelector][sender] = testNonce;
  }

  function getSourceChainConfig(
    uint64 sourceChainSelector
  ) external view returns (OffRamp.SourceChainConfig memory) {
    return s_sourceChainConfigs[sourceChainSelector];
  }

  function setSourceChainConfig(
    uint64 sourceChainSelector,
    OffRamp.SourceChainConfig memory sourceChainConfig
  ) external {
    s_sourceChainConfigs[sourceChainSelector] = sourceChainConfig;
  }

  /// @notice sets the sequence number of the last price update.
  function setLatestPriceSequenceNumber(
    uint64 seqNr
  ) external {
    s_latestPriceSequenceNumber = seqNr;
  }

  /// @notice Returns the sequence number of the last price update.
  /// @return sequenceNumber The latest price update sequence number.
  function getLatestPriceSequenceNumber() external view returns (uint64) {
    return s_latestPriceSequenceNumber;
  }

  function emitCCIPMessageSent(uint64 destChainSelector, Internal.EVM2AnyRampMessage memory message) external {
    emit OnRamp.CCIPMessageSent(destChainSelector, message.header.sequenceNumber, message);
  }

  function emitExecutionStateChanged(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    bytes32 messageId,
    bytes32 messageHash,
    Internal.MessageExecutionState state,
    bytes memory returnData,
    uint256 gasUsed
  ) external {
    emit OffRamp.ExecutionStateChanged(
      sourceChainSelector, sequenceNumber, messageId, messageHash, state, returnData, gasUsed
    );
  }

  function emitCommitReportAccepted(
    OffRamp.CommitReport memory report
  ) external {
    emit OffRamp.CommitReportAccepted(report.blessedMerkleRoots, report.unblessedMerkleRoots, report.priceUpdates);
  }
}
