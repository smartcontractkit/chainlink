// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPConfig} from "../../capability/CCIPConfig.sol";
import {CCIPConfigTypes} from "../../capability/libraries/CCIPConfigTypes.sol";
import {Internal} from "../../libraries/Internal.sol";

contract CCIPConfigHelper is CCIPConfig {
  constructor(address capabilitiesRegistry) CCIPConfig(capabilitiesRegistry) {}

  function stateFromConfigLength(uint256 configLength) public pure returns (CCIPConfigTypes.ConfigState) {
    return _stateFromConfigLength(configLength);
  }

  function validateConfigStateTransition(
    CCIPConfigTypes.ConfigState currentState,
    CCIPConfigTypes.ConfigState newState
  ) public pure {
    _validateConfigStateTransition(currentState, newState);
  }

  function validateConfigTransition(
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig,
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory newConfigWithMeta
  ) public pure {
    _validateConfigTransition(currentConfig, newConfigWithMeta);
  }

  function computeNewConfigWithMeta(
    uint32 donId,
    CCIPConfigTypes.OCR3ConfigWithMeta[] memory currentConfig,
    CCIPConfigTypes.OCR3Config[] memory newConfig,
    CCIPConfigTypes.ConfigState currentState,
    CCIPConfigTypes.ConfigState newState
  ) public view returns (CCIPConfigTypes.OCR3ConfigWithMeta[] memory) {
    return _computeNewConfigWithMeta(donId, currentConfig, newConfig, currentState, newState);
  }

  function groupByPluginType(
    CCIPConfigTypes.OCR3Config[] memory ocr3Configs
  )
    public
    pure
    returns (CCIPConfigTypes.OCR3Config[] memory commitConfigs, CCIPConfigTypes.OCR3Config[] memory execConfigs)
  {
    return _groupByPluginType(ocr3Configs);
  }

  function computeConfigDigest(
    uint32 donId,
    uint64 configCount,
    CCIPConfigTypes.OCR3Config memory ocr3Config
  ) public pure returns (bytes32) {
    return _computeConfigDigest(donId, configCount, ocr3Config);
  }

  function validateConfig(CCIPConfigTypes.OCR3Config memory cfg) public view {
    _validateConfig(cfg);
  }

  function updatePluginConfig(
    uint32 donId,
    Internal.OCRPluginType pluginType,
    CCIPConfigTypes.OCR3Config[] memory newConfig
  ) public {
    _updatePluginConfig(donId, pluginType, newConfig);
  }
}
