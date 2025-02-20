// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {ISystemContext} from "../../vendor/@matter-labs/era-contracts/gas-bound-caller/contracts/ISystemContext.sol";

ISystemContext constant SYSTEM_CONTEXT_CONTRACT = ISystemContext(address(0x800b));

/**
 * @title CallWithExactGasZKSync
 * @notice Library that attempts to call a target contract with exactly `gasAmount` gas on zkSync
 *         and measures how much gas was actually used.
 * Implementation based on the GasBoundCaller contract, https://github.com/matter-labs/era-contracts/blob/main/gas-bound-caller/contracts/GasBoundCaller.sol
 */
library CallWithExactGasZKSync {
  error NoContract();
  error NotEnoughGasForPubdata();
  /// @notice We assume that no more than `CALL_RETURN_OVERHEAD` ergs are used for the O(1) operations at the end of the execution,
  /// as such relaying the return.
  uint256 internal constant CALL_RETURN_OVERHEAD = 400;

  bytes4 internal constant NO_CONTRACT_SIG = 0x0c3b563c;

  /// @notice The function that implements limiting of the total gas expenditure of the call.
  /// @dev On Era, the gas for pubdata is charged at the end of the execution of the entire transaction, meaning
  /// that if a subcall is not trusted, it can consume lots of pubdata in the process. This function ensures that
  /// no more than  `_maxTotalGas` will be allowed to be spent by the call. To be sure, this function uses some margin
  /// (`CALL_ENTRY_OVERHEAD` + `CALL_RETURN_OVERHEAD`) to ensure that the call will not exceed the limit, so it may
  /// actually spend a bit less than `_maxTotalGas` in the end.
  /// @dev The entire `gas` passed to this function could be used, regardless
  /// of the `_maxTotalGas` parameter. In other words, `max(gas(), _maxTotalGas)` is the maximum amount of gas that can be spent by this function.
  /// @dev The function relays the `returndata` returned by the callee. In case the `callee` reverts, it reverts with the same error.
  /// @param _to The address of the contract to call.
  /// @param _maxTotalGas the maximum amount of gas that can be spent by the call.
  /// @param _data The calldata for the call.
  /// @param _maxReturnBytes the maximum amount of bytes that can be returned by the call.
  /// @return success whether the call succeeded
  /// @return retData the return data from the call, capped at maxReturnBytes bytes
  /// @return pubdataGasSpent the pubdata gas used.
  function _callWithExactGasSafeReturnData(
    address _to,
    uint256 _maxTotalGas,
    bytes memory _data,
    uint16 _maxReturnBytes
  ) internal returns (bool success, bytes memory, uint256 pubdataGasSpent) {
    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      // Note we do this check prior to measuring gas.
      if iszero(extcodesize(_to)) {
        mstore(0x0, NO_CONTRACT_SIG)
        revert(0x0, 0x4)
      }
    }

    // We require that `_maxTotalGas` does not exceed the current `gasleft()`.
    // This is a safety check to ensure that a gas limit higher than the available gas is not specified,
    // which would indicate incorrect parameters and could lead to unexpected behavior.
    if (_maxTotalGas > gasleft()) {
      return (false, "", 0);
    }

    uint256 pubdataPublishedBefore = SYSTEM_CONTEXT_CONTRACT.getCurrentPubdataSpent();

    assembly {
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(_maxTotalGas, _to, 0, add(_data, 0x20), mload(_data), 0x0, 0x0)
    }
    bytes memory returnData = new bytes(_maxReturnBytes);
    assembly {
      // limit our copy to maxReturnBytes bytes
      let toCopy := returndatasize()
      if gt(toCopy, _maxReturnBytes) {
        toCopy := _maxReturnBytes
      }
      // Store the length of the copied bytes
      mstore(returnData, toCopy)
      // copy the bytes from retData[0:_toCopy]
      returndatacopy(add(returnData, 0x20), 0x0, toCopy)
    }

    uint256 pubdataPublishedAfter = SYSTEM_CONTEXT_CONTRACT.getCurrentPubdataSpent();

    // It is possible that pubdataPublishedAfter < pubdataPublishedBefore if the call, e.g. removes
    // some of the previously created state diffs
    uint256 pubdataSpent = pubdataPublishedAfter > pubdataPublishedBefore
      ? pubdataPublishedAfter - pubdataPublishedBefore
      : 0;

    uint256 pubdataGasRate = SYSTEM_CONTEXT_CONTRACT.gasPerPubdataByte();

    // In case there is an overflow here, the `_maxTotalGas` wouldn't be able to cover it anyway, so
    // we don't mind the contract panicking here in case of it.
    pubdataGasSpent = pubdataGasRate * pubdataSpent;
    if (pubdataGasSpent != 0) {
      // Here we double check that the additional cost is not higher than the maximum allowed.
      // Note, that the `gasleft()` can be spent on pubdata too.
      if (gasleft() < pubdataGasSpent + CALL_RETURN_OVERHEAD) {
        revert NotEnoughGasForPubdata();
      }
    }
    return (success, returnData, pubdataGasSpent);
  }
}
