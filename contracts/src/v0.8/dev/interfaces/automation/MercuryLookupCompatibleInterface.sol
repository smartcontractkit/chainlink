// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface MercuryLookupCompatibleInterface {
  /**
   * @notice any contract which wants to utilize MercuryLookup feature needs to
   * implement this interface as well as the automation compatible interface.
   * @param values an array of bytes returned from Mercury endpoint.
   * @param extraData context data from Mercury lookup process.
   * @return upkeepNeeded boolean to indicate whether the keeper should call
   * performUpkeep or not.
   * @return performData bytes that the keeper should call performUpkeep with, if
   * upkeep is needed. If you would like to encode data to decode later, try
   * `abi.encode`.
   */
  function mercuryCallback(bytes[] memory values, bytes memory extraData)
    external
    view
    returns (bool upkeepNeeded, bytes memory performData);
}
