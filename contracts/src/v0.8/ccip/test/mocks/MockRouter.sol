// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouter} from "../../interfaces/IRouter.sol";
import {IRouterClient} from "../../interfaces/IRouterClient.sol";
import {IAny2EVMMessageReceiver} from "../../interfaces/IAny2EVMMessageReceiver.sol";

import {Client} from "../../libraries/Client.sol";
import {CallWithExactGas} from "../../../shared/call/CallWithExactGas.sol";
import {Internal} from "../../libraries/Internal.sol";

import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {ERC165Checker} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165Checker.sol";

contract MockCCIPRouter is IRouter, IRouterClient {
  using SafeERC20 for IERC20;
  using ERC165Checker for address;

  error InvalidAddress(bytes encodedAddress);
  error InvalidExtraArgsTag();
  error ReceiverError(bytes error);

  event MessageExecuted(bytes32 messageId, uint64 sourceChainSelector, address offRamp, bytes32 calldataHash);
  event MsgExecuted(bool success, bytes retData, uint256 gasUsed);

  uint16 public constant GAS_FOR_CALL_EXACT_CHECK = 5_000;
  uint64 public constant DEFAULT_GAS_LIMIT = 200_000;

  function routeMessage(
    Client.Any2EVMMessage calldata message,
    uint16 gasForCallExactCheck,
    uint256 gasLimit,
    address receiver
  ) external returns (bool success, bytes memory retData, uint256 gasUsed) {
    return _routeMessage(message, gasForCallExactCheck, gasLimit, receiver);
  }

  function _routeMessage(
    Client.Any2EVMMessage memory message,
    uint16 gasForCallExactCheck,
    uint256 gasLimit,
    address receiver
  ) internal returns (bool success, bytes memory retData, uint256 gasUsed) {
    // Only send through the router if the receiver is a contract and implements the IAny2EVMMessageReceiver interface.
    if (receiver.code.length == 0 || !receiver.supportsInterface(type(IAny2EVMMessageReceiver).interfaceId))
      return (true, "", 0);

    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);

    (success, retData, gasUsed) = CallWithExactGas._callWithExactGasSafeReturnData(
      data,
      receiver,
      gasLimit,
      gasForCallExactCheck,
      Internal.MAX_RET_BYTES
    );

    // Event to assist testing, does not exist on real deployments
    emit MsgExecuted(success, retData, gasUsed);

    // Real router event
    emit MessageExecuted(message.messageId, message.sourceChainSelector, msg.sender, keccak256(data));
    return (success, retData, gasUsed);
  }

  /// @notice Sends the tx locally to the receiver instead of on the destination chain.
  /// @dev Ignores destinationChainSelector
  /// @dev Returns a mock message ID, which is not calculated from the message contents in the
  /// same way as the real message ID.
  function ccipSend(
    uint64, // destinationChainSelector
    Client.EVM2AnyMessage calldata message
  ) external payable returns (bytes32) {
    if (message.receiver.length != 32) revert InvalidAddress(message.receiver);
    uint256 decodedReceiver = abi.decode(message.receiver, (uint256));
    // We want to disallow sending to address(0) and to precompiles, which exist on address(1) through address(9).
    if (decodedReceiver > type(uint160).max || decodedReceiver < 10) revert InvalidAddress(message.receiver);

    address receiver = address(uint160(decodedReceiver));
    uint256 gasLimit = _fromBytes(message.extraArgs).gasLimit;
    bytes32 mockMsgId = keccak256(abi.encode(message));

    Client.Any2EVMMessage memory executableMsg = Client.Any2EVMMessage({
      messageId: mockMsgId,
      sourceChainSelector: 16015286601757825753, // Sepolia
      sender: abi.encode(msg.sender),
      data: message.data,
      destTokenAmounts: message.tokenAmounts
    });

    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      IERC20(message.tokenAmounts[i].token).safeTransferFrom(msg.sender, receiver, message.tokenAmounts[i].amount);
    }

    (bool success, bytes memory retData, ) = _routeMessage(executableMsg, GAS_FOR_CALL_EXACT_CHECK, gasLimit, receiver);

    if (!success) revert ReceiverError(retData);

    return mockMsgId;
  }

  function _fromBytes(bytes calldata extraArgs) internal pure returns (Client.EVMExtraArgsV1 memory) {
    if (extraArgs.length == 0) {
      return Client.EVMExtraArgsV1({gasLimit: DEFAULT_GAS_LIMIT});
    }
    if (bytes4(extraArgs) != Client.EVM_EXTRA_ARGS_V1_TAG) revert InvalidExtraArgsTag();
    return abi.decode(extraArgs[4:], (Client.EVMExtraArgsV1));
  }

  /// @notice Always returns true to make sure this check can be performed on any chain.
  function isChainSupported(uint64) external pure returns (bool supported) {
    return true;
  }

  /// @notice Returns an empty array.
  function getSupportedTokens(uint64) external pure returns (address[] memory tokens) {
    return new address[](0);
  }

  /// @notice Returns 0 as the fee is not supported in this mock contract.
  function getFee(uint64, Client.EVM2AnyMessage memory) external pure returns (uint256 fee) {
    return 0;
  }

  /// @notice Always returns address(1234567890)
  function getOnRamp(uint64 /* destChainSelector */) external pure returns (address onRampAddress) {
    return address(1234567890);
  }

  /// @notice Always returns true
  function isOffRamp(uint64 /* sourceChainSelector */, address /* offRamp */) external pure returns (bool) {
    return true;
  }
}
