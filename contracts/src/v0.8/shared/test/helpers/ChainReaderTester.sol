// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable-next-line gas-struct-packing
struct TestStruct {
  int32 Field;
  string DifferentField;
  uint8 OracleId;
  uint8[32] OracleIds;
  address Account;
  address[] Accounts;
  int192 BigField;
  MidLevelTestStruct NestedStruct;
}

struct MidLevelTestStruct {
  bytes2 FixedBytes;
  InnerTestStruct Inner;
}

struct InnerTestStruct {
  int64 IntVal;
  string S;
}

contract ChainReaderTester {
  event Triggered(
    int32 indexed field,
    string differentField,
    uint8 oracleId,
    uint8[32] oracleIds,
    address Account,
    address[] Accounts,
    int192 bigField,
    MidLevelTestStruct nestedStruct
  );

  event TriggeredEventWithDynamicTopic(string indexed fieldHash, string field);

  // First topic is event hash
  event TriggeredWithFourTopics(int32 indexed field1, int32 indexed field2, int32 indexed field3);

  TestStruct[] private s_seen;
  uint64[] private s_arr;

  constructor() {
    // See chain_reader_interface_tests.go in chainlink-relay
    s_arr.push(3);
    s_arr.push(4);
  }

  function addTestStruct(
    int32 field,
    string calldata differentField,
    uint8 oracleId,
    uint8[32] calldata oracleIds,
    address account,
    address[] calldata accounts,
    int192 bigField,
    MidLevelTestStruct calldata nestedStruct
  ) public {
    s_seen.push(TestStruct(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct));
  }

  function returnSeen(
    int32 field,
    string calldata differentField,
    uint8 oracleId,
    uint8[32] calldata oracleIds,
    address account,
    address[] calldata accounts,
    int192 bigField,
    MidLevelTestStruct calldata nestedStruct
  ) public pure returns (TestStruct memory) {
    return TestStruct(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct);
  }

  function getElementAtIndex(uint256 i) public view returns (TestStruct memory) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return s_seen[i - 1];
  }

  function getPrimitiveValue() public pure returns (uint64) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return 3;
  }

  function getDifferentPrimitiveValue() public pure returns (uint64) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return 1990;
  }

  function getSliceValue() public view returns (uint64[] memory) {
    return s_arr;
  }

  function triggerEvent(
    int32 field,
    string calldata differentField,
    uint8 oracleId,
    uint8[32] calldata oracleIds,
    address account,
    address[] calldata accounts,
    int192 bigField,
    MidLevelTestStruct calldata nestedStruct
  ) public {
    emit Triggered(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct);
  }

  function triggerEventWithDynamicTopic(string calldata field) public {
    emit TriggeredEventWithDynamicTopic(field, field);
  }

  // first topic is the event signature
  function triggerWithFourTopics(int32 field1, int32 field2, int32 field3) public {
    emit TriggeredWithFourTopics(field1, field2, field3);
  }
}
