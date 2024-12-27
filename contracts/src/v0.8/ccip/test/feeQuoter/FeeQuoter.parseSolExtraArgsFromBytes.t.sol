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
    Client.SolanaAccountMeta[] memory solAccounts = new Client.SolanaAccountMeta[](1);
    solAccounts[0] = Client.SolanaAccountMeta({pubKey: VALID_SOL_PUBKEY, isWritable: false});

    Client.SolExtraArgsV1 memory inputArgs = Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accounts: solAccounts});

    bytes memory inputExtraArgs = Client._solArgsToBytes(inputArgs);
    uint256 messageDataLengthBytes = 32;

    Client.SolExtraArgsV1 memory expectedOutputArgs =
      Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accounts: solAccounts});

    vm.assertEq(
      abi.encode(s_feeQuoter.parseSOLExtraArgsFromBytes(inputExtraArgs, s_destChainConfig, messageDataLengthBytes)),
      abi.encode(expectedOutputArgs)
    );
  }

  function test_SolExtraArgsDefault() public view {
    uint256 messageDataLengthBytes = 0;

    Client.SolExtraArgsV1 memory expectedOutputArgs = Client.SolExtraArgsV1({
      computeUnits: s_destChainConfig.defaultTxGasLimit,
      accounts: new Client.SolanaAccountMeta[](0)
    });

    vm.assertEq(
      abi.encode(s_feeQuoter.parseSOLExtraArgsFromBytes("", s_destChainConfig, messageDataLengthBytes)),
      abi.encode(expectedOutputArgs)
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

  function test_SolExtraArgsV1_RevertWhen_SolExtraArgsAccountsCannotBeZero() public {
    bytes memory extraArgs = Client._solArgsToBytes(
      Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accounts: new Client.SolanaAccountMeta[](0)})
    );

    vm.expectRevert(FeeQuoter.SolExtraArgsMustBeProvided.selector);

    s_feeQuoter.parseSOLExtraArgsFromBytes(extraArgs, s_destChainConfig, 1);
  }

  function test_SolExtraArgsV1_RevertWhen_SolAddressCannotBeWritable() public {
    Client.SolanaAccountMeta[] memory solAccounts = new Client.SolanaAccountMeta[](1);
    solAccounts[0] = Client.SolanaAccountMeta({pubKey: VALID_SOL_PUBKEY, isWritable: true});

    bytes memory extraArgs =
      Client._solArgsToBytes(Client.SolExtraArgsV1({computeUnits: GAS_LIMIT, accounts: solAccounts}));

    vm.expectRevert(FeeQuoter.FirstSolExtraArgsAddressCannotBeWritable.selector);

    s_feeQuoter.parseSOLExtraArgsFromBytes(extraArgs, s_destChainConfig, 1);
  }

  function test_SolExtraArgsV1_RevertWhen_MessageGasLimitTooHigh() public {
    Client.SolExtraArgsV1 memory inputArgs = Client.SolExtraArgsV1({
      computeUnits: s_destChainConfig.maxPerMsgGasLimit + 1,
      accounts: new Client.SolanaAccountMeta[](0)
    });

    bytes memory inputExtraArgs = Client._solArgsToBytes(inputArgs);

    // The validity of the extra Args depends on any arbitrary data being sent
    // in the message so we must simulate that here.
    uint256 messageDataLengthBytes = 0;

    vm.expectRevert(FeeQuoter.MessageGasLimitTooHigh.selector);
    s_feeQuoter.parseSOLExtraArgsFromBytes(inputExtraArgs, s_destChainConfig, messageDataLengthBytes);
  }
}
