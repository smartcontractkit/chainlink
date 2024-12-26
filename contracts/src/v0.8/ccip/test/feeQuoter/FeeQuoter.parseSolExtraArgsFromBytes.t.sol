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
  bytes32 internal constant VALID_SOL_PUBKEY = keccak256("valid_sol_pubkey");

  function setUp() public virtual override {
    super.setUp();
    s_destChainConfig = _generateFeeQuoterDestChainConfigArgs()[0].destChainConfig;
    s_destChainConfig.chainFamilySelector = Internal.CHAIN_FAMILY_SELECTOR_SOL;
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
      abi.encode(s_feeQuoter.parseSOLExtraArgsFromBytes("", s_destChainConfig, messageDataLengthBytes)), abi.encode(expectedOutputArgs)
    );
  }

  // Reverts

  function test_SolExtraArgsV1_RevertWhen_SolAddressCannotBeWritable() public {}

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

  // function test_RevertWhen_SolExtraArgsEnforceOutOfOrder() public {
  //   Client.EVMExtraArgsV2 memory inputArgs =
  //     Client.EVMExtraArgsV2({gasLimit: GAS_LIMIT, allowOutOfOrderExecution: false});
  //   bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);
  //   s_destChainConfig.enforceOutOfOrder = true;

  //   vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
  //   s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  // }

  // function test_RevertWhen_SolExtraArgsGasLimitTooHigh() public {
  //   Client.EVMExtraArgsV2 memory inputArgs =
  //     Client.EVMExtraArgsV2({gasLimit: s_destChainConfig.maxPerMsgGasLimit + 1, allowOutOfOrderExecution: true});
  //   bytes memory inputExtraArgs = Client._argsToBytes(inputArgs);

  //   vm.expectRevert(FeeQuoter.MessageGasLimitTooHigh.selector);
  //   s_feeQuoter.parseEVMExtraArgsFromBytes(inputExtraArgs, s_destChainConfig);
  // }
}
