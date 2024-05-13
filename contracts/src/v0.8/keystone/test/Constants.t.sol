// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Constants {
  address internal ADMIN = address(1);
  address internal STRANGER = address(2);
  address internal NODE_OPERATOR_ONE_ADMIN = address(3);
  string internal NODE_OPERATOR_ONE_NAME = "node-operator-one";
  address internal NODE_OPERATOR_ONE_SIGNER_ADDRESS = address(3333);
  address internal NODE_OPERATOR_TWO_ADMIN = address(4);
  string internal NODE_OPERATOR_TWO_NAME = "node-operator-two";
  address internal NODE_OPERATOR_TWO_SIGNER_ADDRESS = address(4444);

  bytes32 internal P2P_ID = hex"e42415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";
  bytes32 internal P2P_ID_TWO = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";
}
