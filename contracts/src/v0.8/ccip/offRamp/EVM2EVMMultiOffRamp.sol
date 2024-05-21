// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "../interfaces/IAny2EVMMessageReceiver.sol";

import {IAny2EVMMultiOffRamp} from "../interfaces/IAny2EVMMultiOffRamp.sol";
import {IAny2EVMOffRamp} from "../interfaces/IAny2EVMOffRamp.sol";
import {ICommitStore} from "../interfaces/ICommitStore.sol";
import {IPool} from "../interfaces/IPool.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRouter} from "../interfaces/IRouter.sol";

import {CallWithExactGas} from "../../shared/call/CallWithExactGas.sol";
import {EnumerableMapAddresses} from "../../shared/enumerable/EnumerableMapAddresses.sol";
import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {OCR2BaseNoChecks} from "../ocr/OCR2BaseNoChecks.sol";

import {ERC165Checker} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165Checker.sol";

/// @notice EVM2EVMOffRamp enables OCR networks to execute multiple messages
/// in an OffRamp in a single transaction.
/// @dev The EVM2EVMOnRamp, CommitStore and EVM2EVMOffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
/// @dev OCR2BaseNoChecks is used to save gas, signatures are not required as the offramp can only execute
/// messages which are committed in the commitStore. We still make use of OCR2 as an executor whitelist
/// and turn-taking mechanism.
contract EVM2EVMMultiOffRamp is IAny2EVMMultiOffRamp, AggregateRateLimiter, ITypeAndVersion, OCR2BaseNoChecks {
  using ERC165Checker for address;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  error AlreadyAttempted(uint64 sourceChainSelector, uint64 sequenceNumber);
  error AlreadyExecuted(uint64 sourceChainSelector, uint64 sequenceNumber);
  error ZeroAddressNotAllowed();
  error CommitStoreAlreadyInUse(uint64 sourceChainSelector);
  error ExecutionError(bytes32 messageId, bytes error);
  error SourceChainNotEnabled(uint64 sourceChainSelector);
  error MessageTooLarge(bytes32 messageId, uint256 maxSize, uint256 actualSize);
  error TokenDataMismatch(uint64 sourceChainSelector, uint64 sequenceNumber);
  error UnexpectedTokenData();
  error UnsupportedNumberOfTokens(uint64 sourceChainSelector, uint64 sequenceNumber);
  error ManualExecutionNotYetEnabled(uint64 sourceChainSelector);
  error ManualExecutionGasLimitMismatch();
  error InvalidManualExecutionGasLimit(uint64 sourceChainSelector, uint256 index, uint256 newLimit);
  error RootNotCommitted(uint64 sourceChainSelector);
  error CanOnlySelfCall();
  error ReceiverError(bytes error);
  error TokenHandlingError(bytes error);
  error EmptyReport();
  error CursedByRMN();
  error InvalidMessageId(bytes32 messageId);
  error NotACompatiblePool(address notPool);
  error InvalidDataLength(uint256 expected, uint256 got);
  error InvalidNewState(uint64 sourceChainSelector, uint64 sequenceNumber, Internal.MessageExecutionState newState);
  error IndexOutOfRange();
  error StaticConfigCannotBeUpdated();

  /// @dev Atlas depends on this event, if changing, please notify Atlas.
  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  // TODO: revisit if fields have to be indexed for skip events
  event SkippedIncorrectNonce(uint64 sourceChainSelector, uint64 nonce, address indexed sender);
  event SkippedSenderWithPreviousRampMessageInflight(
    uint64 indexed sourceChainSelector, uint64 nonce, address indexed sender
  );
  /// @dev RMN depends on this event, if changing, please notify the RMN maintainers.
  event ExecutionStateChanged(
    uint64 indexed sourceChainSelector,
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );
  event TokenAggregateRateLimitAdded(address sourceToken, address destToken);
  event TokenAggregateRateLimitRemoved(address sourceToken, address destToken);
  event SourceChainSelectorAdded(uint64 sourceChainSelector);
  event SourceChainConfigSet(uint64 indexed sourceChainSelector, SourceChainConfig sourceConfig);
  // TODO: index with source chain selector
  event SkippedAlreadyExecutedMessage(uint64 indexed sequenceNumber);

  /// @notice Static offRamp config
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct StaticConfig {
    address commitStore; // ────────╮  CommitStore address on the destination chain
    uint64 chainSelector; // ───────╯  Destination chainSelector
    address rmnProxy; //               RMN proxy address
  }

  /// @notice Per-chain source config (defining a lane from a Source Chain -> Dest OffRamp)
  struct SourceChainConfig {
    // TODO: re-evaluate on removing this (can be controlled by CommitStore)
    // TODO: if used - pack together with onRamp to localise storage slot reads
    bool isEnabled; // ─────────╮  Flag whether the source chain is enabled or not
    address prevOffRamp; // ────╯  Address of previous-version per-lane OffRamp. Used to be able to provide seequencing continuity during a zero downtime upgrade.
    address onRamp; //             OnRamp address on the source chain
    /// @dev Ensures that 2 identical messages sent to 2 different lanes will have a distinct hash.
    /// Must match the metadataHash used in computing leaf hashes offchain for the root committed in
    /// the commitStore and i_metadataHash in the onRamp.
    bytes32 metadataHash; //      Source-chain specific message hash preimage to ensure global uniqueness
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

  /// @notice RateLimitToken struct containing both the source and destination token addresses
  struct RateLimitToken {
    address sourceToken;
    address destToken;
  }

  /// @notice Struct that represents a message route (sender -> receiver and source chain)
  struct Any2EVMMessageRoute {
    bytes sender; //                    Message sender
    uint64 sourceChainSelector; // ───╮ Source chain that the message originates from
    address receiver; // ─────────────╯ Address that receives the message
  }

  /// @notice SourceChainConfig update args scoped to one source chain
  struct SourceChainConfigArgs {
    uint64 sourceChainSelector; //  ───╮  Source chain selector of the config to update
    bool isEnabled; //                 │  Flag whether the source chain is enabled or not
    address prevOffRamp; // ───────────╯  Address of previous-version per-lane OffRamp. Used to be able to provide seequencing continuity during a zero downtime upgrade.
    address onRamp; //                    OnRamp address on the source chain
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMMultiOffRamp 1.6.0-dev";
  /// @dev Commit store address on the destination chain
  address internal immutable i_commitStore;
  /// @dev ChainSelector of this chain
  uint64 internal immutable i_chainSelector;
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;

  // DYNAMIC CONFIG
  DynamicConfig internal s_dynamicConfig;
  /// @dev Tokens that should be included in Aggregate Rate Limiting
  /// An (address => address) map is used for backwards compatability of offchain code
  EnumerableMapAddresses.AddressToAddressMap internal s_rateLimitedTokensDestToSource;

  // TODO: evaluate whether this should be pulled in (since this can be inferred from SourceChainSelectorAdded events instead)
  /// @notice all source chains available in s_sourceChainConfigs
  // uint64[] internal s_sourceChainSelectors;

  /// @notice SourceConfig per chain
  /// (forms lane configurations from sourceChainSelector => StaticConfig.chainSelector)
  mapping(uint64 sourceChainSelector => SourceChainConfig) internal s_sourceChainConfigs;

  // STATE
  /// @dev The expected nonce for a given sender per source chain.
  /// Corresponds to s_senderNonce in the OnRamp for a lane, used to enforce that messages are
  /// executed in the same order they are sent (assuming they are DON). Note that re-execution
  /// of FAILED messages however, can be out of order.
  mapping(uint64 sourceChainSelector => mapping(address sender => uint64 nonce)) internal s_senderNonce;
  /// @dev A mapping of sequence numbers (per source chain) to execution state using a bitmap with each execution
  /// state only taking up 2 bits of the uint256, packing 128 states into a single slot.
  /// Message state is tracked to ensure message can only be executed successfully once.
  mapping(uint64 sourceChainSelector => mapping(uint64 seqNum => uint256 executionStateBitmap)) internal
    s_executionStates;

  constructor(
    StaticConfig memory staticConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs,
    // TODO: convert to array to support per-chain config once multi-ARL is ready
    RateLimiter.Config memory rateLimiterConfig
  ) OCR2BaseNoChecks() AggregateRateLimiter(rateLimiterConfig) {
    if (staticConfig.commitStore == address(0)) revert ZeroAddressNotAllowed();

    i_commitStore = staticConfig.commitStore;
    i_chainSelector = staticConfig.chainSelector;
    i_rmnProxy = staticConfig.rmnProxy;

    _applySourceChainConfigUpdates(sourceChainConfigs);
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  // The size of the execution state in bits
  uint256 private constant MESSAGE_EXECUTION_STATE_BIT_WIDTH = 2;
  // The mask for the execution state bits
  uint256 private constant MESSAGE_EXECUTION_STATE_MASK = (1 << MESSAGE_EXECUTION_STATE_BIT_WIDTH) - 1;

  /// @notice Returns the current execution state of a message based on its sequenceNumber.
  /// @param sourceChainSelector The source chain to get the execution state for
  /// @param sequenceNumber The sequence number of the message to get the execution state for.
  /// @return The current execution state of the message.
  /// @dev we use the literal number 128 because using a constant increased gas usage.
  function getExecutionState(
    uint64 sourceChainSelector,
    uint64 sequenceNumber
  ) public view returns (Internal.MessageExecutionState) {
    return Internal.MessageExecutionState(
      (
        s_executionStates[sourceChainSelector][sequenceNumber / 128]
          >> ((sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH)
      ) & MESSAGE_EXECUTION_STATE_MASK
    );
  }

  /// @notice Sets a new execution state for a given sequence number. It will overwrite any existing state.
  /// @param sourceChainSelector The source chain to set the execution state for
  /// @param sequenceNumber The sequence number for which the state will be saved.
  /// @param newState The new value the state will be in after this function is called.
  /// @dev we use the literal number 128 because using a constant increased gas usage.
  function _setExecutionState(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    Internal.MessageExecutionState newState
  ) internal {
    uint256 offset = (sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH;
    uint256 bitmap = s_executionStates[sourceChainSelector][sequenceNumber / 128];
    // to unset any potential existing state we zero the bits of the section the state occupies,
    // then we do an AND operation to blank out any existing state for the section.
    bitmap &= ~(MESSAGE_EXECUTION_STATE_MASK << offset);
    // Set the new state
    bitmap |= uint256(newState) << offset;

    s_executionStates[sourceChainSelector][sequenceNumber / 128] = bitmap;
  }

  /// @inheritdoc IAny2EVMMultiOffRamp
  function getSenderNonce(uint64 sourceChainSelector, address sender) external view returns (uint64) {
    (uint64 nonce,) = _getSenderNonce(sourceChainSelector, sender);
    return nonce;
  }

  /// @notice Returns the the current nonce for a receiver.
  /// @param sourceChainSelector The source chain to retrieve the nonce for
  /// @param sender The sender address
  /// @return nonce The nonce value belonging to the sender address.
  /// @return isFromPrevRamp True if the nonce was retrieved from the prevOffRamps
  function _getSenderNonce(
    uint64 sourceChainSelector,
    address sender
  ) internal view returns (uint64 nonce, bool isFromPrevRamp) {
    uint64 senderNonce = s_senderNonce[sourceChainSelector][sender];

    if (senderNonce == 0) {
      address prevOffRamp = s_sourceChainConfigs[sourceChainSelector].prevOffRamp;
      if (prevOffRamp != address(0)) {
        // If OffRamp was upgraded, check if sender has a nonce from the previous OffRamp.
        // NOTE: assuming prevOffRamp is always a lane-specific off ramp
        // TODO: on deployment - revisit if this assumption holds
        return (IAny2EVMOffRamp(prevOffRamp).getSenderNonce(sender), true);
      }
    }

    return (senderNonce, false);
  }

  /// @notice Manually executes a set of reports.
  /// @param reports Internal.ExecutionReportSingleChain[] - list of reports to execute
  /// @param gasLimitOverrides New gasLimit for each message per report
  //         The outer array represents each report, inner array represents each message in the report.
  //         i.e. gasLimitOverrides[report1][report1Message1] -> access message1 from report1
  /// @dev We permit gas limit overrides so that users may manually execute messages which failed due to
  /// insufficient gas provided.
  /// The reports do not have to contain all the messages (they can be omitted). Multiple reports can be passed in simultaneously.
  function manuallyExecute(
    Internal.ExecutionReportSingleChain[] memory reports,
    uint256[][] memory gasLimitOverrides
  ) external {
    // We do this here because the other _execute path is already covered OCR2BaseXXX.
    if (i_chainID != block.chainid) revert OCR2BaseNoChecks.ForkedChain(i_chainID, uint64(block.chainid));

    uint256 numReports = reports.length;
    if (numReports != gasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();

    for (uint256 reportIndex = 0; reportIndex < numReports; ++reportIndex) {
      Internal.ExecutionReportSingleChain memory report = reports[reportIndex];

      uint256 numMsgs = report.messages.length;
      uint256[] memory msgGasLimitOverrides = gasLimitOverrides[reportIndex];
      if (numMsgs != msgGasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();

      for (uint256 msgIndex = 0; msgIndex < numMsgs; ++msgIndex) {
        uint256 newLimit = msgGasLimitOverrides[msgIndex];
        // Checks to ensure message cannot be executed with less gas than specified.
        if (newLimit != 0 && newLimit < report.messages[msgIndex].gasLimit) {
          revert InvalidManualExecutionGasLimit(report.sourceChainSelector, msgIndex, newLimit);
        }
      }
    }

    _batchExecute(reports, gasLimitOverrides);
  }

  /// @notice Entrypoint for execution, called by the OCR network
  /// @dev Expects an encoded ExecutionReport
  function _report(bytes calldata report) internal override {
    Internal.ExecutionReportSingleChain[] memory reports = abi.decode(report, (Internal.ExecutionReportSingleChain[]));

    _batchExecute(reports, new uint256[][](0));
  }

  /// @notice Batch executes a set of reports, each report matching one single source chain
  /// @param reports Set of execution reports (one per chain) containing the messages and proofs
  /// @param manualExecGasLimits An array of gas limits to use for manual execution
  //         The outer array represents each report, inner array represents each message in the report.
  //         i.e. gasLimitOverrides[report1][report1Message1] -> access message1 from report1
  /// @dev The manualExecGasLimits array should either be empty, or match the length of the reports array
  /// @dev If called from manual execution, each inner array's length has to match the number of messages.
  function _batchExecute(
    Internal.ExecutionReportSingleChain[] memory reports,
    uint256[][] memory manualExecGasLimits
  ) internal {
    if (reports.length == 0) revert EmptyReport();

    bool areManualGasLimitsEmpty = manualExecGasLimits.length == 0;
    // Cache array for gas savings in the loop's condition
    uint256[] memory emptyGasLimits = new uint256[](0);

    for (uint256 i = 0; i < reports.length; ++i) {
      _execute(reports[i], areManualGasLimitsEmpty ? emptyGasLimits : manualExecGasLimits[i]);
    }
  }

  /// @notice Executes a report, executing each message in order.
  /// @param report The execution report containing the messages and proofs.
  /// @param manualExecGasLimits An array of gas limits to use for manual execution.
  /// @dev If called from the DON, this array is always empty.
  /// @dev If called from manual execution, this array is always same length as messages.
  function _execute(Internal.ExecutionReportSingleChain memory report, uint256[] memory manualExecGasLimits) internal {
    // TODO pass in source chain selector to check for cursed source chain
    if (IRMN(i_rmnProxy).isCursed()) revert CursedByRMN();

    uint256 numMsgs = report.messages.length;
    if (numMsgs == 0) revert EmptyReport();
    if (numMsgs != report.offchainTokenData.length) revert UnexpectedTokenData();

    uint64 sourceChainSelector = report.sourceChainSelector;
    SourceChainConfig storage sourceChainConfig = s_sourceChainConfigs[sourceChainSelector];
    if (!sourceChainConfig.isEnabled) {
      revert SourceChainNotEnabled(sourceChainSelector);
    }

    bytes32[] memory hashedLeaves = new bytes32[](numMsgs);

    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      // We do this hash here instead of in _verifyMessages to avoid two separate loops
      // over the same data, which increases gas cost
      hashedLeaves[i] = Internal._hash(message, sourceChainConfig.metadataHash);
      // For EVM2EVM offramps, the messageID is the leaf hash.
      // Asserting that this is true ensures we don't accidentally commit and then execute
      // a message with an unexpected hash.
      if (hashedLeaves[i] != message.messageId) revert InvalidMessageId(message.messageId);
    }

    // SECURITY CRITICAL CHECK
    // NOTE: This check also verifies that all messages match the report's sourceChainSelector
    // TODO: revisit after MultiCommitStore implementation
    uint256 timestampCommitted = ICommitStore(i_commitStore).verify(hashedLeaves, report.proofs, report.proofFlagBits);
    if (timestampCommitted == 0) revert RootNotCommitted(sourceChainSelector);

    // Execute messages
    bool manualExecution = manualExecGasLimits.length != 0;
    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      uint64 sequenceNumber = message.sequenceNumber;

      Internal.MessageExecutionState originalState = getExecutionState(sourceChainSelector, sequenceNumber);
      if (originalState == Internal.MessageExecutionState.SUCCESS) {
        // If the message has already been executed, we skip it.  We want to not revert on race conditions between
        // executing parties. This will allow us to open up manual exec while also attempting with the DON, without
        // reverting an entire DON batch when a user manually executes while the tx is inflight.
        emit SkippedAlreadyExecutedMessage(message.sequenceNumber);
        continue;
      }
      // Two valid cases here, we either have never touched this message before, or we tried to execute
      // and failed. This check protects against reentry and re-execution because the other state is
      // IN_PROGRESS which should not be allowed to execute.
      if (
        !(
          originalState == Internal.MessageExecutionState.UNTOUCHED
            || originalState == Internal.MessageExecutionState.FAILURE
        )
      ) revert AlreadyExecuted(sourceChainSelector, sequenceNumber);

      if (manualExecution) {
        bool isOldCommitReport =
          (block.timestamp - timestampCommitted) > s_dynamicConfig.permissionLessExecutionThresholdSeconds;
        // Manually execution is fine if we previously failed or if the commit report is just too old
        // Acceptable state transitions: FAILURE->SUCCESS, UNTOUCHED->SUCCESS, FAILURE->FAILURE
        if (!(isOldCommitReport || originalState == Internal.MessageExecutionState.FAILURE)) {
          revert ManualExecutionNotYetEnabled(sourceChainSelector);
        }

        // Manual execution gas limit can override gas limit specified in the message. Value of 0 indicates no override.
        if (manualExecGasLimits[i] != 0) {
          message.gasLimit = manualExecGasLimits[i];
        }
      } else {
        // DON can only execute a message once
        // Acceptable state transitions: UNTOUCHED->SUCCESS, UNTOUCHED->FAILURE
        if (originalState != Internal.MessageExecutionState.UNTOUCHED) {
          revert AlreadyAttempted(sourceChainSelector, sequenceNumber);
        }
      }

      // In the scenario where we upgrade offRamps, we still want to have sequential nonces.
      // Referencing the old offRamp to check the expected nonce if none is set for a
      // given sender allows us to skip the current message if it would not be the next according
      // to the old offRamp. This preserves sequencing between updates.
      (uint64 prevNonce, bool isFromPrevRamp) = _getSenderNonce(sourceChainSelector, message.sender);
      if (isFromPrevRamp) {
        if (prevNonce + 1 != message.nonce) {
          // the starting v2 onramp nonce, i.e. the 1st message nonce v2 offramp is expected to receive,
          // is guaranteed to equal (largest v1 onramp nonce + 1).
          // if this message's nonce isn't (v1 offramp nonce + 1), then v1 offramp nonce != largest v1 onramp nonce,
          // it tells us there are still messages inflight for v1 offramp
          emit SkippedSenderWithPreviousRampMessageInflight(sourceChainSelector, message.nonce, message.sender);
          continue;
        }
        // Otherwise this nonce is indeed the "transitional nonce", that is
        // all messages sent to v1 ramp have been executed by the DON and the sequence can resume in V2.
        // Note if first time user in V2, then prevNonce will be 0, and message.nonce = 1, so this will be a no-op.
        s_senderNonce[sourceChainSelector][message.sender] = prevNonce;
      }

      // UNTOUCHED messages MUST be executed in order always
      if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
        if (prevNonce + 1 != message.nonce) {
          // We skip the message if the nonce is incorrect
          emit SkippedIncorrectNonce(sourceChainSelector, message.nonce, message.sender);
          continue;
        }
      }

      // Although we expect only valid messages will be committed, we check again
      // when executing as a defense in depth measure.
      bytes[] memory offchainTokenData = report.offchainTokenData[i];
      _isWellFormed(
        message.messageId,
        sourceChainSelector,
        sequenceNumber,
        message.tokenAmounts.length,
        message.data.length,
        offchainTokenData.length
      );

      _setExecutionState(sourceChainSelector, sequenceNumber, Internal.MessageExecutionState.IN_PROGRESS);
      (Internal.MessageExecutionState newState, bytes memory returnData) = _trialExecute(message, offchainTokenData);
      _setExecutionState(sourceChainSelector, sequenceNumber, newState);

      // Since it's hard to estimate whether manual execution will succeed, we
      // revert the entire transaction if it fails. This will show the user if
      // their manual exec will fail before they submit it.
      if (manualExecution && newState == Internal.MessageExecutionState.FAILURE) {
        // If manual execution fails, we revert the entire transaction.
        revert ExecutionError(message.messageId, returnData);
      }

      // The only valid prior states are UNTOUCHED and FAILURE (checked above)
      // The only valid post states are FAILURE and SUCCESS (checked below)
      if (newState != Internal.MessageExecutionState.FAILURE && newState != Internal.MessageExecutionState.SUCCESS) {
        revert InvalidNewState(sourceChainSelector, sequenceNumber, newState);
      }

      // Nonce changes per state transition
      // UNTOUCHED -> FAILURE  nonce bump
      // UNTOUCHED -> SUCCESS  nonce bump
      // FAILURE   -> FAILURE  no nonce bump
      // FAILURE   -> SUCCESS  no nonce bump
      if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
        s_senderNonce[sourceChainSelector][message.sender]++;
      }

      emit ExecutionStateChanged(sourceChainSelector, sequenceNumber, message.messageId, newState, returnData);
    }
  }

  /// @notice Does basic message validation. Should never fail.
  /// @param sequenceNumber Sequence number of the message.
  /// @param numberOfTokens Length of tokenAmounts array in the message.
  /// @param dataLength Length of data field in the message.
  /// @param offchainTokenDataLength Length of offchainTokenData array.
  /// @dev reverts on validation failures.
  function _isWellFormed(
    bytes32 messageId,
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    uint256 numberOfTokens,
    uint256 dataLength,
    uint256 offchainTokenDataLength
  ) private view {
    if (numberOfTokens > uint256(s_dynamicConfig.maxNumberOfTokensPerMsg)) {
      revert UnsupportedNumberOfTokens(sourceChainSelector, sequenceNumber);
    }
    if (numberOfTokens != offchainTokenDataLength) revert TokenDataMismatch(sourceChainSelector, sequenceNumber);
    if (dataLength > uint256(s_dynamicConfig.maxDataBytes)) {
      revert MessageTooLarge(messageId, uint256(s_dynamicConfig.maxDataBytes), dataLength);
    }
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
    try this.executeSingleMessage(message, offchainTokenData) {}
    catch (bytes memory err) {
      if (
        ReceiverError.selector == bytes4(err) || TokenHandlingError.selector == bytes4(err)
          || Internal.InvalidEVMAddress.selector == bytes4(err) || InvalidDataLength.selector == bytes4(err)
          || CallWithExactGas.NoContract.selector == bytes4(err) || NotACompatiblePool.selector == bytes4(err)
      ) {
        // If CCIP receiver execution is not successful, bubble up receiver revert data,
        // prepended by the 4 bytes of ReceiverError.selector, TokenHandlingError.selector or InvalidPoolAddress.selector.
        // Max length of revert data is Router.MAX_RET_BYTES, max length of err is 4 + Router.MAX_RET_BYTES
        return (Internal.MessageExecutionState.FAILURE, err);
      } else {
        // If revert is not caused by CCIP receiver, it is unexpected, bubble up the revert.
        revert ExecutionError(message.messageId, err);
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
        Any2EVMMessageRoute({
          sender: abi.encode(message.sender),
          sourceChainSelector: message.sourceChainSelector,
          receiver: message.receiver
        }),
        message.sourceTokenData,
        offchainTokenData
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
      Internal._toAny2EVMMessage(message, destTokenAmounts),
      Internal.GAS_FOR_CALL_EXACT_CHECK,
      message.gasLimit,
      message.receiver
    );
    // If CCIP receiver execution is not successful, revert the call including token transfers
    if (!success) revert ReceiverError(returnData);
  }

  /// @notice creates a unique hash to be used in message hashing.
  function _metadataHash(uint64 sourceChainSelector, address onRamp, bytes32 prefix) internal view returns (bytes32) {
    return keccak256(abi.encode(prefix, sourceChainSelector, i_chainSelector, onRamp));
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static config.
  /// @dev This function will always return the same struct as the contents is static and can never change.
  /// RMN depends on this function, if changing, please notify the RMN maintainers.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({commitStore: i_commitStore, chainSelector: i_chainSelector, rmnProxy: i_rmnProxy});
  }

  /// @notice Returns the current dynamic config.
  /// @return The current config.
  function getDynamicConfig() external view returns (DynamicConfig memory) {
    return s_dynamicConfig;
  }

  /// @notice Returns the source chain config for the provided source chain selector
  /// @param sourceChainSelector chain to retrieve configuration for
  /// @return SourceChainConfig config for the source chain
  function getSourceChainConfig(uint64 sourceChainSelector) external view returns (SourceChainConfig memory) {
    return s_sourceChainConfigs[sourceChainSelector];
  }

  /// @notice Returns all configured source chain selectors
  /// @return sourceChainSelectors source chain selectors
  // function getSourceChainSelectors() external view returns (uint64[] memory) {
  //   return s_sourceChainSelectors;
  // }

  /// @notice Updates source configs
  /// @param sourceChainConfigUpdates Source chain configs
  function applySourceChainConfigUpdates(SourceChainConfigArgs[] memory sourceChainConfigUpdates) external onlyOwner {
    _applySourceChainConfigUpdates(sourceChainConfigUpdates);
  }

  /// @notice Updates source configs
  /// @param sourceChainConfigUpdates Source chain configs
  function _applySourceChainConfigUpdates(SourceChainConfigArgs[] memory sourceChainConfigUpdates) internal {
    for (uint256 i = 0; i < sourceChainConfigUpdates.length; ++i) {
      SourceChainConfigArgs memory sourceConfigUpdate = sourceChainConfigUpdates[i];
      uint64 sourceChainSelector = sourceConfigUpdate.sourceChainSelector;

      if (sourceConfigUpdate.onRamp == address(0)) {
        revert ZeroAddressNotAllowed();
      }

      SourceChainConfig storage currentConfig = s_sourceChainConfigs[sourceChainSelector];

      // OnRamp can never be zero - if it is, then the source chain has been added for the first time
      if (currentConfig.onRamp == address(0)) {
        currentConfig.metadataHash =
          _metadataHash(sourceChainSelector, sourceConfigUpdate.onRamp, Internal.EVM_2_EVM_MESSAGE_HASH);
        currentConfig.onRamp = sourceConfigUpdate.onRamp;
        currentConfig.prevOffRamp = sourceConfigUpdate.prevOffRamp;

        // s_sourceChainSelectors.push(sourceChainSelector);
        emit SourceChainSelectorAdded(sourceChainSelector);
      } else if (
        currentConfig.onRamp != sourceConfigUpdate.onRamp || currentConfig.prevOffRamp != sourceConfigUpdate.prevOffRamp
      ) {
        revert StaticConfigCannotBeUpdated();
      }

      // TODO: re-introduce check when MultiCommitStore is ready
      // Ensures we can never deploy a new offRamp that points to a commitStore that
      // already has roots committed.
      // if (ICommitStore(staticConfig.commitStore).getExpectedNextSequenceNumber() != 1) revert CommitStoreAlreadyInUse();

      // The only dynamic config is the isEnabled flag
      currentConfig.isEnabled = sourceConfigUpdate.isEnabled;
      emit SourceChainConfigSet(sourceChainSelector, currentConfig);
    }
  }

  // TODO: _beforeSetConfig is no longer used in OCR3 - replace this with an external onlyOwner function
  /// @notice Sets the dynamic config. This function is called during `setOCR2Config` flow
  function _beforeSetConfig(bytes memory onchainConfig) internal override {
    DynamicConfig memory dynamicConfig = abi.decode(onchainConfig, (DynamicConfig));

    if (dynamicConfig.router == address(0)) revert ZeroAddressNotAllowed();

    s_dynamicConfig = dynamicConfig;

    emit ConfigSet(
      StaticConfig({commitStore: i_commitStore, chainSelector: i_chainSelector, rmnProxy: i_rmnProxy}), dynamicConfig
    );
  }

  /// @notice Get all tokens which are included in Aggregate Rate Limiting.
  /// @return sourceTokens The source representation of the tokens that are rate limited.
  /// @return destTokens The destination representation of the tokens that are rate limited.
  /// @dev the order of IDs in the list is **not guaranteed**, therefore, if ordering matters when
  /// making successive calls, one should keep the blockheight constant to ensure a consistent result.
  function getAllRateLimitTokens() external view returns (address[] memory sourceTokens, address[] memory destTokens) {
    sourceTokens = new address[](s_rateLimitedTokensDestToSource.length());
    destTokens = new address[](s_rateLimitedTokensDestToSource.length());

    for (uint256 i = 0; i < s_rateLimitedTokensDestToSource.length(); ++i) {
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

  /// @notice Uses pools to release or mint a number of different tokens to a receiver address.
  /// @param sourceTokenAmounts List of tokens and amount values to be released/minted.
  /// @param messageRoute Message route details (original sender, receiver and source chain)
  /// @param encodedSourceTokenData Array of token data returned by token pools on the source chain.
  /// @param offchainTokenData Array of token data fetched offchain by the DON.
  /// @dev This function wrappes the token pool call in a try catch block to gracefully handle
  /// any non-rate limiting errors that may occur. If we encounter a rate limiting related error
  /// we bubble it up. If we encounter a non-rate limiting error we wrap it in a TokenHandlingError.
  function _releaseOrMintTokens(
    Client.EVMTokenAmount[] memory sourceTokenAmounts,
    Any2EVMMessageRoute memory messageRoute,
    bytes[] memory encodedSourceTokenData,
    bytes[] memory offchainTokenData
  ) internal returns (Client.EVMTokenAmount[] memory destTokenAmounts) {
    // Creating a copy is more gas efficient than initializing a new array.
    destTokenAmounts = sourceTokenAmounts;
    uint256 value = 0;
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      // This should never revert as the onRamp creates the sourceTokenData. Only the inner components from
      // this struct come from untrusted sources.
      Internal.SourceTokenData memory sourceTokenData =
        abi.decode(encodedSourceTokenData[i], (Internal.SourceTokenData));
      // We need to safely decode the pool address from the sourceTokenData, as it could be wrong,
      // in which case it doesn't have to be a valid EVM address.
      address localPoolAddress = Internal._validateEVMAddress(sourceTokenData.destPoolAddress);
      // This will call the supportsInterface through the ERC165Checker, and not directly on the pool address.
      // This is done to prevent a pool from reverting the entire transaction if it doesn't support the interface.
      // The call gets a max or 30k gas per instance, of which there are three. This means gas estimations should
      // account for 90k gas overhead due to the interface check.
      if (!localPoolAddress.supportsInterface(Pool.CCIP_POOL_V1)) {
        revert NotACompatiblePool(localPoolAddress);
      }

      // We determined that the pool address is a valid EVM address, but that does not mean the code at this
      // address is a (compatible) pool contract. _callWithExactGasSafeReturnData will check if the location
      // contains a contract. If it doesn't it reverts with a known error, which we catch gracefully.
      // We call the pool with exact gas to increase resistance against malicious tokens or token pools.
      // We protects against return data bombs by capping the return data size at MAX_RET_BYTES.
      (bool success, bytes memory returnData,) = CallWithExactGas._callWithExactGasSafeReturnData(
        abi.encodeWithSelector(
          IPool.releaseOrMint.selector,
          Pool.ReleaseOrMintInV1({
            originalSender: messageRoute.sender,
            receiver: messageRoute.receiver,
            amount: sourceTokenAmounts[i].amount,
            remoteChainSelector: messageRoute.sourceChainSelector,
            sourcePoolAddress: sourceTokenData.sourcePoolAddress,
            sourcePoolData: sourceTokenData.extraData,
            offchainTokenData: offchainTokenData[i]
          })
        ),
        localPoolAddress,
        s_dynamicConfig.maxPoolReleaseOrMintGas,
        Internal.GAS_FOR_CALL_EXACT_CHECK,
        Internal.MAX_RET_BYTES
      );

      // wrap and rethrow the error so we can catch it lower in the stack
      if (!success) revert TokenHandlingError(returnData);

      // If the call was successful, the returnData should be the local token address.
      if (returnData.length != Pool.CCIP_POOL_V1_RET_BYTES) {
        revert InvalidDataLength(Pool.CCIP_POOL_V1_RET_BYTES, returnData.length);
      }
      (uint256 decodedAddress, uint256 amount) = abi.decode(returnData, (uint256, uint256));
      destTokenAmounts[i].token = Internal._validateEVMAddressFromUint256(decodedAddress);
      destTokenAmounts[i].amount = amount;

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
    /* solhint-disable */
    revert();
    /* solhint-enable*/
  }
}
