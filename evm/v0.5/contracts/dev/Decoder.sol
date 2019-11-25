pragma solidity 0.5.0;

import "./CoordinatorInterface.sol";

library Decoder {
  function decodeServiceAgreement(
    bytes memory _serviceAgreementData
  )
    internal
    pure
    returns(CoordinatorInterface.ServiceAgreement memory)
  {
    // solhint-disable indent
    CoordinatorInterface.ServiceAgreement memory agreement;

    ( agreement.payment,
      agreement.expiration,
      agreement.endAt,
      agreement.oracles,
      agreement.requestDigest,
      agreement.aggregator,
      agreement.aggInitiateJobSelector,
      agreement.aggFulfillSelector) =
      abi.decode(
        _serviceAgreementData,
        ( uint256,
        uint256,
        uint256,
        address[],
        bytes32,
        address,
        bytes4,
        bytes4 )
      );

    return agreement;
  }

  function decodeOracleSignatures(
    bytes memory _oracleSignaturesData
  )
    internal
    pure
    returns(CoordinatorInterface.OracleSignatures memory)
  {
    // solhint-disable indent
    CoordinatorInterface.OracleSignatures memory signatures;
    ( signatures.vs, signatures.rs, signatures.ss) =
      abi.decode(_oracleSignaturesData, ( uint8[], bytes32[], bytes32[] ));
    return signatures;
  }
}
