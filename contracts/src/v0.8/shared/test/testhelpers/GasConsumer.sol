pragma solidity ^0.8.0;

contract GasConsumer {
  function consumeAllGas() external view {
    assembly {
      for {

      } lt(0, 1) {

      } {
        //This is as close as we can get to the amount of gas needed to safely consume all and return
        if lt(gas(), 30) {
          pop(add(0, 0))
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
