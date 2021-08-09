// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface CCIPRampInterface{

  event CrossChainSendRequested(
    bytes32 indexed requestId,
    uint256 chainId,
    address sender,
    address receiver,
    bytes message,
    address[] tokens,
    address[] amounts,
    bytes options
  );
  event CrossChainMessagedReceived(
    bytes32 indexed requestId,
    uint256 chainId,
    address sender,
    address receiver,
    bytes message,
    address[] tokens,
    address[] amounts,
    bytes options
  );
  event CrossChainMessageFulfilled(
    bytes32 indexed requestId
  );

  function requestCrossChainSend(
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

  function tokenPool(
    address token
  )
    external
    view
    returns(
      address
    );
  
  function chainId()
    external
    view
    returns(
      uint256
    );
  
  function router()
    external
    view
    returns(
      address
    );
}
