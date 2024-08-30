// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

import {GLOBAL_CURSE_SUBJECT, LIFT_CURSE_VOTE_ADDR, OWNER_CURSE_VOTE_ADDR, RMN} from "../../RMN.sol";
import {RMNSetup, makeCursesHash, makeSubjects} from "./RMNSetup.t.sol";

import {Test} from "forge-std/Test.sol";

bytes28 constant GARBAGE_CURSES_HASH = bytes28(keccak256("GARBAGE_CURSES_HASH"));

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

contract RMN_voteToBless is RMNSetup {
  function _getFirstBlessVoterAndWeight() internal pure returns (address, uint8) {
    RMN.Config memory cfg = rmnConstructorArgs();
    return (cfg.voters[0].blessVoteAddr, cfg.voters[0].blessWeight);
  }

  // Success

  function test_RootSuccess() public {
    uint256 numRoots = 10;

    (address voter, uint8 voterWeight) = _getFirstBlessVoterAndWeight();

    for (uint256 i = 1; i <= numRoots; ++i) {
      vm.expectEmit();
      emit RMN.VotedToBless(1, voter, makeTaggedRoot(i), voterWeight);
    }

    vm.prank(voter);
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, numRoots));

    for (uint256 i = 1; i <= numRoots; ++i) {
      assertFalse(s_rmn.isBlessed(makeTaggedRoot(i)));
      assertEq(voterWeight, getWeightOfVotesToBlessRoot(makeTaggedRoot(i)));
      assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));
    }
  }

  // Reverts

  function test_SenderAlreadyVoted_Revert() public {
    (address voter,) = _getFirstBlessVoterAndWeight();

    vm.startPrank(voter);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertTrue(hasVotedToBlessRoot(voter, makeTaggedRoot(1)));

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    vm.expectRevert(RMN.VoteToBlessNoop.selector);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  function test_IsAlreadyBlessed_Revert() public {
    RMN.Config memory cfg = rmnConstructorArgs();

    // Bless voters 2,3,4 vote to bless
    for (uint256 i = 1; i < cfg.voters.length; i++) {
      vm.startPrank(cfg.voters[i].blessVoteAddr);
      s_rmn.voteToBless(makeTaggedRootSingleton(1));
    }

    uint256 votesToBlessBefore = getWeightOfVotesToBlessRoot(makeTaggedRoot(1));
    vm.startPrank(cfg.voters[0].blessVoteAddr);
    vm.expectRevert(RMN.VoteToBlessNoop.selector);
    s_rmn.voteToBless(makeTaggedRootSingleton(1));
    assertEq(votesToBlessBefore, getWeightOfVotesToBlessRoot(makeTaggedRoot(1)));
  }

  function test_Curse_Revert() public {
    RMN.Config memory cfg = rmnConstructorArgs();

    for (uint256 i = 0; i < cfg.voters.length; i++) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(i), makeSubjects(GLOBAL_CURSE_SUBJECT));
    }

    vm.startPrank(cfg.voters[0].blessVoteAddr);
    vm.expectRevert(RMN.VoteToBlessForbiddenDuringActiveGlobalCurse.selector);
    s_rmn.voteToBless(makeTaggedRootSingleton(12903));
  }

  function test_UnauthorizedVoter_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(RMN.UnauthorizedVoter.selector, STRANGER));
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
  bytes28 internal s_cursesHash;

  function setUp() public override {
    RMNSetup.setUp();
    RMN.Config memory cfg = rmnConstructorArgs();

    s_curser = 0;
    vm.startPrank(cfg.voters[s_curser].curseVoteAddr);
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
    bytes28 expectedCursesHash = makeCursesHash(makeCurseId(1));
    assertFalse(s_rmn.isCursed());
    (address[] memory cursers, bytes28[] memory cursesHashes, uint16 weight, bool cursed) = s_rmn.getCurseProgress(0);
    assertEq(1, cursers.length);
    assertEq(cfg.voters[s_curser].curseVoteAddr, cursers[0]);
    assertEq(cfg.voters[s_curser].curseWeight, weight);
    assertEq(1, cursesHashes.length);
    assertEq(expectedCursesHash, cursesHashes[0]);
    assertFalse(cursed);

    s_cursesHash = expectedCursesHash;
  }

  function test_UnauthorizedVoter() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    // Someone else cannot unvote to curse on the curser's behalf.
    address[] memory unauthorized = new address[](3);
    unauthorized[0] = cfg.voters[s_curser].blessVoteAddr;
    unauthorized[1] = cfg.voters[s_curser ^ 1].blessVoteAddr;
    unauthorized[2] = OWNER;

    for (uint256 i = 0; i < unauthorized.length; ++i) {
      bytes memory expectedRevert = abi.encodeWithSelector(RMN.UnauthorizedVoter.selector, unauthorized[i]);
      vm.startPrank(unauthorized[i]);
      {
        // should fail when using the correct curses hash
        RMN.UnvoteToCurseRequest[] memory reqs = new RMN.UnvoteToCurseRequest[](1);
        reqs[0] = RMN.UnvoteToCurseRequest({subject: 0, cursesHash: s_cursesHash});
        vm.expectRevert(expectedRevert);
        s_rmn.unvoteToCurse(reqs);
      }
      {
        // should fail when using garbage curses hash
        RMN.UnvoteToCurseRequest[] memory reqs = new RMN.UnvoteToCurseRequest[](1);
        reqs[0] = RMN.UnvoteToCurseRequest({subject: 0, cursesHash: GARBAGE_CURSES_HASH});
        vm.expectRevert(expectedRevert);
        s_rmn.unvoteToCurse(reqs);
      }
    }
  }

  function test_InvalidCursesHash() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(cfg.voters[s_curser].curseVoteAddr);
    RMN.UnvoteToCurseRequest[] memory reqs = new RMN.UnvoteToCurseRequest[](1);
    reqs[0] = RMN.UnvoteToCurseRequest({subject: 0, cursesHash: GARBAGE_CURSES_HASH});
    vm.expectRevert(RMN.UnvoteToCurseNoop.selector);
    s_rmn.unvoteToCurse(reqs);
  }

  function test_ValidCursesHash() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(cfg.voters[s_curser].curseVoteAddr);
    RMN.UnvoteToCurseRequest[] memory reqs = new RMN.UnvoteToCurseRequest[](1);
    reqs[0] = RMN.UnvoteToCurseRequest({subject: 0, cursesHash: s_cursesHash});
    s_rmn.unvoteToCurse(reqs); // succeeds
  }

  function test_OwnerSucceeds() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(OWNER);
    RMN.OwnerUnvoteToCurseRequest[] memory reqs = new RMN.OwnerUnvoteToCurseRequest[](1);
    reqs[0] = RMN.OwnerUnvoteToCurseRequest({
      curseVoteAddr: cfg.voters[s_curser].curseVoteAddr,
      unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: s_cursesHash}),
      forceUnvote: false
    });
    s_rmn.ownerUnvoteToCurse(reqs);
  }

  function test_OwnerSkips() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    vm.startPrank(OWNER);
    RMN.OwnerUnvoteToCurseRequest[] memory reqs = new RMN.OwnerUnvoteToCurseRequest[](1);
    reqs[0] = RMN.OwnerUnvoteToCurseRequest({
      curseVoteAddr: cfg.voters[s_curser].curseVoteAddr,
      unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: GARBAGE_CURSES_HASH}),
      forceUnvote: false
    });

    vm.expectEmit();
    emit RMN.SkippedUnvoteToCurse(cfg.voters[s_curser].curseVoteAddr, 0, s_cursesHash, GARBAGE_CURSES_HASH);
    vm.expectRevert(RMN.UnvoteToCurseNoop.selector);
    s_rmn.ownerUnvoteToCurse(reqs);
  }

  function test_VotersCantLiftCurseButOwnerCan() public {
    vm.stopPrank();
    RMN.Config memory cfg = rmnConstructorArgs();
    // s_curser has voted to curse during setUp
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(accWeight, cfg.voters[s_curser].curseWeight);
      assertFalse(cursed);
      assertEq(voters.length, 1);
      assertEq(cursesHashes.length, 1);
      assertEq(voters[0], cfg.voters[s_curser].curseVoteAddr);
      assertEq(cursesHashes[0], makeCursesHash(makeCurseId(1)));
    }
    // everyone else votes now, same curse id, same subject
    {
      for (uint256 i = 0; i < cfg.voters.length; ++i) {
        if (i == s_curser) continue; // already voted to curse
        vm.prank(cfg.voters[i].curseVoteAddr);
        s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
      }
    }
    // subject must be cursed now
    {
      assertTrue(s_rmn.isCursed(0));
    }
    // curse progress should be as full as it can get
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      uint256 allWeights;
      for (uint256 i = 0; i < cfg.voters.length; i++) {
        allWeights += cfg.voters[i].curseWeight;
      }
      assertEq(accWeight, allWeights);
      assertTrue(cursed);
      assertEq(voters.length, cfg.voters.length);
      assertEq(cursesHashes.length, cfg.voters.length);
      for (uint256 i = 0; i < cfg.voters.length; ++i) {
        assertEq(voters[i], cfg.voters[i].curseVoteAddr);
        assertEq(cursesHashes[i], makeCursesHash(makeCurseId(1)));
      }
    }
    // everyone unvotes to curse, successfully
    {
      for (uint256 i = 0; i < cfg.voters.length; ++i) {
        vm.prank(cfg.voters[i].curseVoteAddr);
        RMN.UnvoteToCurseRequest[] memory reqs = new RMN.UnvoteToCurseRequest[](1);
        reqs[0] = RMN.UnvoteToCurseRequest({subject: 0, cursesHash: makeCursesHash(makeCurseId(1))});
        s_rmn.unvoteToCurse(reqs);
      }
    }
    // curse should still be in place as only the owner can lift it
    {
      assertTrue(s_rmn.isCursed(0));
    }
    // curse progress should be empty, expect for the cursed flag
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(accWeight, 0);
      assertTrue(cursed);
      assertEq(voters.length, 0);
      assertEq(cursesHashes.length, 0);
    }
    // owner lifts curse
    {
      RMN.OwnerUnvoteToCurseRequest[] memory ownerReq = new RMN.OwnerUnvoteToCurseRequest[](1);
      ownerReq[0] = RMN.OwnerUnvoteToCurseRequest({
        curseVoteAddr: LIFT_CURSE_VOTE_ADDR,
        unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: 0}),
        forceUnvote: false
      });
      vm.prank(OWNER);
      s_rmn.ownerUnvoteToCurse(ownerReq);
    }
    // curse should now be lifted
    {
      assertFalse(s_rmn.isCursed(0));
    }
  }
}

contract RMN_voteToCurse_2 is RMNSetup {
  function initialConfig() internal pure returns (RMN.Config memory) {
    RMN.Config memory cfg = RMN.Config({voters: new RMN.Voter[](3), blessWeightThreshold: 1, curseWeightThreshold: 3});
    cfg.voters[0] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_1, curseVoteAddr: CURSE_VOTER_1, blessWeight: 1, curseWeight: 1});
    cfg.voters[1] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_2, curseVoteAddr: CURSE_VOTER_2, blessWeight: 1, curseWeight: 1});
    cfg.voters[2] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_3, curseVoteAddr: CURSE_VOTER_3, blessWeight: 1, curseWeight: 1});
    return cfg;
  }

  function setUp() public override {
    vm.prank(OWNER);
    s_rmn = new RMN(initialConfig());
  }

  function test_VotesAreDroppedIfSubjectIsNotCursedDuringConfigChange() public {
    // vote to curse the subject from an insufficient number of voters, one voter
    {
      RMN.Config memory cfg = initialConfig();
      vm.prank(cfg.voters[0].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
    }
    // vote must be in place
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 1);
      assertEq(cursesHashes.length, 1);
      assertEq(accWeight, 1);
      assertFalse(cursed);
    }
    // change config to include only the first voter, i.e., initialConfig().voters[0]
    {
      RMN.Config memory cfg = initialConfig();
      RMN.Voter[] memory voters = cfg.voters;
      assembly {
        mstore(voters, 1)
      }
      cfg.curseWeightThreshold = 1;
      vm.prank(OWNER);
      s_rmn.setConfig(cfg);
    }
    // vote must be dropped
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 0);
      assertEq(cursesHashes.length, 0);
      assertEq(accWeight, 0);
      assertFalse(cursed);
    }
    // cause an owner curse now
    {
      vm.prank(OWNER);
      s_rmn.ownerCurse(makeCurseId(1), makeSubjects(0));
    }
    // only the owner curse must be visible
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 1);
      assertEq(voters[0], OWNER_CURSE_VOTE_ADDR);
      assertEq(cursesHashes.length, 1);
      assertEq(cursesHashes[0], makeCursesHash(makeCurseId(1)));
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }
  }

  function test_VotesAreRetainedIfSubjectIsCursedDuringConfigChange() public {
    uint256 numVotersInitially = initialConfig().voters.length;
    // curse the subject with votes from all voters
    {
      RMN.Config memory cfg = initialConfig();
      for (uint256 i = 0; i < cfg.voters.length; ++i) {
        vm.prank(cfg.voters[i].curseVoteAddr);
        s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
      }
    }
    // subject is now cursed
    {
      assertTrue(s_rmn.isCursed(0));
    }
    // throw in an owner curse
    {
      vm.prank(OWNER);
      s_rmn.ownerCurse(makeCurseId(1), makeSubjects(0));
    }

    uint256 snapshot = vm.snapshot();

    for (uint256 keepVoters = 1; keepVoters <= numVotersInitially; ++keepVoters) {
      vm.revertTo(snapshot);

      // change config to include only the first #keepVoters voters, i.e., initialConfig().voters[0..keepVoters]
      {
        RMN.Config memory cfg = initialConfig();
        RMN.Voter[] memory voters = cfg.voters;
        assembly {
          mstore(voters, keepVoters)
        }
        cfg.curseWeightThreshold = uint16(keepVoters);
        vm.prank(OWNER);
        s_rmn.setConfig(cfg);
      }
      // subject is still cursed
      {
        assertTrue(s_rmn.isCursed(0));
      }
      // all votes from the first keepVoters & owner must be present
      {
        (address[] memory voters, bytes28[] memory cursesHashes, uint16 accWeight, bool cursed) =
          s_rmn.getCurseProgress(0);
        assertEq(voters.length, keepVoters + 1 /* owner */ );
        assertEq(cursesHashes.length, keepVoters + 1 /* owner */ );
        assertEq(accWeight, keepVoters /* 1 per voter */ );
        assertTrue(cursed);
        for (uint256 i = 0; i < keepVoters; ++i) {
          assertEq(voters[i], initialConfig().voters[i].curseVoteAddr);
          assertEq(cursesHashes[i], makeCursesHash(makeCurseId(1)));
        }
        assertEq(voters[voters.length - 1], OWNER_CURSE_VOTE_ADDR);
        assertEq(cursesHashes[cursesHashes.length - 1], makeCursesHash(makeCurseId(1)));
      }
      // the owner unvoting for all is not enough to lift the curse, because remember that the owner has an active vote
      // also
      {
        for (uint256 i = 0; i < keepVoters; ++i) {
          RMN.OwnerUnvoteToCurseRequest[] memory ownerReq = new RMN.OwnerUnvoteToCurseRequest[](1);
          ownerReq[0] = RMN.OwnerUnvoteToCurseRequest({
            curseVoteAddr: initialConfig().voters[i].curseVoteAddr,
            unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: makeCursesHash(makeCurseId(1))}),
            forceUnvote: false
          });
          vm.prank(OWNER);
          s_rmn.ownerUnvoteToCurse(ownerReq);

          assertTrue(s_rmn.isCursed(0));
        }
      }
      // after owner unvotes for themselves, finally, the curse will be lifted
      {
        RMN.OwnerUnvoteToCurseRequest[] memory ownerReq = new RMN.OwnerUnvoteToCurseRequest[](1);
        ownerReq[0] = RMN.OwnerUnvoteToCurseRequest({
          curseVoteAddr: OWNER_CURSE_VOTE_ADDR,
          unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: makeCursesHash(makeCurseId(1))}),
          forceUnvote: false
        });
        vm.prank(OWNER);
        s_rmn.ownerUnvoteToCurse(ownerReq);

        assertFalse(s_rmn.isCursed(0));
      }
    }
  }
}

contract RMN_voteToCurse is RMNSetup {
  function _getFirstCurseVoterAndWeight() internal pure returns (address, uint8) {
    RMN.Config memory cfg = rmnConstructorArgs();
    return (cfg.voters[0].curseVoteAddr, cfg.voters[0].curseWeight);
  }

  // Success

  function test_CurseOnlyWhenThresholdReached_Success() public {
    uint256 numSubjects = 3;
    uint256 maxNumRevotes = 2;

    RMN.Config memory cfg = rmnConstructorArgs();
    bytes16[] memory subjects = new bytes16[](numSubjects);
    for (uint256 i = 0; i < numSubjects; ++i) {
      subjects[i] = bytes16(uint128(i));
    }
    for (uint256 numRevotes = 1; numRevotes <= maxNumRevotes; ++numRevotes) {
      // all voters but the last vote, but can't surpass the curse weight threshold
      for (uint256 i = 0; i < cfg.voters.length - 1; ++i) {
        vm.prank(cfg.voters[i].curseVoteAddr);
        s_rmn.voteToCurse(makeCurseId(numRevotes), subjects);
      }
      // no curse is yet active, last voter also needs to vote for any curse to be active
      {
        // ensure every subject is not cursed
        for (uint256 i = 0; i < numSubjects; ++i) {
          assertFalse(s_rmn.isCursed(subjects[i]));
        }
        // ensure every vote has been recorded
        assertEq(
          s_rmn.getRecordedCurseRelatedOpsCount(),
          1 /* setConfig */ + (cfg.voters.length - 1) * numRevotes * numSubjects
        );
      }
    }

    // last voter now votes
    vm.prank(cfg.voters[cfg.voters.length - 1].curseVoteAddr);
    s_rmn.voteToCurse(makeCurseId(0), subjects);
    // curses should be now active
    {
      // ensure every subject is now cursed
      for (uint256 i = 0; i < numSubjects; ++i) {
        assertTrue(s_rmn.isCursed(subjects[i]));
      }
      // ensure every vote has been recorded
      assertEq(
        s_rmn.getRecordedCurseRelatedOpsCount(),
        1 /* setConfig */ + ((cfg.voters.length - 1) * maxNumRevotes + 1) * numSubjects
      );
    }
  }

  function test_VoteToCurse_NoCurse_Success() public {
    (address voter, uint8 weight) = _getFirstCurseVoterAndWeight();
    vm.startPrank(voter);
    vm.expectEmit();
    emit RMN.VotedToCurse(
      1, // configVersion
      voter,
      GLOBAL_CURSE_SUBJECT,
      makeCurseId(123),
      weight,
      1234567890, // blockTimestamp
      makeCursesHash(makeCurseId(123)), // cursesHash
      weight
    );

    s_rmn.voteToCurse(makeCurseId(123), makeSubjects(GLOBAL_CURSE_SUBJECT));

    (address[] memory voters,, uint16 votes, bool cursed) = s_rmn.getCurseProgress(GLOBAL_CURSE_SUBJECT);
    assertEq(1, voters.length);
    assertEq(voter, voters[0]);
    assertEq(weight, votes);
    assertFalse(cursed);
  }

  function test_VoteToCurse_YesCurse_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    for (uint256 i = 0; i < cfg.voters.length - 1; ++i) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
    }

    vm.expectEmit();
    emit RMN.Cursed(1, 0, uint64(block.timestamp));

    vm.startPrank(cfg.voters[cfg.voters.length - 1].curseVoteAddr);
    vm.resumeGasMetering();
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
  }

  function test_EvenIfAlreadyCursed_Success() public {
    RMN.Config memory cfg = rmnConstructorArgs();
    uint16 weightSum = 0;
    for (uint256 i = 0; i < cfg.voters.length; ++i) {
      vm.startPrank(cfg.voters[i].curseVoteAddr);
      s_rmn.voteToCurse(makeCurseId(i), makeSubjects(0));
      weightSum += cfg.voters[i].curseWeight;
    }

    // Not part of the assertion of this test but good to have as a sanity
    // check. We want a curse to be active in order for the ultimate assertion
    // to make sense.
    assert(s_rmn.isCursed(0));

    vm.expectEmit();
    emit RMN.VotedToCurse(
      1, // configVersion
      cfg.voters[cfg.voters.length - 1].curseVoteAddr,
      0, // subject
      makeCurseId(cfg.voters.length + 1), // this curse id
      cfg.voters[cfg.voters.length - 1].curseWeight,
      uint64(block.timestamp), // blockTimestamp
      makeCursesHash(makeCurseId(cfg.voters.length - 1), makeCurseId(cfg.voters.length + 1)), // cursesHash
      weightSum // accumulatedWeight
    );
    // Asserts that this call to vote with a new curse id goes through with no
    // reverts even when the RMN contract is cursed.
    s_rmn.voteToCurse(makeCurseId(cfg.voters.length + 1), makeSubjects(0));
  }

  function test_OwnerCanCurseAndUncurse() public {
    vm.startPrank(OWNER);
    bytes28 expectedCursesHash = makeCursesHash(makeCurseId(0));
    vm.expectEmit();
    emit RMN.VotedToCurse(
      1, // configVersion
      OWNER_CURSE_VOTE_ADDR, // owner
      0, // subject
      makeCurseId(0), // curse id
      0, // weight
      uint64(block.timestamp), // blockTimestamp
      expectedCursesHash, // cursesHash
      0 // accumulatedWeight
    );
    vm.expectEmit();
    emit RMN.Cursed(
      1, // configVersion
      0, // subject
      uint64(block.timestamp) // blockTimestamp
    );
    s_rmn.ownerCurse(makeCurseId(0), makeSubjects(0));

    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint24 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 1);
      assertEq(voters[0], OWNER_CURSE_VOTE_ADDR /* owner */ );
      assertEq(cursesHashes.length, 1);
      assertEq(cursesHashes[0], expectedCursesHash);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    // ownerCurse again, should cause a vote to appear and a change in curses hash
    expectedCursesHash = makeCursesHash(makeCurseId(0), makeCurseId(1));
    vm.expectEmit();
    emit RMN.VotedToCurse(
      1, // configVersion
      OWNER_CURSE_VOTE_ADDR, // owner
      0, // subject
      makeCurseId(1), // curse id
      0, // weight
      uint64(block.timestamp), // blockTimestamp
      expectedCursesHash, // cursesHash
      0 // accumulatedWeight
    );
    s_rmn.ownerCurse(makeCurseId(1), makeSubjects(0));

    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint24 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 1);
      assertEq(voters[0], OWNER_CURSE_VOTE_ADDR /* owner */ );
      assertEq(cursesHashes.length, 1);
      assertEq(cursesHashes[0], expectedCursesHash);
      assertEq(accWeight, 0);
      assertTrue(cursed);
    }

    RMN.OwnerUnvoteToCurseRequest[] memory unvoteReqs = new RMN.OwnerUnvoteToCurseRequest[](1);
    unvoteReqs[0] = RMN.OwnerUnvoteToCurseRequest({
      curseVoteAddr: OWNER_CURSE_VOTE_ADDR,
      unit: RMN.UnvoteToCurseRequest({subject: 0, cursesHash: 0}),
      forceUnvote: true // TODO: test with forceUnvote false also
    });
    vm.expectEmit();
    emit RMN.CurseLifted(0);
    s_rmn.ownerUnvoteToCurse(unvoteReqs);
    {
      (address[] memory voters, bytes28[] memory cursesHashes, uint24 accWeight, bool cursed) =
        s_rmn.getCurseProgress(0);
      assertEq(voters.length, 0);
      assertEq(cursesHashes.length, 0);
      assertEq(accWeight, 0);
      assertFalse(cursed);
    }
  }

  // Reverts

  function test_UnauthorizedVoter_Revert() public {
    vm.startPrank(STRANGER);

    vm.expectRevert(abi.encodeWithSelector(RMN.UnauthorizedVoter.selector, STRANGER));
    s_rmn.voteToCurse(makeCurseId(12312), makeSubjects(0));
  }

  function test_ReusedCurseId_Revert() public {
    (address voter,) = _getFirstCurseVoterAndWeight();
    vm.startPrank(voter);
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));

    vm.expectRevert(abi.encodeWithSelector(RMN.ReusedCurseId.selector, voter, makeCurseId(1)));
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
  }

  function test_RepeatedSubject_Revert() public {
    (address voter,) = _getFirstCurseVoterAndWeight();
    vm.prank(voter);

    bytes16 subject = bytes16(uint128(1));

    vm.expectRevert(RMN.SubjectsMustBeStrictlyIncreasing.selector);
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(subject, subject));
  }

  function test_EmptySubjects_Revert() public {
    (address voter,) = _getFirstCurseVoterAndWeight();
    vm.prank(voter);

    vm.expectRevert(RMN.VoteToCurseNoop.selector);
    s_rmn.voteToCurse(makeCurseId(1), new bytes16[](0));
  }
}

contract RMN_ownerUnvoteToCurse is RMNSetup {
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
        s_rmn.voteToCurse(makeCurseId(j), makeSubjects(GLOBAL_CURSE_SUBJECT));
      }
    }
  }

  function ownerUnvoteToCurse() internal {
    s_rmn.ownerUnvoteToCurse(makeOwnerUnvoteToCurseRequests());
  }

  function makeOwnerUnvoteToCurseRequests() internal pure returns (RMN.OwnerUnvoteToCurseRequest[] memory) {
    (address[] memory cursers,) = getCursersAndCurseCounts();
    RMN.OwnerUnvoteToCurseRequest[] memory reqs = new RMN.OwnerUnvoteToCurseRequest[](cursers.length);
    for (uint256 i = 0; i < cursers.length; ++i) {
      reqs[i] = RMN.OwnerUnvoteToCurseRequest({
        curseVoteAddr: cursers[i],
        unit: RMN.UnvoteToCurseRequest({subject: GLOBAL_CURSE_SUBJECT, cursesHash: bytes28(0)}),
        forceUnvote: true
      });
    }
    return reqs;
  }

  // Success

  function test_OwnerUnvoteToCurseSuccess_gas() public {
    vm.pauseGasMetering();
    vm.startPrank(OWNER);

    vm.expectEmit();
    emit RMN.CurseLifted(GLOBAL_CURSE_SUBJECT);

    vm.resumeGasMetering();
    ownerUnvoteToCurse();
    vm.pauseGasMetering();

    assertFalse(s_rmn.isCursed());
    (address[] memory voters, bytes28[] memory cursesHashes, uint256 weight, bool cursed) =
      s_rmn.getCurseProgress(GLOBAL_CURSE_SUBJECT);
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
    vm.resumeGasMetering();
  }

  function test_IsIdempotent() public {
    vm.startPrank(OWNER);
    ownerUnvoteToCurse();
    vm.expectRevert(RMN.UnvoteToCurseNoop.selector);
    ownerUnvoteToCurse();

    assertFalse(s_rmn.isCursed());
    (address[] memory voters, bytes28[] memory cursesHashes, uint256 weight, bool cursed) =
      s_rmn.getCurseProgress(GLOBAL_CURSE_SUBJECT);
    assertEq(voters.length, 0);
    assertEq(cursesHashes.length, 0);
    assertEq(weight, 0);
    assertFalse(cursed);
  }

  function test_CanBlessAndCurseAfterGlobalCurseIsLifted() public {
    // Contract is already cursed due to setUp.

    // Owner unvotes to curse.
    vm.startPrank(OWNER);
    vm.expectEmit();
    emit RMN.CurseLifted(GLOBAL_CURSE_SUBJECT);
    ownerUnvoteToCurse();

    // Contract is now uncursed.
    assertFalse(s_rmn.isCursed());

    // Vote to bless should go through.
    vm.startPrank(BLESS_VOTER_1);
    s_rmn.voteToBless(makeTaggedRootSingleton(2387489729));

    // Vote to curse should go through.
    vm.startPrank(CURSE_VOTER_1);
    s_rmn.voteToCurse(makeCurseId(73894728973), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }

  // Reverts

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    ownerUnvoteToCurse();
  }

  function test_UnknownVoter_Revert() public {
    vm.stopPrank();
    RMN.OwnerUnvoteToCurseRequest[] memory reqs = new RMN.OwnerUnvoteToCurseRequest[](1);
    reqs[0] = RMN.OwnerUnvoteToCurseRequest({
      curseVoteAddr: STRANGER,
      unit: RMN.UnvoteToCurseRequest({subject: GLOBAL_CURSE_SUBJECT, cursesHash: bytes28(0)}),
      forceUnvote: true
    });

    vm.prank(OWNER);
    vm.expectEmit();
    emit RMN.SkippedUnvoteToCurse(STRANGER, GLOBAL_CURSE_SUBJECT, bytes28(0), bytes28(0));
    vm.expectRevert(RMN.UnvoteToCurseNoop.selector);
    s_rmn.ownerUnvoteToCurse(reqs);

    // no effect on cursedness
    assertTrue(s_rmn.isCursed(GLOBAL_CURSE_SUBJECT));
  }
}

contract RMN_setConfig is ConfigCompare, RMNSetup {
  /// @notice Test-specific function to use only in setConfig tests
  function getDifferentConfigArgs() private pure returns (RMN.Config memory) {
    RMN.Voter[] memory voters = new RMN.Voter[](2);
    voters[0] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_1,
      curseVoteAddr: CURSE_VOTER_1,
      blessWeight: WEIGHT_1,
      curseWeight: WEIGHT_1
    });
    voters[1] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_2,
      curseVoteAddr: CURSE_VOTER_2,
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
    s_rmn.voteToCurse(makeCurseId(1), makeSubjects(0));
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
    vm.expectRevert(abi.encodeWithSelector(RMN.UnauthorizedVoter.selector, BLESS_VOTER_4));
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

    // Assert that curse votes have been cleared

    (address[] memory curseVoters, bytes28[] memory cursesHashes, uint256 curseWeight, bool cursed) =
      s_rmn.getCurseProgress(0);
    assertEq(0, curseVoters.length);
    assertEq(0, cursesHashes.length);
    assertEq(0, curseWeight);
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

contract RMN_permaBlessing is RMNSetup {
  function addresses() private pure returns (address[] memory) {
    return new address[](0);
  }

  function addresses(address a) private pure returns (address[] memory) {
    address[] memory arr = new address[](1);
    arr[0] = a;
    return arr;
  }

  function addresses(address a, address b) private pure returns (address[] memory) {
    address[] memory arr = new address[](2);
    arr[0] = a;
    arr[1] = b;
    return arr;
  }

  function test_PermaBlessing() public {
    bytes32 SOME_ROOT = bytes32(~uint256(0));
    address COMMIT_STORE_1 = makeAddr("COMMIT_STORE_1");
    address COMMIT_STORE_2 = makeAddr("COMMIT_STORE_2");
    IRMN.TaggedRoot memory taggedRootCommitStore1 = IRMN.TaggedRoot({root: SOME_ROOT, commitStore: COMMIT_STORE_1});
    IRMN.TaggedRoot memory taggedRootCommitStore2 = IRMN.TaggedRoot({root: SOME_ROOT, commitStore: COMMIT_STORE_2});

    assertFalse(s_rmn.isBlessed(taggedRootCommitStore1));
    assertFalse(s_rmn.isBlessed(taggedRootCommitStore2));
    assertEq(s_rmn.getPermaBlessedCommitStores(), addresses());

    // only owner can mutate permaBlessedCommitStores
    vm.prank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_rmn.ownerRemoveThenAddPermaBlessedCommitStores(addresses(), addresses(COMMIT_STORE_1));

    vm.prank(OWNER);
    s_rmn.ownerRemoveThenAddPermaBlessedCommitStores(addresses(), addresses(COMMIT_STORE_1));
    assertTrue(s_rmn.isBlessed(taggedRootCommitStore1));
    assertFalse(s_rmn.isBlessed(taggedRootCommitStore2));
    assertEq(s_rmn.getPermaBlessedCommitStores(), addresses(COMMIT_STORE_1));

    vm.prank(OWNER);
    s_rmn.ownerRemoveThenAddPermaBlessedCommitStores(addresses(COMMIT_STORE_1), addresses(COMMIT_STORE_2));
    assertFalse(s_rmn.isBlessed(taggedRootCommitStore1));
    assertTrue(s_rmn.isBlessed(taggedRootCommitStore2));
    assertEq(s_rmn.getPermaBlessedCommitStores(), addresses(COMMIT_STORE_2));

    vm.prank(OWNER);
    s_rmn.ownerRemoveThenAddPermaBlessedCommitStores(addresses(), addresses(COMMIT_STORE_1));
    assertTrue(s_rmn.isBlessed(taggedRootCommitStore1));
    assertTrue(s_rmn.isBlessed(taggedRootCommitStore2));
    assertEq(s_rmn.getPermaBlessedCommitStores(), addresses(COMMIT_STORE_2, COMMIT_STORE_1));

    vm.prank(OWNER);
    s_rmn.ownerRemoveThenAddPermaBlessedCommitStores(addresses(COMMIT_STORE_1, COMMIT_STORE_2), addresses());
    assertFalse(s_rmn.isBlessed(taggedRootCommitStore1));
    assertFalse(s_rmn.isBlessed(taggedRootCommitStore2));
    assertEq(s_rmn.getPermaBlessedCommitStores(), addresses());
  }
}

contract RMN_getRecordedCurseRelatedOps is RMNSetup {
  function test_OpsPostDeployment() public view {
    // The constructor call includes a setConfig, so that's the only thing we should expect to find.
    assertEq(s_rmn.getRecordedCurseRelatedOpsCount(), 1);
    RMN.RecordedCurseRelatedOp[] memory recordedCurseRelatedOps = s_rmn.getRecordedCurseRelatedOps(0, type(uint256).max);
    assertEq(recordedCurseRelatedOps.length, 1);
    assertEq(uint8(recordedCurseRelatedOps[0].tag), uint8(RMN.RecordedCurseRelatedOpTag.SetConfig));
  }
}
