// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {GasPriceOracle as OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";
import {ChainModuleBase} from "./ChainModuleBase.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";

/**
 * @notice OptimismModuleV2 provides a cost-efficient way to get L1 fee on OP stack.
 * After EIP-4844 is implemented in OP stack, the new OP upgrade includes a new function getL1FeeUpperBound to estimate
 * the upper bound of current transaction's L1 fee.
 */
contract OptimismModuleV2 is ChainModuleBase, ConfirmedOwner {
  error InvalidL1FeeCoefficient(uint8 coefficient);
  event L1FeeCoefficientSet(uint8 coefficient);

  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  /// @dev L1 fee coefficient is used to account for the impact of data compression on the l1 fee
  /// getL1FeeUpperBound returns the upper bound of l1 fee so this configurable coefficient will help
  /// charge a predefined percentage of the upper bound.
  uint8 private s_l1FeeCoefficient = 100;
  uint256 private constant FIXED_GAS_OVERHEAD = 28_000;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 0;

  constructor() ConfirmedOwner(msg.sender) {}

  function getCurrentL1Fee(uint256 dataSize) external view override returns (uint256) {
    return (s_l1FeeCoefficient * _getL1Fee(dataSize)) / 100;
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    return _getL1Fee(dataSize);
  }

  /* @notice this function provides an estimation for L1 fee incurred by calldata of a certain size
   * @dev this function uses the newly provided getL1FeeUpperBound function in OP gas price oracle. this helps
   * estimate L1 fee with much lower cost
   * @param dataSize the size of calldata
   * @return l1Fee the L1 fee
   */
  function _getL1Fee(uint256 dataSize) internal view returns (uint256) {
    return OVM_GASPRICEORACLE.getL1FeeUpperBound(dataSize);
  }

  function getGasOverhead()
    external
    pure
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, PER_CALLDATA_BYTE_GAS_OVERHEAD);
  }

  /* @notice this function sets a new coefficient for L1 fee estimation.
   * @dev this function can only be invoked by contract owner
   * @param coefficient the new coefficient
   */
  function setL1FeeCalculation(uint8 coefficient) external onlyOwner {
    if (coefficient > 100) {
      revert InvalidL1FeeCoefficient(coefficient);
    }

    s_l1FeeCoefficient = coefficient;

    emit L1FeeCoefficientSet(coefficient);
  }

  /* @notice this function returns the s_l1FeeCoefficient
   * @return coefficient the current s_l1FeeCoefficient in effect
   */
  function getL1FeeCoefficient() public view returns (uint256 coefficient) {
    return s_l1FeeCoefficient;
  }
}
