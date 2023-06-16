// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IRouterBase} from "./interfaces/IRouterBase.sol";
import {ConfirmedOwnerWithProposal} from "../../../ConfirmedOwnerWithProposal.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {Pausable} from "../../../shared/vendor/openzeppelin-solidity/v.4.8.0/contracts/security/Pausable.sol";
import {IConfigurable} from "./interfaces/IConfigurable.sol";

abstract contract RouterBase is IRouterBase, Pausable, ITypeAndVersion, ConfirmedOwnerWithProposal {
  // ================================================================
  // |                         Version state                        |
  // ================================================================
  uint16 internal constant s_majorVersion = 1;
  uint16 internal s_minorVersion = 0;
  uint16 internal s_patchVersion = 0;

  // ================================================================
  // |                          Route state                         |
  // ================================================================
  mapping(bytes32 => address) internal s_route; /* jobId => contract address */
  error RouteNotFound(bytes32 jobId);

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================
  uint8 internal constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ProposalSet {
    bytes32[] jobIds;
    address[] from;
    address[] to;
    uint proposedAtBlock;
  }
  ProposalSet internal s_proposalSet;

  event Proposed(
    bytes32[] proposalSetJobIds,
    address[] proposalSetFromAddresses,
    address[] proposalSetToAddresses,
    uint block
  );
  event Upgraded(
    bytes32[] proposalSetJobIds,
    address[] proposalSetFromAddresses,
    address[] proposalSetToAddresses,
    uint block,
    uint16 major,
    uint16 minor,
    uint16 patch
  );

  struct ConfigProposal {
    bytes32 from;
    bytes to;
    uint proposedAtBlock;
  }
  mapping(bytes32 => ConfigProposal) internal s_proposedConfig; /* jobId => ConfigProposal */
  event ConfigProposed(bytes32 jobId, bytes32 from, bytes32 to);
  event ConfigUpdated(bytes32 jobId, bytes32 from, bytes32 to);
  error InvalidProposal();
  error ReservedLabel(bytes32 jobId);

  // ================================================================
  // |                          Timelock state                      |
  // ================================================================
  uint16 internal MAXIMUM_TIMELOCK_BLOCKS;
  uint16 internal s_timelockBlocks;
  struct TimeLockProposal {
    uint16 from;
    uint16 to;
    uint proposedAtBlock;
  }
  TimeLockProposal s_timelockProposal;
  event TimeLockProposed(uint16 from, uint16 to);
  event TimeLockUpdated(uint16 from, uint16 to);
  error ProposedTimelockAboveMaximum();
  error TimelockInEffect();

  // ================================================================
  // |                          Config state                        |
  // ================================================================
  bytes32 internal s_config_hash;

  // ================================================================
  // |                       Initialization                         |
  // ================================================================
  constructor(
    address newOwner,
    uint16 timelockBlocks,
    uint16 maximumTimelockBlocks,
    bytes32[] memory initialJobIds,
    address[] memory initialAddresses,
    bytes memory config
  ) ConfirmedOwnerWithProposal(newOwner, address(0)) Pausable() {
    // Set initial timelock
    s_timelockBlocks = timelockBlocks;
    MAXIMUM_TIMELOCK_BLOCKS = maximumTimelockBlocks;
    // Use a hash of the Router's address to self-identify, since it does not have a jobId
    bytes32 routerLabel = keccak256(abi.encodePacked(address(this)));
    // Set the initial config
    s_route[routerLabel] = address(this);
    _setConfig(config);
    s_config_hash = keccak256(config);
    // Fill initial routes, from empty addresses to current implementation contracts
    address[] memory emptyAddresses = new address[](initialJobIds.length);
    _validateProposal(initialJobIds, emptyAddresses, initialAddresses);
    s_proposalSet = ProposalSet(initialJobIds, emptyAddresses, initialAddresses, block.number);
    for (uint8 i = 0; i < s_proposalSet.jobIds.length; i++) {
      s_route[s_proposalSet.jobIds[i]] = s_proposalSet.to[i];
    }
  }

  // ================================================================
  // |                       Version methods                        |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function version()
    external
    view
    returns (
      uint16,
      uint16,
      uint16
    )
  {
    return (s_majorVersion, s_minorVersion, s_patchVersion);
  }

  // ================================================================
  // |                        Route methods                         |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(bytes32 jobId) public view returns (address) {
    return _getLatestRoute(jobId);
  }

  /**
   * @dev Helper function to get a contract from the current routes
   */
  function _getLatestRoute(bytes32 jobId) internal view returns (address) {
    address currentImplementation = s_route[jobId];
    if (currentImplementation == address(0)) {
      revert RouteNotFound(jobId);
    }
    return currentImplementation;
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(bytes32 jobId, bool useProposed) public view returns (address) {
    if (useProposed == true) {
      for (uint8 i = 0; i < s_proposalSet.jobIds.length; i++) {
        if (jobId == s_proposalSet.jobIds[i]) {
          // NOTE: proposals can be used immediately
          // if (block.number < s_proposalSet.proposedAtBlock + s_timelockBlocks) {
          //   revert TimelockInEffect();
          // }
          return s_proposalSet.to[i];
        }
      }
    }
    return _getLatestRoute(jobId);
  }

  // ================================================================
  // |                 Contract Proposal methods                    |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function getProposalSet()
    external
    view
    returns (
      uint,
      bytes32[] memory,
      address[] memory,
      address[] memory
    )
  {
    return (s_proposalSet.proposedAtBlock, s_proposalSet.jobIds, s_proposalSet.from, s_proposalSet.to);
  }

  /**
   * @dev Helper function to validate a proposal set
   */
  function _validateProposal(
    bytes32[] memory proposalSetJobIds,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) internal view {
    // All arrays must be of equal length
    if (
      proposalSetJobIds.length != proposalSetFromAddresses.length ||
      proposalSetJobIds.length != proposalSetToAddresses.length
    ) {
      revert InvalidProposal();
    }
    // The Proposal Set must not exceed the max length
    if (proposalSetJobIds.length > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < proposalSetJobIds.length; i++) {
      // The Proposed address must be a valid address
      if (proposalSetToAddresses[i] == address(0)) {
        revert InvalidProposal();
      }
      // The Proposed address must point to a different address than what is currently set
      if (s_route[proposalSetJobIds[i]] == proposalSetToAddresses[i]) {
        revert InvalidProposal();
      }
      // The from address must match what is the currently set address
      if (s_route[proposalSetJobIds[i]] != proposalSetFromAddresses[i]) {
        revert InvalidProposal();
      }
      // The Router's id cannot be set
      bytes32 routerLabel = keccak256(abi.encodePacked(address(this)));
      if (proposalSetJobIds[i] == routerLabel) {
        revert ReservedLabel(routerLabel);
      }
    }
  }

  /**
   * @inheritdoc IRouterBase
   */
  function propose(
    bytes32[] memory proposalSetJobIds,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) external override onlyOwner {
    _validateProposal(proposalSetJobIds, proposalSetFromAddresses, proposalSetToAddresses);
    uint currentBlock = block.number;
    s_proposalSet = ProposalSet(proposalSetJobIds, proposalSetFromAddresses, proposalSetToAddresses, currentBlock);
    emit Proposed(proposalSetJobIds, proposalSetFromAddresses, proposalSetToAddresses, currentBlock);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function validateProposal(bytes32 jobId, bytes calldata data) external override {
    _smoke(jobId, data);
  }

  /**
   * @dev Must be implemented by inheriting contract
   * Use to test an end to end request through the system
   */
  function _smoke(bytes32 jobId, bytes calldata data) internal virtual returns (bytes32);

  /**
   * @inheritdoc IRouterBase
   */
  function upgrade() external override onlyOwner {
    if (block.number < s_proposalSet.proposedAtBlock + s_timelockBlocks) {
      revert TimelockInEffect();
    }
    for (uint8 i = 0; i < s_proposalSet.jobIds.length; i++) {
      s_route[s_proposalSet.jobIds[i]] = s_proposalSet.to[i];
    }
    s_minorVersion = s_minorVersion + 1;
    if (s_patchVersion != 0) s_patchVersion = 0;
    emit Upgraded(
      s_proposalSet.jobIds,
      s_proposalSet.from,
      s_proposalSet.to,
      block.number,
      s_majorVersion,
      s_minorVersion,
      s_patchVersion
    );
  }

  // ================================================================
  // |                   Config Proposal methods                    |
  // ================================================================
  /**
   * @notice Get the hash of the Router's current configuration
   * @return config hash of config bytes
   */
  function getConfigHash() external view returns (bytes32 config) {
    return s_config_hash;
  }

  /**
   * @dev Must be implemented by inheriting contract
   * Use to set configuration state
   */
  function _setConfig(bytes memory config) internal virtual;

  /**
   * @inheritdoc IRouterBase
   */
  function proposeConfig(bytes32 jobId, bytes calldata config) external override onlyOwner {
    address implAddr = s_route[jobId];
    bytes32 currentConfigHash = IConfigurable(implAddr).getConfigHash(); // TODO: Does this work on self?
    if (currentConfigHash == keccak256(config)) {
      revert InvalidProposal();
    }
    s_proposedConfig[jobId] = ConfigProposal(currentConfigHash, config, block.number);
    bytes32 proposedHash = keccak256(config);
    emit ConfigProposed(jobId, currentConfigHash, proposedHash);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function updateConfig(bytes32 jobId) external override onlyOwner {
    address implAddr = s_route[jobId];
    ConfigProposal memory proposal = s_proposedConfig[jobId];
    if (block.number < proposal.proposedAtBlock + s_timelockBlocks) {
      revert TimelockInEffect();
    }
    if (jobId == keccak256(abi.encodePacked(address(this)))) {
      _setConfig(proposal.to);
    } else {
      IConfigurable(implAddr).setConfig(proposal.to);
    }
    s_patchVersion = s_patchVersion + 1;
    emit ConfigUpdated(jobId, proposal.from, keccak256(proposal.to));
  }

  // ================================================================
  // |                         Timelock methods                     |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function proposeTimelockBlocks(uint16 blocks) external override onlyOwner {
    if (s_timelockBlocks == blocks) {
      revert InvalidProposal();
    }
    if (blocks > MAXIMUM_TIMELOCK_BLOCKS) {
      revert ProposedTimelockAboveMaximum();
    }
    s_timelockProposal = TimeLockProposal(s_timelockBlocks, blocks, block.number);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function updateTimelockBlocks() external override onlyOwner {
    if (s_timelockBlocks == s_timelockProposal.to) {
      revert InvalidProposal();
    }
    if (block.number < s_timelockProposal.proposedAtBlock + s_timelockBlocks) {
      revert TimelockInEffect();
    }
    s_timelockBlocks = s_timelockProposal.to;
    delete s_timelockProposal;
  }

  // ================================================================
  // |                     Pausable methods                         |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function isPaused() public view override returns (bool) {
    return Pausable.paused();
  }

  /**
   * @inheritdoc IRouterBase
   */
  function pause() external override onlyOwner {
    _pause();
  }

  /**
   * @inheritdoc IRouterBase
   */
  function unpause() external override onlyOwner {
    _unpause();
  }
}
