pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {StreamsLookupCompatibleInterface} from "../interfaces/StreamsLookupCompatibleInterface.sol";
import {ChainSpecificUtil} from "../../shared/util/ChainSpecificUtil.sol";

/*--------------------------------------------------------------------------------------------------------------------+
| Mercury + Automation                                                                                                |
| ________________                                                                                                    |
| This implementation allows for an on-chain registry of price feed data to be maintained and updated by Automation   |
| nodes. The upkeep provides the following advantages:                                                                |
|   - Node operator savings. The single committee of automation nodes is able to update all price feed data using     |
|     off-chain feed data.                                                                                            |
|   - Fetch batches of price data. All price feed data is held on the same contract, so a contract that needs         |
|     multiple sets of feed data can fetch them while paying for only one external call.                              |
|   - Scalability. Feeds can be added or removed from the contract with a single contract call, and the number of     |
|     feeds that the registry can store is unbounded.                                                                 |
|                                                                                                                     |
| Key Contracts:                                                                                                      |
|   - `MercuryRegistry.sol` - stores price feed data and implements core logic.                                       |
|   - `MercuryRegistryBatchUpkeep.sol` - enables batching for the registry.                                           |
|   - `MercuryRegistry.t.sol` - contains foundry tests to demonstrate various flows.                                  |
|                                                                                                                     |
| NOTE: This contract uses Mercury v0.2. Automation will likely upgrade to v0.3 eventually, which may change some     |
| components such as the Report struct, verification, and the StreamsLookup revert.                                   |
|                                                                                                                     |
| TODO:                                                                                                               |
|   - Optimize gas consumption.                                                                                       |
-+---------------------------------------------------------------------------------------------------------------------*/
contract MercuryRegistry is ConfirmedOwner, AutomationCompatibleInterface, StreamsLookupCompatibleInterface {
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

  string private constant FEED_PARAM_KEY = "feedIdHex"; // for Mercury v0.2 - format by which feeds are identified
  string private constant TIME_PARAM_KEY = "blockNumber"; // for Mercury v0.2 - format by which feeds are filtered to be sufficiently recent
  IVerifierProxy public s_verifier; // for Mercury v0.2 - verifies off-chain reports

  int192 private constant SCALE = 1_000_000; // a scalar used for measuring deviation with precision

  string[] public s_feeds; // list of feed Ids
  mapping(string => Feed) public s_feedMapping; // mapping of feed Ids to stored feed data

  constructor(
    string[] memory feedIds,
    string[] memory feedNames,
    int192[] memory deviationPercentagePPMs,
    uint32[] memory stalenessSeconds,
    address verifier
  ) ConfirmedOwner(msg.sender) {
    s_verifier = IVerifierProxy(verifier);

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
    uint256 blockNumber = ChainSpecificUtil._getBlockNumber();
    revert StreamsLookup(FEED_PARAM_KEY, feeds, TIME_PARAM_KEY, blockNumber, "");
  }

  // Filter for feeds that have deviated sufficiently from their respective on-chain values, or where
  // the on-chain values are sufficiently stale.
  function checkCallback(
    bytes[] memory values,
    bytes memory lookupData
  ) external view override returns (bool, bytes memory) {
    bytes[] memory filteredValues = new bytes[](values.length);
    uint256 count = 0;
    for (uint256 i = 0; i < values.length; i++) {
      Report memory report = _getReport(values[i]);
      string memory feedId = _bytes32ToHexString(abi.encodePacked(report.feedId));
      Feed memory feed = s_feedMapping[feedId];
      if (
        (report.observationsTimestamp - feed.observationsTimestamp > feed.stalenessSeconds) ||
        deviationExceedsThreshold(feed.price, report.price, feed.deviationPercentagePPM)
      ) {
        filteredValues[count] = values[i];
        count++;
      }
    }

    // Adjusts the length of the filteredValues array to `count` such that it
    // does not have extra empty slots, in case some items were filtered.
    assembly {
      mstore(filteredValues, count)
    }

    bytes memory performData = abi.encode(filteredValues, lookupData);
    return (filteredValues.length > 0, performData);
  }

  function checkErrorHandler(
    uint256 /* errCode */,
    bytes memory /* extraData */
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // dummy function with default values
    return (false, new bytes(0));
  }

  // Use deviated off-chain values to update on-chain state.
  function performUpkeep(bytes calldata performData) external override {
    (bytes[] memory values /* bytes memory lookupData */, ) = abi.decode(performData, (bytes[], bytes));
    for (uint256 i = 0; i < values.length; i++) {
      // Verify and decode the Mercury report.
      Report memory report = abi.decode(s_verifier.verify(values[i]), (Report));
      string memory feedId = _bytes32ToHexString(abi.encodePacked(report.feedId));

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
  function _getReport(bytes memory signedReport) internal pure returns (Report memory) {
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

  // Check if the off-chain value has deviated sufficiently from the on-chain value to justify an update.
  // `scale` is used to ensure precision is not lost.
  function deviationExceedsThreshold(
    int192 onChain,
    int192 offChain,
    int192 deviationPercentagePPM
  ) public pure returns (bool) {
    // Compute absolute difference between the on-chain and off-chain values.
    int192 scaledDifference = (onChain - offChain) * SCALE;
    if (scaledDifference < 0) {
      scaledDifference = -scaledDifference;
    }

    // Compare to the allowed deviation from the on-chain value.
    int192 deviationMax = ((onChain * SCALE) * deviationPercentagePPM) / SCALE;
    return scaledDifference > deviationMax;
  }

  // Helper function to reconcile a difference in formatting:
  // - Automation passes feedId into their off-chain lookup function as a string.
  // - Mercury stores feedId in their reports as a bytes32.
  function _bytes32ToHexString(bytes memory buffer) internal pure returns (string memory) {
    bytes memory converted = new bytes(buffer.length * 2);
    bytes memory _base = "0123456789abcdef";
    for (uint256 i = 0; i < buffer.length; i++) {
      converted[i * 2] = _base[uint8(buffer[i]) / _base.length];
      converted[i * 2 + 1] = _base[uint8(buffer[i]) % _base.length];
    }
    return string(abi.encodePacked("0x", converted));
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
      _updateFeed(feedId, feedNames[i], deviationPercentagePPMs[i], stalenessSeconds[i]);
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
      _updateFeed(feedId, feedNames[i], deviationPercentagePPMs[i], stalenessSeconds[i]);
      s_feedMapping[feedId].active = true;
    }
    s_feeds = feedIds;
  }

  function _updateFeed(
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
