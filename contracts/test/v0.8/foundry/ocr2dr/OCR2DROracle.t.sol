pragma solidity ^0.8.0;

import {BaseTest} from "../BaseTest.t.sol";
import {OCR2DROracle} from "../../../../src/v0.8/dev/ocr2dr/OCR2DROracle.sol";
import {OCR2DRRegistry} from "../../../../src/v0.8/dev/ocr2dr/OCR2DRRegistry.sol";

// import {LinkToken} from "../../../../src/v0.4/LinkToken.sol";
// import {MockV3Aggregator} from "../../../../src/v0.7/tests/MockV3Aggregator.sol";

contract OCR2DROracleSetup is BaseTest {
  bytes constant DATA = abi.encode("bytes");
  address registryAddress = makeAddr("Registry");

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

contract OCR2DROracle_setRegistry is OCR2DROracleSetup {
  function testSetRegistry_gas() public {
    s_oracle.setRegistry(registryAddress);
  }

  function testSetRegistrySuccess() public {
    address registryAddress = makeAddr("newRegistry");

    // Verify the existing key is different from the new key
    address existingRegistryAddress = s_oracle.getRegistry();
    address expectedRegistryAddress;
    assertEq(existingRegistryAddress, expectedRegistryAddress);

    s_oracle.setRegistry(registryAddress);
    address newRegistryAddress = s_oracle.getRegistry();
    assertEq(registryAddress, newRegistryAddress);
  }

  // Reverts

  function testEmptyPublicKeyReverts() public {
    address registryAddress;

    vm.expectRevert(OCR2DROracle.EmptyBillingRegistry.selector);
    s_oracle.setRegistry(registryAddress);
  }

  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");

    address registryAddress;
    s_oracle.setRegistry(registryAddress);
  }
}

contract OCR2DROracle_sendRequest is OCR2DROracleSetup {
  OCR2DRRegistry s_registry;

  //   LinkToken s_link;
  //   MockV3Aggregator s_linketh;

  function setUp() public virtual override {
    OCR2DROracleSetup.setUp();

    // s_link = new LinkToken();
    // s_linketh = new MockV3Aggregator(0, 5021530000000000);
    s_registry = new OCR2DRRegistry(makeAddr("Link Token"), makeAddr("Link Eth"));
    s_oracle.setRegistry(address(s_registry));
  }

  event OracleRequest(bytes32 requestId, bytes data);

  // TODO: write new ^0.8.0 mocks for LinkToken & MockV3Aggregator
  //   function testSendRequest_gas() public {
  //     s_oracle.sendRequest(0, DATA, 0);
  //   }

  //   function testSendRequestFuzzSuccess(uint64 subscriptionId, bytes calldata data) public {
  //     vm.assume(data.length != 0);

  //     vm.expectEmit(false, false, false, false);
  //     emit OracleRequest(0, data);

  //     s_oracle.sendRequest(subscriptionId, data, 0);
  //   }

  // Reverts

  function testEmptyRequestDataReverts() public {
    bytes memory emptyData;

    vm.expectRevert(OCR2DROracle.EmptyRequestData.selector);
    s_oracle.sendRequest(0, emptyData, 0);
  }
}
