// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;

import "./AccessControllerInterface.sol";
import "./AggregatorV2V3Interface.sol";
import "./AggregatorValidatorInterface.sol";
import "./LinkTokenInterface.sol";
import "./Owned.sol";
import "./OffchainAggregatorBilling.sol";
import "./TypeAndVersionInterface.sol";

/**
  * @notice Onchain verification of reports from the offchain reporting protocol

  * @dev For details on its operation, see the offchain reporting protocol design
  * @dev doc, which refers to this contract as simply the "contract".
*/
contract OffchainAggregator is Owned, OffchainAggregatorBilling, AggregatorV2V3Interface, TypeAndVersionInterface {

  uint256 constant private maxUint32 = (1 << 32) - 1;

  // Storing these fields used on the hot path in a HotVars variable reduces the
  // retrieval of all of them to a single SLOAD. If any further fields are
  // added, make sure that storage of the struct still takes at most 32 bytes.
  struct HotVars {
    // Provides 128 bits of security against 2nd pre-image attacks, but only
    // 64 bits against collisions. This is acceptable, since a malicious owner has
    // easier way of messing up the protocol than to find hash collisions.
    bytes16 latestConfigDigest;
    uint40 latestEpochAndRound; // 32 most sig bits for epoch, 8 least sig bits for round
    // Current bound assumed on number of faulty/dishonest oracles participating
    // in the protocol, this value is referred to as f in the design
    uint8 threshold;
    // Chainlink Aggregators expose a roundId to consumers. The offchain reporting
    // protocol does not use this id anywhere. We increment it whenever a new
    // transmission is made to provide callers with contiguous ids for successive
    // reports.
    uint32 latestAggregatorRoundId;
  }
  HotVars internal s_hotVars;

  // Transmission records the median answer from the transmit transaction at
  // time timestamp
  struct Transmission {
    int192 answer; // 192 bits ought to be enough for anyone
    uint64 timestamp;
  }
  mapping(uint32 /* aggregator round ID */ => Transmission) internal s_transmissions;

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  uint32 internal s_latestConfigBlockNumber; // makes it easier for offchain systems
                                             // to extract config from logs.

  // Lowest answer the system is allowed to report in response to transmissions
  int192 immutable public minAnswer;
  // Highest answer the system is allowed to report in response to transmissions
  int192 immutable public maxAnswer;

  /*
   * @param _maximumGasPrice highest gas price for which transmitter will be compensated
   * @param _reasonableGasPrice transmitter will receive reward for gas prices under this value
   * @param _microLinkPerEth reimbursement per ETH of gas cost, in 1e-6LINK units
   * @param _linkGweiPerObservation reward to oracle for contributing an observation to a successfully transmitted report, in 1e-9LINK units
   * @param _linkGweiPerTransmission reward to transmitter of a successful report, in 1e-9LINK units
   * @param _link address of the LINK contract
   * @param _minAnswer lowest answer the median of a report is allowed to be
   * @param _maxAnswer highest answer the median of a report is allowed to be
   * @param _billingAccessController access controller for billing admin functions
   * @param _requesterAccessController access controller for requesting new rounds
   * @param _decimals answers are stored in fixed-point format, with this many digits of precision
   * @param _description short human-readable description of observable this contract's answers pertain to
   */
  constructor(
    uint32 _maximumGasPrice,
    uint32 _reasonableGasPrice,
    uint32 _microLinkPerEth,
    uint32 _linkGweiPerObservation,
    uint32 _linkGweiPerTransmission,
    LinkTokenInterface _link,
    int192 _minAnswer,
    int192 _maxAnswer,
    AccessControllerInterface _billingAccessController,
    AccessControllerInterface _requesterAccessController,
    uint8 _decimals,
    string memory _description
  )
    OffchainAggregatorBilling(_maximumGasPrice, _reasonableGasPrice, _microLinkPerEth,
      _linkGweiPerObservation, _linkGweiPerTransmission, _link,
      _billingAccessController
    )
  {
    decimals = _decimals;
    s_description = _description;
    setRequesterAccessController(_requesterAccessController);
    setValidatorConfig(AggregatorValidatorInterface(0x0), 0);
    minAnswer = _minAnswer;
    maxAnswer = _maxAnswer;
  }

  /*
   * Versioning
   */
  function typeAndVersion()
    external
    override
    pure
    virtual
    returns (string memory)
  {
    return "OffchainAggregator 4.0.0";
  }

  /*
   * Config logic
   */

  /**
   * @notice triggers a new run of the offchain reporting protocol
   * @param previousConfigBlockNumber block in which the previous config was set, to simplify historic analysis
   * @param configCount ordinal number of this config setting among all config settings over the life of this contract
   * @param signers ith element is address ith oracle uses to sign a report
   * @param transmitters ith element is address ith oracle uses to transmit a report via the transmit method
   * @param threshold maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly
   * @param encodedConfigVersion version of the serialization format used for "encoded" parameter
   * @param encoded serialized data used by oracles to configure their offchain operation
   */
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    uint64 configCount,
    address[] signers,
    address[] transmitters,
    uint8 threshold,
    uint64 encodedConfigVersion,
    bytes encoded
  );

  // Reverts transaction if config args are invalid
  modifier checkConfigValid (
    uint256 _numSigners, uint256 _numTransmitters, uint256 _threshold
  ) {
    require(_numSigners <= maxNumOracles, "too many signers");
    require(_threshold > 0, "threshold must be positive");
    require(
      _numSigners == _numTransmitters,
      "oracle addresses out of registration"
    );
    require(_numSigners > 3*_threshold, "faulty-oracle threshold too high");
    _;
  }

  /**
   * @notice sets offchain reporting protocol configuration incl. participating oracles
   * @param _signers addresses with which oracles sign the reports
   * @param _transmitters addresses oracles use to transmit the reports
   * @param _threshold number of faulty oracles the system can tolerate
   * @param _encodedConfigVersion version number for offchainEncoding schema
   * @param _encoded encoded off-chain oracle configuration
   */
  function setConfig(
    address[] calldata _signers,
    address[] calldata _transmitters,
    uint8 _threshold,
    uint64 _encodedConfigVersion,
    bytes calldata _encoded
  )
    external
    checkConfigValid(_signers.length, _transmitters.length, _threshold)
    onlyOwner()
  {
    while (s_signers.length != 0) { // remove any old signer/transmitter addresses
      uint lastIdx = s_signers.length - 1;
      address signer = s_signers[lastIdx];
      address transmitter = s_transmitters[lastIdx];
      payOracle(transmitter);
      delete s_oracles[signer];
      delete s_oracles[transmitter];
      s_signers.pop();
      s_transmitters.pop();
    }

    for (uint i = 0; i < _signers.length; i++) { // add new signer/transmitter addresses
      require(
        s_oracles[_signers[i]].role == Role.Unset,
        "repeated signer address"
      );
      s_oracles[_signers[i]] = Oracle(uint8(i), Role.Signer);
      require(s_payees[_transmitters[i]] != address(0), "payee must be set");
      require(
        s_oracles[_transmitters[i]].role == Role.Unset,
        "repeated transmitter address"
      );
      s_oracles[_transmitters[i]] = Oracle(uint8(i), Role.Transmitter);
      s_signers.push(_signers[i]);
      s_transmitters.push(_transmitters[i]);
    }
    s_hotVars.threshold = _threshold;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    uint64 configCount = s_configCount;
    {
      s_hotVars.latestConfigDigest = configDigestFromConfigData(
        address(this),
        configCount,
        _signers,
        _transmitters,
        _threshold,
        _encodedConfigVersion,
        _encoded
      );
      s_hotVars.latestEpochAndRound = 0;
    }
    emit ConfigSet(
      previousConfigBlockNumber,
      configCount,
      _signers,
      _transmitters,
      _threshold,
      _encodedConfigVersion,
      _encoded
    );
  }

  function configDigestFromConfigData(
    address _contractAddress,
    uint64 _configCount,
    address[] calldata _signers,
    address[] calldata _transmitters,
    uint8 _threshold,
    uint64 _encodedConfigVersion,
    bytes calldata _encodedConfig
  ) internal pure returns (bytes16) {
    return bytes16(keccak256(abi.encode(_contractAddress, _configCount,
      _signers, _transmitters, _threshold, _encodedConfigVersion, _encodedConfig
    )));
  }

  /**
   * @notice information about current offchain reporting protocol configuration

   * @return configCount ordinal number of current config, out of all configs applied to this contract so far
   * @return blockNumber block at which this config was set
   * @return configDigest domain-separation tag for current config (see configDigestFromConfigData)
   */
  function latestConfigDetails()
    external
    view
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes16 configDigest
    )
  {
    return (s_configCount, s_latestConfigBlockNumber, s_hotVars.latestConfigDigest);
  }

  /**
   * @return list of addresses permitted to transmit reports to this contract

   * @dev The list will match the order used to specify the transmitter during setConfig
   */
  function transmitters()
    external
    view
    returns(address[] memory)
  {
      return s_transmitters;
  }

  /*
   * On-chain validation logc
   */

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
  function validatorConfig()
    external
    view
    returns (AggregatorValidatorInterface validator, uint32 gasLimit)
  {
    ValidatorConfig memory vc = s_validatorConfig;
    return (vc.validator, vc.gasLimit);
  }

  /**
   * @notice sets validator configuration
   * @dev set _newValidator to 0x0 to disable validate calls
   * @param _newValidator address of the new validator contract
   * @param _newGasLimit new gas limit for validate calls
   */
  function setValidatorConfig(AggregatorValidatorInterface _newValidator, uint32 _newGasLimit)
    public
    onlyOwner()
  {
    ValidatorConfig memory previous = s_validatorConfig;

    if (previous.validator != _newValidator || previous.gasLimit != _newGasLimit) {
      s_validatorConfig = ValidatorConfig({
        validator: _newValidator,
        gasLimit: _newGasLimit
      });

      emit ValidatorConfigSet(previous.validator, previous.gasLimit, _newValidator, _newGasLimit);
    }
  }

  function validateAnswer(
    uint32 _aggregatorRoundId,
    int256 _answer
  )
    private
  {
    ValidatorConfig memory vc = s_validatorConfig;

    if (address(vc.validator) == address(0)) {
      return;
    }

    uint32 prevAggregatorRoundId = _aggregatorRoundId - 1;
    int256 prevAggregatorRoundAnswer = s_transmissions[prevAggregatorRoundId].answer;
    require(
      callWithExactGasEvenIfTargetIsNoContract(
        vc.gasLimit,
        address(vc.validator),
        abi.encodeWithSignature(
          "validate(uint256,int256,uint256,int256)",
          uint256(prevAggregatorRoundId),
          prevAggregatorRoundAnswer,
          uint256(_aggregatorRoundId),
          _answer
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
  function callWithExactGasEvenIfTargetIsNoContract(
    uint256 _gasAmount,
    address _target,
    bytes memory _data
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
        // If g - g//64 <= _gasAmount, we don't have enough gas. (We subtract g//64
        // because of EIP-150.)
        if gt(sub(g, div(g, 64)), _gasAmount) {
          // Call and ignore success/return data. Note that we did not check
          // whether a contract actually exists at the _target address.
          pop(call(_gasAmount, _target, 0, add(_data, 0x20), mload(_data), 0, 0))
          sufficientGas := true
        }
      }
    }
  }

  /*
   * requestNewRound logic
   */

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
  event RoundRequested(address indexed requester, bytes16 configDigest, uint32 epoch, uint8 round);

  /**
   * @notice address of the requester access controller contract
   * @return requester access controller address
   */
  function requesterAccessController()
    external
    view
    returns (AccessControllerInterface)
  {
    return s_requesterAccessController;
  }

  /**
   * @notice sets the requester access controller
   * @param _requesterAccessController designates the address of the new requester access controller
   */
  function setRequesterAccessController(AccessControllerInterface _requesterAccessController)
    public
    onlyOwner()
  {
    AccessControllerInterface oldController = s_requesterAccessController;
    if (_requesterAccessController != oldController) {
      s_requesterAccessController = AccessControllerInterface(_requesterAccessController);
      emit RequesterAccessControllerSet(oldController, _requesterAccessController);
    }
  }

  /**
   * @notice immediately requests a new round
   * @return the aggregatorRoundId of the next round. Note: The report for this round may have been
   * transmitted (but not yet mined) *before* requestNewRound() was even called. There is *no*
   * guarantee of causality between the request and the report at aggregatorRoundId.
   */
  function requestNewRound() external returns (uint80) {
    require(msg.sender == owner || s_requesterAccessController.hasAccess(msg.sender, msg.data),
      "Only owner&requester can call");

    HotVars memory hotVars = s_hotVars;

    emit RoundRequested(
      msg.sender,
      hotVars.latestConfigDigest,
      uint32(s_hotVars.latestEpochAndRound >> 8),
      uint8(s_hotVars.latestEpochAndRound)
    );
    return hotVars.latestAggregatorRoundId + 1;
  }

  /*
   * Transmission logic
   */

  /**
   * @notice indicates that a new report was transmitted
   * @param aggregatorRoundId the round to which this report was assigned
   * @param answer median of the observations attached this report
   * @param transmitter address from which the report was transmitted
   * @param observations observations transmitted with this report
   * @param rawReportContext signature-replay-prevention domain-separation tag
   */
  event NewTransmission(
    uint32 indexed aggregatorRoundId,
    int192 answer,
    address transmitter,
    int192[] observations,
    bytes observers,
    bytes32 rawReportContext
  );

  // decodeReport is used to check that the solidity and go code are using the
  // same format. See TestOffchainAggregator.testDecodeReport and TestReportParsing
  function decodeReport(bytes memory _report)
    internal
    pure
    returns (
      bytes32 rawReportContext,
      bytes32 rawObservers,
      int192[] memory observations
    )
  {
    (rawReportContext, rawObservers, observations) = abi.decode(_report,
      (bytes32, bytes32, int192[]));
  }

  // Used to relieve stack pressure in transmit
  struct ReportData {
    HotVars hotVars; // Only read from storage once
    bytes observers; // ith element is the index of the ith observer
    int192[] observations; // ith element is the ith observation
    bytes vs; // jth element is the v component of the jth signature
    bytes32 rawReportContext;
  }

  /*
   * @notice details about the most recent report

   * @return configDigest domain separation tag for the latest report
   * @return epoch epoch in which the latest report was generated
   * @return round OCR round in which the latest report was generated
   * @return latestAnswer median value from latest report
   * @return latestTimestamp when the latest report was transmitted
   */
  function latestTransmissionDetails()
    external
    view
    returns (
      bytes16 configDigest,
      uint32 epoch,
      uint8 round,
      int192 latestAnswer,
      uint64 latestTimestamp
    )
  {
    require(msg.sender == tx.origin, "Only callable by EOA");
    return (
      s_hotVars.latestConfigDigest,
      uint32(s_hotVars.latestEpochAndRound >> 8),
      uint8(s_hotVars.latestEpochAndRound),
      s_transmissions[s_hotVars.latestAggregatorRoundId].answer,
      s_transmissions[s_hotVars.latestAggregatorRoundId].timestamp
    );
  }

  // The constant-length components of the msg.data sent to transmit.
  // See the "If we wanted to call sam" example on for example reasoning
  // https://solidity.readthedocs.io/en/v0.7.2/abi-spec.html
  uint16 private constant TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT =
    4 + // function selector
    32 + // word containing start location of abiencoded _report value
    32 + // word containing location start of abiencoded  _rs value
    32 + // word containing start location of abiencoded _ss value
    32 + // _rawVs value
    32 + // word containing length of _report
    32 + // word containing length _rs
    32 + // word containing length of _ss
    0; // placeholder

  function expectedMsgDataLength(
    bytes calldata _report, bytes32[] calldata _rs, bytes32[] calldata _ss
  ) private pure returns (uint256 length)
  {
    // calldata will never be big enough to make this overflow
    return uint256(TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT) +
      _report.length + // one byte pure entry in _report
      _rs.length * 32 + // 32 bytes per entry in _rs
      _ss.length * 32 + // 32 bytes per entry in _ss
      0; // placeholder
  }

  /**
   * @notice transmit is called to post a new report to the contract
   * @param _report serialized report, which the signatures are signing. See parsing code below for format. The ith element of the observers component must be the index in s_signers of the address for the ith signature
   * @param _rs ith element is the R components of the ith signature on report. Must have at most maxNumOracles entries
   * @param _ss ith element is the S components of the ith signature on report. Must have at most maxNumOracles entries
   * @param _rawVs ith element is the the V component of the ith signature
   */
  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes calldata _report,
    bytes32[] calldata _rs, bytes32[] calldata _ss, bytes32 _rawVs // signatures
  )
    external
  {
    uint256 initialGas = gasleft(); // This line must come first
    // Make sure the transmit message-length matches the inputs. Otherwise, the
    // transmitter could append an arbitrarily long (up to gas-block limit)
    // string of 0 bytes, which we would reimburse at a rate of 16 gas/byte, but
    // which would only cost the transmitter 4 gas/byte. (Appendix G of the
    // yellow paper, p. 25, for G_txdatazero and EIP 2028 for G_txdatanonzero.)
    // This could amount to reimbursement profit of 36 million gas, given a 3MB
    // zero tail.
    require(msg.data.length == expectedMsgDataLength(_report, _rs, _ss),
      "transmit message too long");
    ReportData memory r; // Relieves stack pressure
    {
      r.hotVars = s_hotVars; // cache read from storage

      bytes32 rawObservers;
      (r.rawReportContext, rawObservers, r.observations) = abi.decode(
        _report, (bytes32, bytes32, int192[])
      );

      // rawReportContext consists of:
      // 11-byte zero padding
      // 16-byte configDigest
      // 4-byte epoch
      // 1-byte round

      bytes16 configDigest = bytes16(r.rawReportContext << 88);
      require(
        r.hotVars.latestConfigDigest == configDigest,
        "configDigest mismatch"
      );

      uint40 epochAndRound = uint40(uint256(r.rawReportContext));

      // direct numerical comparison works here, because
      //
      //   ((e,r) <= (e',r')) implies (epochAndRound <= epochAndRound')
      //
      // because alphabetic ordering implies e <= e', and if e = e', then r<=r',
      // so e*256+r <= e'*256+r', because r, r' < 256
      require(r.hotVars.latestEpochAndRound < epochAndRound, "stale report");

      require(_rs.length > r.hotVars.threshold, "not enough signatures");
      require(_rs.length <= maxNumOracles, "too many signatures");
      require(_ss.length == _rs.length, "signatures out of registration");
      require(r.observations.length <= maxNumOracles,
              "num observations out of bounds");
      require(r.observations.length > 2 * r.hotVars.threshold,
              "too few values to trust median");

      // Copy signature parities in bytes32 _rawVs to bytes r.v
      r.vs = new bytes(_rs.length);
      for (uint8 i = 0; i < _rs.length; i++) {
        r.vs[i] = _rawVs[i];
      }

      // Copy observer identities in bytes32 rawObservers to bytes r.observers
      r.observers = new bytes(r.observations.length);
      bool[maxNumOracles] memory seen;
      for (uint8 i = 0; i < r.observations.length; i++) {
        uint8 observerIdx = uint8(rawObservers[i]);
        require(!seen[observerIdx], "observer index repeated");
        seen[observerIdx] = true;
        r.observers[i] = rawObservers[i];
      }

      Oracle memory transmitter = s_oracles[msg.sender];
      require( // Check that sender is authorized to report
        transmitter.role == Role.Transmitter &&
        msg.sender == s_transmitters[transmitter.index],
        "unauthorized transmitter"
      );
      // record epochAndRound here, so that we don't have to carry the local
      // variable in transmit. The change is reverted if something fails later.
      r.hotVars.latestEpochAndRound = epochAndRound;
    }

    { // Verify signatures attached to report
      bytes32 h = keccak256(_report);
      bool[maxNumOracles] memory signed;

      Oracle memory o;
      for (uint i = 0; i < _rs.length; i++) {
        address signer = ecrecover(h, uint8(r.vs[i])+27, _rs[i], _ss[i]);
        o = s_oracles[signer];
        require(o.role == Role.Signer, "address not authorized to sign");
        require(!signed[o.index], "non-unique signature");
        signed[o.index] = true;
      }
    }

    { // Check the report contents, and record the result
      for (uint i = 0; i < r.observations.length - 1; i++) {
        bool inOrder = r.observations[i] <= r.observations[i+1];
        require(inOrder, "observations not sorted");
      }

      int192 median = r.observations[r.observations.length/2];
      require(minAnswer <= median && median <= maxAnswer, "median is out of min-max range");
      r.hotVars.latestAggregatorRoundId++;
      s_transmissions[r.hotVars.latestAggregatorRoundId] =
        Transmission(median, uint64(block.timestamp));

      emit NewTransmission(
        r.hotVars.latestAggregatorRoundId,
        median,
        msg.sender,
        r.observations,
        r.observers,
        r.rawReportContext
      );
      // Emit these for backwards compatability with offchain consumers
      // that only support legacy events
      emit NewRound(
        r.hotVars.latestAggregatorRoundId,
        address(0x0), // use zero address since we don't have anybody "starting" the round here
        block.timestamp
      );
      emit AnswerUpdated(
        median,
        r.hotVars.latestAggregatorRoundId,
        block.timestamp
      );

      validateAnswer(r.hotVars.latestAggregatorRoundId, median);
    }
    s_hotVars = r.hotVars;
    assert(initialGas < maxUint32);
    reimburseAndRewardOracles(uint32(initialGas), r.observers);
  }

  /*
   * v2 Aggregator interface
   */

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
    return s_transmissions[s_hotVars.latestAggregatorRoundId].timestamp;
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
   * @param _roundId the aggregator round of the target report
   */
  function getAnswer(uint256 _roundId)
    public
    override
    view
    virtual
    returns (int256)
  {
    if (_roundId > 0xFFFFFFFF) { return 0; }
    return s_transmissions[uint32(_roundId)].answer;
  }

  /**
   * @notice timestamp of block in which report from given aggregator round was transmitted
   * @param _roundId aggregator round (NOT OCR round) of target report
   */
  function getTimestamp(uint256 _roundId)
    public
    override
    view
    virtual
    returns (uint256)
  {
    if (_roundId > 0xFFFFFFFF) { return 0; }
    return s_transmissions[uint32(_roundId)].timestamp;
  }

  /*
   * v3 Aggregator interface
   */

  string constant private V3_NO_DATA_ERROR = "No data present";

  /**
   * @return answers are stored in fixed-point format, with this many digits of precision
   */
  uint8 immutable public override decimals;

  /**
   * @notice aggregator contract version
   */
  uint256 constant public override version = 4;

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
   * @param _roundId target aggregator round (NOT OCR round). Must fit in uint32
   * @return roundId _roundId
   * @return answer median of report from given _roundId
   * @return startedAt timestamp of block in which report from given _roundId was transmitted
   * @return updatedAt timestamp of block in which report from given _roundId was transmitted
   * @return answeredInRound _roundId
   */
  function getRoundData(uint80 _roundId)
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
    require(_roundId <= 0xFFFFFFFF, V3_NO_DATA_ERROR);
    Transmission memory transmission = s_transmissions[uint32(_roundId)];
    return (
      _roundId,
      transmission.answer,
      transmission.timestamp,
      transmission.timestamp,
      _roundId
    );
  }

  /**
   * @notice aggregator details for the most recently transmitted report
   * @return roundId aggregator round of latest report (NOT OCR round)
   * @return answer median of latest report
   * @return startedAt timestamp of block containing latest report
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
    roundId = s_hotVars.latestAggregatorRoundId;

    // Skipped for compatability with existing FluxAggregator in which latestRoundData never reverts.
    // require(roundId != 0, V3_NO_DATA_ERROR);

    Transmission memory transmission = s_transmissions[uint32(roundId)];
    return (
      roundId,
      transmission.answer,
      transmission.timestamp,
      transmission.timestamp,
      roundId
    );
  }
}
