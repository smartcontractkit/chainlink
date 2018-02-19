pragma solidity ^0.4.18;

import "./Oracle.sol";
import "./Chainlink.sol";


contract Chainlinked {
  using Chainlink for Chainlink.Run;

  Oracle internal oracle;

  function newRun(
    bytes32 _jobId,
    address _cbReceiver,
    string _cbSignature
  ) internal returns (Chainlink.Run) {
    Chainlink.Run memory run;
    run.jobId = _jobId;
    run.receiver = _cbReceiver;
    run.functionSelector = bytes4(keccak256(_cbSignature));
    return run;
  }

  function chainlinkRequest(Chainlink.Run _run) internal returns(uint256) {
    return oracle.requestData(
      _run.jobId, _run.receiver, _run.functionSelector, _run.close());
  }

  function setOracle(address _oracle) internal {
    oracle = Oracle(_oracle);
  }

  modifier onlyOracle() {
    require(msg.sender == address(oracle));
    _;
  }
}
