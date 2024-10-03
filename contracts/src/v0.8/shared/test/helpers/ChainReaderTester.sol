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
  MidLevelDynamicTestStruct NestedDynamicStruct;
  MidLevelStaticTestStruct NestedStaticStruct;
}

struct MidLevelDynamicTestStruct {
  bytes2 FixedBytes;
  InnerDynamicTestStruct Inner;
}

struct InnerDynamicTestStruct {
  int64 IntVal;
  string S;
}

struct MidLevelStaticTestStruct {
  bytes2 FixedBytes;
  InnerStaticTestStruct Inner;
}

struct InnerStaticTestStruct {
  int64 IntVal;
  address A;
}

contract ChainReaderTester {
  event Triggered(
    int32 indexed field,
    uint8 oracleId,
    MidLevelDynamicTestStruct nestedDynamicStruct,
    MidLevelStaticTestStruct nestedStaticStruct,
    uint8[32] oracleIds,
    address Account,
    address[] Accounts,
    string differentField,
    int192 bigField
  );

  event TriggeredEventWithDynamicTopic(string indexed fieldHash, string field);

  // First topic is event hash
  event TriggeredWithFourTopics(int32 indexed field1, int32 indexed field2, int32 indexed field3);

  // first topic is event hash, second and third topics get hashed before getting stored
  event TriggeredWithFourTopicsWithHashed(string indexed field1, uint8[32] indexed field2, bytes32 indexed field3);

  // emits dynamic bytes which encode data in the same way every time.
  event StaticBytes(bytes message);

  TestStruct[] private s_seen;
  uint64[] private s_arr;
  uint64 private s_value;

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
    MidLevelDynamicTestStruct calldata nestedDynamicStruct,
    MidLevelStaticTestStruct calldata nestedStaticStruct
  ) public {
    s_seen.push(
      TestStruct(
        field,
        differentField,
        oracleId,
        oracleIds,
        account,
        accounts,
        bigField,
        nestedDynamicStruct,
        nestedStaticStruct
      )
    );
  }

  function setAlterablePrimitiveValue(uint64 value) public {
    s_value = value;
  }

  function returnSeen(
    int32 field,
    string calldata differentField,
    uint8 oracleId,
    uint8[32] calldata oracleIds,
    address account,
    address[] calldata accounts,
    int192 bigField,
    MidLevelDynamicTestStruct calldata nestedDynamicStruct,
    MidLevelStaticTestStruct calldata nestedStaticStruct
  ) public pure returns (TestStruct memory) {
    return
      TestStruct(
        field,
        differentField,
        oracleId,
        oracleIds,
        account,
        accounts,
        bigField,
        nestedDynamicStruct,
        nestedStaticStruct
      );
  }

  function getElementAtIndex(uint256 i) public view returns (TestStruct memory) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return s_seen[i - 1];
  }

  function getPrimitiveValue() public pure returns (uint64) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return 3;
  }

  function getAlterablePrimitiveValue() public view returns (uint64) {
    // See chain_reader_interface_tests.go in chainlink-relay
    return s_value;
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
    uint8 oracleId,
    MidLevelDynamicTestStruct calldata nestedDynamicStruct,
    MidLevelStaticTestStruct calldata nestedStaticStruct,
    uint8[32] calldata oracleIds,
    address account,
    address[] calldata accounts,
    string calldata differentField,
    int192 bigField
  ) public {
    emit Triggered(
      field,
      oracleId,
      nestedDynamicStruct,
      nestedStaticStruct,
      oracleIds,
      account,
      accounts,
      differentField,
      bigField
    );
  }

  function triggerEventWithDynamicTopic(string calldata field) public {
    emit TriggeredEventWithDynamicTopic(field, field);
  }

  // first topic is the event signature
  function triggerWithFourTopics(int32 field1, int32 field2, int32 field3) public {
    emit TriggeredWithFourTopics(field1, field2, field3);
  }

  // first topic is event hash, second and third topics get hashed before getting stored
  function triggerWithFourTopicsWithHashed(string memory field1, uint8[32] memory field2, bytes32 field3) public {
    emit TriggeredWithFourTopicsWithHashed(field1, field2, field3);
  }

  // emulate CCTP message event.
  function triggerStaticBytes(
    uint32 val1,
    uint32 val2,
    uint32 val3,
    uint64 val4,
    bytes32 val5,
    bytes32 val6,
    bytes32 val7,
    bytes memory raw
  ) public {
    bytes memory _message = abi.encodePacked(val1, val2, val3, val4, val5, val6, val7, raw);
    emit StaticBytes(_message);
  }
}
