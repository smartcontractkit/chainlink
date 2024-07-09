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

  function getRandomAddresses(uint256 length) internal returns (address[] memory) {
    address[] memory addresses = new address[](length);
    for (uint256 i = 0; i < length; ++i) {
      addresses[i] = address(uint160(uint(keccak256(abi.encodePacked(i)))));
    }
    return addresses;
  }

  function addressIsIn(address addr, address[] memory addresses) internal returns (bool) {
    for (uint256 i = 0; i < addresses.length; ++i) {
      if (addresses[i] == addr) {
        return true;
      }
    }
    return false;
  }
}
