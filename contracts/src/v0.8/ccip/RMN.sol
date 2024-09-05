// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IRMN} from "./interfaces/IRMN.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";

import {EnumerableSet} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

// An active curse on this subject will cause isCursed() to return true. Use this subject if there is an issue with a
// remote chain, for which there exists a legacy lane contract deployed on the same chain as this RMN contract is
// deployed, relying on isCursed().
bytes16 constant LEGACY_CURSE_SUBJECT = 0x01000000000000000000000000000000;

// An active curse on this subject will cause isCursed() and isCursed(bytes16) to return true. Use this subject for
// issues affecting all of CCIP chains, or pertaining to the chain that this contract is deployed on, instead of using
// the local chain selector as a subject.
bytes16 constant GLOBAL_CURSE_SUBJECT = 0x01000000000000000000000000000001;

// The curse vote address representing the owner in data structures, events and recorded votes. Remains constant, even
// if the owner changes.
address constant OWNER_CURSE_VOTE_ADDR = address(~uint160(0)); // 0xff...ff

// The curse vote address used in an OwnerUnvoteToCurseRequest to lift a curse, if there is no active curse votes for
// the subject that we are able to unvote, but the conditions for an active curse no longer hold.
address constant LIFT_CURSE_VOTE_ADDR = address(0);

/// @dev This contract is owned by RMN, if changing, please notify the RMN maintainers.
// solhint-disable chainlink-solidity/explicit-returns
contract RMN is IRMN, OwnerIsCreator, ITypeAndVersion {
  using EnumerableSet for EnumerableSet.AddressSet;

  // STATIC CONFIG
  string public constant override typeAndVersion = "RMN 1.5.0";

  uint256 private constant MAX_NUM_VOTERS = 16;

  // MAGIC VALUES
  bytes28 private constant NO_VOTES_CURSES_HASH = bytes28(0);

  // DYNAMIC CONFIG
  /// @notice blessVoteAddr and curseVoteAddr can't be 0. Additionally curseVoteAddr can't be LIFT_CURSE_VOTE_ADDR or
  /// OWNER_CURSE_VOTE_ADDR. At least one of blessWeight & curseWeight must be non-zero, i.e., a voter could only vote
  /// to bless, or only vote to curse, or both vote to bless and vote to curse.
  struct Voter {
    // This is the address the voter should use to call voteToBless.
    address blessVoteAddr;
    // This is the address the voter should use to call voteToCurse.
    address curseVoteAddr;
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
    // When the total weight of voters that have voted to curse a subject reaches or
    // exceeds curseWeightThreshold, the subject becomes cursed.
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
    // This particular ordering saves us ~400 gas per voteToBless call, compared to the bool being at the bottom, even
    // though the size of the struct is the same.
    bool weightThresholdMet;
    // A BlessVoteProgress is considered invalid if weightThresholdMet is false when
    // s_versionedConfig.configVersion changes. we don't want old in-progress
    // votes to continue when we set a new config!
    // The config version at which the bless vote for a tagged root was initiated.
    uint32 configVersion;
    uint16 accumulatedWeight;
    // Care must be taken that the bitmap has at least as many bits as MAX_NUM_VOTERS.
    // uint200 is much larger than we need, but it saves us ~100 gas per voteToBless call to fill the word instead of
    // using a smaller type.
    // _bitmapGet(voterBitmap, i) = true indicates that the i-th voter has voted to bless
    uint200 voterBitmap;
  }

  mapping(bytes32 taggedRootHash => BlessVoteProgress blessVoteProgress) private s_blessVoteProgressByTaggedRootHash;

  // Any tagged root with a commit store included in s_permaBlessedCommitStores will be considered automatically
  // blessed.
  EnumerableSet.AddressSet private s_permaBlessedCommitStores;

  struct CurserRecord {
    bool active;
    uint8 weight;
    mapping(bytes16 curseId => bool used) usedCurseIds; // retained across config changes
  }

  mapping(address curseVoteAddr => CurserRecord curserRecord) private s_curserRecords;

  struct ConfigVersionAndCursesHash {
    uint32 configVersion; // configVersion != s_versionedConfig.configVersion means no active vote
    bytes28 cursesHash; // bytes28(0) means no active vote; truncated so that ConfigVersionAndCursesHash fits in a word
  }

  struct CurseVoteProgress {
    uint32 configVersion; // upon config change, lazy set to new config version
    uint16 curseWeightThreshold; // upon config change, lazy set to new config value
    uint16 accumulatedWeight; // upon config change, lazy set to 0
    // A curse becomes active after either:
    // - sum([voter.weight for voter who voted in current config]) >= curseWeightThreshold
    // - ownerCurse is invoked
    // Once a curse is active, only the owner can lift it.
    bool curseActive; // retained across config changes
    mapping(address => ConfigVersionAndCursesHash) latestVoteToCurseByCurseVoteAddr; // retained across config changes
  }

  mapping(bytes16 subject => CurseVoteProgress curseVoteProgress) private
    s_potentiallyOutdatedCurseVoteProgressBySubject;

  // We intentionally use a struct here, even though it contains a single field, to make it obvious to future editors
  // that there is space for more fields.
  struct CurseHotVars {
    uint64 numSubjectsCursed; // incremented by voteToCurse, ownerCurse; decremented by ownerUnvoteToCurse
  }

  CurseHotVars private s_curseHotVars;

  enum RecordedCurseRelatedOpTag {
    // A vote to curse, through either voteToCurse or ownerCurse.
    VoteToCurse,
    // An unvote to curse, through unvoteToCurse.
    UnvoteToCurse,
    // An unvote to curse, through ownerUnvoteToCurse, which was not forced (forceUnvote=false).
    OwnerUnvoteToCurseUnforced,
    // An unvote to curse, through ownerUnvoteToCurse, which was forced (forceUnvote=true).
    OwnerUnvoteToCurseForced,
    // A configuration change.
    //
    // For subjects that are not cursed when this happens, past votes do not get accounted for in the new configuration.
    // If a voter votes during the new configuration, their curses hash will restart from NO_VOTES_CURSES_HASH.
    //
    // For subjects that are cursed when this happens, past votes get accounted for.
    // If a voter votes during the new configuration, their curses hash will continue from its old value.
    SetConfig
  }

  /// @notice Provides the ability to quickly reconstruct the curse-related state of the contract offchain, without
  /// having to replay all past events. Replaying past events often takes long, and in some cases might even be
  /// infeasible due to log pruning.
  ///
  /// @dev We could save some gas by omitting some fields and instead using them as mapping keys, but we would lose the
  /// cross-voter ordering, or cross-subject ordering, or cross-vote/unvote ordering.
  struct RecordedCurseRelatedOp {
    RecordedCurseRelatedOpTag tag;
    uint64 blockTimestamp;
    bool cursed; // whether the subject is cursed after this op; if tag in {SetConfig}, will be false
    address curseVoteAddr; // if tag in {SetConfig}, will be address(0)
    bytes16 subject; // if tag in {SetConfig}, will be bytes16(0)
    bytes16 curseId; // if tag in {SetConfig, UnvoteToCurse, OwnerUnvoteToCurseUnforced, OwnerUnvoteToCurseForced}, will be bytes16(0)
  }

  RecordedCurseRelatedOp[] private s_recordedCurseRelatedOps;

  /// @dev This function is to _ONLY_ be called in order to determine if a curse should become active upon a
  /// vote-to-curse, or a curse should be deactivated upon an owner-unvote-to-curse.
  /// Other reasons for a curse to be active, which are not covered here:
  /// 1. Cursedness is retained from a prior config.
  /// 2. The curse weight threshold was met at some point, which activated a curse, and enough voters unvoted to curse
  /// such that the curse weight threshold is no longer met.
  function _shouldCurseBeActive(CurseVoteProgress storage sptr_upToDateCurseVoteProgress) internal view returns (bool) {
    return sptr_upToDateCurseVoteProgress.latestVoteToCurseByCurseVoteAddr[OWNER_CURSE_VOTE_ADDR].cursesHash
      != NO_VOTES_CURSES_HASH
      || sptr_upToDateCurseVoteProgress.accumulatedWeight >= sptr_upToDateCurseVoteProgress.curseWeightThreshold;
  }

  /// @dev It might be the case that due to the lazy update of curseVoteProgress, a curse is active even though
  /// _shouldCurseBeActive(curseVoteProgress) is false, i.e., the owner has no active vote to curse and the curse
  /// weight threshold has not been met.
  function _getUpToDateCurseVoteProgress(
    uint32 configVersion,
    bytes16 subject
  ) internal returns (CurseVoteProgress storage) {
    CurseVoteProgress storage sptr_curseVoteProgress = s_potentiallyOutdatedCurseVoteProgressBySubject[subject];
    if (configVersion != sptr_curseVoteProgress.configVersion) {
      sptr_curseVoteProgress.configVersion = configVersion;
      sptr_curseVoteProgress.curseWeightThreshold = s_versionedConfig.config.curseWeightThreshold;
      sptr_curseVoteProgress.accumulatedWeight = 0;

      if (sptr_curseVoteProgress.curseActive) {
        // If a curse was active, count past votes to curse and retain the curses hash for cursers who are part of the
        // new config.
        Config storage sptr_config = s_versionedConfig.config;
        for (uint256 i = 0; i < sptr_config.voters.length; ++i) {
          Voter storage sptr_voter = sptr_config.voters[i];
          ConfigVersionAndCursesHash storage sptr_cvch =
            sptr_curseVoteProgress.latestVoteToCurseByCurseVoteAddr[sptr_voter.curseVoteAddr];
          if (sptr_cvch.configVersion < configVersion && sptr_cvch.cursesHash != NO_VOTES_CURSES_HASH) {
            // `< configVersion` instead of `== configVersion-1`, because there might have been multiple config changes
            // without a lazy update of our subject. This has the side effect of retaining votes from very old configs
            // that we might not really intend to retain, but these can be removed by the owner later.
            sptr_cvch.configVersion = configVersion;
            sptr_curseVoteProgress.accumulatedWeight += sptr_voter.curseWeight;
          }
        }
        // We don't need to think about OWNER_CURSE_VOTE_ADDR here, because its ConfigVersionAndCursesHash counts even
        // if the configVersion is not the current config version, in contrast to regular voters.
        // It's an irregularity, but it saves us > 5k gas (if the owner had previously voted) for the unlucky voter who
        // enters this branch.
      } else {
        // If a curse was not active, we don't count past votes to curse for voters who are part of the new config.
        // Their curses hash will be restart from NO_VOTES_CURSES_HASH when they vote to curse again.
        // We expect that the offchain code will revote to curse in case it voted to curse, and the vote to curse was
        // lost due to any reason, including a config change when the curse was not yet active.
      }
    }
    return sptr_curseVoteProgress;
  }

  // EVENTS, ERRORS

  event ConfigSet(uint32 indexed configVersion, Config config);

  error InvalidConfig();

  event TaggedRootBlessed(uint32 indexed configVersion, IRMN.TaggedRoot taggedRoot, uint16 accumulatedWeight);
  event TaggedRootBlessVotesReset(uint32 indexed configVersion, IRMN.TaggedRoot taggedRoot, bool wasBlessed);
  event VotedToBless(uint32 indexed configVersion, address indexed voter, IRMN.TaggedRoot taggedRoot, uint8 weight);

  event VotedToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    bytes16 subject,
    bytes16 curseId,
    uint8 weight,
    uint64 blockTimestamp,
    bytes28 cursesHash,
    uint16 accumulatedWeight
  );
  event UnvotedToCurse(
    uint32 indexed configVersion,
    address indexed voter,
    bytes16 subject,
    uint8 weight,
    bytes28 cursesHash,
    uint16 remainingAccumulatedWeight
  );
  event SkippedUnvoteToCurse(address indexed voter, bytes16 subject, bytes28 onchainCursesHash, bytes28 cursesHash);
  event Cursed(uint32 indexed configVersion, bytes16 subject, uint64 blockTimestamp);
  event CurseLifted(bytes16 subject);

  // These events make it easier for offchain logic to discover that it performs
  // the same actions multiple times.
  event AlreadyVotedToBless(uint32 indexed configVersion, address indexed voter, IRMN.TaggedRoot taggedRoot);
  event AlreadyBlessed(uint32 indexed configVersion, address indexed voter, IRMN.TaggedRoot taggedRoot);

  // Emitted by ownerRemoveThenAddPermaBlessedCommitStores.
  event PermaBlessedCommitStoreAdded(address commitStore);
  event PermaBlessedCommitStoreRemoved(address commitStore);

  error ReusedCurseId(address voter, bytes16 curseId);
  error UnauthorizedVoter(address voter);
  error VoteToBlessNoop();
  error VoteToCurseNoop();
  error UnvoteToCurseNoop();
  error VoteToBlessForbiddenDuringActiveGlobalCurse();

  /// @notice Thrown when subjects are not a strictly increasing monotone sequence.
  // Prevents a subject from receiving multiple votes to curse with the same curse id.
  error SubjectsMustBeStrictlyIncreasing();

  constructor(Config memory config) {
    {
      // Ensure that the bitmap is large enough to hold MAX_NUM_VOTERS.
      // We do this in the constructor because MAX_NUM_VOTERS is constant.
      BlessVoteProgress memory vp = BlessVoteProgress({
        configVersion: 0,
        voterBitmap: type(uint200).max, // will not compile if it doesn't fit
        accumulatedWeight: 0,
        weightThresholdMet: false
      });
      assert(vp.voterBitmap >> (MAX_NUM_VOTERS - 1) >= 1);
    }
    _setConfig(config);
  }

  function _bitmapGet(uint200 bitmap, uint8 index) internal pure returns (bool) {
    assert(index < MAX_NUM_VOTERS);
    return bitmap & (uint200(1) << index) != 0;
  }

  function _bitmapSet(uint200 bitmap, uint8 index) internal pure returns (uint200) {
    assert(index < MAX_NUM_VOTERS);
    return bitmap | (uint200(1) << index);
  }

  function _bitmapCount(uint200 bitmap) internal pure returns (uint8 oneBits) {
    assert(bitmap < 1 << MAX_NUM_VOTERS);
    // https://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
    for (; bitmap != 0; ++oneBits) {
      bitmap &= bitmap - 1;
    }
  }

  function _taggedRootHash(IRMN.TaggedRoot memory taggedRoot) internal pure returns (bytes32) {
    return keccak256(abi.encode(taggedRoot.commitStore, taggedRoot.root));
  }

  function _cursesHash(bytes28 prevCursesHash, bytes16 curseId) internal pure returns (bytes28) {
    return bytes28(keccak256(abi.encode(prevCursesHash, curseId)));
  }

  function _blockTimestamp() internal view returns (uint64) {
    return uint64(block.timestamp);
  }

  /// @param taggedRoots A tagged root is hashed as `keccak256(abi.encode(taggedRoot.commitStore
  /// /* address */, taggedRoot.root /* bytes32 */))`.
  /// @notice Tagged roots which are already (voted to be) blessed are skipped and emit corresponding events. In case
  /// the call has no effect, i.e., all passed tagged roots are skipped, the function reverts with a `VoteToBlessNoop`.
  function voteToBless(IRMN.TaggedRoot[] calldata taggedRoots) external {
    // If we have an active global curse, something is really wrong. Let's err on the
    // side of caution and not accept further blessings during this time of
    // uncertainty.
    if (isCursed(GLOBAL_CURSE_SUBJECT)) revert VoteToBlessForbiddenDuringActiveGlobalCurse();

    uint32 configVersion = s_versionedConfig.configVersion;
    BlesserRecord memory blesserRecord = s_blesserRecords[msg.sender];
    if (blesserRecord.configVersion != configVersion) revert UnauthorizedVoter(msg.sender);

    bool noop = true;
    for (uint256 i = 0; i < taggedRoots.length; ++i) {
      IRMN.TaggedRoot memory taggedRoot = taggedRoots[i];
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
      } else if (voteProgress.configVersion != configVersion) {
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
      } else if (_bitmapGet(voteProgress.voterBitmap, blesserRecord.index)) {
        // We don't revert here because there might be other tagged roots for
        // which votes might count, and we want to allow that to happen.
        emit AlreadyVotedToBless(configVersion, msg.sender, taggedRoot);
        continue;
      }
      noop = false;
      voteProgress.voterBitmap = _bitmapSet(voteProgress.voterBitmap, blesserRecord.index);
      voteProgress.accumulatedWeight += blesserRecord.weight;
      emit VotedToBless(configVersion, msg.sender, taggedRoot, blesserRecord.weight);
      if (voteProgress.accumulatedWeight >= s_versionedConfig.config.blessWeightThreshold) {
        voteProgress.weightThresholdMet = true;
        emit TaggedRootBlessed(configVersion, taggedRoot, voteProgress.accumulatedWeight);
      }
      s_blessVoteProgressByTaggedRootHash[taggedRootHash] = voteProgress;
    }

    if (noop) {
      revert VoteToBlessNoop();
    }
  }

  /// @notice Can be called by the owner to remove unintentionally voted or even blessed tagged roots in a recovery
  /// scenario. The owner must ensure that there are no in-flight transactions by RMN nodes voting for any of the
  /// taggedRoots before calling this function, as such in-flight transactions could lead to the roots becoming
  /// re-blessed shortly after the call to this function, contrary to the original intention.
  function ownerResetBlessVotes(IRMN.TaggedRoot[] calldata taggedRoots) external onlyOwner {
    uint32 configVersion = s_versionedConfig.configVersion;
    for (uint256 i = 0; i < taggedRoots.length; ++i) {
      IRMN.TaggedRoot memory taggedRoot = taggedRoots[i];
      bytes32 taggedRootHash = _taggedRootHash(taggedRoot);
      BlessVoteProgress memory voteProgress = s_blessVoteProgressByTaggedRootHash[taggedRootHash];
      delete s_blessVoteProgressByTaggedRootHash[taggedRootHash];
      bool wasBlessed = voteProgress.weightThresholdMet;
      if (voteProgress.configVersion == configVersion || wasBlessed) {
        emit TaggedRootBlessVotesReset(configVersion, taggedRoot, wasBlessed);
      }
    }
  }

  struct UnvoteToCurseRequest {
    bytes16 subject;
    bytes28 cursesHash;
  }

  // For use in internal calls.
  enum Privilege {
    Owner,
    Voter
  }

  function _authorizedUnvoteToCurse(
    Privilege priv, // Privilege.Owner during an ownerUnvoteToCurse call, Privilege.Voter during a unvoteToCurse call
    uint32 configVersion,
    address curseVoteAddr,
    UnvoteToCurseRequest memory req,
    bool forceUnvote, // true only during an ownerUnvoteToCurse call, when OwnerUnvoteToCurseRequest.forceUnvote is true
    CurserRecord storage sptr_curserRecord,
    CurseVoteProgress storage sptr_curseVoteProgress
  ) internal returns (bool unvoted, bool curseLifted) {
    {
      assert(priv == Privilege.Voter || priv == Privilege.Owner); // sanity check
      // Check that the supplied arguments are feasible for our privilege.
      if (forceUnvote || curseVoteAddr == OWNER_CURSE_VOTE_ADDR || curseVoteAddr == LIFT_CURSE_VOTE_ADDR) {
        assert(priv == Privilege.Owner);
      }
    }

    ConfigVersionAndCursesHash memory cvch = sptr_curseVoteProgress.latestVoteToCurseByCurseVoteAddr[curseVoteAddr];

    // First, try to unvote.
    if (
      sptr_curserRecord.active && (curseVoteAddr == OWNER_CURSE_VOTE_ADDR || cvch.configVersion == configVersion)
        && cvch.cursesHash != NO_VOTES_CURSES_HASH && (cvch.cursesHash == req.cursesHash || forceUnvote)
    ) {
      unvoted = true;
      delete sptr_curseVoteProgress.latestVoteToCurseByCurseVoteAddr[curseVoteAddr];
      // Assumes: s_curserRecords[OWNER_CURSE_VOTE_ADDR].weight == 0, enforced by _setConfig
      sptr_curseVoteProgress.accumulatedWeight -= sptr_curserRecord.weight;

      emit UnvotedToCurse(
        configVersion,
        curseVoteAddr,
        req.subject,
        sptr_curserRecord.weight,
        req.cursesHash,
        sptr_curseVoteProgress.accumulatedWeight
      );
    }

    // If we have owner privilege, and the conditions for the curse to be active no longer hold, we are able to lift the
    // curse.
    bool shouldTryToLiftCurse = priv == Privilege.Owner && (unvoted || curseVoteAddr == LIFT_CURSE_VOTE_ADDR);

    if (shouldTryToLiftCurse && sptr_curseVoteProgress.curseActive && !_shouldCurseBeActive(sptr_curseVoteProgress)) {
      curseLifted = true;
      sptr_curseVoteProgress.curseActive = false;
      --s_curseHotVars.numSubjectsCursed;
      emit CurseLifted(req.subject);
    }

    if (unvoted || curseLifted) {
      RecordedCurseRelatedOpTag tag;
      if (priv == Privilege.Owner) {
        if (forceUnvote) {
          tag = RecordedCurseRelatedOpTag.OwnerUnvoteToCurseForced;
        } else {
          tag = RecordedCurseRelatedOpTag.OwnerUnvoteToCurseUnforced;
        }
      } else if (priv == Privilege.Voter) {
        tag = RecordedCurseRelatedOpTag.UnvoteToCurse;
      } else {
        // solhint-disable-next-line gas-custom-errors, reason-string
        revert(); // assumption violation
      }
      s_recordedCurseRelatedOps.push(
        RecordedCurseRelatedOp({
          tag: tag,
          cursed: sptr_curseVoteProgress.curseActive,
          curseVoteAddr: curseVoteAddr,
          curseId: bytes16(0),
          subject: req.subject,
          blockTimestamp: _blockTimestamp()
        })
      );
    } else {
      emit SkippedUnvoteToCurse(curseVoteAddr, req.subject, cvch.cursesHash, req.cursesHash);
    }
  }

  /// @notice Can be called by a curser to remove unintentional votes to curse.
  /// We expect this to be called very rarely, e.g. in case of a bug in the
  /// offchain code causing false voteToCurse calls.
  /// @notice Should be called from curser's corresponding curseVoteAddr.
  function unvoteToCurse(UnvoteToCurseRequest[] memory unvoteToCurseRequests) external {
    address curseVoteAddr = msg.sender;
    CurserRecord storage sptr_curserRecord = s_curserRecords[curseVoteAddr];

    if (!sptr_curserRecord.active) revert UnauthorizedVoter(curseVoteAddr);

    uint32 configVersion = s_versionedConfig.configVersion;
    bool anyVoteWasUnvoted = false;
    for (uint256 i = 0; i < unvoteToCurseRequests.length; ++i) {
      UnvoteToCurseRequest memory req = unvoteToCurseRequests[i];
      CurseVoteProgress storage sptr_curseVoteProgress = _getUpToDateCurseVoteProgress(configVersion, req.subject);
      (bool unvoted, bool curseLifted) = _authorizedUnvoteToCurse(
        Privilege.Voter, configVersion, curseVoteAddr, req, false, sptr_curserRecord, sptr_curseVoteProgress
      );
      assert(!curseLifted); // assumption violation: voters can't lift curses
      anyVoteWasUnvoted = anyVoteWasUnvoted || unvoted;
    }

    if (!anyVoteWasUnvoted) {
      revert UnvoteToCurseNoop();
    }
  }

  /// @notice A vote to curse is appropriate during unhealthy blockchain conditions
  /// (eg. finality violations).
  function voteToCurse(bytes16 curseId, bytes16[] memory subjects) external {
    address curseVoteAddr = msg.sender;
    assert(curseVoteAddr != OWNER_CURSE_VOTE_ADDR);
    CurserRecord storage sptr_curserRecord = s_curserRecords[curseVoteAddr];
    if (!sptr_curserRecord.active) revert UnauthorizedVoter(curseVoteAddr);
    _authorizedVoteToCurse(curseVoteAddr, curseId, subjects, sptr_curserRecord);
  }

  function _authorizedVoteToCurse(
    address curseVoteAddr,
    bytes16 curseId,
    bytes16[] memory subjects,
    CurserRecord storage sptr_curserRecord
  ) internal {
    if (subjects.length == 0) revert VoteToCurseNoop();

    if (sptr_curserRecord.usedCurseIds[curseId]) revert ReusedCurseId(curseVoteAddr, curseId);
    sptr_curserRecord.usedCurseIds[curseId] = true;

    // NOTE: We could pack configVersion into CurserRecord that we already load in the beginning of this function to
    // avoid the following extra storage read for it, but since voteToCurse is not on the hot path we'd rather keep
    // things simple.
    uint32 configVersion = s_versionedConfig.configVersion;
    for (uint256 i = 0; i < subjects.length; ++i) {
      if (i >= 1 && !(subjects[i - 1] < subjects[i])) {
        // Prevents a subject from receiving multiple votes to curse with the same curse id.
        revert SubjectsMustBeStrictlyIncreasing();
      }

      bytes16 subject = subjects[i];
      CurseVoteProgress storage sptr_curseVoteProgress = _getUpToDateCurseVoteProgress(configVersion, subject);
      ConfigVersionAndCursesHash memory cvch = sptr_curseVoteProgress.latestVoteToCurseByCurseVoteAddr[curseVoteAddr];
      bytes28 prevCursesHash;
      if (
        (curseVoteAddr != OWNER_CURSE_VOTE_ADDR && cvch.configVersion < configVersion)
          || cvch.cursesHash == NO_VOTES_CURSES_HASH
      ) {
        // if owner's first vote, or if voter's first vote in this config version
        prevCursesHash = NO_VOTES_CURSES_HASH; // start hashchain from scratch, explicit
        sptr_curseVoteProgress.accumulatedWeight += sptr_curserRecord.weight;
      } else {
        // we've already accounted for the weight
        prevCursesHash = cvch.cursesHash;
      }
      sptr_curseVoteProgress.latestVoteToCurseByCurseVoteAddr[curseVoteAddr] = cvch =
        ConfigVersionAndCursesHash({configVersion: configVersion, cursesHash: _cursesHash(prevCursesHash, curseId)});
      emit VotedToCurse(
        configVersion,
        curseVoteAddr,
        subject,
        curseId,
        sptr_curserRecord.weight,
        _blockTimestamp(),
        cvch.cursesHash,
        sptr_curseVoteProgress.accumulatedWeight
      );

      if (
        prevCursesHash == NO_VOTES_CURSES_HASH && !sptr_curseVoteProgress.curseActive
          && _shouldCurseBeActive(sptr_curseVoteProgress)
      ) {
        sptr_curseVoteProgress.curseActive = true;
        ++s_curseHotVars.numSubjectsCursed;
        emit Cursed(configVersion, subject, _blockTimestamp());
      }

      s_recordedCurseRelatedOps.push(
        RecordedCurseRelatedOp({
          tag: RecordedCurseRelatedOpTag.VoteToCurse,
          cursed: sptr_curseVoteProgress.curseActive,
          curseVoteAddr: curseVoteAddr,
          curseId: curseId,
          subject: subject,
          blockTimestamp: _blockTimestamp()
        })
      );
    }
  }

  /// @notice Enables the owner to immediately have the system enter the cursed state.
  function ownerCurse(bytes16 curseId, bytes16[] memory subjects) external onlyOwner {
    address curseVoteAddr = OWNER_CURSE_VOTE_ADDR;
    CurserRecord storage sptr_curserRecord = s_curserRecords[curseVoteAddr];
    // no need to check if sptr_curserRecord.active, we must have the onlyOwner modifier
    _authorizedVoteToCurse(curseVoteAddr, curseId, subjects, sptr_curserRecord);
  }

  // Set curseVoteAddr=LIFT_CURSE_VOTE_ADDR, cursesHash=bytes28(0), to reset curseActive if it can be reset. Useful if
  // all voters have unvoted to curse on their own and the curse can now be lifted without any individual votes that can
  // be unvoted.
  // solhint-disable-next-line gas-struct-packing
  struct OwnerUnvoteToCurseRequest {
    address curseVoteAddr;
    UnvoteToCurseRequest unit;
    bool forceUnvote;
  }

  /// @notice Enables the owner to remove curse votes. After the curse votes are removed,
  /// this function will check whether the curse is still valid and restore the uncursed state if possible.
  /// This function also enables the owner to lift a curse created through ownerCurse.
  function ownerUnvoteToCurse(OwnerUnvoteToCurseRequest[] memory ownerUnvoteToCurseRequests) external onlyOwner {
    bool anyCurseWasLifted = false;
    bool anyVoteWasUnvoted = false;
    uint32 configVersion = s_versionedConfig.configVersion;
    for (uint256 i = 0; i < ownerUnvoteToCurseRequests.length; ++i) {
      OwnerUnvoteToCurseRequest memory req = ownerUnvoteToCurseRequests[i];
      CurseVoteProgress storage sptr_curseVoteProgress = _getUpToDateCurseVoteProgress(configVersion, req.unit.subject);
      (bool unvoted, bool curseLifted) = _authorizedUnvoteToCurse(
        Privilege.Owner,
        configVersion,
        req.curseVoteAddr,
        req.unit,
        req.forceUnvote,
        s_curserRecords[req.curseVoteAddr],
        sptr_curseVoteProgress
      );
      anyVoteWasUnvoted = anyVoteWasUnvoted || unvoted;
      anyCurseWasLifted = anyCurseWasLifted || curseLifted;
    }

    if (anyCurseWasLifted) {
      // Invalidate all in-progress votes to bless or curse by bumping the config version.
      // They might have been based on false information about the source chain
      // (e.g. in case of a finality violation).
      _setConfig(s_versionedConfig.config);
    }

    if (!(anyVoteWasUnvoted || anyCurseWasLifted)) {
      revert UnvoteToCurseNoop();
    }
  }

  function setConfig(Config memory config) external onlyOwner {
    _setConfig(config);
  }

  /// @notice Any tagged root with a commit store included in this array will be considered automatically blessed.
  function getPermaBlessedCommitStores() external view returns (address[] memory) {
    return s_permaBlessedCommitStores.values();
  }

  /// @notice The ordering of parameters is important. First come the commit stores to remove, then the commit stores to
  /// add.
  function ownerRemoveThenAddPermaBlessedCommitStores(
    address[] memory removes,
    address[] memory adds
  ) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      if (s_permaBlessedCommitStores.remove(removes[i])) {
        emit PermaBlessedCommitStoreRemoved(removes[i]);
      }
    }
    for (uint256 i = 0; i < adds.length; ++i) {
      if (s_permaBlessedCommitStores.add(adds[i])) {
        emit PermaBlessedCommitStoreAdded(adds[i]);
      }
    }
  }

  /// @inheritdoc IRMN
  function isBlessed(IRMN.TaggedRoot calldata taggedRoot) external view returns (bool) {
    return s_blessVoteProgressByTaggedRootHash[_taggedRootHash(taggedRoot)].weightThresholdMet
      || s_permaBlessedCommitStores.contains(taggedRoot.commitStore);
  }

  /// @inheritdoc IRMN
  function isCursed() external view returns (bool) {
    if (s_curseHotVars.numSubjectsCursed == 0) {
      return false; // happy path costs a single SLOAD
    } else {
      return s_potentiallyOutdatedCurseVoteProgressBySubject[GLOBAL_CURSE_SUBJECT].curseActive
        || s_potentiallyOutdatedCurseVoteProgressBySubject[LEGACY_CURSE_SUBJECT].curseActive;
    }
  }

  /// @inheritdoc IRMN
  function isCursed(bytes16 subject) public view returns (bool) {
    if (s_curseHotVars.numSubjectsCursed == 0) {
      return false; // happy path costs a single SLOAD
    } else {
      return s_potentiallyOutdatedCurseVoteProgressBySubject[GLOBAL_CURSE_SUBJECT].curseActive
        || s_potentiallyOutdatedCurseVoteProgressBySubject[subject].curseActive;
    }
  }

  /// @notice Config version might be incremented for many reasons, including
  /// the lifting of a curse, or a regular config change.
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
    IRMN.TaggedRoot calldata taggedRoot
  ) external view returns (address[] memory blessVoteAddrs, uint16 accumulatedWeight, bool blessed) {
    bytes32 taggedRootHash = _taggedRootHash(taggedRoot);
    BlessVoteProgress memory progress = s_blessVoteProgressByTaggedRootHash[taggedRootHash];
    blessed = progress.weightThresholdMet;
    if (progress.configVersion == s_versionedConfig.configVersion) {
      accumulatedWeight = progress.accumulatedWeight;
      uint200 bitmap = progress.voterBitmap;
      blessVoteAddrs = new address[](_bitmapCount(bitmap));
      Voter[] memory voters = s_versionedConfig.config.voters;
      uint256 j = 0;
      for (uint8 i = 0; i < voters.length; ++i) {
        if (_bitmapGet(bitmap, i)) {
          blessVoteAddrs[j] = voters[i].blessVoteAddr;
          ++j;
        }
      }
    }
  }

  /// @return curseVoteAddrs the curseVoteAddr of each voter with an active vote to curse
  /// @return cursesHashes the i-th value is the curses hash of curseVoteAddrs[i]
  /// @return accumulatedWeight the accumulated weight of all voters with an active vote to curse who are part of the
  /// current config
  /// @return cursed might be true even if the owner has no active vote and accumulatedWeight < curseWeightThreshold,
  /// due to a retained curse from a prior config
  /// @dev This is a helper method for offchain code so efficiency is not really a concern.
  function getCurseProgress(
    bytes16 subject
  )
    external
    view
    returns (address[] memory curseVoteAddrs, bytes28[] memory cursesHashes, uint16 accumulatedWeight, bool cursed)
  {
    uint32 configVersion = s_versionedConfig.configVersion;
    Config memory config = s_versionedConfig.config;
    // Can't use _getUpToDateCurseVoteProgress here because we can't call a non-view function from within a view.
    // So we get to repeat some accounting.
    CurseVoteProgress storage outdatedCurseVoteProgress = s_potentiallyOutdatedCurseVoteProgressBySubject[subject];

    cursed = outdatedCurseVoteProgress.curseActive;

    // See _getUpToDateCurseVoteProgress for more context.
    bool shouldCountVotesFromOlderConfigs = outdatedCurseVoteProgress.configVersion < configVersion && cursed;

    // A play in two acts, because we can't push to arrays in memory, so we need to precompute the array's length.
    // First act: we count the number of cursers, i.e., voters with active vote.
    // Second act: push the cursers to the arrays, sum their weights.

    uint256 numCursers = 0; // we reuse this variable for writing to perserve stack space
    accumulatedWeight = 0;
    for (uint256 act = 1; act <= 2; ++act) {
      uint256 i = config.voters.length; // not config.voters.length-1 to account for the owner
      while (true) {
        address curseVoteAddr;
        uint8 weight;
        if (i < config.voters.length) {
          curseVoteAddr = config.voters[i].curseVoteAddr;
          weight = config.voters[i].curseWeight;
        } else {
          // Allows us to include the owner's vote and curses hash in the result.
          curseVoteAddr = OWNER_CURSE_VOTE_ADDR;
          weight = 0;
        }

        ConfigVersionAndCursesHash memory cvch =
          outdatedCurseVoteProgress.latestVoteToCurseByCurseVoteAddr[curseVoteAddr];
        bool hasActiveVote = (
          shouldCountVotesFromOlderConfigs || cvch.configVersion == configVersion
            || curseVoteAddr == OWNER_CURSE_VOTE_ADDR
        ) && cvch.cursesHash != NO_VOTES_CURSES_HASH;
        if (hasActiveVote) {
          if (act == 1) {
            ++numCursers;
          } else if (act == 2) {
            accumulatedWeight += weight;
            --numCursers;
            curseVoteAddrs[numCursers] = curseVoteAddr;
            cursesHashes[numCursers] = cvch.cursesHash;
          } else {
            // solhint-disable-next-line gas-custom-errors, reason-string
            revert(); // assumption violation
          }
        }

        if (i > 0) {
          --i;
        } else {
          break;
        }
      }

      if (act == 1) {
        // We are done counting at this point, initialize the arrays for the second act that follows immediately after.
        curseVoteAddrs = new address[](numCursers);
        cursesHashes = new bytes28[](numCursers);
      }
    }
  }

  /// @notice Returns the number of subjects that are currently cursed.
  function getCursedSubjectsCount() external view returns (uint256) {
    return s_curseHotVars.numSubjectsCursed;
  }

  /// @dev This is a helper method for offchain code to know what arguments to use for getRecordedCurseRelatedOps.
  function getRecordedCurseRelatedOpsCount() external view returns (uint256) {
    return s_recordedCurseRelatedOps.length;
  }

  /// @dev This is a helper method for offchain code so efficiency is not really a concern.
  /// @dev Returns s_recordedCurseRelatedOps[offset:offset+limit].
  function getRecordedCurseRelatedOps(
    uint256 offset,
    uint256 limit
  ) external view returns (RecordedCurseRelatedOp[] memory) {
    uint256 pageLen;
    if (offset + limit <= s_recordedCurseRelatedOps.length) {
      pageLen = limit;
    } else if (offset < s_recordedCurseRelatedOps.length) {
      pageLen = s_recordedCurseRelatedOps.length - offset;
    } else {
      pageLen = 0;
    }
    RecordedCurseRelatedOp[] memory page = new RecordedCurseRelatedOp[](pageLen);
    for (uint256 i = 0; i < pageLen; ++i) {
      page[i] = s_recordedCurseRelatedOps[offset + i];
    }
    return page;
  }

  function _validateConfig(Config memory config) internal pure returns (bool) {
    if (
      config.voters.length == 0 || config.voters.length > MAX_NUM_VOTERS || config.blessWeightThreshold == 0
        || config.curseWeightThreshold == 0
    ) {
      return false;
    }

    uint256 totalBlessWeight = 0;
    uint256 totalCurseWeight = 0;
    address[] memory allAddrs = new address[](2 * config.voters.length);
    for (uint256 i = 0; i < config.voters.length; ++i) {
      Voter memory voter = config.voters[i];
      // The owner can always curse using the ownerCurse method, and is not supposed to be included in the voters list.
      // Even though the intent is for the actual owner address to NOT be included in the voters list, we don't
      // explicitly disallow curseVoteAddr == owner() here. Even if we did, the owner could transfer ownership of the
      // contract, and so we couldn't guarantee that the owner is not eventually included in the voters list.
      if (
        voter.blessVoteAddr == address(0) || voter.curseVoteAddr == address(0)
          || voter.curseVoteAddr == LIFT_CURSE_VOTE_ADDR || voter.curseVoteAddr == OWNER_CURSE_VOTE_ADDR
          || (voter.blessWeight == 0 && voter.curseWeight == 0)
      ) {
        return false;
      }
      allAddrs[2 * i + 0] = voter.blessVoteAddr;
      allAddrs[2 * i + 1] = voter.curseVoteAddr;
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
    if (!_validateConfig(config)) revert InvalidConfig();

    // We can't directly assign s_versionedConfig.config to config
    // because copying a memory array into storage is not supported.
    {
      s_versionedConfig.config.blessWeightThreshold = config.blessWeightThreshold;
      s_versionedConfig.config.curseWeightThreshold = config.curseWeightThreshold;
      while (s_versionedConfig.config.voters.length != 0) {
        Voter memory voter = s_versionedConfig.config.voters[s_versionedConfig.config.voters.length - 1];
        delete s_blesserRecords[voter.blessVoteAddr];
        delete s_curserRecords[voter.curseVoteAddr]; // usedCurseIds mapping is retained, as intended
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
      s_blesserRecords[voter.blessVoteAddr] =
        BlesserRecord({configVersion: configVersion, index: i, weight: voter.blessWeight});
      {
        CurserRecord storage sptr_curserRecord = s_curserRecords[voter.curseVoteAddr];
        // Solidity will not let us initialize as CurserRecord({...}) due to the nested mapping
        sptr_curserRecord.active = true;
        sptr_curserRecord.weight = voter.curseWeight;
      }
    }
    {
      // Initialize the owner's CurserRecord
      // We could in principle perform this initialization once in the constructor instead, and save a small bit of gas.
      // But configuration changes are relatively infrequent, and keeping the initialization here makes the contract's
      // correctness easier to reason about.
      CurserRecord storage sptr_ownerCurserRecord = s_curserRecords[OWNER_CURSE_VOTE_ADDR];
      sptr_ownerCurserRecord.active = true; // Assumed by vote/unvote-to-curse logic
      sptr_ownerCurserRecord.weight = 0; // Assumed by vote/unvote-to-curse logic
    }
    s_versionedConfig.blockNumber = uint32(block.number);
    emit ConfigSet(configVersion, config);

    s_recordedCurseRelatedOps.push(
      RecordedCurseRelatedOp({
        tag: RecordedCurseRelatedOpTag.SetConfig,
        blockTimestamp: _blockTimestamp(),
        cursed: false,
        curseVoteAddr: address(0),
        curseId: bytes16(0),
        subject: bytes16(0)
      })
    );
  }
}
