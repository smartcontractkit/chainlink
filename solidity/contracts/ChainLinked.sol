pragma solidity ^0.4.18;

import "./Oracle.sol";
import "./ChainLink.sol";


contract ChainLinked {
  using ChainLink for ChainLink.Run;

  Oracle internal oracle;

  function NewRun(
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

  modifier onlyOracle() {
    require(msg.sender == address(oracle));
    _;
  }
}
