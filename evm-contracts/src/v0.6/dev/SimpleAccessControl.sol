pragma solidity ^0.6.0;

import "../Owned.sol";
import "./AccessController.sol";

/**
 * @title SimpleAccessControl
 * @notice Allows the owner to set access for addresses
 */
contract SimpleAccessControl is AccessController, Owned {

  bool public checkEnabled;
  mapping(address => bool) internal accessList;

  event AddedAccess(address user);
  event RemovedAccess(address user);
  event CheckAccessEnabled();
  event CheckAccessDisabled();

  constructor()
    public
  {
    checkEnabled = true;
  }

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
    return accessList[_user];
  }

  /**
   * @notice Adds an address to the access list
   * @param _user The address to add
   */
  function addAccess(address _user)
    external
    onlyOwner()
  {
    accessList[_user] = true;
    emit AddedAccess(_user);
  }

  /**
   * @notice Removes an address from the access list
   * @param _user The address to remove
   */
  function removeAccess(address _user)
    external
    onlyOwner()
  {
    delete accessList[_user];
    emit RemovedAccess(_user);
  }

  /**
   * @notice makes the access check enforced
   */
  function enableAccessCheck()
    external
    onlyOwner()
  {
    checkEnabled = true;

    emit CheckAccessEnabled();
  }

  /**
   * @notice makes the access check unenforced
   */
  function disableAccessCheck()
    external
    onlyOwner()
  {
    checkEnabled = false;

    emit CheckAccessDisabled();
  }

  /**
   * @dev reverts if the caller does not have access
   */
  modifier checkAccess() virtual override {
    require(hasAccess(msg.sender, msg.data) || !checkEnabled, "No access");
    _;
  }
}
