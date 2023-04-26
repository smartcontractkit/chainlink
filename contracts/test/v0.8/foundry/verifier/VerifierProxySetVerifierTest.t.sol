// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifier} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../../../../src/v0.8/interfaces/IVerifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../src/v0.8/interfaces/AccessControllerInterface.sol";
import {IERC165} from "@openzeppelin/contracts/interfaces/IERC165.sol";

contract VerifierProxyInitializeVerifierTest is BaseTestWithConfiguredVerifier {
  function test_revertsIfNotCorrectVerifier() public {
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.AccessForbidden.selector));
    s_verifierProxy.setVerifier(bytes32("prev-config"), bytes32("new-config"));
  }

  function test_revertsIfDigestAlreadySet() public {
    (, , bytes32 takenDigest) = s_verifier.latestConfigDetails(FEED_ID);

    address maliciousVerifier = address(666);
    bytes32 maliciousDigest = bytes32("malicious-digest");
    vm.mockCall(
      maliciousVerifier,
      abi.encodeWithSelector(IERC165.supportsInterface.selector, IVerifier.verify.selector),
      abi.encode(true)
    );
    s_verifierProxy.initializeVerifier(maliciousVerifier);
    vm.expectRevert(
      abi.encodeWithSelector(VerifierProxy.ConfigDigestAlreadySet.selector, takenDigest, address(s_verifier))
    );
    changePrank(address(maliciousVerifier));
    s_verifierProxy.setVerifier(maliciousDigest, takenDigest);
  }

  function test_updatesVerifierIfVerifier() public {
    (, , bytes32 prevDigest) = s_verifier.latestConfigDetails(FEED_ID);
    changePrank(address(s_verifier));
    s_verifierProxy.setVerifier(prevDigest, bytes32("new-config"));
    assertEq(s_verifierProxy.getVerifier(bytes32("new-config")), address(s_verifier));
    assertEq(s_verifierProxy.getVerifier(prevDigest), address(s_verifier));
  }
}
