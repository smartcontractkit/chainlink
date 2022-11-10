pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

contract BaseTest is Test {
    address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;

    function setUp() public virtual {
        // Set the sender to OWNER permanently
        changePrank(OWNER);
    }
}
