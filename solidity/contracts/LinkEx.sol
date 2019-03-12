pragma solidity 0.4.24;

import "./interfaces/LinkExInterface.sol";

/**
 * @title The LINK exchange contract
 */
contract LinkEx is LinkExInterface {

  uint256 private historicRate;
  uint256 private rate;
  uint256 private rateHeight;

  function currentRate() external view returns (uint256) {
    if (isFutureBlock()) {
      return rate;
    }
    return historicRate;
  }

  function update(uint256 _rate) external {
    if (isFutureBlock()) {
      historicRate = rate;
      rateHeight = block.number;
    }
    rate = _rate;
  }

  function isFutureBlock() internal view returns (bool) {
    return block.number > rateHeight;
  }
}
