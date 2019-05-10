pragma solidity 0.4.24;

import "../ChainlinkClient.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ConversionRate is ChainlinkClient, Ownable {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  uint256 public currentRate;
  bytes32[] public jobIds;
  address[] public oracles;

  constructor(address _link, address[] _oracles, bytes32[] _jobIds) public {
    setChainlinkToken(_link);
    jobIds = _jobIds;
    oracles = _oracles;
  }

  function update() public {
    Chainlink.Request memory request;
    for (uint i = 0; i < oracles.length; i++) {
      request = buildChainlinkRequest(jobIds[i], this, this.chainlinkCallback.selector);
      sendChainlinkRequestTo(oracles[i], request, ORACLE_PAYMENT);
    }
  }

  function chainlinkCallback(bytes32 _requestId, uint256 _rate)
    public
  {
    validateChainlinkCallback(_requestId);
    currentRate = _rate;
  }

}
