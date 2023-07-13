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
  // to calculate: console.logBytes(type(VRFV2PlusConsumerExample).creationCode);
  bytes constant initializeCode =
    hex"60806040523480156200001157600080fd5b50604051620010b6380380620010b6833981016040819052620000349162000262565b818133806000816200008d5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000c057620000c08162000199565b5050506001600160a01b0382166200010a5760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000084565b6001600160a01b038116620001515760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000084565b600280546001600160a01b039384166001600160a01b03199182161790915560038054928416928216929092179091556004805494909216931692909217909155506200029a565b6001600160a01b038116331415620001f45760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000084565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200025d57600080fd5b919050565b600080604083850312156200027657600080fd5b620002818362000245565b9150620002916020840162000245565b90509250929050565b610e0c80620002aa6000396000f3fe6080604052600436106100965760003560e01c80639c7450ff11610069578063a168fa891161004e578063a168fa8914610169578063f2fde38b146101e8578063f793a70e1461020857600080fd5b80639c7450ff146101295780639eccacf61461014957600080fd5b80631fe543e31461009b57806344ff81ce146100bd57806379ba5097146100dd5780638da5cb5b146100f2575b600080fd5b3480156100a757600080fd5b506100bb6100b6366004610bff565b610228565b005b3480156100c957600080fd5b506100bb6100d8366004610b9d565b610294565b3480156100e957600080fd5b506100bb61031c565b3480156100fe57600080fd5b506000546001600160a01b03165b6040516001600160a01b0390911681526020015b60405180910390f35b61013c610137366004610bcd565b6103da565b6040516101209190610d8c565b34801561015557600080fd5b5060045461010c906001600160a01b031681565b34801561017557600080fd5b506101b9610184366004610bcd565b600560205260009081526040902080546002820154600390920154909160ff8116916101009091046001600160a01b03169084565b604051610120949392919093845291151560208401526001600160a01b03166040830152606082015260800190565b3480156101f457600080fd5b506100bb610203366004610b9d565b610602565b34801561021457600080fd5b506100bb610223366004610cd0565b610616565b6002546001600160a01b03163314610286576002546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b0390911660248201526044015b60405180910390fd5b6102908282610811565b5050565b6003546001600160a01b031633146102ed576003546040517f4ae338ff0000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b03909116602482015260440161027d565b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b6001546001600160a01b031633146103765760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161027d565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000818152600560209081526040808320815160a08101835281548152600182018054845181870281018701909552808552606096959294858401939092919083018282801561044957602002820191906000526020600020905b815481526020019060010190808311610435575b5050509183525050600282015460ff81161515602083015261010090046001600160a01b0316604082015260039091015460609091015280519091506104d15760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161027d565b80604001516105225760405162461bcd60e51b815260206004820152601960248201527f72657175657374206e6f742066756c66696c6c65642079657400000000000000604482015260640161027d565b60608101516001600160a01b031633146105a45760405162461bcd60e51b815260206004820152602360248201527f6f6e6c792063616c6c61626c652062792072657175657374696e67206164647260448201527f6573730000000000000000000000000000000000000000000000000000000000606482015260840161027d565b80608001513410156105f85760405162461bcd60e51b815260206004820152601260248201527f696e73756666696369656e742066756e64730000000000000000000000000000604482015260640161027d565b6020015192915050565b61060a610a11565b61061381610a6d565b50565b600480546040517fefcf1d9400000000000000000000000000000000000000000000000000000000815291820184905267ffffffffffffffff8816602483015261ffff8616604483015263ffffffff80881660648401528516608483015282151560a48301526000916001600160a01b039091169063efcf1d949060c401602060405180830381600087803b1580156106ae57600080fd5b505af11580156106c2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106e69190610be6565b905060006040518060a00160405280838152602001600067ffffffffffffffff81111561071557610715610dd0565b60405190808252806020026020018201604052801561073e578160200160208202803683370190505b5081526000602080830182905233604080850191909152606090930182905285825260058152919020825181558282015180519394508493919261078a92600185019290910190610b24565b50604082015160028201805460608501516001600160a01b0316610100027fffffffffffffffffffffff0000000000000000000000000000000000000000ff931515939093167fffffffffffffffffffffff000000000000000000000000000000000000000000909116179190911790556080909101516003909101555050505050505050565b6000828152600560209081526040808320815160a08101835281548152600182018054845181870281018701909552808552919492938584019390929083018282801561087d57602002820191906000526020600020905b815481526020019060010190808311610869575b5050509183525050600282015460ff81161515602083015261010090046001600160a01b0316604082015260039091015460609091015280519091506109055760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161027d565b6000838152600560209081526040909120835161092a92600190920191850190610b24565b5060008381526005602052604090819020600201805460ff191660011790556004805491517ffd26ba4b0000000000000000000000000000000000000000000000000000000081526001600160a01b039092169163fd26ba4b916109949187910190815260200190565b60206040518083038186803b1580156109ac57600080fd5b505afa1580156109c0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109e49190610d5e565b6bffffffffffffffffffffffff166005600085815260200190815260200160002060030181905550505050565b6000546001600160a01b03163314610a6b5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161027d565b565b6001600160a01b038116331415610ac65760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161027d565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610b5f579160200282015b82811115610b5f578251825591602001919060010190610b44565b50610b6b929150610b6f565b5090565b5b80821115610b6b5760008155600101610b70565b803563ffffffff81168114610b9857600080fd5b919050565b600060208284031215610baf57600080fd5b81356001600160a01b0381168114610bc657600080fd5b9392505050565b600060208284031215610bdf57600080fd5b5035919050565b600060208284031215610bf857600080fd5b5051919050565b60008060408385031215610c1257600080fd5b8235915060208084013567ffffffffffffffff80821115610c3257600080fd5b818601915086601f830112610c4657600080fd5b813581811115610c5857610c58610dd0565b8060051b604051601f19603f83011681018181108582111715610c7d57610c7d610dd0565b604052828152858101935084860182860187018b1015610c9c57600080fd5b600095505b83861015610cbf578035855260019590950194938601938601610ca1565b508096505050505050509250929050565b60008060008060008060c08789031215610ce957600080fd5b863567ffffffffffffffff81168114610d0157600080fd5b9550610d0f60208801610b84565b9450604087013561ffff81168114610d2657600080fd5b9350610d3460608801610b84565b92506080870135915060a08701358015158114610d5057600080fd5b809150509295509295509295565b600060208284031215610d7057600080fd5b81516bffffffffffffffffffffffff81168114610bc657600080fd5b6020808252825182820181905260009190848201906040850190845b81811015610dc457835183529284019291840191600101610da8565b50909695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fdfea164736f6c6343000806000a";

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
      consumerCreate2Address := create2(
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
    uint64 subId = s_testCoordinator.createSubscription();
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
    s_testCoordinator.addConsumer(subId, address(s_testConsumer));

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
    s_testConsumer.requestRandomWords(subId, 1_000_000, 0, 1, vrfKeyHash, true);
    (, bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
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
        -pre-seed 2402601893253931681394804483105813402453761377206973393342199339279344687790 \
        -block-hash 0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8 \
        -block-num 10 \
        -sender 0x0CbB6E5072E30A66041d56e4685Bc89A212a823b \
        -native-payment true
        */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        13784231378140317619606979427300809721548326702202651342016612512034650593641,
        91025342104690364980202625669679932716607712211536344002636554500459902852761
      ],
      c: 97075415945014599812189931594026955245996182195983037019563812588493012075991,
      s: 33951078894534651842088067310556634173280653044215477774617519465697015753547,
      seed: 2402601893253931681394804483105813402453761377206973393342199339279344687790,
      uWitness: 0x8fbd33096a3a6a39dEbb34762109B464055Aacd1,
      cGammaWitness: [
        28177521427931507277502946290773551444121733603074496656773443406265408057706,
        110540610919150129837078745035341910221346493421091351037927631903884632132213
      ],
      sHashWitness: [
        35980458851151413200552438295865484617114607586253285552576984433271595449216,
        66854825481520816474481184449404374539360497566069815291484303257959438867971
      ],
      zInv: 39875168368097267053475197549118206942517384276309995362097006085765194590450
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
    (, fulfilled, , ) = s_testConsumer.s_requests(requestId);
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
    (, uint96 ethBalanceAfter, , ) = s_testCoordinator.getSubscription(subId);
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
    s_testConsumer.requestRandomWords(subId, 1_000_000, 0, 1, vrfKeyHash, false);
    (, bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
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
        -pre-seed 2402601893253931681394804483105813402453761377206973393342199339279344687790 \
        -block-hash 0xce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec \
        -block-num 20 \
        -sender 0x0CbB6E5072E30A66041d56e4685Bc89A212a823b
        */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        11773612677411801428844524078239646079508628031403505959606925234712373334422,
        44863111427282037339794671052408205003980838537866235730096014072326667774193
      ],
      c: 54131218972426713512759017952018335937985880462527408888607190130921548483462,
      s: 47283106285452761081548548829297839121712412917995543039893031137084328867889,
      seed: 2402601893253931681394804483105813402453761377206973393342199339279344687790,
      uWitness: 0x4530F7e95cE4E3cF33f3f47B0D6CcFDBa25e86C4,
      cGammaWitness: [
        73990876814552062779748511263047189814846255422274460822173632701676699651305,
        86000075347242335739938120604364304109734977970148052619511528324905514064237
      ],
      sHashWitness: [
        61324270207334644647482802095972159250022424719940739987021108163314472867928,
        86813184643618962721248123980274874788113602016450268393603767390802904351674
      ],
      zInv: 15799126250054815034680311464655418389537057632771456410474424662471490195458
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
    (, fulfilled, , ) = s_testConsumer.s_requests(requestId);
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
    (uint96 linkBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 260_000, 20_000);
  }
}
