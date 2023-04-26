// SPDX-License-Identifier: MIT
import "../../vendor/entrypoint/interfaces/IAccount.sol";
import "./SCALibrary.sol";
import "../../vendor/entrypoint/core/Helpers.sol";

/// TODO: decide on a compiler version. Must not be dynamic, and must be > 0.8.12.
pragma solidity 0.8.15;

/// @dev Smart Contract Account, a contract deployed for a single user and that allows
/// @dev them to invoke meta-transactions.
/// TODO: Consider making the Smart Contract Account upgradeable.
contract SCA is IAccount {
  uint256 public s_nonce;
  address public immutable i_owner;
  address public immutable i_entryPoint;

  error IncorrectNonce(uint256 currentNonce, uint256 nonceGiven);
  error NotAuthorized(address sender);
  error BadFormatOrOOG();
  error TransactionExpired(uint256 deadline, uint256 currentTimestamp);
  error InvalidSignature(bytes32 operationHash, address owner);

  // Assign the owner of this contract upon deployment.
  constructor(address owner, address entryPoint) {
    i_owner = owner;
    i_entryPoint = entryPoint;
  }

  /// @dev Validates the user operation via a signature check.
  /// TODO: Utilize a "validAfter" for a tx to be only valid _after_ a certain time.
  function validateUserOp(
    UserOperation calldata userOp,
    bytes32 userOpHash,
    uint256 /* missingAccountFunds - unused in favor of paymaster */
  ) external returns (uint256 validationData) {
    if (userOp.nonce != s_nonce) {
      // Revert for non-signature errors.
      revert IncorrectNonce(s_nonce, userOp.nonce);
    }

    // Verify signature on hash.
    bytes32 fullHash = SCALibrary.getUserOpFullHash(userOpHash, address(this));
    bytes memory signature = userOp.signature;
    if (SCALibrary.recoverSignature(signature, fullHash) != i_owner) {
      return _packValidationData(true, 0, 0); // signature error
    }
    s_nonce++;

    // Unpack deadline, return successful signature.
    (, , uint48 deadline, ) = abi.decode(userOp.callData[4:], (address, uint256, uint48, bytes));
    return _packValidationData(false, deadline, 0);
  }

  /// @dev Execute a transaction on behalf of the owner. This function can only
  /// @dev be called by the EntryPoint contract, and assumes that `validateUserOp` has succeeded.
  function executeTransactionFromEntryPoint(
    address to,
    uint256 value,
    uint48 deadline,
    bytes calldata data
  ) external {
    if (msg.sender != i_entryPoint) {
      revert NotAuthorized(msg.sender);
    }
    if (deadline != 0 && block.timestamp > deadline) {
      revert TransactionExpired(deadline, block.timestamp);
    }

    // Execute transaction. Bubble up an error if found.
    (bool success, bytes memory returnData) = to.call{value: value}(data);
    if (!success) {
      if (returnData.length == 0) revert BadFormatOrOOG();
      assembly {
        revert(add(32, returnData), mload(returnData))
      }
    }
  }
}
