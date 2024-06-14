// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/introspection/IERC165.sol";

/// @title IOptimismMintableERC20Minimal
/// @dev This interface is a subset of the Optimism ERC20 interface that is defined
/// below. This is done to now have to overwrite the burn and mint functions again in
/// the implementation, as that leads to more complicated, error prone code.
interface IOptimismMintableERC20Minimal is IERC165 {
  /// @notice Returns the address of the token on L1.
  function remoteToken() external view returns (address);

  /// @notice Returns the address of the bridge on L2.
  function bridge() external returns (address);
}

/// @title IOptimismMintableERC20
/// @notice This is the complete interface for the Optimism mintable ERC20 token as defined in
/// https://github.com/ethereum-optimism/optimism/blob/develop/packages/contracts-bedrock/contracts/universal/IOptimismMintableERC20.sol
interface IOptimismMintableERC20 is IERC165, IOptimismMintableERC20Minimal {
  function mint(address _to, uint256 _amount) external;

  function burn(address _from, uint256 _amount) external;
}
