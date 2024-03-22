// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

/**
 * @title Chainable - the contract size limit nullifier
 * @notice Chainable is designed to link together a "chain" of contracts through fallback functions
 * and delegatecalls. All code is executed in the context of the head of the chain, the "master" contract.
 */
contract Chainable {
  /**
   * @dev addresses of the next contract in the chain **have to be immutable/constant** or the system won't work
   */
  address private immutable i_FALLBACK_ADDRESS;

  /**
   * @param fallbackAddress the address of the next contract in the chain
   */
  constructor(address fallbackAddress) {
    i_FALLBACK_ADDRESS = fallbackAddress;
  }

  /**
   * @notice returns the address of the next contract in the chain
   */
  function fallbackTo() external view returns (address) {
    return i_FALLBACK_ADDRESS;
  }

  /**
   * @notice the fallback function routes the call to the next contract in the chain
   * @dev most of the implementation is copied directly from OZ's Proxy contract
   */
  // solhint-disable-next-line no-complex-fallback
  fallback() external payable {
    // copy to memory for assembly access
    address next = i_FALLBACK_ADDRESS;
    // copied directly from OZ's Proxy contract
    assembly {
      // Copy msg.data. We take full control of memory in this inline assembly
      // block because it will not return to Solidity code. We overwrite the
      // Solidity scratch pad at memory position 0.
      calldatacopy(0, 0, calldatasize())

      // Call the next contract.
      // out and outsize are 0 because we don't know the size yet.
      let result := delegatecall(gas(), next, 0, calldatasize(), 0, 0)

      // Copy the returned data.
      returndatacopy(0, 0, returndatasize())

      switch result
      // delegatecall returns 0 on error.
      case 0 {
        revert(0, returndatasize())
      }
      default {
        return(0, returndatasize())
      }
    }
  }
}
