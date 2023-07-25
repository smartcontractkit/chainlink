// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

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
  mapping(bytes32 id => address routableContract) internal s_route;
  error RouteNotFound(bytes32 id);
  // Use empty bytes to self-identify, since it does not have an id
  bytes32 internal constant routerId = bytes32(0);

  // ================================================================
  // |                         Proposal state                       |
  // ================================================================
  uint8 internal constant MAX_PROPOSAL_SET_LENGTH = 8;

  struct ContractProposalSet {
    bytes32[] ids;
    address[] to;
    uint256 timelockEndBlock;
  }
  ContractProposalSet internal s_proposedContractSet;

  event ContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress,
    uint256 timelockEndBlock
  );
  event ContractUpdated(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress,
    uint16 major,
    uint16 minor,
    uint16 patch
  );

  struct ConfigProposal {
    bytes32 fromHash;
    bytes to;
    uint256 timelockEndBlock;
  }
  mapping(bytes32 id => ConfigProposal) internal s_proposedConfig;
  event ConfigProposed(bytes32 id, bytes32 fromHash, bytes toBytes);
  event ConfigUpdated(bytes32 id, bytes32 fromHash, bytes toBytes);
  error InvalidProposal();
  error IdentifierIsReserved(bytes32 id);

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

  error InvalidConfigData();

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
  function version() external view override returns (uint16, uint16, uint16) {
    return (s_majorVersion, s_minorVersion, s_patchVersion);
  }

  // ================================================================
  // |                        Route methods                         |
  // ================================================================

  function _getContractById(bytes32 id, bool useProposed) internal view returns (address) {
    if (!useProposed) {
      address currentImplementation = s_route[id];
      if (currentImplementation != address(0)) {
        return currentImplementation;
      }
    } else {
      for (uint8 i = 0; i < s_proposedContractSet.ids.length; i++) {
        if (id == s_proposedContractSet.ids[i]) {
          // NOTE: proposals can be used immediately
          return s_proposedContractSet.to[i];
        }
      }
    }
    revert RouteNotFound(id);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getContractById(bytes32 id) external view override returns (address routeDestination) {
    routeDestination = _getContractById(id, false);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function getContractById(bytes32 id, bool useProposed) external view override returns (address routeDestination) {
    routeDestination = _getContractById(id, useProposed);
  }

  // ================================================================
  // |                 Contract Proposal methods                    |
  // ================================================================
  /**
   * @inheritdoc IRouterBase
   */
  function getProposedContractSet() external view override returns (uint256, bytes32[] memory, address[] memory) {
    return (s_proposedContractSet.timelockEndBlock, s_proposedContractSet.ids, s_proposedContractSet.to);
  }

  /**
   * @inheritdoc IRouterBase
   */
  function proposeContractsUpdate(
    bytes32[] memory proposedContractSetIds,
    address[] memory proposedContractSetAddresses
  ) external override onlyOwner {
    // All arrays must be of equal length and not must not exceed the max length
    uint256 idsArrayLength = proposedContractSetIds.length;
    if (idsArrayLength != proposedContractSetAddresses.length || idsArrayLength > MAX_PROPOSAL_SET_LENGTH) {
      revert InvalidProposal();
    }
    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < idsArrayLength; i++) {
      bytes32 id = proposedContractSetIds[i];
      address proposedContract = proposedContractSetAddresses[i];
      if (
        proposedContract == address(0) || // The Proposed address must be a valid address
        s_route[id] == proposedContract // The Proposed address must point to a different address than what is currently set
      ) {
        revert InvalidProposal();
      }
      // Reserved ids cannot be set
      if (id == routerId) {
        revert IdentifierIsReserved(id);
      }
    }

    uint256 timelockEndBlock = block.number + s_timelockBlocks;

    s_proposedContractSet = ContractProposalSet(proposedContractSetIds, proposedContractSetAddresses, timelockEndBlock);

    // Iterations will not exceed MAX_PROPOSAL_SET_LENGTH
    for (uint8 i = 0; i < proposedContractSetIds.length; i++) {
      emit ContractProposed(
        proposedContractSetIds[i],
        s_route[proposedContractSetIds[i]],
        proposedContractSetAddresses[i],
        timelockEndBlock
      );
    }
  }

  /**
   * @inheritdoc IRouterBase
   */
  function validateProposedContracts(bytes32 id, bytes calldata data) external override returns (bytes memory) {
    return _validateProposedContracts(id, data);
  }

  /**
   * @dev Must be implemented by the inheriting contract
   * Use to test an end to end request through the system
   */
  function _validateProposedContracts(bytes32 id, bytes calldata data) internal virtual returns (bytes memory);

  /**
   * @inheritdoc IRouterBase
   */
  function updateContracts() external override onlyOwner {
    if (block.number < s_proposedContractSet.timelockEndBlock) {
      revert TimelockInEffect();
    }
    s_minorVersion = s_minorVersion + 1;
    if (s_patchVersion != 0) s_patchVersion = 0;
    for (uint8 i = 0; i < s_proposedContractSet.ids.length; i++) {
      bytes32 id = s_proposedContractSet.ids[i];
      address from = s_route[id];
      address to = s_proposedContractSet.to[i];
      s_route[id] = to;
      emit ContractUpdated(id, from, to, s_majorVersion, s_minorVersion, s_patchVersion);
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
   * Use to set configuration state of the Router
   */
  function _setConfig(bytes memory config) internal virtual;

  /**
   * @inheritdoc IRouterBase
   */
  function proposeConfigUpdate(bytes32 id, bytes calldata config) external override onlyOwner {
    address implAddr = _getContractById(id, false);
    bytes32 currentConfigHash;
    if (implAddr == address(this)) {
      currentConfigHash = s_config_hash;
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
    ConfigProposal memory proposal = s_proposedConfig[id];
    if (block.number < proposal.timelockEndBlock) {
      revert TimelockInEffect();
    }
    if (id == routerId) {
      _setConfig(proposal.to);
      s_config_hash = keccak256(proposal.to);
    } else {
      try IConfigurable(_getContractById(id, false)).setConfig(proposal.to) {} catch {
        revert InvalidConfigData();
      }
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
  function isPaused() external view override returns (bool) {
    return Pausable.paused();
  }

  /**
   * @inheritdoc IRouterBase
   */
  function togglePaused() external override onlyOwner {
    if (Pausable.paused()) {
      _unpause();
    } else {
      _pause();
    }
  }
}
