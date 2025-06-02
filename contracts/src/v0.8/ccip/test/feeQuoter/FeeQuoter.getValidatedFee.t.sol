// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_getValidatedFee is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function test_getValidatedFee_EmptyMessage() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];
      uint64 premiumMultiplierWeiPerEth = s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
      FeeQuoter.DestChainConfig memory destChainConfig = s_feeQuoter.getDestChainConfig(DEST_CHAIN_SELECTOR);

      uint256 feeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD;
      uint256 gasFeeUSD = gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS;
      uint256 messageFeeUSD = _configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth;
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_getValidatedFee_ZeroDataAvailabilityMultiplier() public {
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
    uint256 gasFeeUSD = gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS;
    uint256 messageFeeUSD = _configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth;

    uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD) / s_feeTokenPrice;
    assertEq(totalPriceInFeeToken, feeAmount);
  }

  function test_getValidatedFee_HighGasMessage() public view {
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

      uint256 callDataCostHigh = (customDataSize - DEST_GAS_PER_PAYLOAD_BYTE_THRESHOLD) * DEST_GAS_PER_PAYLOAD_BYTE_HIGH
        + DEST_GAS_PER_PAYLOAD_BYTE_THRESHOLD * DEST_GAS_PER_PAYLOAD_BYTE_BASE;

      uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD + callDataCostHigh;
      uint256 gasFeeUSD = gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS;
      uint256 messageFeeUSD = _configUSDCentToWei(destChainConfig.networkFeeUSDCents) * premiumMultiplierWeiPerEth;
      uint256 dataAvailabilityFeeUSD = s_feeQuoter.getDataAvailabilityCost(
        DEST_CHAIN_SELECTOR, USD_PER_DATA_AVAILABILITY_GAS, message.data.length, message.tokenAmounts.length, 0
      );

      uint256 totalPriceInFeeToken = (gasFeeUSD + messageFeeUSD + dataAvailabilityFeeUSD) / feeTokenPrices[i];
      assertEq(totalPriceInFeeToken, feeAmount);
    }
  }

  function test_getValidatedFee_SingleTokenMessage() public view {
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

      uint256 gasUsed = GAS_LIMIT + DEST_GAS_OVERHEAD + tokenBytesOverhead * DEST_GAS_PER_PAYLOAD_BYTE_BASE
        + s_feeQuoter.getTokenTransferFeeConfig(DEST_CHAIN_SELECTOR, message.tokenAmounts[0].token).destGasOverhead;
      uint256 gasFeeUSD = gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS;
      (uint256 transferFeeUSD,,) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);
      uint256 messageFeeUSD = transferFeeUSD * s_feeQuoter.getPremiumMultiplierWeiPerEth(message.feeToken);
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

  function test_getValidatedFee_MessageWithDataAndTokenTransfer() public view {
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

      (uint256 transferFeeUSD,, uint256 tokenTransferBytesOverhead) =
        s_feeQuoter.getTokenTransferCost(DEST_CHAIN_SELECTOR, message.feeToken, feeTokenPrices[i], message.tokenAmounts);

      uint256 gasFeeUSD;

      {
        uint256 gasUsed = customGasLimit + DEST_GAS_OVERHEAD
          + (message.data.length + tokenTransferBytesOverhead) * DEST_GAS_PER_PAYLOAD_BYTE_BASE + tokenGasOverhead;

        gasFeeUSD = gasUsed * destChainConfig.gasMultiplierWeiPerEth * USD_PER_GAS;
      }

      uint256 messageFeeUSD = transferFeeUSD * premiumMultiplierWeiPerEth;

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

  function testFuzz_getValidatedFee_EnforceOutOfOrder(bool enforce, bool allowOutOfOrderExecution) public {
    // Update config to enforce allowOutOfOrderExecution = defaultVal.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = enforce;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = abi.encodeWithSelector(
      Client.GENERIC_EXTRA_ARGS_V2_TAG,
      Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT * 2, allowOutOfOrderExecution: allowOutOfOrderExecution})
    );

    // If enforcement is on, only true should be allowed.
    if (enforce && !allowOutOfOrderExecution) {
      vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    }
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_SVM() public {
    // Update config to add a Solana chain.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage2SVM();

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_Aptos() public {
    // Update config to add an Aptos chain
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_APTOS;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs =
      Client._argsToBytes(Client.GenericExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: true}));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Reverts

  function test_getValidatedFee_RevertWhen_DestinationChainNotEnabled() public {
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.DestinationChainNotEnabled.selector, DEST_CHAIN_SELECTOR + 1));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR + 1, _generateEmptyMessage());
  }

  function test_getValidatedFee_RevertWhen_EnforceOutOfOrder() public {
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

  function test_getValidatedFee_RevertWhen_MessageTooLarge() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.data = new bytes(MAX_DATA_SIZE + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageTooLarge.selector, MAX_DATA_SIZE, message.data.length));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_TooManyTokens() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 tooMany = MAX_TOKENS_LENGTH + 1;
    message.tokenAmounts = new Client.EVMTokenAmount[](tooMany);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.UnsupportedNumberOfTokens.selector, tooMany, MAX_TOKENS_LENGTH));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  // Asserts gasLimit must be <=maxGasLimit
  function test_getValidatedFee_RevertWhen_MessageGasLimitTooHigh() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: MAX_GAS_LIMIT + 1}));
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.MessageGasLimitTooHigh.selector));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_NotAFeeToken() public {
    address notAFeeToken = address(0x111111);
    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(notAFeeToken, 1);
    message.feeToken = notAFeeToken;

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.FeeTokenNotSupported.selector, notAFeeToken));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_InvalidEVMAddress() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.receiver = abi.encode(type(uint208).max);

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, message.receiver));

    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_SVMMessageWithTokenTransferAndInvalidTokenReceiver() public {
    //setup to set chainFamilySelector for SVM so that token receiver's check flow is enabled
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1);
    // replace with SVM Extra Args
    message.extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: GAS_LIMIT,
        accountIsWritableBitmap: 0,
        allowOutOfOrderExecution: true,
        tokenReceiver: bytes32(0),
        accounts: new bytes32[](0)
      })
    );
    vm.expectRevert(FeeQuoter.InvalidTokenReceiver.selector);
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_TooManySVMExtraArgsAccounts() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 maxAccounts = Client.SVM_EXTRA_ARGS_MAX_ACCOUNTS;

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1);
    message.extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: GAS_LIMIT,
        accountIsWritableBitmap: 0,
        allowOutOfOrderExecution: true,
        tokenReceiver: bytes32(uint256(1)),
        accounts: new bytes32[](maxAccounts + 1)
      })
    );
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.TooManySVMExtraArgsAccounts.selector, maxAccounts + 1, maxAccounts)
    );
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_getValidatedFee_RevertWhen_InvalidSVMExtraArgsWritableBitmap() public {
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 accounts = 4;
    uint64 wrongBitmap = uint64(1 << (accounts + 1));

    Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, 1);
    message.extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: GAS_LIMIT,
        accountIsWritableBitmap: wrongBitmap,
        allowOutOfOrderExecution: true,
        tokenReceiver: bytes32(uint256(1)),
        accounts: new bytes32[](accounts)
      })
    );
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidSVMExtraArgsWritableBitmap.selector, wrongBitmap, accounts));
    s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);
  }
}
