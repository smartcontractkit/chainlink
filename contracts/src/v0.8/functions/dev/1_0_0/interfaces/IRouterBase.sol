// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title Chainlink base Router interface.
 */
interface IRouterBase {
  /**
   * @notice Returns the latest semantic version of the system
   * @dev See https://semver.org/ for more details
   * @return major The current major version number
   * @return minor The current minor version number
   * @return patch The current patch version number
   */
  function version() external view returns (uint16 major, uint16 minor, uint16 patch);

  /**
   * @notice Get the current contract given an ID
   * @param id A bytes32 identifier for the route
   * @return contract The current contract address
   */
  function getRoute(bytes32 id) external view returns (address);

  /**
   * @notice Get the proposed next contract given an ID
   * @param id A bytes32 identifier for the route
   * @return contract The current or proposed contract address
   */
  function getRoute(bytes32 id, bool useProposed) external view returns (address);

  /**
   * @notice Return the latest proprosal set
   * @return proposedAtBlock The block number that the proposal was created at
   * @return ids The identifiers of the contracts to update
   * @return from The addresses of the contracts that will be updated from
   * @return to The addresses of the contracts that will be updated to
   */
  function getProposalSet() external view returns (uint, bytes32[] memory, address[] memory, address[] memory);

  /**
   * @notice Proposes one or more updates to the contract routes
   */
  function propose(
    bytes32[] memory proposalSetIds,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) external;

  /**
   * @notice Tests a proposal for the ability to make a successful upgrade
   */
  function validateProposal(bytes32 id, bytes calldata data) external;

  /**
   * @notice Updates the current contract routes to the proposed contracts once timelock has passed
   */
  function upgrade() external;

  /**
   * @notice Proposes new configuration data that will be given to the contract route
   */
  function proposeConfig(bytes32 id, bytes calldata config) external;

  /**
   * @notice Sends new configuration data to the contract route once timelock has passed
   */
  function updateConfig(bytes32 id) external;

  /**
   * @dev Propose a change to the amount of blocks required for a timelock
   */
  function proposeTimelockBlocks(uint16 blocks) external;

  /**
   * @dev Change the amount of blocks required for the timelock
   * (only after the proposal has gone through the timelock itself)
   */
  function updateTimelockBlocks() external;

  /**
   * @dev Returns true if the contract is paused, and false otherwise.
   */
  function isPaused() external view returns (bool);

  /**
   * @dev Triggers stopped state.
   *
   * Requirements:
   *
   * - The contract must not be paused.
   */
  function pause() external;

  /**
   * @dev Returns to normal state.
   *
   * Requirements:
   *
   * - The contract must be paused.
   */
  function unpause() external;
}
