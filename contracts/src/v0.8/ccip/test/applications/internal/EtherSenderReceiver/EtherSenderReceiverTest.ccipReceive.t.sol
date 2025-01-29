// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPReceiverLegacy} from "../../../../applications/internal/CCIPReceiverLegacy.sol";
import {Client} from "../../../../libraries/Client.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_ccipReceive is EtherSenderReceiverTestSetup {
  uint64 internal constant SOURCE_CHAIN_SELECTOR = 424242;
  address internal constant XCHAIN_SENDER = 0x9951529C13B01E542f7eE3b6D6665D292e9BA2E0;

  error InvalidTokenAmounts(uint256 gotAmounts);
  error InvalidToken(address gotToken, address expectedToken);

  function testFuzz_ccipReceive(
    uint256 tokenAmount
  ) public {
    // cap to 10 ether because OWNER only has 10 ether.
    if (tokenAmount > 10 ether) {
      return;
    }

    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: tokenAmount});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), tokenAmount);

    uint256 balanceBefore = OWNER.balance;

    vm.startPrank(ROUTER);
    s_etherSenderReceiver.ccipReceive(message);

    uint256 balanceAfter = OWNER.balance;
    assertEq(balanceAfter, balanceBefore + tokenAmount, "balance must be correct");
  }

  function test_ccipReceive_happyPath() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), AMOUNT);

    uint256 balanceBefore = OWNER.balance;
    s_etherSenderReceiver.publicCcipReceive(message);
    uint256 balanceAfter = OWNER.balance;
    assertEq(balanceAfter, balanceBefore + AMOUNT, "balance must be correct");
  }

  function test_ccipReceive_fallbackToWethTransfer() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(address(s_linkToken)), // ERC20 cannot receive() ether.
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), AMOUNT);

    uint256 balanceBefore = address(s_linkToken).balance;
    s_etherSenderReceiver.publicCcipReceive(message);
    uint256 balanceAfter = address(s_linkToken).balance;
    assertEq(balanceAfter, balanceBefore, "balance must be unchanged");
    uint256 wethBalance = s_weth.balanceOf(address(s_linkToken));
    assertEq(wethBalance, AMOUNT, "weth balance must be correct");
  }

  // Reverts

  function test_RevertWhen_wrongTokenAmount() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](2);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    destTokenAmounts[1] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    vm.expectRevert(abi.encodeWithSelector(InvalidTokenAmounts.selector, uint256(2)));
    s_etherSenderReceiver.publicCcipReceive(message);
  }

  function test_RevertWhen_wrongToken() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_someOtherWeth), amount: AMOUNT});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    vm.expectRevert(abi.encodeWithSelector(InvalidToken.selector, address(s_someOtherWeth), address(s_weth)));
    s_etherSenderReceiver.publicCcipReceive(message);
  }

  function test_RevertWhen_OnlyRouter() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: 1 ether});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    vm.startPrank(OWNER);
    vm.expectRevert(abi.encodeWithSelector(CCIPReceiverLegacy.InvalidRouter.selector, address(OWNER)));

    s_etherSenderReceiver.ccipReceive(message);
  }
}
