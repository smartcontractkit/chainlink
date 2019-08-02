pragma solidity 0.5.0;

import "./Consumer.sol";

contract UpdatableConsumer is Consumer {

  constructor(bytes32 _specId, address _ens, bytes32 _node) public {
    specId = _specId;
    useChainlinkWithENS(_ens, _node);
  }

  function updateOracle() public {
    updateChainlinkOracleWithENS();
  }

  function getChainlinkToken() public view returns (address) {
    return chainlinkTokenAddress();
  }

  function getOracle() public view returns (address) {
    return chainlinkOracleAddress();
  }

}
