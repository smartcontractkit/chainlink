pragma solidity ^0.6.0;

abstract contract ENSResolver {
  function addr(bytes32 node) public virtual view returns (address);
}
