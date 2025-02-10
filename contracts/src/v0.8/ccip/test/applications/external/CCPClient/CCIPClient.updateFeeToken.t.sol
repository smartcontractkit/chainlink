// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPClient} from "../../../../applications/external/CCIPClient.sol";
import {CCIPClientSetup} from "./CCIPClientSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPClient_updateFeeToken is CCIPClientSetup {
  function test_updateFeeToken_WETH() public {
    // WETH is used as a placeholder for any ERC20 token
    address WETH = s_sourceRouter.getWrappedNative();

    vm.expectEmit();
    emit IERC20.Approval(address(s_sender), address(s_sourceRouter), 0);

    vm.expectEmit();
    emit CCIPClient.FeeTokenUpdated(s_sourceFeeToken, WETH);

    s_sender.updateFeeToken(WETH);

    IERC20 newFeeToken = IERC20(s_sender.getFeeToken());
    assertEq(address(newFeeToken), WETH);
    assertEq(newFeeToken.allowance(address(s_sender), address(s_sourceRouter)), type(uint256).max);
    assertEq(IERC20(s_sourceFeeToken).allowance(address(s_sender), address(s_sourceRouter)), 0);
  }

  function test_updateFeeToken_Native() public {
    vm.expectEmit();
    emit CCIPClient.FeeTokenUpdated(s_sourceFeeToken, address(0));

    s_sender.updateFeeToken(address(0));

    IERC20 newFeeToken = IERC20(s_sender.getFeeToken());
    assertEq(address(newFeeToken), address(0));
  }
}
