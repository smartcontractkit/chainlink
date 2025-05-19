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
  error FailedToUnsetConfig();
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

  function updateConfig(bytes32, address[] calldata, address[] calldata, uint8) external pure {
    revert FailedToUnsetConfig();
  }

  function setConfig(bytes32, address[] calldata, uint8, Common.AddressAndWeight[] calldata) external pure override {
    revert FailedToSetConfig();
  }

  function activateConfig(bytes32) external pure {
    revert FailedToActivateConfig();
  }

  function deactivateConfig(bytes32) external pure {
    revert FailedToDeactivateConfig();
  }

  function latestConfigDetails(bytes32) external pure override returns (uint32) {
    revert FailedToGetLatestConfigDetails();
  }
}
