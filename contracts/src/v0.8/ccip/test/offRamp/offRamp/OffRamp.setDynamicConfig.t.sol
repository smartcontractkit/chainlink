// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_setDynamicConfig is OffRampSetup {
  function test_SetDynamicConfig_Success() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    OffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  function test_SetDynamicConfigWithInterceptor_Success() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));
    dynamicConfig.messageInterceptor = address(s_inboundMessageInterceptor);

    vm.expectEmit();
    emit OffRamp.DynamicConfigSet(dynamicConfig);

    s_offRamp.setDynamicConfig(dynamicConfig);

    OffRamp.DynamicConfig memory newConfig = s_offRamp.getDynamicConfig();
    _assertSameConfig(dynamicConfig, newConfig);
  }

  // Reverts

  function test_NonOwner_Revert() public {
    vm.startPrank(STRANGER);
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(s_feeQuoter));

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_offRamp.setDynamicConfig(dynamicConfig);
  }

  function test_FeeQuoterZeroAddress_Revert() public {
    OffRamp.DynamicConfig memory dynamicConfig = _generateDynamicOffRampConfig(address(0));

    vm.expectRevert(OffRamp.ZeroAddressNotAllowed.selector);

    s_offRamp.setDynamicConfig(dynamicConfig);
  }
}
