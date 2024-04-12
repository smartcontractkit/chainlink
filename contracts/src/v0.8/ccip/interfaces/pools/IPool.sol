// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

// Shared public interface for multiple pool types.
// Each pool type handles a different child token model (lock/unlock, mint/burn.)
interface IPool {
  struct SourceTokenData {
    bytes sourcePoolAddress;
    bytes destPoolAddress;
    bytes extraData;
  }

  /// @notice Lock tokens into the pool or burn the tokens.
  /// @param originalSender Original sender of the tokens.
  /// @param receiver Receiver of the tokens on destination chain.
  /// @param amount Amount to lock or burn.
  /// @param remoteChainSelector Destination chain Id.
  /// @param extraArgs Additional data passed in by sender for lockOrBurn processing
  /// in custom pools on source chain.
  /// @return poolReturnData Versioned, encoded data fields for the processing of tokens
  /// on the destination chain.
  function lockOrBurn(
    address originalSender,
    bytes calldata receiver,
    uint256 amount,
    uint64 remoteChainSelector,
    bytes calldata extraArgs
  ) external returns (bytes memory poolReturnData);

  /// @notice Releases or mints tokens to the receiver address.
  /// @param originalSender Original sender of the tokens.
  /// @param receiver Receiver of the tokens.
  /// @param amount Amount to release or mint.
  /// @param remoteChainSelector Source chain Id.
  /// @param sourceTokenData The source and dest pool addresses, as well as any additional data
  /// from calling lockOrBurn on the source chain.
  /// @param offchainTokenData Additional data supplied offchain for releaseOrMint processing in
  /// custom pools on dest chain. This could be an attestation that was retrieved through a
  /// third party API.
  /// @dev offchainData can come from any untrusted source.
  /// @return localToken The address of the local token.
  function releaseOrMint(
    bytes memory originalSender,
    address receiver,
    uint256 amount,
    uint64 remoteChainSelector,
    SourceTokenData memory sourceTokenData,
    bytes memory offchainTokenData
  ) external returns (address localToken);

  /// @notice Gets the IERC20 token that this pool can lock or burn.
  /// @return token The IERC20 token representation.
  function getToken() external view returns (IERC20 token);

  /// @notice Gets the pool address on the remote chain.
  /// @param remoteChainSelector Destination chain selector.
  function getRemotePool(uint64 remoteChainSelector) external view returns (bytes memory);
}
