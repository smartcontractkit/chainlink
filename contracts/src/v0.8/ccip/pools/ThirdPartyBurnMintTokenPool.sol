// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBurnMintERC20} from "../../shared/token/ERC20/IBurnMintERC20.sol";

import {BurnMintTokenPool} from "./BurnMintTokenPool.sol";
import {Router} from "../Router.sol";

/// @notice This pool mints and burns a 3rd-party token. This pool is not owned by the DON
// and therefor has an additional check on adding offRamps.
contract ThirdPartyBurnMintTokenPool is BurnMintTokenPool {
  error InvalidOffRamp(address offRamp);

  /// @notice the trusted Router address to validate new offRamps through.
  address private s_router;

  constructor(
    IBurnMintERC20 token,
    address[] memory allowlist,
    address router,
    address armProxy
  ) BurnMintTokenPool(token, allowlist, armProxy) {
    s_router = router;
  }

  /// @notice Sets permissions for all on and offRamps.
  /// @dev Only callable by the owner
  /// @param onRamps A list of onRamps and their new permission status
  /// @param offRamps A list of offRamps and their new permission status
  function applyRampUpdates(RampUpdate[] calldata onRamps, RampUpdate[] calldata offRamps) external override onlyOwner {
    // Sanity check the offramps are enabled
    for (uint256 i = 0; i < offRamps.length; ++i) {
      // If the offRamp is being added do an additional check if the offRamp is
      // permission by the router. If not, we revert because we tried to add an
      // invalid offRamp.
      (bool exists, ) = Router(s_router).isOffRamp(offRamps[i].ramp);
      if (!exists) revert InvalidOffRamp(offRamps[i].ramp);
    }
    _applyRampUpdates(onRamps, offRamps);
  }
}
