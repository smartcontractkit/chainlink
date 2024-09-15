// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IRouterClient} from "../interfaces/IRouterClient.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

contract TokenProxy is OwnerIsCreator {
  using SafeERC20 for IERC20;

  error InvalidToken();
  error NoDataAllowed();
  error GasShouldBeZero();

  /// @notice The CCIP router contract
  IRouterClient internal immutable i_ccipRouter;
  /// @notice Only this token is allowed to be sent using this proxy
  address internal immutable i_token;

  constructor(address router, address token) OwnerIsCreator() {
    i_ccipRouter = IRouterClient(router);
    i_token = token;
    // Approve the router to spend an unlimited amount of tokens to reduce
    // gas cost per tx.
    IERC20(token).approve(router, type(uint256).max);
  }

  /// @notice Simply forwards the request to the CCIP router and returns the result.
  /// @param destinationChainSelector The destination chainSelector
  /// @param message The cross-chain CCIP message including data and/or tokens
  /// @return fee returns execution fee for the message delivery to destination chain,
  /// denominated in the feeToken specified in the message.
  /// @dev Reverts with appropriate reason upon invalid message.
  function getFee(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 fee) {
    _validateMessage(message);
    return i_ccipRouter.getFee(destinationChainSelector, message);
  }

  /// @notice Validates the message content, forwards it to the CCIP router and returns the result.
  function ccipSend(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external payable returns (bytes32 messageId) {
    _validateMessage(message);
    if (message.feeToken != address(0)) {
      // This path is probably warmed up already so the extra cost isn't too bad.
      uint256 feeAmount = i_ccipRouter.getFee(destinationChainSelector, message);
      IERC20(message.feeToken).safeTransferFrom(msg.sender, address(this), feeAmount);
      IERC20(message.feeToken).approve(address(i_ccipRouter), feeAmount);
    }

    // Transfer the tokens from the sender to this contract.
    IERC20(message.tokenAmounts[0].token).transferFrom(msg.sender, address(this), message.tokenAmounts[0].amount);

    return i_ccipRouter.ccipSend{value: msg.value}(destinationChainSelector, message);
  }

  /// @notice Validates the message content.
  /// @dev Only allows a single token to be sent, and no data.
  function _validateMessage(Client.EVM2AnyMessage calldata message) internal view {
    if (message.tokenAmounts.length != 1 || message.tokenAmounts[0].token != i_token) revert InvalidToken();
    if (message.data.length > 0) revert NoDataAllowed();

    if (message.extraArgs.length == 0 || bytes4(message.extraArgs) != Client.EVM_EXTRA_ARGS_V1_TAG) {
      revert GasShouldBeZero();
    }

    if (abi.decode(message.extraArgs[4:], (Client.EVMExtraArgsV1)).gasLimit != 0) revert GasShouldBeZero();
  }

  /// @notice Returns the CCIP router contract.
  function getRouter() external view returns (IRouterClient) {
    return i_ccipRouter;
  }

  /// @notice Returns the token that this proxy is allowed to send.
  function getToken() external view returns (address) {
    return i_token;
  }
}
