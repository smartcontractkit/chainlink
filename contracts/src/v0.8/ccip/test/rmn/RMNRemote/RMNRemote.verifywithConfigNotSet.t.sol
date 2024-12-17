// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_verify_withConfigNotSet is RMNRemoteSetup {
  function test_RevertWhen_verifys() public {
    Internal.MerkleRoot[] memory merkleRoots = new Internal.MerkleRoot[](0);
    IRMNRemote.Signature[] memory signatures = new IRMNRemote.Signature[](0);

    vm.expectRevert(RMNRemote.ConfigNotSet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, merkleRoots, signatures);
  }
}
