pragma solidity 0.4.24;


import "./MaliciousChainlinked.sol";


contract MaliciousRequester is MaliciousChainlinked {

  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations
  uint256 private expiration;

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
    chainlinkWithdrawRequest(run, ORACLE_PAYMENT);
  }

  function request(bytes32 _id, address _target, bytes _callbackFunc) public returns (bytes32 requestId) {
    ChainlinkLib.Run memory run = newRun(_id, _target, bytes4(keccak256(_callbackFunc)));
    expiration = now.add(5 minutes);
    requestId = chainlinkRequest(run, ORACLE_PAYMENT);
  }

  function maliciousPrice(bytes32 _id) public returns (bytes32 requestId) {
    ChainlinkLib.Run memory run = newRun(_id, this, this.doesNothing.selector);
    requestId = chainlinkPriceRequest(run, ORACLE_PAYMENT);
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    ChainlinkLib.Run memory run = newRun("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    requestId = chainlinkTargetRequest(_target, run, ORACLE_PAYMENT);
  }

  function maliciousRequestCancel(bytes32 _id, bytes _callbackFunc) public {
    ChainlinkRequestInterface oracle = ChainlinkRequestInterface(oracleAddress());
    oracle.cancel(
      request(_id, this, _callbackFunc),
      ORACLE_PAYMENT,
      this.maliciousRequestCancel.selector,
      expiration
    );
  }

  function doesNothing(bytes32, bytes32) public pure {}
}
