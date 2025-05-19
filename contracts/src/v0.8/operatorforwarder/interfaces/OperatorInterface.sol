// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {OracleInterface} from "./OracleInterface.sol";
import {ChainlinkRequestInterface} from "./ChainlinkRequestInterface.sol";

// solhint-disable-next-line interface-starts-with-i
interface OperatorInterface is OracleInterface, ChainlinkRequestInterface {
  function operatorRequest(
    address sender,
    uint256 payment,
    bytes32 specId,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion,
    bytes calldata data
  ) external;

  function fulfillOracleRequest2(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    bytes calldata data
  ) external returns (bool);

  function ownerTransferAndCall(address to, uint256 value, bytes calldata data) external returns (bool success);

  function distributeFunds(address payable[] calldata receivers, uint256[] calldata amounts) external payable;
}
