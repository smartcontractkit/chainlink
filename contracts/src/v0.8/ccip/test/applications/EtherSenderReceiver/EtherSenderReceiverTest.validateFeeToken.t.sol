// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../../libraries/Client.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_validateFeeToken is EtherSenderReceiverTestSetup {
  error InsufficientMsgValue(uint256 gotAmount, uint256 msgValue);
  error TokenAmountNotEqualToMsgValue(uint256 gotAmount, uint256 msgValue);

  function test_validateFeeToken_valid_native() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(0),
      extraArgs: ""
    });

    s_etherSenderReceiver.validateFeeToken{value: AMOUNT + 1}(message);
  }

  function test_validateFeeToken_valid_feeToken() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    s_etherSenderReceiver.validateFeeToken{value: AMOUNT}(message);
  }

  function test_RevertWhen_validateFeeTokens_feeToken_tokenAmountNotEqualToMsgValue() public {
    Client.EVMTokenAmount[] memory tokenAmount = new Client.EVMTokenAmount[](1);
    tokenAmount[0] = Client.EVMTokenAmount({token: address(s_weth), amount: AMOUNT});
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(XCHAIN_RECEIVER),
      data: "",
      tokenAmounts: tokenAmount,
      feeToken: address(s_weth),
      extraArgs: ""
    });

    vm.expectRevert(abi.encodeWithSelector(TokenAmountNotEqualToMsgValue.selector, AMOUNT, AMOUNT + 1));
    s_etherSenderReceiver.validateFeeToken{value: AMOUNT + 1}(message);
  }
}
