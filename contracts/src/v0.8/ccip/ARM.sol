// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IARM} from "./interfaces/IARM.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";

contract ARM is IARM, OwnerIsCreator, TypeAndVersionInterface {
  // STATIC CONFIG
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ARM 1.0.0";

  uint256 private constant MAX_NUM_VOTERS = 128;

  // DYNAMIC CONFIG
  struct Voter {
    // This is the address the voter should use to call voteToBless.
    address blessVoteAddr;
    // This is the address the voter should use to call voteToCurse.
    address curseVoteAddr;
    // This is the address the voter should use to call unvoteToCurse.
    address curseUnvoteAddr;
    // The weight of this voter's vote for blessing.
    uint8 blessWeight;
    // The weight of this voter's vote for cursing.
    uint8 curseWeight;
  }

  struct Config {
    Voter[] voters;
    // When the total weight of voters that have voted to bless a tagged root reaches
    // or exceeds blessWeightThreshold, the tagged root becomes blessed.
    uint16 blessWeightThreshold;
    // When the total weight of voters that have voted to curse reaches or
    // exceeds curseWeightThreshold, the ARM enters the cursed state.
    uint16 curseWeightThreshold;
  }

  struct VersionedConfig {
    Config config;
    // The version is incremented every time the config changes.
    // The initial configuration on the contract will have configVersion == 1.
    uint32 configVersion;
    // The block number at which the config was last set. Helps the offchain
    // code check that the config was set in a stable block or double-check
    // that it has the correct config by querying logs at that block number.
    uint32 blockNumber;
  }

  VersionedConfig private s_versionedConfig;

  // STATE
  struct BlesserRecord {
    // The config version at which this BlesserRecord was last set. A blesser
    // is considered active iff this configVersion equals
    // s_versionedConfig.configVersion.
    uint32 configVersion;
    uint8 weight;
    uint8 index;
  }

  mapping(address blessVoteAddr => BlesserRecord blesserRecord) private s_blesserRecords;

  struct BlessVoteProgress {
    // A BlessVoteProgress is considered invalid if weightThresholdMet is false when
    // s_versionedConfig.configVersion changes. we don't want old in-progress
    // votes to continue when we set a new config!
    // The config version at which the bless vote for a tagged root was initiated.
    uint32 configVersion;
    uint16 accumulatedWeight;
    // Care must be taken that the bitmap has as many bits as MAX_NUM_VOTERS.
    uint128 voterBitmap;
    bool weightThresholdMet;
  }

  mapping(bytes32 taggedRootHash => BlessVoteProgress blessVoteProgress) private s_blessVoteProgressByTaggedRootHash;

  // voteCount and cursesHash can be reset through unvoteToCurse, and ownerUnvoteToCurse, and may be reset through
  // setConfig if the curser is not part of the new config.
  struct CurserRecord {
    bool active;
    uint8 weight;
    uint32 voteCount;
    address curseUnvoteAddr;
    bytes32 cursesHash;
  }

  mapping(address curseVoteAddr => CurserRecord curserRecord) private s_curserRecords;

  // Maintains a per-curser set of curseIds. Entries from this mapping are
  // never cleared. Once a curseId is used it can never be reused, even after
  // an unvoteToCurse or ownerUnvoteToCurse. This is to prevent accidental
  // re-votes to curse, e.g. caused by TOCTOU issues.
  mapping(address curseVoteAddr => mapping(bytes32 curseId => bool voted)) private s_curseVotes;

  struct CurseVoteProgress {
    uint16 curseWeightThreshold;
    uint16 accumulatedWeight;
    // A curse becomes active after:
    // - accumulatedWeight becomes greater or equal than curseWeightThreshold; or
    // - the owner curses.
    // Once a curse is active, only the owner can lift it.
    bool curseActive;
  }

  CurseVoteProgress private s_curseVoteProgress;

  // AUXILLARY STRUCTS
  struct UnvoteToCurseRecord {
    address curseVoteAddr;
    bytes32 cursesHash;
    bool forceUnvote;
  }

  // EVENTS, ERRORS
  event ConfigSet(uint32 indexed configVersion, Config config);
  error InvalidConfig();

  event TaggedRootBlessed(uint32 indexed configVersion, IARM.TaggedRoot taggedRoot, uint16 accumulatedWeight);
  event TaggedRootBlessVotesReset(uint32 indexed configVersion, IARM.TaggedRoot taggedRoot, bool wasBlessed);
  event VotedToBless(uint32 indexed configVersion, address indexed voter, IARM.TaggedRoot taggedRoot, uint8 weight);

  event VotedToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    uint8 weight,
    uint32 voteCount,
    bytes32 curseId,
    bytes32 cursesHash,
    uint16 accumulatedWeight
  );
  event ReusedVotesToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    uint8 weight,
    uint32 voteCount,
    bytes32 cursesHash,
    uint16 accumulatedWeight
  );
  event UnvotedToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    uint8 weight,
    uint32 voteCount,
    bytes32 cursesHash
  );
  event SkippedUnvoteToCurse(address indexed voter, bytes32 expectedCursesHash, bytes32 actualCursesHash);
  event OwnerCursed(uint256 timestamp);
  event Cursed(uint32 indexed configVersion, uint256 timestamp);

  // These events make it easier for offchain logic to discover that it performs
  // the same actions multiple times.
  event AlreadyVotedToBless(uint32 indexed configVersion, address indexed voter, IARM.TaggedRoot taggedRoot);
  event AlreadyBlessed(uint32 indexed configVersion, address indexed voter, IARM.TaggedRoot taggedRoot);

  event RecoveredFromCurse();

  error AlreadyVotedToCurse(address voter, bytes32 curseId);
  error InvalidVoter(address voter);
  error InvalidCurseState();
  error InvalidCursesHash(bytes32 expectedCursesHash, bytes32 actualCursesHash);
  error MustRecoverFromCurse();

  constructor(Config memory config) {
    {
      // Ensure that the bitmap is large enough to hold MAX_NUM_VOTERS.
      // We do this in the constructor because MAX_NUM_VOTERS is constant.
      BlessVoteProgress memory vp;
      vp.voterBitmap = ~uint128(0);
      assert(vp.voterBitmap >> (MAX_NUM_VOTERS - 1) >= 1);
    }
    _setConfig(config);
  }

  function _bitmapGet(uint128 bitmap, uint8 index) internal pure returns (bool) {
    assert(index < MAX_NUM_VOTERS);
    return bitmap & (uint128(1) << index) != 0;
  }

  function _bitmapSet(uint128 bitmap, uint8 index) internal pure returns (uint128) {
    assert(index < MAX_NUM_VOTERS);
    return bitmap | (uint128(1) << index);
  }

  function _bitmapCount(uint128 bitmap) internal pure returns (uint8 oneBits) {
    // https://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
    for (; bitmap != 0; ++oneBits) {
      bitmap &= bitmap - 1;
    }
  }

  function _taggedRootHash(IARM.TaggedRoot memory taggedRoot) internal pure returns (bytes32) {
    return keccak256(abi.encode(taggedRoot.commitStore, taggedRoot.root));
  }

  /// @param taggedRoots A tagged root is hashed as `keccak256(abi.encode(taggedRoot.commitStore
  /// /* address */, taggedRoot.root /* bytes32 */))`.
  function voteToBless(IARM.TaggedRoot[] calldata taggedRoots) external {
    // If we have an active curse, something is really wrong. Let's err on the
    // side of caution and not accept further blessings during this time of
    // uncertainty.
    if (isCursed()) revert MustRecoverFromCurse();

    uint32 configVersion = s_versionedConfig.configVersion;
    BlesserRecord memory blesserRecord = s_blesserRecords[msg.sender];
    if (blesserRecord.configVersion != configVersion) revert InvalidVoter(msg.sender);

    for (uint256 i = 0; i < taggedRoots.length; ++i) {
      IARM.TaggedRoot memory taggedRoot = taggedRoots[i];
      bytes32 taggedRootHash = _taggedRootHash(taggedRoot);
      BlessVoteProgress memory voteProgress = s_blessVoteProgressByTaggedRootHash[taggedRootHash];
      if (voteProgress.weightThresholdMet) {
        // We don't revert here because it's unreasonable to expect from the
        // voter to know exactly when to stop voting. Most likely when they
        // voted they didn't realize the threshold would be reached by the time
        // their vote was counted.
        // Additionally, there might be other tagged roots for which votes might
        // count, and we want to allow that to happen.
        emit AlreadyBlessed(configVersion, msg.sender, taggedRoot);
        continue;
      }
      if (voteProgress.configVersion != configVersion) {
        // Note that voteProgress.weightThresholdMet must be false at this point

        // If votes were received while an older config was in effect,
        // invalidate them and start from scratch.
        // If votes were never received, set the current config version.
        voteProgress = BlessVoteProgress({
          configVersion: configVersion,
          voterBitmap: 0,
          accumulatedWeight: 0,
          weightThresholdMet: false
        });
      }
      if (_bitmapGet(voteProgress.voterBitmap, blesserRecord.index)) {
        // We don't revert here because there might be other tagged roots for
        // which votes might count, and we want to allow that to happen.
        emit AlreadyVotedToBless(configVersion, msg.sender, taggedRoot);
        continue;
      }
      voteProgress.voterBitmap = _bitmapSet(voteProgress.voterBitmap, blesserRecord.index);
      voteProgress.accumulatedWeight += blesserRecord.weight;
      emit VotedToBless(configVersion, msg.sender, taggedRoot, blesserRecord.weight);
      if (voteProgress.accumulatedWeight >= s_versionedConfig.config.blessWeightThreshold) {
        voteProgress.weightThresholdMet = true;
        emit TaggedRootBlessed(configVersion, taggedRoot, voteProgress.accumulatedWeight);
      }
      s_blessVoteProgressByTaggedRootHash[taggedRootHash] = voteProgress;
    }
  }

  /// @notice Can be called by the owner to remove unintentionally voted or even blessed tagged roots in a recovery
  /// scenario. The owner must ensure that there are no in-flight transactions by ARM nodes voting for any of the
  /// taggedRoots before calling this function, as such in-flight transactions could lead to the roots becoming
  /// re-blessed shortly after the call to this function, contrary to the original intention.
  function ownerResetBlessVotes(IARM.TaggedRoot[] calldata taggedRoots) external onlyOwner {
    uint32 configVersion = s_versionedConfig.configVersion;
    for (uint256 i = 0; i < taggedRoots.length; ++i) {
      IARM.TaggedRoot memory taggedRoot = taggedRoots[i];
      bytes32 taggedRootHash = _taggedRootHash(taggedRoot);
      BlessVoteProgress memory voteProgress = s_blessVoteProgressByTaggedRootHash[taggedRootHash];
      delete s_blessVoteProgressByTaggedRootHash[taggedRootHash];
      bool wasBlessed = voteProgress.weightThresholdMet;
      if (voteProgress.configVersion == configVersion || wasBlessed) {
        emit TaggedRootBlessVotesReset(configVersion, taggedRoot, wasBlessed);
      }
    }
  }

  /// @notice Can be called by a curser to remove unintentional votes to curse.
  /// We expect this to be called very rarely, e.g. in case of a bug in the
  /// offchain code causing false voteToCurse calls.
  /// @notice Should be called from curser's corresponding curseUnvoteAddr.
  function unvoteToCurse(address curseVoteAddr, bytes32 cursesHash) external {
    CurserRecord memory curserRecord = s_curserRecords[curseVoteAddr];

    // If a curse is active, only the owner is allowed to lift it.
    if (isCursed()) revert MustRecoverFromCurse();

    if (msg.sender != curserRecord.curseUnvoteAddr) revert InvalidVoter(msg.sender);

    if (!curserRecord.active || curserRecord.voteCount == 0) revert InvalidCurseState();
    if (curserRecord.cursesHash != cursesHash) revert InvalidCursesHash(curserRecord.cursesHash, cursesHash);

    emit UnvotedToCurse(
      s_versionedConfig.configVersion,
      curseVoteAddr,
      curserRecord.weight,
      curserRecord.voteCount,
      cursesHash
    );
    curserRecord.voteCount = 0;
    curserRecord.cursesHash = 0;
    s_curserRecords[curseVoteAddr] = curserRecord;
    s_curseVoteProgress.accumulatedWeight -= curserRecord.weight;
  }

  /// @notice A vote to curse is appropriate during unhealthy blockchain conditions
  /// (eg. finality violations).
  function voteToCurse(bytes32 curseId) external {
    CurserRecord memory curserRecord = s_curserRecords[msg.sender];
    if (!curserRecord.active) revert InvalidVoter(msg.sender);
    if (s_curseVotes[msg.sender][curseId]) revert AlreadyVotedToCurse(msg.sender, curseId);
    s_curseVotes[msg.sender][curseId] = true;
    ++curserRecord.voteCount;
    curserRecord.cursesHash = keccak256(abi.encode(curserRecord.cursesHash, curseId));
    s_curserRecords[msg.sender] = curserRecord;

    CurseVoteProgress memory curseVoteProgress = s_curseVoteProgress;

    if (curserRecord.voteCount == 1) {
      curseVoteProgress.accumulatedWeight += curserRecord.weight;
    }

    // NOTE: We could pack configVersion into CurserRecord that we already load in the beginning of this function to
    // avoid the following extra storage read for it, but since voteToCurse is not on the hot path we'd rather keep
    // things simple.
    uint32 configVersion = s_versionedConfig.configVersion;
    emit VotedToCurse(
      configVersion,
      msg.sender,
      curserRecord.weight,
      curserRecord.voteCount,
      curseId,
      curserRecord.cursesHash,
      curseVoteProgress.accumulatedWeight
    );
    if (
      !curseVoteProgress.curseActive && curseVoteProgress.accumulatedWeight >= curseVoteProgress.curseWeightThreshold
    ) {
      curseVoteProgress.curseActive = true;
      emit Cursed(configVersion, block.timestamp);
    }
    s_curseVoteProgress = curseVoteProgress;
  }

  /// @notice Enables the owner to immediately have the system enter the cursed state.
  function ownerCurse() external onlyOwner {
    emit OwnerCursed(block.timestamp);
    if (!s_curseVoteProgress.curseActive) {
      s_curseVoteProgress.curseActive = true;
      emit Cursed(s_versionedConfig.configVersion, block.timestamp);
    }
  }

  /// @notice Enables the owner to remove curse votes. After the curse votes are removed,
  /// this function will check whether the curse is still valid and restore the uncursed state if possible.
  /// This function also enables the owner to lift a curse created through ownerCurse.
  function ownerUnvoteToCurse(UnvoteToCurseRecord[] calldata unvoteRecords) external onlyOwner {
    for (uint256 i = 0; i < unvoteRecords.length; ++i) {
      UnvoteToCurseRecord memory unvoteRecord = unvoteRecords[i];
      CurserRecord memory curserRecord = s_curserRecords[unvoteRecord.curseVoteAddr];
      // Owner can avoid the curses hash check by setting forceUnvote to true, in case
      // a malicious curser is flooding the system with votes to curse with the
      // intention to disallow the owner to clear their curse.
      if (!unvoteRecord.forceUnvote && curserRecord.cursesHash != unvoteRecord.cursesHash) {
        emit SkippedUnvoteToCurse(unvoteRecord.curseVoteAddr, curserRecord.cursesHash, unvoteRecord.cursesHash);
        continue;
      }

      if (!curserRecord.active || curserRecord.voteCount == 0) continue;

      emit UnvotedToCurse(
        s_versionedConfig.configVersion,
        unvoteRecord.curseVoteAddr,
        curserRecord.weight,
        curserRecord.voteCount,
        curserRecord.cursesHash
      );
      curserRecord.voteCount = 0;
      curserRecord.cursesHash = 0;
      s_curserRecords[unvoteRecord.curseVoteAddr] = curserRecord;
      s_curseVoteProgress.accumulatedWeight -= curserRecord.weight;
    }

    if (
      s_curseVoteProgress.curseActive &&
      s_curseVoteProgress.accumulatedWeight < s_curseVoteProgress.curseWeightThreshold
    ) {
      s_curseVoteProgress.curseActive = false;
      emit RecoveredFromCurse();
      // Invalidate all in-progress votes to bless by bumping the config version.
      // They might have been based on false information about the source chain
      // (e.g. in case of a finality violation).
      _setConfig(s_versionedConfig.config);
    }
  }

  /// @notice Will revert in case a curse is active. To avoid accidentally invalidating an in-progress curse vote, it
  /// may be advisable to remove voters one-by-one over time, rather than many at once.
  /// @dev The gas use of this function varies depending on the number of curse votes that are active. When calling this
  /// function, be sure to include a gas cushion to account for curse votes that may occur between your transaction
  /// being sent and mined.
  function setConfig(Config memory config) external onlyOwner {
    _setConfig(config);
  }

  /// @inheritdoc IARM
  function isBlessed(IARM.TaggedRoot calldata taggedRoot) external view override returns (bool) {
    return s_blessVoteProgressByTaggedRootHash[_taggedRootHash(taggedRoot)].weightThresholdMet;
  }

  /// @inheritdoc IARM
  function isCursed() public view override returns (bool) {
    return s_curseVoteProgress.curseActive;
  }

  /// @notice Config version might be incremented for many reasons, including
  /// recovery from a curse and a regular config change.
  function getConfigDetails() external view returns (uint32 version, uint32 blockNumber, Config memory config) {
    version = s_versionedConfig.configVersion;
    blockNumber = s_versionedConfig.blockNumber;
    config = s_versionedConfig.config;
  }

  /// @return blessVoteAddrs addresses of voters, will be empty if voting took place with an older config version
  /// @return accumulatedWeight sum of weights of voters, will be zero if voting took place with an older config version
  /// @return blessed will be accurate regardless of when voting took place
  /// @dev This is a helper method for offchain code so efficiency is not really a concern.
  function getBlessProgress(
    IARM.TaggedRoot calldata taggedRoot
  ) external view returns (address[] memory blessVoteAddrs, uint16 accumulatedWeight, bool blessed) {
    bytes32 taggedRootHash = _taggedRootHash(taggedRoot);
    BlessVoteProgress memory progress = s_blessVoteProgressByTaggedRootHash[taggedRootHash];
    blessed = progress.weightThresholdMet;
    if (progress.configVersion == s_versionedConfig.configVersion) {
      accumulatedWeight = progress.accumulatedWeight;
      uint128 bitmap = progress.voterBitmap;
      blessVoteAddrs = new address[](_bitmapCount(bitmap));
      Voter[] memory voters = s_versionedConfig.config.voters;
      uint256 j = 0;
      for (uint256 i = 0; i < voters.length; ++i) {
        if (_bitmapGet(bitmap, s_blesserRecords[voters[i].blessVoteAddr].index)) {
          blessVoteAddrs[j] = voters[i].blessVoteAddr;
          ++j;
        }
      }
    }
  }

  /// @dev This is a helper method for offchain code so efficiency is not really a concern.
  function getCurseProgress()
    external
    view
    returns (
      address[] memory curseVoteAddrs,
      uint32[] memory voteCounts,
      bytes32[] memory cursesHashes,
      uint16 accumulatedWeight,
      bool cursed
    )
  {
    accumulatedWeight = s_curseVoteProgress.accumulatedWeight;
    cursed = s_curseVoteProgress.curseActive;
    uint256 numCursers;
    Voter[] memory voters = s_versionedConfig.config.voters;
    for (uint256 i = 0; i < voters.length; ++i) {
      CurserRecord memory curserRecord = s_curserRecords[voters[i].curseVoteAddr];
      if (curserRecord.voteCount > 0) {
        ++numCursers;
      }
    }
    curseVoteAddrs = new address[](numCursers);
    voteCounts = new uint32[](numCursers);
    cursesHashes = new bytes32[](numCursers);
    uint256 j = 0;
    for (uint256 i = 0; i < voters.length; ++i) {
      address curseVoteAddr = voters[i].curseVoteAddr;
      CurserRecord memory curserRecord = s_curserRecords[curseVoteAddr];
      if (curserRecord.voteCount > 0) {
        curseVoteAddrs[j] = curseVoteAddr;
        voteCounts[j] = curserRecord.voteCount;
        cursesHashes[j] = curserRecord.cursesHash;
        ++j;
      }
    }
  }

  function _validateConfig(Config memory config) internal pure returns (bool) {
    if (
      config.voters.length == 0 ||
      config.voters.length > MAX_NUM_VOTERS ||
      config.blessWeightThreshold == 0 ||
      config.curseWeightThreshold == 0
    ) {
      return false;
    }

    uint256 totalBlessWeight = 0;
    uint256 totalCurseWeight = 0;
    address[] memory allAddrs = new address[](3 * config.voters.length);
    for (uint256 i = 0; i < config.voters.length; ++i) {
      Voter memory voter = config.voters[i];
      if (
        voter.blessVoteAddr == address(0) ||
        voter.curseVoteAddr == address(0) ||
        voter.curseUnvoteAddr == address(0) ||
        (voter.blessWeight == 0 && voter.curseWeight == 0)
      ) {
        return false;
      }
      allAddrs[3 * i + 0] = voter.blessVoteAddr;
      allAddrs[3 * i + 1] = voter.curseVoteAddr;
      allAddrs[3 * i + 2] = voter.curseUnvoteAddr;
      totalBlessWeight += voter.blessWeight;
      totalCurseWeight += voter.curseWeight;
    }
    for (uint256 i = 0; i < allAddrs.length; ++i) {
      address allAddrs_i = allAddrs[i];
      for (uint256 j = i + 1; j < allAddrs.length; ++j) {
        if (allAddrs_i == allAddrs[j]) {
          return false;
        }
      }
    }

    return totalBlessWeight >= config.blessWeightThreshold && totalCurseWeight >= config.curseWeightThreshold;
  }

  function _setConfig(Config memory config) private {
    if (isCursed()) revert MustRecoverFromCurse();
    if (!_validateConfig(config)) revert InvalidConfig();

    Config memory oldConfig = s_versionedConfig.config;

    // We can't directly assign s_versionedConfig.config to config
    // because copying a memory array into storage is not supported.
    {
      s_versionedConfig.config.blessWeightThreshold = config.blessWeightThreshold;
      s_versionedConfig.config.curseWeightThreshold = config.curseWeightThreshold;
      while (s_versionedConfig.config.voters.length != 0) {
        Voter memory voter = s_versionedConfig.config.voters[s_versionedConfig.config.voters.length - 1];
        delete s_blesserRecords[voter.blessVoteAddr];
        s_curserRecords[voter.curseVoteAddr].active = false;
        s_versionedConfig.config.voters.pop();
      }
      for (uint256 i = 0; i < config.voters.length; ++i) {
        s_versionedConfig.config.voters.push(config.voters[i]);
      }
    }

    ++s_versionedConfig.configVersion;
    uint32 configVersion = s_versionedConfig.configVersion;

    for (uint8 i = 0; i < config.voters.length; ++i) {
      Voter memory voter = config.voters[i];
      s_blesserRecords[voter.blessVoteAddr] = BlesserRecord({
        configVersion: configVersion,
        index: i,
        weight: voter.blessWeight
      });
      s_curserRecords[voter.curseVoteAddr] = CurserRecord({
        active: true,
        weight: voter.curseWeight,
        curseUnvoteAddr: voter.curseUnvoteAddr,
        voteCount: s_curserRecords[voter.curseVoteAddr].voteCount,
        cursesHash: s_curserRecords[voter.curseVoteAddr].cursesHash
      });
    }
    s_versionedConfig.blockNumber = uint32(block.number);
    emit ConfigSet(configVersion, config);

    CurseVoteProgress memory newCurseVoteProgress = CurseVoteProgress({
      curseWeightThreshold: config.curseWeightThreshold,
      accumulatedWeight: 0,
      curseActive: false
    });

    // Retain votes for the cursers who are still part of the new config and delete records for the cursers who are not.
    for (uint8 i = 0; i < oldConfig.voters.length; ++i) {
      // We could be more efficient with this but since this is only for
      // setConfig it will do for now.
      address curseVoteAddr = oldConfig.voters[i].curseVoteAddr;
      CurserRecord memory curserRecord = s_curserRecords[curseVoteAddr];
      if (!curserRecord.active) {
        delete s_curserRecords[curseVoteAddr];
      } else if (curserRecord.active && curserRecord.voteCount > 0) {
        newCurseVoteProgress.accumulatedWeight += curserRecord.weight;
        emit ReusedVotesToCurse(
          configVersion,
          curseVoteAddr,
          curserRecord.weight,
          curserRecord.voteCount,
          curserRecord.cursesHash,
          newCurseVoteProgress.accumulatedWeight
        );
      }
    }
    newCurseVoteProgress.curseActive =
      newCurseVoteProgress.accumulatedWeight >= newCurseVoteProgress.curseWeightThreshold;
    if (newCurseVoteProgress.curseActive) {
      emit Cursed(configVersion, block.timestamp);
    }
    s_curseVoteProgress = newCurseVoteProgress;
  }
}
