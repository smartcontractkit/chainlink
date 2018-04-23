pragma solidity ^0.4.23;


import "../Chainlinked.sol";


contract ConcreteChainlinked is Chainlinked {

  function ConcreteChainlinked(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  event Run(
    bytes32 jobId,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );

  function publicNewRun(
    bytes32 _jobId,
    address _address,
    string _fulfillmentSignature
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(_jobId, _address, _fulfillmentSignature);
    emit Run(
      run.jobId,
      run.callbackAddress,
      run.callbackFunctionId,
      run.close()
    );
  }

}
