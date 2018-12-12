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

  function request(bytes32 _id)
    internal
    returns (bytes32 requestId)
  {
    ChainlinkLib.Run memory run = newRun(_id, this, this.doesNothing.selector);
    requestId = chainlinkRequest(run, LINK(1));
  }

  function maliciousPrice(bytes32 _id) public returns (bytes32 requestId) {
    ChainlinkLib.Run memory run = newRun(_id, this, this.doesNothing.selector);
    requestId = chainlinkPriceRequest(run, LINK(1));
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    ChainlinkLib.Run memory run = newRun("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    requestId = chainlinkTargetRequest(_target, run, LINK(1));
  }

  function maliciousRequestCancel(bytes32 _id) public {
    ChainlinkRequestInterface oracle = ChainlinkRequestInterface(oracleAddress());
    oracle.cancel(request(_id));
  }

  function doesNothing(bytes32, bytes32) public pure {}
}
