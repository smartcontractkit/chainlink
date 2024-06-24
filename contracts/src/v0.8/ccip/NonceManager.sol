// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IEVM2AnyOnRamp} from "./interfaces/IEVM2AnyOnRamp.sol";
import {INonceManager} from "./interfaces/INonceManager.sol";

import {AuthorizedCallers} from "../shared/access/AuthorizedCallers.sol";

/// @title NonceManager
/// @notice NonceManager contract that manages sender nonces for the on/off ramps
contract NonceManager is INonceManager, AuthorizedCallers {
  error PreviousRampAlreadySet();

  event PreviousOnRampUpdated(uint64 indexed destChainSelector, address prevOnRamp);

  /// @dev Struct that contains the previous on/off ramp addresses
  // TODO: add prevOffRamp
  struct PreviousRamps {
    address prevOnRamp; // Previous onRamp
  }

  /// @dev Struct that contains the chain selector and the previous on/off ramps, same as PreviousRamps but with the chain selector
  /// so that an array of these can be passed to the applyPreviousRampsUpdates function
  struct PreviousRampsArgs {
    uint64 remoteChainSelector; // Chain selector
    PreviousRamps prevRamps; // Previous on/off ramps
  }

  /// @dev previous ramps
  mapping(uint64 chainSelector => PreviousRamps previousRamps) private s_previousRamps;
  /// @dev The current outbound nonce per sender used on the onramp
  mapping(uint64 destChainSelector => mapping(address sender => uint64 outboundNonce)) private s_outboundNonces;

  constructor(address[] memory authorizedCallers) AuthorizedCallers(authorizedCallers) {}

  /// @inheritdoc INonceManager
  function getIncrementedOutboundNonce(
    uint64 destChainSelector,
    address sender
  ) external onlyAuthorizedCallers returns (uint64) {
    uint64 outboundNonce = _getOutboundNonce(destChainSelector, sender) + 1;
    s_outboundNonces[destChainSelector][sender] = outboundNonce;

    return outboundNonce;
  }

  /// TODO: add incrementInboundNonce

  /// @notice Returns the outbound nonce for the given sender on the given destination chain
  /// @param destChainSelector The destination chain selector
  /// @param sender The sender address
  /// @return The outbound nonce
  function getOutboundNonce(uint64 destChainSelector, address sender) external view returns (uint64) {
    return _getOutboundNonce(destChainSelector, sender);
  }

  function _getOutboundNonce(uint64 destChainSelector, address sender) private view returns (uint64) {
    uint64 outboundNonce = s_outboundNonces[destChainSelector][sender];

    if (outboundNonce == 0) {
      address prevOnRamp = s_previousRamps[destChainSelector].prevOnRamp;
      if (prevOnRamp != address(0)) {
        return IEVM2AnyOnRamp(prevOnRamp).getSenderNonce(sender);
      }
    }

    return outboundNonce;
  }

  /// TODO: add getInboundNonce

  /// @notice Updates the previous ramps addresses
  /// @param previousRampsArgs The previous on/off ramps addresses
  function applyPreviousRampsUpdates(PreviousRampsArgs[] calldata previousRampsArgs) external onlyOwner {
    for (uint256 i = 0; i < previousRampsArgs.length; ++i) {
      PreviousRampsArgs calldata previousRampsArg = previousRampsArgs[i];

      PreviousRamps storage prevRamps = s_previousRamps[previousRampsArg.remoteChainSelector];

      // If the previous onRamp is already set then it should not be updated
      if (prevRamps.prevOnRamp != address(0)) {
        revert PreviousRampAlreadySet();
      }

      prevRamps.prevOnRamp = previousRampsArg.prevRamps.prevOnRamp;
      emit PreviousOnRampUpdated(previousRampsArg.remoteChainSelector, prevRamps.prevOnRamp);
    }
  }

  /// @notice Gets the previous onRamp address for the given chain selector
  /// @param chainSelector The chain selector
  /// @return The previous onRamp address
  function getPreviousRamps(uint64 chainSelector) external view returns (PreviousRamps memory) {
    return s_previousRamps[chainSelector];
  }
}
