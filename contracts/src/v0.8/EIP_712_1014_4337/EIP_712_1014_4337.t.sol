// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import "forge-std/Test.sol";
import "./consumers/Create2Factory.sol";
import "./consumers/SCA.sol";
import "./consumers/Greeter.sol";
import "./interfaces/UserOperation.sol";
import "./contracts/EntryPoint.sol";
import "./interfaces/IEntryPoint.sol";



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
+---------------------------------------------------------------------------------------------------------------------*/

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
        hex"60a060405234801561001057600080fd5b5060405161075c38038061075c83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b6080516106c46100986000396000818160a00152818161011a015261023501526106c46000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633a871cdd146100515780636ca0f814146100775780638b6a5f2b1461008c578063e39782401461009b575b600080fd5b61006461005f366004610439565b6100da565b6040519081526020015b60405180910390f35b61008a6100853660046104a9565b61031a565b005b6040516001815260200161006e565b6100c27f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161006e565b6000805481806100e983610542565b90915550600090506100fe602086018661055b565b61010b606087018761057d565b600054604051610144949392917f00000000000000000000000000000000000000000000000000000000000000009146906020016105cb565b60408051601f198184030181529190528051602090910120905060006101ab61017161014088018861057d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525092506103d2915050565b905060006101fd6101c061014089018961057d565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250602092506103d2915050565b9050600061020f61014089018961057d565b604081811061022057610220610625565b919091013560f81c9150506001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001660018561026384601b61063b565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa1580156102b2573d6000803e3d6000fd5b505050602060405103516001600160a01b03161461030c5760405162461bcd60e51b815260206004820152601260248201527124b73b30b634b21039b4b3b730ba3ab9329760711b60448201526064015b60405180910390fd5b506000979650505050505050565b33730576a174d229e3cfa37253523e645a78a0c91b571461036e5760405162461bcd60e51b815260206004820152600e60248201526d1b9bdd08185d5d1a1bdc9a5e995960921b6044820152606401610303565b826001600160a01b03168282604051610388929190610654565b6000604051808303816000865af19150503d80600081146103c5576040519150601f19603f3d011682016040523d82523d6000602084013e6103ca565b606091505b505050505050565b60008060005b602081101561042f576103ec816008610664565b856103f7838761067b565b8151811061040757610407610625565b01602001516001600160f81b031916901c91909117908061042781610542565b9150506103d8565b5090505b92915050565b60008060006060848603121561044e57600080fd5b833567ffffffffffffffff81111561046557600080fd5b8401610160818703121561047857600080fd5b95602085013595506040909401359392505050565b80356001600160a01b03811681146104a457600080fd5b919050565b6000806000604084860312156104be57600080fd5b6104c78461048d565b9250602084013567ffffffffffffffff808211156104e457600080fd5b818601915086601f8301126104f857600080fd5b81358181111561050757600080fd5b87602082850101111561051957600080fd5b6020830194508093505050509250925092565b634e487b7160e01b600052601160045260246000fd5b6000600182016105545761055461052c565b5060010190565b60006020828403121561056d57600080fd5b6105768261048d565b9392505050565b6000808335601e1984360301811261059457600080fd5b83018035915067ffffffffffffffff8211156105af57600080fd5b6020019150368190038213156105c457600080fd5b9250929050565b6001600160a01b03878116825260a0602083018190528201869052600090868860c0850137600060c0888501810191909152951660408301525060608101929092526080820152601f909201601f19169091010192915050565b634e487b7160e01b600052603260045260246000fd5b60ff81811683821601908111156104335761043361052c565b8183823760009101908152919050565b80820281158282048414176104335761043361052c565b808201808211156104335761043361052c56fea26469706673582212203e005acbad0d577c86296ff90d5c9a9253352072fb338216ca56b3da433f65fb64736f6c63430008110033";

    bytes signature = // Signature by LINK_WHALE. Signs off on "hi" being set as the greeting on the Greeter.sol contract, via their SCA.
        hex"cc503bc5125e8bc43fdeb28931df74c718b433f1b0bec332ffa90bb7be71425b2cc23dc7174253ee4fe29107b38cb589c4fe97ada812381b41b284a0e24c862101";

    bytes signature2 = // Signature by LINK_WHALE_2. Signs off on "bye" being set as the greeting on the Greeter.sol contract, via their SCA.
        hex"8028020b46fc90471b8ed7f2a77ada1ea6ccf6a3be2a153320c0f721b12b8964767554f49ffae39c14c71589b7d4fc4db7919d8a2126b0d551e201fccd4cd61901";

    function setUp() public {
        // Fork Goerli.
        uint256 mainnetFork = vm.createFork(
            "ETH_GOERLI_RPC_URL"
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
            abi.encode(LINK_WHALE)
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
        assertEq(true, SCA(toDeployAddress).isSCA());
        assertEq(LINK_WHALE, SCA(toDeployAddress).s_owner());

        // Create the calldata for a setGreeting call.
        string memory greeting = "hi";
        bytes memory encodedGreetingCall = abi.encodeWithSelector(
            Greeter.setGreeting.selector,
            greeting
        );

        // For developers: log the final hash of the SCA call to easily produce a signature off-chain.
        bytes memory fullEncoding = abi.encodeWithSelector(
            SCA.executeTransaction.selector,
            address(greeter),
            encodedGreetingCall
        );
        bytes32 hashOfEncoding = keccak256(
            abi.encode(
                toDeployAddress,
                fullEncoding,
                LINK_WHALE,
                uint256(1), // nonce++
                block.chainid
            )
        );
        console.logBytes32(hashOfEncoding);

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
            abi.encode(LINK_WHALE_2)
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
            SCA.executeTransaction.selector,
            address(greeter),
            encodedGreetingCall
        );
        bytes32 hashOfEncoding = keccak256(
            abi.encode(
                toDeployAddress,
                fullEncoding,
                LINK_WHALE_2,
                uint256(1), // nonce++
                block.chainid
            )
        );
        console.logBytes32(hashOfEncoding);

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
