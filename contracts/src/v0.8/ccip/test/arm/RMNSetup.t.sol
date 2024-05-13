// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

import {RMN} from "../../RMN.sol";
import {MockRMN} from "../mocks/MockRMN.sol";
import {Test} from "forge-std/Test.sol";

contract RMNSetup is Test {
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant STRANGER = address(999999);
  address internal constant ZERO_ADDRESS = address(0);

  address internal constant BLESS_VOTER_1 = address(1);
  address internal constant CURSE_VOTER_1 = address(10);
  address internal constant CURSE_UNVOTER_1 = address(110);
  address internal constant BLESS_VOTER_2 = address(2);
  address internal constant CURSE_VOTER_2 = address(12);
  address internal constant CURSE_UNVOTER_2 = address(112);
  address internal constant BLESS_VOTER_3 = address(3);
  address internal constant CURSE_VOTER_3 = address(13);
  address internal constant CURSE_UNVOTER_3 = address(113);
  address internal constant BLESS_VOTER_4 = address(4);
  address internal constant CURSE_VOTER_4 = address(14);
  address internal constant CURSE_UNVOTER_4 = address(114);

  uint8 internal constant ZERO = 0;
  uint8 internal constant WEIGHT_1 = 1;
  uint8 internal constant WEIGHT_10 = 10;
  uint8 internal constant WEIGHT_20 = 20;
  uint8 internal constant WEIGHT_40 = 40;

  function makeTaggedRootsInclusive(uint256 from, uint256 to) internal pure returns (IRMN.TaggedRoot[] memory) {
    IRMN.TaggedRoot[] memory votes = new IRMN.TaggedRoot[](to - from + 1);
    for (uint256 i = from; i <= to; ++i) {
      votes[i - from] = IRMN.TaggedRoot({commitStore: address(1), root: bytes32(uint256(i))});
    }
    return votes;
  }

  function makeTaggedRootSingleton(uint256 index) internal pure returns (IRMN.TaggedRoot[] memory) {
    return makeTaggedRootsInclusive(index, index);
  }

  function makeTaggedRoot(uint256 index) internal pure returns (IRMN.TaggedRoot memory) {
    return makeTaggedRootSingleton(index)[0];
  }

  function makeTaggedRootHash(uint256 index) internal pure returns (bytes32) {
    IRMN.TaggedRoot memory taggedRoot = makeTaggedRootSingleton(index)[0];
    return keccak256(abi.encode(taggedRoot.commitStore, taggedRoot.root));
  }

  function makeCurseId(uint256 index) internal pure returns (bytes32) {
    return bytes32(index);
  }

  MockRMN internal s_mockRMN;
  RMN internal s_rmn;

  function setUp() public virtual {
    vm.startPrank(OWNER);
    s_rmn = new RMN(rmnConstructorArgs());
    s_mockRMN = new MockRMN();
    vm.stopPrank();
  }

  function rmnConstructorArgs() internal pure returns (RMN.Config memory) {
    RMN.Voter[] memory voters = new RMN.Voter[](4);
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
    voters[2] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_3,
      curseVoteAddr: CURSE_VOTER_3,
      curseUnvoteAddr: CURSE_UNVOTER_3,
      blessWeight: WEIGHT_20,
      curseWeight: WEIGHT_20
    });
    voters[3] = RMN.Voter({
      blessVoteAddr: BLESS_VOTER_4,
      curseVoteAddr: CURSE_VOTER_4,
      curseUnvoteAddr: CURSE_UNVOTER_4,
      blessWeight: WEIGHT_40,
      curseWeight: WEIGHT_40
    });
    return RMN.Config({
      voters: voters,
      blessWeightThreshold: WEIGHT_10 + WEIGHT_20 + WEIGHT_40,
      curseWeightThreshold: WEIGHT_1 + WEIGHT_10 + WEIGHT_20 + WEIGHT_40
    });
  }

  function hasVotedToBlessRoot(address voter, IRMN.TaggedRoot memory taggedRoot_) internal view returns (bool) {
    (address[] memory voters,,) = s_rmn.getBlessProgress(taggedRoot_);
    for (uint256 i = 0; i < voters.length; ++i) {
      if (voters[i] == voter) {
        return true;
      }
    }
    return false;
  }

  function getWeightOfVotesToBlessRoot(IRMN.TaggedRoot memory taggedRoot_) internal view returns (uint16) {
    (, uint16 weight,) = s_rmn.getBlessProgress(taggedRoot_);
    return weight;
  }
}
