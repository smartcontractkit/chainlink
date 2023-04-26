// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {IVerifier} from "../../../../../src/v0.8/interfaces/IVerifier.sol";

contract ErroredVerifier is IVerifier {
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == this.verify.selector;
  }

  function verify(
    bytes memory, /**signedReport**/
    address /**sender**/
  ) external pure override returns (bytes memory) {
    revert("Failed to verify");
  }

  function setConfig(
    bytes32,
    address[] memory,
    bytes32[] memory,
    uint8,
    bytes memory,
    uint64,
    bytes memory
  ) external pure override {
    revert("Failed to set config");
  }

  function latestConfigDigestAndEpoch(bytes32)
    external
    pure
    override
    returns (
      bool,
      bytes32,
      uint32
    )
  {
    revert("Failed to get latest config digest and epoch");
  }

  function latestConfigDetails(bytes32)
    external
    pure
    override
    returns (
      uint32,
      uint32,
      bytes32
    )
  {
    revert("Failed to get latest config details");
  }
}
