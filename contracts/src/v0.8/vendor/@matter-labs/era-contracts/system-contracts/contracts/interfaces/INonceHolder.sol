// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

/**
 * @author Matter Labs
 * @dev Interface of the nonce holder contract -- a contract used by the system to ensure
 * that there is always a unique identifier for a transaction with a particular account (we call it nonce).
 * In other words, the pair of (address, nonce) should always be unique.
 * @dev Custom accounts should use methods of this contract to store nonces or other possible unique identifiers
 * for the transaction.
 */
interface INonceHolder {
    event ValueSetUnderNonce(address indexed accountAddress, uint256 indexed key, uint256 value);

    /// @dev Returns the current minimal nonce for account.
    function getMinNonce(address _address) external view returns (uint256);

    /// @dev Returns the raw version of the current minimal nonce
    /// (equal to minNonce + 2^128 * deployment nonce).
    function getRawNonce(address _address) external view returns (uint256);

    /// @dev Increases the minimal nonce for the msg.sender.
    function increaseMinNonce(uint256 _value) external returns (uint256);

    /// @dev Sets the nonce value `key` as used.
    function setValueUnderNonce(uint256 _key, uint256 _value) external;

    /// @dev Gets the value stored inside a custom nonce.
    function getValueUnderNonce(uint256 _key) external view returns (uint256);

    /// @dev A convenience method to increment the minimal nonce if it is equal
    /// to the `_expectedNonce`.
    function incrementMinNonceIfEquals(uint256 _expectedNonce) external;

    /// @dev Returns the deployment nonce for the accounts used for CREATE opcode.
    function getDeploymentNonce(address _address) external view returns (uint256);

    /// @dev Increments the deployment nonce for the account and returns the previous one.
    function incrementDeploymentNonce(address _address) external returns (uint256);

    /// @dev Determines whether a certain nonce has been already used for an account.
    function validateNonceUsage(address _address, uint256 _key, bool _shouldBeUsed) external view;

    /// @dev Returns whether a nonce has been used for an account.
    function isNonceUsed(address _address, uint256 _nonce) external view returns (bool);
}
