// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../dev/functions/FunctionsBillingRegistry.sol";

contract FunctionsBillingRegistryWithInit is FunctionsBillingRegistry_v0 {
  constructor(
    address link,
    address linkEthFeed,
    address oracle
  ) {
    initialize(link, linkEthFeed, oracle);
  }
}
