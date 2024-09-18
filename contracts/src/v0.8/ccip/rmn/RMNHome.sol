// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

/// @notice Stores the home configuration for RMN, that is referenced by CCIP oracles, RMN nodes, and the RMNRemote
/// contracts.
contract RMNHome is OwnerIsCreator, ITypeAndVersion {
  error DuplicatePeerId();
  error DuplicateOffchainPublicKey();
  error OutOfOrderSourceChains();
  error OutOfOrderObserverNodeIndices();
  error OutOfBoundsObserverNodeIndex();
  error MinObserversTooHigh();

  event ConfigSet(bytes32 configDigest, VersionedConfig versionedConfig);
  event ConfigRevoked(bytes32 configDigest);

  struct Node {
    string peerId; // used for p2p communication, base58 encoded
    bytes32 offchainPublicKey; // observations are signed with this public key, and are only verified offchain
  }

  struct SourceChain {
    uint64 chainSelector;
    uint64[] observerNodeIndices; // indices into Config.nodes, strictly increasing
    uint64 minObservers; // required to agree on an observation for this source chain
  }

  struct Config {
    // No sorting requirement for nodes, but ensure that SourceChain.observerNodeIndices in the home chain config &
    // Signer.nodeIndex in the remote chain configs are appropriately updated when changing this field
    Node[] nodes;
    // Should be in ascending order of chainSelector
    SourceChain[] sourceChains;
  }

  struct VersionedConfig {
    uint32 version;
    Config config;
  }

  string public constant override typeAndVersion = "RMNHome 1.6.0-dev";
  uint256 public constant CONFIG_RING_BUFFER_SIZE = 2;

  function _configDigest(VersionedConfig memory versionedConfig) internal pure returns (bytes32) {
    uint256 h = uint256(keccak256(abi.encode(versionedConfig)));
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x000b << (256 - 16); // 0x000b00..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }

  // if we were to have VersionedConfig instead of Config in the ring buffer, we couldn't assign directly to it in
  // setConfig without via-ir
  uint32[CONFIG_RING_BUFFER_SIZE] s_configCounts; // s_configCounts[i] == 0 iff s_configs[i] is unusable
  Config[CONFIG_RING_BUFFER_SIZE] s_configs;
  uint256 s_latestConfigIndex;
  bytes32 s_latestConfigDigest;

  /// @param revokePastConfigs if one wants to revoke all past configs, because some past config is faulty
  function setConfig(Config calldata newConfig, bool revokePastConfigs) external onlyOwner {
    // sanity checks
    {
      // no peerId or offchainPublicKey is duplicated
      for (uint256 i = 0; i < newConfig.nodes.length; ++i) {
        for (uint256 j = i + 1; j < newConfig.nodes.length; ++j) {
          if (keccak256(abi.encode(newConfig.nodes[i].peerId)) == keccak256(abi.encode(newConfig.nodes[j].peerId))) {
            revert DuplicatePeerId();
          }
          if (newConfig.nodes[i].offchainPublicKey == newConfig.nodes[j].offchainPublicKey) {
            revert DuplicateOffchainPublicKey();
          }
        }
      }

      for (uint256 i = 0; i < newConfig.sourceChains.length; ++i) {
        // source chains are in strictly increasing order of chain selectors
        if (i > 0 && !(newConfig.sourceChains[i - 1].chainSelector < newConfig.sourceChains[i].chainSelector)) {
          revert OutOfOrderSourceChains();
        }

        // all observerNodeIndices are valid
        for (uint256 j = 0; j < newConfig.sourceChains[i].observerNodeIndices.length; ++j) {
          if (
            j > 0
              && !(newConfig.sourceChains[i].observerNodeIndices[j - 1] < newConfig.sourceChains[i].observerNodeIndices[j])
          ) {
            revert OutOfOrderObserverNodeIndices();
          }
          if (!(newConfig.sourceChains[i].observerNodeIndices[j] < newConfig.nodes.length)) {
            revert OutOfBoundsObserverNodeIndex();
          }
        }

        // minObservers are tenable
        if (!(newConfig.sourceChains[i].minObservers <= newConfig.sourceChains[i].observerNodeIndices.length)) {
          revert MinObserversTooHigh();
        }
      }
    }

    uint256 oldConfigIndex = s_latestConfigIndex;
    uint32 oldConfigCount = s_configCounts[oldConfigIndex];
    uint256 newConfigIndex = (oldConfigIndex + 1) % CONFIG_RING_BUFFER_SIZE;

    for (uint256 i = 0; i < CONFIG_RING_BUFFER_SIZE; ++i) {
      if ((i == newConfigIndex || revokePastConfigs) && s_configCounts[i] > 0) {
        emit ConfigRevoked(_configDigest(VersionedConfig({version: s_configCounts[i], config: s_configs[i]})));
        delete s_configCounts[i];
      }
    }

    uint32 newConfigCount = oldConfigCount + 1;
    VersionedConfig memory newVersionedConfig = VersionedConfig({version: newConfigCount, config: newConfig});
    bytes32 newConfigDigest = _configDigest(newVersionedConfig);
    s_configs[newConfigIndex] = newConfig;
    s_configCounts[newConfigIndex] = newConfigCount;
    s_latestConfigIndex = newConfigIndex;
    s_latestConfigDigest = newConfigDigest;
    emit ConfigSet(newConfigDigest, newVersionedConfig);
  }

  /// @return configDigest will be zero in case no config has been set
  function getLatestConfigDigestAndVersionedConfig()
    external
    view
    returns (bytes32 configDigest, VersionedConfig memory)
  {
    return (
      s_latestConfigDigest,
      VersionedConfig({version: s_configCounts[s_latestConfigIndex], config: s_configs[s_latestConfigIndex]})
    );
  }

  /// @notice The offchain code can use this to fetch an old config which might still be in use by some remotes
  /// @dev Only to be called by offchain code, efficiency is not a concern
  function getConfig(bytes32 configDigest) external view returns (VersionedConfig memory versionedConfig, bool ok) {
    for (uint256 i = 0; i < CONFIG_RING_BUFFER_SIZE; ++i) {
      if (s_configCounts[i] == 0) {
        // unset config
        continue;
      }
      VersionedConfig memory vc = VersionedConfig({version: s_configCounts[i], config: s_configs[i]});
      if (_configDigest(vc) == configDigest) {
        versionedConfig = vc;
        ok = true;
        break;
      }
    }
  }
}
