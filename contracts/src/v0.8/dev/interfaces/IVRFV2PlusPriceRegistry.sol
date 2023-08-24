// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IVRFV2PlusPriceRegistry {
  /**
   * @notice Calculate the payment amount of a VRF request
   * @param startGas starting gas, retreived at the beginning of the fulfillment transaction using gasleft()
   * @param weiPerUnitGas the gas price at which to calculate the payment, typically tx.gasprice
   * @param nativePayment whether the payment is to be calculated in link or in native currency (i.e ether)
  */
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool nativePayment
  ) external view returns (uint96);

  /**
   * @notice Calculates the price of a VRF request with the given callbackGasLimit at the current
   * @notice block.
   *
   * @dev This function relies on the transaction gas price which is not automatically set during
   * @dev simulation. To estimate the price at a specific gas price, use the estimatePrice function.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @return The price of a VRF wrapper request denominated in juels with the given callbackGasLimit and parameters.
   */
  function calculateRequestPriceWrapper(
    uint32 _callbackGasLimit
  ) external view returns (uint256);

  /**
   * @notice Calculates the price of a VRF request with the given callbackGasLimit at the current
   * @notice block.
   *
   * @dev This function relies on the transaction gas price which is not automatically set during
   * @dev simulation. To estimate the price at a specific gas price, use the estimatePrice function.
   *
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @return The price of a VRF wrapper native request denominated in wei with the given callbackGasLimit and parameters.
   */
  function calculateRequestPriceNativeWrapper(
    uint32 _callbackGasLimit
  ) external view returns (uint256);

  /**
   * @notice Estimates the price of a VRF request with the given callbackGasLimit and request
   * @notice gas price.
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @param _requestGasPriceWei is the gas price used to estimate the price.
   * @return The price of a VRF wrapper request denominated in link with the given callbackGasLimit and parameters.
   */
  function estimateRequestPriceWrapper(
    uint32 _callbackGasLimit,
    uint256 _requestGasPriceWei
  ) external view returns (uint256);

  /**
   * @notice Estimates the price of a VRF request with the given callbackGasLimit and request
   * @notice gas price.
   * @param _callbackGasLimit is the gas limit used to estimate the price.
   * @param _requestGasPriceWei is the gas price used to estimate the price.
   * @return The price of a VRF wrapper native request denominated in wei with the given callbackGasLimit and parameters.
   */
  function estimateRequestPriceNativeWrapper(
    uint32 _callbackGasLimit,
    uint256 _requestGasPriceWei
  ) external view returns (uint256);
}
