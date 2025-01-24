// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {IAccountCodeStorage} from "./interfaces/IAccountCodeStorage.sol";
import {Utils} from "./libraries/Utils.sol";
import {DEPLOYER_SYSTEM_CONTRACT, NONCE_HOLDER_SYSTEM_CONTRACT, CURRENT_MAX_PRECOMPILE_ADDRESS} from "./Constants.sol";
import {Unauthorized, InvalidCodeHash, CodeHashReason} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice The storage of this contract serves as a mapping for the code hashes of the 32-byte account addresses.
 * @dev Code hash is not strictly a hash, it's a structure where the first byte denotes the version of the hash,
 * the second byte denotes whether the contract is constructed, and the next two bytes denote the length in 32-byte words.
 * And then the next 28 bytes are the truncated hash.
 * @dev In this version of ZKsync, the first byte of the hash MUST be 1.
 * @dev The length of each bytecode MUST be odd.  It's internal code format requirements, due to padding of SHA256 function.
 * @dev It is also assumed that all the bytecode hashes are *known*, i.e. the full bytecodes
 * were published on L1 as calldata. This contract trusts the ContractDeployer and the KnownCodesStorage
 * system contracts to enforce the invariants mentioned above.
 */
contract AccountCodeStorage is IAccountCodeStorage {
    bytes32 private constant EMPTY_STRING_KECCAK = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;

    modifier onlyDeployer() {
        if (msg.sender != address(DEPLOYER_SYSTEM_CONTRACT)) {
            revert Unauthorized(msg.sender);
        }
        _;
    }

    /// @notice Stores the bytecodeHash of constructing contract.
    /// @param _address The address of the account to set the codehash to.
    /// @param _hash The new bytecode hash of the constructing account.
    /// @dev This method trusts the ContractDeployer to make sure that the bytecode is known and well-formed,
    /// but checks whether the bytecode hash corresponds to the constructing smart contract.
    function storeAccountConstructingCodeHash(address _address, bytes32 _hash) external override onlyDeployer {
        // Check that code hash corresponds to the deploying smart contract
        if (!Utils.isContractConstructing(_hash)) {
            revert InvalidCodeHash(CodeHashReason.NotContractOnConstructor);
        }
        _storeCodeHash(_address, _hash);
    }

    /// @notice Stores the bytecodeHash of constructed contract.
    /// @param _address The address of the account to set the codehash to.
    /// @param _hash The new bytecode hash of the constructed account.
    /// @dev This method trusts the ContractDeployer to make sure that the bytecode is known and well-formed,
    /// but checks whether the bytecode hash corresponds to the constructed smart contract.
    function storeAccountConstructedCodeHash(address _address, bytes32 _hash) external override onlyDeployer {
        // Check that code hash corresponds to the deploying smart contract
        if (!Utils.isContractConstructed(_hash)) {
            revert InvalidCodeHash(CodeHashReason.NotConstructedContract);
        }
        _storeCodeHash(_address, _hash);
    }

    /// @notice Marks the account bytecodeHash as constructed.
    /// @param _address The address of the account to mark as constructed
    function markAccountCodeHashAsConstructed(address _address) external override onlyDeployer {
        bytes32 codeHash = getRawCodeHash(_address);

        if (!Utils.isContractConstructing(codeHash)) {
            revert InvalidCodeHash(CodeHashReason.NotContractOnConstructor);
        }

        // Get the bytecode hash with "isConstructor" flag equal to false
        bytes32 constructedBytecodeHash = Utils.constructedBytecodeHash(codeHash);

        _storeCodeHash(_address, constructedBytecodeHash);
    }

    /// @dev Store the codehash of the account without any checks.
    /// @param _address The address of the account to set the codehash to.
    /// @param _hash The new account bytecode hash.
    function _storeCodeHash(address _address, bytes32 _hash) internal {
        uint256 addressAsKey = uint256(uint160(_address));
        assembly {
            sstore(addressAsKey, _hash)
        }
    }

    /// @notice Get the codehash stored for an address.
    /// @param _address The address of the account of which the codehash to return
    /// @return codeHash The codehash stored for this account.
    function getRawCodeHash(address _address) public view override returns (bytes32 codeHash) {
        uint256 addressAsKey = uint256(uint160(_address));

        assembly {
            codeHash := sload(addressAsKey)
        }
    }

    /// @notice Simulate the behavior of the `extcodehash` EVM opcode.
    /// @param _input The 256-bit account address.
    /// @return codeHash - hash of the bytecode according to the EIP-1052 specification.
    function getCodeHash(uint256 _input) external view override returns (bytes32) {
        // We consider the account bytecode hash of the last 20 bytes of the input, because
        // according to the spec "If EXTCODEHASH of A is X, then EXTCODEHASH of A + 2**160 is X".
        address account = address(uint160(_input));
        if (uint160(account) <= CURRENT_MAX_PRECOMPILE_ADDRESS) {
            return EMPTY_STRING_KECCAK;
        }

        bytes32 codeHash = getRawCodeHash(account);

        // The code hash is equal to the `keccak256("")` if the account is an EOA with at least one transaction.
        // Otherwise, the account is either deployed smart contract or an empty account,
        // for both cases the code hash is equal to the raw code hash.
        if (codeHash == 0x00 && NONCE_HOLDER_SYSTEM_CONTRACT.getRawNonce(account) > 0) {
            codeHash = EMPTY_STRING_KECCAK;
        }
        // The contract is still on the constructor, which means it is not deployed yet,
        // so set `keccak256("")` as a code hash. The EVM has the same behavior.
        else if (Utils.isContractConstructing(codeHash)) {
            codeHash = EMPTY_STRING_KECCAK;
        }

        return codeHash;
    }

    /// @notice Simulate the behavior of the `extcodesize` EVM opcode.
    /// @param _input The 256-bit account address.
    /// @return codeSize - the size of the deployed smart contract in bytes.
    function getCodeSize(uint256 _input) external view override returns (uint256 codeSize) {
        // We consider the account bytecode size of the last 20 bytes of the input, because
        // according to the spec "If EXTCODESIZE of A is X, then EXTCODESIZE of A + 2**160 is X".
        address account = address(uint160(_input));
        bytes32 codeHash = getRawCodeHash(account);

        // If the contract is a default account or is on constructor the code size is zero,
        // otherwise extract the proper value for it from the bytecode hash.
        // NOTE: zero address and precompiles are a special case, they are contracts, but we
        // want to preserve EVM invariants (see EIP-1052 specification). That's why we automatically
        // return `0` length in the following cases:
        // - `codehash(0) == 0`
        // - `account` is a precompile.
        // - `account` is currently being constructed
        if (
            uint160(account) > CURRENT_MAX_PRECOMPILE_ADDRESS &&
            codeHash != 0x00 &&
            !Utils.isContractConstructing(codeHash)
        ) {
            codeSize = Utils.bytecodeLenInBytes(codeHash);
        }
    }
}
