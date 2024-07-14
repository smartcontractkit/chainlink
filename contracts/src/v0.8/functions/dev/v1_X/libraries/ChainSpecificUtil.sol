// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbGasInfo} from "../../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {GasPriceOracle} from "../../../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";

/// @dev A library that abstracts out opcodes that behave differently across chains.
/// @dev The methods below return values that are pertinent to the given chain.
library ChainSpecificUtil {
  // ------------ Start Arbitrum Constants ------------
  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);
  /// @dev ARB_DATA_PADDING_SIZE is the max size of the "static" data on Arbitrum for the transaction which refers to the tx data that is not the calldata (signature, etc.)
  /// @dev reference: https://docs.arbitrum.io/build-decentralized-apps/how-to-estimate-gas#where-do-we-get-all-this-information-from
  uint256 private constant ARB_DATA_PADDING_SIZE = 140;

  uint256 private constant ARB_MAINNET_CHAIN_ID = 42161;
  uint256 private constant ARB_GOERLI_TESTNET_CHAIN_ID = 421613;
  uint256 private constant ARB_SEPOLIA_TESTNET_CHAIN_ID = 421614;

  // ------------ End Arbitrum Constants ------------

  // ------------ Start Optimism Constants ------------
  /// @dev GAS_PRICE_ORACLE_ADDR is the address of the GasPriceOracle precompile on Optimism.
  address private constant GAS_PRICE_ORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  GasPriceOracle private constant GAS_PRICE_ORACLE = GasPriceOracle(GAS_PRICE_ORACLE_ADDR);

  uint256 private constant OP_MAINNET_CHAIN_ID = 10;
  uint256 private constant OP_GOERLI_CHAIN_ID = 420;
  uint256 private constant OP_SEPOLIA_CHAIN_ID = 11155420;

  /// @dev Base is a OP stack based rollup and follows the same L1 pricing logic as Optimism.
  uint256 private constant BASE_MAINNET_CHAIN_ID = 8453;
  uint256 private constant BASE_GOERLI_CHAIN_ID = 84531;
  uint256 private constant BASE_SEPOLIA_CHAIN_ID = 84532;

  // ------------ End Optimism Constants ------------

  /// @notice Returns the upper limit estimate of the L1 fees in wei that will be paid for L2 chains
  /// @notice based on the size of the transaction data and the current gas conditions.
  /// @notice This is an "upper limit" as it assumes the transaction data is uncompressed when posted on L1.
  function _getL1FeeUpperLimit(uint256 calldataSizeBytes) internal view returns (uint256 l1FeeWei) {
    uint256 chainid = block.chainid;
    if (_isArbitrumChainId(chainid)) {
      // https://docs.arbitrum.io/build-decentralized-apps/how-to-estimate-gas#where-do-we-get-all-this-information-from
      (, uint256 l1PricePerByte, , , , ) = ARBGAS.getPricesInWei();
      return l1PricePerByte * (calldataSizeBytes + ARB_DATA_PADDING_SIZE);
    } else if (_isOptimismChainId(chainid)) {
      return GAS_PRICE_ORACLE.getL1FeeUpperBound(calldataSizeBytes);
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
