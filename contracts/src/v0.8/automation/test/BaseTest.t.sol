// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "forge-std/Test.sol";

contract BaseTest is Test {
  address internal OWNER = 0x00007e64E1fB0C487F25dd6D3601ff6aF8d32e4e;
  address internal constant STRANGER = address(999);

  function setUp() public virtual {
    vm.startPrank(OWNER);
    deal(OWNER, 1e20);
  }
}
