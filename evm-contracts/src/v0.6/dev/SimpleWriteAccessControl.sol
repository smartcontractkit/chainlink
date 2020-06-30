pragma solidity ^0.6.0;

import "./SimpleAccessControl.sol";

/**
 * @title SimpleWriteAccessControl
 * @notice Allows the owner to set access for addresses. External accounts are
 * not granted special access.
 */
contract SimpleWriteAccessControl is SimpleAccessControl {

  /**
   * @notice Returns the access of an address
   * @param _user The address to query
   */
  function hasAccess(
    address _user,
    bytes memory
  )
    public
    view
    override
    returns (bool)
  {
    return accessList[_user] || !checkEnabled;
  }

}
