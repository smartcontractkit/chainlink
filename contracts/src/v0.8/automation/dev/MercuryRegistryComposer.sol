pragma solidity 0.8.16;

import "../../shared/access/ConfirmedOwner.sol";
import "../interfaces/AutomationCompatibleInterface.sol";
import "./interfaces/ComposerCompatibleInterfaceV1.sol";
import "../../ChainSpecificUtil.sol";
import "../../vendor/openzeppelin-contracts/contracts/utils/Strings.sol";
import "../../vendor/Strings.sol";

/*--------------------------------------------------------------------------------------------------------------------+
| Composer-Compatible Mercury Milestone 3                                                                             |
| ________________                                                                                                    |
| This implementation allows for an on-chain registry of price feed data to be maintained and updated by Automation   |
| and Functions nodes. The upkeep provides the following advantages:                                                  |                                                                                          |
| TODO:                                                                                                               |
|   - Stop using stringified performData when possible.                                                               |
|   - Optimize gas consumption.                                                                                       |
-+---------------------------------------------------------------------------------------------------------------------*/
contract MercuryRegistryComposer is ConfirmedOwner, AutomationCompatibleInterface, ComposerCompatibleInterfaceV1 {
  using strings for strings.slice;

  error DuplicateFeed(string feedId);
  error FeedNotActive(string feedId);
  error StaleReport(string feedId, uint32 currentTimestamp, uint32 incomingTimestamp);
  error InvalidFeeds();

  // Feed object used for storing feed data.
  // not included but contained in reports:
  // - blocknumberUpperBound
  // - upperBlockhash
  // - blocknumberLowerBound
  // - currentBlockTimestamp
  struct Feed {
    uint32 observationsTimestamp; // the timestamp of the most recent data assigned to this feed
    int192 price; // the current price of the feed
    int192 bid; // the current bid price of the feed
    int192 ask; // the current ask price of the feed
    string feedName; // the name of the feed
    string feedId; // the id of the feed (hex encoded)
    bool active; // true if the feed is being actively updated, otherwise false
    int192 deviationPercentagePPM; // acceptable deviation threshold - 1.5% = 15_000, 100% = 1_000_000, etc..
    uint32 stalenessSeconds; // acceptable staleness threshold - 60 = 1 minute, 300 = 5 minutes, etc..
  }

  // Report object obtained from off-chain Mercury server.
  struct Report {
    bytes32 feedId; // the feed Id of the report
    uint32 observationsTimestamp; // the timestamp of when the data was observed
    int192 price; // the median value of the OCR round
    int192 bid; // the median bid of the OCR round
    int192 ask; // the median ask if the OCR round
    uint64 blocknumberUpperBound; // the highest block observed at the time the report was generated
    bytes32 upperBlockhash; // the blockhash of the highest block observed
    uint64 blocknumberLowerBound; // the lowest block observed at the time the report was generated
    uint64 currentBlockTimestamp; // the timestamp of the highest block observed
  }

  event FeedUpdated(uint32 observationsTimestamp, int192 price, int192 bid, int192 ask, string feedId);

  uint32 private constant MIN_GAS_FOR_PERFORM = 200_000;

  string constant c_feedParamKey = "feedIdHex"; // for Mercury v0.2 - format by which feeds are identified
  string constant c_timeParamKey = "blockNumber"; // for Mercury v0.2 - format by which feeds are filtered to be sufficiently recent
  IVerifierProxy public s_verifier; // for Mercury v0.2 - verifies off-chain reports

  int192 constant scale = 1_000_000; // a scalar used for measuring deviation with precision

  string[] public s_feeds; // list of feed Ids
  mapping(string => Feed) public s_feedMapping; // mapping of feed Ids to stored feed data

  string private s_scriptHash;

  constructor(
    string[] memory feedIds,
    string[] memory feedNames,
    int192[] memory deviationPercentagePPMs,
    uint32[] memory stalenessSeconds,
    address verifier,
    string memory scriptHash
  ) ConfirmedOwner(msg.sender) {
    s_verifier = IVerifierProxy(verifier);
    s_scriptHash = scriptHash;

    // Store desired feeds.
    setFeeds(feedIds, feedNames, deviationPercentagePPMs, stalenessSeconds);
  }

  // Returns a user-defined batch of feed data, based on the on-chain state.
  function getLatestFeedData(string[] memory feedIds) external view returns (Feed[] memory) {
    Feed[] memory feeds = new Feed[](feedIds.length);
    for (uint256 i = 0; i < feedIds.length; i++) {
      feeds[i] = s_feedMapping[feedIds[i]];
    }

    return feeds;
  }

  // Invoke a feed lookup through the checkUpkeep function. Expected to run on a cron schedule.
  function checkUpkeep(bytes calldata /* data */) external view override returns (bool, bytes memory) {
    string[] memory feeds = s_feeds;
    return revertForFeedLookup(feeds);
  }

  // Extracted from `checkUpkeep` for batching purposes.
  function revertForFeedLookup(string[] memory feeds) public view returns (bool, bytes memory) {
    uint256 blockNumber = ChainSpecificUtil.getBlockNumber();
    string[] memory functionsArguments = new string[](1);

    // Pass current on-chain data as an argument to the Functions DON.
    string memory currentMercuryData = "";
    for (uint256 i; i < feeds.length; i++) {
      Feed memory feed = s_feedMapping[feeds[i]];
      string memory entry = string.concat(
        feed.feedId,
        "-",
        Strings.toString(uint192(feed.price)),
        "-",
        Strings.toString(feed.observationsTimestamp),
        "-",
        Strings.toString(uint192(feed.deviationPercentagePPM)),
        "-",
        Strings.toString(uint192(feed.stalenessSeconds))
      );
      currentMercuryData = string.concat(currentMercuryData, entry, i == feeds.length - 1 ? "" : ",");
    }
    functionsArguments[0] = currentMercuryData;

    // Emit Composer request revert.
    revert ComposerRequestV1(
      s_scriptHash,
      functionsArguments,
      true,
      c_feedParamKey,
      feeds,
      c_timeParamKey,
      blockNumber,
      ""
    );
  }

  // Modified checkCallback function that matches the StreamsLookupCompatibleInterface, but
  // accepts the result of a functions call to correctly ABI-encode it. This is a stopgap
  // function only intended to exist while Functions does not yet implement more sophisticated
  // ABI-encoding.
  function checkCallback(
    bytes[] memory data,
    bytes memory lookupData
  ) external view override returns (bool, bytes memory) {
    require(data.length == 1, "should only have one item for abi-decoding");
    string memory values = abi.decode(data[0], (string));

    // Parse the comma separated string of hex-encoded mercury proofs.
    strings.slice memory s = strings.toSlice(values);
    strings.slice memory delim = strings.toSlice(",");
    string[] memory parts = new string[](s.count(delim) + 1);
    for (uint i = 0; i < parts.length; i++) {
      parts[i] = s.split(delim).toString();
    }

    // Convert the hex strings to byte arrays.
    bytes[] memory reports = new bytes[](parts.length);
    for (uint256 i = 0; i < parts.length; i++) {
      // Convert the hex-encoded proof to bytes.
      bytes memory value = fromHex(parts[i]);
      reports[i] = value;
    }

    // Return the well-formatted performData.
    bytes memory performData = abi.encode(reports, lookupData);
    return (reports.length > 0, performData);
  }

  // Use deviated off-chain values to update on-chain state.
  function performUpkeep(bytes calldata performData) external override {
    (bytes[] memory values /* bytes memory lookupData */, ) = abi.decode(performData, (bytes[], bytes));
    for (uint256 i = 0; i < values.length; i++) {
      // Verify and decode the Mercury report.
      Report memory report = abi.decode(s_verifier.verify(values[i]), (Report));
      string memory feedId = bytes32ToHexString(abi.encodePacked(report.feedId));

      // Feeds that have been removed between checkUpkeep and performUpkeep should not be updated.
      if (!s_feedMapping[feedId].active) {
        revert FeedNotActive(feedId);
      }

      // Ensure stale reports do not cause a regression in the registry.
      if (s_feedMapping[feedId].observationsTimestamp > report.observationsTimestamp) {
        revert StaleReport(feedId, s_feedMapping[feedId].observationsTimestamp, report.observationsTimestamp);
      }

      // Assign new values to state.
      s_feedMapping[feedId].bid = report.bid;
      s_feedMapping[feedId].ask = report.ask;
      s_feedMapping[feedId].price = report.price;
      s_feedMapping[feedId].observationsTimestamp = report.observationsTimestamp;

      // Emit log.
      emit FeedUpdated(report.observationsTimestamp, report.price, report.bid, report.ask, feedId);

      // Ensure enough gas remains for the next iteration. Otherwise, stop here.
      if (gasleft() < MIN_GAS_FOR_PERFORM) {
        return;
      }
    }
  }

  // Decodes a mercury respone into an on-chain object. Thanks @mikestone!!
  function getReport(bytes memory signedReport) internal pure returns (Report memory) {
    /*
     * bytes32[3] memory reportContext,
     * bytes memory reportData,
     * bytes32[] memory rs,
     * bytes32[] memory ss,
     * bytes32 rawVs
     **/
    (, bytes memory reportData, , , ) = abi.decode(signedReport, (bytes32[3], bytes, bytes32[], bytes32[], bytes32));

    Report memory report = abi.decode(reportData, (Report));
    return report;
  }

  // Helper function to reconcile a difference in formatting:
  // - Automation passes feedId into their off-chain lookup function as a string.
  // - Mercury stores feedId in their reports as a bytes32.
  function bytes32ToHexString(bytes memory buffer) internal pure returns (string memory) {
    bytes memory converted = new bytes(buffer.length * 2);
    bytes memory _base = "0123456789abcdef";
    for (uint256 i = 0; i < buffer.length; i++) {
      converted[i * 2] = _base[uint8(buffer[i]) / _base.length];
      converted[i * 2 + 1] = _base[uint8(buffer[i]) % _base.length];
    }
    return string(abi.encodePacked("0x", converted));
  }

  // Convert an hexadecimal character to their value
  function fromHexChar(uint8 c) public pure returns (uint8) {
    if (bytes1(c) >= bytes1("0") && bytes1(c) <= bytes1("9")) {
      return c - uint8(bytes1("0"));
    }
    if (bytes1(c) >= bytes1("a") && bytes1(c) <= bytes1("f")) {
      return 10 + c - uint8(bytes1("a"));
    }
    if (bytes1(c) >= bytes1("A") && bytes1(c) <= bytes1("F")) {
      return 10 + c - uint8(bytes1("A"));
    }
    revert("fail");
  }

  // Convert an hexadecimal string to raw bytes
  function fromHex(string memory s) public pure returns (bytes memory) {
    bytes memory ss = bytes(s);
    require(ss.length % 2 == 0); // length must be even
    bytes memory r = new bytes(ss.length / 2);
    for (uint i = 0; i < ss.length / 2; ++i) {
      r[i] = bytes1(fromHexChar(uint8(ss[2 * i])) * 16 + fromHexChar(uint8(ss[2 * i + 1])));
    }
    return r;
  }

  function addFeeds(
    string[] memory feedIds,
    string[] memory feedNames,
    int192[] memory deviationPercentagePPMs,
    uint32[] memory stalenessSeconds
  ) external onlyOwner feedsAreValid(feedIds, feedNames, deviationPercentagePPMs, stalenessSeconds) {
    for (uint256 i = 0; i < feedIds.length; i++) {
      string memory feedId = feedIds[i];
      if (s_feedMapping[feedId].active) {
        revert DuplicateFeed(feedId);
      }
      updateFeed(feedId, feedNames[i], deviationPercentagePPMs[i], stalenessSeconds[i]);
      s_feedMapping[feedId].active = true;

      s_feeds.push(feedId);
    }
  }

  function setFeeds(
    string[] memory feedIds,
    string[] memory feedNames,
    int192[] memory deviationPercentagePPMs,
    uint32[] memory stalenessSeconds
  ) public onlyOwner feedsAreValid(feedIds, feedNames, deviationPercentagePPMs, stalenessSeconds) {
    // Clear prior feeds.
    for (uint256 i = 0; i < s_feeds.length; i++) {
      s_feedMapping[s_feeds[i]].active = false;
    }

    // Assign new feeds.
    for (uint256 i = 0; i < feedIds.length; i++) {
      string memory feedId = feedIds[i];
      if (s_feedMapping[feedId].active) {
        revert DuplicateFeed(feedId);
      }
      updateFeed(feedId, feedNames[i], deviationPercentagePPMs[i], stalenessSeconds[i]);
      s_feedMapping[feedId].active = true;
    }
    s_feeds = feedIds;
  }

  function updateFeed(
    string memory feedId,
    string memory feedName,
    int192 deviationPercentagePPM,
    uint32 stalnessSeconds
  ) internal {
    s_feedMapping[feedId].feedName = feedName;
    s_feedMapping[feedId].deviationPercentagePPM = deviationPercentagePPM;
    s_feedMapping[feedId].stalenessSeconds = stalnessSeconds;
    s_feedMapping[feedId].feedId = feedId;
  }

  function setVerifier(address verifier) external onlyOwner {
    s_verifier = IVerifierProxy(verifier);
  }

  modifier feedsAreValid(
    string[] memory feedIds,
    string[] memory feedNames,
    int192[] memory deviationPercentagePPMs,
    uint32[] memory stalenessSeconds
  ) {
    if (feedIds.length != feedNames.length) {
      revert InvalidFeeds();
    }
    if (feedIds.length != deviationPercentagePPMs.length) {
      revert InvalidFeeds();
    }
    if (feedIds.length != stalenessSeconds.length) {
      revert InvalidFeeds();
    }
    _;
  }
}

interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier, and bills the user if applicable.
   * @param payload The encoded data to be verified, including the signed
   * report and any metadata for billing.
   * @return verifiedReport The encoded report from the verifier.
   */
  function verify(bytes calldata payload) external payable returns (bytes memory verifiedReport);
}
