pragma solidity ^0.4.18;


import "../Chainlinked.sol";


contract ConcreteChainlinked is Chainlinked {

  function ConcreteChainlinked(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  function publicNewRun(
    bytes32 _jobId,
    address _address,
    string _functionSignature
  )
    public
  {
    chainlinkRequest(newRun(_jobId, _address, _functionSignature));
  }

}
