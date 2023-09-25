pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2Plus.sol";
import {SubscriptionAPI} from "../../../../src/v0.8/dev/vrf/SubscriptionAPI.sol";
import {BlockhashStore} from "../../../../src/v0.8/dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../../../../src/v0.8/dev/vrf/testhelpers/VRFV2PlusConsumerExample.sol";
import {VRFV2PlusClient} from "../../../../src/v0.8/dev/vrf/libraries/VRFV2PlusClient.sol";
import {console} from "forge-std/console.sol";
import {VmSafe} from "forge-std/Vm.sol";
import "@openzeppelin/contracts/utils/math/Math.sol"; // for Math.ceilDiv

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
    hex"60806040523480156200001157600080fd5b5060405162001575380380620015758339810160408190526200003491620001cc565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000103565b5050600280546001600160a01b03199081166001600160a01b0394851617909155600580548216958416959095179094555060038054909316911617905562000204565b6001600160a01b0381163314156200015e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c757600080fd5b919050565b60008060408385031215620001e057600080fd5b620001eb83620001af565b9150620001fb60208401620001af565b90509250929050565b61136180620002146000396000f3fe6080604052600436106101445760003560e01c806380980043116100c0578063b96dbba711610074578063de367c8e11610059578063de367c8e14610372578063eff2701714610392578063f2fde38b146103b257600080fd5b8063b96dbba71461034a578063cf62c8ab1461035257600080fd5b80638ea98117116100a55780638ea981171461029d5780639eccacf6146102bd578063a168fa89146102dd57600080fd5b8063809800431461025f5780638da5cb5b1461027f57600080fd5b806336bfffed11610117578063706da1ca116100fc578063706da1ca146101fc5780637725135b1461021257806379ba50971461024a57600080fd5b806336bfffed146101c65780635d7d53e3146101e657600080fd5b80631d2b2afd146101495780631fe543e31461015357806329e5d831146101735780632fa4e442146101a6575b600080fd5b6101516103d2565b005b34801561015f57600080fd5b5061015161016e36600461104c565b61047b565b34801561017f57600080fd5b5061019361018e3660046110f0565b6104e2565b6040519081526020015b60405180910390f35b3480156101b257600080fd5b506101516101c136600461117d565b6105f8565b3480156101d257600080fd5b506101516101e1366004610f59565b6106e1565b3480156101f257600080fd5b5061019360045481565b34801561020857600080fd5b5061019360065481565b34801561021e57600080fd5b50600354610232906001600160a01b031681565b6040516001600160a01b03909116815260200161019d565b34801561025657600080fd5b506101516107e5565b34801561026b57600080fd5b5061015161027a36600461101a565b600655565b34801561028b57600080fd5b506000546001600160a01b0316610232565b3480156102a957600080fd5b506101516102b8366004610f37565b6108a3565b3480156102c957600080fd5b50600254610232906001600160a01b031681565b3480156102e957600080fd5b506103256102f836600461101a565b6007602052600090815260409020805460019091015460ff82169161010090046001600160a01b03169083565b6040805193151584526001600160a01b0390921660208401529082015260600161019d565b610151610962565b34801561035e57600080fd5b5061015161036d36600461117d565b6109a2565b34801561037e57600080fd5b50600554610232906001600160a01b031681565b34801561039e57600080fd5b506101516103ad366004611112565b6109dc565b3480156103be57600080fd5b506101516103cd366004610f37565b610bcf565b6006546104145760405162461bcd60e51b815260206004820152600b60248201526a1cdd58881b9bdd081cd95d60aa1b60448201526064015b60405180910390fd5b60055460065460405163e8509bff60e01b815260048101919091526001600160a01b039091169063e8509bff9034906024015b6000604051808303818588803b15801561046057600080fd5b505af1158015610474573d6000803e3d6000fd5b5050505050565b6002546001600160a01b031633146104d4576002546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b03909116602482015260440161040b565b6104de8282610be3565b5050565b60008281526007602090815260408083208151608081018352815460ff81161515825261010090046001600160a01b0316818501526001820154818401526002820180548451818702810187019095528085528695929460608601939092919083018282801561057157602002820191906000526020600020905b81548152602001906001019080831161055d575b50505050508152505090508060400151600014156105d15760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161040b565b806060015183815181106105e7576105e761131a565b602002602001015191505092915050565b6006546106355760405162461bcd60e51b815260206004820152600b60248201526a1cdd58881b9bdd081cd95d60aa1b604482015260640161040b565b6003546002546006546040805160208101929092526001600160a01b0393841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161068f939291906111f8565b602060405180830381600087803b1580156106a957600080fd5b505af11580156106bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104de9190610ffd565b6006546107305760405162461bcd60e51b815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161040b565b60005b81518110156104de5760055460065483516001600160a01b039092169163bec4c08c91908590859081106107695761076961131a565b60200260200101516040518363ffffffff1660e01b81526004016107a09291909182526001600160a01b0316602082015260400190565b600060405180830381600087803b1580156107ba57600080fd5b505af11580156107ce573d6000803e3d6000fd5b5050505080806107dd906112f1565b915050610733565b6001546001600160a01b0316331461083f5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161040b565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b031633148015906108c957506002546001600160a01b03163314155b1561093357336108e16000546001600160a01b031690565b6002546040517f061db9c10000000000000000000000000000000000000000000000000000000081526001600160a01b039384166004820152918316602483015291909116604482015260640161040b565b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b61096a610c76565b5060055460065460405163e8509bff60e01b815260048101919091526001600160a01b039091169063e8509bff903490602401610447565b6109aa610c76565b506003546002546006546040805160208101929092526001600160a01b0393841693634000aea0931691859101610662565b60006040518060c0016040528084815260200160065481526020018661ffff1681526020018763ffffffff1681526020018563ffffffff168152602001610a54604051806020016040528086151581525060408051825115156020820152606091016040516020818303038152906040529050919050565b90526002546040517f9b1c385e0000000000000000000000000000000000000000000000000000000081529192506000916001600160a01b0390911690639b1c385e90610aa5908590600401611237565b602060405180830381600087803b158015610abf57600080fd5b505af1158015610ad3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610af79190611033565b604080516080810182526000808252336020808401918252838501868152855184815280830187526060860190815287855260078352959093208451815493517fffffffffffffffffffffff0000000000000000000000000000000000000000009094169015157fffffffffffffffffffffff0000000000000000000000000000000000000000ff16176101006001600160a01b039094169390930292909217825591516001820155925180519495509193849392610bbd926002850192910190610ea7565b50505060049190915550505050505050565b610bd7610d94565b610be081610df0565b50565b6004548214610c345760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161040b565b60008281526007602090815260409091208251610c5992600290920191840190610ea7565b50506000908152600760205260409020805460ff19166001179055565b600060065460001415610d8d57600560009054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b158015610cd357600080fd5b505af1158015610ce7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d0b9190611033565b60068190556005546040517fbec4c08c00000000000000000000000000000000000000000000000000000000815260048101929092523060248301526001600160a01b03169063bec4c08c90604401600060405180830381600087803b158015610d7457600080fd5b505af1158015610d88573d6000803e3d6000fd5b505050505b5060065490565b6000546001600160a01b03163314610dee5760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161040b565b565b6001600160a01b038116331415610e495760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161040b565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610ee2579160200282015b82811115610ee2578251825591602001919060010190610ec7565b50610eee929150610ef2565b5090565b5b80821115610eee5760008155600101610ef3565b80356001600160a01b0381168114610f1e57600080fd5b919050565b803563ffffffff81168114610f1e57600080fd5b600060208284031215610f4957600080fd5b610f5282610f07565b9392505050565b60006020808385031215610f6c57600080fd5b823567ffffffffffffffff811115610f8357600080fd5b8301601f81018513610f9457600080fd5b8035610fa7610fa2826112cd565b61129c565b80828252848201915084840188868560051b8701011115610fc757600080fd5b600094505b83851015610ff157610fdd81610f07565b835260019490940193918501918501610fcc565b50979650505050505050565b60006020828403121561100f57600080fd5b8151610f5281611346565b60006020828403121561102c57600080fd5b5035919050565b60006020828403121561104557600080fd5b5051919050565b6000806040838503121561105f57600080fd5b8235915060208084013567ffffffffffffffff81111561107e57600080fd5b8401601f8101861361108f57600080fd5b803561109d610fa2826112cd565b80828252848201915084840189868560051b87010111156110bd57600080fd5b600094505b838510156110e05780358352600194909401939185019185016110c2565b5080955050505050509250929050565b6000806040838503121561110357600080fd5b50508035926020909101359150565b600080600080600060a0868803121561112a57600080fd5b61113386610f23565b9450602086013561ffff8116811461114a57600080fd5b935061115860408701610f23565b925060608601359150608086013561116f81611346565b809150509295509295909350565b60006020828403121561118f57600080fd5b81356bffffffffffffffffffffffff81168114610f5257600080fd5b6000815180845260005b818110156111d1576020818501810151868301820152016111b5565b818111156111e3576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b03841681526bffffffffffffffffffffffff8316602082015260606040820152600061122e60608301846111ab565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261129460e08401826111ab565b949350505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156112c5576112c5611330565b604052919050565b600067ffffffffffffffff8211156112e7576112e7611330565b5060051b60200190565b600060001982141561131357634e487b7160e01b600052601160045260246000fd5b5060010190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b8015158114610be057600080fdfea164736f6c6343000806000a";

  BlockhashStore s_bhs;
  ExposedVRFCoordinatorV2Plus s_testCoordinator;
  ExposedVRFCoordinatorV2Plus s_testCoordinator_noLink;
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
    // Note: adding contract deployments to this section will require the VRF proofs be regenerated.
    s_testCoordinator = new ExposedVRFCoordinatorV2Plus(address(s_bhs));
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

    s_testCoordinator_noLink = new ExposedVRFCoordinatorV2Plus(address(s_bhs));

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKETHFeed(address(s_linkToken), address(s_linkEthFeed));
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
    uint256 subId = s_testCoordinator.createSubscription();
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
  }

  function testCancelSubWithNoLink() public {
    uint256 subId = s_testCoordinator_noLink.createSubscription();
    s_testCoordinator_noLink.fundSubscriptionWithEth{value: 1000 ether}(subId);

    assertEq(LINK_WHALE.balance, 9000 ether);
    s_testCoordinator_noLink.cancelSubscription(subId, LINK_WHALE);
    assertEq(LINK_WHALE.balance, 10_000 ether);

    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_testCoordinator_noLink.getSubscription(subId);
  }

  function testGetActiveSubscriptionIds() public {
    uint numSubs = 40;
    for (uint i = 0; i < numSubs; i++) {
      s_testCoordinator.createSubscription();
    }
    // get all subscriptions, assert length is correct
    uint256[] memory allSubs = s_testCoordinator.getActiveSubscriptionIds(0, 0);
    assertEq(allSubs.length, s_testCoordinator.getActiveSubscriptionIdsLength());

    // paginate through subscriptions, batching by 10.
    // we should eventually get all the subscriptions this way.
    uint256[][] memory subIds = paginateSubscriptions(s_testCoordinator, 10);
    // check that all subscriptions were returned
    uint actualNumSubs = 0;
    for (uint batchIdx = 0; batchIdx < subIds.length; batchIdx++) {
      for (uint subIdx = 0; subIdx < subIds[batchIdx].length; subIdx++) {
        s_testCoordinator.getSubscription(subIds[batchIdx][subIdx]);
        actualNumSubs++;
      }
    }
    assertEq(actualNumSubs, s_testCoordinator.getActiveSubscriptionIdsLength());

    // cancel a bunch of subscriptions, assert that they are not returned
    uint256[] memory subsToCancel = new uint256[](3);
    for (uint i = 0; i < 3; i++) {
      subsToCancel[i] = subIds[0][i];
    }
    for (uint i = 0; i < subsToCancel.length; i++) {
      s_testCoordinator.cancelSubscription(subsToCancel[i], LINK_WHALE);
    }
    uint256[][] memory newSubIds = paginateSubscriptions(s_testCoordinator, 10);
    // check that all subscriptions were returned
    // and assert that none of the canceled subscriptions are returned
    actualNumSubs = 0;
    for (uint batchIdx = 0; batchIdx < newSubIds.length; batchIdx++) {
      for (uint subIdx = 0; subIdx < newSubIds[batchIdx].length; subIdx++) {
        for (uint i = 0; i < subsToCancel.length; i++) {
          assertFalse(newSubIds[batchIdx][subIdx] == subsToCancel[i]);
        }
        s_testCoordinator.getSubscription(newSubIds[batchIdx][subIdx]);
        actualNumSubs++;
      }
    }
    assertEq(actualNumSubs, s_testCoordinator.getActiveSubscriptionIdsLength());
  }

  function paginateSubscriptions(
    ExposedVRFCoordinatorV2Plus coordinator,
    uint256 batchSize
  ) internal view returns (uint256[][] memory) {
    uint arrIndex = 0;
    uint startIndex = 0;
    uint256 numSubs = coordinator.getActiveSubscriptionIdsLength();
    uint256[][] memory subIds = new uint256[][](Math.ceilDiv(numSubs, batchSize));
    while (startIndex < numSubs) {
      subIds[arrIndex] = coordinator.getActiveSubscriptionIds(startIndex, batchSize);
      startIndex += batchSize;
      arrIndex++;
    }
    return subIds;
  }

  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint256 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bytes extraArgs,
    address indexed sender
  );
  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256 outputSeed,
    uint256 indexed subID,
    uint96 payment,
    bytes extraArgs,
    bool success
  );

  function testRequestAndFulfillRandomWordsNative() public {
    uint32 requestBlock = 10;
    vm.roll(requestBlock);
    s_testConsumer.createSubscriptionAndFund(0);
    uint256 subId = s_testConsumer.s_subId();
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
      subId,
      0, // minConfirmations
      1_000_000, // callbackGasLimit
      1, // numWords
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true})), // nativePayment
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
        -pre-seed 53391429126065232382402681707515137895470547057819816488254124798726362946635 \
        -block-hash 0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8 \
        -block-num 10 \
        -sender 0x90A8820424CC8a819d14cBdE54D12fD3fbFa9bb2 \
        -native-payment true
        */
    VRF.Proof memory proof = VRF.Proof({
        pk: [
                72488970228380509287422715226575535698893157273063074627791787432852706183111,
                62070622898698443831883535403436258712770888294397026493185421712108624767191
        ],
        gamma: [
                74288292845972566688571307909097817917223548093079799841483607824668817625540,
                108391008128802524434610390021145143356403959572332866294516804658902583277098
        ],
        c: 65524113346868605510990811885449481340195152860971795410681661812293260632271,
        s: 47558944104507040739369349257537382457591924764169027480240685677723404176714,
        seed: 31303129040197983804803745611870918648646302283240980411046924527959692018756,
        uWitness: 0xB437Ce51bffAbd2D8a094fBc6d134AA5D47C6A6F,
        cGammaWitness: [
                104159033290129901433595573300142948136581807431838207416723974716049917908679,
                106506593565791628113807425931378042388822730040846687633657770982292729032396
        ],
        sHashWitness: [
                57445942468314843583880742803199308241331439830383315972785792658985538927141,
                87094258463026010430958045471156150903043781327286997614138103522186502986254
        ],
        zInv: 90482407711492828870317874309089656639922057521516036393837774147214270591558
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: 1_000_000,
      numWords: 1,
      sender: address(s_testConsumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    (, uint96 ethBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

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
    (, uint96 ethBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(ethBalanceAfter, ethBalanceBefore - 120_000, 10_000);
  }

  function testRequestAndFulfillRandomWordsLINK() public {
    uint32 requestBlock = 20;
    vm.roll(requestBlock);
    s_linkToken.transfer(address(s_testConsumer), 10 ether);
    s_testConsumer.createSubscriptionAndFund(10 ether);
    uint256 subId = s_testConsumer.s_subId();

    // Apply basic configs to contract.
    setConfig(basicFeeConfig);
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, false, true);
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
      subId,
      0, // minConfirmations
      1_000_000, // callbackGasLimit
      1, // numWords
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false})), // nativePayment, // nativePayment
      address(s_testConsumer) // requester
    );
    console.log("R1 LINK");
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
        -pre-seed 14817911724325909152780695848148728017190840227899344848185245004944693487904 \
        -block-hash 0xce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec \
        -block-num 20 \
        -sender 0x90A8820424CC8a819d14cBdE54D12fD3fbFa9bb2
    */
    VRF.Proof memory proof = VRF.Proof({
        pk: [
                72488970228380509287422715226575535698893157273063074627791787432852706183111,
                62070622898698443831883535403436258712770888294397026493185421712108624767191
        ],
        gamma: [
                67609203698595884603387954646200084425014180457141334500789928806415728881358,
                52802431269913868063810516275482898190042726665036960309811705332899263975860
        ],
        c: 45520383383465869860784703980459196740053055527176047141289825570109658348876,
        s: 72770268682222162636913692985737454500823516307489142458985324885436805983384,
        seed: 46982045317828878826154617205966899914846992048581923814565819044948227487859,
        uWitness: 0xbeC558b10D549e4E43674613691b6Cf4F66334A0,
        cGammaWitness: [
                6895540167864847903598845051906252949531604614797639586165082192345087250147,
                70815824792277496802216786153628107893699162930742809458368279367549362144059
        ],
        sHashWitness: [
                2657252779276092620544256906538786324609983597392202104790583667541976835977,
                81784032561723676164068514426324033600760236460834897534944145289677998921371
        ],
        zInv: 112802723767995327592725752552634369149595717490891621239319956942483369202367
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: 1000000,
      numWords: 1,
      sender: address(s_testConsumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}))
    });
    (uint96 linkBalanceBefore, , , , ) = s_testCoordinator.getSubscription(subId);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);

    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

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
    (uint96 linkBalanceAfter, , , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 280_000, 20_000);
  }
}
