// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IScrollL1GasPriceOracle} from "../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";
import {ChainModuleBase} from "./ChainModuleBase.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";

contract ScrollModule is ChainModuleBase, ConfirmedOwner {
  error InvalidL1FeeCoefficient(uint8 coefficient);
  event L1FeeCoefficientSet(uint8 coefficient);

  /// @dev SCROLL_L1_FEE_DATA_PADDING includes 140 bytes for L1 data padding for Scroll
  /// @dev according to testing, this padding allows automation registry to properly estimates L1 data fee with 3-5% buffer
  /// @dev this MAY NOT work for a different product and this may get out of date if transmit function is changed
  bytes private constant SCROLL_L1_FEE_DATA_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev SCROLL_ORACLE_ADDR is the address of the ScrollL1GasPriceOracle precompile on Scroll.
  /// @dev reference: https://docs.scroll.io/en/developers/transaction-fees-on-scroll/#estimating-the-l1-data-fee
  address private constant SCROLL_ORACLE_ADDR = 0x5300000000000000000000000000000000000002;
  IScrollL1GasPriceOracle private constant SCROLL_ORACLE = IScrollL1GasPriceOracle(SCROLL_ORACLE_ADDR);

  /// @dev L1 fee coefficient can be applied to reduce possibly inflated gas cost
  uint8 private s_l1FeeCoefficient = 100;
  uint256 private constant FIXED_GAS_OVERHEAD = 45_000;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 170;

  constructor() ConfirmedOwner(msg.sender) {}

  function getCurrentL1Fee(uint256 dataSize) external view override returns (uint256) {
    return (s_l1FeeCoefficient * _getL1Fee(dataSize)) / 100;
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    return _getL1Fee(dataSize);
  }

  /* @notice this function provides an estimation for L1 fee incurred by calldata of a certain size
   * @param dataSize the size of calldata
   * @return l1Fee the L1 fee
   */
  function _getL1Fee(uint256 dataSize) internal view returns (uint256 l1Fee) {
    // fee is 4 per 0 byte, 16 per non-zero byte. Worst case we can have all non zero-bytes.
    // Instead of setting bytes to non-zero, we initialize 'new bytes' of length 4*dataSize to cover for zero bytes.
    // this is the same as OP.
    bytes memory txCallData = new bytes(4 * dataSize);
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }

  function getGasOverhead()
    external
    view
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
