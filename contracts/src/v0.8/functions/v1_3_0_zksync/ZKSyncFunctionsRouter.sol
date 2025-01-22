// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../v1_0_0/FunctionsRouter.sol";
import {CallWithExactGasZKSync} from "../../shared/call/CallWithExactGasZKSync.sol";

/**
 * @title FunctionsRouterZkSync
 * @notice Specialized version of FunctionsRouter for zkSync that uses
 *         CallWithExactGasZKSync to control callback gas usage.
 */
contract ZKSyncFunctionsRouter is FunctionsRouter {
  constructor(address linkToken, Config memory config) FunctionsRouter(linkToken, config) {}

  /**
   * @dev Override the internal callback function to use CallWithExactGasZKSync
   *      for controlling and measuring gas usage on zkSync.
   */
  function _callback(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint32 callbackGasLimit,
    address client
  ) internal override returns (CallbackResult memory) {
    // 1. Check if client code exists
    bool destinationNoLongerExists;
    assembly {
      destinationNoLongerExists := iszero(extcodesize(client))
    }
    if (destinationNoLongerExists) {
      // If there's no code at `client`, skip the callback
      return CallbackResult({success: false, gasUsed: 0, returnData: new bytes(0)});
    }

    // 2. Encode callback to the client
    bytes memory encodedCallback = abi.encodeWithSelector(
      this.getConfig().handleOracleFulfillmentSelector,
      requestId,
      response,
      err
    );

    // 3. Use our library to enforce an exact gas call
    (bool success, uint256 gasUsed, bytes memory returnData) = CallWithExactGasZKSync._callWithExactGasSafeReturnData(
      client,
      callbackGasLimit,
      encodedCallback,
      MAX_CALLBACK_RETURN_BYTES
    );
    return CallbackResult({success: success, gasUsed: gasUsed, returnData: returnData});
  }
}
