// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRBillingInterface.sol";
import "../../ConfirmedOwner.sol";
import "../../interfaces/AggregatorV3Interface.sol";
import "../../../v0.4/interfaces/ERC20.sol";

/**
 * @title OCR2DR billing abstract.
 */
abstract contract OCR2DRBillingAbstract is ConfirmedOwner, OCR2DRBillingInterface {
  error UnsupportedFeeToken();
  error InsufficientGas();
  error InsufficientApproval();

  BillingConfig internal _billingConfig;

  constructor() {}

  function setFeeToken(
    address feeToken,
    address priceFeed,
    uint32 fee
  ) external onlyOwner {
    _billingConfig.feeTokenPriceOracles[feeToken] = priceFeed;
    _billingConfig.requiredFeeByToken[feeToken] = fee;
  }

  function setGasOverhead(uint32 gasOverhead) external onlyOwner {
    _billingConfig.gasOverhead = gasOverhead;
  }

  function getRequiredFee(bytes calldata data, RequestBillingConfig calldata billing)
    external
    view
    override
    returns (uint32)
  {
    // This is the minimum fee, goes directly to CLL/nop profits
    // We include the whole request for any custom billing we want per req.
    uint32 fee = _additionalRequiredFee(data);
    return _billingConfig.requiredFeeByToken[billing.feeToken] + fee;
  }

  function _additionalRequiredFee(bytes calldata data) internal view returns (uint32) {}

  function estimateExecutionFee(RequestBillingConfig calldata billing) external view override returns (uint256) {
    address feeToken = _billingConfig.feeTokenPriceOracles[billing.feeToken];
    if (feeToken == address(0)) {
      revert UnsupportedFeeToken();
    }
    uint256 timestamp;
    int256 weiPerUnitToken;
    (, weiPerUnitToken, , timestamp, ) = AggregatorV3Interface(feeToken).latestRoundData();

    return (_billingConfig.gasOverhead + billing.gasLimit) * uint256(weiPerUnitToken) * tx.gasprice;
  }

  function _preRequestBilling(bytes calldata data, RequestBillingConfig calldata billing) internal {
    requiredFee = getRequiredFee(data, billing);
    executionFee = estimateExecutionFee(billing);
    totalFee =  requiredFee + executionFee;
    if (billing.totalFee > totalFee) {
      revert InsufficientGas();
    }

    feeToken = ERC20(billing.feeToken);

    allowance = feeToken.allowance(address(this), s_oracle.address);
    if (billing.totalFee > allowance) {
      revert InsufficientApproval();
    }

    feeToken.transferFrom(msg.sender, billing.totalFee, address(this));
  }

  function _postFulfillBilling(address consumer, RequestBillingConfig calldata billing, address transmitter, address[] signers, uint32 initialGas, uint32 callbackGasCost) internal {
    /** 
    * Consumer Refunds *
    * Two options here:
    *   1. You can bill the user for the full gasLimit they requested, which allows you to tell them how much the refund is prior
    *      to the call, so they can take action upon it.
    *   2. You can measure the gas used and then send them the refund but they won't know how much they got refunded which can make
    *      accounting tricky.
    **/ 

    // Using option 1, no gas refund
    // refund = execution_fee - execution_cost
    // feeToken = ERC20(billing.feeToken);
    // feeToken.transferFrom(address(this), refund, consumer);

    /** 
    * Oracle Payment *
    * Two options here:
    *   1. Pay transmitter the full amount. Since the transmitter is chosen OCR, we trust the fairness of their selection algorithm.
    *   2. Reimburse the transmitter for execution cost, then split the requiredFee across all participants.
    **/ 

    // Using option 1, paying transmitter
    feeToken = ERC20(billing.feeToken);
    feeToken.transferFrom(address(this), billing.totalFee, transmitter);
  }
}
