pragma solidity 0.4.24;


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
    chainlinkWithdrawRequest(run, LINK(1));
  }

  function request()
    internal
    returns (bytes32 requestId)
  {
    MaliciousChainlinkLib.Run memory run = newRun("specId", this, this.doesNothing.selector);
    requestId = chainlinkRequest(run, LINK(1));
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    MaliciousChainlinkLib.Run memory run = newRun("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    requestId = chainlinkTargetRequest(_target, run, LINK(1));
  }

  function maliciousRequestCancel() public {
    OracleInterface oracle = OracleInterface(oracleAddress());
    oracle.cancel(request());
  }

  function doesNothing(bytes32, bytes32) public pure {}
}
