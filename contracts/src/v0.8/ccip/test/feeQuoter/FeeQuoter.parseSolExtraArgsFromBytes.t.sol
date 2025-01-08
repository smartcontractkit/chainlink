// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Client} from "../../libraries/Client.sol";

import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_parseSolExtraArgsFromBytes is FeeQuoterSetup {
  FeeQuoter.DestChainConfig private s_destChainConfig;

  /// @dev a Valid pubkey is one that is 32 bytes long, and that's it since no other validation can be performed
  /// within the constraints of the EVM.
  bytes32 internal constant VALID_SOL_PUBKEY = keccak256("SOL_PUBKEY");

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SOL;

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigs[0] =
      FeeQuoter.DestChainConfigArgs({destChainSelector: DEST_CHAIN_SELECTOR, destChainConfig: s_destChainConfig});
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigs);
  }

  function test_SolExtraArgsV1() public view {
    bytes32[] memory solAccounts = new bytes32[](1);
    solAccounts[0] = VALID_SOL_PUBKEY;

    Client.SolExtraArgsV1 memory inputArgs =
      Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accountIsWritableBitmap: 0, accounts: solAccounts});

    bytes memory inputExtraArgs = Client._solArgsToBytes(inputArgs);

    Client.SolExtraArgsV1 memory expectedOutputArgs =
      Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accountIsWritableBitmap: 0, accounts: solAccounts});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseSOLExtraArgsFromBytes(inputExtraArgs, s_destChainConfig)),
      abi.encode(expectedOutputArgs)
    );
  }

  function test_SolExtraArgsDefault() public view {
    Client.SolExtraArgsV1 memory expectedOutputArgs = Client.SolExtraArgsV1({
      computeUnits: s_destChainConfig.defaultTxGasLimit,
      accountIsWritableBitmap: 0,
      accounts: new bytes32[](0)
    });

    vm.assertEq(
      abi.encode(s_feeQuoter.parseSOLExtraArgsFromBytes("", s_destChainConfig)), abi.encode(expectedOutputArgs)
    );
  }

  function test_parseGasLimitFromExtraArgBytes_defaultTxGasLimit() public {
    // Need to apply a chain family selector that does not have an explicit extraArgs parser available
    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = new FeeQuoter.DestChainConfigArgs[](1);
    destChainConfigArgs[0] = _generateFeeQuoterDestChainConfigArgs()[0];
    destChainConfigArgs[0].destChainConfig.isEnabled = false;
    destChainConfigArgs[0].destChainConfig.chainFamilySelector = bytes4(0xdeadbeef);

    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    uint256 defaultTxGasLimit = s_destChainConfig.defaultTxGasLimit;
    uint256 gasLimit = s_feeQuoter.parseGasLimitFromExtraArgBytes("", s_destChainConfig);
    vm.assertEq(gasLimit, defaultTxGasLimit);
  }

  // Reverts

  function test_SolExtraArgsV1_RevertWhen_MessageGasLimitTooHigh() public {
    Client.SolExtraArgsV1 memory inputArgs = Client.SolExtraArgsV1({
      computeUnits: s_destChainConfig.maxPerMsgGasLimit + 1,
      accountIsWritableBitmap: 0,
      accounts: new bytes32[](0)
    });

    bytes memory inputExtraArgs = Client._solArgsToBytes(inputArgs);

    vm.expectRevert(FeeQuoter.MessageGasLimitTooHigh.selector);
    s_feeQuoter.parseSOLExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  }
}
