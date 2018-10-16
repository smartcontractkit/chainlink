pragma solidity ^0.4.24;

import "./Consumer.sol";

contract UpdatableConsumer is Consumer {

  address public publicLinkToken;
  address public publicOracle;

  constructor(bytes32 _specId, address _ens) public {
    specId = _specId;
    address link;
    address oracle;
    (link, oracle) = newChainlinkWithENS(_ens, _ensNode);
    publicLinkToken = link;
    publicOracle = oracle;
  }

  function updateOracle() public {
    publicOracle = updateOracleWithENS();
  }

}
