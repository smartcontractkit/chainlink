// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {IERC165} from "@openzeppelin/contracts/interfaces/IERC165.sol";

interface IVerifier is IERC165 {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier.
   * @param signedReport The encoded data to be verified.
   * @param requester The original address that requested to verify the contract.
   * This is only used for logging purposes.
   * @dev Verification is typically only done through the proxy contract so
   * we can't just use msg.sender to log the requester as the msg.sender
   * contract will always be the proxy.
   * @return response The encoded verified response.
   */
  function verify(bytes memory signedReport, address requester) external returns (bytes memory response);

  /**
   * @notice sets offchain reporting protocol configuration incl. participating oracles
   * @param feedId Feed ID to set config for
   * @param signers addresses with which oracles sign the reports
   * @param offchainTransmitters CSA key for the ith Oracle
   * @param f number of faulty oracles the system can tolerate
   * @param onchainConfig serialized configuration used by the contract (and possibly oracles)
   * @param offchainConfigVersion version number for offchainEncoding schema
   * @param offchainConfig serialized configuration used by the oracles exclusively and only passed through the contract
   */
  function setConfig(
    bytes32 feedId,
    address[] memory signers,
    bytes32[] memory offchainTransmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external;

  /**
   * @notice returns the latest config digest and epoch for a feed
   * @param feedId Feed ID to fetch data for
   * @return scanLogs indicates whether to rely on the configDigest and epoch
   * returned or whether to scan logs for the Transmitted event instead.
   * @return configDigest
   * @return epoch
   */
  function latestConfigDigestAndEpoch(bytes32 feedId)
    external
    view
    returns (
      bool scanLogs,
      bytes32 configDigest,
      uint32 epoch
    );

  /**
   * @notice information about current offchain reporting protocol configuration
   * @param feedId Feed ID to fetch data for
   * @return configCount ordinal number of current config, out of all configs applied to this contract so far
   * @return blockNumber block at which this config was set
   * @return configDigest domain-separation tag for current config
   */
  function latestConfigDetails(bytes32 feedId)
    external
    view
    returns (
      uint32 configCount,
      uint32 blockNumber,
      bytes32 configDigest
    );
}
