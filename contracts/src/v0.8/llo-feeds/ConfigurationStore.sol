// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

struct ChannelDefinition {
    // e.g. evm, solana, CosmWasm, kalechain, etc...
    string reportFormat;
	// Specifies the chain on which this channel can be verified. Currently uses
	// CCIP chain selectors, but lots of other schemes are possible as well.
	uint64 chainSelector;
    // We assume that StreamIDs is always non-empty and that the 0-th stream
	// contains the verification price in LINK and the 1-st stream contains the
	// verification price in the native coin.
	string[] streamIDs;
}

contract ConfigurationStore {
    ////////////////////////
    // protocol instance management
    ////////////////////////

    ChannelDefinition[] private s_channelDefinitions;

    // setProductionConfig() onlyOwner -- the usual OCR way
    // sets config for the production protocol instance

    // setStagingConfig() onlyOwner -- the usual OCR way
    // sets config for the staging protocol instance

    // promoteStagingConfig() onlyOwner 
    // this will trigger the following:
    // - offchain ShouldRetireCache will start returning true for the old (production)
    //   protocol instance
    // - once the old production instance retires it will generate a handover
    //   retirement report
    // - the staging instance will become the new production instance once 
    //   any honest oracle that is on both instances forward the retirement 
    //   report from the old instance to the new instace via the 
    //   PredecessorRetirementReportCache
    //
    // Note: the promotion flow only works if the previous production instance
    // is working correctly & generating reports. If that's not the case, the 
    // owner is expected to "setProductionConfig" directly instead. This will
    // cause "gaps" to be created, but that seems unavoidable in such a scenario.

    ////////////////////////
    // channel management
    ////////////////////////

    addChannel(ChannelDefinition) onlyOwner {
        // TODO
    }

    removeChannel(bytes32 channelId) onlyOwner {
        // TODO

    }

    getChannelDefinitions() onlyEOA public view returns (ChannelDefinition[] memory) {
        // TODO

    }
    // used by ChannelDefinitionCache
}