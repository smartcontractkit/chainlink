// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";
import {IDestinationVerifierProxy} from "./interfaces/IDestinationVerifierProxy.sol";
import {IDestinationVerifier} from "./interfaces/IDestinationVerifier.sol";

/**
 * The verifier proxy contract is the gateway for all report verification requests
 * on a chain.  It is responsible for taking in a verification request and routing
 * it to the correct verifier contract.
 */
contract DestinationVerifierProxy is IDestinationVerifierProxy, ConfirmedOwner, TypeAndVersionInterface {

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
  function verify(
    bytes calldata payload,
    bytes calldata parameterPayload
  ) external payable returns (bytes memory) {
    return s_verifier.verify(payload, parameterPayload, msg.sender);
  }

  /// @inheritdoc IDestinationVerifierProxy
  function verifyBulk(
    bytes[] calldata payloads,
    bytes calldata parameterPayload
  ) external payable returns (bytes[] memory verifiedReports) {
    return s_verifier.verifyBulk(payloads, parameterPayload, msg.sender);
  }


  /// @inheritdoc IDestinationVerifierProxy
  function setVerifier(address verifierAddress) external onlyOwner {
    if(verifierAddress == address(0)) revert ZeroAddress();

    //TODO Selector

    s_verifier = IDestinationVerifier(verifierAddress);
  }

   /// @inheritdoc IDestinationVerifierProxy
  function s_feeManager() external view override returns (address) {
    return s_verifier.getFeeManager();
  }

  function s_accessController() external view override returns (address) {
    return s_verifier.getAccessController();
  }
}
