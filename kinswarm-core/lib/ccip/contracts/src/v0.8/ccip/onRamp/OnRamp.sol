// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IEVM2AnyOnRampClient} from "../interfaces/IEVM2AnyOnRampClient.sol";
import {IFeeQuoter} from "../interfaces/IFeeQuoter.sol";
import {IMessageInterceptor} from "../interfaces/IMessageInterceptor.sol";
import {INonceManager} from "../interfaces/INonceManager.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IRMNRemote} from "../interfaces/IRMNRemote.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../libraries/USDPriceWith18Decimals.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice The OnRamp is a contract that handles lane-specific fee logic
/// @dev The OnRamp and OffRamp form an xchain upgradeable unit. Any change to one of them
/// results in an onchain upgrade of all 3.
contract OnRamp is IEVM2AnyOnRampClient, ITypeAndVersion, OwnerIsCreator {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.AddressSet;
  using USDPriceWith18Decimals for uint224;

  error CannotSendZeroTokens();
  error UnsupportedToken(address token);
  error MustBeCalledByRouter();
  error RouterMustSetOriginalSender();
  error InvalidConfig();
  error CursedByRMN(uint64 sourceChainSelector);
  error GetSupportedTokensFunctionalityRemovedCheckAdminRegistry();
  error InvalidDestChainConfig(uint64 sourceChainSelector);
  error OnlyCallableByOwnerOrAllowlistAdmin();
  error SenderNotAllowed(address sender);
  error InvalidAllowListRequest(uint64 destChainSelector);
  error ReentrancyGuardReentrantCall();

  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event DestChainConfigSet(
    uint64 indexed destChainSelector, uint64 sequenceNumber, IRouter router, bool allowListEnabled
  );
  event FeeTokenWithdrawn(address indexed feeAggregator, address indexed feeToken, uint256 amount);
  /// RMN depends on this event, if changing, please notify the RMN maintainers.
  event CCIPMessageSent(
    uint64 indexed destChainSelector, uint64 indexed sequenceNumber, Internal.EVM2AnyRampMessage message
  );
  event AllowListAdminSet(address indexed allowListAdmin);
  event AllowListSendersAdded(uint64 indexed destChainSelector, address[] senders);
  event AllowListSendersRemoved(uint64 indexed destChainSelector, address[] senders);

  /// @dev Struct that contains the static configuration
  /// RMN depends on this struct, if changing, please notify the RMN maintainers.
  // solhint-disable-next-line gas-struct-packing
  struct StaticConfig {
    uint64 chainSelector; // ─────╮ Source chain selector
    IRMNRemote rmnRemote; // ─────╯ RMN remote address
    address nonceManager; // Nonce manager address
    address tokenAdminRegistry; // Token admin registry address
  }

  /// @dev Struct that contains the dynamic configuration
  // solhint-disable-next-line gas-struct-packing
  struct DynamicConfig {
    address feeQuoter; // FeeQuoter address
    bool reentrancyGuardEntered; // Reentrancy protection
    address messageInterceptor; // Optional message interceptor to validate outbound messages (zero address = no interceptor)
    address feeAggregator; // Fee aggregator address
    address allowListAdmin; // authorized admin to add or remove allowed senders
  }

  /// @dev Struct to hold the configs for a destination chain
  /// @dev sequenceNumber, allowListEnabled, router will all be packed in 1 slot
  struct DestChainConfig {
    // The last used sequence number. This is zero in the case where no messages have yet been sent.
    // 0 is not a valid sequence number for any real transaction.
    uint64 sequenceNumber; // ──────╮ The last used sequence number
    bool allowListEnabled; //       │ boolean indicator to specify if allowList check is enabled
    IRouter router; // ─────────────╯ Local router address  that is allowed to send messages to the destination chain.
    // This is the list of addresses allowed to send messages from onRamp
    EnumerableSet.AddressSet allowedSendersList;
  }

  /// @dev Same as DestChainConfig but with the destChainSelector so that an array of these
  /// can be passed in the constructor and the applyDestChainConfigUpdates function
  //solhint-disable gas-struct-packing
  struct DestChainConfigArgs {
    uint64 destChainSelector; // ─╮ Destination chain selector
    IRouter router; //            │ Source router address
    bool allowListEnabled; //─────╯ Boolean indicator to specify if allowList check is enabled
  }

  /// @dev Struct used to apply AllowList Senders for multiple destChainSelectors
  /// @dev the senders in the AllowlistedSenders here is the user that sends the message
  /// @dev the config restricts the chain to allow only allowedList of senders to send message from this chain to a destChainSelector
  /// @dev destChainSelector, allowListEnabled will be packed in 1 slot
  //solhint-disable gas-struct-packing
  struct AllowListConfigArgs {
    uint64 destChainSelector; // ─────────────╮ Destination chain selector
    //                                        │ destChainSelector and allowListEnabled are packed in the same slot
    bool allowListEnabled; // ────────────────╯ boolean indicator to specify if allowList check is enabled.
    address[] addedAllowlistedSenders; // list of senders to be added to the allowedSendersList
    address[] removedAllowlistedSenders; // list of senders to be removed from the allowedSendersList
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "OnRamp 1.6.0-dev";
  /// @dev The chain ID of the source chain that this contract is deployed to
  uint64 private immutable i_chainSelector;
  /// @dev The rmn contract
  IRMNRemote private immutable i_rmnRemote;
  /// @dev The address of the nonce manager
  address private immutable i_nonceManager;
  /// @dev The address of the token admin registry
  address private immutable i_tokenAdminRegistry;

  // DYNAMIC CONFIG
  /// @dev The dynamic config for the onRamp
  DynamicConfig private s_dynamicConfig;

  /// @dev The destination chain specific configs
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

  /// @notice Gets the next sequence number to be used in the onRamp
  /// @param destChainSelector The destination chain selector
  /// @return nextSequenceNumber The next sequence number to be used
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
    // We rely on a reentrancy guard here due to the untrusted calls performed to the pools
    // This enables some optimizations by not following the CEI pattern
    if (s_dynamicConfig.reentrancyGuardEntered) revert ReentrancyGuardReentrantCall();

    s_dynamicConfig.reentrancyGuardEntered = true;

    DestChainConfig storage destChainConfig = s_destChainConfigs[destChainSelector];

    // NOTE: assumes the message has already been validated through the getFee call
    // Validate message sender is set and allowed. Not validated in `getFee` since it is not user-driven.
    if (originalSender == address(0)) revert RouterMustSetOriginalSender();

    if (destChainConfig.allowListEnabled) {
      if (!destChainConfig.allowedSendersList.contains(originalSender)) {
        revert SenderNotAllowed(originalSender);
      }
    }

    // Router address may be zero intentionally to pause.
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
        // Should be generated after the message is complete
        messageId: "",
        sourceChainSelector: i_chainSelector,
        destChainSelector: destChainSelector,
        // We need the next available sequence number so we increment before we use the value
        sequenceNumber: ++destChainConfig.sequenceNumber,
        // Only bump nonce for messages that specify allowOutOfOrderExecution == false. Otherwise, we
        // may block ordered message nonces, which is not what we want.
        nonce: 0
      }),
      sender: originalSender,
      data: message.data,
      extraArgs: "",
      receiver: message.receiver,
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      feeValueJuels: 0, // calculated later
      // Should be populated via lock / burn pool calls
      tokenAmounts: new Internal.EVM2AnyTokenTransfer[](message.tokenAmounts.length)
    });

    // Lock / burn the tokens as last step. TokenPools may not always be trusted.
    Client.EVMTokenAmount[] memory tokenAmounts = message.tokenAmounts;
    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      newMessage.tokenAmounts[i] =
        _lockOrBurnSingleToken(tokenAmounts[i], destChainSelector, message.receiver, originalSender);
    }

    // Convert message fee to juels and retrieve converted args
    // Validate pool return data after it is populated (view function - no state changes)
    bool isOutOfOrderExecution;
    bytes memory convertedExtraArgs;
    bytes[] memory destExecDataPerToken;
    (newMessage.feeValueJuels, isOutOfOrderExecution, convertedExtraArgs, destExecDataPerToken) = IFeeQuoter(
      s_dynamicConfig.feeQuoter
    ).processMessageArgs(
      destChainSelector, message.feeToken, feeTokenAmount, message.extraArgs, newMessage.tokenAmounts, tokenAmounts
    );

    newMessage.header.nonce = isOutOfOrderExecution
      ? 0
      : INonceManager(i_nonceManager).getIncrementedOutboundNonce(destChainSelector, originalSender);
    newMessage.extraArgs = convertedExtraArgs;

    for (uint256 i = 0; i < newMessage.tokenAmounts.length; ++i) {
      newMessage.tokenAmounts[i].destExecData = destExecDataPerToken[i];
    }

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
    emit CCIPMessageSent(destChainSelector, newMessage.header.sequenceNumber, newMessage);

    s_dynamicConfig.reentrancyGuardEntered = false;

    return newMessage.header.messageId;
  }

  /// @notice Uses a pool to lock or burn a token
  /// @param tokenAndAmount Token address and amount to lock or burn
  /// @param destChainSelector Target destination chain selector of the message
  /// @param receiver Message receiver
  /// @param originalSender Message sender
  /// @return evm2AnyTokenTransfer EVM2Any token and amount data
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

    // NOTE: pool data validations are outsourced to the FeeQuoter to handle family-specific logic handling
    return Internal.EVM2AnyTokenTransfer({
      sourcePoolAddress: address(sourcePool),
      destTokenAddress: poolReturnData.destTokenAddress,
      extraData: poolReturnData.destPoolData,
      amount: tokenAndAmount.amount,
      destExecData: "" // This is set in the processPoolReturnData function
    });
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static onRamp config.
  /// @dev RMN depends on this function, if modified, please notify the RMN maintainers.
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

  /// @notice Gets the source router for a destination chain
  /// @param destChainSelector The destination chain selector
  /// @return router the router for the provided destination chain
  function getRouter(
    uint64 destChainSelector
  ) external view returns (IRouter) {
    return s_destChainConfigs[destChainSelector].router;
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
      destChainConfig.router = destChainConfigArg.router;
      destChainConfig.allowListEnabled = destChainConfigArg.allowListEnabled;

      emit DestChainConfigSet(
        destChainSelector, destChainConfig.sequenceNumber, destChainConfigArg.router, destChainConfig.allowListEnabled
      );
    }
  }

  /// @notice get ChainConfig configured for the DestinationChainSelector
  /// @param destChainSelector The destination chain selector
  /// @return sequenceNumber The last used sequence number
  /// @return allowListEnabled boolean indicator to specify if allowList check is enabled
  /// @return router address of the router
  function getDestChainConfig(
    uint64 destChainSelector
  ) public view returns (uint64 sequenceNumber, bool allowListEnabled, address router) {
    DestChainConfig storage config = s_destChainConfigs[destChainSelector];
    sequenceNumber = config.sequenceNumber;
    allowListEnabled = config.allowListEnabled;
    router = address(config.router);
    return (sequenceNumber, allowListEnabled, router);
  }

  /// @notice get allowedSenders List configured for the DestinationChainSelector
  /// @param destChainSelector The destination chain selector
  /// @return array of allowedList of Senders
  function getAllowedSendersList(
    uint64 destChainSelector
  ) public view returns (address[] memory) {
    return s_destChainConfigs[destChainSelector].allowedSendersList.values();
  }

  // ================================================================
  // │                          Allowlist                           │
  // ================================================================

  /// @notice Updates allowListConfig for Senders
  /// @dev configuration used to set the list of senders who are authorized to send messages
  /// @param allowListConfigArgsItems Array of AllowListConfigArguments where each item is for a destChainSelector
  function applyAllowListUpdates(
    AllowListConfigArgs[] calldata allowListConfigArgsItems
  ) external {
    if (msg.sender != owner()) {
      if (msg.sender != s_dynamicConfig.allowListAdmin) {
        revert OnlyCallableByOwnerOrAllowlistAdmin();
      }
    }

    for (uint256 i = 0; i < allowListConfigArgsItems.length; ++i) {
      AllowListConfigArgs memory allowListConfigArgs = allowListConfigArgsItems[i];

      DestChainConfig storage destChainConfig = s_destChainConfigs[allowListConfigArgs.destChainSelector];
      destChainConfig.allowListEnabled = allowListConfigArgs.allowListEnabled;

      if (allowListConfigArgs.addedAllowlistedSenders.length > 0) {
        if (allowListConfigArgs.allowListEnabled) {
          for (uint256 j = 0; j < allowListConfigArgs.addedAllowlistedSenders.length; ++j) {
            address toAdd = allowListConfigArgs.addedAllowlistedSenders[j];
            if (toAdd == address(0)) {
              revert InvalidAllowListRequest(allowListConfigArgs.destChainSelector);
            }
            destChainConfig.allowedSendersList.add(toAdd);
          }

          emit AllowListSendersAdded(allowListConfigArgs.destChainSelector, allowListConfigArgs.addedAllowlistedSenders);
        } else {
          revert InvalidAllowListRequest(allowListConfigArgs.destChainSelector);
        }
      }

      for (uint256 j = 0; j < allowListConfigArgs.removedAllowlistedSenders.length; ++j) {
        destChainConfig.allowedSendersList.remove(allowListConfigArgs.removedAllowlistedSenders[j]);
      }

      if (allowListConfigArgs.removedAllowlistedSenders.length > 0) {
        emit AllowListSendersRemoved(
          allowListConfigArgs.destChainSelector, allowListConfigArgs.removedAllowlistedSenders
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
    uint64 /*destChainSelector*/
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
  /// @dev This function can be permissionless as it only transfers accepted fee tokens to the fee aggregator which is a trusted address.
  function withdrawFeeTokens() external {
    address[] memory feeTokens = IFeeQuoter(s_dynamicConfig.feeQuoter).getFeeTokens();
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
