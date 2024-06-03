// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "./interfaces/IReceiver.sol";

contract KeystoneFeedsConsumer is IReceiver {
  event MessageReceived(bytes32 indexed workflowId, address indexed workflowOwner, uint256 nReports);
  event FeedReceived(bytes32 indexed feedId, uint256 price, uint64 timestamp);

  constructor() {}

  struct FeedReport {
    bytes32 FeedID;
    uint256 Price;
    uint64 Timestamp;
  }

  function onReport(bytes32 workflowId, address workflowOwner, bytes calldata rawReport) external {
    // TODO: validate sender and workflowOwner

    FeedReport[] memory feeds = abi.decode(rawReport, (FeedReport[]));
    for (uint32 i = 0; i < feeds.length; i++) {
      emit FeedReceived(feeds[i].FeedID, feeds[i].Price, feeds[i].Timestamp);
    }

    emit MessageReceived(workflowId, workflowOwner, feeds.length);
  }
}
