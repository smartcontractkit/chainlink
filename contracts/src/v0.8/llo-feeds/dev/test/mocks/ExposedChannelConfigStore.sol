// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ChannelConfigStore} from "../../ChannelConfigStore.sol";

// Exposed ChannelConfigStore exposes certain internal ChannelConfigStore
// methods/structures so that golang code can access them, and we get
// reliable type checking on their usage
contract ExposedChannelConfigStore is ChannelConfigStore {
  constructor() {}

  function exposedReadChannelDefinitionStates(uint256 donId) public view returns (uint256) {
    return s_channelDefinitionVersions[donId];
  }
}
