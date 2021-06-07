import { BigNumber, BigNumberish } from "@ethersproject/bignumber";
import { assert, expect } from "chai";

/**
 * Check that two big numbers are the same value.
 *
 * @param expected The expected value to match against
 * @param actual The actual value to match against the expected value
 * @param failureMessage Failure message to display if the actual value does not match the expected value.
 */
export function bigNumEquals(expected: BigNumberish, actual: BigNumberish, failureMessage?: string): void {
  const msg = failureMessage ? ": " + failureMessage : "";
  assert(
    BigNumber.from(expected).eq(BigNumber.from(actual)),
    `BigNum (expected)${expected} is not (actual)${actual} ${msg}`,
  );
}

/**
 * Check that an evm operation reverts
 *
 * @param action The asynchronous action to execute, which should cause an evm revert.
 * @param msg The failure message to display if the action __does not__ throw
 */
export async function evmRevert(action: (() => Promise<any>) | Promise<any>, msg?: string) {
  if (msg) {
    await expect(action).to.be.revertedWith(msg);
  } else {
    await expect(action).to.be.reverted;
  }
}
