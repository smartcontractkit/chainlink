// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDemo} from "../../../applications/PingPongDemo.sol";
import {OnRampSetup} from "../../onRamp/OnRamp/OnRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

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
