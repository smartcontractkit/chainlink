pragma solidity 0.4.24;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./interfaces/LinkTokenInterface.sol";
import "./interfaces/CoordinatorInterface.sol";

// Coordinator handles oracle service aggreements between one or more oracles.
contract Coordinator is CoordinatorInterface {
  using SafeMath for uint256;

  LinkTokenInterface internal LINK;

  struct ServiceAgreement {
    uint256 payment;
    uint256 expiration;
    uint256 endAt;
    address[] oracles;
    bytes32 requestDigest;
  }

  struct Callback {
    bytes32 externalId;
    uint256 amount;
    address addr;
    bytes4 functionId;
    uint64 cancelExpiration;
  }

  mapping(uint256 => Callback) private callbacks;
  mapping(bytes32 => ServiceAgreement) public serviceAgreements;

  constructor(address _link) public {
    LINK = LinkTokenInterface(_link);
  }

  event RunRequest(
    bytes32 indexed sAId,
    address indexed requester,
    uint256 indexed amount,
    uint256 internalId,
    uint256 version,
    bytes data
  );

  event NewServiceAgreement(
    bytes32 indexed said,
    bytes32 indexed requestDigest
  );

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
    require(address(this).delegatecall(_data), "Unable to create request"); // calls executeServiceAgreement
  }

  function executeServiceAgreement(
    address _sender,
    uint256 _amount,
    uint256 _version,
    bytes32 _sAId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    bytes32 _externalId,
    bytes _data
  )
    external
    onlyLINK
    sufficientLINK(_amount, _sAId)
  {
    uint256 internalId = uint256(keccak256(abi.encodePacked(_sender, _externalId)));
    callbacks[internalId] = Callback(
      _externalId,
      _amount,
      _callbackAddress,
      _callbackFunctionId,
      uint64(now.add(5 minutes)));
    emit RunRequest(
      _sAId,
      _sender,
      _amount,
      internalId,
      _version,
      _data);
  }

  // This is mostly useful as a sanity check in the #getId test, because the
  // hash value there is illegible by design.
  function getPackedArguments(
    uint256 _payment,
    uint256 _expiration,
    uint256 _endAt,
    address[] _oracles,
    bytes32 _requestDigest
  )
    public
    pure
    returns (bytes)
  {
    return abi.encodePacked(_payment, _expiration, _endAt, _oracles, _requestDigest);
  }

  function getId(
    uint256 _payment,
    uint256 _expiration,
    uint256 _endAt,
    address[] _oracles,
    bytes32 _requestDigest
  )
    public pure returns (bytes32)
  {
    return keccak256(getPackedArguments(_payment, _expiration, _endAt, _oracles, _requestDigest));
  }

  function initiateServiceAgreement(
    uint256 _payment,
    uint256 _expiration,
    uint256 _endAt,
    address[] _oracles,
    uint8[] _vs,
    bytes32[] _rs,
    bytes32[] _ss,
    bytes32 _requestDigest
  )
    public
    returns (bytes32 serviceAgreementID)
  {
    require(_oracles.length == _vs.length && _vs.length == _rs.length && _rs.length == _ss.length, "Must pass in as many signatures as oracles"); /* solium-disable-line max-len */
    require(_endAt > block.timestamp, "End of ServiceAgreement must be in the future");

    serviceAgreementID = getId(_payment, _expiration, _endAt, _oracles, _requestDigest);

    for (uint i = 0; i < _oracles.length; i++) {
      address signer = getOracleAddressFromSASignature(serviceAgreementID, _vs[i], _rs[i], _ss[i]);
      require(_oracles[i] == signer, "Invalid oracle signature specified in SA");
    }

    serviceAgreements[serviceAgreementID] = ServiceAgreement(
      _payment,
      _expiration,
      _endAt,
      _oracles,
      _requestDigest
    );

    emit NewServiceAgreement(serviceAgreementID, _requestDigest);
  }

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

  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token"); 
    _;
  }

  bytes4 constant private permittedFunc =
    bytes4(keccak256("executeServiceAgreement(address,uint256,uint256,bytes32,address,bytes4,bytes32,bytes)")); /* solium-disable-line indentation */

  modifier permittedFunctionsForLINK() {
    bytes4[1] memory funcSelector;
    assembly {
      // solium-disable-next-line security/no-low-level-calls
      calldatacopy(funcSelector, 132, 4) // grab function selector from calldata
    }
    require(funcSelector[0] == permittedFunc, "Must use whitelisted functions");
    _;
  }

  modifier sufficientLINK(uint256 _amount, bytes32 _sAId) {
    require(_amount >= serviceAgreements[_sAId].payment, "Below agreed payment");
    _;
  }

  modifier hasInternalId(uint256 _internalId) {
    require(callbacks[_internalId].addr != address(0), "Must have a valid internalId");
    _;
  }

  function fulfillData(
    uint256 _internalId,
    bytes32 _data
  )
    public
    hasInternalId(_internalId)
    returns (bool)
  {
    Callback memory callback = callbacks[_internalId];
    delete callbacks[_internalId];
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted. See:
    // https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    return callback.addr.call(callback.functionId, callback.externalId, _data); // solium-disable-line security/no-low-level-calls
  }
}
