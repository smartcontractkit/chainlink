pragma solidity 0.8.19;

import "../../shared/interfaces/LinkTokenInterface.sol";

import "./BaseTest.t.sol";
import "../dev/ERC-4337/SmartContractAccountFactory.sol";
import "../dev/testhelpers/SmartContractAccountHelper.sol";
import "../dev/ERC-4337/SCA.sol";
import "../dev/testhelpers/Greeter.sol";
import "../dev/ERC-4337/Paymaster.sol";
import "../../transmission/dev/ERC-4337/SCALibrary.sol";
import "../../mocks/MockLinkToken.sol";
import "../../tests/MockV3Aggregator.sol";
import "../../vrf/mocks/VRFCoordinatorMock.sol";
import "../../vrf/testhelpers/VRFConsumer.sol";

import "../../vendor/entrypoint/interfaces/UserOperation.sol";
import "../../vendor/entrypoint/core/EntryPoint.sol";
import "../../vendor/entrypoint/interfaces/IEntryPoint.sol";

/*--------------------------------------------------------------------------------------------------------------------+
| EIP 712 + 1014 + 4337                                                                                               |
| ________________                                                                                                    |
| This implementation allows for meta-transactions to be signed by end-users and posted on-chain by executors. It     |
| utilizes the following components:                                                                                  |
| - EIP-712: The method by which meta-transactions are authorized.                                                    |
| - EIP-1014: The method by which the Smart Contract Account is generated.                                            |
| - EIP-4337: The method by which meta-transactions are executed.                                                     |
|                                                                                                                     |
| The below tests illustrate end-user flows for interacting with this meta-transaction system. For users with         |
| existing Smart Contract Accounts (SCAs), they simply sign off on the operation, after which  the executor           |
| invokes the EntryPoint that authorizes the operation on the end-user's SCA, and then execute the transaction        |
| as the SCA. For users without existing SCAs, EIP-1014 ensures that the address of an SCA can be known in advance,   |
| so users can sign-off on transactions that will be executed by a not-yet-deployed SCA. The EntryPoint contract      |
| takes advantage of this functionality and allows for the SCA to be created in the same user operation that invokes  |
| it, and the end-user signs off on this creation-and-execution flow. After the initial creation-and-execution, the   |
| SCA is reused for future transactions.                                                                              |
|                                                                                                                     |
| End-Dapps/protocols do not need to be EIP-2771-compliant or accommodate any other kind of transaction standard.     |
| They can be interacted with out-of-the-box through the SCA, which acts in place of the user's EOA as their          |
| immutable identity.                                                                                                 |
|                                                                                                                     |
-+---------------------------------------------------------------------------------------------------------------------*/

contract EIP_712_1014_4337 is BaseTest {
  event RandomnessRequest(address indexed sender, bytes32 indexed keyHash, uint256 indexed seed, uint256 fee);

  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
  address internal ENTRY_POINT;

  Greeter greeter;
  EntryPoint entryPoint;
  MockV3Aggregator linkEthFeed;

  // Randomly generated private/public key pair.
  uint256 END_USER_PKEY = uint256(bytes32(hex"99d518dbfea4b4ec301390f7e26d53d711fa1ca0c1a6e4cbed89617d4c578a8e"));
  address END_USER = 0xB6708257D4E1bf0b8C144793fc2Ff3193C737ed1;

  function setUp() public override {
    BaseTest.setUp();
    // Fund user accounts;
    vm.deal(END_USER, 10_000 ether);
    vm.deal(LINK_WHALE, 10_000 ether);

    // Impersonate a LINK whale.
    changePrank(LINK_WHALE);

    // Create simple greeter contract.
    greeter = new Greeter();
    assertEq("", greeter.getGreeting());

    // Create entry point contract.
    entryPoint = new EntryPoint();
    ENTRY_POINT = address(entryPoint);

    // Deploy link/eth feed.
    linkEthFeed = new MockV3Aggregator(18, 5000000000000000); // .005 ETH
  }

  /// @dev Test case for user that already has a Smart Contract Account.
  /// @dev EntryPoint.sol should use the existing SCA to execute the meta transaction.
  function testEIP712EIP4337WithExistingSmartContractAccount() public {
    // Pre-calculate user smart contract account address.
    SmartContractAccountFactory factory = new SmartContractAccountFactory();
    address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(
      END_USER,
      ENTRY_POINT,
      address(factory)
    );

    // Deploy the end-contract.
    bytes32 salt = bytes32(uint256(uint160(END_USER)) << 96);
    bytes memory fullInitializeCode = SmartContractAccountHelper.getSCAInitCodeWithConstructor(END_USER, ENTRY_POINT);
    factory.deploySmartContractAccount(salt, fullInitializeCode);
    changePrank(END_USER);

    // Ensure a correct deployment and a functioning end-contract.
    uint256 contractCodeSize;
    assembly {
      contractCodeSize := extcodesize(toDeployAddress)
    }
    assertTrue(contractCodeSize > 0);
    assertEq(END_USER, SCA(toDeployAddress).i_owner());

    // Create the calldata for a setGreeting call.
    string memory greeting = "hi";
    bytes memory encodedGreetingCall = bytes.concat(Greeter.setGreeting.selector, abi.encode(greeting)); // abi.encodeWithSelector equivalent

    // Produce the final full end-tx encoding, to be used as calldata in the user operation.
    bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
      address(greeter),
      uint256(0),
      0,
      encodedGreetingCall
    );

    // Construct the user operation.
    UserOperation memory op = UserOperation({
      sender: toDeployAddress,
      nonce: 0,
      initCode: "",
      callData: fullEncoding,
      callGasLimit: 1_000_000,
      verificationGasLimit: 1_000_000,
      preVerificationGas: 10_000,
      maxFeePerGas: 100,
      maxPriorityFeePerGas: 200,
      paymasterAndData: "",
      signature: ""
    });

    // Sign user operation.
    bytes32 userOpHash = entryPoint.getUserOpHash(op);
    bytes32 fullHash = SCALibrary._getUserOpFullHash(userOpHash, toDeployAddress);
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(END_USER_PKEY, fullHash);
    op.signature = abi.encodePacked(r, s, v - 27);

    // Deposit funds for the transaction.
    entryPoint.depositTo{value: 10 ether}(toDeployAddress);

    // Execute the user operation.
    UserOperation[] memory operations = new UserOperation[](1);
    operations[0] = op;
    entryPoint.handleOps(operations, payable(END_USER));

    // Assert that the greeting was set.
    assertEq("hi", Greeter(greeter).getGreeting());
    assertEq(SCA(toDeployAddress).s_nonce(), uint256(1));
  }

  /// @dev Test case for fresh user, EntryPoint.sol should generate a
  /// @dev Smart Contract Account for them and execute the meta transaction.
  function testEIP712EIP4337AndCreateSmartContractAccount() public {
    // Pre-calculate user smart contract account address.
    SmartContractAccountFactory factory = new SmartContractAccountFactory();
    address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(
      END_USER,
      ENTRY_POINT,
      address(factory)
    );

    // Construct initCode byte array.
    bytes memory fullInitializeCode = SmartContractAccountHelper.getInitCode(address(factory), END_USER, ENTRY_POINT);

    // Create the calldata for a setGreeting call.
    string memory greeting = "bye";
    bytes memory encodedGreetingCall = bytes.concat(Greeter.setGreeting.selector, abi.encode(greeting));

    // Produce the final full end-tx encoding, to be used as calldata in the user operation.
    bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
      address(greeter),
      uint256(0),
      0,
      encodedGreetingCall
    );

    // Construct the user operation.
    UserOperation memory op = UserOperation({
      sender: toDeployAddress,
      nonce: 0,
      initCode: fullInitializeCode,
      callData: fullEncoding,
      callGasLimit: 1_000_000,
      verificationGasLimit: 1_000_000,
      preVerificationGas: 10_000,
      maxFeePerGas: 100,
      maxPriorityFeePerGas: 200,
      paymasterAndData: "",
      signature: ""
    });

    // Sign user operation.
    bytes32 userOpHash = entryPoint.getUserOpHash(op);
    bytes32 fullHash = SCALibrary._getUserOpFullHash(userOpHash, toDeployAddress);
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(END_USER_PKEY, fullHash);
    op.signature = abi.encodePacked(r, s, v - 27);

    // Deposit funds for the transaction.
    entryPoint.depositTo{value: 10 ether}(toDeployAddress);

    // Execute the user operation.
    UserOperation[] memory operations = new UserOperation[](1);
    operations[0] = op;
    entryPoint.handleOps(operations, payable(END_USER));

    // Assert that the greeting was set.
    assertEq("bye", Greeter(greeter).getGreeting());
    assertEq(SCA(toDeployAddress).s_nonce(), uint256(1));
    assertEq(SCA(toDeployAddress).i_owner(), END_USER);
  }

  /// @dev Test case for a user executing a setGreeting with a LINK token paymaster.
  function testEIP712EIP4337AndCreateSmartContractAccountWithPaymaster() public {
    // Pre-calculate user smart contract account address.
    SmartContractAccountFactory factory = new SmartContractAccountFactory();
    address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(
      END_USER,
      ENTRY_POINT,
      address(factory)
    );

    // Construct initCode byte array.
    bytes memory fullInitializeCode = SmartContractAccountHelper.getInitCode(address(factory), END_USER, ENTRY_POINT);

    // Create the calldata for a setGreeting call.
    string memory greeting = "good day";
    bytes memory encodedGreetingCall = bytes.concat(Greeter.setGreeting.selector, abi.encode(greeting));

    // Produce the final full end-tx encoding, to be used as calldata in the user operation.
    bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
      address(greeter),
      uint256(0),
      0,
      encodedGreetingCall
    );

    // Create Link token, and deposit into paymaster.
    MockLinkToken linkToken = new MockLinkToken();
    Paymaster paymaster = new Paymaster(LinkTokenInterface(address(linkToken)), linkEthFeed, ENTRY_POINT);
    linkToken.transferAndCall(address(paymaster), 1000 ether, abi.encode(address(toDeployAddress)));

    // Construct the user opeartion.
    UserOperation memory op = UserOperation({
      sender: toDeployAddress,
      nonce: 0,
      initCode: fullInitializeCode,
      callData: fullEncoding,
      callGasLimit: 1_000_000,
      verificationGasLimit: 1_500_000,
      preVerificationGas: 10_000,
      maxFeePerGas: 100,
      maxPriorityFeePerGas: 200,
      paymasterAndData: abi.encodePacked(address(paymaster)),
      signature: ""
    });

    // Sign user operation.
    bytes32 userOpHash = entryPoint.getUserOpHash(op);
    bytes32 fullHash = SCALibrary._getUserOpFullHash(userOpHash, toDeployAddress);
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(END_USER_PKEY, fullHash);
    op.signature = abi.encodePacked(r, s, v - 27);

    // Deposit funds for the transaction.
    entryPoint.depositTo{value: 10 ether}(address(paymaster));

    // Execute the user operation.
    UserOperation[] memory operations = new UserOperation[](1);
    operations[0] = op;
    entryPoint.handleOps(operations, payable(END_USER));

    // Assert that the greeting was set.
    assertEq("good day", Greeter(greeter).getGreeting());
    assertEq(SCA(toDeployAddress).s_nonce(), uint256(1));
  }

  /// @dev Test case for a VRF Request via LINK token paymaster and an SCA.
  function testEIP712EIP4337AndCreateSmartContractAccountWithPaymasterForVRFRequest() public {
    // Pre-calculate user smart contract account address.
    SmartContractAccountFactory factory = new SmartContractAccountFactory();
    address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(
      END_USER,
      ENTRY_POINT,
      address(factory)
    );

    // Construct initCode byte array.
    bytes memory fullInitializeCode = SmartContractAccountHelper.getInitCode(address(factory), END_USER, ENTRY_POINT);

    // Create the calldata for a VRF request.
    bytes32 keyhash = bytes32(uint256(123));
    uint256 fee = 1 ether;
    bytes memory encodedVRFRequestCallData = bytes.concat(
      VRFConsumer.doRequestRandomness.selector,
      abi.encode(keyhash, fee)
    );

    // Create the VRF Contracts
    MockLinkToken linkToken = new MockLinkToken();
    VRFCoordinatorMock vrfCoordinator = new VRFCoordinatorMock(address(linkToken));
    VRFConsumer vrfConsumer = new VRFConsumer(address(vrfCoordinator), address(linkToken));

    // Produce the final full end-tx encoding, to be used as calldata in the user operation.
    bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
      address(vrfConsumer), // end-contract
      uint256(0), // value
      0, // timeout (seconds)
      encodedVRFRequestCallData
    );

    // Create Link token, and deposit into paymaster.
    Paymaster paymaster = new Paymaster(LinkTokenInterface(address(linkToken)), linkEthFeed, ENTRY_POINT);
    linkToken.transferAndCall(address(paymaster), 1000 ether, abi.encode(address(toDeployAddress)));

    // Construct direct funding data.
    SCALibrary.DirectFundingData memory directFundingData = SCALibrary.DirectFundingData({
      recipient: address(vrfConsumer),
      topupThreshold: 1,
      topupAmount: 10 ether
    });

    // Construct the user operation.
    UserOperation memory op = UserOperation({
      sender: toDeployAddress,
      nonce: 0,
      initCode: fullInitializeCode,
      callData: fullEncoding,
      callGasLimit: 200_000,
      verificationGasLimit: 1_000_000,
      preVerificationGas: 10_000,
      maxFeePerGas: 10,
      maxPriorityFeePerGas: 10,
      paymasterAndData: abi.encodePacked(address(paymaster), uint8(0), abi.encode(directFundingData)),
      signature: ""
    });

    // Sign user operation.
    bytes32 fullHash = SCALibrary._getUserOpFullHash(entryPoint.getUserOpHash(op), toDeployAddress);
    op.signature = getSignature(fullHash);

    // Deposit funds for the transaction.
    entryPoint.depositTo{value: 10 ether}(address(paymaster));

    // Assert correct log is emitted for the end-contract vrf request.
    vm.expectEmit(true, true, true, true);
    emit RandomnessRequest(
      address(vrfConsumer),
      keyhash,
      0, // seed - we use a zero seed
      fee
    );

    // Execute the user operation.
    UserOperation[] memory operations = new UserOperation[](1);
    operations[0] = op;

    // Execute user operation and ensure correct outcome.
    entryPoint.handleOps(operations, payable(END_USER));
    assertEq(SCA(toDeployAddress).s_nonce(), uint256(1));
  }

  function getSignature(bytes32 h) internal view returns (bytes memory) {
    (uint8 v, bytes32 r, bytes32 s) = vm.sign(END_USER_PKEY, h);
    return abi.encodePacked(r, s, v - 27);
  }
}
