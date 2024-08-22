// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {KeystoneFeedsPermissionHandler} from "./KeystoneFeedsPermissionHandler.sol";
import {KeystoneFeedDefaultMetadataLib} from "./lib/KeystoneFeedDefaultMetadataLib.sol";

contract KeystoneFeedsConsumer is IReceiver, KeystoneFeedsPermissionHandler {
  using KeystoneFeedDefaultMetadataLib for bytes;

  event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp);

  struct ReceivedFeedReport {
    bytes32 FeedId;
    uint224 Price;
    uint32 Timestamp;
  }

  struct StoredFeedReport {
    uint224 Price;
    uint32 Timestamp;
  }

  mapping(bytes32 feedId => StoredFeedReport feedReport) internal s_feedReports;

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    (bytes10 workflowName, address workflowOwner, bytes2 reportName) = metadata._extractMetadataInfo();

    _validateReportPermission(msg.sender, workflowOwner, workflowName, reportName);

    ReceivedFeedReport[] memory feeds = abi.decode(rawReport, (ReceivedFeedReport[]));
    for (uint256 i = 0; i < feeds.length; ++i) {
      s_feedReports[feeds[i].FeedId] = StoredFeedReport(feeds[i].Price, feeds[i].Timestamp);
      emit FeedReceived(feeds[i].FeedId, feeds[i].Price, feeds[i].Timestamp);
    }
  }

  function getPrice(bytes32 feedId) external view returns (uint224, uint32) {
    StoredFeedReport memory report = s_feedReports[feedId];
    return (report.Price, report.Timestamp);
  }
}
