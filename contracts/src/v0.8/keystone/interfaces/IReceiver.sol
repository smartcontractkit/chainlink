// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IReceiver - receives keystone reports
interface IReceiver {
  function onReport(bytes calldata metadata, bytes calldata report) external;
}
