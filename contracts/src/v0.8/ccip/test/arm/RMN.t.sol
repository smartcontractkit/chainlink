// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

import {RMN} from "../../RMN.sol";
import {RMNSetup} from "./RMNSetup.t.sol";
import {Test} from "forge-std/Test.sol";

contract ConfigCompare is Test {
  function assertConfigEq(RMN.Config memory actualConfig, RMN.Config memory expectedConfig) public pure {
    assertEq(actualConfig.voters.length, expectedConfig.voters.length);
    for (uint256 i = 0; i < expectedConfig.voters.length; ++i) {
      RMN.Voter memory expectedVoter = expectedConfig.voters[i];
      RMN.Voter memory actualVoter = actualConfig.voters[i];
      assertEq(actualVoter.blessVoteAddr, expectedVoter.blessVoteAddr);
      assertEq(actualVoter.curseVoteAddr, expectedVoter.curseVoteAddr);
      assertEq(actualVoter.blessWeight, expectedVoter.blessWeight);
      assertEq(actualVoter.curseWeight, expectedVoter.curseWeight);
    }
    assertEq(actualConfig.blessWeightThreshold, expectedConfig.blessWeightThreshold);
    assertEq(actualConfig.curseWeightThreshold, expectedConfig.curseWeightThreshold);
  }
}

contract RMN_constructor is ConfigCompare, RMNSetup {
  function test_Constructor_Success() public view {
    RMN.Config memory expectedConfig = rmnConstructorArgs();
    (uint32 actualVersion,, RMN.Config memory actualConfig) = s_rmn.getConfigDetails();
    assertEq(actualVersion, 1);
    assertConfigEq(actualConfig, expectedConfig);
  }
}

contract RMN_voteToBlessRoots is RMNSetup {
  event VotedToBless(uint32 indexed configVersion, address indexed voter, IRMN.TaggedRoot taggedRoot, uint8 weight);

  // Success

  function _getFirstBlessVoterAndWeight() internal pure returns (address, uint8) {
    RMN.Config memory cfg = rmnConstructorArgs();
    return (cfg.voters[0].blessVoteAddr, cfg.voters[0].blessWeight);
  }

  function test_1RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    vm.expectEmit();
    emit VotedToBless(1, voter, makeTaggedRoot(1), voterWeight);

    vm.startPrank(voter);
    vm.resumeGasMetering();
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    vm.pauseGasMetering();

    assertFalse(s_rmn.isBlessed(makeTaggedRoot(1)));
    assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
    assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));
    vm.resumeGasMetering();
  }

  function test_3RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    for (uint256 i = 1; i <= 3; ++i) {
      vm.expectEmit();
      emit VotedToBless(1, voter, makeTaggedRoot(i), voterWeight);
    }

    vm.startPrank(voter);
    vm.resumeGasMetering();
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, 3));
    vm.pauseGasMetering();

    for (uint256 i = 1; i <= 3; ++i) {
      assertFalse(s_rmn.isBlessed(makeTaggedRoot(i)));
      assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(i)));
      assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(i)));
    }
    vm.resumeGasMetering();
  }

  function test_5RootSuccess_gas() public {
    vm.pauseGasMetering();
    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    for (uint256 i = 1; i <= 5; ++i) {
      vm.expectEmit();
      emit VotedToBless(1, voter, makeTaggedRoot(i), voterWeight);
    }

    vm.startPrank(voter);
    vm.resumeGasMetering();
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, 5));
    vm.pauseGasMetering();

    for (uint256 i = 1; i <= 5; ++i) {
      assertFalse(s_rmn.isBlessed(makeTaggedRoot(i)));
      assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(i)));
      assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));
    }
    vm.resumeGasMetering();
  }

  function test_IsAlreadyBlessedIgnored_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();

    // Bless voters 2,3,4 vote to bless
    for (uint256 i = 1; i < cfg.voters.length; i++) {
      vm.startPrank(cfg.voters[i].blessVoteAddr);
      s_rmn.voteToBless(makeTaggedRootSingleton(1));
    }

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    vm.startPrank(cfg.voters[0].blessVoteAddr);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  function test_SenderAlreadyVotedIgnored_Success() public {
    (address voter,) = _getFirstBlessVoterAndWeight();

    vm.startPrank(voter);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  // Reverts

  function test_Curse_Revert() public {
    RMN.Config memory cfg = rmnConstructorArgs();

    for (uint256 i = 0; i < cfg.voters.length; i++) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(i));
    }

    vm.startPrank(cfg.voters[0].blessVoteAddr);
    vm.expectRevert(RMN.MustRecoverFromCurse.selector);
    s_rmn.voteToBless(makeTaggedRootSingleton(12903));
  }

  function test_InvalidVoter_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(RMN.InvalidVoter.selector, STRANGER));
    s_rmn.voteToBless(makeTaggedRootSingleton(12321));
  }
}

contract RMN_ownerUnbless is RMNSetup {
  function test_Unbless_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    for (uint256 i = 0; i < cfg.voters.length; ++i) {
      vm.startPrank(cfg.voters[i].blessVoteAddr);
      s_rmn.voteToBless(makeTaggedRootSingleton(1));
    }
    assertTrue(s_rmn.isBlessed(makeTaggedRoot(1)));

    vm.startPrank(OWNER);
    s_rmn.ownerResetBlessVotes(makeTaggedRootSingleton(1));
    assertFalse(s_rmn.isBlessed(makeTaggedRoot(1)));
  }
}

contract RMN_unvoteToCurse is RMNSetup {
  uint256 internal s_curser;
  bytes32 internal s_cursesHash;

  function setUp() public override {
    RMN.Config memory cfg = rmnConstructorArgs();
    RMNSetup.setUp();
    cfg = rmnConstructorArgs();
    s_curser = 0;

    vm.startPrank(cfg.voters[0].curseVoteAddr);
    s_rmn.voteToCurse(makeCurseId(1));
    bytes32 expectedCursesHash = keccak256(abi.encode(bytes32(0), makeCurseId(1)));
    assertFalse(s_rmn.isCursed());
    (address[] memory cursers, uint32[] memory voteCounts, bytes32[] memory cursesHashes, uint16 weight, bool cursed) =
      s_rmn.getCurseProgress();
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

  function test_InvalidVoter() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    // Someone else cannot unvote to curse on the curser's behalf.
    address[] memory unauthorized = new address[](4);
    unauthorized[0] = cfg.voters[s_curser].blessVoteAddr;
    unauthorized[1] = cfg.voters[s_curser].curseVoteAddr;
    unauthorized[2] = OWNER;
    unauthorized[3] = cfg.voters[s_curser ^ 1].curseUnvoteAddr;

    for (uint256 i = 0; i < unauthorized.length; ++i) {
      bytes memory expectedRevert = abi.encodeWithSelector(RMN.InvalidVoter.selector, unauthorized[i]);
      vm.startPrank(unauthorized[i]);
      // should fail when using the correct curses hash
      vm.expectRevert(expectedRevert);
      s_rmn.unvoteToCurse(cfg.voters[s_curser].curseVoteAddr, s_cursesHash);
      // should fail when using garbage curses hash
      vm.expectRevert(expectedRevert);
      s_rmn.unvoteToCurse(
        cfg.voters[s_curser].curseVoteAddr, 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      );
    }
  }

  function test_InvalidCursesHash() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(cfg.voters[s_curser].curseUnvoteAddr);
    vm.expectRevert(
      abi.encodeWithSelector(
        RMN.InvalidCursesHash.selector, s_cursesHash, 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
      )
    );
    s_rmn.unvoteToCurse(
      cfg.voters[s_curser].curseVoteAddr, 0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
    );
  }

  function test_ValidCursesHash() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(cfg.voters[s_curser].curseUnvoteAddr);
    s_rmn.unvoteToCurse(cfg.voters[s_curser].curseVoteAddr, s_cursesHash);
  }

  function test_OwnerSucceeds() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(OWNER);
    RMN.UnvoteToCurseRecord[] memory records = new RMN.UnvoteToCurseRecord[](1);
    records[0] = RMN.UnvoteToCurseRecord({
      curseVoteAddr: cfg.voters[s_curser].curseUnvoteAddr,
      cursesHash: s_cursesHash,
      forceUnvote: false
    });
    s_rmn.ownerUnvoteToCurse(records);
  }

  event SkippedUnvoteToCurse(address indexed voter, bytes32 expectedCursesHash, bytes32 actualCursesHash);

  function test_OwnerSkips() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(OWNER);
    RMN.UnvoteToCurseRecord[] memory records = new RMN.UnvoteToCurseRecord[](1);
    records[0] = RMN.UnvoteToCurseRecord({
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
    s_rmn.ownerUnvoteToCurse(records);
  }

  function test_InvalidCurseState_Revert() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(cfg.voters[1].curseUnvoteAddr);

    vm.expectRevert(RMN.InvalidCurseState.selector);
    s_rmn.unvoteToCurse(cfg.voters[1].curseVoteAddr, "");
  }
}

contract RMN_voteToCurse is RMNSetup {
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
    RMN.Config memory cfg = rmnConstructorArgs();
    return (cfg.voters[0].curseVoteAddr, cfg.voters[0].curseWeight);
  }

  // Success

  function test_VoteToCurseSuccess_gas() public {
    vm.pauseGasMetering();

    (address voter, uint8 weight) = _getFirstCurseVoterAndWeight();
    vm.startPrank(voter);
    vm.expectEmit();
    emit VotedToCurse(
      1, voter, weight, 1, makeCurseId(123), keccak256(abi.encode(bytes32(0), makeCurseId(123))), weight
    );

    vm.resumeGasMetering();
    s_rmn.voteToCurse(makeCurseId(123));
    vm.pauseGasMetering();

    (address[] memory voters,,, uint16 votes, bool cursed) = s_rmn.getCurseProgress();
    assertEq(1, voters.length);
    assertEq(voter, voters[0]);
    assertEq(weight, votes);
    assertFalse(cursed);

    vm.resumeGasMetering();
  }

  function test_EmitCurse_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    for (uint256 i = 0; i < cfg.voters.length - 1; ++i) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(1));
    }

    vm.expectEmit();
    emit Cursed(1, block.timestamp);

    vm.startPrank(cfg.voters[cfg.voters.length - 1].curseVoteAddr);
    s_rmn.voteToCurse(makeCurseId(1));
  }

  function test_EvenIfAlreadyCursed_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    uint16 weightSum = 0;
    for (uint256 i = 0; i < cfg.voters.length; ++i) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(i));
      weightSum += cfg.voters[i].curseWeight;
    }

    // Not part of the assertion of this test but good to have as a sanity
    // check. We want a curse to be active in order for the ultimate assertion
    // to make sense.
    assert(s_rmn.isCursed());

    vm.expectEmit();
    emit VotedToCurse(
      1, // configVersion
      cfg.voters[cfg.voters.length - 1].curseVoteAddr,
      cfg.voters[cfg.voters.length - 1].curseWeight,
      2, // voteCount
      makeCurseId(cfg.voters.length + 1), // this curse id
      keccak256(
        abi.encode(
          keccak256(abi.encode(bytes32(0), makeCurseId(cfg.voters.length - 1))), makeCurseId(cfg.voters.length + 1)
        )
      ), // cursesHash
      weightSum // accumulatedWeight
    );
    // Asserts that this call to vote with a new curse id goes through with no
    // reverts even when the RMN contract is cursed.
    s_rmn.voteToCurse(makeCurseId(cfg.voters.length + 1));
  }

  function test_OwnerCanCurseAndUncurse() public {
    vm.startPrank(OWNER);
    vm.expectEmit();
    emit OwnerCursed(block.timestamp);
    vm.expectEmit();
    emit Cursed(1, block.timestamp);
    s_rmn.ownerCurse();

    {
      (address[] memory voters,,, uint24 accWeight, bool cursed) = s_rmn.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    // ownerCurse again, this time we only get OwnerCursed, but not Cursed
    vm.expectEmit();
    emit OwnerCursed(block.timestamp);
    s_rmn.ownerCurse();

    {
      (address[] memory voters,,, uint24 accWeight, bool cursed) = s_rmn.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    RMN.UnvoteToCurseRecord[] memory unvoteRecords = new RMN.UnvoteToCurseRecord[](0);
    vm.expectEmit();
    emit RecoveredFromCurse();
    s_rmn.ownerUnvoteToCurse(unvoteRecords);
    {
      (address[] memory voters,,, uint24 accWeight, bool cursed) = s_rmn.getCurseProgress();
      assertEq(voters.length, 0);
      assertEq(accWeight, 0);
      assertFalse(cursed);
    }
  }

  // Reverts

  function test_InvalidVoter_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(RMN.InvalidVoter.selector, STRANGER));
    s_rmn.voteToCurse(makeCurseId(12312));
  }

  function test_AlreadyVoted_Revert() public {
    (address voter,) = _getFirstCurseVoterAndWeight();
    vm.startPrank(voter);
    s_rmn.voteToCurse(makeCurseId(1));

    vm.expectRevert(abi.encodeWithSelector(RMN.AlreadyVotedToCurse.selector, voter, makeCurseId(1)));
    s_rmn.voteToCurse(makeCurseId(1));
  }
}

contract RMN_ownerUnvoteToCurse is RMNSetup {
  event RecoveredFromCurse();

  // These cursers are going to curse in setUp curseCount times.
  function getCursersAndCurseCounts() internal pure returns (address[] memory cursers, uint32[] memory curseCounts) {
    // NOTE: Change this when changing setUp or rmnConstructorArgs.
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
    RMNSetup.setUp();
    (address[] memory cursers, uint32[] memory curseCounts) = getCursersAndCurseCounts();
    for (uint256 i = 0; i < cursers.length; ++i) {
      vm.startPrank(cursers[i]);
      for (uint256 j = 0; j < curseCounts[i]; ++j) {
        s_rmn.voteToCurse(makeCurseId(j));
      }
    }
  }

  function ownerUnvoteToCurse() internal {
    s_rmn.ownerUnvoteToCurse(makeUnvoteToCurseRecords());
  }

  function makeUnvoteToCurseRecords() internal pure returns (RMN.UnvoteToCurseRecord[] memory) {
    (address[] memory cursers,) = getCursersAndCurseCounts();
    RMN.UnvoteToCurseRecord[] memory records = new RMN.UnvoteToCurseRecord[](cursers.length);
    for (uint256 i = 0; i < cursers.length; ++i) {
      records[i] =
        RMN.UnvoteToCurseRecord({curseVoteAddr: cursers[i], cursesHash: bytes32(uint256(0)), forceUnvote: true});
    }
    return records;
  }

  // Success

  function test_OwnerUnvoteToCurseSuccess_gas() public {
    vm.pauseGasMetering();
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit RecoveredFromCurse();

    vm.resumeGasMetering();
    ownerUnvoteToCurse();
    vm.pauseGasMetering();

    assertFalse(s_rmn.isCursed());
    (address[] memory voters,, bytes32[] memory cursesHashes, uint256 weight, bool cursed) = s_rmn.getCurseProgress();
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
    vm.resumeGasMetering();
  }

  function test_IsIdempotent() public {
    vm.startPrank(OWNER);
    ownerUnvoteToCurse();
    ownerUnvoteToCurse();

    assertFalse(s_rmn.isCursed());
    (address[] memory voters, uint32[] memory voteCounts, bytes32[] memory cursesHashes, uint256 weight, bool cursed) =
      s_rmn.getCurseProgress();
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(voteCounts.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
  }

  function test_CanBlessAndCurseAfterRecovery() public {
    // Contract is already cursed due to setUp.

    // Owner unvotes to curse.
    vm.startPrank(OWNER);
    vm.expectEmit();
    emit RecoveredFromCurse();
    ownerUnvoteToCurse();

    // Contract is now uncursed.
    assertFalse(s_rmn.isCursed());

    // Vote to bless should go through.
    vm.startPrank(BLESS_VOTER_1);
    s_rmn.voteToBless(makeTaggedRootSingleton(2387489729));

    // Vote to curse should go through.
    vm.startPrank(CURSE_VOTER_1);
    s_rmn.voteToCurse(makeCurseId(73894728973));
  }

  // Reverts

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    ownerUnvoteToCurse();
  }
}

contract RMN_setConfig is ConfigCompare, RMNSetup {
  /// @notice Test-specific function to use only in setConfig tests
  function getDifferentConfigArgs() private pure returns (RMN.Config memory) {
    RMN.Voter[] memory voters = new RMN.Voter[](2);
    voters[0] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_1,
      curseVoteAddr: CURSE_VOTER_1,
      curseUnvoteAddr: CURSE_UNVOTER_1,
      blessWeight: WEIGHT_1,
      curseWeight: WEIGHT_1
    });
    voters[1] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_2,
      curseVoteAddr: CURSE_VOTER_2,
      curseUnvoteAddr: CURSE_UNVOTER_2,
      blessWeight: WEIGHT_10,
      curseWeight: WEIGHT_10
    });
    return RMN.Config({
      voters: voters,
      blessWeightThreshold: WEIGHT_1 + WEIGHT_10,
      curseWeightThreshold: WEIGHT_1 + WEIGHT_10
    });
  }

  function setUp() public virtual override {
    RMNSetup.setUp();
    RMN.Config memory cfg = rmnConstructorArgs();

    // Setup some partial state
    vm.startPrank(cfg.voters[0].blessVoteAddr);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    vm.startPrank(cfg.voters[1].blessVoteAddr);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    vm.startPrank(cfg.voters[1].curseVoteAddr);
    s_rmn.voteToCurse(makeCurseId(1));
  }

  // Success

  event ConfigSet(uint32 indexed configVersion, RMN.Config config);

  function test_VoteToBlessByEjectedVoter_Revert() public {
    // Previous config included BLESS_VOTER_4. Change to new config that doesn't.
    RMN.Config memory cfg = getDifferentConfigArgs();
    vm.startPrank(OWNER);
    s_rmn.setConfig(cfg);

    // BLESS_VOTER_4 is not part of cfg anymore, vote to bless should revert.
    vm.startPrank(BLESS_VOTER_4);
    vm.expectRevert(abi.encodeWithSelector(RMN.InvalidVoter.selector, BLESS_VOTER_4));
    s_rmn.voteToBless(makeTaggedRootSingleton(2));
  }

  function test_SetConfigSuccess_gas() public {
    vm.pauseGasMetering();
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    vm.expectEmit();
    emit ConfigSet(2, cfg);

    (uint32 configVersionBefore,,) = s_rmn.getConfigDetails();
    vm.resumeGasMetering();
    s_rmn.setConfig(cfg);
    vm.pauseGasMetering();
    // Assert VersionedConfig has changed correctly
    (uint32 configVersionAfter,, RMN.Config memory configAfter) = s_rmn.getConfigDetails();
    assertEq(configVersionBefore + 1, configVersionAfter);
    assertConfigEq(configAfter, cfg);

    // Assert that curse votes have been cleared, except for CURSE_VOTER_2 who
    // has already voted and is also part of the new config
    (address[] memory curseVoters,, bytes32[] memory cursesHashes, uint256 curseWeight, bool cursed) =
      s_rmn.getCurseProgress();
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

  function test_NonOwner_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_rmn.setConfig(cfg);
  }

  function test_VotersLengthIsZero_Revert() public {
    vm.startPrank(OWNER);
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(RMN.Config({voters: new RMN.Voter[](0), blessWeightThreshold: 1, curseWeightThreshold: 1}));
  }

  function test_EitherThresholdIsZero_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(
      RMN.Config({voters: cfg.voters, blessWeightThreshold: ZERO, curseWeightThreshold: cfg.curseWeightThreshold})
    );
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(
      RMN.Config({voters: cfg.voters, blessWeightThreshold: cfg.blessWeightThreshold, curseWeightThreshold: ZERO})
    );
  }

  function test_BlessVoterIsZeroAddress_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    cfg.voters[0].blessVoteAddr = ZERO_ADDRESS;
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(cfg);
  }

  function test_WeightIsZeroAddress_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    cfg.voters[0].blessWeight = ZERO;
    cfg.voters[0].curseWeight = ZERO;
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(cfg);
  }

  function test_TotalWeightsSmallerThanEachThreshold_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(
      RMN.Config({voters: cfg.voters, blessWeightThreshold: WEIGHT_40, curseWeightThreshold: cfg.curseWeightThreshold})
    );
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(
      RMN.Config({voters: cfg.voters, blessWeightThreshold: cfg.blessWeightThreshold, curseWeightThreshold: WEIGHT_40})
    );
  }

  function test_RepeatedAddress_Revert() public {
    RMN.Config memory cfg = getDifferentConfigArgs();

    vm.startPrank(OWNER);
    cfg.voters[0].blessVoteAddr = cfg.voters[1].curseVoteAddr;
    vm.expectRevert(RMN.InvalidConfig.selector);
    s_rmn.setConfig(cfg);
  }
}
