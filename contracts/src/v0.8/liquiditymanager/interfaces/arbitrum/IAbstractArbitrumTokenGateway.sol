// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {TokenGateway} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/libraries/gateway/TokenGateway.sol";

/// @dev to generate gethwrappers
abstract contract IAbstractArbitrumTokenGateway is TokenGateway {}
