import "../interfaces/IAccount.sol";

//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.0;

// Smart Contract Account, a contract deployed for a single user and that allows
// them to invoke meta-transactions.
contract SCA is IAccount {
    uint256 s_nonce;
    address public immutable s_owner;

    // TODO: move this logic to the Create2Factory, ditch the hardcoding.
    address constant ENTRY_POINT = 0x0576a174D229E3cFA37253523E645A78A0C91B57;

    // Assign the owner of this contract upon deployment.
    constructor(address owner) {
        s_owner = owner;
    }

    /// @dev Sanity check function.
    /// @dev TODO: remove.
    function isSCA() external pure returns (bool) {
        return true;
    }

    /// @dev Validates the user operation via a signature check.
    /// @dev TODO: Add a domain separator to the signature.
    function validateUserOp(
        UserOperation calldata userOp,
        bytes32 userOpHash,
        uint256 missingAccountFunds
    ) external returns (uint256 validationData) {
        s_nonce++;
        bytes32 h = keccak256(
            abi.encode(
                userOp.callData,
                s_owner,
                s_nonce,
                block.chainid
            )
        );
        bytes32 r = bytesToBytes32(userOp.signature, 0);
        bytes32 s = bytesToBytes32(userOp.signature, 32);
        uint8 v = uint8(userOp.signature[64]);
        require(ecrecover(h, v + 27, r, s) == s_owner, "Invalid signature.");

        return 0;
    }

    /// @dev Execute a transaction on behalf of the owner. This function can only
    /// @dev be called by the EntryPoint contract, and assumes that `validateUserOp` has succeeded.
    function executeTransaction(address to, bytes calldata data) external {
        require(msg.sender == ENTRY_POINT, "not authorized");
        /* (bool success, bytes memory returnData) = */
        to.call(data);
    }

    function bytesToBytes32(bytes memory b, uint256 offset)
        internal
        pure
        returns (bytes32)
    {
        bytes32 out;

        for (uint256 i = 0; i < 32; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}
