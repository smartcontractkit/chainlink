// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";

contract VerifierConstructorTest is BaseTest {
  function test_revertsIfInitializedWithEmptyVerifierProxy() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    new Verifier(address(0));
  }

  function test_setsTheCorrectProperties() public {
    Verifier v = new Verifier(address(s_verifierProxy));
    assertEq(v.owner(), ADMIN);

    (bool scanLogs, bytes32 configDigest, uint32 epoch) = v.latestConfigDigestAndEpoch(FEED_ID);
    assertEq(scanLogs, false);
    assertEq(configDigest, EMPTY_BYTES);
    assertEq(epoch, 0);

    (uint32 configCount, uint32 blockNumber, bytes32 configDigestTwo) = v.latestConfigDetails(FEED_ID);
    assertEq(configCount, 0);
    assertEq(blockNumber, 0);
    assertEq(configDigestTwo, EMPTY_BYTES);

    string memory typeAndVersion = s_verifier.typeAndVersion();
    assertEq(typeAndVersion, "Verifier 0.0.2");
  }
}

contract VerifierSupportsInterfaceTest is BaseTest {
  function test_falseIfIsNotCorrectInterface() public {
    bool isInterface = s_verifier.supportsInterface(bytes4("abcd"));
    assertEq(isInterface, false);
  }

  function test_trueIfIsCorrectInterface() public {
    bool isInterface = s_verifier.supportsInterface(Verifier.verify.selector);
    assertEq(isInterface, true);
  }
}
