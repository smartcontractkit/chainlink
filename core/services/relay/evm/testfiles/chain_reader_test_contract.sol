// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

struct TestStruct {
    int32 Field;
    string DifferentField;
    uint8 OracleId;
    uint8[32] OracleIds;
    bytes32 Account;
    bytes32[] Accounts;
    int192 BigField;
    MidLevelTestStruct NestedStruct;
}

struct MidLevelTestStruct {
    bytes2 FixedBytes;
    InnerTestStruct Inner;
}

struct InnerTestStruct {
    int64 I;
    string S;
}

contract LatestValueHolder {
    event Triggered(
        int32 field,
        string differentField,
        uint8 oracleId,
        uint8[32] oracleIds,
        bytes32 account,
        bytes32[] accounts,
        int192 bigField,
        MidLevelTestStruct nestedStruct);

    TestStruct[] private seen;
    uint64[] private arr;

    constructor() {
        // See chain_reader_interface_tests.go in chainlink-relay
        arr.push(3);
        arr.push(4);
    }

    function AddTestStruct(
        int32 field,
        string calldata differentField,
        uint8 oracleId,
        uint8[32] calldata oracleIds,
        bytes32 account,
        bytes32[] calldata accounts,
        int192 bigField,
        MidLevelTestStruct calldata nestedStruct) public {
        seen.push(TestStruct(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct));
    }

    function ReturnSeen(
        int32 field,
        string calldata differentField,
        uint8 oracleId,
        uint8[32] calldata oracleIds,
        bytes32 account,
        bytes32[] calldata accounts,
        int192 bigField,
        MidLevelTestStruct calldata nestedStruct) pure public returns (TestStruct memory) {
        return TestStruct(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct);
    }

    function GetElementAtIndex(uint256 i) public view returns (TestStruct memory) {
        // See chain_reader_interface_tests.go in chainlink-relay
        return seen[i-1];
    }

    function GetPrimitiveValue() public pure returns (uint64) {
        // See chain_reader_interface_tests.go in chainlink-relay
        return 3;
    }

    function GetSliceValue() public view returns (uint64[] memory) {
        return arr;
    }

    function TriggerEvent(int32 field,
        string calldata differentField,
        uint8 oracleId,
        uint8[32] calldata oracleIds,
        bytes32 account,
        bytes32[] calldata accounts,
        int192 bigField,
        MidLevelTestStruct calldata nestedStruct) public {
        emit Triggered(field, differentField, oracleId, oracleIds, account, accounts, bigField, nestedStruct);
    }
}