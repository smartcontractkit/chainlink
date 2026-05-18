// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";

/// @dev An abstract contract that provides Arbitrum specific L1 fee calculations.
// solhint-disable-next-line contract-name-camelcase
abstract contract ArbitrumL1Fees {
  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  /**
   * @notice Returns the L1 fees that will be paid for the current transaction, given any calldata
   * @notice for the current transaction. It uses ArbGas.getCurrentTxL1GasFees to get the fees.
   */
  function _getL1CostWeiForCalldata() internal view returns (uint256) {
    return ARBGAS.getCurrentTxL1GasFees();
  }

  /**
   * @notice Returns the gas cost in wei of calldataSizeBytes of calldata being posted to L1
   */
  function _getL1CostWeiForCalldataSize(uint256 calldataSizeBytes) internal view returns (uint256) {
    (, uint256 l1PricePerByte, , , , ) = ARBGAS.getPricesInWei();
    // see https://developer.arbitrum.io/devs-how-tos/how-to-estimate-gas#where-do-we-get-all-this-information-from
    // for the justification behind the 140 number.
    return l1PricePerByte * (calldataSizeBytes + 140);
  }
}
