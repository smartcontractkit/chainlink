// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationFeeManager} from "../../../v0.4.0/DestinationFeeManager.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

contract DestinationVerifierProxyInitializeVerifierTest is BaseTest {
  function test_setVerifierCalledByNoOwner() public {
    address STRANGER = address(999);
    changePrank(STRANGER);
    vm.expectRevert(bytes("Only callable by owner"));
    s_verifierProxy.setVerifier(address(s_verifier));
  }

  function test_setVerifierWhichDoesntHonourInterface() public {
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifierProxy.VerifierInvalid.selector, address(rewardManager)));
    s_verifierProxy.setVerifier(address(rewardManager));
  }

  function test_setVerifierOk() public {
    s_verifierProxy.setVerifier(address(s_verifier));
    assertEq(s_verifierProxy.s_feeManager(), s_verifier.s_feeManager());
    assertEq(s_verifierProxy.s_accessController(), s_verifier.s_accessController());
  }

  function test_correctlySetsTheOwner() public {
    DestinationVerifierProxy proxy = new DestinationVerifierProxy();
    assertEq(proxy.owner(), ADMIN);
  }

  function test_correctlySetsVersion() public view {
    string memory version = s_verifierProxy.typeAndVersion();
    assertEq(version, "DestinationVerifierProxy 0.4.0");
  }
}
