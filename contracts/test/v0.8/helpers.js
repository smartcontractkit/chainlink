const { assert } = require("chai");

/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
function publicAbi(contract, expectedPublic) {
  const actualPublic = [];
  for (const m in contract.functions) {
    if (!m.includes("(")) {
      actualPublic.push(m);
    }
  }
  console.log(actualPublic);

  for (const method of actualPublic) {
    const index = expectedPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`);
  }

  for (const method of expectedPublic) {
    const index = actualPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is expected to be public`);
  }
}

module.exports = {
  publicAbi,
};
