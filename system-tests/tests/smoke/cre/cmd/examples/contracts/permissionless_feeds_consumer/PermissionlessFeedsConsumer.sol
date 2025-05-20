// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {IReceiver} from "../dependencies/interfaces/IReceiver.sol";

import {OwnerIsCreator} from "../dependencies/access/OwnerIsCreator.sol";

import {IERC165} from "../dependencies/vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

contract PermissionlessFeedsConsumer is IReceiver, OwnerIsCreator {
  event FeedReceived(bytes32 indexed feedId, uint224 price, uint32 timestamp);

  struct ReceivedFeedReport {
    bytes32 FeedId;
    uint32 Timestamp;
    uint224 Price;
  }

  struct StoredFeedReport {
    uint224 Price;
    uint32 Timestamp;
  }

  mapping(bytes32 feedId => StoredFeedReport feedReport) internal s_feedReports;

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
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

  function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}