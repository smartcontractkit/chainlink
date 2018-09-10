pragma solidity ^0.4.24;

import "../Chainlinked.sol";

contract UpdatableConsumer is Chainlinked {

  constructor(address _link, address _ens, bytes32 _ensNode) public {
    setLinkToken(_link);
    newChainlinkWithENS(_ens, _ensNode);
  }

  function publicOracle() public returns (address) {
    return oracle;
  }

  function publicLinkToken() public returns (address) {
    return link;
  }

}
