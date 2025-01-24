// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

import {SystemContractHelper} from "../libraries/SystemContractHelper.sol";
import {BOOTLOADER_FORMAL_ADDRESS, FORCE_DEPLOYER} from "../Constants.sol";
import {SystemCallFlagRequired, Unauthorized, CallerMustBeSystemContract, CallerMustBeBootloader, CallerMustBeForceDeployer} from "../SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice An abstract contract that is used to reuse modifiers across the system contracts.
 * @dev Solidity does not allow exporting modifiers via libraries, so
 * the only way to do reuse modifiers is to have a base contract
 * @dev Never add storage variables into this contract as some
 * system contracts rely on this abstract contract as on interface!
 */
abstract contract SystemContractBase {
    /// @notice Modifier that makes sure that the method
    /// can only be called via a system call.
    modifier onlySystemCall() {
        if (!SystemContractHelper.isSystemCall() && !SystemContractHelper.isSystemContract(msg.sender)) {
            revert SystemCallFlagRequired();
        }
        _;
    }

    /// @notice Modifier that makes sure that the method
    /// can only be called from a system contract.
    modifier onlyCallFromSystemContract() {
        if (!SystemContractHelper.isSystemContract(msg.sender)) {
            revert CallerMustBeSystemContract();
        }
        _;
    }

    /// @notice Modifier that makes sure that the method
    /// can only be called from a special given address.
    modifier onlyCallFrom(address caller) {
        if (msg.sender != caller) {
            revert Unauthorized(msg.sender);
        }
        _;
    }

    /// @notice Modifier that makes sure that the method
    /// can only be called from the bootloader.
    modifier onlyCallFromBootloader() {
        if (msg.sender != BOOTLOADER_FORMAL_ADDRESS) {
            revert CallerMustBeBootloader();
        }
        _;
    }

    /// @notice Modifier that makes sure that the method
    /// can only be called from the L1 force deployer.
    modifier onlyCallFromForceDeployer() {
        if (msg.sender != FORCE_DEPLOYER) {
            revert CallerMustBeForceDeployer();
        }
        _;
    }
}
