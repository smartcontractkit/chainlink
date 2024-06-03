// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {EVM2EVMMultiOffRamp} from "../../offRamp/EVM2EVMMultiOffRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMMultiOffRampHelper is EVM2EVMMultiOffRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs
  ) EVM2EVMMultiOffRamp(staticConfig, sourceChainConfigs) {}

  function setExecutionStateHelper(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    Internal.MessageExecutionState state
  ) public {
    _setExecutionState(sourceChainSelector, sequenceNumber, state);
  }

  function getExecutionStateBitMap(uint64 sourceChainSelector, uint64 bitmapIndex) public view returns (uint256) {
    return s_executionStates[sourceChainSelector][bitmapIndex];
  }

  function metadataHash(uint64 sourceChainSelector, address onRamp) external view returns (bytes32) {
    return _metadataHash(sourceChainSelector, onRamp, Internal.EVM_2_EVM_MESSAGE_HASH);
  }

  function releaseOrMintTokens(
    Client.EVMTokenAmount[] memory sourceTokenAmounts,
    EVM2EVMMultiOffRamp.Any2EVMMessageRoute memory messageRoute,
    bytes[] calldata sourceTokenData,
    bytes[] calldata offchainTokenData
  ) external returns (Client.EVMTokenAmount[] memory) {
    return _releaseOrMintTokens(sourceTokenAmounts, messageRoute, sourceTokenData, offchainTokenData);
  }

  function trialExecute(
    Internal.EVM2EVMMessage memory message,
    bytes[] memory offchainTokenData
  ) external returns (Internal.MessageExecutionState, bytes memory) {
    return _trialExecute(message, offchainTokenData);
  }

  function report(bytes calldata executableReports) external {
    _report(executableReports);
  }

  function execute(Internal.ExecutionReportSingleChain memory rep, uint256[] memory manualExecGasLimits) external {
    _execute(rep, manualExecGasLimits);
  }

  function batchExecute(
    Internal.ExecutionReportSingleChain[] memory reports,
    uint256[][] memory manualExecGasLimits
  ) external {
    _batchExecute(reports, manualExecGasLimits);
  }
}
