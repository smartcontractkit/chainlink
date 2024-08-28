// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ICapabilityConfiguration} from "../../interfaces/ICapabilityConfiguration.sol";
import {CapabilitiesRegistry} from "../../CapabilitiesRegistry.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Constants} from "../Constants.t.sol";

contract MaliciousConfigurationContract is ICapabilityConfiguration, IERC165, Constants {
  bytes32 internal s_capabilityWithConfigurationContractId;

  constructor(bytes32 capabilityWithConfigContractId) {
    s_capabilityWithConfigurationContractId = capabilityWithConfigContractId;
  }

  function getCapabilityConfiguration(uint32) external pure returns (bytes memory configuration) {
    return bytes("");
  }

  function beforeCapabilityConfigSet(bytes32[] calldata, bytes calldata, uint64, uint32) external {
    CapabilitiesRegistry.NodeParams[] memory nodes = new CapabilitiesRegistry.NodeParams[](2);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);

    hashedCapabilityIds[0] = s_capabilityWithConfigurationContractId;

    // Set node one's signer to another address
    nodes[0] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    nodes[1] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID_THREE,
      signer: NODE_OPERATOR_THREE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    CapabilitiesRegistry(msg.sender).updateNodes(nodes);
  }

  function supportsInterface(bytes4 interfaceId) public pure returns (bool) {
    return interfaceId == type(ICapabilityConfiguration).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
