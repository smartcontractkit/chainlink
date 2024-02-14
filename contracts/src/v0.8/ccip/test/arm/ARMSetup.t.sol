// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IARM} from "../../interfaces/IARM.sol";

import {BaseTest} from "../BaseTest.t.sol";
import {ARM} from "../../ARM.sol";

contract ARMSetup is BaseTest {
  function makeTaggedRootsInclusive(uint256 from, uint256 to) internal pure returns (IARM.TaggedRoot[] memory) {
    IARM.TaggedRoot[] memory votes = new IARM.TaggedRoot[](to - from + 1);
    for (uint256 i = from; i <= to; ++i) {
      votes[i - from] = IARM.TaggedRoot({commitStore: address(1), root: bytes32(uint256(i))});
    }
    return votes;
  }

  function makeTaggedRootSingleton(uint256 index) internal pure returns (IARM.TaggedRoot[] memory) {
    return makeTaggedRootsInclusive(index, index);
  }

  function makeTaggedRoot(uint256 index) internal pure returns (IARM.TaggedRoot memory) {
    return makeTaggedRootSingleton(index)[0];
  }

  function makeTaggedRootHash(uint256 index) internal pure returns (bytes32) {
    IARM.TaggedRoot memory taggedRoot = makeTaggedRootSingleton(index)[0];
    return keccak256(abi.encode(taggedRoot.commitStore, taggedRoot.root));
  }

  function makeCurseId(uint256 index) internal pure returns (bytes32) {
    return bytes32(index);
  }

  ARM internal s_arm;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_arm = new ARM(armConstructorArgs());
  }

  function hasVotedToBlessRoot(address voter, IARM.TaggedRoot memory taggedRoot_) internal view returns (bool) {
    (address[] memory voters, , ) = s_arm.getBlessProgress(taggedRoot_);
    for (uint256 i = 0; i < voters.length; ++i) {
      if (voters[i] == voter) {
        return true;
      }
    }
    return false;
  }

  function getWeightOfVotesToBlessRoot(IARM.TaggedRoot memory taggedRoot_) internal view returns (uint16) {
    (, uint16 weight, ) = s_arm.getBlessProgress(taggedRoot_);
    return weight;
  }
}
