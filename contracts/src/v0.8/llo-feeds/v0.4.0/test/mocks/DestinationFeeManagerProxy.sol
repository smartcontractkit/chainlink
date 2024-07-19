// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "../../interfaces/IDestinationFeeManager.sol";

contract DestinationFeeManagerProxy {
  IDestinationFeeManager internal i_feeManager;

  function processFee(bytes32 poolId, bytes calldata payload, bytes calldata parameterPayload) public payable {
    i_feeManager.processFee{value: msg.value}(poolId, payload, parameterPayload, msg.sender);
  }

  function processFeeBulk(bytes32[] memory poolIds, bytes[] calldata payloads, bytes calldata parameterPayload) public payable {
    i_feeManager.processFeeBulk{value: msg.value}(poolIds, payloads, parameterPayload, msg.sender);
  }

  function setDestinationFeeManager(IDestinationFeeManager feeManager) public {
    i_feeManager = feeManager;
  }
}
