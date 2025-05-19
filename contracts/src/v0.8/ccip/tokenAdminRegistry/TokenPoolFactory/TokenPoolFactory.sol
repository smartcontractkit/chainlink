// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IOwnable} from "../../../shared/interfaces/IOwnable.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {RateLimiter} from "../../libraries/RateLimiter.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {RegistryModuleOwnerCustom} from "../RegistryModuleOwnerCustom.sol";
import {FactoryBurnMintERC20} from "./FactoryBurnMintERC20.sol";

import {Create2} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/Create2.sol";

/// @notice A contract for deploying new tokens and token pools, and configuring them with the token admin registry
/// @dev At the end of the transaction, the ownership transfer process will begin, but the user must accept the
/// ownership transfer in a separate transaction.
/// @dev The address prediction mechanism is only capable of deploying and predicting addresses for EVM based chains.
/// adding compatibility for other chains will require additional offchain computation.
contract TokenPoolFactory is ITypeAndVersion {
  using Create2 for bytes32;

  event RemoteChainConfigUpdated(uint64 indexed remoteChainSelector, RemoteChainConfig remoteChainConfig);

  error InvalidZeroAddress();

  /// @notice The type of pool to deploy. Types may be expanded in future versions
  enum PoolType {
    BURN_MINT,
    LOCK_RELEASE
  }

  /// @dev This struct will only ever exist in memory and as calldata, and therefore does not need to be efficiently packed for storage. The struct is used to pass information to the create2 address generation function.
  struct RemoteTokenPoolInfo {
    uint64 remoteChainSelector; // The CCIP specific selector for the remote chain
    bytes remotePoolAddress; // The address of the remote pool to either deploy or use as is. If empty, address
    // will be predicted
    bytes remotePoolInitCode; // Remote pool creation code if it needs to be deployed, without constructor params
    // appended to the end.
    RemoteChainConfig remoteChainConfig; // The addresses of the remote RMNProxy, Router, factory, and token
    // decimals which are needed for determining the remote address
    PoolType poolType; // The type of pool to deploy, either Burn/Mint or Lock/Release
    bytes remoteTokenAddress; // EVM address for remote token. If empty, the address will be predicted
    bytes remoteTokenInitCode; // The init code to be deployed on the remote chain and includes constructor params
    RateLimiter.Config rateLimiterConfig; // Token Pool rate limit. Values will be applied on incoming an outgoing messages
  }

  // solhint-disable-next-line gas-struct-packing
  struct RemoteChainConfig {
    address remotePoolFactory; // The factory contract on the remote chain which will make the deployment
    address remoteRouter; // The router on the remote chain
    address remoteRMNProxy; // The RMNProxy contract on the remote chain
    uint8 remoteTokenDecimals; // The number of decimals for the token on the remote chain
  }

  string public constant typeAndVersion = "TokenPoolFactory 1.5.1";

  ITokenAdminRegistry private immutable i_tokenAdminRegistry;
  RegistryModuleOwnerCustom private immutable i_registryModuleOwnerCustom;

  address private immutable i_rmnProxy;
  address private immutable i_ccipRouter;

  /// @notice Construct the TokenPoolFactory
  /// @param tokenAdminRegistry The address of the token admin registry
  /// @param tokenAdminModule The address of the token admin module which can register the token via ownership module
  /// @param rmnProxy The address of the RMNProxy contract token pools will be deployed with
  /// @param ccipRouter The address of the CCIPRouter contract token pools will be deployed with
  constructor(
    ITokenAdminRegistry tokenAdminRegistry,
    RegistryModuleOwnerCustom tokenAdminModule,
    address rmnProxy,
    address ccipRouter
  ) {
    if (
      address(tokenAdminRegistry) == address(0) || address(tokenAdminModule) == address(0) || rmnProxy == address(0)
        || ccipRouter == address(0)
    ) revert InvalidZeroAddress();

    i_tokenAdminRegistry = ITokenAdminRegistry(tokenAdminRegistry);
    i_registryModuleOwnerCustom = RegistryModuleOwnerCustom(tokenAdminModule);
    i_rmnProxy = rmnProxy;
    i_ccipRouter = ccipRouter;
  }

  // ================================================================
  // │                   Top-Level Deployment                       │
  // ================================================================

  /// @notice Deploys a token and token pool with the given token information and configures it with remote token pools
  /// @dev The token and token pool are deployed in the same transaction, and the token pool is configured with the
  /// remote token pools. The token pool is then set in the token admin registry. Ownership of the everything is transferred
  /// to the msg.sender, but must be accepted in a separate transaction due to 2-step ownership transfer.
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// or to be predicted if the pool has not been deployed yet on the remote chain
  /// @param localTokenDecimals The amount of decimals to be used in the new token. Since decimals() is not part of the
  /// the ERC20 standard, and thus cannot be certain to exist, the amount must be supplied via user input.
  /// @param tokenInitCode The creation code for the token, which includes the constructor parameters already appended
  /// @param tokenPoolInitCode The creation code for the token pool, without the constructor parameters appended
  /// @param salt The salt to be used in the create2 deployment of the token and token pool to ensure a unique address
  /// @return token The address of the token that was deployed
  /// @return pool The address of the token pool that was deployed
  function deployTokenAndTokenPool(
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    uint8 localTokenDecimals,
    bytes memory tokenInitCode,
    bytes calldata tokenPoolInitCode,
    bytes32 salt
  ) external returns (address, address) {
    // Ensure a unique deployment between senders even if the same input parameter is used to prevent
    // DOS/front running attacks
    salt = keccak256(abi.encodePacked(salt, msg.sender));

    // Deploy the token. The constructor parameters are already provided in the tokenInitCode
    address token = Create2.deploy(0, salt, tokenInitCode);

    // Deploy the token pool
    address pool =
      _createTokenPool(token, localTokenDecimals, remoteTokenPools, tokenPoolInitCode, salt, PoolType.BURN_MINT);

    // Grant the mint and burn roles to the pool for the token
    FactoryBurnMintERC20(token).grantMintAndBurnRoles(pool);

    // Set the token pool for token in the token admin registry since this contract is the token and pool owner
    _setTokenPoolInTokenAdminRegistry(token, pool);

    // Begin the 2 step ownership transfer of the newly deployed token to the msg.sender
    IOwnable(token).transferOwnership(msg.sender);

    return (token, pool);
  }

  /// @notice Deploys a token pool with an existing ERC20 token
  /// @dev Since the token already exists, this contract is not the owner and therefore cannot configure the
  /// token pool in the token admin registry in the same transaction. The user must invoke the calls to the
  /// tokenAdminRegistry manually
  /// @dev since the token already exists, the owner must grant the mint and burn roles to the pool manually
  /// @param token The address of the existing token to be used in the token pool
  /// @param localTokenDecimals The amount of decimals used in the existing token. Since decimals() is not part of the
  /// the ERC20 standard, and thus cannot be certain to exist, the amount must be supplied via user input.
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// @param tokenPoolInitCode The creation code for the token pool
  /// @param salt The salt to be used in the create2 deployment of the token pool
  /// @return poolAddress The address of the token pool that was deployed
  function deployTokenPoolWithExistingToken(
    address token,
    uint8 localTokenDecimals,
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    bytes calldata tokenPoolInitCode,
    bytes32 salt,
    PoolType poolType
  ) external returns (address poolAddress) {
    // Ensure a unique deployment between senders even if the same input parameter is used to prevent
    // DOS/front running attacks
    salt = keccak256(abi.encodePacked(salt, msg.sender));

    // create the token pool and return the address
    return _createTokenPool(token, localTokenDecimals, remoteTokenPools, tokenPoolInitCode, salt, poolType);
  }

  // ================================================================
  // │                Pool Deployment/Configuration                 │
  // ================================================================

  /// @notice Deploys a token pool with the given token information and remote token pools
  /// @param token The token to be used in the token pool
  /// @param remoteTokenPools An array of remote token pools info to be used in the pool's applyChainUpdates function
  /// @param tokenPoolInitCode The creation code for the token pool
  /// @param salt The salt to be used in the create2 deployment of the token pool
  /// @return poolAddress The address of the token pool that was deployed
  function _createTokenPool(
    address token,
    uint8 localTokenDecimals,
    RemoteTokenPoolInfo[] calldata remoteTokenPools,
    bytes calldata tokenPoolInitCode,
    bytes32 salt,
    PoolType poolType
  ) private returns (address) {
    // Create an array of chain updates to apply to the token pool
    TokenPool.ChainUpdate[] memory chainUpdates = new TokenPool.ChainUpdate[](remoteTokenPools.length);

    RemoteTokenPoolInfo memory remoteTokenPool;
    for (uint256 i = 0; i < remoteTokenPools.length; ++i) {
      remoteTokenPool = remoteTokenPools[i];

      // If the user provides an empty byte string, indicated no token has already been deployed,
      // then the address of the token needs to be predicted. Otherwise the address provided will be used.
      if (remoteTokenPool.remoteTokenAddress.length == 0) {
        // The user must provide the initCode for the remote token, so its address can be predicted correctly. It's
        // provided in the remoteTokenInitCode field for the remoteTokenPool
        remoteTokenPool.remoteTokenAddress = abi.encode(
          salt.computeAddress(
            keccak256(remoteTokenPool.remoteTokenInitCode), remoteTokenPool.remoteChainConfig.remotePoolFactory
          )
        );
      }

      // If the user provides an empty byte string parameter, indicating the pool has not been deployed yet,
      // the address of the pool should be predicted. Otherwise use the provided address.
      if (remoteTokenPool.remotePoolAddress.length == 0) {
        // Address is predicted based on the init code hash and the deployer, so the hash must first be computed
        // using the initCode and a concatenated set of constructor parameters.
        bytes32 remotePoolInitcodeHash = _generatePoolInitcodeHash(
          remoteTokenPool.remotePoolInitCode,
          remoteTokenPool.remoteChainConfig,
          abi.decode(remoteTokenPool.remoteTokenAddress, (address)),
          remoteTokenPool.poolType
        );

        // Abi encode the computed remote address so it can be used as bytes in the chain update
        remoteTokenPool.remotePoolAddress =
          abi.encode(salt.computeAddress(remotePoolInitcodeHash, remoteTokenPool.remoteChainConfig.remotePoolFactory));
      }

      bytes[] memory remotePoolAddresses = new bytes[](1);
      remotePoolAddresses[0] = remoteTokenPool.remotePoolAddress;

      chainUpdates[i] = TokenPool.ChainUpdate({
        remoteChainSelector: remoteTokenPool.remoteChainSelector,
        remotePoolAddresses: remotePoolAddresses,
        remoteTokenAddress: remoteTokenPool.remoteTokenAddress,
        outboundRateLimiterConfig: remoteTokenPool.rateLimiterConfig,
        inboundRateLimiterConfig: remoteTokenPool.rateLimiterConfig
      });
    }

    // Construct the initArgs for the token pool using the immutable contracts for CCIP on the local chain
    bytes memory tokenPoolInitArgs;
    if (poolType == PoolType.BURN_MINT) {
      tokenPoolInitArgs = abi.encode(token, localTokenDecimals, new address[](0), i_rmnProxy, i_ccipRouter);
    } else if (poolType == PoolType.LOCK_RELEASE) {
      // Lock/Release pools have an additional boolean constructor parameter that must be accounted for, acceptLiquidity,
      // which is set to true by default in this case. Users wishing to set it to false must deploy the pool manually.
      tokenPoolInitArgs = abi.encode(token, localTokenDecimals, new address[](0), i_rmnProxy, true, i_ccipRouter);
    }

    // Construct the deployment code from the initCode and the initArgs and then deploy
    address poolAddress = Create2.deploy(0, salt, abi.encodePacked(tokenPoolInitCode, tokenPoolInitArgs));

    // Apply the chain updates to the token pool
    TokenPool(poolAddress).applyChainUpdates(new uint64[](0), chainUpdates);

    // Begin the 2 step ownership transfer of the token pool to the msg.sender.
    IOwnable(poolAddress).transferOwnership(address(msg.sender)); // 2 step ownership transfer

    return poolAddress;
  }

  /// @notice Generates the hash of the init code the pool will be deployed with
  /// @dev The init code hash is used with Create2 to predict the address of the pool on the remote chain
  /// @dev ABI-encoding limitations prevent arbitrary constructor parameters from being used, so pool type must be
  /// restricted to those with known types in the constructor. This function should be updated if new pool types are needed.
  /// @param initCode The init code of the pool
  /// @param remoteChainConfig The remote chain config for the pool
  /// @param remoteTokenAddress The address of the remote token
  /// @param poolType The type of pool to deploy
  /// @return bytes32 hash of the init code to be used in the deterministic address calculation
  function _generatePoolInitcodeHash(
    bytes memory initCode,
    RemoteChainConfig memory remoteChainConfig,
    address remoteTokenAddress,
    PoolType poolType
  ) private pure returns (bytes32) {
    if (poolType == PoolType.BURN_MINT) {
      return keccak256(
        abi.encodePacked(
          initCode,
          // constructor(address token, uint8 localTokenDecimals, address[] allowlist, address rmnProxy, address router)
          abi.encode(
            remoteTokenAddress,
            remoteChainConfig.remoteTokenDecimals,
            new address[](0),
            remoteChainConfig.remoteRMNProxy,
            remoteChainConfig.remoteRouter
          )
        )
      );
    } else {
      // if poolType is PoolType.LOCK_RELEASE, but may be expanded in future versions
      return keccak256(
        abi.encodePacked(
          initCode,
          // constructor(address token, uint8 localTokenDecimals, address[] allowList, address rmnProxy, bool acceptLiquidity, address router)
          abi.encode(
            remoteTokenAddress,
            remoteChainConfig.remoteTokenDecimals,
            new address[](0),
            remoteChainConfig.remoteRMNProxy,
            true,
            remoteChainConfig.remoteRouter
          )
        )
      );
    }
  }

  /// @notice Sets the token pool address in the token admin registry for a newly deployed token pool.
  /// @dev this function should only be called when the token is deployed by this contract as well, otherwise
  /// the token pool will not be able to be set in the token admin registry, and this function will revert.
  /// @param token The address of the token to set the pool for
  /// @param pool The address of the pool to set in the token admin registry
  function _setTokenPoolInTokenAdminRegistry(address token, address pool) private {
    i_registryModuleOwnerCustom.registerAdminViaOwner(token);
    i_tokenAdminRegistry.acceptAdminRole(token);
    i_tokenAdminRegistry.setPool(token, pool);

    // Begin the 2 admin transfer process which must be accepted in a separate tx.
    i_tokenAdminRegistry.transferAdminRole(token, msg.sender);
  }
}
