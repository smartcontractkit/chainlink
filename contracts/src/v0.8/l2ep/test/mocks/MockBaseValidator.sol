// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {BaseValidator} from "../../base/BaseValidator.sol";

contract MockBaseValidator is BaseValidator {
  string public constant override typeAndVersion = "MockValidator 1.1.0-dev";

  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit
  ) BaseValidator(l1CrossDomainMessengerAddress, l2UptimeFeedAddr, gasLimit) {}

  function _validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 /* currentAnswer */
  ) internal virtual override returns (bool) {
    return true;
  }
}
