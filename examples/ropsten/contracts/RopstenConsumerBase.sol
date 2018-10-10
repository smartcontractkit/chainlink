pragma solidity ^0.4.24;

import "../../../solidity/contracts/Chainlinked.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract RopstenConsumer is Chainlinked, Ownable {
  uint256 public currentPrice;
  int256 public changeDay;
  bytes32 public lastMarket;

  address constant ROPSTEN_LINK_ADDRESS = 0x20fE562d797A42Dcb3399062AE9546cd06f63280;
  address constant ROPSTEN_ORACLE_ADDRESS = 0x261a3F70acdC85CfC2FFc8badE43b1D42bf75D69;

  bytes32 constant PRICE_SPEC_ID = bytes32("2216dd2bf5464687a05ded0b844e200c");
  bytes32 constant CHANGE_SPEC_ID = bytes32("3f97669ef7c54e1990ce97149d7ada19");
  bytes32 constant MARKET_SPEC_ID = bytes32("3116f19a4b00434ba28b4a5cec349474");

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
    setLinkToken(ROPSTEN_LINK_ADDRESS);
    setOracle(ROPSTEN_ORACLE_ADDRESS);
  }

  function requestEthereumPrice(string _currency) 
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(PRICE_SPEC_ID, this, "fulfillEthereumPrice(bytes32,uint256)");
    run.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    run.addStringArray("path", path);
    run.addInt("times", 100);
    chainlinkRequest(run, LINK(1));
  }

  function requestEthereumChange(string _currency)
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(CHANGE_SPEC_ID, this, "fulfillEthereumChange(bytes32,int256)");
    run.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "CHANGEPCTDAY";
    run.addStringArray("path", path);
    run.addInt("times", 1000000000);
    chainlinkRequest(run, LINK(1));
  }

  function requestEthereumLastMarket(string _currency)
    public
    onlyOwner
  {
    ChainlinkLib.Run memory run = newRun(MARKET_SPEC_ID, this, "fulfillEthereumLastMarket(bytes32,bytes32)");
    run.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "LASTMARKET";
    run.addStringArray("path", path);
    chainlinkRequest(run, LINK(1));
  }

  function fulfillEthereumPrice(bytes32 _requestId, uint256 _price)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumPriceFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  function fulfillEthereumChange(bytes32 _requestId, int256 _change)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumChangeFulfilled(_requestId, _change);
    changeDay = _change;
  }

  function fulfillEthereumLastMarket(bytes32 _requestId, bytes32 _market)
    public
    checkChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumLastMarket(_requestId, _market);
    lastMarket = _market;
  }

  function withdrawLink() public onlyOwner {
    ILinkToken link = ILinkToken(chainlinkToken());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

}
