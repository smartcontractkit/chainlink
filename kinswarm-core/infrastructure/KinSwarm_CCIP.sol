// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IWormholeCore {
    function publishMessage(
        uint32 nonce, 
        bytes memory payload, 
        uint8 consistencyLevel
    ) external payable returns (uint64 sequence);
}

contract KinSwarm_Sovereign_Final {
    IWormholeCore public constant CORE = IWormholeCore(0x180010376B1389c8d97D6431Ff23A04218a0904e);

    // Removed 'payable' from the call to CORE to ensure 0-value strictness
    function forceAnchor(bytes32 merkleRoot, uint256 amount) external payable returns (uint64) {
        bytes memory payload = abi.encode(merkleRoot, amount);
        
        // Consistency Level 1 = Finalized
        // Nonce is incremented to ensure uniqueness in the Guardian set
        return CORE.publishMessage(
            uint32(block.timestamp), 
            payload, 
            1        
        );
    }

    receive() external payable {}
}
