// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {USDPriceWith18Decimals} from "../../libraries/USDPriceWith18Decimals.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_processMessageArgs is FeeQuoterFeeSetup {
  using USDPriceWith18Decimals for uint224;

  function setUp() public virtual override {
    super.setUp();
  }

  function test_processMessageArgs_WithLinkTokenAmount() public view {
    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
      ,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      // LINK
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(msgFeeJuels, MAX_MSG_FEES_JUELS);
  }

  function test_processMessageArgs_WithConvertedTokenAmount() public view {
    address feeToken = s_sourceTokens[1];
    uint256 feeTokenAmount = 10_000 gwei;
    uint256 expectedConvertedAmount = s_feeQuoter.convertTokenAmount(feeToken, feeTokenAmount, s_sourceTokens[0]);

    (
      uint256 msgFeeJuels,
      /* bool isOutOfOrderExecution */
      ,
      /* bytes memory convertedExtraArgs */
      ,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      feeToken,
      feeTokenAmount,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(msgFeeJuels, expectedConvertedAmount);
  }

  function test_processMessageArgs_WithEmptyEVMExtraArgs() public view {
    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, false);
    assertEq(convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes("", DEST_CHAIN_SELECTOR)));
  }

  function test_processMessageArgs_WithEVMExtraArgsV1() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV1({gasLimit: 1000}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      extraArgs,
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, false);
    assertEq(
      convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  function test_processMessageArgs_WitEVMExtraArgsV2() public view {
    bytes memory extraArgs = Client._argsToBytes(Client.EVMExtraArgsV2({gasLimit: 0, allowOutOfOrderExecution: true}));

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      extraArgs,
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertEq(isOutOfOrderExecution, true);
    assertEq(
      convertedExtraArgs, Client._argsToBytes(s_feeQuoter.parseEVMExtraArgsFromBytes(extraArgs, DEST_CHAIN_SELECTOR))
    );
  }

  function test_processMessageArgs_WithSVMExtraArgsV1() public {
    // Apply the chain update to set the chain family selector to SOL
    FeeQuoter.DestChainConfig memory s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SVM;

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: s_destChainConfig});
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigs);

    bytes memory extraArgs = Client._svmArgsToBytes(
      Client.SVMExtraArgsV1({
        computeUnits: 0,
        accountIsWritableBitmap: 0,
        tokenReceiver: bytes32(0),
        allowOutOfOrderExecution: true,
        accounts: new bytes32[](0)
      })
    );

    (
      /* uint256 msgFeeJuels */
      ,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      /* destExecDataPerToken */
    ) = s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      extraArgs,
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );

    assertTrue(isOutOfOrderExecution);
    assertEq(
      convertedExtraArgs, Client._svmArgsToBytes(s_feeQuoter.parseSVMExtraArgsFromBytes(extraArgs, s_destChainConfig))
    );
  }

  // Reverts

  function test_RevertWhen_processMessageArgs_MessageFeeTooHigh() public {
    vm.expectRevert(
      abi.encodeWithSelector(FeeQuoter.MessageFeeTooHigh.selector, MAX_MSG_FEES_JUELS + 1, MAX_MSG_FEES_JUELS)
    );

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS + 1,
      "",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_RevertWhen_processMessageArgs_InvalidExtraArgs() public {
    vm.expectRevert(FeeQuoter.InvalidExtraArgsTag.selector);

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      "wrong extra args",
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_RevertWhen_processMessageArgs_MalformedEVMExtraArgs() public {
    // abi.decode error
    vm.expectRevert();

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      0,
      abi.encodeWithSelector(Client.EVM_EXTRA_ARGS_V2_TAG, Client.EVMExtraArgsV1({gasLimit: 100})),
      new Internal.EVM2AnyTokenTransfer[](0),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_processMessageArgs_WithCorrectPoolReturnData() public view {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](2);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];
    sourceTokenAmounts[1].amount = 1e18;
    sourceTokenAmounts[1].token = CUSTOM_TOKEN_2;

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](2);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    tokenAmounts[1] = _getSourceTokenData(sourceTokenAmounts[1], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    bytes[] memory expectedDestExecData = new bytes[](2);
    expectedDestExecData[0] = abi.encode(
      s_feeQuoterTokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig.destGasOverhead
    );
    expectedDestExecData[1] = abi.encode(DEFAULT_TOKEN_DEST_GAS_OVERHEAD); //expected return data should be abi.encoded  default as isEnabled is false

    // No revert - successful
    ( /* msgFeeJuels */ , /* isOutOfOrderExecution */, /* convertedExtraArgs */, bytes[] memory destExecData) =
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    for (uint256 i = 0; i < destExecData.length; ++i) {
      assertEq(destExecData[i], expectedDestExecData[i]);
    }
  }

  function test_RevertWhen_processMessageArgs_TokenAmountArraysMismatching() public {
    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](2);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);

    // Revert due to index out of bounds access
    vm.expectRevert();

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR,
      s_sourceTokens[0],
      MAX_MSG_FEES_JUELS,
      "",
      new Internal.EVM2AnyTokenTransfer[](1),
      new Client.EVMTokenAmount[](0)
    );
  }

  function test_RevertWhen_applyTokensTransferFeeConfigUpdates_InvalidFeeRange() public {
    address sourceETH = s_sourceTokens[1];

    // Set token config to allow larger data
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 1,
      maxFeeUSDCents: 0,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidFeeRange.selector, 1, 0));

    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );
  }

  function test_RevertWhen_processMessageArgs_SourceTokenDataTooLarge() public {
    address sourceETH = s_sourceTokens[1];

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1000;
    sourceTokenAmounts[0].token = sourceETH;

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);

    // No data set, should succeed
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set max data length, should succeed
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES);
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set data to max length +1, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set token config to allow larger data
    FeeQuoter.TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs = _generateTokenTransferFeeConfigArgs(1, 1);
    tokenTransferFeeConfigArgs[0].destChainSelector = DEST_CHAIN_SELECTOR;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].token = sourceETH;
    tokenTransferFeeConfigArgs[0].tokenTransferFeeConfigs[0].tokenTransferFeeConfig = FeeQuoter.TokenTransferFeeConfig({
      minFeeUSDCents: 0,
      maxFeeUSDCents: 1,
      deciBps: 0,
      destGasOverhead: 0,
      destBytesOverhead: uint32(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) + 32,
      isEnabled: true
    });
    s_feeQuoter.applyTokenTransferFeeConfigUpdates(
      tokenTransferFeeConfigArgs, new FeeQuoter.TokenTransferFeeConfigRemoveArgs[](0)
    );

    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );

    // Set the token data larger than the configured token data, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 32 + 1);

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );
  }

  function test_RevertWhen_processMessageArgs_InvalidEVMAddressDestToken() public {
    bytes memory nonEvmAddress = abi.encode(type(uint208).max);

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    tokenAmounts[0].destTokenAddress = nonEvmAddress;

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, nonEvmAddress));
    s_feeQuoter.processMessageArgs(
      DEST_CHAIN_SELECTOR, s_sourceTokens[0], MAX_MSG_FEES_JUELS, "", tokenAmounts, sourceTokenAmounts
    );
  }
}
