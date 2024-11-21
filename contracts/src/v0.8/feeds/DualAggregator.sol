// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {AccessControllerInterface} from "../shared/interfaces/AccessControllerInterface.sol";
import {AggregatorV2V3Interface} from "../shared/interfaces/AggregatorV2V3Interface.sol";
import {AggregatorValidatorInterface} from "../shared/interfaces/AggregatorValidatorInterface.sol";
import {LinkTokenInterface} from "../shared/interfaces/LinkTokenInterface.sol";

import {SimpleReadAccessController} from "../shared/access/SimpleReadAccessController.sol";
import {CallWithExactGas} from "../shared/call/CallWithExactGas.sol";
import {OCR2Abstract} from "../shared/ocr2/OCR2Abstract.sol";

// This contract is a port of OCR2Aggregator from `libocr` it is being used
// for a new feeds based project that is ongoing, there will be some modernization
// that happens to this contract as the project progresses.
// solhint-disable max-states-count
contract DualAggregator is OCR2Abstract, AggregatorV2V3Interface, SimpleReadAccessController {
  string public constant override typeAndVersion = "DualAggregator 1.0.0";

  // This contract is divided into sections. Each section defines a set of
  // variables, events, and functions that belong together.

  // ================================================================
  // │            Variables used in multiple other sections         │
  // ================================================================

  // Transmitter information.
  struct Transmitter {
    bool active; // ─────────╮ True if active.
    uint8 index; //          │ Index in `s_transmittersList`.
    //                       │
    //                       │ Juels-denominated payment for transmitters, covering gas costs incurred
    //                       │ by the transmitter plus additional rewards. The entire LINK supply (1e9
    uint96 paymentJuels; // ─╯ LINK = 1e27 Juels) will always fit into a uint96.
  }

  // Signer information.
  struct Signer {
    bool active; // ─╮ True if active.
    uint8 index; // ─╯ Index of oracle in `s_signersList`.
  }

  // Storing these fields used on the hot path in a HotVars variable reduces the
  // retrieval of all of them to one SLOAD.
  struct HotVars {
    uint8 f; //  ─────────────────────────╮ Maximum number of faulty oracles.
    //                                    │
    uint40 latestEpochAndRound; //        │ Epoch and round from OCR protocol,
    //                                    │ 32 most significant bits for epoch, 8 least sig bits for round.
    //                                    │
    uint32 latestAggregatorRoundId; //    │ Chainlink Aggregators expose a roundId to consumers. The offchain protocol
    //                                    │ does not use this id anywhere. We increment it whenever a new transmission
    //                                    │ is made to provide callers with contiguous ids for successive reports.
    uint32 latestSecondaryRoundId; //     │ Latest transmission round arrived from the Secondary Proxy.
    uint32 maximumGasPriceGwei; //        │ Highest compensated gas price, in gwei uints.
    uint32 reasonableGasPriceGwei; //     │ If gas price is less (in gwei units), transmitter gets half the savings.
    uint32 observationPaymentGjuels; //   │ Fixed LINK reward for each observer.
    uint32 transmissionPaymentGjuels; //  │ Fixed reward for transmitter.
    bool isLatestSecondary; // ───────────╯ Whether the latest report was secondary or not
  }

  /// @notice mapping containing the transmitter information of each transmitter address.
  mapping(address transmitter => Transmitter) internal s_transmitters;

  /// @notice mapping containing the signer information of each signer address.
  mapping(address signer => Signer) internal s_signers;

  /// @notice s_signersList contains the signing address of each oracle.
  address[] internal s_signersList;

  /// @notice s_transmittersList contains the transmission address of each oracle,
  /// i.e. the address the oracle actually sends transactions to the contract from.
  address[] internal s_transmittersList;

  /// @notice We assume that all oracles contribute observations to all rounds. This
  /// variable tracks (per-oracle) from what round an oracle should be rewarded,
  /// i.e. the oracle gets (latestAggregatorRoundId - rewardFromAggregatorRoundId) * reward.
  uint32[MAX_NUM_ORACLES] internal s_rewardFromAggregatorRoundId;

  /// @notice latest setted config.
  bytes32 internal s_latestConfigDigest;

  /// @notice overhead incurred by accounting logic.
  uint24 internal s_accountingGas;

  /// @notice most common fields used on the hot path.
  HotVars internal s_hotVars;

  /// @notice lowest answer the system is allowed to report in response to transmissions.
  int192 internal immutable i_minAnswer;

  /// @notice highest answer the system is allowed to report in response to transmissions.
  int192 internal immutable i_maxAnswer;

  /// @param link address of the LINK contract.
  /// @param minAnswer_ lowest answer the median of a report is allowed to be.
  /// @param maxAnswer_ highest answer the median of a report is allowed to be.
  /// @param billingAccessController access controller for managing the billing.
  /// @param requesterAccessController access controller for requesting new rounds.
  /// @param decimals_ answers are stored in fixed-point format, with this many digits of precision.
  /// @param description_ short human-readable description of observable this contract's answers pertain to.
  /// @param secondaryProxy_ proxy address to manage the secondary reports.
  /// @param cutoffTime_ timetamp to define the window in which a secondary report is valid.
  /// @param maxSyncIterations_ max iterations the secondary proxy will be able to loop to sync with the primary rounds.
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_,
    address secondaryProxy_,
    uint32 cutoffTime_,
    uint32 maxSyncIterations_
  ) {
    i_decimals = decimals_;
    i_minAnswer = minAnswer_;
    i_maxAnswer = maxAnswer_;
    i_secondaryProxy = secondaryProxy_;

    s_linkToken = link;
    emit LinkTokenSet(LinkTokenInterface(address(0)), link);

    _setBillingAccessController(billingAccessController);
    setRequesterAccessController(requesterAccessController);
    setValidatorConfig(AggregatorValidatorInterface(address(0x0)), 0);

    s_cutoffTime = cutoffTime_;
    s_maxSyncIterations = maxSyncIterations_;
    s_description = description_;
  }

  // ================================================================
  // │                  OCR2Abstract Configuration                  │
  // ================================================================

  // SetConfig information
  struct SetConfigArgs {
    uint64 offchainConfigVersion; // ─╮ OffchainConfig version.
    uint8 f; // ──────────────────────╯ Faulty Oracles amount.
    bytes onchainConfig; //             Onchain configuration.
    bytes offchainConfig; //            Offchain configuration.
    address[] signers; //               Signing addresses of each oracle.
    address[] transmitters; //          Transmitting addresses of each oracle.
  }

  error FMustBePositive();
  error TooManyOracles();
  error OracleLengthMismatch();
  error FaultyOracleFTooHigh();
  error InvalidOnChainConfig();
  error RepeatedSignerAddress();
  error RepeatedTransmitterAddress();

  /// @notice incremented each time a new config is posted. This count is incorporated
  /// into the config digest to prevent replay attacks.
  uint32 internal s_configCount;

  /// @notice makes it easier for offchain systems to extract config from logs.
  uint32 internal s_latestConfigBlockNumber;

  /// @notice check if `f` is a positive number.
  /// @dev left as a function so this check can be disabled in derived contracts.
  /// @param f amount of faulty oracles to check.
  function _requirePositiveF(
    uint256 f
  ) internal pure virtual {
    if (f <= 0) {
      revert FMustBePositive();
    }
  }

  /// @inheritdoc OCR2Abstract
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override onlyOwner {
    if (signers.length > MAX_NUM_ORACLES) {
      revert TooManyOracles();
    }
    if (signers.length != transmitters.length) {
      revert OracleLengthMismatch();
    }
    if (3 * f >= signers.length) {
      revert FaultyOracleFTooHigh();
    }
    _requirePositiveF(f);
    if (keccak256(onchainConfig) != keccak256(abi.encodePacked(uint8(1), /*version*/ i_minAnswer, i_maxAnswer))) {
      revert InvalidOnChainConfig();
    }

    SetConfigArgs memory args = SetConfigArgs({
      signers: signers,
      transmitters: transmitters,
      f: f,
      onchainConfig: onchainConfig,
      offchainConfigVersion: offchainConfigVersion,
      offchainConfig: offchainConfig
    });

    s_hotVars.latestEpochAndRound = 0;
    _payOracles();

    // Remove any old signer/transmitter addresses.
    uint256 oldLength = s_signersList.length;
    for (uint256 i = 0; i < oldLength; ++i) {
      address signer = s_signersList[i];
      address transmitter = s_transmittersList[i];
      delete s_signers[signer];
      delete s_transmitters[transmitter];
    }
    delete s_signersList;
    delete s_transmittersList;

    // Add new signer/transmitter addresses.
    for (uint256 i = 0; i < args.signers.length; ++i) {
      if (s_signers[args.signers[i]].active) {
        revert RepeatedSignerAddress();
      }
      s_signers[args.signers[i]] = Signer({active: true, index: uint8(i)});
      if (s_transmitters[args.transmitters[i]].active) {
        revert RepeatedTransmitterAddress();
      }
      s_transmitters[args.transmitters[i]] = Transmitter({active: true, index: uint8(i), paymentJuels: 0});
    }
    s_signersList = args.signers;
    s_transmittersList = args.transmitters;

    s_hotVars.f = args.f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    s_latestConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );

    emit ConfigSet(
      previousConfigBlockNumber,
      s_latestConfigDigest,
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );

    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    for (uint256 i = 0; i < args.signers.length; ++i) {
      s_rewardFromAggregatorRoundId[i] = latestAggregatorRoundId;
    }
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_configCount, s_latestConfigBlockNumber, s_latestConfigDigest);
  }

  /// @notice get the transmitters list.
  /// @dev The list will match the order used to specify the transmitter during setConfig.
  /// @return s_transmittersList list of addresses permitted to transmit reports to this contract.
  function getTransmitters() external view returns (address[] memory) {
    return s_transmittersList;
  }

  /// @notice Get the mininum answer value.
  /// @return minAnswer the lowest answer the system is allowed to report in a transmission.
  function minAnswer() public view returns (int256) {
    return i_minAnswer;
  }

  /// @notice Get the maximum answer value.
  /// @return maxAnswer the highest answer the system is allowed to report in a transmission.
  function maxAnswer() public view returns (int256) {
    return i_maxAnswer;
  }

  // ================================================================
  // │                      Onchain Validation                      │
  // ================================================================

  // Configuration for validator.
  struct ValidatorConfig {
    AggregatorValidatorInterface validator; // ─╮ Validator contract interface.
    uint32 gasLimit; // ────────────────────────╯ Gas limit defined for the validation call.
  }

  /// @notice indicates that the validator configuration has been set.
  /// @param previousValidator previous validator contract.
  /// @param previousGasLimit previous gas limit for validate calls.
  /// @param currentValidator current validator contract.
  /// @param currentGasLimit current gas limit for validate calls.
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  error InsufficientGas();

  /// @notice contstant exact gas cushion defined to do a call.
  uint16 private constant CALL_WITH_EXACT_GAS_CUSHION = 5_000;

  /// @notice validator configuration.
  ValidatorConfig private s_validatorConfig;

  /// @notice get the validator configuration.
  /// @return validator validator contract.
  /// @return gasLimit gas limit for validate calls.
  function getValidatorConfig() external view returns (AggregatorValidatorInterface validator, uint32 gasLimit) {
    ValidatorConfig memory vc = s_validatorConfig;
    return (vc.validator, vc.gasLimit);
  }

  /// @notice sets validator configuration.
  /// @dev set newValidator to 0x0 to disable validate calls.
  /// @param newValidator address of the new validator contract.
  /// @param newGasLimit new gas limit for validate calls.
  function setValidatorConfig(AggregatorValidatorInterface newValidator, uint32 newGasLimit) public onlyOwner {
    ValidatorConfig memory previous = s_validatorConfig;

    if (previous.validator != newValidator || previous.gasLimit != newGasLimit) {
      s_validatorConfig = ValidatorConfig({validator: newValidator, gasLimit: newGasLimit});

      emit ValidatorConfigSet(previous.validator, previous.gasLimit, newValidator, newGasLimit);
    }
  }

  /// @notice validate the answer against the validator configuration.
  /// @param aggregatorRoundId report round id to validate.
  /// @param answer report answer to validate.
  function _validateAnswer(uint32 aggregatorRoundId, int256 answer) private {
    ValidatorConfig memory vc = s_validatorConfig;

    if (address(vc.validator) == address(0)) {
      return;
    }

    uint32 prevAggregatorRoundId = aggregatorRoundId - 1;
    int256 prevAggregatorRoundAnswer = s_transmissions[prevAggregatorRoundId].answer;

    (, bool sufficientGas) = CallWithExactGas._callWithExactGasEvenIfTargetIsNoContract(
      abi.encodeCall(
        AggregatorValidatorInterface.validate,
        (uint256(prevAggregatorRoundId), prevAggregatorRoundAnswer, uint256(aggregatorRoundId), answer)
      ),
      address(vc.validator),
      vc.gasLimit,
      CALL_WITH_EXACT_GAS_CUSHION
    );

    if (!sufficientGas) {
      revert InsufficientGas();
    }
  }

  // ================================================================
  // │                       RequestNewRound                        │
  // ================================================================

  /// @notice contract address with AccessController Interface.
  AccessControllerInterface internal s_requesterAccessController;

  /// @notice emitted when a new requester access controller contract is set.
  /// @param old the address prior to the current setting.
  /// @param current the address of the new access controller contract.
  event RequesterAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  /// @notice emitted to immediately request a new round.
  /// @param requester the address of the requester.
  /// @param configDigest the latest transmission's configDigest.
  /// @param epoch the latest transmission's epoch.
  /// @param round the latest transmission's round.
  event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round);

  error OnlyOwnerAndRequesterCanCall();

  /// @notice address of the requester access controller contract.
  /// @return s_requesterAccessController requester access controller address.
  function getRequesterAccessController() external view returns (AccessControllerInterface) {
    return s_requesterAccessController;
  }

  /// @notice sets the new requester access controller.
  /// @param requesterAccessController designates the address of the new requester access controller.
  function setRequesterAccessController(
    AccessControllerInterface requesterAccessController
  ) public onlyOwner {
    AccessControllerInterface oldController = s_requesterAccessController;

    if (requesterAccessController != oldController) {
      s_requesterAccessController = AccessControllerInterface(requesterAccessController);
      emit RequesterAccessControllerSet(oldController, requesterAccessController);
    }
  }

  /// @notice immediately requests a new round.
  /// @dev access control provided by requesterAccessController.
  /// @return aggregatorRoundId round id of the next round. Note: The report for this round may have been
  /// transmitted (but not yet mined) *before* requestNewRound() was even called. There is *no*
  /// guarantee of causality between the request and the report at aggregatorRoundId.
  function requestNewRound() external returns (uint80) {
    if (msg.sender != owner() && !s_requesterAccessController.hasAccess(msg.sender, msg.data)) {
      revert OnlyOwnerAndRequesterCanCall();
    }

    uint40 latestEpochAndRound = s_hotVars.latestEpochAndRound;
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;

    emit RoundRequested(msg.sender, s_latestConfigDigest, uint32(latestEpochAndRound >> 8), uint8(latestEpochAndRound));
    return latestAggregatorRoundId + 1;
  }

  // ================================================================
  // │                       Secondary Proxy                        │
  // ================================================================

  // Used to relieve stack pressure in transmit.
  struct Report {
    int192 juelsPerFeeCoin; // ───────╮ Exchange rate between feeCoin (e.g. ETH on Ethereum) and LINK, denominated in juels.
    uint32 observationsTimestamp; // ─╯ Timestamp when the observations were made offchain.
    bytes observers; //                 i-th element is the index of the ith observer.
    int192[] observations; //           i-th element is the ith observation.
  }

  // Transmission records the median answer from the transmit transaction at
  // time timestamp.
  struct Transmission {
    int192 answer; // ───────────────╮ 192 bits ought to be enough for anyone.
    uint32 observationsTimestamp; // │ When were observations made offchain.
    uint32 recordedTimestamp; // ────╯ When was report received onchain.
  }

  /// @notice indicates that a new report arrived from the secondary feed and the round id was updated.
  /// @param secondaryRoundId the new secondary round id.
  event SecondaryRoundIdUpdated(uint32 indexed secondaryRoundId);

  /// @notice indicates that a new report arrived from the primary feed and the report had already been stored .
  /// @param primaryRoundId the new primary round id (if we're at the next block since the report it should be the same).
  event PrimaryFeedUnlocked(uint32 indexed primaryRoundId);

  /// @notice emitted when a new cutoff time is set.
  /// @param cutoffTime the new defined cutoff time.
  event CutoffTimeSet(uint32 cutoffTime);

  /// @notice emitted when a new max sync iterations is set.
  /// @param maxSyncIterations the new defined max sync iterations.
  event MaxSyncIterationsSet(uint32 maxSyncIterations);

  /// @notice revert when the loop reaches the max sync iterations amount.
  error MaxSyncIterations();

  /// @notice mapping containing the Transmission records of each round id.
  mapping(uint32 aggregatorRoundId => Transmission) internal s_transmissions;

  /// @notice secondary proxy address, used to detect who's calling the contract methods.
  address internal immutable i_secondaryProxy;

  /// @notice cutoff time defines the time window in which a secondary report is valid.
  uint32 internal s_cutoffTime;

  /// @notice max iterations the secondary proxy will be able to loop to sync with the primary rounds.
  uint32 internal s_maxSyncIterations;

  /// @notice sets the max time cutoff.
  /// @param _cutoffTime new max cutoff timestamp.
  function setCutoffTime(
    uint32 _cutoffTime
  ) external onlyOwner {
    s_cutoffTime = _cutoffTime;
    emit CutoffTimeSet(_cutoffTime);
  }

  /// @notice sets the max sync iterations.
  /// @param _maxSyncIterations new max sync iterations.
  function setMaxSyncIterations(
    uint32 _maxSyncIterations
  ) external onlyOwner {
    s_maxSyncIterations = _maxSyncIterations;
    emit MaxSyncIterationsSet(_maxSyncIterations);
  }

  /// @notice sync data with the primary rounds, return the freshest valid round id.
  /// @return roundId synced round id with the primary feed.
  function _getSyncPrimaryRound() internal view returns (uint32 roundId) {
    // Get the latest round id and the max iterations.
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    uint32 maxIterations = s_maxSyncIterations;

    // Decreasing loop from the latest primary round id.
    for (uint32 round_ = latestAggregatorRoundId; round_ > 0; --round_) {
      // In case the loop reached the maxIterations, break it.
      if (latestAggregatorRoundId - round_ == maxIterations) {
        break;
      }

      // In case it's the latest secondary round id, return it.
      if (round_ == s_hotVars.latestSecondaryRoundId) {
        return round_;
      }

      // Check if this round does not accomplish the cutoff time condition.
      if (s_transmissions[round_].recordedTimestamp + s_cutoffTime < block.timestamp) {
        return round_;
      }
    }

    // If the loop couldn't find a match, return the latest secondary round id.
    return s_hotVars.latestSecondaryRoundId;
  }

  /// @notice check if a report has already been transmitted.
  /// @param report the report to check.
  /// @return exist wether the report exist or not.
  /// @return roundId the round id where the report was found.
  function _doesReportExist(
    Report memory report
  ) internal view returns (bool exist, uint32 roundId) {
    // Get the latest round id.
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    uint32 maxIterations = s_maxSyncIterations;

    for (uint32 round_ = latestAggregatorRoundId; round_ > 0; --round_) {
      if (latestAggregatorRoundId - round_ == maxIterations) {
        revert MaxSyncIterations();
      }

      Transmission memory transmission = s_transmissions[round_];

      if (transmission.observationsTimestamp < report.observationsTimestamp) {
        return (false, 0);
      }
      // Get median.
      int192 median = report.observations[report.observations.length / 2];

      if (transmission.observationsTimestamp == report.observationsTimestamp && transmission.answer == median) {
        return (true, round_);
      }
    }

    return (false, 0);
  }

  /// @notice aggregator round in which the latest report was conceded depending on the caller.
  /// @return roundId the latest valid round id.
  function _getLatestRound() internal view returns (uint32) {
    // Get the latest round ids.
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    uint32 latestSecondaryRoundId = s_hotVars.latestSecondaryRoundId;
    bool isLatestSecondary = s_hotVars.isLatestSecondary;

    // Check if the message sender is the secondary proxy.
    if (msg.sender == i_secondaryProxy) {
      // In case the latest secondary round does not accomplish the cutoff time condition,
      // get the round id syncing with the primary rounds.
      if (s_transmissions[latestSecondaryRoundId].recordedTimestamp + s_cutoffTime < block.timestamp) {
        return _getSyncPrimaryRound();
      }

      // In case the latest secondary round accomplish the cutoff time condition, return it.
      return latestSecondaryRoundId;
    }
    // In case the report was sent by the secondary proxy.
    if (latestAggregatorRoundId == latestSecondaryRoundId) {
      // In case the transmission was sent in this same block only by the secondary proxy, return the previous round id.
      if (isLatestSecondary && s_transmissions[latestAggregatorRoundId].recordedTimestamp == block.timestamp) {
        return latestAggregatorRoundId - 1;
      }
    }

    return latestAggregatorRoundId;
  }

  // ================================================================
  // │                        Transmission                          │
  // ================================================================

  /// @notice indicates that a new report was transmitted.
  /// @param aggregatorRoundId the round to which this report was assigned.
  /// @param answer median of the observations attached to this report.
  /// @param transmitter address from which the report was transmitted.
  /// @param observationsTimestamp when were observations made offchain.
  /// @param observations observations transmitted with this report.
  /// @param observers i-th element is the oracle id of the oracle that made the i-th observation.
  /// @param juelsPerFeeCoin exchange rate between feeCoin (e.g. ETH on Ethereum) and LINK, denominated in juels.
  /// @param configDigest configDigest of transmission.
  /// @param epochAndRound least-significant byte is the OCR protocol round number, the other bytes give the big-endian OCR protocol epoch number.
  event NewTransmission(
    uint32 indexed aggregatorRoundId,
    int192 answer,
    address transmitter,
    uint32 observationsTimestamp,
    int192[] observations,
    bytes observers,
    int192 juelsPerFeeCoin,
    bytes32 configDigest,
    uint40 epochAndRound
  );

  error CalldataLengthMismatch();
  error StaleReport();
  error UnauthorizedTransmitter();
  error ConfigDigestMismatch();
  error WrongNumberOfSignatures();
  error SignaturesOutOfRegistration();
  error SignatureError();
  error DuplicateSigner();
  error OnlyCallableByEOA();
  error ReportLengthMismatch();
  error NumObservationsOutOfBounds();
  error TooFewValuesToTrustMedian();
  error MedianIsOutOfMinMaxRange();

  /// @notice the constant-length components of the msg.data sent to transmit.
  // See the "If we wanted to call sam" example on for example reasoning
  // https://solidity.readthedocs.io/en/v0.7.2/abi-spec.html
  uint256 private constant TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT = 4 // Function selector.
    + 32 * 3 // 3 words containing reportContext.
    + 32 // Word containing start location of abiencoded report value.
    + 32 // Word containing start location of abiencoded rs value.
    + 32 // Word containing start location of abiencoded ss value.
    + 32 // RawVs value.
    + 32 // Word containing length of report.
    + 32 // Word containing length rs.
    + 32 // Word containing length of ss.
    + 0; // Placeholder.

  /// @notice decodes a serialized report into a Report struct.
  /// @param rawReport serialized report in raw format.
  /// @return report the decoded report in Report struct format.
  function _decodeReport(
    bytes memory rawReport
  ) internal pure returns (Report memory) {
    (uint32 observationsTimestamp, bytes32 rawObservers, int192[] memory observations, int192 juelsPerFeeCoin) =
      abi.decode(rawReport, (uint32, bytes32, int192[], int192));

    _requireExpectedReportLength(rawReport, observations);

    uint256 numObservations = observations.length;
    bytes memory observers = abi.encodePacked(rawObservers);

    assembly {
      // We truncate observers from length 32 to the number of observations.
      mstore(observers, numObservations)
    }

    return Report({
      observationsTimestamp: observationsTimestamp,
      observers: observers,
      observations: observations,
      juelsPerFeeCoin: juelsPerFeeCoin
    });
  }

  /// @notice make sure the calldata length matches the inputs. Otherwise, the
  /// transmitter could append an arbitrarily long (up to gas-block limit)
  /// string of 0 bytes, which we would reimburse at a rate of 16 gas/byte, but
  /// which would only cost the transmitter 4 gas/byte.
  /// @param reportLength the length of the serialized report.
  /// @param rsLength the length of the rs signatures.
  /// @param ssLength the length of the ss signatures.
  function _requireExpectedMsgDataLength(uint256 reportLength, uint256 rsLength, uint256 ssLength) private pure {
    // Calldata will never be big enough to make this overflow.
    uint256 expected = TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT + reportLength // One byte per entry in report.
      + rsLength * 32 // 32 bytes per entry in rs.
      + ssLength * 32 // 32 bytes per entry in ss.
      + 0; // Placeholder.
    if (msg.data.length != expected) {
      revert CalldataLengthMismatch();
    }
  }

  /// @inheritdoc OCR2Abstract
  function transmit(
    // reportContext consists of:
    // reportContext[0]: ConfigDigest.
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round.
    // reportContext[2]: ExtraHash.
    bytes32[3] calldata reportContext,
    bytes calldata report,
    // ECDSA signatures.
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external override {
    // Call the internal transmit function without the isSecondary flag.
    _transmit(reportContext, report, rs, ss, rawVs, false);
  }

  /// @notice secondary proxy transmit entrypoint, call the internal transmit function with the isSecondary flag.
  /// @param reportContext serialized report context containing configDigest, epoch, round and extraHash.
  /// @param report serialized report, which the signatures are signing.
  /// @param rs i-th element is the R components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param ss i-th element is the S components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param rawVs i-th element is the the V component of the i-th signature.
  function transmitSecondary(
    // reportContext consists of:
    // reportContext[0]: ConfigDigest.
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round.
    // reportContext[2]: ExtraHash.
    bytes32[3] calldata reportContext,
    bytes calldata report,
    // ECDSA signatures.
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external {
    _transmit(reportContext, report, rs, ss, rawVs, true);
  }

  /// @notice internal transmit function, is called to post a new report to the contract.
  /// @param reportContext serialized report context containing configDigest, epoch, round and extraHash.
  /// @param report serialized report, which the signatures are signing.
  /// @param rs i-th element is the R components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param ss i-th element is the S components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param rawVs i-th element is the the V component of the i-th signature.
  /// @param isSecondary wether the transmission was sent by the secondary proxy or not.
  function _transmit(
    // reportContext consists of:
    // reportContext[0]: ConfigDigest.
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round.
    // reportContext[2]: ExtraHash.
    bytes32[3] calldata reportContext,
    bytes calldata report,
    // ECDSA signatures.
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs,
    bool isSecondary
  ) internal {
    // NOTE: If the arguments to this function are changed, _requireExpectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly.

    uint256 initialGas = gasleft(); // This line must come first.

    // Validate the report data.
    _validateReport(reportContext, report.length, rs.length, ss.length);

    Report memory report_ = _decodeReport(report); // Decode the report.
    HotVars memory hotVars = s_hotVars; // Load hotVars into memory.

    if (isSecondary) {
      (bool exist, uint32 roundId) = _doesReportExist(report_);
      // In case the report exists, copy the round id and pay the transmitter.
      if (exist) {
        // In case the round has already been processed by the secondary feed.
        if (hotVars.latestSecondaryRoundId >= roundId) {
          revert StaleReport();
        }

        s_hotVars.latestSecondaryRoundId = roundId;
        emit SecondaryRoundIdUpdated(roundId);

        _payTransmitter(hotVars, report_.juelsPerFeeCoin, uint32(initialGas));
        return;
      }
    }

    // Report epoch and round.
    uint40 epochAndRound = uint40(uint256(reportContext[1]));

    // Only skip the report transmission in case the epochAndRound is equal to the latestEpochAndRound
    // and the latest sender was the secondary feed.
    if (epochAndRound != hotVars.latestEpochAndRound || !hotVars.isLatestSecondary) {
      // In case the epochAndRound is lower or equal than the latestEpochAndRound, it's a stale report
      // because it's older or has already been transmitted.
      if (epochAndRound <= hotVars.latestEpochAndRound) {
        revert StaleReport();
      }

      // Verify signatures attached to report.
      _verifySignatures(reportContext, report, rs, ss, rawVs);

      _report(hotVars, reportContext[0], epochAndRound, report_, isSecondary);
    } else {
      // If the report is the same and the latest sender was the secondary feed,
      // we're effectively unlocking the primary feed with this
      emit PrimaryFeedUnlocked(s_hotVars.latestAggregatorRoundId);
    }

    // Store if the latest report was secondary or not.
    s_hotVars.isLatestSecondary = isSecondary;
    _payTransmitter(hotVars, report_.juelsPerFeeCoin, uint32(initialGas));
  }

  /// @notice helper function to validate the report data.
  /// @param reportContext serialized report context containing configDigest, epoch, round and extraHash.
  /// @param reportLength the length of the serialized report.
  /// @param rsLength the length of the rs signatures.
  /// @param ssLength the length of the ss signatures.
  function _validateReport(
    bytes32[3] calldata reportContext,
    uint256 reportLength,
    uint256 rsLength,
    uint256 ssLength
  ) internal view {
    if (!s_transmitters[msg.sender].active) {
      revert UnauthorizedTransmitter();
    }

    if (s_latestConfigDigest != reportContext[0]) {
      revert ConfigDigestMismatch();
    }

    _requireExpectedMsgDataLength(reportLength, rsLength, ssLength);

    if (rsLength != s_hotVars.f + 1) {
      revert WrongNumberOfSignatures();
    }
    if (rsLength != ssLength) {
      revert SignaturesOutOfRegistration();
    }
  }

  /// @notice helper function to verify the report signatures.
  /// @param reportContext serialized report context containing configDigest, epoch, round and extraHash.
  /// @param report serialized report, which the signatures are signing.
  /// @param rs i-th element is the R components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param ss i-th element is the S components of the i-th signature on report. Must have at most maxNumOracles entries.
  /// @param rawVs i-th element is the the V component of the i-th signature.
  function _verifySignatures(
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) internal view {
    bytes32 h = keccak256(abi.encode(keccak256(report), reportContext));

    // i-th byte counts number of sigs made by i-th signer.
    uint256 signedCount = 0;

    Signer memory signer;
    for (uint256 i = 0; i < rs.length; ++i) {
      address signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
      signer = s_signers[signerAddress];
      if (!signer.active) {
        revert SignatureError();
      }
      unchecked {
        signedCount += 1 << (8 * signer.index);
      }
    }

    // The first byte of the mask can be 0, because we only ever have 31 oracles.
    if (signedCount & 0x0001010101010101010101010101010101010101010101010101010101010101 != signedCount) {
      revert DuplicateSigner();
    }
  }

  /// @notice details about the most recent report.
  /// @return configDigest domain separation tag for the latest report.
  /// @return epoch epoch in which the latest report was generated.
  /// @return round OCR round in which the latest report was generated.
  /// @return latestAnswer_ median value from latest report.
  /// @return latestTimestamp_ when the latest report was transmitted.
  function latestTransmissionDetails()
    external
    view
    returns (bytes32 configDigest, uint32 epoch, uint8 round, int192 latestAnswer_, uint64 latestTimestamp_)
  {
    // solhint-disable-next-line avoid-tx-origin
    if (msg.sender != tx.origin) revert OnlyCallableByEOA();
    return (
      s_latestConfigDigest,
      uint32(s_hotVars.latestEpochAndRound >> 8),
      uint8(s_hotVars.latestEpochAndRound),
      s_transmissions[s_hotVars.latestAggregatorRoundId].answer,
      s_transmissions[s_hotVars.latestAggregatorRoundId].recordedTimestamp
    );
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDigestAndEpoch()
    external
    view
    virtual
    override
    returns (bool scanLogs, bytes32 configDigest, uint32 epoch)
  {
    return (false, s_latestConfigDigest, uint32(s_hotVars.latestEpochAndRound >> 8));
  }

  /// @notice evaluate the serialized report length and compare it with the expected length.
  /// @param report serialized report, which the signatures are signing.
  /// @param observations decoded observations from the report.
  function _requireExpectedReportLength(bytes memory report, int192[] memory observations) private pure {
    uint256 expected = 32 // ObservationsTimestamp.
      + 32 // RawObservers.
      + 32 // Observations offset.
      + 32 // JuelsPerFeeCoin.
      + 32 // Observations length.
      + 32 * observations.length // Observations payload.
      + 0;
    if (report.length != expected) revert ReportLengthMismatch();
  }

  /// @notice report a new transmission and emit the necessary events.
  /// @param hotVars most common fields used in the hot path.
  /// @param configDigest digested configuration.
  /// @param epochAndRound report epoch and round.
  /// @param report decoded report in Report struct format.
  /// @param isSecondary wether the report was sent by the secondary proxy or not.
  function _report(
    HotVars memory hotVars,
    bytes32 configDigest,
    uint40 epochAndRound,
    Report memory report,
    bool isSecondary
  ) internal {
    if (report.observations.length > MAX_NUM_ORACLES) revert NumObservationsOutOfBounds();
    // Offchain logic ensures that a quorum of oracles is operating on a matching set of at least
    // 2f+1 observations. By assumption, up to f of those can be faulty, which includes being
    // malformed. Conversely, more than f observations have to be well-formed and sent on chain.
    if (report.observations.length <= hotVars.f) revert TooFewValuesToTrustMedian();

    hotVars.latestEpochAndRound = epochAndRound;

    // Get median, validate its range, store it in new aggregator round.
    int192 median = report.observations[report.observations.length / 2];
    if (i_minAnswer > median || median > i_maxAnswer) revert MedianIsOutOfMinMaxRange();

    hotVars.latestAggregatorRoundId++;
    s_transmissions[hotVars.latestAggregatorRoundId] = Transmission({
      answer: median,
      observationsTimestamp: report.observationsTimestamp,
      recordedTimestamp: uint32(block.timestamp)
    });

    // In case the sender is the secondary proxy, update the latest secondary round id.
    if (isSecondary) {
      hotVars.latestSecondaryRoundId = hotVars.latestAggregatorRoundId;
      emit SecondaryRoundIdUpdated(hotVars.latestSecondaryRoundId);
    }

    // Persist updates to hotVars.
    s_hotVars = hotVars;

    emit NewTransmission(
      hotVars.latestAggregatorRoundId,
      median,
      msg.sender,
      report.observationsTimestamp,
      report.observations,
      report.observers,
      report.juelsPerFeeCoin,
      configDigest,
      epochAndRound
    );
    // Emit these for backwards compatibility with offchain consumers
    // that only support legacy events.
    emit NewRound(
      hotVars.latestAggregatorRoundId,
      address(0x0), // Use zero address since we don't have anybody "starting" the round here.
      report.observationsTimestamp
    );
    emit AnswerUpdated(median, hotVars.latestAggregatorRoundId, block.timestamp);

    _validateAnswer(hotVars.latestAggregatorRoundId, median);
  }

  // ================================================================
  // │                   v2 AggregatorInterface                     │
  // ================================================================

  /// @notice median from the most recent report.
  /// @return answer the latest answer.
  function latestAnswer() public view virtual override returns (int256) {
    return s_transmissions[_getLatestRound()].answer;
  }

  /// @notice timestamp of block in which last report was transmitted.
  /// @return recordedTimestamp the latest recorded timestamp.
  function latestTimestamp() public view virtual override returns (uint256) {
    return s_transmissions[_getLatestRound()].recordedTimestamp;
  }

  /// @notice Aggregator round (NOT OCR round) in which last report was transmitted.
  /// @return roundId the latest round id.
  function latestRound() public view virtual override returns (uint256) {
    return _getLatestRound();
  }

  /// @notice median of report from given aggregator round (NOT OCR round).
  /// @param roundId the aggregator round of the target report.
  /// @return answer the answer of the round id.
  function getAnswer(
    uint256 roundId
  ) public view virtual override returns (int256) {
    if (roundId > _getLatestRound()) return 0;
    return s_transmissions[uint32(roundId)].answer;
  }

  /// @notice timestamp of block in which report from given aggregator round was transmitted.
  /// @param roundId aggregator round (NOT OCR round) of target report.
  /// @return recordedTimestamp the recorded timestamp of the round id.
  function getTimestamp(
    uint256 roundId
  ) public view virtual override returns (uint256) {
    if (roundId > _getLatestRound()) return 0;
    return s_transmissions[uint32(roundId)].recordedTimestamp;
  }

  // ================================================================
  // │                   v3 AggregatorInterface                     │
  // ================================================================

  error RoundNotFound();

  /// @notice amount of decimals.
  uint8 private immutable i_decimals;

  /// @notice aggregator contract version.
  uint256 internal constant VERSION = 6;

  /// @notice human readable description.
  string internal s_description;

  /// @notice get the amount of decimals.
  /// @return i_decimals amount of decimals.
  function decimals() public view virtual override returns (uint8) {
    return i_decimals;
  }

  /// @notice get the contract version.
  /// @return VERSION the contract version.
  function version() public view virtual override returns (uint256) {
    return VERSION;
  }

  /// @notice human-readable description of observable this contract is reporting on.
  /// @return s_description the contract description.
  function description() public view virtual override returns (string memory) {
    return s_description;
  }

  /// @notice details for the given aggregator round.
  /// @param roundId target aggregator round, must fit in uint32.
  /// @return roundId_ roundId.
  /// @return answer median of report from given roundId.
  /// @return startedAt timestamp of when observations were made offchain.
  /// @return updatedAt timestamp of block in which report from given roundId was transmitted.
  /// @return answeredInRound roundId.
  function getRoundData(
    uint80 roundId
  )
    public
    view
    virtual
    override
    returns (uint80 roundId_, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    if (roundId > _getLatestRound()) return (0, 0, 0, 0, 0);
    Transmission memory transmission = s_transmissions[uint32(roundId)];

    return (roundId, transmission.answer, transmission.observationsTimestamp, transmission.recordedTimestamp, roundId);
  }

  /// @notice aggregator details for the most recently transmitted report.
  /// @return roundId aggregator round of latest report (NOT OCR round).
  /// @return answer median of latest report.
  /// @return startedAt timestamp of when observations were made offchain.
  /// @return updatedAt timestamp of block containing latest report.
  /// @return answeredInRound aggregator round of latest report.
  function latestRoundData()
    public
    view
    virtual
    override
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    uint80 latestRoundId = _getLatestRound();
    Transmission memory transmission = s_transmissions[uint32(latestRoundId)];

    return (
      latestRoundId,
      transmission.answer,
      transmission.observationsTimestamp,
      transmission.recordedTimestamp,
      latestRoundId
    );
  }

  // ================================================================
  // │                  Configurable LINK Token                     │
  // ================================================================

  /// @notice emitted when the LINK token contract is set.
  /// @param oldLinkToken the address of the old LINK token contract.
  /// @param newLinkToken the address of the new LINK token contract.
  event LinkTokenSet(LinkTokenInterface indexed oldLinkToken, LinkTokenInterface indexed newLinkToken);

  error TransferRemainingFundsFailed();

  /// @notice we assume that the token contract is correct. This contract is not written
  /// to handle misbehaving ERC20 tokens!
  LinkTokenInterface internal s_linkToken;

  /// @notice sets the LINK token contract used for paying oracles.
  /// @dev this function will return early (without an error) without changing any state
  /// if linkToken equals getLinkToken().
  /// @dev this will trigger a payout so that a malicious owner cannot take from oracles
  /// what is already owed to them.
  /// @dev we assume that the token contract is correct. This contract is not written
  /// to handle misbehaving ERC20 tokens!
  /// @param linkToken the address of the LINK token contract.
  /// @param recipient remaining funds from the previous token contract are transferred
  /// here.
  function setLinkToken(LinkTokenInterface linkToken, address recipient) external onlyOwner {
    LinkTokenInterface oldLinkToken = s_linkToken;
    if (linkToken == oldLinkToken) {
      // No change, nothing to be done.
      return;
    }
    // Call balanceOf as a sanity check on whether we're talking to a token
    // contract.
    linkToken.balanceOf(address(this));
    // We break CEI here, but that's okay because we're dealing with a correct
    // token contract (by assumption).
    _payOracles();
    uint256 remainingBalance = oldLinkToken.balanceOf(address(this));
    if (!oldLinkToken.transfer(recipient, remainingBalance)) revert TransferRemainingFundsFailed();
    // solhint-disable-next-line reentrancy
    s_linkToken = linkToken;
    emit LinkTokenSet(oldLinkToken, linkToken);
  }

  /// @notice gets the LINK token contract used for paying oracles.
  /// @return linkToken the address of the LINK token contract.
  function getLinkToken() external view returns (LinkTokenInterface linkToken) {
    return s_linkToken;
  }

  // ================================================================
  // │             BillingAccessController Management               │
  // ================================================================

  /// @notice emitted when a new access-control contract is set.
  /// @param old the address prior to the current setting.
  /// @param current the address of the new access-control contract.
  event BillingAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  /// @notice controls who can change billing parameters. A billingAdmin is not able to
  /// affect any OCR protocol settings and therefore cannot tamper with the
  /// liveness or integrity of a data feed. However, a billingAdmin can set
  /// faulty billing parameters causing oracles to be underpaid, or causing them
  /// to be paid so much that further calls to setConfig, setBilling,
  /// setLinkToken will always fail due to the contract being underfunded.
  AccessControllerInterface internal s_billingAccessController;

  /// @notice internal function to set a new billingAccessController.
  /// @param billingAccessController new billingAccessController contract address.
  function _setBillingAccessController(
    AccessControllerInterface billingAccessController
  ) internal {
    AccessControllerInterface oldController = s_billingAccessController;
    if (billingAccessController != oldController) {
      s_billingAccessController = billingAccessController;
      emit BillingAccessControllerSet(oldController, billingAccessController);
    }
  }

  /// @notice sets billingAccessController.
  /// @param _billingAccessController new billingAccessController contract address.
  function setBillingAccessController(
    AccessControllerInterface _billingAccessController
  ) external onlyOwner {
    _setBillingAccessController(_billingAccessController);
  }

  /// @notice gets billingAccessController.
  /// @return s_billingAccessController address of billingAccessController contract.
  function getBillingAccessController() external view returns (AccessControllerInterface) {
    return s_billingAccessController;
  }

  // ================================================================
  // │                    Billing Configuration                     │
  // ================================================================

  /// @notice emitted when billing parameters are set.
  /// @param maximumGasPriceGwei highest gas price for which transmitter will be compensated.
  /// @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value.
  /// @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report.
  /// @param transmissionPaymentGjuels reward to transmitter of a successful report.
  /// @param accountingGas gas overhead incurred by accounting logic.
  event BillingSet(
    uint32 maximumGasPriceGwei,
    uint32 reasonableGasPriceGwei,
    uint32 observationPaymentGjuels,
    uint32 transmissionPaymentGjuels,
    uint24 accountingGas
  );

  error OnlyOwnerAndBillingAdminCanCall();

  /// @notice sets billing parameters.
  /// @dev access control provided by billingAccessController.
  /// @param maximumGasPriceGwei highest gas price for which transmitter will be compensated.
  /// @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value.
  /// @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report.
  /// @param transmissionPaymentGjuels reward to transmitter of a successful report.
  /// @param accountingGas gas overhead incurred by accounting logic.
  function setBilling(
    uint32 maximumGasPriceGwei,
    uint32 reasonableGasPriceGwei,
    uint32 observationPaymentGjuels,
    uint32 transmissionPaymentGjuels,
    uint24 accountingGas
  ) external {
    if (!(msg.sender == owner() || s_billingAccessController.hasAccess(msg.sender, msg.data))) {
      revert OnlyOwnerAndBillingAdminCanCall();
    }
    _payOracles();

    s_hotVars.maximumGasPriceGwei = maximumGasPriceGwei;
    s_hotVars.reasonableGasPriceGwei = reasonableGasPriceGwei;
    s_hotVars.observationPaymentGjuels = observationPaymentGjuels;
    s_hotVars.transmissionPaymentGjuels = transmissionPaymentGjuels;
    s_accountingGas = accountingGas;

    emit BillingSet(
      maximumGasPriceGwei, reasonableGasPriceGwei, observationPaymentGjuels, transmissionPaymentGjuels, accountingGas
    );
  }

  /// @notice gets billing parameters.
  /// @param maximumGasPriceGwei highest gas price for which transmitter will be compensated.
  /// @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value.
  /// @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report.
  /// @param transmissionPaymentGjuels reward to transmitter of a successful report.
  /// @param accountingGas gas overhead of the accounting logic.
  function getBilling()
    external
    view
    returns (
      uint32 maximumGasPriceGwei,
      uint32 reasonableGasPriceGwei,
      uint32 observationPaymentGjuels,
      uint32 transmissionPaymentGjuels,
      uint24 accountingGas
    )
  {
    return (
      s_hotVars.maximumGasPriceGwei,
      s_hotVars.reasonableGasPriceGwei,
      s_hotVars.observationPaymentGjuels,
      s_hotVars.transmissionPaymentGjuels,
      s_accountingGas
    );
  }

  // ================================================================
  // │                  Payments and Withdrawals                    │
  // ================================================================

  /// @notice emitted when an oracle has been paid LINK.
  /// @param transmitter address from which the oracle sends reports to the transmit method.
  /// @param payee address to which the payment is sent.
  /// @param amount amount of LINK sent.
  /// @param linkToken address of the LINK token contract.
  event OraclePaid(
    address indexed transmitter, address indexed payee, uint256 amount, LinkTokenInterface indexed linkToken
  );

  error OnlyPayeeCanWithdraw();
  error InsufficientFunds();
  error InsufficientBalance();

  /// @notice withdraws an oracle's payment from the contract.
  /// @param transmitter the transmitter address of the oracle.
  /// @dev must be called by oracle's payee address.
  function withdrawPayment(
    address transmitter
  ) external {
    if (msg.sender != s_payees[transmitter]) revert OnlyPayeeCanWithdraw();
    _payOracle(transmitter);
  }

  /// @notice query an oracle's payment amount, denominated in juels.
  /// @param transmitterAddress the transmitter address of the oracle.
  function owedPayment(
    address transmitterAddress
  ) public view returns (uint256) {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) return 0;
    // Safe from overflow:
    // s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index] <= 2**32.
    // s_hotVars.observationPaymentGjuels <= 2**32.
    // 1 gwei <= 2**32.
    // hence juelsAmount <= 2**96.
    uint256 juelsAmount = uint256(s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index])
      * uint256(s_hotVars.observationPaymentGjuels) * (1 gwei);
    juelsAmount += transmitter.paymentJuels;
    return juelsAmount;
  }

  /// @notice pays out transmitter's oracle balance to the corresponding payee, and zeros it out.
  /// @param transmitterAddress the transmitter address of the oracle.
  function _payOracle(
    address transmitterAddress
  ) internal {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) return;
    uint256 juelsAmount = owedPayment(transmitterAddress);
    if (juelsAmount > 0) {
      address payee = s_payees[transmitterAddress];
      // Poses no re-entrancy issues, because LINK.transfer does not yield
      // control flow.
      if (!s_linkToken.transfer(payee, juelsAmount)) {
        revert InsufficientFunds();
      }
      // solhint-disable-next-line reentrancy
      s_rewardFromAggregatorRoundId[transmitter.index] = s_hotVars.latestAggregatorRoundId;
      // solhint-disable-next-line reentrancy
      s_transmitters[transmitterAddress].paymentJuels = 0;
      emit OraclePaid(transmitterAddress, payee, juelsAmount, s_linkToken);
    }
  }

  /// @notice pays out all transmitters oracles, and zeros out their balances.
  /// It's much more gas-efficient to do this as a single operation, to avoid
  /// hitting storage too much.
  function _payOracles() internal {
    unchecked {
      LinkTokenInterface linkToken = s_linkToken;
      uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;

      uint32[MAX_NUM_ORACLES] memory rewardFromAggregatorRoundId = s_rewardFromAggregatorRoundId;
      address[] memory transmitters = s_transmittersList;

      for (uint256 transmitteridx = 0; transmitteridx < transmitters.length; transmitteridx++) {
        uint256 reimbursementAmountJuels = s_transmitters[transmitters[transmitteridx]].paymentJuels;
        s_transmitters[transmitters[transmitteridx]].paymentJuels = 0;

        uint256 obsCount = latestAggregatorRoundId - rewardFromAggregatorRoundId[transmitteridx];
        uint256 juelsAmount =
          obsCount * uint256(s_hotVars.observationPaymentGjuels) * (1 gwei) + reimbursementAmountJuels;

        if (juelsAmount > 0) {
          address payee = s_payees[transmitters[transmitteridx]];
          // Poses no re-entrancy issues, because LINK.transfer does not yield
          // control flow.
          if (!linkToken.transfer(payee, juelsAmount)) {
            revert InsufficientFunds();
          }
          rewardFromAggregatorRoundId[transmitteridx] = latestAggregatorRoundId;
          emit OraclePaid(transmitters[transmitteridx], payee, juelsAmount, linkToken);
        }
      }
      // "Zero" the accounting storage variables.
      // solhint-disable-next-line reentrancy
      s_rewardFromAggregatorRoundId = rewardFromAggregatorRoundId;
    }
  }

  /// @notice withdraw any available funds left in the contract, up to amount, after accounting for the funds due to participants in past reports.
  /// @dev access control provided by billingAccessController.
  /// @param recipient address to send funds to.
  /// @param amount maximum amount to withdraw, denominated in LINK-wei.
  function withdrawFunds(address recipient, uint256 amount) external {
    if (msg.sender != owner() && !s_billingAccessController.hasAccess(msg.sender, msg.data)) {
      revert OnlyOwnerAndBillingAdminCanCall();
    }
    uint256 linkDue = _totalLinkDue();
    uint256 linkBalance = s_linkToken.balanceOf(address(this));
    if (linkBalance < linkDue) {
      revert InsufficientBalance();
    }
    if (!s_linkToken.transfer(recipient, _min(linkBalance - linkDue, amount))) {
      revert InsufficientFunds();
    }
  }

  /// @notice total LINK due to participants in past reports (denominated in Juels).
  /// @return linkDue total LINK due.
  function _totalLinkDue() internal view returns (uint256 linkDue) {
    // Argument for overflow safety: We do all computations in
    // uint256s. The inputs to linkDue are:
    // - the <= 31 observation rewards each of which has less than
    //   64 bits (32 bits for observationPaymentGjuels, 32 bits
    //   for wei/gwei conversion). Hence 69 bits are sufficient for this part.
    // - the <= 31 gas reimbursements, each of which consists of at most 96
    //   bits. Hence 101 bits are sufficient for this part.
    // So we never need more than 102 bits.

    address[] memory transmitters = s_transmittersList;
    uint256 n = transmitters.length;

    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
    uint32[MAX_NUM_ORACLES] memory rewardFromAggregatorRoundId = s_rewardFromAggregatorRoundId;
    for (uint256 i = 0; i < n; ++i) {
      linkDue += latestAggregatorRoundId - rewardFromAggregatorRoundId[i];
    }
    // Convert observationPaymentGjuels to uint256, or this overflows!
    linkDue *= uint256(s_hotVars.observationPaymentGjuels) * (1 gwei);
    for (uint256 i = 0; i < n; ++i) {
      linkDue += uint256(s_transmitters[transmitters[i]].paymentJuels);
    }

    return linkDue;
  }

  /// @notice allows oracles to check that sufficient LINK balance is available.
  /// @return availableBalance LINK available on this contract, after accounting for outstanding obligations, can become negative.
  function linkAvailableForPayment() external view returns (int256 availableBalance) {
    // There are at most one billion LINK, so this cast is safe.
    int256 balance = int256(s_linkToken.balanceOf(address(this)));
    // According to the argument in the definition of _totalLinkDue,
    // _totalLinkDue is never greater than 2**102, so this cast is safe.
    int256 due = int256(_totalLinkDue());
    // Safe from overflow according to above sizes.
    return int256(balance) - int256(due);
  }

  /// @notice number of observations oracle is due to be reimbursed for.
  /// @param transmitterAddress address used by oracle for signing or transmitting reports.
  /// @return observations difference between the latest oracle reimbursement round id and the latest hotVars round id.
  function oracleObservationCount(
    address transmitterAddress
  ) external view returns (uint32) {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) return 0;
    return s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index];
  }

  // ================================================================
  // │                     Transmitter Payment                      │
  // ================================================================

  error LeftGasCannotExceedInitialGas();

  /// @notice gas price at which the transmitter should be reimbursed, in gwei/gas.
  /// @param txGasPriceGwei transaction gas price in ETH-gwei units.
  /// @param reasonableGasPriceGwei reasonable gas price in ETH-gwei units.
  /// @param maximumGasPriceGwei maximum gas price in ETH-gwei units.
  /// @return gasPrice resulting gas price to reimburse.
  function _reimbursementGasPriceGwei(
    uint256 txGasPriceGwei,
    uint256 reasonableGasPriceGwei,
    uint256 maximumGasPriceGwei
  ) internal pure returns (uint256) {
    // This happens on the path for transmissions. We'd rather pay out
    // a wrong reward than risk a liveness failure due to a revert.
    unchecked {
      // Reward the transmitter for choosing an efficient gas price: if they manage
      // to come in lower than considered reasonable, give them half the savings.
      uint256 gasPriceGwei = txGasPriceGwei;
      if (txGasPriceGwei < reasonableGasPriceGwei) {
        // Give transmitter half the savings for coming in under the reasonable gas price.
        gasPriceGwei += (reasonableGasPriceGwei - txGasPriceGwei) / 2;
      }
      // Don't reimburse a gas price higher than maximumGasPriceGwei.
      return _min(gasPriceGwei, maximumGasPriceGwei);
    }
  }

  /// @notice gas reimbursement due the transmitter, in wei.
  /// @param initialGas initial remaining gas.
  /// @param gasPriceGwei gas price in ETH-gwei units.
  /// @param callDataGas calldata gas cost.
  /// @param accountingGas overhead incurred by accounting logic.
  /// @param leftGas actual remaining gas.
  /// @return fullGasCostWei final calculated gas cost in wei.
  function _transmitterGasCostWei(
    uint256 initialGas,
    uint256 gasPriceGwei,
    uint256 callDataGas,
    uint256 accountingGas,
    uint256 leftGas
  ) internal pure returns (uint256) {
    // This happens on the path for transmissions. We'd rather pay out
    // a wrong reward than risk a liveness failure due to a revert.
    unchecked {
      if (initialGas < leftGas) revert LeftGasCannotExceedInitialGas();
      uint256 usedGas = initialGas - leftGas // Observed gas usage.
        + callDataGas + accountingGas; // Estimated gas usage.
      uint256 fullGasCostWei = usedGas * gasPriceGwei * (1 gwei);
      return fullGasCostWei;
    }
  }

  /// @notice internal function to pay the transmitter on the path for transmissions. Note: We'd rather pay out
  /// a wrong reward than risk a liveness failure due to a revert.
  /// @param hotVars most common fields used in the hot path.
  /// @param juelsPerFeeCoin exchange rate between feeCoin (e.g. ETH on Ethereum) and LINK, denominated in juels.
  /// @param initialGas initial remaining gas.
  function _payTransmitter(HotVars memory hotVars, int192 juelsPerFeeCoin, uint32 initialGas) internal virtual {
    unchecked {
      // We can't deal with negative juelsPerFeeCoin, better to just not pay.
      if (juelsPerFeeCoin < 0) {
        return;
      }

      // Reimburse transmitter of the report for gas usage.
      uint256 gasPriceGwei = _reimbursementGasPriceGwei(
        tx.gasprice / (1 gwei), // Convert to ETH-gwei units.
        hotVars.reasonableGasPriceGwei,
        hotVars.maximumGasPriceGwei
      );
      // The following is only an upper bound, as it ignores the cheaper cost for
      // 0 bytes. Safe from overflow, because calldata just isn't that long.
      uint256 callDataGasCost = 16 * msg.data.length;
      uint256 gasLeft = gasleft();
      uint256 gasCostEthWei =
        _transmitterGasCostWei(uint256(initialGas), gasPriceGwei, callDataGasCost, s_accountingGas, gasLeft);

      // Even if we assume absurdly large values, this still does not overflow, with:
      // - usedGas <= 1'000'000 gas <= 2**20 gas.
      // - weiPerGas <= 1'000'000 gwei <= 2**50 wei.
      // - hence gasCostEthWei <= 2**70.
      // - juelsPerFeeCoin <= 2**96 (more than the entire supply).
      // We still fit into 166 bits.
      uint256 gasCostJuels = (gasCostEthWei * uint192(juelsPerFeeCoin)) / 1e18;

      uint96 oldTransmitterPaymentJuels = s_transmitters[msg.sender].paymentJuels;
      uint96 newTransmitterPaymentJuels = uint96(
        uint256(oldTransmitterPaymentJuels) + gasCostJuels + uint256(hotVars.transmissionPaymentGjuels) * (1 gwei)
      );

      // Overflow *should* never happen, but if it does, let's not persist it.
      if (newTransmitterPaymentJuels < oldTransmitterPaymentJuels) {
        return;
      }
      s_transmitters[msg.sender].paymentJuels = newTransmitterPaymentJuels;
    }
  }

  // ================================================================
  // │                       Payee Management                       │
  // ================================================================

  /// @notice emitted when a transfer of an oracle's payee address has been initiated.
  /// @param transmitter address from which the oracle sends reports to the transmit method.
  /// @param current the payee address for the oracle, prior to this setting.
  /// @param proposed the proposed new payee address for the oracle.
  event PayeeshipTransferRequested(address indexed transmitter, address indexed current, address indexed proposed);

  /// @notice emitted when a transfer of an oracle's payee address has been completed.
  /// @param transmitter address from which the oracle sends reports to the transmit method.
  /// @param previous the previous payee address for the oracle.
  /// @param current the payee address for the oracle, prior to this setting.
  event PayeeshipTransferred(address indexed transmitter, address indexed previous, address indexed current);

  error TransmittersSizeNotEqualPayeeSize();
  error PayeeAlreadySet();
  error OnlyCurrentPayeeCanUpdate();
  error CannotTransferPayeeToSelf();
  error OnlyProposedPayeesCanAccept();

  /// @notice addresses at which oracles want to receive payments, by transmitter address.
  mapping(address transmitter => address paymentAddress) internal s_payees;

  /// @notice payee addresses which must be approved by the owner.
  mapping(address transmitter => address paymentAddress) internal s_proposedPayees;

  /// @notice sets the payees for transmitting addresses.
  /// @dev cannot be used to change payee addresses, only to initially populate them.
  /// @param transmitters addresses oracles use to transmit the reports.
  /// @param payees addresses of payees corresponding to list of transmitters.
  function setPayees(address[] calldata transmitters, address[] calldata payees) external onlyOwner {
    if (transmitters.length != payees.length) revert TransmittersSizeNotEqualPayeeSize();

    for (uint256 i = 0; i < transmitters.length; ++i) {
      address transmitter = transmitters[i];
      address payee = payees[i];
      address currentPayee = s_payees[transmitter];
      bool zeroedOut = currentPayee == address(0);
      if (!zeroedOut && currentPayee != payee) revert PayeeAlreadySet();
      s_payees[transmitter] = payee;

      if (currentPayee != payee) {
        emit PayeeshipTransferred(transmitter, currentPayee, payee);
      }
    }
  }

  /// @notice first step of payeeship transfer (safe transfer pattern).
  /// @dev can only be called by payee addresses.
  /// @param transmitter transmitter address of oracle whose payee is changing.
  /// @param proposed new payee address.
  function transferPayeeship(address transmitter, address proposed) external {
    if (msg.sender != s_payees[transmitter]) {
      revert OnlyCurrentPayeeCanUpdate();
    }
    if (msg.sender == proposed) revert CannotTransferPayeeToSelf();

    address previousProposed = s_proposedPayees[transmitter];
    s_proposedPayees[transmitter] = proposed;

    if (previousProposed != proposed) {
      emit PayeeshipTransferRequested(transmitter, msg.sender, proposed);
    }
  }

  /// @notice second step of payeeship transfer (safe transfer pattern).
  /// @dev can only be called by proposed new payee address.
  /// @param transmitter transmitter address of oracle whose payee is changing.
  function acceptPayeeship(
    address transmitter
  ) external {
    if (msg.sender != s_proposedPayees[transmitter]) revert OnlyProposedPayeesCanAccept();

    address currentPayee = s_payees[transmitter];
    s_payees[transmitter] = msg.sender;
    s_proposedPayees[transmitter] = address(0);

    emit PayeeshipTransferred(transmitter, currentPayee, msg.sender);
  }

  // ================================================================
  // │                       Helper Functions                       │
  // ================================================================

  /// @notice helper function to compare two numbers and return the smallest number.
  /// @param a first number.
  /// @param b second number.
  /// @return result smallest number.
  function _min(uint256 a, uint256 b) internal pure returns (uint256) {
    unchecked {
      if (a < b) return a;
      return b;
    }
  }
}
