pragma solidity ^0.4.23;

import "./Chainlinked.sol";
import "github.com/OpenZeppelin/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ConsumerInt256 is Chainlinked, Ownable {
  bytes32 internal requestId;
  bytes32 internal jobId;
  int256 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    int256 indexed price
  );

  function ConsumerInt256(address _link, address _oracle, bytes32 _jobId) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function requestEthereumPrice(string _currency)
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(jobId, this, "fulfill(bytes32,int256)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    requestId = chainlinkRequest(run, LINK(1));
  }

  function cancelRequest(uint256 _requestId) 
    public 
    onlyOwner
  {
    oracle.cancel(_requestId);
  }

  function fulfill(bytes32 _requestId, int256 _price)
    public
    onlyOracle
    checkRequestId(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  modifier checkRequestId(bytes32 _requestId) {
    require(requestId == _requestId);
    _;
  }

}
