// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;


import "./SimpleReadAccessController.sol";
import "./interfaces/AccessControllerInterface.sol";
import "./interfaces/FlagsInterface.sol";
import "./interfaces/TypeAndVersionInterface.sol";


/**
 * @title The Flags contract
 * @notice Allows flags to signal to any reader on the access control list.
 * The owner can set flags, or designate other addresses to set flags. The
 * owner must turn the flags off, other setters cannot. An expected pattern is
 * to allow addresses to raise flags on themselves, so if you are subscribing to
 * FlagOn events you should filter for addresses you care about.
 */
contract Flags is FlagsInterface, SimpleReadAccessController, TypeAndVersionInterface {

  AccessControllerInterface public raisingAccessController;
  AccessControllerInterface public loweringAccessController;

  mapping(address => bool) private flags;

  event FlagRaised(
    address indexed subject
  );
  event FlagLowered(
    address indexed subject
  );
  event RaisingAccessControllerUpdated(
    address indexed previous,
    address indexed current
  );
  event LoweringAccessControllerUpdated(
    address indexed previous,
    address indexed current
  );

  /**
   * @param racAddress address for the raising access controller.
   * @param lacAddress address for the lowering access controller.
   */
  constructor(
    address racAddress,
    address lacAddress
  )
    public
  {
    setRaisingAccessController(racAddress);
    setLoweringAccessController(lacAddress);
  }

  function typeAndVersion()
    external
    pure
    virtual
    override
    returns (
      string memory
    )
  {
    return "Flags 2.0.0";
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
   * @notice enable the warning flag for an address.
   * Access is controlled by raisingAccessController, except for owner
   * who always has access.
   * @param subject The contract address whose flag is being raised
   */
  function raiseFlag(address subject)
    external
    override
  {
    require(allowedToRaiseFlags(), "Not allowed to raise flags");

    tryToRaiseFlag(subject);
  }

  /**
   * @notice enable the warning flags for multiple addresses.
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
      tryToRaiseFlag(subjects[i]);
    }
  }

  /**
   * @notice allows owner to disable the warning flags for an addresses.
   * Access is controlled by loweringAccessController, except for owner
   * who always has access.
   * @param subject The contract address whose flag is being lowered
   */
  function lowerFlag(address subject)
    external
    override
  {
    require(allowedToLowerFlags(), "Not allowed to lower flags");

    if (flags[subject]) {
      flags[subject] = false;
      emit FlagLowered(subject);
    }
  }

  /**
   * @notice allows owner to disable the warning flags for multiple addresses.
   * Access is controlled by loweringAccessController, except for owner
   * who always has access.
   * @param subjects List of the contract addresses whose flag is being lowered
   */
  function lowerFlags(address[] calldata subjects)
    external
    override
  {
    require(allowedToLowerFlags(), "Not allowed to lower flags");

    for (uint256 i = 0; i < subjects.length; i++) {
      address subject = subjects[i];

      if (flags[subject]) {
        flags[subject] = false;
        emit FlagLowered(subject);
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

    if (previous != racAddress) {
      raisingAccessController = AccessControllerInterface(racAddress);

      emit RaisingAccessControllerUpdated(previous, racAddress);
    }
  }

  function setLoweringAccessController(
    address lacAddress
  )
    public
    override
    onlyOwner()
  {
    address previous = address(loweringAccessController);

    if (previous != lacAddress) {
      loweringAccessController = AccessControllerInterface(lacAddress);

      emit LoweringAccessControllerUpdated(previous, lacAddress);
    }
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

  function allowedToLowerFlags()
    private
    view
    returns (bool)
  {
    return msg.sender == owner ||
      loweringAccessController.hasAccess(msg.sender, msg.data);
  }

  function tryToRaiseFlag(address subject)
    private
  {
    if (!flags[subject]) {
      flags[subject] = true;
      emit FlagRaised(subject);
    }
  }

}
