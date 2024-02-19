// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IPool} from "../interfaces/pools/IPool.sol";
import {IARM} from "../interfaces/IARM.sol";
import {IPriceRegistry} from "../interfaces/IPriceRegistry.sol";
import {IEVM2AnyOnRamp} from "../interfaces/IEVM2AnyOnRamp.sol";
import {IEVM2AnyOnRampClient} from "../interfaces/IEVM2AnyOnRampClient.sol";
import {ILinkAvailable} from "../interfaces/automation/ILinkAvailable.sol";

import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {Client} from "../libraries/Client.sol";
import {Internal} from "../libraries/Internal.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "../libraries/USDPriceWith18Decimals.sol";
import {EnumerableMapAddresses} from "../../shared/enumerable/EnumerableMapAddresses.sol";

import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {EnumerableMap} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableMap.sol";

/// @notice The onRamp is a contract that handles lane-specific fee logic, NOP payments and
/// bridgeable token support.
/// @dev The EVM2EVMOnRamp, CommitStore and EVM2EVMOffRamp form an xchain upgradeable unit. Any change to one of them
/// results an onchain upgrade of all 3.
contract EVM2EVMOnRamp is IEVM2AnyOnRamp, ILinkAvailable, AggregateRateLimiter, ITypeAndVersion {
  using SafeERC20 for IERC20;
  using EnumerableMap for EnumerableMap.AddressToUintMap;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;
  using USDPriceWith18Decimals for uint224;

  error InvalidExtraArgsTag();
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
  error UnsupportedToken(IERC20 token);
  error MustBeCalledByRouter();
  error RouterMustSetOriginalSender();
  error InvalidTokenPoolConfig();
  error PoolAlreadyAdded();
  error PoolDoesNotExist(address token);
  error TokenPoolMismatch();
  error InvalidConfig();
  error InvalidAddress(bytes encodedAddress);
  error BadARMSignal();
  error LinkBalanceNotSettled();
  error InvalidNopAddress(address nop);
  error NotAFeeToken(address token);
  error CannotSendZeroTokens();
  error SourceTokenDataTooLarge(address token);
  error InvalidChainSelector(uint64 chainSelector);

  event ConfigSet(StaticConfig staticConfig, DynamicConfig dynamicConfig);
  event NopPaid(address indexed nop, uint256 amount);
  event FeeConfigSet(FeeTokenConfigArgs[] feeConfig);
  event TokenTransferFeeConfigSet(TokenTransferFeeConfigArgs[] transferFeeConfig);
  /// RMN depends on this event, if changing, please notify the RMN maintainers.
  event CCIPSendRequested(Internal.EVM2EVMMessage message);
  event NopsSet(uint256 nopWeightsTotal, NopAndWeight[] nopsAndWeights);
  event PoolAdded(address token, address pool);
  event PoolRemoved(address token, address pool);

  /// @dev Struct that contains the static configuration
  /// RMN depends on this struct, if changing, please notify the RMN maintainers.
  struct StaticConfig {
    address linkToken; // ────────╮ Link token address
    uint64 chainSelector; // ─────╯ Source chainSelector
    uint64 destChainSelector; // ─╮ Destination chainSelector
    uint64 defaultTxGasLimit; //  │ Default gas limit for a tx
    uint96 maxNopFeesJuels; // ───╯ Max nop fee balance onramp can have
    address prevOnRamp; //          Address of previous-version OnRamp
    address armProxy; //            Address of ARM proxy
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
    uint32 minFeeUSDCents; // ────╮ Minimum fee to charge per token transfer, multiples of 0.01 USD
    uint32 maxFeeUSDCents; //     │ Maximum fee to charge per token transfer, multiples of 0.01 USD
    uint16 deciBps; //            │ Basis points charged on token transfers, multiples of 0.1bps, or 1e-5
    uint32 destGasOverhead; //    │ Gas charged to execute the token transfer on the destination chain
    uint32 destBytesOverhead; // ─╯ Extra data availability bytes on top of fixed transfer data, including sourceTokenData and offchainData
  }

  /// @dev Same as TokenTransferFeeConfig
  /// token included so that an array of these can be passed in to setTokenTransferFeeConfig
  struct TokenTransferFeeConfigArgs {
    address token; // ────────────╮ Token address
    uint32 minFeeUSDCents; //     │ Minimum fee to charge per token transfer, multiples of 0.01 USD
    uint32 maxFeeUSDCents; //     │ Maximum fee to charge per token transfer, multiples of 0.01 USD
    uint16 deciBps; // ───────────╯ Basis points charged on token transfers, multiples of 0.1bps, or 1e-5
    uint32 destGasOverhead; // ───╮ Gas charged to execute the token transfer on the destination chain
    uint32 destBytesOverhead; // ─╯ Extra data availability bytes on top of fixed transfer data, including sourceTokenData and offchainData
  }

  /// @dev Nop address and weight, used to set the nops and their weights
  struct NopAndWeight {
    address nop; // ────╮ Address of the node operator
    uint16 weight; // ──╯ Weight for nop rewards
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "EVM2EVMOnRamp 1.5.0-dev";
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
  /// @dev The address of the arm proxy
  address internal immutable i_armProxy;
  /// @dev the maximum number of nops that can be configured at the same time.
  /// Used to bound gas for loops over nops.
  uint256 private constant MAX_NUMBER_OF_NOPS = 64;

  // DYNAMIC CONFIG
  /// @dev The config for the onRamp
  DynamicConfig internal s_dynamicConfig;
  /// @dev (address nop => uint256 weight)
  EnumerableMap.AddressToUintMap internal s_nops;
  /// @dev source token => token pool
  EnumerableMapAddresses.AddressToAddressMap private s_poolsBySourceToken;

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
    Internal.PoolUpdate[] memory tokensAndPools,
    RateLimiter.Config memory rateLimiterConfig,
    FeeTokenConfigArgs[] memory feeTokenConfigs,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    NopAndWeight[] memory nopsAndWeights
  ) AggregateRateLimiter(rateLimiterConfig) {
    if (
      staticConfig.linkToken == address(0) ||
      staticConfig.chainSelector == 0 ||
      staticConfig.destChainSelector == 0 ||
      staticConfig.defaultTxGasLimit == 0 ||
      staticConfig.armProxy == address(0)
    ) revert InvalidConfig();

    i_metadataHash = keccak256(
      abi.encode(
        Internal.EVM_2_EVM_MESSAGE_HASH,
        staticConfig.chainSelector,
        staticConfig.destChainSelector,
        address(this)
      )
    );
    i_linkToken = staticConfig.linkToken;
    i_chainSelector = staticConfig.chainSelector;
    i_destChainSelector = staticConfig.destChainSelector;
    i_defaultTxGasLimit = staticConfig.defaultTxGasLimit;
    i_maxNopFeesJuels = staticConfig.maxNopFeesJuels;
    i_prevOnRamp = staticConfig.prevOnRamp;
    i_armProxy = staticConfig.armProxy;

    _setDynamicConfig(dynamicConfig);
    _setFeeTokenConfig(feeTokenConfigs);
    _setTokenTransferFeeConfig(tokenTransferFeeConfigArgs);
    _setNops(nopsAndWeights);

    // Set new tokens and pools
    _applyPoolUpdates(new Internal.PoolUpdate[](0), tokensAndPools);
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

    if (senderNonce == 0 && i_prevOnRamp != address(0)) {
      // If OnRamp was upgraded, check if sender has a nonce from the previous OnRamp.
      return IEVM2AnyOnRamp(i_prevOnRamp).getSenderNonce(sender);
    }
    return uint64(senderNonce);
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function forwardFromRouter(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message,
    uint256 feeTokenAmount,
    address originalSender
  ) external whenHealthy returns (bytes32) {
    // Validate message sender is set and allowed. Not validated in `getFee` since it is not user-driven.
    if (originalSender == address(0)) revert RouterMustSetOriginalSender();
    // Router address may be zero intentionally to pause.
    if (msg.sender != s_dynamicConfig.router) revert MustBeCalledByRouter();
    if (destChainSelector != i_destChainSelector) revert InvalidChainSelector(destChainSelector);

    // EVM destination addresses should be abi encoded and therefore always 32 bytes long
    // Not duplicately validated in `getFee`. Invalid address is uncommon, gas cost outweighs UX gain.
    if (message.receiver.length != 32) revert InvalidAddress(message.receiver);
    uint256 decodedReceiver = abi.decode(message.receiver, (uint256));
    // We want to disallow sending to address(0) and to precompiles, which exist on address(1) through address(9).
    if (decodedReceiver > type(uint160).max || decodedReceiver < 10) revert InvalidAddress(message.receiver);

    uint256 gasLimit = _fromBytes(message.extraArgs).gasLimit;
    // Validate the message with various checks
    uint256 numberOfTokens = message.tokenAmounts.length;
    _validateMessage(message.data.length, gasLimit, numberOfTokens);

    // Only check token value if there are tokens
    if (numberOfTokens > 0) {
      for (uint256 i = 0; i < numberOfTokens; ++i) {
        if (message.tokenAmounts[i].amount == 0) revert CannotSendZeroTokens();
      }
      // Rate limit on aggregated token value
      _rateLimitValue(message.tokenAmounts, IPriceRegistry(s_dynamicConfig.priceRegistry));
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

    if (s_senderNonce[originalSender] == 0 && i_prevOnRamp != address(0)) {
      // If this is first time send for a sender in new OnRamp, check if they have a nonce
      // from the previous OnRamp and start from there instead of zero.
      s_senderNonce[originalSender] = IEVM2AnyOnRamp(i_prevOnRamp).getSenderNonce(originalSender);
    }

    // We need the next available sequence number so we increment before we use the value
    Internal.EVM2EVMMessage memory newMessage = Internal.EVM2EVMMessage({
      sourceChainSelector: i_chainSelector,
      sender: originalSender,
      receiver: address(uint160(decodedReceiver)),
      sequenceNumber: ++s_sequenceNumber,
      gasLimit: gasLimit,
      strict: false,
      nonce: ++s_senderNonce[originalSender],
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      data: message.data,
      tokenAmounts: message.tokenAmounts,
      sourceTokenData: new bytes[](numberOfTokens), // will be filled in later
      messageId: ""
    });

    // Lock the tokens as last step. TokenPools may not always be trusted.
    // There should be no state changes after external call to TokenPools.
    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Client.EVMTokenAmount memory tokenAndAmount = message.tokenAmounts[i];
      bytes memory tokenData = getPoolBySourceToken(destChainSelector, IERC20(tokenAndAmount.token)).lockOrBurn(
        originalSender,
        message.receiver,
        tokenAndAmount.amount,
        i_destChainSelector,
        bytes("") // any future extraArgs component would be added here
      );

      // Since the DON has to pay for the tokenData to be included on the destination chain, we cap the length of the tokenData.
      // This prevents gas bomb attacks on the NOPs. We use destBytesOverhead as a proxy to cap the number of bytes we accept.
      // As destBytesOverhead accounts for tokenData + offchainData, this caps the worst case abuse to the number of bytes reserved for offchainData.
      // It therefore fully mitigates gas bombs for most tokens, as most tokens don't use offchainData.
      if (tokenData.length > s_tokenTransferFeeConfig[tokenAndAmount.token].destBytesOverhead)
        revert SourceTokenDataTooLarge(tokenAndAmount.token);

      newMessage.sourceTokenData[i] = tokenData;
    }

    // Hash only after the sourceTokenData has been set
    newMessage.messageId = Internal._hash(newMessage, i_metadataHash);

    // Emit message request
    // Note this must happen after pools, some tokens (eg USDC) emit events that we
    // expect to directly precede this event.
    emit CCIPSendRequested(newMessage);
    return newMessage.messageId;
  }

  /// @dev Convert the extra args bytes into a struct
  /// @param extraArgs The extra args bytes
  /// @return The extra args struct
  function _fromBytes(bytes calldata extraArgs) internal view returns (Client.EVMExtraArgsV1 memory) {
    if (extraArgs.length == 0) {
      return Client.EVMExtraArgsV1({gasLimit: i_defaultTxGasLimit});
    }
    if (bytes4(extraArgs) != Client.EVM_EXTRA_ARGS_V1_TAG) revert InvalidExtraArgsTag();
    // EVMExtraArgsV1 originally included a second boolean (strict) field which we have deprecated entirely.
    // Clients may still send that version but it will be ignored.
    return abi.decode(extraArgs[4:], (Client.EVMExtraArgsV1));
  }

  /// @notice Validate the forwarded message with various checks.
  /// @dev This function can be called multiple times during a CCIPSend,
  /// only common user-driven mistakes are validated here to minimize duplicate validation cost.
  /// @param dataLength The length of the data field of the message.
  /// @param gasLimit The gasLimit set in message for destination execution.
  /// @param numberOfTokens The number of tokens to be sent.
  function _validateMessage(uint256 dataLength, uint256 gasLimit, uint256 numberOfTokens) internal view {
    // Check that payload is formed correctly
    uint256 maxDataBytes = uint256(s_dynamicConfig.maxDataBytes);
    if (dataLength > maxDataBytes) revert MessageTooLarge(maxDataBytes, dataLength);
    if (gasLimit > uint256(s_dynamicConfig.maxPerMsgGasLimit)) revert MessageGasLimitTooHigh();
    if (numberOfTokens > uint256(s_dynamicConfig.maxNumberOfTokensPerMsg)) revert UnsupportedNumberOfTokens();
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static onRamp config.
  /// @dev RMN depends on this function, if changing, please notify the RMN maintainers.
  /// @return the configuration.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return
      StaticConfig({
        linkToken: i_linkToken,
        chainSelector: i_chainSelector,
        destChainSelector: i_destChainSelector,
        defaultTxGasLimit: i_defaultTxGasLimit,
        maxNopFeesJuels: i_maxNopFeesJuels,
        prevOnRamp: i_prevOnRamp,
        armProxy: i_armProxy
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
        armProxy: i_armProxy
      }),
      dynamicConfig
    );
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @inheritdoc IEVM2AnyOnRampClient
  function getSupportedTokens(uint64 /*destChainSelector*/) external view returns (address[] memory) {
    address[] memory sourceTokens = new address[](s_poolsBySourceToken.length());
    for (uint256 i = 0; i < sourceTokens.length; ++i) {
      (sourceTokens[i], ) = s_poolsBySourceToken.at(i);
    }
    return sourceTokens;
  }

  /// @inheritdoc IEVM2AnyOnRampClient
  function getPoolBySourceToken(uint64 /*destChainSelector*/, IERC20 sourceToken) public view returns (IPool) {
    if (!s_poolsBySourceToken.contains(address(sourceToken))) revert UnsupportedToken(sourceToken);
    return IPool(s_poolsBySourceToken.get(address(sourceToken)));
  }

  /// @inheritdoc IEVM2AnyOnRamp
  /// @dev This method can only be called by the owner of the contract.
  function applyPoolUpdates(
    Internal.PoolUpdate[] memory removes,
    Internal.PoolUpdate[] memory adds
  ) external onlyOwner {
    _applyPoolUpdates(removes, adds);
  }

  function _applyPoolUpdates(Internal.PoolUpdate[] memory removes, Internal.PoolUpdate[] memory adds) internal {
    for (uint256 i = 0; i < removes.length; ++i) {
      address token = removes[i].token;
      address pool = removes[i].pool;

      if (!s_poolsBySourceToken.contains(token)) revert PoolDoesNotExist(token);
      if (s_poolsBySourceToken.get(token) != pool) revert TokenPoolMismatch();

      if (s_poolsBySourceToken.remove(token)) {
        emit PoolRemoved(token, pool);
      }
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      address token = adds[i].token;
      address pool = adds[i].pool;

      if (token == address(0) || pool == address(0)) revert InvalidTokenPoolConfig();
      if (token != address(IPool(pool).getToken())) revert TokenPoolMismatch();

      if (s_poolsBySourceToken.set(token, pool)) {
        emit PoolAdded(token, pool);
      } else {
        revert PoolAlreadyAdded();
      }
    }
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

    uint256 gasLimit = _fromBytes(message.extraArgs).gasLimit;
    // Validate the message with various checks
    _validateMessage(message.data.length, gasLimit, message.tokenAmounts.length);

    FeeTokenConfig memory feeTokenConfig = s_feeTokenConfig[message.feeToken];
    if (!feeTokenConfig.enabled) revert NotAFeeToken(message.feeToken);

    (uint224 feeTokenPrice, uint224 packedGasPrice) = IPriceRegistry(s_dynamicConfig.priceRegistry)
      .getTokenAndGasPrices(message.feeToken, destChainSelector);
    uint112 executionGasPrice = uint112(packedGasPrice);

    // Calculate premiumFee in USD with 18 decimals precision first.
    // If message-only and no token transfers, a flat network fee is charged.
    // If there are token transfers, premiumFee is calculated from token transfer fee.
    // If there are both token transfers and message, premiumFee is only calculated from token transfer fee.
    uint256 premiumFee = 0;
    uint32 tokenTransferGas = 0;
    uint32 tokenTransferBytesOverhead = 0;
    if (message.tokenAmounts.length > 0) {
      (premiumFee, tokenTransferGas, tokenTransferBytesOverhead) = _getTokenTransferCost(
        message.feeToken,
        feeTokenPrice,
        message.tokenAmounts
      );
    } else {
      // Convert USD cents with 2 decimals to 18 decimals.
      premiumFee = uint256(feeTokenConfig.networkFeeUSDCents) * 1e16;
    }

    // Apply a feeToken-specific multiplier with 18 decimals, raising the premiumFee to 36 decimals
    premiumFee = premiumFee * feeTokenConfig.premiumMultiplierWeiPerEth;

    // Calculate execution gas fee on destination chain in USD with 36 decimals.
    // We add the message gas limit, the overhead gas, the gas of passing message data to receiver, and token transfer gas together.
    // We then multiply this gas total with the gas multiplier and gas price, converting it into USD with 36 decimals.
    uint256 executionCost = executionGasPrice *
      ((gasLimit +
        s_dynamicConfig.destGasOverhead +
        (message.data.length * s_dynamicConfig.destGasPerPayloadByte) +
        tokenTransferGas) * feeTokenConfig.gasMultiplierWeiPerEth);

    // Calculate data availability cost in USD with 36 decimals. Data availability cost exists on rollups that need to post
    // transaction calldata onto another storage layer, e.g. Eth mainnet, incurring additional storage gas costs.
    uint256 dataAvailabilityCost = 0;
    // Only calculate data availability cost if data availability multiplier is non-zero.
    // The multiplier should be set to 0 if destination chain does not charge data availability cost.
    if (s_dynamicConfig.destDataAvailabilityMultiplierBps > 0) {
      // Parse the data availability gas price stored in the higher-order 112 bits of the encoded gas price.
      uint112 dataAvailabilityGasPrice = uint112(packedGasPrice >> Internal.GAS_PRICE_BITS);

      dataAvailabilityCost = _getDataAvailabilityCost(
        dataAvailabilityGasPrice,
        message.data.length,
        message.tokenAmounts.length,
        tokenTransferBytesOverhead
      );
    }

    // Calculate number of fee tokens to charge.
    // Total USD fee is in 36 decimals, feeTokenPrice is in 18 decimals USD for 1e18 smallest token denominations.
    // Result of the division is the number of smallest token denominations.
    return (premiumFee + executionCost + dataAvailabilityCost) / feeTokenPrice;
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
    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES +
      messageDataLength +
      (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) +
      tokenTransferBytesOverhead;

    // destDataAvailabilityOverheadGas is a separate config value for flexibility to be updated independently of message cost.
    // Its value is determined by CCIP lane implementation, e.g. the overhead data posted for OCR.
    uint256 dataAvailabilityGas = (dataAvailabilityLengthBytes * s_dynamicConfig.destGasPerDataAvailabilityByte) +
      s_dynamicConfig.destDataAvailabilityOverheadGas;

    // dataAvailabilityGasPrice is in 18 decimals, destDataAvailabilityMultiplierBps is in 4 decimals
    // We pad 14 decimals to bring the result to 36 decimals, in line with token bps and execution fee.
    return
      ((dataAvailabilityGas * dataAvailabilityGasPrice) * s_dynamicConfig.destDataAvailabilityMultiplierBps) * 1e14;
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
      TokenTransferFeeConfig memory transferFeeConfig = s_tokenTransferFeeConfig[tokenAmount.token];

      // Validate if the token is supported, do not calculate fee for unsupported tokens.
      if (!s_poolsBySourceToken.contains(tokenAmount.token)) revert UnsupportedToken(IERC20(tokenAmount.token));

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
  function setFeeTokenConfig(FeeTokenConfigArgs[] memory feeTokenConfigArgs) external onlyOwnerOrAdmin {
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
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs
  ) external onlyOwnerOrAdmin {
    _setTokenTransferFeeConfig(tokenTransferFeeConfigArgs);
  }

  /// @notice internal helper to set the token transfer fee config.
  function _setTokenTransferFeeConfig(TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs) internal {
    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      TokenTransferFeeConfigArgs memory configArg = tokenTransferFeeConfigArgs[i];

      s_tokenTransferFeeConfig[configArg.token] = TokenTransferFeeConfig({
        minFeeUSDCents: configArg.minFeeUSDCents,
        maxFeeUSDCents: configArg.maxFeeUSDCents,
        deciBps: configArg.deciBps,
        destGasOverhead: configArg.destGasOverhead,
        destBytesOverhead: configArg.destBytesOverhead
      });
    }
    emit TokenTransferFeeConfigSet(tokenTransferFeeConfigArgs);
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
  function setNops(NopAndWeight[] calldata nopsAndWeights) external onlyOwnerOrAdmin {
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
    if (s_nopWeightsTotal > 0 && s_nopFeesJuels >= s_nopWeightsTotal) {
      payNops();
    }

    // Remove all previous nops, move from end to start to avoid shifting
    for (uint256 i = s_nops.length(); i > 0; --i) {
      (address nop, ) = s_nops.at(i - 1);
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
  function payNops() public onlyOwnerOrAdminOrNop {
    uint256 weightsTotal = s_nopWeightsTotal;
    if (weightsTotal == 0) revert NoNopsToPay();

    uint96 totalFeesToPay = s_nopFeesJuels;
    if (totalFeesToPay < weightsTotal) revert NoFeesToPay();
    if (_linkLeftAfterNopFees() < 0) revert InsufficientBalance();

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
  function withdrawNonLinkFees(address feeToken, address to) external onlyOwnerOrAdmin {
    if (to == address(0)) revert InvalidWithdrawParams();

    // We require the link balance to be settled before allowing withdrawal of non-link fees.
    int256 linkAfterNopFees = _linkLeftAfterNopFees();
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
  /// @return balance if nops were to be paid
  function _linkLeftAfterNopFees() private view returns (int256) {
    // Since LINK caps at uint96, casting to int256 is safe
    return int256(IERC20(i_linkToken).balanceOf(address(this))) - int256(uint256(s_nopFeesJuels));
  }

  /// @notice Allow keeper to monitor funds available for paying nops
  function linkAvailableForPayment() external view returns (int256) {
    return _linkLeftAfterNopFees();
  }

  // ================================================================
  // │                        Access and ARM                        │
  // ================================================================

  /// @dev Require that the sender is the owner or the fee admin or a nop
  modifier onlyOwnerOrAdminOrNop() {
    if (msg.sender != owner() && msg.sender != s_admin && !s_nops.contains(msg.sender))
      revert OnlyCallableByOwnerOrAdminOrNop();
    _;
  }

  /// @dev Require that the sender is the owner or the fee admin
  modifier onlyOwnerOrAdmin() {
    if (msg.sender != owner() && msg.sender != s_admin) revert OnlyCallableByOwnerOrAdmin();
    _;
  }

  /// @notice Ensure that the ARM has not emitted a bad signal, and that the latest heartbeat is not stale.
  modifier whenHealthy() {
    if (IARM(i_armProxy).isCursed()) revert BadARMSignal();
    _;
  }
}
