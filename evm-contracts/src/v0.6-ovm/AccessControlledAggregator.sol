// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "./vendor/InitializableConstants.sol";
import "./vendor/Proxy.sol";
import "./interfaces/AggregatorV2V3Interface.sol";
import "./SimpleReadAccessController.sol";

/**
 * @title AccessControlled FluxAggregator contract
 * @notice This contract requires addresses to be added to a controller
 * in order to read the answers stored in the FluxAggregator contract
 */
contract AccessControlledAggregator is InitializableConstants, Proxy, SimpleReadAccessController {

  // Underlying implementation we delegate calls to
  address immutable public implementation;

  /**
   * @notice Creates an AccessControlledAggregator contract.
   * @dev The underlying aggregator implementation must extend Initializable, as this proxy depends on reserved storage space.
   * @param implementationAddr The address of underlying aggregator implementation used by the proxy.
   */
  constructor(address implementationAddr) public {
    implementation = implementationAddr;
  }

  /**
   * @return The Address of the implementation.
   */
  function _implementation() internal view  override returns (address) {
    return implementation;
  }

  /**
   * @dev Function that is run as the first thing in the fallback function.
   * Can be redefined in derived contracts to add functionality.
   * Redefinitions must call super._willFallback().
   */
  function _willFallback() internal override virtual {
    // We check access for these functions before delegating
    AggregatorV2V3Interface i;
    // Unrolled loop to optimize gas cost
    if (msg.sig == i.getRoundData.selector
      || msg.sig == i.latestRoundData.selector
      // AggregatorInterface
      || msg.sig == i.latestAnswer.selector
      || msg.sig == i.latestRound.selector
      || msg.sig == i.latestTimestamp.selector
      || msg.sig == i.getAnswer.selector
      || msg.sig == i.getTimestamp.selector) {
      require(hasAccess(msg.sender, msg.data), "No access");
    }
  }

  /**
   * @dev Storage slots for the Owned contract storage.
   * We expect the implementation contract to extend both Initializable and Owned.
   * As this proxy contract is also Owned, we need to make sure that Owned storage in this contract
   * aligns with Owned storage in the underlying contract.
   */
  uint8 internal constant STORAGE_SLOT_OWNER = STORAGE_FREE_POINTER;
  uint8 internal constant STORAGE_SLOT_PENDIG_OWNER = STORAGE_SLOT_OWNER + 1;

  /**
   * @return owner - The owner slot.
   */
  function _owner() internal override virtual view returns (address owner) {
    uint8 slot = STORAGE_SLOT_OWNER;
    assembly {
      owner := sload(slot)
    }
  }

  /**
   * @dev Sets the address of the owner.
   * @param newOwner Address of the new owner.
   */
  function _setOwner(address newOwner) internal override virtual {
    uint8 slot = STORAGE_SLOT_OWNER;
    assembly {
      sstore(slot, newOwner)
    }
  }

  /**
   * @return pendingOwner - The pending owner slot.
   */
  function _pendingOwner() internal override virtual view returns (address pendingOwner) {
    uint8 slot = STORAGE_SLOT_PENDIG_OWNER;
    assembly {
      pendingOwner := sload(slot)
    }
  }

  /**
   * @dev Sets the address of the pending owner.
   * @param newPendingOwner Address of the new pending owner.
   */
  function _setPendingOwner(address newPendingOwner) internal override virtual {
    uint8 slot = STORAGE_SLOT_PENDIG_OWNER;
    assembly {
      sstore(slot, newPendingOwner)
    }
  }
}
