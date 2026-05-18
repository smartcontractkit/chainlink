import { BigNumber, BigNumberish } from '@ethersproject/bignumber'
import { ContractReceipt } from '@ethersproject/contracts'
import { assert, expect } from 'chai'

/**
 * Check that two big numbers are the same value.
 *
 * @param expected The expected value to match against
 * @param actual The actual value to match against the expected value
 * @param failureMessage Failure message to display if the actual value does not match the expected value.
 */
export function bigNumEquals(
  expected: BigNumberish,
  actual: BigNumberish,
  failureMessage?: string,
): void {
  const msg = failureMessage ? ': ' + failureMessage : ''
  assert(
    BigNumber.from(expected).eq(BigNumber.from(actual)),
    `BigNum (expected)${expected} is not (actual)${actual} ${msg}`,
  )
}

/**
 * Check that an evm operation reverts
 *
 * @param action The asynchronous action to execute, which should cause an evm revert.
 * @param msg The failure message to display if the action __does not__ throw
 */
export async function evmRevert(
  action: (() => Promise<any>) | Promise<any>,
  msg?: string,
) {
  if (msg) {
    await expect(action).to.be.revertedWith(msg)
  } else {
    await expect(action).to.be.reverted
  }
}

/**
 * Check that an evm operation reverts
 *
 * @param action The asynchronous action to execute, which should cause an evm revert.
 * @param contract The contract where the custom error is defined
 * @param msg The failure message to display if the action __does not__ throw
 */
export async function evmRevertCustomError(
  action: (() => Promise<any>) | Promise<any>,
  contract: any,
  msg?: string,
) {
  if (msg) {
    await expect(action).to.be.revertedWithCustomError(contract, msg)
  } else {
    await expect(action).to.be.reverted
  }
}

/**
 * Assert that an event doesnt exist
 *
 * @param max The maximum allowable gas difference
 * @param receipt1 The contract receipt to compare to
 * @param receipt2 The contract receipt with a gas difference
 */
export function gasDiffLessThan(
  max: number,
  receipt1: ContractReceipt,
  receipt2: ContractReceipt,
) {
  assert(receipt1, 'receipt1 is not present for gas comparison')
  assert(receipt2, 'receipt2 is not present for gas comparison')
  const diff = receipt2.gasUsed?.sub(receipt1.gasUsed || 0)
  assert.isAbove(max, diff?.toNumber() || Infinity)
}
