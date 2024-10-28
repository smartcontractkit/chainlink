// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {IRouterClient} from "../interfaces/IRouterClient.sol";
import {IWrappedNative} from "../interfaces/IWrappedNative.sol";

import {Client} from "./../libraries/Client.sol";
import {CCIPReceiver} from "./CCIPReceiver.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

//solhint-disable interface-starts-with-i
interface CCIPRouter {
  function getWrappedNative() external view returns (address);
}

/// @notice A contract that can send raw ether cross-chain using CCIP.
/// Since CCIP only supports ERC-20 token transfers, this contract accepts
/// normal ether, wraps it, and uses CCIP to send it cross-chain.
/// On the receiving side, the wrapped ether is unwrapped and sent to the final receiver.
/// @notice This contract only supports chains where the wrapped native contract
/// is the WETH contract (i.e not WMATIC, or WAVAX, etc.). This is because the
/// receiving contract will always unwrap the ether using it's local wrapped native contract.
/// @dev This contract is both a sender and a receiver. This same contract can be
/// deployed on source and destination chains to facilitate cross-chain ether transfers
/// and act as a sender and a receiver.
/// @dev This contract is intentionally ownerless and permissionless. This contract
/// will never hold any excess funds, native or otherwise, when used correctly.
contract EtherSenderReceiver is CCIPReceiver, ITypeAndVersion {
  using SafeERC20 for IERC20;

  error InvalidTokenAmounts(uint256 gotAmounts);
  error InvalidToken(address gotToken, address expectedToken);
  error TokenAmountNotEqualToMsgValue(uint256 gotAmount, uint256 msgValue);

  string public constant override typeAndVersion = "EtherSenderReceiver 1.5.0";

  /// @notice The wrapped native token address.
  /// @dev If the wrapped native token address changes on the router, this contract will need to be redeployed.
  IWrappedNative public immutable i_weth;

  /// @param router The CCIP router address.
  constructor(
    address router
  ) CCIPReceiver(router) {
    i_weth = IWrappedNative(CCIPRouter(router).getWrappedNative());
    i_weth.approve(router, type(uint256).max);
  }

  /// @notice Need this in order to unwrap correctly.
  receive() external payable {}

  /// @notice Get the fee for sending a message to a destination chain.
  /// This is mirrored from the router for convenience, construct the appropriate
  /// message and get it's fee.
  /// @param destinationChainSelector The destination chainSelector
  /// @param message The cross-chain CCIP message including data and/or tokens
  /// @return fee returns execution fee for the message
  /// delivery to destination chain, denominated in the feeToken specified in the message.
  /// @dev Reverts with appropriate reason upon invalid message.
  function getFee(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 fee) {
    Client.EVM2AnyMessage memory validatedMessage = _validatedMessage(message);

    return IRouterClient(getRouter()).getFee(destinationChainSelector, validatedMessage);
  }

  /// @notice Send raw native tokens cross-chain.
  /// @param destinationChainSelector The destination chain selector.
  /// @param message The CCIP message with the following fields correctly set:
  /// - bytes receiver: The _contract_ address on the destination chain that will receive the wrapped ether.
  /// The caller must ensure that this contract address is correct, otherwise funds may be lost forever.
  /// - address feeToken: The fee token address. Must be address(0) for native tokens, or a supported CCIP fee token otherwise (i.e, LINK token).
  /// In the event a feeToken is set, we will transferFrom the caller the fee amount before sending the message, in order to forward them to the router.
  /// - EVMTokenAmount[] tokenAmounts: The tokenAmounts array must contain a single element with the following fields:
  ///   - uint256 amount: The amount of ether to send.
  /// There are a couple of cases here that depend on the fee token specified:
  /// 1. If feeToken == address(0), the fee must be included in msg.value. Therefore tokenAmounts[0].amount must be less than msg.value,
  ///    and the difference will be used as the fee.
  /// 2. If feeToken != address(0), the fee is not included in msg.value, and tokenAmounts[0].amount must be equal to msg.value.
  ///    these fees to the CCIP router.
  /// @return messageId The CCIP message ID.
  function ccipSend(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external payable returns (bytes32) {
    _validateFeeToken(message);
    Client.EVM2AnyMessage memory validatedMessage = _validatedMessage(message);

    i_weth.deposit{value: validatedMessage.tokenAmounts[0].amount}();

    uint256 fee = IRouterClient(getRouter()).getFee(destinationChainSelector, validatedMessage);
    if (validatedMessage.feeToken != address(0)) {
      // If the fee token is not native, we need to transfer the fee to this contract and re-approve it to the router.
      // Its not possible to have any leftover tokens in this path because we transferFrom the exact fee that CCIP
      // requires from the caller.
      IERC20(validatedMessage.feeToken).safeTransferFrom(msg.sender, address(this), fee);

      // We gave an infinite approval of weth to the router in the constructor.
      if (validatedMessage.feeToken != address(i_weth)) {
        IERC20(validatedMessage.feeToken).approve(getRouter(), fee);
      }

      return IRouterClient(getRouter()).ccipSend(destinationChainSelector, validatedMessage);
    }

    // We don't want to keep any excess ether in this contract, so we send over the entire address(this).balance as the fee.
    // CCIP will revert if the fee is insufficient, so we don't need to check here.
    return IRouterClient(getRouter()).ccipSend{value: address(this).balance}(destinationChainSelector, validatedMessage);
  }

  /// @notice Validate the message content.
  /// @dev Only allows a single token to be sent. Always overwritten to be address(i_weth)
  /// and receiver is always msg.sender.
  function _validatedMessage(
    Client.EVM2AnyMessage calldata message
  ) internal view returns (Client.EVM2AnyMessage memory) {
    Client.EVM2AnyMessage memory validatedMessage = message;

    if (validatedMessage.tokenAmounts.length != 1) {
      revert InvalidTokenAmounts(validatedMessage.tokenAmounts.length);
    }

    validatedMessage.data = abi.encode(msg.sender);
    validatedMessage.tokenAmounts[0].token = address(i_weth);

    return validatedMessage;
  }

  function _validateFeeToken(
    Client.EVM2AnyMessage calldata message
  ) internal view {
    uint256 tokenAmount = message.tokenAmounts[0].amount;

    if (message.feeToken != address(0)) {
      // If the fee token is NOT native, then the token amount must be equal to msg.value.
      // This is done to ensure that there is no leftover ether in this contract.
      if (msg.value != tokenAmount) {
        revert TokenAmountNotEqualToMsgValue(tokenAmount, msg.value);
      }
    }
  }

  /// @notice Receive the wrapped ether, unwrap it, and send it to the specified EOA in the data field.
  /// @param message The CCIP message containing the wrapped ether amount and the final receiver.
  /// @dev The code below should never revert if the message being is valid according
  /// to the above _validatedMessage and _validateFeeToken functions.
  function _ccipReceive(
    Client.Any2EVMMessage memory message
  ) internal override {
    address receiver = abi.decode(message.data, (address));

    if (message.destTokenAmounts.length != 1) {
      revert InvalidTokenAmounts(message.destTokenAmounts.length);
    }

    if (message.destTokenAmounts[0].token != address(i_weth)) {
      revert InvalidToken(message.destTokenAmounts[0].token, address(i_weth));
    }

    uint256 tokenAmount = message.destTokenAmounts[0].amount;
    i_weth.withdraw(tokenAmount);

    // it is possible that the below call may fail if receiver.code.length > 0 and the contract
    // doesn't e.g have a receive() or a fallback() function.
    (bool success,) = payable(receiver).call{value: tokenAmount}("");
    if (!success) {
      // We have a few options here:
      // 1. Revert: this is bad generally because it may mean that these tokens are stuck.
      // 2. Store the tokens in a mapping and allow the user to withdraw them with another tx.
      // 3. Send WETH to the receiver address.
      // We opt for (3) here because at least the receiver will have the funds and can unwrap them if needed.
      // However it is worth noting that if receiver is actually a contract AND the contract _cannot_ withdraw
      // the WETH, then the WETH will be stuck in this contract.
      i_weth.deposit{value: tokenAmount}();
      i_weth.transfer(receiver, tokenAmount);
    }
  }
}
