pragma solidity ^0.4.23;


import "../Chainlinked.sol";


contract ConcreteChainlinked is Chainlinked {

  constructor(address _link, address _oracle)
    public
  {
    setLinkToken(_link);
    setOracle(_oracle);
  }

  event Run(
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackfunctionSelector,
    bytes data
  );

  function publicNewRun(
    bytes32 _specId,
    address _address,
    string _fulfillmentSignature
  )
    public
  {
    ChainlinkLib.Run memory run = newRun(_specId, _address, _fulfillmentSignature);
    run.close();
    emit Run(
      run.specId,
      run.callbackAddress,
      run.callbackFunctionId,
      run.buf.buf
    );
  }

}
