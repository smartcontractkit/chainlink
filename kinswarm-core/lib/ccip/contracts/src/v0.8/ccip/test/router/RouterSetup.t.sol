// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Router} from "../../Router.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {BaseTest} from "../BaseTest.t.sol";
import {WETH9} from "../WETH9.sol";

contract RouterSetup is BaseTest {
  Router internal s_sourceRouter;
  Router internal s_destRouter;

  function setUp() public virtual override {
    BaseTest.setUp();

    if (address(s_sourceRouter) == address(0)) {
      WETH9 weth = new WETH9();
      s_sourceRouter = new Router(address(weth), address(s_mockRMN));
      vm.label(address(s_sourceRouter), "sourceRouter");
    }
    if (address(s_destRouter) == address(0)) {
      WETH9 weth = new WETH9();
      s_destRouter = new Router(address(weth), address(s_mockRMN));
      vm.label(address(s_destRouter), "destRouter");
    }
  }

  function _generateReceiverMessage(
    uint64 chainSelector
  ) internal pure returns (Client.Any2EVMMessage memory) {
    Client.EVMTokenAmount[] memory ta = new Client.EVMTokenAmount[](0);
    return Client.Any2EVMMessage({
      messageId: bytes32("a"),
      sourceChainSelector: chainSelector,
      sender: bytes("a"),
      data: bytes("a"),
      destTokenAmounts: ta
    });
  }

  function _generateSourceTokenData() internal pure returns (Internal.SourceTokenData memory) {
    return Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(address(12312412312)),
      destTokenAddress: abi.encode(address(9809808909)),
      extraData: "",
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });
  }
}
