// SPDX-License-Identifier: MIT
pragma solidity >=0.6.0 <0.9.0;

pragma experimental ABIEncoderV2;

import {VmSafe} from "./Vm.sol";

// Helpers for parsing and writing TOML files
// To parse:
// ```
// using stdToml for string;
// string memory toml = vm.readFile("<some_path>");
// toml.readUint("<json_path>");
// ```
// To write:
// ```
// using stdToml for string;
// string memory json = "json";
// json.serialize("a", uint256(123));
// string memory semiFinal = json.serialize("b", string("test"));
// string memory finalJson = json.serialize("c", semiFinal);
// finalJson.write("<some_path>");
// ```

library stdToml {
    VmSafe private constant vm = VmSafe(address(uint160(uint256(keccak256("hevm cheat code")))));

    function parseRaw(string memory toml, string memory key) internal pure returns (bytes memory) {
        return vm.parseToml(toml, key);
    }

    function readUint(string memory toml, string memory key) internal pure returns (uint256) {
        return vm.parseTomlUint(toml, key);
    }

    function readUintArray(string memory toml, string memory key) internal pure returns (uint256[] memory) {
        return vm.parseTomlUintArray(toml, key);
    }

    function readInt(string memory toml, string memory key) internal pure returns (int256) {
        return vm.parseTomlInt(toml, key);
    }

    function readIntArray(string memory toml, string memory key) internal pure returns (int256[] memory) {
        return vm.parseTomlIntArray(toml, key);
    }

    function readBytes32(string memory toml, string memory key) internal pure returns (bytes32) {
        return vm.parseTomlBytes32(toml, key);
    }

    function readBytes32Array(string memory toml, string memory key) internal pure returns (bytes32[] memory) {
        return vm.parseTomlBytes32Array(toml, key);
    }

    function readString(string memory toml, string memory key) internal pure returns (string memory) {
        return vm.parseTomlString(toml, key);
    }

    function readStringArray(string memory toml, string memory key) internal pure returns (string[] memory) {
        return vm.parseTomlStringArray(toml, key);
    }

    function readAddress(string memory toml, string memory key) internal pure returns (address) {
        return vm.parseTomlAddress(toml, key);
    }

    function readAddressArray(string memory toml, string memory key) internal pure returns (address[] memory) {
        return vm.parseTomlAddressArray(toml, key);
    }

    function readBool(string memory toml, string memory key) internal pure returns (bool) {
        return vm.parseTomlBool(toml, key);
    }

    function readBoolArray(string memory toml, string memory key) internal pure returns (bool[] memory) {
        return vm.parseTomlBoolArray(toml, key);
    }

    function readBytes(string memory toml, string memory key) internal pure returns (bytes memory) {
        return vm.parseTomlBytes(toml, key);
    }

    function readBytesArray(string memory toml, string memory key) internal pure returns (bytes[] memory) {
        return vm.parseTomlBytesArray(toml, key);
    }

    function serialize(string memory jsonKey, string memory rootObject) internal returns (string memory) {
        return vm.serializeJson(jsonKey, rootObject);
    }

    function serialize(string memory jsonKey, string memory key, bool value) internal returns (string memory) {
        return vm.serializeBool(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, bool[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeBool(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, uint256 value) internal returns (string memory) {
        return vm.serializeUint(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, uint256[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeUint(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, int256 value) internal returns (string memory) {
        return vm.serializeInt(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, int256[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeInt(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, address value) internal returns (string memory) {
        return vm.serializeAddress(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, address[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeAddress(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, bytes32 value) internal returns (string memory) {
        return vm.serializeBytes32(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, bytes32[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeBytes32(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, bytes memory value) internal returns (string memory) {
        return vm.serializeBytes(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, bytes[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeBytes(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, string memory value)
        internal
        returns (string memory)
    {
        return vm.serializeString(jsonKey, key, value);
    }

    function serialize(string memory jsonKey, string memory key, string[] memory value)
        internal
        returns (string memory)
    {
        return vm.serializeString(jsonKey, key, value);
    }

    function write(string memory jsonKey, string memory path) internal {
        vm.writeToml(jsonKey, path);
    }

    function write(string memory jsonKey, string memory path, string memory valueKey) internal {
        vm.writeToml(jsonKey, path, valueKey);
    }
}
