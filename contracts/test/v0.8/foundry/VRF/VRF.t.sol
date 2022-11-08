pragma solidity ^0.8.0;

import {VRFBaseTest} from "./VRFBaseTest.t.sol";

contract VRF is VRFBaseTest {
    function setUp() public virtual override {
        VRFBaseTest.setUp();
    }

    function testV3AggregatorSuccess() public {
        assertEq(LINK_ETH_FEED.decimals(), 18);
    }
}
