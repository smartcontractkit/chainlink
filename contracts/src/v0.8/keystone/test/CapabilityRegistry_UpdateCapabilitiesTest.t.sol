// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

contract CapabilityRegistry_UpdateCapabilitiesTest is BaseTest {
  event CapabilityConfigured(bytes32 indexed hashedCapabilityId);

  CapabilityRegistry.Capability[] internal s_updatedCapabilities;
  CapabilityConfigurationContract internal s_newCapabilityConfig;

  CapabilityRegistry.CapabilityResponseType constant NEW_BASIC_CAPABILITY_RESPONSE_TYPE =
    CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL;
  CapabilityRegistry.CapabilityResponseType constant NEW_CAPABILITY_WITH_CONFIG_CONTRACT_RESPONSE_TYPE =
    CapabilityRegistry.CapabilityResponseType.REPORT;

  function setUp() public override {
    BaseTest.setUp();

    changePrank(ADMIN);
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    s_newCapabilityConfig = new CapabilityConfigurationContract();

    s_updatedCapabilities.push(
      CapabilityRegistry.Capability({
        labelledName: "data-streams-reports",
        version: "1.0.0",
        responseType: NEW_BASIC_CAPABILITY_RESPONSE_TYPE,
        configurationContract: address(s_newCapabilityConfig)
      })
    );
    s_updatedCapabilities.push(
      CapabilityRegistry.Capability({
        labelledName: "read-ethereum-mainnet-gas-price",
        version: "1.0.2",
        responseType: NEW_CAPABILITY_WITH_CONFIG_CONTRACT_RESPONSE_TYPE,
        configurationContract: address(s_newCapabilityConfig)
      })
    );
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);
  }

  function test_RevertWhen_CapabilityDoesNotExists() public {
    bytes32 versionNum = "1.0.3";
    s_updatedCapabilities[0].version = versionNum;
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_updatedCapabilities[0].labelledName,
      versionNum
    );
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityDoesNotExist.selector, hashedCapabilityId));
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);
  }

  function test_RevertWhen_ConfigurationContractNotDeployed() public {
    address nonExistentContract = address(1);
    s_updatedCapabilities[1].configurationContract = nonExistentContract;

    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilityRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        nonExistentContract
      )
    );
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);
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
    s_updatedCapabilities[1].configurationContract = contractWithoutERC165;
    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilityRegistry.InvalidCapabilityConfigurationContractInterface.selector,
        contractWithoutERC165
      )
    );
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);
  }

  function test_UpdateCapabilities_NoConfigurationContract() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      bytes32("data-streams-reports"),
      bytes32("1.0.0")
    );
    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit CapabilityConfigured(hashedCapabilityId);
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);

    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(hashedCapabilityId);

    assertEq(storedCapability.labelledName, s_basicCapability.labelledName);
    assertEq(storedCapability.version, s_basicCapability.version);
    assertEq(uint256(storedCapability.responseType), uint256(NEW_BASIC_CAPABILITY_RESPONSE_TYPE));
    assertEq(storedCapability.configurationContract, address(s_newCapabilityConfig));
  }

  function test_UpdateCapabilities_WithConfiguration() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      bytes32(s_capabilityWithConfigurationContract.labelledName),
      bytes32(s_capabilityWithConfigurationContract.version)
    );
    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit CapabilityConfigured(hashedCapabilityId);
    s_capabilityRegistry.updateCapabilities(s_updatedCapabilities);

    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(hashedCapabilityId);

    assertEq(storedCapability.labelledName, s_capabilityWithConfigurationContract.labelledName);
    assertEq(storedCapability.version, s_capabilityWithConfigurationContract.version);
    assertEq(uint256(storedCapability.responseType), uint256(NEW_CAPABILITY_WITH_CONFIG_CONTRACT_RESPONSE_TYPE));
    assertEq(storedCapability.configurationContract, address(s_newCapabilityConfig));
  }
}
