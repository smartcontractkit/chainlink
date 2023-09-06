// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../../offRamp/EVM2EVMOffRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMOffRampHelper is EVM2EVMOffRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    IERC20[] memory sourceTokens,
    IPool[] memory pools,
    RateLimiter.Config memory rateLimiterConfig
  ) EVM2EVMOffRamp(staticConfig, sourceTokens, pools, rateLimiterConfig) {}

  function setExecutionStateHelper(uint64 sequenceNumber, Internal.MessageExecutionState state) public {
    _setExecutionState(sequenceNumber, state);
  }

  function getExecutionStateBitMap(uint64 bitmapIndex) public view returns (uint256) {
    return s_executionStates[bitmapIndex];
  }

  function releaseOrMintTokens(
    Client.EVMTokenAmount[] memory sourceTokenAmounts,
    bytes calldata originalSender,
    address receiver,
    bytes[] calldata sourceTokenData,
    bytes[] calldata offchainTokenData
  ) external returns (Client.EVMTokenAmount[] memory) {
    return _releaseOrMintTokens(sourceTokenAmounts, originalSender, receiver, sourceTokenData, offchainTokenData);
  }

  function trialExecute(
    Internal.EVM2EVMMessage memory message,
    bytes[] memory offchainTokenData
  ) external returns (Internal.MessageExecutionState, bytes memory) {
    return _trialExecute(message, offchainTokenData);
  }

  function report(bytes calldata executableMessages) external {
    _report(executableMessages);
  }

  function execute(Internal.ExecutionReport memory rep, uint256[] memory manualExecGasLimits) external {
    _execute(rep, manualExecGasLimits);
  }

  function metadataHash() external view returns (bytes32) {
    return _metadataHash(Internal.EVM_2_EVM_MESSAGE_HASH);
  }
}
