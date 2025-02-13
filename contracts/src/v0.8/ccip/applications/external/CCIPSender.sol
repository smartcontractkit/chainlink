// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../interfaces/IRouterClient.sol";

import {Client} from "../../libraries/Client.sol";
import {CCIPBase} from "./CCIPBase.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Example of a client which supports sending messages via CCIP.
/// @dev If chain specific logic is required for different chain families (e.g. particular decoding the bytes sender
/// for authorization checks), it may be required to point to a helper authorization contract unless all chain families
/// are known initially.
/// @dev The ccipSend function does NOT support pre-funding message-fees, and must acquire fee-tokens from the
/// user before the call to the router is made.
contract CCIPSender is CCIPBase {
  using SafeERC20 for IERC20;

  event MessageSent(bytes32 messageId);

  constructor(
    address router
  ) CCIPBase(router) {}

  function typeAndVersion() external pure virtual returns (string memory) {
    return "CCIPSender 1.6.0-dev";
  }

  /// @notice Sends a message through CCIP to the router.
  /// @param destChainSelector A CCIP-Specific and unique chain identifier.
  /// @param tokenAmounts An array of token addresses and amounts to trasfer from the user and forward to the router.
  /// @param data Arbitrary data to be sent through CCIP to the destination contract. If the destination contract is an
  /// EOA, then no calls will be made, and only the tokens forwarded.
  /// @param feeToken The token used to pay for fees in CCIP. address(0) denotes paying fees in native chain tokens.
  /// @return messageId the unique identifier for a message determined by the router when a message is sent.
  function ccipSend(
    uint64 destChainSelector,
    Client.EVMTokenAmount[] memory tokenAmounts,
    bytes memory data,
    address feeToken
  ) public payable isValidChain(destChainSelector) returns (bytes32 messageId) {
    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: s_chainConfigs[destChainSelector].recipient,
      data: data,
      tokenAmounts: tokenAmounts,
      feeToken: feeToken,
      extraArgs: s_chainConfigs[destChainSelector].extraArgsBytes
    });

    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      // Transfer the tokens to this contract and approve the router to take possession during message processing.
      IERC20(tokenAmounts[i].token).safeTransferFrom(msg.sender, address(this), tokenAmounts[i].amount);
      IERC20(tokenAmounts[i].token).safeIncreaseAllowance(s_ccipRouter, tokenAmounts[i].amount);
    }

    uint256 fee = IRouterClient(s_ccipRouter).getFee(destChainSelector, message);

    // If the fee token is not native, transfer the fee amounts from the sender and approve the router.
    if (feeToken != address(0)) {
      IERC20(feeToken).safeTransferFrom(msg.sender, address(this), fee);

      // If the fee token is also included in tokenAmounts, and the amount is very large, this call may fail due to
      // allowance overflowing the max uint256 value. This is a known issue with the ERC20 standard, and is up to the
      // user to handle.
      IERC20(feeToken).safeIncreaseAllowance(s_ccipRouter, fee);
    }

    messageId =
      IRouterClient(s_ccipRouter).ccipSend{value: feeToken == address(0) ? fee : 0}(destChainSelector, message);

    emit MessageSent(messageId);

    return messageId;
  }
}
