// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.addCapability(s_basicCapability);
  }

  function test_RevertWhen_CapabilityExists() public {
    // Successfully add the capability the first time
    s_capabilityRegistry.addCapability(s_basicCapability);

    // Try to add the same capability again
    vm.expectRevert(CapabilityRegistry.CapabilityAlreadyExists.selector);
    s_capabilityRegistry.addCapability(s_basicCapability);
  }

  function test_RevertWhen_ConfigurationContractNotDeployed() public {
    address nonExistentContract = address(1);
    s_capabilityWithConfigurationContract.configurationContract = nonExistentContract;

    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilityRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        nonExistentContract
      )
    );
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_RevertWhen_ConfigurationContractDoesNotMatchInterface() public {
    CapabilityRegistry contractWithoutERC165 = new CapabilityRegistry();

    vm.expectRevert();
    s_capabilityWithConfigurationContract.configurationContract = address(contractWithoutERC165);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_AddCapability_NoConfigurationContract() public {
    s_capabilityRegistry.addCapability(s_basicCapability);

    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      bytes32("data-streams-reports"),
      bytes32("1.0.0")
    );
    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(hashedCapabilityId);

    assertEq(storedCapability.labelledName, s_basicCapability.labelledName);
    assertEq(storedCapability.version, s_basicCapability.version);
    assertEq(uint256(storedCapability.responseType), uint256(s_basicCapability.responseType));
    assertEq(storedCapability.configurationContract, s_basicCapability.configurationContract);
  }

  function test_AddCapability_WithConfiguration() public {
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      bytes32(s_capabilityWithConfigurationContract.labelledName),
      bytes32(s_capabilityWithConfigurationContract.version)
    );
    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(hashedCapabilityId);

    assertEq(storedCapability.labelledName, s_capabilityWithConfigurationContract.labelledName);
    assertEq(storedCapability.version, s_capabilityWithConfigurationContract.version);
    assertEq(uint256(storedCapability.responseType), uint256(s_capabilityWithConfigurationContract.responseType));
    assertEq(storedCapability.configurationContract, s_capabilityWithConfigurationContract.configurationContract);
  }
}
