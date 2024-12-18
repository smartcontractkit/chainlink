// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IAny2EVMMessageReceiver} from "../../../interfaces/IAny2EVMMessageReceiver.sol";

import {ERC165CheckerReverting} from "../../../libraries/ERC165CheckerReverting.sol";
import {ERC165CheckerRevertingSetup} from "./ERC165CheckerRevertingSetup.t.sol";

contract ERC165CheckerReverting_getSupportedInterfaces is ERC165CheckerRevertingSetup {
  using ERC165CheckerReverting for address;

  function test__getSupportedInterfacesReverting() public view {
    bytes4[] memory interfaceIds = new bytes4[](1);

    interfaceIds[0] = type(IAny2EVMMessageReceiver).interfaceId;

    bool[] memory supportedIds = s_receiver._getSupportedInterfacesReverting(interfaceIds);
    assertTrue(supportedIds[0]);
    assertEq(interfaceIds.length, supportedIds.length);
  }
}
