pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2Plus.sol";
import {BlockhashStore} from "../../../../src/v0.8/dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../../../../src/v0.8/dev/vrf/testhelpers/VRFV2PlusConsumerExample.sol";
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
  // to calculate: console.logBytes(type(VRFV2PlusConsumerExample).creationCode);
  bytes constant initializeCode =
    hex"60806040523480156200001157600080fd5b50604051620016fc380380620016fc833981016040819052620000349162000213565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf816200014a565b5050506001600160a01b038116620001095760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000083565b600280546001600160a01b03199081166001600160a01b0393841617909155600480548216948316949094179093556005805490931691161790556200024b565b6001600160a01b038116331415620001a55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200020e57600080fd5b919050565b600080604083850312156200022757600080fd5b6200023283620001f6565b91506200024260208401620001f6565b90509250929050565b6114a1806200025b6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806379ba509711610097578063a168fa8911610066578063a168fa8914610265578063cf62c8ab146102df578063eff27017146102f2578063f2fde38b1461030557600080fd5b806379ba50971461020c5780637ec8773a146102145780638da5cb5b146102275780639eccacf61461024557600080fd5b806344ff81ce116100d357806344ff81ce146101665780635d7d53e314610179578063706da1ca146101825780637725135b146101c757600080fd5b80631fe543e31461010557806329e5d8311461011a5780632fa4e4421461014057806336bfffed14610153575b600080fd5b610118610113366004611126565b610318565b005b61012d6101283660046111ca565b61039e565b6040519081526020015b60405180910390f35b61011861014e366004611281565b6104db565b610118610161366004611033565b6105af565b610118610174366004611011565b610737565b61012d60065481565b6005546101ae9074010000000000000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610137565b6005546101e79073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610137565b6101186107f1565b610118610222366004611011565b6108ee565b60005473ffffffffffffffffffffffffffffffffffffffff166101e7565b6004546101e79073ffffffffffffffffffffffffffffffffffffffff1681565b6102ad6102733660046110f4565b6007602052600090815260409020805460019091015460ff821691610100900473ffffffffffffffffffffffffffffffffffffffff169083565b60408051931515845273ffffffffffffffffffffffffffffffffffffffff909216602084015290820152606001610137565b6101186102ed366004611281565b6108fa565b6101186103003660046111ec565b610ab0565b610118610313366004611011565b610c8b565b60025473ffffffffffffffffffffffffffffffffffffffff163314610390576002546040517f1cf993f400000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff90911660248201526044015b60405180910390fd5b61039a8282610c9c565b5050565b60008281526007602090815260408083208151608081018352815460ff811615158252610100900473ffffffffffffffffffffffffffffffffffffffff16818501526001820154818401526002820180548451818702810187019095528085528695929460608601939092919083018282801561043a57602002820191906000526020600020905b815481526020019060010190808311610426575b50505050508152505090508060400151600014156104b4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610387565b806060015183815181106104ca576104ca611428565b602002602001015191505092915050565b6005546004546040805174010000000000000000000000000000000000000000840467ffffffffffffffff16602082015273ffffffffffffffffffffffffffffffffffffffff93841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161055d939291906112af565b602060405180830381600087803b15801561057757600080fd5b505af115801561058b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061039a91906110d7565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff1661063a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600d60248201527f7375624944206e6f7420736574000000000000000000000000000000000000006044820152606401610387565b60005b815181101561039a57600454600554835173ffffffffffffffffffffffffffffffffffffffff90921691637341c10c9174010000000000000000000000000000000000000000900467ffffffffffffffff16908590859081106106a2576106a2611428565b60200260200101516040518363ffffffff1660e01b81526004016106f292919067ffffffffffffffff92909216825273ffffffffffffffffffffffffffffffffffffffff16602082015260400190565b600060405180830381600087803b15801561070c57600080fd5b505af1158015610720573d6000803e3d6000fd5b50505050808061072f906113c8565b91505061063d565b60035473ffffffffffffffffffffffffffffffffffffffff1633146107aa576003546040517f4ae338ff00000000000000000000000000000000000000000000000000000000815233600482015273ffffffffffffffffffffffffffffffffffffffff9091166024820152604401610387565b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60015473ffffffffffffffffffffffffffffffffffffffff163314610872576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e6572000000000000000000006044820152606401610387565b60008054337fffffffffffffffffffffffff00000000000000000000000000000000000000008083168217845560018054909116905560405173ffffffffffffffffffffffffffffffffffffffff90921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6108f781610d67565b50565b60055474010000000000000000000000000000000000000000900467ffffffffffffffff166104db5760048054604080517fa21a23e4000000000000000000000000000000000000000000000000000000008152905173ffffffffffffffffffffffffffffffffffffffff9092169263a21a23e49282820192602092908290030181600087803b15801561098d57600080fd5b505af11580156109a1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109c59190611257565b600580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff167401000000000000000000000000000000000000000067ffffffffffffffff93841681029190911791829055600480546040517f7341c10c000000000000000000000000000000000000000000000000000000008152929093049093169281019290925230602483015273ffffffffffffffffffffffffffffffffffffffff1690637341c10c90604401600060405180830381600087803b158015610a8f57600080fd5b505af1158015610aa3573d6000803e3d6000fd5b505050506104db30610d67565b600480546005546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815292830185905274010000000000000000000000000000000000000000900467ffffffffffffffff16602483015261ffff8616604483015263ffffffff80881660648401528516608483015282151560a483015260009173ffffffffffffffffffffffffffffffffffffffff9091169063efcf1d949060c401602060405180830381600087803b158015610b6f57600080fd5b505af1158015610b83573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ba7919061110d565b604080516080810182526000808252336020808401918252838501868152855184815280830187526060860190815287855260078352959093208451815493517fffffffffffffffffffffff0000000000000000000000000000000000000000009094169015157fffffffffffffffffffffff0000000000000000000000000000000000000000ff161761010073ffffffffffffffffffffffffffffffffffffffff9094169390930292909217825591516001820155925180519495509193849392610c7a926002850192910190610f74565b505050600691909155505050505050565b610c93610dfb565b6108f781610e7e565b6006548214610d07576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f7265717565737420494420697320696e636f72726563740000000000000000006044820152606401610387565b60008281526007602090815260409091208251610d2c92600290920191840190610f74565b5050600090815260076020526040902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055565b73ffffffffffffffffffffffffffffffffffffffff8116610db4576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600380547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60005473ffffffffffffffffffffffffffffffffffffffff163314610e7c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e6572000000000000000000006044820152606401610387565b565b73ffffffffffffffffffffffffffffffffffffffff8116331415610efe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c660000000000000000006044820152606401610387565b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610faf579160200282015b82811115610faf578251825591602001919060010190610f94565b50610fbb929150610fbf565b5090565b5b80821115610fbb5760008155600101610fc0565b803573ffffffffffffffffffffffffffffffffffffffff81168114610ff857600080fd5b919050565b803563ffffffff81168114610ff857600080fd5b60006020828403121561102357600080fd5b61102c82610fd4565b9392505050565b6000602080838503121561104657600080fd5b823567ffffffffffffffff81111561105d57600080fd5b8301601f8101851361106e57600080fd5b803561108161107c826113a4565b611355565b80828252848201915084840188868560051b87010111156110a157600080fd5b600094505b838510156110cb576110b781610fd4565b8352600194909401939185019185016110a6565b50979650505050505050565b6000602082840312156110e957600080fd5b815161102c81611486565b60006020828403121561110657600080fd5b5035919050565b60006020828403121561111f57600080fd5b5051919050565b6000806040838503121561113957600080fd5b8235915060208084013567ffffffffffffffff81111561115857600080fd5b8401601f8101861361116957600080fd5b803561117761107c826113a4565b80828252848201915084840189868560051b870101111561119757600080fd5b600094505b838510156111ba57803583526001949094019391850191850161119c565b5080955050505050509250929050565b600080604083850312156111dd57600080fd5b50508035926020909101359150565b600080600080600060a0868803121561120457600080fd5b61120d86610ffd565b9450602086013561ffff8116811461122457600080fd5b935061123260408701610ffd565b925060608601359150608086013561124981611486565b809150509295509295909350565b60006020828403121561126957600080fd5b815167ffffffffffffffff8116811461102c57600080fd5b60006020828403121561129357600080fd5b81356bffffffffffffffffffffffff8116811461102c57600080fd5b73ffffffffffffffffffffffffffffffffffffffff84168152600060206bffffffffffffffffffffffff85168184015260606040840152835180606085015260005b8181101561130d578581018301518582016080015282016112f1565b8181111561131f576000608083870101525b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160800195945050505050565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff8111828210171561139c5761139c611457565b604052919050565b600067ffffffffffffffff8211156113be576113be611457565b5060051b60200190565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611421577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b5060010190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b80151581146108f757600080fdfea164736f6c6343000806000a";

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

    // Deploy link token and link/eth feed.
    s_linkToken = new MockLinkToken();
    s_linkEthFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Use create2 to deploy our consumer, so that its address is always the same
    // and surrounding changes do not alter our generated proofs.
    bytes memory consumerInitCode = bytes.concat(
      initializeCode,
      abi.encode(address(s_testCoordinator), address(s_linkToken))
    );
    bytes32 abiEncodedOwnerAddress = bytes32(uint256(uint160(LINK_WHALE)) << 96);
    address consumerCreate2Address;
    assembly {
      consumerCreate2Address := create2(
        0, // value - left at zero here
        add(0x20, consumerInitCode), // initialization bytecode (excluding first memory slot which contains its length)
        mload(consumerInitCode), // length of initialization bytecode
        abiEncodedOwnerAddress // user-defined nonce to ensure unique SCA addresses
      )
    }
    s_testConsumer = VRFV2PlusConsumerExample(consumerCreate2Address);

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
    (uint16 minConfs, uint32 gasLimit, ) = s_testCoordinator.getRequestConfig();
    assertEq(minConfs, 0);
    assertEq(gasLimit, 2_500_000);

    // Test that setting requestConfirmations above MAX_REQUEST_CONFIRMATIONS reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidRequestConfirmations.selector, 500, 500, 200));
    s_testCoordinator.setConfig(500, 2_500_000, 1, 50_000, 50000000000000000, basicFeeConfig);

    // Test that setting fallbackWeiPerUnitLink to zero reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidLinkWeiPrice.selector, 0));
    s_testCoordinator.setConfig(0, 2_500_000, 1, 50_000, 0, basicFeeConfig);
  }

  function testRegisterProvingKey() public {
    // Should set the proving key successfully.
    registerProvingKey();
    (, , bytes32[] memory keyHashes) = s_testCoordinator.getRequestConfig();
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
    s_testConsumer.createSubscriptionAndFund(0);
    s_testConsumer.setSubOwner(LINK_WHALE);
    uint64 subId = s_testConsumer.s_subId();
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);

    // Apply basic configs to contract.
    setConfig(basicFeeConfig);
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_testConsumer),
      subId,
      2
    );
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
    s_testConsumer.requestRandomWords(1_000_000, 0, 1, vrfKeyHash, true);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);

    // Uncomment these console logs to see info about the request:
    // console.log("requestId: ", requestId);
    // console.log("preSeed: ", preSeed);
    // console.log("sender: ", address(s_testConsumer));

    // Move on to the next block.
    // Store the previous block's blockhash, and assert that it is as expected.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);
    assertEq(hex"c65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 9102713253520531303199546283065661072962076495550082047101034889092927570458 \
        -block-hash 0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8 \
        -block-num 10 \
        -sender 0xC6D70F6D290ab2a0d89f1203E3305217bB3FD933 \
        -native-payment true
        */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        33824286797255484980886608537299637201434597794498922407805944093328629727471,
        114147518798412385315080360733297122817966712944449337584222841517425624885032
      ],
      c: 11075547587725277890262854816997880094854735526128882924254657404119227749769,
      s: 113742430391055970116475408112507658669871570640007783965257250797119324048501,
      seed: 9102713253520531303199546283065661072962076495550082047101034889092927570458,
      uWitness: 0xA92d7cb827035787040123bB766beDf480477a78,
      cGammaWitness: [
        37242615258942891536912292988656406880487546847259153285993965992586927126093,
        47856948526976544855067958379568990534673528016392070249700040120883766329482
      ],
      sHashWitness: [
        55268881536185231240621734254713689764818809268897670060268173609158995897844,
        14577046267729827521044535452459796367608781930536222301486918137726608206672
      ],
      zInv: 73999889333953506047206779112265529958991159107579785797818309967593801607035
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: 1,
      callbackGasLimit: 1_000_000,
      numWords: 1,
      sender: address(s_testConsumer),
      nativePayment: true
    });
    (, uint96 ethBalanceBefore, , ) = s_testCoordinator.getSubscription(subId);
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
    (fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 100_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())
    // baseFeeWei = 1 * (50_000 + 100_000)
    // baseFeeWei = 150_000
    // ...
    // billed_fee = baseFeeWei + flatFeeWei + l1CostWei
    // billed_fee = baseFeeWei + 0 + 0
    // billed_fee = 150_000
    (, uint96 ethBalanceAfter, , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(ethBalanceAfter, ethBalanceBefore - 120_000, 10_000);
  }

  function testRequestAndFulfillRandomWordsLINK() public {
    uint32 requestBlock = 20;
    vm.roll(requestBlock);
    s_linkToken.transfer(address(s_testConsumer), 10 ether);
    s_testConsumer.createSubscriptionAndFund(10 ether);
    s_testConsumer.setSubOwner(LINK_WHALE);
    uint64 subId = s_testConsumer.s_subId();

    // Apply basic configs to contract.
    setConfig(basicFeeConfig);
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_testConsumer),
      subId,
      2
    );
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
    s_testConsumer.requestRandomWords(1_000_000, 0, 1, vrfKeyHash, false);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);

    // Uncomment these console logs to see info about the request:
    // console.log("requestId: ", requestId);
    // console.log("preSeed: ", preSeed);
    // console.log("sender: ", address(s_testConsumer));

    // Move on to the next block.
    // Store the previous block's blockhash, and assert that it is as expected.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);
    assertEq(hex"ce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 9102713253520531303199546283065661072962076495550082047101034889092927570458 \
        -block-hash 0xce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec \
        -block-num 20 \
        -sender 0xC6D70F6D290ab2a0d89f1203E3305217bB3FD933
        */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        101527937419961718691748514242384570851058519722940853333656873402677364345322,
        40405823969488665725416120801220014212583219423037643424644297422703492417143
      ],
      c: 6670470949829865058111333234623512607647317166684225455366189704161205978201,
      s: 67388542104554191236311565725149894176731772079037124447314287682272037675724,
      seed: 9102713253520531303199546283065661072962076495550082047101034889092927570458,
      uWitness: 0xe95Ebdf541B8984D419BB3054BEed734684e8950,
      cGammaWitness: [
        100211286782290242854879787352610819187964690432194561980290515242935634799924,
        64713642127768528549886538972900220986585208659313899382628489261291002952236
      ],
      sHashWitness: [
        79422338473188051481994892725427389440945475829088755413573030917218603965589,
        40424759069180895230260830658150778197086554753298835160599247592177404040681
      ],
      zInv: 8304909023005929766582740023641142630720487655878808142854119294837905363168
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: 1,
      callbackGasLimit: 1000000,
      numWords: 1,
      sender: address(s_testConsumer),
      nativePayment: false
    });
    (uint96 linkBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
    (fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 90_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_eth_ratio)
    // paymentNoFee = (1 * (50_000 + 90_000 + 0)) / .5
    // paymentNoFee = 280_000
    // ...
    // billed_fee = paymentNoFee + fulfillmentFlatFeeLinkPPM
    // billed_fee = baseFeeWei + 0
    // billed_fee = 280_000
    // note: delta is doubled from the native test to account for more variance due to the link/eth ratio
    (uint96 linkBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 280_000, 20_000);
  }
}
