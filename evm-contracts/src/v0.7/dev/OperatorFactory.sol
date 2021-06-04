// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./Operator.sol";
import "./AuthorizedForwarder.sol";

/**
 * @title Operator Factory
 * @notice Creates Operator contracts for node operators
 */
contract OperatorFactory {

  address public immutable getChainlinkToken;
  mapping(address => bool) private s_created;

  event OperatorCreated(
    address indexed operator,
    address indexed owner,
    address indexed sender
  );
  event AuthorizedForwarderCreated(
    address indexed forwarder,
    address indexed owner,
    address indexed sender
  );

  /**
   * @param linkAddress address
   */
  constructor(
    address linkAddress
  ) {
    getChainlinkToken = linkAddress;
  }

  /**
   * @notice creates a new Operator contract with the msg.sender as owner
   */
  function deployNewOperator()
    external
    returns (
      address
    )
  {
    Operator operator = new Operator(
      getChainlinkToken,
      msg.sender
    );

    s_created[address(operator)] = true;
    emit OperatorCreated(
      address(operator),
      msg.sender,
      msg.sender
    );

    return address(operator);
  }

  /**
   * @notice creates a new Operator contract with the msg.sender as owner and a
   * new Operator Forwarder with the Operator as the owner
   */
  function deployNewOperatorAndForwarder()
    external
    returns (
      address,
      address
    )
  {
    Operator operator = new Operator(
      getChainlinkToken,
      msg.sender
    );
    s_created[address(operator)] = true;
    emit OperatorCreated(
      address(operator),
      msg.sender,
      msg.sender
    );

    bytes memory tmp = new bytes(0);
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      getChainlinkToken,
      address(operator),
      address(0),
      tmp
    );
    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(
      address(forwarder),
      address(operator),
      msg.sender
    );

    return (address(operator), address(forwarder));
  }

  /**
   * @notice creates a new Forwarder contract with the msg.sender as owner
   */
  function deployNewForwarder()
    external
    returns (
      address
    )
  {
    bytes memory tmp = new bytes(0);
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      getChainlinkToken,
      msg.sender,
      address(0),
      tmp
    );

    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(
      address(forwarder),
      msg.sender,
      msg.sender
    );

    return address(forwarder);
  }

  /**
   * @notice creates a new Forwarder contract with the msg.sender as owner
   */
  function deployNewForwarderAndTransferOwnership(
    address to,
    bytes calldata message
  )
    external
    returns (
      address
    )
  {
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      getChainlinkToken,
      msg.sender,
      to,
      message
    );

    s_created[address(forwarder)] = true;
    emit AuthorizedForwarderCreated(
      address(forwarder),
      msg.sender,
      msg.sender
    );

    return address(forwarder);
  }

  /**
   * @notice indicates whether this factory deployed an address
   */
  function created(
    address query
  )
    external
    view
    returns (bool)
  {
    return s_created[query];
  }

}
