pragma solidity ^0.7.0;

import "../ChainlinkClient.sol";
import "../Chainlink.sol";

contract MultiWordConsumer is ChainlinkClient {
  using Chainlink for Chainlink.Request;

  bytes32 internal specId;
  bytes public currentPrice;

  bytes32 public usd;
  bytes32 public eur;
  bytes32 public jpy;

  uint256 public usdInt;
  uint256 public eurInt;
  uint256 public jpyInt;

  event RequestFulfilled(
    bytes32 indexed requestId, // User-defined ID
    bytes indexed price
  );

  event RequestMultipleFulfilled(bytes32 indexed requestId, bytes32 indexed usd, bytes32 indexed eur, bytes32 jpy);

  event RequestMultipleFulfilledWithCustomURLs(
    bytes32 indexed requestId,
    uint256 indexed usd,
    uint256 indexed eur,
    uint256 jpy
  );

  constructor(
    address _link,
    address _oracle,
    bytes32 _specId
  ) public {
    setChainlinkToken(_link);
    setChainlinkOracle(_oracle);
    specId = _specId;
  }

  function setSpecID(bytes32 _specId) public {
    specId = _specId;
  }

  function requestEthereumPrice(string memory _currency, uint256 _payment) public {
    Chainlink.Request memory req = buildOperatorRequest(specId, this.fulfillBytes.selector);
    sendOperatorRequest(req, _payment);
  }

  function requestMultipleParameters(string memory _currency, uint256 _payment) public {
    Chainlink.Request memory req = buildOperatorRequest(specId, this.fulfillMultipleParameters.selector);
    sendOperatorRequest(req, _payment);
  }

  function requestMultipleParametersWithCustomURLs(
    string memory _urlUSD,
    string memory _pathUSD,
    string memory _urlEUR,
    string memory _pathEUR,
    string memory _urlJPY,
    string memory _pathJPY,
    uint256 _payment
  ) public {
    Chainlink.Request memory req = buildOperatorRequest(specId, this.fulfillMultipleParametersWithCustomURLs.selector);
    req.add("urlUSD", _urlUSD);
    req.add("pathUSD", _pathUSD);
    req.add("urlEUR", _urlEUR);
    req.add("pathEUR", _pathEUR);
    req.add("urlJPY", _urlJPY);
    req.add("pathJPY", _pathJPY);
    sendOperatorRequest(req, _payment);
  }

  function cancelRequest(
    address _oracle,
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  ) public {
    ChainlinkRequestInterface requested = ChainlinkRequestInterface(_oracle);
    requested.cancelOracleRequest(_requestId, _payment, _callbackFunctionId, _expiration);
  }

  function withdrawLink() public {
    LinkTokenInterface _link = LinkTokenInterface(chainlinkTokenAddress());
    require(_link.transfer(msg.sender, _link.balanceOf(address(this))), "Unable to transfer");
  }

  function addExternalRequest(address _oracle, bytes32 _requestId) external {
    addChainlinkExternalRequest(_oracle, _requestId);
  }

  function fulfillMultipleParameters(
    bytes32 _requestId,
    bytes32 _usd,
    bytes32 _eur,
    bytes32 _jpy
  ) public recordChainlinkFulfillment(_requestId) {
    emit RequestMultipleFulfilled(_requestId, _usd, _eur, _jpy);
    usd = _usd;
    eur = _eur;
    jpy = _jpy;
  }

  function fulfillMultipleParametersWithCustomURLs(
    bytes32 _requestId,
    uint256 _usd,
    uint256 _eur,
    uint256 _jpy
  ) public recordChainlinkFulfillment(_requestId) {
    emit RequestMultipleFulfilledWithCustomURLs(_requestId, _usd, _eur, _jpy);
    usdInt = _usd;
    eurInt = _eur;
    jpyInt = _jpy;
  }

  function fulfillBytes(bytes32 _requestId, bytes memory _price) public recordChainlinkFulfillment(_requestId) {
    emit RequestFulfilled(_requestId, _price);
    currentPrice = _price;
  }

  function publicGetNextRequestCount() external view returns (uint256) {
    return getNextRequestCount();
  }
}
