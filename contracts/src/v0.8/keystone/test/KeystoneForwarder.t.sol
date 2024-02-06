// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "forge-std/Test.sol";

import "../KeystoneForwarder.sol";

contract KeystoneForwarderTest is Test {
  function setUp() public virtual {}

  function test_abi_partial_decoding_works() public {
    bytes memory report = hex"0102";
    uint256 amount = 1;
    // bytes memory payload = abi.encodeWithSignature("transfer(bytes,uint256)", report, amount);
    bytes memory payload = abi.encode(report, amount);
    bytes memory decodedReport = abi.decode(payload, (bytes));
    assertEq(decodedReport, report, "not equal");
  }
}
