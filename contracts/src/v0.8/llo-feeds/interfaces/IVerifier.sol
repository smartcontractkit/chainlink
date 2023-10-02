// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC165.sol";
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
   * @notice sets offchain reporting protocol configuration incl. participating oracles
   * @param feedId Feed ID to set config for
   * @param signers addresses with which oracles sign the reports
   * @param offchainTransmitters CSA key for the ith Oracle
   * @param f number of faulty oracles the system can tolerate
   * @param onchainConfig serialized configuration used by the contract (and possibly oracles)
   * @param offchainConfigVersion version number for offchainEncoding schema
   * @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfig(
    bytes32 feedId,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
   * @notice identical to `setConfig` except with args for sourceChainId and sourceAddress
   * @param feedId Feed ID to set config for
   * @param sourceChainId Chain ID of source config
   * @param sourceAddress Address of source config Verifier
   * @param newConfigCount Param to force the new config count
   * @param signers addresses with which oracles sign the reports
   * @param offchainTransmitters CSA key for the ith Oracle
   * @param f number of faulty oracles the system can tolerate
   * @param onchainConfig serialized configuration used by the contract (and possibly oracles)
   * @param offchainConfigVersion version number for offchainEncoding schema
   * @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfigFromSource(
    bytes32 feedId,
    uint256 sourceChainId,
    address sourceAddress,
    uint32 newConfigCount,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
   * @notice Activates the configuration for a config digest
   * @param feedId Feed ID to activate config for
   * @param configDigest The config digest to activate
   * @dev This function can be called by the contract admin to activate a configuration.
   */
  function activateConfig(bytes32 feedId, bytes32 configDigest) external;

  /**
   * @notice Deactivates the configuration for a config digest
   * @param feedId Feed ID to deactivate config for
   * @param configDigest The config digest to deactivate
   * @dev This function can be called by the contract admin to deactivate an incorrect configuration.
   */
  function deactivateConfig(bytes32 feedId, bytes32 configDigest) external;

  /**
   * @notice Activates the given feed
   * @param feedId Feed ID to activated
   * @dev This function can be called by the contract admin to activate a feed
   */
  function activateFeed(bytes32 feedId) external;

  /**
   * @notice Deactivates the given feed
   * @param feedId Feed ID to deactivated
   * @dev This function can be called by the contract admin to deactivate a feed
   */
  function deactivateFeed(bytes32 feedId) external;

  /**
   * @notice returns the latest config digest and epoch for a feed
   * @param feedId Feed ID to fetch data for
   * @return scanLogs indicates whether to rely on the configDigest and epoch
   * returned or whether to scan logs for the Transmitted event instead.
   * @return configDigest
   * @return epoch
   */
  function latestConfigDigestAndEpoch(
    bytes32 feedId
  ) external view returns (bool scanLogs, bytes32 configDigest, uint32 epoch);

  /**
   * @notice information about current offchain reporting protocol configuration
   * @param feedId Feed ID to fetch data for
   * @return configCount ordinal number of current config, out of all configs applied to this contract so far
   * @return blockNumber block at which this config was set
   * @return configDigest domain-separation tag for current config
   */
  function latestConfigDetails(
    bytes32 feedId
  ) external view returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest);
}
