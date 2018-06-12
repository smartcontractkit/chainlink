pragma solidity ^0.4.23;

contract EthLog {
  event LogEvent(bytes32 indexed jobId);

  function logEvent() public {
    emit LogEvent("hello_chainlink");
  }
}
