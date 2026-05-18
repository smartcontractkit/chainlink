// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

ISystemContext constant SYSTEM_CONTEXT_CONTRACT = ISystemContext(address(0x800b));

interface ISystemContext {
  function gasPrice() external view returns (uint256);
  function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
  function getCurrentPubdataSpent() external view returns (uint256 currentPubdataSpent);
}
