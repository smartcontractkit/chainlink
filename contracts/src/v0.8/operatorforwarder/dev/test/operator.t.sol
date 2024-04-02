// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {Operator} from "../Operator.sol";
import {ChainlinkClientHelper} from "./testhelpers/ChainlinkClientHelper.sol";
import {LinkToken} from "../../../shared/token/ERC677/LinkToken.sol";

contract Operator_cancelRequest is Test {
  address public s_link;
  ChainlinkClientHelper public s_client;
  Operator public s_operator;

  function setUp() public {
    s_link = address(new LinkToken());
    s_client = new ChainlinkClientHelper(s_link);

    address[] memory auth = new address[](1);
    auth[0] = address(this);
    s_operator = new Operator(s_link, address(this));
    s_operator.setAuthorizedSenders(auth);
  }

  function test_Success(uint96 payment) public {
    payment = uint96(bound(payment, 1, type(uint96).max));
    deal(s_link, address(s_client), payment);
    // We're going to cancel one request and fulfil the other
    bytes32 requestIdToCancel = s_client.sendRequest(address(s_operator), payment);

    // Nothing withdrawable
    // 1 payment in escrow
    // Client has zero link
    assertEq(s_operator.withdrawable(), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), 0);

    // Advance time so we can cancel
    uint256 expiration = block.timestamp + s_operator.EXPIRYTIME();
    vm.warp(expiration + 1);
    s_client.cancelRequest(requestIdToCancel, payment, expiration);

    // 1 payment has been returned due to the cancellation.
    assertEq(s_operator.withdrawable(), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), payment);
  }

  function test_afterSuccessfulRequestSucess(uint96 payment) public {
    payment = uint96(bound(payment, 1, type(uint96).max) / 2);
    deal(s_link, address(s_client), 2 * payment);

    // Initial state, client has 2 payments, zero in escrow, zero in the operator, zeero withdrawable
    assertEq(s_operator.withdrawable(), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), 2 * payment);

    // We're going to cancel one request and fulfil the other
    bytes32 requestId = s_client.sendRequest(address(s_operator), payment);
    bytes32 requestIdToCancel = s_client.sendRequest(address(s_operator), payment);

    // Nothing withdrawable
    // Operator now has the 2 payments in escrow
    // Client has zero payments
    assertEq(s_operator.withdrawable(), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), 2 * payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), 0);

    // Fulfill one request
    uint256 expiration = block.timestamp + s_operator.EXPIRYTIME();
    s_operator.fulfillOracleRequest(
      requestId,
      payment,
      address(s_client),
      s_client.FULFILSELECTOR(),
      expiration,
      bytes32(hex"01")
    );
    // 1 payment withdrawable from fulfilling `requestId`, 1 payment in escrow
    assertEq(s_operator.withdrawable(), payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), 2 * payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), 0);

    // Advance time so we can cancel
    vm.warp(expiration + 1);
    s_client.cancelRequest(requestIdToCancel, payment, expiration);

    // 1 payment has been returned due to the cancellation, 1 payment should be withdrawable
    assertEq(s_operator.withdrawable(), payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), payment);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), payment);

    // Withdraw the remaining payment
    s_operator.withdraw(address(s_client), payment);

    // End state is exactly the same as the initial state.
    assertEq(s_operator.withdrawable(), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), 0);
    assertEq(LinkToken(s_link).balanceOf(address(s_client)), 2 * payment);
  }
}
