pragma solidity 0.4.24;

import "../Chainlinked.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ConversionRate is Chainlinked, Ownable {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  address public getOracle;
  uint256 public currentRate;
  bytes32 public getJobId;

  constructor(address _link, address _oracle, bytes32 _jobId) public {
    setLinkToken(_link);
    updateOracle(_oracle);
    updateJobId(_jobId);
    getJobId = _jobId;
  }

  function update() public {
    chainlinkRequest(newRequest(getJobId, this, this.updateCallback.selector), ORACLE_PAYMENT);
  }

  function updateCallback(bytes32 _requestId, uint256 _rate)
    public
    recordChainlinkFulfillment(_requestId)
  {
    currentRate = _rate;
  }

  function updateOracle(address _oracle) public onlyOwner {
    getOracle = _oracle;
    setOracle(_oracle);
  }

  function updateJobId(bytes32 _jobId) public onlyOwner {
    getJobId = _jobId;
  }

}
