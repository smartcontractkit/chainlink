// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IScrollL1GasPriceOracle} from "../../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";
import {ChainModuleBase} from "./ChainModuleBase.sol";

contract ScrollModule is ChainModuleBase {
  /// @dev SCROLL_ORACLE_ADDR is the address of the ScrollL1GasPriceOracle precompile on Scroll.
  /// @dev reference: https://docs.scroll.io/en/developers/transaction-fees-on-scroll/#estimating-the-l1-data-fee
  address private constant SCROLL_ORACLE_ADDR = 0x5300000000000000000000000000000000000002;
  IScrollL1GasPriceOracle private constant SCROLL_ORACLE = IScrollL1GasPriceOracle(SCROLL_ORACLE_ADDR);

  uint256 private constant FIXED_GAS_OVERHEAD = 15000;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 20;

  function getCurrentL1Fee() external view override returns (uint256) {
    return SCROLL_ORACLE.getL1Fee(msg.data);
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    // fee is 4 per 0 byte, 16 per non-zero byte. Worst case we can have all non zero-bytes.
    // Instead of setting bytes to non-zero, we initialize 'new bytes' of length 4*dataSize to cover for zero bytes.
    // this is the same as OP.
    bytes memory txCallData = new bytes(4 * dataSize);
    return SCROLL_ORACLE.getL1Fee(txCallData);
  }

  function getGasOverhead()
    external
    view
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, PER_CALLDATA_BYTE_GAS_OVERHEAD);
  }
}
