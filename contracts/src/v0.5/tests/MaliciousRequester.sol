pragma solidity 0.5.0;


import "./MaliciousChainlinkClient.sol";


contract MaliciousRequester is MaliciousChainlinkClient {

  uint256 constant private ORACLE_PAYMENT = 1 * LINK;
  uint256 private expiration;

  constructor(address _link, address _oracle) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_oracle);
  }

  function maliciousWithdraw()
    public
  {
    MaliciousChainlink.WithdrawRequest memory req = newWithdrawRequest(
      "specId", address(this), this.doesNothing.selector);
    chainlinkWithdrawRequest(req, ORACLE_PAYMENT);
  }

  function request(bytes32 _id, address _target, bytes memory _callbackFunc) public returns (bytes32 requestId) {
    Chainlink.Request memory req = buildChainlinkRequest(_id, _target, bytes4(keccak256(_callbackFunc)));
    expiration = now.add(5 minutes); // solhint-disable-line not-rely-on-time
    requestId = sendChainlinkRequest(req, ORACLE_PAYMENT);
  }

  function maliciousPrice(bytes32 _id) public returns (bytes32 requestId) {
    Chainlink.Request memory req = buildChainlinkRequest(_id, address(this), this.doesNothing.selector);
    requestId = chainlinkPriceRequest(req, ORACLE_PAYMENT);
  }

  function maliciousTargetConsumer(address _target) public returns (bytes32 requestId) {
    Chainlink.Request memory req = buildChainlinkRequest("specId", _target, bytes4(keccak256("fulfill(bytes32,bytes32)")));
    requestId = chainlinkTargetRequest(_target, req, ORACLE_PAYMENT);
  }

  function maliciousRequestCancel(bytes32 _id, bytes memory _callbackFunc) public {
    ChainlinkRequestInterface _oracle = ChainlinkRequestInterface(chainlinkOracleAddress());
    _oracle.cancelOracleRequest(
      request(_id, address(this), _callbackFunc),
      ORACLE_PAYMENT,
      this.maliciousRequestCancel.selector,
      expiration
    );
  }

  function doesNothing(bytes32, bytes32) public pure {} // solhint-disable-line no-empty-blocks
}
