// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface CCIPRouterInterface{

  function requestCrossChainSendTo(
    uint256 chainId,
    address receiver,
    bytes calldata message,
    address[] calldata tokens,
    uint256[] calldata amounts,
    bytes calldata options
  )
    external
    returns(
      bytes32
    );
  
  function routeMessage(
    bytes32 requestId,
    uint256 chainId,
    address sender,
    address receiver,
    bytes calldata message
  )
    external
    returns(
      bool
    );
}
