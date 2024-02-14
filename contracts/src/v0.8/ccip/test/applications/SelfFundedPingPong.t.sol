// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {SelfFundedPingPong} from "../../applications/SelfFundedPingPong.sol";
import {EVM2EVMOnRampSetup} from "../onRamp/EVM2EVMOnRampSetup.t.sol";
import {EVM2EVMOnRamp} from "../../onRamp/EVM2EVMOnRamp.sol";
import {Client} from "../../libraries/Client.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract SelfFundedPingPongDappSetup is EVM2EVMOnRampSetup {
  event Ping(uint256 pingPongs);
  event Pong(uint256 pingPongs);
  event CountIncrBeforeFundingSet(uint8 countIncrBeforeFunding);

  SelfFundedPingPong internal s_pingPong;
  IERC20 internal s_feeToken;
  uint8 internal constant s_roundTripsBeforeFunding = 0;

  address internal immutable i_pongContract = address(10);

  function setUp() public virtual override {
    EVM2EVMOnRampSetup.setUp();

    s_feeToken = IERC20(s_sourceTokens[0]);
    s_pingPong = new SelfFundedPingPong(address(s_sourceRouter), s_feeToken, s_roundTripsBeforeFunding);
    s_pingPong.setCounterpart(DEST_CHAIN_SELECTOR, i_pongContract);

    uint256 fundingAmount = 5e18;

    // set ping pong as an onRamp nop to make sure that funding runs
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](1);
    nopsAndWeights[0] = EVM2EVMOnRamp.NopAndWeight({nop: address(s_pingPong), weight: 1});
    s_onRamp.setNops(nopsAndWeights);

    // Fund the contract with LINK tokens
    s_feeToken.transfer(address(s_pingPong), fundingAmount);
  }
}

/// @notice #ccipReceive
contract SelfFundedPingPong_ccipReceive is SelfFundedPingPongDappSetup {
  event Funded();

  function test_FundingSuccess() public {
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: DEST_CHAIN_SELECTOR,
      sender: abi.encode(i_pongContract),
      data: "",
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });

    uint8 countIncrBeforeFunding = 5;

    vm.expectEmit();
    emit CountIncrBeforeFundingSet(countIncrBeforeFunding);

    s_pingPong.setCountIncrBeforeFunding(countIncrBeforeFunding);

    vm.startPrank(address(s_sourceRouter));
    for (uint256 pingPongNumber = 0; pingPongNumber <= countIncrBeforeFunding; ++pingPongNumber) {
      message.data = abi.encode(pingPongNumber);
      if (pingPongNumber == countIncrBeforeFunding - 1) {
        vm.expectEmit();
        emit Funded();
        vm.expectCall(address(s_onRamp), "");
      }
      s_pingPong.ccipReceive(message);
    }
  }

  function test_FundingIfNotANopReverts() public {
    EVM2EVMOnRamp.NopAndWeight[] memory nopsAndWeights = new EVM2EVMOnRamp.NopAndWeight[](0);
    s_onRamp.setNops(nopsAndWeights);

    uint8 countIncrBeforeFunding = 3;
    s_pingPong.setCountIncrBeforeFunding(countIncrBeforeFunding);

    vm.startPrank(address(s_sourceRouter));
    Client.Any2EVMMessage memory message = Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: DEST_CHAIN_SELECTOR,
      sender: abi.encode(i_pongContract),
      data: abi.encode(countIncrBeforeFunding),
      destTokenAmounts: new Client.EVMTokenAmount[](0)
    });

    // because pingPong is not set as a nop
    vm.expectRevert(EVM2EVMOnRamp.OnlyCallableByOwnerOrAdminOrNop.selector);
    s_pingPong.ccipReceive(message);
  }
}

/// @notice #setCountIncrBeforeFunding
contract SelfFundedPingPong_setCountIncrBeforeFunding is SelfFundedPingPongDappSetup {
  function test_setCountIncrBeforeFunding() public {
    uint8 c = s_pingPong.getCountIncrBeforeFunding();

    vm.expectEmit();
    emit CountIncrBeforeFundingSet(c + 1);

    s_pingPong.setCountIncrBeforeFunding(c + 1);
    uint8 c2 = s_pingPong.getCountIncrBeforeFunding();
    assertEq(c2, c + 1);
  }
}
