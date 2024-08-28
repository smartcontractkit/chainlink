// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "../interfaces/IAny2EVMMessageReceiver.sol";
import {IAny2EVMOffRamp} from "../interfaces/IAny2EVMOffRamp.sol";
import {ICommitStore} from "../interfaces/ICommitStore.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../shared/call/CallWithExactGas.sol";
import {EnumerableMapAddresses} from "../../shared/enumerable/EnumerableMapAddresses.sol";
import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {OCR2BaseNoChecks} from "../ocr/OCR2BaseNoChecks.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {ERC165Checker} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/ERC165Checker.sol";

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

  error ZeroAddressNotAllowed();
  error CommitStoreAlreadyInUse();
  error ExecutionError(bytes err);
  error InvalidSourceChain(uint64 sourceChainSelector);
  error MessageTooLarge(uint256 maxSize, uint256 actualSize);
  error TokenDataMismatch(uint64 sequenceNumber);
  error UnexpectedTokenData();
  error UnsupportedNumberOfTokens(uint64 sequenceNumber);
  error ManualExecutionNotYetEnabled();
  error ManualExecutionGasLimitMismatch();
  error DestinationGasAmountCountMismatch(bytes32 messageId, uint64 sequenceNumber);
  error InvalidManualExecutionGasLimit(bytes32 messageId, uint256 oldLimit, uint256 newLimit);
  error InvalidTokenGasOverride(bytes32 messageId, uint256 tokenIndex, uint256 oldLimit, uint256 tokenGasOverride);
  error RootNotCommitted();
  error CanOnlySelfCall();
  error ReceiverError(bytes err);
  error TokenHandlingError(bytes err);
  error ReleaseOrMintBalanceMismatch(uint256 amountReleased, uint256 balancePre, uint256 balancePost);
  error EmptyReport();
  error CursedByRMN();
  error InvalidMessageId();
  error NotACompatiblePool(address notPool);
  error InvalidDataLength(uint256 expected, uint256 got);
  error InvalidNewState(uint64 sequenceNumber, Internal.MessageExecutionState newState);

  /// @dev Atlas depends on this event, if changing, please notify Atlas.
  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event SkippedIncorrectNonce(uint64 indexed nonce, address indexed sender);
  event SkippedSenderWithPreviousRampMessageInflight(uint64 indexed nonce, address indexed sender);
  /// @dev RMN depends on this event, if changing, please notify the RMN maintainers.
  event ExecutionStateChanged(
    uint64 indexed sequenceNumber, bytes32 indexed messageId, Internal.MessageExecutionState state, bytes returnData
  );
  event TokenAggregateRateLimitAdded(address sourceToken, address destToken);
  event TokenAggregateRateLimitRemoved(address sourceToken, address destToken);
  event SkippedAlreadyExecutedMessage(uint64 indexed sequenceNumber);
  event AlreadyAttempted(uint64 sequenceNumber);

  /// @notice Static offRamp config
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  //solhint-disable gas-struct-packing
  struct StaticConfig {
    address commitStore; // ────────╮  CommitStore address on the destination chain
    uint64 chainSelector; // ───────╯  Destination chainSelector
    uint64 sourceChainSelector; // ─╮  Source chainSelector
    address onRamp; // ─────────────╯  OnRamp address on the source chain
    address prevOffRamp; //            Address of previous-version OffRamp
    address rmnProxy; //               RMN proxy address
    address tokenAdminRegistry; //     Token admin registry address
  }

  /// @notice Dynamic offRamp config
  /// @dev since OffRampConfig is part of OffRampConfigChanged event, if changing it, we should update the ABI on Atlas
  struct DynamicConfig {
    uint32 permissionLessExecutionThresholdSeconds; // ─╮ Waiting time before manual execution is enabled
    uint32 maxDataBytes; //                             │ Maximum payload data size in bytes
    uint16 maxNumberOfTokensPerMsg; //                  │ Maximum number of ERC20 token transfers that can be included per message
    address router; // ─────────────────────────────────╯ Router address
    address priceRegistry; //                             Price registry address
  }

  /// @notice RateLimitToken struct containing both the source and destination token addresses
  struct RateLimitToken {
    address sourceToken;
    address destToken;
  }

  /// @notice Gas overrides for manual exec, the number of token overrides must match the number of tokens in the msg.
  struct GasLimitOverride {
    /// @notice Overrides EVM2EVMMessage.gasLimit. A value of zero indicates no override and is valid.
    uint256 receiverExecutionGasLimit;
    /// @notice Overrides EVM2EVMMessage.sourceTokenData.destGasAmount. Must be same length as tokenAmounts. A value
    /// of zero indicates no override and is valid.
    uint32[] tokenGasOverrides;
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMOffRamp 1.5.0";

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
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;
  /// @dev The address of the token admin registry
  address internal immutable i_tokenAdminRegistry;

  // DYNAMIC CONFIG
  DynamicConfig internal s_dynamicConfig;
  /// @dev Tokens that should be included in Aggregate Rate Limiting
  /// An (address => address) map is used for backwards compatability of offchain code
  EnumerableMapAddresses.AddressToAddressMap internal s_rateLimitedTokensDestToSource;

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
    RateLimiter.Config memory rateLimiterConfig
  ) OCR2BaseNoChecks() AggregateRateLimiter(rateLimiterConfig) {
    if (
      staticConfig.onRamp == address(0) || staticConfig.commitStore == address(0)
        || staticConfig.tokenAdminRegistry == address(0)
    ) revert ZeroAddressNotAllowed();
    // Ensures we can never deploy a new offRamp that points to a commitStore that
    // already has roots committed.
    if (ICommitStore(staticConfig.commitStore).getExpectedNextSequenceNumber() != 1) revert CommitStoreAlreadyInUse();

    i_commitStore = staticConfig.commitStore;
    i_sourceChainSelector = staticConfig.sourceChainSelector;
    i_chainSelector = staticConfig.chainSelector;
    i_onRamp = staticConfig.onRamp;
    i_prevOffRamp = staticConfig.prevOffRamp;
    i_rmnProxy = staticConfig.rmnProxy;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;

    i_metadataHash = _metadataHash(Internal.EVM_2_EVM_MESSAGE_HASH);
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
    return Internal.MessageExecutionState(
      (s_executionStates[sequenceNumber / 128] >> ((sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH))
        & MESSAGE_EXECUTION_STATE_MASK
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
  function getSenderNonce(address sender) external view returns (uint64 nonce) {
    uint256 senderNonce = s_senderNonce[sender];

    if (senderNonce == 0) {
      if (i_prevOffRamp != address(0)) {
        // If OffRamp was upgraded, check if sender has a nonce from the previous OffRamp.
        return IAny2EVMOffRamp(i_prevOffRamp).getSenderNonce(sender);
      }
    }
    return uint64(senderNonce);
  }

  /// @notice Manually execute a message.
  /// @param report Internal.ExecutionReport.
  /// @param gasLimitOverrides New gasLimit for each message in the report.
  /// @dev We permit gas limit overrides so that users may manually execute messages which failed due to
  /// insufficient gas provided.
  function manuallyExecute(
    Internal.ExecutionReport memory report,
    GasLimitOverride[] memory gasLimitOverrides
  ) external {
    // We do this here because the other _execute path is already covered OCR2BaseXXX.
    _checkChainForked();

    uint256 numMsgs = report.messages.length;
    if (numMsgs != gasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();
    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      GasLimitOverride memory gasLimitOverride = gasLimitOverrides[i];

      uint256 newLimit = gasLimitOverride.receiverExecutionGasLimit;
      // Checks to ensure message cannot be executed with less gas than specified.
      if (newLimit != 0) {
        if (newLimit < message.gasLimit) {
          revert InvalidManualExecutionGasLimit(message.messageId, message.gasLimit, newLimit);
        }
      }

      if (message.tokenAmounts.length != gasLimitOverride.tokenGasOverrides.length) {
        revert DestinationGasAmountCountMismatch(message.messageId, message.sequenceNumber);
      }

      bytes[] memory encodedSourceTokenData = message.sourceTokenData;

      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        Internal.SourceTokenData memory sourceTokenData =
          abi.decode(encodedSourceTokenData[i], (Internal.SourceTokenData));
        uint256 tokenGasOverride = gasLimitOverride.tokenGasOverrides[j];

        // The gas limit can not be lowered as that could cause the message to fail. If manual execution is done
        // from an UNTOUCHED state and we would allow lower gas limit, anyone could grief by executing the message with
        // lower gas limit than the DON would have used. This results in the message being marked FAILURE and the DON
        // would not attempt it with the correct gas limit.
        if (tokenGasOverride != 0 && tokenGasOverride < sourceTokenData.destGasAmount) {
          revert InvalidTokenGasOverride(message.messageId, j, sourceTokenData.destGasAmount, tokenGasOverride);
        }
      }
    }

    _execute(report, gasLimitOverrides);
  }

  /// @notice Entrypoint for execution, called by the OCR network
  /// @dev Expects an encoded ExecutionReport
  /// @dev Supplies no GasLimitOverrides as the DON will only execute with the original gas limits.
  function _report(bytes calldata report) internal override {
    _execute(abi.decode(report, (Internal.ExecutionReport)), new GasLimitOverride[](0));
  }

  /// @notice Executes a report, executing each message in order.
  /// @param report The execution report containing the messages and proofs.
  /// @param manualExecGasOverrides An array of gas limits to use for manual execution.
  /// @dev If called from the DON, this array is always empty.
  /// @dev If called from manual execution, this array is always same length as messages.
  function _execute(Internal.ExecutionReport memory report, GasLimitOverride[] memory manualExecGasOverrides) internal {
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(i_sourceChainSelector)))) revert CursedByRMN();

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
    bool manualExecution = manualExecGasOverrides.length != 0;

    // SECURITY CRITICAL CHECK
    uint256 timestampCommitted = ICommitStore(i_commitStore).verify(hashedLeaves, report.proofs, report.proofFlagBits);
    if (timestampCommitted == 0) revert RootNotCommitted();

    // Execute messages
    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      Internal.MessageExecutionState originalState = getExecutionState(message.sequenceNumber);
      // Two valid cases here, we either have never touched this message before, or we tried to execute
      // and failed. This check protects against reentry and re-execution because the other state is
      // IN_PROGRESS which should not be allowed to execute.
      if (
        !(
          originalState == Internal.MessageExecutionState.UNTOUCHED
            || originalState == Internal.MessageExecutionState.FAILURE
        )
      ) {
        // If the message has already been executed, we skip it.  We want to not revert on race conditions between
        // executing parties. This will allow us to open up manual exec while also attempting with the DON, without
        // reverting an entire DON batch when a user manually executes while the tx is inflight.
        emit SkippedAlreadyExecutedMessage(message.sequenceNumber);
        continue;
      }
      uint32[] memory tokenGasOverrides;

      if (manualExecution) {
        tokenGasOverrides = manualExecGasOverrides[i].tokenGasOverrides;
        bool isOldCommitReport =
          (block.timestamp - timestampCommitted) > s_dynamicConfig.permissionLessExecutionThresholdSeconds;
        // Manually execution is fine if we previously failed or if the commit report is just too old
        // Acceptable state transitions: FAILURE->SUCCESS, UNTOUCHED->SUCCESS, FAILURE->FAILURE
        if (!(isOldCommitReport || originalState == Internal.MessageExecutionState.FAILURE)) {
          revert ManualExecutionNotYetEnabled();
        }

        // Manual execution gas limit can override gas limit specified in the message. Value of 0 indicates no override.
        if (manualExecGasOverrides[i].receiverExecutionGasLimit != 0) {
          message.gasLimit = manualExecGasOverrides[i].receiverExecutionGasLimit;
        }
      } else {
        // DON can only execute a message once
        // Acceptable state transitions: UNTOUCHED->SUCCESS, UNTOUCHED->FAILURE
        if (originalState != Internal.MessageExecutionState.UNTOUCHED) {
          emit AlreadyAttempted(message.sequenceNumber);
          continue;
        }
      }

      if (message.nonce != 0) {
        // In the scenario where we upgrade offRamps, we still want to have sequential nonces.
        // Referencing the old offRamp to check the expected nonce if none is set for a
        // given sender allows us to skip the current message if it would not be the next according
        // to the old offRamp. This preserves sequencing between updates.
        uint64 prevNonce = s_senderNonce[message.sender];
        if (prevNonce == 0) {
          if (i_prevOffRamp != address(0)) {
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
        }

        // UNTOUCHED messages MUST be executed in order always IF message.nonce > 0.
        if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
          if (prevNonce + 1 != message.nonce) {
            // We skip the message if the nonce is incorrect, since message.nonce > 0.
            emit SkippedIncorrectNonce(message.nonce, message.sender);
            continue;
          }
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
      (Internal.MessageExecutionState newState, bytes memory returnData) =
        _trialExecute(message, offchainTokenData, tokenGasOverrides);
      _setExecutionState(message.sequenceNumber, newState);

      // Since it's hard to estimate whether manual execution will succeed, we
      // revert the entire transaction if it fails. This will show the user if
      // their manual exec will fail before they submit it.
      if (manualExecution) {
        if (newState == Internal.MessageExecutionState.FAILURE) {
          if (originalState != Internal.MessageExecutionState.UNTOUCHED) {
            // If manual execution fails, we revert the entire transaction, unless the originalState is UNTOUCHED as we
            // would still be making progress by changing the state from UNTOUCHED to FAILURE.
            revert ExecutionError(returnData);
          }
        }
      }

      // The only valid prior states are UNTOUCHED and FAILURE (checked above)
      // The only valid post states are SUCCESS and FAILURE (checked below)
      if (newState != Internal.MessageExecutionState.SUCCESS) {
        if (newState != Internal.MessageExecutionState.FAILURE) {
          revert InvalidNewState(message.sequenceNumber, newState);
        }
      }

      // Nonce changes per state transition.
      // These only apply for ordered messages.
      // UNTOUCHED -> FAILURE  nonce bump
      // UNTOUCHED -> SUCCESS  nonce bump
      // FAILURE   -> FAILURE  no nonce bump
      // FAILURE   -> SUCCESS  no nonce bump
      if (message.nonce != 0) {
        if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
          s_senderNonce[message.sender]++;
        }
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
    if (numberOfTokens > uint256(s_dynamicConfig.maxNumberOfTokensPerMsg)) {
      revert UnsupportedNumberOfTokens(sequenceNumber);
    }
    if (numberOfTokens != offchainTokenDataLength) revert TokenDataMismatch(sequenceNumber);
    if (dataLength > uint256(s_dynamicConfig.maxDataBytes)) {
      revert MessageTooLarge(uint256(s_dynamicConfig.maxDataBytes), dataLength);
    }
  }

  /// @notice Try executing a message.
  /// @param message Internal.EVM2EVMMessage memory message.
  /// @param offchainTokenData Data provided by the DON for token transfers.
  /// @return the new state of the message, being either SUCCESS or FAILURE.
  /// @return revert data in bytes if CCIP receiver reverted during execution.
  function _trialExecute(
    Internal.EVM2EVMMessage memory message,
    bytes[] memory offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) internal returns (Internal.MessageExecutionState, bytes memory) {
    try this.executeSingleMessage(message, offchainTokenData, tokenGasOverrides) {}
    catch (bytes memory err) {
      // return the message execution state as FAILURE and the revert data
      // Max length of revert data is Router.MAX_RET_BYTES, max length of err is 4 + Router.MAX_RET_BYTES
      return (Internal.MessageExecutionState.FAILURE, err);
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
  function executeSingleMessage(
    Internal.EVM2EVMMessage calldata message,
    bytes[] calldata offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) external {
    if (msg.sender != address(this)) revert CanOnlySelfCall();
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](0);
    if (message.tokenAmounts.length > 0) {
      destTokenAmounts = _releaseOrMintTokens(
        message.tokenAmounts,
        abi.encode(message.sender),
        message.receiver,
        message.sourceTokenData,
        offchainTokenData,
        tokenGasOverrides
      );
    }
    // There are three cases in which we skip calling the receiver:
    // 1. If the message data is empty AND the gas limit is 0.
    //          This indicates a message that only transfers tokens. It is valid to only send tokens to a contract
    //          that supports the IAny2EVMMessageReceiver interface, but without this first check we would call the
    //          receiver without any gas, which would revert the transaction.
    // 2. If the receiver is not a contract.
    // 3. If the receiver is a contract but it does not support the IAny2EVMMessageReceiver interface.
    //
    // The ordering of these checks is important, as the first check is the cheapest to execute.
    if (
      (message.data.length == 0 && message.gasLimit == 0) || message.receiver.code.length == 0
        || !message.receiver.supportsInterface(type(IAny2EVMMessageReceiver).interfaceId)
    ) return;

    (bool success, bytes memory returnData,) = IRouter(s_dynamicConfig.router).routeMessage(
      Client.Any2EVMMessage({
        messageId: message.messageId,
        sourceChainSelector: message.sourceChainSelector,
        sender: abi.encode(message.sender),
        data: message.data,
        destTokenAmounts: destTokenAmounts
      }),
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
    return StaticConfig({
      commitStore: i_commitStore,
      chainSelector: i_chainSelector,
      sourceChainSelector: i_sourceChainSelector,
      onRamp: i_onRamp,
      prevOffRamp: i_prevOffRamp,
      rmnProxy: i_rmnProxy,
      tokenAdminRegistry: i_tokenAdminRegistry
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
        rmnProxy: i_rmnProxy,
        tokenAdminRegistry: i_tokenAdminRegistry
      }),
      dynamicConfig
    );
  }

  /// @notice Get all tokens which are included in Aggregate Rate Limiting.
  /// @return sourceTokens The source representation of the tokens that are rate limited.
  /// @return destTokens The destination representation of the tokens that are rate limited.
  /// @dev the order of IDs in the list is **not guaranteed**, therefore, if ordering matters when
  /// making successive calls, one should keep the block height constant to ensure a consistent result.
  function getAllRateLimitTokens() external view returns (address[] memory sourceTokens, address[] memory destTokens) {
    uint256 numRateLimitedTokens = s_rateLimitedTokensDestToSource.length();
    sourceTokens = new address[](numRateLimitedTokens);
    destTokens = new address[](numRateLimitedTokens);

    for (uint256 i = 0; i < numRateLimitedTokens; ++i) {
      (address destToken, address sourceToken) = s_rateLimitedTokensDestToSource.at(i);
      sourceTokens[i] = sourceToken;
      destTokens[i] = destToken;
    }
    return (sourceTokens, destTokens);
  }

  /// @notice Adds or removes tokens from being used in Aggregate Rate Limiting.
  /// @param removes - A list of one or more tokens to be removed.
  /// @param adds - A list of one or more tokens to be added.
  function updateRateLimitTokens(RateLimitToken[] memory removes, RateLimitToken[] memory adds) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      if (s_rateLimitedTokensDestToSource.remove(removes[i].destToken)) {
        emit TokenAggregateRateLimitRemoved(removes[i].sourceToken, removes[i].destToken);
      }
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      if (s_rateLimitedTokensDestToSource.set(adds[i].destToken, adds[i].sourceToken)) {
        emit TokenAggregateRateLimitAdded(adds[i].sourceToken, adds[i].destToken);
      }
    }
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @notice Uses a pool to release or mint a token to a receiver address, with balance checks before and after the
  /// transfer. This is done to ensure the exact number of tokens the pool claims to release are actually transferred.
  /// @dev The local token address is validated through the TokenAdminRegistry. If, due to some misconfiguration, the
  /// token is unknown to the registry, the offRamp will revert. The tx, and the tokens, can be retrieved by
  /// registering the token on this chain, and re-trying the msg.
  /// @param sourceAmount The amount of tokens to be released/minted.
  /// @param originalSender The message sender on the source chain.
  /// @param receiver The address that will receive the tokens.
  /// @param sourceTokenData A struct containing the local token address, the source pool address and optional data
  /// returned from the source pool.
  /// @param offchainTokenData Data fetched offchain by the DON.
  function _releaseOrMintToken(
    uint256 sourceAmount,
    bytes memory originalSender,
    address receiver,
    Internal.SourceTokenData memory sourceTokenData,
    bytes memory offchainTokenData
  ) internal returns (Client.EVMTokenAmount memory destTokenAmount) {
    // We need to safely decode the token address from the sourceTokenData, as it could be wrong,
    // in which case it doesn't have to be a valid EVM address.
    address localToken = Internal._validateEVMAddress(sourceTokenData.destTokenAddress);
    // We check with the token admin registry if the token has a pool on this chain.
    address localPoolAddress = ITokenAdminRegistry(i_tokenAdminRegistry).getPool(localToken);
    // This will call the supportsInterface through the ERC165Checker, and not directly on the pool address.
    // This is done to prevent a pool from reverting the entire transaction if it doesn't support the interface.
    // The call gets a max or 30k gas per instance, of which there are three. This means gas estimations should
    // account for 90k gas overhead due to the interface check.
    if (localPoolAddress == address(0) || !localPoolAddress.supportsInterface(Pool.CCIP_POOL_V1)) {
      revert NotACompatiblePool(localPoolAddress);
    }

    // We retrieve the local token balance of the receiver before the pool call.
    (uint256 balancePre, uint256 gasLeft) = _getBalanceOfReceiver(receiver, localToken, sourceTokenData.destGasAmount);

    // We determined that the pool address is a valid EVM address, but that does not mean the code at this
    // address is a (compatible) pool contract. _callWithExactGasSafeReturnData will check if the location
    // contains a contract. If it doesn't it reverts with a known error, which we catch gracefully.
    // We call the pool with exact gas to increase resistance against malicious tokens or token pools.
    // We protects against return data bombs by capping the return data size at MAX_RET_BYTES.
    (bool success, bytes memory returnData, uint256 gasUsedReleaseOrMint) = CallWithExactGas
      ._callWithExactGasSafeReturnData(
      abi.encodeCall(
        IPoolV1.releaseOrMint,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: receiver,
          amount: sourceAmount,
          localToken: localToken,
          remoteChainSelector: i_sourceChainSelector,
          sourcePoolAddress: sourceTokenData.sourcePoolAddress,
          sourcePoolData: sourceTokenData.extraData,
          offchainTokenData: offchainTokenData
        })
      ),
      localPoolAddress,
      gasLeft,
      Internal.GAS_FOR_CALL_EXACT_CHECK,
      Internal.MAX_RET_BYTES
    );

    // wrap and rethrow the error so we can catch it lower in the stack
    if (!success) revert TokenHandlingError(returnData);
    // If the call was successful, the returnData should contain only the local token amount.
    if (returnData.length != Pool.CCIP_POOL_V1_RET_BYTES) {
      revert InvalidDataLength(Pool.CCIP_POOL_V1_RET_BYTES, returnData.length);
    }

    uint256 localAmount = abi.decode(returnData, (uint256));
    // We don't need to do balance checks if the pool is the receiver, as they would always fail in the case
    // of a lockRelease pool.
    if (receiver != localPoolAddress) {
      (uint256 balancePost,) = _getBalanceOfReceiver(receiver, localToken, gasLeft - gasUsedReleaseOrMint);

      // First we check if the subtraction would result in an underflow to ensure we revert with a clear error
      if (balancePost < balancePre || balancePost - balancePre != localAmount) {
        revert ReleaseOrMintBalanceMismatch(localAmount, balancePre, balancePost);
      }
    }

    return Client.EVMTokenAmount({token: localToken, amount: localAmount});
  }

  function _getBalanceOfReceiver(
    address receiver,
    address token,
    uint256 gasLimit
  ) internal returns (uint256 balance, uint256 gasLeft) {
    (bool success, bytes memory returnData, uint256 gasUsed) = CallWithExactGas._callWithExactGasSafeReturnData(
      abi.encodeCall(IERC20.balanceOf, (receiver)),
      token,
      gasLimit,
      Internal.GAS_FOR_CALL_EXACT_CHECK,
      Internal.MAX_RET_BYTES
    );
    if (!success) revert TokenHandlingError(returnData);

    // If the call was successful, the returnData should contain only the balance.
    if (returnData.length != Internal.MAX_BALANCE_OF_RET_BYTES) {
      revert InvalidDataLength(Internal.MAX_BALANCE_OF_RET_BYTES, returnData.length);
    }

    // Return the decoded balance, which cannot fail as we checked the length, and the gas that is left
    // after this call.
    return (abi.decode(returnData, (uint256)), gasLimit - gasUsed);
  }

  /// @notice Uses pools to release or mint a number of different tokens to a receiver address.
  /// @param sourceTokenAmounts List of tokens and amount values to be released/minted.
  /// @param originalSender The message sender.
  /// @param receiver The address that will receive the tokens.
  /// @param encodedSourceTokenData Array of token data returned by token pools on the source chain.
  /// @param offchainTokenData Array of token data fetched offchain by the DON.
  /// @dev This function wrappes the token pool call in a try catch block to gracefully handle
  /// any non-rate limiting errors that may occur. If we encounter a rate limiting related error
  /// we bubble it up. If we encounter a non-rate limiting error we wrap it in a TokenHandlingError.
  function _releaseOrMintTokens(
    Client.EVMTokenAmount[] calldata sourceTokenAmounts,
    bytes memory originalSender,
    address receiver,
    bytes[] calldata encodedSourceTokenData,
    bytes[] calldata offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) internal returns (Client.EVMTokenAmount[] memory destTokenAmounts) {
    // Creating a copy is more gas efficient than initializing a new array.
    destTokenAmounts = sourceTokenAmounts;
    uint256 value = 0;
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      Internal.SourceTokenData memory sourceTokenData =
        abi.decode(encodedSourceTokenData[i], (Internal.SourceTokenData));
      if (tokenGasOverrides.length != 0) {
        if (tokenGasOverrides[i] != 0) {
          sourceTokenData.destGasAmount = tokenGasOverrides[i];
        }
      }
      destTokenAmounts[i] = _releaseOrMintToken(
        sourceTokenAmounts[i].amount,
        originalSender,
        receiver,
        // This should never revert as the onRamp encodes the sourceTokenData struct. Only the inner components from
        // this struct come from untrusted sources.
        sourceTokenData,
        offchainTokenData[i]
      );

      if (s_rateLimitedTokensDestToSource.contains(destTokenAmounts[i].token)) {
        value += _getTokenValue(destTokenAmounts[i], IPriceRegistry(s_dynamicConfig.priceRegistry));
      }
    }

    if (value > 0) _rateLimitValue(value);

    return destTokenAmounts;
  }

  // ================================================================
  // │                            Access                            │
  // ================================================================

  /// @notice Reverts as this contract should not access CCIP messages
  function ccipReceive(Client.Any2EVMMessage calldata) external pure {
    // solhint-disable-next-line
    revert();
  }
}
