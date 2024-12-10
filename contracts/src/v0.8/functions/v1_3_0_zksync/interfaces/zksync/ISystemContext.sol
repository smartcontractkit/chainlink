// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

address constant SYSTEM_CONTEXT = address(0x000000000000000000000000000000000000800B);

interface ISystemContext {
  function gasPrice() external view returns (uint256);
  function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
}
