// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CallWithExactGas} from "../../../../shared/call/CallWithExactGas.sol";
import {ERC165CheckerReverting} from "../../../libraries/ERC165CheckerReverting.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OffRamp_trialExecute is OffRampSetup {
  function setUp() public virtual override {
    super.setUp();
    _setupMultipleOffRamps();
  }

  function test_trialExecute() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 1000;
    amounts[1] = 50;
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);
    IERC20 dstToken0 = IERC20(s_destTokens[0]);

    uint256 startingBalance = dstToken0.balanceOf(message.receiver);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // Check that the tokens were transferred
    assertEq(startingBalance + amounts[0], dstToken0.balanceOf(message.receiver));
  }

  function test_trialExecute_TokenHandlingErrorIsCaught() public {
    Internal.Any2EVMRampMessage memory message = _generateAny2EVMMessageWithMaybeRevertingSingleToken(1, 10);
    address destPool = s_destPoolByToken[message.tokenAmounts[0].destTokenAddress];

    IERC20 dstToken0 = IERC20(s_destTokens[0]);
    uint256 startingBalance = dstToken0.balanceOf(OWNER);

    bytes memory errorMessage = "Random token pool issue";
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, destPool, errorMessage), err);

    // Expect the balance to remain the same
    assertEq(startingBalance, dstToken0.balanceOf(OWNER));
  }

  function test_trialExecute_RateLimitError() public {
    Internal.Any2EVMRampMessage memory message = _generateAny2EVMMessageWithMaybeRevertingSingleToken(1, 10);
    address destPool = s_destPoolByToken[message.tokenAmounts[0].destTokenAddress];

    bytes memory errorMessage = abi.encodeWithSelector(RateLimiter.BucketOverfilled.selector);
    s_maybeRevertingPool.setShouldRevert(errorMessage);

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.TokenHandlingError.selector, destPool, errorMessage), err);
  }

  // TODO test actual pool exists but isn't compatible instead of just no pool
  function test_trialExecute_TokenPoolIsNotAContract() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 10000;
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageWithTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1, amounts);

    // Happy path, pool is correct
    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.SUCCESS), uint256(newState));
    assertEq("", err);

    // address 0 has no contract
    assertEq(address(0).code.length, 0);

    message.tokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(address(0)),
      destTokenAddress: address(0),
      extraData: "",
      amount: message.tokenAmounts[0].amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    message.header.messageId = _hashMessage(message, ON_RAMP_ADDRESS_1);

    // Unhappy path, no revert but marked as failed.
    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)), err);

    address notAContract = makeAddr("not_a_contract");

    message.tokenAmounts[0] = Internal.Any2EVMTokenTransfer({
      sourcePoolAddress: abi.encode(address(0)),
      destTokenAddress: notAContract,
      extraData: "",
      amount: message.tokenAmounts[0].amount,
      destGasAmount: DEFAULT_TOKEN_DEST_GAS_OVERHEAD
    });

    message.header.messageId = _hashMessage(message, ON_RAMP_ADDRESS_1);

    (newState, err) = s_offRamp.trialExecute(message, new bytes[](message.tokenAmounts.length), new uint32[](0));

    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(abi.encodeWithSelector(OffRamp.NotACompatiblePool.selector, address(0)), err);
  }

  function test_trialExecute_SenderIsNotGasEstimator_CallWithExactGasReverts() public {
    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);

    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);

    vm.mockCallRevert(
      address(s_offRamp),
      abi.encodeCall(s_offRamp.executeSingleMessage, (message, offchainTokenData, new uint32[](0))),
      abi.encodeWithSelector(CallWithExactGas.NOT_ENOUGH_GAS_FOR_CALL_SIG, "")
    );

    (Internal.MessageExecutionState newState, bytes memory err) =
      s_offRamp.trialExecute(message, offchainTokenData, new uint32[](0));
    assertEq(uint256(Internal.MessageExecutionState.FAILURE), uint256(newState));
    assertEq(CallWithExactGas.NotEnoughGasForCall.selector, bytes4(err));
  }

  function test_trialExecute_RevertsWhen_SenderIsGasEstimator_InsufficientGasForToCompleteTx() public {
    bytes4[3] memory sigs = [
      CallWithExactGas.NO_GAS_FOR_CALL_EXACT_CHECK_SIG,
      CallWithExactGas.NOT_ENOUGH_GAS_FOR_CALL_SIG,
      ERC165CheckerReverting.InsufficientGasForStaticCall.selector
    ];

    Internal.Any2EVMRampMessage memory message =
      _generateAny2EVMMessageNoTokens(SOURCE_CHAIN_SELECTOR_1, ON_RAMP_ADDRESS_1, 1);
    bytes[] memory offchainTokenData = new bytes[](message.tokenAmounts.length);

    for (uint256 i = 0; i < sigs.length; ++i) {
      vm.mockCallRevert(
        address(s_offRamp),
        abi.encodeCall(s_offRamp.executeSingleMessage, (message, offchainTokenData, new uint32[](0))),
        abi.encodeWithSelector(sigs[i])
      );

      changePrank(Internal.GAS_ESTIMATION_SENDER);

      vm.expectRevert(abi.encodeWithSelector(MultiOCR3Base.InsufficientGasToCompleteTx.selector, sigs[i]));
      s_offRamp.trialExecute(message, offchainTokenData, new uint32[](0));
    }
  }
}
