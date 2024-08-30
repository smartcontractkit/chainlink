// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IDestinationVerifierProxy} from "./interfaces/IDestinationVerifierProxy.sol";
import {IDestinationVerifierProxyVerifier} from "./interfaces/IDestinationVerifierProxyVerifier.sol";

/**
 * @title DestinationVerifierProxy
 * @author Michael Fletcher
 * @notice This contract will be used to route all requests through to the assigned verifier contract. This contract does not support individual feed configurations and is aimed at being a simple proxy for the verifier contract on any destination chain.
 */
contract DestinationVerifierProxy is IDestinationVerifierProxy, ConfirmedOwner, TypeAndVersionInterface {
  /// @notice The active verifier for this proxy
  IDestinationVerifierProxyVerifier private s_verifier;

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when trying to set a verifier address that does not implement the expected interface
  error VerifierInvalid(address verifierAddress);

  constructor() ConfirmedOwner(msg.sender) {}

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "DestinationVerifierProxy 0.4.0";
  }

  /// @inheritdoc IDestinationVerifierProxy
  function verify(bytes calldata payload, bytes calldata parameterPayload) external payable returns (bytes memory) {
    return s_verifier.verify{value: msg.value}(payload, parameterPayload, msg.sender);
  }

  /// @inheritdoc IDestinationVerifierProxy
  function verifyBulk(
    bytes[] calldata payloads,
    bytes calldata parameterPayload
  ) external payable returns (bytes[] memory verifiedReports) {
    return s_verifier.verifyBulk{value: msg.value}(payloads, parameterPayload, msg.sender);
  }

  /// @inheritdoc IDestinationVerifierProxy
  function setVerifier(address verifierAddress) external onlyOwner {
    //check it supports the functions we need
    if (!IERC165(verifierAddress).supportsInterface(type(IDestinationVerifierProxyVerifier).interfaceId))
      revert VerifierInvalid(verifierAddress);

    s_verifier = IDestinationVerifierProxyVerifier(verifierAddress);
  }

  /// @inheritdoc IDestinationVerifierProxy
  // solhint-disable-next-line func-name-mixedcase
  function s_feeManager() external view override returns (address) {
    return s_verifier.s_feeManager();
  }

  /// @inheritdoc IDestinationVerifierProxy
  // solhint-disable-next-line func-name-mixedcase
  function s_accessController() external view override returns (address) {
    return s_verifier.s_accessController();
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == type(IDestinationVerifierProxy).interfaceId;
  }
}
