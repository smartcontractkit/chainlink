// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {OVM_GasPriceOracle} from "./../../../vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol";
import "../interfaces/v2_2/IChainSpecific.sol";

contract OptimismModule is IChainSpecific {
  /// @dev OP_L1_DATA_FEE_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes internal constant OP_L1_DATA_FEE_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  function _blockHash(uint256 blockNumber) external view returns (bytes32) {
    return blockhash(blockNumber);
  }

  function _blockNumber() external view returns (uint256) {
    return block.number;
  }

  function _getL1Fee(bytes calldata txCallData) external view returns (uint256) {
    return OVM_GASPRICEORACLE.getL1Fee(bytes.concat(txCallData, OP_L1_DATA_FEE_PADDING));
  }
}
