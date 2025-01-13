// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "forge-std/Test.sol";

contract BaseTest is Test {
  bool private s_baseTestInitialized;
  address internal constant OWNER = 0x72da681452Ab957d1020c25fFaCA47B43980b7C3;
  address internal constant STRANGER = 0x02e7d5DD1F4dDbC9f512FfA01d30aa190Ae3edBb;

  // Fri May 26 2023 13:49:53 GMT+0000
  uint256 internal constant BLOCK_TIME = 1685108993;

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    vm.label(OWNER, "Owner");
    vm.label(STRANGER, "Stranger");

    // Set the sender to OWNER permanently
    vm.startPrank(OWNER);
    deal(OWNER, 1e20);

    // Set the block time to a constant known value
    vm.warp(BLOCK_TIME);
  }
}
