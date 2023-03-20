// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "./consumers/SmartContractAccountFactory.sol";
import "./utils/SmartContractAccountHelper.sol";
import "./consumers/SCA.sol";
import "./consumers/Greeter.sol";
import "./consumers/Paymaster.sol";
import "./interfaces/UserOperation.sol";
import "./contracts/EntryPoint.sol";
import "./interfaces/IEntryPoint.sol";
import "./consumers/SCALibrary.sol";

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
| invokes the EntryPoint that authorizes the operation on the end-user's SCA, and then exectute the transaction       |
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

/*----------------------------+
| TESTS                       |
| ________________            |
|                             |
+----------------------------*/

contract EIP_712_1014_4337 is Test {
    address internal constant LINK_WHALE =
        0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;
    address internal constant LINK_WHALE_2 =
        0xeFF41C8725be95e66F6B10489B6bF34b08055853;
    address internal constant FEE_TOKEN =
        0x779877A7B0D9E8603169DdbD7836e478b4624789;
    address internal constant ROUTER =
        0xc4f3c3bb9e58AB406450cC1704f20F94E60b02F3;
    address internal constant CONTRACT_ADDRESS =
        0xdF525f47cFE6FC994E1dC9779EbD746880eC8E70;
    address internal constant CREATE2_FACTORY =
        0x93a5d90FBD40fEeBCE2ae5e80A2d7D2EfDbb39B4;
    address internal constant ENTRY_POINT =
        0x0576a174D229E3cFA37253523E645A78A0C91B57;
    address internal constant SENDER_CREATOR =
        0x932a3A220aC2CD48fab18118954601f565f19681;

    Greeter greeter;
    EntryPoint entryPoint;

    bytes signature = // Signature by LINK_WHALE. Signs off on "hi" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"3b8fc44066576110ba9b3f77b479fdec9cc6a3bba5d25f80f27a7dd40e4ece550becaeca8473b82a208bcfe1bf1ad89a32fc2dd5fa642f931f748fd62a8717e601";

    bytes signature2 = // Signature by LINK_WHALE_2. Signs off on "bye" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"cd6317fd718037261501fce54cd7355d646b39f5b75e91d5f01b2c032a56035e43eba05c4da914511fdcb806d75596f1b7a9a762a47f1eb5711168be6f3355ac01";

    bytes signature3 = // Signature #2 by LINK_WHALE_2. Signs off on "bye" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"24e58a624236d33dbd46162752090eb10b3bf2b46e291cb3ea362b92b6da4e833fddb7512ae43ecdefc645cf6f105cb98d73625e44a35db5a57f151f2608fea501";

    function setUp() public {
        // Fork Goerli.
        uint256 mainnetFork = vm.createFork(
            "https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161" // public ETH Goerli RPC URL
        );
        vm.selectFork(mainnetFork);
        vm.rollFork(8598894);

        // Impersonate a LINK whale.
        changePrank(LINK_WHALE);

        // Create simople greeter contract.
        greeter = Greeter(0x5BF01BcFBAf58DC33f2Ca062bf5f6fBe055Ac11b);
        assertEq("", greeter.getGreeting());

        // Use existing entry point contract.
        entryPoint = EntryPoint(payable(ENTRY_POINT));
    }

    /// @dev Test case for user that already has a Smart Contract Account.
    /// @dev EntryPoint.sol should use the existing SCA to execute the meta transaction.
    function testEIP712EIP4337WithExistingSmartContractAccount() public {
        SmartContractAccountFactory factory = new SmartContractAccountFactory();

        // Pre-calculate user smart contract account address.
        address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(LINK_WHALE, ENTRY_POINT, address(factory));

        // Deploy the end-contract.
        changePrank(SENDER_CREATOR);
        bytes32 salt = bytes32(uint256(uint160(LINK_WHALE)) << 96);
        bytes memory fullInitializeCode = SmartContractAccountHelper.getSCAInitCodeWithConstructor(LINK_WHALE, ENTRY_POINT);
        factory.deploySmartContractAccount(salt, fullInitializeCode);
        changePrank(LINK_WHALE);

        // Ensure a correct deployment and a functioning end-contract.
        uint256 contractCodeSize;
        assembly {
            contractCodeSize := extcodesize(toDeployAddress)
        }
        assertTrue(contractCodeSize > 0);
        assertEq(LINK_WHALE, SCA(toDeployAddress).s_owner());

        // Create the calldata for a setGreeting call.
        string memory greeting = "hi";
        bytes memory encodedGreetingCall = bytes.concat( // abi.encodeWithSelector equivalent
            Greeter.setGreeting.selector,
            abi.encode(greeting)
        );

        // Produce the final full end-tx encoding, to be used as calldata in the user operation.
        bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
            address(greeter),
            uint256(0),
            1000,
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
            preVerificationGas: 1_000_000,
            maxFeePerGas: 100,
            maxPriorityFeePerGas: 200,
            paymasterAndData: "",
            signature: signature
        });

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes32 userOpHash = entryPoint.getUserOpHash(op);
        bytes32 fullHash = SmartContractAccountHelper.getFullHashForSigning(userOpHash);
        console.logBytes32(fullHash);

        // Deposit funds for the transaction.
        entryPoint.depositTo{value: 10 ether}(toDeployAddress);

        // Execute the user operation.
        UserOperation[] memory operations = new UserOperation[](1);
        operations[0] = op;
        entryPoint.handleOps(operations, payable(LINK_WHALE));

        // Assert that the greeting was set.
        assertEq("hi", Greeter(greeter).getGreeting());
    }

    /// @dev Test case for fresh user, EntryPoint.sol should generate a 
    /// @dev Smart Contract Account for them and execute the meta transaction.
    function testEIP712EIP4337AndCreateSmartContractAccount() public {
        // Impersonate a different LINK whale.
        changePrank(LINK_WHALE);

        // Pre-calculate user smart contract account address.
        SmartContractAccountFactory factory = new SmartContractAccountFactory();
        address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(LINK_WHALE_2, ENTRY_POINT, address(factory));

        // Construct initCode byte array.
        bytes memory fullInitializeCode = SmartContractAccountHelper.getInitCode(address(factory), LINK_WHALE_2, ENTRY_POINT);

        // Create the calldata for a setGreeting call.
        string memory greeting = "bye";
        bytes memory encodedGreetingCall = bytes.concat(
            Greeter.setGreeting.selector,
            abi.encode(greeting)
        );

        // Produce the final full end-tx encoding, to be used as calldata in the user operation.
        bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
            address(greeter),
            uint256(0),
            1000,
            encodedGreetingCall
        );

        // Construct the user opeartion.
        UserOperation memory op = UserOperation({
            sender: toDeployAddress,
            nonce: 0,
            initCode: fullInitializeCode,
            callData: fullEncoding,
            callGasLimit: 1_000_000,
            verificationGasLimit: 500_000,
            preVerificationGas: 20_000,
            maxFeePerGas: 100,
            maxPriorityFeePerGas: 200,
            paymasterAndData: "",
            signature: signature2
        });

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes32 userOpHash = entryPoint.getUserOpHash(op);
        bytes32 fullHash = SmartContractAccountHelper.getFullHashForSigning(userOpHash);
        console.logBytes32(fullHash);

        // Deposit funds for the transaction.
        entryPoint.depositTo{value: 10 ether}(toDeployAddress);

        // Execute the user operation.
        UserOperation[] memory operations = new UserOperation[](1);
        operations[0] = op;
        entryPoint.handleOps(operations, payable(LINK_WHALE_2));

        // Assert that the greeting was set.
        assertEq("bye", Greeter(greeter).getGreeting());
    }

    /// @dev Test case for fresh user, EntryPoint.sol should generate a 
    /// @dev Smart Contract Account for them and execute the meta transaction.
    function testEIP712EIP4337AndCreateSmartContractAccountWithPaymaster() public {
        // Impersonate a different LINK whale.
        changePrank(LINK_WHALE);

        // Pre-calculate user smart contract account address.
        SmartContractAccountFactory factory = new SmartContractAccountFactory();
        address toDeployAddress = SmartContractAccountHelper.calculateSmartContractAccountAddress(LINK_WHALE_2, ENTRY_POINT, address(factory));

        // Construct initCode byte array.
        bytes memory fullInitializeCode = SmartContractAccountHelper.getInitCode(address(factory), LINK_WHALE_2, ENTRY_POINT);

        // Create the calldata for a setGreeting call.
        string memory greeting = "good day";
        bytes memory encodedGreetingCall = bytes.concat(
            Greeter.setGreeting.selector,
            abi.encode(greeting)
        );

        // Produce the final full end-tx encoding, to be used as calldata in the user operation.
        bytes memory fullEncoding = SmartContractAccountHelper.getFullEndTxEncoding(
            address(greeter),
            uint256(0),
            1000,
            encodedGreetingCall
        );

        Paymaster paymaster = new Paymaster();
        paymaster.deposit(1000 ether, toDeployAddress);

        // Construct the user opeartion.
        UserOperation memory op = UserOperation({
            sender: toDeployAddress,
            nonce: 0,
            initCode: fullInitializeCode,
            callData: fullEncoding,
            callGasLimit: 1_000_000,
            verificationGasLimit: 1_500_000,
            preVerificationGas: 20_000,
            maxFeePerGas: 100,
            maxPriorityFeePerGas: 200,
            paymasterAndData: abi.encodePacked(address(paymaster)),
            signature: signature3
        });

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes32 userOpHash = entryPoint.getUserOpHash(op);
        bytes32 fullHash = SmartContractAccountHelper.getFullHashForSigning(userOpHash);
        console.logBytes32(fullHash);

        // Deposit funds for the transaction.
        entryPoint.depositTo{value: 10 ether}(address(paymaster));

        // Execute the user operation.
        UserOperation[] memory operations = new UserOperation[](1);
        operations[0] = op;
        entryPoint.handleOps(operations, payable(LINK_WHALE_2));

        // Assert that the greeting was set.
        assertEq("good day", Greeter(greeter).getGreeting());
    }

}
