// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPCapabilityConfiguration} from "../../capability/CCIPCapabilityConfiguration.sol";

contract CCIPCapabilityConfigurationHelper is CCIPCapabilityConfiguration {
  constructor(address capabilityRegistry) CCIPCapabilityConfiguration(capabilityRegistry) {}

  function stateFromConfigLength(uint256 configLength) public pure returns (ConfigState) {
    return _stateFromConfigLength(configLength);
  }

  function validateConfigStateTransition(ConfigState currentState, ConfigState newState) public pure {
    _validateConfigStateTransition(currentState, newState);
  }

  function validateConfigTransition(
    OCR3ConfigWithMeta[] memory currentConfig,
    OCR3ConfigWithMeta[] memory newConfigWithMeta
  ) public pure {
    _validateConfigTransition(currentConfig, newConfigWithMeta);
  }

  function computeNewConfigWithMeta(
    uint32 donId,
    OCR3ConfigWithMeta[] memory currentConfig,
    OCR3Config[] memory newConfig,
    ConfigState currentState,
    ConfigState newState
  ) public view returns (OCR3ConfigWithMeta[] memory) {
    return _computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, newState);
  }

  function groupByPluginType(OCR3Config[] memory ocr3Configs)
    public
    pure
    returns (OCR3Config[] memory commitConfigs, OCR3Config[] memory execConfigs)
  {
    return _groupByPluginType(ocr3Configs);
  }

  function computeConfigDigest(
    uint32 donId,
    uint64 configCount,
    OCR3Config memory ocr3Config
  ) public pure returns (bytes32) {
    return _computeConfigDigest(donId, configCount, ocr3Config);
  }

  function validateConfig(OCR3Config memory cfg) public view {
    _validateConfig(cfg);
  }

  function updatePluginConfig(uint32 donId, PluginType pluginType, OCR3Config[] memory newConfig) public {
    _updatePluginConfig(donId, pluginType, newConfig);
  }
}
