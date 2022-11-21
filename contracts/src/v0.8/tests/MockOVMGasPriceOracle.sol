pragma solidity 0.8.6;

contract MockOVMGasPriceOracle {
  function getL1Fee(bytes memory _data) public view returns (uint256) {
    return 2000000;
  }
}
