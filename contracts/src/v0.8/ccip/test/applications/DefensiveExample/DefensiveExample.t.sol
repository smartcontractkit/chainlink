// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {DefensiveExample} from "../../../applications/DefensiveExample.sol";
import {Client} from "../../../libraries/Client.sol";
import {OnRampSetup} from "../../onRamp/OnRamp/OnRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract DefensiveExampleTest is OnRampSetup {
  event MessageFailed(bytes32 indexed messageId, bytes reason);
  event MessageSucceeded(bytes32 indexed messageId);
  event MessageRecovered(bytes32 indexed messageId);

  DefensiveExample internal s_receiver;
  uint64 internal s_sourceChainSelector = 7331;

  function setUp() public virtual override {
    super.setUp();

    s_receiver = new DefensiveExample(s_destRouter, IERC20(s_destFeeToken));
    s_receiver.enableChain(s_sourceChainSelector, abi.encode(""));
  }

  function test_Recovery() public {
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
    emit MessageFailed(messageId, abi.encodeWithSelector(DefensiveExample.ErrorCase.selector));

    s_receiver.ccipReceive(
      Client.Any2EVMMessage({
        messageId: messageId,
        sourceChainSelector: s_sourceChainSelector,
        sender: abi.encode(address(0)), // wrong sender, will revert internally
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );

    address tokenReceiver = address(0x000001337);
    uint256 tokenReceiverBalancePre = IERC20(token).balanceOf(tokenReceiver);
    uint256 receiverBalancePre = IERC20(token).balanceOf(address(s_receiver));

    // Recovery can only be done by the owner.
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit MessageRecovered(messageId);

    s_receiver.retryFailedMessage(messageId, tokenReceiver);

    // Assert the tokens have successfully been rescued from the contract.
    assertEq(IERC20(token).balanceOf(tokenReceiver), tokenReceiverBalancePre + amount);
    assertEq(IERC20(token).balanceOf(address(s_receiver)), receiverBalancePre - amount);
  }

  function test_HappyPath() public {
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
        sourceChainSelector: s_sourceChainSelector,
        sender: abi.encode(address(s_receiver)), // correct sender
        data: "",
        destTokenAmounts: destTokenAmounts
      })
    );
  }
}
