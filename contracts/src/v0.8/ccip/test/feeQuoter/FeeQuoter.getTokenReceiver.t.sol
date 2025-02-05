// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../../ccip/FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getTokenReceiver is FeeQuoterFeeSetup {
  FeeQuoter.DestChainConfig private s_svmDestChainConfig;

  function setUp() public override {
    super.setUp();
    s_svmDestChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_svmDestChainConfig.enforceOutOfOrder = true; // Enforcing out of order execution for messages to SVM
    s_svmDestChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;
  }

  function test_getTokenReceiver_EVM() public {
    uint256 tokenAmount = 10000e18;
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    bytes memory resolvedTokenReceiver = s_feeQuoter.getTokenReceiver(DEST_CHAIN_SELECTOR, message);
    assertEq(resolvedTokenReceiver, abi.encode(OWNER));
  }

  function test_getTokenReceiver_SVM() public {
    // register SVM chain family excludeSelector
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: s_svmDestChainConfig});
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigs);

    // prepare message
    uint256 tokenAmount = 10000e18;
    bytes32 svmTokenReceiver = bytes32("TOKEN RECEIVER");
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    message.extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: GAS_LIMIT,
        accountIsWritableBitmap: 0,
        allowOutOfOrderExecution: true,
        tokenReceiver: svmTokenReceiver,
        accounts: new bytes32[](0)
      })
    );
    bytes memory resolvedTokenReceiver = s_feeQuoter.getTokenReceiver(DEST_CHAIN_SELECTOR, message);
    assertEq(resolvedTokenReceiver, abi.encode(svmTokenReceiver));
  }

  // Reverts
  function test_getTokenReceiver_RevertWhen_InvalidChainFamilySelector() public {
    uint256 tokenAmount = 10000e18;
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidChainFamilySelector.selector, bytes4(0)));
    s_feeQuoter.getTokenReceiver(uint64(vm.randomUint()), message);
  }
}
