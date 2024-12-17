// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getTokenTransferCost is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  address internal s_selfServeTokenDefaultPricing = makeAddr("self-serve-token-default-pricing");

  function test_NoTokenTransferChargesZeroFee() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(0, feeUSDWei);
    assertEq(0, destGasOverhead);
    assertEq(0, destBytesOverhead);
  }

  function test_getTokenTransferCost_selfServeUsesDefaults() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_selfServeTokenDefaultPricing, 1000);

    // Get config to assert it isn't set
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    assertFalse(transferFeeConfig.isEnabled);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // Assert that the default values are used
    assertEq(uint256(DEFAULT_TOKEN_FEE_USD_CENTS) * 1e16, feeUSDWei);
    assertEq(DEFAULT_TOKEN_DEST_GAS_OVERHEAD, destGasOverhead);
    assertEq(DEFAULT_TOKEN_BYTES_OVERHEAD, destBytesOverhead);
  }

  function test_SmallTokenTransferChargesMinFeeAndGas() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1000);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroAmountTokenTransferChargesMinFeeAndGas() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 0);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.minFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_LargeTokenTransferChargesMaxFeeAndGas() public view {
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    assertEq(_configUSDCentToWei(transferFeeConfig.maxFeeUSDCents), feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_FeeTokenBpsFee() public view {
    uint256 tokenAmount = 10000e18;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = _calcUSDValueFromTokenAmount(s_feeTokenPrice, tokenAmount);
    uint256 bpsUSDWei = _applyBpsRatio(
      usdWei, s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_CustomTokenBpsFee() public view {
    uint256 tokenAmount = 200000e18;

    Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
      receiver: abi.encode(OWNER),
      data: "",
      tokenAmounts: new Client.EVMTokenAmount[](1),
      feeToken: s_sourceFeeToken,
      extraArgs: Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: GAS_LIMIT}))
    });
    message.tokenAmounts[0] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: tokenAmount});

    FeeQuoter.TokenTransferFeeConfig memory transferFeeConfig =
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token);

    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    uint256 usdWei = _calcUSDValueFromTokenAmount(CUSTOM_TOKEN_PRICE, tokenAmount);
    uint256 bpsUSDWei = _applyBpsRatio(
      usdWei, s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[1].tokenTransferFeeConfig.deciBps
    );

    assertEq(bpsUSDWei, feeUSDWei);
    assertEq(transferFeeConfig.destGasOverhead, destGasOverhead);
    assertEq(transferFeeConfig.destBytesOverhead, destBytesOverhead);
  }

  function test_ZeroFeeConfigChargesMinFee() public {
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = s_sourceFeeToken;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 0,
      maxFeeUSDCents: 1,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES),
      isEnabled: true
    });
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1e36);
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_feeTokenPrice, message.tokenAmounts);

    // if token charges 0 bps, it should cost minFee to transfer
    assertEq(
      _configUSDCentToWei(
        tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.minFeeUSDCents
      ),
      feeUSDWei
    );
    assertEq(0, destGasOverhead);
    assertEq(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES, destBytesOverhead);
  }

  function testFuzz_TokenTransferFeeDuplicateTokens_Success(uint256 transfers, uint256 amount) public view {
    // It shouldn't be possible to pay materially lower fees by splitting up the transfers.
    // Note it is possible to pay higher fees since the minimum fees are added.
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
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
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, single);
    (uint256 feeMultipleUSDWei, uint32 gasOverheadMultiple, uint32 bytesOverheadMultiple) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, feeToken, s_wrappedTokenPrice, multiple);

    // Note that there can be a rounding error once per split.
    assertGe(feeMultipleUSDWei, (feeSingleUSDWei - destChainConfig.maxNumberOfTokensPerMsg));
    assertEq(gasOverheadMultiple, gasOverheadSingle * transfers);
    assertEq(bytesOverheadMultiple, bytesOverheadSingle * transfers);
  }

  function test_MixedTokenTransferFee() public view {
    address[3] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative(), CUSTOM_TOKEN];
    uint224[3] memory tokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice, CUSTOM_TOKEN_PRICE];
    FeeQuoter.TokenTransferFeeConfig[3] memory tokenTransferFeeConfigs = [
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[0]),
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[1]),
      s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[2])
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
      FeeQuoter.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, testTokens[i]);

      expectedTotalGas += tokenTransferFeeConfig.destGasOverhead == 0
        ? DEFAULT_TOKEN_DEST_GAS_OVERHEAD
        : tokenTransferFeeConfig.destGasOverhead;
      expectedTotalBytes += tokenTransferFeeConfig.destBytesOverhead == 0
        ? DEFAULT_TOKEN_BYTES_OVERHEAD
        : tokenTransferFeeConfig.destBytesOverhead;
    }
    (uint256 feeUSDWei, uint32 destGasOverhead, uint32 destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);

    uint256 expectedFeeUSDWei = 0;
    for (uint256 i = 0; i < testTokens.length; ++i) {
      expectedFeeUSDWei += _configUSDCentToWei(
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

    uint256 token0USDWei = _applyBpsRatio(
      _calcUSDValueFromTokenAmount(tokenPrices[0], message.tokenAmounts[0].amount), tokenTransferFeeConfigs[0].deciBps
    );
    uint256 token1USDWei = _configUSDCentToWei(DEFAULT_TOKEN_FEE_USD_CENTS);

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + _configUSDCentToWei(tokenTransferFeeConfigs[2].minFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 2");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 2");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 2");

    // Set 2nd token transfer to a large amount that is higher than maxFeeUSD
    message.tokenAmounts[2] = Client.EVMTokenAmount({token: testTokens[2], amount: 1e36});

    (feeUSDWei, destGasOverhead, destBytesOverhead) =
      s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, s_wrappedTokenPrice, message.tokenAmounts);
    expectedFeeUSDWei = token0USDWei + token1USDWei + _configUSDCentToWei(tokenTransferFeeConfigs[2].maxFeeUSDCents);

    assertEq(expectedFeeUSDWei, feeUSDWei, "wrong feeUSDWei 3");
    assertEq(expectedTotalGas, destGasOverhead, "wrong destGasOverhead 3");
    assertEq(expectedTotalBytes, destBytesOverhead, "wrong destBytesOverhead 3");
  }

  function _applyBpsRatio(uint256 tokenAmount, uint16 ratio) internal pure returns (uint256) {
    return (tokenAmount * ratio) / 1e5;
  }

  function _calcUSDValueFromTokenAmount(uint224 tokenPrice, uint256 tokenAmount) internal pure returns (uint256) {
    return (tokenPrice * tokenAmount) / 1e18;
  }
}
