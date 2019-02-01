pragma solidity 0.4.24;

import "./Consumer.sol";

contract UpdatableConsumer is Consumer {

  constructor(bytes32 _specId, address _ens, bytes32 _node) public {
    specId = _specId;
    setChainlinkWithENS(_ens, _node);
  }

  function updateOracle() public {
    setOracleWithENS();
  }

  function getChainlinkToken() public view returns (address) {
    return chainlinkToken();
  }

  function getOracle() public view returns (address) {
    return oracleAddress();
  }

}
