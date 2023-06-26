// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

/*
 * @title Common
 * @author Michael Fletcher
 * @notice Common functions and structs
 */
library Common {

    // @notice The asset struct to hold an address of an asset and amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }

    // @notice Struct to hold the address and it's associated weight
    struct AddressAndWeight {
        address addr;
        uint16 weight;
    }
}