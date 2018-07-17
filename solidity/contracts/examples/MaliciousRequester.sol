pragma solidity ^0.4.24;


import "./MaliciousChainlinked.sol";


contract MaliciousRequester is MaliciousChainlinked {

  constructor(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function maliciousWithdraw()
    public
  {
    MaliciousChainlinkLib.Run memory run = newRun("specId", this, "doesNothing(bytes32,bytes32)");
    chainlinkRequest(run, LINK(1));
  }

  function doesNothing(bytes32 _requestId, bytes32 _data) public {}
}
