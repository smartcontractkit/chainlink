pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockHashStoreInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";

import "../vendor/Ownable.sol";

import "./VRF.sol";

contract VRFCoordinatorV2 is VRF, Ownable {

    LinkTokenInterface public immutable LINK;
    AggregatorV3Interface public immutable LINK_ETH_FEED;
    AggregatorV3Interface public immutable FAST_GAS_FEED;
    BlockHashStoreInterface public immutable BLOCKHASH_STORE;

    event SubscriptionCreated(uint256 subId, address owner);
    event SubscriptionFunded(uint256 subId, uint256 amount);
    uint256 currentSubId;
    struct Subscription {
        uint256 subId;
        address owner; // Owner can fund/withdraw/cancel the sub
        address[] consumers; // List of addresses which can consume using this subscription.
        uint256 balance; // Common balance used for all consumer requests.
    }
    mapping(uint256 /* subId */ => Subscription /* subscription */) public s_subscriptions;

    event NewServiceAgreement(bytes32 keyHash, address oracle);
    struct ServiceAgreement {
        address oracle;
        bytes32 keyHash;
    }
    mapping(bytes32 /* keyHash */ => ServiceAgreement) public s_serviceAgreements;
    mapping(bytes32 => mapping(address /* consumer */ => uint256)) s_nonces;

    event RandomWordsRequested(
        bytes32 indexed keyHash,
        uint16 minimumRequestConfirmations,
        uint16 callbackGasLimit,
        uint256 preSeed);
    struct Callback {
        address callbackContract; // Requesting contract, which will receive response
        uint256 callbackGasLimit;
        uint256 numWords;
        uint256 subId;
        bytes32 seedAndBlockNum;
    }
    mapping(uint256 /* requestID */ => Callback) public s_callbacks;


    uint16 internal s_minimumRequestBlockConfirmations = 3;
    uint16 internal s_maxConsumersPerSubscription = 10;

    constructor(address link, address blockHashStore, address linkEthFeed, address fastGasFeed) public {
        LINK = LinkTokenInterface(link);
        LINK_ETH_FEED = AggregatorV3Interface(linkEthFeed);
        FAST_GAS_FEED = AggregatorV3Interface(fastGasFeed);
        BLOCKHASH_STORE = BlockHashStoreInterface(blockHashStore);
    }

    function registerProvingKey(
        address oracle, uint256[2] calldata publicProvingKey
    )
    external
    onlyOwner()
    {
        bytes32 kh = hashOfKey(publicProvingKey);
        require(s_serviceAgreements[kh].oracle == address(0), "cannot re-register the same proving key");
        s_serviceAgreements[kh] = ServiceAgreement({
            oracle: oracle,
            keyHash: kh
        });
        emit NewServiceAgreement(kh, oracle);
    }

    /**
     * @notice Returns the serviceAgreements key associated with this public key
     * @param _publicKey the key to return the address for
     */
    function hashOfKey(uint256[2] memory _publicKey) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(_publicKey));
    }

    function setConfig(
        uint16 minimumRequestBlockConfirmations,
        uint16 maxConsumersPerSubscription
        // TODO: Add fallback fees and timeout params
    )
    external
    onlyOwner()
    {
        s_maxConsumersPerSubscription = maxConsumersPerSubscription;
        s_minimumRequestBlockConfirmations = minimumRequestBlockConfirmations;
    }

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint16  minimumRequestConfirmations,
        uint16  callbackGasLimit,
        uint256 subId,   // A data structure for billing
        uint256 numWords  // Desired number of random words
    )
    external
    returns (uint256 requestId)
    {
       // Sanity check the subscription has enough link? Just
       // accept that gas price fluctuations between request and response could potentially
       // result in request accepted but failed to fulfill
       require(s_subscriptions[subId].owner != address(0), "invalid subId");
       uint256 nonce = s_nonces[keyHash][msg.sender] + 1;
       s_nonces[keyHash][msg.sender] = nonce;
       uint256 preSeedAndRequestId = uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)));
       s_callbacks[preSeedAndRequestId] = Callback({
            callbackContract: msg.sender,
            callbackGasLimit: callbackGasLimit, // For sanity checking
            numWords: numWords,
            subId: subId,
            seedAndBlockNum: keccak256(abi.encodePacked(preSeedAndRequestId, block.number))
       });
       emit RandomWordsRequested(keyHash, minimumRequestConfirmations, callbackGasLimit, preSeedAndRequestId);
       return preSeedAndRequestId;
    }

    // Offsets into fulfillRandomnessRequest's _proof of various values
    //
    // Public key. Skips byte array's length prefix.
    uint256 public constant PUBLIC_KEY_OFFSET = 0x20;
    // Seed is 7th word in proof, plus word for length, (6+1)*0x20=0xe0
    uint256 public constant PRESEED_OFFSET = 0xe0;

    function fulfillRandomWords(
        bytes memory _proof
    )
    external
    {
        // TODO:
        // 1. Verify proof, extract random value, public key and preSeed.
        // 2. Lookup the blockhash (from the store if needed)
        // 3. Get the requestId = hash(abiencode(hash(pk), preseed))
        // 4. Look up the callback = callbacks[requestId] for the callback address
        // 5. Expand the randomness
        // 6. Calculate gas used up to this point, convert to link, charge the subscription and delete callback.
        // 7. Ensure we have the required gasLimit, call the callback with the specified number of words.
        // TODO: maybe fail fast on an invalid keyHash?

        (bytes32 keyHash, Callback memory callback, bytes32 requestId,
        uint256 randomness) = getRandomnessFromProof(_proof);
        // TODO: calculate payment amount and pay s_serviceAgreements[keyHash].oracle
    }

    function getRandomnessFromProof(bytes memory _proof)
    internal view returns (bytes32 currentKeyHash, Callback memory callback,
        bytes32 requestId, uint256 randomness) {
        // blockNum follows proof, which follows length word (only direct-number
        // constants are allowed in assembly, so have to compute this in code)
        uint256 BLOCKNUM_OFFSET = 0x20 + PROOF_LENGTH;
        // _proof.length skips the initial length word, so not including the
        // blocknum in this length check balances out.
        require(_proof.length == BLOCKNUM_OFFSET, "wrong proof length");
        uint256[2] memory publicKey;
        uint256 preSeed;
        uint256 blockNum;
        assembly { // solhint-disable-line no-inline-assembly
            publicKey := add(_proof, PUBLIC_KEY_OFFSET)
            preSeed := mload(add(_proof, PRESEED_OFFSET))
            blockNum := mload(add(_proof, BLOCKNUM_OFFSET))
        }
        currentKeyHash = hashOfKey(publicKey);
//        requestId = makeRequestId(currentKeyHash, preSeed);
        callback = s_callbacks[preSeed];
        require(callback.callbackContract != address(0), "no corresponding request");
        require(callback.seedAndBlockNum == keccak256(abi.encodePacked(preSeed,
            blockNum)), "wrong preSeed or block num");

        bytes32 blockHash = blockhash(blockNum);
        if (blockHash == bytes32(0)) {
            blockHash = BLOCKHASH_STORE.getBlockhash(blockNum);
            require(blockHash != bytes32(0), "please prove blockhash");
        }
        // The seed actually used by the VRF machinery, mixing in the blockhash
        uint256 actualSeed = uint256(keccak256(abi.encodePacked(preSeed, blockHash)));
        // solhint-disable-next-line no-inline-assembly
        assembly { // Construct the actual proof from the remains of _proof
            mstore(add(_proof, PRESEED_OFFSET), actualSeed)
            mstore(_proof, PROOF_LENGTH)
        }
        randomness = VRF.randomValueFromVRFProof(_proof); // Reverts on failure
    }

    function getFeedData()
    private
    view
    returns (
        int256 gasWei,
        int256 linkEth
    )
    {
//        uint32 stalenessSeconds = s_config.stalenessSeconds;
//        bool staleFallback = stalenessSeconds > 0;
        uint256 timestamp;
        (,gasWei,,timestamp,) = FAST_GAS_FEED.latestRoundData();
//        if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
//            gasWei = s_fallbackGasPrice;
//        }
        (,linkEth,,timestamp,) = LINK_ETH_FEED.latestRoundData();
//        if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
//            linkEth = s_fallbackLinkPrice;
//        }
        return (gasWei, linkEth);
    }

    /*
        Subscription management, to be handled by a single account/contract.
    */
    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint256 subId)
    {
        require(consumers.length <= s_maxConsumersPerSubscription, "above max consumers per sub");
        for (uint i = 0; i < consumers.length; i++) {
            require(consumers[i] != address(0), "consumer address must not be zero");
        }
        currentSubId++;
        s_subscriptions[currentSubId] = Subscription({
            owner: msg.sender,
            subId: currentSubId,
            consumers: consumers,
            balance: 0
        });
        emit SubscriptionCreated(currentSubId, msg.sender);
        return subId;
    }

    function updateSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    {
        // TODO: Only the subscription owner, valid sub must exist
        // TODO: No addresses can be zero, set max number of callers
        // subscriptions[currentSubId].consumers = consumers;
        // TODO: emit some logs
        // return currentSubId;
    }

    function fundSubscription(
        uint256 subId,
        uint256 amount
    )
    external
    {
        // TODO check subId, amount is valid, only owner
        s_subscriptions[subId].balance += amount;
        LINK.transferFrom(msg.sender, address(this), amount);
//        // Maybe old and new balance?
//        emit SubscriptionFunded(subId, amount);
    }

    function withdrawFromSubscription(
        uint256 subId,
        address to,
        uint256 amount
    )
    external
    {
        // TODO check subId, amount is valid, only owner
        // subscriptions[subId].balance -= amount;
        // LINK.transfer(address(this), to, amount);
        // TODO: emit some logs
    }

    function cancelSubscription(
        uint256 subId
    )
    external
    {
        // TODO check subId, only owner, must be zeroed
        // Delete the subscription
        // TODO: emit some logs
    }
}
