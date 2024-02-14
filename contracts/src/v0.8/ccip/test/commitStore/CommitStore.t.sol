// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IARM} from "../../interfaces/IARM.sol";
import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";

import {CommitStore} from "../../CommitStore.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";
import {ARM} from "../../ARM.sol";
import {MerkleMultiProof} from "../../libraries/MerkleMultiProof.sol";

import {CommitStoreHelper} from "../helpers/CommitStoreHelper.sol";
import {PriceRegistrySetup} from "../priceRegistry/PriceRegistry.t.sol";
import {OCR2BaseSetup} from "../ocr/OCR2Base.t.sol";

contract CommitStoreSetup is PriceRegistrySetup, OCR2BaseSetup {
  event ConfigSet(CommitStore.StaticConfig, CommitStore.DynamicConfig);

  CommitStoreHelper internal s_commitStore;

  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    s_commitStore = new CommitStoreHelper(
      CommitStore.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        armProxy: address(s_mockARM)
      })
    );
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({
      priceRegistry: address(s_priceRegistry)
    });
    s_commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );

    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = address(s_commitStore);
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
  }
}

contract CommitStoreRealARMSetup is PriceRegistrySetup, OCR2BaseSetup {
  CommitStoreHelper internal s_commitStore;

  ARM internal s_arm;

  address internal constant BLESS_VOTE_ADDR = address(8888);

  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();

    ARM.Voter[] memory voters = new ARM.Voter[](1);
    voters[0] = ARM.Voter({
      blessVoteAddr: BLESS_VOTE_ADDR,
      curseVoteAddr: address(9999),
      curseUnvoteAddr: address(19999),
      blessWeight: 1,
      curseWeight: 1
    });
    // Overwrite base mock arm with real.
    s_arm = new ARM(ARM.Config({voters: voters, blessWeightThreshold: 1, curseWeightThreshold: 1}));
    s_commitStore = new CommitStoreHelper(
      CommitStore.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        sourceChainSelector: SOURCE_CHAIN_SELECTOR,
        onRamp: ON_RAMP_ADDRESS,
        armProxy: address(s_arm)
      })
    );
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({
      priceRegistry: address(s_priceRegistry)
    });
    s_commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
  }
}

/// @notice #constructor
contract CommitStore_constructor is PriceRegistrySetup, OCR2BaseSetup {
  event ConfigSet(CommitStore.StaticConfig, CommitStore.DynamicConfig);

  function setUp() public virtual override(PriceRegistrySetup, OCR2BaseSetup) {
    PriceRegistrySetup.setUp();
    OCR2BaseSetup.setUp();
  }

  function testConstructorSuccess() public {
    CommitStore.StaticConfig memory staticConfig = CommitStore.StaticConfig({
      chainSelector: DEST_CHAIN_SELECTOR,
      sourceChainSelector: SOURCE_CHAIN_SELECTOR,
      onRamp: 0x2C44CDDdB6a900Fa2B585dd299E03D12Fa4293Bc,
      armProxy: address(s_mockARM)
    });
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({
      priceRegistry: address(s_priceRegistry)
    });

    vm.expectEmit();
    emit ConfigSet(staticConfig, dynamicConfig);

    CommitStore commitStore = new CommitStore(staticConfig);
    commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );

    CommitStore.StaticConfig memory gotStaticConfig = commitStore.getStaticConfig();

    assertEq(staticConfig.chainSelector, gotStaticConfig.chainSelector);
    assertEq(staticConfig.sourceChainSelector, gotStaticConfig.sourceChainSelector);
    assertEq(staticConfig.onRamp, gotStaticConfig.onRamp);
    assertEq(staticConfig.armProxy, gotStaticConfig.armProxy);

    CommitStore.DynamicConfig memory gotDynamicConfig = commitStore.getDynamicConfig();

    assertEq(dynamicConfig.priceRegistry, gotDynamicConfig.priceRegistry);

    // CommitStore initial values
    assertEq(0, commitStore.getLatestPriceEpochAndRound());
    assertEq(1, commitStore.getExpectedNextSequenceNumber());
    assertEq(commitStore.typeAndVersion(), "CommitStore 1.2.0");
    assertEq(OWNER, commitStore.owner());
    assertTrue(commitStore.isUnpausedAndARMHealthy());
  }
}

/// @notice #setMinSeqNr
contract CommitStore_setMinSeqNr is CommitStoreSetup {
  function testFuzz_SetMinSeqNrSuccess(uint64 minSeqNr) public {
    s_commitStore.setMinSeqNr(minSeqNr);

    assertEq(s_commitStore.getExpectedNextSequenceNumber(), minSeqNr);
  }

  // Reverts
  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setMinSeqNr(6723);
  }
}

/// @notice #setDynamicConfig
contract CommitStore_setDynamicConfig is CommitStoreSetup {
  function testFuzz_SetDynamicConfigSuccess(address priceRegistry) public {
    vm.assume(priceRegistry != address(0));
    CommitStore.StaticConfig memory staticConfig = s_commitStore.getStaticConfig();
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: priceRegistry});
    bytes memory onchainConfig = abi.encode(dynamicConfig);

    vm.expectEmit();
    emit ConfigSet(staticConfig, dynamicConfig);

    uint32 configCount = 1;

    vm.expectEmit();
    emit ConfigSet(
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
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      onchainConfig,
      s_offchainConfigVersion,
      abi.encode("")
    );

    CommitStore.DynamicConfig memory gotDynamicConfig = s_commitStore.getDynamicConfig();
    assertEq(gotDynamicConfig.priceRegistry, dynamicConfig.priceRegistry);
  }

  function testPriceEpochClearedSuccess() public {
    // Set latest price epoch and round to non-zero.
    uint40 latestEpochAndRound = 1782155;
    s_commitStore.setLatestPriceEpochAndRound(latestEpochAndRound);
    assertEq(latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(1)});
    // New config should clear it.
    s_commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
    // Assert cleared.
    assertEq(0, s_commitStore.getLatestPriceEpochAndRound());
  }

  // Reverts
  function testOnlyOwnerReverts() public {
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(23784264)});

    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
  }

  function testInvalidCommitStoreConfigReverts() public {
    CommitStore.DynamicConfig memory dynamicConfig = CommitStore.DynamicConfig({priceRegistry: address(0)});

    vm.expectRevert(CommitStore.InvalidCommitStoreConfig.selector);
    s_commitStore.setOCR2Config(
      s_valid_signers,
      s_valid_transmitters,
      s_f,
      abi.encode(dynamicConfig),
      s_offchainConfigVersion,
      abi.encode("")
    );
  }
}

/// @notice #resetUnblessedRoots
contract CommitStore_resetUnblessedRoots is CommitStoreRealARMSetup {
  event RootRemoved(bytes32 root);

  function testResetUnblessedRootsSuccess() public {
    bytes32[] memory rootsToReset = new bytes32[](3);
    rootsToReset[0] = "1";
    rootsToReset[1] = "2";
    rootsToReset[2] = "3";

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 2),
      merkleRoot: rootsToReset[0]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(3, 4),
      merkleRoot: rootsToReset[1]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(5, 5),
      merkleRoot: rootsToReset[2]
    });

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    IARM.TaggedRoot[] memory blessedTaggedRoots = new IARM.TaggedRoot[](1);
    blessedTaggedRoots[0] = IARM.TaggedRoot({commitStore: address(s_commitStore), root: rootsToReset[1]});

    changePrank(BLESS_VOTE_ADDR);
    s_arm.voteToBless(blessedTaggedRoots);

    vm.expectEmit(false, false, false, true);
    emit RootRemoved(rootsToReset[0]);

    vm.expectEmit(false, false, false, true);
    emit RootRemoved(rootsToReset[2]);

    changePrank(OWNER);
    s_commitStore.resetUnblessedRoots(rootsToReset);

    assertEq(0, s_commitStore.getMerkleRoot(rootsToReset[0]));
    assertEq(BLOCK_TIME, s_commitStore.getMerkleRoot(rootsToReset[1]));
    assertEq(0, s_commitStore.getMerkleRoot(rootsToReset[2]));
  }

  // Reverts

  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    bytes32[] memory rootToReset;
    s_commitStore.resetUnblessedRoots(rootToReset);
  }
}

/// @notice #report
contract CommitStore_report is CommitStoreSetup {
  event ReportAccepted(CommitStore.CommitReport report);
  event UsdPerTokenUpdated(address indexed feeToken, uint256 value, uint256 timestamp);

  function testReportOnlyRootSuccess_gas() public {
    vm.pauseGasMetering();
    uint64 max1 = 931;
    bytes32 root = "Only a single root";
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, max1),
      merkleRoot: root
    });

    vm.expectEmit();
    emit ReportAccepted(report);

    bytes memory encodedReport = abi.encode(report);

    vm.resumeGasMetering();
    s_commitStore.report(encodedReport, ++s_latestEpochAndRound);
    vm.pauseGasMetering();

    assertEq(max1 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(block.timestamp, s_commitStore.getMerkleRoot(root));
    vm.resumeGasMetering();
  }

  function testReportAndPriceUpdateSuccess() public {
    uint64 max1 = 12;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(1, max1),
      merkleRoot: "test #2"
    });

    vm.expectEmit();
    emit ReportAccepted(report);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    assertEq(max1 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function testStaleReportWithRootSuccess() public {
    uint64 maxSeq = 12;
    uint224 tokenStartPrice = IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry)
      .getTokenPrice(s_sourceFeeToken)
      .value;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report 1"
    });

    vm.expectEmit();
    emit ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(maxSeq + 1, maxSeq * 2),
      merkleRoot: "stale report 2"
    });

    vm.expectEmit();
    emit ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
    assertEq(maxSeq * 2 + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
    assertEq(
      tokenStartPrice,
      IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
  }

  function testOnlyTokenPriceUpdatesSuccess() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function testOnlyGasPriceUpdatesSuccess() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  function testValidPriceUpdateThenStaleReportWithRootSuccess() public {
    uint64 maxSeq = 12;
    uint224 tokenPrice1 = 4e18;
    uint224 tokenPrice2 = 5e18;

    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice1),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: ""
    });

    vm.expectEmit();
    emit UsdPerTokenUpdated(s_sourceFeeToken, tokenPrice1, block.timestamp);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());

    report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, tokenPrice2),
      interval: CommitStore.Interval(1, maxSeq),
      merkleRoot: "stale report"
    });

    vm.expectEmit();
    emit ReportAccepted(report);

    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);

    assertEq(maxSeq + 1, s_commitStore.getExpectedNextSequenceNumber());
    assertEq(
      tokenPrice1,
      IPriceRegistry(s_commitStore.getDynamicConfig().priceRegistry).getTokenPrice(s_sourceFeeToken).value
    );
    assertEq(s_latestEpochAndRound, s_commitStore.getLatestPriceEpochAndRound());
  }

  // Reverts

  function testPausedReverts() public {
    s_commitStore.pause();
    bytes memory report;
    vm.expectRevert(CommitStore.PausedError.selector);
    s_commitStore.report(report, ++s_latestEpochAndRound);
  }

  function testUnhealthyReverts() public {
    s_mockARM.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    vm.expectRevert(CommitStore.BadARMSignal.selector);
    bytes memory report;
    s_commitStore.report(report, ++s_latestEpochAndRound);
  }

  function testInvalidRootRevert() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 4),
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(CommitStore.InvalidRoot.selector);
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function testInvalidIntervalReverts() public {
    CommitStore.Interval memory interval = CommitStore.Interval(2, 2);
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: interval,
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(abi.encodeWithSelector(CommitStore.InvalidInterval.selector, interval));

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function testInvalidIntervalMinLargerThanMaxReverts() public {
    CommitStore.Interval memory interval = CommitStore.Interval(1, 0);
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: interval,
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(abi.encodeWithSelector(CommitStore.InvalidInterval.selector, interval));

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }

  function testZeroEpochAndRoundReverts() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: bytes32(0)
    });

    vm.expectRevert(CommitStore.StaleReport.selector);

    s_commitStore.report(abi.encode(report), 0);
  }

  function testOnlyPriceUpdateStaleReportReverts() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getSingleTokenPriceUpdateStruct(s_sourceFeeToken, 4e18),
      interval: CommitStore.Interval(0, 0),
      merkleRoot: bytes32(0)
    });

    vm.expectEmit();
    emit UsdPerTokenUpdated(s_sourceFeeToken, 4e18, block.timestamp);
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    vm.expectRevert(CommitStore.StaleReport.selector);
    s_commitStore.report(abi.encode(report), s_latestEpochAndRound);
  }

  function testRootAlreadyCommittedReverts() public {
    CommitStore.CommitReport memory report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(1, 2),
      merkleRoot: "Only a single root"
    });
    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);

    report = CommitStore.CommitReport({
      priceUpdates: getEmptyPriceUpdates(),
      interval: CommitStore.Interval(3, 3),
      merkleRoot: "Only a single root"
    });

    vm.expectRevert(CommitStore.RootAlreadyCommitted.selector);

    s_commitStore.report(abi.encode(report), ++s_latestEpochAndRound);
  }
}

/// @notice #verify
contract CommitStore_verify is CommitStoreRealARMSetup {
  function testNotBlessedSuccess() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    s_commitStore.report(
      abi.encode(
        CommitStore.CommitReport({
          priceUpdates: getEmptyPriceUpdates(),
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

  function testBlessedSuccess() public {
    bytes32[] memory leaves = new bytes32[](1);
    leaves[0] = "root";
    s_commitStore.report(
      abi.encode(
        CommitStore.CommitReport({
          priceUpdates: getEmptyPriceUpdates(),
          interval: CommitStore.Interval(1, 2),
          merkleRoot: leaves[0]
        })
      ),
      ++s_latestEpochAndRound
    );
    // Bless that root.
    IARM.TaggedRoot[] memory taggedRoots = new IARM.TaggedRoot[](1);
    taggedRoots[0] = IARM.TaggedRoot({commitStore: address(s_commitStore), root: leaves[0]});
    changePrank(BLESS_VOTE_ADDR);
    s_arm.voteToBless(taggedRoots);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_commitStore.verify(leaves, proofs, 0);
    assertEq(BLOCK_TIME, timestamp);
  }

  // Reverts

  function testPausedReverts() public {
    s_commitStore.pause();

    bytes32[] memory hashedLeaves = new bytes32[](0);
    bytes32[] memory proofs = new bytes32[](0);
    uint256 proofFlagBits = 0;

    vm.expectRevert(CommitStore.PausedError.selector);
    s_commitStore.verify(hashedLeaves, proofs, proofFlagBits);
  }

  function testTooManyLeavesReverts() public {
    bytes32[] memory leaves = new bytes32[](258);
    bytes32[] memory proofs = new bytes32[](0);

    vm.expectRevert(MerkleMultiProof.InvalidProof.selector);

    s_commitStore.verify(leaves, proofs, 0);
  }
}

contract CommitStore_isUnpausedAndARMHealthy is CommitStoreSetup {
  function testARMSuccess() public {
    // Test pausing
    assertFalse(s_commitStore.paused());
    assertTrue(s_commitStore.isUnpausedAndARMHealthy());
    s_commitStore.pause();
    assertTrue(s_commitStore.paused());
    assertFalse(s_commitStore.isUnpausedAndARMHealthy());
    s_commitStore.unpause();
    assertFalse(s_commitStore.paused());
    assertTrue(s_commitStore.isUnpausedAndARMHealthy());

    // Test arm
    assertTrue(s_commitStore.isARMHealthy());
    s_mockARM.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    assertFalse(s_commitStore.isARMHealthy());
    assertFalse(s_commitStore.isUnpausedAndARMHealthy());
    ARM.UnvoteToCurseRecord[] memory records = new ARM.UnvoteToCurseRecord[](1);
    records[0] = ARM.UnvoteToCurseRecord({curseVoteAddr: OWNER, cursesHash: bytes32(uint256(0)), forceUnvote: true});
    s_mockARM.ownerUnvoteToCurse(records);
    assertTrue(s_commitStore.isARMHealthy());
    assertTrue(s_commitStore.isUnpausedAndARMHealthy());

    s_mockARM.voteToCurse(0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff);
    s_commitStore.pause();
    assertFalse(s_commitStore.isUnpausedAndARMHealthy());
  }
}

/// @notice #setLatestPriceEpochAndRound
contract CommitStore_setLatestPriceEpochAndRound is CommitStoreSetup {
  function testSetLatestPriceEpochAndRoundSuccess() public {
    uint40 latestRoundAndEpoch = 1782155;
    s_commitStore.setLatestPriceEpochAndRound(latestRoundAndEpoch);

    assertEq(s_commitStore.getLatestPriceEpochAndRound(), latestRoundAndEpoch);
  }

  // Reverts
  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_commitStore.setLatestPriceEpochAndRound(6723);
  }
}
