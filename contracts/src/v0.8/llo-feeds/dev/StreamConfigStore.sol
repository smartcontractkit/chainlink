// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.16;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IChannelConfigStore} from "./interfaces/IChannelConfigStore.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";

contract StreamConfigStore is ConfirmedOwner, IChannelConfigStore, TypeAndVersionInterface {

    mapping(bytes32 => ChannelDefinition) private s_channelDefinitions;

    mapping(bytes32 => ChannelConfiguration) private s_channelProductionConfigurations;
    mapping(bytes32 => ChannelConfiguration) private s_channelStagingConfigurations;

    event NewChannelDefinition(bytes32 channelId, ChannelDefinition channelDefinition);
    event ChannelDefinitionRemoved(bytes32 channelId);
    event NewProductionConfig(ChannelConfiguration channelConfig);
    event NewStagingConfig(ChannelConfiguration channelConfig);
    event PromoteStagingConfig(bytes32 channelId);

    error OnlyCallableByEOA();
    error StagingConfigAlreadyPromoted();
    error EmptyStreamIDs();
    error ZeroReportFormat();
    error ZeroChainSelector();
    error ChannelDefinitionNotFound();

    constructor() ConfirmedOwner(msg.sender) {}

    function setStagingConfig(bytes32 channelId, ChannelConfiguration calldata channelConfig) external onlyOwner {
        s_channelStagingConfigurations[channelId] = channelConfig;

        emit NewStagingConfig(channelConfig);
    }

    //// this will trigger the following:
    //// - offchain ShouldRetireCache will start returning true for the old (production)
    ////   protocol instance
    //// - once the old production instance retires it will generate a handover
    ////   retirement report
    //// - the staging instance will become the new production instance once
    ////   any honest oracle that is on both instances forward the retirement
    ////   report from the old instance to the new instace via the
    ////   PredecessorRetirementReportCache
    ////
    //// Note: the promotion flow only works if the previous production instance
    //// is working correctly & generating reports. If that's not the case, the
    //// owner is expected to "setProductionConfig" directly instead. This will
    //// cause "gaps" to be created, but that seems unavoidable in such a scenario.
    function promoteStagingConfig(bytes32 channelId) external onlyOwner {
        ChannelConfiguration memory stagingConfig = s_channelStagingConfigurations[channelId];

        if(stagingConfig.channelConfigId.length == 0) {
            revert StagingConfigAlreadyPromoted();
        }

        s_channelProductionConfigurations[channelId] = s_channelStagingConfigurations[channelId];

        emit PromoteStagingConfig(channelId);
    }

    function addChannel(bytes32 channelId, ChannelDefinition calldata channelDefinition) external onlyOwner {

        if(channelDefinition.streamIDs.length == 0) {
            revert EmptyStreamIDs();
        }

        if(channelDefinition.chainSelector == 0) {
            revert ZeroChainSelector();
        }

        if(channelDefinition.reportFormat == 0) {
            revert ZeroReportFormat();
        }

        s_channelDefinitions[channelId] = channelDefinition;

        emit NewChannelDefinition(channelId, channelDefinition);
    }

    function removeChannel(bytes32 channelId) external onlyOwner {
        if(s_channelDefinitions[channelId].streamIDs.length == 0) {
            revert ChannelDefinitionNotFound();
        }

        delete s_channelDefinitions[channelId];

        emit ChannelDefinitionRemoved(channelId);
    }

    function getChannelDefinitions(bytes32 channelId) external view returns (ChannelDefinition memory) {
        if(msg.sender != tx.origin) {
            revert OnlyCallableByEOA();
        }

        return s_channelDefinitions[channelId];
    }

    function typeAndVersion() external override pure returns (string memory) {
        return "StreamConfigStore 0.0.0";
    }

    function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
        return interfaceId == type(IChannelConfigStore).interfaceId;
    }
}

