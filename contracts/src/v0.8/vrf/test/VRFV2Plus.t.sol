pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {VRF} from "../VRF.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2_5} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFCoordinatorV2_5} from "../dev/VRFCoordinatorV2_5.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {BlockhashStore} from "../dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../dev/testhelpers/VRFV2PlusConsumerExample.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {VRFTypes} from "../VRFTypes.sol";
import {console} from "forge-std/console.sol";
import {VmSafe} from "forge-std/Vm.sol";
import {VRFV2PlusLoadTestWithMetrics} from "../dev/testhelpers/VRFV2PlusLoadTestWithMetrics.sol";
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
  uint64 internal constant GAS_LANE_MAX_GAS = 5000 gwei;
  uint16 internal constant MIN_CONFIRMATIONS = 0;
  uint32 internal constant CALLBACK_GAS_LIMIT = 1_000_000;
  uint32 internal constant NUM_WORDS = 1;

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

    vm.txGasPrice(100 gwei);

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

  function setConfig() internal {
    s_testCoordinator.setConfig(
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

  function testSetConfig() public {
    // Should setConfig successfully.
    setConfig();

    // Test that setting requestConfirmations above MAX_REQUEST_CONFIRMATIONS reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidRequestConfirmations.selector, 500, 500, 200));
    s_testCoordinator.setConfig(
      500,
      2_500_000,
      1,
      50_000,
      50000000000000000,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );

    // Test that setting fallbackWeiPerUnitLink to zero reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidLinkWeiPrice.selector, 0));

    s_testCoordinator.setConfig(
      0,
      2_500_000,
      1,
      50_000,
      0,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );

    // Test that setting link discount flat fee higher than native flat fee reverts
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.LinkDiscountTooHigh.selector, uint32(501), uint32(500)));

    s_testCoordinator.setConfig(
      0,
      2_500_000,
      1,
      50_000,
      500,
      500, // fulfillmentFlatFeeNativePPM
      501, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );

    // // Test that setting link discount flat fee equal to native flat fee does not revert
    s_testCoordinator.setConfig(
      0,
      2_500_000,
      1,
      50_000,
      500,
      450, // fulfillmentFlatFeeNativePPM
      450, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );

    // Test that setting native premium percentage higher than 155 will revert
    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidPremiumPercentage.selector, uint8(156), uint8(155))
    );

    s_testCoordinator.setConfig(
      0,
      2_500_000,
      1,
      50_000,
      500,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      156, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );

    // Test that setting LINK premium percentage higher than 155 will revert
    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.InvalidPremiumPercentage.selector, uint8(202), uint8(155))
    );

    s_testCoordinator.setConfig(
      0,
      2_500_000,
      1,
      50_000,
      500,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      202 // linkPremiumPercentage
    );
  }

  function testRegisterProvingKey() public {
    // Should set the proving key successfully.
    registerProvingKey();

    // Should revert when already registered.
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.ProvingKeyAlreadyRegistered.selector, vrfKeyHash));
    s_testCoordinator.registerProvingKey(uncompressedKeyParts, GAS_LANE_MAX_GAS);
  }

  event ProvingKeyRegistered(bytes32 keyHash, uint64 maxGas);
  event ProvingKeyDeregistered(bytes32 keyHash, uint64 maxGas);

  function registerProvingKey() public {
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    bytes32 keyHash = keccak256(abi.encode(uncompressedKeyParts));
    vm.expectEmit(
      false, // no indexed args to check for
      false, // no indexed args to check for
      false, // no indexed args to check for
      true
    ); // check data fields: keyHash and maxGas
    emit ProvingKeyRegistered(keyHash, GAS_LANE_MAX_GAS);
    s_testCoordinator.registerProvingKey(uncompressedKeyParts, GAS_LANE_MAX_GAS);
    (bool exists, uint64 maxGas) = s_testCoordinator.s_provingKeys(keyHash);
    assertTrue(exists);
    assertEq(GAS_LANE_MAX_GAS, maxGas);
    assertEq(s_testCoordinator.s_provingKeyHashes(0), keyHash);
    assertEq(keyHash, vrfKeyHash);
  }

  function testDeregisterProvingKey() public {
    // Should set the proving key successfully.
    registerProvingKey();

    bytes
      memory unregisteredPubKey = hex"6d919e4ed6add6c34b2af77eb6b2d2f5d27db11ba004e70734b23bd4321ea234ff8577a063314bead6d88c1b01849289a5542767a5138924f38fed551a7773db";

    // Should revert when given pubkey is not registered
    uint256[2] memory unregisteredKeyParts = this.getProvingKeyParts(unregisteredPubKey);
    bytes32 unregisterdKeyHash = keccak256(abi.encode(unregisteredKeyParts));
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.NoSuchProvingKey.selector, unregisterdKeyHash));
    s_testCoordinator.deregisterProvingKey(unregisteredKeyParts);

    // correctly deregister pubkey
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    bytes32 keyHash = keccak256(abi.encode(uncompressedKeyParts));
    vm.expectEmit(
      false, // no indexed args to check for
      false, // no indexed args to check for
      false, // no indexed args to check for
      true
    ); // check data fields: keyHash and maxGas
    emit ProvingKeyDeregistered(keyHash, GAS_LANE_MAX_GAS);
    s_testCoordinator.deregisterProvingKey(uncompressedKeyParts);
    (bool exists, uint64 maxGas) = s_testCoordinator.s_provingKeys(keyHash);
    assertFalse(exists);
    assertEq(0, maxGas);
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
  event FallbackWeiPerUnitLinkUsed(uint256 requestId, int256 fallbackWeiPerUnitLink);

  function testRequestAndFulfillRandomWordsNative() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      uint256 subId,
      uint256 requestId
    ) = setupSubAndRequestRandomnessNativePayment();
    (, uint96 nativeBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    uint96 payment = s_testCoordinator.fulfillRandomWords(proof, rc, false);
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 70_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())
    // baseFeeWei = 1e11 * (50_000 + 70_000)
    // baseFeeWei = 1.2e16
    // flatFeeWei = 1e12 * (fulfillmentFlatFeeNativePPM)
    // flatFeeWei = 1e12 * 500_000 = 5e17
    // ...
    // billed_fee = baseFeeWei * (100 + linkPremiumPercentage / 100) + 5e17
    // billed_fee = 1.2e16 * 1.15 + 5e17
    // billed_fee = 5.138e+17
    (, uint96 nativeBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.138 * 1e17, 1e15);
    assertApproxEqAbs(nativeBalanceAfter, nativeBalanceBefore - 5.138 * 1e17, 1e15);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestAndFulfillRandomWordsLINK() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      uint256 subId,
      uint256 requestId
    ) = setupSubAndRequestRandomnessLINKPayment();
    (uint96 linkBalanceBefore, , , , ) = s_testCoordinator.getSubscription(subId);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    uint96 payment = s_testCoordinator.fulfillRandomWords(proof, rc, false);

    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 86_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_native_ratio)
    // paymentNoFee = (1e11 * (50_000 + 86_000 + 0)) / .5
    // paymentNoFee = 2.72e16
    // flatFeeWei = 1e12 * (fulfillmentFlatFeeNativePPM - fulfillmentFlatFeeLinkDiscountPPM)
    // flatFeeWei = 1e12 * (500_000 - 100_000)
    // flatFeeJuels = 1e18 * flatFeeWei / link_native_ratio
    // flatFeeJuels = 4e17 / 0.5 = 8e17
    // billed_fee = paymentNoFee * ((100 + 10) / 100) + 8e17
    // billed_fee = 2.72e16 * 1.1 + 8e17
    // billed_fee = 2.992e16 + 8e17 = 8.2992e17
    // note: delta is doubled from the native test to account for more variance due to the link/native ratio
    (uint96 linkBalanceAfter, , , , ) = s_testCoordinator.getSubscription(subId);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 8.2992 * 1e17, 1e15);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 8.2992 * 1e17, 1e15);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestAndFulfillRandomWordsLINK_FallbackWeiPerUnitLinkUsed() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      ,
      uint256 requestId
    ) = setupSubAndRequestRandomnessLINKPayment();

    (, , , uint32 stalenessSeconds, , , , , ) = s_testCoordinator.s_config();
    int256 fallbackWeiPerUnitLink = s_testCoordinator.s_fallbackWeiPerUnitLink();

    // Set the link feed to be stale.
    (uint80 roundId, int256 answer, uint256 startedAt, , ) = s_linkNativeFeed.latestRoundData();
    uint256 timestamp = block.timestamp - stalenessSeconds - 1;
    s_linkNativeFeed.updateRoundData(roundId, answer, timestamp, startedAt);

    vm.expectEmit(false, false, false, true, address(s_testCoordinator));
    emit FallbackWeiPerUnitLinkUsed(requestId, fallbackWeiPerUnitLink);
    s_testCoordinator.fulfillRandomWords(proof, rc, false);
  }

  function setupSubAndRequestRandomnessLINKPayment()
    internal
    returns (VRF.Proof memory proof, VRFTypes.RequestCommitmentV2Plus memory rc, uint256 subId, uint256 requestId)
  {
    uint32 requestBlock = 20;
    vm.roll(requestBlock);
    s_linkToken.transfer(address(s_testConsumer), 10 ether);
    s_testConsumer.createSubscriptionAndFund(10 ether);
    subId = s_testConsumer.s_subId();

    // Apply basic configs to contract.
    setConfig();
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, false, true);
    uint256 preSeed;
    (requestId, preSeed) = s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(s_testConsumer), subId, 1);
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      subId,
      MIN_CONFIRMATIONS,
      CALLBACK_GAS_LIMIT,
      NUM_WORDS,
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false})), // nativePayment, // nativePayment
      address(s_testConsumer) // requester
    );
    s_testConsumer.requestRandomWords(CALLBACK_GAS_LIMIT, MIN_CONFIRMATIONS, NUM_WORDS, vrfKeyHash, false);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);
    assertTrue(s_testCoordinator.pendingRequestExists(subId));

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
        -pre-seed 58424872742560034068603954318478134981993109073728628043159461959392650534066 \
        -block-hash 0x0000000000000000000000000000000000000000000000000000000000000014 \
        -block-num 20 \
        -sender 0x90A8820424CC8a819d14cBdE54D12fD3fbFa9bb2
    */
    proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        38041205470219573731614166317842050442610096576830191475863676359766283013831,
        31897503406364148988967447112698248795931483458172800286988696482435433838056
      ],
      c: 114706080610174375269579192101772790158458728655229562781479812703475130740224,
      s: 91869928024010088265014058436030407245056128545665425448353233998362687232253,
      seed: 58424872742560034068603954318478134981993109073728628043159461959392650534066,
      uWitness: 0x1514536B09a51E671d070312bcD3653386d5a82b,
      cGammaWitness: [
        90605489216274499662544489893800286859751132311034850249229378789467669572783,
        76568417372883461229305641415175605031997103681542349721251313705711146936024
      ],
      sHashWitness: [
        43417948503950579681520475434461454031791886587406480417092620950034789197994,
        100772571879140362396088596211082924128900752544164141322636815729889228000249
      ],
      zInv: 82374292458278672300647114418593830323283909625362447038989596015264004164958
    });
    rc = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: 1000000,
      numWords: 1,
      sender: address(s_testConsumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}))
    });
    return (proof, rc, subId, requestId);
  }

  function setupSubAndRequestRandomnessNativePayment()
    internal
    returns (VRF.Proof memory proof, VRFTypes.RequestCommitmentV2Plus memory rc, uint256 subId, uint256 requestId)
  {
    uint32 requestBlock = 10;
    vm.roll(requestBlock);
    s_testConsumer.createSubscriptionAndFund(0);
    subId = s_testConsumer.s_subId();
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(subId);

    // Apply basic configs to contract.
    setConfig();
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, true, true);
    uint256 preSeed;
    (requestId, preSeed) = s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(s_testConsumer), subId, 1);
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      subId,
      MIN_CONFIRMATIONS,
      CALLBACK_GAS_LIMIT,
      NUM_WORDS,
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true})), // nativePayment
      address(s_testConsumer) // requester
    );
    s_testConsumer.requestRandomWords(CALLBACK_GAS_LIMIT, MIN_CONFIRMATIONS, NUM_WORDS, vrfKeyHash, true);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);
    assertTrue(s_testCoordinator.pendingRequestExists(subId));

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
        -pre-seed 83266692323404068105564931899467966321583332182309426611016082057597749986430 \
        -block-hash 0x000000000000000000000000000000000000000000000000000000000000000a \
        -block-num 10 \
        -sender 0x90A8820424CC8a819d14cBdE54D12fD3fbFa9bb2 \
        -native-payment true
        */
    proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        47144451677122876068574640250190132179872561942855874114516471019540736524783,
        63001220656590641645486673489302242739512599229187442248048295264418080499391
      ],
      c: 42928477813589729783511577059394077774341588261592343937605454161333818133643,
      s: 14447529458406454898597883219032514356523135029224613793880920230249515634875,
      seed: 83266692323404068105564931899467966321583332182309426611016082057597749986430,
      uWitness: 0x5Ed3bb2AA8EAFe168a23079644d5dfBf892B1038,
      cGammaWitness: [
        40742088032247467257043132769297935807697466810312051815364187117543257089153,
        110399474382135664619186049639190334359061769014381608543009407662815758204131
      ],
      sHashWitness: [
        26556776392895292893984393164594214244553035014769995354896600239759043777485,
        67126706735912782218279556535631175449291035782208850310724682668198932501077
      ],
      zInv: 88742453392918610091640193378775723954629905126315835248392650970979000380325
    });
    rc = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: 1,
      sender: address(s_testConsumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });

    return (proof, rc, subId, requestId);
  }

  function testRequestAndFulfillRandomWords_NetworkGasPriceExceedsGasLane() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      ,

    ) = setupSubAndRequestRandomnessNativePayment();

    // network gas is higher than gas lane max gas
    uint256 networkGasPrice = GAS_LANE_MAX_GAS + 1;
    vm.txGasPrice(networkGasPrice);
    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.GasPriceExceeded.selector, networkGasPrice, GAS_LANE_MAX_GAS)
    );
    s_testCoordinator.fulfillRandomWords(proof, rc, false);
  }

  function testRequestAndFulfillRandomWords_OnlyPremium_NativePayment() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      uint256 subId,
      uint256 requestId
    ) = setupSubAndRequestRandomnessNativePayment();
    (, uint96 nativeBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);

    // network gas is twice the gas lane max gas
    uint256 networkGasPrice = GAS_LANE_MAX_GAS * 2;
    vm.txGasPrice(networkGasPrice);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    uint96 payment = s_testCoordinator.fulfillRandomWords(proof, rc, true /* onlyPremium */);
    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 70_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())
    // network gas price is capped at gas lane max gas (5000 gwei)
    // baseFeeWei = 5e12 * (50_000 + 70_000)
    // baseFeeWei = 6e17
    // flatFeeWei = 1e12 * (fulfillmentFlatFeeNativePPM)
    // flatFeeWei = 1e12 * 500_000 = 5e17
    // ...
    // billed_fee = baseFeeWei * (linkPremiumPercentage / 100) + 5e17
    // billed_fee = 6e17 * 0.15 + 5e17
    // billed_fee = 5.9e+17
    (, uint96 nativeBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 5.9 * 1e17, 1e15);
    assertApproxEqAbs(nativeBalanceAfter, nativeBalanceBefore - 5.9 * 1e17, 1e15);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestAndFulfillRandomWords_OnlyPremium_LinkPayment() public {
    (
      VRF.Proof memory proof,
      VRFTypes.RequestCommitmentV2Plus memory rc,
      uint256 subId,
      uint256 requestId
    ) = setupSubAndRequestRandomnessLINKPayment();
    (uint96 linkBalanceBefore, , , , ) = s_testCoordinator.getSubscription(subId);

    // network gas is twice the gas lane max gas
    uint256 networkGasPrice = GAS_LANE_MAX_GAS * 5;
    vm.txGasPrice(networkGasPrice);

    uint256 outputSeed = s_testCoordinator.getRandomnessFromProofExternal(proof, rc).randomness;
    vm.recordLogs();
    uint96 payment = s_testCoordinator.fulfillRandomWords(proof, rc, true /* onlyPremium */);

    VmSafe.Log[] memory entries = vm.getRecordedLogs();
    assertEq(entries[0].topics[1], bytes32(uint256(requestId)));
    assertEq(entries[0].topics[2], bytes32(uint256(subId)));
    (uint256 loggedOutputSeed, , , bool loggedSuccess) = abi.decode(entries[0].data, (uint256, uint256, bool, bool));
    assertEq(loggedOutputSeed, outputSeed);
    assertEq(loggedSuccess, true);

    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 86_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_native_ratio)
    // network gas price is capped at gas lane max gas (5000 gwei)
    // paymentNoFee = (5e12 * (50_000 + 86_000 + 0)) / .5
    // paymentNoFee = 1.36e+18
    // flatFeeWei = 1e12 * (fulfillmentFlatFeeNativePPM - fulfillmentFlatFeeLinkDiscountPPM)
    // flatFeeWei = 1e12 * (500_000 - 100_000)
    // flatFeeJuels = 1e18 * flatFeeWei / link_native_ratio
    // flatFeeJuels = 4e17 / 0.5 = 8e17
    // billed_fee = paymentNoFee * (10 / 100) + 8e17
    // billed_fee = 1.36e+18 * 0.1 + 8e17
    // billed_fee = 9.36e+17
    // note: delta is doubled from the native test to account for more variance due to the link/native ratio
    (uint96 linkBalanceAfter, , , , ) = s_testCoordinator.getSubscription(subId);
    // 1e15 is less than 1 percent discrepancy
    assertApproxEqAbs(payment, 9.36 * 1e17, 1e15);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 9.36 * 1e17, 1e15);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestRandomWords_InvalidConsumer() public {
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint256 subId = s_testCoordinator.createSubscription();
    VRFV2PlusLoadTestWithMetrics consumer = new VRFV2PlusLoadTestWithMetrics(address(s_testCoordinator));

    // consumer is not added to the subscription
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(consumer)));
    consumer.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      true,
      NUM_WORDS,
      1 /* requestCount */
    );
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestRandomWords_ReAddConsumer_AssertRequestID() public {
    // 1. setup consumer and subscription
    setConfig();
    registerProvingKey();
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint256 subId = s_testCoordinator.createSubscription();
    VRFV2PlusLoadTestWithMetrics consumer = createAndAddLoadTestWithMetricsConsumer(subId);
    uint32 requestBlock = 10;
    vm.roll(requestBlock);
    changePrank(LINK_WHALE);
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(subId);

    // 2. Request random words.
    changePrank(subOwner);
    vm.expectEmit(true, true, false, true);
    uint256 requestId;
    uint256 preSeed;
    (requestId, preSeed) = s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(consumer), subId, 1);
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      subId,
      MIN_CONFIRMATIONS,
      CALLBACK_GAS_LIMIT,
      NUM_WORDS,
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true})),
      address(consumer) // requester
    );
    consumer.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      true /* nativePayment */,
      NUM_WORDS,
      1 /* requestCount */
    );
    assertTrue(s_testCoordinator.pendingRequestExists(subId));

    // 3. Fulfill the request above
    //console.log("requestId: ", requestId);
    //console.log("preSeed: ", preSeed);
    //console.log("sender: ", address(consumer));

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
      -pre-seed 94043941380654896554739370173616551044559721638888689173752661912204412136884 \
      -block-hash 0x000000000000000000000000000000000000000000000000000000000000000a \
      -block-num 10 \
      -sender 0x44CAfC03154A0708F9DCf988681821f648dA74aF \
      -native-payment true
    */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        18593555375562408458806406536059989757338587469093035962641476877033456068708,
        55675218112764789548330682504442195066741636758414578491295297591596761905475
      ],
      c: 56595337384472359782910435918403237878894172750128610188222417200315739516270,
      s: 60666722370046279064490737533582002977678558769715798604164042022636022215663,
      seed: 94043941380654896554739370173616551044559721638888689173752661912204412136884,
      uWitness: 0xEdbE15fd105cfEFb9CCcbBD84403d1F62719E50d,
      cGammaWitness: [
        11752391553651713021860307604522059957920042356542944931263270793211985356642,
        14713353048309058367510422609936133400473710094544154206129568172815229277104
      ],
      sHashWitness: [
        109716108880570827107616596438987062129934448629902940427517663799192095060206,
        79378277044196229730810703755304140279837983575681427317104232794580059801930
      ],
      zInv: 18898957977631212231148068121702167284572066246731769473720131179584458697812
    });
    VRFTypes.RequestCommitmentV2Plus memory rc = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: NUM_WORDS,
      sender: address(consumer),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    s_testCoordinator.fulfillRandomWords(proof, rc, true /* onlyPremium */);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));

    // 4. remove consumer and verify request random words doesn't work
    s_testCoordinator.removeConsumer(subId, address(consumer));
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(consumer)));
    consumer.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      false /* nativePayment */,
      NUM_WORDS,
      1 /* requestCount */
    );

    // 5. re-add consumer and assert requestID nonce starts from 2 (nonce 1 was used before consumer removal)
    s_testCoordinator.addConsumer(subId, address(consumer));
    vm.expectEmit(true, true, false, true);
    uint256 requestId2;
    uint256 preSeed2;
    (requestId2, preSeed2) = s_testCoordinator.computeRequestIdExternal(vrfKeyHash, address(consumer), subId, 2);
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId2,
      preSeed2,
      subId,
      MIN_CONFIRMATIONS,
      CALLBACK_GAS_LIMIT,
      NUM_WORDS,
      VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false})), // nativePayment, // nativePayment
      address(consumer) // requester
    );
    consumer.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      false /* nativePayment */,
      NUM_WORDS,
      1 /* requestCount */
    );
    assertNotEq(requestId, requestId2);
    assertNotEq(preSeed, preSeed2);
    assertTrue(s_testCoordinator.pendingRequestExists(subId));
  }

  function testRequestRandomWords_MultipleConsumers_PendingRequestExists() public {
    // 1. setup consumer and subscription
    setConfig();
    registerProvingKey();
    address subOwner = makeAddr("subOwner");
    changePrank(subOwner);
    uint256 subId = s_testCoordinator.createSubscription();
    VRFV2PlusLoadTestWithMetrics consumer1 = createAndAddLoadTestWithMetricsConsumer(subId);
    VRFV2PlusLoadTestWithMetrics consumer2 = createAndAddLoadTestWithMetricsConsumer(subId);
    uint32 requestBlock = 10;
    vm.roll(requestBlock);
    changePrank(LINK_WHALE);
    s_testCoordinator.fundSubscriptionWithNative{value: 10 ether}(subId);

    // 2. Request random words.
    changePrank(subOwner);
    (uint256 requestId1, uint256 preSeed1) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(consumer1),
      subId,
      1
    );
    (uint256 requestId2, uint256 preSeed2) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(consumer2),
      subId,
      1
    );
    assertNotEq(requestId1, requestId2);
    assertNotEq(preSeed1, preSeed2);
    consumer1.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      true /* nativePayment */,
      NUM_WORDS,
      1 /* requestCount */
    );
    consumer2.requestRandomWords(
      subId,
      MIN_CONFIRMATIONS,
      vrfKeyHash,
      CALLBACK_GAS_LIMIT,
      true /* nativePayment */,
      NUM_WORDS,
      1 /* requestCount */
    );
    assertTrue(s_testCoordinator.pendingRequestExists(subId));

    // Move on to the next block.
    // Store the previous block's blockhash, and assert that it is as expected.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);
    assertEq(hex"000000000000000000000000000000000000000000000000000000000000000a", s_bhs.getBlockhash(requestBlock));

    // 3. Fulfill the 1st request above
    console.log("requestId: ", requestId1);
    console.log("preSeed: ", preSeed1);
    console.log("sender: ", address(consumer1));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
      go run . generate-proof-v2-plus \
      -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
      -pre-seed 94043941380654896554739370173616551044559721638888689173752661912204412136884 \
      -block-hash 0x000000000000000000000000000000000000000000000000000000000000000a \
      -block-num 10 \
      -sender 0x44CAfC03154A0708F9DCf988681821f648dA74aF \
      -native-payment true
    */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        18593555375562408458806406536059989757338587469093035962641476877033456068708,
        55675218112764789548330682504442195066741636758414578491295297591596761905475
      ],
      c: 56595337384472359782910435918403237878894172750128610188222417200315739516270,
      s: 60666722370046279064490737533582002977678558769715798604164042022636022215663,
      seed: 94043941380654896554739370173616551044559721638888689173752661912204412136884,
      uWitness: 0xEdbE15fd105cfEFb9CCcbBD84403d1F62719E50d,
      cGammaWitness: [
        11752391553651713021860307604522059957920042356542944931263270793211985356642,
        14713353048309058367510422609936133400473710094544154206129568172815229277104
      ],
      sHashWitness: [
        109716108880570827107616596438987062129934448629902940427517663799192095060206,
        79378277044196229730810703755304140279837983575681427317104232794580059801930
      ],
      zInv: 18898957977631212231148068121702167284572066246731769473720131179584458697812
    });
    VRFTypes.RequestCommitmentV2Plus memory rc = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: NUM_WORDS,
      sender: address(consumer1),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    s_testCoordinator.fulfillRandomWords(proof, rc, true /* onlyPremium */);
    assertTrue(s_testCoordinator.pendingRequestExists(subId));

    //4. Fulfill the 2nd request
    console.log("requestId: ", requestId2);
    console.log("preSeed: ", preSeed2);
    console.log("sender: ", address(consumer2));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
      go run . generate-proof-v2-plus \
      -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
      -pre-seed 60086281972849674111646805013521068579710860774417505336898013292594859262126 \
      -block-hash 0x000000000000000000000000000000000000000000000000000000000000000a \
      -block-num 10 \
      -sender 0xf5a165378E120f93784395aDF1E08a437e902865 \
      -native-payment true
    */
    proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        8781676794493524976318989249067879326013864868749595045909181134740761572122,
        70144896394968351242907510966944756907625107566821127114847472296460405612124
      ],
      c: 67847193668837615807355025316836592349514589069599294392546721746916067719949,
      s: 114946531382736685625345450298146929067341928840493664822961336014597880904075,
      seed: 60086281972849674111646805013521068579710860774417505336898013292594859262126,
      uWitness: 0xe1de4fD69277D0C5516cAE4d760b1d08BC340A28,
      cGammaWitness: [
        90301582727701442026215692513959255065128476395727596945643431833363167168678,
        61501369717028493801369453424028509804064958915788808540582630993703331669978
      ],
      sHashWitness: [
        98738650825542176387169085844714248077697103572877410412808249468787326424906,
        85647963391545223707301702874240345890884970941786094239896961457539737216630
      ],
      zInv: 29080001901010358083725892808339807464533563010468652346220922643802059192842
    });
    rc = VRFTypes.RequestCommitmentV2Plus({
      blockNum: requestBlock,
      subId: subId,
      callbackGasLimit: CALLBACK_GAS_LIMIT,
      numWords: NUM_WORDS,
      sender: address(consumer2),
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}))
    });
    s_testCoordinator.fulfillRandomWords(proof, rc, true /* onlyPremium */);
    assertFalse(s_testCoordinator.pendingRequestExists(subId));
  }

  function createAndAddLoadTestWithMetricsConsumer(uint256 subId) internal returns (VRFV2PlusLoadTestWithMetrics) {
    VRFV2PlusLoadTestWithMetrics consumer = new VRFV2PlusLoadTestWithMetrics(address(s_testCoordinator));
    s_testCoordinator.addConsumer(subId, address(consumer));
    return consumer;
  }
}
