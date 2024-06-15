pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";

contract BaseTest is Test {
  bool private s_baseTestInitialized;
  address internal constant OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    // Set msg.sender to OWNER until changePrank or stopPrank is called
    vm.startPrank(OWNER);
  }
}
