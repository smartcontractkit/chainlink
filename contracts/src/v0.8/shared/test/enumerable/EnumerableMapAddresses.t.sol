// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {BaseTest} from "../BaseTest.t.sol";
import {EnumerableMapAddresses} from "../../enumerable/EnumerableMapAddresses.sol";

contract EnumerableMapAddressesTest is BaseTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  EnumerableMapAddresses.AddressToAddressMap internal s_addressToAddressMap;
  EnumerableMapAddresses.AddressToBytes32Map internal s_addressToBytes32Map;

  bytes32 internal constant MOCK_BYTES32_VALUE = bytes32(uint256(42));

  function setUp() public virtual override {
    BaseTest.setUp();
  }
}

contract EnumerableMapAddresses_set is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testSetSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
    assertTrue(!s_addressToAddressMap.set(address(this), address(this)));
  }

  function testBytes32SetSuccess() public {
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.contains(address(this)));
    assertTrue(!s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
  }
}

contract EnumerableMapAddresses_remove is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testRemoveSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.remove(address(this)));
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(!s_addressToAddressMap.remove(address(this)));
  }

  function testBytes32RemoveSuccess() public {
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.remove(address(this)));
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(!s_addressToBytes32Map.remove(address(this)));
  }
}

contract EnumerableMapAddresses_contains is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testContainsSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
  }

  function testBytes32ContainsSuccess() public {
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.contains(address(this)));
  }
}

contract EnumerableMapAddresses_length is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testLengthSuccess() public {
    assertTrue(s_addressToAddressMap.length() == 0);
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.length() == 1);
    assertTrue(s_addressToAddressMap.remove(address(this)));
    assertTrue(s_addressToAddressMap.length() == 0);
  }

  function testBytes32LengthSuccess() public {
    assertTrue(s_addressToBytes32Map.length() == 0);
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.length() == 1);
    assertTrue(s_addressToBytes32Map.remove(address(this)));
    assertTrue(s_addressToBytes32Map.length() == 0);
  }
}

contract EnumerableMapAddresses_at is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testAtSuccess() public {
    assertTrue(s_addressToAddressMap.length() == 0);
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.length() == 1);
    (address key, address value) = s_addressToAddressMap.at(0);
    assertTrue(key == address(this));
    assertTrue(value == address(this));
  }

  function testBytes32AtSuccess() public {
    assertTrue(s_addressToBytes32Map.length() == 0);
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.length() == 1);
    (address key, bytes32 value) = s_addressToBytes32Map.at(0);
    assertTrue(key == address(this));
    assertTrue(value == MOCK_BYTES32_VALUE);
  }
}

contract EnumerableMapAddresses_tryGet is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testTryGetSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
    (bool success, address value) = s_addressToAddressMap.tryGet(address(this));
    assertTrue(success);
    assertTrue(value == address(this));
  }

  function testBytes32TryGetSuccess() public {
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.contains(address(this)));
    (bool success, bytes32 value) = s_addressToBytes32Map.tryGet(address(this));
    assertTrue(success);
    assertTrue(value == MOCK_BYTES32_VALUE);
  }
}

contract EnumerableMapAddresses_get is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testGetSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.get(address(this)) == address(this));
  }

  function testBytes32GetSuccess() public {
    assertTrue(!s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.set(address(this), MOCK_BYTES32_VALUE));
    assertTrue(s_addressToBytes32Map.contains(address(this)));
    assertTrue(s_addressToBytes32Map.get(address(this)) == MOCK_BYTES32_VALUE);
  }
}

contract EnumerableMapAddresses_get_errorMessage is EnumerableMapAddressesTest {
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToAddressMap;

  function testGetErrorMessageSuccess() public {
    assertTrue(!s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.set(address(this), address(this)));
    assertTrue(s_addressToAddressMap.contains(address(this)));
    assertTrue(s_addressToAddressMap.get(address(this), "EnumerableMapAddresses: nonexistent key") == address(this));
  }
}
