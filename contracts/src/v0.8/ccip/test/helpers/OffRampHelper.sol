// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract OffRampHelper is OffRamp, IgnoreContractSize {
  using EnumerableSet for EnumerableSet.UintSet;

  mapping(uint64 sourceChainSelector => uint256 overrideTimestamp) private s_sourceChainVerificationOverride;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs
  ) OffRamp(staticConfig, dynamicConfig, sourceChainConfigs) {}

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

  function releaseOrMintSingleToken(
    Internal.Any2EVMTokenTransfer calldata sourceTokenAmount,
    bytes calldata originalSender,
    address receiver,
    uint64 sourceChainSelector,
    bytes calldata offchainTokenData
  ) external returns (Client.EVMTokenAmount memory) {
    return
      _releaseOrMintSingleToken(sourceTokenAmount, originalSender, receiver, sourceChainSelector, offchainTokenData);
  }

  function releaseOrMintTokens(
    Internal.Any2EVMTokenTransfer[] calldata sourceTokenAmounts,
    bytes calldata originalSender,
    address receiver,
    uint64 sourceChainSelector,
    bytes[] calldata offchainTokenData,
    uint32[] calldata tokenGasOverrides
  ) external returns (Client.EVMTokenAmount[] memory) {
    return _releaseOrMintTokens(
      sourceTokenAmounts, originalSender, receiver, sourceChainSelector, offchainTokenData, tokenGasOverrides
    );
  }

  function trialExecute(
    Internal.Any2EVMRampMessage memory message,
    bytes[] memory offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) external returns (Internal.MessageExecutionState, bytes memory) {
    return _trialExecute(message, offchainTokenData, tokenGasOverrides);
  }

  function executeSingleReport(
    Internal.ExecutionReport memory rep,
    GasLimitOverride[] memory manualExecGasExecOverrides
  ) external {
    _executeSingleReport(rep, manualExecGasExecOverrides);
  }

  function batchExecute(
    Internal.ExecutionReport[] memory reports,
    GasLimitOverride[][] memory manualExecGasLimits
  ) external {
    _batchExecute(reports, manualExecGasLimits);
  }

  function verify(
    uint64 sourceChainSelector,
    bytes32[] memory hashedLeaves,
    bytes32[] memory proofs,
    uint256 proofFlagBits
  ) external view returns (uint256 timestamp) {
    return super._verify(sourceChainSelector, hashedLeaves, proofs, proofFlagBits);
  }

  function _verify(
    uint64 sourceChainSelector,
    bytes32[] memory hashedLeaves,
    bytes32[] memory proofs,
    uint256 proofFlagBits
  ) internal view override returns (uint256 timestamp) {
    uint256 overrideTimestamp = s_sourceChainVerificationOverride[sourceChainSelector];

    return overrideTimestamp == 0
      ? super._verify(sourceChainSelector, hashedLeaves, proofs, proofFlagBits)
      : overrideTimestamp;
  }

  /// @dev Test helper to override _verify result for easier exec testing
  function setVerifyOverrideResult(uint64 sourceChainSelector, uint256 overrideTimestamp) external {
    s_sourceChainVerificationOverride[sourceChainSelector] = overrideTimestamp;
  }

  /// @dev Test helper to directly set a root's timestamp
  function setRootTimestamp(uint64 sourceChainSelector, bytes32 root, uint256 timestamp) external {
    s_roots[sourceChainSelector][root] = timestamp;
  }

  function getSourceChainSelectors() external view returns (uint256[] memory chainSelectors) {
    return s_sourceChainSelectors.values();
  }
}
