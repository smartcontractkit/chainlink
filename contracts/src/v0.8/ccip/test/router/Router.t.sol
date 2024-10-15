// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";
import {IRouter} from "../../interfaces/IRouter.sol";
import {IRouterClient} from "../../interfaces/IRouterClient.sol";
import {IWrappedNative} from "../../interfaces/IWrappedNative.sol";

import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {OnRamp} from "../../onRamp/OnRamp.sol";
import {MaybeRevertMessageReceiver} from "../helpers/receivers/MaybeRevertMessageReceiver.sol";
import {OffRampSetup} from "../offRamp/OffRampSetup.t.sol";
import {OnRampSetup} from "../onRamp/OnRampSetup.t.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract Router_constructor is OnRampSetup {
  function test_Constructor_Success() public view {
    assertEq("Router 1.2.0", s_sourceRouter.typeAndVersion());
    assertEq(OWNER, s_sourceRouter.owner());
  }
}

contract Router_recoverTokens is OnRampSetup {
  function test_RecoverTokens_Success() public {
    // Assert we can recover sourceToken
    IERC20 token = IERC20(s_sourceTokens[0]);
    uint256 balanceBefore = token.balanceOf(OWNER);
    token.transfer(address(s_sourceRouter), 1);
    assertEq(token.balanceOf(address(s_sourceRouter)), 1);
    s_sourceRouter.recoverTokens(address(token), OWNER, 1);
    assertEq(token.balanceOf(address(s_sourceRouter)), 0);
    assertEq(token.balanceOf(OWNER), balanceBefore);

    // Assert we can recover native
    balanceBefore = OWNER.balance;
    deal(address(s_sourceRouter), 10);
    assertEq(address(s_sourceRouter).balance, 10);
    s_sourceRouter.recoverTokens(address(0), OWNER, 10);
    assertEq(OWNER.balance, balanceBefore + 10);
    assertEq(address(s_sourceRouter).balance, 0);
  }

  function test_RecoverTokensNonOwner_Revert() public {
    // Reverts if not owner
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_sourceRouter.recoverTokens(address(0), STRANGER, 1);
  }

  function test_RecoverTokensInvalidRecipient_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(Router.InvalidRecipientAddress.selector, address(0)));
    s_sourceRouter.recoverTokens(address(0), address(0), 1);
  }

  function test_RecoverTokensNoFunds_Revert() public {
    // Reverts if no funds present
    vm.expectRevert();
    s_sourceRouter.recoverTokens(address(0), OWNER, 10);
  }

  function test_RecoverTokensValueReceiver_Revert() public {
    MaybeRevertMessageReceiver revertingValueReceiver = new MaybeRevertMessageReceiver(true);
    deal(address(s_sourceRouter), 10);

    // Value receiver reverts
    vm.expectRevert(Router.FailedToSendValue.selector);
    s_sourceRouter.recoverTokens(address(0), address(revertingValueReceiver), 10);
  }
}

contract Router_ccipSend is OnRampSetup {
  event Burned(address indexed sender, uint256 amount);

  function test_CCIPSendLinkFeeOneTokenSuccess_gas() public {
    vm.pauseGasMetering();
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    IERC20 sourceToken1 = IERC20(s_sourceTokens[1]);
    sourceToken1.approve(address(s_sourceRouter), 2 ** 64);

    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 2 ** 64;
    message.tokenAmounts[0].token = s_sourceTokens[1];

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 0);

    uint256 balanceBefore = sourceToken1.balanceOf(OWNER);

    // Assert that the tokens are burned
    vm.expectEmit();
    emit Burned(address(s_onRamp), message.tokenAmounts[0].amount);

    Internal.EVM2AnyRampMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, msgEvent.header.sequenceNumber, msgEvent);

    vm.resumeGasMetering();
    bytes32 messageId = s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    assertEq(msgEvent.header.messageId, messageId);
    // Assert the user balance is lowered by the tokenAmounts sent and the fee amount
    uint256 expectedBalance = balanceBefore - (message.tokenAmounts[0].amount);
    assertEq(expectedBalance, sourceToken1.balanceOf(OWNER));
    vm.resumeGasMetering();
  }

  function test_CCIPSendLinkFeeNoTokenSuccess_gas() public {
    vm.pauseGasMetering();
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 0);

    Internal.EVM2AnyRampMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, msgEvent.header.sequenceNumber, msgEvent);

    vm.resumeGasMetering();
    bytes32 messageId = s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    assertEq(msgEvent.header.messageId, messageId);
    vm.resumeGasMetering();
  }

  function test_ccipSend_nativeFeeOneTokenSuccess_gas() public {
    vm.pauseGasMetering();
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    IERC20 sourceToken1 = IERC20(s_sourceTokens[1]);
    sourceToken1.approve(address(s_sourceRouter), 2 ** 64);

    uint256 balanceBefore = sourceToken1.balanceOf(OWNER);

    message.tokenAmounts = new Client.EVMTokenAmount[](1);
    message.tokenAmounts[0].amount = 2 ** 64;
    message.tokenAmounts[0].token = s_sourceTokens[1];
    // Native fees will be wrapped so we need to calculate the event with
    // the wrapped native feeCoin address.
    message.feeToken = s_sourceRouter.getWrappedNative();
    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 0);

    Internal.EVM2AnyRampMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);
    msgEvent.feeValueJuels = expectedFee * s_sourceTokenPrices[1] / s_sourceTokenPrices[0];

    message.feeToken = address(0);
    // Assert that the tokens are burned
    vm.expectEmit();
    emit Burned(address(s_onRamp), message.tokenAmounts[0].amount);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, msgEvent.header.sequenceNumber, msgEvent);

    vm.resumeGasMetering();
    bytes32 messageId = s_sourceRouter.ccipSend{value: expectedFee}(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    assertEq(msgEvent.header.messageId, messageId);
    // Assert the user balance is lowered by the tokenAmounts sent and the fee amount
    uint256 expectedBalance = balanceBefore - (message.tokenAmounts[0].amount);
    assertEq(expectedBalance, sourceToken1.balanceOf(OWNER));
    vm.resumeGasMetering();
  }

  function test_ccipSend_nativeFeeNoTokenSuccess_gas() public {
    vm.pauseGasMetering();
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    // Native fees will be wrapped so we need to calculate the event with
    // the wrapped native feeCoin address.
    message.feeToken = s_sourceRouter.getWrappedNative();
    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 0);

    Internal.EVM2AnyRampMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee, OWNER);
    msgEvent.feeValueJuels = expectedFee * s_sourceTokenPrices[1] / s_sourceTokenPrices[0];
    // Set it to address(0) to indicate native
    message.feeToken = address(0);

    vm.expectEmit();
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, msgEvent.header.sequenceNumber, msgEvent);

    vm.resumeGasMetering();
    bytes32 messageId = s_sourceRouter.ccipSend{value: expectedFee}(DEST_CHAIN_SELECTOR, message);
    vm.pauseGasMetering();

    assertEq(msgEvent.header.messageId, messageId);
    // Assert the user balance is lowered by the tokenAmounts sent and the fee amount
    vm.resumeGasMetering();
  }

  function test_NonLinkFeeToken_Success() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_feeQuoter.applyFeeTokensUpdates(feeTokens, new address[](0));

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];
    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), 2 ** 64);
    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeToken_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = address(0); // Raw native
    uint256 nativeQuote = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    vm.stopPrank();
    hoax(address(1), 100 ether);
    s_sourceRouter.ccipSend{value: nativeQuote}(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeTokenOverpay_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = address(0); // Raw native
    uint256 nativeQuote = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    vm.stopPrank();
    hoax(address(1), 100 ether);
    s_sourceRouter.ccipSend{value: nativeQuote + 1}(DEST_CHAIN_SELECTOR, message);
    // We expect the overpayment to be taken in full.
    assertEq(address(1).balance, 100 ether - (nativeQuote + 1));
    assertEq(address(s_sourceRouter).balance, 0);
  }

  function test_WrappedNativeFeeToken_Success() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceRouter.getWrappedNative();
    uint256 nativeQuote = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    vm.stopPrank();
    hoax(address(1), 100 ether);
    // Now address(1) has nativeQuote wrapped.
    IWrappedNative(s_sourceRouter.getWrappedNative()).deposit{value: nativeQuote}();
    IWrappedNative(s_sourceRouter.getWrappedNative()).approve(address(s_sourceRouter), nativeQuote);
    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_WhenNotHealthy_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(Router.BadARMSignal.selector);
    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_UnsupportedDestinationChain_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 wrongChain = DEST_CHAIN_SELECTOR + 1;

    vm.expectRevert(abi.encodeWithSelector(IRouterClient.UnsupportedDestinationChain.selector, wrongChain));

    s_sourceRouter.ccipSend(wrongChain, message);
  }

  function test_FeeTokenAmountTooLow_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    IERC20(s_sourceTokens[0]).approve(address(s_sourceRouter), 0);

    vm.expectRevert("ERC20: insufficient allowance");

    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_InvalidMsgValue() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Non-empty feeToken but with msg.value should revert
    vm.stopPrank();
    hoax(address(1), 1);
    vm.expectRevert(IRouterClient.InvalidMsgValue.selector);
    s_sourceRouter.ccipSend{value: 1}(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeTokenZeroValue() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = address(0); // Raw native
    // Include no value, should revert
    vm.expectRevert();
    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeTokenInsufficientValue() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = address(0); // Raw native
    // Include insufficient, should also revert
    vm.stopPrank();

    hoax(address(1), 1);
    vm.expectRevert(IRouterClient.InsufficientFeeTokenAmount.selector);
    s_sourceRouter.ccipSend{value: 1}(DEST_CHAIN_SELECTOR, message);
  }
}

contract Router_getArmProxy is RouterSetup {
  function test_getArmProxy() public view {
    assertEq(s_sourceRouter.getArmProxy(), address(s_mockRMN));
  }
}

contract Router_applyRampUpdates is RouterSetup {
  MaybeRevertMessageReceiver internal s_receiver;

  function setUp() public virtual override(RouterSetup) {
    super.setUp();
    s_receiver = new MaybeRevertMessageReceiver(false);
  }

  function _assertOffRampRouteSucceeds(
    Router.OffRamp memory offRamp
  ) internal {
    vm.startPrank(offRamp.offRamp);

    Client.Any2EVMMessage memory message = _generateReceiverMessage(offRamp.sourceChainSelector);
    vm.expectCall(address(s_receiver), abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message));
    s_sourceRouter.routeMessage(message, GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver));
  }

  function _assertOffRampRouteReverts(
    Router.OffRamp memory offRamp
  ) internal {
    vm.startPrank(offRamp.offRamp);

    vm.expectRevert(IRouter.OnlyOffRamp.selector);
    s_sourceRouter.routeMessage(
      _generateReceiverMessage(offRamp.sourceChainSelector), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );
  }

  function test_Fuzz_OffRampUpdates(
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

  function test_OffRampUpdatesWithRouting() public {
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

  function test_Fuzz_OnRampUpdates(
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

  function test_OnRampDisable() public {
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

  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    Router.OnRamp[] memory onRampUpdates = new Router.OnRamp[](0);
    Router.OffRamp[] memory offRampUpdates = new Router.OffRamp[](0);
    s_sourceRouter.applyRampUpdates(onRampUpdates, offRampUpdates, offRampUpdates);
  }

  function test_OffRampMismatch_Revert() public {
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

contract Router_setWrappedNative is OnRampSetup {
  function test_Fuzz_SetWrappedNative_Success(
    address wrappedNative
  ) public {
    s_sourceRouter.setWrappedNative(wrappedNative);
    assertEq(wrappedNative, s_sourceRouter.getWrappedNative());
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_sourceRouter.setWrappedNative(address(1));
  }
}

contract Router_getSupportedTokens is OnRampSetup {
  function test_GetSupportedTokens_Revert() public {
    vm.expectRevert(OnRamp.GetSupportedTokensFunctionalityRemovedCheckAdminRegistry.selector);
    s_onRamp.getSupportedTokens(DEST_CHAIN_SELECTOR);
  }
}

contract Router_routeMessage is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    vm.startPrank(address(s_offRamp));
  }

  function _generateManualGasLimit(
    uint256 callDataLength
  ) internal view returns (uint256) {
    return ((gasleft() - 2 * (16 * callDataLength + GAS_FOR_CALL_EXACT_CHECK)) * 62) / 64;
  }

  function test_routeMessage_ManualExec_Success() public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    // Manuel execution cannot run out of gas

    (bool success, bytes memory retData, uint256 gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_receiver)
    );
    assertTrue(success);
    assertEq("", retData);
    assertGt(gasUsed, 3_000);
  }

  function test_routeMessage_ExecutionEvent_Success() public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    // Should revert with reason
    bytes memory realError1 = new bytes(2);
    realError1[0] = 0xbe;
    realError1[1] = 0xef;
    s_reverting_receiver.setErr(realError1);

    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (bool success, bytes memory retData, uint256 gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, realError1), retData);
    assertGt(gasUsed, 3_000);

    // Reason is truncated
    // Over the MAX_RET_BYTES limit (including offset and length word since we have a dynamic values), should be ignored
    bytes memory realError2 = new bytes(32 * 2 + 1);
    realError2[32 * 2 - 1] = 0xAA;
    realError2[32 * 2] = 0xFF;
    s_reverting_receiver.setErr(realError2);

    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (success, retData, gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(
      abi.encodeWithSelector(
        MaybeRevertMessageReceiver.CustomError.selector,
        uint256(32),
        uint256(realError2.length),
        uint256(0),
        uint256(0xAA)
      ),
      retData
    );
    assertGt(gasUsed, 3_000);

    // Should emit success
    vm.expectEmit();
    emit Router.MessageExecuted(
      message.messageId,
      message.sourceChainSelector,
      address(s_offRamp),
      keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
    );

    (success, retData, gasUsed) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_receiver)
    );

    assertTrue(success);
    assertEq("", retData);
    assertGt(gasUsed, 3_000);
  }

  function testFuzz_routeMessage_ExecutionEvent_Success(
    bytes calldata error
  ) public {
    Client.Any2EVMMessage memory message = _generateReceiverMessage(SOURCE_CHAIN_SELECTOR);
    s_reverting_receiver.setErr(error);

    bytes memory expectedRetData;

    if (error.length >= 33) {
      uint256 cutOff = error.length > 64 ? 64 : error.length;
      vm.expectEmit();
      emit Router.MessageExecuted(
        message.messageId,
        message.sourceChainSelector,
        address(s_offRamp),
        keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
      );
      expectedRetData = abi.encodeWithSelector(
        MaybeRevertMessageReceiver.CustomError.selector,
        uint256(32),
        uint256(error.length),
        bytes32(error[:32]),
        bytes32(error[32:cutOff])
      );
    } else {
      vm.expectEmit();
      emit Router.MessageExecuted(
        message.messageId,
        message.sourceChainSelector,
        address(s_offRamp),
        keccak256(abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message))
      );
      expectedRetData = abi.encodeWithSelector(MaybeRevertMessageReceiver.CustomError.selector, error);
    }

    (bool success, bytes memory retData,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR),
      GAS_FOR_CALL_EXACT_CHECK,
      _generateManualGasLimit(message.data.length),
      address(s_reverting_receiver)
    );

    assertFalse(success);
    assertEq(expectedRetData, retData);
  }

  function test_routeMessage_AutoExec_Success() public {
    (bool success,,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );

    assertTrue(success);

    (success,,) = s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 1, address(s_receiver)
    );

    // Can run out of gas, should return false
    assertFalse(success);
  }

  // Reverts
  function test_routeMessage_OnlyOffRamp_Revert() public {
    vm.stopPrank();
    vm.startPrank(STRANGER);

    vm.expectRevert(IRouter.OnlyOffRamp.selector);
    s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );
  }

  function test_routeMessage_WhenNotHealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(Router.BadARMSignal.selector);
    s_destRouter.routeMessage(
      _generateReceiverMessage(SOURCE_CHAIN_SELECTOR), GAS_FOR_CALL_EXACT_CHECK, 100_000, address(s_receiver)
    );
  }
}

contract Router_getFee is OnRampSetup {
  function test_GetFeeSupportedChain_Success() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 10e9);
  }

  // Reverts
  function test_UnsupportedDestinationChain_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(abi.encodeWithSelector(IRouterClient.UnsupportedDestinationChain.selector, 999));
    s_sourceRouter.getFee(999, message);
  }
}
