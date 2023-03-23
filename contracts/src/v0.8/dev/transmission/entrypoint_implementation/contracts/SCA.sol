// SPDX-License-Identifier: MIT
import "../../../vendor/entrypoint/interfaces/IAccount.sol";
import "./SCALibrary.sol";
import "../../../vendor/entrypoint/core/Helpers.sol";

/// TODO: decide on a compiler version. Must not be dynamic, and must be > 0.8.12.
pragma solidity 0.8.15; 

// Smart Contract Account, a contract deployed for a single user and that allows
// them to invoke meta-transactions.
contract SCA is IAccount {
  uint256 public s_nonce;
  address public immutable s_owner;
  address public immutable s_entryPoint;

  error NotAuthorized(address sender);
  error TransactionExpired(uint256 deadline, uint256 currentTimestamp);
  error InvalidSignature(bytes32 operationHash, address owner);

  // Assign the owner of this contract upon deployment.
  constructor(address owner, address entryPoint) {
    s_owner = owner;
    s_entryPoint = entryPoint;
  }

  /// @dev Validates the user operation via a signature check.
  function validateUserOp(
    UserOperation calldata userOp,
    bytes32 userOpHash,
    uint256 /* missingAccountFunds - unused in favor of paymaster */
  ) external returns (uint256 validationData) {
    if (userOp.nonce != s_nonce) {
      return _packValidationData(true, 0, 0); // incorrect nonce
    }

    // Verify signature on hash.
    bytes32 fullHash = SCALibrary.getUserOpFullHash(userOpHash);
    bytes memory signature = userOp.signature;
    bytes32 r;
    bytes32 s;
    assembly {
      r := mload(add(signature, 0x20))
      s := mload(add(signature, 0x40))
    }
    uint8 v = uint8(signature[64]);
    if (ecrecover(fullHash, v + 27, r, s) != s_owner) {
      return _packValidationData(true, 0, 0); // signature error
    }

    s_nonce++;
    return _packValidationData(false, 0, 0); // success, with indefinite validity
  }

  /// @dev Execute a transaction on behalf of the owner. This function can only
  /// @dev be called by the EntryPoint contract, and assumes that `validateUserOp` has succeeded.
  function executeTransactionFromEntryPoint(address to, uint256 value, uint256 deadline, bytes calldata data) external {
    if (msg.sender != s_entryPoint) {
      revert NotAuthorized(msg.sender);
    }
    if (block.timestamp > deadline) {
      revert TransactionExpired(deadline, block.timestamp);
    }
    to.call{value: value}(data);
  }
}
