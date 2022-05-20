// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/KeeperCompatibleInterface.sol";
import "../interfaces/KeeperRegistryInterface.sol";

contract CanaryUpkeep is KeeperCompatibleInterface {

    uint private s_keeperIndex;
    uint private s_interval = 300;
    uint private s_timestamp;
    KeeperRegistryInterface private s_keeperRegistry;

    constructor(KeeperRegistryInterface keeperRegistry) {
        s_keeperRegistry = keeperRegistry;
        s_timestamp = block.timestamp;
        s_keeperIndex = 0;
    }

    function getKeeperIndex() external view returns(uint) {
        return s_keeperIndex;
    }

    function getTimestamp() external view returns(uint) {
        return s_timestamp;
    }

    function getInterval() external view returns(uint) {
        return s_interval;
    }

    function setInterval(uint interval) external {
        s_interval = interval;
    }

    function checkUpkeep(bytes calldata checkData) external view returns (bool upkeepNeeded, bytes memory performData) {
        (State memory _s, Config memory _c, address[] memory keepers) = s_keeperRegistry.getState();
        upkeepNeeded = keepers.length != 0 && block.timestamp >= s_interval * 1 seconds + s_timestamp;
        // tx.origin will return 0 in simulated transactions => simulated transactions are often sent to read only functions
    }

    function performUpkeep(bytes calldata performData) external {
        (State memory _s, Config memory _c, address[] memory keepers) = s_keeperRegistry.getState();
        if (keepers.length == 0) {
            revert("no keeper nodes exists");
        }
        if (block.timestamp < s_interval * 1 seconds + s_timestamp) {
            revert("Not enough time has passed after the previous upkeep");
        }
        if (s_keeperIndex >= keepers.length) {
            s_keeperIndex = 0;
        }

        require(tx.origin == keepers[s_keeperIndex], "transaction origin is not the anticipated keeper.");
        s_keeperIndex++;
        s_timestamp = block.timestamp;
    }
}
