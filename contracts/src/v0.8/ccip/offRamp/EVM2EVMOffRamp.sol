// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {ICommitStore} from "../interfaces/ICommitStore.sol";
import {IARM} from "../interfaces/IARM.sol";
import {IPool} from "../interfaces/pools/IPool.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IAny2EVMMessageReceiver} from "../interfaces/IAny2EVMMessageReceiver.sol";
import {IAny2EVMOffRamp} from "../interfaces/IAny2EVMOffRamp.sol";

import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {CallWithExactGas} from "../../shared/call/CallWithExactGas.sol";
import {OCR2BaseNoChecks} from "../ocr/OCR2BaseNoChecks.sol";
import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {EnumerableMapAddresses} from "../../shared/enumerable/EnumerableMapAddresses.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {ERC165Checker} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165Checker.sol";

/// @notice EVM2EVMOffRamp enables OCR networks to execute multiple messages
/// in an OffRamp in a single transaction.
/// @dev The EVM2EVMOnRamp, CommitStore and EVM2EVMOffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
/// @dev OCR2BaseNoChecks is used to save gas, signatures are not required as the offramp can only execute
/// messages which are committed in the commitStore. We still make use of OCR2 as an executor whitelist
/// and turn-taking mechanism.
contract EVM2EVMOffRamp is IAny2EVMOffRamp, AggregateRateLimiter, ITypeAndVersion, OCR2BaseNoChecks {
  using ERC165Checker for address;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  error AlreadyAttempted(uint64 sequenceNumber);
  error AlreadyExecuted(uint64 sequenceNumber);
  error ZeroAddressNotAllowed();
  error CommitStoreAlreadyInUse();
  error ExecutionError(bytes error);
  error InvalidSourceChain(uint64 sourceChainSelector);
  error MessageTooLarge(uint256 maxSize, uint256 actualSize);
  error TokenDataMismatch(uint64 sequenceNumber);
  error UnexpectedTokenData();
  error UnsupportedNumberOfTokens(uint64 sequenceNumber);
  error ManualExecutionNotYetEnabled();
  error ManualExecutionGasLimitMismatch();
  error InvalidManualExecutionGasLimit(uint256 index, uint256 newLimit);
  error RootNotCommitted();
  error UnsupportedToken(IERC20 token);
  error CanOnlySelfCall();
  error ReceiverError(bytes error);
  error TokenHandlingError(bytes error);
  error EmptyReport();
  error BadARMSignal();
  error InvalidMessageId();
  error InvalidTokenPoolConfig();
  error PoolAlreadyAdded();
  error PoolDoesNotExist();
  error TokenPoolMismatch();
  error InvalidNewState(uint64 sequenceNumber, Internal.MessageExecutionState newState);

  event PoolAdded(address token, address pool);
  event PoolRemoved(address token, address pool);
  /// @dev Atlas depends on this event, if changing, please notify Atlas.
  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event SkippedIncorrectNonce(uint64 indexed nonce, address indexed sender);
  event SkippedSenderWithPreviousRampMessageInflight(uint64 indexed nonce, address indexed sender);
  /// @dev RMN depends on this event, if changing, please notify the RMN maintainers.
  event ExecutionStateChanged(
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );

  /// @notice Static offRamp config
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct StaticConfig {
    address commitStore; // ────────╮  CommitStore address on the destination chain
    uint64 chainSelector; // ───────╯  Destination chainSelector
    uint64 sourceChainSelector; // ─╮  Source chainSelector
    address onRamp; // ─────────────╯  OnRamp address on the source chain
    address prevOffRamp; //            Address of previous-version OffRamp
    address armProxy; //               ARM proxy address
  }

  /// @notice Dynamic offRamp config
  /// @dev since OffRampConfig is part of OffRampConfigChanged event, if changing it, we should update the ABI on Atlas
  struct DynamicConfig {
    uint32 permissionLessExecutionThresholdSeconds; // ─╮ Waiting time before manual execution is enabled
    address router; // ─────────────────────────────────╯ Router address
    address priceRegistry; // ──────────╮ Price registry address
    uint16 maxNumberOfTokensPerMsg; //  │ Maximum number of ERC20 token transfers that can be included per message
    uint32 maxDataBytes; //             │ Maximum payload data size in bytes
    uint32 maxPoolReleaseOrMintGas; // ─╯ Maximum amount of gas passed on to token pool when calling releaseOrMint
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMOffRamp 1.5.0-dev";
  /// @dev Commit store address on the destination chain
  address internal immutable i_commitStore;
  /// @dev ChainSelector of the source chain
  uint64 internal immutable i_sourceChainSelector;
  /// @dev ChainSelector of this chain
  uint64 internal immutable i_chainSelector;
  /// @dev OnRamp address on the source chain
  address internal immutable i_onRamp;
  /// @dev metadataHash is a lane-specific prefix for a message hash preimage which ensures global uniqueness.
  /// Ensures that 2 identical messages sent to 2 different lanes will have a distinct hash.
  /// Must match the metadataHash used in computing leaf hashes offchain for the root committed in
  /// the commitStore and i_metadataHash in the onRamp.
  bytes32 internal immutable i_metadataHash;
  /// @dev The address of previous-version OffRamp for this lane.
  /// Used to be able to provide sequencing continuity during a zero downtime upgrade.
  address internal immutable i_prevOffRamp;
  /// @dev The address of the arm proxy
  address internal immutable i_armProxy;

  // DYNAMIC CONFIG
  DynamicConfig internal s_dynamicConfig;
  /// @dev source token => token pool
  EnumerableMapAddresses.AddressToAddressMap private s_poolsBySourceToken;
  /// @dev dest token => token pool
  EnumerableMapAddresses.AddressToAddressMap private s_poolsByDestToken;

  // STATE
  /// @dev The expected nonce for a given sender.
  /// Corresponds to s_senderNonce in the OnRamp, used to enforce that messages are
  /// executed in the same order they are sent (assuming they are DON). Note that re-execution
  /// of FAILED messages however, can be out of order.
  mapping(address sender => uint64 nonce) internal s_senderNonce;
  /// @dev A mapping of sequence numbers to execution state using a bitmap with each execution
  /// state only taking up 2 bits of the uint256, packing 128 states into a single slot.
  /// Message state is tracked to ensure message can only be executed successfully once.
  mapping(uint64 seqNum => uint256 executionStateBitmap) internal s_executionStates;

  constructor(
    StaticConfig memory staticConfig,
    IERC20[] memory sourceTokens,
    IPool[] memory pools,
    RateLimiter.Config memory rateLimiterConfig
  ) OCR2BaseNoChecks() AggregateRateLimiter(rateLimiterConfig) {
    if (sourceTokens.length != pools.length) revert InvalidTokenPoolConfig();
    if (staticConfig.onRamp == address(0) || staticConfig.commitStore == address(0)) revert ZeroAddressNotAllowed();
    // Ensures we can never deploy a new offRamp that points to a commitStore that
    // already has roots committed.
    if (ICommitStore(staticConfig.commitStore).getExpectedNextSequenceNumber() != 1) revert CommitStoreAlreadyInUse();

    i_commitStore = staticConfig.commitStore;
    i_sourceChainSelector = staticConfig.sourceChainSelector;
    i_chainSelector = staticConfig.chainSelector;
    i_onRamp = staticConfig.onRamp;
    i_prevOffRamp = staticConfig.prevOffRamp;
    i_armProxy = staticConfig.armProxy;

    i_metadataHash = _metadataHash(Internal.EVM_2_EVM_MESSAGE_HASH);

    // Set new tokens and pools
    for (uint256 i = 0; i < sourceTokens.length; ++i) {
      s_poolsBySourceToken.set(address(sourceTokens[i]), address(pools[i]));
      s_poolsByDestToken.set(address(pools[i].getToken()), address(pools[i]));
      emit PoolAdded(address(sourceTokens[i]), address(pools[i]));
    }
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  // The size of the execution state in bits
  uint256 private constant MESSAGE_EXECUTION_STATE_BIT_WIDTH = 2;
  // The mask for the execution state bits
  uint256 private constant MESSAGE_EXECUTION_STATE_MASK = (1 << MESSAGE_EXECUTION_STATE_BIT_WIDTH) - 1;

  /// @notice Returns the current execution state of a message based on its sequenceNumber.
  /// @param sequenceNumber The sequence number of the message to get the execution state for.
  /// @return The current execution state of the message.
  /// @dev we use the literal number 128 because using a constant increased gas usage.
  function getExecutionState(uint64 sequenceNumber) public view returns (Internal.MessageExecutionState) {
    return
      Internal.MessageExecutionState(
        (s_executionStates[sequenceNumber / 128] >> ((sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH)) &
          MESSAGE_EXECUTION_STATE_MASK
      );
  }

  /// @notice Sets a new execution state for a given sequence number. It will overwrite any existing state.
  /// @param sequenceNumber The sequence number for which the state will be saved.
  /// @param newState The new value the state will be in after this function is called.
  /// @dev we use the literal number 128 because using a constant increased gas usage.
  function _setExecutionState(uint64 sequenceNumber, Internal.MessageExecutionState newState) internal {
    uint256 offset = (sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH;
    uint256 bitmap = s_executionStates[sequenceNumber / 128];
    // to unset any potential existing state we zero the bits of the section the state occupies,
    // then we do an AND operation to blank out any existing state for the section.
    bitmap &= ~(MESSAGE_EXECUTION_STATE_MASK << offset);
    // Set the new state
    bitmap |= uint256(newState) << offset;

    s_executionStates[sequenceNumber / 128] = bitmap;
  }

  /// @inheritdoc IAny2EVMOffRamp
  function getSenderNonce(address sender) public view returns (uint64 nonce) {
    uint256 senderNonce = s_senderNonce[sender];

    if (senderNonce == 0 && i_prevOffRamp != address(0)) {
      // If OffRamp was upgraded, check if sender has a nonce from the previous OffRamp.
      return IAny2EVMOffRamp(i_prevOffRamp).getSenderNonce(sender);
    }
    return uint64(senderNonce);
  }

  /// @notice Manually execute a message.
  /// @param report Internal.ExecutionReport.
  /// @param gasLimitOverrides New gasLimit for each message in the report.
  /// @dev We permit gas limit overrides so that users may manually execute messages which failed due to
  /// insufficient gas provided.
  function manuallyExecute(Internal.ExecutionReport memory report, uint256[] memory gasLimitOverrides) external {
    // We do this here because the other _execute path is already covered OCR2BaseXXX.
    if (i_chainID != block.chainid) revert OCR2BaseNoChecks.ForkedChain(i_chainID, uint64(block.chainid));

    uint256 numMsgs = report.messages.length;
    if (numMsgs != gasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();
    for (uint256 i = 0; i < numMsgs; ++i) {
      uint256 newLimit = gasLimitOverrides[i];
      // Checks to ensure message cannot be executed with less gas than specified.
      if (newLimit != 0 && newLimit < report.messages[i].gasLimit) revert InvalidManualExecutionGasLimit(i, newLimit);
    }

    _execute(report, gasLimitOverrides);
  }

  /// @notice Entrypoint for execution, called by the OCR network
  /// @dev Expects an encoded ExecutionReport
  function _report(bytes calldata report) internal override {
    _execute(abi.decode(report, (Internal.ExecutionReport)), new uint256[](0));
  }

  /// @notice Executes a report, executing each message in order.
  /// @param report The execution report containing the messages and proofs.
  /// @param manualExecGasLimits An array of gas limits to use for manual execution.
  /// @dev If called from the DON, this array is always empty.
  /// @dev If called from manual execution, this array is always same length as messages.
  function _execute(Internal.ExecutionReport memory report, uint256[] memory manualExecGasLimits) internal whenHealthy {
    uint256 numMsgs = report.messages.length;
    if (numMsgs == 0) revert EmptyReport();
    if (numMsgs != report.offchainTokenData.length) revert UnexpectedTokenData();

    bytes32[] memory hashedLeaves = new bytes32[](numMsgs);

    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      // We do this hash here instead of in _verifyMessages to avoid two separate loops
      // over the same data, which increases gas cost
      hashedLeaves[i] = Internal._hash(message, i_metadataHash);
      // For EVM2EVM offramps, the messageID is the leaf hash.
      // Asserting that this is true ensures we don't accidentally commit and then execute
      // a message with an unexpected hash.
      if (hashedLeaves[i] != message.messageId) revert InvalidMessageId();
    }

    // SECURITY CRITICAL CHECK
    uint256 timestampCommitted = ICommitStore(i_commitStore).verify(hashedLeaves, report.proofs, report.proofFlagBits);
    if (timestampCommitted == 0) revert RootNotCommitted();

    // Execute messages
    bool manualExecution = manualExecGasLimits.length != 0;
    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      Internal.MessageExecutionState originalState = getExecutionState(message.sequenceNumber);
      // Two valid cases here, we either have never touched this message before, or we tried to execute
      // and failed. This check protects against reentry and re-execution because the other states are
      // IN_PROGRESS and SUCCESS, both should not be allowed to execute.
      if (
        !(originalState == Internal.MessageExecutionState.UNTOUCHED ||
          originalState == Internal.MessageExecutionState.FAILURE)
      ) revert AlreadyExecuted(message.sequenceNumber);

      if (manualExecution) {
        bool isOldCommitReport = (block.timestamp - timestampCommitted) >
          s_dynamicConfig.permissionLessExecutionThresholdSeconds;
        // Manually execution is fine if we previously failed or if the commit report is just too old
        // Acceptable state transitions: FAILURE->SUCCESS, UNTOUCHED->SUCCESS, FAILURE->FAILURE
        if (!(isOldCommitReport || originalState == Internal.MessageExecutionState.FAILURE))
          revert ManualExecutionNotYetEnabled();

        // Manual execution gas limit can override gas limit specified in the message. Value of 0 indicates no override.
        if (manualExecGasLimits[i] != 0) {
          message.gasLimit = manualExecGasLimits[i];
        }
      } else {
        // DON can only execute a message once
        // Acceptable state transitions: UNTOUCHED->SUCCESS, UNTOUCHED->FAILURE
        if (originalState != Internal.MessageExecutionState.UNTOUCHED) revert AlreadyAttempted(message.sequenceNumber);
      }

      // In the scenario where we upgrade offRamps, we still want to have sequential nonces.
      // Referencing the old offRamp to check the expected nonce if none is set for a
      // given sender allows us to skip the current message if it would not be the next according
      // to the old offRamp. This preserves sequencing between updates.
      uint64 prevNonce = s_senderNonce[message.sender];
      if (prevNonce == 0 && i_prevOffRamp != address(0)) {
        prevNonce = IAny2EVMOffRamp(i_prevOffRamp).getSenderNonce(message.sender);
        if (prevNonce + 1 != message.nonce) {
          // the starting v2 onramp nonce, i.e. the 1st message nonce v2 offramp is expected to receive,
          // is guaranteed to equal (largest v1 onramp nonce + 1).
          // if this message's nonce isn't (v1 offramp nonce + 1), then v1 offramp nonce != largest v1 onramp nonce,
          // it tells us there are still messages inflight for v1 offramp
          emit SkippedSenderWithPreviousRampMessageInflight(message.nonce, message.sender);
          continue;
        }
        // Otherwise this nonce is indeed the "transitional nonce", that is
        // all messages sent to v1 ramp have been executed by the DON and the sequence can resume in V2.
        // Note if first time user in V2, then prevNonce will be 0, and message.nonce = 1, so this will be a no-op.
        s_senderNonce[message.sender] = prevNonce;
      }

      // UNTOUCHED messages MUST be executed in order always
      if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
        if (prevNonce + 1 != message.nonce) {
          // We skip the message if the nonce is incorrect
          emit SkippedIncorrectNonce(message.nonce, message.sender);
          continue;
        }
      }

      // Although we expect only valid messages will be committed, we check again
      // when executing as a defense in depth measure.
      bytes[] memory offchainTokenData = report.offchainTokenData[i];
      _isWellFormed(
        message.sequenceNumber,
        message.sourceChainSelector,
        message.tokenAmounts.length,
        message.data.length,
        offchainTokenData.length
      );

      _setExecutionState(message.sequenceNumber, Internal.MessageExecutionState.IN_PROGRESS);
      (Internal.MessageExecutionState newState, bytes memory returnData) = _trialExecute(message, offchainTokenData);
      _setExecutionState(message.sequenceNumber, newState);

      // Since it's hard to estimate whether manual execution will succeed, we
      // revert the entire transaction if it fails. This will show the user if
      // their manual exec will fail before they submit it.
      if (manualExecution && newState == Internal.MessageExecutionState.FAILURE) {
        // If manual execution fails, we revert the entire transaction.
        revert ExecutionError(returnData);
      }

      // The only valid prior states are UNTOUCHED and FAILURE (checked above)
      // The only valid post states are FAILURE and SUCCESS (checked below)
      if (newState != Internal.MessageExecutionState.FAILURE && newState != Internal.MessageExecutionState.SUCCESS)
        revert InvalidNewState(message.sequenceNumber, newState);

      // Nonce changes per state transition
      // UNTOUCHED -> FAILURE  nonce bump
      // UNTOUCHED -> SUCCESS  nonce bump
      // FAILURE   -> FAILURE  no nonce bump
      // FAILURE   -> SUCCESS  no nonce bump
      if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
        s_senderNonce[message.sender]++;
      }

      emit ExecutionStateChanged(message.sequenceNumber, message.messageId, newState, returnData);
    }
  }

  /// @notice Does basic message validation. Should never fail.
  /// @param sequenceNumber Sequence number of the message.
  /// @param sourceChainSelector SourceChainSelector of the message.
  /// @param numberOfTokens Length of tokenAmounts array in the message.
  /// @param dataLength Length of data field in the message.
  /// @param offchainTokenDataLength Length of offchainTokenData array.
  /// @dev reverts on validation failures.
  function _isWellFormed(
    uint64 sequenceNumber,
    uint64 sourceChainSelector,
    uint256 numberOfTokens,
    uint256 dataLength,
    uint256 offchainTokenDataLength
  ) private view {
    if (sourceChainSelector != i_sourceChainSelector) revert InvalidSourceChain(sourceChainSelector);
    if (numberOfTokens > uint256(s_dynamicConfig.maxNumberOfTokensPerMsg))
      revert UnsupportedNumberOfTokens(sequenceNumber);
    if (numberOfTokens != offchainTokenDataLength) revert TokenDataMismatch(sequenceNumber);
    if (dataLength > uint256(s_dynamicConfig.maxDataBytes))
      revert MessageTooLarge(uint256(s_dynamicConfig.maxDataBytes), dataLength);
  }

  /// @notice Try executing a message.
  /// @param message Internal.EVM2EVMMessage memory message.
  /// @param offchainTokenData Data provided by the DON for token transfers.
  /// @return the new state of the message, being either SUCCESS or FAILURE.
  /// @return revert data in bytes if CCIP receiver reverted during execution.
  function _trialExecute(
    Internal.EVM2EVMMessage memory message,
    bytes[] memory offchainTokenData
  ) internal returns (Internal.MessageExecutionState, bytes memory) {
    try this.executeSingleMessage(message, offchainTokenData) {} catch (bytes memory err) {
      if (ReceiverError.selector == bytes4(err) || TokenHandlingError.selector == bytes4(err)) {
        // If CCIP receiver execution is not successful, bubble up receiver revert data,
        // prepended by the 4 bytes of ReceiverError.selector or TokenHandlingError.selector
        // Max length of revert data is Router.MAX_RET_BYTES, max length of err is 4 + Router.MAX_RET_BYTES
        return (Internal.MessageExecutionState.FAILURE, err);
      } else {
        // If revert is not caused by CCIP receiver, it is unexpected, bubble up the revert.
        revert ExecutionError(err);
      }
    }
    // If message execution succeeded, no CCIP receiver return data is expected, return with empty bytes.
    return (Internal.MessageExecutionState.SUCCESS, "");
  }

  /// @notice Execute a single message.
  /// @param message The message that will be executed.
  /// @param offchainTokenData Token transfer data to be passed to TokenPool.
  /// @dev We make this external and callable by the contract itself, in order to try/catch
  /// its execution and enforce atomicity among successful message processing and token transfer.
  /// @dev We use ERC-165 to check for the ccipReceive interface to permit sending tokens to contracts
  /// (for example smart contract wallets) without an associated message.
  function executeSingleMessage(Internal.EVM2EVMMessage memory message, bytes[] memory offchainTokenData) external {
    if (msg.sender != address(this)) revert CanOnlySelfCall();
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](0);
    if (message.tokenAmounts.length > 0) {
      destTokenAmounts = _releaseOrMintTokens(
        message.tokenAmounts,
        abi.encode(message.sender),
        message.receiver,
        message.sourceTokenData,
        offchainTokenData
      );
    }
    if (
      message.receiver.code.length == 0 ||
      !message.receiver.supportsInterface(type(IAny2EVMMessageReceiver).interfaceId)
    ) return;

    (bool success, bytes memory returnData, ) = IRouter(s_dynamicConfig.router).routeMessage(
      Internal._toAny2EVMMessage(message, destTokenAmounts),
      Internal.GAS_FOR_CALL_EXACT_CHECK,
      message.gasLimit,
      message.receiver
    );
    // If CCIP receiver execution is not successful, revert the call including token transfers
    if (!success) revert ReceiverError(returnData);
  }

  /// @notice creates a unique hash to be used in message hashing.
  function _metadataHash(bytes32 prefix) internal view returns (bytes32) {
    return keccak256(abi.encode(prefix, i_sourceChainSelector, i_chainSelector, i_onRamp));
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static config.
  /// @dev This function will always return the same struct as the contents is static and can never change.
  /// RMN depends on this function, if changing, please notify the RMN maintainers.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return
      StaticConfig({
        commitStore: i_commitStore,
        chainSelector: i_chainSelector,
        sourceChainSelector: i_sourceChainSelector,
        onRamp: i_onRamp,
        prevOffRamp: i_prevOffRamp,
        armProxy: i_armProxy
      });
  }

  /// @notice Returns the current dynamic config.
  /// @return The current config.
  function getDynamicConfig() external view returns (DynamicConfig memory) {
    return s_dynamicConfig;
  }

  /// @notice Sets the dynamic config. This function is called during `setOCR2Config` flow
  function _beforeSetConfig(bytes memory onchainConfig) internal override {
    DynamicConfig memory dynamicConfig = abi.decode(onchainConfig, (DynamicConfig));

    if (dynamicConfig.router == address(0)) revert ZeroAddressNotAllowed();

    s_dynamicConfig = dynamicConfig;

    emit ConfigSet(
      StaticConfig({
        commitStore: i_commitStore,
        chainSelector: i_chainSelector,
        sourceChainSelector: i_sourceChainSelector,
        onRamp: i_onRamp,
        prevOffRamp: i_prevOffRamp,
        armProxy: i_armProxy
      }),
      dynamicConfig
    );
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @notice Get all supported source tokens
  /// @return sourceTokens of supported source tokens
  function getSupportedTokens() external view returns (IERC20[] memory sourceTokens) {
    sourceTokens = new IERC20[](s_poolsBySourceToken.length());
    for (uint256 i = 0; i < sourceTokens.length; ++i) {
      (address token, ) = s_poolsBySourceToken.at(i);
      sourceTokens[i] = IERC20(token);
    }
    return sourceTokens;
  }

  /// @notice Get a token pool by its source token
  /// @param sourceToken token
  /// @return Token Pool
  function getPoolBySourceToken(IERC20 sourceToken) public view returns (IPool) {
    (bool success, address pool) = s_poolsBySourceToken.tryGet(address(sourceToken));
    if (!success) revert UnsupportedToken(sourceToken);
    return IPool(pool);
  }

  /// @notice Get the destination token from the pool based on a given source token.
  /// @param sourceToken The source token
  /// @return the destination token
  function getDestinationToken(IERC20 sourceToken) external view returns (IERC20) {
    return getPoolBySourceToken(sourceToken).getToken();
  }

  /// @notice Get a token pool by its dest token
  /// @param destToken token
  /// @return Token Pool
  function getPoolByDestToken(IERC20 destToken) external view returns (IPool) {
    (bool success, address pool) = s_poolsByDestToken.tryGet(address(destToken));
    if (!success) revert UnsupportedToken(destToken);
    return IPool(pool);
  }

  /// @notice Get all configured destination tokens
  /// @return destTokens Array of configured destination tokens
  function getDestinationTokens() external view returns (IERC20[] memory destTokens) {
    destTokens = new IERC20[](s_poolsByDestToken.length());
    for (uint256 i = 0; i < destTokens.length; ++i) {
      (address token, ) = s_poolsByDestToken.at(i);
      destTokens[i] = IERC20(token);
    }
    return destTokens;
  }

  /// @notice Adds and removed token pools.
  /// @param removes The tokens and pools to be removed
  /// @param adds The tokens and pools to be added.
  function applyPoolUpdates(
    Internal.PoolUpdate[] calldata removes,
    Internal.PoolUpdate[] calldata adds
  ) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      address token = removes[i].token;
      address pool = removes[i].pool;

      // Check if the pool exists
      if (!s_poolsBySourceToken.contains(token)) revert PoolDoesNotExist();
      // Sanity check
      if (s_poolsBySourceToken.get(token) != pool) revert TokenPoolMismatch();

      s_poolsBySourceToken.remove(token);
      s_poolsByDestToken.remove(address(IPool(pool).getToken()));

      emit PoolRemoved(token, pool);
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      address token = adds[i].token;
      address pool = adds[i].pool;

      if (token == address(0) || pool == address(0)) revert InvalidTokenPoolConfig();
      // Check if the pool is already set
      if (s_poolsBySourceToken.contains(token)) revert PoolAlreadyAdded();

      // Set the s_pools with new config values
      s_poolsBySourceToken.set(token, pool);
      s_poolsByDestToken.set(address(IPool(pool).getToken()), pool);

      emit PoolAdded(token, pool);
    }
  }

  /// @notice Uses pools to release or mint a number of different tokens to a receiver address.
  /// @param sourceTokenAmounts List of tokens and amount values to be released/minted.
  /// @param originalSender The message sender.
  /// @param receiver The address that will receive the tokens.
  /// @param sourceTokenData Array of token data returned by token pools on the source chain.
  /// @param offchainTokenData Array of token data fetched offchain by the DON.
  /// @dev This function wrappes the token pool call in a try catch block to gracefully handle
  /// any non-rate limiting errors that may occur. If we encounter a rate limiting related error
  /// we bubble it up. If we encounter a non-rate limiting error we wrap it in a TokenHandlingError.
  function _releaseOrMintTokens(
    Client.EVMTokenAmount[] memory sourceTokenAmounts,
    bytes memory originalSender,
    address receiver,
    bytes[] memory sourceTokenData,
    bytes[] memory offchainTokenData
  ) internal returns (Client.EVMTokenAmount[] memory) {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](sourceTokenAmounts.length);
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      IPool pool = getPoolBySourceToken(IERC20(sourceTokenAmounts[i].token));
      uint256 sourceTokenAmount = sourceTokenAmounts[i].amount;

      // Call the pool with exact gas to increase resistance against malicious tokens or token pools.
      // _callWithExactGas also protects against return data bombs by capping the return data size
      // at MAX_RET_BYTES.
      (bool success, bytes memory returnData, ) = CallWithExactGas._callWithExactGasSafeReturnData(
        abi.encodeWithSelector(
          pool.releaseOrMint.selector,
          originalSender,
          receiver,
          sourceTokenAmount,
          i_sourceChainSelector,
          abi.encode(sourceTokenData[i], offchainTokenData[i])
        ),
        address(pool),
        s_dynamicConfig.maxPoolReleaseOrMintGas,
        Internal.GAS_FOR_CALL_EXACT_CHECK,
        Internal.MAX_RET_BYTES
      );

      // wrap and rethrow the error so we can catch it lower in the stack
      if (!success) revert TokenHandlingError(returnData);

      destTokenAmounts[i].token = address(pool.getToken());
      destTokenAmounts[i].amount = sourceTokenAmount;
    }
    _rateLimitValue(destTokenAmounts, IPriceRegistry(s_dynamicConfig.priceRegistry));
    return destTokenAmounts;
  }

  // ================================================================
  // │                        Access and ARM                        │
  // ================================================================

  /// @notice Reverts as this contract should not access CCIP messages
  function ccipReceive(Client.Any2EVMMessage calldata) external pure {
    /* solhint-disable */
    revert();
    /* solhint-enable*/
  }

  /// @notice Ensure that the ARM has not emitted a bad signal, and that the latest heartbeat is not stale.
  modifier whenHealthy() {
    if (IARM(i_armProxy).isCursed()) revert BadARMSignal();
    _;
  }
}
