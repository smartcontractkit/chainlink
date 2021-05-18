// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./Operator.sol";
import "./AuthorizedForwarder.sol";

/**
 * @title Operator Factory
 * @notice Creates Operator contracts for node operators
 */
contract OperatorFactory {

  address public immutable link;

  event OperatorCreated(
    address indexed operator,
    address indexed owner
  );
  event AuthorizedForwarderCreated(
    address indexed operator,
    address indexed owner
  );

  /**
   * @param linkAddress address
   */
  constructor(
    address linkAddress
  ) {
    link = linkAddress;
  }

  /**
   * @notice creates a new Operator contract with the msg.sender as owner
   */
  function deployNewOperator()
    external
  {
    Operator operator = new Operator(link, msg.sender);
    emit OperatorCreated(address(operator), msg.sender);
  }

  /**
   * @notice creates a new Operator contract with the msg.sender as owner and a
   * new Operator Forwarder with the Operator as the owner
   */
  function deployNewOperatorAndForwarder()
    external
  {
    Operator operator = new Operator(link, msg.sender);
    emit OperatorCreated(address(operator), msg.sender);

    bytes memory tmp = new bytes(0);
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      link,
      address(operator),
      address(0),
      tmp
    );

    emit AuthorizedForwarderCreated(
      address(forwarder),
      msg.sender
    );
  }

  /**
   * @notice creates a new Forwarder contract with the msg.sender as owner
   */
  function deployNewForwarder()
    external
  {
    bytes memory tmp = new bytes(0);
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      link,
      msg.sender,
      address(0),
      tmp
    );

    emit AuthorizedForwarderCreated(
      address(forwarder),
      msg.sender
    );
  }

  /**
   * @notice creates a new Forwarder contract with the msg.sender as owner
   */
  function deployNewForwarderAndTransferOwnership(
    address to,
    bytes calldata message
  )
    external
  {
    AuthorizedForwarder forwarder = new AuthorizedForwarder(
      link,
      msg.sender,
      to,
      message
    );

    emit AuthorizedForwarderCreated(
      address(forwarder),
      msg.sender
    );
  }

}
