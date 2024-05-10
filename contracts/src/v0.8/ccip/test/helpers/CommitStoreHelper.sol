// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../CommitStore.sol";

contract CommitStoreHelper is CommitStore {
  constructor(StaticConfig memory staticConfig) CommitStore(staticConfig) {}

  /// @dev Expose _report for tests
  function report(bytes calldata commitReport, uint40 epochAndRound) external {
    _report(commitReport, epochAndRound);
  }
}
