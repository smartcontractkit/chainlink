// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../keystone/interfaces/ICapabilityConfiguration.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {ICapabilityRegistry} from "./interfaces/ICapabilityRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice CCIPCapabilityConfiguration stores the configuration for the CCIP capability.
/// We have two classes of configuration: chain configuration and DON (in the CapabilityRegistry sense) configuration.
/// Each chain will have a single configuration which includes information like the router address.
/// Each CR DON will have up to four configurations: for each of (commit, exec), one blue and one green configuration.
/// This is done in order to achieve "blue-green" deployments.
contract CCIPCapabilityConfiguration is ITypeAndVersion, ICapabilityConfiguration, OwnerIsCreator {
  using EnumerableSet for EnumerableSet.UintSet;

  /// @notice Emitted when a chain's configuration is set.
  /// @param chainSelector The chain selector.
  /// @param chainConfig The chain configuration.
  event ChainConfigSet(uint64 chainSelector, ChainConfig chainConfig);

  /// @notice Emitted when a chain's configuration is removed.
  /// @param chainSelector The chain selector.
  event ChainConfigRemoved(uint64 chainSelector);

  error ChainConfigNotSetForChain(uint64 chainSelector);
  error NodeNotInRegistry(bytes32 p2pId);
  error OnlyCapabilityRegistryCanCall();
  error ChainSelectorNotFound(uint64 chainSelector);
  error ChainSelectorNotSet();
  error TooManyOCR3Configs();
  error TooManySigners();
  error TooManyTransmitters();
  error TooManyBootstrapP2PIds();
  error P2PIdsLengthNotMatching(uint256 p2pIdsLength, uint256 signersLength, uint256 transmittersLength);
  error NotEnoughTransmitters(uint256 got, uint256 minimum);
  error FMustBePositive();
  error FChainMustBePositive();
  error FTooHigh();
  error InvalidPluginType();
  error OfframpAddressCannotBeZero();
  error InvalidConfigLength(uint256 length);
  error InvalidConfigStateTransition(ConfigState currentState, ConfigState proposedState);
  error NonExistentConfigTransition();
  error WrongConfigCount(uint64 got, uint64 expected);
  error WrongConfigDigest(bytes32 got, bytes32 expected);
  error WrongConfigDigestBlueGreen(bytes32 got, bytes32 expected);

  /// @notice PluginType indicates the type of plugin that the configuration is for.
  /// @param Commit The configuration is for the commit plugin.
  /// @param Execution The configuration is for the execution plugin.
  enum PluginType {
    Commit,
    Execution
  }

  /// @notice ConfigState indicates the state of the configuration.
  /// A DON's configuration always starts out in the "Init" state - this is the starting state.
  /// The only valid transition from "Init" is to the "Running" state - this is the first ever configuration.
  /// The only valid transition from "Running" is to the "Staging" state - this is a blue/green proposal.
  /// The only valid transition from "Staging" is back to the "Running" state - this is a promotion.
  /// TODO: explain rollbacks?
  enum ConfigState {
    Init,
    Running,
    Staging
  }

  /// @notice Chain configuration.
  /// Changes to chain configuration are detected out-of-band in plugins and decoded offchain.
  struct ChainConfig {
    bytes32[] readers; // The P2P IDs of the readers for the chain. These IDs must be registered in the capability registry.
    uint8 fChain; // The fault tolerance parameter of the chain.
    bytes config; // The chain configuration. This is kept intentionally opaque so as to add fields in the future if needed.
  }

  /// @notice Chain configuration information struct used in applyChainConfigUpdates and getAllChainConfigs.
  struct ChainConfigInfo {
    uint64 chainSelector;
    ChainConfig chainConfig;
  }

  /// @notice OCR3 configuration.
  struct OCR3Config {
    PluginType pluginType; // ────────╮ The plugin that the configuration is for.
    uint64 chainSelector; //          | The (remote) chain that the configuration is for.
    uint8 F; //                       | The "big F" parameter for the role DON.
    uint64 offchainConfigVersion; // ─╯ The version of the offchain configuration.
    bytes offrampAddress; // The remote chain offramp address.
    bytes32[] bootstrapP2PIds; // The bootstrap P2P IDs of the oracles that are part of the role DON.
    // len(p2pIds) == len(signers) == len(transmitters) == 3 * F + 1
    // NOTE: indexes matter here! The p2p ID at index i corresponds to the signer at index i and the transmitter at index i.
    // This is crucial in order to build the oracle ID <-> peer ID mapping offchain.
    bytes32[] p2pIds; // The P2P IDs of the oracles that are part of the role DON.
    bytes[] signers; // The onchain signing keys of nodes in the don.
    bytes[] transmitters; // The onchain transmitter keys of nodes in the don.
    bytes offchainConfig; // The offchain configuration for the OCR3 protocol. Protobuf encoded.
  }

  /// @notice OCR3 configuration with metadata, specifically the config count and the config digest.
  struct OCR3ConfigWithMeta {
    OCR3Config config; // The OCR3 configuration.
    uint64 configCount; // The config count used to compute the config digest.
    bytes32 configDigest; // The config digest of the OCR3 configuration.
  }

  /// @notice Type and version override.
  string public constant override typeAndVersion = "CCIPCapabilityConfiguration 1.6.0-dev";

  /// @notice The canonical capability registry address.
  address internal immutable i_capabilityRegistry;

  /// @notice chain configuration for each chain that CCIP is deployed on.
  mapping(uint64 chainSelector => ChainConfig chainConfig) internal s_chainConfigurations;

  /// @notice All chains that are configured.
  EnumerableSet.UintSet internal s_remoteChainSelectors;

  /// @notice OCR3 configurations for each DON.
  /// Each CR DON will have a commit and execution configuration.
  /// This means that a DON can have up to 4 configurations, since we are implementing blue/green deployments.
  mapping(uint32 donId => mapping(PluginType pluginType => OCR3ConfigWithMeta[] ocr3Configs)) internal s_ocr3Configs;

  /// @notice The DONs that have been configured.
  EnumerableSet.UintSet internal s_donIds;

  uint8 internal constant MAX_OCR3_CONFIGS_PER_PLUGIN = 2;
  uint8 internal constant MAX_OCR3_CONFIGS_PER_DON = 4;
  uint8 internal constant MAX_NUM_ORACLES = 31;

  /// @param capabilityRegistry the canonical capability registry address.
  constructor(address capabilityRegistry) {
    i_capabilityRegistry = capabilityRegistry;
  }

  // ================================================================
  // │                    Config Getters                            │
  // ================================================================

  /// @notice Returns all the chain configurations.
  /// @return The chain configurations.
  // TODO: will this eventually hit the RPC max response size limit?
  function getAllChainConfigs() external view returns (ChainConfigInfo[] memory) {
    uint256[] memory chainSelectors = s_remoteChainSelectors.values();
    ChainConfigInfo[] memory chainConfigs = new ChainConfigInfo[](s_remoteChainSelectors.length());
    for (uint256 i = 0; i < chainSelectors.length; ++i) {
      uint64 chainSelector = uint64(chainSelectors[i]);
      chainConfigs[i] =
        ChainConfigInfo({chainSelector: chainSelector, chainConfig: s_chainConfigurations[chainSelector]});
    }
    return chainConfigs;
  }

  /// @notice Returns the OCR configuration for the given don ID and plugin type.
  /// @param donId The DON ID.
  /// @param pluginType The plugin type.
  /// @return The OCR3 configurations, up to 2 (blue and green).
  function getOCRConfig(uint32 donId, PluginType pluginType) external view returns (OCR3ConfigWithMeta[] memory) {
    return s_ocr3Configs[donId][pluginType];
  }

  // ================================================================
  // │                    Capability Configuration                  │
  // ================================================================

  /// @inheritdoc ICapabilityConfiguration
  /// @dev The CCIP capability will fetch the configuration needed directly from this contract.
  /// The offchain syncer will call this function, however, so its important that it doesn't revert.
  function getCapabilityConfiguration(uint32 /* donId */ ) external pure override returns (bytes memory configuration) {
    return bytes("");
  }

  /// @notice Called by the registry prior to the config being set for a particular DON.
  function beforeCapabilityConfigSet(
    bytes32[] calldata, /* nodes */
    bytes calldata config,
    uint64, /* configCount */
    uint32 donId
  ) external override {
    if (msg.sender != i_capabilityRegistry) {
      revert OnlyCapabilityRegistryCanCall();
    }

    OCR3Config[] memory ocr3Configs = abi.decode(config, (OCR3Config[]));
    (OCR3Config[] memory commitConfigs, OCR3Config[] memory execConfigs) = _groupByPluginType(ocr3Configs);
    if (commitConfigs.length > 0) {
      _updatePluginConfig(donId, PluginType.Commit, commitConfigs);
    }
    if (execConfigs.length > 0) {
      _updatePluginConfig(donId, PluginType.Execution, execConfigs);
    }
  }

  function _updatePluginConfig(uint32 donId, PluginType pluginType, OCR3Config[] memory newConfig) internal {
    OCR3ConfigWithMeta[] memory currentConfig = s_ocr3Configs[donId][pluginType];

    // Validate the state transition being proposed, which is implicitly defined by the combination
    // of lengths of the current and new configurations.
    ConfigState currentState = _stateFromConfigLength(currentConfig.length);
    ConfigState proposedState = _stateFromConfigLength(newConfig.length);
    _validateConfigStateTransition(currentState, proposedState);

    // Build the new configuration with metadata and validate that the transition is valid.
    OCR3ConfigWithMeta[] memory newConfigWithMeta =
      _computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, proposedState);
    _validateConfigTransition(currentConfig, newConfigWithMeta);

    // Update contract state with new configuration if its valid.
    // We won't run out of gas from this delete since the array is at most 2 elements long.
    delete s_ocr3Configs[donId][pluginType];
    for (uint256 i = 0; i < newConfigWithMeta.length; ++i) {
      s_ocr3Configs[donId][pluginType].push(newConfigWithMeta[i]);
    }
  }

  // ================================================================
  // │                    Config State Machine                      │
  // ================================================================

  /// @notice Determine the config state of the configuration from the length of the config.
  /// @param configLen The length of the configuration.
  /// @return The config state.
  function _stateFromConfigLength(uint256 configLen) internal pure returns (ConfigState) {
    if (configLen > 2) {
      revert InvalidConfigLength(configLen);
    }
    return ConfigState(configLen);
  }

  // the only valid state transitions are the following:
  // init    -> running (first ever config)
  // running -> staging (blue/green proposal)
  // staging -> running (promotion)
  // everything else is invalid and should revert.
  function _validateConfigStateTransition(ConfigState currentState, ConfigState newState) internal pure {
    // Calculate the difference between the new state and the current state
    int256 stateDiff = int256(uint256(newState)) - int256(uint256(currentState));

    // Check if the state transition is valid:
    // Valid transitions:
    // 1. currentState -> newState (where stateDiff == 1)
    //    e.g., init -> running or running -> staging
    // 2. staging -> running (where stateDiff == -1)
    if (stateDiff == 1 || (stateDiff == -1 && currentState == ConfigState.Staging)) {
      return;
    }
    revert InvalidConfigStateTransition(currentState, newState);
  }

  function _validateConfigTransition(
    OCR3ConfigWithMeta[] memory currentConfig,
    OCR3ConfigWithMeta[] memory newConfigWithMeta
  ) internal pure {
    uint256 currentConfigLen = currentConfig.length;
    uint256 newConfigLen = newConfigWithMeta.length;
    if (currentConfigLen == 0 && newConfigLen == 1) {
      // Config counts always must start at 1 for the first ever config.
      if (newConfigWithMeta[0].configCount != 1) {
        revert WrongConfigCount(newConfigWithMeta[0].configCount, 1);
      }
      return;
    }

    if (currentConfigLen == 1 && newConfigLen == 2) {
      // On a blue/green proposal:
      // * the config digest of the blue config must remain unchanged.
      // * the green config count must be the blue config count + 1.
      if (newConfigWithMeta[0].configDigest != currentConfig[0].configDigest) {
        revert WrongConfigDigestBlueGreen(newConfigWithMeta[0].configDigest, currentConfig[0].configDigest);
      }
      if (newConfigWithMeta[1].configCount != currentConfig[0].configCount + 1) {
        revert WrongConfigCount(newConfigWithMeta[1].configCount, currentConfig[0].configCount + 1);
      }
      return;
    }

    if (currentConfigLen == 2 && newConfigLen == 1) {
      // On a promotion, the green config digest must become the blue config digest.
      if (newConfigWithMeta[0].configDigest != currentConfig[1].configDigest) {
        revert WrongConfigDigest(newConfigWithMeta[0].configDigest, currentConfig[1].configDigest);
      }
      return;
    }

    revert NonExistentConfigTransition();
  }

  /// @notice Computes a new configuration with metadata based on the current configuration and the new configuration.
  /// @param donId The DON ID.
  /// @param currentConfig The current configuration, including metadata.
  /// @param newConfig The new configuration, without metadata.
  /// @param currentState The current state of the configuration.
  /// @param newState The new state of the configuration.
  /// @return The new configuration with metadata.
  function _computeNewConfigWithMeta(
    uint32 donId,
    OCR3ConfigWithMeta[] memory currentConfig,
    OCR3Config[] memory newConfig,
    ConfigState currentState,
    ConfigState newState
  ) internal view returns (OCR3ConfigWithMeta[] memory) {
    uint64[] memory configCounts = new uint64[](newConfig.length);

    // Set config counts based on the only valid state transitions.
    // Init    -> Running (first ever config)
    // Running -> Staging (blue/green proposal)
    // Staging -> Running (promotion)
    if (currentState == ConfigState.Init && newState == ConfigState.Running) {
      // First ever config starts with config count == 1.
      configCounts[0] = 1;
    } else if (currentState == ConfigState.Running && newState == ConfigState.Staging) {
      // On a blue/green proposal, the config count of the green config is the blue config count + 1.
      configCounts[0] = currentConfig[0].configCount;
      configCounts[1] = currentConfig[0].configCount + 1;
    } else if (currentState == ConfigState.Staging && newState == ConfigState.Running) {
      // On a promotion, the config count of the green config becomes the blue config count.
      configCounts[0] = currentConfig[1].configCount;
    } else {
      revert InvalidConfigStateTransition(currentState, newState);
    }

    OCR3ConfigWithMeta[] memory newConfigWithMeta = new OCR3ConfigWithMeta[](newConfig.length);
    for (uint256 i = 0; i < configCounts.length; ++i) {
      _validateConfig(newConfig[i]);
      newConfigWithMeta[i] = OCR3ConfigWithMeta({
        config: newConfig[i],
        configCount: configCounts[i],
        configDigest: _computeConfigDigest(donId, configCounts[i], newConfig[i])
      });
    }

    return newConfigWithMeta;
  }

  /// @notice Group the OCR3 configurations by plugin type for further processing.
  /// @param ocr3Configs The OCR3 configurations to group.
  function _groupByPluginType(OCR3Config[] memory ocr3Configs)
    internal
    pure
    returns (OCR3Config[] memory commitConfigs, OCR3Config[] memory execConfigs)
  {
    if (ocr3Configs.length > MAX_OCR3_CONFIGS_PER_DON) {
      revert TooManyOCR3Configs();
    }

    // Declare with size 2 since we have a maximum of two configs per plugin type (blue, green).
    // If we have less we will adjust the length later using mstore.
    // If the caller provides more than 2 configs per plugin type, we will revert due to out of bounds
    // access in the for loop below.
    commitConfigs = new OCR3Config[](MAX_OCR3_CONFIGS_PER_PLUGIN);
    execConfigs = new OCR3Config[](MAX_OCR3_CONFIGS_PER_PLUGIN);
    uint256 commitCount;
    uint256 execCount;
    for (uint256 i = 0; i < ocr3Configs.length; ++i) {
      if (ocr3Configs[i].pluginType == PluginType.Commit) {
        commitConfigs[commitCount] = ocr3Configs[i];
        ++commitCount;
      } else {
        execConfigs[execCount] = ocr3Configs[i];
        ++execCount;
      }
    }

    // Adjust the length of the arrays to the actual number of configs.
    assembly {
      mstore(commitConfigs, commitCount)
      mstore(execConfigs, execCount)
    }

    return (commitConfigs, execConfigs);
  }

  function _validateConfig(OCR3Config memory cfg) internal view {
    if (cfg.chainSelector == 0) revert ChainSelectorNotSet();
    if (cfg.pluginType != PluginType.Commit && cfg.pluginType != PluginType.Execution) revert InvalidPluginType();
    // TODO: can we do more sophisticated validation than this?
    if (cfg.offrampAddress.length == 0) revert OfframpAddressCannotBeZero();
    if (!s_remoteChainSelectors.contains(cfg.chainSelector)) revert ChainSelectorNotFound(cfg.chainSelector);

    // Some of these checks below are done in OCR2/3Base config validation, so we do them again here.
    // Role DON OCR configs will have all the Role DON signers but only a subset of transmitters.
    if (cfg.signers.length > MAX_NUM_ORACLES) revert TooManySigners();
    if (cfg.transmitters.length > MAX_NUM_ORACLES) revert TooManyTransmitters();

    // We check for chain config presence above, so fChain here must be non-zero.
    uint256 minTransmittersLength = 3 * s_chainConfigurations[cfg.chainSelector].fChain + 1;
    if (cfg.transmitters.length < minTransmittersLength) {
      revert NotEnoughTransmitters(cfg.transmitters.length, minTransmittersLength);
    }
    if (cfg.F == 0) revert FMustBePositive();
    if (cfg.signers.length <= 3 * cfg.F) revert FTooHigh();

    if (cfg.p2pIds.length != cfg.signers.length || cfg.p2pIds.length != cfg.transmitters.length) {
      revert P2PIdsLengthNotMatching(cfg.p2pIds.length, cfg.signers.length, cfg.transmitters.length);
    }
    if (cfg.bootstrapP2PIds.length > cfg.p2pIds.length) revert TooManyBootstrapP2PIds();

    // Check that the readers are in the capability registry.
    // TODO: check for duplicate signers, duplicate p2p ids, etc.
    // TODO: check that p2p ids in cfg.bootstrapP2PIds are included in cfg.p2pIds.
    for (uint256 i = 0; i < cfg.signers.length; ++i) {
      _ensureInRegistry(cfg.p2pIds[i]);
    }
  }

  /// @notice Computes the digest of the provided configuration.
  /// @dev In traditional OCR config digest computation, block.chainid and address(this) are used
  /// in order to further domain separate the digest. We can't do that here since the digest will
  /// be used on remote chains; so we use the chain selector instead of block.chainid. The don ID
  /// replaces the address(this) in the traditional computation.
  /// @param donId The DON ID.
  /// @param configCount The configuration count.
  /// @param ocr3Config The OCR3 configuration.
  /// @return The computed digest.
  function _computeConfigDigest(
    uint32 donId,
    uint64 configCount,
    OCR3Config memory ocr3Config
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          ocr3Config.chainSelector,
          donId,
          ocr3Config.pluginType,
          ocr3Config.offrampAddress,
          configCount,
          ocr3Config.bootstrapP2PIds,
          ocr3Config.p2pIds,
          ocr3Config.signers,
          ocr3Config.transmitters,
          ocr3Config.F,
          ocr3Config.offchainConfigVersion,
          ocr3Config.offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x000a << (256 - 16); // 0x000a00..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  // ================================================================
  // │                    Chain Configuration                       │
  // ================================================================

  /// @notice Sets and/or removes chain configurations.
  /// @param chainSelectorRemoves The chain configurations to remove.
  /// @param chainConfigAdds The chain configurations to add.
  function applyChainConfigUpdates(
    uint64[] calldata chainSelectorRemoves,
    ChainConfigInfo[] calldata chainConfigAdds
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
      bytes32[] memory readers = chainConfig.readers;
      uint64 chainSelector = chainConfigAdds[i].chainSelector;

      // Verify that the provided readers are present in the capability registry.
      for (uint256 j = 0; j < readers.length; j++) {
        _ensureInRegistry(readers[j]);
      }

      // Verify that fChain is positive.
      if (chainConfig.fChain == 0) {
        revert FChainMustBePositive();
      }

      s_chainConfigurations[chainSelector] = chainConfig;
      s_remoteChainSelectors.add(chainSelector);

      emit ChainConfigSet(chainSelector, chainConfig);
    }
  }

  /// @notice Helper function to ensure that a node is in the capability registry.
  /// @param p2pId The P2P ID of the node to check.
  function _ensureInRegistry(bytes32 p2pId) internal view {
    (ICapabilityRegistry.NodeInfo memory node,) = ICapabilityRegistry(i_capabilityRegistry).getNode(p2pId);
    if (node.p2pId == bytes32("")) {
      revert NodeNotInRegistry(p2pId);
    }
  }
}
