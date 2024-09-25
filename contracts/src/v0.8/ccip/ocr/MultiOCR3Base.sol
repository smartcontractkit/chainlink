// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

/// @notice Onchain verification of reports from the offchain reporting protocol
///         with multiple OCR plugin support.
abstract contract MultiOCR3Base is ITypeAndVersion, OwnerIsCreator {
  // Maximum number of oracles the offchain reporting protocol is designed for
  uint256 internal constant MAX_NUM_ORACLES = 256;

  /// @notice Triggers a new run of the offchain reporting protocol
  /// @param ocrPluginType OCR plugin type for which the config was set
  /// @param configDigest configDigest of this configuration
  /// @param signers ith element is address ith oracle uses to sign a report
  /// @param transmitters ith element is address ith oracle uses to transmit a report via the transmit method
  /// @param F maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly
  event ConfigSet(uint8 ocrPluginType, bytes32 configDigest, address[] signers, address[] transmitters, uint8 F);

  /// @notice Optionally emitted to indicate the latest configDigest and sequence number
  /// for which a report was successfully transmitted. Alternatively, the contract may
  /// use latestConfigDigestAndEpoch with scanLogs set to false.
  event Transmitted(uint8 indexed ocrPluginType, bytes32 configDigest, uint64 sequenceNumber);

  enum InvalidConfigErrorType {
    F_MUST_BE_POSITIVE,
    TOO_MANY_TRANSMITTERS,
    TOO_MANY_SIGNERS,
    F_TOO_HIGH,
    REPEATED_ORACLE_ADDRESS,
    NO_TRANSMITTERS
  }

  error InvalidConfig(InvalidConfigErrorType errorType);
  error WrongMessageLength(uint256 expected, uint256 actual);
  error ConfigDigestMismatch(bytes32 expected, bytes32 actual);
  error ForkedChain(uint256 expected, uint256 actual);
  error WrongNumberOfSignatures();
  error SignaturesOutOfRegistration();
  error UnauthorizedTransmitter();
  error UnauthorizedSigner();
  error NonUniqueSignatures();
  error OracleCannotBeZeroAddress();
  error StaticConfigCannotBeChanged(uint8 ocrPluginType);

  /// @dev Packing these fields used on the hot path in a ConfigInfo variable reduces the
  ///      retrieval of all of them to a minimum number of SLOADs.
  struct ConfigInfo {
    bytes32 configDigest;
    uint8 F; // ──────────────────────────────╮ maximum number of faulty/dishonest oracles the system can tolerate
    uint8 n; //                               │ number of signers / transmitters
    bool isSignatureVerificationEnabled; // ──╯ if true, requires signers and verifies signatures on transmission verification
  }

  /// @notice Used for s_oracles[a].role, where a is an address, to track the purpose
  ///         of the address, or to indicate that the address is unset.
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
    uint8 index; // ───╮ Index of oracle in s_signers/s_transmitters
    Role role; // ─────╯ Role of the address which mapped to this struct
  }

  /// @notice OCR configuration for a single OCR plugin within a DON.
  struct OCRConfig {
    ConfigInfo configInfo; //  latest OCR config
    // NOTE: len(signers) can be different from len(transmitters). There is no index relationship between the two arrays
    address[] signers; //      addresses oracles use to sign the reports
    address[] transmitters; // addresses oracles use to transmit the reports
  }

  /// @notice Args to update an OCR Config.
  struct OCRConfigArgs {
    bytes32 configDigest; // Config digest to update to
    uint8 ocrPluginType; // ──────────────────╮ OCR plugin type to update config for
    uint8 F; //                               │ maximum number of faulty/dishonest oracles
    bool isSignatureVerificationEnabled; // ──╯ if true, requires signers and verifies signatures on transmission verification
    address[] signers; // signing address of each oracle
    address[] transmitters; // transmission address of each oracle (i.e. the address the oracle actually sends transactions to the contract from)
  }

  /// @notice mapping of OCR plugin type -> DON config
  mapping(uint8 ocrPluginType => OCRConfig config) internal s_ocrConfigs;

  /// @notice OCR plugin type => signer OR transmitter address mapping
  mapping(uint8 ocrPluginType => mapping(address signerOrTransmiter => Oracle oracle)) internal s_oracles;

  // Constant-length components of the msg.data sent to transmit.
  // See the "If we wanted to call sam" example on for example reasoning
  // https://solidity.readthedocs.io/en/v0.7.2/abi-spec.html

  /// @notice Constant length component for transmit functions with no signatures.
  /// The signatures are expected to match transmitPlugin(reportContext, report).
  uint16 private constant TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT_NO_SIGNATURES = 4 // function selector
    + 3 * 32 // 3 words containing reportContext
    + 32 // word containing start location of abiencoded report value
    + 32; // word containing length of report

  /// @notice Extra constant length component for transmit functions with signatures (relative to no signatures).
  /// The signatures are expected to match transmitPlugin(reportContext, report, rs, ss, rawVs).
  uint16 private constant TRANSMIT_MSGDATA_EXTRA_CONSTANT_LENGTH_COMPONENT_FOR_SIGNATURES = 32 // word containing location start of abiencoded rs value
    + 32 // word containing start location of abiencoded ss value
    + 32 // rawVs value
    + 32 // word containing length rs
    + 32; // word containing length of ss

  uint256 internal immutable i_chainID;

  constructor() {
    i_chainID = block.chainid;
  }

  /// @notice Sets offchain reporting protocol configuration incl. participating oracles.
  /// NOTE: The OCR3 config must be sanity-checked against the home-chain registry configuration, to ensure
  /// home-chain and remote-chain parity!
  /// @param ocrConfigArgs OCR config update args.
  /// @dev precondition number of transmitters should match the expected F/fChain relationship.
  /// For transmitters, the function only validates that len(transmitters) > 0 && len(transmitters) <= MAX_NUM_ORACLES
  /// && len(transmitters) <= len(signers) [if sig verification is enabled]
  function setOCR3Configs(OCRConfigArgs[] memory ocrConfigArgs) external onlyOwner {
    for (uint256 i; i < ocrConfigArgs.length; ++i) {
      _setOCR3Config(ocrConfigArgs[i]);
    }
  }

  /// @notice Sets offchain reporting protocol configuration incl. participating oracles for a single OCR plugin type.
  /// @param ocrConfigArgs OCR config update args.
  function _setOCR3Config(OCRConfigArgs memory ocrConfigArgs) internal {
    if (ocrConfigArgs.F == 0) revert InvalidConfig(InvalidConfigErrorType.F_MUST_BE_POSITIVE);

    uint8 ocrPluginType = ocrConfigArgs.ocrPluginType;
    OCRConfig storage ocrConfig = s_ocrConfigs[ocrPluginType];
    ConfigInfo storage configInfo = ocrConfig.configInfo;

    // If F is 0, then the config is not yet set
    if (configInfo.F == 0) {
      configInfo.isSignatureVerificationEnabled = ocrConfigArgs.isSignatureVerificationEnabled;
    } else if (configInfo.isSignatureVerificationEnabled != ocrConfigArgs.isSignatureVerificationEnabled) {
      revert StaticConfigCannotBeChanged(ocrPluginType);
    }

    address[] memory transmitters = ocrConfigArgs.transmitters;
    if (transmitters.length > MAX_NUM_ORACLES) revert InvalidConfig(InvalidConfigErrorType.TOO_MANY_TRANSMITTERS);
    if (transmitters.length == 0) revert InvalidConfig(InvalidConfigErrorType.NO_TRANSMITTERS);

    _clearOracleRoles(ocrPluginType, ocrConfig.transmitters);

    if (ocrConfigArgs.isSignatureVerificationEnabled) {
      _clearOracleRoles(ocrPluginType, ocrConfig.signers);

      address[] memory signers = ocrConfigArgs.signers;

      if (signers.length > MAX_NUM_ORACLES) revert InvalidConfig(InvalidConfigErrorType.TOO_MANY_SIGNERS);
      if (signers.length <= 3 * ocrConfigArgs.F) revert InvalidConfig(InvalidConfigErrorType.F_TOO_HIGH);
      // NOTE: Transmitters cannot exceed signers. Transmitters do not have to be >= 3F + 1 because they can match >= 3fChain + 1, where fChain <= F.
      // fChain is not represented in MultiOCR3Base - so we skip this check.
      if (signers.length < transmitters.length) revert InvalidConfig(InvalidConfigErrorType.TOO_MANY_TRANSMITTERS);

      configInfo.n = uint8(signers.length);
      ocrConfig.signers = signers;

      _assignOracleRoles(ocrPluginType, signers, Role.Signer);
    }

    _assignOracleRoles(ocrPluginType, transmitters, Role.Transmitter);

    ocrConfig.transmitters = transmitters;
    configInfo.F = ocrConfigArgs.F;
    configInfo.configDigest = ocrConfigArgs.configDigest;

    emit ConfigSet(
      ocrPluginType, ocrConfigArgs.configDigest, ocrConfig.signers, ocrConfigArgs.transmitters, ocrConfigArgs.F
    );
    _afterOCR3ConfigSet(ocrPluginType);
  }

  /// @notice Hook that is called after a plugin's OCR3 config changes.
  /// @param ocrPluginType Plugin type for which the config changed.
  function _afterOCR3ConfigSet(uint8 ocrPluginType) internal virtual;

  /// @notice Clears oracle roles for the provided oracle addresses.
  /// @param ocrPluginType OCR plugin type to clear roles for.
  /// @param oracleAddresses Oracle addresses to clear roles for.
  function _clearOracleRoles(uint8 ocrPluginType, address[] memory oracleAddresses) internal {
    for (uint256 i = 0; i < oracleAddresses.length; ++i) {
      delete s_oracles[ocrPluginType][oracleAddresses[i]];
    }
  }

  /// @notice Assigns oracles roles for the provided oracle addresses with uniqueness verification.
  /// @param ocrPluginType OCR plugin type to assign roles for.
  /// @param oracleAddresses Oracle addresses to assign roles to.
  /// @param role Role to assign.
  function _assignOracleRoles(uint8 ocrPluginType, address[] memory oracleAddresses, Role role) internal {
    for (uint256 i = 0; i < oracleAddresses.length; ++i) {
      address oracle = oracleAddresses[i];
      if (s_oracles[ocrPluginType][oracle].role != Role.Unset) {
        revert InvalidConfig(InvalidConfigErrorType.REPEATED_ORACLE_ADDRESS);
      }
      if (oracle == address(0)) revert OracleCannotBeZeroAddress();
      s_oracles[ocrPluginType][oracle] = Oracle(uint8(i), role);
    }
  }

  /// @notice _transmit is called to post a new report to the contract.
  ///         The function should be called after the per-DON reporting logic is completed.
  /// @param ocrPluginType OCR plugin type to transmit report for
  /// @param report serialized report, which the signatures are signing.
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries.
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries.
  /// @param rawVs ith element is the the V component of the ith signature.
  function _transmit(
    uint8 ocrPluginType,
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] memory rs,
    bytes32[] memory ss,
    bytes32 rawVs // signatures
  ) internal {
    // reportContext consists of:
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 24 byte padding, 8 byte sequence number
    // reportContext[2]: ExtraHash
    ConfigInfo memory configInfo = s_ocrConfigs[ocrPluginType].configInfo;
    bytes32 configDigest = reportContext[0];

    // Scoping this reduces stack pressure and gas usage
    {
      uint256 expectedDataLength = uint256(TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT_NO_SIGNATURES) + report.length; // one byte pure entry in _report

      if (configInfo.isSignatureVerificationEnabled) {
        expectedDataLength += TRANSMIT_MSGDATA_EXTRA_CONSTANT_LENGTH_COMPONENT_FOR_SIGNATURES + rs.length * 32 // 32 bytes per entry in _rs
          + ss.length * 32; // 32 bytes per entry in _ss)
      }

      if (msg.data.length != expectedDataLength) revert WrongMessageLength(expectedDataLength, msg.data.length);
    }

    if (configInfo.configDigest != configDigest) {
      revert ConfigDigestMismatch(configInfo.configDigest, configDigest);
    }
    // If the cached chainID at time of deployment doesn't match the current chainID, we reject all signed reports.
    // This avoids a (rare) scenario where chain A forks into chain A and A', A' still has configDigest
    // calculated from chain A and so OCR reports will be valid on both forks.
    _whenChainNotForked();

    // Scoping this reduces stack pressure and gas usage
    {
      Oracle memory transmitter = s_oracles[ocrPluginType][msg.sender];
      // Check that sender is authorized to report
      if (
        !(
          transmitter.role == Role.Transmitter
            && msg.sender == s_ocrConfigs[ocrPluginType].transmitters[transmitter.index]
        )
      ) {
        revert UnauthorizedTransmitter();
      }
    }

    if (configInfo.isSignatureVerificationEnabled) {
      // Scoping to reduce stack pressure
      {
        if (rs.length != configInfo.F + 1) revert WrongNumberOfSignatures();
        if (rs.length != ss.length) revert SignaturesOutOfRegistration();
      }

      bytes32 h = keccak256(abi.encodePacked(keccak256(report), reportContext));
      _verifySignatures(ocrPluginType, h, rs, ss, rawVs);
    }

    emit Transmitted(ocrPluginType, configDigest, uint64(uint256(reportContext[1])));
  }

  /// @notice Verifies the signatures of a hashed report value for one OCR plugin type.
  /// @param ocrPluginType OCR plugin type to transmit report for.
  /// @param hashedReport hashed encoded packing of report + reportContext.
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries.
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries.
  /// @param rawVs ith element is the the V component of the ith signature.
  function _verifySignatures(
    uint8 ocrPluginType,
    bytes32 hashedReport,
    bytes32[] memory rs,
    bytes32[] memory ss,
    bytes32 rawVs
  ) internal view {
    // Verify signatures attached to report. Using a uint256 means we can only verify up to 256 oracles.
    uint256 signed = 0;

    uint256 numberOfSignatures = rs.length;
    for (uint256 i; i < numberOfSignatures; ++i) {
      // Safe from ECDSA malleability here since we check for duplicate signers.
      address signer = ecrecover(hashedReport, uint8(rawVs[i]) + 27, rs[i], ss[i]);
      // Since we disallow address(0) as a valid signer address, it can never have a signer role.
      Oracle memory oracle = s_oracles[ocrPluginType][signer];
      if (oracle.role != Role.Signer) revert UnauthorizedSigner();
      if (signed & (0x1 << oracle.index) != 0) revert NonUniqueSignatures();
      signed |= 0x1 << oracle.index;
    }
  }

  /// @notice Validates that the chain ID has not diverged after deployment. Reverts if the chain IDs do not match.
  function _whenChainNotForked() internal view {
    if (i_chainID != block.chainid) revert ForkedChain(i_chainID, block.chainid);
  }

  /// @notice Information about current offchain reporting protocol configuration.
  /// @param ocrPluginType OCR plugin type to return config details for.
  /// @return ocrConfig OCR config for the plugin type.
  function latestConfigDetails(uint8 ocrPluginType) external view returns (OCRConfig memory ocrConfig) {
    return s_ocrConfigs[ocrPluginType];
  }
}
