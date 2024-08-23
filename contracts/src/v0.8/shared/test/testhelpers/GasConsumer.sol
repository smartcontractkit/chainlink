pragma solidity ^0.8.0;

contract GasConsumer {
  function consumeAllGas() external view {
    assembly {
      // While loop that operates indefinitely, written in yul to ensure better granularity over exactly how much gas is spent
      for {
        // Loop will run forever since 0 < 1
      } lt(0, 1) {

      } {
        // If 60 gas is remaining, then exit the loop by returning. 60 was determined by manual binary search to be the minimal amount of gas needed but less than the cost of another loop
        if lt(gas(), 60) {
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
