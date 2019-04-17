pragma solidity 0.4.24;
pragma experimental ABIEncoderV2; // solium-disable-line no-experimental 

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./interfaces/ChainlinkRequestInterface.sol";
import "./interfaces/CoordinatorInterface.sol";
import "./interfaces/LinkTokenInterface.sol";

/**
 * @title The Chainlink Coordinator handles oracle service aggreements between one or more oracles
 */
contract Coordinator is ChainlinkRequestInterface, CoordinatorInterface {
  using SafeMath for uint256;

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
    bytes _data
  )
    external
    onlyLINK
    sufficientLINK(_amount, _sAId)
    checkCallbackAddress(_callbackAddress)
  {
    bytes32 requestId = keccak256(abi.encodePacked(_sender, _nonce));
    require(callbacks[requestId].cancelExpiration == 0, "Must use a unique ID");

    callbacks[requestId].sAId = _sAId;
    callbacks[requestId].amount = _amount;
    callbacks[requestId].addr = _callbackAddress;
    callbacks[requestId].functionId = _callbackFunctionId;
    callbacks[requestId].cancelExpiration = uint64(now.add(EXPIRY_TIME));

    emit OracleRequest(
      _sAId,
      _sender,
      requestId,
      _amount,
      _callbackAddress,
      _callbackFunctionId,
      now.add(EXPIRY_TIME),
      _dataVersion,
      _data);
  }

  /**
   * @notice Stores a Service Agreement which has been signed by the given oracles
   * @dev Validates that each oracle has a valid signature.
   * Emits NewServiceAgreement event.
   * @param _agreement The Service Agreement to be initiated
   * @param _signatures The signatures of the oracles in the agreement
   * @return The Service Agreement ID
   */
  function initiateServiceAgreement(
    ServiceAgreement memory _agreement,
    OracleSignatures memory _signatures
  )
    public 
    returns (bytes32 serviceAgreementID)
  {
    require(
      _agreement.oracles.length == _signatures.vs.length && 
      _signatures.vs.length == _signatures.rs.length && 
      _signatures.rs.length == _signatures.ss.length, 
      "Must pass in as many signatures as oracles"
    ); 
    require(_agreement.endAt > block.timestamp, "End of ServiceAgreement must be in the future");

    serviceAgreementID = getId(
      _agreement.payment, 
      _agreement.expiration, 
      _agreement.endAt, 
      _agreement.oracles, 
      _agreement.requestDigest
    );

    registerOracleSignatures(
      serviceAgreementID, 
      _agreement.oracles, 
      _signatures.vs, 
      _signatures.rs, 
      _signatures.ss
    );

    serviceAgreements[serviceAgreementID] = _agreement; 
    emit NewServiceAgreement(serviceAgreementID, _agreement.requestDigest);
  }

  /**
   * @dev Validates that each signer address matches for the given oracles
   * @param _serviceAgreementID Service agreement ID
   * @param _oracles Array of oracle addresses which agreed to the service agreement
   * @param _vs Array of recovery IDs of the oracle signatures
   * @param _rs Array of first 32 bytes of the oracle signatures
   * @param _ss Array of second 32 bytes of the oracle signatures
   */
  function registerOracleSignatures(
    bytes32 _serviceAgreementID,
    address[] _oracles,
    uint8[] _vs,
    bytes32[] _rs,
    bytes32[] _ss
  )
    private
  {
    for (uint i = 0; i < _oracles.length; i++) {
      address signer = getOracleAddressFromSASignature(_serviceAgreementID, _vs[i], _rs[i], _ss[i]);
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
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
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
  )
    external
    isValidRequest(_requestId)
    returns (bool)
  {
    storeResponse(_requestId, _data);

    Callback memory callback = callbacks[_requestId];
    address[] memory oracles = serviceAgreements[callback.sAId].oracles;
    if (oracles.length != callback.responseCount) {
      return true; // exit early if not all response have been received
    }

    uint256 result = aggregateAndPay(_requestId, callback.amount, oracles);
    return callback.addr.call(callback.functionId, _requestId, result); // solium-disable-line security/no-low-level-calls
  }

  /**
   * @dev Allows the oracle operator to withdraw their LINK
   * @param _recipient is the address the funds will be sent to
   * @param _amount is the amount of LINK transfered from the Coordinator contract
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
  {} // solium-disable-line no-empty-blocks

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
    bytes _data
  )
    public
    onlyLINK
    permittedFunctionsForLINK
  {
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      mstore(add(_data, 36), _sender) // ensure correct sender is passed
      // solium-disable-next-line security/no-low-level-calls
      mstore(add(_data, 68), _amount)    // ensure correct amount is passed
    }
    // solium-disable-next-line security/no-low-level-calls
    require(address(this).delegatecall(_data), "Unable to create request"); // calls oracleRequest
  }

  /**
   * @notice Retrieve the Service Agreement ID for the given parameters
   * @param _payment The amount of payment given (specified in wei)
   * @param _expiration The expiration that nodes should respond by
   * @param _endAt The date which the service agreement is no longer valid
   * @param _oracles Array of oracle addresses which agreed to the service agreement
   * @param _requestDigest Hash of the normalized job specification
   * @return The Service Agreement ID, a keccak256 hash of the input params
   */
  function getId(
    uint256 _payment,
    uint256 _expiration,
    uint256 _endAt,
    address[] _oracles,
    bytes32 _requestDigest
  )
    public pure returns (bytes32)
  {
    return keccak256(abi.encodePacked(_payment, _expiration, _endAt, _oracles, _requestDigest));
  }

  /**
   * @dev Stores an oracle's response and then increments the total number of repsonse received.
   * @param _requestId is the unique identifier generated by the requester and their nonce
   * @param _data is the actual answer to a request submitted by an oracle node
   */
  function storeResponse(bytes32 _requestId, bytes32 _data) private {
    callbacks[_requestId].responses[msg.sender] = uint256(_data);
    callbacks[_requestId].responseCount += 1;
  }

  /**
   * @dev Aggregates the responses in storage and cleans them up.
   * Then pays associated oracles their share of LINK.
   * Deletes the request's callback record, and finally returns the aggregated result.
   * @param _requestId is the unique identifier generated by the requester and their nonce
   * @param _paymentAmount is the amount of LINK shared between the oracles
   * @param _oracles is the list of oracle addresses
   */
  function aggregateAndPay(bytes32 _requestId, uint256 _paymentAmount, address[] memory _oracles) private returns (uint256) {
    uint256 sumQuotients;
    uint256 sumRemainders;
    uint256 oraclePayment = _paymentAmount.div(_oracles.length);
    for (uint i = 0; i < _oracles.length; i++) {
      uint256 response = callbacks[_requestId].responses[_oracles[i]];
      sumQuotients = sumQuotients.add(response.div(_oracles.length)); // aggregate responses and protect from overflows
      sumRemainders = sumRemainders.add(response % _oracles.length); 
      delete callbacks[_requestId].responses[_oracles[i]]; // must explicitly clean-up mappings for gas refund
      withdrawableTokens[_oracles[i]] = withdrawableTokens[_oracles[i]].add(oraclePayment);
    }
    delete callbacks[_requestId];
    return sumQuotients.add(sumRemainders.div(_oracles.length)); // recover lost accuracy from result 
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
    require(callbacks[_requestId].responses[msg.sender] == 0, "Cannot respond twice");
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
   * @dev Reverts if the given data does not begin with the `oracleRequest` function selector
   */
  modifier permittedFunctionsForLINK() {
    bytes4[1] memory funcSelector;
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      calldatacopy(funcSelector, 132, 4) // grab function selector from calldata
    }
    require(funcSelector[0] == this.oracleRequest.selector, "Must use whitelisted functions");
    _;
  }
}
