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
    TestStruct[] private seen;
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

    function GetElementAtIndex(uint256 i) public view returns (TestStruct memory) {
        return seen[i-1];
    }
}