// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "./KinSwarmConsumer.sol";

contract DeployAndSettle is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // 1. Deploy to HyperEVM / Local
        address fakeOracle = 0x0000000000000000000000000000000000000001;
        KinSwarmConsumer consumer = new KinSwarmConsumer(fakeOracle);
        
        console.log("Consumer Deployed at:", address(consumer));

        // 2. Anchor our 10M Verified Root
        bytes32 root = 0x6a99f4d0755e0ce9dba8afb3f3bde5c0c23a364ad47e886ebbaeca8ba75914b2;
        uint256 volume = 350000000; // $350M USD
        
        consumer.anchorSovereignRoot(root, volume);
        
        console.log("Settlement Anchored Successfully.");
        console.log("Total Volume in Contract:", consumer.totalSettledVolume());

        vm.stopBroadcast();
    }
}
