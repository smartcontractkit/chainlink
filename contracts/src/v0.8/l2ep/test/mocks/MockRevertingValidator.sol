// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {MockBaseValidator} from "./MockBaseValidator.sol";

contract MockRevertingValidator is MockBaseValidator {
  constructor(
    address _l1Messenger,
    address _l2UptimeFeed,
    uint32 _gasLimit
  ) MockBaseValidator(_l1Messenger, _l2UptimeFeed, _gasLimit) {}

  function _validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 /* currentAnswer */
  ) internal pure override returns (bool) {
    // solhint-disable-next-line gas-custom-errors
    revert("Mock revert");
  }
}
