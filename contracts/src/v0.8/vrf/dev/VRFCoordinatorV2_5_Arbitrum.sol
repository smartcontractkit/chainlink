// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {VRFCoordinatorV2_5} from "./VRFCoordinatorV2_5.sol";
import {ArbitrumL1Fees} from "./ArbitrumL1Fees.sol";

/// @dev VRFCoordinatorV2_5_Arbitrum combines VRFCoordinatorV2_5 base contract with
/// @dev Arbitrum specific opcodes and L1 gas fee calculations.
// solhint-disable-next-line contract-name-camelcase
contract VRFCoordinatorV2_5_Arbitrum is VRFCoordinatorV2_5, ArbitrumL1Fees {
  /// @dev ARBSYS_ADDR is the address of the ArbSys precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbSys.sol#L10
  address private constant ARBSYS_ADDR = address(0x0000000000000000000000000000000000000064);
  ArbSys private constant ARBSYS = ArbSys(ARBSYS_ADDR);

  constructor(address blockhashStore) VRFCoordinatorV2_5(blockhashStore) {}

  /**
   * @notice Override getBlockhash from VRFCoordinatorV2_5
   * @notice When on a known Arbitrum chain, it uses ArbSys.arbBlockHash to get the blockhash.
   */
  function _getBlockhash(uint64 blockNumber) internal view override returns (bytes32) {
    uint64 currentBlockNumber = uint64(_getBlockNumber());
    if (blockNumber >= currentBlockNumber || (currentBlockNumber - blockNumber) > 256) {
      return "";
    }
    return ARBSYS.arbBlockHash(blockNumber);
  }

  /**
   * @notice Override getBlockNumber from VRFCoordinatorV2_5
   * @notice When on a known Arbitrum chain, it uses ArbSys.arbBlockNumber to get the block number.
   */
  function _getBlockNumber() internal view override returns (uint256) {
    return ARBSYS.arbBlockNumber();
  }

  /// @notice Override getL1CostWei function from VRFCoordinatorV2_5 to activate Arbitrum getL1Fee computation
  function _getL1CostWei(bytes calldata /* data */) internal view override returns (uint256) {
    return _getL1CostWeiForCalldata();
  }
}
