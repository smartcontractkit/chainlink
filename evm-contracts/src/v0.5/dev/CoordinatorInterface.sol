pragma solidity 0.5.0;

contract CoordinatorInterface {

  function initiateServiceAgreement(
    bytes memory _serviceAgreementData,
    bytes memory _oracleSignaturesData)
    public returns (bytes32);

  function fulfillOracleRequest(
    bytes32 _requestId,
    bytes32 _aggregatorArgs)
    external returns (bool);
}
