// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {Internal} from "../../libraries/Internal.sol";

library CCIPConfigTypes {
  /// @notice ConfigState indicates the state of the configuration.
  /// A DON's configuration always starts out in the "Init" state - this is the starting state.
  /// The only valid transition from "Init" is to the "Running" state - this is the first ever configuration.
  /// The only valid transition from "Running" is to the "Staging" state - this is a blue/green proposal.
  /// The only valid transition from "Staging" is back to the "Running" state - this is a promotion.
  /// In order to rollback a configuration, we must therefore do the following:
  /// - Suppose that we have a correct configuration in the "Running" state (V1).
  /// - We propose a new configuration and transition to the "Staging" state (V2).
  /// - V2 turns out to be buggy
  /// - In the same transaction, we must:
  ///   - Promote V2
  ///   - Re-propose V1
  ///   - Promote V1
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
    Internal.OCRPluginType pluginType; // ────────╮ The plugin that the configuration is for.
    uint64 chainSelector; //                      | The (remote) chain that the configuration is for.
    uint8 FRoleDON; //                            | The "big F" parameter for the role DON.
    uint64 offchainConfigVersion; // ─────────────╯ The version of the offchain configuration.
    bytes offrampAddress; // The remote chain offramp address.
    OCR3Node[] nodes; // Keys & IDs of nodes part of the role DON
    bytes offchainConfig; // The offchain configuration for the OCR3 protocol. Protobuf encoded.
  }

  /// @notice OCR3 configuration with metadata, specifically the config count and the config digest.
  struct OCR3ConfigWithMeta {
    OCR3Config config; // The OCR3 configuration.
    uint64 configCount; // The config count used to compute the config digest.
    bytes32 configDigest; // The config digest of the OCR3 configuration.
  }
}
