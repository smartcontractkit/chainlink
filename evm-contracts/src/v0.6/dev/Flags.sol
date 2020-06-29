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
   * @notice allows owner to enable the warning flags for mulitple addresses.
   * @param subjects List of the contract addresses whose flag is being raised
   */
  function setFlagsOn(address[] calldata subjects)
    external
    onlyOwner()
    returns (bool)
  {
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

}
