// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

// Based on code from https://github.com/Uniswap/uniswap-v2-periphery

// a library for handling binary fixed point numbers (https://en.wikipedia.org/wiki/Q_(number_format))
library FixedPoint {
  // range: [0, 2**112 - 1]
  // resolution: 1 / 2**112
  struct uq112x112 {
    uint224 _x;
  }

  // returns a uq112x112 which represents the ratio of the numerator to the denominator
  // equivalent to encode(numerator).div(denominator)
  function fraction(
    uint112 numerator,
    uint112 denominator
  )
    internal
    pure
    returns (
      uq112x112 memory
    )
  {
    require(denominator > 0, "FixedPoint: DIV_BY_ZERO");
    return uq112x112((uint224(numerator) << 112) / denominator);
  }

  // decode a uq112x112 into a uint with 18 decimals of precision
  function decode112with18(
    uq112x112 memory self
  )
    internal
    pure
    returns (
      uint
    )
  {
    // we only have 256 - 224 = 32 bits to spare, so scaling up by ~60 bits is dangerous
    // instead, get close to:
    //  (x * 1e18) >> 112
    // without risk of overflowing, e.g.:
    //  (x) / 2 ** (112 - lg(1e18))
    return uint(self._x) / 5192296858534827;
  }
}

// library with helper methods for oracles that are concerned with computing average prices
library UniswapV2OracleLibrary {
  using FixedPoint for *;

  // helper function that returns the current block timestamp within the range of uint32, i.e. [0, 2**32 - 1]
  function currentBlockTimestamp()
    internal
    view
    returns (
      uint32
    )
  {
    return uint32(block.timestamp % 2 ** 32);
  }

  // produces the cumulative price using counterfactuals to save gas and avoid a call to sync.
  function currentCumulativePrices(
    address pair
  )
    internal
    view
    returns (
      uint price0Cumulative,
      uint price1Cumulative,
      uint32 blockTimestamp
    )
  {
    blockTimestamp = currentBlockTimestamp();
    price0Cumulative = IUniswapV2Pair(pair).price0CumulativeLast();
    price1Cumulative = IUniswapV2Pair(pair).price1CumulativeLast();

    // if time has elapsed since the last update on the pair, mock the accumulated price values
    (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast) = IUniswapV2Pair(pair).getReserves();
    if (blockTimestampLast != blockTimestamp) {
      // subtraction overflow is desired
      uint32 timeElapsed = blockTimestamp - blockTimestampLast;
      // addition overflow is desired
      // counterfactual
      price0Cumulative += uint(FixedPoint.fraction(reserve1, reserve0)._x) * timeElapsed;
      // counterfactual
      price1Cumulative += uint(FixedPoint.fraction(reserve0, reserve1)._x) * timeElapsed;
    }
  }
}

interface IUniswapV2Pair {
  function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast);
  function price0CumulativeLast() external view returns (uint);
  function price1CumulativeLast() external view returns (uint);
}
