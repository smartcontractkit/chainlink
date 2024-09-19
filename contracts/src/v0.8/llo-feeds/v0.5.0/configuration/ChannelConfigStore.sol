// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.19;

import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {IChannelConfigStore} from "./interfaces/IChannelConfigStore.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";

contract ChannelConfigStore is ConfirmedOwner, IChannelConfigStore, TypeAndVersionInterface {
  event NewChannelDefinition(uint256 indexed donId, uint32 version, string url, bytes32 sha);

  constructor() ConfirmedOwner(msg.sender) {}

  /// @notice The version of a channel definition keyed by DON ID
  // Increments by 1 on every update
  mapping(uint256 => uint256) internal s_channelDefinitionVersions;

  function setChannelDefinitions(uint32 donId, string calldata url, bytes32 sha) external onlyOwner {
    uint32 newVersion = uint32(++s_channelDefinitionVersions[uint256(donId)]);
    emit NewChannelDefinition(donId, newVersion, url, sha);
  }

  function typeAndVersion() external pure override returns (string memory) {
    return "ChannelConfigStore 0.0.1";
  }

  function supportsInterface(bytes4 interfaceId) external pure returns (bool) {
    return interfaceId == type(IChannelConfigStore).interfaceId;
  }
}
