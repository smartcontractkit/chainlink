// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {NonceManager} from "../../NonceManager.sol";
import {BaseTest} from "../BaseTest.t.sol";

contract NonceManager_getIncrementedOutboundNonce is BaseTest {
  NonceManager private s_nonceManager;

  function setUp() public override {
    address[] memory authorizedCallers = new address[](1);
    authorizedCallers[0] = address(this);
    s_nonceManager = new NonceManager(authorizedCallers);
  }

  function test_getIncrementedOutboundNonce() public {
    address sender = address(this);

    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 0);

    uint64 outboundNonce = s_nonceManager.getIncrementedOutboundNonce(DEST_CHAIN_SELECTOR, sender);
    assertEq(outboundNonce, 1);
  }

  function test_incrementInboundNonce() public {
    address sender = address(this);

    s_nonceManager.incrementInboundNonce(SOURCE_CHAIN_SELECTOR, 1, abi.encode(sender));

    assertEq(s_nonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR, abi.encode(sender)), 1);
  }

  function test_incrementInboundNonce_SkippedIncorrectNonce() public {
    address sender = address(this);
    uint64 expectedNonce = 2;

    vm.expectEmit();
    emit NonceManager.SkippedIncorrectNonce(SOURCE_CHAIN_SELECTOR, expectedNonce, abi.encode(sender));

    s_nonceManager.incrementInboundNonce(SOURCE_CHAIN_SELECTOR, expectedNonce, abi.encode(sender));

    assertEq(s_nonceManager.getInboundNonce(SOURCE_CHAIN_SELECTOR, abi.encode(sender)), 0);
  }

  function test_incrementNoncesInboundAndOutbound() public {
    address sender = address(this);

    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 0);
    uint64 outboundNonce = s_nonceManager.getIncrementedOutboundNonce(DEST_CHAIN_SELECTOR, sender);
    assertEq(outboundNonce, 1);

    // Inbound nonce unchanged
    assertEq(s_nonceManager.getInboundNonce(DEST_CHAIN_SELECTOR, abi.encode(sender)), 0);

    s_nonceManager.incrementInboundNonce(DEST_CHAIN_SELECTOR, 1, abi.encode(sender));
    assertEq(s_nonceManager.getInboundNonce(DEST_CHAIN_SELECTOR, abi.encode(sender)), 1);

    // Outbound nonce unchanged
    assertEq(s_nonceManager.getOutboundNonce(DEST_CHAIN_SELECTOR, sender), 1);
  }
}
