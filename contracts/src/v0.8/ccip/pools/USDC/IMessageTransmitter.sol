/*
 * Copyright (c) 2022, Circle Internet Financial Limited.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
pragma solidity ^0.8.0;

interface IMessageTransmitter {
  /// @notice Unlocks USDC tokens on the destination chain
  /// @param message The original message on the source chain
  ///     * Message format:
  ///     * Field                 Bytes      Type       Index
  ///     * version               4          uint32     0
  ///     * sourceDomain          4          uint32     4
  ///     * destinationDomain     4          uint32     8
  ///     * nonce                 8          uint64     12
  ///     * sender                32         bytes32    20
  ///     * recipient             32         bytes32    52
  ///     * destinationCaller     32         bytes32    84
  ///     * messageBody           dynamic    bytes      116
  /// param attestation A valid attestation is the concatenated 65-byte signature(s) of
  /// exactly `thresholdSignature` signatures, in increasing order of attester address.
  /// ***If the attester addresses recovered from signatures are not in increasing order,
  /// signature verification will fail.***
  /// If incorrect number of signatures or duplicate signatures are supplied,
  /// signature verification will fail.
  function receiveMessage(bytes calldata message, bytes calldata attestation) external returns (bool success);

  /// Returns domain of chain on which the contract is deployed.
  /// @dev immutable
  function localDomain() external view returns (uint32);

  /// Returns message format version.
  /// @dev immutable
  function version() external view returns (uint32);
}
