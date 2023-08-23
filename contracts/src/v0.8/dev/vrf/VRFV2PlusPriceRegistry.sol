// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/AggregatorV3Interface.sol";
import "../interfaces/IVRFV2PlusPriceRegistry.sol";
import "../../shared/access/ConfirmedOwner.sol";
import "../../ChainSpecificUtil.sol";

contract VRFV2PlusPriceRegistry is ConfirmedOwner, IVRFV2PlusPriceRegistry {
  /// @dev may not be provided upon construction on some chains due to lack of availability
  AggregatorV3Interface public s_linkEthFeed;
  /// @dev may not be provided upon construction on some chains due to lack of availability
  AggregatorV3Interface public s_linkUSDFeed;
  /// @dev may not be provided upon construction on some chains due to lack of availability
  AggregatorV3Interface public s_ethUSDFeed;

  event LinkEthFeedSet(address oldFeed, address newFeed);
  event LinkUSDFeedSet(address oldFeed, address newFeed);
  event EthUSDFeedSet(address oldFeed, address newFeed);

  error InvalidLinkWeiPrice(int256 linkWei);
  error InvalidEthUSDPrice(int256 ethUSD);
  error InvalidLinkUSDPrice(int256 linkUSD);
  error InvalidUSDPrice(address feed, int256 price);
  error PaymentTooLarge();

  struct Config {
    // stalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackWeiPerUnitLink.
    uint32 stalenessSeconds;

    // Gas to cover oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    // The recommended number is below, though it may vary slightly
    // if certain chains do not implement certain EIP's.
    // 21000 + // base cost of the transaction
    // 100 + 5000 + // warm subscription balance read and update. See https://eips.ethereum.org/EIPS/eip-2929
    // 2*2100 + 5000 - // cold read oracle address and oracle balance and first time oracle balance update, note first time will be 20k, but 5k subsequently
    // 4800 + // request delete refund (refunds happen after execution), note pre-london fork was 15k. See https://eips.ethereum.org/EIPS/eip-3529
    // 6685 + // Positive static costs of argument encoding etc. note that it varies by +/- x*12 for every x bytes of non-zero data in the proof.
    // Total: 37,185 gas.
    uint32 gasAfterPaymentCalculation;

    // Flat fee charged per fulfillment in 1e-8 of USD
    // i.e 1 USD == 100000000 "USD-8"
    // in other words, this is USD with 8 decimals rather than 2.
    // So for uint40, the maximum USD we can charge per premium is
    // max(uint40) * 1e-8 USD == 1099511627775 * 1e-8 ~= 10,995 USD
    // which should be more than enough.
    uint40 fulfillmentFlatFeeLinkUSD;

    // Flat fee charged per fulfillment in 1e-8 of USD.
    // i.e 1 USD == 100000000 "USD-8"
    // in other words, this is USD with 8 decimals rather than 2.
    // So for uint40, the maximum USD we can charge per premium is
    // max(uint40) * 1e-8 USD == 1099511627775 * 1e-8 ~= 10,995 USD
    // which should be more than enough.
    uint40 fulfillmentFlatFeeEthUSD;
  }
  Config public s_config;
  event ConfigSet(
    uint32 stalenessSeconds,
    int256 fallbackWeiPerUnitLink,
    int256 fallbackUSDPerUnitEth,
    int256 fallbackUSDPerUnitLink,
    uint40 fulfillmentFlatFeeLinkUSD,
    uint40 fulfillmentFlatFeeEthUSD
  );

  /// @notice fallback link/eth price used when respective feed is stale
  int256 public s_fallbackWeiPerUnitLink;
  /// @notice fallback eth/usd price used when respective feed is stale
  int256 public s_fallbackUSDPerUnitEth;
  /// @notice fallback link/usd price when respective feed is stale
  int256 public s_fallbackUSDPerUnitLink;

  /// @dev this is the number of decimals used in the fee config numbers described
  /// @dev in the FeeConfig struct.
  uint8 public constant USD_FEE_DECIMALS = 8;

  constructor(address linkEthFeed, address linkUSDFeed, address ethUSDFeed) ConfirmedOwner(msg.sender) {
    /// @dev no zero address checks since the provided addresses can legitimately
    /// @dev be zero if there are no feeds on a particular chain
    s_linkEthFeed = AggregatorV3Interface(linkEthFeed);
    s_linkUSDFeed = AggregatorV3Interface(linkUSDFeed);
    s_ethUSDFeed = AggregatorV3Interface(ethUSDFeed);
  }

  /**
   * @notice Set the link-eth feed to be used by the price registry
   * @param linkEthFeed the address of the link-eth feed
  */
  function setLINKETHFeed(address linkEthFeed) external onlyOwner {
    address oldFeed = address(s_linkEthFeed);
    s_linkEthFeed = AggregatorV3Interface(linkEthFeed);
    emit LinkEthFeedSet(oldFeed, linkEthFeed);
  }

  /**
    * @notice Set the link-usd feed to be used by the price registry
    * @param linkUsdFeed the address of the link-usd feed
  */
  function setLINKUSDFeed(address linkUsdFeed) external onlyOwner {
    address oldFeed = address(s_linkUSDFeed);
    s_linkUSDFeed = AggregatorV3Interface(linkUsdFeed);
    emit LinkUSDFeedSet(oldFeed, linkUsdFeed);
  }

  /**
   * @notice Set the eth-usd feed to be used by the price registry
   * @param ethUsdFeed the address of the eth-usd feed
  */
  function setEthUsdFeed(address ethUsdFeed) external onlyOwner {
    address oldFeed = address(s_ethUSDFeed);
    s_ethUSDFeed = AggregatorV3Interface(ethUsdFeed);
    emit EthUSDFeedSet(oldFeed, ethUsdFeed);
  }

  /**
   * @notice Set the config to be used by the price registry
   * @param stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @param gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @param fallbackWeiPerUnitLink fallback link/eth price in the case of a stale feed
   * @param fallbackUSDPerUnitEth fallback eth/usd price in the case of a stale feed
   * @param fallbackUSDPerUnitLink fallback link/usd price in the case of a stale feed
   * @param fulfillmentFlatFeeLinkUSD fulfillment flat fee for LINK fulfillments in USD, denominated in 8 decimals
   * @param fulfillmentFlatFeeEthUSD fulfillment flat fee for ETH (native) fulfillments in USD, denominated in 8 decimals
  */
  function setConfig(
    uint32 stalenessSeconds,
    uint32 gasAfterPaymentCalculation,
    int256 fallbackWeiPerUnitLink,
    int256 fallbackUSDPerUnitEth,
    int256 fallbackUSDPerUnitLink,
    uint40 fulfillmentFlatFeeLinkUSD,
    uint40 fulfillmentFlatFeeEthUSD
  ) external onlyOwner {
    if (fallbackWeiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(fallbackWeiPerUnitLink);
    }
    if (fallbackUSDPerUnitEth <= 0) {
      revert InvalidEthUSDPrice(fallbackUSDPerUnitEth);
    }
    if (fallbackUSDPerUnitLink <= 0) {
      revert InvalidLinkUSDPrice(fallbackUSDPerUnitLink);
    }
    s_fallbackWeiPerUnitLink = fallbackWeiPerUnitLink;
    s_fallbackUSDPerUnitEth = fallbackUSDPerUnitEth;
    s_fallbackUSDPerUnitLink = fallbackUSDPerUnitLink;
    s_config = Config({
      stalenessSeconds: stalenessSeconds,
      gasAfterPaymentCalculation: gasAfterPaymentCalculation,
      fulfillmentFlatFeeLinkUSD: fulfillmentFlatFeeLinkUSD,
      fulfillmentFlatFeeEthUSD: fulfillmentFlatFeeEthUSD
    });
    emit ConfigSet(
      stalenessSeconds,
      fallbackWeiPerUnitLink,
      fallbackUSDPerUnitEth,
      fallbackUSDPerUnitLink,
      fulfillmentFlatFeeLinkUSD,
      fulfillmentFlatFeeEthUSD
    );
  }

  /**
   * @inheritdoc IVRFV2PlusPriceRegistry
  */
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool nativePayment
  ) external view override returns (uint96) {
    if (nativePayment) {
      return
        calculatePaymentAmountEth(
          startGas,
          s_config.gasAfterPaymentCalculation,
          s_config.fulfillmentFlatFeeEthUSD,
          weiPerUnitGas
        );
    }
    return
      calculatePaymentAmountLink(
        startGas,
        s_config.gasAfterPaymentCalculation,
        s_config.fulfillmentFlatFeeLinkUSD,
        weiPerUnitGas
      );
  }

  function calculatePaymentAmountEth(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint40 fulfillmentFlatFeeEthUSD,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = ChainSpecificUtil.getCurrentTxL1GasFees();
    // calculate the payment without the premium
    uint256 baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft());
    // calculate the flat fee in wei, converting from the USD flat fee (denominated in 8 decimals)
    uint256 flatFeeWei = calculateFlatFeeFromUSD(fulfillmentFlatFeeEthUSD, s_ethUSDFeed);
    // return the final fee with the flat fee and l1 cost (if applicable) added
    return uint96(baseFeeWei + flatFeeWei + l1CostWei);
  }

  /**
   * @notice Calculate the flat fee in wei or in juels from the USD fee
   * @notice depending on whether the feed provided is the ETH_USD_FEED or the LINK_ETH_FEED.
   * @notice this is done because there would be unnecessary code duplication and bloat otherwise.
   * @param fulfillmentFlatFeeUSD the flat fee in USD, this is either s_config.fulfillmentFlatFeeEthUSD or s_config.fulfillmentFlatFeeLinkUSD
   * @param feed the feed to use to calculate the fee, this is either ETH_USD_FEED or LINK_ETH_FEED
  */
  function calculateFlatFeeFromUSD(uint40 fulfillmentFlatFeeUSD, AggregatorV3Interface feed) internal view returns (uint256 fee) {
    // if the fee is zero return zero.
    // this is likely the situation where we don't have a feed, therefore the code below would
    // revert due to zero addresses on the feed fields.
    if (fulfillmentFlatFeeUSD == 0) {
      return 0;
    }

    // Note that both LINK and the native token of EVM chains have 18 decimals.
    // Therefore, we can use the same logic for both.
    int256 usdPerUnitCrypto;
    uint8 decimals;
    (usdPerUnitCrypto, decimals) = getUSDFeedData(feed);
    if (usdPerUnitCrypto <= 0) {
      revert InvalidUSDPrice(address(feed), usdPerUnitCrypto);
    }
    if (decimals < USD_FEE_DECIMALS) {
      // because our representation has more decimals, we need to divide by the
      // difference to match the number of decimals in the aggregator contract.
      uint8 decimalDiff = USD_FEE_DECIMALS - decimals;
      // USD / (USD / {ETH|LINK}) = USD * ({ETH|LINK} / USD) = USD * ((1e18 {wei|juels}/{ETH|LINK}) / USD) = {wei|juels}
      // divide additionally by the decimal difference since the premium denomination is
      // in more decimals than the aggregator contract.
      fee = (uint256(fulfillmentFlatFeeUSD) * 1 ether) / (uint256(usdPerUnitCrypto) * uint256(10**decimalDiff));
    } else if (decimals > USD_FEE_DECIMALS) {
      // because our representation has less decimals, we need to multiply by
      // the difference to match the number of decimals in the aggregator contract.
      uint8 decimalDiff = decimals - USD_FEE_DECIMALS;
      // USD / (USD / {ETH|LINK}) = USD * ({ETH|LINK} / USD) = USD * ((1e18 {wei|juels}/{ETH|LINK}) / USD) = {wei|juels}
      // multiply additionally by the decimal difference since the premium denomination is
      // in less decimals than the aggregator contract.
      fee = (uint256(fulfillmentFlatFeeUSD) * 1 ether * uint256(10**decimalDiff)) / uint256(usdPerUnitCrypto);
    } else {
      // our representation is the same as the one in the aggregator contract,
      // so we can just do the conversion right away.
      // USD / (USD / {ETH|LINK}) = USD * ({ETH|LINK} / USD) = USD * ((1e18 {wei|juels}/{ETH|LINK}) / USD) = {wei|juels}
      fee = (uint256(fulfillmentFlatFeeUSD) * 1 ether) / uint256(usdPerUnitCrypto);
    }
  }

  function calculatePaymentAmountLink(
    uint256 startGas,
    uint256 gasAfterPaymentCalculation,
    uint40 fulfillmentFlatFeeLinkUSD,
    uint256 weiPerUnitGas
  ) internal view returns (uint96) {
    int256 weiPerUnitLink;
    weiPerUnitLink = getLINKEthFeedData();
    if (weiPerUnitLink <= 0) {
      revert InvalidLinkWeiPrice(weiPerUnitLink);
    }
    // Will return non-zero on chains that have this enabled
    uint256 l1CostWei = ChainSpecificUtil.getCurrentTxL1GasFees();
    // (1e18 juels/link) ((wei/gas * gas) + l1wei) / (wei/link) = juels
    uint256 paymentNoFee = (1e18 * (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft()) + l1CostWei)) /
      uint256(weiPerUnitLink);
    uint256 fee = calculateFlatFeeFromUSD(fulfillmentFlatFeeLinkUSD, s_linkEthFeed);
    if (paymentNoFee > (1e27 - fee)) {
      revert PaymentTooLarge(); // Payment + fee cannot be more than all of the link in existence.
    }
    return uint96(paymentNoFee + fee);
  }

  function getLINKEthFeedData() private view returns (int256) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 weiPerUnitLink;
    (, weiPerUnitLink, , timestamp, ) = s_linkEthFeed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      weiPerUnitLink = s_fallbackWeiPerUnitLink;
    }
    return weiPerUnitLink;
  }

  function getUSDFeedData(AggregatorV3Interface feed) private view returns (int256 answer, uint8 decimals) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    (, answer, , timestamp, ) = feed.latestRoundData();
    // solhint-disable-next-line not-rely-on-time
    if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
      if (feed == s_ethUSDFeed) {
        answer = s_fallbackUSDPerUnitEth;
      } else if (feed == s_linkEthFeed) {
        answer = s_fallbackWeiPerUnitLink;
      } else {
        revert(); // Should never happen
      }
    }
    decimals = feed.decimals();
  }
}
