// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

contract MockWorkflowRegistryContract is ITypeAndVersion {
  string public constant override typeAndVersion = "MockWorkflowRegistryContract 1.0.0";
}
