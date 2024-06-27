// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {IVerifier} from "../../interfaces/IVerifier.sol";
import {Common} from "../../../libraries/Common.sol";

contract ErroredVerifier is IVerifier {
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == this.verify.selector;
  }

  function verify(
    bytes memory,
    /**
     * signedReport*
     */
    address
  )
    external
    pure
    override
    returns (
      /**
       * sender*
       */
      bytes memory
    )
  {
    revert("Failed to verify");
  }

  function setConfig(
    bytes32,
    address[] memory,
    bytes32[] memory,
    uint8,
    bytes memory,
    uint64,
    bytes memory,
    Common.AddressAndWeight[] memory
  ) external pure override {
    revert("Failed to set config");
  }

  function setConfigFromSource(
    bytes32,
    uint256,
    address,
    uint32,
    address[] memory,
    bytes32[] memory,
    uint8,
    bytes memory,
    uint64,
    bytes memory,
    Common.AddressAndWeight[] memory
  ) external pure override {
    revert("Failed to set config");
  }

  function activateConfig(bytes32, bytes32) external pure {
    revert("Failed to activate config");
  }

  function deactivateConfig(bytes32, bytes32) external pure {
    revert("Failed to deactivate config");
  }

  function activateFeed(bytes32) external pure {
    revert("Failed to activate feed");
  }

  function deactivateFeed(bytes32) external pure {
    revert("Failed to deactivate feed");
  }

  function latestConfigDigestAndEpoch(bytes32) external pure override returns (bool, bytes32, uint32) {
    revert("Failed to get latest config digest and epoch");
  }

  function latestConfigDetails(bytes32) external pure override returns (uint32, uint32, bytes32) {
    revert("Failed to get latest config details");
  }
}
