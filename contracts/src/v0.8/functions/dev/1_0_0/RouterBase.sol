// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IVersioned, Versioned} from "./Versioned.sol";
import {IRouterBase} from "./interfaces/IRouterBase.sol";
import {ConfirmedOwnerWithProposal} from "../../../ConfirmedOwnerWithProposal.sol";
import {Pausable} from "../../../shared/vendor/openzeppelin-solidity/v.4.8.0/contracts/security/Pausable.sol";
import {IConfigurable} from "./interfaces/IConfigurable.sol";

abstract contract RouterBase is IRouterBase, Pausable, Versioned, ConfirmedOwnerWithProposal {
  // ================================================================
  // |                         Version state                        |
  // ================================================================
  uint16 internal constant s_majorVersion = 1;
  uint16 internal s_minorVersion = 0;
  uint16 internal s_patchVersion = 0;

  // ================================================================
  // |                          Route state                         |
  // ================================================================
  mapping(string => address) internal s_route;
  error RouteNotFound(string label);

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================
  uint8 internal constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ProposalSet {
    string[] labels;
    address[] from;
    address[] to;
    uint proposedAtBlock;
  }
  ProposalSet internal s_proposalSet;

  event Proposed(
    string[] proposalSetLabels,
    address[] proposalSetFromAddresses,
    address[] proposalSetToAddresses,
    uint block
  );
  event Upgraded(
    string[] proposalSetLabels,
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
  mapping(string => ConfigProposal) internal s_proposedConfig;
  event ConfigProposed(string label, bytes32 from, bytes32 to);
  event ConfigUpdated(string label, bytes32 from, bytes32 to);
  error InvalidProposal();
  error ReservedLabel(string label);

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
    string memory id,
    address newOwner,
    uint16 timelockBlocks,
    uint16 maximumTimelockBlocks,
    string[] memory initialLabels,
    address[] memory initialAddresses,
    bytes memory config
  ) ConfirmedOwnerWithProposal(newOwner, address(0)) Pausable() Versioned(id, s_majorVersion) {
    s_timelockBlocks = timelockBlocks;
    MAXIMUM_TIMELOCK_BLOCKS = maximumTimelockBlocks;
    s_route[id] = address(this);
    _setConfig(config);
    s_config_hash = keccak256(config);
    address[] memory emptyAddresses = new address[](initialLabels.length);
    _validateProposal(initialLabels, emptyAddresses, initialAddresses);
    s_proposalSet = ProposalSet(initialLabels, emptyAddresses, initialAddresses, block.number);
    for (uint8 i = 0; i < s_proposalSet.labels.length; i++) {
      s_route[s_proposalSet.labels[i]] = s_proposalSet.to[i];
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

  function getLatestRoute(string calldata label) internal view returns (address) {
    address currentImplementation = s_route[label];
    if (currentImplementation == address(0)) {
      revert RouteNotFound(label);
    }
    return currentImplementation;
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(string calldata label) public view returns (address) {
    return getLatestRoute(label);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getRoute(string calldata label, bool useProposed) public view returns (address) {
    if (useProposed == true) {
      for (uint8 i = 0; i < s_proposalSet.labels.length; i++) {
        if (keccak256(abi.encodePacked(label)) == keccak256(abi.encodePacked(s_proposalSet.labels[i]))) {
          // NOTE: proposals can be used immediately
          // if (block.number < s_proposalSet.proposedAtBlock + s_timelockBlocks) {
          //   revert TimelockInEffect();
          // }
          return s_proposalSet.to[i];
        }
      }
    }
    return getLatestRoute(label);
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
      string[] memory,
      address[] memory,
      address[] memory
    )
  {
    return (s_proposalSet.proposedAtBlock, s_proposalSet.labels, s_proposalSet.from, s_proposalSet.to);
  }

  function _validateProposal(
    string[] memory proposalSetLabels,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) internal view {
    // All arrays must be of equal length
    if (
      proposalSetLabels.length != proposalSetFromAddresses.length ||
      proposalSetLabels.length != proposalSetToAddresses.length
    ) {
      revert InvalidProposal();
    }
    // The Proposal Set must not exceed the max length
    if (proposalSetLabels.length > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < proposalSetLabels.length; i++) {
      // The Proposed address must be a valid address
      if (proposalSetToAddresses[i] == address(0)) {
        revert InvalidProposal();
      }
      // The Proposed address must point to a different address than what is currently set
      if (s_route[proposalSetLabels[i]] == proposalSetToAddresses[i]) {
        revert InvalidProposal();
      }
      // The from address must match what is the currently set address
      if (s_route[proposalSetLabels[i]] != proposalSetFromAddresses[i]) {
        revert InvalidProposal();
      }
      // The Router's id cannot be set
      if (keccak256(abi.encodePacked(proposalSetLabels[i])) == keccak256(abi.encodePacked(s_id))) {
        revert ReservedLabel(s_id);
      }
    }
  }

  /**
   * @inheritdoc IRouterBase
   */
  function propose(
    string[] memory proposalSetLabels,
    address[] memory proposalSetFromAddresses,
    address[] memory proposalSetToAddresses
  ) external override onlyOwner {
    _validateProposal(proposalSetLabels, proposalSetFromAddresses, proposalSetToAddresses);
    uint currentBlock = block.number;
    s_proposalSet = ProposalSet(proposalSetLabels, proposalSetFromAddresses, proposalSetToAddresses, currentBlock);
    emit Proposed(proposalSetLabels, proposalSetFromAddresses, proposalSetToAddresses, currentBlock);
  }

  /**
   * @dev Must be implemented by inheriting contract
   * Use to test an end to end request through the system
   */
  function _smoke(bytes calldata data) internal virtual returns (bytes32);

  /**
   * @inheritdoc IRouterBase
   */
  function validateProposal(bytes calldata data) external override {
    _smoke(data);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function upgrade() external override onlyOwner {
    if (block.number < s_proposalSet.proposedAtBlock + s_timelockBlocks) {
      revert TimelockInEffect();
    }
    for (uint8 i = 0; i < s_proposalSet.labels.length; i++) {
      s_route[s_proposalSet.labels[i]] = s_proposalSet.to[i];
    }
    s_minorVersion = s_minorVersion + 1;
    if (s_patchVersion != 0) s_patchVersion = 0;
    emit Upgraded(
      s_proposalSet.labels,
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
   * @notice Get the hash of the current configuration
   * @return config hash of config bytes
   */
  function getConfigHash() external view returns (bytes32 config) {
    return s_config_hash;
  }

  function _setConfig(bytes memory config) internal virtual;

  /**
   * @inheritdoc IRouterBase
   */
  function proposeConfig(string calldata label, bytes calldata config) external override onlyOwner {
    address implAddr = s_route[label];
    bytes32 currentConfigHash = IConfigurable(implAddr).getConfigHash(); // TODO: Does this work on self?
    if (currentConfigHash == keccak256(config)) {
      revert InvalidProposal();
    }
    s_proposedConfig[label] = ConfigProposal(currentConfigHash, config, block.number);
    bytes32 proposedHash = keccak256(config);
    emit ConfigProposed(label, currentConfigHash, proposedHash);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function updateConfig(string calldata label) external override onlyOwner {
    address implAddr = s_route[label];
    ConfigProposal memory proposal = s_proposedConfig[label];
    if (block.number < proposal.proposedAtBlock + s_timelockBlocks) {
      revert TimelockInEffect();
    }
    if (keccak256(abi.encodePacked(label)) == keccak256(abi.encodePacked(s_id))) {
      _setConfig(proposal.to);
    } else {
      IConfigurable(implAddr).setConfig(proposal.to);
    }
    s_patchVersion = s_patchVersion + 1;
    emit ConfigUpdated(label, proposal.from, keccak256(proposal.to));
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
