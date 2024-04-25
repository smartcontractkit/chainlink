// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {Operator} from "../Operator.sol";
import {ChainlinkClientHelper} from "./testhelpers/ChainlinkClientHelper.sol";
import {LinkToken} from "../../shared/token/ERC677/LinkToken.sol";

import "./testhelpers/Deployer.sol";
import "../AuthorizedReceiver.sol";

contract OperatorTest is Deployer, AuthorizedReceiver {
  address public s_link;
  uint256 private dataReceived;
  ChainlinkClientHelper public s_client;
  Operator public s_operator;

  function setUp() public {
    _setUp();
    s_link = address(new LinkToken());
    s_client = new ChainlinkClientHelper(s_link);

    address[] memory auth = new address[](1);
    auth[0] = address(this);
    s_operator = new Operator(s_link, address(this));
    s_operator.setAuthorizedSenders(auth);

    dataReceived = 0;
  }

  // Callback function for oracle request fulfillment
  function callback(bytes32 _requestId) public {
    require(msg.sender == address(s_operator), "Only Operator can call this function");
    dataReceived += 1;
  }

  // @notice concrete implementation of AuthorizedReceiver
  // @return bool of whether sender is authorized
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return true;
  }

  function test_Success(uint96 payment) public {
    payment = uint96(bound(payment, 1, type(uint96).max));
    deal(s_link, address(s_client), payment);
    // We're going to cancel one request and fulfill the other
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

  function test_oracleRequestFlow() public {
    // Define some mock values
    bytes32 specId = keccak256("testSpec");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 0;
    uint256 dataVersion = 1;
    bytes memory data = "";

    uint256 initialLinkBalance = LinkToken(s_link).balanceOf(address(s_operator));
    uint256 payment = 1 ether; // Mock payment value

    uint256 withdrawableBefore = s_operator.withdrawable();

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(s_link, address(alice), payment);
    assertEq(LinkToken(s_link).balanceOf(address(alice)), 1 ether, "balance update failed");

    vm.prank(alice);
    LinkToken(s_link).transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(this),
        payment,
        specId,
        address(this),
        callbackFunctionId,
        nonce,
        dataVersion,
        data
      )
    );

    // Check that the LINK tokens were transferred to the Operator contract
    assertEq(LinkToken(s_link).balanceOf(address(s_operator)), initialLinkBalance + payment);
    // No withdrawable LINK as it's all locked
    assertEq(s_operator.withdrawable(), withdrawableBefore);
  }

  function test_fulfillOracleRequest() public {
    // This test file is the callback target and actual sender contract
    // so we should enable it to set Authorised senders to itself
    address[] memory senders = new address[](2);
    senders[0] = address(this);
    senders[0] = address(bob);

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
    bytes32 data = bytes32(uint256(keccak256(dataBytes)));

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(s_link, address(bob), payment);
    vm.prank(bob);
    LinkToken(s_link).transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(this),
        payment,
        specId,
        address(this),
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    // Fulfill the request using the operator
    bytes32 requestId = keccak256(abi.encodePacked(bob, nonce));
    vm.prank(bob);
    s_operator.fulfillOracleRequest(requestId, payment, address(this), callbackFunctionId, expiration, data);

    assertEq(dataReceived, 1, "Oracle request was not fulfilled");

    // Withdrawable balance
    assertEq(s_operator.withdrawable(), withdrawableBefore + payment, "Internal accounting not updated correctly");
  }

  function test_cancelOracleRequest() public {
    // Define mock values for creating a new oracle request
    bytes32 specId = keccak256("testSpecForCancel");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 2;
    uint256 dataVersion = 1;
    bytes memory dataBytes = "";
    uint256 payment = 1 ether;
    uint256 expiration = block.timestamp + 5 minutes;

    uint256 withdrawableBefore = s_operator.withdrawable();

    // Convert bytes to bytes32
    bytes32 data = bytes32(uint256(keccak256(dataBytes)));

    // Send LINK tokens to the Operator contract using `transferAndCall`
    deal(s_link, address(bob), payment);
    vm.prank(bob);
    LinkToken(s_link).transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(bob),
        payment,
        specId,
        address(bob),
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    // No withdrawable balance as it's all locked
    assertEq(s_operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");

    bytes32 requestId = keccak256(abi.encodePacked(bob, nonce));

    vm.startPrank(alice);
    vm.expectRevert(bytes("Params do not match request ID"));
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);
    vm.stopPrank();

    vm.startPrank(bob);
    vm.expectRevert(bytes("Request is not expired"));
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);

    vm.warp(expiration);
    s_operator.cancelOracleRequest(requestId, payment, callbackFunctionId, expiration);
    vm.stopPrank();

    // Check if the LINK tokens were refunded to the sender (bob in this case)
    assertEq(LinkToken(s_link).balanceOf(address(bob)), 1 ether, "Oracle request was not canceled properly");

    assertEq(s_operator.withdrawable(), withdrawableBefore, "Internal accounting not updated correctly");
  }

  function test_unauthorizedFulfillment() public {
    bytes32 specId = keccak256("unauthorizedFulfillSpec");
    bytes4 callbackFunctionId = bytes4(keccak256("callback(bytes32)"));
    uint256 nonce = 5;
    uint256 dataVersion = 1;
    bytes memory dataBytes = "";
    uint256 payment = 1 ether;
    uint256 expiration = block.timestamp + 5 minutes;

    deal(s_link, address(alice), payment);
    vm.prank(alice);
    LinkToken(s_link).transferAndCall(
      address(s_operator),
      payment,
      abi.encodeWithSignature(
        "oracleRequest(address,uint256,bytes32,address,bytes4,uint256,uint256,bytes)",
        address(alice),
        payment,
        specId,
        address(this),
        callbackFunctionId,
        nonce,
        dataVersion,
        dataBytes
      )
    );

    bytes32 requestId = keccak256(abi.encodePacked(alice, nonce));

    vm.prank(address(bob));
    vm.expectRevert(bytes("Not authorized sender"));
    s_operator.fulfillOracleRequest(
      requestId,
      payment,
      address(this),
      callbackFunctionId,
      expiration,
      bytes32(uint256(keccak256(dataBytes)))
    );
  }
}
