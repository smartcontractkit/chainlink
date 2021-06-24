pragma solidity ^0.4.24;

interface ENS {

  // Logged when the owner of a node assigns a new owner to a subnode.
  event NewOwner(bytes32 indexed node, bytes32 indexed label, address owner);

  // Logged when the owner of a node transfers ownership to a new account.
  event Transfer(bytes32 indexed node, address owner);

  // Logged when the resolver for a node changes.
  event NewResolver(bytes32 indexed node, address resolver);

  // Logged when the TTL of a node changes
  event NewTTL(bytes32 indexed node, uint64 ttl);


  function setSubnodeOwner(bytes32 node, bytes32 label, address owner) public;
  function setResolver(bytes32 node, address resolver) public;
  function setOwner(bytes32 node, address owner) public;
  function setTTL(bytes32 node, uint64 ttl) public;
  function owner(bytes32 node) public view returns (address);
  function resolver(bytes32 node) public view returns (address);
  function ttl(bytes32 node) public view returns (uint64);

}
