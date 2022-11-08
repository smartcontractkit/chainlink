pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

contract BaseTest is Test {
    address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
    uint256 internal constant BLOCK_TIME = 1234567890;

    function setUp() public virtual {
        // Set the sender to OWNER permanently
        changePrank(OWNER);

        // Set the block time to a constant known value
        vm.warp(BLOCK_TIME);
    }
}
