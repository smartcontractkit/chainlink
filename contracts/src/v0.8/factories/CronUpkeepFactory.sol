// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../upkeeps/CronUpkeep.sol";
import "../upkeeps/CronUpkeepDelegate.sol";

/**
 * @title The CronUpkeepFactory contract
 * @notice This contract serves as a delegate for all instances of CronUpkeep. Those contracts
 * delegate their checkUpkeep calls onto this contract. Utilizing this pattern reduces the size
 * of the CronUpkeep contracts.
 */
contract CronUpkeepFactory {
  event NewCronUpkeepCreated(address upkeep, address owner);

  address private immutable s_cronDelegate;

  constructor() {
    s_cronDelegate = address(new CronUpkeepDelegate());
  }

  /**
   * @notice Creates a new CronUpkeep contract, with msg.sender as the owner
   */
  function newCronUpkeep() public {
    emit NewCronUpkeepCreated(address(new CronUpkeep(msg.sender, s_cronDelegate)), msg.sender);
  }

  /**
   * @notice Gets the address of the delegate contract
   * @return the address of the delegate contract
   */
  function cronDelegateAddress() public view returns (address) {
    return s_cronDelegate;
  }
}
