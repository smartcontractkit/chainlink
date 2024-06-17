// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/libraries/Constants.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

/// @title Constants
/// @notice Constants is a library for storing constants. Simple! Don't put everything in here, just
///         the stuff used in multiple contracts. Constants that only apply to a single contract
///         should be defined in that contract instead.
library Constants {
  /// @notice Special address to be used as the tx origin for gas estimation calls in the
  ///         OptimismPortal and CrossDomainMessenger calls. You only need to use this address if
  ///         the minimum gas limit specified by the user is not actually enough to execute the
  ///         given message and you're attempting to estimate the actual necessary gas limit. We
  ///         use address(1) because it's the ecrecover precompile and therefore guaranteed to
  ///         never have any code on any EVM chain.
  address internal constant ESTIMATION_ADDRESS = address(1);

  /// @notice Value used for the L2 sender storage slot in both the OptimismPortal and the
  ///         CrossDomainMessenger contracts before an actual sender is set. This value is
  ///         non-zero to reduce the gas cost of message passing transactions.
  address internal constant DEFAULT_L2_SENDER = 0x000000000000000000000000000000000000dEaD;

  /// @notice The storage slot that holds the address of a proxy implementation.
  /// @dev `bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)`
  bytes32 internal constant PROXY_IMPLEMENTATION_ADDRESS =
    0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

  /// @notice The storage slot that holds the address of the owner.
  /// @dev `bytes32(uint256(keccak256('eip1967.proxy.admin')) - 1)`
  bytes32 internal constant PROXY_OWNER_ADDRESS = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

  /// @notice The address that represents ether when dealing with ERC20 token addresses.
  address internal constant ETHER = 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE;

  /// @notice The address that represents the system caller responsible for L1 attributes
  ///         transactions.
  address internal constant DEPOSITOR_ACCOUNT = 0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAd0001;
}
