// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_applyFeeTokensUpdates is FeeQuoterSetup {
  function test_ApplyFeeTokensUpdates() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];

    vm.expectEmit();
    emit FeeQuoter.FeeTokenAdded(feeTokens[0]);

    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_feeQuoter.getFeeTokens().length, 3);
    assertEq(s_feeQuoter.getFeeTokens()[2], feeTokens[0]);

    // add same feeToken is no-op
    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_feeQuoter.getFeeTokens().length, 3);
    assertEq(s_feeQuoter.getFeeTokens()[2], feeTokens[0]);

    vm.expectEmit();
    emit FeeQuoter.FeeTokenRemoved(feeTokens[0]);

    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_feeQuoter.getFeeTokens().length, 2);

    // removing already removed feeToken is no-op and does not emit an event
    vm.recordLogs();

    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_feeQuoter.getFeeTokens().length, 2);

    vm.assertEq(vm.getRecordedLogs().length, 0);

    // Removing and adding the same fee token is allowed and emits both events
    // Add it first
    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);

    vm.expectEmit();
    emit FeeQuoter.FeeTokenRemoved(feeTokens[0]);
    vm.expectEmit();
    emit FeeQuoter.FeeTokenAdded(feeTokens[0]);

    s_feeQuoter.applyFeeTokensUpdates(feeTokens, feeTokens);
  }

  function test_RevertWhen_OnlyCallableByOwner() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_feeQuoter.applyFeeTokensUpdates(new address[](0), new address[](0));
  }
}
