// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/AccessControllerInterface.sol";
import "./interfaces/AggregatorV2V3Interface.sol";
import "./interfaces/AggregatorValidatorInterface.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/TypeAndVersionInterface.sol";
import "./OCR2Abstract.sol";
import "./OwnerIsCreator.sol";


/**
 * @notice OCR2Aggregator for numerical data with billing support.

 * @dev
 * If you read or change this, be sure to read or adjust the comments. They
 * track the units of the values under consideration, and are crucial to
 * the readability of the operations it specifies.

 * @notice
 * Billing Trust Model:

 * Nothing in this contract prevents a billing admin from setting insane
 * values for the billing parameters in setBilling. Oracles
 * participating in this contract should regularly check that the
 * parameters make sense. Similarly, the outstanding obligations of this
 * contract to the oracles can exceed the funds held by the contract.
 * Oracles participating in this contract should regularly check that it
 * holds sufficient funds and stop interacting with it if funding runs
 * out.

 * This still leaves oracles with some risk due to TOCTOU issues.
 * However, since the sums involved are pretty small (Ethereum
 * transactions aren't that expensive in the end) and an oracle would
 * likely stop participating in a contract it repeatedly lost money on,
 * this risk is deemed acceptable. Oracles should also regularly
 * withdraw any funds in the contract to prevent issues where the
 * contract becomes underfunded at a later time, and different oracles
 * are competing for the left-over funds.

 * Finally, note that any change to the set of oracles or to the billing
 * parameters will trigger payout of all oracles first (using the old
 * parameters), a billing admin cannot take away funds that are already
 * marked for payment.
 */
contract OCR2Aggregator is OCR2Abstract, OwnerIsCreator, AggregatorV2V3Interface {
  // This contract is divided into sections. Each section defines a set of
  // variables, events, and functions that belong together.

  /***************************************************************************
   * Section: Variables used in multiple other sections
   **************************************************************************/

  struct Transmitter {
    bool active;

    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;

    // juels-denominated payment for transmitters, covering gas costs incurred
    // by the transmitter plus additional rewards. The entire LINK supply (1e9
    // LINK = 1e27 Juels) will always fit into a uint96.
    uint96 paymentJuels;
  }
  mapping (address /* transmitter address */ => Transmitter) internal s_transmitters;

  struct Signer {
    bool active;

    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;
  }
  mapping (address /* signer address */ => Signer) internal s_signers;

  // s_signersList contains the signing address of each oracle
  address[] internal s_signersList;

  // s_transmittersList contains the transmission address of each oracle,
  // i.e. the address the oracle actually sends transactions to the contract from
  address[] internal s_transmittersList;

  // We assume that all oracles contribute observations to all rounds. this
  // variable tracks (per-oracle) from what round an oracle should be rewarded,
  // i.e. the oracle gets (latestAggregatorRoundId -
  // rewardFromAggregatorRoundId) * reward
  uint32[maxNumOracles] internal s_rewardFromAggregatorRoundId;

  bytes32 s_latestConfigDigest;

  // Storing these fields used on the hot path in a HotVars variable reduces the
  // retrieval of all of them to a single SLOAD.
  struct HotVars {
    // maximum number of faulty oracles
    uint8 f;

    // epoch and round from OCR protocol.
    // 32 most sig bits for epoch, 8 least sig bits for round
    uint40 latestEpochAndRound;

    // Chainlink Aggregators expose a roundId to consumers. The offchain reporting
    // protocol does not use this id anywhere. We increment it whenever a new
    // transmission is made to provide callers with contiguous ids for successive
    // reports.
    uint32 latestAggregatorRoundId;

    // Highest compensated gas price, in gwei uints
    uint32 maximumGasPriceGwei;

    // If gas price is less (in gwei units), transmitter gets half the savings
    uint32 reasonableGasPriceGwei;

    // Fixed LINK reward for each observer
    uint32 observationPaymentGjuels;

    // Fixed reward for transmitter
    uint32 transmissionPaymentGjuels;

    // Overhead incurred by accounting logic
    uint24 accountingGas;
  }
  HotVars internal s_hotVars;

  // Transmission records the median answer from the transmit transaction at
  // time timestamp
  struct Transmission {
    int192 answer; // 192 bits ought to be enough for anyone
    uint32 observationsTimestamp; // when were observations made offchain
    uint32 transmissionTimestamp; // when was report received onchain
  }
  mapping(uint32 /* aggregator round ID */ => Transmission) internal s_transmissions;

  // Lowest answer the system is allowed to report in response to transmissions
  int192 immutable public minAnswer;
  // Highest answer the system is allowed to report in response to transmissions
  int192 immutable public maxAnswer;

  /***************************************************************************
   * Section: Constructor
   **************************************************************************/

  /**
   * @param link address of the LINK contract
   * @param minAnswer_ lowest answer the median of a report is allowed to be
   * @param maxAnswer_ highest answer the median of a report is allowed to be
   * @param requesterAccessController access controller for requesting new rounds
   * @param decimals_ answers are stored in fixed-point format, with this many digits of precision
   * @param description_ short human-readable description of observable this contract's answers pertain to
   */
  constructor(
    LinkTokenInterface link,
    int192 minAnswer_,
    int192 maxAnswer_,
    AccessControllerInterface billingAccessController,
    AccessControllerInterface requesterAccessController,
    uint8 decimals_,
    string memory description_
  ) {
    s_linkToken = link;
    emit LinkTokenSet(LinkTokenInterface(address(0)), link);
    _setBillingAccessController(billingAccessController);

    decimals = decimals_;
    s_description = description_;
    setRequesterAccessController(requesterAccessController);
    setValidatorConfig(AggregatorValidatorInterface(address(0x0)), 0);
    minAnswer = minAnswer_;
    maxAnswer = maxAnswer_;
  }


  /***************************************************************************
   * Section: OCR2Abstract Configuration
   **************************************************************************/

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest to prevent replay attacks.
  uint32 internal s_configCount;

  // makes it easier for offchain systems to extract config from logs
  uint32 internal s_latestConfigBlockNumber;

  // left as a function so this check can be disabled in derived contracts
  function _requirePositiveF (
    uint256 f
  )
    internal
    pure
    virtual
  {
    require(0 < f, "f must be positive");
  }

  struct SetConfigArgs {
    address[] signers;
    address[] transmitters;
    uint8 f;
    bytes onchainConfig;
    uint64 offchainConfigVersion;
    bytes offchainConfig;
  }

  /// @inheritdoc OCR2Abstract
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  )
    external
    override
    onlyOwner()
  {
    require(signers.length <= maxNumOracles, "too many oracles");
    require(signers.length == transmitters.length, "oracle length mismatch");
    require(3*f < signers.length, "faulty-oracle f too high");
    _requirePositiveF(f);
    require(onchainConfig.length == 0, "onchainConfig must be empty");

    SetConfigArgs memory args = SetConfigArgs({
      signers: signers,
      transmitters: transmitters,
      f: f,
      onchainConfig: abi.encodePacked(uint8(1) /*version*/, minAnswer, maxAnswer),
      offchainConfigVersion: offchainConfigVersion,
      offchainConfig: offchainConfig
    });

    s_hotVars.latestEpochAndRound = 0;
    _payOracles();

    // remove any old signer/transmitter addresses
    uint256 oldLength = s_signersList.length;
    for (uint256 i = 0; i < oldLength; i++) {
      address signer = s_signersList[i];
      address transmitter = s_transmittersList[i];
      delete s_signers[signer];
      delete s_transmitters[transmitter];
    }
    delete s_signersList;
    delete s_transmittersList;

    // add new signer/transmitter addresses
    for (uint i = 0; i < args.signers.length; i++) {
      require(
        !s_signers[args.signers[i]].active,
        "repeated signer address"
      );
      s_signers[args.signers[i]] = Signer({
        active: true,
        index: uint8(i)
      });
      require(
        !s_transmitters[args.transmitters[i]].active,
        "repeated transmitter address"
      );
      s_transmitters[args.transmitters[i]] = Transmitter({
        active: true,
        index: uint8(i),
        paymentJuels: 0
      });
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
    for (uint256 i = 0; i < args.signers.length; i++) {
      s_rewardFromAggregatorRoundId[i] = latestAggregatorRoundId;
    }
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDetails()
    external
    override
    view
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest
    )
  {
    return (s_configCount, s_latestConfigBlockNumber, s_latestConfigDigest);
  }

  /**
   * @return list of addresses permitted to transmit reports to this contract

   * @dev The list will match the order used to specify the transmitter during setConfig
   */
  function getTransmitters()
    external
    view
    returns(address[] memory)
  {
    return s_transmittersList;
  }

  /***************************************************************************
   * Section: Onchain Validation
   **************************************************************************/

  // Configuration for validator
  struct ValidatorConfig {
    AggregatorValidatorInterface validator;
    uint32 gasLimit;
  }
  ValidatorConfig private s_validatorConfig;

  /**
   * @notice indicates that the validator configuration has been set
   * @param previousValidator previous validator contract
   * @param previousGasLimit previous gas limit for validate calls
   * @param currentValidator current validator contract
   * @param currentGasLimit current gas limit for validate calls
   */
  event ValidatorConfigSet(
    AggregatorValidatorInterface indexed previousValidator,
    uint32 previousGasLimit,
    AggregatorValidatorInterface indexed currentValidator,
    uint32 currentGasLimit
  );

  /**
   * @notice validator configuration
   * @return validator validator contract
   * @return gasLimit gas limit for validate calls
   */
  function getValidatorConfig()
    external
    view
    returns (AggregatorValidatorInterface validator, uint32 gasLimit)
  {
    ValidatorConfig memory vc = s_validatorConfig;
    return (vc.validator, vc.gasLimit);
  }

  /**
   * @notice sets validator configuration
   * @dev set newValidator to 0x0 to disable validate calls
   * @param newValidator address of the new validator contract
   * @param newGasLimit new gas limit for validate calls
   */
  function setValidatorConfig(
    AggregatorValidatorInterface newValidator,
    uint32 newGasLimit
  )
    public
    onlyOwner()
  {
    ValidatorConfig memory previous = s_validatorConfig;

    if (previous.validator != newValidator || previous.gasLimit != newGasLimit) {
      s_validatorConfig = ValidatorConfig({
        validator: newValidator,
        gasLimit: newGasLimit
      });

      emit ValidatorConfigSet(previous.validator, previous.gasLimit, newValidator, newGasLimit);
    }
  }

  function _validateAnswer(
    uint32 aggregatorRoundId,
    int256 answer
  )
    private
  {
    ValidatorConfig memory vc = s_validatorConfig;

    if (address(vc.validator) == address(0)) {
      return;
    }

    uint32 prevAggregatorRoundId = aggregatorRoundId - 1;
    int256 prevAggregatorRoundAnswer = s_transmissions[prevAggregatorRoundId].answer;
    require(
      _callWithExactGasEvenIfTargetIsNoContract(
        vc.gasLimit,
        address(vc.validator),
        abi.encodeWithSignature(
          "validate(uint256,int256,uint256,int256)",
          uint256(prevAggregatorRoundId),
          prevAggregatorRoundAnswer,
          uint256(aggregatorRoundId),
          answer
        )
      ),
      "insufficient gas"
    );
  }

  uint256 private constant CALL_WITH_EXACT_GAS_CUSHION = 5_000;

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * or reverts if at least gasAmount gas is not available.
   */
  function _callWithExactGasEvenIfTargetIsNoContract(
    uint256 gasAmount,
    address target,
    bytes memory data
  )
    private
    returns (bool sufficientGas)
  {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      let g := gas()
      // Compute g -= CALL_WITH_EXACT_GAS_CUSHION and check for underflow. We
      // need the cushion since the logic following the above call to gas also
      // costs gas which we cannot account for exactly. So cushion is a
      // conservative upper bound for the cost of this logic.
      if iszero(lt(g, CALL_WITH_EXACT_GAS_CUSHION)) {
        g := sub(g, CALL_WITH_EXACT_GAS_CUSHION)
        // If g - g//64 <= gasAmount, we don't have enough gas. (We subtract g//64
        // because of EIP-150.)
        if gt(sub(g, div(g, 64)), gasAmount) {
          // Call and ignore success/return data. Note that we did not check
          // whether a contract actually exists at the target address.
          pop(call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0))
          sufficientGas := true
        }
      }
    }
  }

  /***************************************************************************
   * Section: RequestNewRound
   **************************************************************************/

  AccessControllerInterface internal s_requesterAccessController;

  /**
   * @notice emitted when a new requester access controller contract is set
   * @param old the address prior to the current setting
   * @param current the address of the new access controller contract
   */
  event RequesterAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  /**
   * @notice emitted to immediately request a new round
   * @param requester the address of the requester
   * @param configDigest the latest transmission's configDigest
   * @param epoch the latest transmission's epoch
   * @param round the latest transmission's round
   */
  event RoundRequested(address indexed requester, bytes32 configDigest, uint32 epoch, uint8 round);

  /**
   * @notice address of the requester access controller contract
   * @return requester access controller address
   */
  function getRequesterAccessController()
    external
    view
    returns (AccessControllerInterface)
  {
    return s_requesterAccessController;
  }

  /**
   * @notice sets the requester access controller
   * @param requesterAccessController designates the address of the new requester access controller
   */
  function setRequesterAccessController(AccessControllerInterface requesterAccessController)
    public
    onlyOwner()
  {
    AccessControllerInterface oldController = s_requesterAccessController;
    if (requesterAccessController != oldController) {
      s_requesterAccessController = AccessControllerInterface(requesterAccessController);
      emit RequesterAccessControllerSet(oldController, requesterAccessController);
    }
  }

  /**
   * @notice immediately requests a new round
   * @return the aggregatorRoundId of the next round. Note: The report for this round may have been
   * transmitted (but not yet mined) *before* requestNewRound() was even called. There is *no*
   * guarantee of causality between the request and the report at aggregatorRoundId.
   */
  function requestNewRound() external returns (uint80) {
    require(msg.sender == owner() || s_requesterAccessController.hasAccess(msg.sender, msg.data),
      "Only owner&requester can call");

    uint40 latestEpochAndRound = s_hotVars.latestEpochAndRound;
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;

    emit RoundRequested(
      msg.sender,
      s_latestConfigDigest,
      uint32(latestEpochAndRound >> 8),
      uint8(latestEpochAndRound)
    );
    return latestAggregatorRoundId + 1;
  }

  /***************************************************************************
   * Section: Transmission
   **************************************************************************/

  /**
   * @notice indicates that a new report was transmitted
   * @param aggregatorRoundId the round to which this report was assigned
   * @param answer median of the observations attached to this report
   * @param transmitter address from which the report was transmitted
   * @param observationsTimestamp when were observations made offchain
   * @param observations observations transmitted with this report
   * @param observers i-th element is the oracle id of the oracle that made the i-th observation
   * @param juelsPerFeeCoin exchange rate between feeCoin (e.g. ETH on Ethereum) and LINK, denominated in juels
   * @param configDigest configDigest of transmission
   * @param epochAndRound least-significant byte is the OCR protocol round number, the other bytes give the big-endian OCR protocol epoch number
   */
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

  // Used to relieve stack pressure in transmit
  struct Report {
    uint32 observationsTimestamp;
    bytes observers; // ith element is the index of the ith observer
    int192[] observations; // ith element is the ith observation
    int192 juelsPerFeeCoin;
  }

  // _decodeReport decodes a serialized report into a Report struct
  function _decodeReport(bytes memory rawReport)
    internal
    pure
    returns (
      Report memory
    )
  {
    uint32 observationsTimestamp;
    bytes32 rawObservers;
    int192[] memory observations;
    int192 juelsPerFeeCoin;
    (observationsTimestamp, rawObservers, observations, juelsPerFeeCoin) = abi.decode(rawReport, (uint32, bytes32, int192[], int192));

    _requireExpectedReportLength(rawReport, observations);

    uint256 numObservations = observations.length;
    bytes memory observers = abi.encodePacked(rawObservers);
    assembly {
      // we truncate observers from length 32 to the number of observations
      mstore(observers, numObservations)
    }

    return Report({
      observationsTimestamp: observationsTimestamp,
      observers: observers,
      observations: observations,
      juelsPerFeeCoin: juelsPerFeeCoin
    });
  }

  // The constant-length components of the msg.data sent to transmit.
  // See the "If we wanted to call sam" example on for example reasoning
  // https://solidity.readthedocs.io/en/v0.7.2/abi-spec.html
  uint256 private constant TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT =
    4 + // function selector
    32 * 3 + // 3 words containing reportContext
    32 + // word containing start location of abiencoded report value
    32 + // word containing location start of abiencoded rs value
    32 + // word containing start location of abiencoded ss value
    32 + // rawVs value
    32 + // word containing length of report
    32 + // word containing length rs
    32 + // word containing length of ss
    0; // placeholder

  // Make sure the calldata length matches the inputs. Otherwise, the
  // transmitter could append an arbitrarily long (up to gas-block limit)
  // string of 0 bytes, which we would reimburse at a rate of 16 gas/byte, but
  // which would only cost the transmitter 4 gas/byte.
  function _requireExpectedMsgDataLength(
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss
  )
    private
    pure
  {
    // calldata will never be big enough to make this overflow
    uint256 expected = TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT +
      report.length + // one byte pure entry in report
      rs.length * 32 + // 32 bytes per entry in rs
      ss.length * 32 + // 32 bytes per entry in ss
      0; // placeholder
    require(msg.data.length == expected, "calldata length mismatch");
  }

  /// @inheritdoc OCR2Abstract
  function transmit(
    // reportContext consists of:
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[2]: ExtraHash
    bytes32[3] calldata reportContext,
    bytes calldata report,
    // ECDSA signatures
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  )
    external
    override
  {
    // NOTE: If the arguments to this function are changed, _requireExpectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly

    uint256 initialGas = gasleft(); // This line must come first

    HotVars memory hotVars = s_hotVars;

    uint40 epochAndRound = uint40(uint256(reportContext[1]));

    require(hotVars.latestEpochAndRound < epochAndRound, "stale report");

    require(s_transmitters[msg.sender].active, "unauthorized transmitter");

    require(s_latestConfigDigest == reportContext[0], "configDigest mismatch");

    _requireExpectedMsgDataLength(report, rs, ss);

    require(rs.length == hotVars.f + 1, "wrong number of signatures");
    require(rs.length == ss.length, "signatures out of registration");

    // Verify signatures attached to report
    {
      bytes32 h = keccak256(abi.encode(keccak256(report), reportContext));

      // i-th byte counts number of sigs made by i-th signer
      uint256 signedCount = 0;

      Signer memory signer;
      for (uint i = 0; i < rs.length; i++) {
        address signerAddress = ecrecover(h, uint8(rawVs[i])+27, rs[i], ss[i]);
        signer = s_signers[signerAddress];
        require(signer.active, "signature error");
        unchecked{
          signedCount += 1 << (8 * signer.index);
        }
      }

      // The first byte of the mask can be 0, because we only ever have 31 oracles
      require(signedCount & 0x0001010101010101010101010101010101010101010101010101010101010101 == signedCount, "duplicate signer");
    }

    int192 juelsPerFeeCoin = _report(hotVars, reportContext[0], epochAndRound, report);

    _payTransmitter(hotVars, juelsPerFeeCoin, uint32(initialGas), msg.sender);
  }

  /**
   * @notice details about the most recent report
   * @return configDigest domain separation tag for the latest report
   * @return epoch epoch in which the latest report was generated
   * @return round OCR round in which the latest report was generated
   * @return latestAnswer_ median value from latest report
   * @return latestTimestamp_ when the latest report was transmitted
   */
  function latestTransmissionDetails()
    external
    view
    returns (
      bytes32 configDigest,
      uint32 epoch,
      uint8 round,
      int192 latestAnswer_,
      uint64 latestTimestamp_
    )
  {
    require(msg.sender == tx.origin, "Only callable by EOA");
    return (
      s_latestConfigDigest,
      uint32(s_hotVars.latestEpochAndRound >> 8),
      uint8(s_hotVars.latestEpochAndRound),
      s_transmissions[s_hotVars.latestAggregatorRoundId].answer,
      s_transmissions[s_hotVars.latestAggregatorRoundId].transmissionTimestamp
    );
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDigestAndEpoch()
    external
    override
    view
    virtual
    returns(
      bool scanLogs,
      bytes32 configDigest,
      uint32 epoch
    )
  {
    return (false, s_latestConfigDigest, uint32(s_hotVars.latestEpochAndRound >> 8));
  }

  function _requireExpectedReportLength(
    bytes memory report,
    int192[] memory observations
  )
    private
    pure
  {
    uint256 expected =
      32 + // observationsTimestamp
      32 + // rawObservers
      32 + // observations offset
      32 + // juelsPerFeeCoin
      32 + // observations length
      32 * observations.length + // observations payload
      0;
    require(report.length == expected, "report length mismatch");
  }

  function _report(
    HotVars memory hotVars,
    bytes32 configDigest,
    uint40 epochAndRound,
    bytes memory rawReport
  )
    internal
    returns (int192 juelsPerFeeCoin)
  {
    Report memory report = _decodeReport(rawReport);

    require(report.observations.length <= maxNumOracles, "num observations out of bounds");
    require(hotVars.f < report.observations.length, "too few values to trust median");

    hotVars.latestEpochAndRound = epochAndRound;

    // get median, validate its range, store it in new aggregator round
    int192 median = report.observations[report.observations.length/2];
    require(minAnswer <= median && median <= maxAnswer, "median is out of min-max range");
    hotVars.latestAggregatorRoundId++;
    s_transmissions[hotVars.latestAggregatorRoundId] =
      Transmission({
        answer: median,
        observationsTimestamp: report.observationsTimestamp,
        transmissionTimestamp: uint32(block.timestamp)
      });

    // persist updates to hotVars
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
    // that only support legacy events
    emit NewRound(
      hotVars.latestAggregatorRoundId,
      address(0x0), // use zero address since we don't have anybody "starting" the round here
      report.observationsTimestamp
    );
    emit AnswerUpdated(
      median,
      hotVars.latestAggregatorRoundId,
      block.timestamp
    );

    _validateAnswer(hotVars.latestAggregatorRoundId, median);

    return report.juelsPerFeeCoin;
  }

  /***************************************************************************
   * Section: v2 AggregatorInterface
   **************************************************************************/

  /**
   * @notice median from the most recent report
   */
  function latestAnswer()
    public
    override
    view
    virtual
    returns (int256)
  {
    return s_transmissions[s_hotVars.latestAggregatorRoundId].answer;
  }

  /**
   * @notice timestamp of block in which last report was transmitted
   */
  function latestTimestamp()
    public
    override
    view
    virtual
    returns (uint256)
  {
    return s_transmissions[s_hotVars.latestAggregatorRoundId].transmissionTimestamp;
  }

  /**
   * @notice Aggregator round (NOT OCR round) in which last report was transmitted
   */
  function latestRound()
    public
    override
    view
    virtual
    returns (uint256)
  {
    return s_hotVars.latestAggregatorRoundId;
  }

  /**
   * @notice median of report from given aggregator round (NOT OCR round)
   * @param roundId the aggregator round of the target report
   */
  function getAnswer(uint256 roundId)
    public
    override
    view
    virtual
    returns (int256)
  {
    if (roundId > 0xFFFFFFFF) { return 0; }
    return s_transmissions[uint32(roundId)].answer;
  }

  /**
   * @notice timestamp of block in which report from given aggregator round was transmitted
   * @param roundId aggregator round (NOT OCR round) of target report
   */
  function getTimestamp(uint256 roundId)
    public
    override
    view
    virtual
    returns (uint256)
  {
    if (roundId > 0xFFFFFFFF) { return 0; }
    return s_transmissions[uint32(roundId)].transmissionTimestamp;
  }

  /***************************************************************************
   * Section: v3 AggregatorInterface
   **************************************************************************/

  /**
   * @return answers are stored in fixed-point format, with this many digits of precision
   */
  uint8 immutable public override decimals;

  /**
   * @notice aggregator contract version
   */
  uint256 constant public override version = 6;

  string internal s_description;

  /**
   * @notice human-readable description of observable this contract is reporting on
   */
  function description()
    public
    override
    view
    virtual
    returns (string memory)
  {
    return s_description;
  }

  /**
   * @notice details for the given aggregator round
   * @param roundId target aggregator round (NOT OCR round). Must fit in uint32
   * @return roundId_ roundId
   * @return answer median of report from given roundId
   * @return startedAt timestamp of when observations were made offchain
   * @return updatedAt timestamp of block in which report from given roundId was transmitted
   * @return answeredInRound roundId
   */
  function getRoundData(uint80 roundId)
    public
    override
    view
    virtual
    returns (
      uint80 roundId_,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    if(roundId > type(uint32).max) { return (0, 0, 0, 0, 0); }
    Transmission memory transmission = s_transmissions[uint32(roundId)];
    return (
      roundId,
      transmission.answer,
      transmission.observationsTimestamp,
      transmission.transmissionTimestamp,
      roundId
    );
  }

  /**
   * @notice aggregator details for the most recently transmitted report
   * @return roundId aggregator round of latest report (NOT OCR round)
   * @return answer median of latest report
   * @return startedAt timestamp of when observations were made offchain
   * @return updatedAt timestamp of block containing latest report
   * @return answeredInRound aggregator round of latest report
   */
  function latestRoundData()
    public
    override
    view
    virtual
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;

    Transmission memory transmission = s_transmissions[latestAggregatorRoundId];
    return (
      latestAggregatorRoundId,
      transmission.answer,
      transmission.observationsTimestamp,
      transmission.transmissionTimestamp,
      latestAggregatorRoundId
    );
  }

  /***************************************************************************
   * Section: Configurable LINK Token
   **************************************************************************/

  // We assume that the token contract is correct. This contract is not written
  // to handle misbehaving ERC20 tokens!
  LinkTokenInterface internal s_linkToken;

  /*
   * @notice emitted when the LINK token contract is set
   * @param oldLinkToken the address of the old LINK token contract
   * @param newLinkToken the address of the new LINK token contract
   */
  event LinkTokenSet(
    LinkTokenInterface indexed oldLinkToken,
    LinkTokenInterface indexed newLinkToken
  );

  /**
   * @notice sets the LINK token contract used for paying oracles
   * @param linkToken the address of the LINK token contract
   * @param recipient remaining funds from the previous token contract are transferred
   * here
   * @dev this function will return early (without an error) without changing any state
   * if linkToken equals getLinkToken().
   * @dev this will trigger a payout so that a malicious owner cannot take from oracles
   * what is already owed to them.
   * @dev we assume that the token contract is correct. This contract is not written
   * to handle misbehaving ERC20 tokens!
   */
  function setLinkToken(
    LinkTokenInterface linkToken,
    address recipient
  ) external
    onlyOwner()
  {
    LinkTokenInterface oldLinkToken = s_linkToken;
    if (linkToken == oldLinkToken) {
      // No change, nothing to be done
      return;
    }
    // call balanceOf as a sanity check on whether we're talking to a token
    // contract
    linkToken.balanceOf(address(this));
    // we break CEI here, but that's okay because we're dealing with a correct
    // token contract (by assumption).
    _payOracles();
    uint256 remainingBalance = oldLinkToken.balanceOf(address(this));
    require(oldLinkToken.transfer(recipient, remainingBalance), "transfer remaining funds failed");
    s_linkToken = linkToken;
    emit LinkTokenSet(oldLinkToken, linkToken);
  }

  /*
   * @notice gets the LINK token contract used for paying oracles
   * @return linkToken the address of the LINK token contract
   */
  function getLinkToken()
    external
    view
    returns(LinkTokenInterface linkToken)
  {
    return s_linkToken;
  }

  /***************************************************************************
   * Section: BillingAccessController Management
   **************************************************************************/

  // Controls who can change billing parameters. A billingAdmin is not able to
  // affect any OCR protocol settings and therefore cannot tamper with the
  // liveness or integrity of a data feed. However, a billingAdmin can set
  // faulty billing parameters causing oracles to be underpaid, or causing them
  // to be paid so much that further calls to setConfig, setBilling,
  // setLinkToken will always fail due to the contract being underfunded.
  AccessControllerInterface internal s_billingAccessController;

  /**
   * @notice emitted when a new access-control contract is set
   * @param old the address prior to the current setting
   * @param current the address of the new access-control contract
   */
  event BillingAccessControllerSet(AccessControllerInterface old, AccessControllerInterface current);

  function _setBillingAccessController(AccessControllerInterface billingAccessController)
    internal
  {
    AccessControllerInterface oldController = s_billingAccessController;
    if (billingAccessController != oldController) {
      s_billingAccessController = billingAccessController;
      emit BillingAccessControllerSet(
        oldController,
        billingAccessController
      );
    }
  }

  /**
   * @notice sets billingAccessController
   * @param _billingAccessController new billingAccessController contract address
   * @dev only owner can call this
   */
  function setBillingAccessController(AccessControllerInterface _billingAccessController)
    external
    onlyOwner
  {
    _setBillingAccessController(_billingAccessController);
  }

  /**
   * @notice gets billingAccessController
   * @return address of billingAccessController contract
   */
  function getBillingAccessController()
    external
    view
    returns (AccessControllerInterface)
  {
    return s_billingAccessController;
  }

  /***************************************************************************
   * Section: Billing Configuration
   **************************************************************************/

  /**
   * @notice emitted when billing parameters are set
   * @param maximumGasPriceGwei highest gas price for which transmitter will be compensated
   * @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value
   * @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report
   * @param transmissionPaymentGjuels reward to transmitter of a successful report
   * @param accountingGas gas overhead incurred by accounting logic
   */
  event BillingSet(
    uint32 maximumGasPriceGwei,
    uint32 reasonableGasPriceGwei,
    uint32 observationPaymentGjuels,
    uint32 transmissionPaymentGjuels,
    uint24 accountingGas
  );

  /**
   * @notice sets billing parameters
   * @param maximumGasPriceGwei highest gas price for which transmitter will be compensated
   * @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value
   * @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report
   * @param transmissionPaymentGjuels reward to transmitter of a successful report
   * @param accountingGas gas overhead incurred by accounting logic
   * @dev access control provided by billingAccessController
   */
  function setBilling(
    uint32 maximumGasPriceGwei,
    uint32 reasonableGasPriceGwei,
    uint32 observationPaymentGjuels,
    uint32 transmissionPaymentGjuels,
    uint24 accountingGas
  )
    external
  {
    AccessControllerInterface access = s_billingAccessController;
    require(msg.sender == owner() || access.hasAccess(msg.sender, msg.data),
      "Only owner&billingAdmin can call");
    _payOracles();

    s_hotVars.maximumGasPriceGwei = maximumGasPriceGwei;
    s_hotVars.reasonableGasPriceGwei = reasonableGasPriceGwei;
    s_hotVars.observationPaymentGjuels = observationPaymentGjuels;
    s_hotVars.transmissionPaymentGjuels = transmissionPaymentGjuels;
    s_hotVars.accountingGas = accountingGas;

    emit BillingSet(maximumGasPriceGwei, reasonableGasPriceGwei,
      observationPaymentGjuels, transmissionPaymentGjuels, accountingGas);
  }

  /**
   * @notice gets billing parameters
   * @param maximumGasPriceGwei highest gas price for which transmitter will be compensated
   * @param reasonableGasPriceGwei transmitter will receive reward for gas prices under this value
   * @param observationPaymentGjuels reward to oracle for contributing an observation to a successfully transmitted report
   * @param transmissionPaymentGjuels reward to transmitter of a successful report
   * @param accountingGas gas overhead of the accounting logic
   */
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
      s_hotVars.accountingGas
    );
  }

  /***************************************************************************
   * Section: Payments and Withdrawals
   **************************************************************************/

  /**
   * @notice withdraws an oracle's payment from the contract
   * @param transmitter the transmitter address of the oracle
   * @dev must be called by oracle's payee address
   */
  function withdrawPayment(address transmitter)
    external
  {
    require(msg.sender == s_payees[transmitter], "Only payee can withdraw");
    _payOracle(transmitter);
  }

  /**
   * @notice query an oracle's payment amount, denominated in juels
   * @param transmitterAddress the transmitter address of the oracle
   */
  function owedPayment(address transmitterAddress)
    public
    view
    returns (uint256)
  {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) { return 0; }
    // safe from overflow:
    // s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index] <= 2**32
    // s_hotVars.observationPaymentGjuels <= 2**32
    // 1 gwei <= 2**32
    // hence juelsAmount <= 2**96
    uint256 juelsAmount =
      uint256(s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index]) *
      uint256(s_hotVars.observationPaymentGjuels) *
      (1 gwei);
    juelsAmount += transmitter.paymentJuels;
    return juelsAmount;
  }

  /**
   * @notice emitted when an oracle has been paid LINK
   * @param transmitter address from which the oracle sends reports to the transmit method
   * @param payee address to which the payment is sent
   * @param amount amount of LINK sent
   * @param linkToken address of the LINK token contract
   */
  event OraclePaid(
    address indexed transmitter,
    address indexed payee,
    uint256 amount,
    LinkTokenInterface indexed linkToken
  );

  // _payOracle pays out transmitter's balance to the corresponding payee, and zeros it out
  function _payOracle(address transmitterAddress)
    internal
  {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) { return; }
    uint256 juelsAmount = owedPayment(transmitterAddress);
    if (juelsAmount > 0) {
      address payee = s_payees[transmitterAddress];
      // Poses no re-entrancy issues, because LINK.transfer does not yield
      // control flow.
      require(s_linkToken.transfer(payee, juelsAmount), "insufficient funds");
      s_rewardFromAggregatorRoundId[transmitter.index] = s_hotVars.latestAggregatorRoundId;
      s_transmitters[transmitterAddress].paymentJuels = 0;
      emit OraclePaid(transmitterAddress, payee, juelsAmount, s_linkToken);
    }
  }

  // _payOracles pays out all transmitters, and zeros out their balances.
  //
  // It's much more gas-efficient to do this as a single operation, to avoid
  // hitting storage too much.
  function _payOracles()
    internal
  {
    unchecked {
      LinkTokenInterface linkToken = s_linkToken;
      uint32 latestAggregatorRoundId = s_hotVars.latestAggregatorRoundId;
      uint32[maxNumOracles] memory rewardFromAggregatorRoundId = s_rewardFromAggregatorRoundId;
      address[] memory transmitters = s_transmittersList;
      for (uint transmitteridx = 0; transmitteridx < transmitters.length; transmitteridx++) {
        uint256 reimbursementAmountJuels = s_transmitters[transmitters[transmitteridx]].paymentJuels;
        s_transmitters[transmitters[transmitteridx]].paymentJuels = 0;
        uint256 obsCount = latestAggregatorRoundId - rewardFromAggregatorRoundId[transmitteridx];
        uint256 juelsAmount =
          obsCount * uint256(s_hotVars.observationPaymentGjuels) * (1 gwei) + reimbursementAmountJuels;
        if (juelsAmount > 0) {
            address payee = s_payees[transmitters[transmitteridx]];
            // Poses no re-entrancy issues, because LINK.transfer does not yield
            // control flow.
            require(linkToken.transfer(payee, juelsAmount), "insufficient funds");
            rewardFromAggregatorRoundId[transmitteridx] = latestAggregatorRoundId;
            emit OraclePaid(transmitters[transmitteridx], payee, juelsAmount, linkToken);
          }
      }
      // "Zero" the accounting storage variables
      s_rewardFromAggregatorRoundId = rewardFromAggregatorRoundId;
    }
  }

  /**
   * @notice withdraw any available funds left in the contract, up to amount, after accounting for the funds due to participants in past reports
   * @param recipient address to send funds to
   * @param amount maximum amount to withdraw, denominated in LINK-wei.
   * @dev access control provided by billingAccessController
   */
  function withdrawFunds(
    address recipient,
    uint256 amount
  )
    external
  {
    require(msg.sender == owner() || s_billingAccessController.hasAccess(msg.sender, msg.data),
      "Only owner&billingAdmin can call");
    uint256 linkDue = _totalLinkDue();
    uint256 linkBalance = s_linkToken.balanceOf(address(this));
    require(linkBalance >= linkDue, "insufficient balance");
    require(s_linkToken.transfer(recipient, _min(linkBalance - linkDue, amount)), "insufficient funds");
  }

  // Total LINK due to participants in past reports (denominated in Juels).
  function _totalLinkDue()
    internal
    view
    returns (uint256 linkDue)
  {
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
    uint32[maxNumOracles] memory rewardFromAggregatorRoundId = s_rewardFromAggregatorRoundId;
    for (uint i = 0; i < n; i++) {
      linkDue += latestAggregatorRoundId - rewardFromAggregatorRoundId[i];
    }
    // Convert observationPaymentGjuels to uint256, or this overflows!
    linkDue *= uint256(s_hotVars.observationPaymentGjuels) * (1 gwei);
    for (uint i = 0; i < n; i++) {
      linkDue += uint256(s_transmitters[transmitters[i]].paymentJuels);
    }
  }

  /**
   * @notice allows oracles to check that sufficient LINK balance is available
   * @return availableBalance LINK available on this contract, after accounting for outstanding obligations. can become negative
   */
  function linkAvailableForPayment()
    external
    view
    returns (int256 availableBalance)
  {
    // there are at most one billion LINK, so this cast is safe
    int256 balance = int256(s_linkToken.balanceOf(address(this)));
    // according to the argument in the definition of _totalLinkDue,
    // _totalLinkDue is never greater than 2**102, so this cast is safe
    int256 due = int256(_totalLinkDue());
    // safe from overflow according to above sizes
    return int256(balance) - int256(due);
  }

  /**
   * @notice number of observations oracle is due to be reimbursed for
   * @param transmitterAddress address used by oracle for signing or transmitting reports
   */
  function oracleObservationCount(address transmitterAddress)
    external
    view
    returns (uint32)
  {
    Transmitter memory transmitter = s_transmitters[transmitterAddress];
    if (!transmitter.active) { return 0; }
    return s_hotVars.latestAggregatorRoundId - s_rewardFromAggregatorRoundId[transmitter.index];
  }

  /***************************************************************************
   * Section: Transmitter Payment
   **************************************************************************/

  // Gas price at which the transmitter should be reimbursed, in gwei/gas
  function _reimbursementGasPriceGwei(
    uint256 txGasPriceGwei,
    uint256 reasonableGasPriceGwei,
    uint256 maximumGasPriceGwei
  )
    internal
    pure
    returns (uint256)
  {
    // this happens on the path for transmissions. we'd rather pay out
    // a wrong reward than risk a liveness failure due to a revert.
    unchecked {
      // Reward the transmitter for choosing an efficient gas price: if they manage
      // to come in lower than considered reasonable, give them half the savings.
      uint256 gasPriceGwei = txGasPriceGwei;
      if (txGasPriceGwei < reasonableGasPriceGwei) {
        // Give transmitter half the savings for coming in under the reasonable gas price
        gasPriceGwei += (reasonableGasPriceGwei - txGasPriceGwei) / 2;
      }
      // Don't reimburse a gas price higher than maximumGasPriceGwei
      return _min(gasPriceGwei, maximumGasPriceGwei);
    }
  }

  // gas reimbursement due the transmitter, in wei
  function _transmitterGasCostWei(
    uint256 initialGas,
    uint256 gasPriceGwei,
    uint256 callDataGas,
    uint256 accountingGas,
    uint256 leftGas
  )
    internal
    pure
    returns (uint256)
  {
    // this happens on the path for transmissions. we'd rather pay out
    // a wrong reward than risk a liveness failure due to a revert.
    unchecked {
      require(initialGas >= leftGas, "leftGas cannot exceed initialGas");
      uint256 usedGas =
        initialGas - leftGas + // observed gas usage
        callDataGas + accountingGas; // estimated gas usage
      uint256 fullGasCostWei = usedGas * gasPriceGwei * (1 gwei);
      return fullGasCostWei;
    }
  }

  function _payTransmitter(
    HotVars memory hotVars,
    int192 juelsPerFeeCoin,
    uint32 initialGas,
    address transmitter
  )
    internal
    virtual
  {
    // this happens on the path for transmissions. we'd rather pay out
    // a wrong reward than risk a liveness failure due to a revert.
    unchecked {
      // we can't deal with negative juelsPerFeeCoin, better to just not pay
      if (juelsPerFeeCoin < 0) {
        return;
      }

      // Reimburse transmitter of the report for gas usage
      uint256 gasPriceGwei = _reimbursementGasPriceGwei(
        tx.gasprice / (1 gwei), // convert to ETH-gwei units
        hotVars.reasonableGasPriceGwei,
        hotVars.maximumGasPriceGwei
      );
      // The following is only an upper bound, as it ignores the cheaper cost for
      // 0 bytes. Safe from overflow, because calldata just isn't that long.
      uint256 callDataGasCost = 16 * msg.data.length;
      uint256 gasLeft = gasleft();
      uint256 gasCostEthWei = _transmitterGasCostWei(
        uint256(initialGas),
        gasPriceGwei,
        callDataGasCost,
        hotVars.accountingGas,
        gasLeft
      );

      // Even if we assume absurdly large values, this still does not overflow. With
      // - usedGas <= 1'000'000 gas <= 2**20 gas
      // - weiPerGas <= 1'000'000 gwei <= 2**50 wei
      // - hence gasCostEthWei <= 2**70
      // - juelsPerFeeCoin <= 2**96 (more than the entire supply)
      // we still fit into 166 bits
      uint256 gasCostJuels = (gasCostEthWei * uint192(juelsPerFeeCoin))/1e18;

      uint96 oldTransmitterPaymentJuels = s_transmitters[transmitter].paymentJuels;
      uint96 newTransmitterPaymentJuels = uint96(uint256(oldTransmitterPaymentJuels) +
        gasCostJuels + uint256(hotVars.transmissionPaymentGjuels) * (1 gwei));

      // overflow *should* never happen, but if it does, let's not persist it.
      if (newTransmitterPaymentJuels < oldTransmitterPaymentJuels) {
        return;
      }
      s_transmitters[transmitter].paymentJuels = newTransmitterPaymentJuels;
    }
  }

  /***************************************************************************
   * Section: Payee Management
   **************************************************************************/

  // Addresses at which oracles want to receive payments, by transmitter address
  mapping (address /* transmitter */ => address /* payment address */)
    internal
    s_payees;

  // Payee addresses which must be approved by the owner
  mapping (address /* transmitter */ => address /* payment address */)
    internal
    s_proposedPayees;

  /**
   * @notice emitted when a transfer of an oracle's payee address has been initiated
   * @param transmitter address from which the oracle sends reports to the transmit method
   * @param current the payee address for the oracle, prior to this setting
   * @param proposed the proposed new payee address for the oracle
   */
  event PayeeshipTransferRequested(
    address indexed transmitter,
    address indexed current,
    address indexed proposed
  );

  /**
   * @notice emitted when a transfer of an oracle's payee address has been completed
   * @param transmitter address from which the oracle sends reports to the transmit method
   * @param current the payee address for the oracle, prior to this setting
   */
  event PayeeshipTransferred(
    address indexed transmitter,
    address indexed previous,
    address indexed current
  );

  /**
   * @notice sets the payees for transmitting addresses
   * @param transmitters addresses oracles use to transmit the reports
   * @param payees addresses of payees corresponding to list of transmitters
   * @dev must be called by owner
   * @dev cannot be used to change payee addresses, only to initially populate them
   */
  function setPayees(
    address[] calldata transmitters,
    address[] calldata payees
  )
    external
    onlyOwner()
  {
    require(transmitters.length == payees.length, "transmitters.size != payees.size");

    for (uint i = 0; i < transmitters.length; i++) {
      address transmitter = transmitters[i];
      address payee = payees[i];
      address currentPayee = s_payees[transmitter];
      bool zeroedOut = currentPayee == address(0);
      require(zeroedOut || currentPayee == payee, "payee already set");
      s_payees[transmitter] = payee;

      if (currentPayee != payee) {
        emit PayeeshipTransferred(transmitter, currentPayee, payee);
      }
    }
  }

  /**
   * @notice first step of payeeship transfer (safe transfer pattern)
   * @param transmitter transmitter address of oracle whose payee is changing
   * @param proposed new payee address
   * @dev can only be called by payee address
   */
  function transferPayeeship(
    address transmitter,
    address proposed
  )
    external
  {
    require(msg.sender == s_payees[transmitter], "only current payee can update");
    require(msg.sender != proposed, "cannot transfer to self");

    address previousProposed = s_proposedPayees[transmitter];
    s_proposedPayees[transmitter] = proposed;

    if (previousProposed != proposed) {
      emit PayeeshipTransferRequested(transmitter, msg.sender, proposed);
    }
  }

  /**
   * @notice second step of payeeship transfer (safe transfer pattern)
   * @param transmitter transmitter address of oracle whose payee is changing
   * @dev can only be called by proposed new payee address
   */
  function acceptPayeeship(
    address transmitter
  )
    external
  {
    require(msg.sender == s_proposedPayees[transmitter], "only proposed payees can accept");

    address currentPayee = s_payees[transmitter];
    s_payees[transmitter] = msg.sender;
    s_proposedPayees[transmitter] = address(0);

    emit PayeeshipTransferred(transmitter, currentPayee, msg.sender);
  }

  /***************************************************************************
   * Section: TypeAndVersionInterface
   **************************************************************************/

  function typeAndVersion()
    external
    override
    pure
    virtual
    returns (string memory)
  {
    return "OCR2Aggregator 1.0.0-alpha";
  }

  /***************************************************************************
   * Section: Helper Functions
   **************************************************************************/

  function _min(
    uint256 a,
    uint256 b
  )
    internal
    pure
    returns (uint256)
  {
    unchecked {
      if (a < b) { return a; }
      return b;
    }
  }
}
