pragma solidity 0.6.6;

import "../VRFConsumerBase.sol";

contract VRFTestnetD20 is VRFConsumerBase {
    using SafeMath_Chainlink for uint;
    uint256[] public d20Results;
    bytes32 internal keyHash;
    uint256 internal fee;

    /**
     * @notice Constructor inherits VRFConsumerBase
     * @dev Ropsten deployment params:
     * @dev   _vrfCoordinator: 0xf720CF1B963e0e7bE9F58fd471EFa67e7bF00cfb
     * @dev   _link:           0x20fE562d797A42Dcb3399062AE9546cd06f63280
     */
    constructor(address _vrfCoordinator, address _link, bytes32 _keyHash)
        VRFConsumerBase(_vrfCoordinator, _link) public
    {
        keyHash = _keyHash;
        fee = 10 ** 18;
    }

    /** 
     * @notice Requests randomness from a user-provided seed
     * @dev This is only an example implementation and not necessarily suitable for mainnet.
     * @dev You must review your implementation details with extreme care.
     */
    function rollDice(uint256 userProvidedSeed) public returns (bytes32 requestId) {
        require(LINK.balanceOf(address(this)) >= fee,
          "Not enough LINK - fill contract with faucet");
        bytes32 _requestId = requestRandomness(keyHash, fee, userProvidedSeed);
        return _requestId;
    }

    /**
     * @notice Callback function used by VRF Coordinator
     * @dev Important! Add a modifier to only allow this function to be called by the VRFCoordinator
     * @dev This is where you do something with randomness!
     * @dev The VRF Coordinator will only send this function verified responses.
     */
    function fulfillRandomness(bytes32 /* requestId */, uint256 randomness) internal override {
        uint256 d20Result = randomness.mod(20).add(1); // Simplified example
        d20Results.push(d20Result);
    }

    /**
     * @notice Convenience function to show the latest roll
     */
    function latestRoll() public view returns (uint256 d20result) {
        return d20Results[d20Results.length - 1];
    }
}
