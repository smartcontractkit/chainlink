// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsBillingConfig} from "../../../dev/v1_X/interfaces/IFunctionsBilling.sol";

contract FunctionsCoordinatorTestHelper is FunctionsCoordinator {
  constructor(
    address router,
    FunctionsBillingConfig memory config,
    address linkToNativeFeed,
    address linkToUsdFeed
  ) FunctionsCoordinator(router, config, linkToNativeFeed, linkToUsdFeed) {}

  function callReport(bytes calldata report) external {
    address[MAX_NUM_ORACLES] memory signers;
    signers[0] = msg.sender;
    (
      bytes32[] memory requestIds,
      bytes[] memory results,
      bytes[] memory errors,
      bytes[] memory onchainMetadata,
      bytes[] memory offchainMetadata
    ) = abi.decode(report, (bytes32[], bytes[], bytes[], bytes[], bytes[]));
    _report(
      DecodedReport({
        requestIds: requestIds,
        results: results,
        errors: errors,
        onchainMetadata: onchainMetadata,
        offchainMetadata: offchainMetadata
      })
    );
  }

  function callReportMultipleSigners(bytes calldata report, address secondSigner) external {
    address[MAX_NUM_ORACLES] memory signers;
    signers[0] = msg.sender;
    signers[1] = secondSigner;
    (
      bytes32[] memory requestIds,
      bytes[] memory results,
      bytes[] memory errors,
      bytes[] memory onchainMetadata,
      bytes[] memory offchainMetadata
    ) = abi.decode(report, (bytes32[], bytes[], bytes[], bytes[], bytes[]));
    _report(
      DecodedReport({
        requestIds: requestIds,
        results: results,
        errors: errors,
        onchainMetadata: onchainMetadata,
        offchainMetadata: offchainMetadata
      })
    );
  }
}
