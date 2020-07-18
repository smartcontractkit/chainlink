pragma solidity 0.6.6;

import "./vendor/SafeMath.sol";

import "./interfaces/LinkTokenInterface.sol";

import "./VRF.sol";
import "./VRFRequestIDBase.sol";
import "./VRFConsumerBase.sol";

/**
 * @title VRFCoordinator coordinates on-chain verifiable-randomness requests
 * @title with off-chain responses
 */
contract VRFCoordinator is VRF, VRFRequestIDBase {

  using SafeMath for uint256;

  LinkTokenInterface internal LINK;

  constructor(address _link) public {
    LINK = LinkTokenInterface(_link);
  }

  struct Callback { // Tracks an ongoing request
    address callbackContract; // Requesting contract, which will receive response
    uint256 randomnessFee; // Amount of LINK paid at request time
    // Seed for the *oracle* to use in generating this random value. It is the
    // hash of the seed provided as input during a randomnessRequest, plus the
    // address of the contract making the request, plus an increasing nonce
    // specific to the VRF proving key and the calling contract. Including this
    // extra data in the VRF input seed helps to prevent unauthorized queries
    // against a VRF by any party who has prior knowledge of the requester's
    // prospective seed. Only the specified contract can make that request.
    uint256 seed;
  }

  struct ServiceAgreement { // Tracks oracle commitments to VRF service
    address vRFOracle; // Oracle committing to respond with VRF service
    bytes32 jobID; // ID of corresponding chainlink job in oracle's DB
    uint256 fee; // Minimum payment for oracle response
  }

  mapping(bytes32 /* (provingKey, seed) */ => Callback) public callbacks;
  mapping(bytes32 /* provingKey */ => ServiceAgreement)
    public serviceAgreements;
  mapping(address /* oracle */ => uint256 /* LINK balance */)
    public withdrawableTokens;
  mapping(bytes32 /* provingKey */ => mapping(address /* consumer */ => uint256))
    private nonces;

  // The oracle only needs the jobID to look up the VRF, but specifying public
  // key as well prevents a malicious oracle from inducing VRF outputs from
  // another oracle by reusing the jobID.
  event RandomnessRequest(
    bytes32 keyHash,
    uint256 seed,
    bytes32 indexed jobID,
    address sender,
    uint256 fee);

  event NewServiceAgreement(bytes32 keyHash, uint256 fee);

  /**
   * @notice Commits calling address to serve randomness
   * @param _fee minimum LINK payment required to serve randomness
   * @param _oracle the address of the Chainlink node with the proving key and job
   * @param _publicProvingKey public key used to prove randomness
   * @param _jobID ID of the corresponding chainlink job in the oracle's db
   */
  function registerProvingKey(
    uint256 _fee, address _oracle, uint256[2] calldata _publicProvingKey, bytes32 _jobID
  )
    external
  {
    bytes32 keyHash = hashOfKey(_publicProvingKey);
    address oldVRFOracle = serviceAgreements[keyHash].vRFOracle;
    require(oldVRFOracle == address(0), "please register a new key");
    require(_oracle != address(0), "_oracle must not be 0x0");
    serviceAgreements[keyHash].vRFOracle = _oracle;
    serviceAgreements[keyHash].jobID = _jobID;
    serviceAgreements[keyHash].fee = _fee;
    emit NewServiceAgreement(keyHash, _fee);
  }

  /**
   * @notice Called by LINK.transferAndCall, on successful LINK transfer
   *
   * @dev To invoke this, use the requestRandomness method in VRFConsumerBase.
   *
   * @dev The VRFCoordinator will call back to the calling contract when the
   * @dev oracle responds, on the method fulfillRandomness. See
   * @dev VRFConsumerBase.fullfilRandomnessRequest for its signature. Your
   * @dev consuming contract should inherit from VRFConsumerBase, and implement
   * @dev fullfilRandomnessRequest.
   *
   * @param _sender address: who sent the LINK (must be a contract)
   * @param _fee amount of LINK sent
   * @param _data abi-encoded call to randomnessRequest
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
   *
   * @param _keyHash ID of the VRF public key against which to generate output
   * @param _seed Input to the VRF, from which randomness is generated
   * @param _feePaid Amount of LINK sent with request. Must exceed fee for key
   * @param _sender Requesting contract; to be called back with VRF output
   */
  function randomnessRequest(
    bytes32 _keyHash,
    uint256 _seed,
    uint256 _feePaid,
    address _sender
  )
    internal
    sufficientLINK(_feePaid, _keyHash)
  {
    uint256 nonce = nonces[_keyHash][_sender];
    uint256 seed = makeVRFInputSeed(_keyHash, _seed, _sender, nonce);
    bytes32 requestId = makeRequestId(_keyHash, seed);
    // Cryptographically guaranteed by seed including an increasing nonce
    assert(callbacks[requestId].callbackContract == address(0));
    callbacks[requestId].callbackContract = _sender;
    callbacks[requestId].randomnessFee = _feePaid;
    callbacks[requestId].seed = seed;
    emit RandomnessRequest(_keyHash, seed, serviceAgreements[_keyHash].jobID,
      _sender, _feePaid);
    nonces[_keyHash][_sender] = nonces[_keyHash][_sender].add(1);
  }

  /**
   * @notice Called by the chainlink node to fullfil requests
   * @param _proof the proof of randomness. Actual random output built from this
   *
   * @dev This is the main entrypoint for chainlink. If you change this, you
   * @dev should also change the solidityABISstring in solidity_proof.go.
   */
  function fulfillRandomnessRequest(bytes memory _proof) public returns (bool) {
    // TODO(alx): Replace the public key in the above proof with an argument
    // specifying the keyHash. Splice the key in here before sending it to
    // VRF.sol. Should be able to save about 2,000 gas that way.
    // https://www.pivotaltracker.com/story/show/170828567
    //
    // TODO(alx): Move this parsing into VRF.sol, where the bytes layout is recorded.
    // https://www.pivotaltracker.com/story/show/170828697
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
    address oadd = serviceAgreements[currentKeyHash].vRFOracle;
    withdrawableTokens[oadd] = withdrawableTokens[oadd].add(callback.randomnessFee);
    // Dummy variable; allows access to method selector in next line. See
    // https://github.com/ethereum/solidity/issues/3506#issuecomment-553727797
    VRFConsumerBase v;
    bytes memory resp = abi.encodeWithSelector(
      v.rawFulfillRandomness.selector, requestId, randomness);
    // solhint-disable-next-line avoid-low-level-calls
    (bool success,) = callback.callbackContract.call(resp);
    delete callbacks[requestId]; // Be a good ethereum citizen
    return success;
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
   * @dev Reverts if amount requested is greater than withdrawable balance
   * @param _amount The given amount to compare to `withdrawableTokens`
   */
  modifier hasAvailableFunds(uint256 _amount) {
    require(withdrawableTokens[msg.sender] >= _amount, "can't withdraw more than balance");
    _;
  }

}
