// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

import "./interfaces/L2FlagInterface.sol";
import "./interfaces/AccessControllerInterface.sol";
import "./SimpleReadAccessController.sol";

/**
 * @title The L2ExplicitFlag contract
 * @notice Explicit Flag that let readers know about
 * every feed state inside a L2 network
 */
contract L2ExplicitFlag is L2FlagInterface, SimpleReadAccessController {
    AccessControllerInterface public raisingAccessController;

    event RaisingAccessControllerUpdated(address indexed previous, address indexed current);
    event FlagLowered(bool _raised);
    event FlagRaised(bool _raised);
    bool private raised;

    /**
     * @param racAddress address for the raising access controller.
     */
    constructor(address racAddress) public {
        setRaisingAccessController(racAddress);
    }

    /**
     * @notice allows owner to change the access controller for raising flags.
     * @param racAddress new address for the raising access controller.
     */
    function setRaisingAccessController(address racAddress) public override onlyOwner() {
        address previous = address(raisingAccessController);

        if (previous != racAddress) {
            raisingAccessController = AccessControllerInterface(racAddress);

            emit RaisingAccessControllerUpdated(previous, racAddress);
        }
    }

    function raiseFlag() external override {
        require(allowedToChangeFlag(), "Not allowed to raise the flag");
        if (!raised) {
            raised = true;
            emit FlagRaised(raised);
        }
    }

    function lowerFlag() external override {
        require(allowedToChangeFlag(), "Not allowed to lower the flag");
        if (raised) {
            raised = false;
            emit FlagLowered(raised);
        }
    }

    function isRaised() external view override returns (bool) {
        return raised;
    }

    // PRIVATE
    function allowedToChangeFlag() private view returns (bool) {
        return msg.sender == owner || raisingAccessController.hasAccess(msg.sender, msg.data);
    }
}
