// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_validateDestFamilyAddress is FeeQuoterSetup {
  function test_ValidEVMAddress_Success() public view {
    bytes memory encodedAddress = abi.encode(address(10000));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, encodedAddress);
  }

  function test_ValidNonEVMAddress_Success() public view {
    s_feeQuoter.validateDestFamilyAddress(bytes4(uint32(1)), abi.encode(type(uint208).max));
  }

  // Reverts

  function test_InvalidEVMAddress_Revert() public {
    bytes memory invalidAddress = abi.encode(type(uint208).max);
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressEncodePacked_Revert() public {
    bytes memory invalidAddress = abi.encodePacked(address(234));
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
  }

  function test_InvalidEVMAddressPrecompiles_Revert() public {
    for (uint160 i = 0; i < Internal.PRECOMPILE_SPACE; ++i) {
      bytes memory invalidAddress = abi.encode(address(i));
      vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
      s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress);
    }

    s_feeQuoter.validateDestFamilyAddress(
      Internal.CHAIN_FAMILY_SELECTOR_EVM, abi.encode(address(uint160(Internal.PRECOMPILE_SPACE)))
    );
  }
}
