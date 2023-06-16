// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract MockCustomerTarget{

    function performUpkeep() external pure returns (int) {
        return 1;
    }
}