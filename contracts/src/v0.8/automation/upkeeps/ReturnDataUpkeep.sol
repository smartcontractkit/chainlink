// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract ReturnDataUpkeep {
    uint256 public testRange = 1000000;
    uint256 public interval = 100;
    uint256 public lastTimestamp;
    uint256 public previousPerformBlock;
    uint256 public initialTimestamp;
    uint256 public length = 300;
    uint256 public size = 999999;
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

    function checkUpkeep(bytes calldata _data) external view returns (bool upkeepNeeded, bytes memory performData) {
        if (initialTimestamp == 0) {
            return (true, _data);
        }

        return ((block.timestamp - initialTimestamp) < testRange && (block.timestamp - lastTimestamp) >= interval, "");
    }

    function performUpkeep(bytes calldata /* performData */) external returns (bytes memory) {
        for (uint256 i = 0; i < length; i++) {
            data[i] = new bytes(length);
        }

        lastTimestamp = block.timestamp;
        previousPerformBlock = block.number;

        return new bytes(size);
    }
}