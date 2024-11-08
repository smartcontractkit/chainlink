// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMNRemote} from "../../../interfaces/IRMNRemote.sol";
import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_verify_withConfigSet is RMNRemoteSetup {
  function setUp() public override {
    super.setUp();
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, f: 3});
    s_rmnRemote.setConfig(config);
    _generatePayloadAndSigs(2, 4);
  }

  function test_verify_success() public view {
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_InvalidSignature_reverts() public {
    IRMNRemote.Signature memory sig = s_signatures[s_signatures.length - 1];
    sig.r = _randomBytes32();
    s_signatures.pop();
    s_signatures.push(sig);

    vm.expectRevert(RMNRemote.InvalidSignature.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_OutOfOrderSignatures_not_sorted_reverts() public {
    IRMNRemote.Signature memory sig1 = s_signatures[s_signatures.length - 1];
    s_signatures.pop();
    IRMNRemote.Signature memory sig2 = s_signatures[s_signatures.length - 1];
    s_signatures.pop();
    s_signatures.push(sig1);
    s_signatures.push(sig2);

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_OutOfOrderSignatures_duplicateSignature_reverts() public {
    IRMNRemote.Signature memory sig = s_signatures[s_signatures.length - 2];
    s_signatures.pop();
    s_signatures.push(sig);

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_UnexpectedSigner_reverts() public {
    _setupSigners(4); // create new signers that aren't configured on RMNRemote
    _generatePayloadAndSigs(2, 4);

    vm.expectRevert(RMNRemote.UnexpectedSigner.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }

  function test_verify_ThresholdNotMet_reverts() public {
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, f: 2}); // 3 = f+1 sigs required
    s_rmnRemote.setConfig(config);

    _generatePayloadAndSigs(2, 2); // 2 sigs generated, but 3 required

    vm.expectRevert(RMNRemote.ThresholdNotMet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures);
  }
}
