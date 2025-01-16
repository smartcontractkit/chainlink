// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPReceiver} from "../../../../applications/external/CCIPReceiver.sol";
import {Client} from "../../../../libraries/Client.sol";
import {CCIPReceiverReverting} from "../../../helpers/receivers/CCIPReceiverReverting.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPReceiver_abandonFailedMessage is CCIPReceiverSetup {
  function test_abandonFailedMessage() public {
    bytes32 messageId = keccak256("messageId");
    address token = address(s_destFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    deal(token, address(s_receiver), amount);

    // Make sure the contract call reverts so we can test recovery.
    s_receiver.setSimRevert(true);

    // The receiver contract will revert if the router is not the sender.
    vm.startPrank(address(s_destRouter));

    vm.expectEmit();
    emit MessageFailed(messageId, abi.encodeWithSelector(CCIPReceiverReverting.ErrorCase.selector));

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: sourceChainSelector,
        sender: abi.encode(address(s_receiver)),
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );

    address tokenReceiver = OWNER;
    uint256 tokenReceiverBalancePre = IERC20(token).balanceOf(tokenReceiver);
    uint256 receiverBalancePre = IERC20(token).balanceOf(address(s_receiver));

    // Recovery can only be done by the owner.
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit CCIPReceiver.MessageAbandoned(messageId, OWNER);

    s_receiver.abandonFailedMessage(messageId, OWNER);

    // Assert the tokens have successfully been rescued from the contract.
    assertEq(
      IERC20(token).balanceOf(tokenReceiver), tokenReceiverBalancePre + amount, "tokens not sent to tokenReceiver"
    );
    assertEq(
      IERC20(token).balanceOf(address(s_receiver)), receiverBalancePre - amount, "tokens not subtracted from receiver"
    );
  }

  function test_abandonFailedMessage_RevertWhen_MessageNotFailed() public {
    bytes32 messageId = keccak256("RANDOM_DATA");

    vm.expectRevert(abi.encodeWithSelector(CCIPReceiver.MessageNotFailed.selector, messageId));

    s_receiver.abandonFailedMessage(messageId, OWNER);
  }
}
