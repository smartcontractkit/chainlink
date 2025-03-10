// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../../../libraries/Client.sol";
import {CCIPClientSetup} from "./CCIPClientSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPClient_ccipSend is CCIPClientSetup {
  function test_ccipSend_NonNativeFeetoken_DestTokens() public {
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

    s_sender.ccipSend({destChainSelector: DEST_CHAIN_SELECTOR, tokenAmounts: destTokenAmounts, data: ""});

    // Assert that tokens were transfered for bridging + fees
    assertEq(IERC20(token).balanceOf(OWNER), feeTokenBalanceBefore - amount - feeTokenAmount);
  }

  function test_ccipSend_NonNativeFeetoken_NoDestTokens() public {
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

    s_sender.ccipSend({destChainSelector: DEST_CHAIN_SELECTOR, tokenAmounts: destTokenAmounts, data: ""});

    // Assert that tokens were transfered for bridging + fees
    assertEq(IERC20(token).balanceOf(OWNER), feeTokenBalanceBefore - feeTokenAmount);
  }

  function test_ccipSend_with_NativeFeeToken_DestTokens() public {
    address token = address(s_sourceFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    s_sender.updateFeeToken(address(0));

    // Make sure we give the receiver contract enough tokens like CCIP would.
    IERC20(token).approve(address(s_sender), type(uint256).max);

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(address(s_sender)),
      data: "",
      tokenAmounts: destTokenAmounts,
      extraArgs: "",
      feeToken: address(s_sourceFeeToken)
    });

    uint256 feeTokenAmount = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    uint256 tokenBalanceBefore = IERC20(token).balanceOf(OWNER);
    uint256 nativeFeeTokenBalanceBefore = OWNER.balance;

    s_sender.ccipSend{value: feeTokenAmount}({
      destChainSelector: DEST_CHAIN_SELECTOR,
      tokenAmounts: destTokenAmounts,
      data: ""
    });

    // Assert that native fees are paid successfully and tokens are transferred
    assertEq(IERC20(token).balanceOf(OWNER), tokenBalanceBefore - amount, "Tokens were not successfully delivered");
    assertEq(
      OWNER.balance, nativeFeeTokenBalanceBefore - feeTokenAmount, "Native fee tokens were not successfully forwarded"
    );
  }
}
