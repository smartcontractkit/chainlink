// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ICommitStore} from "../../interfaces/ICommitStore.sol";

contract MockCommitStore is ICommitStore {
  error PausedError();

  uint64 private s_expectedNextSequenceNumber = 1;

  bool private s_paused = false;

  /// @inheritdoc ICommitStore
  function verify(
    bytes32[] calldata,
    bytes32[] calldata,
    uint256
  ) external view whenNotPaused returns (uint256 timestamp) {
    return 1;
  }

  function getExpectedNextSequenceNumber() external view returns (uint64) {
    return s_expectedNextSequenceNumber;
  }

  function setExpectedNextSequenceNumber(uint64 nextSeqNum) external {
    s_expectedNextSequenceNumber = nextSeqNum;
  }

  modifier whenNotPaused() {
    if (paused()) revert PausedError();
    _;
  }

  function paused() public view returns (bool) {
    return s_paused;
  }

  function pause() external {
    s_paused = true;
  }
}
