// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";
import {IRMN} from "../../interfaces/IRMN.sol";

import {AuthorizedCallers} from "../../../shared/access/AuthorizedCallers.sol";
import {CommitStore} from "../../CommitStore.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {RMN} from "../../RMN.sol";
import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";
import {OCR2Abstract} from "../../ocr/OCR2Abstract.sol";
import {FeeQuoterSetup} from "../feeQuoter/FeeQuoterSetup.t.sol";
import {CommitStoreHelper} from "../helpers/CommitStoreHelper.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";

contract CommitStoreSetup is FeeQuoterSetup, OCR2BaseSetup {
  CommitStoreHelper internal s_commitStore;

  function setUp() public virtual override(FeeQuoterSetup, OCR2BaseSetup) {
    FeeQuoterSetup.setUp();
    OCR2BaseSetup.setUp();

    s_commitStore = new CommitStoreHelper(
      CommitStore.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        rmnProxy: address(s_mockRMN)
      })
    );
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(s_feeQuoter)});
    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_commitStore);
    s_feeQuoter.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: priceUpdaters, removedCallers: new address[](0)})
    );
  }
}

contract CommitStoreRealRMNSetup is FeeQuoterSetup, OCR2BaseSetup {
  CommitStoreHelper internal s_commitStore;

  RMN internal s_rmn;

  address internal constant BLESS_VOTE_ADDR = address(8888);

  function setUp() public virtual override(FeeQuoterSetup, OCR2BaseSetup) {
    FeeQuoterSetup.setUp();
    OCR2BaseSetup.setUp();

    RMN.Voter[] memory voters = new RMN.Voter[](1);
    voters[0] =
      RMN.Voter({blessVoteAddr: BLESS_VOTE_ADDR, curseVoteAddr: address(9999), blessWeight: 1, curseWeight: 1});
    // Overwrite base mock rmn with real.
    s_rmn = new RMN(RMN.Config({voters: voters, blessWeightThreshold: 1, curseWeightThreshold: 1}));
    s_commitStore = new CommitStoreHelper(
      CommitStore.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        rmnProxy: address(s_rmn)
      })
    );
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(s_feeQuoter)});
    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }
}

contract CommitStore_constructor is FeeQuoterSetup, OCR2BaseSetup {
  function setUp() public virtual override(FeeQuoterSetup, OCR2BaseSetup) {
    FeeQuoterSetup.setUp();
    OCR2BaseSetup.setUp();
  }

  function test_Constructor_Success() public {
    CommitStore.StaticConfig memory staticConfig = CommitStore.StaticConfig({
      chainSelector: DEST_CHAIN_SELECTOR,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      onRamp: 0x2C44CDDdB6a900Fa2B585dd299E03D12Fa4293Bc,
      rmnProxy: address(s_mockRMN)
    });
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(s_feeQuoter)});

    vm.expectEmit();
    emit CommitStore.ConfigSet(staticConfig, dynamicConfig);

    CommitStore commitStore = new CommitStore(staticConfig);
    commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );

    CommitStore.StaticConfig memory gotStaticConfig = commitStore.getStaticConfig();

    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.sourceChainSelector, gotStaticConfig.sourceChainSelector);
    assertEq(staticConfig.onRamp, gotStaticConfig.onRamp);
    assertEq(staticConfig.rmnProxy, gotStaticConfig.rmnProxy);

    CommitStore.DynamicConfig memory gotDynamicConfig = commitStore.getDynamicConfig();

    assertEq(dynamicConfig.priceRegistry, gotDynamicConfig.priceRegistry);

    // CommitStore initial values
    assertEq(0, commitStore.getLatestPriceEpochAndRound());
    assertEq(1, commitStore.getExpectedNextSequenceNumber());
    assertEq(commitStore.typeAndVersion(), "CommitStore 1.5.0");
    assertEq(OWNER, commitStore.owner());
    assertTrue(commitStore.isUnpausedAndNotCursed());
  }
}

contract CommitStore_setMinSeqNr is CommitStoreSetup {
  function test_Fuzz_SetMinSeqNr_Success(uint64 minSeqNr) public {
    vm.expectEmit();
    emit CommitStore.SequenceNumberSet(s_commitStore.getExpectedNextSequenceNumber(), minSeqNr);

    s_commitStore.setMinSeqNr(minSeqNr);

    assertEq(s_commitStore.getExpectedNextSequenceNumber(), minSeqNr);
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setMinSeqNr(6723);
  }
}

contract CommitStore_setDynamicConfig is CommitStoreSetup {
  function test_Fuzz_SetDynamicConfig_Success(address priceRegistry) public {
    vm.assume(priceRegistry != address(0));
    CommitStore.StaticConfig memory staticConfig = s_commitStore.getStaticConfig();
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: priceRegistry});
    bytes memory onchainConfig = abi.encode(dynamicConfig);

    vm.expectEmit();
    emit CommitStore.ConfigSet(staticConfig, dynamicConfig);

    uint32 configCount = 1;

    vm.expectEmit();
    emit OCR2Abstract.ConfigSet(
      uint32(block.number),
      getBasicConfigDigest(address(s_commitStore), s_f, configCount, onchainConfig),
      configCount + 1,
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );

    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, onchainConfig, s_offchainConfigVersion, abi.encode("")
    );

    CommitStore.DynamicConfig memory gotDynamicConfig = s_commitStore.getDynamicConfig();
    assertEq(gotDynamicConfig.priceRegistry, dynamicConfig.priceRegistry);
  }

  function test_PriceEpochCleared_Success() public {
    // Set latest price epoch and round to non-zero.
    uint40 latestEpochAndRound = 1782155;
    s_commitStore.setLatestPriceEpochAndRound(latestEpochAndRound);
    assertEq(latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(1)});
    // New config should clear it.
    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
    // Assert cleared.
    assertEq(0, s_commitStore.getLatestPriceEpochAndRound());
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(23784264)});

    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }

  function test_InvalidCommitStoreConfig_Revert() public {
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(0)});

    vm.expectRevert(CommitStore.InvalidCommitStoreConfig.selector);
    s_commitStore.setOCR2Config(
      s_valid_signers, s_valid_transmitters, s_f, abi.encode(dynamicConfig), s_offchainConfigVersion, abi.encode("")
    );
  }
}

contract CommitStore_resetUnblessedRoots is CommitStoreRealRMNSetup {
  function test_ResetUnblessedRoots_Success() public {
    bytes32[] memory rootsToReset = new bytes32[](3);
    rootsToReset[0] = "1";
    rootsToReset[1] = "2";
    rootsToReset[2] = "3";

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 2),
      merkleRoot: rootsToReset[0]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(3, 4),
      merkleRoot: rootsToReset[1]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(5, 5),
      merkleRoot: rootsToReset[2]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    IRMN.TaggedRoot[] memory blessedTaggedRoots = new IRMN.TaggedRoot[](1);
    blessedTaggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_commitStore), root: rootsToReset[1]});

    vm.startPrank(BLESS_VOTE_ADDR);
    s_rmn.voteToBless(blessedTaggedRoots);

    vm.expectEmit(false, false, false, true);
    emit CommitStore.RootRemoved(rootsToReset[0]);

    vm.expectEmit(false, false, false, true);
    emit CommitStore.RootRemoved(rootsToReset[2]);

    vm.startPrank(OWNER);
    s_commitStore.resetUnblessedRoots(rootsToReset);

    assertEq(0, s_commitStore.getMerkleRoot(rootsToReset[0]));
    assertEq(BLOCK_TIME, s_commitStore.getMerkleRoot(rootsToReset[1]));
    assertEq(0, s_commitStore.getMerkleRoot(rootsToReset[2]));
  }

  // Reverts

  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    bytes32[] memory rootToReset;
    s_commitStore.resetUnblessedRoots(rootToReset);
  }
}

contract CommitStore_report is CommitStoreSetup {
  function test_ReportOnlyRootSuccess_gas() public {
    vm.pauseGasMetering();
    uint64 max1 = 931;
    bytes32 root = "Only a single root";
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, max1),
      merkleRoot: root
    });

    vm.expectEmit();
    emit CommitStore.ReportAccepted(report);

    bytes memory encodedReport = abi.encode(report);

    vm.resumeGasMetering();
    s_commitStore.report(encodedReport, ++s_latestEpochAndRound);
    vm.pauseGasMetering();

    assertEq(max1 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(block.timestamp, s_commitStore.getMerkleRoot(root));
    vm.resumeGasMetering();
  }

  function test_ReportAndPriceUpdate_Success() public {
    uint64 max1 = 12;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(1, max1),
      merkleRoot: "test #2"
    });

    vm.expectEmit();
    emit CommitStore.ReportAccepted(report);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    assertEq(max1 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function test_StaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice =
      IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report 1"
    });

    vm.expectEmit();
    emit CommitStore.ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(maxSeq + 1, maxSeq * 2),
      merkleRoot: "stale report 2"
    });

    vm.expectEmit();
    emit CommitStore.ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq * 2 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
    assertEq(
      tokenStartPrice,
      IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
  }

  function test_OnlyTokenPriceUpdates_Success() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function test_OnlyGasPriceUpdates_Success() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function test_ValidPriceUpdateThenStaleReportWithRoot_Success() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2),
      interval: CommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report"
    });

    vm.expectEmit();
    emit CommitStore.ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);

    assertEq(maxSeq + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(
      tokenPrice1, IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  // Reverts

  function test_Paused_Revert() public {
    s_commitStore.pause();
    bytes memory report;
    vm.expectRevert(CommitStore.PausedError.selector);
    s_commitStore.report(report, ++s_latestEpochAndRound);
  }

  function test_Unhealthy_Revert() public {
    s_mockRMN.setGlobalCursed(true);
    vm.expectRevert(CommitStore.CursedByRMN.selector);
    bytes memory report;
    s_commitStore.report(report, ++s_latestEpochAndRound);
  }

  function test_InvalidRootRevert() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 4),
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(CommitStore.InvalidRoot.selector);
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_InvalidInterval_Revert() public {
    CommitStore.Interval memory interval = CommitStore.Interval(2, 2);
    CommitStore.CommitReport memory report =
      CommitStore.CommitReport({priceUpdates: _getEmptyPriceUpdates(), interval: interval, merkleRoot: bytes32(0)});

    vm.expectRevert(abi.encodeWithSelector(CommitStore.InvalidInterval.selector, interval));

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_InvalidIntervalMinLargerThanMax_Revert() public {
    CommitStore.Interval memory interval = CommitStore.Interval(1, 0);
    CommitStore.CommitReport memory report =
      CommitStore.CommitReport({priceUpdates: _getEmptyPriceUpdates(), interval: interval, merkleRoot: bytes32(0)});

    vm.expectRevert(abi.encodeWithSelector(CommitStore.InvalidInterval.selector, interval));

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function test_ZeroEpochAndRound_Revert() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(CommitStore.StaleReport.selector);

    s_commitStore.report(abi.encode(report), 0);
  }

  function test_OnlyPriceUpdateStaleReport_Revert() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: bytes32(0)
    });

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    vm.expectRevert(CommitStore.StaleReport.selector);
    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
  }

  function test_RootAlreadyCommitted_Revert() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 2),
      merkleRoot: "Only a single root"
    });
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: _getEmptyPriceUpdates(),
      interval: CommitStore.Interval(3, 3),
      merkleRoot: "Only a single root"
    });

    vm.expectRevert(CommitStore.RootAlreadyCommitted.selector);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }
}

contract CommitStore_verify is CommitStoreRealRMNSetup {
  function test_NotBlessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    s_commitStore.report(
      abi.encode(
        CommitStore.CommitReport({
          priceUpdates: _getEmptyPriceUpdates(),
          interval: CommitStore.Interval(1, 2),
          merkleRoot: leaves[0]
        })
      ),
      ++s_latestEpochAndRound
    );
    bytes32[] memory proofs = new bytes32[](0);
    // We have not blessed this root, should return 0.
    uint256 timestamp = s_commitStore.verify(leaves, proofs, 0);
    assertEq(uint256(0), timestamp);
  }

  function test_Blessed_Success() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    s_commitStore.report(
      abi.encode(
        CommitStore.CommitReport({
          priceUpdates: _getEmptyPriceUpdates(),
          interval: CommitStore.Interval(1, 2),
          merkleRoot: leaves[0]
        })
      ),
      ++s_latestEpochAndRound
    );
    // Bless that root.
    IRMN.TaggedRoot[] memory taggedRoots = new IRMN.TaggedRoot[](1);
    taggedRoots[0] = IRMN.TaggedRoot({commitStore: address(s_commitStore), root: leaves[0]});
    vm.startPrank(BLESS_VOTE_ADDR);
    s_rmn.voteToBless(taggedRoots);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_commitStore.verify(leaves, proofs, 0);
    assertEq(BLOCK_TIME, timestamp);
  }

  // Reverts

  function test_Paused_Revert() public {
    s_commitStore.pause();

    bytes32[] memory hashedLeaves = new bytes32[](0);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 proofFlagBits = 0;

    vm.expectRevert(CommitStore.PausedError.selector);
    s_commitStore.verify(hashedLeaves, proofs, proofFlagBits);
  }

  function test_TooManyLeaves_Revert() public {
    bytes32[] memory leaves = new bytes32[](258);
    bytes32[] memory proofs = new bytes32[](0);

    vm.expectRevert(MerkleMultiProof.InvalidProof.selector);

    s_commitStore.verify(leaves, proofs, 0);
  }
}

contract CommitStore_isUnpausedAndRMNHealthy is CommitStoreSetup {
  function test_RMN_Success() public {
    // Test pausing
    assertFalse(s_commitStore.paused());
    assertTrue(s_commitStore.isUnpausedAndNotCursed());
    s_commitStore.pause();
    assertTrue(s_commitStore.paused());
    assertFalse(s_commitStore.isUnpausedAndNotCursed());
    s_commitStore.unpause();
    assertFalse(s_commitStore.paused());
    assertTrue(s_commitStore.isUnpausedAndNotCursed());

    // Test rmn
    s_mockRMN.setGlobalCursed(true);
    assertFalse(s_commitStore.isUnpausedAndNotCursed());
    s_mockRMN.setGlobalCursed(false);
    // TODO: also test with s_mockRMN.setChainCursed(sourceChainSelector),
    // also for other similar tests (e.g., OffRamp, OnRamp)
    assertTrue(s_commitStore.isUnpausedAndNotCursed());

    s_mockRMN.setGlobalCursed(true);
    s_commitStore.pause();
    assertFalse(s_commitStore.isUnpausedAndNotCursed());
  }
}

contract CommitStore_setLatestPriceEpochAndRound is CommitStoreSetup {
  function test_SetLatestPriceEpochAndRound_Success() public {
    uint40 latestRoundAndEpoch = 1782155;

    vm.expectEmit();
    emit CommitStore.LatestPriceEpochAndRoundSet(
      uint40(s_commitStore.getLatestPriceEpochAndRound()), latestRoundAndEpoch
    );

    s_commitStore.setLatestPriceEpochAndRound(latestRoundAndEpoch);

    assertEq(uint40(s_commitStore.getLatestPriceEpochAndRound()), latestRoundAndEpoch);
  }

  // Reverts
  function test_OnlyOwner_Revert() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setLatestPriceEpochAndRound(6723);
  }
}
