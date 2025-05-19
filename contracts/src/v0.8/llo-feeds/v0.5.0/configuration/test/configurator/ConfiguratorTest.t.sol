// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseConfiguratorTest.t.sol";
import {Configurator} from "../../Configurator.sol";
import {IConfigurator} from "../../interfaces/IConfigurator.sol";

contract ConfiguratorTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function testTypeAndVersion() public view {
    assertEq(s_configurator.typeAndVersion(), "Configurator 0.5.0");
  }

  function testSupportsInterface() public view {
    assertTrue(s_configurator.supportsInterface(type(IConfigurator).interfaceId));
  }
}
