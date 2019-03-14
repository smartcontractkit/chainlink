pragma solidity 0.4.24;

import "./interfaces/LinkExInterface.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title The LINK exchange contract
 */
contract LinkEx is LinkExInterface, Ownable {

  mapping(address => bool) public authorizedNodes;

  uint256 private historicRate;
  uint256 private rate;
  uint256 private rateHeight;
  address[] private oracles;

  function addOracle(address _oracle) external onlyOwner {
    setFulfillmentPermission(_oracle, true);
    oracles.push(_oracle);
  }

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

  function setFulfillmentPermission(address _oracle, bool _status) private {
    authorizedNodes[_oracle] = _status;
  }
}
