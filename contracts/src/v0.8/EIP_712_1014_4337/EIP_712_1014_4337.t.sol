// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "./consumers/Create2Factory.sol";
import "./consumers/SCA.sol";
import "./consumers/Greeter.sol";
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

    bytes initailizeCode = // Bytecode for SCA.sol (Smart Contract Account). Does not include constructor arguments.
        hex"60c060405234801561001057600080fd5b5060405161088138038061088183398101604081905261002f91610062565b6001600160a01b039182166080521660a052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a0516107b46100cd600039600081816056015261040401526000818160dd0152818161013f01526102bf01526107b46000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80631a4e75de146100515780633a871cdd146100a25780639e045e45146100c3578063e3978240146100d8575b600080fd5b6100787f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b6100b56100b03660046104fe565b6100ff565b604051908152602001610099565b6100d66100d1366004610552565b6103ec565b005b6100787f000000000000000000000000000000000000000000000000000000000000000081565b6000807f23d294a3e6e5266ba3b17997d1a601816663066087b37cad4f06b4d4e30655d961013060608701876105f4565b600054604051610169949392917f0000000000000000000000000000000000000000000000000000000000000000914690602001610660565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201207f1900000000000000000000000000000000000000000000000000000000000000828501527f010000000000000000000000000000000000000000000000000000000000000060218501527f47e79534a245952e8b16893a336b85a3d9ea9fa8c573f3d803afb92a7946921860228501526042808501829052835180860390910181526062909401909252825192019190912090915060006102456101408801886105f4565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920182905250602085015160408087015187519798509196919550919350869250811061029e5761029e6106e2565b016020015160f81c905073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000166001866102ed84601b610740565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa15801561033c573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff16146103c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e76616c6964207369676e61747572652e000000000000000000000000000060448201526064015b60405180910390fd5b6000805490806103d78361075f565b9091555060009b9a5050505050505050505050565b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461048b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6e6f7420617574686f72697a656400000000000000000000000000000000000060448201526064016103bf565b8373ffffffffffffffffffffffffffffffffffffffff168383836040516104b3929190610797565b60006040518083038185875af1925050503d80600081146104f0576040519150601f19603f3d011682016040523d82523d6000602084013e6104f5565b606091505b50505050505050565b60008060006060848603121561051357600080fd5b833567ffffffffffffffff81111561052a57600080fd5b8401610160818703121561053d57600080fd5b95602085013595506040909401359392505050565b6000806000806060858703121561056857600080fd5b843573ffffffffffffffffffffffffffffffffffffffff8116811461058c57600080fd5b935060208501359250604085013567ffffffffffffffff808211156105b057600080fd5b818701915087601f8301126105c457600080fd5b8135818111156105d357600080fd5b8860208285010111156105e557600080fd5b95989497505060200194505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261062957600080fd5b83018035915067ffffffffffffffff82111561064457600080fd5b60200191503681900382131561065957600080fd5b9250929050565b86815260a060208201528460a0820152848660c0830137600060c08683010152600060c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f880116830101905073ffffffffffffffffffffffffffffffffffffffff85166040830152836060830152826080830152979650505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60ff818116838216019081111561075957610759610711565b92915050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361079057610790610711565b5060010190565b818382376000910190815291905056fea164736f6c6343000811000a";

    bytes signature = // Signature by LINK_WHALE. Signs off on "hi" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"fb78470345e472d2362a178a557ed060f3d94cece6d3557862a185ef1317f6ee65c95ca7c9c6a1dfd50899410a7fc5c5c6f3ccc0c0ccfcfc18dce70542d2215700";

    bytes signature2 = // Signature by LINK_WHALE_2. Signs off on "bye" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"f027ddd24c5bd68c87f3dc665e9e51f39fc246d80d00948d4fb0bdcdd2de722f6e8bbc920d0bc746ea2070498a6f37928ae9831f2790136e59b3de9d9f6ca51701";

    function setUp() public {
        // Fork Goerli.
        uint256 mainnetFork = vm.createFork(
            "ETH_RPC_URL_GOERLI"
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
        // Use forked Create2Factory contract.
        Create2Factory factory = Create2Factory(CREATE2_FACTORY);

        // Pre-calculate user smart contract account address.
        bytes32 salt = bytes32(uint256(uint160(LINK_WHALE)) << 96);
        bytes memory fullInitializeCode = bytes.concat(
            initailizeCode,
            abi.encode(LINK_WHALE, ENTRY_POINT)
        );
        bytes32 initializeCodeHash = keccak256(fullInitializeCode);
        address toDeployAddress = factory.findCreate2Address(
            salt,
            initializeCodeHash
        );

        // Deploy the end-contract.
        changePrank(SENDER_CREATOR);
        factory.callCreate2(salt, fullInitializeCode);
        changePrank(LINK_WHALE);

        // Ensure a correct deployment and a functioning end-contract.
        address deployedAddress = factory.findCreate2Address(
            salt,
            initializeCodeHash
        );
        assertEq(deployedAddress, address(0));
        assertEq(LINK_WHALE, SCA(toDeployAddress).s_owner());

        // Create the calldata for a setGreeting call.
        string memory greeting = "hi";
        bytes memory encodedGreetingCall = abi.encodeWithSelector(
            Greeter.setGreeting.selector,
            greeting
        );

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes memory fullEncoding = abi.encodeWithSelector(
            SCA.executeTransactionFromEntryPoint.selector,
            address(greeter),
            uint256(0),
            encodedGreetingCall
        );
        bytes32 hashOfEncoding = keccak256(
            abi.encode(
                SCALibrary.TYPEHASH,
                fullEncoding,
                LINK_WHALE,
                uint256(0),
                block.chainid
            )
        );
        bytes32 fullHash = keccak256(abi.encodePacked(bytes1(0x19), bytes1(0x01), SCALibrary.DOMAIN_SEPARATOR, hashOfEncoding));
        console.logBytes32(fullHash);

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
        bytes32 salt = bytes32(uint256(uint160(LINK_WHALE_2)) << 96);
        bytes memory initializeCodeWithConstructor = bytes.concat(
            initailizeCode,
            abi.encode(LINK_WHALE_2, ENTRY_POINT)
        );
        bytes32 initializeCodeHash = keccak256(initializeCodeWithConstructor);
        Create2Factory factory = Create2Factory(CREATE2_FACTORY);
        address toDeployAddress = factory.findCreate2Address(
            salt,
            initializeCodeHash
        );

        // Construct initCode byte array.
        bytes memory fullInitializeCode = bytes.concat(
            bytes20(CREATE2_FACTORY),
            abi.encodeWithSelector(
                Create2Factory.callCreate2.selector,
                salt,
                initializeCodeWithConstructor
            )
        );

        // Create the calldata for a setGreeting call.
        string memory greeting = "bye";
        bytes memory encodedGreetingCall = abi.encodeWithSelector(
            Greeter.setGreeting.selector,
            greeting
        );

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes memory fullEncoding = abi.encodeWithSelector(
            SCA.executeTransactionFromEntryPoint.selector,
            address(greeter),
            uint256(0),
            encodedGreetingCall
        );

        bytes32 hashOfEncoding = keccak256(
            abi.encode(
                SCALibrary.TYPEHASH,
                fullEncoding,
                LINK_WHALE_2,
                uint256(0),
                block.chainid
            )
        );
        bytes32 fullHash = keccak256(abi.encodePacked(bytes1(0x19), bytes1(0x01), SCALibrary.DOMAIN_SEPARATOR, hashOfEncoding));
        console.logBytes32(fullHash);

        // Construct the user opeartion.
        UserOperation memory op = UserOperation({
            sender: toDeployAddress,
            nonce: 0,
            initCode: fullInitializeCode,
            callData: fullEncoding,
            callGasLimit: 10_000_000,
            verificationGasLimit: 10_000_000,
            preVerificationGas: 10_000_000,
            maxFeePerGas: 100,
            maxPriorityFeePerGas: 200,
            paymasterAndData: "",
            signature: signature2
        });

        // Deposit funds for the transaction.
        entryPoint.depositTo{value: 10 ether}(toDeployAddress);

        // Execute the user operation.
        UserOperation[] memory operations = new UserOperation[](1);
        operations[0] = op;
        entryPoint.handleOps(operations, payable(LINK_WHALE_2));

        // Assert that the greeting was set.
        assertEq("bye", Greeter(greeter).getGreeting());
    }
}

/*----------------------------+
| HELPER FUNCTIONS            |
| ________________            |
|                             |
+----------------------------*/

function bytesToBytes32(bytes memory b, uint256 offset) pure returns (bytes32) {
    bytes32 out;

    for (uint256 i = 0; i < 32; i++) {
        out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
    }
    return out;
}
