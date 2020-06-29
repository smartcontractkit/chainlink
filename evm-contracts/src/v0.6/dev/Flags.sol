pragma solidity 0.6.6;


import "../Owned.sol";
import "../dev/SimpleAccessControl.sol";


contract Flags is Owned, SimpleAccessControl {

  mapping(address => bool) private flags;

  event FlagOn(
    address indexed subject
  );
  event FlagOff(
    address indexed subject
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
   * @notice allows owner to enable the warning flag for an address.
   * @param subject The contract address being checked for a warning.
   */
  function setFlagOn(address subject)
    public
    onlyOwner()
    returns (bool)
  {
    if (!flags[subject]) {
      flags[subject] = true;
      emit FlagOn(subject);
    }
  }

  /**
   * @notice allows owner to disable the warning flag for an address.
   * @param subject The contract address being checked for a warning.
   */
  function setFlagOff(address subject)
    public
    onlyOwner()
    returns (bool)
  {
    if (flags[subject]) {
      flags[subject] = false;
      emit FlagOff(subject);
    }
  }

}
