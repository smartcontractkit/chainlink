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

  address internal s_registry;

  constructor() {}

  /**
   * @inheritdoc OCR2DRBillableInterface
   */
  function getRequiredFee(
    bytes calldata, /* data */
    OCR2DRRegistryInterface.RequestBilling calldata /* billing */
  ) public pure override returns (uint32) {
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
    if (s_registry == address(0)) {
      revert EmptyBillingRegistry();
    }
    uint32 requiredFee = getRequiredFee(data, billing);
    OCR2DRRegistryInterface registry = OCR2DRRegistryInterface(s_registry);
    return registry.estimateCost(data, billing, requiredFee);
  }

  modifier validateRequestId(bytes32 requestId) {
    if (s_registry == address(0)) {
      revert EmptyBillingRegistry();
    }
    OCR2DRRegistryInterface registry = OCR2DRRegistryInterface(s_registry);
    (address consumer, , , ) = registry.getCommitment(requestId);
    if (consumer == address(0)) {
      revert InvalidRequestID();
    }
    _;
  }
}
