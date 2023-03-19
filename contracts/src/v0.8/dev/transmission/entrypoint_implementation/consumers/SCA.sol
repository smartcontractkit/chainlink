import "../interfaces/IAccount.sol";
import "./SCALibrary.sol";

//SPDX-License-Identifier: Unlicense
pragma solidity 0.8.15;

// Smart Contract Account, a contract deployed for a single user and that allows
// them to invoke meta-transactions.
contract SCA is IAccount {
  uint256 s_nonce;
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
    uint256 missingAccountFunds
  ) external returns (uint256 validationData) {
    // Construct hash, consisting of end-transaction details, domain seperators, nonce, and chain ID.
    bytes32 txHash = keccak256(abi.encode(SCALibrary.TYPEHASH, userOp.callData, s_owner, s_nonce, block.chainid));
    bytes32 fullHash = keccak256(abi.encodePacked(bytes1(0x19), bytes1(0x01), SCALibrary.DOMAIN_SEPARATOR, txHash));

    // Verify signature on hash.
    bytes memory signature = userOp.signature;
    bytes32 r;
    bytes32 s;
    assembly {
      r := mload(add(signature, 0x20))
      s := mload(add(signature, 0x40))
    }
    uint8 v = uint8(signature[64]);
    if (ecrecover(fullHash, v + 27, r, s) != s_owner) {
      revert InvalidSignature(fullHash, s_owner);
    }

    s_nonce++;
    return 0; // TOOD: add validationData for billing.
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
