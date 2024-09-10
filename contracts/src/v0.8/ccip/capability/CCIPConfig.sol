// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../keystone/interfaces/ICapabilityConfiguration.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {ICapabilitiesRegistry} from "./interfaces/ICapabilitiesRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Internal} from "../libraries/Internal.sol";
import {CCIPConfigTypes} from "./libraries/CCIPConfigTypes.sol";

import {IERC165} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/interfaces/IERC165.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice CCIPConfig stores the configuration for the CCIP capability.
/// We have two classes of configuration: chain configuration and DON (in the CapabilitiesRegistry sense) configuration.
/// Each chain will have a single configuration which includes information like the router address.
/// Each CR DON will have up to four configurations: for each of (commit, exec), one blue and one green configuration.
/// This is done in order to achieve "blue-green" deployments.
contract CCIPConfig is ITypeAndVersion, ICapabilityConfiguration, OwnerIsCreator, IERC165 {
  using EnumerableSet for EnumerableSet.UintSet;

  /// @notice Emitted when a chain's configuration is set.
  /// @param chainSelector The chain selector.
  /// @param chainConfig The chain configuration.
  event ChainConfigSet(uint64 chainSelector, CCIPConfigTypes.ChainConfig chainConfig);
  /// @notice Emitted when a chain's configuration is removed.
  /// @param chainSelector The chain selector.
  event ChainConfigRemoved(uint64 chainSelector);

  error NodeNotInRegistry(bytes32 p2pId);
  error OnlyCapabilitiesRegistryCanCall();
  error ChainSelectorNotFound(uint64 chainSelector);
  error ChainSelectorNotSet();
  error TooManyOCR3Configs();
  error TooManySigners();
  error P2PIdsLengthNotMatching(uint256 p2pIdsLength, uint256 signersLength, uint256 transmittersLength);
  error NotEnoughTransmitters(uint256 got, uint256 minimum);
  error FMustBePositive();
  error FChainMustBePositive();
  error FTooHigh();
  error InvalidPluginType();
  error OfframpAddressCannotBeZero();
  error InvalidConfigLength(uint256 length);
  error InvalidConfigStateTransition(
    CCIPConfigTypes.ConfigState currentState, CCIPConfigTypes.ConfigState proposedState
  );
  error NonExistentConfigTransition();
  error WrongConfigCount(uint64 got, uint64 expected);
  error WrongConfigDigest(bytes32 got, bytes32 expected);
  error WrongConfigDigestBlueGreen(bytes32 got, bytes32 expected);
  error ZeroAddressNotAllowed();

  /// @notice Type and version override.
  string public constant override typeAndVersion = "CCIPConfig 1.6.0-dev";

  /// @notice The canonical capabilities registry address.
  address internal immutable i_capabilitiesRegistry;

  uint8 internal constant MAX_OCR3_CONFIGS_PER_PLUGIN = 2;
  uint8 internal constant MAX_OCR3_CONFIGS_PER_DON = 4;
  uint256 internal constant CONFIG_DIGEST_PREFIX_MASK = type(uint256).max << (256 - 16); // 0xFFFF00..0
  /// @dev must be equal to libocr multi role: https://github.com/smartcontractkit/libocr/blob/ae747ca5b81236ffdbf1714318c652e923a5ff4d/offchainreporting2plus/types/config_digest.go#L28
  uint256 internal constant CONFIG_DIGEST_PREFIX = 0x000a << (256 - 16); // 0x000a00..00
  bytes32 internal constant EMPTY_ENCODED_ADDRESS_HASH = keccak256(abi.encode(address(0)));
  /// @dev 256 is the hard limit due to the bit encoding of their indexes into a uint256.
  uint256 internal constant MAX_NUM_ORACLES = 256;

  /// @notice chain configuration for each chain that CCIP is deployed on.
  mapping(uint64 chainSelector => CCIPConfigTypes.ChainConfig chainConfig) internal s_chainConfigurations;

  /// @notice All chains that are configured.
  EnumerableSet.UintSet internal s_remoteChainSelectors;

  /// @notice OCR3 configurations for each DON.
  /// Each CR DON will have a commit and execution configuration.
  /// This means that a DON can have up to 4 configurations, since we are implementing blue/green deployments.
  mapping(
    uint32 donId => mapping(Internal.OCRPluginType pluginType => CCIPConfigTypes.OCR3ConfigWithMeta[] ocr3Configs)
  ) internal s_ocr3Configs;

  /// @param capabilitiesRegistry the canonical capabilities registry address.
  constructor(address capabilitiesRegistry) {
    if (capabilitiesRegistry == address(0)) {
      revert ZeroAddressNotAllowed();
    }
    i_capabilitiesRegistry = capabilitiesRegistry;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == type(ICapabilityConfiguration).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  // ================================================================
  // │                    Config Getters                            │
  // ================================================================
  /// @notice Returns the capabilities registry address.
  /// @return The capabilities registry address.
  function getCapabilityRegistry() external view returns (address) {
    return i_capabilitiesRegistry;
  }

  /// @notice Returns the total number of chains configured.
  /// @return The total number of chains configured.
  function getNumChainConfigurations() external view returns (uint256) {
    return s_remoteChainSelectors.length();
  }

  /// @notice Returns all the chain configurations.
  /// @param pageIndex The page index.
  /// @param pageSize The page size.
  /// @return paginatedChainConfigs chain configurations.
  function getAllChainConfigs(
    uint256 pageIndex,
    uint256 pageSize
  ) external view returns (CCIPConfigTypes.ChainConfigInfo[] memory) {
    uint256 totalItems = s_remoteChainSelectors.length(); // Total number of chain selectors
    uint256 startIndex = pageIndex * pageSize;

    if (pageSize == 0 || startIndex >= totalItems) {
      return new CCIPConfigTypes.ChainConfigInfo[](0); // Return an empty array if pageSize is 0 or pageIndex is out of bounds
    }

    uint256 endIndex = startIndex + pageSize;
    if (endIndex > totalItems) {
      endIndex = totalItems;
    }

    CCIPConfigTypes.ChainConfigInfo[] memory paginatedChainConfigs =
      new CCIPConfigTypes.ChainConfigInfo[](endIndex - startIndex);

    uint256[] memory chainSelectors = s_remoteChainSelectors.values();
    for (uint256 i = startIndex; i < endIndex; ++i) {
      uint64 chainSelector = uint64(chainSelectors[i]);
      paginatedChainConfigs[i - startIndex] = CCIPConfigTypes.ChainConfigInfo({
        chainSelector: chainSelector,
        chainConfig: s_chainConfigurations[chainSelector]
      });
    }

    return paginatedChainConfigs;
  }

  /// @notice Returns the OCR configuration for the given don ID and plugin type.
  /// @param donId The DON ID.
  /// @param pluginType The plugin type.
  /// @return The OCR3 configurations, up to 2 (blue and green).
  function getOCRConfig(
    uint32 donId,
    Internal.OCRPluginType pluginType
  ) external view returns (CCIPConfigTypes.OCR3ConfigWithMeta[] memory) {
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
    if (msg.sender != i_capabilitiesRegistry) {
      revert OnlyCapabilitiesRegistryCanCall();
    }

    (CCIPConfigTypes.OCR3Config[] memory commitConfigs, CCIPConfigTypes.OCR3Config[] memory execConfigs) =
      _groupByPluginType(abi.decode(config, (CCIPConfigTypes.OCR3Config[])));
    if (commitConfigs.length > 0) {
      _updatePluginConfig(donId, Internal.OCRPluginType.Commit, commitConfigs);
    }
    if (execConfigs.length > 0) {
      _updatePluginConfig(donId, Internal.OCRPluginType.Execution, execConfigs);
    }
  }

  /// @notice Sets a new OCR3 config for a specific plugin type for a DON.
  /// @param donId The DON ID.
  /// @param pluginType The plugin type.
  /// @param newConfig The new configuration.
  function _updatePluginConfig(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    CCIPConfigTypes.OCR3Config[] memory newConfig
  ) internal {
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig = s_ocr3Configs[donId][pluginType];

    // Validate the state transition being proposed, which is implicitly defined by the combination
    // of lengths of the current and new configurations.
    CCIPConfigTypes.ConfigState currentState = _stateFromConfigLength(currentConfig.length);
    CCIPConfigTypes.ConfigState proposedState = _stateFromConfigLength(newConfig.length);
    _validateConfigStateTransition(currentState, proposedState);

    // Build the new configuration with metadata and validate that the transition is valid.
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta =
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
  function _stateFromConfigLength(uint256 configLen) internal pure returns (CCIPConfigTypes.ConfigState) {
    if (configLen > 2) {
      revert InvalidConfigLength(configLen);
    }
    return CCIPConfigTypes.ConfigState(configLen);
  }

  /// @notice Validates the state transition between two config states.
  /// The only valid state transitions are the following:
  /// Init    -> Running (first ever config)
  /// Running -> Staging (blue/green proposal)
  /// Staging -> Running (promotion)
  /// Everything else is invalid and should revert.
  /// @param currentState The current state.
  /// @param newState The new state.
  function _validateConfigStateTransition(
    CCIPConfigTypes.ConfigState currentState,
    CCIPConfigTypes.ConfigState newState
  ) internal pure {
    // Calculate the difference between the new state and the current state
    int256 stateDiff = int256(uint256(newState)) - int256(uint256(currentState));

    // Check if the state transition is valid:
    // Valid transitions:
    // 1. currentState -> newState (where stateDiff == 1)
    //    e.g., init -> running or running -> staging
    // 2. staging -> running (where stateDiff == -1)
    if (stateDiff == 1 || (stateDiff == -1 && currentState == CCIPConfigTypes.ConfigState.Staging)) {
      return;
    }
    revert InvalidConfigStateTransition(currentState, newState);
  }

  /// @notice Validates the transition between two OCR3 configurations.
  /// @param currentConfig The current configuration with metadata.
  /// @param newConfigWithMeta The new configuration with metadata.
  function _validateConfigTransition(
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig,
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta
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
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig,
    CCIPConfigTypes.OCR3Config[] memory newConfig,
    CCIPConfigTypes.ConfigState currentState,
    CCIPConfigTypes.ConfigState newState
  ) internal view returns (CCIPConfigTypes.OCR3ConfigWithMeta[] memory) {
    uint64[] memory configCounts = new uint64[](newConfig.length);

    // Set config counts based on the only valid state transitions.
    // Init    -> Running (first ever config)
    // Running -> Staging (blue/green proposal)
    // Staging -> Running (promotion)
    if (currentState == CCIPConfigTypes.ConfigState.Init && newState == CCIPConfigTypes.ConfigState.Running) {
      // First ever config starts with config count == 1.
      configCounts[0] = 1;
    } else if (currentState == CCIPConfigTypes.ConfigState.Running && newState == CCIPConfigTypes.ConfigState.Staging) {
      // On a blue/green proposal, the config count of the green config is the blue config count + 1.
      configCounts[0] = currentConfig[0].configCount;
      configCounts[1] = currentConfig[0].configCount + 1;
    } else if (currentState == CCIPConfigTypes.ConfigState.Staging && newState == CCIPConfigTypes.ConfigState.Running) {
      // On a promotion, the config count of the green config becomes the blue config count.
      configCounts[0] = currentConfig[1].configCount;
    } else {
      revert InvalidConfigStateTransition(currentState, newState);
    }

    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta =
      new CCIPConfigTypes.OCR3ConfigWithMeta[](newConfig.length);
    for (uint256 i = 0; i < configCounts.length; ++i) {
      _validateConfig(newConfig[i]);
      newConfigWithMeta[i] = CCIPConfigTypes.OCR3ConfigWithMeta({
        config: newConfig[i],
        configCount: configCounts[i],
        configDigest: _computeConfigDigest(donId, configCounts[i], newConfig[i])
      });
    }

    return newConfigWithMeta;
  }

  /// @notice Group the OCR3 configurations by plugin type for further processing.
  /// @param ocr3Configs The OCR3 configurations to group.
  /// @return commitConfigs The commit configurations.
  /// @return execConfigs The execution configurations.
  function _groupByPluginType(
    CCIPConfigTypes.OCR3Config[] memory ocr3Configs
  )
    internal
    pure
    returns (CCIPConfigTypes.OCR3Config[] memory commitConfigs, CCIPConfigTypes.OCR3Config[] memory execConfigs)
  {
    if (ocr3Configs.length > MAX_OCR3_CONFIGS_PER_DON) {
      revert TooManyOCR3Configs();
    }

    // Declare with size 2 since we have a maximum of two configs per plugin type (blue, green).
    // If we have less we will adjust the length later using mstore.
    // If the caller provides more than 2 configs per plugin type, we will revert due to out of bounds
    // access in the for loop below.
    commitConfigs = new CCIPConfigTypes.OCR3Config[](MAX_OCR3_CONFIGS_PER_PLUGIN);
    execConfigs = new CCIPConfigTypes.OCR3Config[](MAX_OCR3_CONFIGS_PER_PLUGIN);
    uint256 commitCount = 0;
    uint256 execCount = 0;
    for (uint256 i = 0; i < ocr3Configs.length; ++i) {
      if (ocr3Configs[i].pluginType == Internal.OCRPluginType.Commit) {
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

  /// @notice Validates an OCR3 configuration.
  /// @param cfg The OCR3 configuration.
  function _validateConfig(CCIPConfigTypes.OCR3Config memory cfg) internal view {
    if (cfg.chainSelector == 0) revert ChainSelectorNotSet();
    if (cfg.pluginType != Internal.OCRPluginType.Commit && cfg.pluginType != Internal.OCRPluginType.Execution) {
      revert InvalidPluginType();
    }
    if (cfg.offrampAddress.length == 0 || keccak256(cfg.offrampAddress) == EMPTY_ENCODED_ADDRESS_HASH) {
      revert OfframpAddressCannotBeZero();
    }
    if (!s_remoteChainSelectors.contains(cfg.chainSelector)) revert ChainSelectorNotFound(cfg.chainSelector);

    // We check for chain config presence above, so fChain here must be non-zero.
    uint256 minTransmittersLength = 3 * s_chainConfigurations[cfg.chainSelector].fChain + 1;
    if (cfg.transmitters.length < minTransmittersLength) {
      revert NotEnoughTransmitters(cfg.transmitters.length, minTransmittersLength);
    }
    uint256 numberOfSigners = cfg.signers.length;
    if (numberOfSigners > MAX_NUM_ORACLES) revert TooManySigners();
    if (numberOfSigners != cfg.p2pIds.length || numberOfSigners != cfg.transmitters.length) {
      revert P2PIdsLengthNotMatching(cfg.p2pIds.length, cfg.signers.length, cfg.transmitters.length);
    }
    if (cfg.F == 0) revert FMustBePositive();
    if (numberOfSigners <= 3 * cfg.F) revert FTooHigh();

    // Check that the readers are in the capabilities registry.
    _ensureInRegistry(cfg.p2pIds);
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
    CCIPConfigTypes.OCR3Config memory ocr3Config
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          ocr3Config.chainSelector,
          donId,
          ocr3Config.pluginType,
          ocr3Config.offrampAddress,
          configCount,
          ocr3Config.p2pIds,
          ocr3Config.signers,
          ocr3Config.transmitters,
          ocr3Config.F,
          ocr3Config.offchainConfigVersion,
          ocr3Config.offchainConfig
        )
      )
    );

    return bytes32((CONFIG_DIGEST_PREFIX & CONFIG_DIGEST_PREFIX_MASK) | (h & ~CONFIG_DIGEST_PREFIX_MASK));
  }

  // ================================================================
  // │                    Chain Configuration                       │
  // ================================================================

  /// @notice Sets and/or removes chain configurations.
  /// @param chainSelectorRemoves The chain configurations to remove.
  /// @param chainConfigAdds The chain configurations to add.
  function applyChainConfigUpdates(
    uint64[] calldata chainSelectorRemoves,
    CCIPConfigTypes.ChainConfigInfo[] calldata chainConfigAdds
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
      CCIPConfigTypes.ChainConfig memory chainConfig = chainConfigAdds[i].chainConfig;
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
