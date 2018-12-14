pragma solidity ^0.4.23;

contract EthLog {
  event LogEvent(bytes32 indexed jobId);
  event Fulfillment(bytes32 data);

  function logEvent() public {
    emit LogEvent("hello_chainlink");
  }

  function fulfill(bytes32 _externalId, bytes32 _data) public {
      emit Fulfillment(_data);
  }
}
