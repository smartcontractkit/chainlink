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
 * @notice This contract will be used to reward NOPs for generating reports and any additional Orchestrators such as Chainlink Labs. Any reports verified by end users will have the fee split amongst the Payees and Orchestrators that are part of the DON that produced the report.
 */
contract RewardBank is IRewardBank, ConfirmedOwner, TypeAndVersionInterface{
    using SafeERC20 for IERC20;

    // @dev The mapping of total fees collected for a particular pot: totalPayeeFees[configDigest][asset]
    mapping(bytes32 => mapping(address => uint256)) private totalPayeeFees;

    // @dev The mapping of Payee configs for a particular feed config: Payees[configDigest][Payee]
    mapping(bytes32 => PayeeAndWeight[]) private PayeesAndWeights;

    // @dev The address of the verifier proxy
    address private immutable i_verifierProxyAddr;

    // The total weight of all Payees
    uint8 private constant TOTAL_WEIGHT = 100;

    // @notice Thrown whenever a zero address is passed
    error ZeroAddress();

    // @notice Thrown whenever the Payee weights are invalid
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
            totalPayeeFees[configDigest][fee.assetAddress] += fee.amount;
        }

        //transfer the fee to the contract
        IERC20(fee.assetAddress).safeTransferFrom(msg.sender, address(this), fee.amount);
    }

    // @inheritdoc IRewardBank
    function payPayees(bytes32 configDigest, address assetAddress) external override {
        //get the total fees collected for this pot
        uint256 totalFees = totalPayeeFees[configDigest][assetAddress];

        //ensure there are fees to pay
        if (totalFees == 0) revert ZeroFees();

        //variables used within the loop
        uint256 totalPayees = PayeesAndWeights[configDigest].length;
        uint256 PayeeShare;
        PayeeAndWeight memory Payee;

        //Dust will be lost setting this to 0, but the value of that dust is significantly less than the gas would be if we were to account for the dust.
        //Setting this here will prevent us needing to use the nonReentrancy modifier
        totalPayeeFees[configDigest][assetAddress] = 0;

        //loop all the Payees and pay them their share of the fees
        for (uint256 i; i < totalPayees;) {
            //get the Payee for address and weight
            Payee = PayeesAndWeights[configDigest][i];

            //calculate the Payee share of the fees.
            PayeeShare = (totalFees * Payee.weight) / TOTAL_WEIGHT;

            //transfer the Payee share of the fees to the Payee
            IERC20(assetAddress).safeTransfer(Payee.PayeeAddress, PayeeShare);

            unchecked {
                // there will never be enough ops in a config for i to overflow
                ++i;
            }
        }
    }


    // @inheritdoc IRewardBank
    function setPayees(bytes32 configDigest, PayeeAndWeight[] calldata PayeeAndWeights) external override onlyAdminOrProxy {
        if (PayeeAndWeights.length == 0) revert ZeroAddress();

        //loop all the Payees and validate the weight and address
        uint256 totalWeight;
        for (uint256 i; i < PayeeAndWeights.length;) {
            unchecked {
                //keep track of the cumulative weight, this cannot overflow as the passed in weight is 16 bits
                totalWeight += PayeeAndWeights[i].weight;

                // there will never be enough ops in a config for i to overflow
                ++i;
            }

            //ensure the Payee address is not zero
            if (PayeeAndWeights[i].PayeeAddress == address(0)) revert ZeroAddress();

            //ensure Payee weight is not 0
            if (PayeeAndWeights[i].weight == 0) revert InvalidWeights();

            //copy to contract storage
            PayeesAndWeights[configDigest][i] = PayeeAndWeight(PayeeAndWeights[i].PayeeAddress, PayeeAndWeights[i].weight);
        }

        // If total weight is not met, the fees will either be under or over distributed
        if(totalWeight != TOTAL_WEIGHT) revert InvalidWeights();
    }
}