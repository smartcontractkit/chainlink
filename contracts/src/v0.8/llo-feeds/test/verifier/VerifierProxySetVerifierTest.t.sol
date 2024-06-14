// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../../interfaces/IVerifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {AccessControllerInterface} from "../../../shared/interfaces/AccessControllerInterface.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC165.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierProxyInitializeVerifierTest is BaseTestWithConfiguredVerifierAndFeeManager {
  function test_revertsIfNotCorrectVerifier() public {
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.AccessForbidden.selector));
    s_verifierProxy.setVerifier(bytes32("prev-config"), bytes32("new-config"), new Common.AddressAndWeight[](0));
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
    s_verifierProxy.setVerifier(maliciousDigest, takenDigest, new Common.AddressAndWeight[](0));
  }

  function test_updatesVerifierIfVerifier() public {
    (, , bytes32 prevDigest) = s_verifier.latestConfigDetails(FEED_ID);
    changePrank(address(s_verifier));
    s_verifierProxy.setVerifier(prevDigest, bytes32("new-config"), new Common.AddressAndWeight[](0));
    assertEq(s_verifierProxy.getVerifier(bytes32("new-config")), address(s_verifier));
    assertEq(s_verifierProxy.getVerifier(prevDigest), address(s_verifier));
  }
}
