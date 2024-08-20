// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import "../../onRamp/OnRamp.sol";
import {IgnoreContractSize} from "./IgnoreContractSize.sol";

contract OnRampHelper is OnRamp, IgnoreContractSize {
  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    DestChainConfigArgs[] memory destChainConfigArgs
  ) OnRamp(staticConfig, dynamicConfig, destChainConfigArgs) {}
}
