pragma solidity 0.4.24;


import "../Chainlinked.sol";
import "./MaliciousConsumer.sol";


contract MaliciousServiceAgreementConsumer is Chainlinked, MaliciousConsumer {

  constructor(address _link, address _oracle)
    public
    MaliciousConsumer(_link, _oracle)
  {}

  function requestData(string _callbackFunc) public {
    bytes4 callbackFID = bytes4(keccak256(bytes(_callbackFunc)));
    ChainlinkLib.Run memory run = newRun("specId", this, callbackFID);
    chainlinkRequest(run, LINK(1));
  }

}
