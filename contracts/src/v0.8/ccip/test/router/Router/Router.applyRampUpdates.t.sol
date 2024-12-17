// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";
import {IRouter} from "../../../interfaces/IRouter.sol";

import {Router} from "../../../Router.sol";
import {Client} from "../../../libraries/Client.sol";
import {BaseTest} from "../../BaseTest.t.sol";
import {MaybeRevertMessageReceiver} from "../../helpers/receivers/MaybeRevertMessageReceiver.sol";

contract Router_applyRampUpdates is BaseTest {
  MaybeRevertMessageReceiver internal s_receiver;

  function setUp() public virtual override {
    super.setUp();
    s_receiver = new MaybeRevertMessageReceiver(false);
  }

  function _assertOffRampRouteSucceeds(
    Router.OffRamp memory offRamp
  ) internal {
    vm.startPrank(offRamp.offRamp);

    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: offRamp.sourceChainSelector,
      sender: bytes("a"),
      data: bytes("a"),
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });

    vm.expectCall(address(s_receiver), abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message));
    s_sourceRouter.routeMessage(message, GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver));
  }

  function _assertOffRampRouteReverts(
    Router.OffRamp memory offRamp
  ) internal {
    vm.startPrank(offRamp.offRamp);

    vm.expectRevert(IRouter.OnlyOffRamp.selector);
    s_sourceRouter.routeMessage(
      Client.Any2EVMMessage({
        messageId: bytes32("a"),
        sourceChainSelector: offRamp.sourceChainSelector,
        sender: bytes("a"),
        data: bytes("a"),
        destTokenAmounts: new Client.EVMTokenAmount[](0)
      }),
      GAS_FOR_CALL_EXACT_CHECK,
      100_000,
      address(s_receiver)
    );
  }

  function testFuzz_applyRampUpdates_OffRampUpdates(
    address[20] memory offRampsInput
  ) public {
    Router.OffRamp[] memory offRamps = new Router.OffRamp[](20);

    for (uint256 i = 0; i < offRampsInput.length; ++i) {
      offRamps[i] = Router.OffRamp({sourceChainSelector: uint64(i), offRamp: offRampsInput[i]});
    }

    // Test adding offRamps
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRamps);

    // There is no uniqueness guarantee on fuzz input, offRamps will not emit in case of a duplicate,
    // hence cannot assert on number of offRamps event emissions, we need to use isOffRa
    for (uint256 i = 0; i < offRamps.length; ++i) {
      assertTrue(s_sourceRouter.isOffRamp(offRamps[i].sourceChainSelector, offRamps[i].offRamp));
    }

    // Test removing offRamps
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), s_sourceRouter.getOffRamps(), new Router.OffRamp[](0));

    assertEq(0, s_sourceRouter.getOffRamps().length);
    for (uint256 i = 0; i < offRamps.length; ++i) {
      assertFalse(s_sourceRouter.isOffRamp(offRamps[i].sourceChainSelector, offRamps[i].offRamp));
    }

    // Testing removing and adding in same call
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), new Router.OffRamp[](0), offRamps);
    s_sourceRouter.applyRampUpdates(new Router.OnRamp[](0), offRamps, offRamps);
    for (uint256 i = 0; i < offRamps.length; ++i) {
      assertTrue(s_sourceRouter.isOffRamp(offRamps[i].sourceChainSelector, offRamps[i].offRamp));
    }
  }

  function test_applyRampUpdates_OffRampUpdatesWithRouting() public {
    // Explicitly construct chain selectors and ramp addresses so we have ramp uniqueness for the various test scenarios.
    uint256 numberOfSelectors = 10;
    uint64[] memory sourceChainSelectors = new uint64[](numberOfSelectors);
    for (uint256 i = 0; i < numberOfSelectors; ++i) {
      sourceChainSelectors[i] = uint64(i);
    }

    uint256 numberOfOffRamps = 5;
    address[] memory offRamps = new address[](numberOfOffRamps);
    for (uint256 i = 0; i < numberOfOffRamps; ++i) {
      offRamps[i] = address(uint160(i * 10));
    }

    // 1st test scenario: add offramps.
    // Check all the offramps are added correctly, and can route messages.
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](numberOfSelectors * numberOfOffRamps);

    // Ensure there are multi-offramp source and multi-source offramps
    for (uint256 i = 0; i < numberOfSelectors; ++i) {
      for (uint256 j = 0; j < numberOfOffRamps; ++j) {
        offRampUpdates[(i * numberOfOffRamps) + j] = Router.OffRamp(sourceChainSelectors[i], offRamps[j]);
      }
    }

    for (uint256 i = 0; i < offRampUpdates.length; ++i) {
      vm.expectEmit();
      emit Router.OffRampAdded(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp);
    }
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    Router.OffRamp[] memory gotOffRamps = s_sourceRouter.getOffRamps();
    assertEq(offRampUpdates.length, gotOffRamps.length);

    for (uint256 i = 0; i < offRampUpdates.length; ++i) {
      assertEq(offRampUpdates[i].offRamp, gotOffRamps[i].offRamp);
      assertTrue(s_sourceRouter.isOffRamp(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp));
      _assertOffRampRouteSucceeds(offRampUpdates[i]);
    }

    vm.startPrank(OWNER);

    // 2nd test scenario: partially remove existing offramps, add new offramps.
    // Check offramps are removed correctly. Removed offramps cannot route messages.
    // Check new offramps are added correctly. New offramps can route messages.
    // Check unmodified offramps remain correct, and can still route messages.
    uint256 numberOfPartialUpdates = offRampUpdates.length / 2;
    Router.OffRamp[] memory partialOffRampRemoves = new Router.OffRamp[](numberOfPartialUpdates);
    Router.OffRamp[] memory partialOffRampAdds = new Router.OffRamp[](numberOfPartialUpdates);
    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      partialOffRampRemoves[i] = offRampUpdates[i];
      partialOffRampAdds[i] = Router.OffRamp({
        sourceChainSelector: offRampUpdates[i].sourceChainSelector,
        offRamp: address(uint160(offRampUpdates[i].offRamp) + 1e18) // Ensure unique new offRamps addresses
      });
    }

    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      vm.expectEmit();
      emit Router.OffRampRemoved(partialOffRampRemoves[i].sourceChainSelector, partialOffRampRemoves[i].offRamp);
    }
    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      vm.expectEmit();
      emit Router.OffRampAdded(partialOffRampAdds[i].sourceChainSelector, partialOffRampAdds[i].offRamp);
    }
    s_sourceRouter.applyRampUpdates(onRampUpdates, partialOffRampRemoves, partialOffRampAdds);

    gotOffRamps = s_sourceRouter.getOffRamps();
    assertEq(offRampUpdates.length, gotOffRamps.length);

    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      assertFalse(
        s_sourceRouter.isOffRamp(partialOffRampRemoves[i].sourceChainSelector, partialOffRampRemoves[i].offRamp)
      );
      _assertOffRampRouteReverts(partialOffRampRemoves[i]);

      assertTrue(s_sourceRouter.isOffRamp(partialOffRampAdds[i].sourceChainSelector, partialOffRampAdds[i].offRamp));
      _assertOffRampRouteSucceeds(partialOffRampAdds[i]);
    }
    for (uint256 i = numberOfPartialUpdates; i < offRampUpdates.length; ++i) {
      assertTrue(s_sourceRouter.isOffRamp(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp));
      _assertOffRampRouteSucceeds(offRampUpdates[i]);
    }

    vm.startPrank(OWNER);

    // 3rd test scenario: remove all offRamps.
    // Check all offramps have been removed, no offramp is able to route messages.
    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      vm.expectEmit();
      emit Router.OffRampRemoved(partialOffRampAdds[i].sourceChainSelector, partialOffRampAdds[i].offRamp);
    }
    s_sourceRouter.applyRampUpdates(onRampUpdates, partialOffRampAdds, new Router.OffRamp[](0));

    uint256 numberOfRemainingOfframps = offRampUpdates.length - numberOfPartialUpdates;
    Router.OffRamp[] memory remainingOffRampRemoves = new Router.OffRamp[](numberOfRemainingOfframps);
    for (uint256 i = 0; i < numberOfRemainingOfframps; ++i) {
      remainingOffRampRemoves[i] = offRampUpdates[i + numberOfPartialUpdates];
    }

    for (uint256 i = 0; i < numberOfRemainingOfframps; ++i) {
      vm.expectEmit();
      emit Router.OffRampRemoved(remainingOffRampRemoves[i].sourceChainSelector, remainingOffRampRemoves[i].offRamp);
    }
    s_sourceRouter.applyRampUpdates(onRampUpdates, remainingOffRampRemoves, new Router.OffRamp[](0));

    // Check there are no offRamps.
    assertEq(0, s_sourceRouter.getOffRamps().length);

    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      assertFalse(s_sourceRouter.isOffRamp(partialOffRampAdds[i].sourceChainSelector, partialOffRampAdds[i].offRamp));
      _assertOffRampRouteReverts(partialOffRampAdds[i]);
    }
    for (uint256 i = 0; i < offRampUpdates.length; ++i) {
      assertFalse(s_sourceRouter.isOffRamp(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp));
      _assertOffRampRouteReverts(offRampUpdates[i]);
    }

    vm.startPrank(OWNER);

    // 4th test scenario: add initial onRamps back.
    // Check the offramps are added correctly, and can route messages.
    // Check offramps that were not added back remain unset, and cannot route messages.
    for (uint256 i = 0; i < offRampUpdates.length; ++i) {
      vm.expectEmit();
      emit Router.OffRampAdded(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp);
    }
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    // Check initial offRamps are added back and can route to receiver.
    gotOffRamps = s_sourceRouter.getOffRamps();
    assertEq(offRampUpdates.length, gotOffRamps.length);

    for (uint256 i = 0; i < offRampUpdates.length; ++i) {
      assertEq(offRampUpdates[i].offRamp, gotOffRamps[i].offRamp);
      assertTrue(s_sourceRouter.isOffRamp(offRampUpdates[i].sourceChainSelector, offRampUpdates[i].offRamp));
      _assertOffRampRouteSucceeds(offRampUpdates[i]);
    }

    // Check offramps that were not added back remain unset.
    for (uint256 i = 0; i < numberOfPartialUpdates; ++i) {
      assertFalse(s_sourceRouter.isOffRamp(partialOffRampAdds[i].sourceChainSelector, partialOffRampAdds[i].offRamp));
      _assertOffRampRouteReverts(partialOffRampAdds[i]);
    }
  }

  function testFuzz_applyRampUpdates_OnRampUpdates(
    Router.OnRamp[] memory onRamps
  ) public {
    // Test adding onRamps
    for (uint256 i = 0; i < onRamps.length; ++i) {
      vm.expectEmit();
      emit Router.OnRampSet(onRamps[i].destChainSelector, onRamps[i].onRamp);
    }

    s_sourceRouter.applyRampUpdates(onRamps, new Router.OffRamp[](0), new Router.OffRamp[](0));

    // Test setting onRamps to unsupported
    for (uint256 i = 0; i < onRamps.length; ++i) {
      onRamps[i].onRamp = address(0);

      vm.expectEmit();
      emit Router.OnRampSet(onRamps[i].destChainSelector, onRamps[i].onRamp);
    }
    s_sourceRouter.applyRampUpdates(onRamps, new Router.OffRamp[](0), new Router.OffRamp[](0));
    for (uint256 i = 0; i < onRamps.length; ++i) {
      assertEq(address(0), s_sourceRouter.getOnRamp(onRamps[i].destChainSelector));
      assertFalse(s_sourceRouter.isChainSupported(onRamps[i].destChainSelector));
    }
  }

  function test_applyRampUpdates_OnRampDisable() public {
    // Add onRamp
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](1);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](0);
    address onRamp = address(uint160(2));
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);
    assertEq(onRamp, s_sourceRouter.getOnRamp(DEST_CHAIN_SELECTOR));
    assertTrue(s_sourceRouter.isChainSupported(DEST_CHAIN_SELECTOR));

    // Disable onRamp
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: address(0)});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));
    assertEq(address(0), s_sourceRouter.getOnRamp(DEST_CHAIN_SELECTOR));
    assertFalse(s_sourceRouter.isChainSupported(DEST_CHAIN_SELECTOR));

    // Re-enable onRamp
    onRampUpdates[0] = Router.OnRamp({destChainSelector: DEST_CHAIN_SELECTOR, onRamp: onRamp});
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), new Router.OffRamp[](0));
    assertEq(onRamp, s_sourceRouter.getOnRamp(DEST_CHAIN_SELECTOR));
    assertTrue(s_sourceRouter.isChainSupported(DEST_CHAIN_SELECTOR));
  }

  function test_RevertWhen_applyRampUpdatesWhen_OnlyOwner() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](0);
    s_sourceRouter.applyRampUpdates(onRampUpdates, offRampUpdates, offRampUpdates);
  }

  function test_RevertWhen_applyRampUpdatesWhen_OffRampMismatch() public {
    address offRamp = address(uint160(2));

    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](1);
    offRampUpdates[0] = Router.OffRamp(DEST_CHAIN_SELECTOR, offRamp);

    vm.expectEmit();
    emit Router.OffRampAdded(DEST_CHAIN_SELECTOR, offRamp);
    s_sourceRouter.applyRampUpdates(onRampUpdates, new Router.OffRamp[](0), offRampUpdates);

    offRampUpdates[0] = Router.OffRamp(SOURCE_CHAIN_SELECTOR, offRamp);

    vm.expectRevert(abi.encodeWithSelector(Router.OffRampMismatch.selector, SOURCE_CHAIN_SELECTOR, offRamp));
    s_sourceRouter.applyRampUpdates(onRampUpdates, offRampUpdates, offRampUpdates);
  }
}
