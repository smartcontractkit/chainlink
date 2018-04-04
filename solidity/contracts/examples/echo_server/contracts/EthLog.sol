pragma solidity ^0.4.18;

contract EthLog {
  event LogEvent(bytes32 indexed jobId);

  function logEvent() public {
    LogEvent("hello_chainlink");
  }
}
