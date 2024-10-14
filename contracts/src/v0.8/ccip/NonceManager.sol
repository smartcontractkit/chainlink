// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IEVM2AnyOnRamp} from "./interfaces/IEVM2AnyOnRamp.sol";
import {INonceManager} from "./interfaces/INonceManager.sol";

import {AuthorizedCallers} from "../shared/access/AuthorizedCallers.sol";

/// @title NonceManager
/// @notice NonceManager contract that manages sender nonces for the on/off ramps
contract NonceManager is INonceManager, AuthorizedCallers, ITypeAndVersion {
  error PreviousRampAlreadySet();

  event PreviousRampsUpdated(uint64 indexed remoteChainSelector, PreviousRamps prevRamp);
  event SkippedIncorrectNonce(uint64 sourceChainSelector, uint64 nonce, bytes sender);

  /// @dev Struct that contains the previous on/off ramp addresses
  struct PreviousRamps {
    address prevOnRamp; // Previous onRamp
    address prevOffRamp; // Previous offRamp
  }

  /// @dev Struct that contains the chain selector and the previous on/off ramps, same as PreviousRamps but with the chain selector
  /// so that an array of these can be passed to the applyPreviousRampsUpdates function
  struct PreviousRampsArgs {
    uint64 remoteChainSelector; // Chain selector
    PreviousRamps prevRamps; // Previous on/off ramps
  }

  string public constant override typeAndVersion = "NonceManager 1.6.0-dev";

  /// @dev The previous on/off ramps per chain selector
  mapping(uint64 chainSelector => PreviousRamps previousRamps) private s_previousRamps;
  /// @dev The current outbound nonce per sender used on the onramp
  mapping(uint64 destChainSelector => mapping(address sender => uint64 outboundNonce)) private s_outboundNonces;
  /// @dev The current inbound nonce per sender used on the offramp
  /// Eventually in sync with the outbound nonce in the remote source chain NonceManager, used to enforce that messages are
  /// executed in the same order they are sent (assuming they are DON)
  mapping(uint64 sourceChainSelector => mapping(bytes sender => uint64 inboundNonce)) private s_inboundNonces;

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

  /// @notice Returns the outbound nonce for a given sender on a given destination chain.
  /// @param destChainSelector The destination chain selector.
  /// @param sender The sender address.
  /// @return outboundNonce The outbound nonce.
  function getOutboundNonce(uint64 destChainSelector, address sender) external view returns (uint64) {
    return _getOutboundNonce(destChainSelector, sender);
  }

  function _getOutboundNonce(uint64 destChainSelector, address sender) private view returns (uint64) {
    uint64 outboundNonce = s_outboundNonces[destChainSelector][sender];

    // When introducing the NonceManager with existing lanes, we still want to have sequential nonces.
    // Referencing the old onRamp preserves sequencing between updates.
    if (outboundNonce == 0) {
      address prevOnRamp = s_previousRamps[destChainSelector].prevOnRamp;
      if (prevOnRamp != address(0)) {
        return IEVM2AnyOnRamp(prevOnRamp).getSenderNonce(sender);
      }
    }

    return outboundNonce;
  }

  /// @inheritdoc INonceManager
  function incrementInboundNonce(
    uint64 sourceChainSelector,
    uint64 expectedNonce,
    bytes calldata sender
  ) external onlyAuthorizedCallers returns (bool) {
    uint64 inboundNonce = _getInboundNonce(sourceChainSelector, sender) + 1;

    if (inboundNonce != expectedNonce) {
      // If the nonce is not the expected one, this means that there are still messages in flight so we skip
      // the nonce increment
      emit SkippedIncorrectNonce(sourceChainSelector, expectedNonce, sender);
      return false;
    }

    s_inboundNonces[sourceChainSelector][sender] = inboundNonce;

    return true;
  }

  /// @notice Returns the inbound nonce for a given sender on a given source chain.
  /// @param sourceChainSelector The source chain selector.
  /// @param sender The encoded sender address.
  /// @return inboundNonce The inbound nonce.
  function getInboundNonce(uint64 sourceChainSelector, bytes calldata sender) external view returns (uint64) {
    return _getInboundNonce(sourceChainSelector, sender);
  }

  function _getInboundNonce(uint64 sourceChainSelector, bytes calldata sender) private view returns (uint64) {
    uint64 inboundNonce = s_inboundNonces[sourceChainSelector][sender];

    // When introducing the NonceManager with existing lanes, we still want to have sequential nonces.
    // Referencing the old offRamp to check the expected nonce if none is set for a
    // given sender allows us to skip the current message in the current offRamp if it would not be the next according
    // to the old offRamp. This preserves sequencing between updates.
    if (inboundNonce == 0) {
      address prevOffRamp = s_previousRamps[sourceChainSelector].prevOffRamp;
      if (prevOffRamp != address(0)) {
        // We only expect EVM previous offRamps here so we can safely decode the sender
        return IEVM2AnyOnRamp(prevOffRamp).getSenderNonce(abi.decode(sender, (address)));
      }
    }

    return inboundNonce;
  }

  /// @notice Updates the previous ramps addresses.
  /// @param previousRampsArgs The previous on/off ramps addresses.
  function applyPreviousRampsUpdates(PreviousRampsArgs[] calldata previousRampsArgs) external onlyOwner {
    for (uint256 i = 0; i < previousRampsArgs.length; ++i) {
      PreviousRampsArgs calldata previousRampsArg = previousRampsArgs[i];

      PreviousRamps storage prevRamps = s_previousRamps[previousRampsArg.remoteChainSelector];

      // If the previous ramps are already set then they should not be updated.
      // In versions prior to the introduction of the NonceManager contract, nonces were tracked in the on/off ramps.
      // This config does a 1-time migration to move the nonce from on/off ramps into NonceManager
      if (prevRamps.prevOnRamp != address(0) || prevRamps.prevOffRamp != address(0)) {
        revert PreviousRampAlreadySet();
      }

      prevRamps.prevOnRamp = previousRampsArg.prevRamps.prevOnRamp;
      prevRamps.prevOffRamp = previousRampsArg.prevRamps.prevOffRamp;

      emit PreviousRampsUpdated(previousRampsArg.remoteChainSelector, previousRampsArg.prevRamps);
    }
  }

  /// @notice Gets the previous onRamp address for the given chain selector
  /// @param chainSelector The chain selector
  /// @return previousRamps The previous on/offRamp addresses
  function getPreviousRamps(uint64 chainSelector) external view returns (PreviousRamps memory) {
    return s_previousRamps[chainSelector];
  }
}
