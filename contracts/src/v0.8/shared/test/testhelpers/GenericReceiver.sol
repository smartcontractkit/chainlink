// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract GenericReceiver {
  bool public s_toRevert;
  bytes private s_err;

  constructor(bool toRevert) {
    s_toRevert = toRevert;
  }

  function setRevert(bool toRevert) external {
    s_toRevert = toRevert;
  }

  function setErr(bytes memory err) external {
    s_err = err;
  }

  // solhint-disable-next-line payable-fallback
  fallback() external {
    if (s_toRevert) {
      bytes memory reason = s_err;
      assembly {
        revert(add(32, reason), mload(reason))
      }
    }
  }
}

contract GasConsumer {
  function consumeAllGas() external view {
    assembly {
      for {} lt(0, 1) {} {
        //This is as close as I can get to the amount of gas needed to safely consume all and return
        if lt(gas(), 30) {
          pop(add(0,0))
          return(0x0, 0x0)
        }
      }
    }
  }

  function throwOutOfGasError() external pure {
    while (true) {
      //Intentionally consume all gas and throw an OOG error.
    }
  }
}