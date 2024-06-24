// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

contract CapabilitiesRegistry_AddCapabilitiesTest is BaseTest {
  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);

    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_basicCapability;

    vm.expectRevert("Only callable by owner");
    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_CapabilityExists() public {
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_basicCapability;

    // Successfully add the capability the first time
    s_CapabilitiesRegistry.addCapabilities(capabilities);

    // Try to add the same capability again
    vm.expectRevert(
      abi.encodeWithSelector(CapabilitiesRegistry.CapabilityAlreadyExists.selector, s_basicHashedCapabilityId)
    );
    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_ConfigurationContractNotDeployed() public {
    address nonExistentContract = address(1);
    s_capabilityWithConfigurationContract.configurationContract = nonExistentContract;

    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_capabilityWithConfigurationContract;

    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilitiesRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        nonExistentContract
      )
    );
    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_ConfigurationContractDoesNotMatchInterface() public {
    address contractWithoutERC165 = address(9999);
    vm.mockCall(
      contractWithoutERC165,
      abi.encodeWithSelector(
        IERC165.supportsInterface.selector,
        ICapabilityConfiguration.getCapabilityConfiguration.selector ^
          ICapabilityConfiguration.beforeCapabilityConfigSet.selector
      ),
      abi.encode(false)
    );
    s_capabilityWithConfigurationContract.configurationContract = contractWithoutERC165;
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_capabilityWithConfigurationContract;

    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilitiesRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        contractWithoutERC165
      )
    );
    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_AddCapability_NoConfigurationContract() public {
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_basicCapability;

    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId("data-streams-reports", "1.0.0");
    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.CapabilityConfigured(hashedCapabilityId);
    s_CapabilitiesRegistry.addCapabilities(capabilities);
    CapabilitiesRegistry.CapabilityInfo memory storedCapability = s_CapabilitiesRegistry.getCapability(
      hashedCapabilityId
    );

    assertEq(storedCapability.labelledName, s_basicCapability.labelledName);
    assertEq(storedCapability.version, s_basicCapability.version);
    assertEq(uint256(storedCapability.responseType), uint256(s_basicCapability.responseType));
    assertEq(storedCapability.configurationContract, s_basicCapability.configurationContract);
  }

  function test_AddCapability_WithConfiguration() public {
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](1);
    capabilities[0] = s_capabilityWithConfigurationContract;

    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_capabilityWithConfigurationContract.labelledName,
      s_capabilityWithConfigurationContract.version
    );
    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.CapabilityConfigured(hashedCapabilityId);
    s_CapabilitiesRegistry.addCapabilities(capabilities);

    CapabilitiesRegistry.CapabilityInfo memory storedCapability = s_CapabilitiesRegistry.getCapability(
      hashedCapabilityId
    );

    assertEq(storedCapability.labelledName, s_capabilityWithConfigurationContract.labelledName);
    assertEq(storedCapability.version, s_capabilityWithConfigurationContract.version);
    assertEq(uint256(storedCapability.responseType), uint256(s_capabilityWithConfigurationContract.responseType));
    assertEq(storedCapability.configurationContract, s_capabilityWithConfigurationContract.configurationContract);
  }
}
