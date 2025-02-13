// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC165} from "../../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/introspection/IERC165.sol";
import {IAny2EVMMessageReceiver} from "../../../../interfaces/IAny2EVMMessageReceiver.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_supportsInterface is EtherSenderReceiverTestSetup {
  function test_supportsInterface() public view {
    assertTrue(s_etherSenderReceiver.supportsInterface(type(IERC165).interfaceId));
    assertTrue(s_etherSenderReceiver.supportsInterface(type(IAny2EVMMessageReceiver).interfaceId));
  }
}
