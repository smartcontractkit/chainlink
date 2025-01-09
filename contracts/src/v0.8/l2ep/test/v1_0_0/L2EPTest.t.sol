// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Greeter} from "../Greeter.sol";

import {MultiSend} from "../../../vendor/MultiSend.sol";
import {Test} from "forge-std/Test.sol";

contract L2EPTest is Test {
  /// Helper variable(s)
  address internal s_strangerAddr = vm.addr(0x1);
  address internal s_l1OwnerAddr = vm.addr(0x2);
  address internal s_eoaValidator = vm.addr(0x3);
  address internal s_deployerAddr = vm.addr(0x4);

  /// @param selector - the function selector
  /// @param greeterAddr - the address of the Greeter contract
  /// @param message - the new greeting message, which will be passed as an argument to Greeter#setGreeting
  /// @return a 2-layer encoding such that decoding the first layer provides the CrossDomainForwarder#forward
  ///         function selector and the corresponding arguments to the forward function, and decoding the
  ///         second layer provides the Greeter#setGreeting function selector and the corresponding
  ///         arguments to the set greeting function (which in this case is the input message)
  function encodeCrossDomainSetGreetingMsg(
    bytes4 selector,
    address greeterAddr,
    string memory message
  ) public pure returns (bytes memory) {
    return abi.encodeWithSelector(selector, greeterAddr, abi.encodeWithSelector(Greeter.setGreeting.selector, message));
  }

  /// @param selector - the function selector
  /// @param multiSendAddr - the address of the MultiSend contract
  /// @param encodedTxs - an encoded list of transactions (e.g. abi.encodePacked(encodeMultiSendTx("some data"), ...))
  /// @return a 2-layer encoding such that decoding the first layer provides the CrossDomainGoverner#forwardDelegate
  ///         function selector and the corresponding arguments to the forwardDelegate function, and decoding the
  ///         second layer provides the MultiSend#multiSend function selector and the corresponding
  ///         arguments to the multiSend function (which in this case is the input encodedTxs)
  function encodeCrossDomainMultiSendMsg(
    bytes4 selector,
    address multiSendAddr,
    bytes memory encodedTxs
  ) public pure returns (bytes memory) {
    return
      abi.encodeWithSelector(selector, multiSendAddr, abi.encodeWithSelector(MultiSend.multiSend.selector, encodedTxs));
  }

  /// @param greeterAddr - the address of the greeter contract
  /// @param data - the transaction data string
  /// @return an encoded transaction structured as specified in the MultiSend#multiSend comments
  function encodeMultiSendTx(address greeterAddr, bytes memory data) public pure returns (bytes memory) {
    bytes memory txData = abi.encodeWithSelector(Greeter.setGreeting.selector, data);
    return
      abi.encodePacked(
        uint8(0), // operation
        greeterAddr, // to
        uint256(0), // value
        uint256(txData.length), // data length
        txData // data as bytes
      );
  }

  /// @param l1Address - Address on L1
  /// @return an Arbitrum L2 address
  function toArbitrumL2AliasAddress(address l1Address) public pure returns (address) {
    return address(uint160(l1Address) + uint160(0x1111000000000000000000000000000000001111));
  }
}
