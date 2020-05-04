pragma solidity ^0.6.0;

import "../Owned.sol";

/**
 * @title Whitelisted
 * @notice Allows the owner to add and remove addresses from a whitelist
 */
contract Whitelisted is Owned {

  bool public whitelistEnabled;
  mapping(address => bool) public whitelisted;

  event AddedToWhitelist(address user);
  event RemovedFromWhitelist(address user);
  event WhitelistEnabled();
  event WhitelistDisabled();

  constructor()
    public
  {
    whitelistEnabled = true;
  }

  /**
   * @notice Adds an address to the whitelist
   * @param _user The address to whitelist
   */
  function addToWhitelist(address _user) external onlyOwner() {
    whitelisted[_user] = true;
    emit AddedToWhitelist(_user);
  }

  /**
   * @notice Removes an address from the whitelist
   * @param _user The address to remove
   */
  function removeFromWhitelist(address _user) external onlyOwner() {
    delete whitelisted[_user];
    emit RemovedFromWhitelist(_user);
  }

  /**
   * @notice makes the whitelist check enforced
   */
  function enableWhitelist()
    external
    onlyOwner()
  {
    whitelistEnabled = true;

    emit WhitelistEnabled();
  }

  /**
   * @notice makes the whitelist check unenforced
   */
  function disableWhitelist()
    external
    onlyOwner()
  {
    whitelistEnabled = false;

    emit WhitelistDisabled();
  }

  /**
   * @dev reverts if the caller is not whitelisted
   */
  modifier isWhitelisted() {
    require(whitelisted[msg.sender] || !whitelistEnabled, "Not whitelisted");
    _;
  }
}
