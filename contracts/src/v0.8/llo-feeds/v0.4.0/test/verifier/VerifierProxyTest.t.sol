// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationFeeManager} from "../../../v0.4.0/DestinationFeeManager.sol";
import {IERC165} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

contract VerifierProxyInitializeVerifierTest is BaseTest {
    function test_verifyCallWithNoVerifierSet() public {
        bytes memory somebytes = "something";
        vm.expectRevert(abi.encodeWithSelector(DestinationVerifierProxy.ZeroAddress.selector));
        s_verifierProxy.verify(somebytes, somebytes);
    }

    function test_verifyBulkCallWithNoVerifierSet() public {
        bytes memory somebytes = "something";
        bytes[] memory arrayData = new bytes[](1);
        arrayData[0] = somebytes;
        vm.expectRevert(abi.encodeWithSelector(DestinationVerifierProxy.ZeroAddress.selector));
        s_verifierProxy.verifyBulk(arrayData, somebytes);
    }

    function test_setVerifierCalledByNoOwner() public {
        address STRANGER = address(999);
        changePrank(STRANGER);
        vm.expectRevert(bytes("Only callable by owner"));
        s_verifierProxy.setVerifier(address(s_verifier));
    }

    function test_setVerifierZeroVerifier() public {
        vm.expectRevert(abi.encodeWithSelector(DestinationVerifierProxy.ZeroAddress.selector));
        s_verifierProxy.setVerifier(address(0));
    }

    function test_setVerifierWhichDoesntHonourInterface() public {
        vm.expectRevert(
            abi.encodeWithSelector(DestinationVerifierProxy.VerifierInvalid.selector, address(rewardManager))
        );
        s_verifierProxy.setVerifier(address(rewardManager));
    }

    function test_setVerifierOk() public {
        s_verifierProxy.setVerifier(address(s_verifier));
    }
}
