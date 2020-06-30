pragma solidity 0.6.6;


import "./SimpleReadAccessController.sol";
import "./AccessControllerInterface.sol";


/**
 * @title The Flags contract
 * @notice Allows flags to signal to any reader on the access control list.
 * The owner can set flags, or designate other addresses to set flags. The
 * owner must turn the flags off, other setters cannot.
 */
contract Flags is SimpleReadAccessController {

  AccessControllerInterface public flaggingAccessController;

  mapping(address => bool) private flags;

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

  constructor(
    address acAddress
  )
    public
  {
    setFlaggingAccessController(acAddress);
  }

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
  {
    require(msg.sender == owner ||
      flaggingAccessController.hasAccess(msg.sender, msg.data),
      "No access");

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
   * @notice allows owner to change the access controller for flagging addresses.
   * @param acAddress new address for flagging access controller.
   */
  function setFlaggingAccessController(
    address acAddress
  )
    public
    onlyOwner()
  {
    flaggingAccessController = AccessControllerInterface(acAddress);
  }

}
