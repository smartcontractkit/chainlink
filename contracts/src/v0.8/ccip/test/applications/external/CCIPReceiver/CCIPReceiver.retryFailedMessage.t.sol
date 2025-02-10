// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPReceiver} from "../../../../applications/external/CCIPReceiver.sol";
import {Client} from "../../../../libraries/Client.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

contract CCIPSender_retryFailedMessage is CCIPReceiverSetup {
  function test_retryFailedMessage() public {
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
        sourceChainSelector: s_sourceChainSelector,
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
    assertEq(failedMessage.sourceChainSelector, s_sourceChainSelector);
    assertEq(failedMessage.destTokenAmounts[0].token, token);
    assertEq(failedMessage.destTokenAmounts[0].amount, amount);

    // Check that message status is failed
    assertTrue(s_receiver.isFailedMessage(messageId), "Message should be marked as failed");

    vm.startPrank(OWNER);

    // The message failed initially because the sender was not approved. Now we approve it and retry processing. Because retryFailedMessage() calls processMessage normally, it should execute successfully now.
    CCIPBase.ApprovedSenderUpdate[] memory senderUpdates = new CCIPBase.ApprovedSenderUpdate[](1);

    senderUpdates[0] =
      CCIPBase.ApprovedSenderUpdate({destChainSelector: s_sourceChainSelector, sender: abi.encode(address(1))});

    s_receiver.updateApprovedSenders(senderUpdates, new CCIPBase.ApprovedSenderUpdate[](0));

    vm.expectEmit();
    emit CCIPReceiver.MessageRecovered(messageId);

    s_receiver.retryFailedMessage(messageId);
    assertFalse(s_receiver.isFailedMessage(messageId), "Message should be marked as resolved");
  }

  function test_retryMessage_RevertWhen_MessageHasNotAlreadyFailed() public {
    bytes32 messageId = keccak256("RANDOM_DATA");

    vm.expectRevert(abi.encodeWithSelector(CCIPReceiver.MessageNotFailed.selector, messageId));

    s_receiver.retryFailedMessage(messageId);
  }
}
