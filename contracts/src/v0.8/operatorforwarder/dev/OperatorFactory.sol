// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Operator} from "./Operator.sol";
import {AuthorizedForwarder} from "./AuthorizedForwarder.sol";

// @title Operator Factory
// @notice Creates Operator contracts for node operators
// solhint-disable gas-custom-errors
contract OperatorFactory {
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable linkToken;
  mapping(address => bool) private s_created;

  event OperatorCreated(address indexed operator, address indexed owner, address indexed sender);
  event AuthorizedForwarderCreated(address indexed forwarder, address indexed owner, address indexed sender);

  // @param linkAddress address
  constructor(address linkAddress) {
    linkToken = linkAddress;
  }

  string public constant typeAndVersion = "OperatorFactory 1.0.0";

  // @notice creates a new Operator contract with the msg.sender as owner
  function deployNewOperator() external returns (address) {
    Operator operator = new Operator(linkToken, msg.sender);

    s_created[address(operator)] = true;
    emit OperatorCreated(address(operator), msg.sender, msg.sender);

    return address(operator);
  }

  // @notice creates a new Operator contract with the msg.sender as owner and a
  // new Operator Forwarder with the OperatorFactory as the owner
  function deployNewOperatorAndForwarder() external returns (address, address) {
    Operator operator = new Operator(linkToken, msg.sender);
    s_created[address(operator)] = true;
    emit OperatorCreated(address(operator), msg.sender, msg.sender);

    AuthorizedForwarder forwarder = new AuthorizedForwarder(linkToken, address(this), address(operator), new bytes(0));
    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(address(forwarder), address(this), msg.sender);

    return (address(operator), address(forwarder));
  }

  // @notice creates a new Forwarder contract with the msg.sender as owner
  function deployNewForwarder() external returns (address) {
    AuthorizedForwarder forwarder = new AuthorizedForwarder(linkToken, msg.sender, address(0), new bytes(0));

    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(address(forwarder), msg.sender, msg.sender);

    return address(forwarder);
  }

  // @notice creates a new Forwarder contract with the msg.sender as owner
  function deployNewForwarderAndTransferOwnership(address to, bytes calldata message) external returns (address) {
    AuthorizedForwarder forwarder = new AuthorizedForwarder(linkToken, msg.sender, to, message);

    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(address(forwarder), msg.sender, msg.sender);

    return address(forwarder);
  }

  // @notice indicates whether this factory deployed an address
  function created(address query) external view returns (bool) {
    return s_created[query];
  }
}
