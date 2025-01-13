// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {NonceManager} from "../../NonceManager.sol";
import {OnRampSetup} from "../onRamp/OnRamp/OnRampSetup.t.sol";

contract NonceManager_applyPreviousRampsUpdates is OnRampSetup {
  function test_SingleRampUpdate() public {
    address prevOnRamp = makeAddr("prevOnRamp");
    address prevOffRamp = makeAddr("prevOffRamp");
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp, prevOffRamp),
      overrideExistingRamps: false
    });

    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR, previousRamps[0].prevRamps);

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
  }

  function test_MultipleRampsUpdates() public {
    address prevOnRamp1 = makeAddr("prevOnRamp1");
    address prevOnRamp2 = makeAddr("prevOnRamp2");
    address prevOffRamp1 = makeAddr("prevOffRamp1");
    address prevOffRamp2 = makeAddr("prevOffRamp2");
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](2);
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp1, prevOffRamp1),
      overrideExistingRamps: false
    });
    previousRamps[1] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR + 1,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp2, prevOffRamp2),
      overrideExistingRamps: false
    });

    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR, previousRamps[0].prevRamps);
    vm.expectEmit();
    emit NonceManager.PreviousRampsUpdated(DEST_CHAIN_SELECTOR + 1, previousRamps[1].prevRamps);

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    _assertPreviousRampsEqual(s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR), previousRamps[0].prevRamps);
    _assertPreviousRampsEqual(
      s_outboundNonceManager.getPreviousRamps(DEST_CHAIN_SELECTOR + 1), previousRamps[1].prevRamps
    );
  }

  function test_PreviousRampAlreadySet_overrideAllowed() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOffRamp = makeAddr("prevOffRamp");
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(address(0), prevOffRamp),
      overrideExistingRamps: true
    });

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(address(0), prevOffRamp),
      overrideExistingRamps: true
    });

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function test_ZeroInput() public {
    vm.recordLogs();
    s_outboundNonceManager.applyPreviousRampsUpdates(new NonceManager.PreviousRampsArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_RevertWhen_applyPreviousRampsUpdatesWhen_PreviousRampAlreadySetOnRamp() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOnRamp = makeAddr("prevOnRamp");
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp, address(0)),
      overrideExistingRamps: false
    });

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp, address(0)),
      overrideExistingRamps: false
    });

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function test_RevertWhen_applyPreviousRampsUpdatesWhen_PreviousRampAlreadySetOffRamp() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOffRamp = makeAddr("prevOffRamp");
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(address(0), prevOffRamp),
      overrideExistingRamps: false
    });

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(address(0), prevOffRamp),
      overrideExistingRamps: false
    });

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function test_RevertWhen_applyPreviousRampsUpdatesWhen_PreviousRampAlreadySetOnRampAndOffRamp_Revert() public {
    NonceManager.PreviousRampsArgs[] memory previousRamps = new NonceManager.PreviousRampsArgs[](1);
    address prevOnRamp = makeAddr("prevOnRamp");
    address prevOffRamp = makeAddr("prevOffRamp");
    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp, prevOffRamp),
      overrideExistingRamps: false
    });

    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);

    previousRamps[0] = NonceManager.PreviousRampsArgs({
      remoteChainSelector: DEST_CHAIN_SELECTOR,
      prevRamps: NonceManager.PreviousRamps(prevOnRamp, prevOffRamp),
      overrideExistingRamps: false
    });

    vm.expectRevert(NonceManager.PreviousRampAlreadySet.selector);
    s_outboundNonceManager.applyPreviousRampsUpdates(previousRamps);
  }

  function _assertPreviousRampsEqual(
    NonceManager.PreviousRamps memory a,
    NonceManager.PreviousRamps memory b
  ) internal pure {
    assertEq(a.prevOnRamp, b.prevOnRamp);
    assertEq(a.prevOffRamp, b.prevOffRamp);
  }
}
