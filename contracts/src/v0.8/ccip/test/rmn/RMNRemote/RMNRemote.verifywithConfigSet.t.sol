// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_verify_withConfigSet is RMNRemoteSetup {
  function setUp() public override {
    super.setUp();

    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 3});
    s_rmnRemote.setConfig(config);
    _generatePayloadAndSigs(2, 4);
  }

  function test_verify() public view {
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_InvalidSignature() public {
    s_signatures[s_signatures.length - 1].r = 0x0;

    vm.expectRevert(RMNRemote.InvalidSignature.selector);

    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_OutOfOrderSignatures_not_sorted() public {
    IRMNRemote.Signature memory sig1 = s_signatures[s_signatures.length - 1];
    IRMNRemote.Signature memory sig2 = s_signatures[s_signatures.length - 2];

    s_signatures[s_signatures.length - 1] = sig2;
    s_signatures[s_signatures.length - 2] = sig1;

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_OutOfOrderSignatures_duplicateSignature() public {
    s_signatures[s_signatures.length - 1] = s_signatures[s_signatures.length - 2];

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_UnexpectedSigner() public {
    _setupSigners(4); // create new signers that aren't configured on RMNRemote
    _generatePayloadAndSigs(2, 4);

    vm.expectRevert(RMNRemote.UnexpectedSigner.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_RevertWhen_ThresholdNotMet() public {
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 2}); // 3 = f+1 sigs required
    s_rmnRemote.setConfig(config);

    _generatePayloadAndSigs(2, 2); // 2 sigs generated, but 3 required

    vm.expectRevert(RMNRemote.ThresholdNotMet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }
}
