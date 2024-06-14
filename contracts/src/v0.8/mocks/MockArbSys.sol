// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

/**
 * @dev this contract mocks the arbitrum precompiled ArbSys contract
 * https://developer.arbitrum.io/arbos/precompiles#ArbSys
 */
contract MockArbSys {
  function arbBlockNumber() public view returns (uint256) {
    return block.number;
  }

  function arbBlockHash(uint256 arbBlockNum) external view returns (bytes32) {
    return blockhash(arbBlockNum);
  }
}
