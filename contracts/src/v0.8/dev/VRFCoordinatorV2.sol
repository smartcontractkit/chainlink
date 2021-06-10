pragma solidity ^0.8.0;

import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/BlockHashStoreInterface.sol";

import "../vendor/Ownable.sol";

import "./VRF.sol";

contract VRFCoordinatorV2 is VRF, Ownable {

    LinkTokenInterface internal LINK;
    BlockHashStoreInterface internal blockHashStore;

    event SubscriptionCreated(uint256 subId, address owner);
    uint256 currentSubId;
    struct Subscription {
        uint256 subId;
        address owner; // Owner can fund/withdraw/cancel the sub
        address[] consumers; // List of addresses which can consume using this subscription.
        uint256 balance; // Common balance used for all consumer requests.
    }
    mapping(uint256 /* subId */ => Subscription /* subscription */) public subscriptions;

    struct Oracle {
        address oracle;
        bytes32 keyHash;
        mapping(address => uint256) nonces;
    }
    mapping(bytes32 /* keyHash */ => Oracle) public oracles;

    struct Callback {
        address callbackContract; // Requesting contract, which will receive response
        bytes32 requestBlockNum; // Request block number
        uint256 numWords;
        uint256 subId;
    }
    mapping(bytes32 /* requestID */ => Callback) public callbacks;


    constructor(address _link, address _blockHashStore) public {
        LINK = LinkTokenInterface(_link);
        blockHashStore = BlockHashStoreInterface(_blockHashStore);
    }

    function registerProvingKey(
        address _oracle, uint256[2] calldata _publicProvingKey
    )
    external
    onlyOwner()
    {
        // Must be unique key
        // I don't think we need the jobID - the serialized PK function as the
        // identifier of the offchain job to run, since its a strict 1-1 map?
        // TODO: Save a mapping of keyHash -> {oracle}
        // Oracle is the address of who gets paid for fulfilling requests.
    }

    function requestRandomWords(
        bytes32 keyHash,  // Corresponds to a particular offchain job which uses that key for the proofs
        uint16  minimumRequestConfirmations,
        uint16  callbackGasLimit,
        uint256 subId,   // A data structure for billing
        uint256 numWords  // Desired number of random words
    )
    external
    returns (bytes32 requestId)
    {
       // TODO:
       // Validate all inputs
       // Sanity check the subscription has enough link? Just
       // accept that gas price fluctuations between request and response could potentially
       // result in request accepted but failed to fulfill.
       /*
       nonce = oracles[_keyHash].nonces[msg.sender] + 1;
       oracles[_keyHash].nonces[msg.sender] = nonce;
       preSeed = keccak256(abi.encode(keyHash, msg.sender, nonce)); // seed with no blockhash
       requestId = keccak256(abi.encode(keyHash, preSeed)); // could this just be the preseed to simplify?
       callbacks[requestId] = Callback({
            callbackContract: msg.sender,
            callbackGasLimit: callbackGasLimit, // For sanity checking
            requestBlockNum: block.number,
            numWords: numWords,
            subId: subId
       });
       */
       // emit log including (indexed keyHash, requestConfs, gasLimit, preSeed)
       // block num/hash is implicit in the log
    }

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
        // TODO: No addresses can be zero, set max number of callers, etc.
        currentSubId++;
        subscriptions[currentSubId] = Subscription({
            owner: msg.sender,
            subId: currentSubId,
            consumers: consumers,
            balance: 0
        });
        emit SubscriptionCreated(subId, msg.sender);
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
        // subscriptions[subId].balance += amount
        // LINK.transferFrom(msg.sender, address(this), amount);
        // TODO: emit some logs
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
