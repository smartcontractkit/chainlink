pragma solidity 0.4.24;

import "./interfaces/LinkExInterface.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title The LINK exchange contract
 */
contract LinkEx is LinkExInterface, Ownable {
  using SafeMath for uint256;

  struct Rate {
    uint256 blockNumber;
    uint256 rate;
  }

  uint256 public constant VALID_BLOCKS = 25;

  mapping(address => bool) public authorizedNodes;

  mapping(address => Rate) private rates;
  uint256 private historicRate;
  uint256 private rate;
  uint256 private rateHeight;
  address[] private oracles;

  function addOracle(address _oracle) external onlyOwner {
    require(!authorizedNodes[_oracle], "Oracle is already added");
    setFulfillmentPermission(_oracle, true);
    oracles.push(_oracle);
  }

  function currentRate() external view returns (uint256) {
    if (isFutureBlock()) {
      return rate;
    }
    return historicRate;
  }

  function removeOracle(address _oracle) external onlyOwner {
    require(_removeOracle(_oracle), "Oracle does not exist");
    if (oracles.length > 0) {
      _update();
    }
  }

  function _removeOracle(address _oracle) private returns (bool) {
    delete authorizedNodes[_oracle];
    uint256 size = oracles.length.sub(1);
    for (uint i = 0; i < oracles.length; i++) {
      if (oracles[i] == _oracle) {
        delete rates[oracles[i]];
        delete oracles[i];
        if (i < size) {
          oracles[i] = oracles[size];
        }
        oracles.length = size;
        return true;
      }
    }
    return false;
  }

  function update(uint256 _rate) external onlyAuthorizedNode {
    rates[msg.sender] = Rate(block.number, _rate);
    _update();
  }

  function _update() private {
    if (isFutureBlock()) {
      historicRate = rate;
      rateHeight = block.number;
    }

    uint256 rateSum;
    uint256 count;
    uint256 fromBlock = block.number.sub(VALID_BLOCKS);
    for (uint i = 0; i < oracles.length; i++) {
      if(rates[oracles[i]].blockNumber > fromBlock) {
        rateSum = rateSum.add(rates[oracles[i]].rate);
        count++;
      }
    }
    rate = rateSum.div(count);
  }

  function isFutureBlock() internal view returns (bool) {
    return block.number > rateHeight;
  }

  function setFulfillmentPermission(address _oracle, bool _status) private {
    authorizedNodes[_oracle] = _status;
  }

  modifier onlyAuthorizedNode() {
    require(authorizedNodes[msg.sender], "Only an authorized node may call this function");
    _;
  }
}
