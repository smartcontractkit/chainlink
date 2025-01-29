// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";

contract VerifierConstructorTestV05 is BaseTest {
  function test_revertsIfInitializedWithEmptyVerifierProxy() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    new Verifier(address(0));
  }

  function test_setsTheCorrectProperties() public {
    Verifier v = new Verifier(address(s_verifierProxy));
    assertEq(v.owner(), ADMIN);

    uint32 blockNumber = v.latestConfigDetails(FEED_ID);
    assertEq(blockNumber, 0);

    string memory typeAndVersion = s_verifier.typeAndVersion();
    assertEq(typeAndVersion, "Verifier 2.0.0");
  }
}

contract VerifierSupportsInterfaceTest is BaseTest {
  function test_falseIfIsNotCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(bytes4("abcd"));
    assertEq(isInterface, false);
  }

  function test_trueIfIsCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(Verifier.verify.selector);
    assertEq(isInterface, true);
  }
}
