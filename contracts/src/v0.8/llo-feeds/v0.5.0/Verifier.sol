// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IVerifier} from "./interfaces/IVerifier.sol";
import {IVerifierProxy} from "./interfaces/IVerifierProxy.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";

// OCR2 standard
uint256 constant MAX_NUM_ORACLES = 31;

/*
 * The verifier contract is used to verify offchain reports signed
 * by DONs. A report consists of a price, block number and feed Id. It
 * represents the observed price of an asset at a specified block number for
 * a feed. The verifier contract is used to verify that such reports have
 * been signed by the correct signers.
 **/
contract Verifier is IVerifier, ConfirmedOwner, ITypeAndVersion {
  // The first byte of the mask can be 0, because we only ever have 31 oracles
  uint256 internal constant ORACLE_MASK = 0x0001010101010101010101010101010101010101010101010101010101010101;

  enum Role {
    // Default role for an oracle address.  This means that the oracle address
    // is not a signer
    Unset,
    // Role given to an oracle address that is allowed to sign a report
    Signer
  }

  struct Signer {
    // Index of oracle in a configuration
    uint8 index;
    // The oracle's role
    Role role;
  }

  struct VerifierState {
    // The block number of the block the last time the configuration was updated.
    uint32 latestConfigBlockNumber;
    // Whether the config is deactivated
    bool isActive;
    // Fault tolerance
    uint8 f;
    // Number of signers
    uint8 oracleCount;
    // Map of signer addresses to oracles
    mapping(address => Signer) oracles;
  }

  /// @notice This event is emitted when a new report is verified.
  /// It is used to keep a historical record of verified reports.
  event ReportVerified(bytes32 indexed feedId, address requester);

  /// @notice This event is emitted whenever a new DON configuration is set.
  event ConfigSet(bytes32 indexed configDigest, address[] signers, uint8 f);

  /// @notice This event is
  event ConfigUpdated(bytes32 indexed configDigest, address[] prevSigners, address[] newSigners);

  /// @notice This event is emitted whenever a configuration is deactivated
  event ConfigDeactivated(bytes32 indexed configDigest);

  /// @notice This event is emitted whenever a configuration is activated
  event ConfigActivated(bytes32 indexed configDigest);

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

  /// @notice This error is thrown whenever a config digest is already set when setting the configuration
  error ConfigDigestAlreadySet();

  /// @notice The address of the verifier proxy
  address private immutable i_verifierProxyAddr;

  /// @notice Verifier states keyed on config digest
  mapping(bytes32 => VerifierState) internal s_verifierStates;

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

  /// @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "Verifier 2.0.0";
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

    // reportContext consists of:
    // reportContext[0]: ConfigDigest
    // reportContext[1]: 27 byte padding, 4-byte epoch and 1-byte round
    // reportContext[2]: ExtraHash
    bytes32 configDigest = reportContext[0];

    VerifierState storage verifierState = s_verifierStates[configDigest];

    _validateReport(configDigest, rs, ss, verifierState);

    bytes32 hashedReport = keccak256(reportData);

    _verifySignatures(hashedReport, reportContext, rs, ss, rawVs, verifierState);
    emit ReportVerified(bytes32(reportData), sender);

    return reportData;
  }

  /// @notice Validates parameters of the report
  /// @param configDigest Config digest from the report
  /// @param rs R components from the report
  /// @param ss S components from the report
  /// @param config Config for the given digest
  function _validateReport(
    bytes32 configDigest,
    bytes32[] memory rs,
    bytes32[] memory ss,
    VerifierState storage config
  ) private view {
    uint8 expectedNumSignatures = config.f + 1;

    if (!config.isActive) revert DigestInactive(configDigest);
    if (rs.length != expectedNumSignatures) revert IncorrectSignatureCount(rs.length, expectedNumSignatures);
    if (rs.length != ss.length) revert MismatchedSignatures(rs.length, ss.length);
  }

  /// @notice Verifies that a report has been signed by the correct
  /// signers and that enough signers have signed the reports.
  /// @param hashedReport The keccak256 hash of the raw report's bytes
  /// @param reportContext The context the report was signed in
  /// @param rs ith element is the R components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param ss ith element is the S components of the ith signature on report. Must have at most MAX_NUM_ORACLES entries
  /// @param rawVs ith element is the the V component of the ith signature
  /// @param config The config digest the report was signed for
  function _verifySignatures(
    bytes32 hashedReport,
    bytes32[3] memory reportContext,
    bytes32[] memory rs,
    bytes32[] memory ss,
    bytes32 rawVs,
    VerifierState storage config
  ) private view {
    bytes32 h = keccak256(abi.encodePacked(hashedReport, reportContext));
    // i-th byte counts number of sigs made by i-th signer
    uint256 signedCount;

    Signer memory o;
    address signerAddress;
    uint256 numSigners = rs.length;
    for (uint256 i; i < numSigners; ++i) {
      signerAddress = ecrecover(h, uint8(rawVs[i]) + 27, rs[i], ss[i]);
      o = config.oracles[signerAddress];
      if (o.role != Role.Signer) revert BadVerification();
      unchecked {
        signedCount += 1 << (8 * o.index);
      }
    }

    if (signedCount & ORACLE_MASK != signedCount) revert BadVerification();
  }

  /// @inheritdoc IVerifier
  function updateConfig(
    bytes32 configDigest,
    address[] calldata prevSigners,
    address[] calldata newSigners,
    uint8 f
  ) external override checkConfigValid(newSigners.length, f) onlyOwner {
    VerifierState storage config = s_verifierStates[configDigest];

    if (config.f == 0) revert DigestNotSet(configDigest);

    // We must be removing the number of signers that were originally set
    if (config.oracleCount != prevSigners.length) {
      revert NonUniqueSignatures();
    }

    for (uint256 i; i < prevSigners.length; ++i) {
      // Check the signers being removed are not zero address or duplicates
      if (config.oracles[prevSigners[i]].role == Role.Unset) {
        revert NonUniqueSignatures();
      }

      delete config.oracles[prevSigners[i]];
    }

    // Once signers have been cleared we can set the new signers
    _setConfig(configDigest, newSigners, f, new Common.AddressAndWeight[](0), true);

    emit ConfigUpdated(configDigest, prevSigners, newSigners);
  }

  /// @inheritdoc IVerifier
  function setConfig(
    bytes32 configDigest,
    address[] calldata signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external override checkConfigValid(signers.length, f) onlyOwner {
    _setConfig(configDigest, signers, f, recipientAddressesAndWeights, false);
  }

  function _setConfig(
    bytes32 configDigest,
    address[] calldata signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights,
    bool _updateConfig
  ) internal {
    VerifierState storage verifierState = s_verifierStates[configDigest];

    if (verifierState.f > 0 && !_updateConfig) {
      revert ConfigDigestAlreadySet();
    }

    verifierState.latestConfigBlockNumber = uint32(block.number);
    verifierState.f = f;
    verifierState.isActive = true;
    verifierState.oracleCount = uint8(signers.length);

    for (uint8 i; i < signers.length; ++i) {
      address signerAddr = signers[i];
      if (signerAddr == address(0)) revert ZeroAddress();

      // All signer roles are unset by default for a new config digest.
      // Here the contract checks to see if a signer's address has already
      // been set to ensure that the group of signer addresses that will
      // sign reports with the config digest are unique.
      bool isSignerAlreadySet = verifierState.oracles[signerAddr].role != Role.Unset;
      if (isSignerAlreadySet) revert NonUniqueSignatures();
      verifierState.oracles[signerAddr] = Signer({role: Role.Signer, index: i});
    }

    if (!_updateConfig) {
      IVerifierProxy(i_verifierProxyAddr).setVerifier(bytes32(0), configDigest, recipientAddressesAndWeights);

      emit ConfigSet(configDigest, signers, f);
    }
  }

  /// @inheritdoc IVerifier
  function activateConfig(bytes32 configDigest) external onlyOwner {
    VerifierState storage verifierState = s_verifierStates[configDigest];

    if (configDigest == bytes32("")) revert DigestEmpty();
    if (verifierState.f == 0) revert DigestNotSet(configDigest);
    verifierState.isActive = true;
    emit ConfigActivated(configDigest);
  }

  /// @inheritdoc IVerifier
  function deactivateConfig(bytes32 configDigest) external onlyOwner {
    VerifierState storage verifierState = s_verifierStates[configDigest];

    if (configDigest == bytes32("")) revert DigestEmpty();
    if (verifierState.f == 0) revert DigestNotSet(configDigest);
    verifierState.isActive = false;
    emit ConfigDeactivated(configDigest);
  }

  /// @inheritdoc IVerifier
  function latestConfigDetails(bytes32 configDigest) external view override returns (uint32 blockNumber) {
    VerifierState storage verifierState = s_verifierStates[configDigest];
    return (verifierState.latestConfigBlockNumber);
  }
}
