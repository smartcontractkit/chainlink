pragma solidity 0.4.24;


import "../Chainlinked.sol";
import "./MaliciousConsumer.sol";


contract MaliciousServiceAgreementConsumer is Chainlinked, MaliciousConsumer {

  constructor(address _link, address _oracle)
    public
    MaliciousConsumer(_link, _oracle)
  {}

  function requestData(string _callbackFunc) public {
    ChainlinkLib.Run memory run = newRun("specId", this, _callbackFunc);
    serviceRequest(run, LINK(1));
  }

}
