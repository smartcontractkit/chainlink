// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../../../interfaces/IRouterClient.sol";

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPSender} from "../../../../applications/external/CCIPSender.sol";
import {Client} from "../../../../libraries/Client.sol";
import {OnRampSetup} from "../../../OnRamp/OnRamp/OnRampSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPSender_ccipSend is OnRampSetup {
  bytes32 public constant MESSAGE_ID = keccak256("FAKE_MESSAGE_ID");

  CCIPSender internal s_sender;

  function setUp() public virtual override {
    OnRampSetup.setUp();

    vm.mockCall(
      address(s_sourceRouter), abi.encodeWithSelector(IRouterClient.ccipSend.selector), abi.encode(MESSAGE_ID)
    );

    s_sender = new CCIPSender(address(s_sourceRouter));

    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      recipient: abi.encode(address(s_sender)),
      extraArgsBytes: ""
    });
    s_sender.applyChainUpdates(chainUpdates);
  }

  function test_ccipSend_NonNativeFeetoken_andDestTokens() public {
    address token = address(s_sourceFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    IERC20(token).approve(address(s_sender), type(uint256).max);

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(address(s_sender)),
      data: "",
      tokenAmounts: destTokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: ""
    });

    uint256 feeTokenAmount = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    uint256 feeTokenBalanceBefore = IERC20(s_sourceFeeToken).balanceOf(OWNER);

    vm.expectEmit();
    emit CCIPSender.MessageSent(MESSAGE_ID);

    bytes32 messageId = s_sender.ccipSend({
      destChainSelector: DEST_CHAIN_SELECTOR,
      tokenAmounts: destTokenAmounts,
      data: "",
      feeToken: address(s_sourceFeeToken)
    });

    assertEq(messageId, MESSAGE_ID);
    assertEq(IERC20(token).balanceOf(OWNER), feeTokenBalanceBefore - amount - feeTokenAmount);
  }

  function test_ccipSend_withNonNativeFeetoken_andNoDestTokens() public {
    address token = address(s_sourceFeeToken);
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](0);

    // Make sure we give the receiver contract enough tokens like CCIP would.
    IERC20(token).approve(address(s_sender), type(uint256).max);

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(address(s_sender)),
      data: "",
      tokenAmounts: destTokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: ""
    });

    uint256 feeTokenAmount = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    uint256 feeTokenBalanceBefore = IERC20(s_sourceFeeToken).balanceOf(OWNER);

    vm.expectEmit();
    emit CCIPSender.MessageSent(MESSAGE_ID);

    s_sender.ccipSend({
      destChainSelector: DEST_CHAIN_SELECTOR,
      tokenAmounts: destTokenAmounts,
      data: "",
      feeToken: address(s_sourceFeeToken)
    });

    // Assert that tokens were transfered for bridging + fees
    assertEq(IERC20(token).balanceOf(OWNER), feeTokenBalanceBefore - feeTokenAmount);
  }

  function test_ccipSend_NativeFeeToken_andDestTokens_Success() public {
    address token = address(s_sourceFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    IERC20(token).approve(address(s_sender), type(uint256).max);

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(address(s_sender)),
      data: "",
      tokenAmounts: destTokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    uint256 feeTokenAmount = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    uint256 tokenBalanceBefore = IERC20(token).balanceOf(OWNER);
    uint256 nativeFeeTokenBalanceBefore = OWNER.balance;

    vm.expectEmit();
    emit CCIPSender.MessageSent(MESSAGE_ID);

    s_sender.ccipSend{value: feeTokenAmount}({
      destChainSelector: DEST_CHAIN_SELECTOR,
      tokenAmounts: destTokenAmounts,
      data: "",
      feeToken: address(0)
    });

    // Assert that native fees are paid successfully and tokens are transferred
    assertEq(IERC20(token).balanceOf(OWNER), tokenBalanceBefore - amount);
    assertEq(OWNER.balance, nativeFeeTokenBalanceBefore - feeTokenAmount);
  }
}
