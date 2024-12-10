// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {SYSTEM_CONTEXT, ISystemContext} from "../interfaces/zksync/ISystemContext.sol";

/// @dev A library that abstracts out opcodes that behave differently across chains.
/// @dev The methods below return values that are pertinent to the given chain.
/// @dev The exisiting library used within other functions contracts contains opcodes not supported by ZKSync
library ZKSyncChainSpecificUtil {
  function _getCurrentTxL1GasFees() internal view returns (uint256 l1FeeWei) {
    ISystemContext systemContext = ISystemContext(SYSTEM_CONTEXT);

    // blob_gas_price_on_l1 * gas_per_byte ~= gas_price_on_l2 * l2_gas_PerPubdataByte
    return systemContext.gasPrice() * systemContext.gasPerPubdataByte();
  }
}
