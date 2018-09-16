pragma solidity ^0.4.24;

//import "../Chainlinked.sol";
import "./Consumer.sol";

contract UpdatableConsumer is Consumer {

  constructor(bytes32 _specId, address _ens, bytes32 _ensNode) public {
    specId = _specId;
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
