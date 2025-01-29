// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";

interface IVerifier is IERC165 {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier.
   * @param signedReport The encoded data to be verified.
   * @param sender The address that requested to verify the contract.
   * This is only used for logging purposes.
   * @dev Verification is typically only done through the proxy contract so
   * we can't just use msg.sender to log the requester as the msg.sender
   * contract will always be the proxy.
   * @return verifierResponse The encoded verified response.
   */
  function verify(bytes calldata signedReport, address sender) external returns (bytes memory verifierResponse);

  /**
   * @notice sets a configuration and its associated keys and f
   * @param configDigest The digest of the configuration we're setting
   * @param signers addresses with which oracles sign the reports
   * @param f number of faulty oracles the system can tolerate
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfig(
    bytes32 configDigest,
    address[] calldata signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
   * @notice updates a configuration that has been set
   * @param configDigest The digest of the configuration we're updating
   * @param prevSigners the existing signers that need to be removed
   * @param newSigners the signers to be added
   * @param f the newnumber of faulty oracles the system can tolerate
   */
  function updateConfig(
    bytes32 configDigest,
    address[] calldata prevSigners,
    address[] calldata newSigners,
    uint8 f
  ) external;

  /**
   * @notice Activates the configuration for a config digest
   * @param configDigest The config digest to activate
   * @dev This function can be called by the contract admin to activate a configuration.
   */
  function activateConfig(bytes32 configDigest) external;

  /**
   * @notice Deactivates the configuration for a config digest
   * @param configDigest The config digest to deactivate
   * @dev This function can be called by the contract admin to deactivate an incorrect configuration.
   */
  function deactivateConfig(bytes32 configDigest) external;

  /**
   * @notice information about current offchain reporting protocol configuration
   * @param configDigest Config Digest to fetch data for
   * @return blockNumber block at which this config was set
   */
  function latestConfigDetails(bytes32 configDigest) external view returns (uint32 blockNumber);
}
