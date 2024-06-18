// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import "../../shared/access/ConfirmedOwner.sol";
import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {AccessControl} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/AccessControl.sol";
import {LinkAvailableBalanceMonitor} from "../upkeeps/LinkAvailableBalanceMonitor.sol";


struct Log {
    uint256 index; // Index of the log in the block
    uint256 timestamp; // Timestamp of the block containing the log
    bytes32 txHash; // Hash of the transaction containing the log
    uint256 blockNumber; // Number of the block containing the log
    bytes32 blockHash; // Hash of the block containing the log
    address source; // Address of the contract that emitted the log
    bytes32[] topics; // Indexed topics of the log
    bytes data; // Data of the log
}

interface ILogAutomation {
    function checkLog(
        Log calldata log,
        bytes memory checkData
    ) external returns (bool upkeepNeeded, bytes memory performData);

    function performUpkeep(bytes calldata performData) external;
}

 contract CCIPOnRampAutoWatchlist is ILogAutomation,ConfirmedOwner {


    event setWatchlistOnMonitor(uint64 indexed _dstChainSelector,address indexed _onRamp);

    LinkAvailableBalanceMonitor public LinkAvailableBalanceMonitorAddress;
    address public routerAddress;


    constructor(address contractaddress,address _routerAddress) ConfirmedOwner(msg.sender){
        LinkAvailableBalanceMonitorAddress = LinkAvailableBalanceMonitor(contractaddress);
        routerAddress = _routerAddress;
    }

    // Option to change the Balance Monitor Address
    function setBalanceMonitorAddress(address _newBalMonAddress) external onlyOwner{
        LinkAvailableBalanceMonitorAddress = LinkAvailableBalanceMonitor(_newBalMonAddress);
    }

    // Option to change the Router Address
    // Changing Router Address would also require setting up new Upkeep as Upkeep is linked with the Router and cannot be changed
    function setRouterAddress(address _newRouterAddress) external onlyOwner{
        routerAddress = _newRouterAddress;
    }

    function updateWatchList(address _targetAddress, uint64 _dstChainSelector) internal{
        LinkAvailableBalanceMonitorAddress.addToWatchListOrDecommission(_targetAddress, _dstChainSelector);
    }

    function checkLog(
    Log calldata log,
    bytes memory checkData
) external view override returns (bool upkeepNeeded, bytes memory performData) {
    
    // Ensure Router address is set
    require(routerAddress != address(0), "Router address not set");

    // Define the event signature for OnRampSet(uint64,address)
    bytes32 eventSignature = keccak256(abi.encodePacked("OnRampSet(uint64,address)"));

    // Check if the log source matches router contract and topics contain the event signature
    if (log.source == routerAddress && log.topics.length > 0 && log.topics[0] == eventSignature) {
        // Extract the indexed parameter from the log
        uint64 destChainSelector = uint64(uint256(log.topics[1])); // cast to uint64
        address onRamp = address(uint160(uint256(bytes32(log.data)))); // extract from log.data

        // Determine if upkeep is needed based on the emitted log
        // For simplicity, always return true to trigger performUpkeep
        return (true, abi.encode(destChainSelector, onRamp));
    }

    // If the event signature doesn't match or log source is not Router contract, no upkeep is needed
    return (false, "");
}


    function performUpkeep(bytes memory performData) external override{
        // Decode the data received from checkLog
        (uint64 destChainSelector, address onRamp) = abi.decode(performData, (uint64, address));
        // Perform the necessary upkeep actions based on the decoded data
        updateWatchList(onRamp,destChainSelector);
        emit setWatchlistOnMonitor(destChainSelector,onRamp);
    }
}