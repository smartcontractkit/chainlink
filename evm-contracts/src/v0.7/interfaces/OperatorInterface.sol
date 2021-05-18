// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./ChainlinkRequestInterface.sol";
import "./OracleInterface.sol";

interface OperatorInterface is
  ChainlinkRequestInterface,
  OracleInterface
{

  function requestOracleData(
    address sender,
    uint256 payment,
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion,
    bytes calldata data
  )
    external;

  function fulfillOracleRequest2(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    bytes calldata data
  )
    external
    returns (
      bool
    );

  function operatorTransferAndCall(
    address to,
    uint256 value,
    bytes calldata data
  )
    external
    returns (
      bool success
    );

  function distributeFunds(
    address payable[] calldata receivers,
    uint[] calldata amounts
  )
    external
    payable;

  function setAuthorizedSendersOn(
    address[] calldata senders,
    address[] calldata targets
  )
    external;

  function transferOwnableContracts(
    address[] calldata ownable,
    address newOwner
  )
    external;

  function acceptOwnableContracts(
    address[] calldata ownable
  )
    external;

}
