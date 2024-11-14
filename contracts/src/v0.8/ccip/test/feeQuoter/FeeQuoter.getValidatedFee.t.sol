// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getValidatedFee is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage_Success() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_ZeroDataAvailabilityMultiplier_Success() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
    destChainConfigArgs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: destChainConfig});
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = 0;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);

    uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

    uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
    uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
    uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);

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

      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + customDataSize * DEST_GAS_PER_PAYLOAD_BYTE;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      uint256 messageFeeUSD = (_configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
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
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);
      uint32 destBytesOverhead =
        s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destBytesOverhead;
      uint32 tokenBytesOverhead =
        destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD
        + s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken));
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
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
      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      message.tokenAmounts[0] = Client.EVMTokenAmount({token: s_sourceFeeToken, amount: 10000e18}); // feeTokenAmount
      message.tokenAmounts[1] = Client.EVMTokenAmount({token: CUSTOM_TOKEN, amount: 200000e18}); // customTokenAmount
      message.data = "random bits and bytes that should be factored into the cost of the message";

      uint32 tokenGasOverhead = 0;
      uint32 tokenBytesOverhead = 0;
      for (uint256 j = 0; j < message.tokenAmounts.length; ++j) {
        tokenGasOverhead +=
          s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destGasOverhead;
        uint32 destBytesOverhead =
          s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[j].token).destBytesOverhead;
        tokenBytesOverhead += destBytesOverhead == 0 ? uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) : destBytesOverhead;
      }

      uint256 gasUsed =
        customGasLimit + DEST_GAS_OVERHEAD + message.data.length * DEST_GAS_PER_PAYLOAD_BYTE + tokenGasOverhead;
      uint256 gasFeeUSD = (gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS);
      (uint256 transferFeeUSD,,) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = (transferFeeUSD * premiumMultiplierWeiPerEth);
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR,
        USD_PER_DATA_AVAILABILITY_GAS,
        message.data.length,
        message.tokenAmounts.length,
        tokenBytesOverhead
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message));
    }
  }

  function testFuzz_EnforceOutOfOrder(bool enforce, bool allowOutOfOrderExecution) public {
    // Update config to enforce allowOutOfOrderExecution = defaultVal.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = enforce;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.EVM_EXTRA_ARGS_V2_TAG,
      Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: allowOutOfOrderExecution})
    );

    // If enforcement is on, only true should be allowed.
    if (enforce && !allowOutOfOrderExecution) {
      vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    }
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_DestinationChainNotEnabled_Revert() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.DestinationChainNotEnabled.selector, DEST_CHAIN_SELECTOR + 1));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR + 1, _generateEmptyMessage());
  }

  function test_EnforceOutOfOrder_Revert() public {
    // Update config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = true;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_MessageTooLarge_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_TooManyTokens_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.UnsupportedNumberOfTokens.selector, tooMany, MAX_TOKENS_LENGTH));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_MessageGasLimitTooHigh_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageGasLimitTooHigh.selector));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_NotAFeeToken_Revert() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.FeeTokenNotSupported.selector, notAFeeToken));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_InvalidEVMAddress_Revert() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }
}
