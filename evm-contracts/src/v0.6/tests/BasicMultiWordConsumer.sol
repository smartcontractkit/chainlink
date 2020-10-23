pragma solidity ^0.6.0;

import "./MultiWordConsumer.sol";

contract BasicMultiWordConsumer is MultiWordConsumer {

  constructor(address _link, address _oracle, bytes32 _specId) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_oracle);
    specId = _specId;
  }

}