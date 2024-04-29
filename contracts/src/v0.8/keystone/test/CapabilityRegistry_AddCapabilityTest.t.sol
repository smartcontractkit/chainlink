// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  CapabilityRegistry.Capability private basicCapability =
    CapabilityRegistry.Capability({
      capabilityType: "data-streams-reports",
      version: "1.0.0",
      responseType: CapabilityRegistry.CapabilityResponseType.REPORT,
      configurationContract: address(0)
    });

  CapabilityRegistry.Capability private capabilityWithConfigurationContract =
    CapabilityRegistry.Capability({
      capabilityType: "read-ethereum-mainnet-gas-price",
      version: "1.0.2",
      responseType: CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL,
      configurationContract: address(s_capabilityConfigurationContract)
    });

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.addCapability(basicCapability);
  }

  function test_RevertWhen_CapabilityExists() public {
    // Successfully add the capability the first time
    s_capabilityRegistry.addCapability(basicCapability);

    // Try to add the same capability again
    vm.expectRevert(CapabilityRegistry.CapabilityAlreadyExists.selector);
    s_capabilityRegistry.addCapability(basicCapability);
  }

  function test_RevertWhen_ConfigurationContractNotDeployed() public {
    address nonExistentContract = address(1);
    capabilityWithConfigurationContract.configurationContract = nonExistentContract;

    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilityRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        nonExistentContract
      )
    );
    s_capabilityRegistry.addCapability(capabilityWithConfigurationContract);
  }

  function test_RevertWhen_ConfigurationContractDoesNotMatchInterface() public {
    CapabilityRegistry contractWithoutERC165 = new CapabilityRegistry();

    vm.expectRevert();
    capabilityWithConfigurationContract.configurationContract = address(contractWithoutERC165);
    s_capabilityRegistry.addCapability(capabilityWithConfigurationContract);
  }

  function test_AddCapability_NoConfigurationContract() public {
    s_capabilityRegistry.addCapability(basicCapability);

    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(bytes32("data-streams-reports"), bytes32("1.0.0"));
    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(capabilityId);

    assertEq(storedCapability.capabilityType, basicCapability.capabilityType);
    assertEq(storedCapability.version, basicCapability.version);
    assertEq(uint256(storedCapability.responseType), uint256(basicCapability.responseType));
    assertEq(storedCapability.configurationContract, basicCapability.configurationContract);
  }

  function test_AddCapability_WithConfiguration() public {
    s_capabilityRegistry.addCapability(capabilityWithConfigurationContract);

    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(
      bytes32(capabilityWithConfigurationContract.capabilityType),
      bytes32(capabilityWithConfigurationContract.version)
    );
    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(capabilityId);

    assertEq(storedCapability.capabilityType, capabilityWithConfigurationContract.capabilityType);
    assertEq(storedCapability.version, capabilityWithConfigurationContract.version);
    assertEq(uint256(storedCapability.responseType), uint256(capabilityWithConfigurationContract.responseType));
    assertEq(storedCapability.configurationContract, capabilityWithConfigurationContract.configurationContract);
  }
}
