// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbGasInfo} from "../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {L1Block} from "../../../vendor/@eth-optimism/contracts-bedrock/v0.17.1/src/L2/L1Block.sol";

/// @dev A library that abstracts out opcodes that behave differently across chains.
/// @dev The methods below return values that are pertinent to the given chain.
library ChainSpecificUtil {
  // ------------ Start Arbitrum Constants ------------
  /// @dev ARB_L1_FEE_DATA_PADDING_SIZE is the L1 data padding for Optimism
  uint256 private const ARB_L1_FEE_DATA_PADDING_SIZE = 140;
  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  uint256 private constant ARB_MAINNET_CHAIN_ID = 42161;
  uint256 private constant ARB_GOERLI_TESTNET_CHAIN_ID = 421613;
  uint256 private constant ARB_SEPOLIA_TESTNET_CHAIN_ID = 421614;
  // ------------ End Arbitrum Constants ------------

  // ------------ Start Optimism Constants ------------
  /// @dev OP_L1_FEE_DATA_PADDING_SIZE is the L1 data padding for Optimism
  uint256 private const OP_L1_FEE_DATA_PADDING_SIZE = 35;
  /// @dev L1BLOCK_ADDR is the address of the L1Block precompile on Optimism.
  address private constant L1BLOCK_ADDR = address(0x4200000000000000000000000000000000000015);
  L1Block private constant L1BLOCK = L1Block(L1BLOCK_ADDR);

  uint256 private constant OP_MAINNET_CHAIN_ID = 10;
  uint256 private constant OP_GOERLI_CHAIN_ID = 420;
  uint256 private constant OP_SEPOLIA_CHAIN_ID = 11155420;

  /// @dev Base is a OP stack based rollup and follows the same L1 pricing logic as Optimism.
  uint256 private constant BASE_MAINNET_CHAIN_ID = 8453;
  uint256 private constant BASE_GOERLI_CHAIN_ID = 84531;
  uint256 private constant BASE_SEPOLIA_CHAIN_ID = 84532;
  // ------------ End Optimism Constants ------------

  /// @notice Returns the L1 fees in wei that will be paid for L2 chains
  /// @notice based on the size of the transaction data and the current gas conditions.
  function _getL1FeeUpperLimit(uint256 dataSizeBytes) internal view returns (uint256 l1FeeWei) {
    uint256 chainid = block.chainid;
    if (_isArbitrumChainId(chainid)) {
      // https://docs.arbitrum.io/build-decentralized-apps/how-to-estimate-gas#where-do-we-get-all-this-information-from
      (, uint256 l1PricePerByte, , , , ) = ARBGAS.getPricesInWei();
      return l1PricePerByte * (dataSizeBytes + ARB_L1_FEE_DATA_PADDING_SIZE);
    } else if (_isOptimismChainId(chainid)) {
      // https://docs.optimism.io/stack/transactions/fees#ecotone
      // note we conservatively assume all non-zero bytes: tx_compressed_size = tx_data_size_bytes
      uint256 l1BaseFeeWei = L1BLOCK.baseFee();
      uint256 l1BaseFeeScalar = L1BLOCK.baseFeeScalar();
      uint256 l1BlobBaseFeeWei = L1BLOCK.blobBaseFee();
      uint256 l1BlobBaseFeeScalar = L1BLOCK.blobBaseFeeScalar();
      uint256 weightedGasPrice = 16 * l1BaseFeeScalar * l1BaseFee + l1BlobBaseFeeScalar * l1BlobBaseFeeWei;
      return weightedGasPrice * dataSizeBytes;
    }
    return 0;
  }

  /// @notice Return true if and only if the provided chain ID is an Arbitrum chain ID.
  function _isArbitrumChainId(uint256 chainId) internal pure returns (bool) {
    return
      chainId == ARB_MAINNET_CHAIN_ID ||
      chainId == ARB_GOERLI_TESTNET_CHAIN_ID ||
      chainId == ARB_SEPOLIA_TESTNET_CHAIN_ID;
  }

  /// @notice Return true if and only if the provided chain ID is an Optimism (or Base) chain ID.
  /// @notice Note that optimism chain id's are also OP stack chain id's.
  function _isOptimismChainId(uint256 chainId) internal pure returns (bool) {
    return
      chainId == OP_MAINNET_CHAIN_ID ||
      chainId == OP_GOERLI_CHAIN_ID ||
      chainId == OP_SEPOLIA_CHAIN_ID ||
      chainId == BASE_MAINNET_CHAIN_ID ||
      chainId == BASE_GOERLI_CHAIN_ID ||
      chainId == BASE_SEPOLIA_CHAIN_ID;
  }
}
