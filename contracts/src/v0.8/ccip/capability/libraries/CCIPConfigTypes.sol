// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {Internal} from "../../libraries/Internal.sol";

library CCIPConfigTypes {
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
    bytes32[] readers; // The P2P IDs of the readers for the chain. These IDs must be registered in the capabilities registry.
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
    Internal.OCRPluginType pluginType; // ────────╮ The plugin that the configuration is for.
    uint64 chainSelector; //                      | The (remote) chain that the configuration is for.
    uint8 F; //                                   | The "big F" parameter for the role DON.
    uint64 offchainConfigVersion; // ─────────────╯ The version of the offchain configuration.
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
}
