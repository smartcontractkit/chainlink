// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../../ChainSpecificUtil.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title BlockhashStore
 * @notice This contract provides a way to access blockhashes older than
 *   the 256 block limit imposed by the BLOCKHASH opcode.
 *   You may assume that any blockhash stored by the contract is correct.
 *   Note that the contract depends on the format of serialized Ethereum
 *   blocks. If a future hardfork of Ethereum changes that format, the
 *   logic in this contract may become incorrect and an updated version
 *   would have to be deployed.
 */
contract TrustedBlockhashStore is ConfirmedOwner {
    error NotInWhitelist();
    error InvalidTrustedBlockhashes();
    error InvalidRecentBlockhash();

    mapping(uint256 => bytes32) public s_blockhashes;
    mapping(uint256 => bool) public s_whitelistStatus;
    address[] public s_whitelist;

    constructor(address[] memory whitelist) ConfirmedOwner(msg.sender) {
        s_whitelist = whitelist;
    }

    /**
     * @notice sets the whitelist of addresses that can store blockhashes
     * @param whitelist the whitelist of addresses that can store blockhashes
     */
    function setWhitelist(address[] calldata whitelist) external onlyOwner {
        address[] memory previousWhitelist = s_whitelist;
        s_whitelist = whitelist;

        // Unset whitelist status for all addresses in the previous whitelist,
        // and set whitelist status for all addresses in the new whitelist.
        for (uint256 i = 0; i < previousWhitelist.length; i++) {
            s_whitelistStatus[previousWhitelist[i]] = false;
        }
        for (uint256 i = 0; i < whitelist.length; i++) {
            s_whitelistStatus[whitelist[i]] = true;
        }
    }

    /**
     * @notice stores blockhash of a given block, assuming it is available through BLOCKHASH
     * @param n the number of the block whose blockhash should be stored
     */

    function store(uint256 n) public {
        bytes32 h = ChainSpecificUtil.getBlockhash(uint64(n));
        require(h != 0x0, "blockhash(n) failed");
        s_blockhashes[n] = h;
    }

    /**
     * @notice stores a list of blockhashes and their respective blocks, only callable
     * by a whitelisted address
     * @param blockhashes the list of blockhashes and their respective blocks
     */
    function storeTrusted(
        uint256[] calldata blockNums,
        bytes32[] calldata blockhashes,
        uint256 recentBlockNumber,
        bytes32 recentBlockhash
    ) external {
        bytes32 onChainHash = ChainSpecificUtil.getBlockhash(recentBlockNumber);
        if (onChainHash != recentBlockhash) {
            revert InvalidRecentBlockhash();
        }

        if (!s_whitelistStatus[msg.sender]) {
            revert NotInWhitelist();
        }

        if (blockNums.length != blockhashes.length) {
            revert InvalidTrustedBlockhashes();
        }

        for (uint256 i = 0; i < blockNums.length; i++) {
            s_blockhashes[blockNums[i]] = blockhashes[i];
        }
    }

    /**
     * @notice gets a blockhash from the store. If no hash is known, this function reverts.
     * @param n the number of the block whose blockhash should be returned
     */
    function getBlockhash(uint256 n) external view returns (bytes32) {
        bytes32 h = s_blockhashes[n];
        require(h != 0x0, "blockhash not found in store");
        return h;
    }
}
