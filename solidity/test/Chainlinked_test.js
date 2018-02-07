'use strict';

require('./support/helpers.js')

contract('Chainlinked', () => {
  let Oracle = artifacts.require("./contracts/Oracle.sol");
  let Chainlinked = artifacts.require("./contracts/Chainlinked.sol");

  it("has a limited public interface", () => {
    checkPublicABI(Chainlinked, []);
  });
});
