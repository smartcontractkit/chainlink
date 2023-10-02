pragma solidity 0.8.6;

import {Cron as CronInternal, Spec} from "../internal/Cron.sol";

/**
 * @title The Cron library
 * @notice A utility contract for encoding/decoding cron strings (ex: 0 0 * * *) into an
 * abstraction called a Spec. The library also includes a spec function, nextTick(), which
 * determines the next time a cron job should fire based on the current block timestamp.
 * @dev this is the external version of the library, which relies on the internal library
 * by the same name.
 */
library Cron {
  using CronInternal for Spec;
  using CronInternal for string;

  /**
   * @notice nextTick calculates the next datetime that a spec "ticks", starting
   * from the current block timestamp. This is gas-intensive and therefore should
   * only be called off-chain.
   * @param spec the spec to evaluate
   * @return the next tick
   */
  function nextTick(Spec calldata spec) public view returns (uint256) {
    return spec.nextTick();
  }

  /**
   * @notice lastTick calculates the previous datetime that a spec "ticks", starting
   * from the current block timestamp. This is gas-intensive and therefore should
   * only be called off-chain.
   * @param spec the spec to evaluate
   * @return the next tick
   */
  function lastTick(Spec calldata spec) public view returns (uint256) {
    return spec.lastTick();
  }

  /**
   * @notice matches evaluates whether or not a spec "ticks" at a given timestamp
   * @param spec the spec to evaluate
   * @param timestamp the timestamp to compare against
   * @return true / false if they match
   */
  function matches(Spec calldata spec, uint256 timestamp) public view returns (bool) {
    return spec.matches(timestamp);
  }

  /**
   * @notice toSpec converts a cron string to a spec struct. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param cronString the cron string
   * @return the spec struct
   */
  function toSpec(string calldata cronString) public pure returns (Spec memory) {
    return cronString.toSpec();
  }

  /**
   * @notice toEncodedSpec converts a cron string to an abi-encoded spec. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param cronString the cron string
   * @return the abi-encoded spec
   */
  function toEncodedSpec(string calldata cronString) public pure returns (bytes memory) {
    return cronString.toEncodedSpec();
  }

  /**
   * @notice toCronString converts a cron spec to a human-readable cron string. This is gas-intensive
   * and therefore should only be called off-chain.
   * @param spec the cron spec
   * @return the corresponding cron string
   */
  function toCronString(Spec calldata spec) public pure returns (string memory) {
    return spec.toCronString();
  }
}
