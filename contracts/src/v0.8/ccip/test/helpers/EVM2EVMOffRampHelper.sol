// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../offRamp/EVM2EVMOffRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMOffRampHelper is EVM2EVMOffRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    RateLimiter.Config memory rateLimiterConfig
  ) EVM2EVMOffRamp(staticConfig, rateLimiterConfig) {}

  function setExecutionStateHelper(uint64 sequenceNumber, Internal.MessageExecutionState state) public {
    _setExecutionState(sequenceNumber, state);
  }

  function getExecutionStateBitMap(uint64 bitmapIndex) public view returns (uint256) {
    return s_executionStates[bitmapIndex];
  }

  function releaseOrMintToken(
    uint256 sourceTokenAmount,
    bytes calldata originalSender,
    address receiver,
    Internal.SourceTokenData calldata sourceTokenData,
    bytes calldata offchainTokenData
  ) external returns (Client.EVMTokenAmount memory) {
    return _releaseOrMintToken(sourceTokenAmount, originalSender, receiver, sourceTokenData, offchainTokenData);
  }

  function releaseOrMintTokens(
    Client.EVMTokenAmount[] calldata sourceTokenAmounts,
    bytes calldata originalSender,
    address receiver,
    bytes[] calldata sourceTokenData,
    bytes[] calldata offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) external returns (Client.EVMTokenAmount[] memory) {
    return _releaseOrMintTokens(
      sourceTokenAmounts, originalSender, receiver, sourceTokenData, offchainTokenData, tokenGasOverrides
    );
  }

  function trialExecute(
    Internal.EVM2EVMMessage memory message,
    bytes[] memory offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) external returns (Internal.MessageExecutionState, bytes memory) {
    return _trialExecute(message, offchainTokenData, tokenGasOverrides);
  }

  function report(bytes calldata executableMessages) external {
    _report(executableMessages);
  }

  function execute(Internal.ExecutionReport memory rep, GasLimitOverride[] memory gasLimitOverrides) external {
    _execute(rep, gasLimitOverrides);
  }

  function metadataHash() external view returns (bytes32) {
    return _metadataHash(Internal.EVM_2_EVM_MESSAGE_HASH);
  }
}
