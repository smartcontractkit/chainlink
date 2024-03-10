// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "./CronUpkeep.sol";
import "./CronUpkeepDelegate.sol";
import "../../shared/access/ConfirmedOwner.sol";
import {Spec, Cron as CronExternal} from "../libraries/external/Cron.sol";

/**
 * @title The CronUpkeepFactory contract
 * @notice This contract serves as a delegate for all instances of CronUpkeep. Those contracts
 * delegate their checkUpkeep calls onto this contract. Utilizing this pattern reduces the size
 * of the CronUpkeep contracts.
 */
contract CronUpkeepFactory is ConfirmedOwner {
  event NewCronUpkeepCreated(address upkeep, address owner);

  address private immutable s_cronDelegate;
  uint256 public s_maxJobs = 5;

  constructor() ConfirmedOwner(msg.sender) {
    s_cronDelegate = address(new CronUpkeepDelegate());
  }

  /**
   * @notice Creates a new CronUpkeep contract, with msg.sender as the owner
   */
  function newCronUpkeep() external {
    newCronUpkeepWithJob(bytes(""));
  }

  /**
   * @notice Creates a new CronUpkeep contract, with msg.sender as the owner, and registers a cron job
   */
  function newCronUpkeepWithJob(bytes memory encodedJob) public {
    emit NewCronUpkeepCreated(address(new CronUpkeep(msg.sender, s_cronDelegate, s_maxJobs, encodedJob)), msg.sender);
  }

  /**
   * @notice Sets the max job limit on new cron upkeeps
   */
  function setMaxJobs(uint256 maxJobs) external onlyOwner {
    s_maxJobs = maxJobs;
  }

  /**
   * @notice Gets the address of the delegate contract
   * @return the address of the delegate contract
   */
  function cronDelegateAddress() external view returns (address) {
    return s_cronDelegate;
  }

  /**
   * @notice Converts a cron string to a Spec, validates the spec, and encodes the spec.
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to convert and encode
   * @return the abi encoding of the Spec struct representing the cron string
   */
  function encodeCronString(string memory cronString) external pure returns (bytes memory) {
    return CronExternal.toEncodedSpec(cronString);
  }

  /**
   * @notice Converts, validates, and encodes a full cron spec. This payload is then passed to newCronUpkeepWithJob.
   * @param target the destination contract of a cron job
   * @param handler the function signature on the target contract to call
   * @param cronString the cron string to convert and encode
   * @return the abi encoding of the entire cron job
   */
  function encodeCronJob(
    address target,
    bytes memory handler,
    string memory cronString
  ) external pure returns (bytes memory) {
    Spec memory spec = CronExternal.toSpec(cronString);
    return abi.encode(target, handler, spec);
  }
}
