// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwner} from "../../../../shared/access/ConfirmedOwner.sol";
import {OCR2Abstract} from "./OCR2Abstract.sol";

/**
 * @notice Onchain verification of reports from the offchain reporting protocol
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 * @dev For details on its operation, see the offchain reporting protocol design
 * doc, which refers to this contract as simply the "contract".
 * @dev This contract is meant to aid rapid development of new applications based on OCR2.
 * However, for actual production contracts, it is expected that most of the logic of this contract
 * will be folded directly into the application contract. Inheritance prevents us from doing lots
 * of juicy storage layout optimizations, leading to a substantial increase in gas cost.
 */
abstract contract OCR2Base is ConfirmedOwner, OCR2Abstract {
  error ReportInvalid();

  bool internal immutable i_uniqueReports;

  constructor(bool uniqueReports) ConfirmedOwner(msg.sender) {
    i_uniqueReports = uniqueReports;
  }

  uint256 private constant maxUint32 = (1 << 32) - 1;

  // Storing these fields used on the hot path in a ConfigInfo variable reduces the
  // retrieval of all of them to a single SLOAD. If any further fields are
  // added, make sure that storage of the struct still takes at most 32 bytes.
  struct ConfigInfo {
    bytes32 latestConfigDigest;
    uint8 f; // TODO: could be optimized by squeezing into one slot
    uint8 n;
  }
  ConfigInfo internal s_configInfo;

  // incremented each time a new config is posted. This count is incorporated
  // into the config digest, to prevent replay attacks.
  uint32 internal s_configCount;
  uint32 internal s_latestConfigBlockNumber; // makes it easier for offchain systems
  // to extract config from logs.

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

  mapping(address => Oracle) /* signer OR transmitter address */ internal s_oracles;

  // s_signers contains the signing address of each oracle
  address[] internal s_signers;

  // s_transmitters contains the transmission address of each oracle,
  // i.e. the address the oracle actually sends transactions to the contract from
  address[] internal s_transmitters;

  /*
   * Config logic
   */

  // Reverts transaction if config args are invalid
  modifier checkConfigValid(
    uint256 _numSigners,
    uint256 _numTransmitters,
    uint256 _f
  ) {
    require(_numSigners <= maxNumOracles, "too many signers");
    require(_f > 0, "f must be positive");
    require(_numSigners == _numTransmitters, "oracle addresses out of registration");
    require(_numSigners > 3 * _f, "faulty-oracle f too high");
    _;
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

    for (uint256 i = 0; i < args.signers.length; ++i) {
      // add new signer/transmitter addresses
      require(s_oracles[args.signers[i]].role == Role.Unset, "repeated signer address");
      s_oracles[args.signers[i]] = Oracle(uint8(i), Role.Signer);
      require(s_oracles[args.transmitters[i]].role == Role.Unset, "repeated transmitter address");
      s_oracles[args.transmitters[i]] = Oracle(uint8(i), Role.Transmitter);
      s_signers.push(args.signers[i]);
      s_transmitters.push(args.transmitters[i]);
    }
    s_configInfo.f = args.f;
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;
    {
      s_configInfo.latestConfigDigest = configDigestFromConfigData(
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

    _afterSetConfig(args.f, args.onchainConfig);
  }

  function configDigestFromConfigData(
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
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
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

  function _afterSetConfig(uint8 _f, bytes memory _onchainConfig) internal virtual;

  /**
   * @dev hook to allow additional validation of the report by the extending contract
   * @param configDigest separation tag for current config (see configDigestFromConfigData)
   * @param epochAndRound 27 byte padding, 4-byte epoch and 1-byte round
   * @param report serialized report
   */
  function _validateReport(
    bytes32 configDigest,
    uint40 epochAndRound,
    bytes memory report
  ) internal virtual returns (bool);

  /**
   * @dev hook called after the report has been fully validated
   * for the extending contract to handle additional logic, such as oracle payment
   * @param initialGas the amount of gas before validation
   * @param transmitter the address of the account that submitted the report
   * @param signers the addresses of all signing accounts
   * @param report serialized report
   */
  function _report(
    uint256 initialGas,
    address transmitter,
    uint8 signerCount,
    address[maxNumOracles] memory signers,
    bytes calldata report
  ) internal virtual;

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

  function requireExpectedMsgDataLength(
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
    require(msg.data.length == expected, "calldata length mismatch");
  }

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
    uint256 initialGas = gasleft(); // This line must come first

    {
      // reportContext consists of:
      // reportContext[0]: ConfigDigest
      // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
      // reportContext[2]: ExtraHash
      bytes32 configDigest = reportContext[0];
      uint32 epochAndRound = uint32(uint256(reportContext[1]));

      if (!_validateReport(configDigest, epochAndRound, report)) {
        revert ReportInvalid();
      }

      emit Transmitted(configDigest, uint32(epochAndRound >> 8));

      ConfigInfo memory configInfo = s_configInfo;
      require(configInfo.latestConfigDigest == configDigest, "configDigest mismatch");

      requireExpectedMsgDataLength(report, rs, ss);

      uint256 expectedNumSignatures;
      if (i_uniqueReports) {
        expectedNumSignatures = (configInfo.n + configInfo.f) / 2 + 1;
      } else {
        expectedNumSignatures = configInfo.f + 1;
      }

      require(rs.length == expectedNumSignatures, "wrong number of signatures");
      require(rs.length == ss.length, "signatures out of registration");

      Oracle memory transmitter = s_oracles[msg.sender];
      require( // Check that sender is authorized to report
        transmitter.role == Role.Transmitter && msg.sender == s_transmitters[transmitter.index],
        "unauthorized transmitter"
      );
    }

    address[maxNumOracles] memory signed;
    uint8 signerCount = 0;

    {
      // Verify signatures attached to report
      bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));

      Oracle memory o;
      for (uint256 i = 0; i < rs.length; ++i) {
        address signer = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
        o = s_oracles[signer];
        require(o.role == Role.Signer, "address not authorized to sign");
        require(signed[o.index] == address(0), "non-unique signature");
        signed[o.index] = signer;
        signerCount += 1;
      }
    }

    _report(initialGas, msg.sender, signerCount, signed, report);
  }
}
