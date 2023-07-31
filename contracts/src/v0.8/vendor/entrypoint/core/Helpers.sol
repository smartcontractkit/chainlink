// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.12;

/**
 * returned data from validateUserOp.
 * validateUserOp returns a uint256, with is created by `_packedValidationData` and parsed by `_parseValidationData`
 * @param aggregator - address(0) - the account validated the signature by itself.
 *              address(1) - the account failed to validate the signature.
 *              otherwise - this is an address of a signature aggregator that must be used to validate the signature.
 * @param validAfter - this UserOp is valid only after this timestamp.
 * @param validaUntil - this UserOp is valid only up to this timestamp.
 */
struct ValidationData {
  address aggregator;
  uint48 validAfter;
  uint48 validUntil;
}

//extract sigFailed, validAfter, validUntil.
// also convert zero validUntil to type(uint48).max
function _parseValidationData(uint validationData) pure returns (ValidationData memory data) {
  address aggregator = address(uint160(validationData));
  uint48 validUntil = uint48(validationData >> 160);
  if (validUntil == 0) {
    validUntil = type(uint48).max;
  }
  uint48 validAfter = uint48(validationData >> (48 + 160));
  return ValidationData(aggregator, validAfter, validUntil);
}

// intersect account and paymaster ranges.
function _intersectTimeRange(
  uint256 validationData,
  uint256 paymasterValidationData
) pure returns (ValidationData memory) {
  ValidationData memory accountValidationData = _parseValidationData(validationData);
  ValidationData memory pmValidationData = _parseValidationData(paymasterValidationData);
  address aggregator = accountValidationData.aggregator;
  if (aggregator == address(0)) {
    aggregator = pmValidationData.aggregator;
  }
  uint48 validAfter = accountValidationData.validAfter;
  uint48 validUntil = accountValidationData.validUntil;
  uint48 pmValidAfter = pmValidationData.validAfter;
  uint48 pmValidUntil = pmValidationData.validUntil;

  if (validAfter < pmValidAfter) validAfter = pmValidAfter;
  if (validUntil > pmValidUntil) validUntil = pmValidUntil;
  return ValidationData(aggregator, validAfter, validUntil);
}

/**
 * helper to pack the return value for validateUserOp
 * @param data - the ValidationData to pack
 */
function _packValidationData(ValidationData memory data) pure returns (uint256) {
  return uint160(data.aggregator) | (uint256(data.validUntil) << 160) | (uint256(data.validAfter) << (160 + 48));
}

/**
 * helper to pack the return value for validateUserOp, when not using an aggregator
 * @param sigFailed - true for signature failure, false for success
 * @param validUntil last timestamp this UserOperation is valid (or zero for infinite)
 * @param validAfter first timestamp this UserOperation is valid
 */
function _packValidationData(bool sigFailed, uint48 validUntil, uint48 validAfter) pure returns (uint256) {
  return (sigFailed ? 1 : 0) | (uint256(validUntil) << 160) | (uint256(validAfter) << (160 + 48));
}
