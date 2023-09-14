// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {IStreamsConfig} from "./interfaces/IStreamsConfig.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "./libraries/Common.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/*
 * The verifier contract is used to verify offchain reports signed
 * by DONs. A report consists of a price, block number and feed Id. It
 * represents the observed price of an asset at a specified block number for
 * a feed. The verifier contract is used to verify that such reports have
 * been signed by the correct signers.
 **/
contract StreamsConfig is IStreamsConfig, ConfirmedOwner, TypeAndVersionInterface {
  // The first byte of the mask can be 0, because we only ever have 31 oracles
  uint256 internal constant ORACLE_MASK = 0x0001010101010101010101010101010101010101010101010101010101010101;

  enum Role {
    // Default role for an oracle address.  This means that the oracle address
    // is not a signer
    Unset,
    // Role given to an oracle address that is allowed to sign feed data
    Signer
  }

  struct Signer {
    // Index of oracle in a configuration
    uint8 index;
    // The oracle's role
    Role role;
  }

  struct Config {
    // Fault tolerance
    uint8 f;
    // Marks whether or not a configuration is active
    bool isActive;
    // Map of signer addresses to oracles
    mapping(address => Signer) oracles;
  }

  struct VerifierState {
    // The number of times a new configuration
    /// has been set
    uint32 configCount;
    // The block number of the block the last time
    /// the configuration was updated.
    uint32 latestConfigBlockNumber;
    // The latest epoch a report was verified for
    uint32 latestEpoch;
    // Whether or not the verifier for this feed has been deactivated
    bool isDeactivated;
    /// The latest config digest set
    bytes32 latestConfigDigest;
    /// The previously set config
    Config s_verificationDataConfig;
  }

  /// @notice This event is emitted when a new report is verified.
  /// It is used to keep a historical record of verified reports.
  event ReportVerified(address requester);

  /// @notice This event is emitted whenever a new configuration is set.  It triggers a new run of the offchain reporting protocol.
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    bytes32[] offchainTransmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );

  /// @notice This event is emitted whenever a configuration is deactivated
  event ConfigDeactivated(bytes32 configDigest);

  /// @notice This event is emitted whenever a configuration is activated
  event ConfigActivated(bytes32 configDigest);

  /// @notice This error is thrown whenever an address tries
  /// to exeecute a transaction that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown whenever the config digest
  /// is empty
  error DigestEmpty();

  /// @notice This error is thrown whenever the config digest
  /// passed in has not been set in this verifier
  /// @param configDigest The config digest that has not been set
  error DigestNotSet(bytes32 configDigest);

  /// @notice This error is thrown whenever the config digest
  /// has been deactivated
  /// @param configDigest The config digest that is inactive
  error DigestInactive(bytes32 configDigest);

  /// @notice This error is thrown whenever trying to set a config
  /// with a fault tolerance of 0
  error FaultToleranceMustBePositive();

  /// @notice This error is thrown whenever a report is signed
  /// with more than the max number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param maxSigners The maximum number of signers that can sign a report
  error ExcessSigners(uint256 numSigners, uint256 maxSigners);

  /// @notice This error is thrown whenever a report is signed
  /// with less than the minimum number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param minSigners The minimum number of signers that need to sign a report
  error InsufficientSigners(uint256 numSigners, uint256 minSigners);

  /// @notice This error is thrown whenever a report is signed
  /// with an incorrect number of signers
  /// @param numSigners The number of signers who have signed the report
  /// @param expectedNumSigners The expected number of signers that need to sign
  /// a report
  error IncorrectSignatureCount(uint256 numSigners, uint256 expectedNumSigners);

  /// @notice This error is thrown whenever the R and S signer components
  /// have different lengths
  /// @param rsLength The number of r signature components
  /// @param ssLength The number of s signature components
  error MismatchedSignatures(uint256 rsLength, uint256 ssLength);

  /// @notice This error is thrown whenever setting a config with duplicate signatures
  error NonUniqueSignatures();

  /// @notice This error is thrown whenever a report fails to verify due to bad or duplicate signatures
  error BadVerification();

  /// @notice This error is thrown whenever the admin tries to deactivate
  /// the latest config digest
  /// @param configDigest The latest config digest
  error CannotDeactivateLatestConfig(bytes32 configDigest);

  /// @notice The address of the verifier proxy
  address private immutable i_verifierProxyAddr;

  /// @notice Verifier state
  VerifierState internal s_feedVerifierState;

  /// @param verifierProxyAddr The address of the VerifierProxy contract
  constructor(address verifierProxyAddr) ConfirmedOwner(msg.sender) {
    if (verifierProxyAddr == address(0)) revert ZeroAddress();
    i_verifierProxyAddr = verifierProxyAddr;
  }

  modifier checkConfigValid(uint256 numSigners, uint256 f) {
    if (f == 0) revert FaultToleranceMustBePositive();
    if (numSigners > MAX_NUM_ORACLES) revert ExcessSigners(numSigners, MAX_NUM_ORACLES);
    if (numSigners <= 3 * f) revert InsufficientSigners(numSigners, 3 * f + 1);
    _;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool isVerifier) {
    return interfaceId == this.verify.selector;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "Verifier 1.2.0";
  }

  /// @inheritdoc IVerifier
  function verify(
    bytes calldata signedReport,
    address sender
  ) external override returns (bytes memory verifierResponse) {
    if (msg.sender != i_verifierProxyAddr) revert AccessForbidden();
    (
      bytes32[3] memory reportContext,
      bytes memory reportData,
      bytes32[] memory rs,
      bytes32[] memory ss,
      bytes32 rawVs
    ) = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

    VerifierState storage feedVerifierState = s_feedVerifierState;

    // If the feed has been deactivated, do not verify the report
    if (feedVerifierState.isDeactivated) {
      revert DigestInactive();
    }

    // reportContext consists of:
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[2]: ExtraHash
    bytes32 configDigest = reportContext[0];
    Config storage s_config = feedVerifierState.s_verificationDataConfigs[configDigest];

    _validateReport( configDigest, rs, ss, s_config);
    _updateEpoch(reportContext, feedVerifierState);

    bytes32 hashedReport = keccak256(reportData);

    _verifySignatures(hashedReport, reportContext, rs, ss, rawVs, s_config);
    emit ReportVerified( sender);

    return reportData;
  }

  /// @notice Validates parameters of the report
  /// @param configDigest Config digest from the report
  /// @param rs R components from the report
  /// @param ss S components from the report
  /// @param config Config for the given feed ID keyed on the config digest
  function _validateReport(
    bytes32 configDigest,
    bytes32[] memory rs,
    bytes32[] memory ss,
    Config storage config
  ) private view {
    uint8 expectedNumSignatures = config.f + 1;

    if (!config.isActive) revert DigestInactive( configDigest);
    if (rs.length != expectedNumSignatures) revert IncorrectSignatureCount(rs.length, expectedNumSignatures);
    if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);
  }

  /**
   * @notice Conditionally update the epoch for a feed
   * @param reportContext Report context containing the epoch and round
   * @param feedVerifierState Feed verifier state to conditionally update
   */
  function _updateEpoch(bytes32[3] memory reportContext, VerifierState storage feedVerifierState) private {
    uint40 epochAndRound = uint40(uint256(reportContext[1]));
    uint32 epoch = uint32(epochAndRound >> 8);
    if (epoch > feedVerifierState.latestEpoch) {
      feedVerifierState.latestEpoch = epoch;
    }
  }

  /// @notice Verifies that a report has been signed by the correct
  /// signers and that enough signers have signed the reports.
  /// @param hashedReport The keccak256 hash of the raw report's bytes
  /// @param reportContext The context the report was signed in
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param rawVs ith element is the the V component of the ith signature
  /// @param s_config The config digest the report was signed for
  function _verifySignatures(
    bytes32 hashedReport,
    bytes32[3] memory reportContext,
    bytes32[] memory rs,
    bytes32[] memory ss,
    bytes32 rawVs,
    Config storage s_config
  ) private view {
    bytes32 h = keccak256(abi.encodePacked(hashedReport, reportContext));
    // i-th byte counts number of sigs made by i-th signer
    uint256 signedCount;

    Signer memory o;
    address signerAddress;
    uint256 numSigners = rs.length;
    for (uint256 i; i < numSigners; ++i) {
      signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
      o = s_config.oracles[signerAddress];
      if (o.role != Role.Signer) revert BadVerification();
      unchecked {
        signedCount += 1 << (8 * o.index);
      }
    }

    if (signedCount & ORACLE_MASK != signedCount) revert BadVerification();
  }

  /// @inheritdoc IVerifier
  function setConfig(
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(
      block.chainid,
      address(this),
      0, // 0 defaults to feedConfig.configCount + 1
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig,
      recipientAddressesAndWeights
    );
  }

  /// @inheritdoc IVerifier
  function setConfigFromSource(
    uint256 sourceChainId,
    address sourceAddress,
    uint32 newConfigCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(
      sourceChainId,
      sourceAddress,
      newConfigCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig,
      recipientAddressesAndWeights
    );
  }

  /// @notice Sets config based on the given arguments
  /// @param sourceChainId Chain ID of source config
  /// @param sourceAddress Address of source config Verifier
  /// @param newConfigCount Optional param to force the new config count
  /// @param signers addresses with which oracles sign the reports
  /// @param offchainTransmitters CSA key for the ith Oracle
  /// @param f number of faulty oracles the system can tolerate
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version number for offchainEncoding schema
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  /// @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
  function _setConfig(
    uint256 sourceChainId,
    address sourceAddress,
    uint32 newConfigCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) internal {
    VerifierState storage feedVerifierState = s_feedVerifierState;

    // Increment the number of times a config has been set first
    if (newConfigCount > 0) feedVerifierState.configCount = newConfigCount;
    else feedVerifierState.configCount++;

    bytes32 configDigest = _configDigestFromConfigData(
      sourceChainId,
      sourceAddress,
      feedVerifierState.configCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    feedVerifierState.s_verificationDataConfigs[configDigest].f = f;
    feedVerifierState.s_verificationDataConfigs[configDigest].isActive = true;
    for (uint8 i; i < signers.length; ++i) {
      address signerAddr = signers[i];
      if (signerAddr == address(0)) revert ZeroAddress();

      // All signer roles are unset by default for a new config digest.
      // Here the contract checks to see if a signer's address has already
      // been set to ensure that the group of signer addresses that will
      // sign reports with the config digest are unique.
      bool isSignerAlreadySet = feedVerifierState.s_verificationDataConfigs[configDigest].oracles[signerAddr].role !=
        Role.Unset;
      if (isSignerAlreadySet) revert NonUniqueSignatures();
      feedVerifierState.s_verificationDataConfigs[configDigest].oracles[signerAddr] = Signer({
        role: Role.Signer,
        index: i
      });
    }

    IVerifierProxy(i_verifierProxyAddr).setVerifier(
      feedVerifierState.latestConfigDigest,
      configDigest,
      recipientAddressesAndWeights
    );

    emit ConfigSet(
      feedVerifierState.latestConfigBlockNumber,
      configDigest,
      feedVerifierState.configCount,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    feedVerifierState.latestEpoch = 0;
    feedVerifierState.latestConfigBlockNumber = uint32(block.number);
    feedVerifierState.latestConfigDigest = configDigest;
  }

  /// @notice Generates the config digest from config data
  /// @param sourceChainId Chain ID of source config
  /// @param sourceAddress Address of source config Verifier
  /// @param configCount ordinal number of this config setting among all config settings over the life of this contract
  /// @param signers ith element is address ith oracle uses to sign a report
  /// @param offchainTransmitters ith element is address ith oracle used to transmit reports (in this case used for flexible additional field, such as CSA pub keys)
  /// @param f maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly
  /// @param onchainConfig serialized configuration used by the contract (and possibly oracles)
  /// @param offchainConfigVersion version of the serialization format used for "offchainConfig" parameter
  /// @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
  /// @dev This function is a modified version of the method from OCR2Abstract
  function _configDigestFromConfigData(
    uint256 sourceChainId,
    address sourceAddress,
    uint64 configCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          sourceChainId,
          sourceAddress,
          configCount,
          signers,
          offchainTransmitters,
          f,
          onchainConfig,
          offchainConfigVersion,
          offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    // 0x0006 corresponds to ConfigDigestPrefixMercuryV02 in libocr
    uint256 prefix = 0x0006 << (256 - 16); // 0x000600..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  /// @inheritdoc IVerifier
  function activateConfig(bytes32 configDigest) external onlyOwner {
    VerifierState storage feedVerifierState = s_feedVerifierState;

    if (configDigest == bytes32("")) revert DigestEmpty();
    if (feedVerifierState.s_verificationDataConfigs[configDigest].f == 0) revert DigestNotSet(configDigest);
    feedVerifierState.s_verificationDataConfigs[configDigest].isActive = true;
    emit ConfigActivated(configDigest);
  }

  /// @inheritdoc IVerifier
  function deactivateConfig(bytes32 configDigest) external onlyOwner {
    VerifierState storage feedVerifierState = s_feedVerifierState;

    if (configDigest == bytes32("")) revert DigestEmpty();
    if (feedVerifierState.s_verificationDataConfigs[configDigest].f == 0) revert DigestNotSet(configDigest);
    if (configDigest == feedVerifierState.latestConfigDigest) {
      revert CannotDeactivateLatestConfig(configDigest);
    }
    feedVerifierState.s_verificationDataConfigs[configDigest].isActive = false;
    emit ConfigDeactivated(configDigest);
  }

  /// @inheritdoc IVerifier
  function latestConfigDigestAndEpoch(
  ) external view override returns (bool scanLogs, bytes32 configDigest, uint32 epoch) {
    VerifierState storage feedVerifierState = s_feedVerifierState;
    return (false, feedVerifierState.latestConfigDigest, feedVerifierState.latestEpoch);
  }

  /// @inheritdoc IVerifier
  function latestConfigDetails(
  ) external view override returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest) {
    VerifierState storage feedVerifierState = s_feedVerifierState;
    return (
      feedVerifierState.configCount,
      feedVerifierState.latestConfigBlockNumber,
      feedVerifierState.latestConfigDigest
    );
  }
}
