// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsCoordinator} from "../../../dev/v1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../../dev/v1_0_0/FunctionsBilling.sol";

contract FunctionsCoordinatorTestHelper is FunctionsCoordinator {
  constructor(
    address router,
    FunctionsBilling.Config memory config,
    address linkToNativeFeed
  ) FunctionsCoordinator(router, config, linkToNativeFeed) {}

  function callReport(bytes calldata report) external {
    address[MAX_NUM_ORACLES] memory signers;
    signers[0] = msg.sender;
    _report(gasleft(), msg.sender, 1, signers, report);
  }

  function callReportMultipleSigners(bytes calldata report, address secondSigner) external {
    address[MAX_NUM_ORACLES] memory signers;
    signers[0] = msg.sender;
    signers[1] = secondSigner;
    _report(gasleft(), msg.sender, 2, signers, report);
  }

  function callReportWithSigners(bytes calldata report, address[MAX_NUM_ORACLES] memory signers) external {
    _report(gasleft(), msg.sender, 2, signers, report);
  }
}
