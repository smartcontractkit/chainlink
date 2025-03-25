// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRouter} from "../../../interfaces/IRouter.sol";

import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

import {Vm} from "forge-std/Vm.sol";

contract OffRamp_applySourceChainConfigUpdates is OffRampSetup {
  function test_ApplyZeroUpdates() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](0);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No logs emitted
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 0);

    assertEq(s_offRamp.getSourceChainSelectors().length, 0);
  }

  function test_AddNewChain() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    OffRamp.SourceChainConfig memory expectedSourceChainConfig = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: ON_RAMP_ADDRESS_1,
      isRMNVerificationDisabled: false
    });

    vm.expectEmit();
    emit OffRamp.SourceChainSelectorAdded(SOURCE_CHAIN_SELECTOR_1);

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);
  }

  function test_ReplaceExistingChain() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].isEnabled = false;
    OffRamp.SourceChainConfig memory expectedSourceChainConfig = OffRamp.SourceChainConfig({
      router: s_destRouter,
      isEnabled: false,
      minSeqNr: 1,
      onRamp: ON_RAMP_ADDRESS_1,
      isRMNVerificationDisabled: false
    });

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(SOURCE_CHAIN_SELECTOR_1, expectedSourceChainConfig);

    vm.recordLogs();
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    // No log emitted for chain selector added (only for setting the config)
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    _assertSourceChainConfigEquality(s_offRamp.getSourceChainConfig(SOURCE_CHAIN_SELECTOR_1), expectedSourceChainConfig);

    uint256[] memory resultSourceChainSelectors = s_offRamp.getSourceChainSelectors();
    assertEq(resultSourceChainSelectors.length, 1);
  }

  function test_AddMultipleChains() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](3);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 0),
      isEnabled: true,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 1,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 1),
      isEnabled: false,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[2] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1 + 2,
      onRamp: abi.encode(ON_RAMP_ADDRESS_1, 2),
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    OffRamp.SourceChainConfig[] memory expectedSourceChainConfigs = new OffRamp.SourceChainConfig[](3);
    for (uint256 i = 0; i < 3; ++i) {
      expectedSourceChainConfigs[i] = OffRamp.SourceChainConfig({
        router: s_destRouter,
        isEnabled: sourceChainConfigs[i].isEnabled,
        minSeqNr: 1,
        onRamp: abi.encode(ON_RAMP_ADDRESS_1, i),
        isRMNVerificationDisabled: false
      });

      vm.expectEmit();
      emit OffRamp.SourceChainSelectorAdded(sourceChainConfigs[i].sourceChainSelector);

      vm.expectEmit();
      emit OffRamp.SourceChainConfigSet(sourceChainConfigs[i].sourceChainSelector, expectedSourceChainConfigs[i]);
    }

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    for (uint256 i = 0; i < 3; ++i) {
      _assertSourceChainConfigEquality(
        s_offRamp.getSourceChainConfig(sourceChainConfigs[i].sourceChainSelector), expectedSourceChainConfigs[i]
      );
    }
  }

  // Setting lower fuzz run as 256 runs was sometimes resulting in flakes.
  /// forge-config: default.fuzz.runs = 32
  /// forge-config: ccip.fuzz.runs = 32
  function testFuzz_applySourceChainConfigUpdate_Success(
    OffRamp.SourceChainConfigArgs memory sourceChainConfigArgs
  ) public {
    // Skip invalid inputs
    vm.assume(sourceChainConfigArgs.sourceChainSelector != 0);
    vm.assume(sourceChainConfigArgs.onRamp.length != 0);
    vm.assume(address(sourceChainConfigArgs.router) != address(0));
    vm.assume(keccak256(sourceChainConfigArgs.onRamp) != keccak256(abi.encode(address(0))));

    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](2);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });
    sourceChainConfigs[1] = sourceChainConfigArgs;

    // Handle cases when an update occurs
    bool isNewChain = sourceChainConfigs[1].sourceChainSelector != SOURCE_CHAIN_SELECTOR_1;
    if (!isNewChain) {
      sourceChainConfigs[1].onRamp = sourceChainConfigs[0].onRamp;
    }

    OffRamp.SourceChainConfig memory expectedSourceChainConfig = OffRamp.SourceChainConfig({
      router: sourceChainConfigArgs.router,
      isEnabled: sourceChainConfigArgs.isEnabled,
      minSeqNr: 1,
      onRamp: sourceChainConfigArgs.onRamp,
      isRMNVerificationDisabled: sourceChainConfigArgs.isRMNVerificationDisabled
    });

    if (isNewChain) {
      vm.expectEmit();
      emit OffRamp.SourceChainSelectorAdded(sourceChainConfigArgs.sourceChainSelector);
    }

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(sourceChainConfigArgs.sourceChainSelector, expectedSourceChainConfig);

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    _assertSourceChainConfigEquality(
      s_offRamp.getSourceChainConfig(sourceChainConfigArgs.sourceChainSelector), expectedSourceChainConfig
    );
  }

  function test_ReplaceExistingChainOnRamp() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS_2;

    vm.expectEmit();
    emit OffRamp.SourceChainConfigSet(
      SOURCE_CHAIN_SELECTOR_1,
      OffRamp.SourceChainConfig({
        router: s_destRouter,
        isEnabled: true,
        minSeqNr: 1,
        onRamp: ON_RAMP_ADDRESS_2,
        isRMNVerificationDisabled: false
      })
    );
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_allowNonOnRampUpdateAfterLaneIsUsed() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "test #2"
    });

    _commit(
      OffRamp.CommitReport({
        priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
        blessedMerkleRoots: roots,
        unblessedMerkleRoots: new Internal.MerkleRoot[](0),
        rmnSignatures: s_rmnSignatures
      }),
      s_latestSequenceNumber
    );

    vm.startPrank(OWNER);

    // Allow changes to the Router even after the seqNum is not 1
    assertGt(s_offRamp.getSourceChainConfig(sourceChainConfigs[0].sourceChainSelector).minSeqNr, 1);

    sourceChainConfigs[0].router = IRouter(makeAddr("newRouter"));

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  // Reverts

  function test_RevertWhen_ZeroOnRampAddress() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: new bytes(0),
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    sourceChainConfigs[0].onRamp = abi.encode(address(0));
    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_RevertWhen_RouterAddress() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: IRouter(address(0)),
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_RevertWhen_ZeroSourceChainSelector() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: 0,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    vm.expectRevert(OffRamp.ZeroChainSelectorNotAllowed.selector);
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }

  function test_RevertWhen_InvalidOnRampUpdate() public {
    OffRamp.SourceChainConfigArgs[] memory sourceChainConfigs = new OffRamp.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = OffRamp.SourceChainConfigArgs({
      router: s_destRouter,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRamp: ON_RAMP_ADDRESS_1,
      isEnabled: true,
      isRMNVerificationDisabled: false
    });

    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);

    Internal.MerkleRoot[] memory roots = new Internal.MerkleRoot[](1);
    roots[0] = Internal.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR_1,
      onRampAddress: ON_RAMP_ADDRESS_1,
      minSeqNr: 1,
      maxSeqNr: 2,
      merkleRoot: "test #2"
    });

    _commit(
      OffRamp.CommitReport({
        priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
        blessedMerkleRoots: roots,
        unblessedMerkleRoots: new Internal.MerkleRoot[](0),
        rmnSignatures: s_rmnSignatures
      }),
      s_latestSequenceNumber
    );

    vm.stopPrank();
    vm.startPrank(OWNER);

    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS_2;

    vm.expectRevert(abi.encodeWithSelector(OffRamp.InvalidOnRampUpdate.selector, SOURCE_CHAIN_SELECTOR_1));
    s_offRamp.applySourceChainConfigUpdates(sourceChainConfigs);
  }
}
