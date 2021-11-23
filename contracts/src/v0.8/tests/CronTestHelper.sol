// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import {Cron as CronInternal, Spec} from "../libraries/internal/Cron.sol";
import {Cron as CronExternal} from "../libraries/external/Cron.sol";

/**
 * @title The CronInternalTestHelper contract
 * @notice This contract exposes core functionality of the internal/Cron library.
 * It is only intended for use in tests.
 */
contract CronInternalTestHelper {
  /**
   * @notice Converts a cron string to a Spec, validates the spec, and encodes the spec.
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to convert and encode
   * @return the abi encoding of the Spec struct representing the cron string
   */
  function encodeCronString(string memory cronString) external pure returns (bytes memory) {
    return CronInternal.toEncodedSpec(cronString);
  }

  /**
   * @notice encodedSpecToString is a helper function for turning an
   * encoded spec back into a string. There is limited or no use for this outside
   * of tests.
   */
  function encodedSpecToString(bytes memory encodedSpec) public pure returns (string memory) {
    Spec memory spec = abi.decode(encodedSpec, (Spec));
    return CronInternal.toCronString(spec);
  }

  /**
   * @notice encodedSpecToString is a helper function for turning a string
   * into a spec struct.
   */
  function cronStringtoEncodedSpec(string memory cronString) public pure returns (Spec memory) {
    return CronInternal.toSpec(cronString);
  }

  /**
   * @notice calculateNextTick calculates the next time a cron job should "tick".
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to consider
   * @return the timestamp in UTC of the next "tick"
   */
  function calculateNextTick(string memory cronString) external view returns (uint256) {
    return CronInternal.nextTick(CronInternal.toSpec(cronString));
  }

  /**
   * @notice calculateLastTick calculates the last time a cron job "ticked".
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to consider
   * @return the timestamp in UTC of the last "tick"
   */
  function calculateLastTick(string memory cronString) external view returns (uint256) {
    return CronInternal.lastTick(CronInternal.toSpec(cronString));
  }
}

/**
 * @title The CronExternalTestHelper contract
 * @notice This contract exposes core functionality of the external/Cron library.
 * It is only intended for use in tests.
 */
contract CronExternalTestHelper {
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
   * @notice encodedSpecToString is a helper function for turning an
   * encoded spec back into a string. There is limited or no use for this outside
   * of tests.
   */
  function encodedSpecToString(bytes memory encodedSpec) public pure returns (string memory) {
    Spec memory spec = abi.decode(encodedSpec, (Spec));
    return CronExternal.toCronString(spec);
  }

  /**
   * @notice encodedSpecToString is a helper function for turning a string
   * into a spec struct.
   */
  function cronStringtoEncodedSpec(string memory cronString) public pure returns (Spec memory) {
    return CronExternal.toSpec(cronString);
  }

  /**
   * @notice calculateNextTick calculates the next time a cron job should "tick".
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to consider
   * @return the timestamp in UTC of the next "tick"
   */
  function calculateNextTick(string memory cronString) external view returns (uint256) {
    return CronExternal.nextTick(CronExternal.toSpec(cronString));
  }

  /**
   * @notice calculateLastTick calculates the last time a cron job "ticked".
   * This should only be called off-chain, as it is gas expensive!
   * @param cronString the cron string to consider
   * @return the timestamp in UTC of the last "tick"
   */
  function calculateLastTick(string memory cronString) external view returns (uint256) {
    return CronExternal.lastTick(CronExternal.toSpec(cronString));
  }
}
