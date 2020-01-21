pragma solidity 0.5.0; // solhint-disable-line compiler-version

////////////////////////////////////////////////////////////////////////////////
////////// DO NOT USE THIS IN PRODUCTION UNTIL IT HAS BEEN AUDITED /////////////
////////////////////////////////////////////////////////////////////////////////

import "../interfaces/LinkTokenInterface.sol";
import "./VRF.sol";

/**
 * @title VRFCoordinator coordinates on-chain verifiable-randomness requests
 * @title with off-chain responses
 */
contract VRFCoordinator is VRF {

  LinkTokenInterface internal LINK;

  constructor(address _link) public {
    LINK = LinkTokenInterface(_link);
  }

  struct Callback { // Tracks an ongoing request
    address callbackContract; // Requesting contract, which will receive response
    uint256 randomnessFee; // Amount of LINK paid at request time
    uint256 seed; // Seed to use in generating this random value
  }

  struct ServiceAgreement { // Tracks oracle commitments to VRF service
    address vRFOracle; // Oracle committing to respond with VRF service
    bytes32 jobID; // ID of corresponding chainlink job the oracle's database
    uint256 fee; // Minimum payment for oracle response
  }

  mapping(bytes32 /* (provingKey, seed) */ => Callback) public callbacks;
  mapping(bytes32 /* provingKey */ => ServiceAgreement)
    public serviceAgreements;
  mapping(address => uint256) public withdrawableTokens; // Oracle LINK balances
  mapping(bytes32 /* provingKey */ => mapping(uint256 /* seed */ => bool))
    private observedSeeds;

  // The oracle only needs the jobID to look up the VRF, but specifying public
  // key as well prevents a malicious oracle from inducing VRF outputs from
  // another oracle by reusing the jobID. The sender will be useful later, for
  // whitelisting randomness requests.
  event RandomnessRequest(bytes32 keyHash, uint256 seed, bytes32 jobID,
    address sender, uint256 fee);
  event NewServiceAgreement(bytes32 keyHash, uint256 fee);

  /**
   * @notice Commits calling address to serve randomness
   * @param _fee minimum LINK payment required to serve randomness
   * @param _publicProvingKey public key used to prove randomness
   * @param _jobID ID of the corresponding chainlink job in the oracle's db
   */
  function registerProvingKey(
    uint256 _fee, uint256[2] calldata _publicProvingKey, bytes32 _jobID
  )
    external returns (bytes32 keyHash, address oracle, uint256 fee)
  {
    keyHash = hashOfKey(_publicProvingKey);
    address oldVRFOracle = serviceAgreements[keyHash].vRFOracle;
    require(oldVRFOracle == address(0), "please register a new key");
    serviceAgreements[keyHash].vRFOracle = msg.sender;
    serviceAgreements[keyHash].jobID = _jobID;
    serviceAgreements[keyHash].fee = _fee;
    emit NewServiceAgreement(keyHash, _fee);
    return (keyHash, serviceAgreements[keyHash].vRFOracle, _fee);
  }

  /**
   * @notice Called by LINK.transferAndCall, on successful LINK transfer
   *
   * @notice To invoke this, send LINK using transferAndCall. E.g.
   * @notice
   * @notice   LINK.transferAndCall(vrfCoordinator, _fee, abi.encode(_keyHash, _seed));
   * @notice
   * @notice where LINK is the address of the LINK contract, wrapped in
   * @notice LinkTokenInterface.
   * @notice
   * @notice The VRFCoordinator will call back to the calling contract when the
   * @notice oracle responds, on the method fulfillRandomness. See
   * @notice callbackMethod for its signature. Make sure to implement
   * @notice fulfillRandomness on your calling contract, or your request will
   * @notice fail.
   *
   * @dev TODO(alx): Make a VRFClient to take care of the above for the user.
   *
   * @param _sender address: who sent the LINK (must be a contract)
   * @param _fee amount of LINK sent
   * @param _data abi-encoded call to randomnessRequest
   *
   * @dev Memory layout of _data, as an abi-encoding of a call to randomnessRequest:
   *
   * @dev uint256 _data.length (32 bytes)
   * @dev bytes32 keyHash (32 bytes)
   * @dev uint256 seed (32 bytes)
   */
  function onTokenTransfer(address _sender, uint256 _fee, bytes memory _data)
    public
    onlyLINK
  {
    (bytes32 keyHash, uint256 seed) = abi.decode(_data, (bytes32, uint256));
    randomnessRequest(keyHash, seed, _fee, _sender);
  }

  /**
   * @notice creates the chainlink request for randomness
   */
  function randomnessRequest(
    bytes32 _keyHash,
    uint256 _seed,
    uint256 _feePaid,
    address _sender
  )
    internal
    sufficientLINK(_feePaid, _keyHash)
    isNewSeed(_keyHash, _seed)
  {
    bytes32 requestId = makeRequestId(_keyHash, _seed);
    assert(callbacks[requestId].callbackContract == address(0)); // Guaranteed by isNewSeed
    callbacks[requestId].callbackContract = _sender;
    callbacks[requestId].randomnessFee = _feePaid;
    callbacks[requestId].seed = _seed;
    emit RandomnessRequest(_keyHash, _seed, serviceAgreements[_keyHash].jobID,
      _sender, _feePaid);
  }

  /**
   * @notice Called by the chainlink node to fullfil requests
   * @param _proof the proof of randomness. Actual random output built from this
   *
   * @dev This is the main entrypoint for chainlink. If you change this, you 
   * @dev should also change the solidityABISstring in solidity_proof.go.
   */
  function fulfillRandomnessRequest(bytes memory _proof) public returns (bool) {
    // TODO(alx): Replace the public key out in the above proof with an argument
    // specifying the keyHash. Splice the key in here before sending it to
    // VRF.sol. Should be able to save about 2,000 gas that way.
    //
    // TODO(alx): Move this parsing into VRF.sol, where the bytes layout is recorded.
    uint256[2] memory publicKey;
    uint256 seed;
    // solhint-disable-next-line no-inline-assembly
    assembly { // Extract the public key and seed from proof
      publicKey := add(_proof, 0x20) // Skip length word for first 64 bytes
      seed := mload(add(_proof, 0xe0)) // Seed is 7th word in proof, plus word for length
    }
    bytes32 currentKeyHash = hashOfKey(publicKey);
    bytes32 requestId = makeRequestId(currentKeyHash, seed);
    Callback memory callback = callbacks[requestId];
    require(callback.callbackContract != address(0), "no corresponding request");
    uint256 randomness = VRF.randomValueFromVRFProof(_proof); // Reverts on failure
    observedSeeds[currentKeyHash][seed] = true;
    withdrawableTokens[serviceAgreements[currentKeyHash].vRFOracle] += callback.randomnessFee;
    bytes memory resp = abi.encodeWithSelector(callbackMethod, requestId, randomness);
    // solhint-disable-next-line avoid-low-level-calls
    (bool success,) = callback.callbackContract.call(resp);
    return success;
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
    withdrawableTokens[msg.sender] -= _amount;
    assert(LINK.transfer(_recipient, _amount));
  }

  // web3.utils.sha3("fulfillRandomness(bytes32,uint256)").slice(0, 10)
  bytes4 callbackMethod = 0x1f1f897f;

  /**
   * @notice Returns the id for this request
   * @param _keyHash The serviceAgreement ID to be used for this request
   * @param _seed The seed to be used in generating this randomness.
   */
  function makeRequestId(
    bytes32 _keyHash, uint256 _seed) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(_keyHash, _seed));
  }

  /**
   * @notice Returns the serviceAgreements key associated with this public key
   * @param _publicKey the key to return the address for
   */
  function hashOfKey(uint256[2] memory _publicKey) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(_publicKey));
  }

  /**
   * @dev Reverts if amount is not at least what was agreed upon in the service agreement
   * @param _feePaid The payment for the request
   * @param _keyHash The key which the request is for
   */
  modifier sufficientLINK(uint256 _feePaid, bytes32 _keyHash) {
    require(_feePaid >= serviceAgreements[_keyHash].fee, "Below agreed payment");
    _;
  }

/**
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    require(msg.sender == address(LINK), "Must use LINK token");
    _;
  }

  /**
   * @dev Reverts if the seed has been seen before, for this proving key.
   * @param _keyHash on which to check for prior request
   * @param _seed to check for prior request
   */
  modifier isNewSeed(bytes32 _keyHash, uint256 _seed) {
    require(!observedSeeds[_keyHash][_seed], "please request a novel seed");
    _;
  }

  /**
   * @dev Reverts if amount requested is greater than withdrawable balance
   * @param _amount The given amount to compare to `withdrawableTokens`
   */
  modifier hasAvailableFunds(uint256 _amount) {
    require(withdrawableTokens[msg.sender] >= _amount, "can't withdraw more than balance");
    _;
  }

}
