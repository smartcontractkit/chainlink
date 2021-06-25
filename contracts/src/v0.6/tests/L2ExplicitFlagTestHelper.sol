// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../L2ExplicitFlag.sol";

contract L2ExplicitFlagTestHelper {
    L2ExplicitFlag public flag;

    constructor(address flagContract) public {
        flag = L2ExplicitFlag(flagContract);
    }

    function isRaised() external view returns (bool) {
        return flag.isRaised();
    }
}
