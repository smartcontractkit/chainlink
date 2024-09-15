// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IEVM2AnyOnRamp} from "../interfaces/IEVM2AnyOnRamp.sol";
import {IEVM2AnyOnRampClient} from "../interfaces/IEVM2AnyOnRampClient.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";
import {ILinkAvailable} from "../interfaces/automation/ILinkAvailable.sol";

import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "../libraries/USDPriceWith18Decimals.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableMap} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableMap.sol";

/// @notice The onRamp is a contract that handles lane-specific fee logic, NOP payments and
/// bridgeable token support.
/// @dev The EVM2EVMOnRamp, CommitStore and EVM2EVMOffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
contract EVM2EVMOnRamp is IEVM2AnyOnRamp, ILinkAvailable, AggregateRateLimiter, ITypeAndVersion {
  using SafeERC20 for IERC20;
  using EnumerableMap for EnumerableMap.AddressToUintMap;
  using USDPriceWith18Decimals for uint224;

  error InvalidExtraArgsTag();
  error ExtraArgOutOfOrderExecutionMustBeTrue();
  error OnlyCallableByOwnerOrAdmin();
  error OnlyCallableByOwnerOrAdminOrNop();
  error InvalidWithdrawParams();
  error NoFeesToPay();
  error NoNopsToPay();
  error InsufficientBalance();
  error TooManyNops();
  error MaxFeeBalanceReached();
  error MessageTooLarge(uint256 maxSize, uint256 actualSize);
  error MessageGasLimitTooHigh();
  error UnsupportedNumberOfTokens();
  error UnsupportedToken(address token);
  error MustBeCalledByRouter();
  error RouterMustSetOriginalSender();
  error InvalidConfig();
  error CursedByRMN();
  error LinkBalanceNotSettled();
  error InvalidNopAddress(address nop);
  error NotAFeeToken(address token);
  error CannotSendZeroTokens();
  error SourceTokenDataTooLarge(address token);
  error InvalidChainSelector(uint64 chainSelector);
  error GetSupportedTokensFunctionalityRemovedCheckAdminRegistry();
  error InvalidDestBytesOverhead(address token, uint32 destBytesOverhead);

  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event NopPaid(address indexed nop, uint256 amount);
  event FeeConfigSet(FeeTokenConfigArgs[] feeConfig);
  event TokenTransferFeeConfigSet(TokenTransferFeeConfigArgs[] transferFeeConfig);
  event TokenTransferFeeConfigDeleted(address[] tokens);
  /// RMN depends on this event, if changing, please notify the RMN maintainers.
  event CCIPSendRequested(Internal.EVM2EVMMessage message);
  event NopsSet(uint256 nopWeightsTotal, NopAndWeight[] nopsAndWeights);

  /// @dev Struct that contains the static configuration
  /// RMN depends on this struct, if changing, please notify the RMN maintainers.
  //solhint-disable gas-struct-packing
  struct StaticConfig {
    address linkToken; // ────────╮ Link token address
    uint64 chainSelector; // ─────╯ Source chainSelector
    uint64 destChainSelector; // ─╮ Destination chainSelector
    uint64 defaultTxGasLimit; //  │ Default gas limit for a tx
    uint96 maxNopFeesJuels; // ───╯ Max nop fee balance onramp can have
    address prevOnRamp; //          Address of previous-version OnRamp
    address rmnProxy; //            Address of RMN proxy
    address tokenAdminRegistry; //  Address of the token admin registry
  }

  /// @dev Struct to contains the dynamic configuration
  struct DynamicConfig {
    address router; // ──────────────────────────╮ Router address
    uint16 maxNumberOfTokensPerMsg; //           │ Maximum number of distinct ERC20 token transferred per message
    uint32 destGasOverhead; //                   │ Gas charged on top of the gasLimit to cover destination chain costs
    uint16 destGasPerPayloadByte; //             │ Destination chain gas charged for passing each byte of `data` payload to receiver
    uint32 destDataAvailabilityOverheadGas; // ──╯ Extra data availability gas charged on top of the message, e.g. for OCR
    uint16 destGasPerDataAvailabilityByte; // ───╮ Amount of gas to charge per byte of message data that needs availability
    uint16 destDataAvailabilityMultiplierBps; // │ Multiplier for data availability gas, multiples of bps, or 0.0001
    address priceRegistry; //                    │ Price registry address
    uint32 maxDataBytes; //                      │ Maximum payload data size in bytes
    uint32 maxPerMsgGasLimit; // ────────────────╯ Maximum gas limit for messages targeting EVMs
    //                                           │
    // The following three properties are defaults, they can be overridden by setting the TokenTransferFeeConfig for a token
    uint16 defaultTokenFeeUSDCents; // ──────────╮ Default token fee charged per token transfer
    uint32 defaultTokenDestGasOverhead; //       │ Default gas charged to execute the token transfer on the destination chain
    bool enforceOutOfOrder; // ──────────────────╯ Whether to enforce the allowOutOfOrderExecution extraArg value to be true.
  }

  /// @dev Struct to hold the execution fee configuration for a fee token
  struct FeeTokenConfig {
    uint32 networkFeeUSDCents; // ─────────╮ Flat network fee to charge for messages,  multiples of 0.01 USD
    uint64 gasMultiplierWeiPerEth; //      │ Multiplier for gas costs, 1e18 based so 11e17 = 10% extra cost.
    uint64 premiumMultiplierWeiPerEth; //  │ Multiplier for fee-token-specific premiums
    bool enabled; // ──────────────────────╯ Whether this fee token is enabled
  }

  /// @dev Struct to hold the fee configuration for a fee token, same as the FeeTokenConfig but with
  /// token included so that an array of these can be passed in to setFeeTokenConfig to set the mapping
  struct FeeTokenConfigArgs {
    address token; // ─────────────────────╮ Token address
    uint32 networkFeeUSDCents; //          │ Flat network fee to charge for messages,  multiples of 0.01 USD
    uint64 gasMultiplierWeiPerEth; // ─────╯ Multiplier for gas costs, 1e18 based so 11e17 = 10% extra cost
    uint64 premiumMultiplierWeiPerEth; // ─╮ Multiplier for fee-token-specific premiums, 1e18 based
    bool enabled; // ──────────────────────╯ Whether this fee token is enabled
  }

  /// @dev Struct to hold the transfer fee configuration for token transfers
  struct TokenTransferFeeConfig {
    uint32 minFeeUSDCents; // ──────────╮ Minimum fee to charge per token transfer, multiples of 0.01 USD
    uint32 maxFeeUSDCents; //           │ Maximum fee to charge per token transfer, multiples of 0.01 USD
    uint16 deciBps; //                  │ Basis points charged on token transfers, multiples of 0.1bps, or 1e-5
    uint32 destGasOverhead; //          │ Gas charged to execute the token transfer on the destination chain
    //                                  │ Extra data availability bytes that are returned from the source pool and sent
    uint32 destBytesOverhead; //        │ to the destination pool. Must be >= Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES
    bool aggregateRateLimitEnabled; //  │ Whether this transfer token is to be included in Aggregate Rate Limiting
    bool isEnabled; // ─────────────────╯ Whether this token has custom transfer fees
  }

  /// @dev Same as TokenTransferFeeConfig
  /// token included so that an array of these can be passed in to setTokenTransferFeeConfig
  struct TokenTransferFeeConfigArgs {
    address token; // ──────────────────╮ Token address
    uint32 minFeeUSDCents; //           │ Minimum fee to charge per token transfer, multiples of 0.01 USD
    uint32 maxFeeUSDCents; //           │ Maximum fee to charge per token transfer, multiples of 0.01 USD
    uint16 deciBps; // ─────────────────╯ Basis points charged on token transfers, multiples of 0.1bps, or 1e-5
    uint32 destGasOverhead; // ─────────╮ Gas charged to execute the token transfer on the destination chain
    //                                  │ Extra data availability bytes that are returned from the source pool and sent
    uint32 destBytesOverhead; //        │ to the destination pool. Must be >= Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES
    bool aggregateRateLimitEnabled; // ─╯ Whether this transfer token is to be included in Aggregate Rate Limiting
  }

  /// @dev Nop address and weight, used to set the nops and their weights
  struct NopAndWeight {
    address nop; // ────╮ Address of the node operator
    uint16 weight; // ──╯ Weight for nop rewards
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMOnRamp 1.5.0";
  /// @dev metadataHash is a lane-specific prefix for a message hash preimage which ensures global uniqueness
  /// Ensures that 2 identical messages sent to 2 different lanes will have a distinct hash.
  /// Must match the metadataHash used in computing leaf hashes offchain for the root committed in
  /// the commitStore and i_metadataHash in the offRamp.
  bytes32 internal immutable i_metadataHash;
  /// @dev Default gas limit for a transactions that did not specify
  /// a gas limit in the extraArgs.
  uint64 internal immutable i_defaultTxGasLimit;
  /// @dev Maximum nop fee that can accumulate in this onramp
  uint96 internal immutable i_maxNopFeesJuels;
  /// @dev The link token address - known to pay nops for their work
  address internal immutable i_linkToken;
  /// @dev The chain ID of the source chain that this contract is deployed to
  uint64 internal immutable i_chainSelector;
  /// @dev The chain ID of the destination chain
  uint64 internal immutable i_destChainSelector;
  /// @dev The address of previous-version OnRamp for this lane
  /// Used to be able to provide sequencing continuity during a zero downtime upgrade.
  address internal immutable i_prevOnRamp;
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;
  /// @dev The address of the token admin registry
  address internal immutable i_tokenAdminRegistry;
  /// @dev the maximum number of nops that can be configured at the same time.
  /// Used to bound gas for loops over nops.
  uint256 private constant MAX_NUMBER_OF_NOPS = 64;

  // DYNAMIC CONFIG
  /// @dev The config for the onRamp
  DynamicConfig internal s_dynamicConfig;
  /// @dev (address nop => uint256 weight)
  EnumerableMap.AddressToUintMap internal s_nops;

  /// @dev The execution fee token config that can be set by the owner or fee admin
  mapping(address token => FeeTokenConfig feeTokenConfig) internal s_feeTokenConfig;
  /// @dev The token transfer fee config that can be set by the owner or fee admin
  mapping(address token => TokenTransferFeeConfig tranferFeeConfig) internal s_tokenTransferFeeConfig;

  // STATE
  /// @dev The current nonce per sender.
  /// The offramp has a corresponding s_senderNonce mapping to ensure messages
  /// are executed in the same order they are sent.
  mapping(address sender => uint64 nonce) internal s_senderNonce;
  /// @dev The amount of LINK available to pay NOPS
  uint96 internal s_nopFeesJuels;
  /// @dev The combined weight of all NOPs weights
  uint32 internal s_nopWeightsTotal;
  /// @dev The last used sequence number. This is zero in the case where no
  /// messages has been sent yet. 0 is not a valid sequence number for any
  /// real transaction.
  uint64 internal s_sequenceNumber;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    RateLimiter.Config memory rateLimiterConfig,
    FeeTokenConfigArgs[] memory feeTokenConfigs,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    NopAndWeight[] memory nopsAndWeights
  ) AggregateRateLimiter(rateLimiterConfig) {
    if (
      staticConfig.linkToken == address(0) || staticConfig.chainSelector == 0 || staticConfig.destChainSelector == 0
        || staticConfig.defaultTxGasLimit == 0 || staticConfig.rmnProxy == address(0)
        || staticConfig.tokenAdminRegistry == address(0)
    ) revert InvalidConfig();

    i_metadataHash = keccak256(
      abi.encode(
        Internal.EVM_2_EVM_MESSAGE_HASH, staticConfig.chainSelector, staticConfig.destChainSelector, address(this)
      )
    );
    i_linkToken = staticConfig.linkToken;
    i_chainSelector = staticConfig.chainSelector;
    i_destChainSelector = staticConfig.destChainSelector;
    i_defaultTxGasLimit = staticConfig.defaultTxGasLimit;
    i_maxNopFeesJuels = staticConfig.maxNopFeesJuels;
    i_prevOnRamp = staticConfig.prevOnRamp;
    i_rmnProxy = staticConfig.rmnProxy;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;

    _setDynamicConfig(dynamicConfig);
    _setFeeTokenConfig(feeTokenConfigs);
    _setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, new address[](0));
    _setNops(nopsAndWeights);
  }

  // ================================================================
  // │                          Messaging                           │
  // ================================================================

  /// @inheritdoc IEVM2AnyOnRamp
  function getExpectedNextSequenceNumber() external view returns (uint64) {
    return s_sequenceNumber + 1;
  }

  /// @inheritdoc IEVM2AnyOnRamp
  function getSenderNonce(address sender) external view returns (uint64) {
    uint256 senderNonce = s_senderNonce[sender];

    if (i_prevOnRamp != address(0)) {
      if (senderNonce == 0) {
        // If OnRamp was upgraded, check if sender has a nonce from the previous OnRamp.
        return IEVM2AnyOnRamp(i_prevOnRamp).getSenderNonce(sender);
      }
    }
    return uint64(senderNonce);
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function forwardFromRouter(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message,
    uint256 feeTokenAmount,
    address originalSender
  ) external returns (bytes32) {
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(destChainSelector)))) revert CursedByRMN();
    // Validate message sender is set and allowed. Not validated in `getFee` since it is not user-driven.
    if (originalSender == address(0)) revert RouterMustSetOriginalSender();
    // Router address may be zero intentionally to pause.
    if (msg.sender != s_dynamicConfig.router) revert MustBeCalledByRouter();
    if (destChainSelector != i_destChainSelector) revert InvalidChainSelector(destChainSelector);

    Client.EVMExtraArgsV2 memory extraArgs = _fromBytes(message.extraArgs);
    // Validate the message with various checks
    uint256 numberOfTokens = message.tokenAmounts.length;
    _validateMessage(message.data.length, extraArgs.gasLimit, numberOfTokens, extraArgs.allowOutOfOrderExecution);

    // Only check token value if there are tokens
    if (numberOfTokens > 0) {
      uint256 value;
      for (uint256 i = 0; i < numberOfTokens; ++i) {
        if (message.tokenAmounts[i].amount == 0) revert CannotSendZeroTokens();
        if (s_tokenTransferFeeConfig[message.tokenAmounts[i].token].aggregateRateLimitEnabled) {
          value += _getTokenValue(message.tokenAmounts[i], IPriceRegistry(s_dynamicConfig.priceRegistry));
        }
      }
      // Rate limit on aggregated token value
      if (value > 0) _rateLimitValue(value);
    }

    // Convert feeToken to link if not already in link
    if (message.feeToken == i_linkToken) {
      // Since there is only 1b link this is safe
      s_nopFeesJuels += uint96(feeTokenAmount);
    } else {
      // the cast from uint256 to uint96 is considered safe, uint96 can store more than max supply of link token
      s_nopFeesJuels += uint96(
        IPriceRegistry(s_dynamicConfig.priceRegistry).convertTokenAmount(message.feeToken, feeTokenAmount, i_linkToken)
      );
    }
    if (s_nopFeesJuels > i_maxNopFeesJuels) revert MaxFeeBalanceReached();

    // Get the current nonce if the message is an ordered message. If it's not ordered, we don't have to make the
    // external call.
    if (!extraArgs.allowOutOfOrderExecution) {
      if (i_prevOnRamp != address(0)) {
        if (s_senderNonce[originalSender] == 0) {
          // If this is first time send for a sender in new OnRamp, check if they have a nonce
          // from the previous OnRamp and start from there instead of zero.
          s_senderNonce[originalSender] = IEVM2AnyOnRamp(i_prevOnRamp).getSenderNonce(originalSender);
        }
      }
    }

    // We need the next available sequence number so we increment before we use the value
    Internal.EVM2EVMMessage memory newMessage = Internal.EVM2EVMMessage({
      sourceChainSelector: i_chainSelector,
      sender: originalSender,
      // EVM destination addresses should be abi encoded and therefore always 32 bytes long
      // Not duplicately validated in `getFee`. Invalid address is uncommon, gas cost outweighs UX gain.
      receiver: Internal._validateEVMAddress(message.receiver),
      sequenceNumber: ++s_sequenceNumber,
      gasLimit: extraArgs.gasLimit,
      strict: false,
      // Only bump nonce for messages that specify allowOutOfOrderExecution == false. Otherwise, we
      // may block ordered message nonces, which is not what we want.
      nonce: extraArgs.allowOutOfOrderExecution ? 0 : ++s_senderNonce[originalSender],
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      data: message.data,
      tokenAmounts: message.tokenAmounts,
      sourceTokenData: new bytes[](numberOfTokens), // will be populated below
      messageId: ""
    });

    // Lock the tokens as last step. TokenPools may not always be trusted.
    // There should be no state changes after external call to TokenPools.
    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Client.EVMTokenAmount memory tokenAndAmount = message.tokenAmounts[i];
      IPoolV1 sourcePool = getPoolBySourceToken(destChainSelector, IERC20(tokenAndAmount.token));
      // We don't have to check if it supports the pool version in a non-reverting way here because
      // if we revert here, there is no effect on CCIP. Therefore we directly call the supportsInterface
      // function and not through the ERC165Checker.
      if (address(sourcePool) == address(0) || !sourcePool.supportsInterface(Pool.CCIP_POOL_V1)) {
        revert UnsupportedToken(tokenAndAmount.token);
      }

      Pool.LockOrBurnOutV1 memory poolReturnData = sourcePool.lockOrBurn(
        Pool.LockOrBurnInV1({
          receiver: message.receiver,
          remoteChainSelector: i_destChainSelector,
          originalSender: originalSender,
          amount: tokenAndAmount.amount,
          localToken: tokenAndAmount.token
        })
      );

      // Since the DON has to pay for the extraData to be included on the destination chain, we cap the length of the
      // extraData. This prevents gas bomb attacks on the NOPs. As destBytesOverhead accounts for both
      // extraData and offchainData, this caps the worst case abuse to the number of bytes reserved for offchainData.
      if (poolReturnData.destPoolData.length > Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) {
        // If TokenTransferFeeConfig.enabled is false, there is no config. That means destBytesOverhead is zero and
        // this check is always true. That ensures that a pool without config cannot send more than
        // Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes of data.
        if (poolReturnData.destPoolData.length > s_tokenTransferFeeConfig[tokenAndAmount.token].destBytesOverhead) {
          revert SourceTokenDataTooLarge(tokenAndAmount.token);
        }
      }

      // We validate the token address to ensure it is a valid EVM address
      Internal._validateEVMAddress(poolReturnData.destTokenAddress);

      newMessage.sourceTokenData[i] = abi.encode(
        Internal.SourceTokenData({
          sourcePoolAddress: abi.encode(sourcePool),
          destTokenAddress: poolReturnData.destTokenAddress,
          extraData: poolReturnData.destPoolData,
          // The user will be billed either the default or the override, so we send the exact amount that we billed for
          // to the destination chain to be used for the token releaseOrMint and transfer.
          destGasAmount: s_tokenTransferFeeConfig[tokenAndAmount.token].isEnabled
            ? s_tokenTransferFeeConfig[tokenAndAmount.token].destGasOverhead
            : s_dynamicConfig.defaultTokenDestGasOverhead
        })
      );
    }

    // Hash only after the sourceTokenData has been set
    newMessage.messageId = Internal._hash(newMessage, i_metadataHash);

    // Emit message request
    // This must happen after any pool events as some tokens (e.g. USDC) emit events that we expect to precede this
    // event in the offchain code.
    emit CCIPSendRequested(newMessage);
    return newMessage.messageId;
  }

  /// @dev Convert the extra args bytes into a struct
  /// @param extraArgs The extra args bytes
  /// @return The extra args struct
  function _fromBytes(bytes calldata extraArgs) internal view returns (Client.EVMExtraArgsV2 memory) {
    if (extraArgs.length == 0) {
      return Client.EVMExtraArgsV2({gasLimit: i_defaultTxGasLimit, allowOutOfOrderExecution: false});
    }

    bytes4 extraArgsTag = bytes4(extraArgs);
    if (extraArgsTag == Client.EVM_EXTRA_ARGS_V2_TAG) {
      return abi.decode(extraArgs[4:], (Client.EVMExtraArgsV2));
    } else if (extraArgsTag == Client.EVM_EXTRA_ARGS_V1_TAG) {
      // EVMExtraArgsV1 originally included a second boolean (strict) field which has been deprecated.
      // Clients may still include it but it will be ignored.
      return Client.EVMExtraArgsV2({gasLimit: abi.decode(extraArgs[4:], (uint256)), allowOutOfOrderExecution: false});
    }

    revert InvalidExtraArgsTag();
  }

  /// @notice Validate the forwarded message with various checks.
  /// @dev This function can be called multiple times during a CCIPSend,
  /// only common user-driven mistakes are validated here to minimize duplicate validation cost.
  /// @param dataLength The length of the data field of the message.
  /// @param gasLimit The gasLimit set in message for destination execution.
  /// @param numberOfTokens The number of tokens to be sent.
  function _validateMessage(
    uint256 dataLength,
    uint256 gasLimit,
    uint256 numberOfTokens,
    bool allowOutOfOrderExecution
  ) internal view {
    uint256 maxDataBytes = uint256(s_dynamicConfig.maxDataBytes);
    if (dataLength > maxDataBytes) revert MessageTooLarge(maxDataBytes, dataLength);
    if (gasLimit > uint256(s_dynamicConfig.maxPerMsgGasLimit)) revert MessageGasLimitTooHigh();
    if (numberOfTokens > uint256(s_dynamicConfig.maxNumberOfTokensPerMsg)) revert UnsupportedNumberOfTokens();
    if (!allowOutOfOrderExecution) {
      if (s_dynamicConfig.enforceOutOfOrder) {
        revert ExtraArgOutOfOrderExecutionMustBeTrue();
      }
    }
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static onRamp config.
  /// @dev RMN depends on this function, if changing, please notify the RMN maintainers.
  /// @return the configuration.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({
      linkToken: i_linkToken,
      chainSelector: i_chainSelector,
      destChainSelector: i_destChainSelector,
      defaultTxGasLimit: i_defaultTxGasLimit,
      maxNopFeesJuels: i_maxNopFeesJuels,
      prevOnRamp: i_prevOnRamp,
      rmnProxy: i_rmnProxy,
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

  /// @notice Internal version of setDynamicConfig to allow for reuse in the constructor.
  function _setDynamicConfig(DynamicConfig memory dynamicConfig) internal {
    // We permit router to be set to zero as a way to pause the contract.
    if (dynamicConfig.priceRegistry == address(0)) revert InvalidConfig();
    s_dynamicConfig = dynamicConfig;

    emit ConfigSet(
      StaticConfig({
        linkToken: i_linkToken,
        chainSelector: i_chainSelector,
        destChainSelector: i_destChainSelector,
        defaultTxGasLimit: i_defaultTxGasLimit,
        maxNopFeesJuels: i_maxNopFeesJuels,
        prevOnRamp: i_prevOnRamp,
        rmnProxy: i_rmnProxy,
        tokenAdminRegistry: i_tokenAdminRegistry
      }),
      dynamicConfig
    );
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @inheritdoc IEVM2AnyOnRampClient
  function getPoolBySourceToken(uint64, /*destChainSelector*/ IERC20 sourceToken) public view returns (IPoolV1) {
    return IPoolV1(ITokenAdminRegistry(i_tokenAdminRegistry).getPool(address(sourceToken)));
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function getSupportedTokens(uint64) external pure returns (address[] memory) {
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
    if (destChainSelector != i_destChainSelector) revert InvalidChainSelector(destChainSelector);

    Client.EVMExtraArgsV2 memory extraArgs = _fromBytes(message.extraArgs);
    // Validate the message with various checks
    _validateMessage(
      message.data.length, extraArgs.gasLimit, message.tokenAmounts.length, extraArgs.allowOutOfOrderExecution
    );

    FeeTokenConfig memory feeTokenConfig = s_feeTokenConfig[message.feeToken];
    if (!feeTokenConfig.enabled) revert NotAFeeToken(message.feeToken);

    (uint224 feeTokenPrice, uint224 packedGasPrice) =
      IPriceRegistry(s_dynamicConfig.priceRegistry).getTokenAndGasPrices(message.feeToken, destChainSelector);

    // Calculate premiumFee in USD with 18 decimals precision first.
    // If message-only and no token transfers, a flat network fee is charged.
    // If there are token transfers, premiumFee is calculated from token transfer fee.
    // If there are both token transfers and message, premiumFee is only calculated from token transfer fee.
    uint256 premiumFee = 0;
    uint32 tokenTransferGas = 0;
    uint32 tokenTransferBytesOverhead = 0;
    if (message.tokenAmounts.length > 0) {
      (premiumFee, tokenTransferGas, tokenTransferBytesOverhead) =
        _getTokenTransferCost(message.feeToken, feeTokenPrice, message.tokenAmounts);
    } else {
      // Convert USD cents with 2 decimals to 18 decimals.
      premiumFee = uint256(feeTokenConfig.networkFeeUSDCents) * 1e16;
    }

    // Calculate data availability cost in USD with 36 decimals. Data availability cost exists on rollups that need to post
    // transaction calldata onto another storage layer, e.g. Eth mainnet, incurring additional storage gas costs.
    uint256 dataAvailabilityCost = 0;
    // Only calculate data availability cost if data availability multiplier is non-zero.
    // The multiplier should be set to 0 if destination chain does not charge data availability cost.
    if (s_dynamicConfig.destDataAvailabilityMultiplierBps > 0) {
      dataAvailabilityCost = _getDataAvailabilityCost(
        // Parse the data availability gas price stored in the higher-order 112 bits of the encoded gas price.
        uint112(packedGasPrice >> Internal.GAS_PRICE_BITS),
        message.data.length,
        message.tokenAmounts.length,
        tokenTransferBytesOverhead
      );
    }

    // Calculate execution gas fee on destination chain in USD with 36 decimals.
    // We add the message gas limit, the overhead gas, the gas of passing message data to receiver, and token transfer gas together.
    // We then multiply this gas total with the gas multiplier and gas price, converting it into USD with 36 decimals.
    // uint112(packedGasPrice) = executionGasPrice
    uint256 executionCost = uint112(packedGasPrice)
      * (
        extraArgs.gasLimit + s_dynamicConfig.destGasOverhead
          + (message.data.length * s_dynamicConfig.destGasPerPayloadByte) + tokenTransferGas
      ) * feeTokenConfig.gasMultiplierWeiPerEth;

    // Calculate number of fee tokens to charge.
    // Total USD fee is in 36 decimals, feeTokenPrice is in 18 decimals USD for 1e18 smallest token denominations.
    // Result of the division is the number of smallest token denominations.
    return
      ((premiumFee * feeTokenConfig.premiumMultiplierWeiPerEth) + executionCost + dataAvailabilityCost) / feeTokenPrice;
  }

  /// @notice Returns the estimated data availability cost of the message.
  /// @dev To save on gas, we use a single destGasPerDataAvailabilityByte value for both zero and non-zero bytes.
  /// @param dataAvailabilityGasPrice USD per data availability gas in 18 decimals.
  /// @param messageDataLength length of the data field in the message.
  /// @param numberOfTokens number of distinct token transfers in the message.
  /// @param tokenTransferBytesOverhead additional token transfer data passed to destination, e.g. USDC attestation.
  /// @return dataAvailabilityCostUSD36Decimal total data availability cost in USD with 36 decimals.
  function _getDataAvailabilityCost(
    uint112 dataAvailabilityGasPrice,
    uint256 messageDataLength,
    uint256 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) internal view returns (uint256 dataAvailabilityCostUSD36Decimal) {
    // dataAvailabilityLengthBytes sums up byte lengths of fixed message fields and dynamic message fields.
    // Fixed message fields do account for the offset and length slot of the dynamic fields.
    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    // destDataAvailabilityOverheadGas is a separate config value for flexibility to be updated independently of message cost.
    // Its value is determined by CCIP lane implementation, e.g. the overhead data posted for OCR.
    uint256 dataAvailabilityGas = (dataAvailabilityLengthBytes * s_dynamicConfig.destGasPerDataAvailabilityByte)
      + s_dynamicConfig.destDataAvailabilityOverheadGas;

    // dataAvailabilityGasPrice is in 18 decimals, destDataAvailabilityMultiplierBps is in 4 decimals
    // We pad 14 decimals to bring the result to 36 decimals, in line with token bps and execution fee.
    return ((dataAvailabilityGas * dataAvailabilityGasPrice) * s_dynamicConfig.destDataAvailabilityMultiplierBps) * 1e14;
  }

  /// @notice Returns the token transfer cost parameters.
  /// A basis point fee is calculated from the USD value of each token transfer.
  /// For each individual transfer, this fee is between [minFeeUSD, maxFeeUSD].
  /// Total transfer fee is the sum of each individual token transfer fee.
  /// @dev Assumes that tokenAmounts are validated to be listed tokens elsewhere.
  /// @dev Splitting one token transfer into multiple transfers is discouraged,
  /// as it will result in a transferFee equal or greater than the same amount aggregated/de-duped.
  /// @param feeToken address of the feeToken.
  /// @param feeTokenPrice price of feeToken in USD with 18 decimals.
  /// @param tokenAmounts token transfers in the message.
  /// @return tokenTransferFeeUSDWei total token transfer bps fee in USD with 18 decimals.
  /// @return tokenTransferGas total execution gas of the token transfers.
  /// @return tokenTransferBytesOverhead additional token transfer data passed to destination, e.g. USDC attestation.
  function _getTokenTransferCost(
    address feeToken,
    uint224 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts
  ) internal view returns (uint256 tokenTransferFeeUSDWei, uint32 tokenTransferGas, uint32 tokenTransferBytesOverhead) {
    uint256 numberOfTokens = tokenAmounts.length;

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Client.EVMTokenAmount memory tokenAmount = tokenAmounts[i];

      // Validate if the token is supported, do not calculate fee for unsupported tokens.
      if (address(getPoolBySourceToken(i_destChainSelector, IERC20(tokenAmount.token))) == address(0)) {
        revert UnsupportedToken(tokenAmount.token);
      }

      TokenTransferFeeConfig memory transferFeeConfig = s_tokenTransferFeeConfig[tokenAmount.token];

      // If the token has no specific overrides configured, we use the global defaults.
      if (!transferFeeConfig.isEnabled) {
        tokenTransferFeeUSDWei += uint256(s_dynamicConfig.defaultTokenFeeUSDCents) * 1e16;
        tokenTransferGas += s_dynamicConfig.defaultTokenDestGasOverhead;
        tokenTransferBytesOverhead += Pool.CCIP_POOL_V1_RET_BYTES;
        continue;
      }

      uint256 bpsFeeUSDWei = 0;
      // Only calculate bps fee if ratio is greater than 0. Ratio of 0 means no bps fee for a token.
      // Useful for when the PriceRegistry cannot return a valid price for the token.
      if (transferFeeConfig.deciBps > 0) {
        uint224 tokenPrice = 0;
        if (tokenAmount.token != feeToken) {
          tokenPrice = IPriceRegistry(s_dynamicConfig.priceRegistry).getValidatedTokenPrice(tokenAmount.token);
        } else {
          tokenPrice = feeTokenPrice;
        }

        // Calculate token transfer value, then apply fee ratio
        // ratio represents multiples of 0.1bps, or 1e-5
        bpsFeeUSDWei = (tokenPrice._calcUSDValueFromTokenAmount(tokenAmount.amount) * transferFeeConfig.deciBps) / 1e5;
      }

      tokenTransferGas += transferFeeConfig.destGasOverhead;
      tokenTransferBytesOverhead += transferFeeConfig.destBytesOverhead;

      // Bps fees should be kept within range of [minFeeUSD, maxFeeUSD].
      // Convert USD values with 2 decimals to 18 decimals.
      uint256 minFeeUSDWei = uint256(transferFeeConfig.minFeeUSDCents) * 1e16;
      if (bpsFeeUSDWei < minFeeUSDWei) {
        tokenTransferFeeUSDWei += minFeeUSDWei;
        continue;
      }

      uint256 maxFeeUSDWei = uint256(transferFeeConfig.maxFeeUSDCents) * 1e16;
      if (bpsFeeUSDWei > maxFeeUSDWei) {
        tokenTransferFeeUSDWei += maxFeeUSDWei;
        continue;
      }

      tokenTransferFeeUSDWei += bpsFeeUSDWei;
    }

    return (tokenTransferFeeUSDWei, tokenTransferGas, tokenTransferBytesOverhead);
  }

  /// @notice Gets the fee configuration for a token
  /// @param token The token to get the fee configuration for
  /// @return feeTokenConfig FeeTokenConfig struct
  function getFeeTokenConfig(address token) external view returns (FeeTokenConfig memory feeTokenConfig) {
    return s_feeTokenConfig[token];
  }

  /// @notice Sets the fee configuration for a token
  /// @param feeTokenConfigArgs Array of FeeTokenConfigArgs structs.
  function setFeeTokenConfig(FeeTokenConfigArgs[] memory feeTokenConfigArgs) external {
    _onlyOwnerOrAdmin();
    _setFeeTokenConfig(feeTokenConfigArgs);
  }

  /// @dev Set the fee config
  /// @param feeTokenConfigArgs The fee token configs.
  function _setFeeTokenConfig(FeeTokenConfigArgs[] memory feeTokenConfigArgs) internal {
    for (uint256 i = 0; i < feeTokenConfigArgs.length; ++i) {
      FeeTokenConfigArgs memory configArg = feeTokenConfigArgs[i];

      s_feeTokenConfig[configArg.token] = FeeTokenConfig({
        networkFeeUSDCents: configArg.networkFeeUSDCents,
        gasMultiplierWeiPerEth: configArg.gasMultiplierWeiPerEth,
        premiumMultiplierWeiPerEth: configArg.premiumMultiplierWeiPerEth,
        enabled: configArg.enabled
      });
    }
    emit FeeConfigSet(feeTokenConfigArgs);
  }

  /// @notice Gets the transfer fee config for a given token.
  function getTokenTransferFeeConfig(
    address token
  ) external view returns (TokenTransferFeeConfig memory tokenTransferFeeConfig) {
    return s_tokenTransferFeeConfig[token];
  }

  /// @notice Sets the transfer fee config.
  /// @dev only callable by the owner or admin.
  function setTokenTransferFeeConfig(
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    address[] memory tokensToUseDefaultFeeConfigs
  ) external {
    _onlyOwnerOrAdmin();
    _setTokenTransferFeeConfig(tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs);
  }

  /// @notice internal helper to set the token transfer fee config.
  function _setTokenTransferFeeConfig(
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    address[] memory tokensToUseDefaultFeeConfigs
  ) internal {
    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      TokenTransferFeeConfigArgs memory configArg = tokenTransferFeeConfigArgs[i];

      if (configArg.destBytesOverhead < Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) {
        revert InvalidDestBytesOverhead(configArg.token, configArg.destBytesOverhead);
      }

      s_tokenTransferFeeConfig[configArg.token] = TokenTransferFeeConfig({
        minFeeUSDCents: configArg.minFeeUSDCents,
        maxFeeUSDCents: configArg.maxFeeUSDCents,
        deciBps: configArg.deciBps,
        destGasOverhead: configArg.destGasOverhead,
        destBytesOverhead: configArg.destBytesOverhead,
        aggregateRateLimitEnabled: configArg.aggregateRateLimitEnabled,
        isEnabled: true
      });
    }
    emit TokenTransferFeeConfigSet(tokenTransferFeeConfigArgs);

    // Remove the custom fee configs for the tokens that are in the tokensToUseDefaultFeeConfigs array
    for (uint256 i = 0; i < tokensToUseDefaultFeeConfigs.length; ++i) {
      delete s_tokenTransferFeeConfig[tokensToUseDefaultFeeConfigs[i]];
    }
    if (tokensToUseDefaultFeeConfigs.length > 0) {
      emit TokenTransferFeeConfigDeleted(tokensToUseDefaultFeeConfigs);
    }
  }

  // ================================================================
  // │                         NOP payments                         │
  // ================================================================

  /// @notice Get the total amount of fees to be paid to the Nops (in LINK)
  /// @return totalNopFees
  function getNopFeesJuels() external view returns (uint96) {
    return s_nopFeesJuels;
  }

  /// @notice Gets the Nops and their weights
  /// @return nopsAndWeights Array of NopAndWeight structs
  /// @return weightsTotal The sum weight of all Nops
  function getNops() external view returns (NopAndWeight[] memory nopsAndWeights, uint256 weightsTotal) {
    uint256 length = s_nops.length();
    nopsAndWeights = new NopAndWeight[](length);
    for (uint256 i = 0; i < length; ++i) {
      (address nopAddress, uint256 nopWeight) = s_nops.at(i);
      nopsAndWeights[i] = NopAndWeight({nop: nopAddress, weight: uint16(nopWeight)});
    }
    weightsTotal = s_nopWeightsTotal;
    return (nopsAndWeights, weightsTotal);
  }

  /// @notice Sets the Nops and their weights
  /// @param nopsAndWeights Array of NopAndWeight structs
  function setNops(NopAndWeight[] calldata nopsAndWeights) external {
    _onlyOwnerOrAdmin();
    _setNops(nopsAndWeights);
  }

  /// @param nopsAndWeights New set of nops and weights
  /// @dev Clears existing nops, sets new nops and weights
  /// @dev We permit fees to accrue before nops are configured, in which case
  /// they will go to the first set of configured nops.
  function _setNops(NopAndWeight[] memory nopsAndWeights) internal {
    uint256 numberOfNops = nopsAndWeights.length;
    if (numberOfNops > MAX_NUMBER_OF_NOPS) revert TooManyNops();

    // Make sure all nops have been paid before removing nops
    // We only have to pay when there are nops and there is enough
    // outstanding NOP balance to trigger a payment.
    if (s_nopWeightsTotal > 0) {
      if (s_nopFeesJuels >= s_nopWeightsTotal) {
        payNops();
      }
    }

    // Remove all previous nops, move from end to start to avoid shifting
    for (uint256 i = s_nops.length(); i > 0; --i) {
      (address nop,) = s_nops.at(i - 1);
      s_nops.remove(nop);
    }

    // Add new
    uint32 nopWeightsTotal = 0;
    // nopWeightsTotal is bounded by the MAX_NUMBER_OF_NOPS and the weight of
    // a single nop being of type uint16. This ensures nopWeightsTotal will
    // always fit into the uint32 type.
    for (uint256 i = 0; i < numberOfNops; ++i) {
      // Make sure the LINK token is not a nop because the link token doesn't allow
      // self transfers. If set as nop, payNops would always revert. Since setNops
      // calls payNops, we can never remove the LINK token as a nop.
      address nop = nopsAndWeights[i].nop;
      uint16 weight = nopsAndWeights[i].weight;
      if (nop == i_linkToken || nop == address(0)) revert InvalidNopAddress(nop);
      s_nops.set(nop, weight);
      nopWeightsTotal += weight;
    }
    s_nopWeightsTotal = nopWeightsTotal;
    emit NopsSet(nopWeightsTotal, nopsAndWeights);
  }

  /// @notice Pays the Node Ops their outstanding balances.
  /// @dev some balance can remain after payments are done. This is at most the sum
  /// of the weight of all nops. Since nop weights are uint16s and we can have at
  /// most MAX_NUMBER_OF_NOPS NOPs, the highest possible value is 2**22 or 0.04 gjuels.
  function payNops() public {
    if (msg.sender != owner()) {
      if (msg.sender != s_admin) {
        if (!s_nops.contains(msg.sender)) {
          revert OnlyCallableByOwnerOrAdminOrNop();
        }
      }
    }
    uint256 weightsTotal = s_nopWeightsTotal;
    if (weightsTotal == 0) revert NoNopsToPay();

    uint96 totalFeesToPay = s_nopFeesJuels;
    if (totalFeesToPay < weightsTotal) revert NoFeesToPay();
    if (linkAvailableForPayment() < 0) revert InsufficientBalance();

    uint96 fundsLeft = totalFeesToPay;
    uint256 numberOfNops = s_nops.length();
    for (uint256 i = 0; i < numberOfNops; ++i) {
      (address nop, uint256 weight) = s_nops.at(i);
      // amount can never be higher than totalFeesToPay so the cast to uint96 is safe
      uint96 amount = uint96((totalFeesToPay * weight) / weightsTotal);
      fundsLeft -= amount;
      IERC20(i_linkToken).safeTransfer(nop, amount);
      emit NopPaid(nop, amount);
    }
    // Some funds can remain, since this is an incredibly small
    // amount we consider this OK.
    s_nopFeesJuels = fundsLeft;
  }

  /// @notice Allows the owner to withdraw any ERC20 token from the contract.
  /// The NOP link balance is not withdrawable.
  /// @param feeToken The token to withdraw
  /// @param to The address to send the tokens to
  function withdrawNonLinkFees(address feeToken, address to) external {
    _onlyOwnerOrAdmin();
    if (to == address(0)) revert InvalidWithdrawParams();

    // We require the link balance to be settled before allowing withdrawal of non-link fees.
    int256 linkAfterNopFees = linkAvailableForPayment();
    if (linkAfterNopFees < 0) revert LinkBalanceNotSettled();

    if (feeToken == i_linkToken) {
      // Withdraw only the left over link balance
      IERC20(feeToken).safeTransfer(to, uint256(linkAfterNopFees));
    } else {
      // Withdrawal all non-link tokens in the contract
      IERC20(feeToken).safeTransfer(to, IERC20(feeToken).balanceOf(address(this)));
    }
  }

  // ================================================================
  // │                        Link monitoring                       │
  // ================================================================

  /// @notice Calculate remaining LINK balance after paying nops
  /// @dev Allow keeper to monitor funds available for paying nops
  /// @return balance if nops were to be paid
  function linkAvailableForPayment() public view returns (int256) {
    // Since LINK caps at uint96, casting to int256 is safe
    return int256(IERC20(i_linkToken).balanceOf(address(this))) - int256(uint256(s_nopFeesJuels));
  }

  // ================================================================
  // │                           Access                             │
  // ================================================================

  /// @dev Require that the sender is the owner or the fee admin
  /// Not a modifier to save on contract size
  function _onlyOwnerOrAdmin() internal view {
    if (msg.sender != owner()) {
      if (msg.sender != s_admin) {
        revert OnlyCallableByOwnerOrAdmin();
      }
    }
  }
}
