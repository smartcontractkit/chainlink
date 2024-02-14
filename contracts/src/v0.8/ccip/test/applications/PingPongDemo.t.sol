// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import "../../applications/PingPongDemo.sol";
import "../onRamp/EVM2EVMOnRampSetup.t.sol";
import "../../libraries/Client.sol";

// setup
contract PingPongDappSetup is EVM2EVMOnRampSetup {
  event Ping(uint256 pingPongs);
  event Pong(uint256 pingPongs);

  PingPongDemo internal s_pingPong;
  IERC20 internal s_feeToken;

  address internal immutable i_pongContract = address(10);

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

/// @notice #startPingPong
contract PingPong_startPingPong is PingPongDappSetup {
  event ConfigPropagated(uint64 chainSelector, address contractAddress);

  function testStartPingPongSuccess() public {
    uint256 pingPongNumber = 1;
    bytes memory data = abi.encode(pingPongNumber);

    Client.EVM2AnyMessage memory sentMessage = Client.EVM2AnyMessage({
      receiver: abi.encode(i_pongContract),
      data: data,
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 2e5}))
    });

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, sentMessage);

    Internal.EVM2EVMMessage memory message = Internal.EVM2EVMMessage({
      sequenceNumber: 1,
      feeTokenAmount: expectedFee,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      sender: address(s_pingPong),
      receiver: i_pongContract,
      nonce: 1,
      data: data,
      tokenAmounts: sentMessage.tokenAmounts,
      sourceTokenData: new bytes[](sentMessage.tokenAmounts.length),
      gasLimit: 2e5,
      feeToken: sentMessage.feeToken,
      strict: false,
      messageId: ""
    });
    message.messageId = Internal._hash(message, s_metadataHash);

    vm.expectEmit();
    emit Ping(pingPongNumber);

    vm.expectEmit();
    emit CCIPSendRequested(message);

    s_pingPong.startPingPong();
  }
}

/// @notice #ccipReceive
contract PingPong_ccipReceive is PingPongDappSetup {
  function testCcipReceiveSuccess() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](0);

    uint256 pingPongNumber = 5;

    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: DEST_CHAIN_SELECTOR,
      sender: abi.encode(i_pongContract),
      data: abi.encode(pingPongNumber),
      destTokenAmounts: tokenAmounts
    });

    changePrank(address(s_sourceRouter));

    vm.expectEmit();
    emit Pong(pingPongNumber + 1);

    s_pingPong.ccipReceive(message);
  }
}

contract PingPong_plumbing is PingPongDappSetup {
  function testFuzz_CounterPartChainSelectorSuccess(uint64 chainSelector) public {
    s_pingPong.setCounterpartChainSelector(chainSelector);

    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function testFuzz_CounterPartAddressSuccess(address counterpartAddress) public {
    s_pingPong.setCounterpartAddress(counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
  }

  function testFuzz_CounterPartAddressSuccess(uint64 chainSelector, address counterpartAddress) public {
    s_pingPong.setCounterpart(chainSelector, counterpartAddress);

    assertEq(s_pingPong.getCounterpartAddress(), counterpartAddress);
    assertEq(s_pingPong.getCounterpartChainSelector(), chainSelector);
  }

  function testPausingSuccess() public {
    assertFalse(s_pingPong.isPaused());

    s_pingPong.setPaused(true);

    assertTrue(s_pingPong.isPaused());
  }
}
