pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract UpdatableConsumer is Chainlinked {

  constructor(address _ens, bytes32 _ensNode) public {
    newChainlinkWithENS(_ens, _ensNode);
  }

  function publicOracle() public returns (address) {
    return oracle;
  }

  function publicLinkToken() public returns (address) {
    return link;
  }

  function updateOracle() public {
    updateOracleWithENS();
  }

}
