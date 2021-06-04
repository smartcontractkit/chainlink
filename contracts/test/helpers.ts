import { BigNumber, Contract } from "ethers";
import { assert } from "chai";

export const constants = {
  ZERO_ADDRESS: "0x0000000000000000000000000000000000000000",
  ZERO_BYTES32: "0x0000000000000000000000000000000000000000000000000000000000000000",
  MAX_UINT256: BigNumber.from("2").pow(BigNumber.from("256")).sub(BigNumber.from("1")),
  MAX_INT256: BigNumber.from("2").pow(BigNumber.from("255")).sub(BigNumber.from("1")),
  MIN_INT256: BigNumber.from("2").pow(BigNumber.from("255")).mul(BigNumber.from("-1")),
};

/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
export function publicAbi(contract: Contract, expectedPublic: string[]) {
  const actualPublic = [];
  for (const m in contract.functions) {
    if (!m.includes("(")) {
      actualPublic.push(m);
    }
  }

  for (const method of actualPublic) {
    const index = expectedPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`);
  }

  for (const method of expectedPublic) {
    const index = actualPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is expected to be public`);
  }
}
