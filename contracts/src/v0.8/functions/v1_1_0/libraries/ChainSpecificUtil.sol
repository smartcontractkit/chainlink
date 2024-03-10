// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbGasInfo} from "../../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import {GasPriceOracle} from "../../../vendor/@eth-optimism/contracts-bedrock/v0.16.2/src/L2/GasPriceOracle.sol";

/// @dev A library that abstracts out opcodes that behave differently across chains.
/// @dev The methods below return values that are pertinent to the given chain.
library ChainSpecificUtil {
  // ------------ Start Arbitrum Constants ------------

  /// @dev ARBGAS_ADDR is the address of the ArbGasInfo precompile on Arbitrum.
  /// @dev reference: https://github.com/OffchainLabs/nitro/blob/v2.0.14/contracts/src/precompiles/ArbGasInfo.sol#L10
  address private constant ARBGAS_ADDR = address(0x000000000000000000000000000000000000006C);
  ArbGasInfo private constant ARBGAS = ArbGasInfo(ARBGAS_ADDR);

  uint256 private constant ARB_MAINNET_CHAIN_ID = 42161;
  uint256 private constant ARB_GOERLI_TESTNET_CHAIN_ID = 421613;
  uint256 private constant ARB_SEPOLIA_TESTNET_CHAIN_ID = 421614;

  // ------------ End Arbitrum Constants ------------

  // ------------ Start Optimism Constants ------------
  /// @dev L1_FEE_DATA_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes internal constant L1_FEE_DATA_PADDING =
    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  /// @dev OVM_GASPRICEORACLE_ADDR is the address of the GasPriceOracle precompile on Optimism.
  /// @dev reference: https://community.optimism.io/docs/developers/build/transaction-fees/#estimating-the-l1-data-fee
  address private constant OVM_GASPRICEORACLE_ADDR = address(0x420000000000000000000000000000000000000F);
  GasPriceOracle private constant OVM_GASPRICEORACLE = GasPriceOracle(OVM_GASPRICEORACLE_ADDR);

  uint256 private constant OP_MAINNET_CHAIN_ID = 10;
  uint256 private constant OP_GOERLI_CHAIN_ID = 420;
  uint256 private constant OP_SEPOLIA_CHAIN_ID = 11155420;

  /// @dev Base is a OP stack based rollup and follows the same L1 pricing logic as Optimism.
  uint256 private constant BASE_MAINNET_CHAIN_ID = 8453;
  uint256 private constant BASE_GOERLI_CHAIN_ID = 84531;
  uint256 private constant BASE_SEPOLIA_CHAIN_ID = 84532;

  // ------------ End Optimism Constants ------------

  /// @notice Returns the L1 fees in wei that will be paid for the current transaction, given any calldata
  /// @notice for the current transaction.
  /// @notice When on a known Arbitrum chain, it uses ArbGas.getCurrentTxL1GasFees to get the fees.
  /// @notice On Arbitrum, the provided calldata is not used to calculate the fees.
  /// @notice On Optimism, the provided calldata is passed to the GasPriceOracle predeploy
  /// @notice and getL1Fee is called to get the fees.
  function _getCurrentTxL1GasFees(bytes memory txCallData) internal view returns (uint256 l1FeeWei) {
    uint256 chainid = block.chainid;
    if (_isArbitrumChainId(chainid)) {
      return ARBGAS.getCurrentTxL1GasFees();
    } else if (_isOptimismChainId(chainid)) {
      return OVM_GASPRICEORACLE.getL1Fee(bytes.concat(txCallData, L1_FEE_DATA_PADDING));
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
