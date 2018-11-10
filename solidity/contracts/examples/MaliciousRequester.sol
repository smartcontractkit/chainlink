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
    MaliciousChainlinkLib.WithdrawRun memory run = newWithdrawRun(
      "specId", this, this.doesNothing.selector);
    run.amount = link.balanceOf(address(oracle));
    chainlinkWithdrawRequest(run, LINK(1));
  }

  function request()
    internal
    returns (bytes32 requestId)
  {
    MaliciousChainlinkLib.Run memory run = newRun(
      "specId", this, this.doesNothing.selector);
    requestId = chainlinkRequest(run, LINK(1));
  }

  function maliciousRequestCancel() public {
    oracle.cancel(request());
  }

  function doesNothing(bytes32, bytes32) public pure {}
}
