pragma solidity 0.4.24;

import "./ChainlinkLib.sol";
import "./ENSResolver.sol";
import "./interfaces/ENSInterface.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/ChainlinkRequestInterface.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract Chainlinked {
  using ChainlinkLib for ChainlinkLib.Run;
  using SafeMath for uint256;

  uint256 constant internal LINK = 10**18;
  uint256 constant private ARGS_VERSION = 1;
  bytes32 constant private ENS_TOKEN_SUBNAME = keccak256("link");
  bytes32 constant private ENS_ORACLE_SUBNAME = keccak256("oracle");

  ENSInterface private ens;
  bytes32 private ensNode;
  LinkTokenInterface private link;
  ChainlinkRequestInterface private oracle;
  uint256 private requests = 1;
  mapping(bytes32 => address) private pendingRequests;

  event ChainlinkRequested(bytes32 indexed id);
  event ChainlinkFulfilled(bytes32 indexed id);
  event ChainlinkCancelled(bytes32 indexed id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _payment)
    internal
    returns (bytes32)
  {
    return chainlinkRequestTo(oracle, _run, _payment);
  }

  function chainlinkRequestTo(address _oracle, ChainlinkLib.Run memory _run, uint256 _payment)
    internal
    returns (bytes32 requestId)
  {
    requestId = keccak256(abi.encodePacked(this, requests));
    _run.nonce = requests;
    pendingRequests[requestId] = _oracle;
    emit ChainlinkRequested(requestId);
    require(link.transferAndCall(_oracle, _payment, encodeRequest(_run)), "unable to transferAndCall to oracle");
    requests += 1;

    return requestId;
  }

  function cancelChainlinkRequest(
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunc,
    uint256 _expiration
  )
    internal
  {
    ChainlinkRequestInterface requested = ChainlinkRequestInterface(pendingRequests[_requestId]);
    delete pendingRequests[_requestId];
    emit ChainlinkCancelled(_requestId);
    requested.cancel(_requestId, _payment, _callbackFunc, _expiration);
  }

  function setOracle(address _oracle) internal {
    oracle = ChainlinkRequestInterface(_oracle);
  }

  function setLinkToken(address _link) internal {
    link = LinkTokenInterface(_link);
  }

  function chainlinkToken()
    internal
    view
    returns (address)
  {
    return address(link);
  }

  function oracleAddress()
    internal
    view
    returns (address)
  {
    return address(oracle);
  }

  function addExternalRequest(address _oracle, bytes32 _requestId)
    internal
    isPendingRequest(_requestId)
  {
    pendingRequests[_requestId] = _oracle;
  }

  function setChainlinkWithENS(address _ens, bytes32 _node)
    internal
  {
    ens = ENSInterface(_ens);
    ensNode = _node;
    ENSResolver resolver = ENSResolver(ens.resolver(ensNode));
    bytes32 linkSubnode = keccak256(abi.encodePacked(ensNode, ENS_TOKEN_SUBNAME));
    setLinkToken(resolver.addr(linkSubnode));
    setOracleWithENS();
  }

  function setOracleWithENS()
    internal
  {
    ENSResolver resolver = ENSResolver(ens.resolver(ensNode));
    bytes32 oracleSubnode = keccak256(abi.encodePacked(ensNode, ENS_ORACLE_SUBNAME));
    setOracle(resolver.addr(oracleSubnode));
  }

  function encodeRequest(ChainlinkLib.Run memory _run)
    internal
    view
    returns (bytes memory)
  {
    return abi.encodeWithSelector(
      oracle.requestData.selector,
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      ARGS_VERSION,
      _run.id,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.nonce,
      _run.buf.buf);
  }

  function completeChainlinkFulfillment(bytes32 _requestId)
    internal
    checkChainlinkFulfillment(_requestId)
  {}

  modifier checkChainlinkFulfillment(bytes32 _requestId) {
    require(msg.sender == pendingRequests[_requestId], "source must be the oracle of the request");
    delete pendingRequests[_requestId];
    emit ChainlinkFulfilled(_requestId);
    _;
  }

  modifier isPendingRequest(bytes32 _requestId) {
    require(pendingRequests[_requestId] == address(0), "Request is already fulfilled");
    _;
  }
}
