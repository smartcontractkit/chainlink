// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

import {Router} from "../../../Router.sol";
import {IRouterClient} from "../../../interfaces/IRouterClient.sol";
import {IWrappedNative} from "../../../interfaces/IWrappedNative.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";

import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampSetup} from "../../onRamp/OnRamp/OnRampSetup.t.sol";

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

  function test_NonLinkFeeToken() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];
    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];
    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), 2 ** 64);
    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeToken() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = address(0); // Raw native
    uint256 nativeQuote = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    vm.stopPrank();
    hoax(address(1), 100 ether);
    s_sourceRouter.ccipSend{value: nativeQuote}(DEST_CHAIN_SELECTOR, message);
  }

  function test_NativeFeeTokenOverpay() public {
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

  function test_WrappedNativeFeeToken() public {
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

  function test_RevertWhen_WhenNotHealthy() public {
    vm.mockCall(address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encode(true));

    vm.expectRevert(Router.BadARMSignal.selector);

    s_sourceRouter.ccipSend(DEST_CHAIN_SELECTOR, _generateEmptyMessage());
  }

  function test_RevertWhen_UnsupportedDestinationChain() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 wrongChain = DEST_CHAIN_SELECTOR + 1;

    vm.expectRevert(abi.encodeWithSelector(IRouterClient.UnsupportedDestinationChain.selector, wrongChain));

    s_sourceRouter.ccipSend(wrongChain, message);
  }

  function test_RevertWhen_FeeTokenAmountTooLow() public {
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
