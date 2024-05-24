// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IFunctionsRouter} from "./IFunctionsRouter.sol";
import {IOwnable} from "../../../shared/interfaces/IOwnable.sol";

/// @title Chainlink Functions Router interface with Ownability.
interface IOwnableFunctionsRouter is IOwnable, IFunctionsRouter {}
