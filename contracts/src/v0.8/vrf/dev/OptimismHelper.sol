// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IOptimismHelper} from "./interfaces/IOptimismHelper.sol";
import {OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol";
import {L1Block} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.1/src/L2/L1Block.sol";

contract OptimismHelper is IOptimismHelper {
  // sets the mode for calculating L1 gas
  uint8 s_l1_gas_calculation_mode;

  /// @dev L1_FEE_DATA_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes internal constant L1_FEE_DATA_PADDING =
    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);
  IGetL1FeeUpperBoundClient private constant GET_L1_FEE_UPPER_BOUND_CLIENT = IGetL1FeeUpperBoundClient(OVM_GASPRICEORACLE_ADDR);
  /// @dev L1BLOCK_ADDR is the address of the L1Block precompile on Optimism.
  address private constant L1BLOCK_ADDR = address(0x4200000000000000000000000000000000000015);
  L1Block private constant L1BLOCK = L1Block(L1BLOCK_ADDR);
  /// @dev included L1_FEE_DATA_PADDING and 68 bytes for the transaction signature
  /// @dev reference: https://github.com/ethereum-optimism/optimism/blob/233ede59d16cb01bdd8e7ff662a153a4c3178bdd/packages/contracts-bedrock/contracts/L2/GasPriceOracle.sol#L110
  uint256 private constant OP_L1_TX_PADDING_BYTES = 103;


  error InvalidMode();

  constructor(uint8 l1_gas_calculation_mode) {
    s_l1_gas_calculation_mode = l1_gas_calculation_mode;
  }

  function getTxL1GasFees(
    bytes memory data
  ) external view override returns (uint256) {
    if (s_l1_gas_calculation_mode == 0) {
      return OVM_GASPRICEORACLE.getL1Fee(bytes.concat(data, L1_FEE_DATA_PADDING));
    } else if (s_l1_gas_calculation_mode == 1) {
      return _calculateOptimismL1DataFee(data.length);
    } else if (s_l1_gas_calculation_mode == 2) {
      return GET_L1_FEE_UPPER_BOUND_CLIENT.getL1FeeUpperBound(data.length);
    }
    revert InvalidMode();
  }

  function _calculateOptimismL1DataFee(uint256 calldataSizeBytes) internal view returns (uint256) {
    // reference: https://docs.optimism.io/stack/transactions/fees#ecotone
    // also: https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/exec-engine.md#ecotone-l1-cost-fee-changes-eip-4844-da
    // we assume the worst-case scenario and treat all bytes in the calldata payload as non-zero bytes (cost: 16 gas)
    uint256 scaledBaseFee = 16 * L1BLOCK.baseFeeScalar() * OVM_GASPRICEORACLE.l1BaseFee();
    uint256 scaledBlobBaseFee = L1BLOCK.blobBaseFeeScalar() * L1BLOCK.blobBaseFee();
    uint256 fee = (calldataSizeBytes + OP_L1_TX_PADDING_BYTES) * (scaledBaseFee + scaledBlobBaseFee);
    return fee / (10 ** OVM_GASPRICEORACLE.decimals());
  }
}

interface IGetL1FeeUpperBoundClient {
  function getL1FeeUpperBound(uint256 _unsignedTxSize) external view returns (uint256);
}