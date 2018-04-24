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
  bytes4 constant oracleFid = bytes4(keccak256("requestData(uint256,bytes32,address,bytes4,bytes32,bytes)"));

  LinkToken internal link;
  Oracle internal oracle;
  uint256 internal requests = 1;

  function newRun(
    bytes32 _jobId,
    address _callbackAddress,
    string _callbackFunctionSignature
  ) internal returns (ChainlinkLib.Run memory) {
    ChainlinkLib.Run memory run;
    Buffer.init(run.buf, 128);
    run.jobId = _jobId;
    run.callbackAddress = _callbackAddress;
    run.callbackFunctionId = bytes4(keccak256(_callbackFunctionSignature));
    run.buf.startMap();

    return run;
  }

  function chainlinkRequest(ChainlinkLib.Run memory _run, uint256 _wei)
    internal
    returns(bytes32)
  {
    bytes32 requestId = keccak256(this, requests++);
    bytes memory requestDataABI = abi.encodeWithSelector(
      oracleFid,
      clArgsVersion,
      _run.jobId,
      _run.callbackAddress,
      _run.callbackFunctionId,
      requestId,
      _run.close());
    link.transferAndCall(oracle, _wei, requestDataABI);

    return requestId;
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
