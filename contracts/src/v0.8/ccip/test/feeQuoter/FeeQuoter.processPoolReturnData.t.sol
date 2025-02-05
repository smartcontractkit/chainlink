// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";

import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";
import {Pool} from "../../libraries/Pool.sol";
import {FeeQuoterFeeSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_processPoolReturnData is FeeQuoterFeeSetup {
  function test_processPoolReturnData() public view {
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
    bytes[] memory destExecData =
      s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);

    for (uint256 i = 0; i < destExecData.length; ++i) {
      assertEq(destExecData[i], expectedDestExecData[i]);
    }
  }

  function test_processPoolReturnData_RevertWhen_SourceTokenDataTooLarge() public {
    address sourceETH = s_sourceTokens[1];

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1000;
    sourceTokenAmounts[0].token = sourceETH;

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);

    // No data set, should succeed
    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);

    // Set max data length, should succeed
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES);
    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);

    // Set data to max length +1, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 1);
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);

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

    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);

    // Set the token data larger than the configured token data, should revert
    tokenAmounts[0].extraData = new bytes(Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES + 32 + 1);

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.SourceTokenDataTooLarge.selector, sourceETH));
    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);
  }

  function test_processPoolReturnData_RevertWhen_InvalidEVMAddressDestToken() public {
    bytes memory nonEvmAddress = abi.encode(type(uint208).max);

    Client.EVMTokenAmount[] memory sourceTokenAmounts = new Client.EVMTokenAmount[](1);
    sourceTokenAmounts[0].amount = 1e18;
    sourceTokenAmounts[0].token = s_sourceTokens[0];

    Internal.EVM2AnyTokenTransfer[] memory tokenAmounts = new Internal.EVM2AnyTokenTransfer[](1);
    tokenAmounts[0] = _getSourceTokenData(sourceTokenAmounts[0], s_tokenAdminRegistry, DEST_CHAIN_SELECTOR);
    tokenAmounts[0].destTokenAddress = nonEvmAddress;

    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, nonEvmAddress));
    s_feeQuoter.processPoolReturnData(DEST_CHAIN_SELECTOR, tokenAmounts, sourceTokenAmounts);
  }
}
