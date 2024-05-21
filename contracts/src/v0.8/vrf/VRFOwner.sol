// SPDX-License-Identifier: MIT
// solhint-disable-next-line one-contract-per-file
pragma solidity ^0.8.6;

import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {AuthorizedReceiver} from "./AuthorizedReceiver.sol";
import {VRFTypes} from "./VRFTypes.sol";

// Taken from VRFCoordinatorV2.sol
// Must be abi-compatible with what's there
struct FeeConfig {
  // Flat fee charged per fulfillment in millionths of link
  // So fee range is [0, 2^32/10^6].
  uint32 fulfillmentFlatFeeLinkPPMTier1;
  uint32 fulfillmentFlatFeeLinkPPMTier2;
  uint32 fulfillmentFlatFeeLinkPPMTier3;
  uint32 fulfillmentFlatFeeLinkPPMTier4;
  uint32 fulfillmentFlatFeeLinkPPMTier5;
  uint24 reqsForTier2;
  uint24 reqsForTier3;
  uint24 reqsForTier4;
  uint24 reqsForTier5;
}

// Taken from VRFCoordinatorV2.sol
// Must be abi-compatible with what's there
struct Config {
  uint16 minimumRequestConfirmations;
  uint32 maxGasLimit;
  // stalenessSeconds is how long before we consider the feed price to be stale
  // and fallback to fallbackWeiPerUnitLink.
  uint32 stalenessSeconds;
  // Gas to cover oracle payment after we calculate the payment.
  // We make it configurable in case those operations are repriced.
  uint32 gasAfterPaymentCalculation;
  int256 fallbackWeiPerUnitLink;
  FeeConfig feeConfig;
}

/// @dev IVRFCoordinatorV2 is the set of functions on the VRF coordinator V2
/// @dev that are used in the VRF Owner contract below.
interface IVRFCoordinatorV2 {
  function acceptOwnership() external;

  function transferOwnership(address to) external;

  function registerProvingKey(address oracle, uint256[2] calldata publicProvingKey) external;

  function deregisterProvingKey(uint256[2] calldata publicProvingKey) external;

  function setConfig(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    FeeConfig memory feeConfig
  ) external;

  function getConfig()
    external
    view
    returns (
      uint16 minimumRequestConfirmations,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation
    );

  function getFeeConfig()
    external
    view
    returns (
      uint32 fulfillmentFlatFeeLinkPPMTier1,
      uint32 fulfillmentFlatFeeLinkPPMTier2,
      uint32 fulfillmentFlatFeeLinkPPMTier3,
      uint32 fulfillmentFlatFeeLinkPPMTier4,
      uint32 fulfillmentFlatFeeLinkPPMTier5,
      uint24 reqsForTier2,
      uint24 reqsForTier3,
      uint24 reqsForTier4,
      uint24 reqsForTier5
    );

  function getFallbackWeiPerUnitLink() external view returns (int256);

  function ownerCancelSubscription(uint64 subId) external;

  function recoverFunds(address to) external;

  function hashOfKey(uint256[2] memory publicKey) external pure returns (bytes32);

  function fulfillRandomWords(
    VRFTypes.Proof memory proof,
    VRFTypes.RequestCommitment memory rc
  ) external returns (uint96);
}

/**
 * @notice VRFOwner is a contract that acts as the owner of the VRF
 * @notice coordinator, with some useful utilities in the event extraordinary
 * @notice things happen on-chain (i.e ETH/LINK price wildly fluctuates, and
 * @notice a VRF fulfillment reverts on-chain).
 */
contract VRFOwner is ConfirmedOwner, AuthorizedReceiver {
  int256 private constant MAX_INT256 = 0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff;

  IVRFCoordinatorV2 internal s_vrfCoordinator;

  event RandomWordsForced(uint256 indexed requestId, uint64 indexed subId, address indexed sender);

  constructor(address _vrfCoordinator) ConfirmedOwner(msg.sender) {
    // solhint-disable-next-line gas-custom-errors
    require(_vrfCoordinator != address(0), "vrf coordinator address must be non-zero");
    s_vrfCoordinator = IVRFCoordinatorV2(_vrfCoordinator);
  }

  /**
   * @notice Accepts ownership of the VRF coordinator if transferred to us.
   */
  function acceptVRFOwnership() external onlyOwner {
    s_vrfCoordinator.acceptOwnership();
  }

  /**
   * @notice Transfers ownership of the VRF coordinator to the specified address.
   * @param to the address to transfer ownership of the VRF Coordinator to.
   */
  function transferVRFOwnership(address to) external onlyOwner {
    s_vrfCoordinator.transferOwnership(to);
  }

  /**
   * @notice Returns the address of the VRF coordinator reference in this contract.
   * @return The address of the VRF coordinator reference in this contract.
   */
  function getVRFCoordinator() public view returns (address) {
    return address(s_vrfCoordinator);
  }

  /**
   * @notice Registers a proving key to an oracle.
   * @param oracle address of the oracle
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function registerProvingKey(address oracle, uint256[2] calldata publicProvingKey) external onlyOwner {
    s_vrfCoordinator.registerProvingKey(oracle, publicProvingKey);
  }

  /**
   * @notice Deregisters a proving key to an oracle.
   * @param publicProvingKey key that oracle can use to submit vrf fulfillments
   */
  function deregisterProvingKey(uint256[2] calldata publicProvingKey) external onlyOwner {
    s_vrfCoordinator.deregisterProvingKey(publicProvingKey);
  }

  /**
   * @notice Sets the configuration of the vrfv2 coordinator
   * @param minimumRequestConfirmations global min for request confirmations
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @param feeConfig fee tier configuration
   */
  function setConfig(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    FeeConfig memory feeConfig
  ) public onlyOwner {
    s_vrfCoordinator.setConfig(
      minimumRequestConfirmations,
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      feeConfig
    );
  }

  /**
   * @notice Sets the configuration of the vrfv2 coordinator - only to be used from within fulfillRandomWords.
   * @dev The reason plain setConfig cannot be used is that it is marked as onlyOwner. Since fulfillRandomWords
   * @dev is gated by authorized senders, and the authorized senders are not necessarily owners, the call will
   * @dev always fail if the caller of fulfillRandomWords is not the owner, which is not what we want.
   * @param minimumRequestConfirmations global min for request confirmations
   * @param maxGasLimit global max for request gas limit
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @param feeConfig fee tier configuration
   */
  function _setConfig(
    uint16 minimumRequestConfirmations,
    uint32 maxGasLimit,
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    FeeConfig memory feeConfig
  ) private {
    s_vrfCoordinator.setConfig(
      minimumRequestConfirmations,
      maxGasLimit,
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      feeConfig
    );
  }

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @param subId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subId) external onlyOwner {
    s_vrfCoordinator.ownerCancelSubscription(subId);
  }

  /**
   * @notice Recover link sent with transfer instead of transferAndCall.
   * @param to address to send link to
   */
  function recoverFunds(address to) external onlyOwner {
    s_vrfCoordinator.recoverFunds(to);
  }

  /**
   * @notice Get all relevant configs from the VRF coordinator.
   * @dev This is done in a separate function to avoid the "stack too deep" issue
   * @dev when too many local variables are in the same scope.
   * @return Config struct containing all relevant configs from the VRF coordinator.
   */
  function _getConfigs() private view returns (Config memory) {
    (
      uint16 minimumRequestConfirmations,
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint32 gasAfterPaymentCalculation
    ) = s_vrfCoordinator.getConfig();
    (
      uint32 fulfillmentFlatFeeLinkPPMTier1,
      uint32 fulfillmentFlatFeeLinkPPMTier2,
      uint32 fulfillmentFlatFeeLinkPPMTier3,
      uint32 fulfillmentFlatFeeLinkPPMTier4,
      uint32 fulfillmentFlatFeeLinkPPMTier5,
      uint24 reqsForTier2,
      uint24 reqsForTier3,
      uint24 reqsForTier4,
      uint24 reqsForTier5
    ) = s_vrfCoordinator.getFeeConfig();
    int256 fallbackWeiPerUnitLink = s_vrfCoordinator.getFallbackWeiPerUnitLink();
    return
      Config({
        minimumRequestConfirmations: minimumRequestConfirmations,
        maxGasLimit: maxGasLimit,
        stalenessSeconds: stalenessSeconds,
        gasAfterPaymentCalculation: gasAfterPaymentCalculation,
        fallbackWeiPerUnitLink: fallbackWeiPerUnitLink,
        feeConfig: FeeConfig({
          fulfillmentFlatFeeLinkPPMTier1: fulfillmentFlatFeeLinkPPMTier1,
          fulfillmentFlatFeeLinkPPMTier2: fulfillmentFlatFeeLinkPPMTier2,
          fulfillmentFlatFeeLinkPPMTier3: fulfillmentFlatFeeLinkPPMTier3,
          fulfillmentFlatFeeLinkPPMTier4: fulfillmentFlatFeeLinkPPMTier4,
          fulfillmentFlatFeeLinkPPMTier5: fulfillmentFlatFeeLinkPPMTier5,
          reqsForTier2: reqsForTier2,
          reqsForTier3: reqsForTier3,
          reqsForTier4: reqsForTier4,
          reqsForTier5: reqsForTier5
        })
      });
  }

  /**
   * @notice Fulfill a randomness request
   * @param proof contains the proof and randomness
   * @param rc request commitment pre-image, committed to at request time
   */
  function fulfillRandomWords(
    VRFTypes.Proof memory proof,
    VRFTypes.RequestCommitment memory rc
  ) external validateAuthorizedSender {
    uint256 requestId = _requestIdFromProof(proof.pk, proof.seed);

    // Get current configs to restore them to original values after
    // calling _setConfig.
    Config memory cfg = _getConfigs();

    // call _setConfig with the appropriate params in order to fulfill
    // an accidentally-underfunded request.
    _setConfig(
      cfg.minimumRequestConfirmations,
      cfg.maxGasLimit,
      1, // stalenessSeconds
      0, // gasAfterPaymentCalculation
      MAX_INT256, // fallbackWeiPerUnitLink
      FeeConfig({
        fulfillmentFlatFeeLinkPPMTier1: 0,
        fulfillmentFlatFeeLinkPPMTier2: 0,
        fulfillmentFlatFeeLinkPPMTier3: 0,
        fulfillmentFlatFeeLinkPPMTier4: 0,
        fulfillmentFlatFeeLinkPPMTier5: 0,
        reqsForTier2: 0,
        reqsForTier3: 0,
        reqsForTier4: 0,
        reqsForTier5: 0
      })
    );

    s_vrfCoordinator.fulfillRandomWords(proof, rc);

    // reset configuration back to old values.
    _setConfig(
      cfg.minimumRequestConfirmations,
      cfg.maxGasLimit,
      cfg.stalenessSeconds,
      cfg.gasAfterPaymentCalculation,
      cfg.fallbackWeiPerUnitLink,
      cfg.feeConfig
    );

    emit RandomWordsForced(requestId, rc.subId, rc.sender);
  }

  /**
   * @notice Concrete implementation of AuthorizedReceiver
   * @return bool of whether sender is authorized
   */
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return owner() == msg.sender;
  }

  /**
   * @notice Returns the request for corresponding to the given public key and proof seed.
   * @param publicKey the VRF public key associated with the proof
   * @param proofSeed the proof seed
   * @dev Refer to VRFCoordinatorV2.getRandomnessFromProof for original implementation.
   */
  function _requestIdFromProof(uint256[2] memory publicKey, uint256 proofSeed) private view returns (uint256) {
    bytes32 keyHash = s_vrfCoordinator.hashOfKey(publicKey);
    uint256 requestId = uint256(keccak256(abi.encode(keyHash, proofSeed)));
    return requestId;
  }
}
