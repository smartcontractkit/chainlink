pragma solidity 0.4.24;

import "./ChainlinkLib.sol";
import "./ENSResolver.sol";
import "./interfaces/ENSInterface.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/OracleInterface.sol";
import "./interfaces/CoordinatorInterface.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract Chainlinked {
  using ChainlinkLib for ChainlinkLib.Run;
  using SafeMath for uint256;

  uint256 constant private clArgsVersion = 1;
  uint256 constant private linkDivisibility = 10**18;

  LinkTokenInterface private link;
  OracleInterface private oracle;
  uint256 private requests = 1;
  mapping(bytes32 => address) private unfulfilledRequests;

  ENSInterface private ens;
  bytes32 private ensNode;
  bytes32 constant private ensTokenSubname = keccak256("link");
  bytes32 constant private ensOracleSubname = keccak256("oracle");

  event ChainlinkRequested(bytes32 id);
  event ChainlinkFulfilled(bytes32 id);
  event ChainlinkCancelled(bytes32 id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    bytes4 _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _amount)
    internal
    returns (bytes32)
  {
    return chainlinkRequestFrom(oracle, _run, _amount);
  }

  function chainlinkRequestFrom(address _oracle, ChainlinkLib.Run memory _run, uint256 _amount)
    internal
    returns (bytes32)
  {
    _run.requestId = bytes32(requests);
    requests += 1;
    _run.close();
    unfulfilledRequests[_run.requestId] = _oracle;
    emit ChainlinkRequested(_run.requestId);
    require(link.transferAndCall(_oracle, _amount, encodeForOracle(_run)), "unable to transferAndCall to oracle");

    return _run.requestId;
  }

  function cancelChainlinkRequest(bytes32 _requestId)
    internal
  {
    OracleInterface requested = OracleInterface(unfulfilledRequests[_requestId]);
    delete unfulfilledRequests[_requestId];
    emit ChainlinkCancelled(_requestId);
    requested.cancel(_requestId);
  }

  function LINK(uint256 _amount) internal pure returns (uint256) {
    return _amount.mul(linkDivisibility);
  }

  function setOracle(address _oracle) internal {
    oracle = OracleInterface(_oracle);
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

  function newChainlinkWithENS(address _ens, bytes32 _node)
    internal
    returns (address, address)
  {
    ens = ENSInterface(_ens);
    ensNode = _node;
    ENSResolver resolver = ENSResolver(ens.resolver(ensNode));
    bytes32 linkSubnode = keccak256(abi.encodePacked(ensNode, ensTokenSubname));
    setLinkToken(resolver.addr(linkSubnode));
    return (link, updateOracleWithENS());
  }

  function updateOracleWithENS()
    internal
    returns (address)
  {
    ENSResolver resolver = ENSResolver(ens.resolver(ensNode));
    bytes32 oracleSubnode = keccak256(abi.encodePacked(ensNode, ensOracleSubname));
    setOracle(resolver.addr(oracleSubnode));
    return oracle;
  }

  function encodeForOracle(ChainlinkLib.Run memory _run)
    internal
    view
    returns (bytes memory)
  {
    return abi.encodeWithSelector(
      oracle.requestData.selector,
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      clArgsVersion,
      _run.specId,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.requestId,
      _run.buf.buf);
  }

  function encodeForCoordinator(ChainlinkLib.Run memory _run)
    internal
    view
    returns (bytes memory)
  {
    return abi.encodeWithSelector(
      CoordinatorInterface(oracle).executeServiceAgreement.selector,
      0, // overridden by onTokenTransfer
      0, // overridden by onTokenTransfer
      clArgsVersion,
      _run.specId,
      _run.callbackAddress,
      _run.callbackFunctionId,
      _run.requestId,
      _run.buf.buf);
  }

  function serviceRequest(ChainlinkLib.Run memory _run, uint256 _amount)
    internal
    returns (bytes32)
  {
    _run.requestId = bytes32(requests);
    requests += 1;
    _run.close();
    unfulfilledRequests[_run.requestId] = oracle;
    emit ChainlinkRequested(_run.requestId);
    require(link.transferAndCall(oracle, _amount, encodeForCoordinator(_run)), "unable to transferAndCall to oracle");

    return _run.requestId;
  }

  modifier checkChainlinkFulfillment(bytes32 _requestId) {
    require(msg.sender == unfulfilledRequests[_requestId], "source must be the oracle of the request");
    delete unfulfilledRequests[_requestId];
    emit ChainlinkFulfilled(_requestId);
    _;
  }
}
