// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

address constant GAS_BOUND_CALLER = address(0xc706EC7dfA5D4Dc87f29f859094165E8290530f5);

interface IGasBoundCaller {
  function gasBoundCall(address _to, uint256 _maxTotalGas, bytes calldata _data) external payable;
}
