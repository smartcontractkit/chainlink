// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMultiCommitStore} from "../../interfaces/IMultiCommitStore.sol";
import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";
import {IRMN} from "../../interfaces/IRMN.sol";

import {MultiCommitStore} from "../../MultiCommitStore.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";
import {RMN} from "../../RMN.sol";
import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";
import {OCR2Abstract} from "../../ocr/OCR2Abstract.sol";
import {MultiCommitStoreHelper} from "../helpers/MultiCommitStoreHelper.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";

contract MultiCommitStoreSetup is PriceRegistrySetup, OCR2BaseSetup {
  MultiCommitStoreHelper internal s_multiCommitStore;

  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigs = new MultiCommitStore.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: ON_RAMP_ADDRESS
    });

    s_multiCommitStore = new MultiCommitStoreHelper(
      MultiCommitStore.StaticConfig({chainSelector: DEST_CHAIN_SELECTOR, rmnProxy: address(s_mockRMN)}),
      sourceChainConfigs
    );
    MultiCommitStore.DynamicConfig memory dynamicConfig =
      MultiCommitStore.DynamicConfig({priceRegistry: address(s_priceRegistry)});
    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_multiCommitStore);
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
  }
}

contract MultiCommitStoreRealRMNSetup is PriceRegistrySetup, OCR2BaseSetup {
  MultiCommitStoreHelper internal s_multiCommitStore;

  RMN internal s_rmn;

  address internal constant BLESS_VOTE_ADDR = address(8888);

  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    RMN.Voter[] memory voters = new RMN.Voter[](1);
    voters[0] = RMN.Voter({
      blessVoteAddr: BLESS_VOTE_ADDR,
      curseVoteAddr: address(9999),
      curseUnvoteAddr: address(19999),
      blessWeight: 1,
      curseWeight: 1
    });
    // Overwrite base mock rmn with real.
    s_rmn = new RMN(RMN.Config({voters: voters, blessWeightThreshold: 1, curseWeightThreshold: 1}));

    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigs = new MultiCommitStore.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: ON_RAMP_ADDRESS
    });

    s_multiCommitStore = new MultiCommitStoreHelper(
      MultiCommitStore.StaticConfig({chainSelector: DEST_CHAIN_SELECTOR, rmnProxy: address(s_rmn)}), sourceChainConfigs
    );
    MultiCommitStore.DynamicConfig memory dynamicConfig =
      MultiCommitStore.DynamicConfig({priceRegistry: address(s_priceRegistry)});
    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }
}

contract MultiCommitStore_constructor is PriceRegistrySetup, OCR2BaseSetup {
  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();
  }

  function test_Constructor_Success() public {
    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigs = new MultiCommitStore.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: 0x2C44CDDdB6a900Fa2B585dd299E03D12Fa4293Bc
    });
    MultiCommitStore.StaticConfig memory staticConfig =
      MultiCommitStore.StaticConfig({chainSelector: DEST_CHAIN_SELECTOR, rmnProxy: address(s_mockRMN)});
    MultiCommitStore.DynamicConfig memory dynamicConfig =
      MultiCommitStore.DynamicConfig({priceRegistry: address(s_priceRegistry)});

    vm.expectEmit();
    emit MultiCommitStore.SourceChainConfigUpdated(
      sourceChainConfigs[0].sourceChainSelector,
      IMultiCommitStore.SourceChainConfig({isEnabled: true, minSeqNr: 1, onRamp: sourceChainConfigs[0].onRamp})
    );
    vm.expectEmit();
    emit MultiCommitStore.ConfigSet(staticConfig, dynamicConfig);

    MultiCommitStore multiCommitStore = new MultiCommitStore(staticConfig, sourceChainConfigs);
    multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );

    MultiCommitStore.StaticConfig memory gotStaticConfig = multiCommitStore.getStaticConfig();

    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.rmnProxy, gotStaticConfig.rmnProxy);
    assertEq(multiCommitStore.getOnRamp(sourceChainConfigs[0].sourceChainSelector), sourceChainConfigs[0].onRamp);

    MultiCommitStore.DynamicConfig memory gotDynamicConfig = multiCommitStore.getDynamicConfig();

    assertEq(dynamicConfig.priceRegistry, gotDynamicConfig.priceRegistry);

    MultiCommitStore.SourceChainConfig memory sourceChainConfig =
      multiCommitStore.getSourceChainConfig(sourceChainConfigs[0].sourceChainSelector);

    // MultiCommitStore initial values
    assertEq(0, multiCommitStore.getLatestPriceEpochAndRound());
    assertTrue(sourceChainConfig.isEnabled);
    assertEq(sourceChainConfigs[0].onRamp, sourceChainConfig.onRamp);
    assertEq(1, sourceChainConfig.minSeqNr);
    assertEq(multiCommitStore.typeAndVersion(), "MultiCommitStore 1.6.0-dev");
    assertEq(OWNER, multiCommitStore.owner());
    assertTrue(multiCommitStore.isUnpausedAndNotCursed(sourceChainConfigs[0].sourceChainSelector));
  }

  function test_Constructor_Failure() public {
    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigs = new MultiCommitStore.SourceChainConfigArgs[](1);

    // Invalid chain selector
    sourceChainConfigs[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: ON_RAMP_ADDRESS
    });
    MultiCommitStore.StaticConfig memory staticConfig =
      MultiCommitStore.StaticConfig({chainSelector: 0, rmnProxy: address(s_mockRMN)});

    vm.expectRevert(MultiCommitStore.InvalidCommitStoreConfig.selector);
    new MultiCommitStore(staticConfig, sourceChainConfigs);

    // Invalid rmn proxy
    staticConfig.chainSelector = DEST_CHAIN_SELECTOR;
    staticConfig.rmnProxy = address(0);

    vm.expectRevert(MultiCommitStore.InvalidCommitStoreConfig.selector);
    new MultiCommitStore(staticConfig, sourceChainConfigs);

    // Invalid source chain selector
    staticConfig.rmnProxy = address(s_mockRMN);
    sourceChainConfigs[0].sourceChainSelector = 0;

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigs[0].sourceChainSelector
      )
    );
    new MultiCommitStore(staticConfig, sourceChainConfigs);

    // Invalid onRamp
    sourceChainConfigs[0].sourceChainSelector = SOURCE_CHAIN_SELECTOR + 1;
    sourceChainConfigs[0].onRamp = address(0);

    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigs[0].sourceChainSelector
      )
    );
    new MultiCommitStore(staticConfig, sourceChainConfigs);

    // Invalid minSeqNr
    sourceChainConfigs[0].sourceChainSelector = SOURCE_CHAIN_SELECTOR + 1;
    sourceChainConfigs[0].onRamp = ON_RAMP_ADDRESS;
    sourceChainConfigs[0].minSeqNr = 2;
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigs[0].sourceChainSelector
      )
    );
    new MultiCommitStore(staticConfig, sourceChainConfigs);
  }
}

contract MultiCommitStore_applySourceChainConfigUpdates is MultiCommitStoreSetup {
  function test_Fuzz_applySourceChainConfigUpdates_Success(
    MultiCommitStore.SourceChainConfigArgs memory sourceChainConfig
  ) public {
    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigs = new MultiCommitStore.SourceChainConfigArgs[](1);
    sourceChainConfigs[0] = sourceChainConfig;
    bool shouldRevert;

    if (sourceChainConfig.onRamp == address(0) || sourceChainConfig.sourceChainSelector == 0) {
      shouldRevert = true;
    } else {
      address currentOnRamp = s_multiCommitStore.getSourceChainConfig(sourceChainConfig.sourceChainSelector).onRamp;

      if (currentOnRamp == address(0)) {
        if (sourceChainConfig.minSeqNr != 1) shouldRevert = true;
      } else {
        if (currentOnRamp != sourceChainConfig.onRamp) shouldRevert = true;
      }
    }

    if (shouldRevert) {
      vm.expectRevert(
        abi.encodeWithSelector(
          MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfig.sourceChainSelector
        )
      );
      s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigs);
    } else {
      vm.expectEmit();
      emit MultiCommitStore.SourceChainConfigUpdated(
        sourceChainConfig.sourceChainSelector,
        IMultiCommitStore.SourceChainConfig({
          isEnabled: sourceChainConfig.isEnabled,
          minSeqNr: sourceChainConfig.minSeqNr,
          onRamp: sourceChainConfig.onRamp
        })
      );
      s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigs);

      MultiCommitStore.SourceChainConfig memory setSourceChainConfig =
        s_multiCommitStore.getSourceChainConfig(sourceChainConfig.sourceChainSelector);
      assertEq(sourceChainConfig.isEnabled, setSourceChainConfig.isEnabled);
      assertEq(sourceChainConfig.minSeqNr, setSourceChainConfig.minSeqNr);
      assertEq(sourceChainConfig.onRamp, setSourceChainConfig.onRamp);
    }
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigUpdate =
      new MultiCommitStore.SourceChainConfigArgs[](1);
    sourceChainConfigUpdate[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      minSeqNr: 6723,
      onRamp: DUMMY_CONTRACT_ADDRESS
    });
    vm.expectRevert("Only callable by owner");
    s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigUpdate);
  }

  function test_InvalidSourceChainConfig_Revert() public {
    MultiCommitStore.SourceChainConfigArgs[] memory sourceChainConfigUpdate =
      new MultiCommitStore.SourceChainConfigArgs[](1);
    // Set new source chain onRamp to address 0
    sourceChainConfigUpdate[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: address(0)
    });
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigUpdate[0].sourceChainSelector
      )
    );
    s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigUpdate);

    // Set new source chain minSeqNr to other than 1
    sourceChainConfigUpdate[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR + 1,
      isEnabled: true,
      minSeqNr: 2,
      onRamp: DUMMY_CONTRACT_ADDRESS
    });
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigUpdate[0].sourceChainSelector
      )
    );
    s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigUpdate);

    // Update already set onRamp
    sourceChainConfigUpdate[0] = MultiCommitStore.SourceChainConfigArgs({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      isEnabled: true,
      minSeqNr: 1,
      onRamp: DUMMY_CONTRACT_ADDRESS
    });
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.InvalidSourceChainConfig.selector, sourceChainConfigUpdate[0].sourceChainSelector
      )
    );
    s_multiCommitStore.applySourceChainConfigUpdates(sourceChainConfigUpdate);
  }
}

contract MultiCommitStore_setDynamicConfig is MultiCommitStoreSetup {
  function test_Fuzz_SetDynamicConfig_Success(address priceRegistry) public {
    vm.assume(priceRegistry != address(0));
    MultiCommitStore.StaticConfig memory staticConfig = s_multiCommitStore.getStaticConfig();
    MultiCommitStore.DynamicConfig memory dynamicConfig = MultiCommitStore.DynamicConfig({priceRegistry: priceRegistry});
    bytes memory onchainConfig = abi.encode(dynamicConfig);

    vm.expectEmit();
    emit MultiCommitStore.ConfigSet(staticConfig, dynamicConfig);

    uint32 configCount = 1;

    vm.expectEmit();
    emit OCR2Abstract.ConfigSet(
      uint32(block.number),
      getBasicConfigDigest(address(s_multiCommitStore), s_f, configCount, onchainConfig),
      configCount + 1,
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );

    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, onchainConfig, s_offchainConfigVersion, abi.encode("")
    );

    MultiCommitStore.DynamicConfig memory gotDynamicConfig = s_multiCommitStore.getDynamicConfig();
    assertEq(gotDynamicConfig.priceRegistry, dynamicConfig.priceRegistry);
  }

  function test_PriceEpochCleared_Success() public {
    // Set latest price epoch and round to non-zero.
    uint40 latestEpochAndRound = 1782155;
    s_multiCommitStore.setLatestPriceEpochAndRound(latestEpochAndRound);
    assertEq(latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());

    MultiCommitStore.DynamicConfig memory dynamicConfig = MultiCommitStore.DynamicConfig({priceRegistry: address(1)});
    // New config should clear it.
    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
    // Assert cleared.
    assertEq(0, s_multiCommitStore.getLatestPriceEpochAndRound());
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    MultiCommitStore.DynamicConfig memory dynamicConfig =
      MultiCommitStore.DynamicConfig({priceRegistry: address(23784264)});

    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }

  function test_InvalidCommitStoreConfig_Revert() public {
    MultiCommitStore.DynamicConfig memory dynamicConfig = MultiCommitStore.DynamicConfig({priceRegistry: address(0)});

    vm.expectRevert(MultiCommitStore.InvalidCommitStoreConfig.selector);
    s_multiCommitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }
}

contract MultiCommitStore_resetUnblessedRoots is MultiCommitStoreRealRMNSetup {
  function test_ResetUnblessedRoots_Success() public {
    MultiCommitStore.UnblessedRoot[] memory rootsToReset = new MultiCommitStore.UnblessedRoot[](3);
    rootsToReset[0] = MultiCommitStore.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "1"});
    rootsToReset[1] = MultiCommitStore.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "2"});
    rootsToReset[2] = MultiCommitStore.UnblessedRoot({sourceChainSelector: SOURCE_CHAIN_SELECTOR, merkleRoot: "3"});

    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](3);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: rootsToReset[0].merkleRoot
    });
    roots[1] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(3, 4),
      merkleRoot: rootsToReset[1].merkleRoot
    });
    roots[2] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(5, 5),
      merkleRoot: rootsToReset[2].merkleRoot
    });

    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    IRMN.TaggedRoot[] memory blessedTaggedRoots = new IRMN.TaggedRoot[](1);
    blessedTaggedRoots[0] =
      IRMN.TaggedRoot({commitStore: address(s_multiCommitStore), root: rootsToReset[1].merkleRoot});

    vm.startPrank(BLESS_VOTE_ADDR);
    s_rmn.voteToBless(blessedTaggedRoots);

    vm.expectEmit(false, false, false, true);
    emit MultiCommitStore.RootRemoved(rootsToReset[0].merkleRoot);

    vm.expectEmit(false, false, false, true);
    emit MultiCommitStore.RootRemoved(rootsToReset[2].merkleRoot);

    vm.startPrank(OWNER);
    s_multiCommitStore.resetUnblessedRoots(rootsToReset);

    assertEq(0, s_multiCommitStore.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[0].merkleRoot));
    assertEq(BLOCK_TIME, s_multiCommitStore.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[1].merkleRoot));
    assertEq(0, s_multiCommitStore.getMerkleRoot(SOURCE_CHAIN_SELECTOR, rootsToReset[2].merkleRoot));
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    MultiCommitStore.UnblessedRoot[] memory rootsToReset = new MultiCommitStore.UnblessedRoot[](0);
    s_multiCommitStore.resetUnblessedRoots(rootsToReset);
  }
}

contract MultiCommitStore_report is MultiCommitStoreSetup {
  function test_ReportOnlyRootSuccess_gas() public {
    vm.pauseGasMetering();
    uint64 max1 = 931;
    bytes32 root = "Only a single root";

    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, max1),
      merkleRoot: root
    });

    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectEmit();
    emit MultiCommitStore.ReportAccepted(report);

    bytes memory encodedReport = abi.encode(report);

    vm.resumeGasMetering();
    s_multiCommitStore.report(encodedReport, ++s_latestEpochAndRound);
    vm.pauseGasMetering();

    assertEq(max1 + 1, s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(block.timestamp, s_multiCommitStore.getMerkleRoot(SOURCE_CHAIN_SELECTOR, root));
    vm.resumeGasMetering();
  }

  function test_ReportAndPriceUpdate_Success() public {
    uint64 max1 = 12;

    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, max1),
      merkleRoot: "test #2"
    });

    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });

    vm.expectEmit();
    emit MultiCommitStore.ReportAccepted(report);

    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    assertEq(max1 + 1, s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());
  }

  function test_StaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice =
      IPriceRegistry(s_multiCommitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value;

    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report 1"
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectEmit();
    emit MultiCommitStore.ReportAccepted(report);
    s_multiCommitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq + 1, s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());

    report.merkleRoots[0].interval = MultiCommitStore.Interval(maxSeq + 1, maxSeq * 2);
    report.merkleRoots[0].merkleRoot = "stale report 2";
    vm.expectEmit();
    emit MultiCommitStore.ReportAccepted(report);
    s_multiCommitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq * 2 + 1, s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());
    assertEq(
      tokenStartPrice,
      IPriceRegistry(s_multiCommitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
  }

  function test_OnlyTokenPriceUpdates_Success() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](0);
    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());
  }

  function test_OnlyGasPriceUpdates_Success() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](0);
    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());
  }

  function test_ValidPriceUpdateThenStaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](0);
    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      merkleRoots: roots
    });
    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());

    roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report"
    });
    report.priceUpdates = getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2);
    report.merkleRoots = roots;

    vm.expectEmit();
    emit MultiCommitStore.ReportAccepted(report);
    s_multiCommitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq + 1, s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR).minSeqNr);
    assertEq(
      tokenPrice1,
      IPriceRegistry(s_multiCommitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
    assertEq(s_latestEpochAndRound, s_multiCommitStore.getLatestPriceEpochAndRound());
  }
  // Reverts

  function test_Paused_Revert() public {
    s_multiCommitStore.pause();
    bytes memory report;
    vm.expectRevert(MultiCommitStore.PausedError.selector);
    s_multiCommitStore.report(report, ++s_latestEpochAndRound);
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: "Only a single root"
    });

    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    vm.expectRevert(abi.encodeWithSelector(MultiCommitStore.CursedByRMN.selector, roots[0].sourceChainSelector));
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_InvalidRootRevert() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 4),
      merkleRoot: bytes32(0)
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    vm.expectRevert(MultiCommitStore.InvalidRoot.selector);
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_InvalidInterval_Revert() public {
    MultiCommitStore.Interval memory interval = MultiCommitStore.Interval(2, 2);
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: interval,
      merkleRoot: bytes32(0)
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    vm.expectRevert(
      abi.encodeWithSelector(MultiCommitStore.InvalidInterval.selector, roots[0].sourceChainSelector, interval)
    );
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_InvalidIntervalMinLargerThanMax_Revert() public {
    s_multiCommitStore.getSourceChainConfig(SOURCE_CHAIN_SELECTOR);
    MultiCommitStore.Interval memory interval = MultiCommitStore.Interval(1, 0);
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: interval,
      merkleRoot: bytes32(0)
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    vm.expectRevert(
      abi.encodeWithSelector(MultiCommitStore.InvalidInterval.selector, roots[0].sourceChainSelector, interval)
    );
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_ZeroEpochAndRound_Revert() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](0);
    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
    vm.expectRevert(MultiCommitStore.StaleReport.selector);
    s_multiCommitStore.report(abi.encode(report), 0);
  }

  function test_OnlyPriceUpdateStaleReport_Revert() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](0);
    MultiCommitStore.CommitReport memory report = MultiCommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      merkleRoots: roots
    });
    vm.expectEmit();
    emit PriceRegistry.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    vm.expectRevert(MultiCommitStore.StaleReport.selector);
    s_multiCommitStore.report(abi.encode(report), s_latestEpochAndRound);
  }

  function test_SourceChainNotEnabled_Revert() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: 0,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: "Only a single root"
    });

    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});

    vm.expectRevert(abi.encodeWithSelector(MultiCommitStore.SourceChainNotEnabled.selector, 0));
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_RootAlreadyCommitted_Revert() public {
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: "Only a single root"
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    report.merkleRoots[0].interval = MultiCommitStore.Interval(3, 3);
    vm.expectRevert(
      abi.encodeWithSelector(
        MultiCommitStore.RootAlreadyCommitted.selector, roots[0].sourceChainSelector, roots[0].merkleRoot
      )
    );
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }
}

contract MultiCommitStore_verify is MultiCommitStoreRealRMNSetup {
  function test_NotBlessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";

    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: leaves[0]
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    bytes32[] memory proofs = new bytes32[](0);
    // We have not blessed this root, should return 0.
    uint256 timestamp = s_multiCommitStore.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
    assertEq(uint256(0), timestamp);
  }

  function test_Blessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: leaves[0]
    });
    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    // Bless that root.
    IRMN.TaggedRoot[] memory taggedRoots = new IRMN.TaggedRoot[](1);
    taggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_multiCommitStore), root: leaves[0]});
    vm.startPrank(BLESS_VOTE_ADDR);
    s_rmn.voteToBless(taggedRoots);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_multiCommitStore.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
    assertEq(BLOCK_TIME, timestamp);
  }

  function test_NotBlessedWrongChainSelector_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    MultiCommitStore.MerkleRoot[] memory roots = new MultiCommitStore.MerkleRoot[](1);
    roots[0] = MultiCommitStore.MerkleRoot({
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      interval: MultiCommitStore.Interval(1, 2),
      merkleRoot: leaves[0]
    });

    MultiCommitStore.CommitReport memory report =
      MultiCommitStore.CommitReport({priceUpdates: getEmptyPriceUpdates(), merkleRoots: roots});
    s_multiCommitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    // Bless that root.
    IRMN.TaggedRoot[] memory taggedRoots = new IRMN.TaggedRoot[](1);
    taggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_multiCommitStore), root: leaves[0]});
    vm.startPrank(BLESS_VOTE_ADDR);
    s_rmn.voteToBless(taggedRoots);

    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_multiCommitStore.verify(SOURCE_CHAIN_SELECTOR + 1, leaves, proofs, 0);
    assertEq(uint256(0), timestamp);
  }

  // Reverts

  function test_Paused_Revert() public {
    s_multiCommitStore.pause();
    bytes32[] memory hashedLeaves = new bytes32[](0);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 proofFlagBits = 0;
    vm.expectRevert(MultiCommitStore.PausedError.selector);
    s_multiCommitStore.verify(SOURCE_CHAIN_SELECTOR, hashedLeaves, proofs, proofFlagBits);
  }

  function test_TooManyLeaves_Revert() public {
    bytes32[] memory leaves = new bytes32[](258);
    bytes32[] memory proofs = new bytes32[](0);
    vm.expectRevert(MerkleMultiProof.InvalidProof.selector);
    s_multiCommitStore.verify(SOURCE_CHAIN_SELECTOR, leaves, proofs, 0);
  }
}

contract MultiCommitStore_isUnpausedAndRMNHealthy is MultiCommitStoreSetup {
  function test_RMN_Success() public {
    // Test pausing
    assertFalse(s_multiCommitStore.paused());
    assertTrue(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
    s_multiCommitStore.pause();
    assertTrue(s_multiCommitStore.paused());
    assertFalse(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
    s_multiCommitStore.unpause();
    assertFalse(s_multiCommitStore.paused());
    assertTrue(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
    // Test rmn
    s_mockRMN.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    assertFalse(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
    RMN.UnvoteToCurseRecord[] memory records = new RMN.UnvoteToCurseRecord[](1);
    records[0] = RMN.UnvoteToCurseRecord({curseVoteAddr: OWNER, cursesHash: bytes32(uint256(0)), forceUnvote: true});
    s_mockRMN.ownerUnvoteToCurse(records);
    assertTrue(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
    s_mockRMN.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    s_multiCommitStore.pause();
    assertFalse(s_multiCommitStore.isUnpausedAndNotCursed(SOURCE_CHAIN_SELECTOR));
  }
}

contract MultiCommitStore_setLatestPriceEpochAndRound is MultiCommitStoreSetup {
  function test_SetLatestPriceEpochAndRound_Success() public {
    uint40 latestRoundAndEpoch = 1782155;
    s_multiCommitStore.setLatestPriceEpochAndRound(latestRoundAndEpoch);
    assertEq(s_multiCommitStore.getLatestPriceEpochAndRound(), latestRoundAndEpoch);
  }
  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_multiCommitStore.setLatestPriceEpochAndRound(6723);
  }
}
