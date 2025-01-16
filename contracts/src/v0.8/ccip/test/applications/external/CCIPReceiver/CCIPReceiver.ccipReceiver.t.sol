// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {Client} from "../../../../libraries/Client.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

contract CCIPReceiver_ccipReceive is CCIPReceiverSetup {
  function test_ccipReceive() public {
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
    emit MessageSucceeded(messageId);

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: sourceChainSelector,
        sender: abi.encode(address(s_receiver)), // correct sender
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );
  }

  function test_ccipReceive_RevertWhen_InvalidChainSelector() public {
    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: sourceChainSelector,
      allowed: false,
      recipient: abi.encode(address(s_receiver)),
      extraArgsBytes: ""
    });

    vm.expectEmit();
    emit CCIPBase.ChainRemoved(sourceChainSelector);

    s_receiver.applyChainUpdates(chainUpdates);

    bytes32 messageId = keccak256("messageId");
    address token = address(s_destFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    deal(token, address(s_receiver), amount);

    // The receiver contract will revert if the router is not the sender.
    vm.startPrank(address(s_destRouter));

    vm.expectRevert(abi.encodeWithSelector(CCIPBase.InvalidChain.selector, sourceChainSelector));

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: sourceChainSelector,
        sender: abi.encode(address(1)),
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );
  }

  function test_ccipReceive_RevertWhen_InvalidSender() public {
    CCIPBase.ApprovedSenderUpdate[] memory senderUpdates = new CCIPBase.ApprovedSenderUpdate[](1);
    senderUpdates[0] =
      CCIPBase.ApprovedSenderUpdate({destChainSelector: sourceChainSelector, sender: abi.encode(address(s_receiver))});

    s_receiver.updateApprovedSenders(new CCIPBase.ApprovedSenderUpdate[](0), senderUpdates);

    // assertFalse(s_receiver.s_approvedSenders(sourceChainSelector, abi.encode(address(s_receiver))));
    assertFalse(s_receiver.isApprovedSender(sourceChainSelector, abi.encode(address(s_receiver))));

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
      messageId, abi.encodeWithSelector(bytes4(CCIPBase.InvalidSender.selector), abi.encode(address(s_receiver)))
    );

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: sourceChainSelector,
        sender: abi.encode(address(s_receiver)),
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );
  }
}
