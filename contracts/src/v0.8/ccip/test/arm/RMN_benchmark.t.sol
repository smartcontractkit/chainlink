// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {GLOBAL_CURSE_SUBJECT, OWNER_CURSE_VOTE_ADDR, RMN} from "../../RMN.sol";
import {RMNSetup, makeCursesHash, makeSubjects} from "./RMNSetup.t.sol";

contract RMN_voteToBless_Benchmark is RMNSetup {
  function test_RootSuccess_gas(uint256 n) internal {
    vm.prank(BLESS_VOTER_1);
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, n));
  }

  function test_1RootSuccess_gas() public {
    test_RootSuccess_gas(1);
  }

  function test_3RootSuccess_gas() public {
    test_RootSuccess_gas(3);
  }

  function test_5RootSuccess_gas() public {
    test_RootSuccess_gas(5);
  }
}

contract RMN_voteToBless_Blessed_Benchmark is RMN_voteToBless_Benchmark {
  function setUp() public virtual override {
    RMNSetup.setUp();
    vm.prank(BLESS_VOTER_2);
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, 1));
    vm.prank(BLESS_VOTER_3);
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, 1));
  }

  function test_1RootSuccessBecameBlessed_gas() public {
    vm.prank(BLESS_VOTER_4);
    s_rmn.voteToBless(makeTaggedRootsInclusive(1, 1));
  }
}

abstract contract RMN_voteToCurse_Benchmark is RMNSetup {
  struct PreVote {
    address voter;
    bytes16 subject;
  }

  PreVote[] internal s_preVotes;

  function setUp() public virtual override {
    // Intentionally does not inherit RMNSetup setUp(), because we set up a simpler config here.
    // The only way to ensure that storage slots are cold for the actual functions to be benchmarked is to perform the
    // setup in setUp().

    RMN.Config memory cfg = RMN.Config({voters: new RMN.Voter[](3), blessWeightThreshold: 3, curseWeightThreshold: 3});
    cfg.voters[0] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_1, curseVoteAddr: CURSE_VOTER_1, blessWeight: 1, curseWeight: 1});
    cfg.voters[1] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_2, curseVoteAddr: CURSE_VOTER_2, blessWeight: 1, curseWeight: 1});
    cfg.voters[2] =
      RMN.Voter({blessVoteAddr: BLESS_VOTER_3, curseVoteAddr: CURSE_VOTER_3, blessWeight: 1, curseWeight: 1});
    vm.prank(OWNER);
    s_rmn = new RMN(cfg);

    for (uint256 i = 0; i < s_preVotes.length; ++i) {
      vm.prank(s_preVotes[i].voter);
      s_rmn.voteToCurse(makeCurseId(i), makeSubjects(s_preVotes[i].subject));
    }
  }
}

contract RMN_voteToCurse_Benchmark_1 is RMN_voteToCurse_Benchmark {
  constructor() {
    // some irrelevant subject & voter so that we don't pay for the nonzero->zero SSTORE of
    // s_recordedVotesToCurse.length in the benchmark below
    s_preVotes.push(PreVote({voter: CURSE_VOTER_3, subject: bytes16(~uint128(0))}));
  }

  function test_VoteToCurse_NewSubject_NewVoter_NoCurse_gas() public {
    vm.prank(CURSE_VOTER_1);
    s_rmn.voteToCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }

  function test_VoteToCurse_NewSubject_NewVoter_YesCurse_gas() public {
    vm.prank(OWNER);
    s_rmn.ownerCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }
}

contract RMN_voteToCurse_Benchmark_2 is RMN_voteToCurse_Benchmark {
  constructor() {
    s_preVotes.push(PreVote({voter: CURSE_VOTER_1, subject: GLOBAL_CURSE_SUBJECT}));
  }

  function test_VoteToCurse_OldSubject_OldVoter_NoCurse_gas() public {
    vm.prank(CURSE_VOTER_1);
    s_rmn.voteToCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }

  function test_VoteToCurse_OldSubject_NewVoter_NoCurse_gas() public {
    vm.prank(CURSE_VOTER_2);
    s_rmn.voteToCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }
}

contract RMN_voteToCurse_Benchmark_3 is RMN_voteToCurse_Benchmark {
  constructor() {
    s_preVotes.push(PreVote({voter: CURSE_VOTER_1, subject: GLOBAL_CURSE_SUBJECT}));
    s_preVotes.push(PreVote({voter: CURSE_VOTER_2, subject: GLOBAL_CURSE_SUBJECT}));
  }

  function test_VoteToCurse_OldSubject_NewVoter_YesCurse_gas() public {
    vm.prank(CURSE_VOTER_3);
    s_rmn.voteToCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }
}

contract RMN_lazyVoteToCurseUpdate_Benchmark is RMN_voteToCurse_Benchmark {
  constructor() {
    s_preVotes.push(PreVote({voter: CURSE_VOTER_1, subject: GLOBAL_CURSE_SUBJECT}));
    s_preVotes.push(PreVote({voter: CURSE_VOTER_2, subject: GLOBAL_CURSE_SUBJECT}));
    s_preVotes.push(PreVote({voter: CURSE_VOTER_3, subject: GLOBAL_CURSE_SUBJECT}));
  }

  function setUp() public override {
    RMN_voteToCurse_Benchmark.setUp(); // sends the prevotes
    // initial config includes voters CURSE_VOTER_1, CURSE_VOTER_2, CURSE_VOTER_3
    // include a new voter in the config
    {
      (,, RMN.Config memory cfg) = s_rmn.getConfigDetails();
      RMN.Voter[] memory newVoters = new RMN.Voter[](cfg.voters.length + 1);
      for (uint256 i = 0; i < cfg.voters.length; ++i) {
        newVoters[i] = cfg.voters[i];
      }
      newVoters[newVoters.length - 1] =
        RMN.Voter({blessVoteAddr: BLESS_VOTER_4, curseVoteAddr: CURSE_VOTER_4, blessWeight: 1, curseWeight: 1});
      cfg.voters = newVoters;

      vm.prank(OWNER);
      s_rmn.setConfig(cfg);
    }
  }

  function test_VoteToCurseLazilyRetain3VotersUponConfigChange_gas() public {
    // send a vote as the new voter, should cause a lazy update and votes from CURSE_VOTER_1, CURSE_VOTER_2,
    // CURSE_VOTER_3 to be retained, which is the worst case for the prior config
    vm.prank(CURSE_VOTER_4);
    s_rmn.voteToCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }
}

contract RMN_setConfig_Benchmark is RMNSetup {
  uint256 s_numVoters;

  function configWithVoters(uint256 numVoters) internal pure returns (RMN.Config memory) {
    RMN.Config memory cfg =
      RMN.Config({voters: new RMN.Voter[](numVoters), blessWeightThreshold: 1, curseWeightThreshold: 1});
    for (uint256 i = 1; i <= numVoters; ++i) {
      cfg.voters[i - 1] = RMN.Voter({
        blessVoteAddr: address(uint160(2 * i)),
        curseVoteAddr: address(uint160(2 * i + 1)),
        blessWeight: 1,
        curseWeight: 1
      });
    }
    return cfg;
  }

  function setUp() public virtual override {
    vm.prank(OWNER);
    s_rmn = new RMN(configWithVoters(s_numVoters));
  }
}

contract RMN_setConfig_Benchmark_1 is RMN_setConfig_Benchmark {
  constructor() {
    s_numVoters = 1;
  }

  function test_SetConfig_7Voters_gas() public {
    vm.prank(OWNER);
    s_rmn.setConfig(configWithVoters(7));
  }
}

contract RMN_setConfig_Benchmark_2 is RMN_setConfig_Benchmark {
  constructor() {
    s_numVoters = 7;
  }

  function test_ResetConfig_7Voters_gas() public {
    vm.prank(OWNER);
    s_rmn.setConfig(configWithVoters(7));
  }
}

contract RMN_ownerUnvoteToCurse_Benchmark is RMN_setConfig_Benchmark {
  constructor() {
    s_numVoters = 7;
  }

  function setUp() public override {
    RMN_setConfig_Benchmark.setUp();
    vm.prank(OWNER);
    s_rmn.ownerCurse(makeCurseId(0xffff), makeSubjects(GLOBAL_CURSE_SUBJECT));
  }

  function test_OwnerUnvoteToCurse_1Voter_LiftsCurse_gas() public {
    RMN.OwnerUnvoteToCurseRequest[] memory reqs = new RMN.OwnerUnvoteToCurseRequest[](1);
    reqs[0] = RMN.OwnerUnvoteToCurseRequest({
      curseVoteAddr: OWNER_CURSE_VOTE_ADDR,
      unit: RMN.UnvoteToCurseRequest({cursesHash: makeCursesHash(makeCurseId(0xffff)), subject: GLOBAL_CURSE_SUBJECT}),
      forceUnvote: false
    });
    vm.prank(OWNER);
    s_rmn.ownerUnvoteToCurse(reqs);
  }
}
