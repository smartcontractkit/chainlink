// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract Constants {
  address internal constant ADMIN = address(1);
  address internal constant STRANGER = address(2);

  uint32 internal constant TEST_NODE_OPERATOR_ONE_ID = 1;
  address internal constant NODE_OPERATOR_ONE_ADMIN = address(3);
  string internal constant NODE_OPERATOR_ONE_NAME = "node-operator-one";
  bytes32 internal constant NODE_OPERATOR_ONE_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(3333)));
  bytes32 internal constant P2P_ID = hex"e42415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";

  uint32 internal constant TEST_NODE_OPERATOR_TWO_ID = 2;
  address internal constant NODE_OPERATOR_TWO_ADMIN = address(4);
  string internal constant NODE_OPERATOR_TWO_NAME = "node-operator-two";
  bytes32 internal constant NODE_OPERATOR_TWO_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(4444)));
  bytes32 internal constant P2P_ID_TWO = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a05";

  uint32 internal constant TEST_NODE_OPERATOR_THREE_ID = 3;
  address internal constant NODE_OPERATOR_THREE = address(4);
  string internal constant NODE_OPERATOR_THREE_NAME = "node-operator-three";
  bytes32 internal constant NODE_OPERATOR_THREE_SIGNER_ADDRESS = bytes32(abi.encodePacked(address(5555)));
  bytes32 internal constant P2P_ID_THREE = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a06";

  uint32 internal constant F_VALUE = 1;
  uint32 internal constant DON_ID = 1;
  uint32 internal constant DON_ID_TWO = 2;

  bytes32 internal constant INVALID_P2P_ID = bytes32("fake-p2p");
  bytes32 internal constant NEW_NODE_SIGNER = hex"f53415859707d90ed4dc534ad730f187a17b0c368e1beec2e9b995587c4b0a07";

  bytes internal constant BASIC_CAPABILITY_CONFIG = bytes("basic-capability-config");
  bytes internal constant CONFIG_CAPABILITY_CONFIG = bytes("config-capability-config");
}
