// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IOwner {
  /// @notice Returns the owner of the contract.
  /// @dev This method is named to match with the OpenZeppelin Ownable contract.
  function owner() external view returns (address);
}
