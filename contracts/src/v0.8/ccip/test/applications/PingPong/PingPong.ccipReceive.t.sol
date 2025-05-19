// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {PingPongDemo} from "../../../applications/PingPongDemo.sol";
import {Client} from "../../../libraries/Client.sol";

import {PingPongDappSetup} from "./PingPongDappSetup.t.sol";

contract PingPong_ccipReceive is PingPongDappSetup {
  function test_CcipReceive() public {
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
