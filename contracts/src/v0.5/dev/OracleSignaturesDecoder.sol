pragma solidity 0.5.0;

contract OracleSignaturesDecoder {

  struct OracleSignatures {
    uint8[] vs;
    bytes32[] rs;
    bytes32[] ss;
  }

  function decodeOracleSignatures(
    bytes memory _oracleSignaturesData
  )
    internal
    pure
    returns(OracleSignatures memory)
  {
    // solhint-disable indent
    OracleSignatures memory signatures;
    ( signatures.vs, signatures.rs, signatures.ss) =
      abi.decode(_oracleSignaturesData, ( uint8[], bytes32[], bytes32[] ));
    return signatures;
  }
}
