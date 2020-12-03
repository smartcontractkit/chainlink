pragma solidity 0.6.6;

import "../VRFConsumerBase.sol";
import "../Owned.sol";

/**
 * @notice A Chainlink VRF consumer which uses randomness to mimic the rolling
 * of a 20 sided die
 */
contract VRFD20 is VRFConsumerBase, Owned {

    bytes32 private s_keyHash;
    uint256 private s_fee;
    mapping(address => bytes32) private s_rollers;
    mapping(bytes32 => uint256) private s_results;
    mapping(uint256 => string) private s_houses;

    event DiceRolled(bytes32 indexed requestId, address indexed roller);
    event DiceLanded(bytes32 indexed requestId, uint256 indexed result);

    /**
     * @notice Constructor inherits VRFConsumerBase
     *
     * @dev NETWORK: KOVAN
     * @dev   Chainlink VRF Coordinator address: 0xdD3782915140c8f3b190B5D67eAc6dc5760C46E9
     * @dev   LINK token address:                0xa36085F69e2889c224210F603D836748e7dC0088
     * @dev   Key Hash:   0x6c3699283bda56ad74f6b855546325b68d482e983852a7a82979cc4807b641f4
     * @dev   Fee:        0.1 LINK (100000000000000000)
     *
     * @param vrfCoordinator address of the VRF Coordinator
     * @param link address of the LINK token
     * @param keyHash bytes32 representing the hash of the VRF job
     * @param fee uint256 fee to pay the VRF oracle
     */
    constructor(address vrfCoordinator, address link, bytes32 keyHash, uint256 fee)
        public
        VRFConsumerBase(vrfCoordinator, link)
    {
        s_keyHash = keyHash;
        s_fee = fee;
        s_houses[1] = "Targaryen";
        s_houses[2] = "Lannister";
        s_houses[3] = "Stark";
        s_houses[4] = "Tyrell";
        s_houses[5] = "Baratheon";
        s_houses[6] = "Martell";
        s_houses[7] = "Tully";
        s_houses[8] = "Bolton";
        s_houses[9] = "Greyjoy";
        s_houses[10] = "Arryn";
        s_houses[11] = "Frey";
        s_houses[12] = "Mormont";
        s_houses[13] = "Tarley";
        s_houses[14] = "Dayne";
        s_houses[15] = "Umber";
        s_houses[16] = "Valeryon";
        s_houses[17] = "Manderly";
        s_houses[18] = "Clegane";
        s_houses[19] = "Glover";
        s_houses[20] = "Karstark";
    }

    /**
     * @notice Requests randomness from a user-provided seed
     * @dev This is only an example implementation and not necessarily suitable for mainnet.
     * @dev You must review your implementation details with extreme care.
     *
     * @param userProvidedSeed uint256 unpredictable seed
     * @param roller address of the roller
     */
    function rollDice(uint256 userProvidedSeed, address roller) public onlyOwner returns (bytes32 requestId) {
        require(LINK.balanceOf(address(this)) >= s_fee, "Not enough LINK to pay fee");
        require(s_rollers[roller] == bytes32(0), "Already rolled");
        requestId = requestRandomness(s_keyHash, s_fee, userProvidedSeed);
        s_rollers[roller] = requestId;
        emit DiceRolled(requestId, roller);
    }

    /**
     * @notice Get the house assigned to the player once the address has rolled
     * @param player address
     * @return house as a string
     */
    function house(address player) public view returns (string memory) {
        require(s_rollers[player] != bytes32(0), "Dice not rolled");
        require(s_results[s_rollers[player]] != 0, "Roll in progress");
        return s_houses[s_results[s_rollers[player]]];
    }

    /**
     * @notice Withdraw LINK from this contract.
     * @param to the address to withdraw LINK to
     * @param value the amount of LINK to withdraw
     */
    function withdrawLINK(address to, uint256 value) public onlyOwner {
        require(LINK.transfer(to, value), "Not enough LINK");
    }

    /**
     * @notice Set the key hash for the oracle
     *
     * @param keyHash bytes32
     */
    function setKeyHash(bytes32 keyHash) public onlyOwner {
        s_keyHash = keyHash;
    }

    /**
     * @notice Get the current key hash
     *
     * @return bytes32
     */
    function keyHash() public view returns (bytes32) {
        return s_keyHash;
    }

    /**
     * @notice Set the oracle fee for requesting randomness
     *
     * @param fee uint256
     */
    function setFee(uint256 fee) public onlyOwner {
        s_fee = fee;
    }

    /**
     * @notice Get the current fee
     *
     * @return uint256
     */
    function fee() public view returns (uint256) {
        return s_fee;
    }

    /**
     * @notice Callback function used by VRF Coordinator to return the random number
     * to this contract.
     * @dev This is where you do something with randomness!
     * @dev The VRF Coordinator will only send this function verified responses.
     *
     * @param requestId bytes32
     * @param randomness The random result returned by the oracle
     */
    function fulfillRandomness(bytes32 requestId, uint256 randomness) internal override {
        uint256 result = randomness.mod(20).add(1);
        s_results[requestId] = result;
        emit DiceLanded(requestId, result);
    }
}
