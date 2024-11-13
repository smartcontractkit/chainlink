// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../../libraries/Client.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_validatedMessage is EtherSenderReceiverTestSetup {
  error InvalidDestinationReceiver(bytes destReceiver);
  error InvalidTokenAmounts(uint256 gotAmounts);
  error InvalidWethAddress(address want, address got);
  error GasLimitTooLow(uint256 minLimit, uint256 gotLimit);

  function testFuzz_validatedMessage_msgSenderOverwrite(
    bytes memory data
  ) public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
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
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function testFuzz_validatedMessage_tokenAddressOverwrite(
    address token
  ) public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: AMOUNT});
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
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_emptyDataOverwrittenToMsgSender() public view {
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
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_dataOverwrittenToMsgSender() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(0), // callers may not specify this.
      amount: AMOUNT
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
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_tokenOverwrittenToWeth() public view {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({
      token: address(42), // incorrect token.
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
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(validatedMessage.extraArgs, bytes(""), "extraArgs must be empty");
  }

  function test_validatedMessage_validMessage_extraArgs() public view {
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
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000}))
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);
    assertEq(validatedMessage.receiver, abi.encode(XCHAIN_RECEIVER), "receiver must be XCHAIN_RECEIVER");
    assertEq(validatedMessage.data, abi.encode(OWNER), "data must be msg.sender");
    assertEq(validatedMessage.tokenAmounts[0].token, address(s_weth), "token must be weth");
    assertEq(validatedMessage.tokenAmounts[0].amount, AMOUNT, "amount must be correct");
    assertEq(validatedMessage.feeToken, address(0), "feeToken must be 0");
    assertEq(
      validatedMessage.extraArgs,
      Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000})),
      "extraArgs must be correct"
    );
  }

  function test_validatedMessage_invalidTokenAmounts() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](2);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(0), amount: AMOUNT});
    tokenAmounts[1] = Client.EVMTokenAmount({token: address(0), amount: AMOUNT});
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
