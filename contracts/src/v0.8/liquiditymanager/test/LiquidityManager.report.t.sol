// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import "../interfaces/ILiquidityManager.sol";
import "./LiquidityManagerSetup.t.sol";

contract LiquidityManager_report is LiquidityManagerSetup {
  function test_report_RevertWhen_EmptyReport() external {
    ILiquidityManager.LiquidityInstructions memory instructions = ILiquidityManager.LiquidityInstructions({
      sendLiquidityParams: new ILiquidityManager.SendLiquidityParams[](0),
      receiveLiquidityParams: new ILiquidityManager.ReceiveLiquidityParams[](0)
    });

    vm.expectRevert(LiquidityManager.EmptyReport.selector);

    s_liquidityManager.report(abi.encode(instructions), 123);
  }
}
