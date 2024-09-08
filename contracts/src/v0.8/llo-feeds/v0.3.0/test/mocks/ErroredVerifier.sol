// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {IVerifier} from "../../interfaces/IVerifier.sol";
import {Common} from "../../../libraries/Common.sol";

contract ErroredVerifier is IVerifier {
  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == this.verify.selector;
  }

  //define each of the errors thrown in the revert below

  error FailedToVerify();
  error FailedToSetConfig();
  error FailedToActivateConfig();
  error FailedToDeactivateConfig();
  error FailedToActivateFeed();
  error FailedToDeactivateFeed();
  error FailedToGetLatestConfigDigestAndEpoch();
  error FailedToGetLatestConfigDetails();

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
    revert FailedToVerify();
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
    revert FailedToSetConfig();
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
    revert FailedToSetConfig();
  }

  function activateConfig(bytes32, bytes32) external pure {
    revert FailedToActivateConfig();
  }

  function deactivateConfig(bytes32, bytes32) external pure {
    revert FailedToDeactivateConfig();
  }

  function activateFeed(bytes32) external pure {
    revert FailedToActivateFeed();
  }

  function deactivateFeed(bytes32) external pure {
    revert FailedToDeactivateFeed();
  }

  function latestConfigDigestAndEpoch(bytes32) external pure override returns (bool, bytes32, uint32) {
    revert FailedToGetLatestConfigDigestAndEpoch();
  }

  function latestConfigDetails(bytes32) external pure override returns (uint32, uint32, bytes32) {
    revert FailedToGetLatestConfigDetails();
  }
}
