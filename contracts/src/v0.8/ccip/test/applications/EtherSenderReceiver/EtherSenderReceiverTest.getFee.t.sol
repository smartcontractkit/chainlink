// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../../interfaces/IRouterClient.sol";
import {Client} from "../../../libraries/Client.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_getFee is EtherSenderReceiverTestSetup {
  uint64 internal constant DESTINATION_CHAIN_SELECTOR = 424242;
  uint256 internal constant FEE_WEI = 121212;

  function test_getFee() public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: address(0), amount: AMOUNT});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: address(0),
      extraArgs: ""
    });

    Client.EVM2AnyMessage memory validatedMessage = s_etherSenderReceiver.validatedMessage(message);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IRouterClient.getFee.selector, DESTINATION_CHAIN_SELECTOR, validatedMessage),
      abi.encode(FEE_WEI)
    );

    uint256 fee = s_etherSenderReceiver.getFee(DESTINATION_CHAIN_SELECTOR, message);
    assertEq(fee, FEE_WEI, "fee must be feeWei");
  }
}
