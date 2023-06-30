// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IRewardManager} from "./interfaces/IRewardManager.sol";
import {IERC20} from "../shared/vendor/IERC20.sol";
import {SafeERC20} from "../shared/vendor/SafeERC20.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";
import {Common} from "../libraries/internal/Common.sol";

/*
 * @title FeeManager
 * @author Michael Fletcher
 * @notice This contract will be used to reward any configured recipients within a pool. Recipients will receive a share of their reward-manager relative to their configured weight.
 */
contract RewardManager is IRewardManager, ConfirmedOwner, TypeAndVersionInterface {
    using SafeERC20 for IERC20;

    // @dev The mapping of total fees collected for a particular pot: totalRewardRecipientFees[poolId]
    mapping(bytes32 => uint256) private totalRewardRecipientFees;

    // @dev The mapping of fee balances for each pot last time the recipient claimed: totalRewardRecipientFeesLastClaimedAmounts[poolId][oracle]
    mapping(bytes32 => mapping(address => uint256)) private totalRewardRecipientFeesLastClaimedAmounts;

    // @dev The mapping of RewardRecipient weights for a particular poolId: rewardRecipientWeights[poolId][rewardRecipient]. Weights are stored in uint256 to optimize on calculations
    mapping(bytes32 => mapping(address => uint256)) private rewardRecipientWeights;

    // @dev Only authorized contracts can configure reward recipients
    mapping(address => bool) private authorizedContracts;

    // @dev Keep track of the reward recipient weights that have been set to calculate relative share of the pool
    mapping(bytes32 => bool) private rewardRecipientWeightsSet;

    // @dev Store a list of pool ids that have been registered, to make off chain lookups easier
    bytes32[] private registeredPoolIds;

    // @dev The address for the link contract
    address private immutable LINK_ADDRESS;

    // The total weight of all RewardRecipients. 1000 = 10% of the pool fees
    uint16 private constant PERCENTAGE_SCALAR = 10000;

    // The verifier proxy address
    address private verifierProxyAddress;

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
     * @param linkAddr address of the wrapped link token
      */
    constructor(address linkAddr) ConfirmedOwner(msg.sender) {
        //ensure that the addresses are not zero
        if (linkAddr == address(0)) revert InvalidAddress();

        LINK_ADDRESS = linkAddr;
    }

    // @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external pure override returns (string memory) {
        return "RewardManager 0.0.1";
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

    // @inheritdoc IRewardManager
    function addAuthorizedContract(address contractAddress) external onlyOwner {
        authorizedContracts[contractAddress] = true;
    }

    // @inheritdoc IRewardManager
    function removeAuthorizedContract(address contractAddress) external onlyOwner {
        delete authorizedContracts[contractAddress];
    }

    modifier onlyOwnerOrProxy() {
        if (msg.sender != verifierProxyAddress && msg.sender != owner()) revert Unauthorized();
        _;
    }

    modifier onlyOwnerOrRecipientInPool(bytes32 poolId) {
        if (msg.sender != owner() && rewardRecipientWeights[poolId][msg.sender] == 0) revert Unauthorized();
        _;
    }

    // @inheritdoc IRewardManager
    function onFeePaid(bytes32 poolId, address payee, Common.Asset calldata fee) external override {
        //fee must be link
        if (fee.assetAddress != LINK_ADDRESS) revert InvalidAddress();

        //update the total fees collected for this pot
        unchecked {
            //the total amount for any ERC20 asset cannot exceed 2^256 - 1
            totalRewardRecipientFees[poolId] += fee.amount;
        }

        //transfer the fee to this contract
        IERC20(LINK_ADDRESS).safeTransferFrom(payee, address(this), fee.amount);
    }

    // @inheritdoc IRewardManager
    function updateRewardRecipients(bytes32[] calldata poolIds, Common.AddressAndWeight[] calldata newRewardRecipients) external override onlyOwner {
        //the interface to _claimRewards and _setRewardRecipientWeights requires an array of calldata poolIds. Forcing this to be passed as an array will improve code readability and optimize on gas
        if (newRewardRecipients.length != 1) revert InvalidPoolId();

        //get the poolId
        bytes32 poolId = poolIds[0];

        //loop all the reward recipients. The sum of the existing weights must equal to the new collective weight
        uint256 existingTotalWeight;
        for (uint256 i; i < newRewardRecipients.length;) {
            //get the address
            address recipientAddress = newRewardRecipients[i].addr;
            //get the existing weight
            uint256 existingWeight = rewardRecipientWeights[poolId][recipientAddress];

            //if the existing weight is 0, the recipient isn't part of this configuration
            if(existingWeight == 0) revert InvalidAddress();

            //if existing weight is 0, a new recipient is being added so we must set their totalRewardRecipientFeesLastClaimedAmounts to the current amount in the pot, which will prevent them having a claim over previous fees
            if (existingWeight == 0) {
                totalRewardRecipientFeesLastClaimedAmounts[poolId][recipientAddress] = totalRewardRecipientFees[poolId];
            } else {
                //if their existing weight is set, then their weight is being updated, so we should claim reward-manager to ensure the new weight doesn't apply to historic fees
                _claimRewards(newRewardRecipients[i].addr, poolIds);
            }

            //if the new weight is 0, the recipient is being removed, so we can remove their totalRewardRecipientFeesLastClaimedAmounts as their reward-manager have been claimed
            if (newRewardRecipients[i].weight == 0) {
                delete totalRewardRecipientFeesLastClaimedAmounts[poolId][recipientAddress];
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


    // @inheritdoc IRewardManager
    function setRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] calldata rewardRecipientAndWeights) external override onlyOwnerOrProxy {
        //revert if there's no recipients to set
        if (rewardRecipientAndWeights.length == 0) revert InvalidAddress();

        //keep track of the registered poolIds to make offchain lookups easier
        registeredPoolIds.push(poolId);

        //set the reward recipient, this will only be called once and contain the full set of RewardRecipients with a total weight of 100%
        _setRewardRecipientWeights(poolId, rewardRecipientAndWeights, PERCENTAGE_SCALAR);
    }

    // @inheritdoc IRewardManager
    function payRecipients(bytes32 poolId, address[] calldata recipients) external onlyOwnerOrRecipientInPool(poolId) {

        //convert poolIds to an array to match the interface of _claimRewards
        bytes32[] memory poolIdsArray = new bytes32[](1);
        poolIdsArray[0] = poolId;

        //loop each recipient and claim the reward-manager for each of the pools and assets
        for (uint256 i; i < recipients.length;) {
            _claimRewards(recipients[i], poolIdsArray);

            unchecked {
                //there will never be enough recipients for i to overflow
                ++i;
            }
        }
    }

    // @inheritdoc IRewardManager
    function claimRewards(bytes32[] memory poolIds) external override {
        _claimRewards(msg.sender, poolIds);
    }


    // wrapper impl for setRewardRecipients
    function _setRewardRecipientWeights(bytes32 poolId, Common.AddressAndWeight[] calldata rewardRecipientAndWeights, uint256 expectedWeight) internal {
        //loop all the reward recipients and validate the weight and address
        uint256 totalWeight;
        for (uint256 i; i < rewardRecipientAndWeights.length; ++i) {
            //get the weight a uint256 to save multiple autoboxing
            uint256 recipientWeight = rewardRecipientAndWeights[i].weight;
            //get the address
            address recipientAddress = rewardRecipientAndWeights[i].addr;

            //ensure the reward recipient address is not zero
            if (recipientAddress == address(0)) revert InvalidAddress();

            //ensure the weight is not zero
            if(recipientWeight == 0) revert InvalidWeights();

            //save/overwrite the weight for the reward recipient
            rewardRecipientWeights[poolId][recipientAddress] = recipientWeight;

            unchecked {
                //keep track of the cumulative weight, this cannot overflow as the passed in weight is 16 bits
                totalWeight += recipientWeight;
            }
        }

        //if total weight is not met, the fees will either be under or over distributed
        if (totalWeight != expectedWeight) revert InvalidWeights();
    }


    // wrapper impl for claimRewards
    function _claimRewards(address recipient, bytes32[] memory poolIds) internal {
        //get the total amount claimable for this recipient
        uint256 claimAmount;

        //loop and claim all the reward-manager in the poolId pot
        for (uint256 i; i < poolIds.length; ++i) {
            //get the poolId we're claiming for
            bytes32 poolId = poolIds[i];

            //get the total fees for the pot
            uint256 totalFeesInPot = totalRewardRecipientFees[poolId];

            unchecked {
                //get the claimable amount for this recipient, this calculation will never exceed the amount in the pot
                uint256 claimableAmount = totalFeesInPot - totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient];

                //if there's no fees to claim, continue as there's nothing to update
                if (claimableAmount == 0) continue;

                //calculate the recipients share of the fees, which is their weighted share of the difference between the last amount they claimed and the current amount in the pot. This can never be more than the total amount in existence
                uint256 recipientShare = claimableAmount * rewardRecipientWeights[poolId][recipient] / PERCENTAGE_SCALAR;

                //keep track of the total amount claimable, this can never be more than the total amount in existence
                claimAmount += recipientShare;

                //set the current total amount of fees in the pot as it's used to calculate future claims
                totalRewardRecipientFeesLastClaimedAmounts[poolId][recipient] = totalFeesInPot;
            }
        }

        //check if there's any reward-manager to claim in the given poolId
        if (claimAmount > 0) {
            //transfer the reward to the recipient
            IERC20(LINK_ADDRESS).safeTransfer(recipient, claimAmount);
        }
    }

    // @inheritdoc IRewardManager
    function getAvailableRewardPoolIds(address recipient) external view returns (bytes32[] memory) {
        //get the length of the pool ids which we will loop through and potentially return
        uint256 registeredPoolIdsLength = registeredPoolIds.length;

        //create a new array with the maximum amount of potential pool ids
        bytes32[] memory claimablePoolIds = new bytes32[](registeredPoolIdsLength);
        //we want the pools which a recipient has funds for to be sequential, so we need to keep track of the index
        uint256 poolIdArrayIndex;

        //loop all the pool ids, and check if the recipient has a registered weight and a claimable amount
        for (uint256 i; i < registeredPoolIdsLength;) {
            //get the poolId
            bytes32 poolId = registeredPoolIds[i];
            //if the recipient has a weight, they are a recipient of this poolId
            if (rewardRecipientWeights[poolId][recipient] > 0) {
                //if the recipient has any link, then add the poolId to the array
                if (totalRewardRecipientFees[poolId] > 0) {
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

    // @inheritdoc IRewardManager
    function updateBillingAddress(address newBillingAddress) external {
        //loop each pool and update each recipients weight to reflect the new billing address
        uint256 registeredPoolIdsLength = registeredPoolIds.length;
        for (uint256 i; i < registeredPoolIdsLength;) {
            //get the poolId
            bytes32 poolId = registeredPoolIds[i];

            //get the recipients weight
            uint256 recipientWeight = rewardRecipientWeights[poolId][msg.sender];

            //if the recipient has a weight, they are a recipient of this poolId
            if (recipientWeight > 0) {
                //update the recipients new billing address weight to the existing weight
                rewardRecipientWeights[poolId][newBillingAddress] = recipientWeight;
                //delete the old weight
                delete rewardRecipientWeights[poolId][msg.sender];
            }

            //loop the claimed reward for this recipients existing billing address and update against the new billing address
            totalRewardRecipientFeesLastClaimedAmounts[poolId][newBillingAddress] = totalRewardRecipientFeesLastClaimedAmounts[poolId][msg.sender];

            //delete the old claimed reward
            delete totalRewardRecipientFeesLastClaimedAmounts[poolId][msg.sender];

            unchecked {
            //there will never be enough poolIds for i to overflow
                ++i;
            }
        }
    }

    // @inheritdoc IRewardManager
    function setVerifierProxy(address newVerifierProxyAddress) external onlyOwner {
        verifierProxyAddress = newVerifierProxyAddress;
    }
}