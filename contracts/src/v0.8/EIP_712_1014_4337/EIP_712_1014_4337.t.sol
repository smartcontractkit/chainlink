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
        hex"60a060405234801561001057600080fd5b5060405161084a38038061084a83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b6080516107b26100986000396000818160a001528181610119015261025f01526107b26000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633a871cdd146100515780636ca0f814146100775780638b6a5f2b1461008c578063e39782401461009b575b600080fd5b61006461005f3660046104e3565b6100e7565b6040519081526020015b60405180910390f35b61008a610085366004610537565b610376565b005b6040516001815260200161006e565b6100c27f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161006e565b6000805481806100f683610604565b909155506000905061010b606086018661063c565b6000546040516101439392917f00000000000000000000000000000000000000000000000000000000000000009146906020016106a8565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190528051602090910120905060006101c861018e61014088018861063c565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052509250610464915050565b9050600061021a6101dd61014089018961063c565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060209250610464915050565b9050600061022c61014089018961063c565b604081811061023d5761023d610723565b919091013560f81c91505073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001660018561028d84601b610752565b6040805160008152602081018083529390935260ff90911690820152606081018690526080810185905260a0016020604051602081039080840390855afa1580156102dc573d6000803e3d6000fd5b5050506020604051035173ffffffffffffffffffffffffffffffffffffffff1614610368576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f496e76616c6964207369676e61747572652e000000000000000000000000000060448201526064015b60405180910390fd5b506000979650505050505050565b33730576a174d229e3cfa37253523e645a78a0c91b57146103f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600e60248201527f6e6f7420617574686f72697a6564000000000000000000000000000000000000604482015260640161035f565b8273ffffffffffffffffffffffffffffffffffffffff16828260405161041a92919061076b565b6000604051808303816000865af19150503d8060008114610457576040519150601f19603f3d011682016040523d82523d6000602084013e61045c565b606091505b505050505050565b60008060005b60208110156104d95761047e81600861077b565b856104898387610792565b8151811061049957610499610723565b01602001517fff0000000000000000000000000000000000000000000000000000000000000016901c9190911790806104d181610604565b91505061046a565b5090505b92915050565b6000806000606084860312156104f857600080fd5b833567ffffffffffffffff81111561050f57600080fd5b8401610160818703121561052257600080fd5b95602085013595506040909401359392505050565b60008060006040848603121561054c57600080fd5b833573ffffffffffffffffffffffffffffffffffffffff8116811461057057600080fd5b9250602084013567ffffffffffffffff8082111561058d57600080fd5b818601915086601f8301126105a157600080fd5b8135818111156105b057600080fd5b8760208285010111156105c257600080fd5b6020830194508093505050509250925092565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610635576106356105d5565b5060010190565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261067157600080fd5b83018035915067ffffffffffffffff82111561068c57600080fd5b6020019150368190038213156106a157600080fd5b9250929050565b60808152846080820152848660a0830137600060a08683010152600060a07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f880116830101905073ffffffffffffffffffffffffffffffffffffffff851660208301528360408301528260608301529695505050505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60ff81811683821601908111156104dd576104dd6105d5565b8183823760009101908152919050565b80820281158282048414176104dd576104dd6105d5565b808201808211156104dd576104dd6105d556fea164736f6c6343000811000a";

    bytes signature = // Signature by LINK_WHALE. Signs off on "hi" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"b4c3671358afba4454e1de3cf0865d496b92d77be42f8dbecf178a3f05e73d1e42ae8bc1380b34fdefaee2883edeeac001ef35be2720360ff1ec72f71ba0b90400";

    bytes signature2 = // Signature by LINK_WHALE_2. Signs off on "bye" being set as the greeting on the Greeter.sol contract, without knowing the address of their SCA.
        hex"b1bfcf380055c4b8e30f8c1032d5ec8f09c31633d58f27aec84000e77f7fd86b69e125f4e009415988c62397567c65dd87a36bce398788a5336ec91feabb21eb01";

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
