// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "./interfaces/IAny2EVMMessageReceiver.sol";
import {IEVM2AnyOnRamp} from "./interfaces/IEVM2AnyOnRamp.sol";
import {IRMN} from "./interfaces/IRMN.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {IRouterClient} from "./interfaces/IRouterClient.sol";
import {IWrappedNative} from "./interfaces/IWrappedNative.sol";

import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {CallWithExactGas} from "../shared/call/CallWithExactGas.sol";
import {Client} from "./libraries/Client.sol";
import {Internal} from "./libraries/Internal.sol";

import {IERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @title Router
/// @notice This is the entry point for the end user wishing to send data across chains.
/// @dev This contract is used as a router for both on-ramps and off-ramps
contract Router is IRouter, IRouterClient, ITypeAndVersion, OwnerIsCreator {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.UintSet;

  error FailedToSendValue();
  error InvalidRecipientAddress(address to);
  error OffRampMismatch(uint64 chainSelector, address offRamp);
  error BadARMSignal();

  event OnRampSet(uint64 indexed destChainSelector, address onRamp);
  event OffRampAdded(uint64 indexed sourceChainSelector, address offRamp);
  event OffRampRemoved(uint64 indexed sourceChainSelector, address offRamp);
  event MessageExecuted(bytes32 messageId, uint64 sourceChainSelector, address offRamp, bytes32 calldataHash);

  struct OnRamp {
    uint64 destChainSelector;
    address onRamp;
  }

  struct OffRamp {
    uint64 sourceChainSelector;
    address offRamp;
  }

  string public constant override typeAndVersion = "Router 1.2.0";
  // We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 public constant MAX_RET_BYTES = 4 + 4 * 32;
  // STATIC CONFIG
  // Address of RMN proxy contract (formerly known as ARM)
  address private immutable i_armProxy;

  // DYNAMIC CONFIG
  address private s_wrappedNative;
  // destChainSelector => onRamp address
  // Only ever one onRamp enabled at a time for a given destChainSelector.
  mapping(uint256 destChainSelector => address onRamp) private s_onRamps;
  // Stores [sourceChainSelector << 160 + offramp] as a pair to allow for
  // lookups for specific chain/offramp pairs.
  EnumerableSet.UintSet private s_chainSelectorAndOffRamps;

  constructor(address wrappedNative, address armProxy) {
    // Zero address indicates unsupported auto-wrapping, therefore, unsupported
    // native fee token payments.
    s_wrappedNative = wrappedNative;
    i_armProxy = armProxy;
  }

  // ================================================================
  // │                       Message sending                        │
  // ================================================================

  /// @inheritdoc IRouterClient
  function getFee(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage memory message
  ) external view returns (uint256 fee) {
    if (message.feeToken == address(0)) {
      // For empty feeToken return native quote.
      message.feeToken = address(s_wrappedNative);
    }
    address onRamp = s_onRamps[destinationChainSelector];
    if (onRamp == address(0)) revert UnsupportedDestinationChain(destinationChainSelector);
    return IEVM2AnyOnRamp(onRamp).getFee(destinationChainSelector, message);
  }

  /// @notice This functionality has been removed and will revert when called.
  function getSupportedTokens(uint64 chainSelector) external view returns (address[] memory) {
    if (!isChainSupported(chainSelector)) {
      return new address[](0);
    }
    return IEVM2AnyOnRamp(s_onRamps[uint256(chainSelector)]).getSupportedTokens(chainSelector);
  }

  /// @inheritdoc IRouterClient
  function isChainSupported(uint64 chainSelector) public view returns (bool) {
    return s_onRamps[chainSelector] != address(0);
  }

  /// @inheritdoc IRouterClient
  function ccipSend(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage memory message
  ) external payable whenNotCursed returns (bytes32) {
    address onRamp = s_onRamps[destinationChainSelector];
    if (onRamp == address(0)) revert UnsupportedDestinationChain(destinationChainSelector);
    uint256 feeTokenAmount;
    // address(0) signals payment in true native
    if (message.feeToken == address(0)) {
      // for fee calculation we check the wrapped native price as we wrap
      // as part of the native fee coin payment.
      message.feeToken = s_wrappedNative;
      // We rely on getFee to validate that the feeToken is whitelisted.
      feeTokenAmount = IEVM2AnyOnRamp(onRamp).getFee(destinationChainSelector, message);
      // Ensure sufficient native.
      if (msg.value < feeTokenAmount) revert InsufficientFeeTokenAmount();
      // Wrap and send native payment.
      // Note we take the whole msg.value regardless if its larger.
      feeTokenAmount = msg.value;
      IWrappedNative(message.feeToken).deposit{value: feeTokenAmount}();
      IERC20(message.feeToken).safeTransfer(onRamp, feeTokenAmount);
    } else {
      if (msg.value > 0) revert InvalidMsgValue();
      // We rely on getFee to validate that the feeToken is whitelisted.
      feeTokenAmount = IEVM2AnyOnRamp(onRamp).getFee(destinationChainSelector, message);
      IERC20(message.feeToken).safeTransferFrom(msg.sender, onRamp, feeTokenAmount);
    }

    // Transfer the tokens to the token pools.
    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      IERC20 token = IERC20(message.tokenAmounts[i].token);
      // We rely on getPoolBySourceToken to validate that the token is whitelisted.
      token.safeTransferFrom(
        msg.sender,
        address(IEVM2AnyOnRamp(onRamp).getPoolBySourceToken(destinationChainSelector, token)),
        message.tokenAmounts[i].amount
      );
    }

    return IEVM2AnyOnRamp(onRamp).forwardFromRouter(destinationChainSelector, message, feeTokenAmount, msg.sender);
  }

  // ================================================================
  // │                      Message execution                       │
  // ================================================================

  /// @inheritdoc IRouter
  /// @dev _callWithExactGas protects against return data bombs by capping the return data size at MAX_RET_BYTES.
  function routeMessage(
    Client.Any2EVMMessage calldata message,
    uint16 gasForCallExactCheck,
    uint256 gasLimit,
    address receiver
  ) external override whenNotCursed returns (bool success, bytes memory retData, uint256 gasUsed) {
    // We only permit offRamps to call this function.
    if (!isOffRamp(message.sourceChainSelector, msg.sender)) revert OnlyOffRamp();

    // We encode here instead of the offRamps to constrain specifically what functions
    // can be called from the router.
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);

    (success, retData, gasUsed) = CallWithExactGas._callWithExactGasSafeReturnData(
      data, receiver, gasLimit, gasForCallExactCheck, Internal.MAX_RET_BYTES
    );

    emit MessageExecuted(message.messageId, message.sourceChainSelector, msg.sender, keccak256(data));
    return (success, retData, gasUsed);
  }

  // @notice Merges a chain selector and offRamp address into a single uint256 by shifting the
  // chain selector 160 bits to the left.
  function _mergeChainSelectorAndOffRamp(
    uint64 sourceChainSelector,
    address offRampAddress
  ) internal pure returns (uint256) {
    return (uint256(sourceChainSelector) << 160) + uint160(offRampAddress);
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Gets the wrapped representation of the native fee coin.
  /// @return The address of the ERC20 wrapped native.
  function getWrappedNative() external view returns (address) {
    return s_wrappedNative;
  }

  /// @notice Sets a new wrapped native token.
  /// @param wrappedNative The address of the new wrapped native ERC20 token.
  function setWrappedNative(address wrappedNative) external onlyOwner {
    s_wrappedNative = wrappedNative;
  }

  /// @notice Gets the RMN address, formerly known as ARM
  /// @return The address of the RMN proxy contract, formerly known as ARM
  function getArmProxy() external view returns (address) {
    return i_armProxy;
  }

  /// @inheritdoc IRouter
  function getOnRamp(uint64 destChainSelector) external view returns (address) {
    return s_onRamps[destChainSelector];
  }

  function getOffRamps() external view returns (OffRamp[] memory) {
    uint256[] memory encodedOffRamps = s_chainSelectorAndOffRamps.values();
    OffRamp[] memory offRamps = new OffRamp[](encodedOffRamps.length);
    for (uint256 i = 0; i < encodedOffRamps.length; ++i) {
      uint256 encodedOffRamp = encodedOffRamps[i];
      offRamps[i] =
        OffRamp({sourceChainSelector: uint64(encodedOffRamp >> 160), offRamp: address(uint160(encodedOffRamp))});
    }
    return offRamps;
  }

  /// @inheritdoc IRouter
  function isOffRamp(uint64 sourceChainSelector, address offRamp) public view returns (bool) {
    // We have to encode the sourceChainSelector and offRamp into a uint256 to use as a key in the set.
    return s_chainSelectorAndOffRamps.contains(_mergeChainSelectorAndOffRamp(sourceChainSelector, offRamp));
  }

  /// @notice applyRampUpdates applies a set of ramp changes which provides
  /// the ability to add new chains and upgrade ramps.
  function applyRampUpdates(
    OnRamp[] calldata onRampUpdates,
    OffRamp[] calldata offRampRemoves,
    OffRamp[] calldata offRampAdds
  ) external onlyOwner {
    // Apply egress updates.
    // We permit zero address as way to disable egress.
    for (uint256 i = 0; i < onRampUpdates.length; ++i) {
      OnRamp memory onRampUpdate = onRampUpdates[i];
      s_onRamps[onRampUpdate.destChainSelector] = onRampUpdate.onRamp;
      emit OnRampSet(onRampUpdate.destChainSelector, onRampUpdate.onRamp);
    }

    // Apply ingress updates.
    for (uint256 i = 0; i < offRampRemoves.length; ++i) {
      uint64 sourceChainSelector = offRampRemoves[i].sourceChainSelector;
      address offRampAddress = offRampRemoves[i].offRamp;

      // If the selector-offRamp pair does not exist, revert.
      if (!s_chainSelectorAndOffRamps.remove(_mergeChainSelectorAndOffRamp(sourceChainSelector, offRampAddress))) {
        revert OffRampMismatch(sourceChainSelector, offRampAddress);
      }

      emit OffRampRemoved(sourceChainSelector, offRampAddress);
    }

    for (uint256 i = 0; i < offRampAdds.length; ++i) {
      uint64 sourceChainSelector = offRampAdds[i].sourceChainSelector;
      address offRampAddress = offRampAdds[i].offRamp;

      if (s_chainSelectorAndOffRamps.add(_mergeChainSelectorAndOffRamp(sourceChainSelector, offRampAddress))) {
        emit OffRampAdded(sourceChainSelector, offRampAddress);
      }
    }
  }

  /// @notice Provides the ability for the owner to recover any tokens accidentally
  /// sent to this contract.
  /// @dev Must be onlyOwner to avoid malicious token contract calls.
  /// @param tokenAddress ERC20-token to recover
  /// @param to Destination address to send the tokens to.
  function recoverTokens(address tokenAddress, address to, uint256 amount) external onlyOwner {
    if (to == address(0)) revert InvalidRecipientAddress(to);

    if (tokenAddress == address(0)) {
      (bool success,) = to.call{value: amount}("");
      if (!success) revert FailedToSendValue();
      return;
    }
    IERC20(tokenAddress).safeTransfer(to, amount);
  }

  // ================================================================
  // │                           Access                             │
  // ================================================================

  /// @notice Ensure that the RMN has not cursed the network.
  modifier whenNotCursed() {
    if (IRMN(i_armProxy).isCursed()) revert BadARMSignal();
    _;
  }
}
