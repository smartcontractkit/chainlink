// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../onRamp/EVM2EVMMultiOnRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract EVM2EVMMultiOnRampHelper is EVM2EVMMultiOnRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig
  ) EVM2EVMMultiOnRamp(staticConfig, dynamicConfig) {}
}
