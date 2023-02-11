// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MockTimelock {
  struct Transaction {
    address target;
    uint256 value;
    string signature;
    bytes data;
    uint256 eta;
  }
  Transaction s_queuedTransaction;
  Transaction s_cancelledTransaction;
  Transaction s_executedTransaction;

  function getQueuedTransaction() external view returns (Transaction memory) {
    return s_queuedTransaction;
  }

  function getCancelledTransaction() external view returns (Transaction memory) {
    return s_cancelledTransaction;
  }

  function getExecutedTransaction() external view returns (Transaction memory) {
    return s_executedTransaction;
  }

  function delay() external pure returns (uint256) {
    return 0;
  }

  function GRACE_PERIOD() external pure returns (uint256) {
    return type(uint256).max - 1;
  }

  function queueTransaction(
    address target,
    uint256 value,
    string calldata signature,
    bytes calldata data,
    uint256 eta
  ) external returns (bytes32) {
    s_queuedTransaction = Transaction({target: target, value: value, signature: signature, data: data, eta: eta});
    return "";
  }

  function queuedTransactions(bytes32 hash) external view returns (bool) {
    return false;
  }

  function cancelTransaction(
    address target,
    uint256 value,
    string calldata signature,
    bytes calldata data,
    uint256 eta
  ) external {
    s_cancelledTransaction = Transaction({target: target, value: value, signature: signature, data: data, eta: eta});
  }

  function executeTransaction(
    address target,
    uint256 value,
    string calldata signature,
    bytes calldata data,
    uint256 eta
  ) external payable returns (bytes memory) {
    s_executedTransaction = Transaction({target: target, value: value, signature: signature, data: data, eta: eta});
    return "";
  }
}
