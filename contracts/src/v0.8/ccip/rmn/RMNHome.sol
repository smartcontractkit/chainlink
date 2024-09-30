// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

/// @notice Stores the home configuration for RMN, that is referenced by CCIP oracles, RMN nodes, and the RMNRemote
/// contracts.
/// @dev This contract is a state machine with the following states:
/// - Init: The initial state of the contract, no config has been set, or all configs have been revoked.
///   [0, 0]
///
/// - Candidate: A new config has been set, but it has not been promoted yet, or all active configs have been revoked.
///   [0, 1]
///
/// - Active: A non-zero config has been promoted and is active, there is no candidate configured.
///   [1, 0]
///
/// - ActiveAndCandidate: A non-zero config has been promoted and is active, and a new config has been set as candidate.
///   [1, 1]
///
/// The following state transitions are allowed:
/// - Init -> Candidate: setCandidate()
/// - Candidate -> Active: promoteCandidateAndRevokeActive()
/// - Candidate -> Candidate: setCandidate()
/// - Candidate -> Init: revokeCandidate()
/// - Active -> ActiveAndCandidate: setCandidate()
/// - Active -> Init: promoteCandidateAndRevokeActive()
/// - ActiveAndCandidate -> Active: promoteCandidateAndRevokeActive()
/// - ActiveAndCandidate -> Active: revokeCandidate()
/// - ActiveAndCandidate -> ActiveAndCandidate: setCandidate()
///
/// This means the following calls are not allowed at the following states:
/// - Init: promoteCandidateAndRevokeActive(), as there is no config to promote.
/// - Init: revokeCandidate(), as there is no config to revoke
/// - Active: revokeCandidate(), as there is no candidate to revoke
/// Note that we explicitly do allow promoteCandidateAndRevokeActive() to be called when there is an active config but
/// no candidate config. This is the only way to remove the active config. The alternative would be to set some unusable
/// config as candidate and promote that, but fully clearing it is cleaner.
///
///       ┌─────────────┐   setCandidate     ┌─────────────┐
///       │             ├───────────────────►│             │ setCandidate
///       │    Init     │   revokeCandidate  │  Candidate  │◄───────────┐
///       │    [0,0]    │◄───────────────────┤    [0,1]    │────────────┘
///       │             │  ┌─────────────────┤             │
///       └─────────────┘  │  promote-       └─────────────┘
///                  ▲     │  Candidate
///        promote-  │     │
///        Candidate │     │
///                  │     │
///       ┌──────────┴──┐  │  promote-       ┌─────────────┐
///       │             │◄─┘  Candidate OR   │  Active &   │ setCandidate
///       │    Active   │    revokeCandidate │  Candidate  │◄───────────┐
///       │    [1,0]    │◄───────────────────┤    [1,1]    │────────────┘
///       │             ├───────────────────►│             │
///       └─────────────┘    setSecondary    └─────────────┘
///
contract RMNHome is OwnerIsCreator, ITypeAndVersion {
  event ConfigSet(bytes32 indexed configDigest, uint32 version, StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event ActiveConfigRevoked(bytes32 indexed configDigest);
  event CandidateConfigRevoked(bytes32 indexed configDigest);
  event DynamicConfigSet(bytes32 indexed configDigest, DynamicConfig dynamicConfig);
  event ConfigPromoted(bytes32 indexed configDigest);

  error OutOfBoundsNodesLength();
  error DuplicatePeerId();
  error DuplicateOffchainPublicKey();
  error DuplicateSourceChain();
  error OutOfBoundsObserverNodeIndex();
  error MinObserversTooHigh();
  error ConfigDigestMismatch(bytes32 expectedConfigDigest, bytes32 gotConfigDigest);
  error DigestNotFound(bytes32 configDigest);
  error RevokingZeroDigestNotAllowed();
  error NoOpStateTransitionNotAllowed();

  struct Node {
    bytes32 peerId; //            Used for p2p communication.
    bytes32 offchainPublicKey; // Observations are signed with this public key, and are only verified offchain.
  }

  struct SourceChain {
    uint64 chainSelector; // ─────╮ The Source chain selector.
    uint64 minObservers; // ──────╯ Required number of observers to agree on an observation for this source chain.
    // ObserverNodesBitmap & (1<<i) == (1<<i) iff StaticConfig.nodes[i] is an observer for this source chain.
    uint256 observerNodesBitmap;
  }

  struct StaticConfig {
    // No sorting requirement for nodes, but ensure that SourceChain.observerNodeIndices in the home chain config &
    // Signer.nodeIndex in the remote chain configs are appropriately updated when changing this field.
    Node[] nodes;
    bytes offchainConfig; // Offchain configuration for RMN nodes.
  }

  struct DynamicConfig {
    // No sorting requirement for source chains, it is most gas efficient to append new source chains to the right.
    SourceChain[] sourceChains;
    bytes offchainConfig; // Offchain configuration for RMN nodes.
  }

  /// @notice The main struct stored in the contract, containing the static and dynamic parts of the config as well as
  /// the version and the digest of the config.
  struct VersionedConfig {
    uint32 version;
    bytes32 configDigest;
    StaticConfig staticConfig;
    DynamicConfig dynamicConfig;
  }

  string public constant override typeAndVersion = "RMNHome 1.6.0-dev";

  /// @notice Used for encoding the config digest prefix, unique per Home contract implementation.
  uint256 private constant PREFIX = 0x000b << (256 - 16); // 0x000b00..00
  /// @notice Used for encoding the config digest prefix
  uint256 private constant PREFIX_MASK = type(uint256).max << (256 - 16); // 0xFFFF00..00
  /// @notice The max number of configs that can be active at the same time.
  uint256 private constant MAX_CONCURRENT_CONFIGS = 2;
  /// @notice Helper to identify the zero config digest with less casting.
  bytes32 private constant ZERO_DIGEST = bytes32(uint256(0));
  // @notice To ensure that observerNodesBitmap can be bit-encoded into a uint256.
  uint256 private constant MAX_NODES = 256;

  /// @notice This array holds the configs.
  /// @dev Value i in this array is valid iff s_configs[i].configDigest != 0.
  VersionedConfig[MAX_CONCURRENT_CONFIGS] private s_configs;

  /// @notice The latest version set, incremented by one for each new config.
  uint32 private s_currentVersion = 0;
  /// @notice The index of the active config. Used to determine which config is active. Adding the configs to a list
  /// with two items and using this index to determine which one is active is a gas efficient way to handle this. Having
  /// a set place for the active config would mean we have to copy the candidate config to the active config when it is
  /// promoted, which would be more expensive. This index allows us to flip the configs around using `XOR 1`, which
  /// flips 0 to 1 and 1 to 0.
  uint32 private s_activeConfigIndex = 0;

  // ================================================================
  // │                          Getters                             │
  // ================================================================

  /// @notice Returns the current active and candidate config digests.
  /// @dev Can be bytes32(0) if no config has been set yet or it has been revoked.
  /// @return activeConfigDigest The digest of the active config.
  /// @return candidateConfigDigest The digest of the candidate config.
  function getConfigDigests() external view returns (bytes32 activeConfigDigest, bytes32 candidateConfigDigest) {
    return (s_configs[_getActiveIndex()].configDigest, s_configs[_getCandidateIndex()].configDigest);
  }

  /// @notice Returns the active config digest
  function getActiveDigest() external view returns (bytes32) {
    return s_configs[_getActiveIndex()].configDigest;
  }

  /// @notice Returns the candidate config digest
  function getCandidateDigest() public view returns (bytes32) {
    return s_configs[_getCandidateIndex()].configDigest;
  }

  /// @notice The offchain code can use this to fetch an old config which might still be in use by some remotes. Use
  /// in case one of the configs is too large to be returnable by one of the other getters.
  /// @param configDigest The digest of the config to fetch.
  /// @return versionedConfig The config and its version.
  /// @return ok True if the config was found, false otherwise.
  function getConfig(bytes32 configDigest) external view returns (VersionedConfig memory versionedConfig, bool ok) {
    for (uint256 i = 0; i < MAX_CONCURRENT_CONFIGS; ++i) {
      // We never want to return true for a zero digest, even if the caller is asking for it, as this can expose old
      // config state that is invalid.
      if (s_configs[i].configDigest == configDigest && configDigest != ZERO_DIGEST) {
        return (s_configs[i], true);
      }
    }
    return (versionedConfig, false);
  }

  function getAllConfigs()
    external
    view
    returns (VersionedConfig memory activeConfig, VersionedConfig memory candidateConfig)
  {
    VersionedConfig memory storedActiveConfig = s_configs[_getActiveIndex()];
    if (storedActiveConfig.configDigest != ZERO_DIGEST) {
      activeConfig = storedActiveConfig;
    }

    VersionedConfig memory storedCandidateConfig = s_configs[_getCandidateIndex()];
    if (storedCandidateConfig.configDigest != ZERO_DIGEST) {
      candidateConfig = storedCandidateConfig;
    }

    return (activeConfig, candidateConfig);
  }

  // ================================================================
  // │                     State transitions                        │
  // ================================================================

  /// @notice Sets a new config as the candidate config. Does not influence the active config.
  /// @param staticConfig The static part of the config.
  /// @param dynamicConfig The dynamic part of the config.
  /// @param digestToOverwrite The digest of the config to overwrite, or ZERO_DIGEST if no config is to be overwritten.
  /// This is done to prevent accidental overwrites.
  /// @return newConfigDigest The digest of the new config.
  function setCandidate(
    StaticConfig calldata staticConfig,
    DynamicConfig calldata dynamicConfig,
    bytes32 digestToOverwrite
  ) external onlyOwner returns (bytes32 newConfigDigest) {
    _validateStaticAndDynamicConfig(staticConfig, dynamicConfig);

    bytes32 existingDigest = getCandidateDigest();

    if (existingDigest != digestToOverwrite) {
      revert ConfigDigestMismatch(existingDigest, digestToOverwrite);
    }

    // are we going to overwrite a config? If so, emit an event.
    if (existingDigest != ZERO_DIGEST) {
      emit CandidateConfigRevoked(digestToOverwrite);
    }

    uint32 newVersion = ++s_currentVersion;
    newConfigDigest = _calculateConfigDigest(abi.encode(staticConfig), newVersion);

    VersionedConfig storage existingConfig = s_configs[_getCandidateIndex()];
    existingConfig.configDigest = newConfigDigest;
    existingConfig.version = newVersion;
    existingConfig.staticConfig = staticConfig;
    existingConfig.dynamicConfig = dynamicConfig;

    emit ConfigSet(newConfigDigest, newVersion, staticConfig, dynamicConfig);

    return newConfigDigest;
  }

  /// @notice Revokes a specific config by digest. This is used when the candidate config turns out to be incorrect to
  /// remove it without it ever having to be promoted. It's also possible to revoke the candidate config by setting a
  /// newer candidate config using `setCandidate`.
  /// @param configDigest The digest of the config to revoke. This is done to prevent accidental revokes.
  function revokeCandidate(bytes32 configDigest) external onlyOwner {
    if (configDigest == ZERO_DIGEST) {
      revert RevokingZeroDigestNotAllowed();
    }

    uint256 candidateConfigIndex = _getCandidateIndex();
    if (s_configs[candidateConfigIndex].configDigest != configDigest) {
      revert ConfigDigestMismatch(s_configs[candidateConfigIndex].configDigest, configDigest);
    }

    emit CandidateConfigRevoked(configDigest);
    // Delete only the digest, as that's what's used to determine if a config is active. This means the actual
    // config stays in storage which should significantly reduce the gas cost of overwriting that storage space in
    // the future.
    delete s_configs[candidateConfigIndex].configDigest;
  }

  /// @notice Promotes the candidate config to the active config and revokes the active config.
  /// @param digestToPromote The digest of the config to promote.
  /// @param digestToRevoke The digest of the config to revoke.
  /// @dev No config is changed in storage, the only storage changes that happen are
  /// - The activeConfigIndex is flipped.
  /// - The digest of the old active config is deleted.
  function promoteCandidateAndRevokeActive(bytes32 digestToPromote, bytes32 digestToRevoke) external onlyOwner {
    if (digestToPromote == ZERO_DIGEST && digestToRevoke == ZERO_DIGEST) {
      revert NoOpStateTransitionNotAllowed();
    }

    uint256 candidateConfigIndex = _getCandidateIndex();
    if (s_configs[candidateConfigIndex].configDigest != digestToPromote) {
      revert ConfigDigestMismatch(s_configs[candidateConfigIndex].configDigest, digestToPromote);
    }

    VersionedConfig storage activeConfig = s_configs[_getActiveIndex()];
    if (activeConfig.configDigest != digestToRevoke) {
      revert ConfigDigestMismatch(activeConfig.configDigest, digestToRevoke);
    }

    delete activeConfig.configDigest;

    s_activeConfigIndex ^= 1;
    if (digestToRevoke != ZERO_DIGEST) {
      emit ActiveConfigRevoked(digestToRevoke);
    }
    emit ConfigPromoted(digestToPromote);
  }

  /// @notice Sets the dynamic config for a specific config.
  /// @param newDynamicConfig The new dynamic config.
  /// @param currentDigest The digest of the config to update.
  /// @dev This does not update the config digest as only the static config is part of the digest.
  function setDynamicConfig(DynamicConfig calldata newDynamicConfig, bytes32 currentDigest) external onlyOwner {
    for (uint256 i = 0; i < MAX_CONCURRENT_CONFIGS; ++i) {
      if (s_configs[i].configDigest == currentDigest && currentDigest != ZERO_DIGEST) {
        _validateDynamicConfig(newDynamicConfig, s_configs[i].staticConfig.nodes.length);
        // Since the static config doesn't change we don't have to update the digest or version.
        s_configs[i].dynamicConfig = newDynamicConfig;

        emit DynamicConfigSet(currentDigest, newDynamicConfig);
        return;
      }
    }

    revert DigestNotFound(currentDigest);
  }

  /// @notice Calculates the config digest for a given plugin key, static config, and version.
  /// @param staticConfig The static part of the config.
  /// @param version The version of the config.
  /// @return The calculated config digest.
  function _calculateConfigDigest(bytes memory staticConfig, uint32 version) internal view returns (bytes32) {
    return bytes32(
      (PREFIX & PREFIX_MASK)
        | (
          uint256(
            keccak256(bytes.concat(abi.encode(bytes32("EVM"), block.chainid, address(this), version), staticConfig))
          ) & ~PREFIX_MASK
        )
    );
  }

  function _getActiveIndex() private view returns (uint32) {
    return s_activeConfigIndex;
  }

  function _getCandidateIndex() private view returns (uint32) {
    return s_activeConfigIndex ^ 1;
  }

  // ================================================================
  // │                         Validation                           │
  // ================================================================

  /// @notice Validates the static and dynamic config. Reverts when the config is invalid.
  /// @param staticConfig The static part of the config.
  /// @param dynamicConfig The dynamic part of the config.
  function _validateStaticAndDynamicConfig(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig
  ) internal pure {
    // Ensure that observerNodesBitmap can be bit-encoded into a uint256.
    if (staticConfig.nodes.length > MAX_NODES) {
      revert OutOfBoundsNodesLength();
    }

    // Ensure no peerId or offchainPublicKey is duplicated.
    for (uint256 i = 0; i < staticConfig.nodes.length; ++i) {
      for (uint256 j = i + 1; j < staticConfig.nodes.length; ++j) {
        if (staticConfig.nodes[i].peerId == staticConfig.nodes[j].peerId) {
          revert DuplicatePeerId();
        }
        if (staticConfig.nodes[i].offchainPublicKey == staticConfig.nodes[j].offchainPublicKey) {
          revert DuplicateOffchainPublicKey();
        }
      }
    }

    _validateDynamicConfig(dynamicConfig, staticConfig.nodes.length);
  }

  /// @notice Validates the dynamic config. Reverts when the config is invalid.
  /// @param dynamicConfig The dynamic part of the config.
  /// @param numberOfNodes The number of nodes in the static config.
  function _validateDynamicConfig(DynamicConfig memory dynamicConfig, uint256 numberOfNodes) internal pure {
    uint256 numberOfSourceChains = dynamicConfig.sourceChains.length;
    for (uint256 i = 0; i < numberOfSourceChains; ++i) {
      SourceChain memory currentSourceChain = dynamicConfig.sourceChains[i];
      // Ensure the source chain is unique.
      for (uint256 j = i + 1; j < numberOfSourceChains; ++j) {
        if (currentSourceChain.chainSelector == dynamicConfig.sourceChains[j].chainSelector) {
          revert DuplicateSourceChain();
        }
      }

      // all observer node indices are valid
      uint256 bitmap = currentSourceChain.observerNodesBitmap;
      // Check if there are any bits set for indexes outside of the expected range.
      if (bitmap & (type(uint256).max >> (256 - numberOfNodes)) != bitmap) {
        revert OutOfBoundsObserverNodeIndex();
      }

      uint256 observersCount = 0;
      for (; bitmap != 0; ++observersCount) {
        bitmap &= bitmap - 1;
      }

      // minObservers are tenable
      if (currentSourceChain.minObservers > observersCount) {
        revert MinObserversTooHigh();
      }
    }
  }
}
