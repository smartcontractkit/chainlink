// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseConfiguratorTest.t.sol";
import {Configurator} from "../../Configurator.sol";

contract ConfiguratorSetProductionConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    bytes[] memory signers = _getSigners(MAX_ORACLES);

    vm.startPrank(USER);
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      OFFCHAIN_CONFIG_VERSION,
      bytes("")
    );
  }

  function test_revertsIfSetWithTooManySigners() public {
    bytes[] memory signers = new bytes[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(Configurator.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      OFFCHAIN_CONFIG_VERSION,
      bytes("")
    );
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(Configurator.FaultToleranceMustBePositive.selector));
    bytes[] memory signers = _getSigners(MAX_ORACLES);
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      s_offchaintransmitters,
      0,
      bytes(""),
      OFFCHAIN_CONFIG_VERSION,
      bytes("")
    );
  }

  function test_revertsIfNotEnoughSigners() public {
    bytes[] memory signers = _getSigners(2);

    vm.expectRevert(
      abi.encodeWithSelector(Configurator.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      OFFCHAIN_CONFIG_VERSION,
      bytes("")
    );
  }

  function test_revertsIfOnchainConfigIsInvalid() public {
    bytes[] memory signers = _getSigners(4);
    bytes32[] memory offchainTransmitters = _getOffchainTransmitters(4);
    bytes memory onchainConfig = bytes("");
    uint8 f = 1;
    bytes memory offchainConfig = abi.encodePacked(keccak256("offchainConfig"));

    vm.expectRevert(abi.encodeWithSelector(Configurator.InvalidOnchainLength.selector, onchainConfig.length));
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    onchainConfig = abi.encode(uint256(0), bytes32(0));

    vm.expectRevert(abi.encodeWithSelector(Configurator.UnsupportedOnchainConfigVersion.selector, uint256(0)));
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    onchainConfig = abi.encode(
      uint256(1),
      keccak256("any non-zero predecessor config digest is invalid for production")
    );

    vm.expectRevert(
      abi.encodeWithSelector(
        Configurator.NonZeroPredecessorConfigDigest.selector,
        keccak256("any non-zero predecessor config digest is invalid for production")
      )
    );
    s_configurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );
  }

  function test_supportsHigherVersionsIgnoringExcessOnchainConfig() public {
    bytes[] memory signers = _getSigners(4);
    bytes32[] memory offchainTransmitters = _getOffchainTransmitters(4);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint256(2), bytes32(0), keccak256("some rubbish"));
    bytes memory offchainConfig = abi.encodePacked(keccak256("offchainConfig"));

    // initial block number
    vm.roll(5);

    bytes32 cd1 = s_exposedConfigurator.exposedConfigDigestFromConfigData(
      CONFIG_ID_1,
      block.chainid,
      address(s_exposedConfigurator),
      1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    // when isGreenProduction=false

    vm.expectEmit();
    emit ProductionConfigSet(
      CONFIG_ID_1,
      0,
      cd1,
      1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig,
      false
    );

    s_exposedConfigurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );
  }

  function test_correctlyUpdatesTheConfig() public {
    bytes[] memory signers = _getSigners(4);
    bytes32[] memory offchainTransmitters = _getOffchainTransmitters(4);
    uint8 f = 1;
    bytes memory onchainConfig = abi.encodePacked(uint256(1), bytes32(0));
    bytes memory offchainConfig = abi.encodePacked(keccak256("offchainConfig"));

    // initial block number
    vm.roll(5);

    bytes32 cd1 = s_exposedConfigurator.exposedConfigDigestFromConfigData(
      CONFIG_ID_1,
      block.chainid,
      address(s_exposedConfigurator),
      1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    // when isGreenProduction=false

    vm.expectEmit();
    emit ProductionConfigSet(
      CONFIG_ID_1,
      0,
      cd1,
      1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig,
      false
    );

    s_exposedConfigurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    Configurator.ConfigurationState memory configurationState = s_exposedConfigurator.exposedReadConfigurationStates(
      CONFIG_ID_1
    );
    assertEq(configurationState.configDigest[0], cd1);
    assertEq(configurationState.configDigest[1], 0); // no staging config yet
    assertEq(configurationState.configCount, 1);
    assertEq(configurationState.isGreenProduction, false);
    assertEq(configurationState.latestConfigBlockNumber, block.number);

    // go to new block
    vm.roll(10);

    // set it again, configCount=2

    bytes32 cd2 = s_exposedConfigurator.exposedConfigDigestFromConfigData(
      CONFIG_ID_1,
      block.chainid,
      address(s_exposedConfigurator),
      2,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    vm.expectEmit();
    emit ProductionConfigSet(
      CONFIG_ID_1,
      5,
      cd2,
      2,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig,
      false
    );

    s_exposedConfigurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    configurationState = s_exposedConfigurator.exposedReadConfigurationStates(CONFIG_ID_1);
    assertEq(configurationState.configDigest[0], cd2);
    assertEq(configurationState.configDigest[1], 0); // no staging config yet
    assertEq(configurationState.configCount, 2);
    assertEq(configurationState.isGreenProduction, false);
    assertEq(configurationState.latestConfigBlockNumber, block.number);

    // when isGreenProduction=true
    s_exposedConfigurator.exposedSetIsGreenProduction(CONFIG_ID_1, true);

    // go to new block
    vm.roll(15);

    // set it again, configCount=3
    bytes32 cd3 = s_exposedConfigurator.exposedConfigDigestFromConfigData(
      CONFIG_ID_1,
      block.chainid,
      address(s_exposedConfigurator),
      3,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    vm.expectEmit();
    emit ProductionConfigSet(
      CONFIG_ID_1,
      10,
      cd3,
      3,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig,
      true
    );

    s_exposedConfigurator.setProductionConfig(
      CONFIG_ID_1,
      signers,
      offchainTransmitters,
      f,
      onchainConfig,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfig
    );

    configurationState = s_exposedConfigurator.exposedReadConfigurationStates(CONFIG_ID_1);
    assertEq(configurationState.configDigest[0], cd2); // the previous config left unchanged
    assertEq(configurationState.configDigest[1], cd3); // new config is on green now because green is production due to isGreenProduction=true
    assertEq(configurationState.configCount, 3);
    assertEq(configurationState.isGreenProduction, true);
    assertEq(configurationState.latestConfigBlockNumber, block.number);
  }
}
