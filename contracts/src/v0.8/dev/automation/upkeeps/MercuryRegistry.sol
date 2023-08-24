pragma solidity 0.8.6;

import "../../../automation/interfaces/AutomationCompatibleInterface.sol";
import "../2_1/interfaces/FeedLookupCompatibleInterface.sol";
import "../../../ChainSpecificUtil.sol";

/*--------------------------------------------------------------------------------------------------------------------+
| Mercury + Automation                                                                                                |
| ________________                                                                                                    |
| This implementation allows for an on-chain registry of price feed data to be maintained and updated by Automation   |
| nodes. The upkeep provides the following advantages:                                                                |
|   - Node operator savings. The single committee of automation nodes is able to update all price feed data using     |
|     off-chain feed data.                                                                                            |
|   - Fetch batches of price data. All price feed data is held on the same contract, so a contract that needs         |
|     multiple sets of feed data can fetch them while paying for only one external call.                              |
|   - Scalabiliy. Feeds can be added or removed from the contract with a single contract call, and the number of      |
|     feeds that the registry can store is unbounded.                                                                 |
|                                                                                                                     |
| Key Contracts:                                                                                                      |
|   - `MercuryRegistry.sol` - stores price feed data and implements core logic.                                       |
|   - `MercuryRegistryBatchUpkeep.sol` - enables batching for the registry.                                           |
|   - `MercuryRegistry.t.soll` - contains foundry tests to demonstrate various flows.                                 |
|                                                                                                                     |
| TODO:                                                                                                               |
|   - Convert s_feedIds to an enumerable set.                                                                         |
|   - Enable an owner to add and remove feeds after the contract is deployed.                                         |
|   - Access control. Specifically, the the ability to execute `performUpkeep`.                                       |
|   - Optimize gas consumption.                                                                                       |
-+---------------------------------------------------------------------------------------------------------------------*/
contract MercuryRegistry is AutomationCompatibleInterface, FeedLookupCompatibleInterface {
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

  string constant c_feedParamKey = "feedIdHex"; // for Mercury v0.2 - format by which feeds are identified
  string constant c_timeParamKey = "blockNumber"; // for Mercury v0.2 - format by which feeds are filtered to be sufficiently recent
  IVerifierProxy immutable i_verifier; // for Mercury v0.2 - verifies off-chain reports

  int192 constant scale = 1_000_000; // a scalar used for measuring deviation with precision
  int192 s_deviationPercentagePPM; // acceptable deviatoin threshold - 1.5% = 15_000, 100% = 1_000_000, etc..
  uint32 s_stalenessSeconds; // acceptable staleness threshold - 60 = 1 minute, 300 = 5 minutes, etc..

  string[] public s_feeds; // list of feed Ids
  mapping(string => Feed) public s_feedMapping; // mapping of feed Ids to stored feed data

  constructor(
    string[] memory feedIds,
    string[] memory feedNames,
    address verifier,
    int192 deviationPercentagePPM,
    uint32 stalenessSeconds
  ) {
    i_verifier = IVerifierProxy(verifier);

    // Ensure correctly formatted constructor arguments.
    require(feedIds.length == feedNames.length, "incorrect constructor args");

    // Store desired deviation threshold and staleness seconds.
    s_deviationPercentagePPM = deviationPercentagePPM;
    s_stalenessSeconds = stalenessSeconds;

    // Store desired feeds.
    s_feeds = feedIds;
    for (uint256 i = 0; i < feedIds.length; i++) {
      s_feedMapping[s_feeds[i]] = Feed({
        feedName: feedNames[i],
        feedId: feedIds[i],
        price: 0,
        bid: 0,
        ask: 0,
        observationsTimestamp: 0
      });
    }
  }

  // Returns a user-defined batch of feed data, based on the on-chain state.
  function getLatestFeedData(string[] memory feedIds) external view returns (Feed[] memory) {
    Feed[] memory feeds = new Feed[](feedIds.length);
    for (uint256 i = 0; i < feedIds.length; i++) {
      feeds[i] = s_feedMapping[feedIds[i]];
    }

    return feeds;
  }

  // Invoke a feed lookup through the checkUpkeep function. Expected to run on a chron schedule.
  function checkUpkeep(bytes calldata /* data */) external view override returns (bool, bytes memory) {
    string[] memory feeds = s_feeds;
    return revertForFeedLookup(feeds);
  }

  // Extracted from `checkUpkeep` for batching purposes.
  function revertForFeedLookup(string[] memory feeds) public view returns (bool, bytes memory) {
    uint256 blockNumber = ChainSpecificUtil.getBlockNumber();
    revert FeedLookup(c_feedParamKey, feeds, c_timeParamKey, blockNumber, "EXTRA_DATA_FOR_FUTURE_FUNCTIONS_CALLS");
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
      Report memory report = getReport(values[i]);
      string memory feedId = bytes32ToHextString(abi.encodePacked(report.feedId));
      Feed memory feed = s_feedMapping[feedId];
      if (
        (report.observationsTimestamp - feed.observationsTimestamp > s_stalenessSeconds) ||
        deviationExceedsThreshold(feed.price, report.price)
      ) {
        filteredValues[count] = values[i];
        count++;
      }
    }

    // Adjusts the lenght of the filteredValues array to `count` such that it
    // does not have extra empty slots, in case some items were filtered.
    assembly {
      mstore(filteredValues, count)
    }

    bytes memory performData = abi.encode(filteredValues, lookupData);
    return (filteredValues.length > 0, performData);
  }

  // Use deviated off-chain values to update on-chain state.
  // TODO:
  // - The implementation provided here is readable but crude. Remaining gas should be checked between iterations
  // of the for-loop, and the failure of a single item should not cause the entire batch to revert.
  function performUpkeep(bytes calldata performData) external override {
    (bytes[] memory values /* bytes memory lookupData */, ) = abi.decode(performData, (bytes[], bytes));
    for (uint256 i = 0; i < values.length; i++) {
      // Verify and decode report.
      i_verifier.verify(values[i]);
      Report memory report = getReport(values[i]);
      string memory feedId = bytes32ToHextString(abi.encodePacked(report.feedId));

      // Feeds that have been removed between checkUpkeep and performUpkeep should not be updated.
      require(bytes(s_feedMapping[feedId].feedId).length > 0, "feed removed");

      // Sanity check. Stale reports should not get through, but ensure they do not cause a regression
      // in the registry.
      require(s_feedMapping[feedId].observationsTimestamp <= report.observationsTimestamp, "stale report");

      // Assign new values to state.
      s_feedMapping[feedId].bid = report.bid;
      s_feedMapping[feedId].ask = report.ask;
      s_feedMapping[feedId].price = report.price;
      s_feedMapping[feedId].observationsTimestamp = report.observationsTimestamp;

      // Emit log (not gas efficient to do this for each update).
      emit FeedUpdated(report.observationsTimestamp, report.price, report.bid, report.ask, feedId);
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

  // Check if the off-chain value has deviated sufficiently from the on-chain value to justify an update.
  // `scale` is used to ensure precision is not lost.
  function deviationExceedsThreshold(int192 onChain, int192 offChain) public view returns (bool) {
    // Compute absolute difference between the on-chain and off-chain values.
    int192 scaledDifference = (onChain - offChain) * scale;
    if (scaledDifference < 0) {
      scaledDifference = -scaledDifference;
    }

    // Compare to the allowed deviation from the on-chain value.
    int192 deviationMax = ((onChain * scale) * s_deviationPercentagePPM) / scale;
    return scaledDifference > deviationMax;
  }

  // Helper function to reconcile a difference in formatting:
  // - Automation passes feedId into their off-chain lookup function as a string.
  // - Mercury stores feedId in their reports as a bytes32.
  function bytes32ToHextString(bytes memory buffer) internal pure returns (string memory) {
    bytes memory converted = new bytes(buffer.length * 2);
    bytes memory _base = "0123456789abcdef";
    for (uint256 i = 0; i < buffer.length; i++) {
      converted[i * 2] = _base[uint8(buffer[i]) / _base.length];
      converted[i * 2 + 1] = _base[uint8(buffer[i]) % _base.length];
    }
    return string(abi.encodePacked("0x", converted));
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
