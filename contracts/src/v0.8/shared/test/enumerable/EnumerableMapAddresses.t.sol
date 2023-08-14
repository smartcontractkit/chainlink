// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "../BaseTest.t.sol";
import {EnumerableMapAddressesHelper} from "./EnumerableMapAddressesHelper.sol";

contract EnumerableMapAddressesTest is BaseTest {
  EnumerableMapAddressesHelper internal s_helper;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_helper = new EnumerableMapAddressesHelper();
  }
}

contract EnumerableMapAddresses_set is EnumerableMapAddressesTest {
  function testSetSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(!s_helper.set(address(this), address(this)));
  }
}

contract EnumerableMapAddresses_remove is EnumerableMapAddressesTest {
  function testRemoveSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(s_helper.remove(address(this)));
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(!s_helper.remove(address(this)));
  }
}

contract EnumerableMapAddresses_contains is EnumerableMapAddressesTest {
  function testContainsSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
  }
}

contract EnumerableMapAddresses_length is EnumerableMapAddressesTest {
  function testLengthSuccess() public {
    assertTrue(s_helper.length() == 0);
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.length() == 1);
    assertTrue(s_helper.remove(address(this)));
    assertTrue(s_helper.length() == 0);
  }
}

contract EnumerableMapAddresses_at is EnumerableMapAddressesTest {
  function testAtSuccess() public {
    assertTrue(s_helper.length() == 0);
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.length() == 1);
    (address key, address value) = s_helper.at(0);
    assertTrue(key == address(this));
    assertTrue(value == address(this));
  }
}

contract EnumerableMapAddresses_tryGet is EnumerableMapAddressesTest {
  function testTryGetSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    (bool success, address value) = s_helper.tryGet(address(this));
    assertTrue(success);
    assertTrue(value == address(this));
  }
}

contract EnumerableMapAddresses_get is EnumerableMapAddressesTest {
  function testGetSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(s_helper.get(address(this)) == address(this));
  }
}

contract EnumerableMapAddresses_get_errorMessage is EnumerableMapAddressesTest {
  function testGetErrorMessageSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(s_helper.get(address(this), "EnumerableMapAddresses: nonexistent key") == address(this));
  }
}
