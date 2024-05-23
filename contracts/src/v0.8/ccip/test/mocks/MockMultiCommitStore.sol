// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMultiCommitStore} from "../../interfaces/IMultiCommitStore.sol";

contract MockMultiCommitStore is IMultiCommitStore {
  error PausedError();

  bool private s_paused = false;
  mapping(uint64 sourceChainSelector => bool shouldVerify) private s_verifiedSourceChains;
  mapping(uint64 sourceChainSelector => IMultiCommitStore.SourceChainConfig config) s_sourceChainConfigs;

  function getSourceChainConfig(uint64 sourceChainSelector)
    external
    view
    returns (IMultiCommitStore.SourceChainConfig memory)
  {
    return s_sourceChainConfigs[sourceChainSelector];
  }

  function setSourceChainConfig(
    uint64 sourceChainSelector,
    IMultiCommitStore.SourceChainConfig memory sourceChainConfig
  ) external {
    s_sourceChainConfigs[sourceChainSelector] = sourceChainConfig;
  }

  /// @inheritdoc IMultiCommitStore
  function verify(
    uint64 sourceChainSelector,
    bytes32[] calldata,
    bytes32[] calldata,
    uint256
  ) external view whenNotPaused returns (uint256 timestamp) {
    return s_verifiedSourceChains[sourceChainSelector] ? 1 : 0;
  }

  function setVerifyResult(uint64 sourceChainSelector, bool shouldVerify) external {
    s_verifiedSourceChains[sourceChainSelector] = shouldVerify;
  }

  modifier whenNotPaused() {
    if (paused()) revert PausedError();
    _;
  }

  function paused() public view returns (bool) {
    return s_paused;
  }

  function pause() external {
    s_paused = true;
  }
}
