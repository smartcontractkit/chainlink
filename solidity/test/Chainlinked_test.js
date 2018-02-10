'use strict';

require('./support/helpers.js')

contract('Chainlinked', () => {
  let Oracle = artifacts.require("Oracle.sol");
  let Chainlinked = artifacts.require("Chainlinked.sol");

  it("has a limited public interface", () => {
    checkPublicABI(Chainlinked, []);
  });
});
