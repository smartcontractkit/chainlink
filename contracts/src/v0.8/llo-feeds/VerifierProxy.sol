// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IVerifierProxy} from "./interfaces/IVerifierProxy.sol";
import {IVerifier} from "./interfaces/IVerifier.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {AccessControllerInterface} from "../interfaces/AccessControllerInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";
import {IRewardManager} from "./interfaces/IRewardManager.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {IERC20} from "../shared/vendor/IERC20.sol";
import {Common} from "../libraries/internal/Common.sol";
import {IWERC20} from "../shared/vendor/IWERC20.sol";

/**
 * The verifier proxy contract is the gateway for all report verification requests
 * on a chain.  It is responsible for taking in a verification request and routing
 * it to the correct verifier contract.
 */
contract VerifierProxy is IVerifierProxy, ConfirmedOwner, TypeAndVersionInterface {
  /// @notice This event is emitted whenever a new verifier contract is set
  /// @param oldConfigDigest The config digest that was previously the latest config
  /// digest of the verifier contract at the verifier address.
  /// @param oldConfigDigest The latest config digest of the verifier contract
  /// at the verifier address.
  /// @param verifierAddress The address of the verifier contract that verifies reports for
  /// a given digest
  event VerifierSet(bytes32 oldConfigDigest, bytes32 newConfigDigest, address verifierAddress);

  /// @notice This event is emitted whenever a new verifier contract is initialized
  /// @param verifierAddress The address of the verifier contract that verifies reports
  event VerifierInitialized(address verifierAddress);

  /// @notice This event is emitted whenever a verifier is unset
  /// @param configDigest The config digest that was unset
  /// @param verifierAddress The Verifier contract address unset
  event VerifierUnset(bytes32 configDigest, address verifierAddress);

  /// @notice This event is emitted when a new access controller is set
  /// @param oldAccessController The old access controller address
  /// @param newAccessController The new access controller address
  event AccessControllerSet(address oldAccessController, address newAccessController);

  /// @notice This error is thrown whenever an address tries
  /// to exeecute a transaction that it is not authorized to do so
  error AccessForbidden();

  /// @notice This error is thrown whenever a zero address is passed
  error ZeroAddress();

  /// @notice This error is thrown when trying to set a verifier address
  /// for a digest that has already been initialized
  /// @param configDigest The digest for the verifier that has
  /// already been set
  /// @param verifier The address of the verifier the digest was set for
  error ConfigDigestAlreadySet(bytes32 configDigest, address verifier);

  /// @notice This error is thrown when trying to set a verifier address that has already been initialized
  error VerifierAlreadyInitialized(address verifier);

  /// @notice This error is thrown when the verifier at an address does
  /// not conform to the verifier interface
  error VerifierInvalid();

  /// @notice This error is thrown whenever a verifier is not found
  /// @param configDigest The digest for which a verifier is not found
  error VerifierNotFound(bytes32 configDigest);

  /// @notice This error is thrown when the verifier does not include the correct amount or quote to retrieve the correct deposit
  error InvalidDeposit();

  /// @notice Mapping of authorized verifiers
  mapping(address => bool) private s_initializedVerifiers;

  /// @notice Mapping between config digests and verifiers
  mapping(bytes32 => address) private s_verifiersByConfig;

  /// @notice The contract to control addresses that are allowed to verify reports
  AccessControllerInterface private s_accessController;

  /// @notice The contract to control fees for report verification
  IFeeManager private immutable s_feeManager;

  /// @notice The contract to control reward distribution for report verification
  IRewardManager private immutable s_rewardsManager;

  /// @notice The Wrapped Native contract
  address private immutable s_wrappedNative;

  constructor(
    AccessControllerInterface accessController,
    IFeeManager feeManager,
    IRewardManager rewardManager,
    address wrappedNativeAddress
  ) ConfirmedOwner(msg.sender) {
    s_accessController = accessController;
    s_feeManager = feeManager;
    s_rewardsManager = rewardManager;

    if (address(wrappedNativeAddress) == address(0)) revert ZeroAddress();
    s_wrappedNative = wrappedNativeAddress;
  }

  /// @dev reverts if the caller does not have access by the accessController contract or is the contract itself.
  modifier checkAccess() {
    AccessControllerInterface ac = s_accessController;
    if (address(ac) != address(0) && !ac.hasAccess(msg.sender, msg.data)) revert AccessForbidden();
    _;
  }

  /// @dev only allow verified addresses to call this function
  modifier onlyInitializedVerifier() {
    if (!s_initializedVerifiers[msg.sender]) revert AccessForbidden();
    _;
  }

  modifier onlyValidVerifier(address verifierAddress) {
    if (verifierAddress == address(0)) revert ZeroAddress();
    if (!IERC165(verifierAddress).supportsInterface(IVerifier.verify.selector)) revert VerifierInvalid();
    _;
  }

  /// @notice Reverts if the config digest has already been assigned
  /// a verifier
  modifier onlyUnsetConfigDigest(bytes32 configDigest) {
    address configDigestVerifier = s_verifiersByConfig[configDigest];
    if (configDigestVerifier != address(0)) revert ConfigDigestAlreadySet(configDigest, configDigestVerifier);
    _;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "VerifierProxy 1.0.0";
  }

  //***************************//
  //       Admin Functions     //
  //***************************//

  /// @notice This function can be called by the contract admin to set
  /// the proxy's access controller contract
  /// @param accessController The new access controller to set
  /// @dev The access controller can be set to the zero address to allow
  /// all addresses to verify reports
  function setAccessController(AccessControllerInterface accessController) external onlyOwner {
    address oldAccessController = address(s_accessController);
    s_accessController = accessController;
    emit AccessControllerSet(oldAccessController, address(accessController));
  }

  /// @notice Returns the current access controller
  /// @return accessController The current access controller contract
  /// the proxy is using to gate access
  function getAccessController() external view returns (AccessControllerInterface accessController) {
    return s_accessController;
  }

  //***************************//
  //  Verification Functions   //
  //***************************//

  /// @inheritdoc IVerifierProxy
  /// @dev Contract skips checking whether or not the current verifier
  /// is valid as it checks this before a new verifier is set.
  function verify(
    bytes calldata payload
  ) external payable override checkAccess returns (bytes memory verifierResponse) {
    // First 32 bytes of the signed report is the config digest.
    bytes32 configDigest = bytes32(payload);
    address verifierAddress = s_verifiersByConfig[configDigest];
    if (verifierAddress == address(0)) revert VerifierNotFound(configDigest);
    (bytes memory verifiedReport, bytes memory quoteData) = IVerifier(verifierAddress).verify(payload, msg.sender);

    //if we have a registered fee and reward-manager manager, bill the verifier.
    if (address(s_feeManager) != address(0) && address(s_rewardsManager) != address(0)) {
      //decode the fee
      Common.Asset memory asset = s_feeManager.getFee(msg.sender, verifiedReport, quoteData);

      //some users might not be billed
      if (asset.amount > 0) {
        //get the sender of the funds, wrapping it will turn the proxy into the sender
        address transferFeeFromAddress = msg.sender;

        //if native has been sent in, calculate the amount to wrap and return the rest
        if (msg.value > 0) {
          if (asset.assetAddress != s_wrappedNative) revert InvalidDeposit();
          if (msg.value < asset.amount) revert InvalidDeposit();

          //wrap the amount required to pay the fee & approve
          IWERC20(s_wrappedNative).deposit{value: asset.amount}();
          IERC20(s_wrappedNative).approve(address(s_rewardsManager), asset.amount);
          transferFeeFromAddress = address(this);

          unchecked {
            //msg.value is always >= to asset.amount
            uint256 change = msg.value - asset.amount;

            //return the change
            if (change > 0) {
              payable(msg.sender).transfer(change);
            }
          }
        }

        //bill the payee
        s_rewardsManager.onFeePaid(configDigest, transferFeeFromAddress, asset);
      }
    }

    return verifiedReport;
  }

  /// @inheritdoc IVerifierProxy
  function initializeVerifier(address verifierAddress) external override onlyOwner onlyValidVerifier(verifierAddress) {
    if (s_initializedVerifiers[verifierAddress]) revert VerifierAlreadyInitialized(verifierAddress);

    s_initializedVerifiers[verifierAddress] = true;
    emit VerifierInitialized(verifierAddress);
  }

  /// @inheritdoc IVerifierProxy
  function setVerifier(
    bytes32 currentConfigDigest,
    bytes32 newConfigDigest,
    Common.AddressAndWeight[] memory addressAndWeights
  ) external override onlyUnsetConfigDigest(newConfigDigest) onlyInitializedVerifier {
    s_verifiersByConfig[newConfigDigest] = msg.sender;

    //empty recipients array will be ignored and will need to be set off chain
    if (addressAndWeights.length > 0) {
      //Set the reward recipients for this digest
      s_rewardsManager.setRewardRecipients(newConfigDigest, addressAndWeights);
    }

    emit VerifierSet(currentConfigDigest, newConfigDigest, msg.sender);
  }

  /// @inheritdoc IVerifierProxy
  function unsetVerifier(bytes32 configDigest) external override onlyOwner {
    address verifierAddress = s_verifiersByConfig[configDigest];
    if (verifierAddress == address(0)) revert VerifierNotFound(configDigest);
    delete s_verifiersByConfig[configDigest];
    emit VerifierUnset(configDigest, verifierAddress);
  }

  /// @inheritdoc IVerifierProxy
  function getVerifier(bytes32 configDigest) external view override returns (address) {
    return s_verifiersByConfig[configDigest];
  }
}
