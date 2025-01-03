// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../shared/access/Ownable2Step.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_applyPremiumMultiplierWeiPerEthUpdates is FeeQuoterSetup {
  function testFuzz_applyPremiumMultiplierWeiPerEthUpdates_Success(
    FeeQuoter.PremiumMultiplierWeiPerEthArgs memory premiumMultiplierWeiPerEthArg
  ) public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = premiumMultiplierWeiPerEthArg;

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      premiumMultiplierWeiPerEthArg.token, premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(premiumMultiplierWeiPerEthArg.token)
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesSingleToken() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = s_feeQuoterPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      s_feeQuoterPremiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesMultipleTokens() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](2);
    premiumMultiplierWeiPerEthArgs[0] = s_feeQuoterPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);
    premiumMultiplierWeiPerEthArgs[1].token = vm.addr(2);

    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );
    vm.expectEmit();
    emit FeeQuoter.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(2), premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth
    );

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
    assertEq(
      premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth,
      s_feeQuoter.getPremiumMultiplierWeiPerEth(vm.addr(2))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesZeroInput() public {
    vm.recordLogs();
    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_RevertWhen_OnlyCallableByOwnerOrAdmin() public {
    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs;
    vm.startPrank(STRANGER);

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
  }
}
