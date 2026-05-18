// SPDX-License-Identifier: MIT

pragma solidity ^0.8.6;

contract KeeperRegistryCheckUpkeepGasUsageWrapper1_2Mock {
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);

  function emitOwnershipTransferRequested(address from, address to) public {
    emit OwnershipTransferRequested(from, to);
  }

  function emitOwnershipTransferred(address from, address to) public {
    emit OwnershipTransferred(from, to);
  }

  bool public s_mockResult;
  bytes public s_mockPayload;
  uint256 public s_mockGas;

  // Function to set mock return data for the measureCheckGas function
  function setMeasureCheckGasResult(bool result, bytes memory payload, uint256 gas) external {
    s_mockResult = result;
    s_mockPayload = payload;
    s_mockGas = gas;
  }

  // Mock measureCheckGas function
  function measureCheckGas(uint256 id, address from) external returns (bool, bytes memory, uint256) {
    return (s_mockResult, s_mockPayload, s_mockGas);
  }
}
