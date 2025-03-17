// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MockV3Aggregator} from "../../../shared/mocks/MockV3Aggregator.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {TokenAdminRegistry} from "../../tokenAdminRegistry/TokenAdminRegistry.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {FeeQuoterHelper} from "../helpers/FeeQuoterHelper.sol";

contract FeeQuoterSetup is TokenSetup {
  uint112 internal constant USD_PER_GAS = 1e6; // 0.001 gwei
  uint112 internal constant USD_PER_DATA_AVAILABILITY_GAS = 1e9; // 1 gwei

  address internal constant DUMMY_CONTRACT_ADDRESS = 0x1111111111111111111111111111111111111112;
  address internal constant CUSTOM_TOKEN = address(12345);
  address internal constant CUSTOM_TOKEN_2 = address(bytes20(keccak256("CUSTOM_TOKEN_2")));

  uint32 internal constant MAX_DATA_SIZE = 30_000;
  uint16 internal constant MAX_TOKENS_LENGTH = 5;
  uint32 internal constant MAX_GAS_LIMIT = 4_000_000;

  // OnRamp
  uint96 internal constant MAX_MSG_FEES_JUELS = 1_000e18;
  uint32 internal constant DEST_GAS_OVERHEAD = 300_000;
  uint8 internal constant DEST_GAS_PER_PAYLOAD_BYTE_BASE = 16;
  uint8 internal constant DEST_GAS_PER_PAYLOAD_BYTE_HIGH = 40;
  uint16 internal constant DEST_GAS_PER_PAYLOAD_BYTE_THRESHOLD = 3000;

  uint16 internal constant DEFAULT_TOKEN_FEE_USD_CENTS = 50;
  uint32 internal constant DEFAULT_TOKEN_BYTES_OVERHEAD = 32;

  // Use 16 gas per data availability byte in our tests.
  // This is an overestimation in OP stack, it ignores 4 gas per 0 byte rule.
  // Arbitrum on the other hand, does always use 16 gas per data availability byte.
  // This value may be substantially decreased after EIP 4844.
  uint16 internal constant DEST_GAS_PER_DATA_AVAILABILITY_BYTE = 16;

  // Total L1 data availability overhead estimate is 33_596 gas.
  // This value includes complete CommitStore and OffRamp call data.
  uint32 internal constant DEST_DATA_AVAILABILITY_OVERHEAD_GAS = 188 // Fixed data availability overhead in OP stack.
    + (32 * 31 + 4) * DEST_GAS_PER_DATA_AVAILABILITY_BYTE // CommitStore single-root transmission takes up about 31 slots, plus selector.
    + (32 * 34 + 4) * DEST_GAS_PER_DATA_AVAILABILITY_BYTE; // OffRamp transmission excluding EVM2EVMMessage takes up about 34 slots, plus selector.

  // Multiples of bps, or 0.0001, use 6840 to be same as OP mainnet compression factor of 0.684.
  uint16 internal constant DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS = 6840;

  uint224 internal constant CUSTOM_TOKEN_PRICE = 1e17; // $0.1 CUSTOM

  // Encode L1 gas price and L2 gas price into a packed price.
  // L1 gas price is left-shifted to the higher-order bits.
  uint224 internal constant PACKED_USD_PER_GAS =
    (uint224(USD_PER_DATA_AVAILABILITY_GAS) << Internal.GAS_PRICE_BITS) + USD_PER_GAS;

  bytes internal constant MESSAGE_RECEIVER = abi.encode(bytes20(keccak256("MESSAGE_RECEIVER")));

  FeeQuoterHelper internal s_feeQuoter;
  // Cheat to store the price updates in storage since struct arrays aren't supported.
  bytes internal s_encodedInitialPriceUpdates;
  address internal s_weth;

  address[] internal s_sourceFeeTokens;
  uint224[] internal s_sourceTokenPrices;

  FeeQuoter.PremiumMultiplierWeiPerEthArgs[] internal s_feeQuoterPremiumMultiplierWeiPerEthArgs;
  FeeQuoter.TokenTransferFeeConfigArgs[] internal s_feeQuoterTokenTransferFeeConfigArgs;

  mapping(address token => address dataFeedAddress) internal s_dataFeedByToken;

  function setUp() public virtual override {
    super.setUp();

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

    uint224[] memory destTokenPrices = new uint224[](3);
    destTokenPrices[0] = 5e18;
    destTokenPrices[1] = 2000e18;
    destTokenPrices[2] = 2000e18;

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

    Internal.PriceUpdates memory priceUpdates = _getPriceUpdatesStruct(pricedTokens, tokenPrices);
    priceUpdates.gasPriceUpdates = new Internal.GasPriceUpdate[](1);
    priceUpdates.gasPriceUpdates[0] =
      Internal.GasPriceUpdate({destChainSelector: DEST_CHAIN_SELECTOR, usdPerUnitGas: PACKED_USD_PER_GAS});

    s_encodedInitialPriceUpdates = abi.encode(priceUpdates);

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = OWNER;
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_weth;

    s_feeQuoterPremiumMultiplierWeiPerEthArgs.push(
      FeeQuoter.PremiumMultiplierWeiPerEthArgs({
        token: s_sourceFeeToken,
        premiumMultiplierWeiPerEth: 5e17 // 0.5x
      })
    );
    s_feeQuoterPremiumMultiplierWeiPerEthArgs.push(
      FeeQuoter.PremiumMultiplierWeiPerEthArgs({
        token: s_sourceRouter.getWrappedNative(),
        premiumMultiplierWeiPerEth: 2e18 // 2x
      })
    );

    s_feeQuoterTokenTransferFeeConfigArgs.push();
    s_feeQuoterTokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs.push(
      FeeQuoter.TokenTransferFeeConfigSingleTokenArgs({
        token: s_sourceFeeToken,
        tokenTransferFeeConfig: FeeQuoter.TokenTransferFeeConfig({
          minFeeUSDCents: 1_00, // 1 USD
          maxFeeUSDCents: 1000_00, // 1,000 USD
          deciBps: 2_5, // 2.5 bps, or 0.025%
          destGasOverhead: 100_000,
          destBytesOverhead: 32,
          isEnabled: true
        })
      })
    );
    s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs.push(
      FeeQuoter.TokenTransferFeeConfigSingleTokenArgs({
        token: CUSTOM_TOKEN,
        tokenTransferFeeConfig: FeeQuoter.TokenTransferFeeConfig({
          minFeeUSDCents: 2_00, // 1 USD
          maxFeeUSDCents: 2000_00, // 1,000 USD
          deciBps: 10_0, // 10 bps, or 0.1%
          destGasOverhead: 95_000,
          destBytesOverhead: 200,
          isEnabled: true
        })
      })
    );
    s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs.push(
      FeeQuoter.TokenTransferFeeConfigSingleTokenArgs({
        token: CUSTOM_TOKEN_2,
        tokenTransferFeeConfig: FeeQuoter.TokenTransferFeeConfig({
          minFeeUSDCents: 2_00, // 1 USD
          maxFeeUSDCents: 2000_00, // 1,000 USD
          deciBps: 10_0, // 10 bps, or 0.1%
          destGasOverhead: 1,
          destBytesOverhead: 200,
          isEnabled: false
        })
      })
    );

    //setting up the destination token for CUSTOM_TOKEN_2 here as it is specific to these tests
    s_destTokenBySourceToken[CUSTOM_TOKEN_2] = address(bytes20(keccak256("CUSTOM_TOKEN_2_DEST")));

    s_feeQuoter = new FeeQuoterHelper(
      FeeQuoter.StaticConfig({
        linkToken: s_sourceTokens[0],
        maxFeeJuelsPerMsg: MAX_MSG_FEES_JUELS,
        tokenPriceStalenessThreshold: uint32(TWELVE_HOURS)
      }),
      priceUpdaters,
      feeTokens,
      new FeeQuoter.TokenPriceFeedUpdate[](0),
      s_feeQuoterTokenTransferFeeConfigArgs,
      s_feeQuoterPremiumMultiplierWeiPerEthArgs,
      _generateFeeQuoterDestChainConfigArgs()
    );
    s_feeQuoter.updatePrices(priceUpdates);
  }

  function _deployTokenPriceDataFeed(address token, uint8 decimals, int256 initialAnswer) internal returns (address) {
    MockV3Aggregator dataFeed = new MockV3Aggregator(decimals, initialAnswer);
    s_dataFeedByToken[token] = address(dataFeed);
    return address(dataFeed);
  }

  function _getPriceUpdatesStruct(
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

  function _getSingleTokenPriceFeedUpdateStruct(
    address sourceToken,
    address dataFeedAddress,
    uint8 tokenDecimals
  ) internal pure returns (FeeQuoter.TokenPriceFeedUpdate memory) {
    return FeeQuoter.TokenPriceFeedUpdate({
      sourceToken: sourceToken,
      feedConfig: FeeQuoter.TokenPriceFeedConfig({
        dataFeedAddress: dataFeedAddress,
        tokenDecimals: tokenDecimals,
        isEnabled: true
      })
    });
  }

  function _initialiseSingleTokenPriceFeed() internal returns (address) {
    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates = new FeeQuoter.TokenPriceFeedUpdate[](1);
    tokenPriceFeedUpdates[0] =
      _getSingleTokenPriceFeedUpdateStruct(s_sourceTokens[0], s_dataFeedByToken[s_sourceTokens[0]], 18);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeedUpdates);
    return s_sourceTokens[0];
  }

  function _generateTokenTransferFeeConfigArgs(
    uint256 destChainSelectorLength,
    uint256 tokenLength
  ) internal pure returns (FeeQuoter.TokenTransferFeeConfigArgs[] memory) {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs =
      new FeeQuoter.TokenTransferFeeConfigArgs[](destChainSelectorLength);
    for (uint256 i = 0; i < destChainSelectorLength; ++i) {
      tokenTransferFeeConfigArgs[i].tokenTransferFeeConfigs =
        new FeeQuoter.TokenTransferFeeConfigSingleTokenArgs[](tokenLength);
    }
    return tokenTransferFeeConfigArgs;
  }

  function _generateFeeQuoterDestChainConfigArgs() internal pure returns (FeeQuoter.DestChainConfigArgs[] memory) {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] = FeeQuoter.DestChainConfigArgs({
      destChainSelector: DEST_CHAIN_SELECTOR,
      destChainConfig: FeeQuoter.DestChainConfig({
        isEnabled: true,
        maxNumberOfTokensPerMsg: MAX_TOKENS_LENGTH,
        destGasOverhead: DEST_GAS_OVERHEAD,
        destGasPerPayloadByteBase: DEST_GAS_PER_PAYLOAD_BYTE_BASE,
        destGasPerPayloadByteHigh: DEST_GAS_PER_PAYLOAD_BYTE_HIGH,
        destGasPerPayloadByteThreshold: DEST_GAS_PER_PAYLOAD_BYTE_THRESHOLD,
        destDataAvailabilityOverheadGas: DEST_DATA_AVAILABILITY_OVERHEAD_GAS,
        destGasPerDataAvailabilityByte: DEST_GAS_PER_DATA_AVAILABILITY_BYTE,
        destDataAvailabilityMultiplierBps: DEST_GAS_DATA_AVAILABILITY_MULTIPLIER_BPS,
        maxDataBytes: MAX_DATA_SIZE,
        maxPerMsgGasLimit: MAX_GAS_LIMIT,
        defaultTokenFeeUSDCents: DEFAULT_TOKEN_FEE_USD_CENTS,
        defaultTokenDestGasOverhead: DEFAULT_TOKEN_DEST_GAS_OVERHEAD,
        defaultTxGasLimit: GAS_LIMIT,
        gasMultiplierWeiPerEth: 5e17,
        networkFeeUSDCents: 1_00,
        gasPriceStalenessThreshold: uint32(TWELVE_HOURS),
        enforceOutOfOrder: false,
        chainFamilySelector: Internal.CHAIN_FAMILY_SELECTOR_EVM
      })
    });
    return destChainConfigs;
  }

  function _assertTokenPriceFeedConfigEquality(
    FeeQuoter.TokenPriceFeedConfig memory config1,
    FeeQuoter.TokenPriceFeedConfig memory config2
  ) internal pure virtual {
    assertEq(config1.dataFeedAddress, config2.dataFeedAddress);
    assertEq(config1.tokenDecimals, config2.tokenDecimals);
    assertEq(config1.isEnabled, config2.isEnabled);
  }

  function _assertTokenTransferFeeConfigEqual(
    FeeQuoter.TokenTransferFeeConfig memory a,
    FeeQuoter.TokenTransferFeeConfig memory b
  ) internal pure {
    assertEq(a.minFeeUSDCents, b.minFeeUSDCents);
    assertEq(a.maxFeeUSDCents, b.maxFeeUSDCents);
    assertEq(a.deciBps, b.deciBps);
    assertEq(a.destGasOverhead, b.destGasOverhead);
    assertEq(a.destBytesOverhead, b.destBytesOverhead);
    assertEq(a.isEnabled, b.isEnabled);
  }

  function _assertFeeQuoterDestChainConfigsEqual(
    FeeQuoter.DestChainConfig memory a,
    FeeQuoter.DestChainConfig memory b
  ) internal pure {
    assertEq(a.isEnabled, b.isEnabled);
    assertEq(a.maxNumberOfTokensPerMsg, b.maxNumberOfTokensPerMsg);
    assertEq(a.maxDataBytes, b.maxDataBytes);
    assertEq(a.maxPerMsgGasLimit, b.maxPerMsgGasLimit);
    assertEq(a.destGasOverhead, b.destGasOverhead);
    assertEq(a.destGasPerPayloadByteBase, b.destGasPerPayloadByteBase);
    assertEq(a.destGasPerPayloadByteHigh, b.destGasPerPayloadByteHigh);
    assertEq(a.destGasPerPayloadByteThreshold, b.destGasPerPayloadByteThreshold);
    assertEq(a.destDataAvailabilityOverheadGas, b.destDataAvailabilityOverheadGas);
    assertEq(a.destGasPerDataAvailabilityByte, b.destGasPerDataAvailabilityByte);
    assertEq(a.destDataAvailabilityMultiplierBps, b.destDataAvailabilityMultiplierBps);
    assertEq(a.defaultTokenFeeUSDCents, b.defaultTokenFeeUSDCents);
    assertEq(a.defaultTokenDestGasOverhead, b.defaultTokenDestGasOverhead);
    assertEq(a.defaultTxGasLimit, b.defaultTxGasLimit);
  }
}

contract FeeQuoterFeeSetup is FeeQuoterSetup {
  uint224 internal s_feeTokenPrice;
  uint224 internal s_wrappedTokenPrice;

  function setUp() public virtual override {
    super.setUp();

    s_feeTokenPrice = s_sourceTokenPrices[0];
    s_wrappedTokenPrice = s_sourceTokenPrices[2];

    s_feeQuoter.updatePrices(_getSingleTokenPriceUpdateStruct(CUSTOM_TOKEN, CUSTOM_TOKEN_PRICE));
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

  // Used to generate a message with a specific extraArgs tag for SVM
  function _generateEmptyMessage2SVM() public view returns (Client.EVM2AnyMessage memory) {
    return Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](0),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._svmArgsToBytes(
        Client.SVMExtraArgsV1({
          computeUnits: GAS_LIMIT,
          accountIsWritableBitmap: 0,
          allowOutOfOrderExecution: true,
          tokenReceiver: bytes32(0),
          accounts: new bytes32[](0)
        })
      )
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

  function _getSourceTokenData(
    Client.EVMTokenAmount memory tokenAmount,
    TokenAdminRegistry tokenAdminRegistry,
    uint64 destChainSelector
  ) internal view returns (Internal.EVM2AnyTokenTransfer memory) {
    address destToken = s_destTokenBySourceToken[tokenAmount.token];

    uint32 expectedDestGasAmount;
    FeeQuoter.TokenTransferFeeConfig memory tokenTransferFeeConfig =
      FeeQuoter(s_feeQuoter).getTokenTransferFeeConfig(destChainSelector, tokenAmount.token);
    expectedDestGasAmount =
      tokenTransferFeeConfig.isEnabled ? tokenTransferFeeConfig.destGasOverhead : DEFAULT_TOKEN_DEST_GAS_OVERHEAD;

    return Internal.EVM2AnyTokenTransfer({
      sourcePoolAddress: tokenAdminRegistry.getTokenConfig(tokenAmount.token).tokenPool,
      destTokenAddress: abi.encode(destToken),
      extraData: abi.encode(18),
      amount: tokenAmount.amount,
      destExecData: abi.encode(expectedDestGasAmount)
    });
  }

  function _configUSDCentToWei(
    uint256 usdCent
  ) internal pure returns (uint256) {
    return usdCent * 1e16;
  }
}
