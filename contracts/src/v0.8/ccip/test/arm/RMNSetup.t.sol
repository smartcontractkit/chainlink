// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

import {RMN} from "../../RMN.sol";
import {BaseTest} from "../BaseTest.t.sol";

contract RMNSetup is BaseTest {
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

  RMN internal s_rmn;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_rmn = new RMN(rmnConstructorArgs());
    vm.stopPrank();
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
