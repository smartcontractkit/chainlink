// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITokenMessenger} from "../../../../pools/USDC/ITokenMessenger.sol";

import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {USDCSetup} from "../USDCSetup.t.sol";

contract USDCTokenPool_constructor is USDCSetup {
  function test_constructor() public {
    new USDCTokenPool(
      s_mockUSDC, s_cctpMessageTransmitterProxy, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );
  }

  function test_constructor_RevertWhen_TokenMessangerAddressZero() public {
    vm.expectRevert(USDCTokenPool.InvalidConfig.selector);
    new USDCTokenPool(
      ITokenMessenger(address(0)),
      s_cctpMessageTransmitterProxy,
      s_token,
      new address[](0),
      address(s_mockRMNRemote),
      address(s_router)
    );
  }

  function test_constructor_RevertWhen_TransmitterVersionDoesNotMatchSupportedUSDCVersion() public {
    uint32 transmitterVersion = uint32(vm.randomUint());
    vm.mockCall(
      address(s_mockUSDCTransmitter), abi.encodeCall(s_mockUSDCTransmitter.version, ()), abi.encode(transmitterVersion)
    );
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidMessageVersion.selector, transmitterVersion));
    new USDCTokenPool(
      s_mockUSDC, s_cctpMessageTransmitterProxy, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );
  }

  function test_constructor_RevertWhen_TokenMessengerVersionDoesNotMatchSupportedUSDCVersion() public {
    uint32 tokenMessengerVersion = s_mockUSDC.messageBodyVersion() + 1;
    vm.mockCall(
      address(s_mockUSDC), abi.encodeCall(s_mockUSDC.messageBodyVersion, ()), abi.encode(tokenMessengerVersion)
    );
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidTokenMessengerVersion.selector, tokenMessengerVersion));
    new USDCTokenPool(
      s_mockUSDC, s_cctpMessageTransmitterProxy, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );
  }

  function test_constructor_RevertWhen_InvalidTransmitterVersionInProxy() public {
    address transmitterAddress = makeAddr("INVALID_TRANSMITTER");
    vm.mockCall(
      address(s_cctpMessageTransmitterProxy),
      abi.encodeCall(s_cctpMessageTransmitterProxy.i_cctpTransmitter, ()),
      abi.encode(transmitterAddress)
    );
    vm.expectRevert(abi.encodeWithSelector(USDCTokenPool.InvalidTransmitterInProxy.selector));
    new USDCTokenPool(
      s_mockUSDC, s_cctpMessageTransmitterProxy, s_token, new address[](0), address(s_mockRMNRemote), address(s_router)
    );
  }
}
