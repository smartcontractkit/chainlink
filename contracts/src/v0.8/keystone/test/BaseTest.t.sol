// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Constants} from "./Constants.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract BaseTest is Test, Constants {
  CapabilityRegistry internal s_capabilityRegistry;
  CapabilityConfigurationContract internal s_capabilityConfigurationContract;
  CapabilityRegistry.Capability internal s_basicCapability;
  CapabilityRegistry.Capability internal s_capabilityWithConfigurationContract;
  bytes32 internal s_basicHashedCapabilityId;
  bytes32 internal s_capabilityWithConfigurationContractId;
  bytes32 internal s_nonExistentHashedCapabilityId;

  function setUp() public virtual {
    vm.startPrank(ADMIN);
    s_capabilityRegistry = new CapabilityRegistry();
    s_capabilityConfigurationContract = new CapabilityConfigurationContract();

    s_basicCapability = CapabilityRegistry.Capability({
      labelledName: "data-streams-reports",
      version: "1.0.0",
      responseType: CapabilityRegistry.CapabilityResponseType.REPORT,
      configurationContract: address(0)
    });
    s_capabilityWithConfigurationContract = CapabilityRegistry.Capability({
      labelledName: "read-ethereum-mainnet-gas-price",
      version: "1.0.2",
      responseType: CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL,
      configurationContract: address(s_capabilityConfigurationContract)
    });

    s_basicHashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    s_capabilityWithConfigurationContractId = s_capabilityRegistry.getHashedCapabilityId(
      s_capabilityWithConfigurationContract.labelledName,
      s_capabilityWithConfigurationContract.version
    );
    s_nonExistentHashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId("non-existent-capability", "1.0.0");
  }

  function _getNodeOperators() internal view returns (CapabilityRegistry.NodeOperator[] memory) {
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](2);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NODE_OPERATOR_ONE_ADMIN, name: NODE_OPERATOR_ONE_NAME});
    nodeOperators[1] = CapabilityRegistry.NodeOperator({admin: NODE_OPERATOR_TWO_ADMIN, name: NODE_OPERATOR_TWO_NAME});
    return nodeOperators;
  }
}
