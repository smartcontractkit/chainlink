// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ChainSpecificUtil} from "../../shared/util/ChainSpecificUtil.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {BlockhashStore} from "./BlockhashStore.sol";

contract TrustedBlockhashStore is ConfirmedOwner, BlockhashStore {
  error NotInWhitelist();
  error InvalidTrustedBlockhashes();
  error InvalidRecentBlockhash();

  mapping(address => bool) public s_whitelistStatus;
  address[] public s_whitelist;

  constructor(address[] memory whitelist) ConfirmedOwner(msg.sender) {
    setWhitelist(whitelist);
  }

  /**
   * @notice sets the whitelist of addresses that can store blockhashes
   * @param whitelist the whitelist of addresses that can store blockhashes
   */
  function setWhitelist(address[] memory whitelist) public onlyOwner {
    address[] memory previousWhitelist = s_whitelist;
    s_whitelist = whitelist;

    // Unset whitelist status for all addresses in the previous whitelist,
    // and set whitelist status for all addresses in the new whitelist.
    for (uint256 i = 0; i < previousWhitelist.length; i++) {
      s_whitelistStatus[previousWhitelist[i]] = false;
    }
    for (uint256 i = 0; i < whitelist.length; i++) {
      s_whitelistStatus[whitelist[i]] = true;
    }
  }

  /**
   * @notice stores a list of blockhashes and their respective blocks, only callable
   * by a whitelisted address
   * @param blockhashes the list of blockhashes and their respective blocks
   */
  function storeTrusted(
    uint256[] calldata blockNums,
    bytes32[] calldata blockhashes,
    uint256 recentBlockNumber,
    bytes32 recentBlockhash
  ) external {
    bytes32 onChainHash = ChainSpecificUtil._getBlockhash(uint64(recentBlockNumber));
    if (onChainHash != recentBlockhash) {
      revert InvalidRecentBlockhash();
    }

    if (!s_whitelistStatus[msg.sender]) {
      revert NotInWhitelist();
    }

    if (blockNums.length != blockhashes.length) {
      revert InvalidTrustedBlockhashes();
    }

    for (uint256 i = 0; i < blockNums.length; i++) {
      s_blockhashes[blockNums[i]] = blockhashes[i];
    }
  }
}
