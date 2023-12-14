// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IRouterClient} from "../../../interfaces/IRouterClient.sol";

import {OwnerIsCreator} from "../../../../shared/access/OwnerIsCreator.sol";
import {Client} from "../../../libraries/Client.sol";
import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @title FacadeClient - A simple proxy for calling Router
contract FacadeClient {
  address private immutable i_router;
  uint64 private immutable i_destChainSelector;
  IERC20 private immutable i_sourceToken;
  IERC20 private immutable i_feeToken;

  uint256 private s_msg_sequence = 1;

  constructor(address router, uint64 destChainSelector, IERC20 sourceToken, IERC20 feeToken) {
    i_router = router;
    i_destChainSelector = destChainSelector;
    i_sourceToken = sourceToken;
    i_feeToken = feeToken;

    sourceToken.approve(address(router), 2 ** 256 - 1);
    feeToken.approve(address(router), 2 ** 256 - 1);
  }

  /// @dev Calls Router to initiate CCIP send.
  /// The expectation is that s_msg_sequence will alway match the sequence in emitted CCIP messages.
  function send(uint256 amount) public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = address(i_sourceToken);
    tokenAmounts[0].amount = amount;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(address(100)),
      data: abi.encodePacked(s_msg_sequence),
      tokenAmounts: tokenAmounts,
      extraArgs: "",
      feeToken: address(i_feeToken)
    });

    s_msg_sequence++;

    IRouterClient(i_router).ccipSend(i_destChainSelector, message);
  }

  function getSequence() public view returns (uint256) {
    return s_msg_sequence;
  }
}
