// SPDX-License-Identifier: MIT OR Apache-2.0
pragma solidity >=0.8.13 <0.9.0;

import {Vm} from "./Vm.sol";

struct FindData {
    uint256 slot;
    uint256 offsetLeft;
    uint256 offsetRight;
    bool found;
}

struct StdStorage {
    mapping(address => mapping(bytes4 => mapping(bytes32 => FindData))) finds;
    bytes32[] _keys;
    bytes4 _sig;
    uint256 _depth;
    address _target;
    bytes32 _set;
    bool _enable_packed_slots;
    bytes _calldata;
}

library stdStorageSafe {
    event SlotFound(address who, bytes4 fsig, bytes32 keysHash, uint256 slot);
    event WARNING_UninitedSlot(address who, uint256 slot);

    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));
    uint256 constant UINT256_MAX = 115792089237316195423570985008687907853269984665640564039457584007913129639935;

    /// @notice Returns the 4-byte function selector for `sigStr`.
    function sigs(string memory sigStr) internal pure returns (bytes4) {
        return bytes4(keccak256(bytes(sigStr)));
    }

    /// @notice Returns the encoded call parameters (keys or raw calldata) for the configured target function.
    function getCallParams(StdStorage storage self) internal view returns (bytes memory) {
        if (self._calldata.length == 0) {
            return _flatten(self._keys);
        } else {
            return self._calldata;
        }
    }

    /// @notice Calls the target contract with the configured parameters and returns the success flag and return value.
    function callTarget(StdStorage storage self) internal view returns (bool, bytes32) {
        bytes memory cd = abi.encodePacked(self._sig, getCallParams(self));
        (bool success, bytes memory rdat) = self._target.staticcall(cd);
        bytes32 result = _bytesToBytes32(rdat, 32 * self._depth);

        return (success, result);
    }

    /// @notice Returns whether mutating `slot` changes the return value of the configured target call.
    /// @dev Temporarily writes either `type(uint256).max` or `0` to the slot to detect sensitivity.
    function checkSlotMutatesCall(StdStorage storage self, bytes32 slot) internal returns (bool) {
        bytes32 prevSlotValue = vm.load(self._target, slot);
        (bool success, bytes32 prevReturnValue) = callTarget(self);

        bytes32 testVal = prevReturnValue == bytes32(0) ? bytes32(UINT256_MAX) : bytes32(0);
        vm.store(self._target, slot, testVal);

        (, bytes32 newReturnValue) = callTarget(self);

        vm.store(self._target, slot, prevSlotValue);

        return (success && (prevReturnValue != newReturnValue));
    }

    /// @notice Searches for the bit offset of the packed variable within `slot` from the left or right side.
    function findOffset(StdStorage storage self, bytes32 slot, bool left) internal returns (bool, uint256) {
        for (uint256 offset = 0; offset < 256; offset++) {
            uint256 valueToPut = left ? (1 << (255 - offset)) : (1 << offset);
            vm.store(self._target, slot, bytes32(valueToPut));

            (bool success, bytes32 data) = callTarget(self);

            if (success && (uint256(data) > 0)) {
                return (true, offset);
            }
        }
        return (false, 0);
    }

    /// @notice Returns whether both offsets were found, along with the left and right bit offsets of the packed variable within `slot`.
    function findOffsets(StdStorage storage self, bytes32 slot) internal returns (bool, uint256, uint256) {
        bytes32 prevSlotValue = vm.load(self._target, slot);

        (bool foundLeft, uint256 offsetLeft) = findOffset(self, slot, true);
        (bool foundRight, uint256 offsetRight) = findOffset(self, slot, false);

        // `findOffset` may mutate slot value, so we are setting it to initial value
        vm.store(self._target, slot, prevSlotValue);
        return (foundLeft && foundRight, offsetLeft, offsetRight);
    }

    /// @notice Finds the storage slot for the configured target and returns its data, clearing the configuration.
    function find(StdStorage storage self) internal returns (FindData storage) {
        return find(self, true);
    }

    /// @notice find an arbitrary storage slot given a function sig, input data, address of the contract and a value to check against
    // slot complexity:
    //  if flat, will be bytes32(uint256(uint));
    //  if map, will be keccak256(abi.encode(key, uint(slot)));
    //  if deep map, will be keccak256(abi.encode(key1, keccak256(abi.encode(key0, uint(slot)))));
    //  if map struct, will be bytes32(uint256(keccak256(abi.encode(key1, keccak256(abi.encode(key0, uint(slot)))))) + structFieldDepth);
    function find(StdStorage storage self, bool _clear) internal returns (FindData storage) {
        address who = self._target;
        bytes4 fsig = self._sig;
        uint256 field_depth = self._depth;
        bytes memory params = getCallParams(self);

        // calldata to test against
        if (self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))].found) {
            if (_clear) {
                clear(self);
            }
            return self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))];
        }
        vm.record();
        (, bytes32 callResult) = callTarget(self);
        (bytes32[] memory reads,) = vm.accesses(address(who));

        if (reads.length == 0) {
            revert("stdStorage find(StdStorage): No storage use detected for target.");
        } else {
            for (uint256 i = reads.length; i > 0;) {
                --i;
                bytes32 prev = vm.load(who, reads[i]);
                if (prev == bytes32(0)) {
                    emit WARNING_UninitedSlot(who, uint256(reads[i]));
                }

                if (!checkSlotMutatesCall(self, reads[i])) {
                    continue;
                }

                (uint256 offsetLeft, uint256 offsetRight) = (0, 0);

                if (self._enable_packed_slots) {
                    bool found;
                    (found, offsetLeft, offsetRight) = findOffsets(self, reads[i]);
                    if (!found) {
                        continue;
                    }
                }

                // Check that value between found offsets is equal to the current call result
                uint256 curVal = (uint256(prev) & getMaskByOffsets(offsetLeft, offsetRight)) >> offsetRight;

                if (uint256(callResult) != curVal) {
                    continue;
                }

                emit SlotFound(who, fsig, keccak256(abi.encodePacked(params, field_depth)), uint256(reads[i]));
                self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))] =
                    FindData(uint256(reads[i]), offsetLeft, offsetRight, true);
                break;
            }
        }

        require(
            self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))].found,
            "stdStorage find(StdStorage): Slot(s) not found."
        );

        if (_clear) {
            clear(self);
        }
        return self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))];
    }

    /// @notice Sets the target contract address for the storage lookup.
    function target(StdStorage storage self, address _target) internal returns (StdStorage storage) {
        self._target = _target;
        return self;
    }

    /// @notice Sets the target function selector for the storage lookup.
    function sig(StdStorage storage self, bytes4 _sig) internal returns (StdStorage storage) {
        self._sig = _sig;
        return self;
    }

    /// @notice Sets the target function selector from a signature string for the storage lookup.
    function sig(StdStorage storage self, string memory _sig) internal returns (StdStorage storage) {
        self._sig = sigs(_sig);
        return self;
    }

    /// @notice Sets raw calldata to use instead of ABI-encoded keys for the target call.
    function with_calldata(StdStorage storage self, bytes memory _calldata) internal returns (StdStorage storage) {
        self._calldata = _calldata;
        return self;
    }

    /// @notice Adds an address mapping key to the storage lookup path.
    function with_key(StdStorage storage self, address who) internal returns (StdStorage storage) {
        self._keys.push(bytes32(uint256(uint160(who))));
        return self;
    }

    /// @notice Adds a uint256 mapping key to the storage lookup path.
    function with_key(StdStorage storage self, uint256 amt) internal returns (StdStorage storage) {
        self._keys.push(bytes32(amt));
        return self;
    }

    /// @notice Adds a bytes32 mapping key to the storage lookup path.
    function with_key(StdStorage storage self, bytes32 key) internal returns (StdStorage storage) {
        self._keys.push(key);
        return self;
    }

    /// @notice Enables detection and handling of values packed into shared storage slots.
    function enable_packed_slots(StdStorage storage self) internal returns (StdStorage storage) {
        self._enable_packed_slots = true;
        return self;
    }

    /// @notice Sets the struct field depth for storage lookups into nested structs.
    function depth(StdStorage storage self, uint256 _depth) internal returns (StdStorage storage) {
        self._depth = _depth;
        return self;
    }

    function _read(StdStorage storage self) private returns (bytes memory) {
        FindData storage data = find(self, false);
        uint256 mask = getMaskByOffsets(data.offsetLeft, data.offsetRight);
        uint256 value = (uint256(vm.load(self._target, bytes32(data.slot))) & mask) >> data.offsetRight;
        clear(self);
        return abi.encode(value);
    }

    /// @notice Reads the found storage slot value as bytes32.
    function read_bytes32(StdStorage storage self) internal returns (bytes32) {
        return abi.decode(_read(self), (bytes32));
    }

    /// @notice Reads the found storage slot value as bool.
    /// @dev Reverts if the stored value is neither `0` nor `1`.
    function read_bool(StdStorage storage self) internal returns (bool) {
        int256 v = read_int(self);
        if (v == 0) return false;
        if (v == 1) return true;
        revert("stdStorage read_bool(StdStorage): Cannot decode. Make sure you are reading a bool.");
    }

    /// @notice Reads the found storage slot value as address.
    function read_address(StdStorage storage self) internal returns (address) {
        return abi.decode(_read(self), (address));
    }

    /// @notice Reads the found storage slot value as uint256.
    function read_uint(StdStorage storage self) internal returns (uint256) {
        return abi.decode(_read(self), (uint256));
    }

    /// @notice Reads the found storage slot value as int256.
    function read_int(StdStorage storage self) internal returns (int256) {
        return abi.decode(_read(self), (int256));
    }

    /// @notice Returns the parent mapping slot index and the key used to reach the found slot.
    function parent(StdStorage storage self) internal returns (uint256, bytes32) {
        address who = self._target;
        uint256 field_depth = self._depth;
        vm.startMappingRecording();
        uint256 child = find(self, true).slot - field_depth;
        (bool found, bytes32 key, bytes32 parent_slot) = vm.getMappingKeyAndParentOf(who, bytes32(child));
        if (!found) {
            revert(
                "stdStorage parent(StdStorage): Cannot find parent. Make sure you give a slot and startMappingRecording() has been called."
            );
        }
        return (uint256(parent_slot), key);
    }

    /// @notice Returns the root mapping slot index by traversing the mapping parent chain.
    function root(StdStorage storage self) internal returns (uint256) {
        address who = self._target;
        uint256 field_depth = self._depth;
        vm.startMappingRecording();
        uint256 child = find(self, true).slot - field_depth;
        bool found;
        bytes32 root_slot;
        bytes32 parent_slot;
        (found,, parent_slot) = vm.getMappingKeyAndParentOf(who, bytes32(child));
        if (!found) {
            revert(
                "stdStorage root(StdStorage): Cannot find parent. Make sure you give a slot and startMappingRecording() has been called."
            );
        }
        while (found) {
            root_slot = parent_slot;
            (found,, parent_slot) = vm.getMappingKeyAndParentOf(who, bytes32(root_slot));
        }
        return uint256(root_slot);
    }

    function _bytesToBytes32(bytes memory b, uint256 offset) private pure returns (bytes32) {
        bytes32 out;

        // Cap read length by remaining bytes from `offset`, and at most 32 bytes to avoid out-of-bounds
        uint256 max = b.length > offset ? b.length - offset : 0;
        if (max > 32) {
            max = 32;
        }
        for (uint256 i = 0; i < max; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }

    function _flatten(bytes32[] memory b) private pure returns (bytes memory) {
        bytes memory result = new bytes(b.length * 32);
        for (uint256 i = 0; i < b.length; i++) {
            bytes32 k = b[i];
            assembly ("memory-safe") {
                mstore(add(result, add(32, mul(32, i))), k)
            }
        }

        return result;
    }

    /// @notice Resets all configured parameters on `self`.
    function clear(StdStorage storage self) internal {
        delete self._target;
        delete self._sig;
        delete self._keys;
        delete self._depth;
        delete self._enable_packed_slots;
        delete self._calldata;
    }

    /// @notice Returns a bitmask with ones in the bit range `[offsetRight, 255 - offsetLeft]`.
    /// @dev `(slotValue & mask) >> offsetRight` extracts the packed variable's value.
    function getMaskByOffsets(uint256 offsetLeft, uint256 offsetRight) internal pure returns (uint256 mask) {
        // mask = ((1 << (256 - (offsetRight + offsetLeft))) - 1) << offsetRight;
        // using assembly because (1 << 256) causes overflow
        assembly {
            mask := shl(offsetRight, sub(shl(sub(256, add(offsetRight, offsetLeft)), 1), 1))
        }
    }

    /// @notice Returns `curValue` with the packed variable at `[offsetRight, 255 - offsetLeft]` replaced by `varValue`.
    function getUpdatedSlotValue(bytes32 curValue, uint256 varValue, uint256 offsetLeft, uint256 offsetRight)
        internal
        pure
        returns (bytes32 newValue)
    {
        return bytes32((uint256(curValue) & ~getMaskByOffsets(offsetLeft, offsetRight)) | (varValue << offsetRight));
    }
}

library stdStorage {
    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    /// @notice Returns the 4-byte function selector for `sigStr`.
    function sigs(string memory sigStr) internal pure returns (bytes4) {
        return stdStorageSafe.sigs(sigStr);
    }

    /// @notice Finds the storage slot index for the configured target, clearing the configuration.
    function find(StdStorage storage self) internal returns (uint256) {
        return find(self, true);
    }

    /// @notice Finds the storage slot index for the configured target, optionally clearing the configuration.
    function find(StdStorage storage self, bool _clear) internal returns (uint256) {
        return stdStorageSafe.find(self, _clear).slot;
    }

    /// @notice Sets the target contract address for the storage lookup.
    function target(StdStorage storage self, address _target) internal returns (StdStorage storage) {
        return stdStorageSafe.target(self, _target);
    }

    /// @notice Sets the target function selector for the storage lookup.
    function sig(StdStorage storage self, bytes4 _sig) internal returns (StdStorage storage) {
        return stdStorageSafe.sig(self, _sig);
    }

    /// @notice Sets the target function selector from a signature string for the storage lookup.
    function sig(StdStorage storage self, string memory _sig) internal returns (StdStorage storage) {
        return stdStorageSafe.sig(self, _sig);
    }

    /// @notice Adds an address mapping key to the storage lookup path.
    function with_key(StdStorage storage self, address who) internal returns (StdStorage storage) {
        return stdStorageSafe.with_key(self, who);
    }

    /// @notice Adds a uint256 mapping key to the storage lookup path.
    function with_key(StdStorage storage self, uint256 amt) internal returns (StdStorage storage) {
        return stdStorageSafe.with_key(self, amt);
    }

    /// @notice Adds a bytes32 mapping key to the storage lookup path.
    function with_key(StdStorage storage self, bytes32 key) internal returns (StdStorage storage) {
        return stdStorageSafe.with_key(self, key);
    }

    /// @notice Sets raw calldata to use instead of ABI-encoded keys for the target call.
    function with_calldata(StdStorage storage self, bytes memory _calldata) internal returns (StdStorage storage) {
        return stdStorageSafe.with_calldata(self, _calldata);
    }

    /// @notice Enables detection and handling of values packed into shared storage slots.
    function enable_packed_slots(StdStorage storage self) internal returns (StdStorage storage) {
        return stdStorageSafe.enable_packed_slots(self);
    }

    /// @notice Sets the struct field depth for storage lookups into nested structs.
    function depth(StdStorage storage self, uint256 _depth) internal returns (StdStorage storage) {
        return stdStorageSafe.depth(self, _depth);
    }

    /// @notice Resets all configured parameters on `self`.
    function clear(StdStorage storage self) internal {
        stdStorageSafe.clear(self);
    }

    /// @notice Writes `who` to the found storage slot and verifies the value was applied correctly.
    function checked_write(StdStorage storage self, address who) internal {
        checked_write(self, bytes32(uint256(uint160(who))));
    }

    /// @notice Writes `amt` to the found storage slot and verifies the value was applied correctly.
    function checked_write(StdStorage storage self, uint256 amt) internal {
        checked_write(self, bytes32(amt));
    }

    /// @notice Writes `val` to the found storage slot and verifies the value was applied correctly.
    function checked_write_int(StdStorage storage self, int256 val) internal {
        checked_write(self, bytes32(uint256(val)));
    }

    /// @notice Writes `write` to the found storage slot and verifies the value was applied correctly.
    function checked_write(StdStorage storage self, bool write) internal {
        bytes32 t;
        assembly ("memory-safe") {
            t := write
        }
        checked_write(self, t);
    }

    /// @notice Writes `set` to the found storage slot and verifies the value was applied correctly.
    function checked_write(StdStorage storage self, bytes32 set) internal {
        address who = self._target;
        bytes4 fsig = self._sig;
        uint256 field_depth = self._depth;
        bytes memory params = stdStorageSafe.getCallParams(self);

        if (!self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))].found) {
            find(self, false);
        }
        FindData storage data = self.finds[who][fsig][keccak256(abi.encodePacked(params, field_depth))];
        if ((data.offsetLeft + data.offsetRight) > 0) {
            uint256 maxVal = 2 ** (256 - (data.offsetLeft + data.offsetRight));
            require(
                uint256(set) < maxVal,
                string(
                    abi.encodePacked(
                        "stdStorage checked_write(StdStorage): Packed slot. We can't fit value greater than ",
                        vm.toString(maxVal)
                    )
                )
            );
        }
        bytes32 curVal = vm.load(who, bytes32(data.slot));
        bytes32 valToSet = stdStorageSafe.getUpdatedSlotValue(curVal, uint256(set), data.offsetLeft, data.offsetRight);

        vm.store(who, bytes32(data.slot), valToSet);

        (bool success, bytes32 callResult) = stdStorageSafe.callTarget(self);

        if (!success || callResult != set) {
            vm.store(who, bytes32(data.slot), curVal);
            revert("stdStorage checked_write(StdStorage): Failed to write value.");
        }
        clear(self);
    }

    /// @notice Reads the found storage slot value as bytes32.
    function read_bytes32(StdStorage storage self) internal returns (bytes32) {
        return stdStorageSafe.read_bytes32(self);
    }

    /// @notice Reads the found storage slot value as bool.
    function read_bool(StdStorage storage self) internal returns (bool) {
        return stdStorageSafe.read_bool(self);
    }

    /// @notice Reads the found storage slot value as address.
    function read_address(StdStorage storage self) internal returns (address) {
        return stdStorageSafe.read_address(self);
    }

    /// @notice Reads the found storage slot value as uint256.
    function read_uint(StdStorage storage self) internal returns (uint256) {
        return stdStorageSafe.read_uint(self);
    }

    /// @notice Reads the found storage slot value as int256.
    function read_int(StdStorage storage self) internal returns (int256) {
        return stdStorageSafe.read_int(self);
    }

    /// @notice Returns the parent mapping slot index and the key used to reach the found slot.
    function parent(StdStorage storage self) internal returns (uint256, bytes32) {
        return stdStorageSafe.parent(self);
    }

    /// @notice Returns the root mapping slot index by traversing the mapping parent chain.
    function root(StdStorage storage self) internal returns (uint256) {
        return stdStorageSafe.root(self);
    }
}
