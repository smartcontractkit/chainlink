import { ethers } from 'ethers'
import { createFundedWallet } from './wallet'
import { assert } from 'chai'

export interface Roles {
  defaultAccount: ethers.Wallet
  oracleNode: ethers.Wallet
  oracleNode1: ethers.Wallet
  oracleNode2: ethers.Wallet
  oracleNode3: ethers.Wallet
  stranger: ethers.Wallet
  consumer: ethers.Wallet
}

export interface Personas {
  Default: ethers.Wallet
  Neil: ethers.Wallet
  Ned: ethers.Wallet
  Nelly: ethers.Wallet
  Carol: ethers.Wallet
  Eddy: ethers.Wallet
}

interface RolesAndPersonasV2 {
  roles: Roles
  personas: Personas
}

/**
 * Generate roles and personas for tests along with their corrolated account addresses
 */
export async function initializeRolesAndPersonas(
  provider: ethers.providers.AsyncSendable,
): Promise<RolesAndPersonasV2> {
  const accounts = await Promise.all(
    Array(6)
      .fill(null)
      .map(async (_, i) => createFundedWallet(provider, i).then(w => w.wallet)),
  )

  const personas: Personas = {
    Default: accounts[0],
    Neil: accounts[1],
    Ned: accounts[2],
    Nelly: accounts[3],
    Carol: accounts[4],
    Eddy: accounts[5],
  }

  const roles: Roles = {
    defaultAccount: accounts[0],
    oracleNode: accounts[1],
    oracleNode1: accounts[1],
    oracleNode2: accounts[2],
    oracleNode3: accounts[3],
    stranger: accounts[4],
    consumer: accounts[5],
  }

  return { personas, roles }
}

type AsyncFunction = () => Promise<void>
export async function assertActionThrows(action: AsyncFunction) {
  let e: Error | undefined = undefined
  try {
    await action()
  } catch (error) {
    e = error
  }
  if (!e) {
    assert.exists(e, 'Expected an error to be raised')
    return
  }

  const { message } = e
  assert(message, 'Expected an error to contain a message')
  const invalidOpcode = message.includes('invalid opcode')
  const reverted = message.includes(
    'VM Exception while processing transaction: revert',
  )
  assert(
    invalidOpcode || reverted,
    'expected following error message to include "invalid JUMP" or ' +
      `"revert": "${message}"`,
  )
  // see https://github.com/ethereumjs/testrpc/issues/39
  // for why the "invalid JUMP" is the throw related error when using TestRPC
}

export function checkPublicABI(
  contract: ethers.Contract,
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
