// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IEVM2AnyOnRampClient} from "../interfaces/IEVM2AnyOnRampClient.sol";
import {IMessageInterceptor} from "../interfaces/IMessageInterceptor.sol";
import {INonceManager} from "../interfaces/INonceManager.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../libraries/USDPriceWith18Decimals.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice The OnRamp is a contract that handles lane-specific fee logic
/// @dev The OnRamp, MultiCommitStore and OffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
contract OnRamp is IEVM2AnyOnRampClient, ITypeAndVersion, OwnerIsCreator {
  using SafeERC20 for IERC20;
  using USDPriceWith18Decimals for uint224;

  error CannotSendZeroTokens();
  error UnsupportedToken(address token);
  error MustBeCalledByRouter();
  error RouterMustSetOriginalSender();
  error InvalidConfig();
  error CursedByRMN(uint64 sourceChainSelector);
  error GetSupportedTokensFunctionalityRemovedCheckAdminRegistry();
  error InvalidDestChainConfig(uint64 sourceChainSelector);

  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event DestChainConfigSet(uint64 indexed destChainSelector, DestChainConfig destChainConfig);
  event FeePaid(address indexed feeToken, uint256 feeValueJuels);
  event FeeTokenWithdrawn(address indexed feeAggregator, address indexed feeToken, uint256 amount);
  /// RMN depends on this event, if changing, please notify the RMN maintainers.
  event CCIPSendRequested(uint64 indexed destChainSelector, Internal.EVM2AnyRampMessage message);

  /// @dev Struct that contains the static configuration
  /// RMN depends on this struct, if changing, please notify the RMN maintainers.
  // solhint-disable-next-line gas-struct-packing
  struct StaticConfig {
    uint64 chainSelector; // ─────╮ Source chainSelector
    address rmnProxy; // ─────────╯ Address of RMN proxy
    address nonceManager; // Address of the nonce manager
    address tokenAdminRegistry; // Token admin registry address
  }

  /// @dev Struct to contains the dynamic configuration
  // solhint-disable-next-line gas-struct-packing
  struct DynamicConfig {
    address priceRegistry; // Price registry address
    address messageValidator; // Optional message validator to validate outbound messages (zero address = no validator)
    address feeAggregator; // Fee aggregator address
  }

  /// @dev Struct to hold the configs for a destination chain
  struct DestChainConfig {
    // The last used sequence number. This is zero in the case where no messages has been sent yet.
    // 0 is not a valid sequence number for any real transaction.
    uint64 sequenceNumber;
    // This is the local router address that is allowed to send messages to the destination chain.
    // This is NOT the receiving router address on the destination chain.
    IRouter router;
  }

  /// @dev Same as DestChainConfig but with the destChainSelector so that an array of these
  /// can be passed in the constructor and the applyDestChainConfigUpdates function
  //solhint-disable gas-struct-packing
  struct DestChainConfigArgs {
    uint64 destChainSelector; // Destination chain selector
    IRouter router; // Source router address
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "OnRamp 1.6.0-dev";
  /// @dev The chain ID of the source chain that this contract is deployed to
  uint64 internal immutable i_chainSelector;
  /// @dev The address of the rmn proxy
  address internal immutable i_rmnProxy;
  /// @dev The address of the nonce manager
  address internal immutable i_nonceManager;
  /// @dev The address of the token admin registry
  address internal immutable i_tokenAdminRegistry;
  /// @dev the maximum number of nops that can be configured at the same time.

  // DYNAMIC CONFIG
  /// @dev The config for the onRamp
  DynamicConfig internal s_dynamicConfig;

  /// @dev The destination chain specific configs
  mapping(uint64 destChainSelector => DestChainConfig destChainConfig) internal s_destChainConfigs;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigArgs
  ) {
    if (
      staticConfig.chainSelector == 0 || staticConfig.rmnProxy == address(0) || staticConfig.nonceManager == address(0)
        || staticConfig.tokenAdminRegistry == address(0)
    ) {
      revert InvalidConfig();
    }

    i_chainSelector = staticConfig.chainSelector;
    i_rmnProxy = staticConfig.rmnProxy;
    i_nonceManager = staticConfig.nonceManager;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;

    _setDynamicConfig(dynamicConfig);
    _applyDestChainConfigUpdates(destChainConfigArgs);
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  /// @notice Gets the next sequence number to be used in the onRamp
  /// @param destChainSelector The destination chain selector
  /// @return the next sequence number to be used
  function getExpectedNextSequenceNumber(uint64 destChainSelector) external view returns (uint64) {
    return s_destChainConfigs[destChainSelector].sequenceNumber + 1;
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function forwardFromRouter(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message,
    uint256 feeTokenAmount,
    address originalSender
  ) external returns (bytes32) {
    DestChainConfig storage destChainConfig = s_destChainConfigs[destChainSelector];

    // NOTE: assumes the message has already been validated through the getFee call
    // Validate message sender is set and allowed. Not validated in `getFee` since it is not user-driven.
    if (originalSender == address(0)) revert RouterMustSetOriginalSender();
    // Router address may be zero intentionally to pause.
    if (msg.sender != address(destChainConfig.router)) revert MustBeCalledByRouter();

    {
      // scoped to reduce stack usage
      address messageValidator = s_dynamicConfig.messageValidator;
      if (messageValidator != address(0)) {
        IMessageInterceptor(messageValidator).onOutboundMessage(destChainSelector, message);
      }
    }

    // Convert message fee to juels and retrieve converted args
    (uint256 msgFeeJuels, bool isOutOfOrderExecution, bytes memory convertedExtraArgs) = IPriceRegistry(
      s_dynamicConfig.priceRegistry
    ).processMessageArgs(destChainSelector, message.feeToken, feeTokenAmount, message.extraArgs);

    emit FeePaid(message.feeToken, msgFeeJuels);

    Internal.EVM2AnyRampMessage memory newMessage = Internal.EVM2AnyRampMessage({
      header: Internal.RampMessageHeader({
        // Should be generated after the message is complete
        messageId: "",
        sourceChainSelector: i_chainSelector,
        destChainSelector: destChainSelector,
        // We need the next available sequence number so we increment before we use the value
        sequenceNumber: ++destChainConfig.sequenceNumber,
        // Only bump nonce for messages that specify allowOutOfOrderExecution == false. Otherwise, we
        // may block ordered message nonces, which is not what we want.
        nonce: isOutOfOrderExecution
          ? 0
          : INonceManager(i_nonceManager).getIncrementedOutboundNonce(destChainSelector, originalSender)
      }),
      sender: originalSender,
      data: message.data,
      extraArgs: message.extraArgs,
      receiver: message.receiver,
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      // Should be populated via lock / burn pool calls
      tokenAmounts: new Internal.RampTokenAmount[](message.tokenAmounts.length)
    });

    // Lock the tokens as last step. TokenPools may not always be trusted.
    // There should be no state changes after external call to TokenPools.
    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      newMessage.tokenAmounts[i] =
        _lockOrBurnSingleToken(message.tokenAmounts[i], destChainSelector, message.receiver, originalSender);
    }

    // Validate pool return data after it is populated (view function - no state changes)
    IPriceRegistry(s_dynamicConfig.priceRegistry).validatePoolReturnData(
      destChainSelector, newMessage.tokenAmounts, message.tokenAmounts
    );

    // Override extraArgs with latest version
    newMessage.extraArgs = convertedExtraArgs;

    // Hash only after all fields have been set
    newMessage.header.messageId = Internal._hash(
      newMessage,
      // Metadata hash preimage to ensure global uniqueness, ensuring 2 identical messages sent to 2 different
      // lanes will have a distinct hash.
      keccak256(abi.encode(Internal.EVM_2_ANY_MESSAGE_HASH, i_chainSelector, destChainSelector, address(this)))
    );

    // Emit message request
    // This must happen after any pool events as some tokens (e.g. USDC) emit events that we expect to precede this
    // event in the offchain code.
    emit CCIPSendRequested(destChainSelector, newMessage);
    return newMessage.header.messageId;
  }

  /// @notice Uses a pool to lock or burn a token
  /// @param tokenAndAmount Token address and amount to lock or burn
  /// @param destChainSelector Target dest chain selector of the message
  /// @param receiver Message receiver
  /// @param originalSender Message sender
  /// @return rampTokenAndAmount Ramp token and amount data
  function _lockOrBurnSingleToken(
    Client.EVMTokenAmount memory tokenAndAmount,
    uint64 destChainSelector,
    bytes memory receiver,
    address originalSender
  ) internal returns (Internal.RampTokenAmount memory) {
    if (tokenAndAmount.amount == 0) revert CannotSendZeroTokens();

    IPoolV1 sourcePool = getPoolBySourceToken(destChainSelector, IERC20(tokenAndAmount.token));
    // We don't have to check if it supports the pool version in a non-reverting way here because
    // if we revert here, there is no effect on CCIP. Therefore we directly call the supportsInterface
    // function and not through the ERC165Checker.
    if (address(sourcePool) == address(0) || !sourcePool.supportsInterface(Pool.CCIP_POOL_V1)) {
      revert UnsupportedToken(tokenAndAmount.token);
    }

    Pool.LockOrBurnOutV1 memory poolReturnData = sourcePool.lockOrBurn(
      Pool.LockOrBurnInV1({
        receiver: receiver,
        remoteChainSelector: destChainSelector,
        originalSender: originalSender,
        amount: tokenAndAmount.amount,
        localToken: tokenAndAmount.token
      })
    );

    // NOTE: pool data validations are outsourced to the PriceRegistry to handle family-specific logic handling

    return Internal.RampTokenAmount({
      sourcePoolAddress: abi.encode(sourcePool),
      destTokenAddress: poolReturnData.destTokenAddress,
      extraData: poolReturnData.destPoolData,
      amount: tokenAndAmount.amount
    });
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static onRamp config.
  /// @dev RMN depends on this function, if changing, please notify the RMN maintainers.
  /// @return the configuration.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({
      chainSelector: i_chainSelector,
      rmnProxy: i_rmnProxy,
      nonceManager: i_nonceManager,
      tokenAdminRegistry: i_tokenAdminRegistry
    });
  }

  /// @notice Returns the dynamic onRamp config.
  /// @return dynamicConfig the configuration.
  function getDynamicConfig() external view returns (DynamicConfig memory dynamicConfig) {
    return s_dynamicConfig;
  }

  /// @notice Sets the dynamic configuration.
  /// @param dynamicConfig The configuration.
  function setDynamicConfig(DynamicConfig memory dynamicConfig) external onlyOwner {
    _setDynamicConfig(dynamicConfig);
  }

  /// @notice Gets the source router for a destination chain
  /// @param destChainSelector The destination chain selector
  /// @return router the router for the provided destination chain
  function getRouter(uint64 destChainSelector) external view returns (IRouter) {
    return s_destChainConfigs[destChainSelector].router;
  }

  /// @notice Internal version of setDynamicConfig to allow for reuse in the constructor.
  function _setDynamicConfig(DynamicConfig memory dynamicConfig) internal {
    if (dynamicConfig.priceRegistry == address(0) || dynamicConfig.feeAggregator == address(0)) revert InvalidConfig();

    s_dynamicConfig = dynamicConfig;

    emit ConfigSet(
      StaticConfig({
        chainSelector: i_chainSelector,
        rmnProxy: i_rmnProxy,
        nonceManager: i_nonceManager,
        tokenAdminRegistry: i_tokenAdminRegistry
      }),
      dynamicConfig
    );
  }

  /// @notice Updates the destination chain specific config.
  /// @param destChainConfigArgs Array of source chain specific configs.
  function applyDestChainConfigUpdates(DestChainConfigArgs[] memory destChainConfigArgs) external onlyOwner {
    _applyDestChainConfigUpdates(destChainConfigArgs);
  }

  /// @notice Internal version of applyDestChainConfigUpdates.
  function _applyDestChainConfigUpdates(DestChainConfigArgs[] memory destChainConfigArgs) internal {
    for (uint256 i = 0; i < destChainConfigArgs.length; ++i) {
      DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[i];
      uint64 destChainSelector = destChainConfigArgs[i].destChainSelector;

      if (destChainSelector == 0) {
        revert InvalidDestChainConfig(destChainSelector);
      }

      DestChainConfig memory newDestChainConfig = DestChainConfig({
        sequenceNumber: s_destChainConfigs[destChainSelector].sequenceNumber,
        router: destChainConfigArg.router
      });
      s_destChainConfigs[destChainSelector] = newDestChainConfig;

      emit DestChainConfigSet(destChainSelector, newDestChainConfig);
    }
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @inheritdoc IEVM2AnyOnRampClient
  function getPoolBySourceToken(uint64, /*destChainSelector*/ IERC20 sourceToken) public view returns (IPoolV1) {
    return IPoolV1(ITokenAdminRegistry(i_tokenAdminRegistry).getPool(address(sourceToken)));
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function getSupportedTokens(uint64 /*destChainSelector*/ ) external pure returns (address[] memory) {
    revert GetSupportedTokensFunctionalityRemovedCheckAdminRegistry();
  }

  // ================================================================
  // │                             Fees                             │
  // ================================================================

  /// @inheritdoc IEVM2AnyOnRampClient
  /// @dev getFee MUST revert if the feeToken is not listed in the fee token config, as the router assumes it does.
  /// @param destChainSelector The destination chain selector.
  /// @param message The message to get quote for.
  /// @return feeTokenAmount The amount of fee token needed for the fee, in smallest denomination of the fee token.
  function getFee(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 feeTokenAmount) {
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(destChainSelector)))) revert CursedByRMN(destChainSelector);

    return IPriceRegistry(s_dynamicConfig.priceRegistry).getValidatedFee(destChainSelector, message);
  }

  /// @notice Withdraws the outstanding fee token balances to the fee aggregator.
  /// @dev This function can be permissionless as it only transfers accepted fee tokens to the fee aggregator which is a trusted address.
  function withdrawFeeTokens() external {
    address[] memory feeTokens = IPriceRegistry(s_dynamicConfig.priceRegistry).getFeeTokens();
    address feeAggregator = s_dynamicConfig.feeAggregator;

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      IERC20 feeToken = IERC20(feeTokens[i]);
      uint256 feeTokenBalance = feeToken.balanceOf(address(this));

      if (feeTokenBalance > 0) {
        feeToken.safeTransfer(feeAggregator, feeTokenBalance);

        emit FeeTokenWithdrawn(feeAggregator, address(feeToken), feeTokenBalance);
      }
    }
  }
}
