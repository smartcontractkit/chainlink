'use strict';

require('./support/helpers.js')

contract('ChainLinked', () => {
  let Oracle = artifacts.require("./contracts/Oracle.sol");
  let ChainLinked = artifacts.require("./contracts/ChainLinked.sol");

  it("has a limited public interface", () => {
    checkPublicABI(ChainLinked, []);
  });
});
