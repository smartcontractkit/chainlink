// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ICapabilityConfiguration} from "../../interfaces/ICapabilityConfiguration.sol";
import {ERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/ERC165.sol";

contract CapabilityConfigurationContract is ICapabilityConfiguration, ERC165 {
  mapping(uint256 => bytes) private s_donConfiguration;

  function getCapabilityConfiguration(uint32 donId) external view returns (bytes memory configuration) {
    return s_donConfiguration[donId];
  }

  function beforeCapabilityConfigSet(bytes32[] calldata, bytes calldata config, uint64, uint32 donId) external {
    // In practice a real capability configuration contract will do more than just store
    // the config for a DON
    s_donConfiguration[donId] = config;
  }

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == this.getCapabilityConfiguration.selector ^ this.beforeCapabilityConfigSet.selector;
  }
}
