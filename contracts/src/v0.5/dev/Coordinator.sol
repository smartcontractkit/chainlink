pragma solidity 0.5.0;

import "./CoordinatorInterface.sol";
import "../interfaces/ChainlinkRequestInterface.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../vendor/SafeMathChainlink.sol";
import "./ServiceAgreementDecoder.sol";
import "./OracleSignaturesDecoder.sol";


/**
 * @title The Chainlink Coordinator handles oracle service agreements between one or more oracles
 */
contract Coordinator is ChainlinkRequestInterface, CoordinatorInterface, ServiceAgreementDecoder, OracleSignaturesDecoder {
  using SafeMathChainlink for uint256;

  uint256 constant public EXPIRY_TIME = 5 minutes;
  LinkTokenInterface internal LINK;

  struct Callback {
    bytes32 sAId;
    uint256 amount;
    address addr;
    bytes4 functionId;
    uint64 cancelExpiration;
    uint8 responseCount;
    mapping(address => uint256) responses;
  }

  mapping(bytes32 => Callback) private callbacks;
  mapping(bytes32 => mapping(address => bool)) private allowedOracles;
  mapping(bytes32 => ServiceAgreement) public serviceAgreements;
  mapping(address => uint256) public withdrawableTokens;

  /**
   * @notice Deploy with the address of the LINK token
   * @dev Sets the LinkToken address for the imported LinkTokenInterface
   * @param _link The address of the LINK token
   */
  constructor(address _link) public {
    LINK = LinkTokenInterface(_link);
  }

  event OracleRequest(
    bytes32 indexed sAId,
    address requester,
    bytes32 requestId,
    uint256 payment,
    address callbackAddr,
    bytes4 callbackFunctionId,
    uint256 cancelExpiration,
    uint256 dataVersion,
    bytes data
  );

  event NewServiceAgreement(
    bytes32 indexed said,
    bytes32 indexed requestDigest
  );

  event CancelOracleRequest(
    bytes32 internalId
  );

  /**
   * @notice Creates the Chainlink request
   * @dev Stores the params on-chain in a callback for the request.
   * Emits OracleRequest event for Chainlink nodes to detect.
   * @param _sender The sender of the request
   * @param _amount The amount of payment given (specified in wei)
   * @param _sAId The Service Agreement ID
   * @param _callbackAddress The callback address for the response
   * @param _callbackFunctionId The callback function ID for the response
   * @param _nonce The nonce sent by the requester
   * @param _dataVersion The specified data version
   * @param _data The CBOR payload of the request
   */
  function oracleRequest(
    address _sender,
    uint256 _amount,
    bytes32 _sAId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    uint256 _nonce,
    uint256 _dataVersion,
    bytes calldata _data
  )
    external
    onlyLINK
    sufficientLINK(_amount, _sAId)
    checkCallbackAddress(_callbackAddress)
    // checkServiceAgreementPresence(_sAId) // TODO: exhausts the stack
  {
    bytes32 requestId = keccak256(abi.encodePacked(_sender, _nonce));
    require(callbacks[requestId].cancelExpiration == 0, "Must use a unique ID");
    callbacks[requestId].sAId = _sAId;
    callbacks[requestId].amount = _amount;
    callbacks[requestId].addr = _callbackAddress;
    callbacks[requestId].functionId = _callbackFunctionId;
    // solhint-disable-next-line not-rely-on-time
    callbacks[requestId].cancelExpiration = uint64(now.add(EXPIRY_TIME));

    emit OracleRequest(
      _sAId,
      _sender,
      requestId,
      _amount,
      _callbackAddress,
      _callbackFunctionId,
      now.add(EXPIRY_TIME), // solhint-disable-line not-rely-on-time
      _dataVersion,
      _data);
  }

  /**
   * @notice Stores a Service Agreement which has been signed by the given oracles
   * @dev Validates that each oracle has a valid signature.
   * Emits NewServiceAgreement event.
   * @return The Service Agreement ID
   */
  function initiateServiceAgreement(
    bytes memory _serviceAgreementData,
    bytes memory _oracleSignaturesData
  )
    public
    returns (bytes32 serviceAgreementID)
  {

    ServiceAgreement memory _agreement = decodeServiceAgreement(_serviceAgreementData);
    OracleSignatures memory _signatures = decodeOracleSignatures(_oracleSignaturesData);

    require(
      _agreement.oracles.length == _signatures.vs.length &&
      _signatures.vs.length == _signatures.rs.length &&
      _signatures.rs.length == _signatures.ss.length,
      "Must pass in as many signatures as oracles"
    );
    // solhint-disable-next-line not-rely-on-time
    require(_agreement.endAt > block.timestamp,
      "ServiceAgreement must end in the future");
    require(serviceAgreements[serviceAgreementID].endAt == 0,
      "serviceAgreement already initiated");
    serviceAgreementID = getId(_agreement);

    registerOracleSignatures(
      serviceAgreementID,
      _agreement.oracles,
      _signatures
    );

    serviceAgreements[serviceAgreementID] = _agreement;
    emit NewServiceAgreement(serviceAgreementID, _agreement.requestDigest);
    // solhint-disable-next-line avoid-low-level-calls
    (bool ok, bytes memory response) = _agreement.aggregator.call(
      abi.encodeWithSelector(
        _agreement.aggInitiateJobSelector,
        serviceAgreementID,
        _serviceAgreementData
      )
    );
    require(ok, "Aggregator failed to initiate Service Agreement");
    require(response.length > 0, "probably wrong address/selector");
    (bool success, bytes memory message) = abi.decode(response, (bool, bytes));
    if ((!success) && message.length == 0) {
      // Revert with a non-empty message to give user a hint where to look
      require(success, "initiation failed; empty message");
    }
    require(success, string(message));
  }

  /**
   * @dev Validates that each signer address matches for the given oracles
   * @param _serviceAgreementID Service agreement ID
   * @param _oracles Array of oracle addresses which agreed to the service agreement
   * @param _signatures contains the collected parts(v, r, and s) of each oracle's signature.
   */
  function registerOracleSignatures(
    bytes32 _serviceAgreementID,
    address[] memory _oracles,
    OracleSignatures memory _signatures
  )
    private
  {
    for (uint i = 0; i < _oracles.length; i++) {
      address signer = getOracleAddressFromSASignature(
        _serviceAgreementID,
        _signatures.vs[i],
        _signatures.rs[i],
        _signatures.ss[i]
      );
      require(_oracles[i] == signer, "Invalid oracle signature specified in SA");
      allowedOracles[_serviceAgreementID][_oracles[i]] = true;
    }

  }

  /**
   * @dev Recovers the address of the signer for a service agreement
   * @param _serviceAgreementID Service agreement ID
   * @param _v Recovery ID of the oracle signature
   * @param _r First 32 bytes of the oracle signature
   * @param _s Second 32 bytes of the oracle signature
   * @return The address of the signer
   */
  function getOracleAddressFromSASignature(
    bytes32 _serviceAgreementID,
    uint8 _v,
    bytes32 _r,
    bytes32 _s
  )
    private pure returns (address)
  {
    bytes32 prefixedHash = keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", _serviceAgreementID));
    return ecrecover(prefixedHash, _v, _r, _s);
  }

  /**
   * @notice Called by the Chainlink node to fulfill requests
   * @dev Response must have a valid callback, and will delete the associated callback storage
   * before calling the external contract.
   * @param _requestId The fulfillment request ID that must match the requester's
   * @param _data The data to return to the consuming contract
   * @return Status if the external call was successful
   */
  function fulfillOracleRequest(
    bytes32 _requestId,
    bytes32 _data
  ) external isValidRequest(_requestId) returns (bool) {
    Callback memory callback = callbacks[_requestId];
    ServiceAgreement memory sA = serviceAgreements[callback.sAId];
    // solhint-disable-next-line avoid-low-level-calls
    (bool ok, bytes memory aggResponse) = sA.aggregator.call(
      abi.encodeWithSelector(
        sA.aggFulfillSelector, _requestId, callback.sAId, msg.sender, _data));
    require(ok, "aggregator.fulfill failed");
    require(aggResponse.length > 0, "probably wrong address/selector");
    (bool aggSuccess, bool aggComplete, bytes memory response, int256[] memory paymentAmounts) = abi.decode( // solhint-disable-line
      aggResponse, (bool, bool, bytes, int256[]));
    require(aggSuccess, string(response));
    if (aggComplete) {
      require(paymentAmounts.length == sA.oracles.length, "wrong paymentAmounts.length");
      for (uint256 oIdx = 0; oIdx < sA.oracles.length; oIdx++) { // pay oracles
        withdrawableTokens[sA.oracles[oIdx]] = uint256(int256(
          withdrawableTokens[sA.oracles[oIdx]]) + paymentAmounts[oIdx]);
      } // solhint-disable-next-line avoid-low-level-calls
      (bool success,) = callback.addr.call(abi.encodeWithSelector( // report final result
        callback.functionId, _requestId, abi.decode(response, (bytes32))));
      return success;
    }
    return true;
  }

  /**
   * @dev Allows the oracle operator to withdraw their LINK
   * @param _recipient is the address the funds will be sent to
   * @param _amount is the amount of LINK transferred from the Coordinator contract
   */
  function withdraw(address _recipient, uint256 _amount)
    external
    hasAvailableFunds(_amount)
  {
    withdrawableTokens[msg.sender] = withdrawableTokens[msg.sender].sub(_amount);
    assert(LINK.transfer(_recipient, _amount));
  }

  /**
   * @dev Necessary to implement ChainlinkRequestInterface
   */
  function cancelOracleRequest(bytes32, uint256, bytes4, uint256)
    external
  {} // solhint-disable-line no-empty-blocks

  /**
   * @notice Called when LINK is sent to the contract via `transferAndCall`
   * @dev The data payload's first 2 words will be overwritten by the `_sender` and `_amount`
   * values to ensure correctness. Calls oracleRequest.
   * @param _sender Address of the sender
   * @param _amount Amount of LINK sent (specified in wei)
   * @param _data Payload of the transaction
   */
  function onTokenTransfer(
    address _sender,
    uint256 _amount,
    bytes memory _data
  )
    public
    onlyLINK
    permittedFunctionsForLINK
  {
    assembly { // solhint-disable-line no-inline-assembly
      mstore(add(_data, 36), _sender) // ensure correct sender is passed
      mstore(add(_data, 68), _amount) // ensure correct amount is passed
    }
    // solhint-disable-next-line avoid-low-level-calls
    (bool success,) = address(this).delegatecall(_data); // calls oracleRequest or depositFunds
    require(success, "Unable to create request");
  }

  /**
   * @notice Retrieve the Service Agreement ID for the given parameters
   * @param _agreementData contains all of the terms of the service agreement that can be verified on-chain.
   * @return The Service Agreement ID, a keccak256 hash of the input params
   */
  function getId(bytes memory _agreementData) public pure returns (bytes32)
  {
    ServiceAgreement memory _agreement = decodeServiceAgreement(_agreementData);
    return getId(_agreement);
  }

  function getId(ServiceAgreement memory _agreement) internal pure returns (bytes32)
  {
    return keccak256(
      abi.encodePacked(
        _agreement.payment,
        _agreement.expiration,
        _agreement.endAt,
        _agreement.oracles,
        _agreement.requestDigest,
        _agreement.aggregator,
        _agreement.aggInitiateJobSelector,
        _agreement.aggFulfillSelector
    ));
  }

  /**
   * @notice Called when LINK is sent to the contract via `transferAndCall`
   * @param _sender Address of the sender
   * @param _amount Amount of LINK sent (specified in wei)
   */
  function depositFunds(address _sender, uint256 _amount) external onlyLINK
  {
    withdrawableTokens[_sender] = withdrawableTokens[_sender].add(_amount);
  }

  /**
   * @param _account Address to check balance of
   * @return Balance of account (specified in wei)
   */
  function balanceOf(address _account) public view returns (uint256)
  {
    return withdrawableTokens[_account];
  }

  /**
   * @dev Reverts if the callback address is the LINK token
   * @param _to The callback address
   */
  modifier checkCallbackAddress(address _to) {
    require(_to != address(LINK), "Cannot callback to LINK");
    _;
  }

  /**
   * @dev Reverts if amount requested is greater than withdrawable balance
   * @param _amount The given amount to compare to `withdrawableTokens`
   */
  modifier hasAvailableFunds(uint256 _amount) {
    require(withdrawableTokens[msg.sender] >= _amount, "Amount requested is greater than withdrawable balance");
    _;
  }

  /**
   * @dev Reverts if request ID does not exist
   * @param _requestId The given request ID to check in stored `callbacks`
   */
  modifier isValidRequest(bytes32 _requestId) {
    require(callbacks[_requestId].addr != address(0), "Must have a valid requestId");
    require(allowedOracles[callbacks[_requestId].sAId][msg.sender], "Oracle not recognized on service agreement");
    _;
  }

  /**
   * @dev Reverts if amount is not at least what was agreed upon in the service agreement
   * @param _amount The payment for the request
   * @param _sAId The service agreement ID which the request is for
   */
  modifier sufficientLINK(uint256 _amount, bytes32 _sAId) {
    require(_amount >= serviceAgreements[_sAId].payment, "Below agreed payment");
    _;
  }

  /**
   * @dev Reverts if the given data does not begin with the `oracleRequest` or
   * `depositFunds` function selector
   */
  modifier permittedFunctionsForLINK() {
    bytes4[1] memory funcSelector;
    assembly { // solhint-disable-line no-inline-assembly
      calldatacopy(funcSelector, 132, 4) // grab function selector from calldata
    }
    require(
      funcSelector[0] == this.oracleRequest.selector || funcSelector[0] == this.depositFunds.selector,
      "Must use whitelisted functions"
    );
    _;
  }

  modifier checkServiceAgreementPresence(bytes32 _sAId) {
    require(uint256(serviceAgreements[_sAId].requestDigest) != 0,
            "Must reference an existing ServiceAgreement");
    _;
  }

  /**
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }
}
