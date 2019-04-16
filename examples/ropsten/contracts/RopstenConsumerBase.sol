pragma solidity 0.4.24;

import "../../../evm/contracts/Chainlinked.sol";
import "../../../node_modules/openzeppelin-evm/contracts/ownership/Ownable.sol";

contract ARopstenConsumer is Chainlinked, Ownable {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  uint256 public currentPrice;
  int256 public changeDay;
  bytes32 public lastMarket;

  address constant ROPSTEN_ENS = 0x112234455C3a32FD11230C42E7Bccd4A84e02010;
  bytes32 constant ROPSTEN_CHAINLINK_ENS = 0xead9c0180f6d685e43522fcfe277c2f0465fe930fb32b5b415826eacf9803727;

  event RequestEthereumPriceFulfilled(
    bytes32 indexed requestId,
    uint256 indexed price
  );

  event RequestEthereumChangeFulfilled(
    bytes32 indexed requestId,
    int256 indexed change
  );

  event RequestEthereumLastMarket(
    bytes32 indexed requestId,
    bytes32 indexed market
  );

  constructor() Ownable() public {
    setChainlinkWithENS(ROPSTEN_ENS, ROPSTEN_CHAINLINK_ENS);
  }

  function requestEthereumPrice(string _jobId, string _currency) 
    public
    onlyOwner
  {
    Chainlink.Request memory req = newRequest(stringToBytes32(_jobId), this, this.fulfillEthereumPrice.selector);
    req.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    req.addStringArray("path", path);
    req.addInt("times", 100);
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function requestEthereumChange(string _jobId, string _currency)
    public
    onlyOwner
  {
    Chainlink.Request memory req = newRequest(stringToBytes32(_jobId), this, this.fulfillEthereumChange.selector);
    req.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "CHANGEPCTDAY";
    req.addStringArray("path", path);
    req.addInt("times", 1000000000);
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function requestEthereumLastMarket(string _jobId, string _currency)
    public
    onlyOwner
  {
    Chainlink.Request memory req = newRequest(stringToBytes32(_jobId), this, this.fulfillEthereumLastMarket.selector);
    req.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "LASTMARKET";
    req.addStringArray("path", path);
    chainlinkRequest(req, ORACLE_PAYMENT);
  }

  function fulfillEthereumPrice(bytes32 _requestId, uint256 _price)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumPriceFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  function fulfillEthereumChange(bytes32 _requestId, int256 _change)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumChangeFulfilled(_requestId, _change);
    changeDay = _change;
  }

  function fulfillEthereumLastMarket(bytes32 _requestId, bytes32 _market)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumLastMarket(_requestId, _market);
    lastMarket = _market;
  }

  function updateChainlinkAddresses() public onlyOwner {
    setChainlinkWithENS(ROPSTEN_ENS, ROPSTEN_CHAINLINK_ENS);
  }

  function getChainlinkToken() public view returns (address) {
    return chainlinkToken();
  }
  
  function getOracle() public view returns (address) {
    return oracleAddress();
  }

  function withdrawLink() public onlyOwner {
    LinkTokenInterface link = LinkTokenInterface(chainlinkToken());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

  function stringToBytes32(string memory source) private pure returns (bytes32 result) {
    bytes memory tempEmptyStringTest = bytes(source);
    if (tempEmptyStringTest.length == 0) {
      return 0x0;
    }

    assembly {
      result := mload(add(source, 32))
    }
  }

}
