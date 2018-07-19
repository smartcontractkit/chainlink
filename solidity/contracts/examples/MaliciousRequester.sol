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
    MaliciousChainlinkLib.WithdrawRun memory run = newWithdrawRun("specId", this, "doesNothing(bytes32,bytes32)");
    run.amount = link.balanceOf(address(oracle));
    chainlinkWithdrawRequest(run, LINK(1));
  }

  function doesNothing(bytes32 _requestId, bytes32 _data) public {}
}
