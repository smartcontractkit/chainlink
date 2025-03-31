// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../FeeQuoter.sol";
import {Internal} from "../../libraries/Internal.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_validateDestFamilyAddress is FeeQuoterSetup {
  function test_validateDestFamilyAddress_EVMs() public view {
    bytes memory encodedAddress = abi.encode(address(10000));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, encodedAddress, 0);
  }

  function test_validateDestFamilyAddress_SVM() public view {
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_SVM, abi.encode(type(uint208).max), 0);
  }

  function test_validateDestFamilyAddress_Aptos() public view {
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_APTOS, abi.encode(type(uint208).max), 0);
  }

  // Reverts
  function test_validateDestFamilyAddress_RevertWhen_InvalidChainFamilySelector() public {
    bytes4 selector = bytes4(0xdeadbeef);
    bytes memory encodedAddress = abi.encode(address(10000));
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.InvalidChainFamilySelector.selector, selector));
    s_feeQuoter.validateDestFamilyAddress(selector, encodedAddress, 0);
  }

  function test_validateDestFamilyAddress_EVM_RevertWhen_InvalidEVMAddress() public {
    bytes memory invalidAddress = abi.encode(type(uint208).max);
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress, 0);
  }

  function test_validateDestFamilyAddress_EVM_RevertWhen_InvalidEVMAddressEncodePacked() public {
    bytes memory invalidAddress = abi.encodePacked(address(234));
    vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress, 0);
  }

  function test_validateDestFamilyAddress_EVM_RevertWhen_InvalidEVMAddressPrecompiles() public {
    for (uint160 i = 0; i < Internal.PRECOMPILE_SPACE; ++i) {
      bytes memory invalidAddress = abi.encode(address(i));
      vm.expectRevert(abi.encodeWithSelector(Internal.InvalidEVMAddress.selector, invalidAddress));
      s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_EVM, invalidAddress, 0);
    }

    s_feeQuoter.validateDestFamilyAddress(
      Internal.CHAIN_FAMILY_SELECTOR_EVM, abi.encode(address(uint160(Internal.PRECOMPILE_SPACE))), 0
    );
  }

  function test_validateDestFamilyAddress_SVM_RevertWhen_Invalid32ByteAddress() public {
    bytes memory invalidAddress = abi.encode(address(234), address(234));

    vm.expectRevert(abi.encodeWithSelector(Internal.Invalid32ByteAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_SVM, invalidAddress, 0);
  }

  function test_validateDestFamilyAddress_Aptos_RevertWhen_Invalid32ByteAddress() public {
    bytes memory invalidAddress = abi.encode(address(234), address(234));

    vm.expectRevert(abi.encodeWithSelector(Internal.Invalid32ByteAddress.selector, invalidAddress));
    s_feeQuoter.validateDestFamilyAddress(Internal.CHAIN_FAMILY_SELECTOR_APTOS, invalidAddress, 0);
  }
}
