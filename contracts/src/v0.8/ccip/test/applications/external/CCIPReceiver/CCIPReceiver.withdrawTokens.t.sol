// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {Client} from "../../../../libraries/Client.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPSender_withdrawToekns is CCIPReceiverSetup {
  function test_withdrawTokens_ERC20() public {
    bytes32 messageId = keccak256("messageId");
    address token = address(s_destFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    deal(token, address(s_receiver), amount);

    // The receiver contract will revert if the router is not the sender.
    vm.startPrank(address(s_destRouter));

    vm.expectEmit();
    emit MessageFailed(
      messageId, abi.encodeWithSelector(bytes4(CCIPBase.InvalidSender.selector), abi.encode(address(1)))
    );

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: sourceChainSelector,
        sender: abi.encode(address(1)),
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );

    vm.stopPrank();

    // Check that the message was stored properly by comparing each of the fields.
    // There's no way to check that a function internally will revert from a top-level test, so we need to check state differences
    Client.Any2EVMMessage memory failedMessage = s_receiver.getMessageContents(messageId);
    assertEq(failedMessage.sender, abi.encode(address(1)));
    assertEq(failedMessage.sourceChainSelector, sourceChainSelector);
    assertEq(failedMessage.destTokenAmounts[0].token, token);
    assertEq(failedMessage.destTokenAmounts[0].amount, amount);

    // Check that message status is failed
    assertTrue(s_receiver.isFailedMessage(messageId), "Message should be marked as failed");

    uint256 tokenBalanceBefore = IERC20(token).balanceOf(OWNER);

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_receiver), OWNER, amount);
    s_receiver.withdrawTokens(token, OWNER, amount);

    assertEq(IERC20(token).balanceOf(OWNER), tokenBalanceBefore + amount);
    assertGt(IERC20(token).balanceOf(OWNER), 0);
  }

  function test_withdrawTokens_Native() public {
    uint256 amount = 100 ether;
    deal(address(s_receiver), amount);

    uint256 balanceBefore = OWNER.balance;

    vm.startPrank(OWNER);

    s_receiver.withdrawTokens(address(0), payable(OWNER), amount);

    assertEq(OWNER.balance, balanceBefore + amount);
  }
}
