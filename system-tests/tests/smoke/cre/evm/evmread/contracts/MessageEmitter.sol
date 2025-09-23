// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract MessageEmitter {
    event MessageEmitted(string message);

    function emitMessage(string calldata message) external {
        emit MessageEmitted(message);
    }

    function getMessage(string memory str) public pure returns (string memory) {
        return string(abi.encodePacked("getMessage returns: ", str));
    }

    function onReport(bytes calldata, bytes calldata report) external {
         emit MessageEmitted(string(report));
    }

}

