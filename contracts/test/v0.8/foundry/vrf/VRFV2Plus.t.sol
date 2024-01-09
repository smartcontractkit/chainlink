pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5} from "../../../../src/v0.8/vrf/dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFCoordinatorV2_5} from "../../../../src/v0.8/vrf/dev/VRFCoordinatorV2_5.sol";
import {SubscriptionAPI} from "../../../../src/v0.8/vrf/dev/SubscriptionAPI.sol";
import {BlockhashStore} from "../../../../src/v0.8/vrf/dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../../../../src/v0.8/vrf/dev/testhelpers/VRFV2PlusConsumerExample.sol";
import {VRFV2PlusClient} from "../../../../src/v0.8/vrf/dev/libraries/VRFV2PlusClient.sol";
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
    hex"60806040523480156200001157600080fd5b5060405162001377380380620013778339810160408190526200003491620001cc565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf8162000103565b5050600280546001600160a01b03199081166001600160a01b0394851617909155600580548216958416959095179094555060038054909316911617905562000204565b6001600160a01b0381163314156200015e5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b0381168114620001c757600080fd5b919050565b60008060408385031215620001e057600080fd5b620001eb83620001af565b9150620001fb60208401620001af565b90509250929050565b61116380620002146000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80638098004311610097578063cf62c8ab11610066578063cf62c8ab14610242578063de367c8e14610255578063eff2701714610268578063f2fde38b1461027b57600080fd5b806380980043146101ab5780638da5cb5b146101be5780638ea98117146101cf578063a168fa89146101e257600080fd5b80635d7d53e3116100d35780635d7d53e314610166578063706da1ca1461016f5780637725135b1461017857806379ba5097146101a357600080fd5b80631fe543e31461010557806329e5d8311461011a5780632fa4e4421461014057806336bfffed14610153575b600080fd5b610118610113366004610e4e565b61028e565b005b61012d610128366004610ef2565b6102fa565b6040519081526020015b60405180910390f35b61011861014e366004610f7f565b610410565b610118610161366004610d5b565b6104bc565b61012d60045481565b61012d60065481565b60035461018b906001600160a01b031681565b6040516001600160a01b039091168152602001610137565b6101186105c0565b6101186101b9366004610e1c565b600655565b6000546001600160a01b031661018b565b6101186101dd366004610d39565b61067e565b61021d6101f0366004610e1c565b6007602052600090815260409020805460019091015460ff82169161010090046001600160a01b03169083565b6040805193151584526001600160a01b03909216602084015290820152606001610137565b610118610250366004610f7f565b61073d565b60055461018b906001600160a01b031681565b610118610276366004610f14565b610880565b610118610289366004610d39565b610a51565b6002546001600160a01b031633146102ec576002546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b0390911660248201526044015b60405180910390fd5b6102f68282610a65565b5050565b60008281526007602090815260408083208151608081018352815460ff81161515825261010090046001600160a01b0316818501526001820154818401526002820180548451818702810187019095528085528695929460608601939092919083018282801561038957602002820191906000526020600020905b815481526020019060010190808311610375575b50505050508152505090508060400151600014156103e95760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016102e3565b806060015183815181106103ff576103ff61111c565b602002602001015191505092915050565b6003546002546006546040805160208101929092526001600160a01b0393841693634000aea09316918591015b6040516020818303038152906040526040518463ffffffff1660e01b815260040161046a93929190610ffa565b602060405180830381600087803b15801561048457600080fd5b505af1158015610498573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102f69190610dff565b60065461050b5760405162461bcd60e51b815260206004820152600d60248201527f7375624944206e6f74207365740000000000000000000000000000000000000060448201526064016102e3565b60005b81518110156102f65760055460065483516001600160a01b039092169163bec4c08c91908590859081106105445761054461111c565b60200260200101516040518363ffffffff1660e01b815260040161057b9291909182526001600160a01b0316602082015260400190565b600060405180830381600087803b15801561059557600080fd5b505af11580156105a9573d6000803e3d6000fd5b5050505080806105b8906110f3565b91505061050e565b6001546001600160a01b0316331461061a5760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e65720000000000000000000060448201526064016102e3565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b6000546001600160a01b031633148015906106a457506002546001600160a01b03163314155b1561070e57336106bc6000546001600160a01b031690565b6002546040517f061db9c10000000000000000000000000000000000000000000000000000000081526001600160a01b03938416600482015291831660248301529190911660448201526064016102e3565b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b60065461041057600560009054906101000a90046001600160a01b03166001600160a01b031663a21a23e46040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561079457600080fd5b505af11580156107a8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107cc9190610e35565b60068190556005546040517fbec4c08c00000000000000000000000000000000000000000000000000000000815260048101929092523060248301526001600160a01b03169063bec4c08c90604401600060405180830381600087803b15801561083557600080fd5b505af1158015610849573d6000803e3d6000fd5b505050506003546002546006546040516001600160a01b0393841693634000aea0931691859161043d919060200190815260200190565b60006040518060c0016040528084815260200160065481526020018661ffff1681526020018763ffffffff1681526020018563ffffffff1681526020016108d66040518060200160405280861515815250610af8565b90526002546040517f9b1c385e0000000000000000000000000000000000000000000000000000000081529192506000916001600160a01b0390911690639b1c385e90610927908590600401611039565b602060405180830381600087803b15801561094157600080fd5b505af1158015610955573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109799190610e35565b604080516080810182526000808252336020808401918252838501868152855184815280830187526060860190815287855260078352959093208451815493517fffffffffffffffffffffff0000000000000000000000000000000000000000009094169015157fffffffffffffffffffffff0000000000000000000000000000000000000000ff16176101006001600160a01b039094169390930292909217825591516001820155925180519495509193849392610a3f926002850192910190610ca9565b50505060049190915550505050505050565b610a59610b96565b610a6281610bf2565b50565b6004548214610ab65760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f727265637400000000000000000060448201526064016102e3565b60008281526007602090815260409091208251610adb92600290920191840190610ca9565b50506000908152600760205260409020805460ff19166001179055565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610b3191511515815260200190565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b6000546001600160a01b03163314610bf05760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e65720000000000000000000060448201526064016102e3565b565b6001600160a01b038116331415610c4b5760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c6600000000000000000060448201526064016102e3565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610ce4579160200282015b82811115610ce4578251825591602001919060010190610cc9565b50610cf0929150610cf4565b5090565b5b80821115610cf05760008155600101610cf5565b80356001600160a01b0381168114610d2057600080fd5b919050565b803563ffffffff81168114610d2057600080fd5b600060208284031215610d4b57600080fd5b610d5482610d09565b9392505050565b60006020808385031215610d6e57600080fd5b823567ffffffffffffffff811115610d8557600080fd5b8301601f81018513610d9657600080fd5b8035610da9610da4826110cf565b61109e565b80828252848201915084840188868560051b8701011115610dc957600080fd5b600094505b83851015610df357610ddf81610d09565b835260019490940193918501918501610dce565b50979650505050505050565b600060208284031215610e1157600080fd5b8151610d5481611148565b600060208284031215610e2e57600080fd5b5035919050565b600060208284031215610e4757600080fd5b5051919050565b60008060408385031215610e6157600080fd5b8235915060208084013567ffffffffffffffff811115610e8057600080fd5b8401601f81018613610e9157600080fd5b8035610e9f610da4826110cf565b80828252848201915084840189868560051b8701011115610ebf57600080fd5b600094505b83851015610ee2578035835260019490940193918501918501610ec4565b5080955050505050509250929050565b60008060408385031215610f0557600080fd5b50508035926020909101359150565b600080600080600060a08688031215610f2c57600080fd5b610f3586610d25565b9450602086013561ffff81168114610f4c57600080fd5b9350610f5a60408701610d25565b9250606086013591506080860135610f7181611148565b809150509295509295909350565b600060208284031215610f9157600080fd5b81356bffffffffffffffffffffffff81168114610d5457600080fd5b6000815180845260005b81811015610fd357602081850181015186830182015201610fb7565b81811115610fe5576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b03841681526bffffffffffffffffffffffff831660208201526060604082015260006110306060830184610fad565b95945050505050565b60208152815160208201526020820151604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c08084015261109660e0840182610fad565b949350505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156110c7576110c7611132565b604052919050565b600067ffffffffffffffff8211156110e9576110e9611132565b5060051b60200190565b600060001982141561111557634e487b7160e01b600052601160045260246000fd5b5060010190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b8015158114610a6257600080fdfea164736f6c6343000806000a";

  BlockhashStore s_bhs;
  ExposedVRFCoordinatorV2_5 s_testCoordinator;
  ExposedVRFCoordinatorV2_5 s_testCoordinator_noLink;
  VRFV2PlusConsumerExample s_testConsumer;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkNativeFeed;

  VRFCoordinatorV2_5.FeeConfig basicFeeConfig =
    VRFCoordinatorV2_5.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeNativePPM: 0});

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
    s_testCoordinator = new ExposedVRFCoordinatorV2_5(address(s_bhs));
    s_linkToken = new MockLinkToken();
    s_linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

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

    s_testCoordinator_noLink = new ExposedVRFCoordinatorV2_5(address(s_bhs));

    // Configure the coordinator.
    s_testCoordinator.setLINKAndLINKNativeFeed(address(s_linkToken), address(s_linkNativeFeed));
  }

  function setConfig(VRFCoordinatorV2_5.FeeConfig memory feeConfig) internal {
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
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidRequestConfirmations.selector, 500, 500, 200));
    s_testCoordinator.setConfig(500, 2_500_000, 1, 50_000, 50000000000000000, basicFeeConfig);

    // Test that setting fallbackWeiPerUnitLink to zero reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidLinkWeiPrice.selector, 0));
    s_testCoordinator.setConfig(0, 2_500_000, 1, 50_000, 0, basicFeeConfig);
  }

  function testRegisterProvingKey() public {
    // Should set the proving key successfully.
    registerProvingKey();
    (, , bytes32[] memory keyHashes) = s_testCoordinator.getRequestConfig();
    assertEq(keyHashes[0], vrfKeyHash);

    // Should revert when already registered.
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.ProvingKeyAlreadyRegistered.selector, vrfKeyHash));
    s_testCoordinator.registerProvingKey(uncompressedKeyParts);
  }

  function registerProvingKey() public {
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    s_testCoordinator.registerProvingKey(uncompressedKeyParts);
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
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(subId);
  }

  function testCancelSubWithNoLink() public {
    uint256 subId = s_testCoordinator_noLink.createSubscription();
    s_testCoordinator_noLink.fundSubscriptionWithNative{value: 1000 ether}(subId);

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
    ExposedVRFCoordinatorV2_5 coordinator,
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
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(subId);

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
    assertEq(hex"000000000000000000000000000000000000000000000000000000000000000a", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
       go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 93724884573574303181157854277074121673523280784530506403108144933983063023487 \
        -block-hash 0x000000000000000000000000000000000000000000000000000000000000000a \
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
        51111463251706978184511913295560024261167135799300172382907308330135472647507,
        41885656274025752055847945432737871864088659248922821023734315208027501951872
      ],
      c: 96917856581077810363012153828220232197567408835708926581335248000925197916153,
      s: 103298896676233752268329042222773891728807677368628421408380318882272184455566,
      seed: 93724884573574303181157854277074121673523280784530506403108144933983063023487,
      uWitness: 0xFCaA10875C6692f6CcC86c64300eb0b52f2D4323,
      cGammaWitness: [
        61463607927970680172418313129927007099021056249775757132623753443657677198526,
        48686021866486086188742596461341782400160109177829661164208082534005682984658
      ],
      sHashWitness: [
        91508089836242281395929619352465003226819385335975246221498243754781593857533,
        63571625936444669399167157725633389238098818902162172059681813608664564703308
      ],
      zInv: 97568175302326019383632009699686265453584842953005404815285123863099260038246
    });
    VRFCoordinatorV2_5.RequestCommitment memory rc = VRFCoordinatorV2_5.RequestCommitment({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: 1_000_000,
      numWords: 1,
      sender: address(s_testConsumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    (, uint96 nativeBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);

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
    (, uint96 nativeBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(nativeBalanceAfter, nativeBalanceBefore - 120_000, 10_000);
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
    assertEq(hex"0000000000000000000000000000000000000000000000000000000000000014", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 108233140904510496268355288815996296196427471042093167619305836589216327096601 \
        -block-hash 0x0000000000000000000000000000000000000000000000000000000000000014 \
        -block-num 20 \
        -sender 0x90A8820424CC8a819d14cBdE54D12fD3fbFa9bb2
    */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        49785247270467418393187938018746488660500261614113251546613288843777654841004,
        8320717868018488740308781441198484312662094766876176838868269181386589318272
      ],
      c: 41596204381278553342984662603150353549780558761307588910860350083645227536604,
      s: 81592778991188138734863787790226463602813498664606420860910885269124681994753,
      seed: 108233140904510496268355288815996296196427471042093167619305836589216327096601,
      uWitness: 0x56920892EE71E624d369dCc8dc63B6878C85Ca70,
      cGammaWitness: [
        28250667431035633903490940933503696927659499415200427260709034207157951953043,
        105660182690338773283351292037478192732977803900032569393220726139772041021018
      ],
      sHashWitness: [
        18420263847278540234821121001488166570853056146131705862117248292063859054211,
        15740432967529684573970722302302642068194042971767150190061244675457227502736
      ],
      zInv: 100579074451139970455673776933943662313989441807178260211316504761358492254052
    });
    VRFCoordinatorV2_5.RequestCommitment memory rc = VRFCoordinatorV2_5.RequestCommitment({
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
    // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_native_ratio)
    // paymentNoFee = (1 * (50_000 + 90_000 + 0)) / .5
    // paymentNoFee = 280_000
    // ...
    // billed_fee = paymentNoFee + fulfillmentFlatFeeLinkPPM
    // billed_fee = baseFeeWei + 0
    // billed_fee = 280_000
    // note: delta is doubled from the native test to account for more variance due to the link/native ratio
    (uint96 linkBalanceAfter, , , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 280_000, 20_000);
  }
}
