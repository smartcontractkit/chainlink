// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title Chainlink Proof-of-Reserve address list interface.
 * @notice This interface enables Chainlink nodes to get the list addresses to be used in a PoR feed. A single
 * contract that implements this interface can only store an address list for a single PoR feed.
 * @dev All functions in this interface are expected to be called off-chain, so gas usage is not a big concern.
 * This makes it possible to store addresses in optimized data types and convert them to human-readable strings
 * in `getPoRAddressList()`.
 */
interface PoRAddressList {
  /// @notice Get total number of addresses in the list.
  function getPoRAddressListLength() external view returns (uint256);

  /**
   * @notice Get a batch of human-readable addresses from the address list. The requested batch size can be greater
   * than the actual address list size, in which the full address list will be returned.
   * @dev Due to limitations of gas usage in off-chain calls, we need to support fetching the addresses in batches.
   * EVM addresses need to be converted to human-readable strings. The address strings need to be in the same format
   * that would be used when querying the balance of that address.
   * @param startIndex The index of the first address in the batch.
   * @param endIndex The index of the last address in the batch. If `endIndex > getPoRAddressListLength()-1`,
   * endIndex need to default to `getPoRAddressListLength()-1`.
   * @return Array of addresses as strings.
   */
  function getPoRAddressList(uint256 startIndex, uint256 endIndex) external view returns (string[] memory);
}
