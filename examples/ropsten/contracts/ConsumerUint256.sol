pragma solidity ^0.4.23;

import "./Chainlinked.sol";
import "github.com/OpenZeppelin/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ConsumerUint256 is Chainlinked, Ownable {
  bytes32 internal jobId;
  uint256 public currentPrice;

  event RequestFulfilled(
    bytes32 indexed requestId,
    uint256 indexed price
  );

  constructor(address _link, address _oracle, bytes32 _jobId) Ownable() public {
    setLinkToken(_link);
    setOracle(_oracle);
    jobId = _jobId;
  }

  function requestEthereumPrice(string _currency)
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(jobId, this, "fulfill(bytes32,uint256)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    chainlinkRequest(run, LINK(1));
  }

  function cancelRequest(uint256 _requestId) 
  public 
  onlyOwner
  {
    oracle.cancel(_requestId);
  }

  function fulfill(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkRequest(_requestId)
  {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

}
