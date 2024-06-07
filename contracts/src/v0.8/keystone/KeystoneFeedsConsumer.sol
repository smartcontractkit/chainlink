// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiver} from "./interfaces/IReceiver.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";

contract KeystoneFeedsConsumer is IReceiver, ConfirmedOwner {
  event FeedReceived(bytes32 indexed feedId, int192 price, uint32 timestamp);

  error UnauthorizedSender(address sender);
  error UnauthorizedWorkflowOwner(address workflowOwner);
  error UnauthorizedWorkflowName(bytes10 workflowName);

  constructor() ConfirmedOwner(msg.sender) {}

  struct ReceivedFeedReport {
    bytes32 FeedId;
    int192 Price;
    uint32 Timestamp;
  }

  struct StoredFeedReport {
    int192 Price;
    uint32 Timestamp;
  }

  mapping(bytes32 feedId => StoredFeedReport feedReport) internal s_feedReports;
  address[] internal s_allowedSendersList;
  mapping(address => bool) internal s_allowedSenders;
  address[] internal s_allowedWorkflowOwnersList;
  mapping(address => bool) internal s_allowedWorkflowOwners;
  bytes10[] internal s_allowedWorkflowNamesList;
  mapping(bytes10 => bool) internal s_allowedWorkflowNames;

  function setConfig(
    address[] calldata _allowedSendersList,
    address[] calldata _allowedWorkflowOwnersList,
    bytes10[] calldata _allowedWorkflowNamesList
  ) external onlyOwner {
    for (uint32 i = 0; i < s_allowedSendersList.length; i++) {
      s_allowedSenders[s_allowedSendersList[i]] = false;
    }
    for (uint32 i = 0; i < _allowedSendersList.length; i++) {
      s_allowedSenders[_allowedSendersList[i]] = true;
    }
    s_allowedSendersList = _allowedSendersList;
    for (uint32 i = 0; i < s_allowedWorkflowOwnersList.length; i++) {
      s_allowedWorkflowOwners[s_allowedWorkflowOwnersList[i]] = false;
    }
    for (uint32 i = 0; i < _allowedWorkflowOwnersList.length; i++) {
      s_allowedWorkflowOwners[_allowedWorkflowOwnersList[i]] = true;
    }
    s_allowedWorkflowOwnersList = _allowedWorkflowOwnersList;
    for (uint32 i = 0; i < s_allowedWorkflowNamesList.length; i++) {
      s_allowedWorkflowNames[s_allowedWorkflowNamesList[i]] = false;
    }
    for (uint32 i = 0; i < _allowedWorkflowNamesList.length; i++) {
      s_allowedWorkflowNames[_allowedWorkflowNamesList[i]] = true;
    }
    s_allowedWorkflowNamesList = _allowedWorkflowNamesList;
  }

  function onReport(bytes calldata metadata, bytes calldata rawReport) external {
    if (s_allowedSenders[msg.sender] == false) {
      revert UnauthorizedSender(msg.sender);
    }

    (bytes10 workflowName, address workflowOwner) = _getInfo(metadata);
    if (s_allowedWorkflowNames[workflowName] == false) {
      revert UnauthorizedWorkflowName(workflowName);
    }
    if (s_allowedWorkflowOwners[workflowOwner] == false) {
      revert UnauthorizedWorkflowOwner(workflowOwner);
    }

    ReceivedFeedReport[] memory feeds = abi.decode(rawReport, (ReceivedFeedReport[]));
    for (uint32 i = 0; i < feeds.length; i++) {
      s_feedReports[feeds[i].FeedId] = StoredFeedReport(feeds[i].Price, feeds[i].Timestamp);
      emit FeedReceived(feeds[i].FeedId, feeds[i].Price, feeds[i].Timestamp);
    }
  }

  // solhint-disable-next-line chainlink-solidity/explicit-returns
  function _getInfo(bytes memory metadata) internal pure returns (bytes10 workflowName, address workflowOwner) {
    // (first 32 bytes contain length of the byte array)
    // workflow_cid             // offset 32, size 32
    // workflow_name            // offset 64, size 10
    // workflow_owner           // offset 74, size 20
    // report_name              // offset 94, size  2
    assembly {
      // no shifting needed for bytes10 type
      workflowName := mload(add(metadata, 64))
      // shift right by 12 bytes to get the actual value
      workflowOwner := shr(mul(12, 8), mload(add(metadata, 74)))
    }
  }

  function getPrice(bytes32 feedId) external view returns (int192, uint32) {
    StoredFeedReport memory report = s_feedReports[feedId];
    return (report.Price, report.Timestamp);
  }
}
