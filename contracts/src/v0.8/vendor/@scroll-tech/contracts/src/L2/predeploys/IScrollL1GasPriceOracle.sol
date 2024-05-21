// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

interface IScrollL1GasPriceOracle {
    /**********
     * Events *
     **********/

    /// @notice Emitted when current fee overhead is updated.
    /// @param overhead The current fee overhead updated.
    event OverheadUpdated(uint256 overhead);

    /// @notice Emitted when current fee scalar is updated.
    /// @param scalar The current fee scalar updated.
    event ScalarUpdated(uint256 scalar);

    /// @notice Emitted when current l1 base fee is updated.
    /// @param l1BaseFee The current l1 base fee updated.
    event L1BaseFeeUpdated(uint256 l1BaseFee);

    /*************************
     * Public View Functions *
     *************************/

    /// @notice Return the current l1 fee overhead.
    function overhead() external view returns (uint256);

    /// @notice Return the current l1 fee scalar.
    function scalar() external view returns (uint256);

    /// @notice Return the latest known l1 base fee.
    function l1BaseFee() external view returns (uint256);

    /// @notice Computes the L1 portion of the fee based on the size of the rlp encoded input
    ///         transaction, the current L1 base fee, and the various dynamic parameters.
    /// @param data Unsigned fully RLP-encoded transaction to get the L1 fee for.
    /// @return L1 fee that should be paid for the tx
    function getL1Fee(bytes memory data) external view returns (uint256);

    /// @notice Computes the amount of L1 gas used for a transaction. Adds the overhead which
    ///         represents the per-transaction gas overhead of posting the transaction and state
    ///         roots to L1. Adds 74 bytes of padding to account for the fact that the input does
    ///         not have a signature.
    /// @param data Unsigned fully RLP-encoded transaction to get the L1 gas for.
    /// @return Amount of L1 gas used to publish the transaction.
    function getL1GasUsed(bytes memory data) external view returns (uint256);

    /*****************************
     * Public Mutating Functions *
     *****************************/

    /// @notice Allows whitelisted caller to modify the l1 base fee.
    /// @param _l1BaseFee New l1 base fee.
    function setL1BaseFee(uint256 _l1BaseFee) external;
}