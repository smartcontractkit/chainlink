pragma solidity ^0.4.23;

import "./ChainlinkLib.sol";
import "./LinkToken.sol";
import "./Oracle.sol";
import "./Buffer.sol";
import "./CBOR.sol";


contract Chainlinked {
  using ChainlinkLib for ChainlinkLib.Run;
  using CBOR for Buffer.buffer;

  uint256 constant clArgsVersion = 1;

  LinkToken internal link;
  Oracle internal oracle;
  uint256 internal requests = 1;

  function newRun(
    bytes32 _jobId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal pure returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    return run.initialize(_jobId, _callbackAddress, _callbackFunctionSignature);
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    requests += 1;
    _run.requestId = bytes32(requests);
    _run.close();
    require(link.transferAndCall(
      oracle,
      _wei,
      _run.encodeForOracle(clArgsVersion)));

    return _run.requestId;
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
