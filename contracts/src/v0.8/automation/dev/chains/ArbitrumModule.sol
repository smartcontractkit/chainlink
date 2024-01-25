// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {ArbSys} from "./../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {ArbGasInfo} from "./../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import "../interfaces/v2_2/IChainSpecific.sol";

contract ArbitrumModule is IChainSpecific {
  /// @dev ARBSYS_ADDR is the address of the ArbSys precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbSys.sol#L10
  address private constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);
  ArbSys private constant ARBSYS = ArbSys(ARBSYS_ADDR);

  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  function _blockHash(uint256 blockNumber) external view returns (bytes32) {
    return ARBSYS.arbBlockHash(blockNumber);
  }

  function _blockNumber() external view returns (uint256) {
    return ARBSYS.arbBlockNumber();
  }

  function _getL1FeeForTransaction(bytes calldata) external view returns (uint256) {
    return ARBGAS.getCurrentTxL1GasFees();
  }

  function _getL1FeeForSimulation(bytes calldata txCallData) external view returns (uint256) {
    // fee is 4 per 0 byte, 16 per non-zero byte - we assume all non-zero and
    // max data size to calculate max payment
    (, uint256 perL1CalldataUnit, , , , ) = ARBGAS.getPricesInWei();
    l1CostWei = perL1CalldataUnit * s_storage.maxPerformDataSize * 16;
    return 0;
  }
}
