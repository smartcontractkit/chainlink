pragma solidity ^0.4.23;

import "github.com/smartcontractkit/chainlink/solidity/contracts/Chainlinked.sol";

contract RopstenConsumer is Chainlinked, Ownable {
  uint256 public currentPrice;
  int256 public changeDay;
  bytes32 public lastMarket;

  address constant ROPSTEN_LINK_ADDRESS = 0x20fE562d797A42Dcb3399062AE9546cd06f63280;
  address constant ROPSTEN_ORACLE_ADDRESS = 0x8D96353145F2e889890191e029ad63755Da173A5;

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
    string[] memory tasks = new string[](5);
    tasks[0] = "httpget";
    tasks[1] = "jsonparse";
    tasks[2] = "multiply";
    tasks[3] = "ethuint256";
    tasks[4] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfillEthereumPrice(bytes32,uint256)");
    spec.add("url", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](1);
    path[0] = _currency;
    spec.addStringArray("path", path);
    spec.addInt("times", 100);
    chainlinkRequest(spec, LINK(1));
  }

  function requestEthereumChange(string _currency)
    public
    onlyOwner
  {
    string[] memory tasks = new string[](5);
    tasks[0] = "httpget";
    tasks[1] = "jsonparse";
    tasks[2] = "multiply";
    tasks[3] = "ethint256";
    tasks[4] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfillEthereumChange(bytes32,int256)");
    spec.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "CHANGEPCTDAY";
    spec.addStringArray("path", path);
    spec.addInt("times", 1000000000);
    chainlinkRequest(spec, LINK(1));
  }

  function requestEthereumLastMarket(string _currency)
    public
    onlyOwner
  {
    string[] memory tasks = new string[](4);
    tasks[0] = "httpget";
    tasks[1] = "jsonparse";
    tasks[2] = "ethbytes32";
    tasks[3] = "ethtx";

    ChainlinkLib.Spec memory spec = newSpec(tasks, this, "fulfillEthereumLastMarket(bytes32,bytes32)");
    spec.add("url", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "LASTMARKET";
    spec.addStringArray("path", path);
    chainlinkRequest(spec, LINK(1));
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
    require(link.transfer(owner, link.balanceOf(address(this))));
  }

}
