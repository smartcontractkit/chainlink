pragma solidity ^0.8.24;

import {OwnerIsCreator} from "../../../shared/access/OwnerIsCreator.sol";
import {IBurnMintERC20} from "../../../shared/token/ERC20/IBurnMintERC20.sol";

import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

import {Router} from "../../Router.sol";

/// @notice Allows migration of a lane in a token pool from Lock/Release to CCTP supported Burn/Mint. Contract
/// functionality is based on hard requirements defined by Circle to allow future CCTP compatibility
/// @dev Once a migration for a lane has occured, it can never be reversed, and CCTP will be the mechanism forever. This makes the assumption that Circle will continue to support that lane indefinitely.
abstract contract USDCBridgeMigrator is OwnerIsCreator {
  using EnumerableSet for EnumerableSet.UintSet;

  event CCTPMigrationProposed(uint64 remoteChainSelector);
  event CCTPMigrationExecuted(uint64 remoteChainSelector, uint256 USDCBurned);
  event CCTPMigrationCancelled(uint64 existingProposalSelector);
  event CircleMigratorAddressSet(address migratorAddress);

  error onlyCircle();
  error ExistingMigrationProposal();
  error NoExistingMigrationProposal();
  error NoMigrationProposalPending();
  error InvalidChainSelector(uint64 remoteChainSelector);

  IBurnMintERC20 internal immutable i_USDC;
  Router internal immutable i_router;

  address internal s_circleUSDCMigrator;
  uint64 internal s_proposedUSDCMigrationChain;

  mapping(uint64 chainSelector => uint256 lockedBalance) internal s_lockedTokensByChainSelector;

  mapping(uint64 chainSelector => bool shouldUseLockRelease) internal s_shouldUseLockRelease;

  constructor(address token, address router) {
    i_USDC = IBurnMintERC20(token);
    i_router = Router(router);
  }

  /// @notice Burn USDC locked for a specific lane so that destination USDC can be converted from
  /// non-canonical to canonical USDC.
  /// @dev This function can only be called by an address specified by the owner to be controlled by circle
  /// @dev proposeCCTPMigration must be called first on an approved lane to execute properly.
  /// @dev This function signature should NEVER be overwritten, otherwise it will be unable to be called by
  /// circle to properly migrate USDC over to CCTP.
  function burnLockedUSDC() public {
    if (msg.sender != s_circleUSDCMigrator) revert onlyCircle();
    if (s_proposedUSDCMigrationChain == 0) revert ExistingMigrationProposal();

    uint64 burnChainSelector = s_proposedUSDCMigrationChain;
    uint256 tokensToBurn = s_lockedTokensByChainSelector[burnChainSelector];

    // Even though USDC is a trusted call, ensure CEI by updating state first
    delete s_lockedTokensByChainSelector[burnChainSelector];
    delete s_proposedUSDCMigrationChain;

    // This should only be called after this contract has been granted a "zero allowance minter role" on USDC by Circle,
    // otherwise the call will revert. Executing this burn will functionally convert all USDC on the destination chain
    // to canonical USDC by removing the canonical USDC backing it from circulation.
    i_USDC.burn(tokensToBurn);

    // Disable L/R automatically on burned chain and enable CCTP
    delete s_shouldUseLockRelease[burnChainSelector];

    emit CCTPMigrationExecuted(burnChainSelector, tokensToBurn);
  }

  /// @notice Propose a destination chain to migrate from lock/release mechanism to CCTP enabled burn/mint
  /// through a Circle controlled burn.
  /// @param remoteChainSelector the CCIP specific selector for the remote chain currently using a
  /// non-canonical form of USDC which they wish to update to canonical. Function will revert if the chain
  /// selector is zero, or if a migration has already occured for the specified selector.
  /// @dev This function can only be called by the owner
  function proposeCCTPMigration(uint64 remoteChainSelector) external onlyOwner {
    // Prevent overwriting existing migration proposals until the current one is finished
    if (s_proposedUSDCMigrationChain != 0) revert ExistingMigrationProposal();

    s_proposedUSDCMigrationChain = remoteChainSelector;

    emit CCTPMigrationProposed(remoteChainSelector);
  }

  /// @notice Cancel an existing proposal to migrate a lane to CCTP.
  function cancelExistingCCTPMigrationProposal() external onlyOwner {
    if (s_proposedUSDCMigrationChain == 0) revert NoExistingMigrationProposal();

    uint64 currentProposalChainSelector = s_proposedUSDCMigrationChain;
    delete s_proposedUSDCMigrationChain;

    emit CCTPMigrationCancelled(currentProposalChainSelector);
  }

  /// @notice retrieve the chain selector for an ongoing CCTP migration in progress.
  /// @return uint64 the chain selector of the lane to be migrated. Will be zero if no proposal currently
  /// exists
  function getCurrentProposedCCTPChainMigration() public view returns (uint64) {
    return s_proposedUSDCMigrationChain;
  }

  /// @notice Set the address of the circle-controlled wallet which will execute a CCTP lane migration
  /// @dev The function should only be invoked once the address has been confirmed by Circle prior to
  /// chain expansion.
  function setCircleMigratorAddress(address migrator) external onlyOwner {
    s_circleUSDCMigrator = migrator;

    emit CircleMigratorAddressSet(migrator);
  }

  /// @notice Retrieve the amount of canonical USDC locked into this lane and minted on the destination
  /// @param remoteChainSelector the CCIP specific destination chain implementing a mintable and
  /// non-canonical form of USDC at present.
  /// @return uint256 the amount of USDC locked into the specified lane. If non-zero, the number
  /// should match the current circulating supply of USDC on the destination chain
  function getLockedTokensForChain(uint64 remoteChainSelector) public view returns (uint256) {
    return s_lockedTokensByChainSelector[remoteChainSelector];
  }
}
