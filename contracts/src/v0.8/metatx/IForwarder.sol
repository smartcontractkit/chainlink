// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;
pragma abicoder v2;

import {IERC165} from "../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/introspection/IERC165.sol";

/// @title The Forwarder Interface
/// @notice The contracts implementing this interface take a role of authorization, authentication and replay protection
/// for contracts that choose to trust a `Forwarder`, instead of relying on a mechanism built into the Ethereum protocol.
///
/// @notice if the `Forwarder` contract decides that an incoming `ForwardRequest` is valid, it must append 20 bytes that
/// represent the caller to the `data` field of the request and send this new data to the target address (the `to` field)
///
/// @dev This implementation has been ported from OpenGSN's Forwarder.sol and modified in following ways:
/// @dev 1. execute() does not accept "gas" parameter which allows caller to specify max gas limit for the forwarded call
/// @dev 2. execute() does not accept "value" parameter which allows caller to pass native token to the forwarded call
/// @dev 3. renamed field: "address to" => "address target"
///
/// :warning: **Warning** :warning: The Forwarder can have a full control over a `Recipient` contract.
/// Any vulnerability in a `Forwarder` implementation can make all of its `Recipient` contracts susceptible!
/// Recipient contracts should only trust forwarders that passed through security audit,
/// otherwise they are susceptible to identity theft.
interface IForwarder is IERC165 {
  /// @notice A representation of a request for a `Forwarder` to send `data` on behalf of a `from` to a target (`to`).
  struct ForwardRequest {
    address from;
    address target;
    uint256 nonce;
    bytes data;
    uint256 validUntilTime;
  }

  event DomainRegistered(bytes32 indexed domainSeparator, bytes domainValue);

  event RequestTypeRegistered(bytes32 indexed typeHash, string typeStr);

  /// @notice Verify the transaction is valid and can be executed.
  /// @dev Implementations must validate the signature and the nonce of the request are correct.
  /// @dev Does not revert and returns successfully if the input is valid.
  /// @dev Reverts if any validation has failed. For instance, if either signature or nonce are incorrect.
  /// @dev Reverts if `domainSeparator` or `requestTypeHash` are not registered as well.
  function verify(
    ForwardRequest calldata forwardRequest,
    bytes32 domainSeparator,
    bytes32 requestTypeHash,
    bytes calldata suffixData,
    bytes calldata signature
  ) external view;

  /// @notice Executes a transaction specified by the `ForwardRequest`.
  /// The transaction is first verified and then executed.
  /// The success flag and returned bytes array of the `CALL` are returned as-is.
  ///
  /// This method would revert only in case of a verification error.
  ///
  /// All the target errors are reported using the returned success flag and returned bytes array.
  ///
  /// @param forwardRequest All requested transaction parameters.
  /// @param domainSeparator The domain used when signing this request.
  /// @param requestTypeHash The request type used when signing this request.
  /// @param suffixData The ABI-encoded extension data for the current `RequestType` used when signing this request.
  /// @param signature The client signature to be validated.
  function execute(
    ForwardRequest calldata forwardRequest,
    bytes32 domainSeparator,
    bytes32 requestTypeHash,
    bytes calldata suffixData,
    bytes calldata signature
  ) external payable returns (bool success, bytes memory ret);

  /// @notice Register a new Request typehash.
  /// @notice This is necessary for the Forwarder to be able to verify the signatures conforming to the ERC-712.
  /// @param typeName The name of the request type.
  /// @param typeSuffix Any extra data after the generic params. Must contain add at least one param.
  /// The generic ForwardRequest type is always registered by the constructor.
  function registerRequestType(string calldata typeName, string calldata typeSuffix) external;

  /// @notice Register a new domain separator.
  /// @notice This is necessary for the Forwarder to be able to verify the signatures conforming to the ERC-712.
  /// @notice The domain separator must have the following fields: `name`, `version`, `chainId`, `verifyingContract`.
  /// The `chainId` is the current network's `chainId`, and the `verifyingContract` is this Forwarder's address.
  /// This method accepts the domain name and version to create and register the domain separator value.
  /// @param name The domain's display name.
  /// @param version The domain/protocol version.
  function registerDomainSeparator(string calldata name, string calldata version) external;
}
