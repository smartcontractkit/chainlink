// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

//import { KeeperCompatibleInterface } from "../interfaces/KeeperCompatibleInterface.sol";

contract DataBombContract {
    uint256 public testRange;
    uint256 public interval;
    uint256 public lastTimestamp;
    uint256 public previousPerformBlock;
    uint256 public initialTimestamp;
    uint256 public length;
    uint256 public size;
    mapping(uint256 => bytes) public data;

    function setConfig(uint256 _testRange, uint256 _interval, uint256 _length) external {
        testRange = _testRange;
        interval = _interval;
        lastTimestamp = block.timestamp;
        previousPerformBlock = 0;
        initialTimestamp = 0;
        length = _length;
    }

    function setLargeDataSize(uint256 _size) external {
        size = _size;
    }

    function checkUpkeep(bytes calldata data) external view returns (bool upkeepNeeded, bytes memory performData) {
        if (initialTimestamp == 0) {
            return (true, data);
        }

        return ((block.timestamp - initialTimestamp) < testRange && (block.timestamp - lastTimestamp) >= interval, "");
    }

    function performUpkeep(bytes calldata /* performData */) external returns (bytes memory) {
        for (uint256 i = 0; i < length; i++) {
            data[i] = new bytes(length);
        }

        return new bytes(size);
    }
}