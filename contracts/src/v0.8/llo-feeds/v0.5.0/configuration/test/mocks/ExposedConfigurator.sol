// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Configurator} from "../../Configurator.sol";

// Exposed ChannelConfigStore exposes certain internal ChannelConfigStore
// methods/structures so that golang code can access them, and we get
// reliable type checking on their usage
contract ExposedConfigurator is Configurator {
  constructor() {}

  function exposedReadConfigurationStates(bytes32 donId) public view returns (ConfigurationState memory) {
    return s_configurationStates[donId];
  }
}
