pragma solidity ^0.6.0;

import "./AccessControllerInterface.sol";

/**
 * @title Controller
 * @notice Allows inheriting contracts to control access to functions
 */
abstract contract AccessController is AccessControllerInterface {

  /**
   * @notice Returns the access of an address
   * @param _user The address to query
   */
  function hasAccess(address _user, bytes memory _data) public view virtual override returns (bool);

  /**
   * @dev reverts if the caller does not have access
   */
  modifier checkAccess() virtual {
    require(hasAccess(msg.sender, msg.data), "No Access");
    _;
  }
}
