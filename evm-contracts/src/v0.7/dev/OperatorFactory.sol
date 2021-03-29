// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./Operator.sol";

/**
 * @title Operator Factory
 * @notice Creates Operator contracts for node operators
 */
contract OperatorFactory {

  address public link;

  event OperatorCreated(
    address indexed operator,
    address indexed owner
  );

  /**
   * @param linkAddress address
   */
  constructor(
    address linkAddress
  ) {
    link = linkAddress;
  }

  /**
   * @notice fallback to create a new Operator contract with the msg.sender as owner
   */
  fallback() external {
    Operator operator = new Operator(link, msg.sender);
    emit OperatorCreated(address(operator), msg.sender);
  }
}
