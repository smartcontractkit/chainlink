// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {PingPongDemo} from "../../applications/PingPongDemo.sol";
import {Client} from "../../libraries/Client.sol";
import "../onRamp/EVM2EVMOnRampSetup.t.sol";

// setup
contract PingPongDappSetup is EVM2EVMOnRampSetup {
  PingPongDemo internal s_pingPong;
  IERC20 internal s_feeToken;

  address internal immutable i_pongContract = makeAddr("ping_pong_counterpart");

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    s_feeToken = IERC20(s_sourceTokens[0]);
    s_pingPong = new PingPongDemo(address(s_sourceRouter), s_feeToken);
    s_pingPong.setCounterpart(DEST_CHAIN_SELECTOR, i_pongContract);

    uint256 fundingAmount = 1e18;

    // Fund the contract with LINK tokens
    s_feeToken.transfer(address(s_pingPong), fundingAmount);
  }
}

contract PingPong_startPingPong is PingPongDappSetup {
  uint256 internal pingPongNumber = 1;

  function test_StartPingPong_With_Sequenced_Ordered_Success() public {
    Client.EVM2AnyMessage memory sentMessage = Client.EVM2AnyMessage({
      receiver: abi.encode(i_pongContract),
      data: abi.encode(pingPongNumber),
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 200_000}))
    });

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, sentMessage);

    Internal.EVM2EVMMessage memory message = Internal.EVM2EVMMessage({
      sequenceNumber: 1,
      feeTokenAmount: expectedFee,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      sender: address(s_pingPong),
      receiver: i_pongContract,
      nonce: 1,
      data: abi.encode(pingPongNumber),
      tokenAmounts: sentMessage.tokenAmounts,
      sourceTokenData: new bytes[](sentMessage.tokenAmounts.length),
      gasLimit: 200_000,
      feeToken: sentMessage.feeToken,
      strict: false,
      messageId: ""
    });

    _assertPingPongSuccess(message);
  }

  function test_StartPingPong_With_OOO_Success() public {
    s_pingPong.setOutOfOrderExecution(true);

    Client.EVM2AnyMessage memory sentMessage = Client.EVM2AnyMessage({
      receiver: abi.encode(i_pongContract),
      data: abi.encode(pingPongNumber),
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV2({gasLimit: 200_000, allowOutOfOrderExecution: true}))
    });

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, sentMessage);

    Internal.EVM2EVMMessage memory message = Internal.EVM2EVMMessage({
      sequenceNumber: 1,
      feeTokenAmount: expectedFee,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      sender: address(s_pingPong),
      receiver: i_pongContract,
      nonce: 0,
      data: abi.encode(pingPongNumber),
      tokenAmounts: sentMessage.tokenAmounts,
      sourceTokenData: new bytes[](sentMessage.tokenAmounts.length),
      gasLimit: 200_000,
      feeToken: sentMessage.feeToken,
      strict: false,
      messageId: ""
    });

    _assertPingPongSuccess(message);
  }

  function _assertPingPongSuccess(Internal.EVM2EVMMessage memory message) internal {
    message.messageId = Internal._hash(message, s_metadataHash);

    vm.expectEmit();
    emit PingPongDemo.Ping(pingPongNumber);

    vm.expectEmit();
    emit EVM2EVMOnRamp.CCIPSendRequested(message);

    s_pingPong.startPingPong();
  }
}

contract PingPong_ccipReceive is PingPongDappSetup {
  function test_CcipReceive_Success() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](0);

    uint256 pingPongNumber = 5;

    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: DEST_CHAIN_SELECTOR,
      sender: abi.encode(i_pongContract),
      data: abi.encode(pingPongNumber),
      destTokenAmounts: tokenAmounts
    });

    vm.startPrank(address(s_sourceRouter));

    vm.expectEmit();
    emit PingPongDemo.Pong(pingPongNumber + 1);

    s_pingPong.ccipReceive(message);
  }
}

contract PingPong_plumbing is PingPongDappSetup {
  function test_Fuzz_CounterPartChainSelector_Success(uint64 chainSelector) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_Fuzz_CounterPartAddress_Success(address counterpartAddress) public {
    s_pingPong.setCounterpartAddress(counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
  }

  function test_Fuzz_CounterPartAddress_Success(uint64 chainSelector, address counterpartAddress) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    s_pingPong.setCounterpart(chainSelector, counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_Pausing_Success() public {
    assertFalse(s_pingPong.isPaused());

    s_pingPong.setPaused(true);

    assertTrue(s_pingPong.isPaused());
  }

  function test_OutOfOrderExecution_Success() public {
    assertFalse(s_pingPong.getOutOfOrderExecution());

    vm.expectEmit();
    emit PingPongDemo.OutOfOrderExecutionChange(true);

    s_pingPong.setOutOfOrderExecution(true);

    assertTrue(s_pingPong.getOutOfOrderExecution());
  }
}
