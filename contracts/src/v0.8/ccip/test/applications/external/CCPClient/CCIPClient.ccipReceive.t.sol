// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../../../libraries/Client.sol";
import {CCIPClientSetup} from "./CCIPClientSetup.t.sol";

contract CCIPClient_ccipReceive is CCIPClientSetup {
  function test_ccipReceive() public {
    bytes32 messageId = keccak256("messageId");
    address token = address(s_destFeeToken);
    uint256 amount = 111333333777;
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    // Make sure we give the receiver contract enough tokens like CCIP would.
    deal(token, address(s_sender), amount);

    // The receiver contract will revert if the router is not the sender.
    vm.startPrank(address(s_sourceRouter));

    vm.expectEmit();
    emit MessageSucceeded(messageId);

    s_sender.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: DEST_CHAIN_SELECTOR,
        sender: abi.encode(address(s_sender)), // correct sender
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );
  }
}
