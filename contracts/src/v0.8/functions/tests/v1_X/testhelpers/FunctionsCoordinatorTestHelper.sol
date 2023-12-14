// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsBillingConfig} from "../../../dev/v1_X/interfaces/IFunctionsBilling.sol";

contract FunctionsCoordinatorTestHelper is FunctionsCoordinator {
  constructor(
    address router,
    FunctionsBillingConfig memory config,
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
