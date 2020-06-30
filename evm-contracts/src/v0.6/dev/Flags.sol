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

  AccessControllerInterface public raisingAccessController;

  mapping(address => bool) private flags;

  event FlagOn(
    address indexed subject
  );
  event FlagOff(
    address indexed subject
  );

  /**
   * @param racAddress address for the raising access controller.
   */
  constructor(
    address racAddress
  )
    public
  {
    setRaisingAccessController(racAddress);
  }

  /**
   * @notice read the warning flag status of a contract address.
   * @param subject The contract address being checked for a flag.
   * A true value indicates that a flag was raised and a false value
   * indicates that no flag was raised.
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
  function raiseFlags(address[] calldata subjects)
    external
  {
    require(allowedToRaiseFlags(), "Not allowed to raise flags");

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
  function lowerFlags(address[] calldata subjects)
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
   * @param racAddress new address for the raising access controller.
   */
  function setRaisingAccessController(
    address racAddress
  )
    public
    onlyOwner()
  {
    raisingAccessController = AccessControllerInterface(racAddress);
  }


  // PRIVATE

  function allowedToRaiseFlags()
    private
    returns (bool)
  {
    return msg.sender == owner ||
      raisingAccessController.hasAccess(msg.sender, msg.data);
  }

}
