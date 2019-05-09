pragma solidity 0.4.24;

import "../../../evm/contracts/ChainlinkClient.sol";
import "../../../node_modules/openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title An externally-connected Chainlink example contract
 * @notice This contract can be used on any public network to create Chainlink requests
 * @dev Use our docs at https://docs.chain.link
 */
contract ATestnetConsumer is ChainlinkClient, Ownable {
  // Requests on the test networks are priced at 1 request = 1 LINK
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  uint256 public currentPrice;
  int256 public changeDay;
  bytes32 public lastMarket;

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

  /**
   * @notice Deploys the contract and automatically detects the LINK token address
   * @dev Calling setChainlinkToken with address(0) will call setPublicChainlinkToken
   * in the ChainlinkClient contract to automatically set the correct LINK token address
   */
  constructor() Ownable() public {
    setPublicChainlinkToken();
  }

  /**
   * @notice Creates a Chainlink request to the specified oracle address for the specified JobID and currency
   * @dev The fulfillment function expects a uint256 type
   * @param _oracle The oracle address to send the request to
   * @param _jobId The JobID on the oracle to execute
   * @param _currency The currency that the price should be converted to (supports USD, EUR, & JPY)
   */
  function requestEthereumPrice(address _oracle, string _jobId, string _currency)
    public
    onlyOwner
  {
    Chainlink.Request memory req = buildChainlinkRequest(stringToBytes32(_jobId), this, this.fulfillEthereumPrice.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY");
    req.add("path", _currency);
    req.addInt("times", 100);
    sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
  }

  /**
   * @notice Creates a Chainlink request to the specified oracle address for the specified JobID and currency
   * @dev The fulfillment function expects an int256 type
   * @param _oracle The oracle address to send the request to
   * @param _jobId The JobID on the oracle to execute
   * @param _currency The currency that the price should be converted to (supports USD, EUR, & JPY)
   */
  function requestEthereumChange(address _oracle, string _jobId, string _currency)
    public
    onlyOwner
  {
    Chainlink.Request memory req = buildChainlinkRequest(stringToBytes32(_jobId), this, this.fulfillEthereumChange.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "CHANGEPCTDAY";
    req.addStringArray("path", path);
    req.addInt("times", 1000000000);
    sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
  }

  /**
   * @notice Creates a Chainlink request to the specified oracle address for the specified JobID and currency
   * @dev The fulfillment function expects a bytes32 type
   * @param _oracle The oracle address to send the request to
   * @param _jobId The JobID on the oracle to execute
   * @param _currency The currency that the price should be converted to (supports USD, EUR, & JPY)
   */
  function requestEthereumLastMarket(address _oracle, string _jobId, string _currency)
    public
    onlyOwner
  {
    Chainlink.Request memory req = buildChainlinkRequest(stringToBytes32(_jobId), this, this.fulfillEthereumLastMarket.selector);
    req.add("get", "https://min-api.cryptocompare.com/data/pricemultifull?fsyms=ETH&tsyms=USD,EUR,JPY");
    string[] memory path = new string[](4);
    path[0] = "RAW";
    path[1] = "ETH";
    path[2] = _currency;
    path[3] = "LASTMARKET";
    req.addStringArray("path", path);
    sendChainlinkRequestTo(_oracle, req, ORACLE_PAYMENT);
  }

  /**
   * @notice The fulfill method from requests created by requestEthereumPrice
   * @dev The recordChainlinkFulfillment protects this function from being called
   * by anyone other than the oracle address that the request was sent to
   * @param _requestId The ID that was generated for the request
   * @param _price The answer provided by the oracle
   */
  function fulfillEthereumPrice(bytes32 _requestId, uint256 _price)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumPriceFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  /**
   * @notice The fulfill method from requests created by requestEthereumChange
   * @dev The recordChainlinkFulfillment protects this function from being called
   * by anyone other than the oracle address that the request was sent to
   * @param _requestId The ID that was generated for the request
   * @param _change The answer provided by the oracle
   */
  function fulfillEthereumChange(bytes32 _requestId, int256 _change)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumChangeFulfilled(_requestId, _change);
    changeDay = _change;
  }

  /**
   * @notice The fulfill method from requests created by requestEthereumLastMarket
   * @dev The recordChainlinkFulfillment protects this function from being called
   * by anyone other than the oracle address that the request was sent to
   * @param _requestId The ID that was generated for the request
   * @param _market The answer provided by the oracle
   */
  function fulfillEthereumLastMarket(bytes32 _requestId, bytes32 _market)
    public
    recordChainlinkFulfillment(_requestId)
  {
    emit RequestEthereumLastMarket(_requestId, _market);
    lastMarket = _market;
  }

  /**
   * @notice Returns the address of the LINK token
   * @dev This is the public implementation for chainlinkTokenAddress, which is
   * an internal method of the ChainlinkClient contract
   */
  function getChainlinkToken() public view returns (address) {
    return chainlinkTokenAddress();
  }

  /**
   * @notice Allows the owner to withdraw any LINK balance on the contract
   */
  function withdrawLink() public onlyOwner {
    LinkTokenInterface link = LinkTokenInterface(getChainlinkToken());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
  }

  /**
   * @notice Helper method to convert strings to bytes32
   */
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
