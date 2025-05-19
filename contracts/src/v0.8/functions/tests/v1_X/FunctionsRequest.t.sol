// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";

import {Test} from "forge-std/Test.sol";

/// @notice #REQUEST_DATA_VERSION
contract FunctionsRequest_REQUEST_DATA_VERSION is Test {
  function test_REQUEST_DATA_VERSION() public {
    // Exposes REQUEST_DATA_VERSION
    assertEq(FunctionsRequest.REQUEST_DATA_VERSION, 1);
  }
}

/// @notice #DEFAULT_BUFFER_SIZE
contract FunctionsRequest_DEFAULT_BUFFER_SIZE is Test {
  function test_DEFAULT_BUFFER_SIZE() public {
    // Exposes DEFAULT_BUFFER_SIZE
    assertEq(FunctionsRequest.DEFAULT_BUFFER_SIZE, 256);
  }
}

/// @notice #encodeCBOR
contract FunctionsRequest_EncodeCBOR is Test {
  function test_EncodeCBOR_Success() public {
    // Exposes DEFAULT_BUFFER_SIZE
    assertEq(FunctionsRequest.DEFAULT_BUFFER_SIZE, 256);
  }
}

/// @notice #initializeRequest
contract FunctionsRequest_InitializeRequest is Test {}

/// @notice #initializeRequestForInlineJavaScript
contract FunctionsRequest_InitializeRequestForInlineJavaScript is Test {}

/// @notice #addSecretsReference
contract FunctionsRequest_AddSecretsReference is Test {}

/// @notice #addDONHostedSecrets
contract FunctionsRequest_AddDONHostedSecrets is Test {}

/// @notice #setArgs
contract FunctionsRequest_SetArgs is Test {}

/// @notice #setBytesArgs
contract FunctionsRequest_SetBytesArgs is Test {}
