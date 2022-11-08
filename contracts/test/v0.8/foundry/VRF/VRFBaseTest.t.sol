pragma solidity ^0.8.0;

import {BaseTest} from "../BaseTest.t.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";

contract VRFBaseTest is BaseTest {
    MockV3Aggregator internal LINK_ETH_FEED = new MockV3Aggregator(18, 10000000000000000); // .01 ETH default answer

    function setUp() public virtual override {
        BaseTest.setUp();
    }
}
