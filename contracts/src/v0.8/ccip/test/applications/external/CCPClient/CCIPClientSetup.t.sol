// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPClient} from "../../../../applications/external/CCIPClient.sol";
import {Client} from "../../../../libraries/Client.sol";
import {OnRampSetup} from "../../../onRamp/OnRamp/OnRampSetup.t.sol";

import {IERC20} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract CCIPClientSetup is OnRampSetup {
  event MessageFailed(bytes32 indexed messageId, bytes reason);
  event MessageSucceeded(bytes32 indexed messageId);
  event MessageRecovered(bytes32 indexed messageId);

  CCIPClient internal s_sender;

  function setUp() public virtual override {
    OnRampSetup.setUp();

    s_sender = new CCIPClient(address(s_sourceRouter), IERC20(s_sourceFeeToken), false);

    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      recipient: abi.encode(address(s_sender)),
      extraArgsBytes: ""
    });
    s_sender.applyChainUpdates(chainUpdates);

    CCIPBase.ApprovedSenderUpdate[] memory senderUpdates = new CCIPBase.ApprovedSenderUpdate[](1);
    senderUpdates[0] =
      CCIPBase.ApprovedSenderUpdate({destChainSelector: DEST_CHAIN_SELECTOR, sender: abi.encode(address(s_sender))});

    s_sender.updateApprovedSenders(senderUpdates, new CCIPBase.ApprovedSenderUpdate[](0));
  }

  function test_HappyPath_Success() public {
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
