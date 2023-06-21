// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IRewardBank} from "./interfaces/IRewardBank.sol";
import {IERC20} from "../shared/vendor/IERC20.sol";
import {SafeERC20} from "../shared/vendor/SafeERC20.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";

/*
 * @title FeeManager
 * @author Michael Fletcher
 * @notice This contract will be used to reward any configured recipients within a pool. Recipients will receive a share of their rewards relative to their configured weight.
 */
contract RewardBank is IRewardBank, ConfirmedOwner, TypeAndVersionInterface{
    using SafeERC20 for IERC20;

    // @dev The mapping of total fees collected for a particular pot: totalRewardRecipientFees[poolId][asset]
    mapping(bytes32 => mapping(address => uint256)) private totalRewardRecipientFees;

    // @dev The mapping of fee balances for each pot last time the recipient claimed: totalRewardRecipientFeesLastClaimedAmounts[poolId][asset][oracle]
    mapping(bytes32 => mapping(address => mapping(address => uint256))) private totalRewardRecipientFeesLastClaimedAmounts;

    // @dev The mapping of RewardRecipient weights for a particular poolId: rewardRecipientWeights[poolId][rewardRecipient]. Weights are stored in uint256 to optimize on calculations
    mapping(bytes32 => mapping(address => uint256)) private rewardRecipientWeights;

    // @dev Only authorized contracts can configure reward recipients
    mapping(address => bool) private authorizedContracts;

    // @dev Keep track of the reward recipient weights that have been set to calculate relative share of the pool
    mapping(bytes32 => bool) private rewardRecipientWeightsSet;

    // @dev Store a list of pool ids that have been registered, to make off chain lookups easier
    bytes32[] private registeredPoolIds;

    // @dev The address for the wrapped native contract
    address private immutable wrappedNativeAddr;

    // @dev The address for the link contract
    address private immutable linkAddr;

    // The total weight of all RewardRecipients. 1000 = 10% of the pool fees
    uint16 private constant PERCENTAGE_SCALAR = 10000;

    // @notice Thrown whenever the RewardRecipient weights are invalid
    error InvalidWeights();

    // @notice Thrown when any given address is invalid
    error InvalidAddress();

    // @notice Thrown when the pool id is invalid
    error InvalidPoolId();

    // @notice Thrown when the calling contract is not within the authorized contracts
    error Unauthorized();

    /**
     * @notice Constructor
     * @param _wrappedNativeAddr address of the wrapped native token
     * @param _wrappedLinkAddr address of the wrapped link token
      */
    constructor(address _wrappedNativeAddr, address _wrappedLinkAddr) ConfirmedOwner(msg.sender) {
        //ensure that the addresses are not zero
        if(_wrappedNativeAddr == address(0)) revert InvalidAddress();
        if(_wrappedLinkAddr == address(0)) revert InvalidAddress();

        wrappedNativeAddr = _wrappedNativeAddr;
        linkAddr = _wrappedLinkAddr;
    }

    modifier onlyAuthorizedContracts() {
        if(!authorizedContracts[msg.sender]) revert Unauthorized();
        _;
    }

    // @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external pure override returns (string memory) {
        return "RewardBank 0.0.1";
    }

    // @inheritdoc IERC165
    function supportsInterface(bytes4 interfaceId)
    external
    pure
    override
    returns (bool)
    {
        return interfaceId == this.onFeePaid.selector;
    }

    // @inheritdoc IRewardBank
    function addAuthorizedContract(address contractAddress) external onlyOwner {
        authorizedContracts[contractAddress] = true;
    }

    // @inheritdoc IRewardBank
    function removeAuthorizedContract(address contractAddress) external onlyOwner {
        delete authorizedContracts[contractAddress];
    }

    // @inheritdoc IRewardBank
    function onFeePaid(bytes32 poolId, Asset calldata fee) external override {
        //fee must either be wrapped native or link
        if(fee.assetAddress != wrappedNativeAddr && fee.assetAddress != linkAddr) revert InvalidAddress();

        //update the total fees collected for this pot
        unchecked {
            //the total amount for any ERC20 asset cannot exceed 2^256 - 1
            totalRewardRecipientFees[poolId][fee.assetAddress] += fee.amount;
        }

        //transfer the fee to this contract
        IERC20(fee.assetAddress).safeTransferFrom(msg.sender, address(this), fee.amount);
    }

    // @inheritdoc IRewardBank
    function updateRewardRecipients(bytes32[] calldata poolIds, RewardRecipientAndWeight[] calldata newRewardRecipients) external override onlyOwner {
        //the interface to _claimRewards and _setRewardRecipientWeights requires an array of calldata poolIds. Forcing this to be passed as an array will improve code readability and optimize on gas
        if(newRewardRecipients.length != 1) revert InvalidPoolId();

        //get the poolId
        bytes32 poolId = poolIds[0];

        //loop all the reward recipients with the following rules:
        //if existing weight is 0, a new user is being added
        //if the existing weight is greater than 0, the user is being updated or removed
        //if the new weight is 0, the user will be removed
        //if the new weight is greater than 0, the users weight will be updated
        //the sum of the existing weights must equal to the new collective weight
        uint256 existingTotalWeight;
        for (uint256 i; i < newRewardRecipients.length;) {
            //get the address
            address recipientAddress = newRewardRecipients[i].rewardRecipientAddress;
            //get the existing weight
            uint256 existingWeight = rewardRecipientWeights[poolId][recipientAddress];

            //if existing weight is 0, a new recipient is being added so we must set their totalRewardRecipientFeesLastClaimedAmounts to the current amount in the pot, which will prevent them having a claim over previous fees
            if(existingWeight == 0) {
                totalRewardRecipientFeesLastClaimedAmounts[poolId][wrappedNativeAddr][recipientAddress] = totalRewardRecipientFees[poolId][wrappedNativeAddr];
                totalRewardRecipientFeesLastClaimedAmounts[poolId][linkAddr][recipientAddress] = totalRewardRecipientFees[poolId][linkAddr];
            } else {
                //if their existing weight is set, then their weight is being updated, so we should claim rewards to ensure the new weight doesn't apply to historic fees
                _claimRewards(newRewardRecipients[i].rewardRecipientAddress, poolIds, wrappedNativeAddr);
                _claimRewards(newRewardRecipients[i].rewardRecipientAddress, poolIds, linkAddr);
            }

            //if the new weight is 0, the recipient is being removed, so we can remove their totalRewardRecipientFeesLastClaimedAmounts as their rewards have been claimed
            if(newRewardRecipients[i].weight == 0) {
                delete totalRewardRecipientFeesLastClaimedAmounts[poolId][wrappedNativeAddr][recipientAddress];
                delete totalRewardRecipientFeesLastClaimedAmounts[poolId][linkAddr][recipientAddress];
            }

            unchecked {
                //keep tally of the weights so we know the expected collective weight
                existingTotalWeight += existingWeight;
                //there will never be enough reward recipients for i to overflow
                ++i;
            }
        }

        //update the reward recipients, if the new collective weight isn't equal to the previous collective weight, the tx will be reverted
        _setRewardRecipientWeights(poolIds[0], newRewardRecipients, existingTotalWeight);
    }


    // @inheritdoc IRewardBank
    function setRewardRecipients(bytes32 poolId, RewardRecipientAndWeight[] calldata rewardRecipientAndWeights) external override onlyAuthorizedContracts {
        //keep track of the registered poolIds to make offchain lookups easier
        registeredPoolIds.push(poolId);

        //set the reward recipient, this will only be called once and contain the full set of RewardRecipients with a total weight of 100%
        _setRewardRecipientWeights(poolId, rewardRecipientAndWeights, PERCENTAGE_SCALAR);
    }

    // @inheritdoc IRewardBank
    function claimRewards(bytes32[] calldata poolIds) external override {
        //claim the rewards for each asset
        _claimRewards(msg.sender, poolIds, wrappedNativeAddr);
        _claimRewards(msg.sender, poolIds, linkAddr);
    }


    // wrapper impl for setRewardRecipients
    function _setRewardRecipientWeights(bytes32 poolId, RewardRecipientAndWeight[] calldata rewardRecipientAndWeights, uint256 expectedWeight) internal {
        //loop all the reward recipients and validate the weight and address
        uint256 totalWeight;
        for (uint256 i; i < rewardRecipientAndWeights.length;) {
            //get the weight a uint256 to save multiple autoboxing
            uint256 recipientWeight = rewardRecipientAndWeights[i].weight;
            //get the address
            address recipientAddress = rewardRecipientAndWeights[i].rewardRecipientAddress;

            //ensure the reward recipient address is not zero
            if (recipientAddress == address(0)) revert InvalidAddress();

            //save/overwrite the weight for the reward recipient
            rewardRecipientWeights[poolId][recipientAddress] = recipientWeight;

            unchecked {
                //keep track of the cumulative weight, this cannot overflow as the passed in weight is 16 bits
                totalWeight += recipientWeight;
                //there will never be enough recipients for i to overflow
                ++i;
            }
        }

        //if total weight is not met, the fees will either be under or over distributed
        if(totalWeight != expectedWeight) revert InvalidWeights();
    }


    // wrapper impl for claimRewards
    function _claimRewards(address claimant, bytes32[] calldata poolIds, address assetAddress) internal {
        //get the total amount claimable for this recipient
        uint256 claimAmount;

        //loop and claim all the rewards in the poolId pot for the given asset
        for (uint256 i; i < poolIds.length;) {
            //get the poolId we're claiming for
            bytes32 poolId = poolIds[i];

            //get the total fees for the pot
            uint256 totalFeesInPot = totalRewardRecipientFees[poolId][assetAddress];

            unchecked {
                //get the claimable amount for this recipient, this calculation will never exceed the amount in the pot
                uint256 claimableAmount = totalFeesInPot - totalRewardRecipientFeesLastClaimedAmounts[poolId][assetAddress][claimant];

                //if there's no fees to claim, continue as there's nothing to update
                if(claimableAmount == 0) continue;

                //calculate the recipients share of the fees, which is their weighted share of the difference between the last amount they claimed and the current amount in the pot. This can never be more than the total amount in existence
                uint256 recipientShare = claimableAmount * rewardRecipientWeights[poolId][msg.sender] / PERCENTAGE_SCALAR;

                //keep track of the total amount claimable, this can never be more than the total amount in existence
                claimAmount += recipientShare;

                //set the current total amount of fees in the pot as it's used to calculate future claims
                totalRewardRecipientFeesLastClaimedAmounts[poolId][assetAddress][claimant] = totalFeesInPot;

                //there will never be enough poolIds for i to overflow
                ++i;
            }
        }

        //check if there's any rewards to claim in the given poolId
        if(claimAmount > 0) {
            //transfer the reward to the recipient
            IERC20(assetAddress).safeTransfer(claimant, claimAmount);
        }
    }

    // @inheritdoc IRewardBank
    function getAvailableRewardPoolIds(address recipient) external view returns (bytes32[] memory) {
        //get the length of the pool ids which we will loop through and potentially return
        uint256 registeredPoolIdsLength = registeredPoolIds.length;

        //create a new array with the maximum amount of potential pool ids
        bytes32[] memory claimablePoolIds = new bytes32[](registeredPoolIdsLength);
        //we want the pools which a recipient has funds for to be sequential, so we need to keep track of the index
        uint256 poolIdArrayIndex;

        //loop all the pool ids, and check if the recipient has a registered weight and a claimable amount
        for(uint256 i; i <registeredPoolIdsLength;) {
            //get the poolId
            bytes32 poolId = registeredPoolIds[i];
            //if the recipient has a weight, they are a recipient of this poolId
            if(rewardRecipientWeights[poolId][recipient] > 0) {
                //if the recipient has any link or native, then add the poolId to the array
                if(totalRewardRecipientFees[poolId][linkAddr] > 0 || totalRewardRecipientFees[poolId][wrappedNativeAddr] > 0) {
                    claimablePoolIds[poolIdArrayIndex] = poolId;
                    unchecked {
                        //there will never be enough pool ids for i to overflow
                        ++poolIdArrayIndex;
                    }
                }
            }

            unchecked {
                //there will never be enough poolIds for i to overflow
                ++i;
            }
        }

        return claimablePoolIds;
    }
}