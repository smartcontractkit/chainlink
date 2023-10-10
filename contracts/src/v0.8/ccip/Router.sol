// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IRouterClient} from "./interfaces/IRouterClient.sol";
import {IRouter} from "./interfaces/IRouter.sol";
import {IEVM2AnyOnRamp} from "./interfaces/IEVM2AnyOnRamp.sol";
import {IARM} from "./interfaces/IARM.sol";
import {IWrappedNative} from "./interfaces/IWrappedNative.sol";
import {IAny2EVMMessageReceiver} from "./interfaces/IAny2EVMMessageReceiver.sol";

import {Client} from "./libraries/Client.sol";
import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";

import {EnumerableMap} from "../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/structs/EnumerableMap.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "../vendor/openzeppelin-solidity/v4.8.0/contracts/token/ERC20/IERC20.sol";

/// @title Router
/// @notice This is the entry point for the end user wishing to send data across chains.
/// @dev This contract is used as a router for both on-ramps and off-ramps
contract Router is IRouter, IRouterClient, ITypeAndVersion, OwnerIsCreator {
  using SafeERC20 for IERC20;
  using EnumerableMap for EnumerableMap.AddressToUintMap;

  error FailedToSendValue();
  error InvalidRecipientAddress(address to);
  error OffRampMismatch();
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

  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "Router 1.0.0";
  // We limit return data to a selector plus 4 words. This is to avoid
  // malicious contracts from returning large amounts of data and causing
  // repeated out-of-gas scenarios.
  uint16 public constant MAX_RET_BYTES = 4 + 4 * 32;
  // STATIC CONFIG
  // Address of arm proxy contract.
  address private immutable i_armProxy;

  // DYNAMIC CONFIG
  address private s_wrappedNative;
  // destChainSelector => onRamp address
  // Only ever one onRamp enabled at a time for a given destChainSelector.
  mapping(uint256 destChainSelector => address onRamp) private s_onRamps;
  // Mapping of offRamps to source chain ids
  // Can be multiple offRamps enabled at a time for a given sourceChainSelector,
  // for example during an no downtime upgrade while v1 messages are being flushed.
  EnumerableMap.AddressToUintMap private s_offRamps;

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
    return IEVM2AnyOnRamp(onRamp).getFee(message);
  }

  /// @inheritdoc IRouterClient
  function getSupportedTokens(uint64 chainSelector) external view returns (address[] memory) {
    if (!isChainSupported(chainSelector)) {
      return new address[](0);
    }
    return IEVM2AnyOnRamp(s_onRamps[uint256(chainSelector)]).getSupportedTokens();
  }

  /// @inheritdoc IRouterClient
  function isChainSupported(uint64 chainSelector) public view returns (bool) {
    return s_onRamps[chainSelector] != address(0);
  }

  /// @inheritdoc IRouterClient
  function ccipSend(
    uint64 destinationChainSelector,
    Client.EVM2AnyMessage memory message
  ) external payable whenHealthy returns (bytes32) {
    address onRamp = s_onRamps[destinationChainSelector];
    if (onRamp == address(0)) revert UnsupportedDestinationChain(destinationChainSelector);
    uint256 feeTokenAmount;
    // address(0) signals payment in true native
    if (message.feeToken == address(0)) {
      // for fee calculation we check the wrapped native price as we wrap
      // as part of the native fee coin payment.
      message.feeToken = s_wrappedNative;
      // We rely on getFee to validate that the feeToken is whitelisted.
      feeTokenAmount = IEVM2AnyOnRamp(onRamp).getFee(message);
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
      feeTokenAmount = IEVM2AnyOnRamp(onRamp).getFee(message);
      IERC20(message.feeToken).safeTransferFrom(msg.sender, onRamp, feeTokenAmount);
    }

    // Transfer the tokens to the token pools.
    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      IERC20 token = IERC20(message.tokenAmounts[i].token);
      // We rely on getPoolBySourceToken to validate that the token is whitelisted.
      token.safeTransferFrom(
        msg.sender,
        address(IEVM2AnyOnRamp(onRamp).getPoolBySourceToken(token)),
        message.tokenAmounts[i].amount
      );
    }

    return IEVM2AnyOnRamp(onRamp).forwardFromRouter(message, feeTokenAmount, msg.sender);
  }

  // ================================================================
  // │                      Message execution                       │
  // ================================================================

  /// @inheritdoc IRouter
  /// @dev Handles the edge case where we want to pass a specific amount of gas,
  /// @dev but EIP-150 sends all but 1/64 of the remaining gas instead so the user gets
  /// @dev less gas than they paid for. The other 2 parts of EIP-150 do not apply since
  /// @dev a) we hard code value=0 and b) we ensure code already exists.
  /// @dev If we revert instead, then that will never happen.
  /// @dev Separately we capture the return data up to a maximum size to avoid return bombs,
  /// @dev borrowed from https://github.com/nomad-xyz/ExcessivelySafeCall/blob/main/src/ExcessivelySafeCall.sol.
  function routeMessage(
    Client.Any2EVMMessage calldata message,
    uint16 gasForCallExactCheck,
    uint256 gasLimit,
    address receiver
  )
    external
    override
    onlyOffRamp(message.sourceChainSelector)
    whenHealthy
    returns (bool success, bytes memory retData)
  {
    // We encode here instead of the offRamps to constrain specifically what functions
    // can be called from the router.
    bytes memory data = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);
    // allocate retData memory ahead of time
    retData = new bytes(MAX_RET_BYTES);

    // solhint-disable-next-line no-inline-assembly
    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      // Note we do this check prior to measuring gas so gasForCallExactCheck (our "cushion")
      // doesn't need to account for it.
      if iszero(extcodesize(receiver)) {
        revert(0, 0)
      }

      let g := gas()
      // Compute g -= gasForCallExactCheck and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas. gasForCallExactCheck ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, gasForCallExactCheck) {
        revert(0, 0)
      }
      g := sub(g, gasForCallExactCheck)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasLimit)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasLimit, receiver, 0, add(data, 0x20), mload(data), 0, 0)

      // limit our copy to MAX_RET_BYTES bytes
      let toCopy := returndatasize()
      if gt(toCopy, MAX_RET_BYTES) {
        toCopy := MAX_RET_BYTES
      }
      // Store the length of the copied bytes
      mstore(retData, toCopy)
      // copy the bytes from retData[0:_toCopy]
      returndatacopy(add(retData, 0x20), 0, toCopy)
    }
    emit MessageExecuted(message.messageId, message.sourceChainSelector, msg.sender, keccak256(data));
    return (success, retData);
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

  /// @notice Gets the arm address
  /// @return The address of the ARM proxy contract.
  function getArmProxy() external view returns (address) {
    return i_armProxy;
  }

  /// @notice Return the configured onramp for specific a destination chain.
  /// @param destChainSelector The destination chain Id to get the onRamp for.
  /// @return The address of the onRamp.
  function getOnRamp(uint64 destChainSelector) external view returns (address) {
    return s_onRamps[destChainSelector];
  }

  /// @notice Return a full list of configured offRamps.
  function getOffRamps() external view returns (OffRamp[] memory) {
    OffRamp[] memory offRamps = new OffRamp[](s_offRamps.length());
    for (uint256 i = 0; i < offRamps.length; ++i) {
      (address offRamp, uint256 sourceChainSelector) = s_offRamps.at(i);
      offRamps[i] = OffRamp({sourceChainSelector: uint64(sourceChainSelector), offRamp: offRamp});
    }
    return offRamps;
  }

  /// @notice Returns true if the given address is a permissioned offRamp
  /// and sourceChainSelector if so.
  function isOffRamp(address offRamp) external view returns (bool, uint64) {
    (bool exists, uint256 sourceChainSelector) = s_offRamps.tryGet(offRamp);
    return (exists, uint64(sourceChainSelector));
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
    // We permit an empty list as a way to disable ingress.
    for (uint256 i = 0; i < offRampRemoves.length; ++i) {
      uint64 rampSelector = offRampRemoves[i].sourceChainSelector;
      address rampAddress = offRampRemoves[i].offRamp;

      if (s_offRamps.get(rampAddress) != uint256(rampSelector)) revert OffRampMismatch();

      if (s_offRamps.remove(rampAddress)) {
        emit OffRampRemoved(rampSelector, rampAddress);
      }
    }
    for (uint256 i = 0; i < offRampAdds.length; ++i) {
      uint64 rampSelector = offRampAdds[i].sourceChainSelector;
      address rampAddress = offRampAdds[i].offRamp;

      if (s_offRamps.set(rampAddress, rampSelector)) {
        emit OffRampAdded(rampSelector, rampAddress);
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
      (bool success, ) = to.call{value: amount}("");
      if (!success) revert FailedToSendValue();
      return;
    }
    IERC20(tokenAddress).safeTransfer(to, amount);
  }

  // ================================================================
  // │                           Access                             │
  // ================================================================

  /// @notice only lets permissioned offRamps execute
  /// @dev We additionally restrict offRamps to specific source chains for defense in depth.
  modifier onlyOffRamp(uint64 expectedSourceChainSelector) {
    (bool exists, uint256 sourceChainSelector) = s_offRamps.tryGet(msg.sender);
    if (!exists || expectedSourceChainSelector != uint64(sourceChainSelector)) revert OnlyOffRamp();
    _;
  }

  /// @notice Ensure that the ARM has not emitted a bad signal, and that the latest heartbeat is not stale.
  modifier whenHealthy() {
    if (IARM(i_armProxy).isCursed()) revert BadARMSignal();
    _;
  }
}
