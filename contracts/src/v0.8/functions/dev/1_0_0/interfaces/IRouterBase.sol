// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Chainlink base Router interface.
 */
interface IRouterBase {
  /**
   * @notice Get the current contract given an ID
   * @param id A bytes32 identifier for the route
   * @return contract The current contract address
   */
  function getContractById(bytes32 id) external view returns (address);

  /**
   * @notice Get the proposed next contract given an ID
   * @param id A bytes32 identifier for the route
   * @return contract The current or proposed contract address
   */
  function getProposedContractById(bytes32 id) external view returns (address);

  /**
   * @notice Return the latest proprosal set
   * @return timelockEndBlock The block number that the proposal is able to be merged at
   * @return ids The identifiers of the contracts to update
   * @return to The addresses of the contracts that will be updated to
   */
  function getProposedContractSet() external view returns (uint, bytes32[] memory, address[] memory);

  /**
   * @notice Proposes one or more updates to the contract routes
   * @dev Only callable by owner
   */
  function proposeContractsUpdate(bytes32[] memory proposalSetIds, address[] memory proposalSetAddresses) external;

  /**
   * @notice Updates the current contract routes to the proposed contracts
   * @dev Only callable once timelock has passed
   * @dev Only callable by owner
   */
  function updateContracts() external;

  /**
   * @notice Proposes new configuration data for the Router contract itself
   * @dev Only callable by owner
   */
  function proposeConfigUpdateSelf(bytes calldata config) external;

  /**
   * @notice Updates configuration data for the Router contract itself
   * @dev Only callable once timelock has passed
   * @dev Only callable by owner
   */
  function updateConfigSelf() external;

  /**
   * @notice Proposes new configuration data for the current (not proposed) contract
   * @dev Only callable by owner
   */
  function proposeConfigUpdate(bytes32 id, bytes calldata config) external;

  /**
   * @notice Sends new configuration data to the contract along a route route
   * @dev Only callable once timelock has passed
   * @dev Only callable by owner
   */
  function updateConfig(bytes32 id) external;

  /**
   * @notice Propose a change to the amount of blocks of the timelock
   * (the amount of blocks that are required to pass before a change can be applied)
   * @dev Only callable by owner
   */
  function proposeTimelockBlocks(uint16 blocks) external;

  /**
   * @notice Apply a proposed change to the amount of blocks required for the timelock
   * (the amount of blocks that are required to pass before a change can be applied)
   * @dev Only callable after the timelock blocks proposal has gone through the timelock itself
   * @dev Only callable by owner
   */
  function updateTimelockBlocks() external;

  /**
   * @dev Puts the system into an emergency stopped state.
   * @dev Only callable by owner
   */
  function pause() external;

  /**
   * @dev Takes the system out of an emergency stopped state.
   * @dev Only callable by owner
   */
  function unpause() external;
}
