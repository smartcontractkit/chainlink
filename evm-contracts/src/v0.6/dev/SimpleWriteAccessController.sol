pragma solidity ^0.6.0;

import "./SimpleReadAccessController.sol";

/**
 * @title SimpleWriteAccessController
 * @notice Allows the owner to set access for addresses. External accounts are
 * not granted special access.
 */
contract SimpleWriteAccessController is SimpleReadAccessController {

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
