// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ComposerCompatibleInterfaceV1 {
    error ComposerRequestV1(
        string scriptHash, // functions: keccack256 of Functions script to run
        string[] functionsArguments, // functions: additional arguments to be passed to the Functions DON
        bool useMercury, // mercury: determines if mercury data should be fetched
        string feedParamKey, // mercury: key to use for Mercury lookup
        string[] feeds, // mercury: feed Ids to fetch
        string timeParamKey, // mercury: type of time key to use
        uint256 time, // mercury: specific value of the time key
        bytes extraData // any extra data to be used
    );
}
