// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IScrollL1GasPriceOracle} from "../../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";
import {ChainModuleBase} from "./ChainModuleBase.sol";

contract ScrollModule is ChainModuleBase {
  /// @dev SCROLL_L1_FEE_DATA_PADDING includes 120 bytes for L1 data padding for Optimism
  /// @dev according to testing, this padding allows automation registry to properly estimates L1 data fee with 3-5% buffer
  /// @dev this MAY NOT work for a different product and this may get out of date if transmit function is changed
  bytes private constant SCROLL_L1_FEE_DATA_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev SCROLL_ORACLE_ADDR is the address of the L1GasPriceOracle precompile on Optimism.
  /// @dev reference: https://docs.scroll.io/en/developers/transaction-fees-on-scroll/#estimating-the-l1-data-fee
  address private constant SCROLL_ORACLE_ADDR = 0x5300000000000000000000000000000000000002;
  IScrollL1GasPriceOracle private constant SCROLL_ORACLE = IScrollL1GasPriceOracle(SCROLL_ORACLE_ADDR);

  function getCurrentL1Fee() external view override returns (uint256) {
    // TODO: Verify this is accruate calcualtion with appropriate padding
    return SCROLL_ORACLE.getL1Fee(bytes.concat(msg.data, SCROLL_L1_FEE_DATA_PADDING));
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    // fee is 4 per 0 byte, 16 per non-zero byte. Worst case we can have all non zero-bytes.
    // Instead of setting bytes to non-zero, we initialize 'new bytes' of length 4*dataSize to cover for zero bytes.
    // this is the same as OP.
    // TODO: Verify this is an accurate estimate
    bytes memory txCallData = new bytes(4 * dataSize);
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }

  function getGasOverhead()
    external
    view
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    // TODO: Calculate
    return (0, 0);
  }
}
