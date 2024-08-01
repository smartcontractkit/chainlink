// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";
import {ITokenAdminRegistry} from "../../interfaces/ITokenAdminRegistry.sol";

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";

import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {LockReleaseTokenPool} from "../../pools/LockReleaseTokenPool.sol";
import {TokenPool} from "../../pools/TokenPool.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";

import {TokenSetup} from "../TokenSetup.t.sol";
import {MaybeRevertingBurnMintTokenPool} from "../helpers/MaybeRevertingBurnMintTokenPool.sol";
import {PriceRegistryHelper} from "../helpers/PriceRegistryHelper.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

import {Vm} from "forge-std/Vm.sol";
import {console} from "forge-std/console.sol";

contract PriceRegistrySetup is TokenSetup {
  uint112 internal constant USD_PER_GAS = 1e6; // 0.001 gwei
  uint112 internal constant USD_PER_DATA_AVAILABILITY_GAS = 1e9; // 1 gwei

  address internal constant CUSTOM_TOKEN = address(12345);
  uint224 internal constant CUSTOM_TOKEN_PRICE = 1e17; // $0.1 CUSTOM

  // Encode L1 gas price and L2 gas price into a packed price.
  // L1 gas price is left-shifted to the higher-order bits.
  uint224 internal constant PACKED_USD_PER_GAS =
    (uint224(USD_PER_DATA_AVAILABILITY_GAS) << Internal.GAS_PRICE_BITS) + USD_PER_GAS;

  PriceRegistryHelper internal s_priceRegistry;
  // Cheat to store the price updates in storage since struct arrays aren't supported.
  bytes internal s_encodedInitialPriceUpdates;
  address internal s_weth;

  address[] internal s_sourceFeeTokens;
  uint224[] internal s_sourceTokenPrices;
  address[] internal s_destFeeTokens;
  uint224[] internal s_destTokenPrices;

  PriceRegistry.PremiumMultiplierWeiPerEthArgs[] internal s_priceRegistryPremiumMultiplierWeiPerEthArgs;
  PriceRegistry.TokenTransferFeeConfigArgs[] internal s_priceRegistryTokenTransferFeeConfigArgs;

  mapping(address token => address dataFeedAddress) internal s_dataFeedByToken;

  function setUp() public virtual override {
    TokenSetup.setUp();

    _deployTokenPriceDataFeed(s_sourceFeeToken, 8, 1e8);

    s_weth = s_sourceRouter.getWrappedNative();
    _deployTokenPriceDataFeed(s_weth, 8, 1e11);

    address[] memory sourceFeeTokens = new address[](3);
    sourceFeeTokens[0] = s_sourceTokens[0];
    sourceFeeTokens[1] = s_sourceTokens[1];
    sourceFeeTokens[2] = s_sourceRouter.getWrappedNative();
    s_sourceFeeTokens = sourceFeeTokens;

    uint224[] memory sourceTokenPrices = new uint224[](3);
    sourceTokenPrices[0] = 5e18;
    sourceTokenPrices[1] = 2000e18;
    sourceTokenPrices[2] = 2000e18;
    s_sourceTokenPrices = sourceTokenPrices;

    address[] memory destFeeTokens = new address[](3);
    destFeeTokens[0] = s_destTokens[0];
    destFeeTokens[1] = s_destTokens[1];
    destFeeTokens[2] = s_destRouter.getWrappedNative();
    s_destFeeTokens = destFeeTokens;

    uint224[] memory destTokenPrices = new uint224[](3);
    destTokenPrices[0] = 5e18;
    destTokenPrices[1] = 2000e18;
    destTokenPrices[2] = 2000e18;
    s_destTokenPrices = destTokenPrices;

    uint256 sourceTokenCount = sourceFeeTokens.length;
    uint256 destTokenCount = destFeeTokens.length;
    address[] memory pricedTokens = new address[](sourceTokenCount + destTokenCount);
    uint224[] memory tokenPrices = new uint224[](sourceTokenCount + destTokenCount);
    for (uint256 i = 0; i < sourceTokenCount; ++i) {
      pricedTokens[i] = sourceFeeTokens[i];
      tokenPrices[i] = sourceTokenPrices[i];
    }
    for (uint256 i = 0; i < destTokenCount; ++i) {
      pricedTokens[i + sourceTokenCount] = destFeeTokens[i];
      tokenPrices[i + sourceTokenCount] = destTokenPrices[i];
    }

    Internal.PriceUpdates memory priceUpdates = getPriceUpdatesStruct(pricedTokens, tokenPrices);
    priceUpdates.gasPriceUpdates =
      getSingleGasPriceUpdateStruct(DEST_CHAIN_SELECTOR, PACKED_USD_PER_GAS).gasPriceUpdates;

    s_encodedInitialPriceUpdates = abi.encode(priceUpdates);

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = OWNER;
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_weth;
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](0);

    s_priceRegistryPremiumMultiplierWeiPerEthArgs.push(
      PriceRegistry.PremiumMultiplierWeiPerEthArgs({
        token: s_sourceFeeToken,
        premiumMultiplierWeiPerEth: 5e17 // 0.5x
      })
    );
    s_priceRegistryPremiumMultiplierWeiPerEthArgs.push(
      PriceRegistry.PremiumMultiplierWeiPerEthArgs({
        token: s_sourceRouter.getWrappedNative(),
        premiumMultiplierWeiPerEth: 2e18 // 2x
      })
    );

    s_priceRegistryTokenTransferFeeConfigArgs.push();
    s_priceRegistryTokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    s_priceRegistryTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs.push(
      PriceRegistry.TokenTransferFeeConfigSingleTokenArgs({
        token: s_sourceFeeToken,
        tokenTransferFeeConfig: PriceRegistry.TokenTransferFeeConfig({
          minFeeUSDCents: 1_00, // 1 USD
          maxFeeUSDCents: 1000_00, // 1,000 USD
          deciBps: 2_5, // 2.5 bps, or 0.025%
          destGasOverhead: 40_000,
          destBytesOverhead: 32,
          isEnabled: true
        })
      })
    );
    s_priceRegistryTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs.push(
      PriceRegistry.TokenTransferFeeConfigSingleTokenArgs({
        token: CUSTOM_TOKEN,
        tokenTransferFeeConfig: PriceRegistry.TokenTransferFeeConfig({
          minFeeUSDCents: 2_00, // 1 USD
          maxFeeUSDCents: 2000_00, // 1,000 USD
          deciBps: 10_0, // 10 bps, or 0.1%
          destGasOverhead: 1,
          destBytesOverhead: 200,
          isEnabled: true
        })
      })
    );

    s_priceRegistry = new PriceRegistryHelper(
      PriceRegistry.StaticConfig({
        linkToken: s_sourceTokens[0],
        maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
        stalenessThreshold: uint32(TWELVE_HOURS)
      }),
      priceUpdaters,
      feeTokens,
      tokenPriceFeedUpdates,
      s_priceRegistryTokenTransferFeeConfigArgs,
      s_priceRegistryPremiumMultiplierWeiPerEthArgs,
      _generatePriceRegistryDestChainConfigArgs()
    );
    s_priceRegistry.updatePrices(priceUpdates);
  }

  function _deployTokenPriceDataFeed(address token, uint8 decimals, int256 initialAnswer) internal returns (address) {
    MockV3Aggregator dataFeed = new MockV3Aggregator(decimals, initialAnswer);
    s_dataFeedByToken[token] = address(dataFeed);
    return address(dataFeed);
  }

  function getPriceUpdatesStruct(
    address[] memory tokens,
    uint224[] memory prices
  ) internal pure returns (Internal.PriceUpdates memory) {
    uint256 length = tokens.length;

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](length);
    for (uint256 i = 0; i < length; ++i) {
      tokenPriceUpdates[i] = Internal.TokenPriceUpdate({sourceToken: tokens[i], usdPerToken: prices[i]});
    }
    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: new Internal.GasPriceUpdate[](0)});

    return priceUpdates;
  }

  function getEmptyPriceUpdates() internal pure returns (Internal.PriceUpdates memory priceUpdates) {
    return Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
  }

  function getSingleTokenPriceFeedUpdateStruct(
    address sourceToken,
    address dataFeedAddress,
    uint8 tokenDecimals
  ) internal pure returns (PriceRegistry.TokenPriceFeedUpdate memory) {
    return PriceRegistry.TokenPriceFeedUpdate({
      sourceToken: sourceToken,
      feedConfig: IPriceRegistry.TokenPriceFeedConfig({dataFeedAddress: dataFeedAddress, tokenDecimals: tokenDecimals})
    });
  }

  function _initialiseSingleTokenPriceFeed() internal returns (address) {
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    return s_sourceTokens[0];
  }

  function _generateTokenTransferFeeConfigArgs(
    uint256 destChainSelectorLength,
    uint256 tokenLength
  ) internal pure returns (PriceRegistry.TokenTransferFeeConfigArgs[] memory) {
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new PriceRegistry.TokenTransferFeeConfigArgs[](destChainSelectorLength);
    for (uint256 i = 0; i < destChainSelectorLength; ++i) {
      tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs =
        new PriceRegistry.TokenTransferFeeConfigSingleTokenArgs[](tokenLength);
    }
    return tokenTransferFeeConfigArgs;
  }

  function _generatePriceRegistryDestChainConfigArgs()
    internal
    pure
    returns (PriceRegistry.DestChainConfigArgs[] memory)
  {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigs = new PriceRegistry.DestChainConfigArgs[](1);
    destChainConfigs[0] = PriceRegistry.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      destChainConfig: PriceRegistry.DestChainConfig({
        isEnabled: true,
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT,
        defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
        defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        defaultTokenDestBytesOverhead: DEFAULT_TOKEN_BYTES_OVERHEAD,
        defaultTxGasLimit: GAS_LIMIT,
        gasMultiplierWeiPerEth: 5e17,
        networkFeeUSDCents: 1_00,
        enforceOutOfOrder: false,
        chainFamilySelector: Internal.CHAIN_FAMILY_SELECTOR_EVM
      })
    });
    return destChainConfigs;
  }

  function _assertTokenPriceFeedConfigEquality(
    IPriceRegistry.TokenPriceFeedConfig memory config1,
    IPriceRegistry.TokenPriceFeedConfig memory config2
  ) internal pure virtual {
    assertEq(config1.dataFeedAddress, config2.dataFeedAddress);
    assertEq(config1.tokenDecimals, config2.tokenDecimals);
  }

  function _assertTokenPriceFeedConfigUnconfigured(IPriceRegistry.TokenPriceFeedConfig memory config)
    internal
    pure
    virtual
  {
    _assertTokenPriceFeedConfigEquality(
      config, IPriceRegistry.TokenPriceFeedConfig({dataFeedAddress: address(0), tokenDecimals: 0})
    );
  }

  function _assertTokenTransferFeeConfigEqual(
    PriceRegistry.TokenTransferFeeConfig memory a,
    PriceRegistry.TokenTransferFeeConfig memory b
  ) internal pure {
    assertEq(a.minFeeUSDCents, b.minFeeUSDCents);
    assertEq(a.maxFeeUSDCents, b.maxFeeUSDCents);
    assertEq(a.deciBps, b.deciBps);
    assertEq(a.destGasOverhead, b.destGasOverhead);
    assertEq(a.destBytesOverhead, b.destBytesOverhead);
    assertEq(a.isEnabled, b.isEnabled);
  }

  function _assertPriceRegistryStaticConfigsEqual(
    PriceRegistry.StaticConfig memory a,
    PriceRegistry.StaticConfig memory b
  ) internal pure {
    assertEq(a.linkToken, b.linkToken);
    assertEq(a.maxFeeJuelsPerMsg, b.maxFeeJuelsPerMsg);
  }

  function _assertPriceRegistryDestChainConfigsEqual(
    PriceRegistry.DestChainConfig memory a,
    PriceRegistry.DestChainConfig memory b
  ) internal pure {
    assertEq(a.isEnabled, b.isEnabled);
    assertEq(a.maxNumberOfTokensPerMsg, b.maxNumberOfTokensPerMsg);
    assertEq(a.maxDataBytes, b.maxDataBytes);
    assertEq(a.maxPerMsgGasLimit, b.maxPerMsgGasLimit);
    assertEq(a.destGasOverhead, b.destGasOverhead);
    assertEq(a.destGasPerPayloadByte, b.destGasPerPayloadByte);
    assertEq(a.destDataAvailabilityOverheadGas, b.destDataAvailabilityOverheadGas);
    assertEq(a.destGasPerDataAvailabilityByte, b.destGasPerDataAvailabilityByte);
    assertEq(a.destDataAvailabilityMultiplierBps, b.destDataAvailabilityMultiplierBps);
    assertEq(a.defaultTokenFeeUSDCents, b.defaultTokenFeeUSDCents);
    assertEq(a.defaultTokenDestGasOverhead, b.defaultTokenDestGasOverhead);
    assertEq(a.defaultTokenDestBytesOverhead, b.defaultTokenDestBytesOverhead);
    assertEq(a.defaultTxGasLimit, b.defaultTxGasLimit);
  }
}

contract PriceRegistryFeeSetup is PriceRegistrySetup {
  uint224 internal s_feeTokenPrice;
  uint224 internal s_wrappedTokenPrice;
  uint224 internal s_customTokenPrice;

  address internal s_selfServeTokenDefaultPricing = makeAddr("self-serve-token-default-pricing");

  address internal s_destTokenPool = makeAddr("destTokenPool");
  address internal s_destToken = makeAddr("destToken");

  function setUp() public virtual override {
    super.setUp();

    s_feeTokenPrice = s_sourceTokenPrices[0];
    s_wrappedTokenPrice = s_sourceTokenPrices[2];
    s_customTokenPrice = CUSTOM_TOKEN_PRICE;

    s_priceRegistry.updatePrices(getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, CUSTOM_TOKEN_PRICE));
  }

  function _generateEmptyMessage() public view returns (Client.EVM2AnyMessage memory) {
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _generateSingleTokenMessage(
    address token,
    uint256 amount
  ) public view returns (Client.EVM2AnyMessage memory) {
    Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
    tokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: tokenAmounts,
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
  }

  function _messageToEvent(
    Client.EVM2AnyMessage memory message,
    uint64 sourceChainSelector,
    uint64 destChainSelector,
    uint64 seqNum,
    uint64 nonce,
    uint256 feeTokenAmount,
    address originalSender,
    bytes32 metadataHash,
    TokenAdminRegistry tokenAdminRegistry
  ) internal view returns (Internal.EVM2AnyRampMessage memory) {
    Client.EVMExtraArgsV2 memory extraArgs =
      s_priceRegistry.parseEVMExtraArgsFromBytes(message.extraArgs, destChainSelector);

    Internal.EVM2AnyRampMessage memory messageEvent = Internal.EVM2AnyRampMessage({
      header: Internal.RampMessageHeader({
        messageId: "",
        sourceChainSelector: sourceChainSelector,
        destChainSelector: destChainSelector,
        sequenceNumber: seqNum,
        nonce: extraArgs.allowOutOfOrderExecution ? 0 : nonce
      }),
      sender: originalSender,
      data: message.data,
      receiver: message.receiver,
      extraArgs: Client._argsToBytes(extraArgs),
      feeToken: message.feeToken,
      feeTokenAmount: feeTokenAmount,
      tokenAmounts: new Internal.RampTokenAmount[](message.tokenAmounts.length)
    });

    for (uint256 i = 0; i < message.tokenAmounts.length; ++i) {
      messageEvent.tokenAmounts[i] = _getSourceTokenData(message.tokenAmounts[i], tokenAdminRegistry);
    }

    messageEvent.header.messageId = Internal._hash(messageEvent, metadataHash);
    return messageEvent;
  }

  function _getSourceTokenData(
    Client.EVMTokenAmount memory tokenAmount,
    TokenAdminRegistry tokenAdminRegistry
  ) internal view returns (Internal.RampTokenAmount memory) {
    address destToken = s_destTokenBySourceToken[tokenAmount.token];

    return Internal.RampTokenAmount({
      sourcePoolAddress: abi.encode(tokenAdminRegistry.getTokenConfig(tokenAmount.token).tokenPool),
      destTokenAddress: abi.encode(destToken),
      extraData: "",
      amount: tokenAmount.amount
    });
  }

  function calcUSDValueFromTokenAmount(uint224 tokenPrice, uint256 tokenAmount) internal pure returns (uint256) {
    return (tokenPrice * tokenAmount) / 1e18;
  }

  function applyBpsRatio(uint256 tokenAmount, uint16 ratio) internal pure returns (uint256) {
    return (tokenAmount * ratio) / 1e5;
  }

  function configUSDCentToWei(uint256 usdCent) internal pure returns (uint256) {
    return usdCent * 1e16;
  }
}

contract PriceRegistry_constructor is PriceRegistrySetup {
  function test_Setup_Success() public virtual {
    address[] memory priceUpdaters = new address[](2);
    priceUpdaters[0] = STRANGER;
    priceUpdaters[1] = OWNER;
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_sourceTokens[1];
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](2);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);
    tokenPriceFeedUpdates[1] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[1], s_dataFeedByToken[s_sourceTokens[1]], 6);

    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();

    PriceRegistry.StaticConfig memory staticConfig = PriceRegistry.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });
    s_priceRegistry = new PriceRegistryHelper(
      staticConfig,
      priceUpdaters,
      feeTokens,
      tokenPriceFeedUpdates,
      s_priceRegistryTokenTransferFeeConfigArgs,
      s_priceRegistryPremiumMultiplierWeiPerEthArgs,
      destChainConfigArgs
    );

    _assertPriceRegistryStaticConfigsEqual(s_priceRegistry.getStaticConfig(), staticConfig);
    assertEq(feeTokens, s_priceRegistry.getFeeTokens());
    assertEq(priceUpdaters, s_priceRegistry.getAllAuthorizedCallers());
    assertEq(s_priceRegistry.typeAndVersion(), "PriceRegistry 1.6.0-dev");

    _assertTokenPriceFeedConfigEquality(
      tokenPriceFeedUpdates[0].feedConfig, s_priceRegistry.getTokenPriceFeedConfig(s_sourceTokens[0])
    );

    _assertTokenPriceFeedConfigEquality(
      tokenPriceFeedUpdates[1].feedConfig, s_priceRegistry.getTokenPriceFeedConfig(s_sourceTokens[1])
    );

    assertEq(
      s_priceRegistryPremiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(s_priceRegistryPremiumMultiplierWeiPerEthArgs[0].token)
    );

    assertEq(
      s_priceRegistryPremiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(s_priceRegistryPremiumMultiplierWeiPerEthArgs[1].token)
    );

    PriceRegistry.TokenTransferFeeConfigArgs memory tokenTransferFeeConfigArg =
      s_priceRegistryTokenTransferFeeConfigArgs[0];
    for (uint256 i = 0; i < tokenTransferFeeConfigArg.tokenTransferFeeConfigs.length; ++i) {
      PriceRegistry.TokenTransferFeeConfigSingleTokenArgs memory tokenFeeArgs =
        s_priceRegistryTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i];

      _assertTokenTransferFeeConfigEqual(
        tokenFeeArgs.tokenTransferFeeConfig,
        s_priceRegistry.getTokenTransferFeeConfig(tokenTransferFeeConfigArg.destChainSelector, tokenFeeArgs.token)
      );
    }

    for (uint256 i = 0; i < destChainConfigArgs.length; ++i) {
      PriceRegistry.DestChainConfig memory expectedConfig = destChainConfigArgs[i].destChainConfig;
      uint64 destChainSelector = destChainConfigArgs[i].destChainSelector;

      _assertPriceRegistryDestChainConfigsEqual(expectedConfig, s_priceRegistry.getDestChainConfig(destChainSelector));
    }
  }

  function test_InvalidStalenessThreshold_Revert() public {
    PriceRegistry.StaticConfig memory staticConfig = PriceRegistry.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: 0
    });

    vm.expectRevert(PriceRegistry.InvalidStaticConfig.selector);

    s_priceRegistry = new PriceRegistryHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new PriceRegistry.TokenPriceFeedUpdate[](0),
      s_priceRegistryTokenTransferFeeConfigArgs,
      s_priceRegistryPremiumMultiplierWeiPerEthArgs,
      new PriceRegistry.DestChainConfigArgs[](0)
    );
  }

  function test_InvalidLinkTokenEqZeroAddress_Revert() public {
    PriceRegistry.StaticConfig memory staticConfig = PriceRegistry.StaticConfig({
      linkToken: address(0),
      maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });

    vm.expectRevert(PriceRegistry.InvalidStaticConfig.selector);

    s_priceRegistry = new PriceRegistryHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new PriceRegistry.TokenPriceFeedUpdate[](0),
      s_priceRegistryTokenTransferFeeConfigArgs,
      s_priceRegistryPremiumMultiplierWeiPerEthArgs,
      new PriceRegistry.DestChainConfigArgs[](0)
    );
  }

  function test_InvalidMaxFeeJuelsPerMsg_Revert() public {
    PriceRegistry.StaticConfig memory staticConfig = PriceRegistry.StaticConfig({
      linkToken: s_sourceTokens[0],
      maxFeeJuelsPerMsg: 0,
      stalenessThreshold: uint32(TWELVE_HOURS)
    });

    vm.expectRevert(PriceRegistry.InvalidStaticConfig.selector);

    s_priceRegistry = new PriceRegistryHelper(
      staticConfig,
      new address[](0),
      new address[](0),
      new PriceRegistry.TokenPriceFeedUpdate[](0),
      s_priceRegistryTokenTransferFeeConfigArgs,
      s_priceRegistryPremiumMultiplierWeiPerEthArgs,
      new PriceRegistry.DestChainConfigArgs[](0)
    );
  }
}

contract PriceRegistry_getTokenPrices is PriceRegistrySetup {
  function test_GetTokenPrices_Success() public view {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    address[] memory tokens = new address[](3);
    tokens[0] = s_sourceTokens[0];
    tokens[1] = s_sourceTokens[1];
    tokens[2] = s_weth;

    Internal.TimestampedPackedUint224[] memory tokenPrices = s_priceRegistry.getTokenPrices(tokens);

    assertEq(tokenPrices.length, 3);
    assertEq(tokenPrices[0].value, priceUpdates.tokenPriceUpdates[0].usdPerToken);
    assertEq(tokenPrices[1].value, priceUpdates.tokenPriceUpdates[1].usdPerToken);
    assertEq(tokenPrices[2].value, priceUpdates.tokenPriceUpdates[2].usdPerToken);
  }
}

contract PriceRegistry_getTokenPrice is PriceRegistrySetup {
  function test_GetTokenPriceFromFeed_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Below staleness threshold
    vm.warp(originalTimestampValue + 1 hours);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    Internal.TimestampedPackedUint224 memory tokenPriceAnswer = s_priceRegistry.getTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer.value, uint224(1e18));
    assertEq(tokenPriceAnswer.timestamp, uint32(block.timestamp));
  }
}

contract PriceRegistry_getValidatedTokenPrice is PriceRegistrySetup {
  function test_GetValidatedTokenPrice_Success() public view {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_priceRegistry.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  function test_GetValidatedTokenPriceFromFeed_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right below staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedOverStalenessPeriod_Success() public {
    uint256 originalTimestampValue = block.timestamp;

    // Right above staleness threshold
    vm.warp(originalTimestampValue + TWELVE_HOURS + 1);

    address sourceToken = _initialiseSingleTokenPriceFeed();
    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(sourceToken);

    // Price answer is 1e8 (18 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedMaxInt224Value_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max)));

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is: uint224.MAX_VALUE * (10 ** (36 - 18 - 18))
    assertEq(tokenPriceAnswer, uint224(type(uint224).max));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Below18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 6);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 6);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e6) -> expected 1e30
    assertEq(tokenPriceAnswer, uint224(1e30));
  }

  function test_GetValidatedTokenPriceFromFeedErc20Above18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 24);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 8, 1e8);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 24);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e24) -> expected 1e12
    assertEq(tokenPriceAnswer, uint224(1e12));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt18Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, 1e18);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e18) -> expected 1e18
    assertEq(tokenPriceAnswer, uint224(1e18));
  }

  function test_GetValidatedTokenPriceFromFeedFeedAt0Decimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 0);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 0, 1e31);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 0);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e31 (0 decimal token) - unit is (1e18 * 1e18 / 1e0) -> expected 1e36
    assertEq(tokenPriceAnswer, uint224(1e67));
  }

  function test_GetValidatedTokenPriceFromFeedFlippedDecimals_Success() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 20);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 20, 1e18);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 20);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    uint224 tokenPriceAnswer = s_priceRegistry.getValidatedTokenPrice(tokenAddress);

    // Price answer is 1e8 (6 decimal token) - unit is (1e18 * 1e18 / 1e20) -> expected 1e14
    assertEq(tokenPriceAnswer, uint224(1e14));
  }

  function test_StaleFeeToken_Success() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    address token = priceUpdates.tokenPriceUpdates[0].sourceToken;

    uint224 tokenPrice = s_priceRegistry.getValidatedTokenPrice(token);

    assertEq(priceUpdates.tokenPriceUpdates[0].usdPerToken, tokenPrice);
  }

  // Reverts

  function test_OverflowFeedPrice_Revert() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, int256(uint256(type(uint224).max) + 1));

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(PriceRegistry.DataFeedValueOutOfUint224Range.selector);
    s_priceRegistry.getValidatedTokenPrice(tokenAddress);
  }

  function test_UnderflowFeedPrice_Revert() public {
    address tokenAddress = _deploySourceToken("testToken", 0, 18);
    address feedAddress = _deployTokenPriceDataFeed(tokenAddress, 18, -1);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] = getSingleTokenPriceFeedUpdateStruct(tokenAddress, feedAddress, 18);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    vm.expectRevert(PriceRegistry.DataFeedValueOutOfUint224Range.selector);
    s_priceRegistry.getValidatedTokenPrice(tokenAddress);
  }

  function test_TokenNotSupported_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.getValidatedTokenPrice(DUMMY_CONTRACT_ADDRESS);
  }

  function test_TokenNotSupportedFeed_Revert() public {
    address sourceToken = _initialiseSingleTokenPriceFeed();
    MockV3Aggregator(s_dataFeedByToken[sourceToken]).updateAnswer(0);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, sourceToken));
    s_priceRegistry.getValidatedTokenPrice(sourceToken);
  }
}

contract PriceRegistry_applyFeeTokensUpdates is PriceRegistrySetup {
  function test_ApplyFeeTokensUpdates_Success() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];

    vm.expectEmit();
    emit PriceRegistry.FeeTokenAdded(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    // add same feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    vm.expectEmit();
    emit PriceRegistry.FeeTokenRemoved(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);

    // removing already removed feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);
  }

  function test_OnlyCallableByOwner_Revert() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = STRANGER;
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
  }
}

contract PriceRegistry_updatePrices is PriceRegistrySetup {
  function test_OnlyTokenPrice_Success() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    update.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(
      update.tokenPriceUpdates[0].sourceToken, update.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );

    s_priceRegistry.updatePrices(update);

    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[0]).value, update.tokenPriceUpdates[0].usdPerToken);
  }

  function test_OnlyGasPrice_Success() public {
    Internal.PriceUpdates memory update = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](1)
    });
    update.gasPriceUpdates[0] =
      Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});

    vm.expectEmit();
    emit PriceRegistry.UsdPerUnitGasUpdated(
      update.gasPriceUpdates[0].destChainSelector, update.gasPriceUpdates[0].usdPerUnitGas, block.timestamp
    );

    s_priceRegistry.updatePrices(update);

    assertEq(
      s_priceRegistry.getDestinationChainGasPrice(DEST_CHAIN_SELECTOR).value, update.gasPriceUpdates[0].usdPerUnitGas
    );
  }

  function test_UpdateMultiplePrices_Success() public {
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](3);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[1], usdPerToken: 1800e18});
    tokenPriceUpdates[2] = Internal.TokenPriceUpdate({sourceToken: address(12345), usdPerToken: 1e18});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](3);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: 2e6});
    gasPriceUpdates[1] = Internal.GasPriceUpdate({destChainSelector: SOURCE_CHAIN_SELECTOR, usdPerUnitGas: 2000e18});
    gasPriceUpdates[2] = Internal.GasPriceUpdate({destChainSelector: 12345, usdPerUnitGas: 1e18});

    Internal.PriceUpdates memory update =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: gasPriceUpdates});

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit PriceRegistry.UsdPerTokenUpdated(
        update.tokenPriceUpdates[i].sourceToken, update.tokenPriceUpdates[i].usdPerToken, block.timestamp
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      vm.expectEmit();
      emit PriceRegistry.UsdPerUnitGasUpdated(
        update.gasPriceUpdates[i].destChainSelector, update.gasPriceUpdates[i].usdPerUnitGas, block.timestamp
      );
    }

    s_priceRegistry.updatePrices(update);

    for (uint256 i = 0; i < tokenPriceUpdates.length; ++i) {
      assertEq(
        s_priceRegistry.getTokenPrice(update.tokenPriceUpdates[i].sourceToken).value, tokenPriceUpdates[i].usdPerToken
      );
    }
    for (uint256 i = 0; i < gasPriceUpdates.length; ++i) {
      assertEq(
        s_priceRegistry.getDestinationChainGasPrice(update.gasPriceUpdates[i].destChainSelector).value,
        gasPriceUpdates[i].usdPerUnitGas
      );
    }
  }

  function test_UpdatableByAuthorizedCaller_Success() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](1),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });
    priceUpdates.tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});

    // Revert when caller is not authorized
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_priceRegistry.updatePrices(priceUpdates);

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;
    vm.startPrank(OWNER);
    s_priceRegistry.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );

    // Stranger is now an authorized caller to update prices
    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(
      priceUpdates.tokenPriceUpdates[0].sourceToken, priceUpdates.tokenPriceUpdates[0].usdPerToken, block.timestamp
    );
    s_priceRegistry.updatePrices(priceUpdates);

    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[0]).value, priceUpdates.tokenPriceUpdates[0].usdPerToken);

    vm.startPrank(OWNER);
    s_priceRegistry.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: new address[](0), removedCallers: priceUpdaters})
    );

    // Revert when authorized caller is removed
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_priceRegistry.updatePrices(priceUpdates);
  }

  // Reverts

  function test_OnlyCallableByUpdater_Revert() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      gasPriceUpdates: new Internal.GasPriceUpdate[](0)
    });

    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(AuthorizedCallers.UnauthorizedCaller.selector, STRANGER));
    s_priceRegistry.updatePrices(priceUpdates);
  }
}

contract PriceRegistry_convertTokenAmount is PriceRegistrySetup {
  function test_ConvertTokenAmount_Success() public view {
    Internal.PriceUpdates memory initialPriceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));
    uint256 amount = 3e16;
    uint256 conversionRate = (uint256(initialPriceUpdates.tokenPriceUpdates[2].usdPerToken) * 1e18)
      / uint256(initialPriceUpdates.tokenPriceUpdates[0].usdPerToken);
    uint256 expected = (amount * conversionRate) / 1e18;
    assertEq(s_priceRegistry.convertTokenAmount(s_weth, amount, s_sourceTokens[0]), expected);
  }

  function test_Fuzz_ConvertTokenAmount_Success(
    uint256 feeTokenAmount,
    uint224 usdPerFeeToken,
    uint160 usdPerLinkToken,
    uint224 usdPerUnitGas
  ) public {
    vm.assume(usdPerFeeToken > 0);
    vm.assume(usdPerLinkToken > 0);
    // We bound the max fees to be at most uint96.max link.
    feeTokenAmount = bound(feeTokenAmount, 0, (uint256(type(uint96).max) * usdPerLinkToken) / usdPerFeeToken);

    address feeToken = address(1);
    address linkToken = address(2);
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = feeToken;
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](2);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: feeToken, usdPerToken: usdPerFeeToken});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: linkToken, usdPerToken: usdPerLinkToken});

    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: usdPerUnitGas});

    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: tokenPriceUpdates, gasPriceUpdates: gasPriceUpdates});

    s_priceRegistry.updatePrices(priceUpdates);

    uint256 linkFee = s_priceRegistry.convertTokenAmount(feeToken, feeTokenAmount, linkToken);
    assertEq(linkFee, (feeTokenAmount * usdPerFeeToken) / usdPerLinkToken);
  }

  // Reverts

  function test_LinkTokenNotSupported_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(DUMMY_CONTRACT_ADDRESS, 3e16, s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(s_sourceTokens[0], 3e16, DUMMY_CONTRACT_ADDRESS);
  }
}

contract PriceRegistry_getTokenAndGasPrices is PriceRegistrySetup {
  function test_GetFeeTokenAndGasPrices_Success() public view {
    (uint224 feeTokenPrice, uint224 gasPrice) =
      s_priceRegistry.getTokenAndGasPrices(s_sourceFeeToken, DEST_CHAIN_SELECTOR);

    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    assertEq(feeTokenPrice, s_sourceTokenPrices[0]);
    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function test_ZeroGasPrice_Success() public {
    uint64 zeroGasDestChainSelector = 345678;
    Internal.GasPriceUpdate[] memory gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    gasPriceUpdates[0] = Internal.GasPriceUpdate({destChainSelector: zeroGasDestChainSelector, usdPerUnitGas: 0});

    Internal.PriceUpdates memory priceUpdates =
      Internal.PriceUpdates({tokenPriceUpdates: new Internal.TokenPriceUpdate[](0), gasPriceUpdates: gasPriceUpdates});
    s_priceRegistry.updatePrices(priceUpdates);

    (, uint224 gasPrice) = s_priceRegistry.getTokenAndGasPrices(s_sourceFeeToken, zeroGasDestChainSelector);

    assertEq(gasPrice, priceUpdates.gasPriceUpdates[0].usdPerUnitGas);
  }

  function test_UnsupportedChain_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.ChainNotSupported.selector, DEST_CHAIN_SELECTOR + 1));
    s_priceRegistry.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR + 1);
  }

  function test_StaleGasPrice_Revert() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.StaleGasPrice.selector, DEST_CHAIN_SELECTOR, TWELVE_HOURS, diff)
    );
    s_priceRegistry.getTokenAndGasPrices(s_sourceTokens[0], DEST_CHAIN_SELECTOR);
  }
}

contract PriceRegistry_updateTokenPriceFeeds is PriceRegistrySetup {
  function test_ZeroFeeds_Success() public {
    Vm.Log[] memory logEntries = vm.getRecordedLogs();

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](0);
    vm.recordLogs();
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    // Verify no log emissions
    assertEq(logEntries.length, 0);
  }

  function test_SingleFeedUpdate_Success() public {
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    _assertTokenPriceFeedConfigUnconfigured(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken)
    );

    vm.expectEmit();
    emit PriceRegistry.PriceFeedPerTokenUpdated(
      tokenPriceFeedUpdates[0].sourceToken, tokenPriceFeedUpdates[0].feedConfig
    );

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
  }

  function test_MultipleFeedUpdate_Success() public {
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](2);

    for (uint256 i = 0; i < 2; ++i) {
      tokenPriceFeedUpdates[i] =
        getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[i], s_dataFeedByToken[s_sourceTokens[i]], 18);

      _assertTokenPriceFeedConfigUnconfigured(
        s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[i].sourceToken)
      );

      vm.expectEmit();
      emit PriceRegistry.PriceFeedPerTokenUpdated(
        tokenPriceFeedUpdates[i].sourceToken, tokenPriceFeedUpdates[i].feedConfig
      );
    }

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[1].sourceToken), tokenPriceFeedUpdates[1].feedConfig
    );
  }

  function test_FeedUnset_Success() public {
    Internal.TimestampedPackedUint224 memory priceQueryInitial = s_priceRegistry.getTokenPrice(s_sourceTokens[0]);
    assertFalse(priceQueryInitial.value == 0);
    assertFalse(priceQueryInitial.timestamp == 0);

    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );

    tokenPriceFeedUpdates[0].feedConfig.dataFeedAddress = address(0);
    vm.expectEmit();
    emit PriceRegistry.PriceFeedPerTokenUpdated(
      tokenPriceFeedUpdates[0].sourceToken, tokenPriceFeedUpdates[0].feedConfig
    );

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );

    // Price data should remain after a feed has been set->unset
    Internal.TimestampedPackedUint224 memory priceQueryPostUnsetFeed = s_priceRegistry.getTokenPrice(s_sourceTokens[0]);
    assertEq(priceQueryPostUnsetFeed.value, priceQueryInitial.value);
    assertEq(priceQueryPostUnsetFeed.timestamp, priceQueryInitial.timestamp);
  }

  function test_FeedNotUpdated() public {
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);

    _assertTokenPriceFeedConfigEquality(
      s_priceRegistry.getTokenPriceFeedConfig(tokenPriceFeedUpdates[0].sourceToken), tokenPriceFeedUpdates[0].feedConfig
    );
  }

  // Reverts

  function test_FeedUpdatedByNonOwner_Revert() public {
    PriceRegistry.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new PriceRegistry.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);

    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");

    s_priceRegistry.updateTokenPriceFeeds(tokenPriceFeedUpdates);
  }
}

contract PriceRegistry_applyDestChainConfigUpdates is PriceRegistrySetup {
  function test_Fuzz_applyDestChainConfigUpdates_Success(PriceRegistry.DestChainConfigArgs memory destChainConfigArgs)
    public
  {
    vm.assume(destChainConfigArgs.destChainSelector != 0);
    vm.assume(destChainConfigArgs.destChainConfig.maxPerMsgGasLimit != 0);
    destChainConfigArgs.destChainConfig.defaultTxGasLimit = uint32(
      bound(
        destChainConfigArgs.destChainConfig.defaultTxGasLimit, 1, destChainConfigArgs.destChainConfig.maxPerMsgGasLimit
      )
    );
    destChainConfigArgs.destChainConfig.defaultTokenDestBytesOverhead = uint32(
      bound(
        destChainConfigArgs.destChainConfig.defaultTokenDestBytesOverhead,
        Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES,
        type(uint32).max
      )
    );
    destChainConfigArgs.destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;

    bool isNewChain = destChainConfigArgs.destChainSelector != DEST_CHAIN_SELECTOR;

    PriceRegistry.DestChainConfigArgs[] memory newDestChainConfigArgs = new PriceRegistry.DestChainConfigArgs[](1);
    newDestChainConfigArgs[0] = destChainConfigArgs;

    if (isNewChain) {
      vm.expectEmit();
      emit PriceRegistry.DestChainAdded(destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig);
    } else {
      vm.expectEmit();
      emit PriceRegistry.DestChainConfigUpdated(
        destChainConfigArgs.destChainSelector, destChainConfigArgs.destChainConfig
      );
    }

    s_priceRegistry.applyDestChainConfigUpdates(newDestChainConfigArgs);

    _assertPriceRegistryDestChainConfigsEqual(
      destChainConfigArgs.destChainConfig, s_priceRegistry.getDestChainConfig(destChainConfigArgs.destChainSelector)
    );
  }

  function test_applyDestChainConfigUpdates_Success() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = new PriceRegistry.DestChainConfigArgs[](2);
    destChainConfigArgs[0] = _generatePriceRegistryDestChainConfigArgs()[0];
    destChainConfigArgs[0].destChainConfig.isEnabled = false;
    destChainConfigArgs[1] = _generatePriceRegistryDestChainConfigArgs()[0];
    destChainConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    vm.expectEmit();
    emit PriceRegistry.DestChainConfigUpdated(DEST_CHAIN_SELECTOR, destChainConfigArgs[0].destChainConfig);
    vm.expectEmit();
    emit PriceRegistry.DestChainAdded(DEST_CHAIN_SELECTOR + 1, destChainConfigArgs[1].destChainConfig);

    vm.recordLogs();
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    PriceRegistry.DestChainConfig memory gotDestChainConfig0 = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);
    PriceRegistry.DestChainConfig memory gotDestChainConfig1 =
      s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);

    assertEq(vm.getRecordedLogs().length, 2);
    _assertPriceRegistryDestChainConfigsEqual(destChainConfigArgs[0].destChainConfig, gotDestChainConfig0);
    _assertPriceRegistryDestChainConfigsEqual(destChainConfigArgs[1].destChainConfig, gotDestChainConfig1);
  }

  function test_applyDestChainConfigUpdatesZeroIntput_Success() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = new PriceRegistry.DestChainConfigArgs[](0);

    vm.recordLogs();
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitEqZero_Revert() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    PriceRegistry.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.defaultTxGasLimit = 0;
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_applyDestChainConfigUpdatesDefaultTxGasLimitGtMaxPerMessageGasLimit_Revert() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    PriceRegistry.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    // Allow setting to the max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit;
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    // Revert when exceeding max value
    destChainConfigArg.destChainConfig.defaultTxGasLimit = destChainConfigArg.destChainConfig.maxPerMsgGasLimit + 1;
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestChainConfigDestChainSelectorEqZero_Revert() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    PriceRegistry.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainSelector = 0;
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidDestBytesOverhead_Revert() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    PriceRegistry.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.defaultTokenDestBytesOverhead = uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.InvalidDestChainConfig.selector, DEST_CHAIN_SELECTOR));

    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
  }

  function test_InvalidChainFamilySelector_Revert() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    PriceRegistry.DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[0];

    destChainConfigArg.destChainConfig.chainFamilySelector = bytes4(uint32(1));

    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.InvalidDestChainConfig.selector, destChainConfigArg.destChainSelector)
    );
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
  }
}

contract PriceRegistry_getDataAvailabilityCost is PriceRegistrySetup {
  function test_EmptyMessageCalculatesDataAvailabilityCost_Success() public {
    uint256 dataAvailabilityCostUSD =
      s_priceRegistry.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);

    PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);

    // Test that the cost is destnation chain specific
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    destChainConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR + 1;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas =
      destChainConfig.destDataAvailabilityOverheadGas * 2;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte =
      destChainConfig.destGasPerDataAvailabilityByte * 2;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps =
      destChainConfig.destDataAvailabilityMultiplierBps * 2;
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR + 1);
    uint256 dataAvailabilityCostUSD2 =
      s_priceRegistry.getDataAvailabilityCost(DEST_CHAIN_SELECTOR + 1, USD_PER_DATA_AVAILABILITY_GAS, 0, 0, 0);
    dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES;
    expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD2);
    assertFalse(dataAvailabilityCostUSD == dataAvailabilityCostUSD2);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCost_Success() public view {
    uint256 dataAvailabilityCostUSD =
      s_priceRegistry.getDataAvailabilityCost(DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);

    uint256 dataAvailabilityLengthBytes =
      Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES + 100 + (5 * Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES_PER_TOKEN) + 50;
    uint256 dataAvailabilityGas = destChainConfig.destDataAvailabilityOverheadGas
      + destChainConfig.destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      USD_PER_DATA_AVAILABILITY_GAS * dataAvailabilityGas * destChainConfig.destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }

  function test_SimpleMessageCalculatesDataAvailabilityCostUnsupportedDestChainSelector_Success() public view {
    uint256 dataAvailabilityCostUSD =
      s_priceRegistry.getDataAvailabilityCost(0, USD_PER_DATA_AVAILABILITY_GAS, 100, 5, 50);

    assertEq(dataAvailabilityCostUSD, 0);
  }

  function test_Fuzz_ZeroDataAvailabilityGasPriceAlwaysCalculatesZeroDataAvailabilityCost_Success(
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public view {
    uint256 dataAvailabilityCostUSD = s_priceRegistry.getDataAvailabilityCost(
      DEST_CHAIN_SELECTOR, 0, messageDataLength, numberOfTokens, tokenTransferBytesOverhead
    );

    assertEq(0, dataAvailabilityCostUSD);
  }

  function test_Fuzz_CalculateDataAvailabilityCost_Success(
    uint64 destChainSelector,
    uint32 destDataAvailabilityOverheadGas,
    uint16 destGasPerDataAvailabilityByte,
    uint16 destDataAvailabilityMultiplierBps,
    uint112 dataAvailabilityGasPrice,
    uint64 messageDataLength,
    uint32 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) public {
    vm.assume(destChainSelector != 0);
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = new PriceRegistry.DestChainConfigArgs[](1);
    PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(destChainSelector);
    destChainConfigArgs[0] =
      PriceRegistry.DestChainConfigArgs({destChainSelector: destChainSelector, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityOverheadGas = destDataAvailabilityOverheadGas;
    destChainConfigArgs[0].destChainConfig.destGasPerDataAvailabilityByte = destGasPerDataAvailabilityByte;
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = destDataAvailabilityMultiplierBps;
    destChainConfigArgs[0].destChainConfig.defaultTxGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.maxPerMsgGasLimit = GAS_LIMIT;
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_EVM;
    destChainConfigArgs[0].destChainConfig.defaultTokenDestBytesOverhead = DEFAULT_TOKEN_BYTES_OVERHEAD;

    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 dataAvailabilityCostUSD = s_priceRegistry.getDataAvailabilityCost(
      destChainConfigArgs[0].destChainSelector,
      dataAvailabilityGasPrice,
      messageDataLength,
      numberOfTokens,
      tokenTransferBytesOverhead
    );

    uint256 dataAvailabilityLengthBytes = Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.ANY_2_EVM_MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    uint256 dataAvailabilityGas =
      destDataAvailabilityOverheadGas + destGasPerDataAvailabilityByte * dataAvailabilityLengthBytes;
    uint256 expectedDataAvailabilityCostUSD =
      dataAvailabilityGasPrice * dataAvailabilityGas * destDataAvailabilityMultiplierBps * 1e14;

    assertEq(expectedDataAvailabilityCostUSD, dataAvailabilityCostUSD);
  }
}

contract PriceRegistry_applyPremiumMultiplierWeiPerEthUpdates is PriceRegistrySetup {
  function test_Fuzz_applyPremiumMultiplierWeiPerEthUpdates_Success(
    PriceRegistry.PremiumMultiplierWeiPerEthArgs memory premiumMultiplierWeiPerEthArg
  ) public {
    PriceRegistry.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new PriceRegistry.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = premiumMultiplierWeiPerEthArg;

    vm.expectEmit();
    emit PriceRegistry.PremiumMultiplierWeiPerEthUpdated(
      premiumMultiplierWeiPerEthArg.token, premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth
    );

    s_priceRegistry.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArg.premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(premiumMultiplierWeiPerEthArg.token)
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesSingleToken_Success() public {
    PriceRegistry.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new PriceRegistry.PremiumMultiplierWeiPerEthArgs[](1);
    premiumMultiplierWeiPerEthArgs[0] = s_priceRegistryPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);

    vm.expectEmit();
    emit PriceRegistry.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );

    s_priceRegistry.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      s_priceRegistryPremiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesMultipleTokens_Success() public {
    PriceRegistry.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs =
      new PriceRegistry.PremiumMultiplierWeiPerEthArgs[](2);
    premiumMultiplierWeiPerEthArgs[0] = s_priceRegistryPremiumMultiplierWeiPerEthArgs[0];
    premiumMultiplierWeiPerEthArgs[0].token = vm.addr(1);
    premiumMultiplierWeiPerEthArgs[1].token = vm.addr(2);

    vm.expectEmit();
    emit PriceRegistry.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(1), premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth
    );
    vm.expectEmit();
    emit PriceRegistry.PremiumMultiplierWeiPerEthUpdated(
      vm.addr(2), premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth
    );

    s_priceRegistry.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);

    assertEq(
      premiumMultiplierWeiPerEthArgs[0].premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(vm.addr(1))
    );
    assertEq(
      premiumMultiplierWeiPerEthArgs[1].premiumMultiplierWeiPerEth,
      s_priceRegistry.getPremiumMultiplierWeiPerEth(vm.addr(2))
    );
  }

  function test_applyPremiumMultiplierWeiPerEthUpdatesZeroInput() public {
    vm.recordLogs();
    s_priceRegistry.applyPremiumMultiplierWeiPerEthUpdates(new PriceRegistry.PremiumMultiplierWeiPerEthArgs[](0));

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    PriceRegistry.PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs;
    vm.startPrank(STRANGER);

    vm.expectRevert("Only callable by owner");

    s_priceRegistry.applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
  }
}

contract PriceRegistry_applyTokenTransferFeeConfigUpdates is PriceRegistrySetup {
  function test_Fuzz_ApplyTokenTransferFeeConfig_Success(
    PriceRegistry.TokenTransferFeeConfig[2] memory tokenTransferFeeConfigs
  ) public {
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(2, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[1].destChainSelector = DEST_CHAIN_SELECTOR + 1;

    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      for (uint256 j = 0; j < tokenTransferFeeConfigs.length; ++j) {
        tokenTransferFeeConfigs[j].destBytesOverhead = uint32(
          bound(tokenTransferFeeConfigs[j].destBytesOverhead, Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, type(uint32).max)
        );
        address feeToken = s_sourceTokens[j];
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].token = feeToken;
        tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs[j].tokenTransferFeeConfig = tokenTransferFeeConfigs[j];

        vm.expectEmit();
        emit PriceRegistry.TokenTransferFeeConfigUpdated(
          tokenTransferFeeConfigArgs[i].destChainSelector, feeToken, tokenTransferFeeConfigs[j]
        );
      }
    }

    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );

    for (uint256 i = 0; i < tokenTransferFeeConfigs.length; ++i) {
      _assertTokenTransferFeeConfigEqual(
        tokenTransferFeeConfigs[i],
        s_priceRegistry.getTokenTransferFeeConfig(
          tokenTransferFeeConfigArgs[0].destChainSelector,
          tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[i].token
        )
      );
    }
  }

  function test_ApplyTokenTransferFeeConfig_Success() public {
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 2);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = PriceRegistry
      .TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: 312,
      isEnabled: true
    });
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token = address(11);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig = PriceRegistry
      .TokenTransferFeeConfig({
      minFeeUSDCents: 12,
      maxFeeUSDCents: 13,
      deciBps: 14,
      destGasOverhead: 15,
      destBytesOverhead: 394,
      isEnabled: true
    });

    vm.expectEmit();
    emit PriceRegistry.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig
    );
    vm.expectEmit();
    emit PriceRegistry.TokenTransferFeeConfigUpdated(
      tokenTransferFeeConfigArgs[0].destChainSelector,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token,
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig
    );

    PriceRegistry.TokenTransferFeeConfigRemoveArgs[] memory tokensToRemove =
      new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0);
    s_priceRegistry.applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, tokensToRemove);

    PriceRegistry.TokenTransferFeeConfig memory config0 = s_priceRegistry.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    PriceRegistry.TokenTransferFeeConfig memory config1 = s_priceRegistry.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig, config0
    );
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );

    // Remove only the first token and validate only the first token is removed
    tokensToRemove = new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](1);
    tokensToRemove[0] = PriceRegistry.TokenTransferFeeConfigRemoveArgs({
      destChainSelector: tokenTransferFeeConfigArgs[0].destChainSelector,
      token: tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    });

    vm.expectEmit();
    emit PriceRegistry.TokenTransferFeeConfigDeleted(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );

    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      new PriceRegistry.TokenTransferFeeConfigArgs[](0), tokensToRemove
    );

    config0 = s_priceRegistry.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token
    );
    config1 = s_priceRegistry.getTokenTransferFeeConfig(
      tokenTransferFeeConfigArgs[0].destChainSelector, tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].token
    );

    PriceRegistry.TokenTransferFeeConfig memory emptyConfig;

    _assertTokenTransferFeeConfigEqual(emptyConfig, config0);
    _assertTokenTransferFeeConfigEqual(
      tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig, config1
    );
  }

  function test_ApplyTokenTransferFeeZeroInput() public {
    vm.recordLogs();
    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      new PriceRegistry.TokenTransferFeeConfigArgs[](0), new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );

    assertEq(vm.getRecordedLogs().length, 0);
  }

  // Reverts

  function test_OnlyCallableByOwnerOrAdmin_Revert() public {
    vm.startPrank(STRANGER);
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs;

    vm.expectRevert("Only callable by owner");

    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }

  function test_InvalidDestBytesOverhead_Revert() public {
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = address(5);
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = PriceRegistry
      .TokenTransferFeeConfig({
      minFeeUSDCents: 6,
      maxFeeUSDCents: 7,
      deciBps: 8,
      destGasOverhead: 9,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES - 1),
      isEnabled: true
    });

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.InvalidDestBytesOverhead.selector,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token,
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destBytesOverhead
      )
    );

    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }
}

contract PriceRegistry_getTokenTransferCost is PriceRegistryFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_NoTokenTransferChargesZeroFee_Success() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function test_getTokenTransferCost_selfServeUsesDefaults_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_selfServeTokenDefaultPricing, 1000);

    // Get config to assert it isn't set
    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    assertFalse(transferFeeConfig.isEnabled);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // Assert that the default values are used
    assertEq(uint256(DEFAULT_TOKEN_FEE_USD_CENTS) * 1e16, feeUSDWei);
    assertEq(DEFAULT_TOKEN_DEST_GAS_OVERHEAD, destGasOverhead);
    assertEq(DEFAULT_TOKEN_BYTES_OVERHEAD, destBytesOverhead);
  }

  function test_SmallTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroAmountTokenTransferChargesMinFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_LargeTokenTransferChargesMaxFeeAndGas_Success() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_FeeTokenBpsFee_Success() public view {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(
      usdWei, s_priceRegistryTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_CustomTokenBpsFee_Success() public view {
    uint256 tokenAmount = 200000e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: tokenAmount});

    PriceRegistry.TokenTransferFeeConfig memory transferFeeConfig =
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = calcUSDValueFromTokenAmount(s_customTokenPrice, tokenAmount);
    uint256 bpsUSDWei = applyBpsRatio(
      usdWei, s_priceRegistryTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroFeeConfigChargesMinFee_Success() public {
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = s_sourceFeeToken;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = PriceRegistry
      .TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      isEnabled: true
    });
    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(
      configUSDCentToWei(tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.minFeeUSDCents),
      feeUSDWei
    );
    assertEq(0, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function test_Fuzz_TokenTransferFeeDuplicateTokens_Success(uint256 transfers, uint256 amount) public view {
    // It shouldn't be possible to pay materially lower fees by splitting up the transfers.
    // Note it is possible to pay higher fees since the minimum fees are added.
    PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);
    transfers = bound(transfers, 1, destChainConfig.maxNumberOfTokensPerMsg);
    // Cap amount to avoid overflow
    amount = bound(amount, 0, 1e36);
    Client.EVMTokenAmount[] memory multiple = new Client.EVMTokenAmount[](transfers);
    for (uint256 i = 0; i < transfers; ++i) {
      multiple[i] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount});
    }
    Client.EVMTokenAmount[] memory single = new Client.EVMTokenAmount[](1);
    single[0] = Client.EVMTokenAmount({token: s_sourceTokens[0], amount: amount * transfers});

    address feeToken = s_sourceRouter.getWrappedNative();

    (uint256 feeSingleUSDWei, uint32 gasOverheadSingle, uint32 bytesOverheadSingle) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, single);
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) =
      s_priceRegistry.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertGe(feeMultipleUSDWei, (feeSingleUSDWei - destChainConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function test_MixedTokenTransferFee_Success() public view {
    address[3] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative(), CUSTOM_TOKEN];
    uint224[3] memory tokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice, s_customTokenPrice];
    PriceRegistry.TokenTransferFeeConfig[3] memory tokenTransferFeeConfigs = [
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[0]),
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[1]),
      s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[2])
    ];

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](3),
      feeToken: s_sourceRouter.getWrappedNative(),
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    uint256 expectedTotalGas = 0;
    uint256 expectedTotalBytes = 0;

    // Start with small token transfers, total bps fee is lower than min token transfer fee
    for (uint256 i = 0; i < testTokens.length; ++i) {
      message.tokenAmounts[i] = Client.EVMTokenAmount({token: testTokens[i], amount: 1e14});
      PriceRegistry.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[i]);

      expectedTotalGas += tokenTransferFeeConfig.destGasOverhead == 0
        ? DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        : tokenTransferFeeConfig.destGasOverhead;
      expectedTotalBytes += tokenTransferFeeConfig.destBytesOverhead == 0
        ? DEFAULT_TOKEN_BYTES_OVERHEAD
        : tokenTransferFeeConfig.destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) = s_priceRegistry.getTokenTransferCost(
      DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts
    );

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += configUSDCentToWei(
        tokenTransferFeeConfigs[i].minFeeUSDCents == 0
          ? DEFAULT_TOKEN_FEE_USD_CENTS
          : tokenTransferFeeConfigs[i].minFeeUSDCents
      );
    }

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 1");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 1");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 1");

    // Set 1st token transfer to a meaningful amount so its bps fee is now between min and max fee
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: testTokens[0], amount: 10000e18});

    uint256 token0USDWei = applyBpsRatio(
      calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    uint256 token1USDWei = configUSDCentToWei(DEFAULT_TOKEN_FEE_USD_CENTS);

    (feeUSDWei, destGasOverhead, destBytesOverhead) = s_priceRegistry.getTokenTransferCost(
      DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts
    );
    expectedFeeUSDWei = token0USDWei + token1USDWei + configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 2");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 2");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 2");

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[2] = Client.EVMTokenAmount({token: testTokens[2], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) = s_priceRegistry.getTokenTransferCost(
      DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts
    );
    expectedFeeUSDWei = token0USDWei + token1USDWei + configUSDCentToWei(tokenTransferFeeConfigs[2].maxFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 3");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 3");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 3");
  }
}

contract PriceRegistry_getValidatedFee is PriceRegistryFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      uint64 premiumMultiplierWeiPerEth = s_priceRegistry.getPremiumMultiplierWeiPerEth(message.feeToken);
      PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_priceRegistry.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_ZeroDataAvailabilityMultiplier_Success() public {
    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = new PriceRegistry.DestChainConfigArgs[](1);
    PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);
    destChainConfigArgs[0] =
      PriceRegistry.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = 0;
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 premiumMultiplierWeiPerEth = s_priceRegistry.getPremiumMultiplierWeiPerEth(message.feeToken);

    uint256 feeAmount = s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD = (configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);

    uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD) / s_feeTokenPrice;
    assertEq(totalPriceInFeeToken, feeAmount);
  }

  function test_HighGasMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = MAX_GAS_LIMIT;
    uint256 customDataSize = MAX_DATA_SIZE;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: new bytes(customDataSize),
        tokenAmounts: new Client.EVMTokenAmount[](0),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });

      uint64 premiumMultiplierWeiPerEth = s_priceRegistry.getPremiumMultiplierWeiPerEth(message.feeToken);
      PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_priceRegistry.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_SingleTokenMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];
      PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);
      uint32 destBytesOverhead =
        s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD
        + s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) = s_priceRegistry.getTokenTransferCost(
        DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts
      );
      uint256 messageFeeUSD = (transferFeeUSD * s_priceRegistry.getPremiumMultiplierWeiPerEth(message.feeToken));
      uint256 dataAvailabilityFeeUSD = s_priceRegistry.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_MessageWithDataAndTokenTransfer_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 customGasLimit = 1_000_000;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
        receiver: abi.encode(OWNER),
        data: "",
        tokenAmounts: new Client.EVMTokenAmount[](2),
        feeToken: testTokens[i],
        extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: customGasLimit}))
      });
      uint64 premiumMultiplierWeiPerEth = s_priceRegistry.getPremiumMultiplierWeiPerEth(message.feeToken);
      PriceRegistry.DestChainConfig memory destChainConfig = s_priceRegistry.getDestChainConfig(DEST_CHAIN_SELECTOR);

      message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: 10000e18}); // feeTokenAmount
      message.tokenAmounts[1] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: 200000e18}); // customTokenAmount
      message.data = "random bits and bytes that should be factored into the cost of the message";

      uint32 tokenGasOverhead = 0;
      uint32 tokenBytesOverhead = 0;
      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        tokenGasOverhead +=
          s_priceRegistry.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destGasOverhead;
        uint32 destBytesOverhead = s_priceRegistry.getTokenTransferFeeConfig(
          DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token
        ).destBytesOverhead;
        tokenBytesOverhead += destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;
      }

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) = s_priceRegistry.getTokenTransferCost(
        DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts
      );
      uint256 messageFeeUSD = (transferFeeUSD * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_priceRegistry.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message));
    }
  }

  function test_Fuzz_EnforceOutOfOrder(bool enforce, bool allowOutOfOrderExecution) public {
    // Update config to enforce allowOutOfOrderExecution = defaultVal.
    vm.stopPrank();
    vm.startPrank(OWNER);

    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = enforce;
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG,
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: allowOutOfOrderExecution})
    );

    // If enforcement is on, only true should be allowed.
    if (enforce && !allowOutOfOrderExecution) {
      vm.expectRevert(PriceRegistry.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    }
    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_DestinationChainNotEnabled_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.DestinationChainNotEnabled.selector, DEST_CHAIN_SELECTOR + 1));
    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR + 1, _generateEmptyMessage());
  }

  function test_EnforceOutOfOrder_Revert() public {
    // Update config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);

    PriceRegistry.DestChainConfigArgs[] memory destChainConfigArgs = _generatePriceRegistryDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = true;
    s_priceRegistry.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";

    vm.expectRevert(PriceRegistry.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(PriceRegistry.UnsupportedNumberOfTokens.selector);
    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.MessageGasLimitTooHigh.selector));
    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_NotAFeeToken_Revert() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, notAFeeToken));

    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_InvalidEVMAddress_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_priceRegistry.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }
}

contract PriceRegistry_processMessageArgs is PriceRegistryFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function setUp() public virtual override {
    super.setUp();
  }

  function test_WithLinkTokenAmount_Success() public view {
    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
    ) = s_priceRegistry.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      // LINK
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS,
      ""
    );

    assertEq(msgFeeJuels, MAX_MSG_FEES_JUELS);
  }

  function test_WithConvertedTokenAmount_Success() public view {
    address feeToken = s_sourceTokens[1];
    uint256 feeTokenAmount = 10_000 gwei;
    uint256 expectedConvertedAmount = s_priceRegistry.convertTokenAmount(feeToken, feeTokenAmount, s_sourceTokens[0]);

    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
    ) = s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, feeToken, feeTokenAmount, "");

    assertEq(msgFeeJuels, expectedConvertedAmount);
  }

  function test_WithEmptyEVMExtraArgs_Success() public view {
    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs
    ) = s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, s_sourceTokens[0], 0, "");

    assertEq(isOutOfOrderExecution, false);
    assertEq(
      convertedExtraArgs, Client._argsToBytes(s_priceRegistry.parseEVMExtraArgsFromBytes("", DEST_CHAIN_SELECTOR))
    );
  }

  function test_WithEVMExtraArgsV1_Success() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 1000}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs
    ) = s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, s_sourceTokens[0], 0, extraArgs);

    assertEq(isOutOfOrderExecution, false);
    assertEq(
      convertedExtraArgs,
      Client._argsToBytes(s_priceRegistry.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  function test_WitEVMExtraArgsV2_Success() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV2({gasLimit: 0, allowOutOfOrderExecution: true}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs
    ) = s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, s_sourceTokens[0], 0, extraArgs);

    assertEq(isOutOfOrderExecution, true);
    assertEq(
      convertedExtraArgs,
      Client._argsToBytes(s_priceRegistry.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  // Reverts

  function test_MessageFeeTooHigh_Revert() public {
    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS + 1, "");
  }

  function test_InvalidExtraArgs_Revert() public {
    vm.expectRevert(PriceRegistry.InvalidExtraArgsTag.selector);

    s_priceRegistry.processMessageArgs(DEST_CHAIN_SELECTOR, s_sourceTokens[0], 0, "abcde");
  }

  function test_MalformedEVMExtraArgs_Revert() public {
    // abi.decode error
    vm.expectRevert();

    s_priceRegistry.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      abi.encodeWithSelector(Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV1({gasLimit: 100}))
    );
  }
}

contract PriceRegistry_validatePoolReturnData is PriceRegistryFeeSetup {
  function test_WithSingleToken_Success() public view {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.RampTokenAmount[] memory rampTokenAmounts = new Internal.RampTokenAmount[](1);
    rampTokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry);

    // No revert - successful
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);
  }

  function test_TokenAmountArraysMismatching_Revert() public {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.RampTokenAmount[] memory rampTokenAmounts = new Internal.RampTokenAmount[](1);
    rampTokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry);

    // Revert due to index out of bounds access
    vm.expectRevert();

    s_priceRegistry.validatePoolReturnData(
      DEST_CHAIN_SELECTOR, new Internal.RampTokenAmount[](1), new Client.EVMTokenAmount[](0)
    );
  }

  function test_SourceTokenDataTooLarge_Revert() public {
    address sourceETH = s_sourceTokens[1];

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1000;
    sourceTokenAmounts[0].token = sourceETH;

    Internal.RampTokenAmount[] memory rampTokenAmounts = new Internal.RampTokenAmount[](1);
    rampTokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry);

    // No data set, should succeed
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);

    // Set max data length, should succeed
    rampTokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES);
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);

    // Set data to max length +1, should revert
    rampTokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 1);
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.SourceTokenDataTooLarge.selector, sourceETH));
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);

    // Set token config to allow larger data
    PriceRegistry.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = PriceRegistry
      .TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });
    s_priceRegistry.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new PriceRegistry.TokenTransferFeeConfigRemoveArgs[](0)
    );

    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);

    // Set the token data larger than the configured token data, should revert
    rampTokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 32 + 1);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.SourceTokenDataTooLarge.selector, sourceETH));
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);
  }

  function test_InvalidEVMAddressDestToken_Revert() public {
    bytes memory nonEvmAddress = abi.encode(type(uint208).max);

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.RampTokenAmount[] memory rampTokenAmounts = new Internal.RampTokenAmount[](1);
    rampTokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry);
    rampTokenAmounts[0].destTokenAddress = nonEvmAddress;

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, nonEvmAddress));
    s_priceRegistry.validatePoolReturnData(DEST_CHAIN_SELECTOR, rampTokenAmounts, sourceTokenAmounts);
  }
}

contract PriceRegistry_validateDestFamilyAddress is PriceRegistrySetup {
  function test_ValidEVMAddress_Success() public view {
    bytes memory encodedAddress = abi.encode(address(10000));
    s_priceRegistry.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, encodedAddress);
  }

  function test_ValidNonEVMAddress_Success() public view {
    s_priceRegistry.validateDestFamilyAddress(bytes4(uint32(1)), abi.encode(type(uint208).max));
  }

  // Reverts

  function test_InvalidEVMAddress_Revert() public {
    bytes memory invalidAddress = abi.encode(type(uint208).max);
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_priceRegistry.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressEncodePacked_Revert() public {
    bytes memory invalidAddress = abi.encodePacked(address(234));
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_priceRegistry.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressPrecompiles_Revert() public {
    for (uint160 i = 0; i < Internal.PRECOMPILE_SPACE; ++i) {
      bytes memory invalidAddress = abi.encode(address(i));
      vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
      s_priceRegistry.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
    }

    s_priceRegistry.validateDestFamilyAddress(
      Internal.CHAIN_FAMILY_SELECTOR_EVM, abi.encode(address(uint160(Internal.PRECOMPILE_SPACE)))
    );
  }
}

contract PriceRegistry_parseEVMExtraArgsFromBytes is PriceRegistrySetup {
  PriceRegistry.DestChainConfig private s_destChainConfig;

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generatePriceRegistryDestChainConfigArgs()[0].destChainConfig;
  }

  function test_EVMExtraArgsV1_Success() public view {
    Client.EVMExtraArgsV1 memory inputArgs = Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    Client.EVMExtraArgsV2 memory expectedOutputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_priceRegistry.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)),
      abi.encode(expectedOutputArgs)
    );
  }

  function test_EVMExtraArgsV2_Success() public view {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.assertEq(
      abi.encode(s_priceRegistry.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)), abi.encode(inputArgs)
    );
  }

  function test_EVMExtraArgsDefault_Success() public view {
    Client.EVMExtraArgsV2 memory expectedOutputArgs =
      Client.EVMExtraArgsV2({gasLimit: s_destChainConfig.defaultTxGasLimit, allowOutOfOrderExecution: false});

    vm.assertEq(
      abi.encode(s_priceRegistry.parseEVMExtraArgsFromBytes("", s_destChainConfig)), abi.encode(expectedOutputArgs)
    );
  }

  // Reverts

  function test_EVMExtraArgsInvalidExtraArgsTag_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    // Invalidate selector
    inputExtraArgs[0] = bytes1(uint8(0));

    vm.expectRevert(PriceRegistry.InvalidExtraArgsTag.selector);
    s_priceRegistry.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_EVMExtraArgsEnforceOutOfOrder_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
    s_destChainConfig.enforceOutOfOrder = true;

    vm.expectRevert(PriceRegistry.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_priceRegistry.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }

  function test_EVMExtraArgsGasLimitTooHigh_Revert() public {
    Client.EVMExtraArgsV2 memory inputArgs =
      Client.EVMExtraArgsV2({gasLimit: s_destChainConfig.maxPerMsgGasLimit + 1, allowOutOfOrderExecution: true});
    bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

    vm.expectRevert(PriceRegistry.MessageGasLimitTooHigh.selector);
    s_priceRegistry.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }
}
