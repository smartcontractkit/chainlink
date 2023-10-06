// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IVRFV2PlusMigration {
  function migrate(uint256 subId, address newCoordinator) external;
}
