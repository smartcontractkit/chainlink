// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BaseSequencerUptimeFeed} from "../../base/BaseSequencerUptimeFeed.sol";

contract MockBaseSequencerUptimeFeed is BaseSequencerUptimeFeed {
  string public constant override typeAndVersion = "MockSequencerUptimeFeed 1.1.0-dev";

  /// @dev this will be used for internal testing
  bool private s_validateSenderShouldPass;

  constructor(
    address l1SenderAddress,
    bool initialStatus,
    bool validateSenderShouldPass
  ) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {
    s_validateSenderShouldPass = validateSenderShouldPass;
  }

  function _validateSender(address /* l1Sender */) internal view override {
    if (!s_validateSenderShouldPass) {
      revert InvalidSender();
    }
  }
}
