// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";

contract Constants {
  address internal ADMIN = address(1);
  address internal STRANGER = address(2);
  address internal NODE_OPERATOR_ONE_ADMIN = address(3);
  string internal NODE_OPERATOR_ONE_NAME = "node-operator-one";
  address internal NODE_OPERATOR_TWO_ADMIN = address(4);
  string internal NODE_OPERATOR_TWO_NAME = "node-operator-two";
}
