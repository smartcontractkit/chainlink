pragma solidity ^0.8.19;

import {console} from "forge-std/console.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {VRF} from "../VRF.sol";
import {VRFTypes} from "../VRFTypes.sol";
import {BlockhashStore} from "../dev/BlockhashStore.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {ExposedVRFCoordinatorV2_5} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFV2PlusConsumerExample} from "../dev/testhelpers/VRFV2PlusConsumerExample.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import "./BaseTest.t.sol";

contract FixtureVRFCoordinatorV2_5 is BaseTest, VRF {
  address internal SUBSCRIPTION_OWNER = makeAddr("SUBSCRIPTION_OWNER");

  uint64 internal constant GAS_LANE_MAX_GAS = 5000 gwei;
  uint16 internal constant MIN_CONFIRMATIONS = 0;
  uint32 internal constant CALLBACK_GAS_LIMIT = 1_000_000;
  uint32 internal constant NUM_WORDS = 1;

  // VRF KeyV2 generated from a node; not sensitive information.
  // The secret key used to generate this key is: 10.
  bytes internal constant VRF_UNCOMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
  bytes internal constant VRF_COMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
  bytes32 internal constant VRF_KEY_HASH = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

  BlockhashStore internal s_bhs;
  ExposedVRFCoordinatorV2_5 internal s_coordinator;

  // Use multiple consumers because VRFV2PlusConsumerExample cannot have multiple pending requests.
  uint256 internal s_subId;
  VRFV2PlusConsumerExample internal s_consumer;
  VRFV2PlusConsumerExample internal s_consumer1;

  MockLinkToken internal s_linkToken;
  MockV3Aggregator internal s_linkNativeFeed;

  function setUp() public virtual override {
    BaseTest.setUp();
    vm.stopPrank();

    vm.startPrank(OWNER);
    s_bhs = new BlockhashStore();

    // Deploy coordinator.
    s_coordinator = new ExposedVRFCoordinatorV2_5(address(s_bhs));
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Configure the coordinator.
    s_coordinator.setLINKAndLINKNativeFeed(address(s_linkToken), address(s_linkNativeFeed));
    vm.stopPrank();

    // Deploy consumers.
    vm.startPrank(SUBSCRIPTION_OWNER);
    s_consumer = new VRFV2PlusConsumerExample(address(s_coordinator), address(s_linkToken));
    s_consumer1 = new VRFV2PlusConsumerExample(address(s_coordinator), address(s_linkToken));
    vm.stopPrank();
  }

  function _setUpConfig() internal {
    vm.prank(OWNER);
    s_coordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitLink
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );
  }

  function _setUpProvingKey() internal {
    uint256[2] memory uncompressedKeyParts = this._getProvingKeyParts(VRF_UNCOMPRESSED_PUBLIC_KEY);
    vm.prank(OWNER);
    s_coordinator.registerProvingKey(uncompressedKeyParts, GAS_LANE_MAX_GAS);
  }

  function _setUpSubscription() internal {
    vm.startPrank(SUBSCRIPTION_OWNER);
    s_subId = s_coordinator.createSubscription();
    s_coordinator.addConsumer(s_subId, address(s_consumer));
    s_consumer.setSubId(s_subId);
    s_coordinator.addConsumer(s_subId, address(s_consumer1));
    s_consumer1.setSubId(s_subId);
    vm.stopPrank();
  }

  // note: Call this function via this.getProvingKeyParts to be able to pass memory as calldata and
  // index over the byte array.
  function _getProvingKeyParts(bytes calldata uncompressedKey) public pure returns (uint256[2] memory) {
    uint256 keyPart1 = uint256(bytes32(uncompressedKey[0:32]));
    uint256 keyPart2 = uint256(bytes32(uncompressedKey[32:64]));
    return [keyPart1, keyPart2];
  }

  /**
   * Prints the command to generate a proof for a VRF request.
   *
   * This function provides a convenient way to generate the proof off-chain to be copied into the tests.
   *
   * An example of the command looks like this:
   * go run . generate-proof-v2-plus \
   *   -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
   *   -pre-seed 76568185840201037774581758921393822690942290841865097674309745036496166431060 \
   *   -block-hash 0x14 \
   *   -block-num 20 \
   *   -sender 0x2f1c0761d6e4b1e5f01968d6c746f695e5f3e25d \
   *   -native-payment false
   */
  function _printGenerateProofV2PlusCommand(
    address sender,
    uint64 nonce,
    uint256 requestBlock,
    bool nativePayment
  ) internal {
    (, uint256 preSeed) = s_coordinator.computeRequestIdExternal(VRF_KEY_HASH, sender, s_subId, nonce);

    console.log("go run . generate-proof-v2-plus \\");
    console.log(string.concat("  -key-hash ", Strings.toHexString(uint256(VRF_KEY_HASH)), " \\"));
    console.log(string.concat("  -pre-seed ", Strings.toString(preSeed), " \\"));
    console.log(string.concat("  -block-hash ", Strings.toHexString(uint256(blockhash(requestBlock))), " \\"));
    console.log(string.concat("  -block-num ", Strings.toString(requestBlock), " \\"));
    console.log(string.concat("  -sender ", Strings.toHexString(sender), " \\"));
    console.log(string.concat("  -native-payment ", nativePayment ? "true" : "false"));
  }
}
