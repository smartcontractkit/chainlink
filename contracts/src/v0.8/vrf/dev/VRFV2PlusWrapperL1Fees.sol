// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbitrumL1Fees} from "./ArbitrumL1Fees.sol";
import {OptimismL1Fees} from "./OptimismL1Fees.sol";
import {GasPriceOracle as OVM_GasPriceOracle} from "../../vendor/@eth-optimism/contracts-bedrock/v0.17.3/src/L2/GasPriceOracle.sol";

/// @dev This abstract contract combines all L2 specific operations necessary for L1 gas fee computation.
/// @dev It hides away the L1 gas fee computation from the VRFV2PlusWrapper contract.
// solhint-disable-next-line contract-name-camelcase
abstract contract VRFV2PlusWrapperL1Fees is ArbitrumL1Fees, OptimismL1Fees {
  uint256 private constant ARB_MAINNET_CHAIN_ID = 42161;
  uint256 private constant ARB_GOERLI_TESTNET_CHAIN_ID = 421613;
  uint256 private constant ARB_SEPOLIA_TESTNET_CHAIN_ID = 421614;

  uint256 private constant OP_MAINNET_CHAIN_ID = 10;
  uint256 private constant OP_GOERLI_CHAIN_ID = 420;
  uint256 private constant OP_SEPOLIA_CHAIN_ID = 11155420;

  /// @dev Base is a OP stack based rollup and follows the same L1 pricing logic as Optimism.
  uint256 private constant BASE_MAINNET_CHAIN_ID = 8453;
  uint256 private constant BASE_GOERLI_CHAIN_ID = 84531;
  uint256 private constant BASE_SEPOLIA_CHAIN_ID = 84532;

  /// @dev this is the size of a VRF v2 fulfillment's calldata abi-encoded in bytes.
  /// @dev proofSize = 13 words = 13 * 256 = 3328 bits
  /// @dev commitmentSize = 10 words = 10 * 256 = 2560 bits
  /// @dev onlyPremiumParameterSize = 256 bits
  /// @dev dataSize = proofSize + commitmentSize + onlyPremiumParameterSize = 6144 bits
  /// @dev function selector = 32 bits
  /// @dev total data size = 6144 bits + 32 bits = 6176 bits = 772 bytes
  uint32 public s_fulfillmentTxSizeBytes = 772;

  error UnsupportedChainId(uint256 chainId);
  error UnsupportedFunction();

  event FulfillmentTxSizeSet(uint32 size);

  /**
   * @notice setFulfillmentTxSize sets the size of the fulfillment transaction in bytes.
   * @param _size is the size of the fulfillment transaction in bytes.
   */
  function setFulfillmentTxSize(uint32 _size) external onlyOwner {
    s_fulfillmentTxSizeBytes = _size;

    emit FulfillmentTxSizeSet(_size);
  }

  /**
   * @notice Return true if and only if the provided chain ID is an Arbitrum chain ID.
   */
  function _isArbitrumChainId(uint256 chainId) internal pure returns (bool) {
    return
      chainId == ARB_MAINNET_CHAIN_ID ||
      chainId == ARB_GOERLI_TESTNET_CHAIN_ID ||
      chainId == ARB_SEPOLIA_TESTNET_CHAIN_ID;
  }

  /**
   * @notice Return true if and only if the provided chain ID is an Optimism chain ID.
   * @notice Note that optimism chain id's are also OP stack chain id's (e.g. Base).
   */
  function _isOptimismChainId(uint256 chainId) internal pure returns (bool) {
    return
      chainId == OP_MAINNET_CHAIN_ID ||
      chainId == OP_GOERLI_CHAIN_ID ||
      chainId == OP_SEPOLIA_CHAIN_ID ||
      chainId == BASE_MAINNET_CHAIN_ID ||
      chainId == BASE_GOERLI_CHAIN_ID ||
      chainId == BASE_SEPOLIA_CHAIN_ID;
  }

  /**
   * @notice This is only necessary to avoid compiler issue and to prevent anyone
   * @notice from accidentially using this function.
   */
  function _getL1CostWeiForCalldataSize(
    uint256 calldataSizeBytes
  ) internal view override(ArbitrumL1Fees, OptimismL1Fees) returns (uint256) {
    revert UnsupportedFunction();
  }

  /**
   * @notice Returns estimated L1 gas fee cost for fulfillment calldata payload once
   * @notice the request has been made through VRFV2PlusWrapper (direct funding model).
   */
  function _getL1CostWei() internal view returns (uint256) {
    uint256 chainid = block.chainid;
    if (_isArbitrumChainId(chainid)) {
      return ArbitrumL1Fees._getL1CostWeiForCalldataSize(s_fulfillmentTxSizeBytes);
    } else if (_isOptimismChainId(chainid)) {
      return OptimismL1Fees._getL1CostWeiForCalldataSize(s_fulfillmentTxSizeBytes);
    }
    revert UnsupportedChainId(chainid);
  }
}
