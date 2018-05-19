pragma solidity ^0.4.23;

import "./ChainlinkLib.sol";
import "./Oracle.sol";
import "linkToken/contracts/LinkToken.sol";

contract Chainlinked {
  using ChainlinkLib for ChainlinkLib.Run;
  using ChainlinkLib for ChainlinkLib.Spec;

  uint256 constant clArgsVersion = 1;

  LinkToken internal link;
  Oracle internal oracle;
  uint256 internal requests = 1;

  event ChainlinkRequest(bytes32 id);

  function newRun(
    bytes32 _specId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    return run.initialize(_specId, _callbackAddress, _callbackFunctionSignature);
  }

  function newSpec(
    string[] _tasks,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Spec memory) {
    ChainlinkLib.Spec memory spec;
    return spec.initialize(_tasks, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.close();
    require(link.transferAndCall(oracle, _wei, _run.encodeForOracle(clArgsVersion)));
    emit ChainlinkRequest(_run.requestId);
    return _run.requestId;
  }

  function chainlinkRequest(ChainlinkLib.Spec memory _spec, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _spec.requestId = bytes32(requests);
    _spec.close();
    require(link.transferAndCall(oracle, _wei, _spec.encodeForOracle(clArgsVersion)));
    emit ChainlinkRequest(_spec.requestId);
    return _spec.requestId;
  }

  function LINK(uint256 _amount) internal pure returns (uint256) {
    return _amount * 10**18;
  }

  function setOracle(address _oracle) internal {
    oracle = Oracle(_oracle);
  }

  function setLinkToken(address _link) internal {
    link = LinkToken(_link);
  }

  modifier onlyOracle() {
    require(msg.sender == address(oracle));
    _;
  }
}
