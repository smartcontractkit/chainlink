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

  /// @dev L1 fee coefficient can be applied to reduce possibly inflated gas price
  uint8 public s_l1FeeCoefficient = 100;
  uint256 private constant FIXED_GAS_OVERHEAD = 60_000;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 270;

  /// @dev This is the padding size for unsigned RLP-encoded transaction without the signature data
  /// @dev The padding size was estimated based on hypothetical max RLP-encoded transaction size
  uint256 private constant L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE = 71;

  constructor() ConfirmedOwner(msg.sender) {}

  function getCurrentL1Fee() external view override returns (uint256) {
    return (s_l1FeeCoefficient * _getL1Fee(msg.data.length)) / 100;
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    return _getL1Fee(dataSize);
  }

  function _getL1Fee(uint256 dataSize) internal view returns (uint256) {
    // getL1FeeUpperBound expects unsigned fully RLP-encoded transaction size so we have to account for paddding bytes as well
    return OVM_GASPRICEORACLE.getL1FeeUpperBound(dataSize + L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE);
  }

  function getGasOverhead()
    external
    pure
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, PER_CALLDATA_BYTE_GAS_OVERHEAD);
  }

  function setL1FeeCalculation(uint8 coefficient) external onlyOwner {
    if (coefficient > 100) {
      revert InvalidL1FeeCoefficient(coefficient);
    }

    s_l1FeeCoefficient = coefficient;

    emit L1FeeCoefficientSet(coefficient);
  }
}
