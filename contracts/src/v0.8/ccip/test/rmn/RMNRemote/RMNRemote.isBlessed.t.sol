// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMN} from "../../../interfaces/IRMN.sol";

import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_isBlessed is RMNRemoteSetup {
  function test_isBlessed() public {
    IRMN.TaggedRoot memory taggedRoot = IRMN.TaggedRoot({root: keccak256("root"), commitStore: makeAddr("commitStore")});

    vm.mockCall(
      address(s_legacyRMN), abi.encodeWithSelector(s_legacyRMN.isBlessed.selector, taggedRoot), abi.encode(true)
    );

    assertTrue(s_rmnRemote.isBlessed(taggedRoot));

    vm.mockCall(
      address(s_legacyRMN), abi.encodeWithSelector(s_legacyRMN.isBlessed.selector, taggedRoot), abi.encode(false)
    );

    assertFalse(s_rmnRemote.isBlessed(taggedRoot));
  }

  function test_RevertWhen_isBlessedWhen_IsBlessedNotAvailable() public {
    IRMN.TaggedRoot memory taggedRoot = IRMN.TaggedRoot({root: keccak256("root"), commitStore: makeAddr("commitStore")});

    s_rmnRemote = new RMNRemote(100, IRMN(address(0)));

    vm.expectRevert(RMNRemote.IsBlessedNotAvailable.selector);
    s_rmnRemote.isBlessed(taggedRoot);
  }
}
