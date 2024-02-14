// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IARM} from "../../interfaces/IARM.sol";

import {Test} from "forge-std/Test.sol";
import {ARMSetup} from "./ARMSetup.t.sol";
import {ARM} from "../../ARM.sol";

contract ConfigCompare is Test {
  function assertConfigEq(ARM.Config memory actualConfig, ARM.Config memory expectedConfig) public {
    assertEq(actualConfig.voters.length, expectedConfig.voters.length);
    for (uint256 i = 0; i < expectedConfig.voters.length; ++i) {
      ARM.Voter memory expectedVoter = expectedConfig.voters[i];
      ARM.Voter memory actualVoter = actualConfig.voters[i];
      assertEq(actualVoter.blessVoteAddr, expectedVoter.blessVoteAddr);
      assertEq(actualVoter.curseVoteAddr, expectedVoter.curseVoteAddr);
      assertEq(actualVoter.blessWeight, expectedVoter.blessWeight);
      assertEq(actualVoter.curseWeight, expectedVoter.curseWeight);
    }
    assertEq(actualConfig.blessWeightThreshold, expectedConfig.blessWeightThreshold);
    assertEq(actualConfig.curseWeightThreshold, expectedConfig.curseWeightThreshold);
  }
}

contract ARM_constructor is ConfigCompare, ARMSetup {
  function testConstructorSuccess() public {
    ARM.Config memory expectedConfig = armConstructorArgs();
    (uint32 actualVersion, , ARM.Config memory actualConfig) = s_arm.getConfigDetails();
    assertEq(actualVersion, 1);
    assertConfigEq(actualConfig, expectedConfig);
  }
}

contract ARM_voteToBlessRoots is ARMSetup {
  event VotedToBless(uint32 indexed configVersion, address indexed voter, IARM.TaggedRoot taggedRoot, uint8 weight);

  // Success

  function _getFirstBlessVoterAndWeight() internal pure returns (address, uint8) {
    ARM.Config memory cfg = armConstructorArgs();
    return (cfg.voters[0].blessVoteAddr, cfg.voters[0].blessWeight);
  }

  function test1RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    vm.expectEmit();
    emit VotedToBless(1, voter, makeTaggedRoot(1), voterWeight);

    changePrank(voter);
    vm.resumeGasMetering();
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    vm.pauseGasMetering();

    assertFalse(s_arm.isBlessed(makeTaggedRoot(1)));
    assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
    assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));
    vm.resumeGasMetering();
  }

  function test3RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    for (uint256 i = 1; i <= 3; ++i) {
      vm.expectEmit();
      emit VotedToBless(1, voter, makeTaggedRoot(i), voterWeight);
    }

    changePrank(voter);
    vm.resumeGasMetering();
    s_arm.voteToBless(makeTaggedRootsInclusive(1, 3));
    vm.pauseGasMetering();

    for (uint256 i = 1; i <= 3; ++i) {
      assertFalse(s_arm.isBlessed(makeTaggedRoot(i)));
      assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(i)));
      assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(i)));
    }
    vm.resumeGasMetering();
  }

  function test5RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    for (uint256 i = 1; i <= 5; ++i) {
      vm.expectEmit();
      emit VotedToBless(1, voter, makeTaggedRoot(i), voterWeight);
    }

    changePrank(voter);
    vm.resumeGasMetering();
    s_arm.voteToBless(makeTaggedRootsInclusive(1, 5));
    vm.pauseGasMetering();

    for (uint256 i = 1; i <= 5; ++i) {
      assertFalse(s_arm.isBlessed(makeTaggedRoot(i)));
      assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(i)));
      assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));
    }
    vm.resumeGasMetering();
  }

  function testIsAlreadyBlessedIgnoredSuccess() public {
    ARM.Config memory cfg = armConstructorArgs();

    // Bless voters 2,3,4 vote to bless
    for (uint256 i = 1; i < cfg.voters.length; i++) {
      changePrank(cfg.voters[i].blessVoteAddr);
      s_arm.voteToBless(makeTaggedRootSingleton(1));
    }

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    changePrank(cfg.voters[0].blessVoteAddr);
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  function testSenderAlreadyVotedIgnoredSuccess() public {
    (address voter, ) = _getFirstBlessVoterAndWeight();

    changePrank(voter);
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  // Reverts

  function testCurseReverts() public {
    ARM.Config memory cfg = armConstructorArgs();

    for (uint256 i = 0; i < cfg.voters.length; i++) {
      changePrank(cfg.voters[i].curseVoteAddr);
      s_arm.voteToCurse(makeCurseId(i));
    }

    changePrank(cfg.voters[0].blessVoteAddr);
    vm.expectRevert(ARM.MustRecoverFromCurse.selector);
    s_arm.voteToBless(makeTaggedRootSingleton(12903));
  }

  function testInvalidVoterReverts() public {
    changePrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(ARM.InvalidVoter.selector, STRANGER));
    s_arm.voteToBless(makeTaggedRootSingleton(12321));
  }
}

contract ARM_ownerUnbless is ARMSetup {
  function testUnblessSuccess() public {
    ARM.Config memory cfg = armConstructorArgs();
    for (uint256 i = 0; i < cfg.voters.length; ++i) {
      changePrank(cfg.voters[i].blessVoteAddr);
      s_arm.voteToBless(makeTaggedRootSingleton(1));
    }
    assertTrue(s_arm.isBlessed(makeTaggedRoot(1)));

    changePrank(OWNER);
    s_arm.ownerResetBlessVotes(makeTaggedRootSingleton(1));
    assertFalse(s_arm.isBlessed(makeTaggedRoot(1)));
  }
}

contract ARM_unvoteToCurse is ARMSetup {
  uint256 s_curser;
  bytes32 s_cursesHash;

  function setUp() public override {
    ARM.Config memory cfg = armConstructorArgs();
    ARMSetup.setUp();
    cfg = armConstructorArgs();
    s_curser = 0;

    changePrank(cfg.voters[0].curseVoteAddr);
    s_arm.voteToCurse(makeCurseId(1));
    bytes32 expectedCursesHash = keccak256(abi.encode(bytes32(0), makeCurseId(1)));
    assertFalse(s_arm.isCursed());
    (
      address[] memory cursers,
      uint32[] memory voteCounts,
      bytes32[] memory cursesHashes,
      uint16 weight,
      bool cursed
    ) = s_arm.getCurseProgress();
    assertEq(1, cursers.length);
    assertEq(1, voteCounts.length);
    assertEq(cfg.voters[s_curser].curseVoteAddr, cursers[0]);
    assertEq(1, voteCounts[0]);
    assertEq(cfg.voters[s_curser].curseWeight, weight);
    assertEq(1, cursesHashes.length);
    assertEq(expectedCursesHash, cursesHashes[0]);
    assertFalse(cursed);

    s_cursesHash = expectedCursesHash;
  }

  function testInvalidVoter() public {
    ARM.Config memory cfg = armConstructorArgs();
    // Someone else cannot unvote to curse on the curser's behalf.
    address[] memory unauthorized = new address[](4);
    unauthorized[0] = cfg.voters[s_curser].blessVoteAddr;
    unauthorized[1] = cfg.voters[s_curser].curseVoteAddr;
    unauthorized[2] = OWNER;
    unauthorized[3] = cfg.voters[s_curser ^ 1].curseUnvoteAddr;

    for (uint256 i = 0; i < unauthorized.length; ++i) {
      bytes memory expectedRevert = abi.encodeWithSelector(ARM.InvalidVoter.selector, unauthorized[i]);
      changePrank(unauthorized[i]);
      // should fail when using the correct curses hash
      vm.expectRevert(expectedRevert);
      s_arm.unvoteToCurse(cfg.voters[s_curser].curseVoteAddr, s_cursesHash);
      // should fail when using garbage curses hash
      vm.expectRevert(expectedRevert);
      s_arm.unvoteToCurse(
        cfg.voters[s_curser].curseVoteAddr,
        0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      );
    }
  }

  function testInvalidCursesHash() public {
    ARM.Config memory cfg = armConstructorArgs();
    changePrank(cfg.voters[s_curser].curseUnvoteAddr);
    vm.expectRevert(
      abi.encodeWithSelector(
        ARM.InvalidCursesHash.selector,
        s_cursesHash,
        0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      )
    );
    s_arm.unvoteToCurse(
      cfg.voters[s_curser].curseVoteAddr,
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
    );
  }

  function testValidCursesHash() public {
    ARM.Config memory cfg = armConstructorArgs();
    changePrank(cfg.voters[s_curser].curseUnvoteAddr);
    s_arm.unvoteToCurse(cfg.voters[s_curser].curseVoteAddr, s_cursesHash);
  }

  function testOwnerSucceeds() public {
    ARM.Config memory cfg = armConstructorArgs();
    changePrank(OWNER);
    ARM.UnvoteToCurseRecord[] memory records = new ARM.UnvoteToCurseRecord[](1);
    records[0] = ARM.UnvoteToCurseRecord({
      curseVoteAddr: cfg.voters[s_curser].curseUnvoteAddr,
      cursesHash: s_cursesHash,
      forceUnvote: false
    });
    s_arm.ownerUnvoteToCurse(records);
  }

  event SkippedUnvoteToCurse(address indexed voter, bytes32 expectedCursesHash, bytes32 actualCursesHash);

  function testOwnerSkips() public {
    ARM.Config memory cfg = armConstructorArgs();
    changePrank(OWNER);
    ARM.UnvoteToCurseRecord[] memory records = new ARM.UnvoteToCurseRecord[](1);
    records[0] = ARM.UnvoteToCurseRecord({
      curseVoteAddr: cfg.voters[s_curser].curseVoteAddr,
      cursesHash: 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff,
      forceUnvote: false
    });
    vm.expectEmit();
    emit SkippedUnvoteToCurse(
      cfg.voters[s_curser].curseVoteAddr,
      s_cursesHash,
      0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
    );
    s_arm.ownerUnvoteToCurse(records);
  }

  function testInvalidCurseStateReverts() public {
    ARM.Config memory cfg = armConstructorArgs();
    changePrank(cfg.voters[1].curseUnvoteAddr);

    vm.expectRevert(ARM.InvalidCurseState.selector);
    s_arm.unvoteToCurse(cfg.voters[1].curseVoteAddr, "");
  }
}

contract ARM_voteToCurse is ARMSetup {
  event VotedToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    uint8 weight,
    uint32 voteCount,
    bytes32 curseId,
    bytes32 cursesHash,
    uint16 accumulatedWeight
  );
  event Cursed(uint32 indexed configVersion, uint256 timestamp);
  event OwnerCursed(uint256 timestamp);
  event RecoveredFromCurse();

  function _getFirstCurseVoterAndWeight() internal pure returns (address, uint8) {
    ARM.Config memory cfg = armConstructorArgs();
    return (cfg.voters[0].curseVoteAddr, cfg.voters[0].curseWeight);
  }

  // Success

  function testVoteToCurseSuccess_gas() public {
    vm.pauseGasMetering();

    (address voter, uint8 weight) = _getFirstCurseVoterAndWeight();
    changePrank(voter);
    vm.expectEmit();
    emit VotedToCurse(
      1,
      voter,
      weight,
      1,
      makeCurseId(123),
      keccak256(abi.encode(bytes32(0), makeCurseId(123))),
      weight
    );

    vm.resumeGasMetering();
    s_arm.voteToCurse(makeCurseId(123));
    vm.pauseGasMetering();

    (address[] memory voters, , , uint16 votes, bool cursed) = s_arm.getCurseProgress();
    assertEq(1, voters.length);
    assertEq(voter, voters[0]);
    assertEq(weight, votes);
    assertFalse(cursed);

    vm.resumeGasMetering();
  }

  function testEmitCurseSuccess() public {
    ARM.Config memory cfg = armConstructorArgs();
    for (uint256 i = 0; i < cfg.voters.length - 1; ++i) {
      changePrank(cfg.voters[i].curseVoteAddr);
      s_arm.voteToCurse(makeCurseId(1));
    }

    vm.expectEmit();
    emit Cursed(1, block.timestamp);

    changePrank(cfg.voters[cfg.voters.length - 1].curseVoteAddr);
    s_arm.voteToCurse(makeCurseId(1));
  }

  function testEvenIfAlreadyCursedSuccess() public {
    ARM.Config memory cfg = armConstructorArgs();
    uint16 weightSum = 0;
    for (uint256 i = 0; i < cfg.voters.length; ++i) {
      changePrank(cfg.voters[i].curseVoteAddr);
      s_arm.voteToCurse(makeCurseId(i));
      weightSum += cfg.voters[i].curseWeight;
    }

    // Not part of the assertion of this test but good to have as a sanity
    // check. We want a curse to be active in order for the ultimate assertion
    // to make sense.
    assert(s_arm.isCursed());

    vm.expectEmit();
    emit VotedToCurse(
      1, // configVersion
      cfg.voters[cfg.voters.length - 1].curseVoteAddr,
      cfg.voters[cfg.voters.length - 1].curseWeight,
      2, // voteCount
      makeCurseId(cfg.voters.length + 1), // this curse id
      keccak256(
        abi.encode(
          keccak256(abi.encode(bytes32(0), makeCurseId(cfg.voters.length - 1))),
          makeCurseId(cfg.voters.length + 1)
        )
      ), // cursesHash
      weightSum // accumulatedWeight
    );
    // Asserts that this call to vote with a new curse id goes through with no
    // reverts even when the ARM contract is cursed.
    s_arm.voteToCurse(makeCurseId(cfg.voters.length + 1));
  }

  function testOwnerCanCurseAndUncurse() public {
    changePrank(OWNER);
    vm.expectEmit();
    emit OwnerCursed(block.timestamp);
    vm.expectEmit();
    emit Cursed(1, block.timestamp);
    s_arm.ownerCurse();

    {
      (address[] memory voters, , , uint24 accWeight, bool cursed) = s_arm.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    // ownerCurse again, this time we only get OwnerCursed, but not Cursed
    vm.expectEmit();
    emit OwnerCursed(block.timestamp);
    s_arm.ownerCurse();

    {
      (address[] memory voters, , , uint24 accWeight, bool cursed) = s_arm.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    ARM.UnvoteToCurseRecord[] memory unvoteRecords = new ARM.UnvoteToCurseRecord[](0);
    vm.expectEmit();
    emit RecoveredFromCurse();
    s_arm.ownerUnvoteToCurse(unvoteRecords);
    {
      (address[] memory voters, , , uint24 accWeight, bool cursed) = s_arm.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertFalse(cursed);
    }
  }

  // Reverts

  function testInvalidVoterReverts() public {
    changePrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(ARM.InvalidVoter.selector, STRANGER));
    s_arm.voteToCurse(makeCurseId(12312));
  }

  function testAlreadyVotedReverts() public {
    (address voter, ) = _getFirstCurseVoterAndWeight();
    changePrank(voter);
    s_arm.voteToCurse(makeCurseId(1));

    vm.expectRevert(abi.encodeWithSelector(ARM.AlreadyVotedToCurse.selector, voter, makeCurseId(1)));
    s_arm.voteToCurse(makeCurseId(1));
  }
}

contract ARM_ownerUnvoteToCurse is ARMSetup {
  event RecoveredFromCurse();

  // These cursers are going to curse in setUp curseCount times.
  function getCursersAndCurseCounts() internal pure returns (address[] memory cursers, uint32[] memory curseCounts) {
    // NOTE: Change this when changing setUp or armConstructorArgs.
    // This is a bit ugly and error prone but if we read from storage we would
    // not get an accurate gas reading for ownerUnvoteToCurse when we need it.
    cursers = new address[](4);
    cursers[0] = CURSE_VOTER_1;
    cursers[1] = CURSE_VOTER_2;
    cursers[2] = CURSE_VOTER_3;
    cursers[3] = CURSE_VOTER_4;
    curseCounts = new uint32[](cursers.length);
    for (uint256 i = 0; i < cursers.length; ++i) {
      curseCounts[i] = 1;
    }
  }

  function setUp() public virtual override {
    ARMSetup.setUp();
    (address[] memory cursers, uint32[] memory curseCounts) = getCursersAndCurseCounts();
    for (uint256 i = 0; i < cursers.length; ++i) {
      changePrank(cursers[i]);
      for (uint256 j = 0; j < curseCounts[i]; ++j) {
        s_arm.voteToCurse(makeCurseId(j));
      }
    }
  }

  function ownerUnvoteToCurse() internal {
    s_arm.ownerUnvoteToCurse(makeUnvoteToCurseRecords());
  }

  function makeUnvoteToCurseRecords() internal pure returns (ARM.UnvoteToCurseRecord[] memory) {
    (address[] memory cursers, ) = getCursersAndCurseCounts();
    ARM.UnvoteToCurseRecord[] memory records = new ARM.UnvoteToCurseRecord[](cursers.length);
    for (uint256 i = 0; i < cursers.length; ++i) {
      records[i] = ARM.UnvoteToCurseRecord({
        curseVoteAddr: cursers[i],
        cursesHash: bytes32(uint256(0)),
        forceUnvote: true
      });
    }
    return records;
  }

  // Success

  function testOwnerUnvoteToCurseSuccess_gas() public {
    vm.pauseGasMetering();
    changePrank(OWNER);

    vm.expectEmit();
    emit RecoveredFromCurse();

    vm.resumeGasMetering();
    ownerUnvoteToCurse();
    vm.pauseGasMetering();

    assertFalse(s_arm.isCursed());
    (address[] memory voters, , bytes32[] memory cursesHashes, uint256 weight, bool cursed) = s_arm.getCurseProgress();
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
    vm.resumeGasMetering();
  }

  function testIsIdempotent() public {
    changePrank(OWNER);
    ownerUnvoteToCurse();
    ownerUnvoteToCurse();

    assertFalse(s_arm.isCursed());
    (
      address[] memory voters,
      uint32[] memory voteCounts,
      bytes32[] memory cursesHashes,
      uint256 weight,
      bool cursed
    ) = s_arm.getCurseProgress();
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(voteCounts.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
  }

  function testCanBlessAndCurseAfterRecovery() public {
    // Contract is already cursed due to setUp.

    // Owner unvotes to curse.
    changePrank(OWNER);
    vm.expectEmit();
    emit RecoveredFromCurse();
    ownerUnvoteToCurse();

    // Contract is now uncursed.
    assertFalse(s_arm.isCursed());

    // Vote to bless should go through.
    changePrank(BLESS_VOTER_1);
    s_arm.voteToBless(makeTaggedRootSingleton(2387489729));

    // Vote to curse should go through.
    changePrank(CURSE_VOTER_1);
    s_arm.voteToCurse(makeCurseId(73894728973));
  }

  // Reverts

  function testNonOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    ownerUnvoteToCurse();
  }
}

contract ARM_setConfig is ConfigCompare, ARMSetup {
  /// @notice Test-specific function to use only in setConfig tests
  function getDifferentConfigArgs() private pure returns (ARM.Config memory) {
    ARM.Voter[] memory voters = new ARM.Voter[](2);
    voters[0] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_1,
      curseVoteAddr: CURSE_VOTER_1,
      curseUnvoteAddr: CURSE_UNVOTER_1,
      blessWeight: WEIGHT_1,
      curseWeight: WEIGHT_1
    });
    voters[1] = ARM.Voter({
      blessVoteAddr: BLESS_VOTER_2,
      curseVoteAddr: CURSE_VOTER_2,
      curseUnvoteAddr: CURSE_UNVOTER_2,
      blessWeight: WEIGHT_10,
      curseWeight: WEIGHT_10
    });
    return
      ARM.Config({
        voters: voters,
        blessWeightThreshold: WEIGHT_1 + WEIGHT_10,
        curseWeightThreshold: WEIGHT_1 + WEIGHT_10
      });
  }

  function setUp() public virtual override {
    ARMSetup.setUp();
    ARM.Config memory cfg = armConstructorArgs();

    // Setup some partial state
    changePrank(cfg.voters[0].blessVoteAddr);
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    changePrank(cfg.voters[1].blessVoteAddr);
    s_arm.voteToBless(makeTaggedRootSingleton(1));
    changePrank(cfg.voters[1].curseVoteAddr);
    s_arm.voteToCurse(makeCurseId(1));
  }

  // Success

  event ConfigSet(uint32 indexed configVersion, ARM.Config config);

  function testVoteToBlessByEjectedVoterReverts() public {
    // Previous config included BLESS_VOTER_4. Change to new config that doesn't.
    ARM.Config memory cfg = getDifferentConfigArgs();
    changePrank(OWNER);
    s_arm.setConfig(cfg);

    // BLESS_VOTER_4 is not part of cfg anymore, vote to bless should revert.
    changePrank(BLESS_VOTER_4);
    vm.expectRevert(abi.encodeWithSelector(ARM.InvalidVoter.selector, BLESS_VOTER_4));
    s_arm.voteToBless(makeTaggedRootSingleton(2));
  }

  function testSetConfigSuccess_gas() public {
    vm.pauseGasMetering();
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    vm.expectEmit();
    emit ConfigSet(2, cfg);

    (uint32 configVersionBefore, , ) = s_arm.getConfigDetails();
    vm.resumeGasMetering();
    s_arm.setConfig(cfg);
    vm.pauseGasMetering();
    // Assert VersionedConfig has changed correctly
    (uint32 configVersionAfter, , ARM.Config memory configAfter) = s_arm.getConfigDetails();
    assertEq(configVersionBefore + 1, configVersionAfter);
    assertConfigEq(configAfter, cfg);

    // Assert that curse votes have been cleared, except for CURSE_VOTER_2 who
    // has already voted and is also part of the new config
    (address[] memory curseVoters, , bytes32[] memory cursesHashes, uint256 curseWeight, bool cursed) = s_arm
      .getCurseProgress();
    assertEq(1, curseVoters.length);
    assertEq(WEIGHT_10, curseWeight);
    assertEq(1, cursesHashes.length);
    assertEq(keccak256(abi.encode(bytes32(0), makeCurseId(1))), cursesHashes[0]);
    assertFalse(cursed);

    // Assert that good votes have been cleared
    uint256 votesToBlessRoot = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    assertEq(ZERO, votesToBlessRoot);
    assertFalse(hasVotedToBlessRoot(cfg.voters[0].blessVoteAddr, makeTaggedRoot(1)));
    assertFalse(hasVotedToBlessRoot(cfg.voters[1].blessVoteAddr, makeTaggedRoot(1)));
    vm.resumeGasMetering();
  }

  // Reverts

  function testNonOwnerReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_arm.setConfig(cfg);
  }

  function testVotersLengthIsZeroReverts() public {
    changePrank(OWNER);
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(ARM.Config({voters: new ARM.Voter[](0), blessWeightThreshold: 1, curseWeightThreshold: 1}));
  }

  function testEitherThresholdIsZeroReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(
      ARM.Config({voters: cfg.voters, blessWeightThreshold: ZERO, curseWeightThreshold: cfg.curseWeightThreshold})
    );
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(
      ARM.Config({voters: cfg.voters, blessWeightThreshold: cfg.blessWeightThreshold, curseWeightThreshold: ZERO})
    );
  }

  function testBlessVoterIsZeroAddressReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    cfg.voters[0].blessVoteAddr = ZERO_ADDRESS;
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(cfg);
  }

  function testWeightIsZeroAddressReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    cfg.voters[0].blessWeight = ZERO;
    cfg.voters[0].curseWeight = ZERO;
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(cfg);
  }

  function testTotalWeightsSmallerThanEachThresholdReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(
      ARM.Config({voters: cfg.voters, blessWeightThreshold: WEIGHT_40, curseWeightThreshold: cfg.curseWeightThreshold})
    );
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(
      ARM.Config({voters: cfg.voters, blessWeightThreshold: cfg.blessWeightThreshold, curseWeightThreshold: WEIGHT_40})
    );
  }

  function testRepeatedAddressReverts() public {
    ARM.Config memory cfg = getDifferentConfigArgs();

    changePrank(OWNER);
    cfg.voters[0].blessVoteAddr = cfg.voters[1].curseVoteAddr;
    vm.expectRevert(ARM.InvalidConfig.selector);
    s_arm.setConfig(cfg);
  }
}
