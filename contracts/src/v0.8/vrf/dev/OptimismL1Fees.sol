// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {GasPriceOracle as OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";

/// @dev An abstract contract that provides Optimism specific L1 fee calculations.
// solhint-disable-next-line contract-name-camelcase
abstract contract OptimismL1Fees is ConfirmedOwner {
  /// @dev This is the padding size for unsigned RLP-encoded transaction without the signature data
  /// @dev The padding size was estimated based on looking at existing transactions on Optimism
  uint256 internal constant L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE = 71;
  /// @dev Signature data size used in the GasPriceOracle predeploy
  /// @dev reference: https://github.com/ethereum-optimism/optimism/blob/a96cbe7c8da144d79d4cec1303d8ae60a64e681e/packages/contracts-bedrock/contracts/L2/GasPriceOracle.sol#L145
  uint256 internal constant L1_TX_SIGNATURE_DATA_BYTES_SIZE = 68;
  /// @dev L1_FEE_DATA_PADDING includes 71 bytes for L1 data padding for Optimism
  bytes internal constant L1_FEE_DATA_PADDING =
    hex"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the OVM_GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  OVM_GasPriceOracle private constant OVM_GASPRICEORACLE = OVM_GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  /// @dev Option 1: getL1Fee() function from predeploy GasPriceOracle contract with the fulfillment calldata payload
  /// @dev This option is only available for the Coordinator contract
  uint8 internal constant L1GasFeesMode = 0;
  /// @dev Option 2: our own implementation of getL1Fee() function (Ecotone version) with projected
  /// @dev fulfillment calldata payload (number of non-zero bytes estimated based on historical data)
  /// @dev This option is available for the Coordinator and the Wrapper contract
  uint8 internal constant L1CalldataGasCostMode = 1;
  /// @dev Option 3: getL1FeeUpperBound() function from predeploy GasPriceOracle contract (available after Fjord upgrade)
  /// @dev This option is available for the Coordinator and the Wrapper contract
  uint8 internal constant L1GasFeesUpperBoundMode = 2;

  uint8 public s_L1FeeCalculationMode = L1GasFeesMode;

  error InvalidL1FeeCalculationMode(uint8 mode);

  event L1FeeCalculationModeSet(uint8 mode);

  function setL1FeeCalculationMode(uint8 mode) external onlyOwner {
    if (mode >= 3) {
      revert InvalidL1FeeCalculationMode(mode);
    }

    s_L1FeeCalculationMode = mode;

    emit L1FeeCalculationModeSet(mode);
  }

  function _getL1CostWeiForCalldata(bytes calldata data) internal view virtual returns (uint256) {
    if (s_L1FeeCalculationMode == L1GasFeesMode) {
      return OVM_GASPRICEORACLE.getL1Fee(bytes.concat(data, L1_FEE_DATA_PADDING));
    }
    return _getL1CostWeiForCalldataSize(data.length);
  }

  function _getL1CostWeiForCalldataSize(uint256 calldataSizeBytes) internal view virtual returns (uint256) {
    if (s_L1FeeCalculationMode == L1CalldataGasCostMode) {
      return _calculateOptimismL1DataFee(calldataSizeBytes);
    } else if (s_L1FeeCalculationMode == L1GasFeesUpperBoundMode) {
      // getL1FeeUpperBound expects unsigned fully RLP-encoded transaction size so we have to account for paddding bytes as well
      return OVM_GASPRICEORACLE.getL1FeeUpperBound(calldataSizeBytes + L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE);
    }
    revert InvalidL1FeeCalculationMode(s_L1FeeCalculationMode);
  }

  function _calculateOptimismL1DataFee(uint256 calldataSizeBytes) internal view returns (uint256) {
    // reference: https://docs.optimism.io/stack/transactions/fees#ecotone
    // also: https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/exec-engine.md#ecotone-l1-cost-fee-changes-eip-4844-da
    // we treat all bytes in the calldata payload as non-zero bytes (cost: 16 gas) because accurate estimation is too expensive
    // we also have to account for the padding data and the signature data
    uint256 l1GasUsed = (calldataSizeBytes + L1_TX_SIGNATURE_DATA_BYTES_SIZE + L1_UNSIGNED_RLP_ENC_TX_DATA_BYTES_SIZE) *
      16;
    uint256 scaledBaseFee = OVM_GASPRICEORACLE.baseFeeScalar() * 16 * OVM_GASPRICEORACLE.l1BaseFee();
    uint256 scaledBlobBaseFee = OVM_GASPRICEORACLE.blobBaseFeeScalar() * OVM_GASPRICEORACLE.blobBaseFee();
    uint256 fee = l1GasUsed * (scaledBaseFee + scaledBlobBaseFee);
    return fee / (16 * 10 ** OVM_GASPRICEORACLE.decimals());
  }
}
