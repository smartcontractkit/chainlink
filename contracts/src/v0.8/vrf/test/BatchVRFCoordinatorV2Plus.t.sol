pragma solidity 0.8.19;

import {console} from "forge-std/console.sol";
import {VRF} from "../VRF.sol";
import {VRFTypes} from "../VRFTypes.sol";
import {BatchVRFCoordinatorV2Plus} from "../dev/BatchVRFCoordinatorV2Plus.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {VRFCoordinatorV2_5} from "../dev/VRFCoordinatorV2_5.sol";
import "./BaseTest.t.sol";
import {FixtureVRFCoordinatorV2_5} from "./FixtureVRFCoordinatorV2_5.t.sol";

contract BatchVRFCoordinatorV2PlusTest is FixtureVRFCoordinatorV2_5 {
  BatchVRFCoordinatorV2Plus private s_batchCoordinator;

  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256 outputSeed,
    uint256 indexed subId,
    uint96 payment,
    bool nativePayment,
    bool success,
    bool onlyPremium
  );

  function setUp() public override {
    FixtureVRFCoordinatorV2_5.setUp();

    s_batchCoordinator = new BatchVRFCoordinatorV2Plus(address(s_coordinator));
  }

  function test_fulfillRandomWords() public {
    _setUpConfig();
    _setUpProvingKey();
    _setUpSubscription();

    uint32 requestBlock = 10;
    vm.roll(requestBlock);

    vm.startPrank(SUBSCRIPTION_OWNER);
    vm.deal(SUBSCRIPTION_OWNER, 10 ether);
    s_coordinator.fundSubscriptionWithNative{value: 10 ether}(s_subId);

    // Request random words.
    s_consumer.requestRandomWords(CALLBACK_GAS_LIMIT, MIN_CONFIRMATIONS, NUM_WORDS, VRF_KEY_HASH, true);
    vm.stopPrank();

    // Move on to the next block.
    // Store the previous block's blockhash.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);

    VRFTypes.Proof[] memory proofs = new VRFTypes.Proof[](2);
    VRFTypes.RequestCommitmentV2Plus[] memory rcs = new VRFTypes.RequestCommitmentV2Plus[](2);

    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    // _printGenerateProofV2PlusCommand(address(s_consumer), 1, requestBlock, true);
    /*
       go run . generate-proof-v2-plus \
         -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
         -pre-seed 33855227690351884611579800220581891477580182035146587491531555927634180294480 \
         -block-hash 0x0a \
         -block-num 10 \
         -sender 0xdc90e8ce61c1af8a638b95264037c8e67ee5765c \
         -native-payment true

    */
    proofs[0] = VRFTypes.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        80420391742429647505172101941811820476888293644816377569181566466584288434705,
        24046736031266889997051641830469514057863365715722268340801477580836256044582
      ],
      c: 74775128390693502914275156263410881155583102046081919417827483535122161050585,
      s: 69563235412360165148368009853509434870917653835330501139204071967997764190111,
      seed: 33855227690351884611579800220581891477580182035146587491531555927634180294480,
      uWitness: 0xfB0663eaf48785540dE0FD0F837FD9c09BF4B80A,
      cGammaWitness: [
        53711159452748734758194447734939737695995909567499536035707522847057731697403,
        113650002631484103366420937668971311744887820666944514581352028601506700116835
      ],
      sHashWitness: [
        89656531714223714144489731263049239277719465105516547297952288438117443488525,
        90859682705760125677895017864538514058733199985667976488434404721197234427011
      ],
      zInv: 97275608653505690744303242942631893944856831559408852202478373762878300587548
    });
    rcs[0] = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: s_subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: 1,
      sender: address(s_consumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });

    VRFCoordinatorV2_5.Output memory output = s_coordinator.getRandomnessFromProofExternal(
      abi.decode(abi.encode(proofs[0]), (VRF.Proof)),
      rcs[0]
    );

    requestBlock = 20;
    vm.roll(requestBlock);

    vm.startPrank(SUBSCRIPTION_OWNER);
    s_linkToken.setBalance(address(SUBSCRIPTION_OWNER), 10 ether);
    s_linkToken.transferAndCall(address(s_coordinator), 10 ether, abi.encode(s_subId));

    // Request random words.
    s_consumer1.requestRandomWords(CALLBACK_GAS_LIMIT, MIN_CONFIRMATIONS, NUM_WORDS, VRF_KEY_HASH, false);
    vm.stopPrank();

    // Move on to the next block.
    // Store the previous block's blockhash.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);

    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    // _printGenerateProofV2PlusCommand(address(s_consumer1), 1, requestBlock, false);
    /*
       go run . generate-proof-v2-plus \
         -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
         -pre-seed 76568185840201037774581758921393822690942290841865097674309745036496166431060 \
         -block-hash 0x14 \
         -block-num 20 \
         -sender 0x2f1c0761d6e4b1e5f01968d6c746f695e5f3e25d \
         -native-payment false
    */
    proofs[1] = VRFTypes.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        21323932463597506192387578758854201988004673105893105492473194972397109828006,
        96834737826889397196571646974355352644437196500310392203712129010026003355112
      ],
      c: 8775807990949224376582975115621037245862755412370175152581490650310350359728,
      s: 6805708577951013810918872616271445638109899206333819877111740872779453350091,
      seed: 76568185840201037774581758921393822690942290841865097674309745036496166431060,
      uWitness: 0xE82fF24Fecfbe73d682f38308bE3E039Dfabdf5c,
      cGammaWitness: [
        92810770919624535241476539842820168209710445519252592382122118536598338376923,
        17271305664006119131434661141858450289379246199095231636439133258170648418554
      ],
      sHashWitness: [
        29540023305939374439696120003978246982707698669656874393367212257432197207536,
        93902323936532381028323379401739289810874348405259732508442252936582467730050
      ],
      zInv: 88845170436601946907659333156418518556235340365885668267853966404617557948692
    });
    rcs[1] = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: s_subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: 1,
      sender: address(s_consumer1),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}))
    });

    VRFCoordinatorV2_5.Output memory output1 = s_coordinator.getRandomnessFromProofExternal(
      abi.decode(abi.encode(proofs[1]), (VRF.Proof)),
      rcs[1]
    );

    // The payments are NOT pre-calculated and simply copied from the actual event.
    // We can assert and ignore the payment field but the code will be considerably longer.
    vm.expectEmit(true, true, false, true, address(s_coordinator));
    emit RandomWordsFulfilled(output.requestId, output.randomness, s_subId, 500000000000143283, true, true, false);
    vm.expectEmit(true, true, false, true, address(s_coordinator));
    emit RandomWordsFulfilled(output1.requestId, output1.randomness, s_subId, 800000000000306143, false, true, false);

    // Fulfill the requests.
    s_batchCoordinator.fulfillRandomWords(proofs, rcs);
  }
}
