// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";
import {OCR2Abstract} from "./OCR2Abstract.sol";

/**
 * @notice Onchain verification of reports from the offchain reporting protocol
 * @dev For details on its operation, see the offchain reporting protocol design
 * doc, which refers to this contract as simply the "contract".
 */
abstract contract OCR2Base is ConfirmedOwner, OCR2Abstract {
  error ReportInvalid(string message);
  error InvalidConfig(string message);

  constructor() ConfirmedOwner(msg.sender) {}

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  uint32 internal s_latestConfigBlockNumber; // makes it easier for offchain systems
  // to extract config from logs.

  // Storing these fields used on the hot path in a ConfigInfo variable reduces the
  // retrieval of all of them into two SLOADs. If any further fields are
  // added, make sure that storage of the struct still takes at most 64 bytes.
  struct ConfigInfo {
    bytes32 latestConfigDigest;
    uint8 f; // ───╮
    uint8 n; // ───╯
  }
  ConfigInfo internal s_configInfo;

  // Used for s_oracles[a].role, where a is an address, to track the purpose
  // of the address, or to indicate that the address is unset.
  enum Role {
    // No oracle role has been set for address a
    Unset,
    // Signing address for the s_oracles[a].index'th oracle. I.e., report
    // signatures from this oracle should ecrecover back to address a.
    Signer,
    // Transmission address for the s_oracles[a].index'th oracle. I.e., if a
    // report is received by OCR2Aggregator.transmit in which msg.sender is
    // a, it is attributed to the s_oracles[a].index'th oracle.
    Transmitter
  }

  struct Oracle {
    uint8 index; // Index of oracle in s_signers/s_transmitters
    Role role; // Role of the address which mapped to this struct
  }

  mapping(address signerOrTransmitter => Oracle) internal s_oracles;

  // s_signers contains the signing address of each oracle
  address[] internal s_signers;

  // s_transmitters contains the transmission address of each oracle,
  // i.e. the address the oracle actually sends transactions to the contract from
  address[] internal s_transmitters;

  struct DecodedReport {
    bytes32[] requestIds;
    bytes[] results;
    bytes[] errors;
    bytes[] onchainMetadata;
    bytes[] offchainMetadata;
  }

  /*
   * Config logic
   */

  // Reverts transaction if config args are invalid
  modifier checkConfigValid(
    uint256 numSigners,
    uint256 numTransmitters,
    uint256 f
  ) {
    if (numSigners > MAX_NUM_ORACLES) revert InvalidConfig("too many signers");
    if (f == 0) revert InvalidConfig("f must be positive");
    if (numSigners != numTransmitters) revert InvalidConfig("oracle addresses out of registration");
    if (numSigners <= 3 * f) revert InvalidConfig("faulty-oracle f too high");
    _;
  }

  // solhint-disable-next-line gas-struct-packing
  struct SetConfigArgs {
    address[] signers;
    address[] transmitters;
    uint8 f;
    bytes onchainConfig;
    uint64 offchainConfigVersion;
    bytes offchainConfig;
  }

  /// @inheritdoc OCR2Abstract
  function latestConfigDigestAndEpoch()
    external
    view
    virtual
    override
    returns (bool scanLogs, bytes32 configDigest, uint32 epoch)
  {
    return (true, bytes32(0), uint32(0));
  }

  /**
   * @notice sets offchain reporting protocol configuration incl. participating oracles
   * @param _signers addresses with which oracles sign the reports
   * @param _transmitters addresses oracles use to transmit the reports
   * @param _f number of faulty oracles the system can tolerate
   * @param _onchainConfig encoded on-chain contract configuration
   * @param _offchainConfigVersion version number for offchainEncoding schema
   * @param _offchainConfig encoded off-chain oracle configuration
   */
  function setConfig(
    address[] memory _signers,
    address[] memory _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _offchainConfigVersion,
    bytes memory _offchainConfig
  ) external override checkConfigValid(_signers.length, _transmitters.length, _f) onlyOwner {
    SetConfigArgs memory args = SetConfigArgs({
      signers: _signers,
      transmitters: _transmitters,
      f: _f,
      onchainConfig: _onchainConfig,
      offchainConfigVersion: _offchainConfigVersion,
      offchainConfig: _offchainConfig
    });

    _beforeSetConfig(args.f, args.onchainConfig);

    while (s_signers.length != 0) {
      // remove any old signer/transmitter addresses
      uint256 lastIdx = s_signers.length - 1;
      address signer = s_signers[lastIdx];
      address transmitter = s_transmitters[lastIdx];
      delete s_oracles[signer];
      delete s_oracles[transmitter];
      s_signers.pop();
      s_transmitters.pop();
    }

    // Bounded by MAX_NUM_ORACLES in OCR2Abstract.sol
    for (uint256 i = 0; i < args.signers.length; i++) {
      if (args.signers[i] == address(0)) revert InvalidConfig("signer must not be empty");
      if (args.transmitters[i] == address(0)) revert InvalidConfig("transmitter must not be empty");
      // add new signer/transmitter addresses
      if (s_oracles[args.signers[i]].role != Role.Unset) revert InvalidConfig("repeated signer address");
      s_oracles[args.signers[i]] = Oracle(uint8(i), Role.Signer);
      if (s_oracles[args.transmitters[i]].role != Role.Unset) revert InvalidConfig("repeated transmitter address");
      s_oracles[args.transmitters[i]] = Oracle(uint8(i), Role.Transmitter);
      s_signers.push(args.signers[i]);
      s_transmitters.push(args.transmitters[i]);
    }
    s_configInfo.f = args.f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    {
      s_configInfo.latestConfigDigest = _configDigestFromConfigData(
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
    }
    s_configInfo.n = uint8(args.signers.length);

    emit ConfigSet(
      previousConfigBlockNumber,
      s_configInfo.latestConfigDigest,
      s_configCount,
      args.signers,
      args.transmitters,
      args.f,
      args.onchainConfig,
      args.offchainConfigVersion,
      args.offchainConfig
    );
  }

  function _configDigestFromConfigData(
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    address[] memory _signers,
    address[] memory _transmitters,
    uint8 _f,
    bytes memory _onchainConfig,
    uint64 _encodedConfigVersion,
    bytes memory _encodedConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          _chainId,
          _contractAddress,
          _configCount,
          _signers,
          _transmitters,
          _f,
          _onchainConfig,
          _encodedConfigVersion,
          _encodedConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32(prefix | (h & ~prefixMask));
  }

  /**
   * @notice information about current offchain reporting protocol configuration
   * @return configCount ordinal number of current config, out of all configs applied to this contract so far
   * @return blockNumber block at which this config was set
   * @return configDigest domain-separation tag for current config (see __configDigestFromConfigData)
   */
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_configCount, s_latestConfigBlockNumber, s_configInfo.latestConfigDigest);
  }

  /**
   * @return list of addresses permitted to transmit reports to this contract
   * @dev The list will match the order used to specify the transmitter during setConfig
   */
  function transmitters() external view returns (address[] memory) {
    return s_transmitters;
  }

  function _beforeSetConfig(uint8 _f, bytes memory _onchainConfig) internal virtual;

  /**
   * @dev hook called after the report has been fully validated
   * for the extending contract to handle additional logic, such as oracle payment
   * @param decodedReport decodedReport
   */
  function _report(DecodedReport memory decodedReport) internal virtual;

  // The constant-length components of the msg.data sent to transmit.
  // See the "If we wanted to call sam" example on for example reasoning
  // https://solidity.readthedocs.io/en/v0.7.2/abi-spec.html
  uint16 private constant TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT =
    4 + // function selector
      32 *
      3 + // 3 words containing reportContext
      32 + // word containing start location of abiencoded report value
      32 + // word containing location start of abiencoded rs value
      32 + // word containing start location of abiencoded ss value
      32 + // rawVs value
      32 + // word containing length of report
      32 + // word containing length rs
      32 + // word containing length of ss
      0; // placeholder

  function _requireExpectedMsgDataLength(
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss
  ) private pure {
    // calldata will never be big enough to make this overflow
    uint256 expected = uint256(TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT) +
      report.length + // one byte pure entry in _report
      rs.length *
      32 + // 32 bytes per entry in _rs
      ss.length *
      32 + // 32 bytes per entry in _ss
      0; // placeholder
    if (msg.data.length != expected) revert ReportInvalid("calldata length mismatch");
  }

  function _beforeTransmit(
    bytes calldata report
  ) internal virtual returns (bool shouldStop, DecodedReport memory decodedReport);

  /**
   * @notice transmit is called to post a new report to the contract
   * @param report serialized report, which the signatures are signing.
   * @param rs ith element is the R components of the ith signature on report. Must have at most maxNumOracles entries
   * @param ss ith element is the S components of the ith signature on report. Must have at most maxNumOracles entries
   * @param rawVs ith element is the the V component of the ith signature
   */
  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) external override {
    (bool shouldStop, DecodedReport memory decodedReport) = _beforeTransmit(report);

    if (shouldStop) {
      return;
    }

    {
      // reportContext consists of:
      // reportContext[0]: ConfigDigest
      // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
      // reportContext[2]: ExtraHash
      bytes32 configDigest = reportContext[0];
      uint32 epochAndRound = uint32(uint256(reportContext[1]));

      emit Transmitted(configDigest, uint32(epochAndRound >> 8));

      // The following check is disabled to allow both current and proposed routes to submit reports using the same OCR config digest
      // Chainlink Functions uses globally unique request IDs. Metadata about the request is stored and checked in the Coordinator and Router
      // require(configInfo.latestConfigDigest == configDigest, "configDigest mismatch");

      _requireExpectedMsgDataLength(report, rs, ss);

      uint256 expectedNumSignatures = (s_configInfo.n + s_configInfo.f) / 2 + 1;

      if (rs.length != expectedNumSignatures) revert ReportInvalid("wrong number of signatures");
      if (rs.length != ss.length) revert ReportInvalid("report rs and ss must be of equal length");

      Oracle memory transmitter = s_oracles[msg.sender];
      if (transmitter.role != Role.Transmitter && msg.sender != s_transmitters[transmitter.index])
        revert ReportInvalid("unauthorized transmitter");
    }

    address[MAX_NUM_ORACLES] memory signed;

    {
      // Verify signatures attached to report
      bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));

      Oracle memory o;
      // Bounded by MAX_NUM_ORACLES in OCR2Abstract.sol
      for (uint256 i = 0; i < rs.length; ++i) {
        address signer = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
        o = s_oracles[signer];
        if (o.role != Role.Signer) revert ReportInvalid("address not authorized to sign");
        if (signed[o.index] != address(0)) revert ReportInvalid("non-unique signature");
        signed[o.index] = signer;
      }
    }

    _report(decodedReport);
  }
}
