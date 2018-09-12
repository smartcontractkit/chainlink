pragma solidity ^0.4.24;

import "./ChainlinkLib.sol";
import "./ENSResolver.sol";
import "./Oracle.sol";
import "ens/contracts/ENS.sol";
import "linkToken/contracts/LinkToken.sol";

contract Chainlinked {
  using ChainlinkLib for ChainlinkLib.Run;
  using SafeMath for uint256;

  uint256 constant clArgsVersion = 1;

  LinkToken internal link;
  Oracle internal oracle;
  uint256 internal requests = 1;
  mapping(bytes32 => bool) internal unfulfilledRequests;

  ENS internal ens;
  bytes32 internal ensNode;

  event ChainlinkRequested(bytes32 id);
  event ChainlinkFulfilled(bytes32 id);
  event ChainlinkCancelled(bytes32 id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.close();
    require(link.transferAndCall(oracle, _wei, _run.encodeForOracle(clArgsVersion)), "Unable to transferAndCall to oracle");
    emit ChainlinkRequested(_run.requestId);
    unfulfilledRequests[_run.requestId] = true;
    return _run.requestId;
  }

  function cancelChainlinkRequest(bytes32 _requestId)
    internal
  {
    oracle.cancel(_requestId);
    unfulfilledRequests[_requestId] = false;
    emit ChainlinkCancelled(_requestId);
  }

  function LINK(uint256 _amount) internal view returns (uint256) {
    return _amount.mul(10**18);
  }

  function setOracle(address _oracle) internal {
    oracle = Oracle(_oracle);
  }

  function setLinkToken(address _link) internal {
    link = LinkToken(_link);
  }

  function newChainlinkWithENS(address _ens, bytes32 _node) internal {
    ens = ENS(_ens);
    ensNode = _node;
    ENSResolver resolver = ENSResolver(ens.resolver(_node));
    setOracle(resolver.addr(_node));
  }

  modifier checkChainlinkFulfillment(bytes32 _requestId) {
    require(msg.sender == address(oracle) && unfulfilledRequests[_requestId], "Source must be oracle with a valid requestId");
    _;
    unfulfilledRequests[_requestId] = false;
    emit ChainlinkFulfilled(_requestId);
  }
}