// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseConfiguratorTest.t.sol";
import {Configurator} from "../../Configurator.sol";

contract ConfiguratorPromoteStagingConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.startPrank(USER);

    vm.expectRevert("Only callable by owner");

    s_configurator.promoteStagingConfig(CONFIG_ID_1, false);
  }

  function test_revertsIfIsGreenProductionDoesNotMatchContractState() public {
    vm.expectRevert(
      abi.encodeWithSelector(Configurator.IsGreenProductionMustMatchContractState.selector, CONFIG_ID_1, false)
    );
    s_configurator.promoteStagingConfig(CONFIG_ID_1, true);
  }

  function test_revertsIfNoConfigHasEverBeenSetWithThisConfigId() public {
    vm.expectRevert(abi.encodeWithSelector(Configurator.ConfigUnset.selector, keccak256("nonExistentConfigId")));
    s_configurator.promoteStagingConfig(keccak256("nonExistentConfigId"), false);
  }

  function test_revertsIfStagingConfigDigestIsZero() public {
    // isGreenProduction = false
    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(1, uint32(block.number), false, [bytes32(0), bytes32(0)])
    );

    vm.expectRevert(abi.encodeWithSelector(Configurator.ConfigUnsetStaging.selector, CONFIG_ID_1, false));
    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, false);

    // isGreenProduction = true
    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(1, uint32(block.number), true, [bytes32(0), bytes32(0)])
    );

    vm.expectRevert(abi.encodeWithSelector(Configurator.ConfigUnsetStaging.selector, CONFIG_ID_1, true));
    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, true);
  }

  function test_revertsIfProductionConfigDigestIsZero() public {
    // isGreenProduction = false
    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(1, uint32(block.number), false, [bytes32(0), keccak256("stagingConfigDigest")])
    );

    vm.expectRevert(abi.encodeWithSelector(Configurator.ConfigUnsetProduction.selector, CONFIG_ID_1, false));
    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, false);

    // isGreenProduction = true

    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(1, uint32(block.number), true, [keccak256("stagingConfigDigest"), bytes32(0)])
    );

    vm.expectRevert(abi.encodeWithSelector(Configurator.ConfigUnsetProduction.selector, CONFIG_ID_1, true));
    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, true);
  }

  function test_promotesStagingConfig() public {
    // isGreenProduction = false
    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(
        1,
        uint32(block.number),
        false,
        [keccak256("productionConfigDigest"), keccak256("stagingConfigDigest")]
      )
    );

    vm.expectEmit();
    emit PromoteStagingConfig(CONFIG_ID_1, keccak256("productionConfigDigest"), true);

    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, false);
    assertEq(s_exposedConfigurator.exposedReadConfigurationStates(CONFIG_ID_1).isGreenProduction, true);

    // isGreenProduction = true

    s_exposedConfigurator.exposedSetConfigurationState(
      CONFIG_ID_1,
      Configurator.ConfigurationState(
        1,
        uint32(block.number),
        true,
        [keccak256("stagingConfigDigest"), keccak256("productionConfigDigest")]
      )
    );

    vm.expectEmit();
    emit PromoteStagingConfig(CONFIG_ID_1, keccak256("productionConfigDigest"), false);

    s_exposedConfigurator.promoteStagingConfig(CONFIG_ID_1, true);
    assertEq(s_exposedConfigurator.exposedReadConfigurationStates(CONFIG_ID_1).isGreenProduction, false);
  }
}
