// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../MultiCommitStore.sol";

contract MultiCommitStoreHelper is MultiCommitStore {
  constructor(
    StaticConfig memory staticConfig,
    SourceChainConfigArgs[] memory srcConfigs
  ) MultiCommitStore(staticConfig, srcConfigs) {}

  /// @dev Expose _report for tests
  function report(bytes calldata commitReport, uint40 epochAndRound) external {
    _report(commitReport, epochAndRound);
  }
}
