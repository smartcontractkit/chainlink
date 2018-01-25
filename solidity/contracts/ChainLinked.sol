pragma solidity ^0.4.18;

import "./Oracle.sol";
import "./ChainLink.sol";


contract ChainLinked {
  using ChainLink for ChainLink.Run;

  Oracle internal oracle;

  function newRun(
    bytes32 _jobId,
    address _cbReceiver,
    string _cbSignature
  ) internal returns (ChainLink.Run) {
    ChainLink.Run memory run;
    run.jobId = _jobId;
    run.receiver = _cbReceiver;
    run.functionHash = bytes4(keccak256(_cbSignature));
    return run;
  }

  function chainlinkRequest(ChainLink.Run _run) internal returns(uint256) {
    return oracle.requestData(_run.receiver, _run.functionHash, _run.close());
  }

  modifier onlyOracle() {
    require(msg.sender == address(oracle));
    _;
  }
}
