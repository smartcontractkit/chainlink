pragma solidity ^0.8.0;

/**
 * @title Create2 factory contract
 * @author 0age - modified by vreff for use with EIP-4337.
 */
contract Create2Factory {
    event ContractCreated(address);

    /**
     * @dev Create a contract using CREATE2 by submitting a given salt or nonce
     * along with the initialization code for the contract.
     * @param salt bytes32 The nonce that will be passed into the CREATE2 call.
     * @param initializationCode bytes The initialization code that will be passed
     * into the CREATE2 call.
     * @return deploymentAddress of the contract that will be created, or the null address
     * if a contract already exists at that address.
     */
    function callCreate2(bytes32 salt, bytes calldata initializationCode)
        external
        payable
        returns (address deploymentAddress)
    {
        // move the initialization code from calldata to memory, and append the entry point contract address.
        bytes memory initCode = initializationCode;

        // using inline assembly: load data and length of data, then call CREATE2.
        assembly {
            // solhint-disable-line
            let encoded_data := add(0x20, initCode) // load initialization code.
            let encoded_size := mload(initCode) // load the init code's length.
            deploymentAddress := create2(
                // call CREATE2 with 4 arguments.
                callvalue(), // forward any attached value.
                encoded_data, // pass in initialization code.
                encoded_size, // pass in init code's length.
                salt // pass in the salt value.
            )
        }

        emit ContractCreated(deploymentAddress);

        // ensure that the contract address is not equal to the null address.
        require(
            deploymentAddress != address(0),
            "Failed to deploy contract using provided salt and initialization code."
        );
    }

    /**
     * @dev Compute the address of the contract that will be created when
     * submitting a given salt or nonce to the contract. The CREATE2 address is
     * computed in accordance with EIP-1014, and adheres to the formula therein of
     * `keccak256( 0xff ++ address ++ salt ++ keccak256(init_code)))[12:]` when
     * performing the computation. The computed address is then checked for any
     * existing contract code - if so, the null address will be returned instead.
     * @param salt bytes32 The nonce passed into the CREATE2 address calculation.
     * @param initCodeHash bytes32 The keccak256 hash of the initialization code
     * that will be passed into the CREATE2 address calculation.
     * @return deploymentAddress of the contract that will be created, or the null address
     * if a contract already exists at that address.
     */
    function findCreate2Address(bytes32 salt, bytes32 initCodeHash)
        external
        view
        returns (address deploymentAddress)
    {
        // variable for checking code size of any pre-existing contract at address.
        uint256 existingContractSize;

        // determine the address where the contract will be deployed.
        deploymentAddress = address(
            uint160( // downcast to match the address type.
                uint256( // convert to uint to truncate upper digits.
                    keccak256( // compute the CREATE2 hash using 4 inputs.
                        abi.encodePacked( // pack all inputs to the hash together.
                            hex"ff", // start with 0xff to distinguish from RLP.
                            address(this), // this contract will be the caller.
                            salt, // pass in the supplied salt value.
                            initCodeHash // pass in the hash of initialization code.
                        )
                    )
                )
            )
        );

        // determine if any contract code already exists at the computed address.
        assembly {
            // solhint-disable-line
            existingContractSize := extcodesize(deploymentAddress)
        }

        // if so, return null address to signify failure. (detect selfdestructed?)
        if (existingContractSize > 0) {
            return address(0);
        }
    }
}
