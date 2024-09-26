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
    assertEq(hex"1a192fabce13988b84994d4296e6cdc418d55e2f1d7f942188d4040b94fc57ac", s_bhs.getBlockhash(requestBlock));

    VRFTypes.Proof[] memory proofs = new VRFTypes.Proof[](2);
    VRFTypes.RequestCommitmentV2Plus[] memory rcs = new VRFTypes.RequestCommitmentV2Plus[](2);

    // Proof generated via the generate-proof-v2-plus script command.
    // 1st step: Uncomment the print command below and run the test to print the output.
    // _printGenerateProofV2PlusCommand(address(s_consumer1), 1, requestBlock, false);
    // 2nd step: export the following environment variables to run the generate-proof-v2-plus script.
    // export ETH_URL=https://ethereum-sepolia-rpc.publicnode.com # or any other RPC provider you prefer
    // export ETH_CHAIN_ID=11155111 # or switch to any other chain
    // export ACCOUNT_KEY=<your test EOA private key>
    // 3rd step: copy the output from the 1st step and update the command below, then run the command
    // and copy the command output in the proof section below
    /*
       Run from this folder: chainlink/core/scripts/vrfv2plus/testnet
       go run . generate-proof-v2-plus \
         -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
         -pre-seed 4430852740828987645228960511496023658059009607317025880962658187812299131155 \
         -block-hash 0x1a192fabce13988b84994d4296e6cdc418d55e2f1d7f942188d4040b94fc57ac \
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
        26762213923453052192184693334574145607290366984305044804336172347176490943606,
        70503534560525619072578237689732581746976650376431765635714023643649039207077
      ],
      c: 10992233996918874905152274435276937088064589467016709044984819613170049539489,
      s: 79662863379962724455809192044326025082567113176696761949197261107120333769102,
      seed: 4430852740828987645228960511496023658059009607317025880962658187812299131155,
      uWitness: 0x421A52Fb797d76Fb610aA1a0c020346fC1Ee2DeB,
      cGammaWitness: [
        50748523246052507241857300891945475679319243536065937584940024494820365165901,
        85746856994474260612851047426766648416105284284185975301552792881940939754570
      ],
      sHashWitness: [
        78637275871978664522379716948105702461748200460627087255706483027519919611423,
        82219236913923465822780520561305604064850823877720616893986252854976640396959
      ],
      zInv: 60547558497534848069125896511700272238016171243048151035528198622956754542730
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
    assertEq(hex"731dc163f73d31d8c68f9917ce4ff967753939f70432973c04fd2c2a48148607", s_bhs.getBlockhash(requestBlock));

    // Proof generated via the generate-proof-v2-plus script command.
    // 1st step: Uncomment the print command below and run the test to print the output.
    // _printGenerateProofV2PlusCommand(address(s_consumer1), 1, requestBlock, false);
    // 2nd step: export the following environment variables to run the generate-proof-v2-plus script.
    // export ETH_URL=https://ethereum-sepolia-rpc.publicnode.com # or any other RPC provider you prefer
    // export ETH_CHAIN_ID=11155111 # or switch to any other chain
    // export ACCOUNT_KEY=<your test EOA private key>
    // 3rd step: copy the output from the 1st step and update the command below, then run the command
    // and copy the command output in the proof section below
    /*
       Run from this folder: chainlink/core/scripts/vrfv2plus/testnet
       go run . generate-proof-v2-plus \
         -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
         -pre-seed 14541556911652758131165474365357244907354309169650401973525070879190071151266 \
         -block-hash 0x731dc163f73d31d8c68f9917ce4ff967753939f70432973c04fd2c2a48148607 \
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
        97658842840420719674383370910135023062422561858595941631054490821636116883585,
        44255438468488339528368406358785988551798314198954634050943346751039644360856
      ],
      c: 5233652943248967403606766735502925802264855214922758107203237169366748118852,
      s: 87931642435666855739510477620068257005869145374865238974094299759068218698655,
      seed: 14541556911652758131165474365357244907354309169650401973525070879190071151266,
      uWitness: 0x0A87a9CB71983cE0F2C4bA41D0c1A6Fb1785c46A,
      cGammaWitness: [
        54062743217909816783918413821204010151082432359411822104552882037459289383418,
        67491004534731980264926765871774299056809003077448271411776926359153820235981
      ],
      sHashWitness: [
        7745933951617569731026754652291310837540252155195826133994719499558406927394,
        58405861596456412358325504621101233475720292237067230796670629212111423924259
      ],
      zInv: 44253513765558903217330502897662324213800000485156126961643960636269885275795
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
    emit RandomWordsFulfilled(output.requestId, output.randomness, s_subId, 500000000000143261, true, true, false);
    vm.expectEmit(true, true, false, true, address(s_coordinator));
    emit RandomWordsFulfilled(output1.requestId, output1.randomness, s_subId, 800000000000312358, false, true, false);

    // Fulfill the requests.
    s_batchCoordinator.fulfillRandomWords(proofs, rcs);
  }
}
