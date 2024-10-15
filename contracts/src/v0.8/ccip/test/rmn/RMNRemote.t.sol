// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMNRemote} from "../../interfaces/IRMNRemote.sol";
import {Internal} from "../../libraries/Internal.sol";
import {GLOBAL_CURSE_SUBJECT, LEGACY_CURSE_SUBJECT, RMNRemote} from "../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_constructor is RMNRemoteSetup {
  function test_constructor_success() public view {
    assertEq(s_rmnRemote.getLocalChainSelector(), 1);
  }

  function test_constructor_zeroChainSelector_reverts() public {
    vm.expectRevert(RMNRemote.ZeroValueNotAllowed.selector);
    new RMNRemote(0);
  }
}

contract RMNRemote_setConfig is RMNRemoteSetup {
  function test_setConfig_minSignersIs0_success() public {
    // Initially there is no config, the version is 0
    uint32 currentConfigVersion = 0;
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    (uint32 version, RMNRemote.Config memory gotConfig) = s_rmnRemote.getVersionedConfig();
    assertEq(gotConfig.minSigners, 0);
    assertEq(version, currentConfigVersion);

    // A new config should increment the version
    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);
  }

  function test_setConfig_addSigner_removeSigner_success() public {
    uint32 currentConfigVersion = 0;
    uint256 numSigners = s_signers.length;
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    // add a signer
    address newSigner = makeAddr("new signer");
    s_signers.push(RMNRemote.Signer({onchainPublicKey: newSigner, nodeIndex: uint64(numSigners)}));
    config = RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    (uint32 version, RMNRemote.Config memory gotConfig) = s_rmnRemote.getVersionedConfig();
    assertEq(gotConfig.signers.length, s_signers.length);
    assertEq(gotConfig.signers[numSigners].onchainPublicKey, newSigner);
    assertEq(gotConfig.signers[numSigners].nodeIndex, uint64(numSigners));
    assertEq(version, currentConfigVersion);

    // remove two signers
    s_signers.pop();
    s_signers.pop();
    config = RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    (version, gotConfig) = s_rmnRemote.getVersionedConfig();
    assertEq(gotConfig.signers.length, s_signers.length);
    assertEq(version, currentConfigVersion);
  }

  function test_setConfig_invalidSignerOrder_reverts() public {
    s_signers.push(RMNRemote.Signer({onchainPublicKey: address(4), nodeIndex: 0}));
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectRevert(RMNRemote.InvalidSignerOrder.selector);
    s_rmnRemote.setConfig(config);
  }

  function test_setConfig_minSignersTooHigh_reverts() public {
    RMNRemote.Config memory config = RMNRemote.Config({
      rmnHomeContractConfigDigest: _randomBytes32(),
      signers: s_signers,
      minSigners: uint64(s_signers.length + 1)
    });

    vm.expectRevert(RMNRemote.MinSignersTooHigh.selector);
    s_rmnRemote.setConfig(config);
  }

  function test_setConfig_duplicateOnChainPublicKey_reverts() public {
    s_signers.push(RMNRemote.Signer({onchainPublicKey: s_signerWallets[0].addr, nodeIndex: uint64(s_signers.length)}));
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0});

    vm.expectRevert(RMNRemote.DuplicateOnchainPublicKey.selector);
    s_rmnRemote.setConfig(config);
  }
}

contract RMNRemote_verify_withConfigNotSet is RMNRemoteSetup {
  function test_verify_reverts() public {
    Internal.MerkleRoot[] memory merkleRoots = new Internal.MerkleRoot[](0);
    IRMNRemote.Signature[] memory signatures = new IRMNRemote.Signature[](0);

    vm.expectRevert(RMNRemote.ConfigNotSet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, merkleRoots, signatures, 0);
  }
}

contract RMNRemote_verify_withConfigSet is RMNRemoteSetup {
  function setUp() public override {
    super.setUp();
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 2});
    s_rmnRemote.setConfig(config);
    _generatePayloadAndSigs(2, 2);
  }

  function test_verify_success() public view {
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }

  function test_verify_minSignersIsZero_success() public {
    vm.stopPrank();
    vm.prank(OWNER);
    s_rmnRemote.setConfig(
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, minSigners: 0})
    );

    vm.stopPrank();
    vm.prank(OFF_RAMP_ADDRESS);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, new IRMNRemote.Signature[](0), s_v);
  }

  function test_verify_InvalidSignature_reverts() public {
    IRMNRemote.Signature memory sig = s_signatures[s_signatures.length - 1];
    sig.r = _randomBytes32();
    s_signatures.pop();
    s_signatures.push(sig);

    vm.expectRevert(RMNRemote.InvalidSignature.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }

  function test_verify_OutOfOrderSignatures_not_sorted_reverts() public {
    IRMNRemote.Signature memory sig1 = s_signatures[s_signatures.length - 1];
    s_signatures.pop();
    IRMNRemote.Signature memory sig2 = s_signatures[s_signatures.length - 1];
    s_signatures.pop();
    s_signatures.push(sig1);
    s_signatures.push(sig2);

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }

  function test_verify_OutOfOrderSignatures_duplicateSignature_reverts() public {
    IRMNRemote.Signature memory sig = s_signatures[s_signatures.length - 2];
    s_signatures.pop();
    s_signatures.push(sig);

    vm.expectRevert(RMNRemote.OutOfOrderSignatures.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }

  function test_verify_UnexpectedSigner_reverts() public {
    _setupSigners(2); // create 2 new signers that aren't configured on RMNRemote
    _generatePayloadAndSigs(2, 2);

    vm.expectRevert(RMNRemote.UnexpectedSigner.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }

  function test_verify_ThresholdNotMet_reverts() public {
    _generatePayloadAndSigs(2, 1); // 1 sig requested, but 2 required

    vm.expectRevert(RMNRemote.ThresholdNotMet.selector);
    s_rmnRemote.verify(OFF_RAMP_ADDRESS, s_merkleRoots, s_signatures, s_v);
  }
}

contract RMNRemote_curse is RMNRemoteSetup {
  function test_curse_success() public {
    vm.expectEmit();
    emit RMNRemote.Cursed(s_curseSubjects);

    s_rmnRemote.curse(s_curseSubjects);

    assertEq(abi.encode(s_rmnRemote.getCursedSubjects()), abi.encode(s_curseSubjects));
    assertTrue(s_rmnRemote.isCursed(curseSubj1));
    assertTrue(s_rmnRemote.isCursed(curseSubj2));
    // Should not have cursed a random subject
    assertFalse(s_rmnRemote.isCursed(bytes16(keccak256("subject 3"))));
  }

  function test_curse_AlreadyCursed_duplicateSubject_reverts() public {
    s_curseSubjects.push(curseSubj1);

    vm.expectRevert(abi.encodeWithSelector(RMNRemote.AlreadyCursed.selector, curseSubj1));
    s_rmnRemote.curse(s_curseSubjects);
  }

  function test_curse_calledByNonOwner_reverts() public {
    vm.expectRevert("Only callable by owner");
    vm.stopPrank();
    vm.prank(STRANGER);
    s_rmnRemote.curse(s_curseSubjects);
  }
}

contract RMNRemote_uncurse is RMNRemoteSetup {
  function setUp() public override {
    super.setUp();
    s_rmnRemote.curse(s_curseSubjects);
  }

  function test_uncurse_success() public {
    vm.expectEmit();
    emit RMNRemote.Uncursed(s_curseSubjects);

    s_rmnRemote.uncurse(s_curseSubjects);

    assertEq(s_rmnRemote.getCursedSubjects().length, 0);
    assertFalse(s_rmnRemote.isCursed(curseSubj1));
    assertFalse(s_rmnRemote.isCursed(curseSubj2));
  }

  function test_uncurse_NotCursed_duplicatedUncurseSubject_reverts() public {
    s_curseSubjects.push(curseSubj1);

    vm.expectRevert(abi.encodeWithSelector(RMNRemote.NotCursed.selector, curseSubj1));
    s_rmnRemote.uncurse(s_curseSubjects);
  }

  function test_uncurse_calledByNonOwner_reverts() public {
    vm.expectRevert("Only callable by owner");
    vm.stopPrank();
    vm.prank(STRANGER);
    s_rmnRemote.uncurse(s_curseSubjects);
  }
}

contract RMNRemote_global_and_legacy_curses is RMNRemoteSetup {
  function test_global_and_legacy_curses_success() public {
    bytes16 randSubject = bytes16(keccak256("random subject"));
    assertFalse(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject));

    s_rmnRemote.curse(GLOBAL_CURSE_SUBJECT);
    assertTrue(s_rmnRemote.isCursed());
    assertTrue(s_rmnRemote.isCursed(randSubject));

    s_rmnRemote.uncurse(GLOBAL_CURSE_SUBJECT);
    assertFalse(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject));

    s_rmnRemote.curse(LEGACY_CURSE_SUBJECT);
    assertTrue(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject)); // legacy curse doesn't affect specific subjects

    s_rmnRemote.uncurse(LEGACY_CURSE_SUBJECT);
    assertFalse(s_rmnRemote.isCursed());
    assertFalse(s_rmnRemote.isCursed(randSubject));
  }
}
