// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Router} from "../Router.sol";
import {Client} from "../libraries/Client.sol";
import {EVM2EVMOnRamp} from "../onRamp/EVM2EVMOnRamp.sol";
import {PingPongDemo} from "./PingPongDemo.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract SelfFundedPingPong is PingPongDemo {
  string public constant override typeAndVersion = "SelfFundedPingPong 1.5.0";

  event Funded();
  event CountIncrBeforeFundingSet(uint8 countIncrBeforeFunding);

  // Defines the increase in ping pong count before self-funding is attempted.
  // Set to 0 to disable auto-funding, auto-funding only works for ping-pongs that are set as NOPs in the onRamp.
  uint8 private s_countIncrBeforeFunding;

  constructor(address router, IERC20 feeToken, uint8 roundTripsBeforeFunding) PingPongDemo(router, feeToken) {
    // PingPong count increases by 2 for each round trip.
    s_countIncrBeforeFunding = roundTripsBeforeFunding * 2;
  }

  function _respond(uint256 pingPongCount) internal override {
    if (pingPongCount & 1 == 1) {
      emit Ping(pingPongCount);
    } else {
      emit Pong(pingPongCount);
    }

    fundPingPong(pingPongCount);

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(s_counterpartAddress),
      data: abi.encode(pingPongCount),
      tokenAmounts: new Client.EVMTokenAmount[](0),
      extraArgs: "",
      feeToken: address(s_feeToken)
    });
    Router(getRouter()).ccipSend(s_counterpartChainSelector, message);
  }

  /// @notice A function that is responsible for funding this contract.
  /// The contract can only be funded if it is set as a nop in the target onRamp.
  /// In case your contract is not a nop you can prevent this function from being called by setting s_countIncrBeforeFunding=0.
  function fundPingPong(uint256 pingPongCount) public {
    // If selfFunding is disabled, or ping pong count has not reached s_countIncrPerFunding, do not attempt funding.
    if (s_countIncrBeforeFunding == 0 || pingPongCount < s_countIncrBeforeFunding) return;

    // Ping pong on one side will always be even, one side will always to odd.
    if (pingPongCount % s_countIncrBeforeFunding <= 1) {
      EVM2EVMOnRamp(Router(getRouter()).getOnRamp(s_counterpartChainSelector)).payNops();
      emit Funded();
    }
  }

  function getCountIncrBeforeFunding() external view returns (uint8) {
    return s_countIncrBeforeFunding;
  }

  function setCountIncrBeforeFunding(uint8 countIncrBeforeFunding) external onlyOwner {
    s_countIncrBeforeFunding = countIncrBeforeFunding;
    emit CountIncrBeforeFundingSet(countIncrBeforeFunding);
  }
}
