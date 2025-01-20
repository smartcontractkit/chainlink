// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IRewardManager} from "./interfaces/IRewardManager.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {Common} from "../libraries/Common.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title RewardManager
 * @author Michael Fletcher
 * @author Austin Born
 * @notice This contract will be used to reward any configured recipients within a pool. Recipients will receive a share of their pool relative to their configured weight.
 */
contract RewardManager is IRewardManager, ConfirmedOwner, ITypeAndVersion {
  using SafeERC20 for IERC20;

  // @dev The mapping of total fees collected for a particular pot: s_totalRewardRecipientFees[poolId]
  mapping(bytes32 => uint256) public s_totalRewardRecipientFees;

  // @dev The mapping of fee balances for each pot last time the recipient claimed: s_totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient]
  mapping(bytes32 => mapping(address => uint256)) public s_totalRewardRecipientFeesLastClaimedAmounts;

  // @dev The mapping of RewardRecipient weights for a particular poolId: s_rewardRecipientWeights[poolId][rewardRecipient].
  mapping(bytes32 => mapping(address => uint256)) public s_rewardRecipientWeights;

  // @dev Keep track of the reward recipient weights that have been set to prevent duplicates
  mapping(bytes32 => bool) public s_rewardRecipientWeightsSet;

  // @dev Store a list of pool ids that have been registered, to make off chain lookups easier
  bytes32[] public s_registeredPoolIds;

  // @dev The address for the LINK contract
  address public immutable i_linkAddress;

  // The total weight of all RewardRecipients. 1e18 = 100% of the pool fees
  uint64 private constant PERCENTAGE_SCALAR = 1e18;

  // The fee manager address
  address public s_feeManagerAddress;

  // @notice Thrown whenever the RewardRecipient weights are invalid
  error InvalidWeights();

  // @notice Thrown when any given address is invalid
  error InvalidAddress();

  // @notice Thrown when the pool id is invalid
  error InvalidPoolId();

  // @notice Thrown when the calling contract is not within the authorized contracts
  error Unauthorized();

  // @notice Thrown when getAvailableRewardPoolIds parameters are incorrectly set
  error InvalidPoolLength();

  // Events emitted upon state change
  event RewardRecipientsUpdated(bytes32 indexed poolId, Common.AddressAndWeight[] newRewardRecipients);
  event RewardsClaimed(bytes32 indexed poolId, address indexed recipient, uint192 quantity);
  event FeeManagerUpdated(address newFeeManagerAddress);
  event FeePaid(FeePayment[] payments, address payer);

  /**
   * @notice Constructor
   * @param linkAddress address of the wrapped LINK token
   */
  constructor(address linkAddress) ConfirmedOwner(msg.sender) {
    //ensure that the address ia not zero
    if (linkAddress == address(0)) revert InvalidAddress();

    i_linkAddress = linkAddress;
  }

  // @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "RewardManager 1.1.0";
  }

  // @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == this.onFeePaid.selector;
  }

  modifier onlyOwnerOrFeeManager() {
    if (msg.sender != owner() && msg.sender != s_feeManagerAddress) revert Unauthorized();
    _;
  }

  modifier onlyOwnerOrRecipientInPool(bytes32 poolId) {
    if (msg.sender != owner() && s_rewardRecipientWeights[poolId][msg.sender] == 0) revert Unauthorized();
    _;
  }

  modifier onlyFeeManager() {
    if (msg.sender != s_feeManagerAddress) revert Unauthorized();
    _;
  }

  /// @inheritdoc IRewardManager
  function onFeePaid(FeePayment[] calldata payments, address payer) external override onlyFeeManager {
    uint256 totalFeeAmount;
    for (uint256 i; i < payments.length; ++i) {
      unchecked {
        //the total amount for any ERC-20 asset cannot exceed 2^256 - 1
        //see https://github.com/OpenZeppelin/openzeppelin-contracts/blob/36bf1e46fa811f0f07d38eb9cfbc69a955f300ce/contracts/token/ERC20/ERC20.sol#L266
        //for example implementation.
        s_totalRewardRecipientFees[payments[i].poolId] += payments[i].amount;

        //tally the total payable fees
        totalFeeAmount += payments[i].amount;
      }
    }

    //transfer the fees to this contract
    IERC20(i_linkAddress).safeTransferFrom(payer, address(this), totalFeeAmount);

    emit FeePaid(payments, payer);
  }

  /// @inheritdoc IRewardManager
  function claimRewards(bytes32[] memory poolIds) external override {
    _claimRewards(msg.sender, poolIds);
  }

  // wrapper impl for claimRewards
  function _claimRewards(address recipient, bytes32[] memory poolIds) internal returns (uint256) {
    //get the total amount claimable for this recipient
    uint256 claimAmount;

    //loop and claim all the rewards in the poolId pot
    for (uint256 i; i < poolIds.length; ++i) {
      //get the poolId to be claimed
      bytes32 poolId = poolIds[i];

      //get the total fees for the pot
      uint256 totalFeesInPot = s_totalRewardRecipientFees[poolId];

      unchecked {
        //avoid unnecessary storage reads if there's no fees in the pot
        if (totalFeesInPot == 0) continue;

        //get the claimable amount for this recipient, this calculation will never exceed the amount in the pot
        uint256 claimableAmount = totalFeesInPot - s_totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient];

        //calculate the recipients share of the fees, which is their weighted share of the difference between the last amount they claimed and the current amount in the pot. This can never be more than the total amount in existence
        uint256 recipientShare = (claimableAmount * s_rewardRecipientWeights[poolId][recipient]) / PERCENTAGE_SCALAR;

        //if there's no fees to claim, continue as there's nothing to update
        if (recipientShare == 0) continue;

        //keep track of the total amount claimable, this can never be more than the total amount in existence
        claimAmount += recipientShare;

        //set the current total amount of fees in the pot as it's used to calculate future claims
        s_totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient] = totalFeesInPot;

        //emit event if the recipient has rewards to claim
        emit RewardsClaimed(poolIds[i], recipient, uint192(recipientShare));
      }
    }

    //check if there's any rewards to claim in the given poolId
    if (claimAmount != 0) {
      //transfer the reward to the recipient
      IERC20(i_linkAddress).safeTransfer(recipient, claimAmount);
    }

    return claimAmount;
  }

  /// @inheritdoc IRewardManager
  function setRewardRecipients(
    bytes32 poolId,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external override onlyOwnerOrFeeManager {
    //revert if there are no recipients to set
    if (rewardRecipientAndWeights.length == 0) revert InvalidAddress();

    //check that the weights have not been previously set
    if (s_rewardRecipientWeightsSet[poolId]) revert InvalidPoolId();

    //keep track of the registered poolIds to make off chain lookups easier
    s_registeredPoolIds.push(poolId);

    //keep track of which pools have had their reward recipients set
    s_rewardRecipientWeightsSet[poolId] = true;

    //set the reward recipients, this will only be called once and contain the full set of RewardRecipients with a total weight of 100%
    _setRewardRecipientWeights(poolId, rewardRecipientAndWeights, PERCENTAGE_SCALAR);

    emit RewardRecipientsUpdated(poolId, rewardRecipientAndWeights);
  }

  function _setRewardRecipientWeights(
    bytes32 poolId,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights,
    uint256 expectedWeight
  ) internal {
    //we can't update the weights if it contains duplicates
    if (Common._hasDuplicateAddresses(rewardRecipientAndWeights)) revert InvalidAddress();

    //loop all the reward recipients and validate the weight and address
    uint256 totalWeight;
    for (uint256 i; i < rewardRecipientAndWeights.length; ++i) {
      //get the weight
      uint256 recipientWeight = rewardRecipientAndWeights[i].weight;
      //get the address
      address recipientAddress = rewardRecipientAndWeights[i].addr;

      //ensure the reward recipient address is not zero
      if (recipientAddress == address(0)) revert InvalidAddress();

      //save/overwrite the weight for the reward recipient
      s_rewardRecipientWeights[poolId][recipientAddress] = recipientWeight;

      unchecked {
        //keep track of the cumulative weight, this cannot overflow as the total weight is restricted at 1e18
        totalWeight += recipientWeight;
      }
    }

    //if total weight is not met, the fees will either be under or over distributed
    if (totalWeight != expectedWeight) revert InvalidWeights();
  }

  /// @inheritdoc IRewardManager
  function updateRewardRecipients(
    bytes32 poolId,
    Common.AddressAndWeight[] calldata newRewardRecipients
  ) external override onlyOwner {
    //create an array of poolIds to pass to _claimRewards if required
    bytes32[] memory poolIds = new bytes32[](1);
    poolIds[0] = poolId;

    //loop all the reward recipients and claim their rewards before updating their weights
    uint256 existingTotalWeight;
    for (uint256 i; i < newRewardRecipients.length; ++i) {
      //get the address
      address recipientAddress = newRewardRecipients[i].addr;
      //get the existing weight
      uint256 existingWeight = s_rewardRecipientWeights[poolId][recipientAddress];

      //if a recipient is updated, the rewards must be claimed first as they can't claim previous fees at the new weight
      _claimRewards(newRewardRecipients[i].addr, poolIds);

      unchecked {
        //keep tally of the weights so that the expected collective weight is known
        existingTotalWeight += existingWeight;
      }
    }

    //update the reward recipients, if the new collective weight isn't equal to the previous collective weight, the fees will either be under or over distributed
    _setRewardRecipientWeights(poolId, newRewardRecipients, existingTotalWeight);

    //emit event
    emit RewardRecipientsUpdated(poolId, newRewardRecipients);
  }

  /// @inheritdoc IRewardManager
  function payRecipients(bytes32 poolId, address[] calldata recipients) external onlyOwnerOrRecipientInPool(poolId) {
    //convert poolIds to an array to match the interface of _claimRewards
    bytes32[] memory poolIdsArray = new bytes32[](1);
    poolIdsArray[0] = poolId;

    //loop each recipient and claim the rewards for each of the pools and assets
    for (uint256 i; i < recipients.length; ++i) {
      _claimRewards(recipients[i], poolIdsArray);
    }
  }

  /// @inheritdoc IRewardManager
  function setFeeManager(address newFeeManagerAddress) external onlyOwner {
    if (newFeeManagerAddress == address(0)) revert InvalidAddress();

    s_feeManagerAddress = newFeeManagerAddress;

    emit FeeManagerUpdated(newFeeManagerAddress);
  }

  /// @inheritdoc IRewardManager
  function getAvailableRewardPoolIds(
    address recipient,
    uint256 startIndex,
    uint256 endIndex
  ) external view returns (bytes32[] memory) {
    //get the length of the pool ids which we will loop through and potentially return
    uint256 registeredPoolIdsLength = s_registeredPoolIds.length;

    uint256 lastIndex = endIndex > registeredPoolIdsLength ? registeredPoolIdsLength : endIndex;

    if (startIndex > lastIndex) revert InvalidPoolLength();

    //create a new array with the maximum amount of potential pool ids
    bytes32[] memory claimablePoolIds = new bytes32[](lastIndex - startIndex);
    //we want the pools which a recipient has funds for to be sequential, so we need to keep track of the index
    uint256 poolIdArrayIndex;

    //loop all the pool ids, and check if the recipient has a registered weight and a claimable amount
    for (uint256 i = startIndex; i < lastIndex; ++i) {
      //get the poolId
      bytes32 poolId = s_registeredPoolIds[i];

      //if the recipient has a weight, they are a recipient of this poolId
      if (s_rewardRecipientWeights[poolId][recipient] != 0) {
        //get the total in this pool
        uint256 totalPoolAmount = s_totalRewardRecipientFees[poolId];
        //if the recipient has any LINK, then add the poolId to the array
        unchecked {
          //s_totalRewardRecipientFeesLastClaimedAmounts can never exceed total pool amount, and the number of pools can't exceed the max array length
          if (totalPoolAmount - s_totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient] != 0) {
            claimablePoolIds[poolIdArrayIndex++] = poolId;
          }
        }
      }
    }

    return claimablePoolIds;
  }
}
