// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

interface IStreamConfigStore is IERC165 {

    function setStagingConfig(bytes32 channelId, ChannelConfiguration calldata channelConfig) external;

    function promoteStagingConfig(bytes32 channelId) external;

    function addChannel(bytes32 channelId, ChannelDefinition calldata channelDefinition) external;

    function removeChannel(bytes32 channelId) external;

    function getChannelDefinitions(bytes32 channelId) external view returns (ChannelDefinition memory);

    struct ChannelConfiguration {
        bytes32 channelConfigId;
    }

    struct ChannelDefinition {
        // e.g. evm, solana, CosmWasm, kalechain, etc...
        bytes8 reportFormat;
        // Specifies the chain on which this channel can be verified. Currently uses
        // CCIP chain selectors, but lots of other schemes are possible as well.
        uint64 chainSelector;
        // We assume that StreamIDs is always non-empty and that the 0-th stream
        // contains the verification price in LINK and the 1-st stream contains the
        // verification price in the native coin.
        bytes32[] streamIDs;
    }

}
