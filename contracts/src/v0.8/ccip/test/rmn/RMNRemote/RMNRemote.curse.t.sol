// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_curse is RMNRemoteSetup {
  function test_curse() public {
    vm.expectEmit();
    emit RMNRemote.Cursed(s_curseSubjects);

    s_rmnRemote.curse(s_curseSubjects);

    assertEq(abi.encode(s_rmnRemote.getCursedSubjects()), abi.encode(s_curseSubjects));
    assertTrue(s_rmnRemote.isCursed(CURSE_SUBJ_1));
    assertTrue(s_rmnRemote.isCursed(CURSE_SUBJ_2));
    // Should not have cursed a random subject
    assertFalse(s_rmnRemote.isCursed(bytes16(keccak256("subject 3"))));
  }

  function test_RevertWhen_curse_AlreadyCursed_duplicateSubject() public {
    s_curseSubjects.push(CURSE_SUBJ_1);

    vm.expectRevert(abi.encodeWithSelector(RMNRemote.AlreadyCursed.selector, CURSE_SUBJ_1));
    s_rmnRemote.curse(s_curseSubjects);
  }

  function test_RevertWhen_curse_calledByNonOwner() public {
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    vm.stopPrank();
    vm.prank(STRANGER);
    s_rmnRemote.curse(s_curseSubjects);
  }
}
