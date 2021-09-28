// NOTICE: pragma change from original (>=0.4.21 <0.7.0)
pragma solidity >=0.4.21 <0.9.0;

/**
 * @title Precompiled contract that exists in every Arbitrum chain at address(100), 0x0000000000000000000000000000000000000064. Exposes a variety of system-level functionality.
 */
interface ArbSys {
  /**
   * @notice Get internal version number identifying an ArbOS build
   * @return version number as int
   */
  function arbOSVersion() external pure returns (uint256);

  function arbChainID() external view returns (uint256);

  /**
   * @notice Get Arbitrum block number (distinct from L1 block number; Arbitrum genesis block has block number 0)
   * @return block number as int
   */
  function arbBlockNumber() external view returns (uint256);

  /**
   * @notice Send given amount of Eth to dest from sender.
   * This is a convenience function, which is equivalent to calling sendTxToL1 with empty calldataForL1.
   * @param destination recipient address on L1
   * @return unique identifier for this L2-to-L1 transaction.
   */
  function withdrawEth(address destination) external payable returns (uint256);

  /**
   * @notice Send a transaction to L1
   * @param destination recipient address on L1
   * @param calldataForL1 (optional) calldata for L1 contract call
   * @return a unique identifier for this L2-to-L1 transaction.
   */
  function sendTxToL1(address destination, bytes calldata calldataForL1) external payable returns (uint256);

  /**
   * @notice get the number of transactions issued by the given external account or the account sequence number of the given contract
   * @param account target account
   * @return the number of transactions issued by the given external account or the account sequence number of the given contract
   */
  function getTransactionCount(address account) external view returns (uint256);

  /**
   * @notice get the value of target L2 storage slot
   * This function is only callable from address 0 to prevent contracts from being able to call it
   * @param account target account
   * @param index target index of storage slot
   * @return stotage value for the given account at the given index
   */
  function getStorageAt(address account, uint256 index) external view returns (uint256);

  /**
   * @notice check if current call is coming from l1
   * @return true if the caller of this was called directly from L1
   */
  function isTopLevelCall() external view returns (bool);

  event L2ToL1Transaction(
    address caller,
    address indexed destination,
    uint256 indexed uniqueId,
    uint256 indexed batchNumber,
    uint256 indexInBatch,
    uint256 arbBlockNum,
    uint256 ethBlockNum,
    uint256 timestamp,
    uint256 callvalue,
    bytes data
  );
}
