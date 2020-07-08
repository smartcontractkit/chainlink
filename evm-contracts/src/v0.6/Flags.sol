pragma solidity ^0.6.0;


import "./SimpleReadAccessController.sol";
import "./interfaces/AccessControllerInterface.sol";
import "./interfaces/FlagsInterface.sol";


/**
 * @title The Flags contract
 * @notice Allows flags to signal to any reader on the access control list.
 * The owner can set flags, or designate other addresses to set flags. The
 * owner must turn the flags off, other setters cannot.
 */
contract Flags is FlagsInterface, SimpleReadAccessController {

  AccessControllerInterface public raisingAccessController;

  mapping(address => bool) private flags;

  event FlagOn(
    address indexed subject
  );
  event FlagOff(
    address indexed subject
  );
  event RaisingAccessControllerChanged(
    address indexed previous,
    address indexed current
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
   * @return A true value indicates that a flag was raised and a
   * false value indicates that no flag was raised.
   */
  function getFlag(address subject)
    external
    view
    override
    checkAccess()
    returns (bool)
  {
    return flags[subject];
  }

  /**
   * @notice read the warning flag status of a contract address.
   * @param subjects An array of addresses being checked for a flag.
   * @return An array of bools where a true value for any flag indicates that
   * a flag was raised and a false value indicates that no flag was raised.
   */
  function getFlags(address[] calldata subjects)
    external
    view
    override
    checkAccess()
    returns (bool[] memory)
  {
    bool[] memory responses = new bool[](subjects.length);
    for (uint256 i = 0; i < subjects.length; i++) {
      responses[i] = flags[subjects[i]];
    }
    return responses;
  }

  /**
   * @notice enable the warning flags for mulitple addresses.
   * Access is controlled by raisingAccessController, except for owner
   * who always has access.
   * @param subjects List of the contract addresses whose flag is being raised
   */
  function raiseFlags(address[] calldata subjects)
    external
    override
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
    override
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
   * @notice allows owner to change the access controller for raising flags.
   * @param racAddress new address for the raising access controller.
   */
  function setRaisingAccessController(
    address racAddress
  )
    public
    override
    onlyOwner()
  {
    address previous = address(raisingAccessController);

    raisingAccessController = AccessControllerInterface(racAddress);

    emit RaisingAccessControllerChanged(previous, racAddress);
  }


  // PRIVATE

  function allowedToRaiseFlags()
    private
    view
    returns (bool)
  {
    return msg.sender == owner ||
      raisingAccessController.hasAccess(msg.sender, msg.data);
  }

}
