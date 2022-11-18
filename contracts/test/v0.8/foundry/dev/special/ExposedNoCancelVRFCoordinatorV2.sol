// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {NoCancelVRFCoordinatorV2} from "../../../../../src/v0.8/dev/special/NoCancelVRFCoordinatorV2.sol";

contract ExposedNoCancelVRFCoordinatorV2 is NoCancelVRFCoordinatorV2 {
  constructor(
    address link,
    address blockhashStore,
    address linkEthFeed
  )
    // solhint-disable-next-line no-empty-blocks
    NoCancelVRFCoordinatorV2(link, blockhashStore, linkEthFeed)
  {
    /* empty */
  }

  function calculatePaymentAmountTest(
    uint256 gasAfterPaymentCalculation,
    uint32 fulfillmentFlatFeeLinkPPM,
    uint256 weiPerUnitGas
  ) external returns (uint96) {
    return calculatePaymentAmount(
      gasleft(),
      gasAfterPaymentCalculation,
      fulfillmentFlatFeeLinkPPM,
      weiPerUnitGas
    );
  }
}
