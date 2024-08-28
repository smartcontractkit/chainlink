// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";

contract CCIPReaderTester {
  event CCIPSendRequested(uint64 indexed destChainSelector, Internal.EVM2AnyRampMessage message);

  mapping(uint64 sourceChainSelector => OffRamp.SourceChainConfig sourceChainConfig) internal s_sourceChainConfigs;

  function getSourceChainConfig(uint64 sourceChainSelector) external view returns (OffRamp.SourceChainConfig memory) {
    return s_sourceChainConfigs[sourceChainSelector];
  }

  function setSourceChainConfig(
    uint64 sourceChainSelector,
    OffRamp.SourceChainConfig memory sourceChainConfig
  ) external {
    s_sourceChainConfigs[sourceChainSelector] = sourceChainConfig;
  }

  function emitCCIPSendRequested(uint64 destChainSelector, Internal.EVM2AnyRampMessage memory message) external {
    emit CCIPSendRequested(destChainSelector, message);
  }

  event ExecutionStateChanged(
    uint64 indexed sourceChainSelector,
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );

  function emitExecutionStateChanged(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    bytes32 messageId,
    Internal.MessageExecutionState state,
    bytes memory returnData
  ) external {
    emit ExecutionStateChanged(sourceChainSelector, sequenceNumber, messageId, state, returnData);
  }

  event CommitReportAccepted(OffRamp.CommitReport report);

  function emitCommitReportAccepted(OffRamp.CommitReport memory report) external {
    emit CommitReportAccepted(report);
  }
}
