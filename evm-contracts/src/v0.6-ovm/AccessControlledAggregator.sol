pragma solidity ^0.6.0;

import "./vendor/Proxy.sol";
import "./interfaces/AggregatorV2V3Interface.sol";
import "./SimpleReadAccessController.sol";

/**
 * @title AccessControlled FluxAggregator contract
 * @notice This contract requires addresses to be added to a controller
 * in order to read the answers stored in the FluxAggregator contract
 */
contract AccessControlledAggregator is Proxy, SimpleReadAccessController {

  // Underlying implementation we delegate calls to
  address immutable public implementation;
  // We check access for these functions before delegating
  mapping(bytes4 => bool) private guardedFunctions;

  /**
   * @notice Creates an AccessControlledAggregator contract.
   * @dev The underlying aggregator implementation must extend Initializable, as this proxy depends on reserved ([0,49] slots) storage space.
   * @param implementationAddr The address of underlying aggregator implementation used by the proxy.
   */
  constructor(address implementationAddr) public {
    implementation = implementationAddr;
    _setUpGuards();
  }

  /**
   * @dev Set up guarded functions
   */
  function _setUpGuards() internal virtual {
    AggregatorV2V3Interface i;
    // AggregatorV3Interface
    guardedFunctions[i.getRoundData.selector] = true;
    guardedFunctions[i.latestRoundData.selector] = true;
    // AggregatorInterface
    guardedFunctions[i.latestAnswer.selector] = true;
    guardedFunctions[i.latestRound.selector] = true;
    guardedFunctions[i.latestTimestamp.selector] = true;
    guardedFunctions[i.getAnswer.selector] = true;
    guardedFunctions[i.getTimestamp.selector] = true;
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
    super._willFallback();
    // Check if one of the guarded functions
    if (guardedFunctions[msg.sig]) {
      require(hasAccess(msg.sender, msg.data), "No access");
    }
  }

  /**
   * @dev Storage slots for the Owned contract storage.
   * As we expect the implementation contract to extend both Initializable ([0,49] slots reserved) and Owned.
   * As this proxy contract is also Owned, we need to make sure that Owned storage in this contract
   * alligns with Owned storage in the underlying contract.
   */
  uint8 internal constant SLOT_OWNER = 50;
  uint8 internal constant SLOT_PENDIG_OWNER = 51; // SLOT_OWNER + 1

  /**
   * @return owner - The owner slot.
   */
  function _owner() internal override virtual view returns (address owner) {
    assembly {
      owner := sload(SLOT_OWNER)
    }
  }

  /**
   * @dev Sets the address of the owner.
   * @param newOwner Address of the new owner.
   */
  function _setOwner(address newOwner) internal override virtual {
    assembly {
      sstore(SLOT_OWNER, newOwner)
    }
  }

  /**
   * @return pendingOwner - The pending owner slot.
   */
  function _pendingOwner() internal override virtual view returns (address pendingOwner) {
    assembly {
      pendingOwner := sload(SLOT_PENDIG_OWNER)
    }
  }

  /**
   * @dev Sets the address of the pending owner.
   * @param newPendingOwner Address of the new pending owner.
   */
  function _setPendingOwner(address newPendingOwner) internal override virtual {
    assembly {
      sstore(SLOT_PENDIG_OWNER, newPendingOwner)
    }
  }
}
