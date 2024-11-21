// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Operator} from "../Operator.sol";
import {Callback} from "./testhelpers/Callback.sol";
import {ChainlinkClientHelper} from "./testhelpers/ChainlinkClientHelper.sol";
import {Deployer} from "./testhelpers/Deployer.t.sol";

contract OperatorTest is Deployer {
  ChainlinkClientHelper private s_client;
  Callback private s_callback;
  Operator private s_operator;

  function setUp() public {
    _setUp();
    s_client = new ChainlinkClientHelper(address(s_link));

    address[] memory auth = new address[](1);
    auth[0] = address(this);
    s_operator = new Operator(address(s_link), address(this));
    s_operator.setAuthorizedSenders(auth);

    s_callback = new Callback(address(s_operator));
  }

  function testFuzz_SendRequest_Success(uint96 payment) public {
    vm.assume(payment > 0);
    deal(address(s_link), address(s_client), payment);
    // We're going to cancel one request and fulfill the other
    bytes32 requestIdToCancel = s_client.sendRequest(address(s_operator), payment);

    // Nothing withdrawable
    // 1 payment in escrow
    // Client has zero link
    assertEq(s_operator.withdrawable(), 0);
    assertEq(s_link.balanceOf(address(s_operator)), payment);
    assertEq(s_link.balanceOf(address(s_client)), 0);

    // Advance time so we can cancel
    uint256 expiration = block.timestamp + s_operator.EXPIRYTIME();
    vm.warp(expiration + 1);
    s_client.cancelRequest(requestIdToCancel, payment, expiration);

    // 1 payment has been returned due to the cancellation.
    assertEq(s_operator.withdrawable(), 0);
    assertEq(s_link.balanceOf(address(s_operator)), 0);
    assertEq(s_link.balanceOf(address(s_client)), payment);
  }

  function testFuzz_SendRequestAndCancelRequest_Success(uint96 payment) public {
    vm.assume(payment > 1);
    payment /= payment;

    deal(address(s_link), address(s_client), 2 * payment);

    // Initial state, client has 2 payments, zero in escrow, zero in the operator, zeero withdrawable
    assertEq(s_operator.withdrawable(), 0);
    assertEq(s_link.balanceOf(address(s_operator)), 0);
    assertEq(s_link.balanceOf(address(s_client)), 2 * payment);

    // We're going to cancel one request and fulfill the other
    bytes32 requestId = s_client.sendRequest(address(s_operator), payment);
    bytes32 requestIdToCancel = s_client.sendRequest(address(s_operator), payment);

    // Nothing withdrawable
    // Operator now has the 2 payments in escrow
    // Client has zero payments
    assertEq(s_operator.withdrawable(), 0);
    assertEq(s_link.balanceOf(address(s_operator)), 2 * payment);
    assertEq(s_link.balanceOf(address(s_client)), 0);

    // Fulfill one request
    uint256 expiration = block.timestamp + s_operator.EXPIRYTIME();
    s_operator.fulfillOracleRequest(
      requestId,
      payment,
      address(s_client),
      s_client.FULFILL_SELECTOR(),
      expiration,
      bytes32(hex"01")
    );
    // 1 payment withdrawable from fulfilling `requestId`, 1 payment in escrow
    assertEq(s_operator.withdrawable(), payment);
    assertEq(s_link.balanceOf(address(s_operator)), 2 * payment);
    assertEq(s_link.balanceOf(address(s_client)), 0);

    // Advance time so we can cancel
    vm.warp(expiration + 1);
    s_client.cancelRequest(requestIdToCancel, payment, expiration);

    // 1 payment has been returned due to the cancellation, 1 payment should be withdrawable
    assertEq(s_operator.withdrawable(), payment);
    assertEq(s_link.balanceOf(address(s_operator)), payment);
    assertEq(s_link.balanceOf(address(s_client)), payment);

    // Withdraw the remaining payment
    s_operator.withdraw(address(s_client), payment);

    // End state is exactly the same as the initial state.
    assertEq(s_operator.withdrawable(), 0);
    assertEq(s_link.balanceOf(address(s_operator)), 0);
    assertEq(s_link.balanceOf(address(s_client)), 2 * payment);
  }

  function test_OracleRequest_Success() public {
    // Define some mock values
    bytes32 specId = keccak256("testSpec");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 0;
    uint256 dataVersion = 1;
    bytes memory data = "";

    uint256 initialLinkBalance = s_link.balanceOf(address(s_operator));
    uint256 payment = 1 ether; // Mock payment value

    uint256 withdrawableBefore = s_operator.withdrawable();

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(address(s_link), ALICE, payment);
    assertEq(s_link.balanceOf(ALICE), 1 ether, "balance update failed");

    vm.prank(ALICE);
    s_link.transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(this),
        payment,
        specId,
        address(s_callback),
        callbackFunctionId,
        nonce,
        dataVersion,
        data
      )
    );

    // Check that the LINK tokens were transferred to the Operator contract
    assertEq(s_link.balanceOf(address(s_operator)), initialLinkBalance + payment);
    // No withdrawable LINK as it's all locked
    assertEq(s_operator.withdrawable(), withdrawableBefore);
  }

  function test_FulfillOracleRequest_Success() public {
    // This test file is the callback target and actual sender contract
    // so we should enable it to set Authorised senders to itself
    address[] memory senders = new address[](2);
    senders[0] = address(this);
    senders[0] = BOB;

    s_operator.setAuthorizedSenders(senders);

    uint256 withdrawableBefore = s_operator.withdrawable();

    // Define mock values for creating a new oracle request
    bytes32 specId = keccak256("testSpecForFulfill");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 1;
    uint256 dataVersion = 1;
    bytes memory dataBytes = "";
    uint256 payment = 1 ether;
    uint256 expiration = block.timestamp + 5 minutes;

    // Convert bytes to bytes32
    bytes32 data = bytes32(keccak256(dataBytes));

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(address(s_link), BOB, payment);
    vm.prank(BOB);
    s_link.transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(this),
        payment,
        specId,
        address(s_callback),
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    // Fulfill the request using the operator
    bytes32 requestId = keccak256(abi.encodePacked(BOB, nonce));
    vm.prank(BOB);
    s_operator.fulfillOracleRequest(requestId, payment, address(s_callback), callbackFunctionId, expiration, data);

    assertEq(s_callback.getCallbacksReceived(), 1, "Oracle request was not fulfilled");

    // Withdrawable balance
    assertEq(s_operator.withdrawable(), withdrawableBefore + payment, "Internal accounting not updated correctly");
  }

  function test_CancelOracleRequest_Success() public {
    // Define mock values for creating a new oracle request
    bytes32 specId = keccak256("testSpecForCancel");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 2;
    uint256 dataVersion = 1;
    bytes memory dataBytes = "";
    uint256 payment = 1 ether;
    uint256 expiration = block.timestamp + 5 minutes;

    uint256 withdrawableBefore = s_operator.withdrawable();

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(address(s_link), BOB, payment);
    vm.prank(BOB);
    s_link.transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        BOB,
        payment,
        specId,
        BOB,
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    // No withdrawable balance as it's all locked
    assertEq(s_operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");

    bytes32 requestId = keccak256(abi.encodePacked(BOB, nonce));

    vm.startPrank(ALICE);
    vm.expectRevert(bytes("Params do not match request ID"));
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);

    vm.startPrank(BOB);
    vm.expectRevert(bytes("Request is not expired"));
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);

    vm.warp(expiration);
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);

    // Check if the LINK tokens were refunded to the sender (bob in this case)
    assertEq(s_link.balanceOf(BOB), 1 ether, "Oracle request was not canceled properly");

    assertEq(s_operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");
  }

  function test_NotAuthorizedSender_Revert() public {
    bytes32 specId = keccak256("unauthorizedFulfillSpec");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 5;
    uint256 dataVersion = 1;
    bytes memory dataBytes = "";
    uint256 payment = 1 ether;
    uint256 expiration = block.timestamp + 5 minutes;

    deal(address(s_link), ALICE, payment);
    vm.prank(ALICE);
    s_link.transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        ALICE,
        payment,
        specId,
        address(s_callback),
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    bytes32 requestId = keccak256(abi.encodePacked(ALICE, nonce));

    vm.prank(BOB);
    vm.expectRevert(bytes("Not authorized sender"));
    s_operator.fulfillOracleRequest(
      requestId,
      payment,
      address(s_callback),
      callbackFunctionId,
      expiration,
      bytes32(keccak256(dataBytes))
    );
  }
}
