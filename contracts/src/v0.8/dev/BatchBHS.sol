// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

/**
 * @title BatchBHS
 * @notice The BatchBHS contract acts as a proxy to write many blockhashes to the
 *   provided BlockhashStore contract efficiently in a single transaction. This results
 *   in plenty of gas savings and higher throughput of blockhash storage, which is desirable
 *   in times of high network congestion. 
*/
contract BatchBHS {
    BlockhashStore public immutable BHS;

    constructor(address blockhashStoreAddr) {
        BHS = BlockhashStore(blockhashStoreAddr);
    }

    /**
     * @notice stores blockhashes of the given block numbers in the configured blockhash store, assuming
     *   they are availble though the blockhash() instruction.
     * @param blockNumbers the block numbers to store the blockhashes of. Must be available via the
     *   blockhash() instruction, otherwise this function call will revert.
    */
    function store(uint256[] memory blockNumbers) public {
        for (uint256 i = 0; i < blockNumbers.length; i++) {
            BHS.store(blockNumbers[i]);
        }
    }

    /**
     * @notice stores blockhashes after verifying blockheader of child/subsequent block
     * @param blockNumbers the block numbers whose blockhashes should be stored, in decreasing order
     * @param headers the rlp-encoded block headers of blockNumbers[i] + 1.
    */
    function storeVerifyHeader(uint256[] memory blockNumbers, bytes[] memory headers) public {
        require(blockNumbers.length == headers.length, "input array arg lengths mismatch");
        for (uint256 i = 0; i < blockNumbers.length; i++) {
            BHS.storeVerifyHeader(blockNumbers[i], headers[i]);
        }
    }

    /**
     * @notice retrieves blockhashes of all the given block numbers from the blockhash store, if available.
     * @param blockNumbers array of block numbers to fetch blockhashes for
    */
    function getBlockhashes(uint256[] memory blockNumbers) external view returns (bytes32[] memory) {
        bytes32[] memory blockHashes = new bytes32[](blockNumbers.length);
        for (uint256 i = 0; i < blockNumbers.length; i++) {
            blockHashes[i] = BHS.getBlockhash(blockNumbers[i]);
        }
        return blockHashes;
    }
}

interface BlockhashStore {
    function storeVerifyHeader(uint256 n, bytes memory header) external;
    function store(uint256 n) external;
    function getBlockhash(uint256 n) external view returns (bytes32);
}
