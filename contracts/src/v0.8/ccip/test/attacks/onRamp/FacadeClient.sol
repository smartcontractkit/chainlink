// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRouterClient} from "../../../interfaces/IRouterClient.sol";

import {Client} from "../../../libraries/Client.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/// @title FacadeClient - A simple proxy for calling Router
contract FacadeClient {
  address private immutable i_router;
  uint64 private immutable i_destChainSelector;
  IERC20 private immutable i_sourceToken;
  IERC20 private immutable i_feeToken;
  address private immutable i_receiver;

  uint256 private s_msg_sequence = 1;

  constructor(address router, uint64 destChainSelector, IERC20 sourceToken, IERC20 feeToken, address receiver) {
    i_router = router;
    i_destChainSelector = destChainSelector;
    i_sourceToken = sourceToken;
    i_feeToken = feeToken;
    i_receiver = receiver;

    sourceToken.approve(address(router), 2 ** 256 - 1);
    feeToken.approve(address(router), 2 ** 256 - 1);
  }

  /// @dev Calls Router to initiate CCIP send.
  /// The expectation is that s_msg_sequence will always match the sequence in emitted CCIP messages.
  function send(uint256 amount) public {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0].token = address(i_sourceToken);
    tokenAmounts[0].amount = amount;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(i_receiver),
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
