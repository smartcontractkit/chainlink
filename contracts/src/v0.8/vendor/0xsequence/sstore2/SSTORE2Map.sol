// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../Create3/Create3.sol";

import "./utils/Bytecode.sol";


/**
  @title A write-once key-value storage for storing chunks of data with a lower write & read cost.
  @author Agustin Aguilar <aa@horizon.io>

  Readme: https://github.com/0xsequence/sstore2#readme
*/
library SSTORE2Map {
    error WriteError();

    //                                         keccak256(bytes('@0xSequence.SSTORE2Map.slot'))
    bytes32 private constant SLOT_KEY_PREFIX = 0xd351a9253491dfef66f53115e9e3afda3b5fdef08a1de6937da91188ec553be5;

    function internalKey(bytes32 _key) internal pure returns (bytes32) {
        // Mutate the key so it doesn't collide
        // if the contract is also using CREATE3 for other things
        return keccak256(abi.encode(SLOT_KEY_PREFIX, _key));
    }

    /**
@notice Stores `_data` and returns `pointer` as key for later retrieval
    @dev The pointer is a contract address with `_data` as code
    @param _data To be written
    @param _key unique string key for accessing the written data (can only be used once)
    @return pointer Pointer to the written `_data`
  */
    function write(string memory _key, bytes memory _data) internal returns (address pointer) {
        return write(keccak256(bytes(_key)), _data);
    }

    /**
@notice Stores `_data` and returns `pointer` as key for later retrieval
    @dev The pointer is a contract address with `_data` as code
    @param _data to be written
    @param _key unique bytes32 key for accessing the written data (can only be used once)
    @return pointer Pointer to the written `_data`
  */
    function write(bytes32 _key, bytes memory _data) internal returns (address pointer) {
        // Append 00 to _data so contract can't be called
        // Build init code
        bytes memory code = Bytecode.creationCodeFor(
            abi.encodePacked(
                hex'00',
                _data
            )
        );

        // Deploy contract using create3
        pointer = Create3.create3(internalKey(_key), code);
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key string key that constains the data
    @return data read from contract associated with `_key`
  */
    function read(string memory _key) internal view returns (bytes memory) {
        return read(keccak256(bytes(_key)));
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key string key that constains the data
    @param _start number of bytes to skip
    @return data read from contract associated with `_key`
  */
    function read(string memory _key, uint256 _start) internal view returns (bytes memory) {
        return read(keccak256(bytes(_key)), _start);
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key string key that constains the data
    @param _start number of bytes to skip
    @param _end index before which to end extraction
    @return data read from contract associated with `_key`
  */
    function read(string memory _key, uint256 _start, uint256 _end) internal view returns (bytes memory) {
        return read(keccak256(bytes(_key)), _start, _end);
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key bytes32 key that constains the data
    @return data read from contract associated with `_key`
  */
    function read(bytes32 _key) internal view returns (bytes memory) {
        return Bytecode.codeAt(Create3.addressOf(internalKey(_key)), 1, type(uint256).max);
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key bytes32 key that constains the data
    @param _start number of bytes to skip
    @return data read from contract associated with `_key`
  */
    function read(bytes32 _key, uint256 _start) internal view returns (bytes memory) {
        return Bytecode.codeAt(Create3.addressOf(internalKey(_key)), _start + 1, type(uint256).max);
    }

    /**
@notice Reads the contents for a given `_key`, it maps to a contract code as data, skips the first byte
    @dev The function is intended for reading pointers first written by `write`
    @param _key bytes32 key that constains the data
    @param _start number of bytes to skip
    @param _end index before which to end extraction
    @return data read from contract associated with `_key`
  */
    function read(bytes32 _key, uint256 _start, uint256 _end) internal view returns (bytes memory) {
        return Bytecode.codeAt(Create3.addressOf(internalKey(_key)), _start + 1, _end + 1);
    }
}