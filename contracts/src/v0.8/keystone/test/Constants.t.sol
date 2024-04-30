// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Constants {
  address internal ADMIN = address(1);
  address internal STRANGER = address(2);
  address internal NODE_OPERATOR_ONE_ADMIN = address(3);
  string internal NODE_OPERATOR_ONE_NAME = "node-operator-one";
  address internal NODE_OPERATOR_TWO_ADMIN = address(4);
  string internal NODE_OPERATOR_TWO_NAME = "node-operator-two";

  bytes internal P2P_ID = bytes("12D3KooWRAw36ARW7T81yb7Ss5WPqGV7AnLcTmK1nApkbMS6s9cx");
}
