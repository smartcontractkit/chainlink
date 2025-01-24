// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {IImmutableSimulator, ImmutableData} from "./interfaces/IImmutableSimulator.sol";
import {DEPLOYER_SYSTEM_CONTRACT} from "./Constants.sol";
import {Unauthorized} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice System smart contract that simulates the behavior of immutable variables in Solidity.
 * @dev The contract stores the immutable variables created during deployment by other contracts on his storage.
 * @dev This simulator is needed so that smart contracts with the same Solidity code but different
 * constructor parameters have the same bytecode.
 * @dev The users are not expected to call this contract directly, only indirectly via the compiler simulations
 * for the immutable variables in Solidity.
 */
contract ImmutableSimulator is IImmutableSimulator {
    /// @dev mapping (contract address) => (index of immutable variable) => value
    /// @notice that address uses `uint256` type to leave the option to introduce 32-byte address space in future.
    mapping(uint256 contractAddress => mapping(uint256 index => bytes32 value)) internal immutableDataStorage;

    /// @notice Method that returns the immutable with a certain index for a user.
    /// @param _dest The address which the immutable belongs to.
    /// @param _index The index of the immutable.
    /// @return The value of the immutables.
    function getImmutable(address _dest, uint256 _index) external view override returns (bytes32) {
        return immutableDataStorage[uint256(uint160(_dest))][_index];
    }

    /// @notice Method used by the contract deployer to store the immutables for an account
    /// @param _dest The address which to store the immutables for.
    /// @param _immutables The list of the immutables.
    function setImmutables(address _dest, ImmutableData[] calldata _immutables) external override {
        if (msg.sender != address(DEPLOYER_SYSTEM_CONTRACT)) {
            revert Unauthorized(msg.sender);
        }
        unchecked {
            uint256 immutablesLength = _immutables.length;
            for (uint256 i = 0; i < immutablesLength; ++i) {
                uint256 index = _immutables[i].index;
                bytes32 value = _immutables[i].value;
                immutableDataStorage[uint256(uint160(_dest))][index] = value;
            }
        }
    }
}
