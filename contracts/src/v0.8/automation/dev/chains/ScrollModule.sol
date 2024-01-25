// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {IScrollL1GasPriceOracle} from "../../../vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol";
import "../interfaces/v2_2/IChainSpecific.sol";

contract ScrollModule is IChainSpecific {
  /// @dev SCROLL_L1_FEE_DATA_PADDING includes 120 bytes for L1 data padding for Optimism
  /// @dev according to testing, this padding allows automation registry to properly estimates L1 data fee with 3-5% buffer
  /// @dev this MAY NOT work for a different product and this may get out of date if transmit function is changed
  bytes internal constant SCROLL_L1_FEE_DATA_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev SCROLL_ORACLE_ADDR is the address of the L1GasPriceOracle precompile on Optimism.
  /// @dev reference: https://docs.scroll.io/en/developers/transaction-fees-on-scroll/#estimating-the-l1-data-fee
  address private constant SCROLL_ORACLE_ADDR = address(0x5300000000000000000000000000000000000002);
  IScrollL1GasPriceOracle internal constant SCROLL_ORACLE = IScrollL1GasPriceOracle(SCROLL_ORACLE_ADDR);

  function blockHash(uint256 blocknumber) external view returns (bytes32) {
    return blockhash(blocknumber);
  }

  function blockNumber() external view returns (uint256) {
    return block.number;
  }

  function getL1Fee(bytes calldata txCallData) external view returns (uint256) {
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }

  function getMaxL1Fee(uint256 dataSize) external view returns (uint256) {
    // fee is 4 per 0 byte, 16 per non-zero byte. Worst case we can have all non zero-bytes.
    // Instead of setting bytes to non-zero, we initialize 'new bytes' of length 4*dataSize to cover for zero bytes.
    bytes memory txCallData = new bytes(4 * dataSize);
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }
}
