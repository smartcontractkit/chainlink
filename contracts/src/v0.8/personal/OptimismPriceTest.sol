pragma solidity 0.8.6;

import "../vendor/@eth-optimism/contracts/0.8.6/contracts/L2/predeploys/OVM_GasPriceOracle.sol";

contract OptimismPriceTest {
  OVM_GasPriceOracle public immutable OPTIMISM_ORACLE = OVM_GasPriceOracle(0x420000000000000000000000000000000000000F);

  uint256 public gas0;
  uint256 public gas1;
  uint256 public l1CostWei;

  function test(bytes calldata performData) external {
    gas0 = gasleft();

    l1CostWei = OPTIMISM_ORACLE.getL1Fee(performData);

    gas1 = gasleft();
  }
}
