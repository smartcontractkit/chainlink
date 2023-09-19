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

   /**
   * @notice any contract which wants to utilize ComposerRequestV1 feature needs to
   * implement this interface as well as the automation compatible interface.
   * @param values an array containing the abi encoded result of a Functions call.
   * @param extraData context data from composer process.
   * @return upkeepNeeded boolean to indicate whether the keeper should call performUpkeep or not.
   * @return performData bytes that the keeper should call performUpkeep with, if
   * upkeep is needed. If you would like to encode data to decode later, try `abi.encode`.
   */
  function checkCallback(
    bytes[] memory values,
    bytes memory extraData
  ) external view returns (bool upkeepNeeded, bytes memory performData);
}
