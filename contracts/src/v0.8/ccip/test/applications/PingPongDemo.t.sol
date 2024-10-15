// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {PingPongDemo} from "../../applications/PingPongDemo.sol";
import {Client} from "../../libraries/Client.sol";
import "../onRamp/OnRampSetup.t.sol";

// setup
contract PingPongDappSetup is OnRampSetup {
  PingPongDemo internal s_pingPong;
  IERC20 internal s_feeToken;

  address internal immutable i_pongContract = makeAddr("ping_pong_counterpart");

  function setUp() public virtual override {
    super.setUp();

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
    _assertPingPongSuccess();
  }

  function test_StartPingPong_With_OOO_Success() public {
    s_pingPong.setOutOfOrderExecution(true);

    _assertPingPongSuccess();
  }

  function _assertPingPongSuccess() internal {
    vm.expectEmit();
    emit PingPongDemo.Ping(pingPongNumber);

    Internal.EVM2AnyRampMessage memory message;

    vm.expectEmit(false, false, false, false);
    emit OnRamp.CCIPMessageSent(DEST_CHAIN_SELECTOR, 1, message);

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
  function test_Fuzz_CounterPartChainSelector_Success(
    uint64 chainSelector
  ) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function test_Fuzz_CounterPartAddress_Success(
    address counterpartAddress
  ) public {
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
