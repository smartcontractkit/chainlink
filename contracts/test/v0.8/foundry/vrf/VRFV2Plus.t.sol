pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/vrf/VRFCoordinatorV2Plus.sol";
import {BlockhashStore} from "../../../../src/v0.8/dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../../../../src/v0.8/vrf/testhelpers/VRFV2PlusConsumerExample.sol";
import {console} from "forge-std/console.sol";

/*
 * USAGE INSTRUCTIONS:
 * To add new tests/proofs, uncomment the "console.sol" import from foundry, and gather key fields
 * from your VRF request.
 * Then, pass your request info into the generate-proof-v2-plus script command
 * located in /core/scripts/vrfv2/testnet/proofs.go to generate a proof that can be tested on-chain.
**/

contract VRFV2Plus is BaseTest {
    address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;

    // Bytecode for a VRFV2PlusConsumerExample contract.
    bytes constant initializeCode = type(VRFV2PlusConsumerExample).creationCode;

    BlockhashStore s_bhs;
    ExposedVRFCoordinatorV2Plus s_testCoordinator;
    VRFV2PlusConsumerExample s_testConsumer;
    MockLinkToken s_linkToken;
    MockV3Aggregator s_linkEthFeed;

    VRFCoordinatorV2Plus.FeeConfig basicFeeConfig =
        VRFCoordinatorV2Plus.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeEthPPM: 0});

    // VRF KeyV2 generated from a node; not sensitive information.
    // The secret key used to generate this key is: 10.
    bytes vrfUncompressedPublicKey =
        hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
    bytes vrfCompressedPublicKey = hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
    bytes32 vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

    function setUp() public override {
        BaseTest.setUp();

        // Fund our users.
        vm.roll(1);
        vm.deal(LINK_WHALE, 10_000 ether);
        changePrank(LINK_WHALE);

        // Instantiate BHS.
        s_bhs = new BlockhashStore();

        // Deploy coordinator and consumer.
        s_testCoordinator = new ExposedVRFCoordinatorV2Plus(address(s_bhs));

        // Use create2 to deploy our consumer, so that its address is always the same
        // and surrounding changes do not alter our generated proofs.
        bytes memory consumerInitCode = bytes.concat(initializeCode, abi.encode(address(s_testCoordinator), LINK_WHALE));
        bytes32 abiEncodedOwnerAddress = bytes32(uint256(uint160(LINK_WHALE)) << 96);
        address consumerCreate2Address;
        assembly {
            consumerCreate2Address :=
                create2(
                    0, // value - left at zero here
                    add(0x20, consumerInitCode), // initialization bytecode (excluding first memory slot which contains its length)
                    mload(consumerInitCode), // length of initialization bytecode
                    abiEncodedOwnerAddress // user-defined nonce to ensure unique SCA addresses
                )
        }
        s_testConsumer = VRFV2PlusConsumerExample(consumerCreate2Address);

        // Deploy link token and link/eth feed.
        s_linkToken = new MockLinkToken();
        s_linkEthFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

        // Configure the coordinator.
        s_testCoordinator.setLINK(address(s_linkToken));
        s_testCoordinator.setLinkEthFeed(address(s_linkEthFeed));
    }

    function setConfig(VRFCoordinatorV2Plus.FeeConfig memory feeConfig) internal {
        s_testCoordinator.setConfig(
            0, // minRequestConfirmations
            2_500_000, // maxGasLimit
            1, // stalenessSeconds
            50_000, // gasAfterPaymentCalculation
            50000000000000000, // fallbackWeiPerUnitLink
            feeConfig
        );
    }

    function testSetConfig() public {
        // Should setConfig successfully.
        setConfig(basicFeeConfig);
        (uint16 minConfs, uint32 gasLimit,) = s_testCoordinator.getRequestConfig();
        assertEq(minConfs, 0);
        assertEq(gasLimit, 2_500_000);

        // Test that setting requestConfirmations above MAX_REQUEST_CONFIRMATIONS reverts.
        vm.expectRevert(
            abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidRequestConfirmations.selector, 500, 500, 200)
        );
        s_testCoordinator.setConfig(500, 2_500_000, 1, 50_000, 50000000000000000, basicFeeConfig);

        // Test that setting fallbackWeiPerUnitLink to zero reverts.
        vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidLinkWeiPrice.selector, 0));
        s_testCoordinator.setConfig(0, 2_500_000, 1, 50_000, 0, basicFeeConfig);
    }

    function testRegisterProvingKey() public {
        // Should set the proving key successfully.
        registerProvingKey();
        (,, bytes32[] memory keyHashes) = s_testCoordinator.getRequestConfig();
        assertEq(keyHashes[0], vrfKeyHash);

        // Should revert when already registered.
        uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
        vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.ProvingKeyAlreadyRegistered.selector, vrfKeyHash));
        s_testCoordinator.registerProvingKey(LINK_WHALE, uncompressedKeyParts);
    }

    function registerProvingKey() public {
        uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
        s_testCoordinator.registerProvingKey(LINK_WHALE, uncompressedKeyParts);
    }

    // note: Call this function via this.getProvingKeyParts to be able to pass memory as calldata and
    // index over the byte array.
    function getProvingKeyParts(bytes calldata uncompressedKey) public pure returns (uint256[2] memory) {
        uint256 keyPart1 = uint256(bytes32(uncompressedKey[0:32]));
        uint256 keyPart2 = uint256(bytes32(uncompressedKey[32:64]));
        return [keyPart1, keyPart2];
    }

    function testCreateSubscription() public {
        uint64 subId = s_testCoordinator.createSubscription();
        s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
    }

    event RandomWordsRequested(
        bytes32 indexed keyHash,
        uint256 requestId,
        uint256 preSeed,
        uint64 indexed subId,
        uint16 minimumRequestConfirmations,
        uint32 callbackGasLimit,
        uint32 numWords,
        bool nativePayment,
        address indexed sender
    );

    function testRequestAndFulfillRandomWordsNative() public {
        uint32 requestBlock = 10;
        vm.roll(requestBlock);
        uint64 subId = s_testCoordinator.createSubscription();
        s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
        s_testCoordinator.addConsumer(subId, address(s_testConsumer));

        // Apply basic configs to contract.
        setConfig(basicFeeConfig);
        registerProvingKey();

        // Request random words.
        vm.expectEmit(true, true, true, true);
        (uint256 requestId, uint256 preSeed) =
            s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(s_testConsumer), subId, 2);
        emit RandomWordsRequested(
            vrfKeyHash,
            requestId,
            preSeed,
            1, // subId
            0, // minConfirmations
            1_000_000, // callbackGasLimit
            1, // numWords
            true, // nativePayment
            address(s_testConsumer) // requester
        );
        s_testConsumer.requestRandomWords(subId, 1_000_000, 0, 1, vrfKeyHash, true);
        (, bool fulfilled,,) = s_testConsumer.s_requests(requestId);
        assertEq(fulfilled, false);

        // Uncomment these console logs to see info about the request:
        // console.log("requestId: ", requestId);
        // console.log("preSeed: ", preSeed);
        // console.log("sender: ", address(s_testConsumer));

        // Move on to the next block.
        // Store the previous block's blockhash, and assert that it is as expected.
        vm.roll(requestBlock + 1);
        s_bhs.store(requestBlock);
        assertEq(
            hex"c65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8", s_bhs.getBlockhash(requestBlock)
        );

        // Fulfill the request.
        // Proof generated via the generate-proof-v2-plus script command. Example usage:
        /* 
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 76695080434118067850784566306514319596947576573091324284020286914149628112655 \
        -block-hash 0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8 \
        -block-num 10 \
        -sender 0xDcd88e478961f13eE0f3D4d7603a1ce89A7deE5C \
        -native-payment true
        */
        VRF.Proof memory proof = VRF.Proof({
            pk: [
                72488970228380509287422715226575535698893157273063074627791787432852706183111,
                62070622898698443831883535403436258712770888294397026493185421712108624767191
            ],
            gamma: [
                99855868970690593100315840615826042784813834664657591642539612483566803407811,
                37686020045276405818480671014897961508737555042389532071178702882521706486919
            ],
            c: 83423337586161416418818868628889498023418194267454390668593887134109299119865,
            s: 12886344513267760439642474429537414504789505873912508066060416512946909154608,
            seed: 76695080434118067850784566306514319596947576573091324284020286914149628112655,
            uWitness: 0x48801970bdf833c0a0E3D2067990F1A5A7f15360,
            cGammaWitness: [
                27280727028699827665170906060616498676379003586802733167641553392283452188931,
                47224949111743613594829992135584623925361956364059179352835430742298341360882
            ],
            sHashWitness: [
                48409693222960225380327985320351632554835811410218908004051700701461559728114,
                23793109531067333764470684214158149909437606745416349416689447109216869181602
            ],
            zInv: 60401217999121404690117102209537470385156607700802185866816657329608429337918
        });
        VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
            blockNum: requestBlock,
            subId: 1,
            callbackGasLimit: 1_000_000,
            numWords: 1,
            sender: address(s_testConsumer),
            nativePayment: true
        });
        (, uint96 ethBalanceBefore,,) = s_testCoordinator.getSubscription(subId);
        s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
        (, fulfilled,,) = s_testConsumer.s_requests(requestId);
        assertEq(fulfilled, true);

        // The cost of fulfillRandomWords is approximately 75_000 gas.
        // gasAfterPaymentCalculation is 50_000.
        //
        // The cost of the VRF fulfillment charged to the user is:
        // baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())
        // baseFeeWei = 1 * (50_000 + 75_000)
        // baseFeeWei = 125_000
        // ...
        // billed_fee = baseFeeWei + flatFeeWei + l1CostWei
        // billed_fee = baseFeeWei + 0 + 0
        // billed_fee = 125_000
        (, uint96 ethBalanceAfter,,) = s_testCoordinator.getSubscription(subId);
        assertApproxEqAbs(ethBalanceAfter, ethBalanceBefore - 130_000, 10_000);
    }

    function testRequestAndFulfillRandomWordsLINK() public {
        uint32 requestBlock = 20;
        vm.roll(requestBlock);
        uint64 subId = s_testCoordinator.createSubscription();
        s_linkToken.transferAndCall(address(s_testCoordinator), 10 ether, abi.encode(subId));
        s_testCoordinator.addConsumer(subId, address(s_testConsumer));

        // Apply basic configs to contract.
        setConfig(basicFeeConfig);
        registerProvingKey();

        // Request random words.
        vm.expectEmit(true, true, true, true);
        (uint256 requestId, uint256 preSeed) =
            s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(s_testConsumer), subId, 2);
        emit RandomWordsRequested(
            vrfKeyHash,
            requestId,
            preSeed,
            1, // subId
            0, // minConfirmations
            1_000_000, // callbackGasLimit
            1, // numWords
            false, // nativePayment
            address(s_testConsumer) // requester
        );
        s_testConsumer.requestRandomWords(subId, 1_000_000, 0, 1, vrfKeyHash, false);
        (, bool fulfilled,,) = s_testConsumer.s_requests(requestId);
        assertEq(fulfilled, false);

        // Uncomment these console logs to see info about the request:
        // console.log("requestId: ", requestId);
        // console.log("preSeed: ", preSeed);
        // console.log("sender: ", address(s_testConsumer));

        // Move on to the next block.
        // Store the previous block's blockhash, and assert that it is as expected.
        vm.roll(requestBlock + 1);
        s_bhs.store(requestBlock);
        assertEq(
            hex"ce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec", s_bhs.getBlockhash(requestBlock)
        );

        // Fulfill the request.
        // Proof generated via the generate-proof-v2-plus script command. Example usage:
        /* 
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 76695080434118067850784566306514319596947576573091324284020286914149628112655 \
        -block-hash 0xce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec \
        -block-num 20 \
        -sender 0xDcd88e478961f13eE0f3D4d7603a1ce89A7deE5C
        */
        VRF.Proof memory proof = VRF.Proof({
            pk: [
                72488970228380509287422715226575535698893157273063074627791787432852706183111,
                62070622898698443831883535403436258712770888294397026493185421712108624767191
            ],
            gamma: [
                75789895770583553412289127485749412742338088061501341383147055533298475722968,
                85528689556526601823253976155521693953338623099543915074527082453568228650360
            ],
            c: 65704621016703078595306213270746247497212253596291446686711361795088275912814,
            s: 92835376495010816653664761462661719802013857123275599192835906932901543062722,
            seed: 76695080434118067850784566306514319596947576573091324284020286914149628112655,
            uWitness: 0xa991f7c34aA0ADb3326B64b634B12Caac8e1C73f,
            cGammaWitness: [
                83589836960561653027214426880941250519885073687366914143245327555949149929672,
                54008725737004974862566361542939795327487421510439605844675809913986013096957
            ],
            sHashWitness: [
                87349481569327046689512192412234143651805181404367424067902567052614493103926,
                87704428761539194493686286322791549547667816990048295974554312592241468947702
            ],
            zInv: 37964763331511074394165802914121445169662177828265364200490386480929331087409
        });
        VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
            blockNum: requestBlock,
            subId: 1,
            callbackGasLimit: 1000000,
            numWords: 1,
            sender: address(s_testConsumer),
            nativePayment: false
        });
        (uint96 linkBalanceBefore,,,) = s_testCoordinator.getSubscription(subId);
        s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
        (, fulfilled,,) = s_testConsumer.s_requests(requestId);
        assertEq(fulfilled, true);

        // The cost of fulfillRandomWords is approximately 75_000 gas.
        // gasAfterPaymentCalculation is 50_000.
        //
        // The cost of the VRF fulfillment charged to the user is:
        // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_eth_ratio)
        // paymentNoFee = (1 * (50_000 + 80_000 + 0)) / .5
        // paymentNoFee = 260_000
        // ...
        // billed_fee = paymentNoFee + fulfillmentFlatFeeLinkPPM
        // billed_fee = baseFeeWei + 0
        // billed_fee = 260_000
        // note: delta is doubled from the native test to account for more variance due to the link/eth ratio
        (uint96 linkBalanceAfter,,,) = s_testCoordinator.getSubscription(subId);
        assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 260_000, 20_000);
    }
}
