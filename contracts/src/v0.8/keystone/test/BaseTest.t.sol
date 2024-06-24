// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Constants} from "./Constants.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract BaseTest is Test, Constants {
  CapabilitiesRegistry internal s_CapabilitiesRegistry;
  CapabilityConfigurationContract internal s_capabilityConfigurationContract;
  CapabilitiesRegistry.Capability internal s_basicCapability;
  CapabilitiesRegistry.Capability internal s_capabilityWithConfigurationContract;
  bytes32 internal s_basicHashedCapabilityId;
  bytes32 internal s_capabilityWithConfigurationContractId;
  bytes32 internal s_nonExistentHashedCapabilityId;

  function setUp() public virtual {
    vm.startPrank(ADMIN);
    s_CapabilitiesRegistry = new CapabilitiesRegistry();
    s_capabilityConfigurationContract = new CapabilityConfigurationContract();

    s_basicCapability = CapabilitiesRegistry.Capability({
      labelledName: "data-streams-reports",
      version: "1.0.0",
      responseType: CapabilitiesRegistry.CapabilityResponseType.REPORT,
      configurationContract: address(0),
      capabilityType: CapabilitiesRegistry.CapabilityType.TRIGGER
    });
    s_capabilityWithConfigurationContract = CapabilitiesRegistry.Capability({
      labelledName: "read-ethereum-mainnet-gas-price",
      version: "1.0.2",
      responseType: CapabilitiesRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL,
      configurationContract: address(s_capabilityConfigurationContract),
      capabilityType: CapabilitiesRegistry.CapabilityType.ACTION
    });

    s_basicHashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    s_capabilityWithConfigurationContractId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_capabilityWithConfigurationContract.labelledName,
      s_capabilityWithConfigurationContract.version
    );
    s_nonExistentHashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId("non-existent-capability", "1.0.0");
  }

  function _getNodeOperators() internal pure returns (CapabilitiesRegistry.NodeOperator[] memory) {
    CapabilitiesRegistry.NodeOperator[] memory nodeOperators = new CapabilitiesRegistry.NodeOperator[](3);
    nodeOperators[0] = CapabilitiesRegistry.NodeOperator({
      admin: NODE_OPERATOR_ONE_ADMIN,
      name: NODE_OPERATOR_ONE_NAME
    });
    nodeOperators[1] = CapabilitiesRegistry.NodeOperator({
      admin: NODE_OPERATOR_TWO_ADMIN,
      name: NODE_OPERATOR_TWO_NAME
    });
    nodeOperators[2] = CapabilitiesRegistry.NodeOperator({admin: NODE_OPERATOR_THREE, name: NODE_OPERATOR_THREE_NAME});
    return nodeOperators;
  }
}
