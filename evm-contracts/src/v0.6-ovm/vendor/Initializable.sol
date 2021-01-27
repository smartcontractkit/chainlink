pragma solidity ^0.6.0;

import "./InitializableConstants.sol";

/**
 * @title Initializable
 *
 * @dev Helper contract to support initializer functions. To use it, replace
 * the constructor with a function that has the `initializer` modifier.
 * WARNING: Unlike constructors, initializer functions must be manually
 * invoked. This applies both to deploying an Initializable contract, as well
 * as extending an Initializable contract via inheritance.
 * WARNING: When used with inheritance, manual care must be taken to not invoke
 * a parent initializer twice, or ensure that all initializers are idempotent,
 * because this is not dealt with automatically as with constructors.
 * WARNING: The contract reserves part of storage which is useful when used
 * in combination with the Proxy contract:
 *   - Slots [0..29] reserved for the Proxy contract (+ any children)
 *   - Slots [30..49] reserved for future layout changes of this Initializable contract
 *   - Slots [50..] reserved for the implementation contract (child of Initializable)
 */
contract Initializable is InitializableConstants {

  // If this Initializable contract storage is expanded, this constant
  // needs to be appropriately updated
  uint8 internal constant STORAGE_USED_INITIALIZABLE = 1;
  // Amount of reserved storage still available before the free pointer
  uint8 internal constant STORAGE_RESERVED_INITIALIZABLE = STORAGE_FREE_POINTER - STORAGE_RESERVED_PROXY - STORAGE_USED_INITIALIZABLE;

  // Reserved storage space for the Proxy contract
  uint256[STORAGE_RESERVED_PROXY] private _____gap;

  /**
   * @dev Indicates that the contract has been initialized.
   */
  bool private initialized;

  /**
   * @dev Indicates that the contract is in the process of being initialized.
   */
  bool private initializing;

  /**
   * @dev Modifier to use in the initializer function of a contract.
   */
  modifier initializer() {
    require(initializing || isConstructor() || !initialized, "Contract instance has already been initialized");

    bool isTopLevelCall = !initializing;
    if (isTopLevelCall) {
      initializing = true;
      initialized = true;
    }

    _;

    if (isTopLevelCall) {
      initializing = false;
    }
  }

  /// @dev Returns true if and only if the function is running in the constructor
  function isConstructor() private view returns (bool) {
    // extcodesize checks the size of the code stored in an address, and
    // address returns the current address. Since the code is still not
    // deployed when running a constructor, any checks on its code size will
    // yield zero, making it an effective way to detect if a contract is
    // under construction or not.
    address self = address(this);
    uint256 cs;
    assembly { cs := extcodesize(self) }
    return cs == 0;
  }

  // Reserved storage space to allow for layout changes in the future.
  uint256[STORAGE_RESERVED_INITIALIZABLE] private ______gap;
}
