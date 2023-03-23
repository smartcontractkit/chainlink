pragma solidity ^0.7.0;

import "./KeeperCompatibleInterface.sol";
import "./KeeperBase.sol";

contract KeeperConsumer is KeeperCompatibleInterface, KeeperBase {
    uint public counter;
    uint public immutable interval;
    uint public lastTimeStamp;


    constructor(uint updateInterval) public {
        interval = updateInterval;
        lastTimeStamp = block.timestamp;
        counter = 0;
    }

    function checkUpkeep(bytes calldata checkData) 
    external 
    override 
    view
    cannotExecute
    returns (bool upkeepNeeded, bytes memory performData) {
        return (true, checkData);
    }

    function performUpkeep(bytes calldata performData) external override {
        counter = counter + 1;
    }
}

