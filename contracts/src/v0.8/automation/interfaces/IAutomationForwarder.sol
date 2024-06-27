// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IAutomationRegistryConsumer} from "./IAutomationRegistryConsumer.sol";

interface IAutomationForwarder is ITypeAndVersion {
  function forward(uint256 gasAmount, bytes memory data) external returns (bool success, uint256 gasUsed);

  function updateRegistry(address newRegistry) external;

  function getRegistry() external view returns (IAutomationRegistryConsumer);

  function getTarget() external view returns (address);
}
