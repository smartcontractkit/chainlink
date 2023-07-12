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
  mapping(bytes32 => address) internal s_route; /* id => contract address */
  error RouteNotFound(bytes32 id);
  // Use empty bytes to self-identify, since it does not have an id
  bytes32 internal constant routerId = bytes32(0);

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================
  uint8 internal constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ProposalSet {
    bytes32[] ids;
    address[] from;
    address[] to;
    uint256 timelockEndBlock;
  }
  ProposalSet internal s_proposalSet;

  event Proposed(
    bytes32 proposalSetId,
    address proposalSetFromAddress,
    address proposalSetToAddress,
    uint256 timelockEndBlock
  );
  event Upgraded(
    bytes32 proposalSetId,
    address proposalSetFromAddress,
    address proposalSetToAddress,
    uint16 major,
    uint16 minor,
    uint16 patch
  );

  struct ConfigProposal {
    bytes32 fromHash;
    bytes to;
    uint256 timelockEndBlock;
  }
  mapping(bytes32 => ConfigProposal) internal s_proposedConfig; /* id => ConfigProposal */
  event ConfigProposed(bytes32 id, bytes32 fromHash, bytes toBytes);
  event ConfigUpdated(bytes32 id, bytes32 fromHash, bytes toBytes);
  error InvalidProposal();
  error ReservedLabel(bytes32 id);

  // ================================================================
  // |                          Timelock state                      |
  // ================================================================
  uint16 internal MAXIMUM_TIMELOCK_BLOCKS;
  uint16 internal s_timelockBlocks;
  struct TimeLockProposal {
    uint16 from;
    uint16 to;
    uint256 timelockEndBlock;
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
    bytes memory selfConfig
  ) ConfirmedOwnerWithProposal(newOwner, address(0)) Pausable() {
    // Set initial value for the number of blocks of the timelock
    s_timelockBlocks = timelockBlocks;
    // Set maximum number of blocks that the timelock can be
    // NOTE: this cannot be later modified
    MAXIMUM_TIMELOCK_BLOCKS = maximumTimelockBlocks;
    // Set the initial configuration for the Router
    s_route[routerId] = address(this);
    _setConfig(selfConfig);
    s_config_hash = keccak256(selfConfig);
  }

  // ================================================================
  // |                       Version methods                        |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function version() external view returns (uint16, uint16, uint16) {
    return (s_majorVersion, s_minorVersion, s_patchVersion);
  }

  // ================================================================
  // |                        Route methods                         |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(bytes32 id) public view returns (address) {
    address currentImplementation = s_route[id];
    if (currentImplementation == address(0)) {
      revert RouteNotFound(id);
    }
    return currentImplementation;
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(bytes32 id, bool useProposed) public view returns (address) {
    if (useProposed == true) {
      return getRoute(id);
    }

    for (uint8 i = 0; i < s_proposalSet.ids.length; i++) {
      if (id == s_proposalSet.ids[i]) {
        // NOTE: proposals can be used immediately
        return s_proposalSet.to[i];
      }
    }
    revert RouteNotFound(id);
  }

  // ================================================================
  // |                 Contract Proposal methods                    |
  // ================================================================

  /**
   * @inheritdoc IRouterBase
   */
  function getProposalSet() external view returns (uint256, bytes32[] memory, address[] memory, address[] memory) {
    return (s_proposalSet.timelockEndBlock, s_proposalSet.ids, s_proposalSet.from, s_proposalSet.to);
  }

  /**
   * @dev Helper function to validate a proposal set
   */
  function _validateProposalSet(
    bytes32[] memory proposalSetIds,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) internal view {
    // All arrays must be of equal length
    if (
      proposalSetIds.length != proposalSetFromAddresses.length || proposalSetIds.length != proposalSetToAddresses.length
    ) {
      revert InvalidProposal();
    }
    // The Proposal Set must not exceed the max length
    if (proposalSetIds.length > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < proposalSetIds.length; i++) {
      // The Proposed address must be a valid address
      if (proposalSetToAddresses[i] == address(0)) {
        revert InvalidProposal();
      }
      // The Proposed address must point to a different address than what is currently set
      if (s_route[proposalSetIds[i]] == proposalSetToAddresses[i]) {
        revert InvalidProposal();
      }
      // The from address must match what is the currently set address
      if (s_route[proposalSetIds[i]] != proposalSetFromAddresses[i]) {
        revert InvalidProposal();
      }
      // The Router's id cannot be set
      if (proposalSetIds[i] == routerId) {
        revert ReservedLabel(routerId);
      }
    }
  }

  /**
   * @inheritdoc IRouterBase
   */
  function propose(
    bytes32[] memory proposalSetIds,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) external override onlyOwner {
    _validateProposalSet(proposalSetIds, proposalSetFromAddresses, proposalSetToAddresses);
    uint timelockEndBlock = block.number + s_timelockBlocks;
    s_proposalSet = ProposalSet(proposalSetIds, proposalSetFromAddresses, proposalSetToAddresses, timelockEndBlock);
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < proposalSetIds.length; i++) {
      emit Proposed(proposalSetIds[i], proposalSetFromAddresses[i], proposalSetToAddresses[i], timelockEndBlock);
    }
  }

  /**
   * @inheritdoc IRouterBase
   */
  function validateProposal(bytes32 id, bytes calldata data) external override {
    _validateProposal(id, data);
  }

  /**
   * @dev Must be implemented by the inheriting contract
   * Use to test an end to end request through the system
   */
  function _validateProposal(bytes32 id, bytes calldata data) internal virtual returns (bytes32);

  /**
   * @inheritdoc IRouterBase
   */
  function upgrade() external override onlyOwner {
    if (block.number < s_proposalSet.timelockEndBlock) {
      revert TimelockInEffect();
    }
    s_minorVersion = s_minorVersion + 1;
    if (s_patchVersion != 0) s_patchVersion = 0;
    for (uint8 i = 0; i < s_proposalSet.ids.length; i++) {
      s_route[s_proposalSet.ids[i]] = s_proposalSet.to[i];
      emit Upgraded(
        s_proposalSet.ids[i],
        s_proposalSet.from[i],
        s_proposalSet.to[i],
        s_majorVersion,
        s_minorVersion,
        s_patchVersion
      );
    }
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
  function proposeConfig(bytes32 id, bytes calldata config) external override onlyOwner {
    address implAddr = s_route[id];
    bytes32 currentConfigHash;
    if (implAddr == address(this)) {
      currentConfigHash = this.getConfigHash();
    } else {
      currentConfigHash = IConfigurable(implAddr).getConfigHash();
    }
    if (currentConfigHash == keccak256(config)) {
      revert InvalidProposal();
    }
    s_proposedConfig[id] = ConfigProposal(currentConfigHash, config, block.number + s_timelockBlocks);
    emit ConfigProposed(id, currentConfigHash, config);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function updateConfig(bytes32 id) external override onlyOwner {
    address implAddr = s_route[id];
    ConfigProposal memory proposal = s_proposedConfig[id];
    if (block.number < proposal.timelockEndBlock) {
      revert TimelockInEffect();
    }
    if (id == routerId) {
      _setConfig(proposal.to);
    } else {
      IConfigurable(implAddr).setConfig(proposal.to);
    }
    s_patchVersion = s_patchVersion + 1;
    emit ConfigUpdated(id, proposal.fromHash, proposal.to);
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
    s_timelockProposal = TimeLockProposal(s_timelockBlocks, blocks, block.number + s_timelockBlocks);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function updateTimelockBlocks() external override onlyOwner {
    if (block.number < s_timelockProposal.timelockEndBlock) {
      revert TimelockInEffect();
    }
    s_timelockBlocks = s_timelockProposal.to;
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
