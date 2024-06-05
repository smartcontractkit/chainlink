// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Constants {
  address internal ADMIN = address(1);
  address internal STRANGER = address(2);

  address internal NODE_OPERATOR_ONE_ADMIN = address(3);
  string internal NODE_OPERATOR_ONE_NAME = "node-operator-one";
  bytes32 internal NODE_OPERATOR_ONE_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(3333)));
  bytes32 internal P2P_ID = hex"e42415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";

  address internal NODE_OPERATOR_TWO_ADMIN = address(4);
  string internal NODE_OPERATOR_TWO_NAME = "node-operator-two";
  bytes32 internal NODE_OPERATOR_TWO_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(4444)));
  bytes32 internal P2P_ID_TWO = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";

  address internal NODE_OPERATOR_THREE = address(4);
  string internal NODE_OPERATOR_THREE_NAME = "node-operator-three";
  bytes32 internal NODE_OPERATOR_THREE_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(5555)));
  bytes32 internal P2P_ID_THREE = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a06";

  uint32 internal F_VALUE = 1;
}
