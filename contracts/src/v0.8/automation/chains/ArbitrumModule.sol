// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {ArbGasInfo} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {ChainModuleBase} from "./ChainModuleBase.sol";

contract ArbitrumModule is ChainModuleBase {
  /// @dev ARB_SYS_ADDR is the address of the ArbSys precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbSys.sol#L10
  address private constant ARB_SYS_ADDR = 0x0000000000000000000000000000000000000064;
  ArbSys private constant ARB_SYS = ArbSys(ARB_SYS_ADDR);

  /// @dev ARB_GAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARB_GAS_ADDR = 0x000000000000000000000000000000000000006C;
  ArbGasInfo private constant ARB_GAS = ArbGasInfo(ARB_GAS_ADDR);

  uint256 private constant FIXED_GAS_OVERHEAD = 5000;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 0;

  function blockHash(uint256 n) external view override returns (bytes32) {
    uint256 blockNum = ARB_SYS.arbBlockNumber();
    if (n >= blockNum || blockNum - n > 256) {
      return "";
    }
    return ARB_SYS.arbBlockHash(n);
  }

  function blockNumber() external view override returns (uint256) {
    return ARB_SYS.arbBlockNumber();
  }

  function getCurrentL1Fee() external view override returns (uint256) {
    return ARB_GAS.getCurrentTxL1GasFees();
  }

  function getMaxL1Fee(uint256 dataSize) external view override returns (uint256) {
    (, uint256 perL1CalldataByte, , , , ) = ARB_GAS.getPricesInWei();
    return perL1CalldataByte * dataSize;
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
