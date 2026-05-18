// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

import {CCIPRouter} from "../../applications/EtherSenderReceiver.sol";

import {IRouterClient} from "../../interfaces/IRouterClient.sol";
import {Client} from "../../libraries/Client.sol";
import {WETH9} from "../WETH9.sol";
import {EtherSenderReceiverHelper} from "./../helpers/EtherSenderReceiverHelper.sol";

import {ERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/ERC20.sol";

contract EtherSenderReceiverTest is Test {
  EtherSenderReceiverHelper internal s_etherSenderReceiver;
  WETH9 internal s_weth;
  WETH9 internal s_someOtherWeth;
  ERC20 internal s_linkToken;

  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant ROUTER = 0x0F3779ee3a832D10158073ae2F5e61ac7FBBF880;
  address internal constant XCHAIN_RECEIVER = 0xBd91b2073218AF872BF73b65e2e5950ea356d147;

  function setUp() public {
    vm.startPrank(OWNER);

    s_linkToken = new ERC20("Chainlink Token", "LINK");
    s_someOtherWeth = new WETH9();
    s_weth = new WETH9();
    vm.mockCall(ROUTER, abi.encodeWithSelector(CCIPRouter.getWrappedNative.selector), abi.encode(address(s_weth)));
    s_etherSenderReceiver = new EtherSenderReceiverHelper(ROUTER);

    deal(OWNER, 1_000_000 ether);
    deal(address(s_linkToken), OWNER, 1_000_000 ether);

    // deposit some eth into the weth contract.
    s_weth.deposit{value: 10 ether}();
    uint256 wethSupply = s_weth.totalSupply();
    assertEq(wethSupply, 10 ether, "total weth supply must be 10 ether");
  }
}

contract EtherSenderReceiverTest_constructor is EtherSenderReceiverTest {
  function test_constructor() public view {
    assertEq(s_etherSenderReceiver.getRouter(), ROUTER, "router must be set correctly");
    uint256 allowance = s_weth.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(allowance, type(uint256).max, "allowance must be set infinite");
  }
}

contract EtherSenderReceiverTest_validateFeeToken is EtherSenderReceiverTest {
  uint256 internal constant amount = 100;

  error InsufficientMsgValue(uint256 gotAmount, uint256 msgValue);
  error TokenAmountNotEqualToMsgValue(uint256 gotAmount, uint256 msgValue);

  function test_validateFeeToken_valid_native() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(0),
      extraArgs: ""
    });

    s_etherSenderReceiver.validateFeeToken{value: amount + 1}(message);
  }

  function test_validateFeeToken_valid_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    s_etherSenderReceiver.validateFeeToken{value: amount}(message);
  }

  function test_validateFeeToken_reverts_feeToken_tokenAmountNotEqualToMsgValue() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    vm.expectRevert(abi.encodeWithSelector(TokenAmountNotEqualToMsgValue.selector, amount, amount + 1));
    s_etherSenderReceiver.validateFeeToken{value: amount + 1}(message);
  }
}

contract EtherSenderReceiverTest_validatedMessage is EtherSenderReceiverTest {
  error InvalidDestinationReceiver(bytes destReceiver);
  error InvalidTokenAmounts(uint256 gotAmounts);
  error InvalidWethAddress(address want, address got);
  error GasLimitTooLow(uint256 minLimit, uint256 gotLimit);

  uint256 internal constant amount = 100;

  function test_Fuzz_validatedMessage_msgSenderOverwrite(
    bytes memory data
  ) public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: data,
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_Fuzz_validatedMessage_tokenAddressOverwrite(
    address token
  ) public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_emptyDataOverwrittenToMsgSender() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_dataOverwrittenToMsgSender() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: abi.encode(address(42)),
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_tokenOverwrittenToWeth() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(42), // incorrect token.
      amount: amount
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_validMessage_extraArgs() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
    });
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000}))
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, amount, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(
      validatedMessage.extraArgs,
      Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000})),
      "extraArgs must be correct"
    );
  }

  function test_validatedMessage_invalidTokenAmounts() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(0), amount: amount});
    tokenAmounts[1] = Client.EVMTokenAmount({token: address(0), amount: amount});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    vm.expectRevert(abi.encodeWithSelector(InvalidTokenAmounts.selector, uint256(2)));
    s_etherSenderReceiver.validatedMessage(message);
  }
}

contract EtherSenderReceiverTest_getFee is EtherSenderReceiverTest {
  uint64 internal constant destinationChainSelector = 424242;
  uint256 internal constant feeWei = 121212;
  uint256 internal constant amount = 100;

  function test_getFee() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(0), amount: amount});
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );

    uint256 fee = s_etherSenderReceiver.getFee(destinationChainSelector, message);
    assertEq(fee, feeWei, "fee must be feeWei");
  }
}

contract EtherSenderReceiverTest_ccipReceive is EtherSenderReceiverTest {
  uint256 internal constant amount = 100;
  uint64 internal constant sourceChainSelector = 424242;
  address internal constant XCHAIN_SENDER = 0x9951529C13B01E542f7eE3b6D6665D292e9BA2E0;

  error InvalidTokenAmounts(uint256 gotAmounts);
  error InvalidToken(address gotToken, address expectedToken);

  function test_Fuzz_ccipReceive(
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
      sourceChainSelector: sourceChainSelector,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), tokenAmount);

    uint256 balanceBefore = OWNER.balance;
    s_etherSenderReceiver.publicCcipReceive(message);
    uint256 balanceAfter = OWNER.balance;
    assertEq(balanceAfter, balanceBefore + tokenAmount, "balance must be correct");
  }

  function test_ccipReceive_happyPath() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(OWNER),
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), amount);

    uint256 balanceBefore = OWNER.balance;
    s_etherSenderReceiver.publicCcipReceive(message);
    uint256 balanceAfter = OWNER.balance;
    assertEq(balanceAfter, balanceBefore + amount, "balance must be correct");
  }

  function test_ccipReceive_fallbackToWethTransfer() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: keccak256(abi.encode("ccip send")),
      sourceChainSelector: 424242,
      sender: abi.encode(XCHAIN_SENDER),
      data: abi.encode(address(s_linkToken)), // ERC20 cannot receive() ether.
      destTokenAmounts: destTokenAmounts
    });

    // simulate a cross-chain token transfer, just transfer the weth to s_etherSenderReceiver.
    s_weth.transfer(address(s_etherSenderReceiver), amount);

    uint256 balanceBefore = address(s_linkToken).balance;
    s_etherSenderReceiver.publicCcipReceive(message);
    uint256 balanceAfter = address(s_linkToken).balance;
    assertEq(balanceAfter, balanceBefore, "balance must be unchanged");
    uint256 wethBalance = s_weth.balanceOf(address(s_linkToken));
    assertEq(wethBalance, amount, "weth balance must be correct");
  }

  function test_ccipReceive_wrongTokenAmount() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](2);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
    destTokenAmounts[1] = Client.EVMTokenAmount({token: address(s_weth), amount: amount});
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

  function test_ccipReceive_wrongToken() public {
    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](1);
    destTokenAmounts[0] = Client.EVMTokenAmount({token: address(s_someOtherWeth), amount: amount});
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
}

contract EtherSenderReceiverTest_ccipSend is EtherSenderReceiverTest {
  error InsufficientFee(uint256 gotFee, uint256 fee);

  uint256 internal constant amount = 100;
  uint64 internal constant destinationChainSelector = 424242;
  uint256 internal constant feeWei = 121212;
  uint256 internal constant feeJuels = 232323;

  function test_Fuzz_ccipSend(uint256 feeFromRouter, uint256 feeSupplied) public {
    // cap the fuzzer because OWNER only has a million ether.
    vm.assume(feeSupplied < 1_000_000 ether - amount);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeFromRouter)
    );

    if (feeSupplied < feeFromRouter) {
      vm.expectRevert();
      s_etherSenderReceiver.ccipSend{value: amount + feeSupplied}(destinationChainSelector, message);
    } else {
      bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
      vm.mockCall(
        ROUTER,
        feeSupplied,
        abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
        abi.encode(expectedMsgId)
      );

      bytes32 actualMsgId =
        s_etherSenderReceiver.ccipSend{value: amount + feeSupplied}(destinationChainSelector, message);
      assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    }
  }

  function test_Fuzz_ccipSend_feeToken(uint256 feeFromRouter, uint256 feeSupplied) public {
    // cap the fuzzer because OWNER only has a million LINK.
    vm.assume(feeSupplied < 1_000_000 ether - amount);

    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeFromRouter)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), feeSupplied);

    if (feeSupplied < feeFromRouter) {
      vm.expectRevert();
      s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
    } else {
      bytes32 expectedMsgId = keccak256(abi.encode("ccip send"));
      vm.mockCall(
        ROUTER,
        abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
        abi.encode(expectedMsgId)
      );

      bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
      assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    }
  }

  function test_ccipSend_reverts_insufficientFee_weth() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );

    s_weth.approve(address(s_etherSenderReceiver), feeWei - 1);

    vm.expectRevert("SafeERC20: low-level call failed");
    s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
  }

  function test_ccipSend_reverts_insufficientFee_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeJuels)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), feeJuels - 1);

    vm.expectRevert("ERC20: insufficient allowance");
    s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
  }

  function test_ccipSend_reverts_insufficientFee_native() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );

    vm.expectRevert();
    s_etherSenderReceiver.ccipSend{value: amount + feeWei - 1}(destinationChainSelector, message);
  }

  function test_ccipSend_success_nativeExcess() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );

    // we assert that the correct value is sent to the router call, which should be
    // the msg.value - feeWei.
    vm.mockCall(
      ROUTER,
      feeWei + 1,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
      abi.encode(expectedMsgId)
    );

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: amount + feeWei + 1}(destinationChainSelector, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
  }

  function test_ccipSend_success_native() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );
    vm.mockCall(
      ROUTER,
      feeWei,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
      abi.encode(expectedMsgId)
    );

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: amount + feeWei}(destinationChainSelector, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
  }

  function test_ccipSend_success_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeJuels)
    );
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
      abi.encode(expectedMsgId)
    );

    s_linkToken.approve(address(s_etherSenderReceiver), feeJuels);

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    uint256 routerAllowance = s_linkToken.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(routerAllowance, feeJuels, "router allowance must be feeJuels");
  }

  function test_ccipSend_success_weth() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: amount
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
      abi.encodeWithSelector(IRouterClient.getFee.selector, destinationChainSelector, validatedMessage),
      abi.encode(feeWei)
    );
    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.ccipSend.selector, destinationChainSelector, validatedMessage),
      abi.encode(expectedMsgId)
    );

    s_weth.approve(address(s_etherSenderReceiver), feeWei);

    bytes32 actualMsgId = s_etherSenderReceiver.ccipSend{value: amount}(destinationChainSelector, message);
    assertEq(actualMsgId, expectedMsgId, "message id must be correct");
    uint256 routerAllowance = s_weth.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(routerAllowance, type(uint256).max, "router allowance must be max for weth");
  }
}
