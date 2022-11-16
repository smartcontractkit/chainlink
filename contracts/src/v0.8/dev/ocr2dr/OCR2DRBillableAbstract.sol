// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRBillableInterface.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title OCR2DR billable oracle abstract.
 */
abstract contract OCR2DRBillableAbstract is OCR2DRBillableInterface {
  error EmptyBillingRegistry();
  error InvalidRequestID();

  OCR2DRRegistryInterface internal s_registry;

  constructor() {}

  /**
   * @inheritdoc OCR2DRBillableInterface
   */
  function getRequiredFee(
    bytes calldata, /* data */
    OCR2DRRegistryInterface.RequestBilling calldata /* billing */
  ) public pure override returns (uint96) {
    // NOTE: Optionally, compute additional fee here
    return 0;
  }

  /**
   * @inheritdoc OCR2DRBillableInterface
   */
  function estimateCost(bytes calldata data, OCR2DRRegistryInterface.RequestBilling calldata billing)
    external
    view
    override
    returns (uint96)
  {
    if (address(s_registry) == address(0)) {
      revert EmptyBillingRegistry();
    }
    uint96 requiredFee = getRequiredFee(data, billing);
    return s_registry.estimateCost(data, billing, requiredFee);
  }

  modifier validateRequestId(bytes32 requestId) {
    if (address(s_registry) == address(0)) {
      revert EmptyBillingRegistry();
    }
    (address consumer, , ) = s_registry.getCommitment(requestId);
    if (consumer == address(0)) {
      revert InvalidRequestID();
    }
    _;
  }
}
