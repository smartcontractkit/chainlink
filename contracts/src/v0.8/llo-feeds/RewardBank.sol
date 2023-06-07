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
 * @notice This contract will be used to reward NOPs for generating reports. Any reports verified by end users will have the fee split amongst the NOPs that are part of the DON that produced the report.
 */
contract RewardBank is IRewardBank, ConfirmedOwner, TypeAndVersionInterface{
    using SafeERC20 for IERC20;

    // @dev The mapping of total fees collected for a particular pot: totalNOPFees[configDigest][asset]
    mapping(bytes32 => mapping(address => uint256)) private totalNOPFees;

    // @dev The mapping of NOP configs for a particular feed config: NOPs[configDigest][NOP]
    mapping(bytes32 => NOPAndWeight[]) private NOPsAndWeights;

    // @dev The address of the verifier proxy
    address private immutable i_verifierProxyAddr;

    // The total weight of all NOPs in a DON
    uint8 private constant TOTAL_WEIGHT = 100;

    // @notice Thrown whenever a zero address is passed
    error ZeroAddress();

    // @notice Thrown whenever the NOP weights are invalid
    error InvalidWeights();

    // @notice Thrown whenever their are no fees to claim
    error ZeroFees();

    // @param verifierProxyAddr The address of the VerifierProxy contract
    constructor(address verifierProxyAddr) ConfirmedOwner(msg.sender) {
        if (verifierProxyAddr == address(0)) revert ZeroAddress();
        i_verifierProxyAddr = verifierProxyAddr;
    }

    // @notice This modifier is used to ensure that only the VerifierProxy contract can call the function
    modifier onlyAdminOrProxy() {
        if (msg.sender != i_verifierProxyAddr && msg.sender != owner()) revert ZeroAddress();
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
    function onFeePaid(bytes32 configDigest, Asset calldata fee) external override onlyAdminOrProxy {
        //update the total fees collected for this pot
        unchecked {
            // the total amount for any ERC20 token cannot exceed 2^256 - 1
            totalNOPFees[configDigest][fee.assetAddress] += fee.amount;
        }

        //transfer the fee to the contract
        IERC20(fee.assetAddress).safeTransferFrom(msg.sender, address(this), fee.amount);
    }

    // @inheritdoc IRewardBank
    function payNOPs(bytes32 configDigest, address assetAddress) external override {
        //get the total fees collected for this pot
        uint256 totalFees = totalNOPFees[configDigest][assetAddress];

        //ensure there are fees to pay
        if (totalFees == 0) revert ZeroFees();

        //variables used within the loop
        uint256 totalNOPs = NOPsAndWeights[configDigest].length;
        uint256 NOPShare;
        NOPAndWeight storage NOP;

        //Dust will be lost setting this to 0, but the value of that dust is significantly less than the gas would be if we were to account for the dust.
        //Setting this here will prevent us needing to use the nonReentrancy modifier
        totalNOPFees[configDigest][assetAddress] = 0;

        //loop all the NOPs and pay them their share of the fees
        for (uint256 i; i < totalNOPs;) {
            //get the NOP for address and weight
            NOP = NOPsAndWeights[configDigest][i];

            //calculate the NOP share of the fees.
            NOPShare = (totalFees * NOP.weight) / TOTAL_WEIGHT;

            //transfer the NOP share of the fees to the NOP
            IERC20(assetAddress).safeTransfer(NOP.NOPAddress, NOPShare);

            unchecked {
                // there will never be enough ops in a config for i to overflow
                ++i;
            }
        }
    }


    // @inheritdoc IRewardBank
    function setNOPs(bytes32 configDigest, NOPAndWeight[] calldata NOPAndWeights) external override onlyAdminOrProxy {
        if (NOPAndWeights.length == 0) revert ZeroAddress();

        //loop all the NOPs and validate the weight and address
        uint256 totalWeight;
        for (uint256 i; i < NOPAndWeights.length;) {
            unchecked {
                //keep track of the cumulative weight, this cannot overflow as the passed in weight is 16 bits
                totalWeight += NOPAndWeights[i].weight;

                // there will never be enough ops in a config for i to overflow
                ++i;
            }

            //ensure the NOP address is not zero
            if (NOPAndWeights[i].NOPAddress == address(0)) revert ZeroAddress();

            //ensure NOP weight is not 0
            if (NOPAndWeights[i].weight == 0) revert InvalidWeights();

            //copy to contract storage
            NOPsAndWeights[configDigest][i] = NOPAndWeight(NOPAndWeights[i].NOPAddress, NOPAndWeights[i].weight);
        }

        // If total weight is not met, the fees will either be under or over distributed
        if(totalWeight != TOTAL_WEIGHT) revert InvalidWeights();
    }
}