pragma solidity ^0.8.0;

contract GasConsumer {
  function consumeAllGas() external view {
    assembly {

      // While loop that operates indefinitely, written in yul to ensure better granularity over exactly how much gas is spent
      for {

      // Loop will run forever since 0 < 1
      } lt(0, 1) {

      } {
        // If 30 gas is remaining, then exit the loop by returning
        if lt(gas(), 30) {
          pop(add(0, 0)) // Add two numbers but don't push result onto the stack. Safely consume any residual gas from an odd number remaining.
          return(0x0, 0x0) // Return with no return data
        }
      }
    }
  }

  function throwOutOfGasError() external pure {
    while (true) {
      // Intentionally consume all gas to throw an OOG error.
    }
  }
}
