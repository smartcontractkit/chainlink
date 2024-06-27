// SPDX-License-Identifier: MIT
// solhint-disable-next-line one-contract-per-file
pragma solidity 0.8.19;

import {ChainSpecificUtil} from "../ChainSpecificUtil.sol";

/**
 * @title BatchBlockhashStore
 * @notice The BatchBlockhashStore contract acts as a proxy to write many blockhashes to the
 *   provided BlockhashStore contract efficiently in a single transaction. This results
 *   in plenty of gas savings and higher throughput of blockhash storage, which is desirable
 *   in times of high network congestion.
 */
contract BatchBlockhashStore {
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
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
      // skip the block if it's not storeable, the caller will have to check
      // after the transaction is mined to see if the blockhash was truly stored.
      if (!_storeableBlock(blockNumbers[i])) {
        continue;
      }
      BHS.store(blockNumbers[i]);
    }
  }

  /**
   * @notice stores blockhashes after verifying blockheader of child/subsequent block
   * @param blockNumbers the block numbers whose blockhashes should be stored, in decreasing order
   * @param headers the rlp-encoded block headers of blockNumbers[i] + 1.
   */
  function storeVerifyHeader(uint256[] memory blockNumbers, bytes[] memory headers) public {
    // solhint-disable-next-line gas-custom-errors
    require(blockNumbers.length == headers.length, "input array arg lengths mismatch");
    for (uint256 i = 0; i < blockNumbers.length; i++) {
      BHS.storeVerifyHeader(blockNumbers[i], headers[i]);
    }
  }

  /**
   * @notice retrieves blockhashes of all the given block numbers from the blockhash store, if available.
   * @param blockNumbers array of block numbers to fetch blockhashes for
   * @return blockhashes array of block hashes corresponding to each block number provided in the `blockNumbers`
   *   param. If the blockhash is not found, 0x0 is returned instead of the real blockhash, indicating
   *   that it is not in the blockhash store.
   */
  function getBlockhashes(uint256[] memory blockNumbers) external view returns (bytes32[] memory) {
    bytes32[] memory blockHashes = new bytes32[](blockNumbers.length);
    for (uint256 i = 0; i < blockNumbers.length; i++) {
      try BHS.getBlockhash(blockNumbers[i]) returns (bytes32 bh) {
        blockHashes[i] = bh;
      } catch Error(string memory /* reason */) {
        blockHashes[i] = 0x0;
      }
    }
    return blockHashes;
  }

  /**
   * @notice returns true if and only if the given block number's blockhash can be retrieved
   *   using the blockhash() instruction.
   * @param blockNumber the block number to check if it's storeable with blockhash()
   */
  function _storeableBlock(uint256 blockNumber) private view returns (bool) {
    // handle edge case on simulated chains which possibly have < 256 blocks total.
    return
      ChainSpecificUtil._getBlockNumber() <= 256 ? true : blockNumber >= (ChainSpecificUtil._getBlockNumber() - 256);
  }
}

// solhint-disable-next-line interface-starts-with-i
interface BlockhashStore {
  function storeVerifyHeader(uint256 n, bytes memory header) external;

  function store(uint256 n) external;

  function getBlockhash(uint256 n) external view returns (bytes32);
}
