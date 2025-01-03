// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_uncurse is RMNRemoteSetup {
  function setUp() public override {
    super.setUp();
    s_rmnRemote.curse(s_curseSubjects);
  }

  function test_uncurse() public {
    vm.expectEmit();
    emit RMNRemote.Uncursed(s_curseSubjects);

    s_rmnRemote.uncurse(s_curseSubjects);

    assertEq(s_rmnRemote.getCursedSubjects().length, 0);
    assertFalse(s_rmnRemote.isCursed(CURSE_SUBJ_1));
    assertFalse(s_rmnRemote.isCursed(CURSE_SUBJ_2));
  }

  function test_RevertWhen_uncurse_NotCursed_duplicatedUncurseSubject() public {
    s_curseSubjects.push(CURSE_SUBJ_1);

    vm.expectRevert(abi.encodeWithSelector(RMNRemote.NotCursed.selector, CURSE_SUBJ_1));
    s_rmnRemote.uncurse(s_curseSubjects);
  }

  function test_RevertWhen_uncurse_calledByNonOwner() public {
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    vm.stopPrank();
    vm.prank(STRANGER);
    s_rmnRemote.uncurse(s_curseSubjects);
  }
}
