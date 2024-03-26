// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ChainSpecificUtil} from "../../ChainSpecificUtil.sol";

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
contract BlockhashStore {
  mapping(uint256 => bytes32) internal s_blockhashes;

  /**
   * @notice stores blockhash of a given block, assuming it is available through BLOCKHASH
   * @param n the number of the block whose blockhash should be stored
   */
  function store(uint256 n) public {
    bytes32 h = ChainSpecificUtil._getBlockhash(uint64(n));
    // solhint-disable-next-line gas-custom-errors
    require(h != 0x0, "blockhash(n) failed");
    s_blockhashes[n] = h;
  }

  /**
   * @notice stores blockhash of the earliest block still available through BLOCKHASH.
   */
  function storeEarliest() external {
    store(ChainSpecificUtil._getBlockNumber() - 256);
  }

  /**
   * @notice stores blockhash after verifying blockheader of child/subsequent block
   * @param n the number of the block whose blockhash should be stored
   * @param header the rlp-encoded blockheader of block n+1. We verify its correctness by checking
   *   that it hashes to a stored blockhash, and then extract parentHash to get the n-th blockhash.
   */
  function storeVerifyHeader(uint256 n, bytes memory header) public {
    // solhint-disable-next-line gas-custom-errors
    require(keccak256(header) == s_blockhashes[n + 1], "header has unknown blockhash");

    // At this point, we know that header is the correct blockheader for block n+1.

    // The header is an rlp-encoded list. The head item of that list is the 32-byte blockhash of the parent block.
    // Based on how rlp works, we know that blockheaders always have the following form:
    // 0xf9____a0PARENTHASH...
    //   ^ ^   ^
    //   | |   |
    //   | |   +--- PARENTHASH is 32 bytes. rlpenc(PARENTHASH) is 0xa || PARENTHASH.
    //   | |
    //   | +--- 2 bytes containing the sum of the lengths of the encoded list items
    //   |
    //   +--- 0xf9 because we have a list and (sum of lengths of encoded list items) fits exactly into two bytes.
    //
    // As a consequence, the PARENTHASH is always at offset 4 of the rlp-encoded block header.

    bytes32 parentHash;
    assembly {
      parentHash := mload(add(header, 36)) // 36 = 32 byte offset for length prefix of ABI-encoded array
      //    +  4 byte offset of PARENTHASH (see above)
    }

    s_blockhashes[n] = parentHash;
  }

  /**
   * @notice gets a blockhash from the store. If no hash is known, this function reverts.
   * @param n the number of the block whose blockhash should be returned
   */
  function getBlockhash(uint256 n) external view returns (bytes32) {
    bytes32 h = s_blockhashes[n];
    // solhint-disable-next-line gas-custom-errors
    require(h != 0x0, "blockhash not found in store");
    return h;
  }
}
