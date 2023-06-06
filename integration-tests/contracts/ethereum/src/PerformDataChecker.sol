// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./KeeperCompatibleInterface.sol";

contract PerformDataChecker is KeeperCompatibleInterface {
    uint256 public counter;
    bytes public s_expectedData;

    constructor(bytes memory expectedData) {
        s_expectedData = expectedData;
    }

    function setExpectedData(bytes calldata expectedData) external {
        s_expectedData = expectedData;
    }

    function checkUpkeep(bytes calldata checkData)
    external
    view
    override
    returns (bool upkeepNeeded, bytes memory performData)
    {
        return (keccak256(checkData) == keccak256(s_expectedData), checkData);
    }

    function performUpkeep(bytes calldata performData) external override {
        if (keccak256(performData) == keccak256(s_expectedData)) {
            counter++;
        }
    }
}
