pragma solidity ^0.8.0;

import {BaseTest} from "../BaseTest.t.sol";
import {OCR2DROracle} from "../../../../src/v0.8/dev/ocr2dr/OCR2DROracle.sol";

contract OCR2DROracleSetup is BaseTest {
    OCR2DROracle s_oracle;

    function setUp() public virtual override {
        BaseTest.setUp();

        s_oracle = new OCR2DROracle();
    }
}

contract OCR2DROracle_typeAndVersion is OCR2DROracleSetup {
    function testTypeAndVersionSuccess() public {
       assertEq(  s_oracle.typeAndVersion(), "OCR2DROracle 0.0.0");
    }
}

contract OCR2DROracle_setDONPublicKey is OCR2DROracleSetup {
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
        vm.expectRevert(OCR2DROracle.EmptyPublicKey.selector);

        bytes memory donPublicKey;
        s_oracle.setDONPublicKey(donPublicKey);
    }

    function testOnlyOwnerReverts() public {
        vm.stopPrank();
        vm.expectRevert("Only callable by owner");

        bytes memory donPublicKey;
        s_oracle.setDONPublicKey(donPublicKey);
    }
}
