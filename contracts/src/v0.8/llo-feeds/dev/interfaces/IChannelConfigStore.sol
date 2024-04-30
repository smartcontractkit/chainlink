// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

interface IChannelConfigStore is IERC165 {
  // function setStagingConfig(bytes32 configDigest, ChannelConfiguration calldata channelConfig) external;

  // function promoteStagingConfig(bytes32 configDigest) external;

  function addChannel(uint32 channelId, ChannelDefinition calldata channelDefinition) external;

  function removeChannel(uint32 channelId) external;

  function getChannelDefinitions(uint32 channelId) external view returns (ChannelDefinition memory);

  // struct ChannelConfiguration {
  //     bytes32 configDigest;
  // }

  struct ChannelDefinition {
    // e.g. evm, solana, CosmWasm, kalechain, etc...
    uint32 reportFormat;
    // Specifies the chain on which this channel can be verified. Currently uses
    // CCIP chain selectors, but lots of other schemes are possible as well.
    uint64 chainSelector;
    // We assume that StreamIDs is always non-empty and that the 0-th stream
    // contains the verification price in LINK and the 1-st stream contains the
    // verification price in the native coin.
    uint32[] streamIDs;
  }
}
