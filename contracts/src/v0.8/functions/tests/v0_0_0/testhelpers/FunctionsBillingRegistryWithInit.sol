// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsBillingRegistry} from "../../../dev/v0_0_0/FunctionsBillingRegistry.sol";

contract FunctionsBillingRegistryWithInit is FunctionsBillingRegistry {
  constructor(address link, address linkEthFeed, address oracle) {
    initialize(link, linkEthFeed, oracle);
  }
}
