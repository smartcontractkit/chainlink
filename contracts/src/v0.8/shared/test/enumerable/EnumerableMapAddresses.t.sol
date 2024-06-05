// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "../BaseTest.t.sol";
import {EnumerableMapAddressesBytes32Helper} from "./EnumerableMapAddressesBytes32Helper.sol";
import {EnumerableMapAddressesHelper} from "./EnumerableMapAddressesHelper.sol";

contract EnumerableMapAddressesTest is BaseTest {
  EnumerableMapAddressesHelper internal s_helper;
  EnumerableMapAddressesBytes32Helper internal s_bytes32Helper;

  bytes32 internal constant MOCK_BYTES32_VALUE = bytes32(uint256(42));

  function setUp() public virtual override {
    BaseTest.setUp();
    s_helper = new EnumerableMapAddressesHelper();
    s_bytes32Helper = new EnumerableMapAddressesBytes32Helper();
  }
}

contract EnumerableMapAddresses_set is EnumerableMapAddressesTest {
  function testSetSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(!s_helper.set(address(this), address(this)));
  }

  function testBytes32SetSuccess() public {
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.contains(address(this)));
    assertTrue(!s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
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

  function testBytes32RemoveSuccess() public {
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.remove(address(this)));
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(!s_bytes32Helper.remove(address(this)));
  }
}

contract EnumerableMapAddresses_contains is EnumerableMapAddressesTest {
  function testContainsSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
  }

  function testBytes32ContainsSuccess() public {
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.contains(address(this)));
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

  function testBytes32LengthSuccess() public {
    assertTrue(s_bytes32Helper.length() == 0);
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.length() == 1);
    assertTrue(s_bytes32Helper.remove(address(this)));
    assertTrue(s_bytes32Helper.length() == 0);
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

  function testBytes32AtSuccess() public {
    assertTrue(s_bytes32Helper.length() == 0);
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.length() == 1);
    (address key, bytes32 value) = s_bytes32Helper.at(0);
    assertTrue(key == address(this));
    assertTrue(value == MOCK_BYTES32_VALUE);
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

  function testBytes32TryGetSuccess() public {
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.contains(address(this)));
    (bool success, bytes32 value) = s_bytes32Helper.tryGet(address(this));
    assertTrue(success);
    assertTrue(value == MOCK_BYTES32_VALUE);
  }
}

contract EnumerableMapAddresses_get is EnumerableMapAddressesTest {
  function testGetSuccess() public {
    assertTrue(!s_helper.contains(address(this)));
    assertTrue(s_helper.set(address(this), address(this)));
    assertTrue(s_helper.contains(address(this)));
    assertTrue(s_helper.get(address(this)) == address(this));
  }

  function testBytes32GetSuccess() public {
    assertTrue(!s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_bytes32Helper.contains(address(this)));
    assertTrue(s_bytes32Helper.get(address(this)) == MOCK_BYTES32_VALUE);
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
