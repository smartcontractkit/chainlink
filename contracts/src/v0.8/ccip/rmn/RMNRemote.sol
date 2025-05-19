// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRMNRemote} from "../interfaces/IRMNRemote.sol";

import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";
import {EnumerableSet} from "../../shared/enumerable/EnumerableSetWithBytes16.sol";
import {Internal} from "../libraries/Internal.sol";

/// @dev An active curse on this subject will cause isCursed() and isCursed(bytes16) to return true. Use this subject
/// for issues affecting all of CCIP chains, or pertaining to the chain that this contract is deployed on, instead of
/// using the local chain selector as a subject.
bytes16 constant GLOBAL_CURSE_SUBJECT = 0x01000000000000000000000000000001;

/// @notice This contract supports verification of RMN reports for any Any2EVM OffRamp.
/// @dev This contract implements both the new IRMNRemote interface and the legacy IRMN interface. This is to allow for
/// a seamless migration from the legacy RMN contract to this one. The only function that has been dropped in the newer
/// interface is `isBlessed`. For the `isBlessed` function, this contract relays the call to the legacy RMN contract.
contract RMNRemote is Ownable2StepMsgSender, ITypeAndVersion, IRMNRemote, IRMN {
  using EnumerableSet for EnumerableSet.Bytes16Set;

  error AlreadyCursed(bytes16 subject);
  error ConfigNotSet();
  error DuplicateOnchainPublicKey();
  error InvalidSignature();
  error InvalidSignerOrder();
  error NotEnoughSigners();
  error NotCursed(bytes16 subject);
  error OutOfOrderSignatures();
  error ThresholdNotMet();
  error UnexpectedSigner();
  error ZeroValueNotAllowed();
  error IsBlessedNotAvailable();

  event ConfigSet(uint32 indexed version, Config config);
  event Cursed(bytes16[] subjects);
  event Uncursed(bytes16[] subjects);

  /// @dev the configuration of an RMN signer.
  struct Signer {
    address onchainPublicKey; // ─╮ For signing reports.
    uint64 nodeIndex; // ─────────╯ Maps to nodes in home chain config, should be strictly increasing.
  }

  /// @dev the contract config.
  struct Config {
    bytes32 rmnHomeContractConfigDigest; // Digest of the RMNHome contract config.
    Signer[] signers; // List of signers.
    uint64 fSign; // Max number of faulty RMN nodes; f+1 signers are required to verify a report, must configure 2f+1 signers in total.
  }

  /// @dev part of the payload that RMN nodes sign: keccak256(abi.encode(RMN_V1_6_ANY2EVM_REPORT, report)).
  /// @dev this struct is only ever abi-encoded and hashed; it is never stored.
  struct Report {
    uint256 destChainId; //                 To guard against chain selector misconfiguration.
    uint64 destChainSelector; //  ────────╮ The chain selector of the destination chain.
    address rmnRemoteContractAddress; // ─╯ The address of this contract.
    address offrampAddress; //              The address of the offramp on the same chain as this contract.
    bytes32 rmnHomeContractConfigDigest; // The digest of the RMNHome contract config.
    Internal.MerkleRoot[] merkleRoots; //   The dest lane updates.
  }

  /// @dev this is included in the preimage of the digest that RMN nodes sign.
  bytes32 private constant RMN_V1_6_ANY2EVM_REPORT = keccak256("RMN_V1_6_ANY2EVM_REPORT");

  string public constant override typeAndVersion = "RMNRemote 1.6.0";
  uint64 internal immutable i_localChainSelector;
  IRMN internal immutable i_legacyRMN;

  Config private s_config;
  uint32 private s_configCount;

  /// @dev RMN nodes only generate sigs with v=27; making this constant allows us to save gas by not transmitting v.
  /// @dev Any valid ECDSA sig (r, s, v) can be "flipped" into (r, s*, v*) without knowing the private key (where v=27 or 28 for secp256k1)
  /// https://github.com/kadenzipfel/smart-contract-vulnerabilities/blob/master/vulnerabilities/signature-malleability.md.
  uint8 private constant ECDSA_RECOVERY_V = 27;

  EnumerableSet.Bytes16Set private s_cursedSubjects;
  mapping(address signer => bool exists) private s_signers; // for more gas efficient verify.

  /// @param localChainSelector the chain selector of the chain this contract is deployed to.
  constructor(uint64 localChainSelector, IRMN legacyRMN) {
    if (localChainSelector == 0) revert ZeroValueNotAllowed();
    i_localChainSelector = localChainSelector;

    i_legacyRMN = legacyRMN;
  }

  // ================================================================
  // │                         Verification                         │
  // ================================================================

  /// @inheritdoc IRMNRemote
  function verify(
    address offRampAddress,
    Internal.MerkleRoot[] calldata merkleRoots,
    Signature[] calldata signatures
  ) external view {
    if (s_configCount == 0) {
      revert ConfigNotSet();
    }
    if (signatures.length < s_config.fSign + 1) revert ThresholdNotMet();

    bytes32 digest = keccak256(
      abi.encode(
        RMN_V1_6_ANY2EVM_REPORT,
        Report({
          destChainId: block.chainid,
          destChainSelector: i_localChainSelector,
          rmnRemoteContractAddress: address(this),
          offrampAddress: offRampAddress,
          rmnHomeContractConfigDigest: s_config.rmnHomeContractConfigDigest,
          merkleRoots: merkleRoots
        })
      )
    );

    address prevAddress;
    address signerAddress;
    for (uint256 i = 0; i < signatures.length; ++i) {
      signerAddress = ecrecover(digest, ECDSA_RECOVERY_V, signatures[i].r, signatures[i].s);
      if (signerAddress == address(0)) revert InvalidSignature();
      if (prevAddress >= signerAddress) revert OutOfOrderSignatures();
      if (!s_signers[signerAddress]) revert UnexpectedSigner();
      prevAddress = signerAddress;
    }
  }

  // ================================================================
  // │                            Config                            │
  // ================================================================

  /// @notice Sets the configuration of the contract.
  /// @param newConfig the new configuration.
  /// @dev setting config is atomic; we delete all pre-existing config and set everything from scratch.
  function setConfig(
    Config calldata newConfig
  ) external onlyOwner {
    if (newConfig.rmnHomeContractConfigDigest == bytes32(0)) {
      revert ZeroValueNotAllowed();
    }

    // signers are in ascending order of nodeIndex.
    for (uint256 i = 1; i < newConfig.signers.length; ++i) {
      if (!(newConfig.signers[i - 1].nodeIndex < newConfig.signers[i].nodeIndex)) {
        revert InvalidSignerOrder();
      }
    }

    // min signers requirement is tenable.
    if (newConfig.signers.length < 2 * newConfig.fSign + 1) {
      revert NotEnoughSigners();
    }

    // clear the old signers.
    for (uint256 i = s_config.signers.length; i > 0; --i) {
      delete s_signers[s_config.signers[i - 1].onchainPublicKey];
    }

    // set the new signers.
    for (uint256 i = 0; i < newConfig.signers.length; ++i) {
      if (s_signers[newConfig.signers[i].onchainPublicKey]) {
        revert DuplicateOnchainPublicKey();
      }
      s_signers[newConfig.signers[i].onchainPublicKey] = true;
    }

    s_config = newConfig;
    uint32 newConfigCount = ++s_configCount;
    emit ConfigSet(newConfigCount, newConfig);
  }

  /// @notice Returns the current configuration of the contract and a version number.
  /// @return version the current configs version.
  /// @return config the current config.
  function getVersionedConfig() external view returns (uint32 version, Config memory config) {
    return (s_configCount, s_config);
  }

  /// @notice Returns the chain selector configured at deployment time.
  /// @return localChainSelector the chain selector, not the chain ID.
  function getLocalChainSelector() external view returns (uint64 localChainSelector) {
    return i_localChainSelector;
  }

  /// @notice Returns the 32 byte header used in computing the report digest.
  /// @return digestHeader the digest header.
  function getReportDigestHeader() external pure returns (bytes32 digestHeader) {
    return RMN_V1_6_ANY2EVM_REPORT;
  }

  // ================================================================
  // │                           Cursing                            │
  // ================================================================

  /// @notice Curse a single subject.
  /// @param subject the subject to curse.
  function curse(
    bytes16 subject
  ) external {
    bytes16[] memory subjects = new bytes16[](1);
    subjects[0] = subject;
    curse(subjects);
  }

  /// @notice Curse an array of subjects.
  /// @param subjects the subjects to curse.
  /// @dev reverts if any of the subjects are already cursed or if there is a duplicate.
  function curse(
    bytes16[] memory subjects
  ) public onlyOwner {
    for (uint256 i = 0; i < subjects.length; ++i) {
      if (!s_cursedSubjects.add(subjects[i])) {
        revert AlreadyCursed(subjects[i]);
      }
    }
    emit Cursed(subjects);
  }

  /// @notice Uncurse a single subject.
  /// @param subject the subject to uncurse.
  function uncurse(
    bytes16 subject
  ) external {
    bytes16[] memory subjects = new bytes16[](1);
    subjects[0] = subject;
    uncurse(subjects);
  }

  /// @notice Uncurse an array of subjects.
  /// @param subjects the subjects to uncurse.
  /// @dev reverts if any of the subjects are not cursed or if there is a duplicate.
  function uncurse(
    bytes16[] memory subjects
  ) public onlyOwner {
    for (uint256 i = 0; i < subjects.length; ++i) {
      if (!s_cursedSubjects.remove(subjects[i])) {
        revert NotCursed(subjects[i]);
      }
    }
    emit Uncursed(subjects);
  }

  /// @inheritdoc IRMNRemote
  function getCursedSubjects() external view returns (bytes16[] memory subjects) {
    return s_cursedSubjects.values();
  }

  /// @inheritdoc IRMNRemote
  function isCursed() external view override(IRMN, IRMNRemote) returns (bool) {
    // There are zero curses under normal circumstances, which means it's cheaper to check for the absence of curses.
    // than to check the subject list for the global curse subject.
    if (s_cursedSubjects.length() == 0) {
      return false;
    }
    return s_cursedSubjects.contains(GLOBAL_CURSE_SUBJECT);
  }

  /// @inheritdoc IRMNRemote
  function isCursed(
    bytes16 subject
  ) external view override(IRMN, IRMNRemote) returns (bool) {
    // There are zero curses under normal circumstances, which means it's cheaper to check for the absence of curses.
    // than to check the subject list twice, as we have to check for both the given and global curse subjects.
    if (s_cursedSubjects.length() == 0) {
      return false;
    }
    return s_cursedSubjects.contains(subject) || s_cursedSubjects.contains(GLOBAL_CURSE_SUBJECT);
  }

  // ================================================================
  // │                     Legacy pass through                      │
  // ================================================================

  /// @inheritdoc IRMN
  /// @dev This function is only expected to be used for messages from CCIP versions below 1.6.
  function isBlessed(
    TaggedRoot calldata taggedRoot
  ) external view returns (bool) {
    if (i_legacyRMN == IRMN(address(0))) {
      revert IsBlessedNotAvailable();
    }

    return i_legacyRMN.isBlessed(taggedRoot);
  }
}
