// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

interface IChainModule {
  /* @notice this function provides the block number of current chain.
   * @dev certain chains have its own function to retrieve block number, e.g. Arbitrum
   * @return blockNumber the block number of the current chain.
   */
  function blockNumber() external view returns (uint256 blockNumber);

  /* @notice this function provides the block hash of a block number.
   * @dev this function can usually look back 256 blocks at most, unless otherwise specified
   * @param blockNumber the block number
   * @return blockHash the block hash of the input block number
   */
  function blockHash(uint256 blockNumber) external view returns (bytes32 blockHash);

  /* @notice this function provides the L1 fee of current transaction.
   * @dev retrieve the L1 data fee for a L2 transaction. it should return 0 for L1 chains. it should
   * return 0 for L2 chains if they don't have L1 fee component.
   * @param dataSize the calldata size of the current transaction
   * @return l1Fee the L1 fee in wei incurred by calldata of this data size
   */
  function getCurrentL1Fee(uint256 dataSize) external view returns (uint256 l1Fee);

  /* @notice this function provides the max possible L1 fee of current transaction.
   * @dev retrieve the max possible L1 data fee for a L2 transaction. it should return 0 for L1 chains. it should
   * return 0 for L2 chains if they don't have L1 fee component.
   * @param dataSize the calldata size of the current transaction
   * @return maxL1Fee the max possible L1 fee in wei incurred by calldata of this data size
   */
  function getMaxL1Fee(uint256 dataSize) external view returns (uint256 maxL1Fee);

  /* @notice this function provides the overheads of calling this chain module.
   * @return chainModuleFixedOverhead the fixed overhead incurred by calling this chain module
   * @return chainModulePerByteOverhead the fixed overhead per byte incurred by calling this chain module with calldata
   */
  function getGasOverhead()
    external
    view
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead);
}
