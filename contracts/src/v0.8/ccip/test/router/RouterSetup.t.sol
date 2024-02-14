// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "../BaseTest.t.sol";
import {Router} from "../../Router.sol";
import {WETH9} from "../WETH9.sol";
import {Client} from "../../libraries/Client.sol";

contract RouterSetup is BaseTest {
  Router internal s_sourceRouter;
  Router internal s_destRouter;

  function setUp() public virtual override {
    BaseTest.setUp();

    if (address(s_sourceRouter) == address(0)) {
      WETH9 weth = new WETH9();
      s_sourceRouter = new Router(address(weth), address(s_mockARM));
      vm.label(address(s_sourceRouter), "sourceRouter");
    }
    if (address(s_destRouter) == address(0)) {
      WETH9 weth = new WETH9();
      s_destRouter = new Router(address(weth), address(s_mockARM));
      vm.label(address(s_destRouter), "destRouter");
    }
  }

  function generateReceiverMessage(uint64 chainSelector) internal pure returns (Client.Any2EVMMessage memory) {
    Client.EVMTokenAmount[] memory ta = new Client.EVMTokenAmount[](0);
    return
      Client.Any2EVMMessage({
        messageId: bytes32("a"),
        sourceChainSelector: chainSelector,
        sender: bytes("a"),
        data: bytes("a"),
        destTokenAmounts: ta
      });
  }
}
