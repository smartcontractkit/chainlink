// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IEVM2AnyOnRampClient} from "../interfaces/IEVM2AnyOnRampClient.sol";
import {IFeeQuoter} from "../interfaces/IFeeQuoter.sol";
import {IMessageInterceptor} from "../interfaces/IMessageInterceptor.sol";
import {INonceManager} from "../interfaces/INonceManager.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IRMNRemote} from "../interfaces/IRMNRemote.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../libraries/USDPriceWith18Decimals.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice The OnRamp is a contract that handles lane-specific fee logic.
/// @dev The OnRamp and OffRamp form a cross chain upgradeable unit. Any change to one of them results in an onchain
/// upgrade of both contracts.
contract OnRamp is IEVM2AnyOnRampClient, ITypeAndVersion, Ownable2StepMsgSender {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.AddressSet;
  using USDPriceWith18Decimals for uint224;

  error CannotSendZeroTokens();
  error UnsupportedToken(address token);
  error MustBeCalledByRouter();
  error RouterMustSetOriginalSender();
  error InvalidConfig();
  error CursedByRMN(uint64 destChainSelector);
  error GetSupportedTokensFunctionalityRemovedCheckAdminRegistry();
  error InvalidDestChainConfig(uint64 destChainSelector);
  error OnlyCallableByOwnerOrAllowlistAdmin();
  error SenderNotAllowed(address sender);
  error InvalidAllowListRequest(uint64 destChainSelector);
  error ReentrancyGuardReentrantCall();

  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event DestChainConfigSet(
    uint64 indexed destChainSelector, uint64 sequenceNumber, IRouter router, bool allowlistEnabled
  );
  event FeeTokenWithdrawn(address indexed feeAggregator, address indexed feeToken, uint256 amount);

  /// RMN depends on this event, if changing, please notify the RMN maintainers.
  event CCIPMessageSent(
    uint64 indexed destChainSelector, uint64 indexed sequenceNumber, Internal.EVM2AnyRampMessage message
  );
  event AllowListAdminSet(address indexed allowlistAdmin);
  event AllowListSendersAdded(uint64 indexed destChainSelector, address[] senders);
  event AllowListSendersRemoved(uint64 indexed destChainSelector, address[] senders);

  /// @dev Struct that contains the static configuration.
  // solhint-disable-next-line gas-struct-packing
  struct StaticConfig {
    uint64 chainSelector; // ────╮ Source chain selector.
    IRMNRemote rmnRemote; // ────╯ RMN remote address.
    address nonceManager; //       Nonce manager address.
    address tokenAdminRegistry; // Token admin registry address.
  }

  /// @dev Struct that contains the dynamic configuration
  // solhint-disable-next-line gas-struct-packing
  struct DynamicConfig {
    address feeQuoter; // FeeQuoter address.
    bool reentrancyGuardEntered; // Reentrancy protection.
    address messageInterceptor; // Optional message interceptor to validate messages. Zero address = no interceptor.
    address feeAggregator; // Fee aggregator address.
    address allowlistAdmin; // authorized admin to add or remove allowed senders.
  }

  /// @dev Struct to hold the configs for a single destination chain.
  struct DestChainConfig {
    // The last used sequence number. This is zero in the case where no messages have yet been sent.
    // 0 is not a valid sequence number for any real transaction as this value will be incremented before use.
    uint64 sequenceNumber; // ──╮ The last used sequence number.
    bool allowlistEnabled; //   │ True if the allowlist is enabled.
    IRouter router; // ─────────╯ Local router address  that is allowed to send messages to the destination chain.
    EnumerableSet.AddressSet allowedSendersList; // The list of addresses allowed to send messages.
  }

  /// @dev Same as DestChainConfig but with the destChainSelector so that an array of these can be passed in the
  /// constructor and the applyDestChainConfigUpdates function.
  // solhint-disable gas-struct-packing
  struct DestChainConfigArgs {
    uint64 destChainSelector; // ─╮ Destination chain selector.
    IRouter router; //            │ Source router address.
    bool allowlistEnabled; // ────╯ True if the allowlist is enabled.
  }

  /// @dev Struct to hold the allowlist configuration args per dest chain.
  struct AllowlistConfigArgs {
    uint64 destChainSelector; // ──╮ Destination chain selector.
    bool allowlistEnabled; // ─────╯ True if the allowlist is enabled.
    address[] addedAllowlistedSenders; // list of senders to be added to the allowedSendersList.
    address[] removedAllowlistedSenders; // list of senders to be removed from the allowedSendersList.
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "OnRamp 1.6.0";
  /// @dev The chain ID of the source chain that this contract is deployed to.
  uint64 private immutable i_chainSelector;
  /// @dev The rmn contract.
  IRMNRemote private immutable i_rmnRemote;
  /// @dev The address of the nonce manager.
  address private immutable i_nonceManager;
  /// @dev The address of the token admin registry.
  address private immutable i_tokenAdminRegistry;

  // DYNAMIC CONFIG
  /// @dev The dynamic config for the onRamp.
  DynamicConfig private s_dynamicConfig;

  /// @dev The destination chain specific configs.
  mapping(uint64 destChainSelector => DestChainConfig destChainConfig) private s_destChainConfigs;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigArgs
  ) {
    if (
      staticConfig.chainSelector == 0 || address(staticConfig.rmnRemote) == address(0)
        || staticConfig.nonceManager == address(0) || staticConfig.tokenAdminRegistry == address(0)
    ) {
      revert InvalidConfig();
    }

    i_chainSelector = staticConfig.chainSelector;
    i_rmnRemote = staticConfig.rmnRemote;
    i_nonceManager = staticConfig.nonceManager;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;

    _setDynamicConfig(dynamicConfig);
    _applyDestChainConfigUpdates(destChainConfigArgs);
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  /// @notice Gets the next sequence number to be used in the onRamp.
  /// @param destChainSelector The destination chain selector.
  /// @return nextSequenceNumber The next sequence number to be used.
  function getExpectedNextSequenceNumber(
    uint64 destChainSelector
  ) external view returns (uint64) {
    return s_destChainConfigs[destChainSelector].sequenceNumber + 1;
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function forwardFromRouter(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message,
    uint256 feeTokenAmount,
    address originalSender
  ) external returns (bytes32) {
    // We rely on a reentrancy guard here due to the untrusted calls performed to the pools. This enables some
    // optimizations by not following the CEI pattern.
    if (s_dynamicConfig.reentrancyGuardEntered) revert ReentrancyGuardReentrantCall();
    s_dynamicConfig.reentrancyGuardEntered = true;

    DestChainConfig storage destChainConfig = s_destChainConfigs[destChainSelector];

    // NOTE: assumes the message has already been validated through the getFee call.
    // Validate originalSender is set and allowed. Not validated in `getFee` since it is not user-driven.
    if (originalSender == address(0)) revert RouterMustSetOriginalSender();

    if (destChainConfig.allowlistEnabled) {
      if (!destChainConfig.allowedSendersList.contains(originalSender)) {
        revert SenderNotAllowed(originalSender);
      }
    }

    // Router address may be zero intentionally to pause, which should stop all messages.
    if (msg.sender != address(destChainConfig.router)) revert MustBeCalledByRouter();

    {
      // scoped to reduce stack usage
      address messageInterceptor = s_dynamicConfig.messageInterceptor;
      if (messageInterceptor != address(0)) {
        IMessageInterceptor(messageInterceptor).onOutboundMessage(destChainSelector, message);
      }
    }

    Internal.EVM2AnyRampMessage memory newMessage = Internal.EVM2AnyRampMessage({
      header: Internal.RampMessageHeader({
        // Should be generated after the message is complete.
        messageId: "",
        sourceChainSelector: i_chainSelector,
        destChainSelector: destChainSelector,
        // We need the next available sequence number so we increment before we use the value.
        sequenceNumber: ++destChainConfig.sequenceNumber,
        // Only bump nonce for messages that specify allowOutOfOrderExecution == false. Otherwise, we may block ordered
        // message nonces, which is not what we want.
        nonce: 0
      }),
      sender: originalSender,
      data: message.data,
      extraArgs: "",
      receiver: message.receiver,
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      feeValueJuels: 0, // calculated later.
      // Should be populated via lock / burn pool calls.
      tokenAmounts: new Internal.EVM2AnyTokenTransfer[](message.tokenAmounts.length)
    });

    // Convert message fee to juels and retrieve converted args.
    // Validate pool return data after it is populated (view function - no state changes).
    bool isOutOfOrderExecution;
    bytes memory tokenReceiver;
    (newMessage.feeValueJuels, isOutOfOrderExecution, newMessage.extraArgs, tokenReceiver) = IFeeQuoter(
      s_dynamicConfig.feeQuoter
    ).processMessageArgs(destChainSelector, message.feeToken, feeTokenAmount, message.extraArgs, message.receiver);

    Client.EVMTokenAmount[] memory tokenAmounts = message.tokenAmounts;
    // Lock / burn the tokens as last step. TokenPools may not always be trusted.
    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      newMessage.tokenAmounts[i] =
        _lockOrBurnSingleToken(tokenAmounts[i], destChainSelector, tokenReceiver, originalSender);
    }

    bytes[] memory destExecDataPerToken = IFeeQuoter(s_dynamicConfig.feeQuoter).processPoolReturnData(
      destChainSelector, newMessage.tokenAmounts, tokenAmounts
    );
    newMessage.header.nonce = isOutOfOrderExecution
      ? 0
      : INonceManager(i_nonceManager).getIncrementedOutboundNonce(destChainSelector, originalSender);

    for (uint256 i = 0; i < newMessage.tokenAmounts.length; ++i) {
      newMessage.tokenAmounts[i].destExecData = destExecDataPerToken[i];
    }

    newMessage = _postProcessMessage(newMessage);

    // Hash only after all fields have been set.
    newMessage.header.messageId = Internal._hash(
      newMessage,
      // Metadata hash preimage to ensure global uniqueness, ensuring 2 identical messages sent to 2 different lanes
      // will have a distinct hash.
      keccak256(abi.encode(Internal.EVM_2_ANY_MESSAGE_HASH, i_chainSelector, destChainSelector, address(this)))
    );

    // Emit message request.
    // This must happen after any pool events as some tokens (e.g. USDC) emit events that we expect to precede this
    // event in the offchain code.
    emit CCIPMessageSent(destChainSelector, newMessage.header.sequenceNumber, newMessage);

    s_dynamicConfig.reentrancyGuardEntered = false;

    return newMessage.header.messageId;
  }

  /// @notice hook for applying custom logic to the message from the router
  /// @param message router message
  /// @return transformedMessage modified message
  function _postProcessMessage(
    Internal.EVM2AnyRampMessage memory message
  ) internal virtual returns (Internal.EVM2AnyRampMessage memory transformedMessage) {
    return message;
  }

  /// @notice Uses a pool to lock or burn a token.
  /// @param tokenAndAmount Token address and amount to lock or burn.
  /// @param destChainSelector Target destination chain selector of the message.
  /// @param receiver Message receiver.
  /// @param originalSender Message sender.
  /// @return evm2AnyTokenTransfer EVM2Any token and amount data.
  function _lockOrBurnSingleToken(
    Client.EVMTokenAmount memory tokenAndAmount,
    uint64 destChainSelector,
    bytes memory receiver,
    address originalSender
  ) internal returns (Internal.EVM2AnyTokenTransfer memory) {
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

    // NOTE: pool data validations are outsourced to the FeeQuoter to handle family-specific logic handling.
    return Internal.EVM2AnyTokenTransfer({
      sourcePoolAddress: address(sourcePool),
      destTokenAddress: poolReturnData.destTokenAddress,
      extraData: poolReturnData.destPoolData,
      amount: tokenAndAmount.amount,
      destExecData: "" // This is set in the processPoolReturnData function.
    });
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static onRamp config.
  /// @return staticConfig the static configuration.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({
      chainSelector: i_chainSelector,
      rmnRemote: i_rmnRemote,
      nonceManager: i_nonceManager,
      tokenAdminRegistry: i_tokenAdminRegistry
    });
  }

  /// @notice Returns the dynamic onRamp config.
  /// @return dynamicConfig the dynamic configuration.
  function getDynamicConfig() external view returns (DynamicConfig memory dynamicConfig) {
    return s_dynamicConfig;
  }

  /// @notice Sets the dynamic configuration.
  /// @param dynamicConfig The configuration.
  function setDynamicConfig(
    DynamicConfig memory dynamicConfig
  ) external onlyOwner {
    _setDynamicConfig(dynamicConfig);
  }

  /// @notice Internal version of setDynamicConfig to allow for reuse in the constructor.
  function _setDynamicConfig(
    DynamicConfig memory dynamicConfig
  ) internal {
    if (
      dynamicConfig.feeQuoter == address(0) || dynamicConfig.feeAggregator == address(0)
        || dynamicConfig.reentrancyGuardEntered
    ) revert InvalidConfig();

    s_dynamicConfig = dynamicConfig;

    emit ConfigSet(
      StaticConfig({
        chainSelector: i_chainSelector,
        rmnRemote: i_rmnRemote,
        nonceManager: i_nonceManager,
        tokenAdminRegistry: i_tokenAdminRegistry
      }),
      dynamicConfig
    );
  }

  /// @notice Updates destination chains specific configs.
  /// @param destChainConfigArgs Array of destination chain specific configs.
  function applyDestChainConfigUpdates(
    DestChainConfigArgs[] memory destChainConfigArgs
  ) external onlyOwner {
    _applyDestChainConfigUpdates(destChainConfigArgs);
  }

  /// @notice Internal version of applyDestChainConfigUpdates.
  function _applyDestChainConfigUpdates(
    DestChainConfigArgs[] memory destChainConfigArgs
  ) internal {
    for (uint256 i = 0; i < destChainConfigArgs.length; ++i) {
      DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[i];
      uint64 destChainSelector = destChainConfigArgs[i].destChainSelector;

      if (destChainSelector == 0) {
        revert InvalidDestChainConfig(destChainSelector);
      }

      DestChainConfig storage destChainConfig = s_destChainConfigs[destChainSelector];
      // The router can be zero to pause the destination chain
      destChainConfig.router = destChainConfigArg.router;
      destChainConfig.allowlistEnabled = destChainConfigArg.allowlistEnabled;

      emit DestChainConfigSet(
        destChainSelector, destChainConfig.sequenceNumber, destChainConfigArg.router, destChainConfig.allowlistEnabled
      );
    }
  }

  /// @notice get ChainConfig configured for the DestinationChainSelector.
  /// @param destChainSelector The destination chain selector.
  /// @return sequenceNumber The last used sequence number.
  /// @return allowlistEnabled boolean indicator to specify if allowlist check is enabled.
  /// @return router address of the router.
  function getDestChainConfig(
    uint64 destChainSelector
  ) external view returns (uint64 sequenceNumber, bool allowlistEnabled, address router) {
    DestChainConfig storage config = s_destChainConfigs[destChainSelector];
    sequenceNumber = config.sequenceNumber;
    allowlistEnabled = config.allowlistEnabled;
    router = address(config.router);
    return (sequenceNumber, allowlistEnabled, router);
  }

  /// @notice get allowedSenders List configured for the DestinationChainSelector.
  /// @param destChainSelector The destination chain selector.
  /// @return isEnabled True if allowlist is enabled.
  /// @return configuredAddresses This is always populated with the list of allowed senders, even if the allowlist
  /// is turned off. This is because the only way to know what addresses are configured is through this function. If
  /// it would return an empty list when the allowlist is disabled, it would be impossible to know what addresses are
  /// configured.
  function getAllowedSendersList(
    uint64 destChainSelector
  ) external view returns (bool isEnabled, address[] memory configuredAddresses) {
    return (
      s_destChainConfigs[destChainSelector].allowlistEnabled,
      s_destChainConfigs[destChainSelector].allowedSendersList.values()
    );
  }

  // ================================================================
  // │                          Allowlist                           │
  // ================================================================

  /// @notice Updates allowlistConfig for Senders.
  /// @dev configuration used to set the list of senders who are authorized to send messages.
  /// @param allowlistConfigArgsItems Array of AllowlistConfigArguments where each item is for a destChainSelector.
  function applyAllowlistUpdates(
    AllowlistConfigArgs[] calldata allowlistConfigArgsItems
  ) external {
    if (msg.sender != owner()) {
      if (msg.sender != s_dynamicConfig.allowlistAdmin) {
        revert OnlyCallableByOwnerOrAllowlistAdmin();
      }
    }

    for (uint256 i = 0; i < allowlistConfigArgsItems.length; ++i) {
      AllowlistConfigArgs memory allowlistConfigArgs = allowlistConfigArgsItems[i];

      DestChainConfig storage destChainConfig = s_destChainConfigs[allowlistConfigArgs.destChainSelector];
      destChainConfig.allowlistEnabled = allowlistConfigArgs.allowlistEnabled;

      if (allowlistConfigArgs.addedAllowlistedSenders.length > 0) {
        if (allowlistConfigArgs.allowlistEnabled) {
          for (uint256 j = 0; j < allowlistConfigArgs.addedAllowlistedSenders.length; ++j) {
            address toAdd = allowlistConfigArgs.addedAllowlistedSenders[j];
            if (toAdd == address(0)) {
              revert InvalidAllowListRequest(allowlistConfigArgs.destChainSelector);
            }
            destChainConfig.allowedSendersList.add(toAdd);
          }

          emit AllowListSendersAdded(allowlistConfigArgs.destChainSelector, allowlistConfigArgs.addedAllowlistedSenders);
        } else {
          revert InvalidAllowListRequest(allowlistConfigArgs.destChainSelector);
        }
      }

      for (uint256 j = 0; j < allowlistConfigArgs.removedAllowlistedSenders.length; ++j) {
        destChainConfig.allowedSendersList.remove(allowlistConfigArgs.removedAllowlistedSenders[j]);
      }

      if (allowlistConfigArgs.removedAllowlistedSenders.length > 0) {
        emit AllowListSendersRemoved(
          allowlistConfigArgs.destChainSelector, allowlistConfigArgs.removedAllowlistedSenders
        );
      }
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
  function getSupportedTokens(
    uint64 // destChainSelector
  ) external pure returns (address[] memory) {
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
    if (i_rmnRemote.isCursed(bytes16(uint128(destChainSelector)))) revert CursedByRMN(destChainSelector);

    return IFeeQuoter(s_dynamicConfig.feeQuoter).getValidatedFee(destChainSelector, message);
  }

  /// @notice Withdraws the outstanding fee token balances to the fee aggregator.
  /// @param feeTokens The fee tokens to withdraw.
  /// @dev This function can be permissionless as it only transfers tokens to the fee aggregator which is a trusted address.
  function withdrawFeeTokens(
    address[] calldata feeTokens
  ) external {
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
