// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ILogAutomation, Log} from "../2_1/interfaces/ILogAutomation.sol";
import "../2_1/interfaces/FeedLookupCompatibleInterface.sol";
import {ArbSys} from "../../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";

contract LogTriggeredFeedLookup is ILogAutomation, FeedLookupCompatibleInterface {

    ArbSys internal constant ARB_SYS = ArbSys(0x0000000000000000000000000000000000000064);

    function checkLog(Log calldata log) external override returns (bool upkeepNeeded, bytes memory performData) {

    }

    function performUpkeep(bytes calldata performData) external override {

    }

    function checkCallback(
        bytes[] memory values,
        bytes memory extraData
    ) external view returns (bool upkeepNeeded, bytes memory performData) {

    }
}