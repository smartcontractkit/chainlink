// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNRemoteSetup} from "./RMNRemoteSetup.t.sol";

contract RMNRemote_setConfig is RMNRemoteSetup {
  function test_RevertWhen_setConfig_ZeroValueNotAllowed() public {
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: bytes32(0), signers: s_signers, fSign: 1});

    vm.expectRevert(RMNRemote.ZeroValueNotAllowed.selector);

    s_rmnRemote.setConfig(config);
  }

  function test_setConfig_addSigner_removeSigner() public {
    uint32 currentConfigVersion = 0;
    uint256 numSigners = s_signers.length;
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 1});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    // add a signer
    address newSigner = makeAddr("new signer");
    s_signers.push(RMNRemote.Signer({onchainPublicKey: newSigner, nodeIndex: uint64(numSigners)}));
    config = RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 1});

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
    config = RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 1});

    vm.expectEmit();
    emit RMNRemote.ConfigSet(++currentConfigVersion, config);

    s_rmnRemote.setConfig(config);

    (version, gotConfig) = s_rmnRemote.getVersionedConfig();
    assertEq(gotConfig.signers.length, s_signers.length);
    assertEq(version, currentConfigVersion);
  }

  function test_RevertWhen_setConfig_invalidSignerOrder() public {
    s_signers.push(RMNRemote.Signer({onchainPublicKey: address(4), nodeIndex: 0}));
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 1});

    vm.expectRevert(RMNRemote.InvalidSignerOrder.selector);
    s_rmnRemote.setConfig(config);
  }

  function test_RevertWhen_setConfig_notEnoughSigners() public {
    RMNRemote.Config memory config = RMNRemote.Config({
      rmnHomeContractConfigDigest: _randomBytes32(),
      signers: s_signers,
      fSign: uint64(s_signers.length / 2) // at least 2f+1 is required
    });

    vm.expectRevert(RMNRemote.NotEnoughSigners.selector);
    s_rmnRemote.setConfig(config);
  }

  function test_RevertWhen_setConfig_duplicateOnChainPublicKey() public {
    s_signers.push(RMNRemote.Signer({onchainPublicKey: s_signerWallets[0].addr, nodeIndex: uint64(s_signers.length)}));
    RMNRemote.Config memory config =
      RMNRemote.Config({rmnHomeContractConfigDigest: _randomBytes32(), signers: s_signers, fSign: 1});

    vm.expectRevert(RMNRemote.DuplicateOnchainPublicKey.selector);
    s_rmnRemote.setConfig(config);
  }
}
