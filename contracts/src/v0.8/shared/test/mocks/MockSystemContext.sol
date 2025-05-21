// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ISystemContext} from "../../../vendor/@matter-labs/era-contracts/gas-bound-caller/contracts/ISystemContext.sol";

///
/// @notice A minimal mock for ISystemContext to satisfy all interface functions.
///         This can be deployed (not abstract) so you can reference it in tests.
///
contract MockSystemContext is ISystemContext {
  // ---------------------------------------
  // Storage variables for testing
  // ---------------------------------------
  uint256 private s_currentPubdataSpent;
  uint256 private s_gasPerPubdataByte = 10;

  // Example placeholders for block number & timestamp
  uint128 private s_mockBlockNumber = 1000;
  uint128 private s_mockBlockTimestamp = 123456789;

  // ---------------------------------------
  // Functions required by ISystemContext
  // ---------------------------------------

  function chainId() external view override returns (uint256) {
    // Return the current chain ID or a fixed mock
    return block.chainid;
  }

  function origin() external view override returns (address) {
    // Return the tx.origin or a mock address
    // solhint-disable-next-line avoid-tx-origin
    return tx.origin;
  }

  function gasPrice() external pure override returns (uint256) {
    // Return a dummy gas price
    return 1000000000; // 1 gwei, for example
  }

  function blockGasLimit() external pure override returns (uint256) {
    // Return a dummy block gas limit
    return 30_000_000;
  }

  function coinbase() external view override returns (address) {
    // Return the current block.coinbase or a mock
    return block.coinbase;
  }

  function difficulty() external view override returns (uint256) {
    // Return the block.difficulty or a fixed mock
    return block.prevrandao;
  }

  function baseFee() external view override returns (uint256) {
    // Return the current block.basefee or a mock
    return block.basefee;
  }

  function txNumberInBlock() external pure override returns (uint16) {
    // Return a dummy txNumberInBlock
    return 1;
  }

  function getBlockHashEVM(uint256 _block) external view override returns (bytes32) {
    // Return a dummy value (or actual blockhash if you prefer)
    return blockhash(_block);
  }

  function getBatchHash(uint256 _batchNumber) external pure override returns (bytes32) {
    // Return dummy
    return keccak256(abi.encodePacked("BatchHashMock", _batchNumber));
  }

  function getBlockNumber() external view override returns (uint128) {
    // Return your stored mock or real block.number cast to uint128
    return s_mockBlockNumber;
  }

  function getBlockTimestamp() external view override returns (uint128) {
    // Return your stored mock or real block.timestamp cast to uint128
    return s_mockBlockTimestamp;
  }

  function getBatchNumberAndTimestamp() external view override returns (uint128 blockNumber, uint128 blockTimestamp) {
    // Return dummy or relevant block info
    return (s_mockBlockNumber, s_mockBlockTimestamp);
  }

  function getL2BlockNumberAndTimestamp() external view override returns (uint128 blockNumber, uint128 blockTimestamp) {
    // Return dummy or relevant block info
    return (s_mockBlockNumber, s_mockBlockTimestamp);
  }

  function gasPerPubdataByte() external pure override returns (uint256 gasPerPubdataByte_) {
    return gasPerPubdataByte_;
  }

  function getCurrentPubdataSpent() external pure override returns (uint256 currentPubdataSpent) {
    return currentPubdataSpent;
  }

  // ---------------------------------------
  // Extra helpers for testing
  // ---------------------------------------

  /// @notice Lets you set the mock pubdata spent for testing
  function setCurrentPubdataSpent(uint256 newVal) external {
    s_currentPubdataSpent = newVal;
  }

  /// @notice Lets you set the mock gas per pubdata byte for testing
  function setGasPerPubdataByte(uint256 newVal) external {
    s_gasPerPubdataByte = newVal;
  }

  /// @notice Lets you set the mock block number
  function setMockBlockNumber(uint128 newVal) external {
    s_mockBlockNumber = newVal;
  }

  /// @notice Lets you set the mock block timestamp
  function setMockBlockTimestamp(uint128 newVal) external {
    s_mockBlockTimestamp = newVal;
  }
}
