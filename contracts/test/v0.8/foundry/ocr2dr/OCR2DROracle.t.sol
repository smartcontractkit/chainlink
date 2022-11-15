pragma solidity ^0.8.0;

import {BaseTest} from "../BaseTest.t.sol";
import {OCR2DROracle} from "../../../../src/v0.8/dev/ocr2dr/OCR2DROracle.sol";

contract OCR2DROracleSetup is BaseTest {
    bytes constant DATA = abi.encode("bytes");

    OCR2DROracle s_oracle;

    function setUp() public virtual override {
        BaseTest.setUp();

        s_oracle = new OCR2DROracle();
    }
}

contract OCR2DROracle_typeAndVersion is OCR2DROracleSetup {
    function testTypeAndVersionSuccess() public {
       assertEq(s_oracle.typeAndVersion(), "OCR2DROracle 0.0.0");
    }
}

contract OCR2DROracle_setDONPublicKey is OCR2DROracleSetup {
    function testSetDONPublicKey_gas() public {
        s_oracle.setDONPublicKey(DATA);
    }

    function testSetDONPublicKeySuccess() public {
        bytes memory donPublicKey = abi.encode("newKey");

        // Verify the existing key is different from the new key
        bytes memory existingDonPublicKey = s_oracle.getDONPublicKey();
        bytes memory expectedExistingDonPublicKey;
        assertEq(existingDonPublicKey, expectedExistingDonPublicKey);
        // If they have different lengths they are not the same
        assertFalse(donPublicKey.length == expectedExistingDonPublicKey.length);

        s_oracle.setDONPublicKey(donPublicKey);
        bytes memory newDonPublicKey = s_oracle.getDONPublicKey();
        assertEq(newDonPublicKey, donPublicKey);
    }

    // Reverts

    function testEmptyPublicKeyReverts() public {
        bytes memory donPublicKey;

        vm.expectRevert(OCR2DROracle.EmptyPublicKey.selector);
        s_oracle.setDONPublicKey(donPublicKey);
    }

    function testOnlyOwnerReverts() public {
        vm.stopPrank();
        vm.expectRevert("Only callable by owner");

        bytes memory donPublicKey;
        s_oracle.setDONPublicKey(donPublicKey);
    }
}

contract OCR2DROracle_sendRequest is OCR2DROracleSetup {
    event OracleRequest(bytes32 requestId, bytes data);

    function testSendRequest_gas() public {
        s_oracle.sendRequest(0, DATA);
    }

    function testSendRequestFuzzSuccess(uint256 subscriptionId, bytes calldata data) public {
        vm.assume(data.length != 0);

        vm.expectEmit(false, false, false, false);
        emit OracleRequest(0, data);

        s_oracle.sendRequest(subscriptionId, data);
    }

    // Reverts

    function testEmptyRequestDataReverts() public {
        bytes memory emptyData;

        vm.expectRevert(OCR2DROracle.EmptyRequestData.selector);
        s_oracle.sendRequest(0, emptyData);
    }
}