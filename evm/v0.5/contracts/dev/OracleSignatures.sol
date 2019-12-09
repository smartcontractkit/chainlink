pragma solidity 0.5.0;

library OracleSignatures {

  struct Signatures {
    uint8[] vs;
    bytes32[] rs;
    bytes32[] ss;
  }

  function decodeOracleSignatures(
    bytes memory _oracleSignaturesData
  )
    internal
    pure
    returns(Signatures memory)
  {
    // solhint-disable indent
    Signatures memory signatures;
    ( signatures.vs, signatures.rs, signatures.ss) =
      abi.decode(_oracleSignaturesData, ( uint8[], bytes32[], bytes32[] ));
    return signatures;
  }
}
