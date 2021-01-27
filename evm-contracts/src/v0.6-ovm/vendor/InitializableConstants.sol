pragma solidity ^0.6.0;

/**
 * @dev Workaround to expose these storage constants, used to define/compute
 * reserved storage slots for Initializable contract, to any Proxy children
 * delegating calls to a contract that is Initializable.
 */
contract InitializableConstants {

  // First free storage slot for Initializable contract child
  uint8 internal constant STORAGE_FREE_POINTER = 50;
  // Amount of reserved storage starting from slot 0, to be used by Proxy
  // contract which is delegating calls to a contract that is Initializable.
  uint8 internal constant STORAGE_RESERVED_PROXY = 30;
}
