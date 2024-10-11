// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {OffRamp} from "../../offRamp/OffRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMOffRampHelper {
  uint64 public s_nonce;
  mapping(address sender => uint64 nonce) public s_nonces;

  function execute(Internal.ExecutionReport memory report, OffRamp.GasLimitOverride[] memory) external {
    for (uint256 i; i < report.messages.length; i++) {
      Internal.EVM2EVMMessage memory message = report.messages[i];
      s_nonces[message.sender]++;
    }
  }

  function metadataHash() external pure returns (bytes32) {
    return 0x0;
  }

  function getSenderNonce(address sender) external view returns (uint64 nonce) {
    return s_nonces[sender];
  }
}
