// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {OCR2Abstract} from "./OCR2Abstract.sol";

/// @notice Onchain verification of reports from the offchain reporting protocol
/// @dev For details on its operation, see the offchain reporting protocol design
/// doc, which refers to this contract as simply the "contract".
abstract contract OCR2Base is OwnerIsCreator, OCR2Abstract {
  error InvalidConfig(string message);
  error WrongMessageLength(uint256 expected, uint256 actual);
  error ConfigDigestMismatch(bytes32 expected, bytes32 actual);
  error ForkedChain(uint256 expected, uint256 actual);
  error WrongNumberOfSignatures();
  error SignaturesOutOfRegistration();
  error UnauthorizedTransmitter();
  error UnauthorizedSigner();
  error NonUniqueSignatures();
  error OracleCannotBeZeroAddress();

  // Packing these fields used on the hot path in a ConfigInfo variable reduces the
  // retrieval of all of them to a minimum number of SLOADs.
  struct ConfigInfo {
    bytes32 latestConfigDigest;
    uint8 f;
    uint8 n;
  }

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

  // The current config
  ConfigInfo internal s_configInfo;

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  // makes it easier for offchain systems to extract config from logs.
  uint32 internal s_latestConfigBlockNumber;

  // signer OR transmitter address
  mapping(address signerOrTransmitter => Oracle oracle) internal s_oracles;

  // s_signers contains the signing address of each oracle
  address[] internal s_signers;

  // s_transmitters contains the transmission address of each oracle,
  // i.e. the address the oracle actually sends transactions to the contract from
  address[] internal s_transmitters;

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
      32; // word containing length of ss

  bool internal immutable i_uniqueReports;
  uint256 internal immutable i_chainID;

  constructor(bool uniqueReports) {
    i_uniqueReports = uniqueReports;
    i_chainID = block.chainid;
  }

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

  /// @notice sets offchain reporting protocol configuration incl. participating oracles
  /// @param signers addresses with which oracles sign the reports
  /// @param transmitters addresses oracles use to transmit the reports
  /// @param f number of faulty oracles the system can tolerate
  /// @param onchainConfig encoded on-chain contract configuration
  /// @param offchainConfigVersion version number for offchainEncoding schema
  /// @param offchainConfig encoded off-chain oracle configuration
  function setOCR2Config(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override checkConfigValid(signers.length, transmitters.length, f) onlyOwner {
    _beforeSetConfig(onchainConfig);
    uint256 oldSignerLength = s_signers.length;
    for (uint256 i = 0; i < oldSignerLength; ++i) {
      delete s_oracles[s_signers[i]];
      delete s_oracles[s_transmitters[i]];
    }

    uint256 newSignersLength = signers.length;
    for (uint256 i = 0; i < newSignersLength; ++i) {
      // add new signer/transmitter addresses
      address signer = signers[i];
      if (s_oracles[signer].role != Role.Unset) revert InvalidConfig("repeated signer address");
      if (signer == address(0)) revert OracleCannotBeZeroAddress();
      s_oracles[signer] = Oracle(uint8(i), Role.Signer);

      address transmitter = transmitters[i];
      if (s_oracles[transmitter].role != Role.Unset) revert InvalidConfig("repeated transmitter address");
      if (transmitter == address(0)) revert OracleCannotBeZeroAddress();
      s_oracles[transmitter] = Oracle(uint8(i), Role.Transmitter);
    }

    s_signers = signers;
    s_transmitters = transmitters;

    s_configInfo.f = f;
    s_configInfo.n = uint8(newSignersLength);
    s_configInfo.latestConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      ++s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);

    emit ConfigSet(
      previousConfigBlockNumber,
      s_configInfo.latestConfigDigest,
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  /// @dev Hook that is run from setOCR2Config() right after validating configuration.
  /// Empty by default, please provide an implementation in a child contract if you need additional configuration processing
  function _beforeSetConfig(bytes memory _onchainConfig) internal virtual {}

  /// @return list of addresses permitted to transmit reports to this contract
  /// @dev The list will match the order used to specify the transmitter during setConfig
  function getTransmitters() external view returns (address[] memory) {
    return s_transmitters;
  }

  /// @notice transmit is called to post a new report to the contract
  /// @param report serialized report, which the signatures are signing.
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param rawVs ith element is the the V component of the ith signature
  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs // signatures
  ) external override {
    // Scoping this reduces stack pressure and gas usage
    {
      // report and epochAndRound
      _report(report, uint40(uint256(reportContext[1])));
    }

    // reportContext consists of:
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[2]: ExtraHash
    bytes32 configDigest = reportContext[0];
    ConfigInfo memory configInfo = s_configInfo;

    if (configInfo.latestConfigDigest != configDigest)
      revert ConfigDigestMismatch(configInfo.latestConfigDigest, configDigest);
    // If the cached chainID at time of deployment doesn't match the current chainID, we reject all signed reports.
    // This avoids a (rare) scenario where chain A forks into chain A and A', A' still has configDigest
    // calculated from chain A and so OCR reports will be valid on both forks.
    if (i_chainID != block.chainid) revert ForkedChain(i_chainID, block.chainid);

    emit Transmitted(configDigest, uint32(uint256(reportContext[1]) >> 8));

    uint256 expectedNumSignatures;
    if (i_uniqueReports) {
      expectedNumSignatures = (configInfo.n + configInfo.f) / 2 + 1;
    } else {
      expectedNumSignatures = configInfo.f + 1;
    }
    if (rs.length != expectedNumSignatures) revert WrongNumberOfSignatures();
    if (rs.length != ss.length) revert SignaturesOutOfRegistration();

    // Scoping this reduces stack pressure and gas usage
    {
      Oracle memory transmitter = s_oracles[msg.sender];
      // Check that sender is authorized to report
      if (!(transmitter.role == Role.Transmitter && msg.sender == s_transmitters[transmitter.index]))
        revert UnauthorizedTransmitter();
    }
    // Scoping this reduces stack pressure and gas usage
    {
      uint256 expectedDataLength = uint256(TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT) +
        report.length + // one byte pure entry in _report
        rs.length *
        32 + // 32 bytes per entry in _rs
        ss.length *
        32; // 32 bytes per entry in _ss)
      if (msg.data.length != expectedDataLength) revert WrongMessageLength(expectedDataLength, msg.data.length);
    }

    // Verify signatures attached to report
    bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));
    bool[MAX_NUM_ORACLES] memory signed;

    uint256 numberOfSignatures = rs.length;
    for (uint256 i = 0; i < numberOfSignatures; ++i) {
      // Safe from ECDSA malleability here since we check for duplicate signers.
      address signer = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
      // Since we disallow address(0) as a valid signer address, it can
      // never have a signer role.
      Oracle memory oracle = s_oracles[signer];
      if (oracle.role != Role.Signer) revert UnauthorizedSigner();
      if (signed[oracle.index]) revert NonUniqueSignatures();
      signed[oracle.index] = true;
    }
  }

  /// @notice information about current offchain reporting protocol configuration
  /// @return configCount ordinal number of current config, out of all configs applied to this contract so far
  /// @return blockNumber block at which this config was set
  /// @return configDigest domain-separation tag for current config (see _configDigestFromConfigData)
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_configCount, s_latestConfigBlockNumber, s_configInfo.latestConfigDigest);
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

  function _report(bytes calldata report, uint40 epochAndRound) internal virtual;
}
