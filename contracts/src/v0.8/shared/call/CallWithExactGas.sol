// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @dev Handles the edge case where we want to pass a specific amount of gas,
/// @dev but EIP-150 sends all but 1/64 of the remaining gas instead so the user gets
/// @dev less gas than they paid for. The other 2 parts of EIP-150 do not apply since
/// @dev a) we hard code value=0 and b) we ensure code already exists.
/// @dev If we revert instead, then that will never happen.
/// @dev Separately we capture the return data up to a maximum size to avoid return bombs,
/// @dev borrowed from https://github.com/nomad-xyz/ExcessivelySafeCall/blob/main/src/ExcessivelySafeCall.sol.
library CallWithExactGas {
  error NoContract();
  error NoGasForCallExactCheck();
  error NotEnoughGasForCall();

  function _callWithExactGasSafeReturnData(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck,
    uint16 maxReturnBytes
  ) internal returns (bool success, bytes memory retData) {
    // allocate retData memory ahead of time
    retData = new bytes(maxReturnBytes);

    bytes4 noContract = NoContract.selector;
    bytes4 noGasForCallExactCheck = NoGasForCallExactCheck.selector;
    bytes4 notEnoughGasForCall = NotEnoughGasForCall.selector;

    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      // Note we do this check prior to measuring gas so gasForCallExactCheck (our "cushion")
      // doesn't need to account for it.
      if iszero(extcodesize(target)) {
        mstore(0, noContract)
        revert(0, 0x4)
      }

      let g := gas()
      // Compute g -= gasForCallExactCheck and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas. gasForCallExactCheck ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, gasForCallExactCheck) {
        mstore(0, noGasForCallExactCheck)
        revert(0, 0x4)
      }
      g := sub(g, gasForCallExactCheck)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasLimit)) {
        mstore(0, notEnoughGasForCall)
        revert(0, 0x4)
      }

      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasLimit, target, 0, add(payload, 0x20), mload(payload), 0, 0)

      // limit our copy to maxReturnBytes bytes
      let toCopy := returndatasize()
      if gt(toCopy, maxReturnBytes) {
        toCopy := maxReturnBytes
      }
      // Store the length of the copied bytes
      mstore(retData, toCopy)
      // copy the bytes from retData[0:_toCopy]
      returndatacopy(add(retData, 0x20), 0, toCopy)
    }
    return (success, retData);
  }

  function _callWithExactGas(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck
  ) internal returns (bool success) {
    bytes4 noContract = NoContract.selector;
    bytes4 noGasForCallExactCheck = NoGasForCallExactCheck.selector;
    bytes4 notEnoughGasForCall = NotEnoughGasForCall.selector;

    assembly {
      // solidity calls check that a contract actually exists at the destination, so we do the same
      // Note we do this check prior to measuring gas so gasForCallExactCheck (our "cushion")
      // doesn't need to account for it.
      if iszero(extcodesize(target)) {
        mstore(0, noContract)
        revert(0, 0x4)
      }

      let g := gas()
      // Compute g -= gasForCallExactCheck and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas. gasForCallExactCheck ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, gasForCallExactCheck) {
        mstore(0, noGasForCallExactCheck)
        revert(0, 0x4)
      }
      g := sub(g, gasForCallExactCheck)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasLimit)) {
        mstore(0, notEnoughGasForCall)
        revert(0, 0x4)
      }

      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasLimit, target, 0, add(payload, 0x20), mload(payload), 0, 0)
    }
    return success;
  }

  function _callWithExactGasEvenIfTargetIsNoContract(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck
  ) internal returns (bool sufficientGas) {
    assembly {
      let g := gas()
      // Compute g -= CALL_WITH_EXACT_GAS_CUSHION and check for underflow. We
      // need the cushion since the logic following the above call to gas also
      // costs gas which we cannot account for exactly. So cushion is a
      // conservative upper bound for the cost of this logic.
      if iszero(lt(g, gasForCallExactCheck)) {
        g := sub(g, gasForCallExactCheck)
        // If g - g//64 <= gasAmount, we don't have enough gas. (We subtract g//64
        // because of EIP-150.)
        if gt(sub(g, div(g, 64)), gasLimit) {
          // Call and ignore success/return data. Note that we did not check
          // whether a contract actually exists at the target address.
          pop(call(gasLimit, target, 0, add(payload, 0x20), mload(payload), 0, 0))
          sufficientGas := true
        }
      }
    }
    return sufficientGas;
  }
}
