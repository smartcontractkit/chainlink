// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/f707883038d527cbf1e9f8ea513fe33255deadbc/packages/contracts-bedrock/src/L1/L1StandardBridge.sol
pragma solidity ^0.8.0;

interface IOptimismL1StandardBridge {
  /// @custom:legacy
  /// @notice Deposits some amount of ETH into a target account on L2.
  ///         Note that if ETH is sent to a contract on L2 and the call fails, then that ETH will
  ///         be locked in the L2StandardBridge. ETH may be recoverable if the call can be
  ///         successfully replayed by increasing the amount of gas supplied to the call. If the
  ///         call will fail for any amount of gas, then the ETH will be locked permanently.
  /// @param _to          Address of the recipient on L2.
  /// @param _minGasLimit Minimum gas limit for the deposit message on L2.
  /// @param _extraData   Optional data to forward to L2.
  ///                     Data supplied here will not be used to execute any code on L2 and is
  ///                     only emitted as extra data for the convenience of off-chain tooling.
  function depositETHTo(address _to, uint32 _minGasLimit, bytes calldata _extraData) external payable;
}
