// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../keystone/interfaces/ICapabilityConfiguration.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {ICapabilitiesRegistry} from "../interfaces/ICapabilitiesRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Internal} from "../libraries/Internal.sol";

import {IERC165} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice CCIPHome stores the configuration for the CCIP capability.
/// We have two classes of configuration: chain configuration and DON (in the CapabilitiesRegistry sense) configuration.
/// Each chain will have a single configuration which includes information like the router address.
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
contract CCIPHome is OwnerIsCreator, ITypeAndVersion, ICapabilityConfiguration, IERC165 {
  using EnumerableSet for EnumerableSet.UintSet;

  event ChainConfigRemoved(uint64 chainSelector);
  event ChainConfigSet(uint64 chainSelector, ChainConfig chainConfig);
  event ConfigSet(bytes32 indexed configDigest, uint32 version, OCR3Config config);
  event ActiveConfigRevoked(bytes32 indexed configDigest);
  event CandidateConfigRevoked(bytes32 indexed configDigest);
  event ConfigPromoted(bytes32 indexed configDigest);

  error NodeNotInRegistry(bytes32 p2pId);
  error ChainSelectorNotFound(uint64 chainSelector);
  error FChainMustBePositive();
  error ChainSelectorNotSet();
  error InvalidPluginType();
  error OfframpAddressCannotBeZero();
  error FChainTooHigh(uint256 fChain, uint256 FRoleDON);
  error TooManySigners();
  error FTooHigh();
  error RMNHomeAddressCannotBeZero();
  error InvalidNode(OCR3Node node);
  error NotEnoughTransmitters(uint256 got, uint256 minimum);
  error OnlyCapabilitiesRegistryCanCall();
  error ZeroAddressNotAllowed();
  error ConfigDigestMismatch(bytes32 expectedConfigDigest, bytes32 gotConfigDigest);
  error CanOnlySelfCall();
  error RevokingZeroDigestNotAllowed();
  error NoOpStateTransitionNotAllowed();
  error InvalidSelector(bytes4 selector);
  error DONIdMismatch(uint32 callDonId, uint32 capabilityRegistryDonId);

  error InvalidStateTransition(
    bytes32 currentActiveDigest,
    bytes32 currentCandidateDigest,
    bytes32 proposedActiveDigest,
    bytes32 proposedCandidateDigest
  );

  /// @notice Represents an oracle node in OCR3 configs part of the role DON.
  /// Every configured node should be a signer, but does not have to be a transmitter.
  struct OCR3Node {
    bytes32 p2pId; // Peer2Peer connection ID of the oracle
    bytes signerKey; // On-chain signer public key
    bytes transmitterKey; // On-chain transmitter public key. Can be set to empty bytes to represent that the node is a signer but not a transmitter.
  }

  /// @notice OCR3 configuration.
  /// Note that FRoleDON >= fChain, since FRoleDON represents the role DON, and fChain represents sub-committees.
  /// FRoleDON values are typically identical across multiple OCR3 configs since the chains pertain to one role DON,
  /// but FRoleDON values can change across OCR3 configs to indicate role DON splits.
  struct OCR3Config {
    Internal.OCRPluginType pluginType; // ─╮ The plugin that the configuration is for.
    uint64 chainSelector; //               │ The (remote) chain that the configuration is for.
    uint8 FRoleDON; //                     │ The "big F" parameter for the role DON.
    uint64 offchainConfigVersion; // ──────╯ The version of the exec offchain configuration.
    bytes offrampAddress; // The remote chain offramp address.
    bytes rmnHomeAddress; // The home chain RMN home address.
    OCR3Node[] nodes; // Keys & IDs of nodes part of the role DON
    bytes offchainConfig; // The offchain configuration for the OCR3 plugin. Protobuf encoded.
  }

  struct VersionedConfig {
    uint32 version;
    bytes32 configDigest;
    OCR3Config config;
  }

  /// @notice Chain configuration.
  /// Changes to chain configuration are detected out-of-band in plugins and decoded offchain.
  struct ChainConfig {
    bytes32[] readers; // The P2P IDs of the readers for the chain. These IDs must be registered in the capabilities registry.
    uint8 fChain; // The fault tolerance parameter of the chain.
    bytes config; // The chain configuration. This is kept intentionally opaque so as to add fields in the future if needed.
  }

  /// @notice Chain configuration information struct used in applyChainConfigUpdates and getAllChainConfigs.
  struct ChainConfigArgs {
    uint64 chainSelector;
    ChainConfig chainConfig;
  }

  string public constant override typeAndVersion = "CCIPHome 1.6.0-dev";

  /// @dev A prefix added to all config digests that is unique to the implementation
  uint256 private constant PREFIX = 0x000a << (256 - 16); // 0x000a00..00
  bytes32 internal constant EMPTY_ENCODED_ADDRESS_HASH = keccak256(abi.encode(address(0)));
  /// @dev 256 is the hard limit due to the bit encoding of their indexes into a uint256.
  uint256 internal constant MAX_NUM_ORACLES = 256;

  /// @notice Used for encoding the config digest prefix
  uint256 private constant PREFIX_MASK = type(uint256).max << (256 - 16); // 0xFFFF00..00
  /// @notice The max number of configs that can be active at the same time.
  uint256 private constant MAX_CONCURRENT_CONFIGS = 2;
  /// @notice Helper to identify the zero config digest with less casting.
  bytes32 private constant ZERO_DIGEST = bytes32(uint256(0));

  /// @dev The canonical capabilities registry address.
  address internal immutable i_capabilitiesRegistry;

  /// @dev chain configuration for each chain that CCIP is deployed on.
  mapping(uint64 chainSelector => ChainConfig chainConfig) private s_chainConfigurations;

  /// @dev All chains that are configured.
  EnumerableSet.UintSet private s_remoteChainSelectors;

  /// @notice This array holds the configs.
  /// @dev A DonID covers a single chain, and the plugin type is used to differentiate between the commit and execution
  mapping(uint32 donId => mapping(Internal.OCRPluginType pluginType => VersionedConfig[MAX_CONCURRENT_CONFIGS])) private
    s_configs;

  /// @notice The total number of configs ever set, used for generating the version of the configs.
  /// @dev Used to ensure unique digests across all configurations.
  uint32 private s_currentVersion = 0;
  /// @notice The index of the active config on a per-don and per-plugin basis.
  mapping(uint32 donId => mapping(Internal.OCRPluginType pluginType => uint32)) private s_activeConfigIndexes;

  /// @notice Constructor for the CCIPHome contract takes in the address of the capabilities registry. This address
  /// is the only allowed caller to mutate the configuration through beforeCapabilityConfigSet.
  constructor(address capabilitiesRegistry) {
    if (capabilitiesRegistry == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    i_capabilitiesRegistry = capabilitiesRegistry;
  }

  // ================================================================
  // │                    Capability Registry                       │
  // ================================================================

  /// @notice Returns the capabilities registry address.
  /// @return The capabilities registry address.
  function getCapabilityRegistry() external view returns (address) {
    return i_capabilitiesRegistry;
  }

  /// @inheritdoc IERC165
  /// @dev Required for the capabilities registry to recognize this contract.
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == type(ICapabilityConfiguration).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  /// @notice Called by the registry prior to the config being set for a particular DON.
  /// @dev precondition Requires destination chain config to be set
  function beforeCapabilityConfigSet(
    bytes32[] calldata, // nodes
    bytes calldata update,
    // Config count is unused because we don't want to invalidate a config on blue/green promotions so we keep track of
    // the actual newly submitted configs instead of the number of config mutations.
    uint64, // config count
    uint32 donId
  ) external override {
    if (msg.sender != i_capabilitiesRegistry) {
      revert OnlyCapabilitiesRegistryCanCall();
    }

    bytes4 selector = bytes4(update[:4]);
    // We only allow self-calls to the following approved methods
    if (
      selector != this.setCandidate.selector && selector != this.revokeCandidate.selector
        && selector != this.promoteCandidateAndRevokeActive.selector
    ) {
      revert InvalidSelector(selector);
    }

    // We validate that the call contains the correct DON ID. The DON ID is always the first function argument.
    uint256 callDonId = abi.decode(update[4:36], (uint256));
    if (callDonId != donId) {
      revert DONIdMismatch(uint32(callDonId), donId);
    }

    // solhint-disable-next-line avoid-low-level-calls
    (bool success, bytes memory retData) = address(this).call(update);
    // if not successful, revert with the original revert
    if (!success) {
      assembly {
        revert(add(retData, 0x20), returndatasize())
      }
    }
  }

  /// @inheritdoc ICapabilityConfiguration
  /// @dev The CCIP capability will fetch the configuration needed directly from this contract.
  /// The offchain syncer will call this function, so its important that it doesn't revert.
  function getCapabilityConfiguration(uint32 /* donId */ ) external pure override returns (bytes memory configuration) {
    return bytes("");
  }

  // ================================================================
  // │                          Getters                             │
  // ================================================================

  /// @notice Returns the current active and candidate config digests.
  /// @dev Can be bytes32(0) if no config has been set yet or it has been revoked.
  /// @param donId The key of the plugin to get the config digests for.
  /// @return activeConfigDigest The digest of the active config.
  /// @return candidateConfigDigest The digest of the candidate config.
  function getConfigDigests(
    uint32 donId,
    Internal.OCRPluginType pluginType
  ) public view returns (bytes32 activeConfigDigest, bytes32 candidateConfigDigest) {
    return (
      s_configs[donId][pluginType][_getActiveIndex(donId, pluginType)].configDigest,
      s_configs[donId][pluginType][_getCandidateIndex(donId, pluginType)].configDigest
    );
  }

  /// @notice Returns the active config digest for for a given key.
  /// @param donId The key of the plugin to get the config digests for.
  function getActiveDigest(uint32 donId, Internal.OCRPluginType pluginType) public view returns (bytes32) {
    return s_configs[donId][pluginType][_getActiveIndex(donId, pluginType)].configDigest;
  }

  /// @notice Returns the candidate config digest for for a given key.
  /// @param donId The key of the plugin to get the config digests for.
  function getCandidateDigest(uint32 donId, Internal.OCRPluginType pluginType) public view returns (bytes32) {
    return s_configs[donId][pluginType][_getCandidateIndex(donId, pluginType)].configDigest;
  }

  /// @notice The offchain code can use this to fetch an old config which might still be in use by some remotes. Use
  /// in case one of the configs is too large to be returnable by one of the other getters.
  /// @param donId The unique key for the DON that the configuration applies to.
  /// @param configDigest The digest of the config to fetch.
  /// @return versionedConfig The config and its version.
  /// @return ok True if the config was found, false otherwise.
  function getConfig(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    bytes32 configDigest
  ) external view returns (VersionedConfig memory versionedConfig, bool ok) {
    for (uint256 i = 0; i < MAX_CONCURRENT_CONFIGS; ++i) {
      // We never want to return true for a zero digest, even if the caller is asking for it, as this can expose old
      // config state that is invalid.
      if (s_configs[donId][pluginType][i].configDigest == configDigest && configDigest != ZERO_DIGEST) {
        return (s_configs[donId][pluginType][i], true);
      }
    }
    // versionConfig is uninitialized so it contains default values.
    return (versionedConfig, false);
  }

  /// @notice Returns the active and candidate configuration for a given plugin key.
  /// @param donId The unique key for the DON that the configuration applies to.
  /// @return activeConfig The active configuration.
  /// @return candidateConfig The candidate configuration.
  function getAllConfigs(
    uint32 donId,
    Internal.OCRPluginType pluginType
  ) external view returns (VersionedConfig memory activeConfig, VersionedConfig memory candidateConfig) {
    VersionedConfig memory storedActiveConfig = s_configs[donId][pluginType][_getActiveIndex(donId, pluginType)];
    if (storedActiveConfig.configDigest != ZERO_DIGEST) {
      activeConfig = storedActiveConfig;
    }

    VersionedConfig memory storedCandidateConfig = s_configs[donId][pluginType][_getCandidateIndex(donId, pluginType)];
    if (storedCandidateConfig.configDigest != ZERO_DIGEST) {
      candidateConfig = storedCandidateConfig;
    }

    return (activeConfig, candidateConfig);
  }

  // ================================================================
  // │                     State transitions                        │
  // ================================================================

  /// @notice Sets a new config as the candidate config. Does not influence the active config.
  /// @param donId The key of the plugin to set the config for.
  /// @return newConfigDigest The digest of the new config.
  function setCandidate(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    OCR3Config calldata config,
    bytes32 digestToOverwrite
  ) external returns (bytes32 newConfigDigest) {
    _onlySelfCall();
    _validateConfig(config);

    bytes32 existingDigest = getCandidateDigest(donId, pluginType);

    if (existingDigest != digestToOverwrite) {
      revert ConfigDigestMismatch(existingDigest, digestToOverwrite);
    }

    // are we going to overwrite a config? If so, emit an event.
    if (existingDigest != ZERO_DIGEST) {
      emit CandidateConfigRevoked(digestToOverwrite);
    }

    uint32 newVersion = ++s_currentVersion;
    newConfigDigest = _calculateConfigDigest(donId, pluginType, abi.encode(config), newVersion);

    VersionedConfig storage existingConfig = s_configs[donId][pluginType][_getCandidateIndex(donId, pluginType)];
    existingConfig.configDigest = newConfigDigest;
    existingConfig.version = newVersion;
    existingConfig.config = config;

    emit ConfigSet(newConfigDigest, newVersion, config);

    return newConfigDigest;
  }

  /// @notice Revokes a specific config by digest.
  /// @param donId The key of the plugin to revoke the config for.
  /// @param configDigest The digest of the config to revoke. This is done to prevent accidental revokes.
  function revokeCandidate(uint32 donId, Internal.OCRPluginType pluginType, bytes32 configDigest) external {
    _onlySelfCall();

    if (configDigest == ZERO_DIGEST) {
      revert RevokingZeroDigestNotAllowed();
    }

    uint256 candidateConfigIndex = _getCandidateIndex(donId, pluginType);
    if (s_configs[donId][pluginType][candidateConfigIndex].configDigest != configDigest) {
      revert ConfigDigestMismatch(s_configs[donId][pluginType][candidateConfigIndex].configDigest, configDigest);
    }

    emit CandidateConfigRevoked(configDigest);
    // Delete only the digest, as that's what's used to determine if a config is active. This means the actual
    // config stays in storage which should significantly reduce the gas cost of overwriting that storage space in
    // the future.
    delete s_configs[donId][pluginType][candidateConfigIndex].configDigest;
  }

  /// @notice Promotes the candidate config to the active config and revokes the active config.
  /// @param donId The key of the plugin to promote the config for.
  /// @param digestToPromote The digest of the config to promote.
  function promoteCandidateAndRevokeActive(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    bytes32 digestToPromote,
    bytes32 digestToRevoke
  ) external {
    _onlySelfCall();

    if (digestToPromote == ZERO_DIGEST && digestToRevoke == ZERO_DIGEST) {
      revert NoOpStateTransitionNotAllowed();
    }

    uint256 candidateConfigIndex = _getCandidateIndex(donId, pluginType);
    if (s_configs[donId][pluginType][candidateConfigIndex].configDigest != digestToPromote) {
      revert ConfigDigestMismatch(s_configs[donId][pluginType][candidateConfigIndex].configDigest, digestToPromote);
    }

    VersionedConfig storage activeConfig = s_configs[donId][pluginType][_getActiveIndex(donId, pluginType)];
    if (activeConfig.configDigest != digestToRevoke) {
      revert ConfigDigestMismatch(activeConfig.configDigest, digestToRevoke);
    }

    delete activeConfig.configDigest;

    s_activeConfigIndexes[donId][pluginType] ^= 1;
    if (digestToRevoke != ZERO_DIGEST) {
      emit ActiveConfigRevoked(digestToRevoke);
    }
    emit ConfigPromoted(digestToPromote);
  }

  /// @notice Calculates the config digest for a given plugin key, static config, and version.
  /// @param donId The key of the plugin to calculate the digest for.
  /// @param staticConfig The static part of the config.
  /// @param version The version of the config.
  /// @return The calculated config digest.
  function _calculateConfigDigest(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    bytes memory staticConfig,
    uint32 version
  ) internal view returns (bytes32) {
    return bytes32(
      (PREFIX & PREFIX_MASK)
        | (
          uint256(
            keccak256(
              bytes.concat(
                abi.encode(bytes32("EVM"), block.chainid, address(this), donId, pluginType, version), staticConfig
              )
            )
          ) & ~PREFIX_MASK
        )
    );
  }

  function _getActiveIndex(uint32 donId, Internal.OCRPluginType pluginType) private view returns (uint32) {
    return s_activeConfigIndexes[donId][pluginType];
  }

  function _getCandidateIndex(uint32 donId, Internal.OCRPluginType pluginType) private view returns (uint32) {
    return s_activeConfigIndexes[donId][pluginType] ^ 1;
  }

  // ================================================================
  // │                         Validation                           │
  // ================================================================

  function _validateConfig(OCR3Config memory cfg) internal view {
    if (cfg.chainSelector == 0) revert ChainSelectorNotSet();
    if (cfg.pluginType != Internal.OCRPluginType.Commit && cfg.pluginType != Internal.OCRPluginType.Execution) {
      revert InvalidPluginType();
    }
    if (cfg.offrampAddress.length == 0 || keccak256(cfg.offrampAddress) == EMPTY_ENCODED_ADDRESS_HASH) {
      revert OfframpAddressCannotBeZero();
    }
    if (cfg.rmnHomeAddress.length == 0 || keccak256(cfg.rmnHomeAddress) == EMPTY_ENCODED_ADDRESS_HASH) {
      revert RMNHomeAddressCannotBeZero();
    }
    if (!s_remoteChainSelectors.contains(cfg.chainSelector)) revert ChainSelectorNotFound(cfg.chainSelector);

    // fChain cannot exceed FRoleDON, since it is a subcommittee in the larger DON
    uint256 FRoleDON = cfg.FRoleDON;
    uint256 fChain = s_chainConfigurations[cfg.chainSelector].fChain;
    // fChain > 0 is enforced in applyChainConfigUpdates, and the presence of a chain config is checked above
    // FRoleDON != 0 because FRoleDON >= fChain is enforced here
    if (fChain > FRoleDON) {
      revert FChainTooHigh(fChain, FRoleDON);
    }

    // len(nodes) >= 3 * FRoleDON + 1
    // len(nodes) == numberOfSigners
    uint256 numberOfNodes = cfg.nodes.length;
    if (numberOfNodes > MAX_NUM_ORACLES) revert TooManySigners();
    if (numberOfNodes <= 3 * FRoleDON) revert FTooHigh();

    uint256 nonZeroTransmitters = 0;
    bytes32[] memory p2pIds = new bytes32[](numberOfNodes);
    for (uint256 i = 0; i < numberOfNodes; ++i) {
      OCR3Node memory node = cfg.nodes[i];

      // 3 * fChain + 1 <= nonZeroTransmitters <= 3 * FRoleDON + 1
      // Transmitters can be set to 0 since there can be more signers than transmitters,
      if (node.transmitterKey.length != 0) {
        nonZeroTransmitters++;
      }

      // Signer key and p2pIds must always be present
      if (node.signerKey.length == 0 || node.p2pId == bytes32(0)) {
        revert InvalidNode(node);
      }

      p2pIds[i] = node.p2pId;
    }

    // We check for chain config presence above, so fChain here must be non-zero. fChain <= FRoleDON due to the checks above.
    // There can be less transmitters than signers - so they can be set to zero (which indicates that a node is a signer, but not a transmitter).
    uint256 minTransmittersLength = 3 * fChain + 1;
    if (nonZeroTransmitters < minTransmittersLength) {
      revert NotEnoughTransmitters(nonZeroTransmitters, minTransmittersLength);
    }

    // Check that the readers are in the capabilities registry.
    _ensureInRegistry(p2pIds);
  }

  function _onlySelfCall() internal view {
    if (msg.sender != address(this)) {
      revert CanOnlySelfCall();
    }
  }

  // ================================================================
  // │                    Chain Configuration                       │
  // ================================================================

  /// @notice Returns the total number of chains configured.
  /// @return The total number of chains configured.
  function getNumChainConfigurations() external view returns (uint256) {
    return s_remoteChainSelectors.length();
  }

  /// @notice Returns all the chain configurations.
  /// @param pageIndex The page index.
  /// @param pageSize The page size.
  /// @return paginatedChainConfigs chain configurations.
  function getAllChainConfigs(uint256 pageIndex, uint256 pageSize) external view returns (ChainConfigArgs[] memory) {
    uint256 numberOfChains = s_remoteChainSelectors.length();
    uint256 startIndex = pageIndex * pageSize;

    if (pageSize == 0 || startIndex >= numberOfChains) {
      return new ChainConfigArgs[](0); // Return an empty array if pageSize is 0 or pageIndex is out of bounds
    }

    uint256 endIndex = startIndex + pageSize;
    if (endIndex > numberOfChains) {
      endIndex = numberOfChains;
    }

    ChainConfigArgs[] memory paginatedChainConfigs = new ChainConfigArgs[](endIndex - startIndex);

    uint256[] memory chainSelectors = s_remoteChainSelectors.values();
    for (uint256 i = startIndex; i < endIndex; ++i) {
      uint64 chainSelector = uint64(chainSelectors[i]);
      paginatedChainConfigs[i - startIndex] =
        ChainConfigArgs({chainSelector: chainSelector, chainConfig: s_chainConfigurations[chainSelector]});
    }

    return paginatedChainConfigs;
  }

  /// @notice Sets and/or removes chain configurations.
  /// Does not validate that fChain <= FRoleDON and relies on OCR3Configs to be changed in case fChain becomes larger than the FRoleDON value.
  /// @param chainSelectorRemoves The chain configurations to remove.
  /// @param chainConfigAdds The chain configurations to add.
  function applyChainConfigUpdates(
    uint64[] calldata chainSelectorRemoves,
    ChainConfigArgs[] calldata chainConfigAdds
  ) external onlyOwner {
    // Process removals first.
    for (uint256 i = 0; i < chainSelectorRemoves.length; ++i) {
      // check if the chain selector is in s_remoteChainSelectors first.
      if (!s_remoteChainSelectors.contains(chainSelectorRemoves[i])) {
        revert ChainSelectorNotFound(chainSelectorRemoves[i]);
      }

      delete s_chainConfigurations[chainSelectorRemoves[i]];
      s_remoteChainSelectors.remove(chainSelectorRemoves[i]);

      emit ChainConfigRemoved(chainSelectorRemoves[i]);
    }

    // Process additions next.
    for (uint256 i = 0; i < chainConfigAdds.length; ++i) {
      ChainConfig memory chainConfig = chainConfigAdds[i].chainConfig;
      uint64 chainSelector = chainConfigAdds[i].chainSelector;

      // Verify that the provided readers are present in the capabilities registry.
      _ensureInRegistry(chainConfig.readers);

      // Verify that fChain is positive.
      if (chainConfig.fChain == 0) {
        revert FChainMustBePositive();
      }

      s_chainConfigurations[chainSelector] = chainConfig;
      s_remoteChainSelectors.add(chainSelector);

      emit ChainConfigSet(chainSelector, chainConfig);
    }
  }

  /// @notice Helper function to ensure that a node is in the capabilities registry.
  /// @param p2pIds The P2P IDs of the node to check.
  function _ensureInRegistry(bytes32[] memory p2pIds) internal view {
    for (uint256 i = 0; i < p2pIds.length; ++i) {
      // TODO add a method that does the validation in the ICapabilitiesRegistry contract
      if (ICapabilitiesRegistry(i_capabilitiesRegistry).getNode(p2pIds[i]).p2pId == bytes32("")) {
        revert NodeNotInRegistry(p2pIds[i]);
      }
    }
  }
}
