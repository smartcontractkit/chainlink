// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {Common} from "./libraries/Common.sol";
import {IERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {SafeERC20} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {FeeManager} from "./FeeManager.sol";

/**
 * @title FeeManagerNoNative
 * @author Michael Fletcher
 * @author Austin Born
 * @author TJ Cunha
 * @notice This contract is a variation of the FeeManager contract, and adds a modifier blocks native billing to fee handling functions
 */
contract FeeManagerNoNative is FeeManager {
  using SafeERC20 for IERC20;

  /// @notice thrown when trying to pay nativeFee with native, which is disallowed when we force wETH billing for nativeFee
  error NativeBillingDisallowed();

  modifier blockNativeBilling() {
    if (msg.value != 0) revert NativeBillingDisallowed();
    _;
  }
  constructor(
    address _linkAddress,
    address _nativeAddress,
    address _proxyAddress,
    address _rewardManagerAddress
  ) FeeManager(_linkAddress, _nativeAddress, _proxyAddress, _rewardManagerAddress) {}

  /// @inheritdoc FeeManager
  function processFee(
    bytes calldata payload,
    bytes calldata parameterPayload,
    address subscriber
  ) external payable override onlyProxy blockNativeBilling {
    (Common.Asset memory fee, Common.Asset memory reward, uint256 appliedDiscount) = _processFee(
      payload,
      parameterPayload,
      subscriber
    );

    IFeeManager.FeeAndReward[] memory feeAndReward = new IFeeManager.FeeAndReward[](1);
    feeAndReward[0] = IFeeManager.FeeAndReward(bytes32(payload), fee, reward, appliedDiscount);

    if (fee.assetAddress == i_linkAddress) {
      _handleFeesAndRewards(subscriber, feeAndReward, 1, 0);
    } else {
      _handleFeesAndRewards(subscriber, feeAndReward, 0, 1);
    }
  }

  /// @inheritdoc FeeManager
  function processFeeBulk(
    bytes[] calldata payloads,
    bytes calldata parameterPayload,
    address subscriber
  ) external payable override onlyProxy blockNativeBilling {
    FeeAndReward[] memory feesAndRewards = new IFeeManager.FeeAndReward[](payloads.length);

    //keep track of the number of fees to prevent over initialising the FeePayment array within _convertToLinkAndNativeFees
    uint256 numberOfLinkFees;
    uint256 numberOfNativeFees;

    uint256 feesAndRewardsIndex;
    for (uint256 i; i < payloads.length; ++i) {
      (Common.Asset memory fee, Common.Asset memory reward, uint256 appliedDiscount) = _processFee(
        payloads[i],
        parameterPayload,
        subscriber
      );

      if (fee.amount != 0) {
        feesAndRewards[feesAndRewardsIndex++] = IFeeManager.FeeAndReward(
          bytes32(payloads[i]),
          fee,
          reward,
          appliedDiscount
        );

        unchecked {
          //keep track of some tallys to make downstream calculations more efficient
          if (fee.assetAddress == i_linkAddress) {
            ++numberOfLinkFees;
          } else {
            ++numberOfNativeFees;
          }
        }
      }
    }

    if (numberOfLinkFees != 0 || numberOfNativeFees != 0) {
      _handleFeesAndRewards(subscriber, feesAndRewards, numberOfLinkFees, numberOfNativeFees);
    }
  }
}
