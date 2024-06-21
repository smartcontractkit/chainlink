// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import "../../shared/access/ConfirmedOwner.sol";
import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {AccessControl} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/AccessControl.sol";
import {LinkAvailableBalanceMonitor} from "../upkeeps/LinkAvailableBalanceMonitor.sol";
import "../interfaces/ILogAutomation.sol";

contract CCIPOnRampAutoWatchlist is ILogAutomation, ConfirmedOwner {
  event setWatchlistOnMonitor(uint64 indexed _dstChainSelector, address indexed _onRamp);

  LinkAvailableBalanceMonitor public linkAvailableBalanceMonitorAddress;
  address public routerAddress;

  // keccak256 signature for OnRampSet event set as constant
  bytes32 private constant EVENT_SIGNATURE = 0x1f7d0ec248b80e5c0dde0ee531c4fc8fdb6ce9a2b3d90f560c74acd6a7202f23;

  constructor(address contractaddress, address _routerAddress) ConfirmedOwner(msg.sender) {
    linkAvailableBalanceMonitorAddress = LinkAvailableBalanceMonitor(contractaddress);
    routerAddress = _routerAddress;
  }

  // Option to change the Balance Monitor Address
  function setBalanceMonitorAddress(address _newBalMonAddress) external onlyOwner {
    linkAvailableBalanceMonitorAddress = LinkAvailableBalanceMonitor(_newBalMonAddress);
  }

  // Option to change the Router Address
  // Changing Router Address would also require setting up new Upkeep as Upkeep is linked with the Router and cannot be changed
  function setRouterAddress(address _newRouterAddress) external onlyOwner {
    routerAddress = _newRouterAddress;
  }

  function updateWatchList(address _targetAddress, uint64 _dstChainSelector) internal {
    linkAvailableBalanceMonitorAddress.addToWatchListOrDecommission(_targetAddress, _dstChainSelector);
  }

  function checkLog(
    Log calldata log,
    bytes memory checkData
  ) external view override returns (bool upkeepNeeded, bytes memory performData) {
    // Ensure Router address is set
    require(routerAddress != address(0), "Router address not set");

    // Define the event signature for OnRampSet(uint64,address)
    //bytes32 eventSignature = keccak256(abi.encodePacked("OnRampSet(uint64,address)"));

    // Check if the log source matches router contract and topics contain the event signature
    if (log.source == routerAddress && log.topics.length > 0 && log.topics[0] == EVENT_SIGNATURE) {
      // Extract the indexed parameter from the log
      uint64 destChainSelector = uint64(uint256(log.topics[1])); // cast to uint64
      address onRamp = address(uint160(uint256(bytes32(log.data)))); // extract from log.data

      // No sanity checks necessary as the would on to relay information to LinkAvailableBalanceMonitor as it is
      // Checking is enabled in LinkAvailableBalanceMonitor contract
      return (true, abi.encode(destChainSelector, onRamp));
    }

    // If the event signature doesn't match or log source is not Router contract, no upkeep is needed
    return (false, "");
  }

  function performUpkeep(bytes memory performData) external override {
    // Decode the data received from checkLog
    (uint64 destChainSelector, address onRamp) = abi.decode(performData, (uint64, address));
    // Perform the necessary upkeep actions based on the decoded data
    updateWatchList(onRamp, destChainSelector);
    emit setWatchlistOnMonitor(destChainSelector, onRamp);
  }
}
