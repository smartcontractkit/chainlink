// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IFeeManager} from "../../interfaces/IFeeManager.sol";

contract FeeManagerProxy {
  IFeeManager internal s_feeManager;

  function processFee(bytes calldata payload, bytes calldata parameterPayload) public payable {
    s_feeManager.processFee{value: msg.value}(payload, parameterPayload, msg.sender);
  }

  function processFeeBulk(bytes[] calldata payloads, bytes calldata parameterPayload) public payable {
    s_feeManager.processFeeBulk{value: msg.value}(payloads, parameterPayload, msg.sender);
  }

  function setFeeManager(IFeeManager feeManager) public {
    s_feeManager = feeManager;
  }
}
