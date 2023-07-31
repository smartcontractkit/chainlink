// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouterBase} from "./IRouterBase.sol";
import {IOwnable} from "../../../../shared/interfaces/IOwnable.sol";

/**
 * @title Chainlink base Router interface with Ownable.
 */
interface IOwnableRouter is IOwnable, IRouterBase {}