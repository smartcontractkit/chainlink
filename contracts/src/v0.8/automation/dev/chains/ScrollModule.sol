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

  function _blockHash(uint256 blockNumber) external view returns (bytes32) {
    return blockhash(blockNumber);
  }

  function _blockNumber() external view returns (uint256) {
    return block.number;
  }

  function _getL1FeeForTransaction(bytes calldata txCallData) external view returns (uint256) {
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }

  function _getL1FeeForSimulation(bytes calldata txCallData) external view returns (uint256) {
    return SCROLL_ORACLE.getL1Fee(bytes.concat(txCallData, SCROLL_L1_FEE_DATA_PADDING));
  }
}
