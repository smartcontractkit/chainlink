// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IPubdataChunkPublisher} from "./interfaces/IPubdataChunkPublisher.sol";
import {SystemContractBase} from "./abstract/SystemContractBase.sol";
import {L1_MESSENGER_CONTRACT, BLOB_SIZE_BYTES, MAX_NUMBER_OF_BLOBS, SystemLogKey} from "./Constants.sol";
import {SystemContractHelper} from "./libraries/SystemContractHelper.sol";
import {TooMuchPubdata} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice Smart contract for chunking pubdata into the appropriate size for EIP-4844 blobs.
 */
contract PubdataChunkPublisher is IPubdataChunkPublisher, SystemContractBase {
    /// @notice Chunks pubdata into pieces that can fit into blobs.
    /// @param _pubdata The total l2 to l1 pubdata that will be sent via L1 blobs.
    /// @dev Note: This is an early implementation, in the future we plan to support up to 16 blobs per l1 batch.
    /// @dev We always publish 6 system logs even if our pubdata fits into a single blob. This makes processing logs on L1 easier.
    function chunkAndPublishPubdata(bytes calldata _pubdata) external onlyCallFrom(address(L1_MESSENGER_CONTRACT)) {
        if (_pubdata.length > BLOB_SIZE_BYTES * MAX_NUMBER_OF_BLOBS) {
            revert TooMuchPubdata(BLOB_SIZE_BYTES * MAX_NUMBER_OF_BLOBS, _pubdata.length);
        }

        bytes32[] memory blobHashes = new bytes32[](MAX_NUMBER_OF_BLOBS);

        // We allocate to the full size of MAX_NUMBER_OF_BLOBS * BLOB_SIZE_BYTES because we need to pad
        // the data on the right with 0s if it doesn't take up the full blob
        bytes memory totalBlobs = new bytes(BLOB_SIZE_BYTES * MAX_NUMBER_OF_BLOBS);

        assembly {
            // The pointer to the allocated memory above. We skip 32 bytes to avoid overwriting the length.
            let ptr := add(totalBlobs, 0x20)
            calldatacopy(ptr, _pubdata.offset, _pubdata.length)
        }

        for (uint256 i = 0; i < MAX_NUMBER_OF_BLOBS; ++i) {
            uint256 start = BLOB_SIZE_BYTES * i;

            // We break if the pubdata isn't enough to cover all 6 blobs. On L1 it is expected that the hash
            // will be bytes32(0) if a blob isn't going to be used.
            if (start >= _pubdata.length) {
                break;
            }

            bytes32 blobHash;
            assembly {
                // The pointer to the allocated memory above skipping the length.
                let ptr := add(totalBlobs, 0x20)
                blobHash := keccak256(add(ptr, start), BLOB_SIZE_BYTES)
            }

            blobHashes[i] = blobHash;
        }

        for (uint8 i = 0; i < MAX_NUMBER_OF_BLOBS; ++i) {
            SystemContractHelper.toL1(
                true,
                bytes32(uint256(SystemLogKey(i + uint256(SystemLogKey.BLOB_ONE_HASH_KEY)))),
                blobHashes[i]
            );
        }
    }
}
