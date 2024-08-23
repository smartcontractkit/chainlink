// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {IDestinationVerifierProxy} from "./interfaces/IDestinationVerifierProxy.sol";
import {IDestinationVerifier} from "./interfaces/IDestinationVerifier.sol";

/**
 * @title DestinationVerifierProxy
 * @author Michael Fletcher
 * @notice This contract will be used to route all requests through to the assigned verifier contract. This contract does not support individual feed configurations and is aimed at being a simple proxy for the verifier contract on any destination chain.
 */
contract DestinationVerifierProxy is IDestinationVerifierProxy, ConfirmedOwner, TypeAndVersionInterface, IERC165 {
  /// @notice The active verifier for this proxy
  IDestinationVerifier private s_verifier;

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when trying to set a verifier address that does not implement the expected interface
  error VerifierInvalid(address verifierAddress);

  constructor() ConfirmedOwner(msg.sender) {}

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "DestinationVerifierProxy 1.0.0";
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
    if (
      !IERC165(verifierAddress).supportsInterface(IDestinationVerifier.s_accessController.selector) ||
      !IERC165(verifierAddress).supportsInterface(IDestinationVerifier.s_feeManager.selector) ||
      !IERC165(verifierAddress).supportsInterface(IDestinationVerifier.verify.selector) ||
      !IERC165(verifierAddress).supportsInterface(IDestinationVerifier.verifyBulk.selector)
    ) revert VerifierInvalid(verifierAddress);

    s_verifier = IDestinationVerifier(verifierAddress);
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
    return
      interfaceId == this.setVerifier.selector ||
      interfaceId == this.verify.selector ||
      interfaceId == this.verifyBulk.selector ||
      interfaceId == this.s_feeManager.selector ||
      interfaceId == this.s_accessController.selector;
  }
}
