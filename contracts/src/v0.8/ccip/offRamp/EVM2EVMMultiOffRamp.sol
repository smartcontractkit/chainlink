// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "../interfaces/IAny2EVMMessageReceiver.sol";
import {IAny2EVMMultiOffRamp} from "../interfaces/IAny2EVMMultiOffRamp.sol";
import {IAny2EVMOffRamp} from "../interfaces/IAny2EVMOffRamp.sol";
import {IMessageInterceptor} from "../interfaces/IMessageInterceptor.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../shared/call/CallWithExactGas.sol";
import {EnumerableMapAddresses} from "../../shared/enumerable/EnumerableMapAddresses.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {MerkleMultiProof} from "../libraries/MerkleMultiProof.sol";
import {Pool} from "../libraries/Pool.sol";
import {MultiOCR3Base} from "../ocr/MultiOCR3Base.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {ERC165Checker} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165Checker.sol";

/// @notice EVM2EVMOffRamp enables OCR networks to execute multiple messages
/// in an OffRamp in a single transaction.
/// @dev The EVM2EVMOnRamp, CommitStore and EVM2EVMOffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
/// @dev MultiOCR3Base is used to store multiple OCR configs for both the OffRamp and the CommitStore.
/// The execution plugin type has to be configured without signature verification, and the commit
/// plugin type with verification.
contract EVM2EVMMultiOffRamp is IAny2EVMMultiOffRamp, ITypeAndVersion, MultiOCR3Base {
  using ERC165Checker for address;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  error AlreadyAttempted(uint64 sourceChainSelector, uint64 sequenceNumber);
  error AlreadyExecuted(uint64 sourceChainSelector, uint64 sequenceNumber);
  error ZeroAddressNotAllowed();
  error ZeroChainSelectorNotAllowed();
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
  error RootAlreadyCommitted(uint64 sourceChainSelector, bytes32 merkleRoot);
  error InvalidRoot();
  error CanOnlySelfCall();
  error ReceiverError(bytes error);
  error TokenHandlingError(bytes error);
  error EmptyReport();
  error CursedByRMN(uint64 sourceChainSelector);
  error InvalidMessageId(bytes32 messageId);
  error NotACompatiblePool(address notPool);
  error InvalidDataLength(uint256 expected, uint256 got);
  error InvalidNewState(uint64 sourceChainSelector, uint64 sequenceNumber, Internal.MessageExecutionState newState);
  error InvalidStaticConfig(uint64 sourceChainSelector);
  error StaleCommitReport();
  error InvalidInterval(uint64 sourceChainSelector, Interval interval);
  error PausedError();

  /// @dev Atlas depends on this event, if changing, please notify Atlas.
  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event SkippedIncorrectNonce(uint64 sourceChainSelector, uint64 nonce, address sender);
  event SkippedSenderWithPreviousRampMessageInflight(uint64 sourceChainSelector, uint64 nonce, address sender);
  /// @dev RMN depends on this event, if changing, please notify the RMN maintainers.
  event ExecutionStateChanged(
    uint64 indexed sourceChainSelector,
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    Internal.MessageExecutionState state,
    bytes returnData
  );
  event SourceChainSelectorAdded(uint64 sourceChainSelector);
  event SourceChainConfigSet(uint64 indexed sourceChainSelector, SourceChainConfig sourceConfig);
  event SkippedAlreadyExecutedMessage(uint64 sourceChainSelector, uint64 sequenceNumber);
  event Paused(address account);
  event Unpaused(address account);
  /// @dev RMN depends on this event, if changing, please notify the RMN maintainers.
  event CommitReportAccepted(CommitReport report);
  event RootRemoved(bytes32 root);

  /// @notice Static offRamp config
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct StaticConfig {
    uint64 chainSelector; // ───╮  Destination chainSelector
    address rmnProxy; // ───────╯  RMN proxy address
    address tokenAdminRegistry; // Token admin registry address
  }

  /// @notice Per-chain source config (defining a lane from a Source Chain -> Dest OffRamp)
  struct SourceChainConfig {
    bool isEnabled; // ─────────╮  Flag whether the source chain is enabled or not
    uint64 minSeqNr; //         |  The min sequence number expected for future messages
    address prevOffRamp; // ────╯  Address of previous-version per-lane OffRamp. Used to be able to provide seequencing continuity during a zero downtime upgrade.
    address onRamp; //             OnRamp address on the source chain
    /// @dev Ensures that 2 identical messages sent to 2 different lanes will have a distinct hash.
    /// Must match the metadataHash used in computing leaf hashes offchain for the root committed in
    /// the commitStore and i_metadataHash in the onRamp.
    bytes32 metadataHash; //      Source-chain specific message hash preimage to ensure global uniqueness
  }

  /// @notice SourceChainConfig update args scoped to one source chain
  struct SourceChainConfigArgs {
    uint64 sourceChainSelector; //  ───╮  Source chain selector of the config to update
    bool isEnabled; //                 │  Flag whether the source chain is enabled or not
    address prevOffRamp; // ───────────╯  Address of previous-version per-lane OffRamp. Used to be able to provide sequencing continuity during a zero downtime upgrade.
    address onRamp; //                    OnRamp address on the source chain
  }

  /// @notice Dynamic offRamp config
  /// @dev since OffRampConfig is part of OffRampConfigChanged event, if changing it, we should update the ABI on Atlas
  struct DynamicConfig {
    address router; // ─────────────────────────────────╮ Router address
    uint32 permissionLessExecutionThresholdSeconds; //  │ Waiting time before manual execution is enabled
    uint32 maxTokenTransferGas; //                      │ Maximum amount of gas passed on to token `transfer` call
    uint32 maxPoolReleaseOrMintGas; // ─────────────────╯ Maximum amount of gas passed on to token pool when calling releaseOrMint
    uint16 maxNumberOfTokensPerMsg; // ──╮ Maximum number of ERC20 token transfers that can be included per message
    uint32 maxDataBytes; //              │ Maximum payload data size in bytes
    address messageValidator; // ────────╯ Optional message validator to validate incoming messages (zero address = no validator)
    address priceRegistry; // Price registry address on the local chain
  }

  /// @notice Struct that represents a message route (sender -> receiver and source chain)
  struct Any2EVMMessageRoute {
    bytes sender; //                    Message sender
    uint64 sourceChainSelector; // ───╮ Source chain that the message originates from
    address receiver; // ─────────────╯ Address that receives the message
  }

  /// @notice a sequenceNumber interval
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct Interval {
    uint64 min; // ───╮ Minimum sequence number, inclusive
    uint64 max; // ───╯ Maximum sequence number, inclusive
  }

  /// @dev Struct to hold a merkle root and an interval for a source chain so that an array of these can be passed in the CommitReport.
  struct MerkleRoot {
    uint64 sourceChainSelector; // Remote source chain selector that the Merkle Root is scoped to
    Interval interval; // Report interval of the merkle root
    bytes32 merkleRoot; // Merkle root covering the interval & source chain messages
  }

  /// @notice Report that is committed by the observing DON at the committing phase
  /// @dev RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct CommitReport {
    Internal.PriceUpdates priceUpdates; // Collection of gas and price updates to commit
    MerkleRoot[] merkleRoots; // Collection of merkle roots per source chain to commit
  }

  /// @dev Struct to hold a merkle root for a source chain so that an array of these can be passed in the resetUblessedRoots function.
  struct UnblessedRoot {
    uint64 sourceChainSelector; // Remote source chain selector that the Merkle Root is scoped to
    bytes32 merkleRoot; // Merkle root of a single remote source chain
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMMultiOffRamp 1.6.0-dev";
  /// @dev ChainSelector of this chain
  uint64 internal immutable i_chainSelector;
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;
  /// @dev The address of the token admin registry
  address internal immutable i_tokenAdminRegistry;

  // DYNAMIC CONFIG
  DynamicConfig internal s_dynamicConfig;

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

  // sourceChainSelector => merkleRoot => timestamp when received
  mapping(uint64 sourceChainSelector => mapping(bytes32 merkleRoot => uint256 timestamp)) internal s_roots;
  /// @dev The epoch and round of the last report
  uint40 private s_latestPriceEpochAndRound;
  /// @dev Whether this OffRamp is paused or not
  bool private s_paused = false;

  constructor(StaticConfig memory staticConfig, SourceChainConfigArgs[] memory sourceChainConfigs) MultiOCR3Base() {
    if (staticConfig.rmnProxy == address(0) || staticConfig.tokenAdminRegistry == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    if (staticConfig.chainSelector == 0) {
      revert ZeroChainSelectorNotAllowed();
    }

    i_chainSelector = staticConfig.chainSelector;
    i_rmnProxy = staticConfig.rmnProxy;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;

    _applySourceChainConfigUpdates(sourceChainConfigs);
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  // The size of the execution state in bits
  uint256 private constant MESSAGE_EXECUTION_STATE_BIT_WIDTH = 2;
  // The mask for the execution state bits
  uint256 private constant MESSAGE_EXECUTION_STATE_MASK = (1 << MESSAGE_EXECUTION_STATE_BIT_WIDTH) - 1;

  // ================================================================
  // │                           Execution                          │
  // ================================================================

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
    // We do this here because the other _execute path is already covered by MultiOCR3Base.
    // TODO: contract size golfing - split to internal function
    if (i_chainID != block.chainid) revert MultiOCR3Base.ForkedChain(i_chainID, uint64(block.chainid));

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

  /// @notice Transmit function for execution reports. The function takes no signatures,
  /// and expects the exec plugin type to be configured with no signatures.
  /// @param report serialized execution report
  function execute(bytes32[3] calldata reportContext, bytes calldata report) external {
    _reportExec(report);

    // TODO: gas / contract size saving from CONSTANT?
    bytes32[] memory emptySigs = new bytes32[](0);
    _transmit(uint8(Internal.OCRPluginType.Execution), reportContext, report, emptySigs, emptySigs, bytes32(""));
  }

  /// @notice Reporting function for the execution plugin
  /// @param encodedReport encoded ExecutionReport
  function _reportExec(bytes calldata encodedReport) internal {
    _batchExecute(abi.decode(encodedReport, (Internal.ExecutionReportSingleChain[])), new uint256[][](0));
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
      _executeSingleReport(reports[i], areManualGasLimitsEmpty ? emptyGasLimits : manualExecGasLimits[i]);
    }
  }

  /// @notice Executes a report, executing each message in order.
  /// @param report The execution report containing the messages and proofs.
  /// @param manualExecGasLimits An array of gas limits to use for manual execution.
  /// @dev If called from the DON, this array is always empty.
  /// @dev If called from manual execution, this array is always same length as messages.
  function _executeSingleReport(
    Internal.ExecutionReportSingleChain memory report,
    uint256[] memory manualExecGasLimits
  ) internal {
    uint64 sourceChainSelector = report.sourceChainSelector;
    // TODO: re-use isCursed / isUnpaused check from _verify here
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(sourceChainSelector)))) revert CursedByRMN(sourceChainSelector);

    uint256 numMsgs = report.messages.length;
    if (numMsgs == 0) revert EmptyReport();
    if (numMsgs != report.offchainTokenData.length) revert UnexpectedTokenData();

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
    uint256 timestampCommitted = _verify(sourceChainSelector, hashedLeaves, report.proofs, report.proofFlagBits);
    if (timestampCommitted == 0) revert RootNotCommitted(sourceChainSelector);

    // Execute messages
    bool manualExecution = manualExecGasLimits.length != 0;
    for (uint256 i = 0; i < numMsgs; ++i) {
      Internal.EVM2EVMMessage memory message = report.messages[i];

      Internal.MessageExecutionState originalState = getExecutionState(sourceChainSelector, message.sequenceNumber);
      if (originalState == Internal.MessageExecutionState.SUCCESS) {
        // If the message has already been executed, we skip it.  We want to not revert on race conditions between
        // executing parties. This will allow us to open up manual exec while also attempting with the DON, without
        // reverting an entire DON batch when a user manually executes while the tx is inflight.
        emit SkippedAlreadyExecutedMessage(sourceChainSelector, message.sequenceNumber);
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
      ) revert AlreadyExecuted(sourceChainSelector, message.sequenceNumber);

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
          revert AlreadyAttempted(sourceChainSelector, message.sequenceNumber);
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
      // TODO: GAS GOLF - evaluate caching sequenceNumber instead of offchainTokenData
      bytes[] memory offchainTokenData = report.offchainTokenData[i];
      _isWellFormed(
        message.messageId,
        sourceChainSelector,
        message.sequenceNumber,
        message.tokenAmounts.length,
        message.data.length,
        offchainTokenData.length
      );

      _setExecutionState(sourceChainSelector, message.sequenceNumber, Internal.MessageExecutionState.IN_PROGRESS);
      (Internal.MessageExecutionState newState, bytes memory returnData) = _trialExecute(message, offchainTokenData);
      _setExecutionState(sourceChainSelector, message.sequenceNumber, newState);

      // Since it's hard to estimate whether manual execution will succeed, we
      // revert the entire transaction if it fails. This will show the user if
      // their manual exec will fail before they submit it.
      if (
        manualExecution && newState == Internal.MessageExecutionState.FAILURE
          && originalState != Internal.MessageExecutionState.UNTOUCHED
      ) {
        // If manual execution fails, we revert the entire transaction, unless the originalState is UNTOUCHED as we
        // would still be making progress by changing the state from UNTOUCHED to FAILURE.
        revert ExecutionError(message.messageId, returnData);
      }

      // The only valid prior states are UNTOUCHED and FAILURE (checked above)
      // The only valid post states are FAILURE and SUCCESS (checked below)
      if (newState != Internal.MessageExecutionState.FAILURE && newState != Internal.MessageExecutionState.SUCCESS) {
        revert InvalidNewState(sourceChainSelector, message.sequenceNumber, newState);
      }

      // Nonce changes per state transition
      // UNTOUCHED -> FAILURE  nonce bump
      // UNTOUCHED -> SUCCESS  nonce bump
      // FAILURE   -> FAILURE  no nonce bump
      // FAILURE   -> SUCCESS  no nonce bump
      if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
        s_senderNonce[sourceChainSelector][message.sender]++;
      }

      emit ExecutionStateChanged(sourceChainSelector, message.sequenceNumber, message.messageId, newState, returnData);
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
    // TODO: move maxNumberOfTokens & data length validation offchain
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
      bytes4 errorSelector = bytes4(err);
      if (
        ReceiverError.selector == errorSelector || TokenHandlingError.selector == errorSelector
          || Internal.InvalidEVMAddress.selector == errorSelector || InvalidDataLength.selector == errorSelector
          || CallWithExactGas.NoContract.selector == errorSelector || NotACompatiblePool.selector == errorSelector
          || IMessageInterceptor.MessageValidationError.selector == errorSelector
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

    Client.Any2EVMMessage memory any2EvmMessage = Internal._toAny2EVMMessage(message, destTokenAmounts);

    address messageValidator = s_dynamicConfig.messageValidator;
    if (messageValidator != address(0)) {
      try IMessageInterceptor(messageValidator).onIncomingMessage(any2EvmMessage) {}
      catch (bytes memory err) {
        revert IMessageInterceptor.MessageValidationError(err);
      }
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
      any2EvmMessage, Internal.GAS_FOR_CALL_EXACT_CHECK, message.gasLimit, message.receiver
    );
    // If CCIP receiver execution is not successful, revert the call including token transfers
    if (!success) revert ReceiverError(returnData);
  }

  /// @notice creates a unique hash to be used in message hashing.
  function _metadataHash(uint64 sourceChainSelector, address onRamp, bytes32 prefix) internal view returns (bytes32) {
    return keccak256(abi.encode(prefix, sourceChainSelector, i_chainSelector, onRamp));
  }

  // ================================================================
  // │                           Commit                             │
  // ================================================================

  /// @notice Transmit function for commit reports. The function requires signatures,
  /// and expects the commit plugin type to be configured with signatures.
  /// @param report serialized commit report
  /// @dev A commitReport can have two distinct parts (batched together to amortize the cost of checking sigs):
  /// 1. Price updates
  /// 2. A batch of merkle root and sequence number intervals (per-source)
  /// Both have their own, separate, staleness checks, with price updates using the epoch and round
  /// number of the latest price update. The merkle root checks for staleness based on the seqNums.
  /// They need to be separate because a price report for round t+2 might be included before a report
  /// containing a merkle root for round t+1. This merkle root report for round t+1 is still valid
  /// and should not be rejected. When a report with a stale root but valid price updates is submitted,
  /// we are OK to revert to preserve the invariant that we always revert on invalid sequence number ranges.
  /// If that happens, prices will be updates in later rounds.
  function commit(
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) external {
    _reportCommit(report, uint40(uint256(reportContext[1])));
    _transmit(uint8(Internal.OCRPluginType.Commit), reportContext, report, rs, ss, rawVs);
  }

  /// @notice Reporting function for the commit plugin
  /// @param encodedReport encoded CommitReport
  /// @param epochAndRound Epoch and round of the report
  function _reportCommit(bytes calldata encodedReport, uint40 epochAndRound) internal whenNotPaused {
    CommitReport memory report = abi.decode(encodedReport, (CommitReport));

    // Check if the report contains price updates
    if (report.priceUpdates.tokenPriceUpdates.length > 0 || report.priceUpdates.gasPriceUpdates.length > 0) {
      // Check for price staleness based on the epoch and round
      if (s_latestPriceEpochAndRound < epochAndRound) {
        // If prices are not stale, update the latest epoch and round
        s_latestPriceEpochAndRound = epochAndRound;
        // And update the prices in the price registry
        IPriceRegistry(s_dynamicConfig.priceRegistry).updatePrices(report.priceUpdates);

        // If there is no root, the report only contained fee updated and
        // we return to not revert on the empty root check below.
        if (report.merkleRoots.length == 0) return;
      } else {
        // If prices are stale and the report doesn't contain a root, this report
        // does not have any valid information and we revert.
        // If it does contain a merkle root, continue to the root checking section.
        if (report.merkleRoots.length == 0) revert StaleCommitReport();
      }
    }

    for (uint256 i = 0; i < report.merkleRoots.length; ++i) {
      MerkleRoot memory root = report.merkleRoots[i];
      uint64 sourceChainSelector = root.sourceChainSelector;

      if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(sourceChainSelector)))) revert CursedByRMN(sourceChainSelector);

      SourceChainConfig storage sourceChainConfig = s_sourceChainConfigs[sourceChainSelector];

      if (!sourceChainConfig.isEnabled) revert SourceChainNotEnabled(sourceChainSelector);
      // If we reached this section, the report should contain a valid root
      if (sourceChainConfig.minSeqNr != root.interval.min || root.interval.min > root.interval.max) {
        revert InvalidInterval(root.sourceChainSelector, root.interval);
      }

      // TODO: confirm how RMN offchain blessing impacts commit report
      if (root.merkleRoot == bytes32(0)) revert InvalidRoot();
      // Disallow duplicate roots as that would reset the timestamp and
      // delay potential manual execution.
      if (s_roots[root.sourceChainSelector][root.merkleRoot] != 0) {
        revert RootAlreadyCommitted(root.sourceChainSelector, root.merkleRoot);
      }

      sourceChainConfig.minSeqNr = root.interval.max + 1;
      s_roots[root.sourceChainSelector][root.merkleRoot] = block.timestamp;
    }

    emit CommitReportAccepted(report);
  }

  /// @notice Returns the epoch and round of the last price update.
  /// @return the latest price epoch and round.
  function getLatestPriceEpochAndRound() public view returns (uint64) {
    return s_latestPriceEpochAndRound;
  }

  /// @notice Sets the latest epoch and round for price update.
  /// @param latestPriceEpochAndRound The new epoch and round for prices.
  function setLatestPriceEpochAndRound(uint40 latestPriceEpochAndRound) external onlyOwner {
    s_latestPriceEpochAndRound = latestPriceEpochAndRound;
  }

  /// @notice Returns the timestamp of a potentially previously committed merkle root.
  /// If the root was never committed 0 will be returned.
  /// @param sourceChainSelector The source chain selector.
  /// @param root The merkle root to check the commit status for.
  /// @return the timestamp of the committed root or zero in the case that it was never
  /// committed.
  function getMerkleRoot(uint64 sourceChainSelector, bytes32 root) external view returns (uint256) {
    return s_roots[sourceChainSelector][root];
  }

  /// @notice Returns if a root is blessed or not.
  /// @param root The merkle root to check the blessing status for.
  /// @return whether the root is blessed or not.
  function isBlessed(bytes32 root) public view returns (bool) {
    // TODO: update RMN to also consider the source chain selector for blessing
    return IRMN(i_rmnProxy).isBlessed(IRMN.TaggedRoot({commitStore: address(this), root: root}));
  }

  /// @notice Used by the owner in case an invalid sequence of roots has been
  /// posted and needs to be removed. The interval in the report is trusted.
  /// @param rootToReset The roots that will be reset. This function will only
  /// reset roots that are not blessed.
  function resetUnblessedRoots(UnblessedRoot[] calldata rootToReset) external onlyOwner {
    for (uint256 i = 0; i < rootToReset.length; ++i) {
      UnblessedRoot memory root = rootToReset[i];
      if (!isBlessed(root.merkleRoot)) {
        delete s_roots[root.sourceChainSelector][root.merkleRoot];
        emit RootRemoved(root.merkleRoot);
      }
    }
  }

  /// @notice Returns timestamp of when root was accepted or 0 if verification fails.
  /// @dev This method uses a merkle tree within a merkle tree, with the hashedLeaves,
  /// proofs and proofFlagBits being used to get the root of the inner tree.
  /// This root is then used as the singular leaf of the outer tree.
  function _verify(
    uint64 sourceChainSelector,
    bytes32[] memory hashedLeaves,
    bytes32[] memory proofs,
    uint256 proofFlagBits
  ) internal view virtual whenNotPaused returns (uint256 timestamp) {
    bytes32 root = MerkleMultiProof.merkleRoot(hashedLeaves, proofs, proofFlagBits);
    // Only return non-zero if present and blessed.
    if (!isBlessed(root)) {
      return 0;
    }
    return s_roots[sourceChainSelector][root];
  }

  /// @inheritdoc MultiOCR3Base
  function _afterOCR3ConfigSet(uint8 ocrPluginType) internal override {
    if (ocrPluginType == uint8(Internal.OCRPluginType.Commit)) {
      // When the OCR config changes, we reset the price epoch and round
      // since epoch and rounds are scoped per config digest.
      // Note that s_minSeqNr/roots do not need to be reset as the roots persist
      // across reconfigurations and are de-duplicated separately.
      s_latestPriceEpochAndRound = 0;
    }
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static config.
  /// @dev This function will always return the same struct as the contents is static and can never change.
  /// RMN depends on this function, if changing, please notify the RMN maintainers.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return
      StaticConfig({chainSelector: i_chainSelector, rmnProxy: i_rmnProxy, tokenAdminRegistry: i_tokenAdminRegistry});
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

      if (sourceChainSelector == 0) {
        revert ZeroChainSelectorNotAllowed();
      }

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
        currentConfig.minSeqNr = 1;

        emit SourceChainSelectorAdded(sourceChainSelector);
      } else if (
        currentConfig.onRamp != sourceConfigUpdate.onRamp || currentConfig.prevOffRamp != sourceConfigUpdate.prevOffRamp
      ) {
        revert InvalidStaticConfig(sourceChainSelector);
      }

      // The only dynamic config is the isEnabled flag
      currentConfig.isEnabled = sourceConfigUpdate.isEnabled;
      emit SourceChainConfigSet(sourceChainSelector, currentConfig);
    }
  }

  /// @notice Sets the dynamic config.
  function setDynamicConfig(DynamicConfig memory dynamicConfig) external onlyOwner {
    if (dynamicConfig.priceRegistry == address(0)) revert ZeroAddressNotAllowed();
    if (dynamicConfig.router == address(0)) revert ZeroAddressNotAllowed();

    s_dynamicConfig = dynamicConfig;

    // TODO: contract size golfing - is StaticConfig needed in the event?
    emit ConfigSet(
      StaticConfig({chainSelector: i_chainSelector, rmnProxy: i_rmnProxy, tokenAdminRegistry: i_tokenAdminRegistry}),
      dynamicConfig
    );
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
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      // This should never revert as the onRamp creates the sourceTokenData. Only the inner components from
      // this struct come from untrusted sources.
      Internal.SourceTokenData memory sourceTokenData =
        abi.decode(encodedSourceTokenData[i], (Internal.SourceTokenData));
      // We need to safely decode the pool address from the sourceTokenData, as it could be wrong,
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

      // We determined that the pool address is a valid EVM address, but that does not mean the code at this
      // address is a (compatible) pool contract. _callWithExactGasSafeReturnData will check if the location
      // contains a contract. If it doesn't it reverts with a known error, which we catch gracefully.
      // We call the pool with exact gas to increase resistance against malicious tokens or token pools.
      // We protects against return data bombs by capping the return data size at MAX_RET_BYTES.
      (bool success, bytes memory returnData,) = CallWithExactGas._callWithExactGasSafeReturnData(
        abi.encodeWithSelector(
          IPoolV1.releaseOrMint.selector,
          Pool.ReleaseOrMintInV1({
            originalSender: messageRoute.sender,
            receiver: messageRoute.receiver,
            amount: sourceTokenAmounts[i].amount,
            localToken: localToken,
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
      uint256 amount = abi.decode(returnData, (uint256));

      (success, returnData,) = CallWithExactGas._callWithExactGasSafeReturnData(
        abi.encodeWithSelector(IERC20.transfer.selector, messageRoute.receiver, amount),
        localToken,
        s_dynamicConfig.maxTokenTransferGas,
        Internal.GAS_FOR_CALL_EXACT_CHECK,
        Internal.MAX_RET_BYTES
      );

      if (!success) revert TokenHandlingError(returnData);

      destTokenAmounts[i].token = localToken;
      destTokenAmounts[i].amount = amount;
    }

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

  /// @notice Single function to check the status of the commitStore.
  function isUnpausedAndNotCursed(uint64 sourceChainSelector) external view returns (bool) {
    return !IRMN(i_rmnProxy).isCursed(bytes16(uint128(sourceChainSelector))) && !s_paused;
  }

  // TODO: global pausing can be removed delegated to the i_rmnProxy
  /// @notice Modifier to make a function callable only when the contract is not paused.
  modifier whenNotPaused() {
    if (paused()) revert PausedError();
    _;
  }

  /// @notice Returns true if the contract is paused, and false otherwise.
  function paused() public view returns (bool) {
    return s_paused;
  }

  /// @notice Pause the contract
  /// @dev only callable by the owner
  function pause() external onlyOwner {
    s_paused = true;
    emit Paused(msg.sender);
  }

  /// @notice Unpause the contract
  /// @dev only callable by the owner
  function unpause() external onlyOwner {
    s_paused = false;
    emit Unpaused(msg.sender);
  }
}
