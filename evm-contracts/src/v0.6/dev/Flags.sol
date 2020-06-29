pragma solidity 0.6.6;


import "../Owned.sol";
import "../dev/SimpleAccessControl.sol";


/**
 * @title The Flags contract
 * @notice Allows flags to signal to any reader on the access control list.
 * The owner can set flags, or designate other addresses to set flags. The
 * owner must turn the flags off, other setters cannot.
 */
contract Flags is Owned, SimpleAccessControl {

  mapping(address => bool) private flags;
  mapping(address => bool) private setters;

  event FlagOn(
    address indexed subject
  );
  event FlagOff(
    address indexed subject
  );
  event SetterEnabled(
    address indexed setter
  );
  event SetterDisabled(
    address indexed setter
  );

  /**
   * @notice read the warning flag status of a contract address.
   * @param subject The contract address being checked for a warning.
   */
  function getFlag(address subject)
    public
    view
    checkAccess()
    returns (bool)
  {
    return flags[subject];
  }

  /**
   * @notice allows owner to enable the warning flags for mulitple addresses.
   * @param subjects List of the contract addresses whose flag is being raised
   */
  function setFlagsOn(address[] calldata subjects)
    external
    returns (bool)
  {
    require(msg.sender == owner || setters[msg.sender], "Only callable by enabled setters");

    for (uint256 i = 0; i < subjects.length; i++) {
      address subject = subjects[i];

      if (!flags[subject]) {
        flags[subject] = true;
        emit FlagOn(subject);
      }
    }
  }

  /**
   * @notice allows owner to disable the warning flags for mulitple addresses.
   * @param subjects List of the contract addresses whose flag is being lowered
   */
  function setFlagsOff(address[] calldata subjects)
    external
    onlyOwner()
    returns (bool)
  {
    for (uint256 i = 0; i < subjects.length; i++) {
      address subject = subjects[i];

      if (flags[subject]) {
        flags[subject] = false;
        emit FlagOff(subject);
      }
    }
  }

  /**
   * @notice allows owner to give other addresses permission to set flags on.
   * @param added List of the addresses of setters to be enabled.
   */
  function enableSetters(address[] calldata added)
    external
    onlyOwner()
    returns (bool)
  {
    for (uint256 i = 0; i < added.length; i++) {
      address setter = added[i];

      if (!setters[setter]) {
        setters[setter] = true;
        emit SetterEnabled(setter);
      }
    }
  }

  /**
   * @notice allows owner to remove addresses with permission to set flags on.
   * @param removed List of the addresses of setters to be enabled.
   */
  function disableSetters(address[] calldata removed)
    external
    onlyOwner()
    returns (bool)
  {
    for (uint256 i = 0; i < removed.length; i++) {
      address setter = removed[i];

      if (setters[setter]) {
        setters[setter] = false;
        emit SetterDisabled(setter);
      }
    }
  }

}
