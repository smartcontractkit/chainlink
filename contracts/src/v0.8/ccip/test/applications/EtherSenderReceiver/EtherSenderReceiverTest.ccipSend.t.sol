// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../../interfaces/IRouterClient.sol";
import {Client} from "../../../libraries/Client.sol";

import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_ccipSend is EtherSenderReceiverTestSetup {
  error InsufficientFee(uint256 gotFee, uint256 fee);

  uint64 internal constant DESTINATION_CHAIN_SELECTOR = 424242;
  uint256 internal constant FEE_WEI = 121212;
  uint256 internal constant FEE_JUELS = 232323;

  function testFuzz_ccipSend(uint256 feeFromRouter, uint256 feeSupplied) public {
    // cap the fuzzer because OWNER only has a million ether.
    vm.assume(feeSupplied < 1_000_000 ether - AMOUNT);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(feeFromRouter)
    );

    if (feeSupplied < feeFromRouter) {
      vm.expectRevert();
      s_etherSenderReceiver.ccipSend{value: AMOUNT + feeSupplied}(DESTINATION_CHAIN_SELECTOR, message);
    } else {
      bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
      vm.mockCall(
        ROUTER,
        feeSupplied,
        abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
        abi.encode(expectedMsgId)
      );

      bytes32 actualMsgId =
        s_etherSenderReceiver.ccipSend{value: AMOUNT + feeSupplied}(DESTINATION_CHAIN_SELECTOR, message);
      assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    }
  }

  function testFuzz_ccipSend_feeToken(uint256 feeFromRouter, uint256 feeSupplied) public {
    // cap the fuzzer because OWNER only has a million LINK.
    vm.assume(feeSupplied < 1_000_000 ether - AMOUNT);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(s_linkToken),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(feeFromRouter)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), feeSupplied);

    if (feeSupplied < feeFromRouter) {
      vm.expectRevert();
      s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
    } else {
      bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
      vm.mockCall(
        ROUTER,
        abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
        abi.encode(expectedMsgId)
      );

      bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
      assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    }
  }

  function test_RevertWhen_ccipSends_insufficientFee_weth() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );

    s_weth.approve(address(s_etherSenderReceiver), FEE_WEI - 1);

    vm.expectRevert("SafeERC20: low-level call failed");
    s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
  }

  function test_RevertWhen_ccipSends_insufficientFee_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(s_linkToken),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_JUELS)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), FEE_JUELS - 1);

    vm.expectRevert("ERC20: insufficient allowance");
    s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
  }

  function test_RevertWhen_ccipSends_insufficientFee_native() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );

    vm.expectRevert();
    s_etherSenderReceiver.ccipSend{value: AMOUNT + FEE_WEI - 1}(DESTINATION_CHAIN_SELECTOR, message);
  }

  function test_ccipSend_nativeExcess() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );

    // we assert that the correct value is sent to the router call, which should be
    // the msg.value - feeWei.
    vm.mockCall(
      ROUTER,
      FEE_WEI + 1,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(expectedMsgId)
    );

    bytes32 actualMsgId =
      s_etherSenderReceiver.ccipSend{value: AMOUNT + FEE_WEI + 1}(DESTINATION_CHAIN_SELECTOR, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
  }

  function test_ccipSend_native() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );
    vm.mockCall(
      ROUTER,
      FEE_WEI,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(expectedMsgId)
    );

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: AMOUNT + FEE_WEI}(DESTINATION_CHAIN_SELECTOR, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
  }

  function test_ccipSend_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(s_linkToken),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_JUELS)
    );
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(expectedMsgId)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), FEE_JUELS);

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    uint256 routerAllowance = s_linkToken.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(routerAllowance, FEE_JUELS, "router allowance must be feeJuels");
  }

  function test_ccipSend_weth() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(expectedMsgId)
    );

    s_weth.approve(address(s_etherSenderReceiver), FEE_WEI);

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: AMOUNT}(DESTINATION_CHAIN_SELECTOR, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    uint256 routerAllowance = s_weth.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(routerAllowance, type(uint256).max, "router allowance must be max for weth");
  }
}
