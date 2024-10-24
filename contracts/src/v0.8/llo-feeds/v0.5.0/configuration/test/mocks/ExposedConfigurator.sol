// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Configurator} from "../../Configurator.sol";

// Exposed ChannelConfigStore exposes certain internal ChannelConfigStore
// methods/structures so that golang code can access them, and we get
// reliable type checking on their usage
contract ExposedConfigurator is Configurator {
  constructor() {}

  function exposedReadConfigurationStates(bytes32 configId) public view returns (ConfigurationState memory) {
    return s_configurationStates[configId];
  }

  function exposedSetIsGreenProduction(bytes32 configId, bool isGreenProduction) public {
    s_configurationStates[configId].isGreenProduction = isGreenProduction;
  }

  function exposedSetConfigurationState(bytes32 configId, ConfigurationState memory state) public {
    s_configurationStates[configId] = state;
  }

  function exposedConfigDigestFromConfigData(
    bytes32 _configId,
    uint256 _chainId,
    address _contractAddress,
    uint64 _configCount,
    bytes[] memory _signers,
    bytes32[] memory _offchainTransmitters,
    uint8 _f,
    bytes calldata _onchainConfig,
    uint64 _encodedConfigVersion,
    bytes memory _encodedConfig
  ) public pure returns (bytes32) {
    return
      _configDigestFromConfigData(
        _configId,
        _chainId,
        _contractAddress,
        _configCount,
        _signers,
        _offchainTransmitters,
        _f,
        _onchainConfig,
        _encodedConfigVersion,
        _encodedConfig
      );
  }
}
