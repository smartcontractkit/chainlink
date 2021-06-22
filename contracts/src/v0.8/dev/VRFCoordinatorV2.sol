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

    event SubscriptionCreated(uint64 subId, address owner, address[] consumers);
    event SubscriptionFundsAdded(uint64 subId, uint256 oldBalance, uint256 newBalance);
    event SubscriptionConsumersUpdated(uint64 subId, address[] oldConsumers, address[] newConsumers);
    event SubscriptionFundsWithdrawn(uint64 subId, uint256 oldBalance, uint256 newBalance);
    event SubscriptionCanceled(uint64 subId);
    uint64 currentSubId;
    struct Subscription {
        address owner; // Owner can fund/withdraw/cancel the sub
        address[] consumers; // List of addresses which can consume using this subscription.
        uint256 balance; // Common balance used for all consumer requests.
    }
    mapping(uint64 /* subId */ => Subscription /* subscription */) private s_subscriptions;

    event NewServiceAgreement(bytes32 keyHash, address oracle);
    mapping(bytes32 /* keyHash */ => address /* oracle */) private s_serviceAgreements;
    mapping(address /* oracle */ => uint256 /* LINK balance */) private s_withdrawableTokens;
    mapping(bytes32 => mapping(address /* consumer */ => uint256)) private s_nonces;

    event RandomWordsRequested(
        bytes32 indexed keyHash,
        uint256 preSeedAndRequestId,
        uint64 subId,
        uint64 minimumRequestConfirmations,
        uint64 callbackGasLimit,
        uint64 numWords,
        address sender);
    event RandomWordsFulfilled(
        uint256 requestId,
        uint256[] output,
        bool success);
    // Just to relieve stack pressure
    struct FulfillmentParams {
        uint64 subId;
        uint64 callbackGasLimit;
        uint64 numWords;
        address sender;
    }
    mapping(uint256 /* requestID */ => bytes32) private s_callbacks;

    bytes4 constant private FULFILL_RANDOM_WORDS_SELECTOR = bytes4(keccak256("fulfillRandomWords(uint256,uint256[])"));

    struct Config {
        uint16 minimumRequestBlockConfirmations;
        uint16 maxConsumersPerSubscription;
        uint32 stalenessSeconds;
        // Gas to cover oracle payment + fulfillment log after we calculate the payment.
        // Should be approximately: 5k (change oracle balance) + 5k (change sub balance) + ~1k (log)
        // However we make it configurable in case those operations are repriced.
        uint32 gasAfterPaymentCalculation;
    }
    Config private s_config;
    int256 private s_fallbackGasPrice;
    int256 private s_fallbackLinkPrice;
    event ConfigSet(
        uint16 minimumRequestBlockConfirmations,
        uint16 maxConsumersPerSubscription,
        uint32 stalenessSeconds,
        uint32 gasAfterPaymentCalculation,
        int256 fallbackGasPrice,
        int256 fallbackLinkPrice
    );

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
        require(s_serviceAgreements[kh] == address(0), "cannot re-register the same proving key");
        s_serviceAgreements[kh] = oracle;
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
        uint32 gasAfterPaymentCalculation,
        int256 fallbackGasPrice,
        int256 fallbackLinkPrice
    )
    external
    onlyOwner()
    {
        s_config = Config({
            minimumRequestBlockConfirmations: minimumRequestBlockConfirmations,
            maxConsumersPerSubscription: maxConsumersPerSubscription,
            stalenessSeconds: stalenessSeconds,
            gasAfterPaymentCalculation: gasAfterPaymentCalculation
        });
        s_fallbackGasPrice = fallbackGasPrice;
        s_fallbackLinkPrice = fallbackLinkPrice;
        emit ConfigSet(minimumRequestBlockConfirmations,
            maxConsumersPerSubscription,
            stalenessSeconds,
            gasAfterPaymentCalculation,
            fallbackGasPrice,
            fallbackLinkPrice
        );
    }

    /**
     * @notice read the current configuration of the coordinator.
     */
    function getConfig()
    external
    view
    returns (
        uint16 minimumRequestBlockConfirmations,
        uint16 maxConsumersPerSubscription,
        uint32 stalenessSeconds,
        uint32 gasAfterPaymentCalculation,
        int256 fallbackGasPrice,
        int256 fallbackLinkPrice
    )
    {
        Config memory config = s_config;
        return (
            config.minimumRequestBlockConfirmations,
            config.maxConsumersPerSubscription,
            config.stalenessSeconds,
            config.gasAfterPaymentCalculation,
            s_fallbackGasPrice,
            s_fallbackLinkPrice
        );
    }

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint64  subId,
        uint64  minimumRequestConfirmations,
        uint64  callbackGasLimit,
        uint64  numWords  // Desired number of random words
    )
    external
    returns (uint256 requestId)
    {
       // One would consider doing a sanity check that the subscription has enough link here,
       // but:
       //   1) you'd need to either store a cached price or do a lookup which are both costly
       //   2) set a fixed min which would be hard to estimate
       //   3) price could fluctuate anyways between request and response, so it could still fail to fulfill
       //   4) if it did fail to fulfill, you could simply top up the subscription and re-try and it would succeed
       require(s_subscriptions[subId].owner != address(0), "invalid subId");
       bool validConsumer;
       for (uint16 i = 0; i < s_subscriptions[subId].consumers.length; i++) {
           if (s_subscriptions[subId].consumers[i] == msg.sender) {
               validConsumer = true;
               break;
           }
       }
       require(validConsumer, "invalid consumer");
       require(s_serviceAgreements[keyHash] != address(0), "must be a registered key");
       uint256 nonce = s_nonces[keyHash][msg.sender] + 1;
       s_nonces[keyHash][msg.sender] = nonce;
       uint256 preSeedAndRequestId = uint256(keccak256(abi.encode(keyHash, msg.sender, nonce)));
       // Min req confirmations not needed as part of fulfillment, leave out of the commitment
       s_callbacks[preSeedAndRequestId] = keccak256(abi.encodePacked(preSeedAndRequestId, block.number, subId, callbackGasLimit, numWords, msg.sender));
       emit RandomWordsRequested(keyHash, preSeedAndRequestId, subId, minimumRequestConfirmations, callbackGasLimit, numWords, msg.sender);
       return preSeedAndRequestId;
    }

    function getCallback(uint256 requestId)
    external
    view
    returns (bytes32){
        return s_callbacks[requestId];
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
        uint256 startGas = gasleft();
        (bytes32 keyHash, uint256 requestId,
        uint256 randomness, FulfillmentParams memory fp) = getRandomnessFromProof(_proof);

        uint256[] memory randomWords = new uint256[](fp.numWords);
        for (uint256 i = 0; i < fp.numWords; i++) {
            randomWords[i] = uint256(keccak256(abi.encode(randomness, i)));
        }

        // Prevent re-entrancy. The user callback cannot call fulfillRandomWords again
        // with the same proof because this getRandomnessFromProof will revert because the requestId
        // is gone.
        delete s_callbacks[requestId];
        bytes memory resp = abi.encodeWithSelector(FULFILL_RANDOM_WORDS_SELECTOR, requestId, randomWords);
        require(gasleft() > fp.callbackGasLimit, "not enough gas for consumer");
        (bool success,) = fp.sender.call(resp);
        // Avoid unused-local-variable warning. (success is only present to prevent
        // a warning that the return value of consumerContract.call is unused.)
        (success);

        // We want to charge users exactly for how much gas they use in their callback.
        // The gasAfterPaymentCalculation is meant to cover these to additional operations where we
        // decrement the subscription balance and increment the oracles withdrawable balance.
        uint256 payment = calculatePaymentAmount(startGas, s_config.gasAfterPaymentCalculation);

        s_subscriptions[fp.subId].balance -= payment;
        s_withdrawableTokens[s_serviceAgreements[keyHash]] += payment;
        emit RandomWordsFulfilled(requestId, randomWords, success);
    }

    function calculatePaymentAmount(
        uint256 startGas,
        uint256 gasBuffer
    )
    private
    returns (uint256)
    {
        // Get the amount of gas used for (fulfillment + request)
        uint256 gasWei; // wei/gas i.e. gasPrice
        uint256 linkWei; // link/wei i.e. link price in wei.
        (gasWei, linkWei) = getFeedData();
        // (1e18 linkWei/link) (wei/gas * gas) / (wei/link) = linkWei
        return 1e18*gasWei*(gasBuffer + startGas - gasleft()) / linkWei;
    }

    function getRandomnessFromProof(bytes memory _proof)
    public view returns (bytes32 currentKeyHash,
        uint256 requestId, uint256 randomness, FulfillmentParams memory fp) {
        // blockNum follows proof, which follows length word (only direct-number
        // constants are allowed in assembly, so have to compute this in code)
        uint256 BLOCKNUM_OFFSET = 0x20 + PROOF_LENGTH;
        // Note that _proof.length skips the initial length word.
        // We expected the total length to be proof + 5 words (blocknum, subId, callbackLimit, nw, sender)
        require(_proof.length == PROOF_LENGTH + 0x20*5, "wrong proof length");
        uint256[2] memory publicKey;
        uint256 preSeed;
        uint256 blockNum;
        address sender;
        assembly { // solhint-disable-line no-inline-assembly
            publicKey := add(_proof, PUBLIC_KEY_OFFSET)
            preSeed := mload(add(_proof, PRESEED_OFFSET))
            blockNum := mload(add(_proof, BLOCKNUM_OFFSET))
            // We use a struct to limit local variables to avoid stack depth errors.
            mstore(fp, mload(add(add(_proof, BLOCKNUM_OFFSET), 0x20)))
            mstore(add(fp, 0x20), mload(add(add(_proof, BLOCKNUM_OFFSET), 0x40)))
            mstore(add(fp, 0x40), mload(add(add(_proof, BLOCKNUM_OFFSET), 0x60)))
            sender := mload(add(add(_proof, BLOCKNUM_OFFSET), 0x80))
        }
        currentKeyHash = hashOfKey(publicKey);
        bytes32 callback = s_callbacks[preSeed];
        requestId = preSeed;
        require(callback != 0, "no corresponding request");
        require(callback == keccak256(abi.encodePacked(requestId, blockNum, fp.subId, fp.callbackGasLimit, fp.numWords, sender)), "incorrect commitment");
        fp.sender = sender;

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
        uint32 stalenessSeconds = s_config.stalenessSeconds;
        bool staleFallback = stalenessSeconds > 0;
        uint256 timestamp;
        int256 gasWei;
        int256 linkEth;
        // Fallback to the fallback price if the feed is too stale.
        // Maybe need to optimize to avoid this contract call and used a cached
        // price, say in the case of a large number of fulfillments in short succession?
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

    function withdraw(address _recipient, uint256 _amount)
    external
    {
        // Will revert if insufficient funds
        s_withdrawableTokens[msg.sender] = s_withdrawableTokens[msg.sender] - _amount;
        assert(LINK.transfer(_recipient, _amount));
    }

    function getSubscription(uint64 subId)
    external
    view
    returns (Subscription memory)
    {
        return s_subscriptions[subId];
    }

    function createSubscription(
        address[] memory consumers // permitted consumers of the subscription
    )
    external
    returns (uint64)
    {
        require(consumers.length <= s_config.maxConsumersPerSubscription, "above max consumers per sub");
        allConsumersValid(consumers);
        currentSubId++;
        s_subscriptions[currentSubId] = Subscription({
            owner: msg.sender,
            consumers: consumers,
            balance: 0
        });
        emit SubscriptionCreated(currentSubId, msg.sender, consumers);
        return currentSubId;
    }

    function allConsumersValid(address[] memory consumers) internal {
        require(consumers.length <= s_config.maxConsumersPerSubscription, "above max consumers per sub");
        for (uint i = 0; i < consumers.length; i++) {
            require(consumers[i] != address(0), "consumer address must not be zero");
        }
    }

    function updateSubscription(
        uint64 subId,
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
        uint64 subId,
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
        uint64 subId,
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

    // Keep this separate from zeroing, perhaps there is a use case where consumers
    // want to keep the subId, but withdraw all the link.
    function cancelSubscription(
        uint64 subId
    )
    external
    {
        require(msg.sender == s_subscriptions[subId].owner, "only subscription owner can cancel");
        require(s_subscriptions[subId].balance == 0, "balance must be zero to cancel");
        delete s_subscriptions[subId];
        emit SubscriptionCanceled(subId);
    }
}
