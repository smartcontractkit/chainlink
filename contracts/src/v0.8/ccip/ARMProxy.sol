// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IARM} from "./interfaces/IARM.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";

contract ARMProxy is OwnerIsCreator, TypeAndVersionInterface {
  error ZeroAddressNotAllowed();

  event ARMSet(address arm);

  // STATIC CONFIG
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ARMProxy 1.0.0";

  // DYNAMIC CONFIG
  address private s_arm;

  constructor(address arm) {
    setARM(arm);
  }

  /// @notice SetARM sets the ARM implementation contract address.
  /// @param arm The address of the arm implementation contract.
  function setARM(address arm) public onlyOwner {
    if (arm == address(0)) revert ZeroAddressNotAllowed();
    s_arm = arm;
    emit ARMSet(arm);
  }

  /// @notice getARM gets the ARM implementation contract address.
  /// @return arm The address of the arm implementation contract.
  function getARM() external view returns (address) {
    return s_arm;
  }

  // We use a fallback function instead of explicit implementations of the functions
  // defined in IARM.sol to preserve compatibility with future additions to the IARM
  // interface. Calling IARM interface methods in ARMProxy should be transparent, i.e.
  // their input/output behaviour should be identical to calling the proxied s_arm
  // contract directly. (If s_arm doesn't point to a contract, we always revert.)
  fallback() external {
    address arm = s_arm;
    assembly {
      // Revert if no contract present at destination address, otherwise call
      // might succeed unintentionally.
      if iszero(extcodesize(arm)) {
        revert(0, 0)
      }
      // We use memory starting at zero, overwriting anything that might already
      // be stored there. This messes with Solidity's expectations around memory
      // layout, but it's fine because we always exit execution of this contract
      // inside this assembly block, i.e. we don't cede control to code generated
      // by the Solidity compiler that might have expectations around memory
      // layout.
      // Copy calldatasize() bytes from calldata offset 0 to memory offset 0.
      calldatacopy(0, 0, calldatasize())
      // Call the underlying ARM implementation. out and outsize are 0 because
      // we don't know the size yet. We hardcode value to zero.
      let success := call(gas(), arm, 0, 0, calldatasize(), 0, 0)
      // Copy the returned data.
      returndatacopy(0, 0, returndatasize())
      // Pass through successful return or revert and associated data.
      if success {
        return(0, returndatasize())
      }
      revert(0, returndatasize())
    }
  }
}
