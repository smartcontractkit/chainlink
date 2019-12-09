pragma solidity 0.5.0;

library OracleSignatures {

  struct Instance {
    uint8[] vs;
    bytes32[] rs;
    bytes32[] ss;
  }

  function decode(
    bytes memory _oracleSignaturesData
  )
    internal
    pure
    returns(Instance memory)
  {
    // solhint-disable indent
    Instance memory signatures;
    ( signatures.vs, signatures.rs, signatures.ss) =
      abi.decode(_oracleSignaturesData, ( uint8[], bytes32[], bytes32[] ));
    return signatures;
  }
}
