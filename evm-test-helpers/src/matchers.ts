/**
 * @packageDocumentation
 *
 * This file contains a number of matcher functions meant to perform common assertions for
 * ethereum based tests. Specific assertion functions targeting chainlink smart contracts live in
 * their respective contracts/<contract>.ts file.
 */
import { assert } from 'chai'
import { ethers } from 'ethers'
import { ContractReceipt } from 'ethers/contract'
import { BigNumberish, EventDescription } from 'ethers/utils'
import { makeDebug } from './debug'
import { findEventIn } from './helpers'
const debug = makeDebug('helpers')

/**
 * Check that two big numbers are the same value.
 *
 * @param expected The expected value to match against
 * @param actual The actual value to match against the expected value
 * @param failureMessage Failure message to display if the actual value does not match the expected value.
 */
export function bigNum(
  expected: BigNumberish,
  actual: BigNumberish,
  failureMessage?: string,
): void {
  const msg = failureMessage ? ': ' + failureMessage : ''
  assert(
    ethers.utils.bigNumberify(expected).eq(ethers.utils.bigNumberify(actual)),
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
  const d = debug.extend('assertActionThrows')
  let e: Error | undefined = undefined

  try {
    if (typeof action === 'function') {
      await action()
    } else {
      await action
    }
  } catch (error) {
    e = error
  }
  d(e)
  if (!e) {
    assert.exists(e, 'Expected an error to be raised')
    return
  }

  assert(e.message, 'Expected an error to contain a message')

  const ERROR_MESSAGES = ['invalid opcode', 'revert']
  const hasErrored = ERROR_MESSAGES.some(msg => e?.message?.includes(msg))

  if (msg) {
    expect(e.message).toMatch(msg)
  }

  assert(
    hasErrored,
    `expected following error message to include ${ERROR_MESSAGES.join(
      ' or ',
    )}. Got: "${e.message}"`,
  )
}

/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
export function publicAbi(
  contract: ethers.Contract | ethers.ContractFactory,
  expectedPublic: string[],
) {
  const actualPublic = []
  for (const method of contract.interface.abi) {
    if (method.type === 'function') {
      actualPublic.push(method.name)
    }
  }

  for (const method of actualPublic) {
    const index = expectedPublic.indexOf(method)
    assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`)
  }

  for (const method of expectedPublic) {
    const index = actualPublic.indexOf(method)
    assert.isAtLeast(index, 0, `#${method} is expected to be public`)
  }
}

/**
 * Assert that an event exists
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
export function eventExists(
  receipt: ContractReceipt,
  eventDescription: EventDescription,
): ethers.Event {
  const event = findEventIn(receipt, eventDescription)
  if (!event) {
    throw Error(
      `Unable to find ${eventDescription.name} in transaction receipt`,
    )
  }

  return event
}

/**
 * Assert that an event doesnt exist
 *
 * @param receipt The contract receipt to find the event in
 * @param eventDescription A description of the event to search by
 */
export function eventDoesNotExist(
  receipt: ContractReceipt,
  eventDescription: EventDescription,
) {
  const event = findEventIn(receipt, eventDescription)
  if (event) {
    throw Error(
      `Found ${eventDescription.name} in transaction receipt, when expecting no instances`,
    )
  }
}
