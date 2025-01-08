// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {GLOBAL_CURSE_SUBJECT} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_global_curses is RMNRemoteSetup {
  function test_isCursed_globalCurseSubject() public {
    bytes16 randSubject = bytes16(keccak256("random subject"));
    assertFalse(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject));

    s_rmnRemote.curse(GLOBAL_CURSE_SUBJECT);
    assertTrue(s_rmnRemote.isCursed());
    assertTrue(s_rmnRemote.isCursed(randSubject));

    s_rmnRemote.uncurse(GLOBAL_CURSE_SUBJECT);
    assertFalse(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject));
  }
}
