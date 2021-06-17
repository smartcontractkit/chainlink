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

    // TODO: struct packing
    event Payment(uint256 seed, uint256 v2);
    event SubscriptionCreated(uint256 subId, address owner, address[] consumers);
    event SubscriptionFundsAdded(uint256 subId, uint256 oldBalance, uint256 newBalance);
    event SubscriptionConsumersUpdated(uint256 subId, address[] oldConsumers, address[] newConsumers);
    event SubscriptionFundsWithdrawn(uint256 subId, uint256 oldBalance, uint256 newBalance);
    event SubscriptionCanceled(uint256 subId);
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
    mapping(address /* oracle */ => uint256 /* LINK balance */) public s_withdrawableTokens;
    mapping(bytes32 => mapping(address /* consumer */ => uint256)) s_nonces;

    event RandomWordsRequested(
        bytes32 indexed keyHash,
        uint16 minimumRequestConfirmations,
        uint16 callbackGasLimit,
        uint256 preSeed,
        uint256 subId);
    event RandomWordsFulfilled(
        uint256 requestId, uint256[] output);
    struct Callback {
        address callbackContract; // Requesting contract, which will receive response
        uint256 callbackGasLimit;
        uint256 numWords;
        uint256 subId;
        bytes32 seedAndBlockNum;
    }
    mapping(uint256 /* requestID */ => Callback) public s_callbacks;


    bytes4 constant private FULFILL_RANDOM_WORDS_SELECTOR = bytes4(keccak256("fulfillRandomWords(bytes32,[]uint256)"));


    uint16 private s_minimumRequestBlockConfirmations = 3;
    uint16 private s_maxConsumersPerSubscription = 10;
    uint32 private s_stalenessSeconds = 0;
    int256 private s_fallbackGasPrice = 200;
    int256 private s_fallbackLinkPrice = 200000;

    constructor(address link, address blockHashStore, address linkEthFeed, address fastGasFeed) {
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
        uint16 maxConsumersPerSubscription,
        uint32 stalenessSeconds,
        int256 fallbackGasPrice,
        int256 fallbackLinkPrice
    )
    external
    onlyOwner()
    {
        s_maxConsumersPerSubscription = maxConsumersPerSubscription;
        s_minimumRequestBlockConfirmations = minimumRequestBlockConfirmations;
        s_stalenessSeconds = stalenessSeconds;
        s_fallbackGasPrice = fallbackGasPrice;
        s_fallbackLinkPrice = fallbackLinkPrice;
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
       require(s_serviceAgreements[keyHash].oracle != address(0), "must be a registered key");
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
       emit RandomWordsRequested(keyHash, minimumRequestConfirmations, callbackGasLimit, preSeedAndRequestId, subId);
       return preSeedAndRequestId;
    }

    // Offsets into fulfillRandomnessRequest's _proof of various values
    //
    // Public key. Skips byte array's length prefix.
    uint256 public constant PUBLIC_KEY_OFFSET = 0x20;
    // Seed is 7th word in proof, plus word for length, (6+1)*0x20=0xe0
    uint256 public constant PRESEED_OFFSET = 0xe0;
    // TODO: Gas for making payment itself
    uint256 public constant GAS_BUFFER = 2_000;

    function fulfillRandomWords(
        bytes memory _proof
    )
    external
    {
        // TODO: maybe fail fast on an invalid keyHash?
        uint256 startGas = gasleft();
        (bytes32 keyHash, Callback memory callback, uint256 requestId,
        uint256 randomness) = getRandomnessFromProof(_proof);
        uint256[] memory randomWords = new uint256[](callback.numWords);
        for (uint256 i = 0; i < callback.numWords; i++) {
            randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
        }

        bytes memory resp = abi.encodeWithSelector(FULFILL_RANDOM_WORDS_SELECTOR, requestId, randomWords);
        // TODO: Make more exact
        uint256 payment = calculatePaymentAmount(startGas);
        s_subscriptions[callback.subId].balance -= payment;
        s_withdrawableTokens[s_serviceAgreements[keyHash].oracle] += payment;
//
//        require(gasleft() > callback.callbackGasLimit, "not enough gas for consumer");
//
//        (bool success,) = callback.callbackContract.call(resp);
//        // Avoid unused-local-variable warning. (success is only present to prevent
//        // a warning that the return value of consumerContract.call is unused.)
//        (success);
        emit RandomWordsFulfilled(requestId, randomWords);
    }

    function calculatePaymentAmount(
        uint256 startGas
    )
    public
    returns (uint256)
    {
        // Get the amount of gas used for (fulfillment + request)
        uint256 gasWei; // wei/gas i.e. gasPrice
        uint256 linkWei; // link/wei i.e. link price in wei.
        (gasWei, linkWei) = getFeedData();
        // (1e18 linkWei/link) (wei/gas * gas) / (wei/link) = linkWei
        return 1e18*gasWei*(GAS_BUFFER + startGas - gasleft()) / linkWei;
    }

    function getRandomnessFromProof(bytes memory _proof)
    internal view returns (bytes32 currentKeyHash, Callback memory callback,
        uint256 requestId, uint256 randomness) {
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
        callback = s_callbacks[preSeed];
        requestId = preSeed;
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
        uint256,
        uint256
    )
    {
        uint32 stalenessSeconds = s_stalenessSeconds;
        bool staleFallback = stalenessSeconds > 0;
        uint256 timestamp;
        int256 gasWei;
        int256 linkEth;
        (,gasWei,,timestamp,) = FAST_GAS_FEED.latestRoundData();
        if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
            gasWei = s_fallbackGasPrice;
        }
        (,linkEth,,timestamp,) = LINK_ETH_FEED.latestRoundData();
        if (staleFallback && stalenessSeconds < block.timestamp - timestamp) {
            linkEth = s_fallbackLinkPrice;
        }
        return (uint256(gasWei), uint256(linkEth));
    }

    /*
        Subscription management, to be handled by a single account/contract.
    */
    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint256)
    {
        require(consumers.length <= s_maxConsumersPerSubscription, "above max consumers per sub");
        allConsumersValid(consumers);
        currentSubId++;
        s_subscriptions[currentSubId] = Subscription({
            owner: msg.sender,
            subId: currentSubId,
            consumers: consumers,
            balance: 0
        });
        emit SubscriptionCreated(currentSubId, msg.sender, consumers);
        // TODO: optionally fund also in the creation transaction? We'd still need a separate
        // fund tx anyways to top it up, but we'd save a tx
        return currentSubId;
    }

    function allConsumersValid(address[] memory consumers) internal {
        require(consumers.length <= s_maxConsumersPerSubscription, "above max consumers per sub");
        for (uint i = 0; i < consumers.length; i++) {
            require(consumers[i] != address(0), "consumer address must not be zero");
        }
    }

    function updateSubscription(
        uint256 subId,
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    {
        require(msg.sender == s_subscriptions[subId].owner, "only subscription owner can update");
        allConsumersValid(consumers);
        address[] memory oldConsumers = s_subscriptions[subId].consumers;
        s_subscriptions[subId].consumers = consumers;
        emit SubscriptionConsumersUpdated(subId, oldConsumers, consumers);
    }

    function fundSubscription(
        uint256 subId,
        uint256 amount
    )
    external
    {
        require(s_subscriptions[subId].owner != address(0), "subID doesnt exist");
        require(msg.sender == s_subscriptions[subId].owner, "only subscription owner can fund");
        uint256 oldBalance = s_subscriptions[subId].balance;
        s_subscriptions[subId].balance += amount;
        LINK.transferFrom(msg.sender, address(this), amount);
        emit SubscriptionFundsAdded(subId, oldBalance, s_subscriptions[subId].balance);
    }

    function withdrawFromSubscription(
        uint256 subId,
        address to,
        uint256 amount
    )
    external
    {
        require(msg.sender == s_subscriptions[subId].owner, "only subscription owner can withdraw");
        require(s_subscriptions[subId].balance >= amount, "insufficient balance");
        uint256 oldBalance = s_subscriptions[subId].balance;
        s_subscriptions[subId].balance -= amount;
        LINK.transfer(to, amount);
        emit SubscriptionFundsWithdrawn(subId, oldBalance, s_subscriptions[subId].balance);
    }

// CONTRACT TOO LARGE IF THIS IS INCLUDED
//    // Keep this separate from zeroing, perhaps there is a use case where consumers
//    // want to keep the subId, but withdraw all the link.
//    function cancelSubscription(
//        uint256 subId
//    )
//    external
//    {
//        require(msg.sender == s_subscriptions[subId].owner, "only subscription owner can cancel");
//        require(s_subscriptions[subId].balance == 0, "balance must be zero to cancel");
//        delete s_subscriptions[subId];
//        emit SubscriptionCanceled(subId);
//    }
}
