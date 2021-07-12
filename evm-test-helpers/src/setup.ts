/**
 * @packageDocumentation
 *
 * This file provides utility functions related to test setup, such as creating a test provider,
 * optimizing test times via snapshots, and making test accounts.
 */

import {
  RevertTraceSubprovider,
  SolCompilerArtifactAdapter,
  Web3ProviderEngine,
} from '@0x/sol-trace'
import {
  FakeGasEstimateSubprovider,
  GanacheSubprovider,
} from '@0x/subproviders'
import { ethers } from 'ethers'
import * as path from 'path'
import { makeDebug } from './debug'
import { createFundedWallet } from './wallet'

const debug = makeDebug('helpers')

/**
 * Create a test provider which uses an in-memory, in-process chain
 */
export function provider(): ethers.providers.JsonRpcProvider {
  const providerEngine = new Web3ProviderEngine()
  providerEngine.addProvider(new FakeGasEstimateSubprovider(5 * 10 ** 6)) // Ganache does a poor job of estimating gas, so just crank it up for testing.

  if (process.env.DEBUG) {
    debug('Debugging enabled, using sol-trace module...')
    const defaultFromAddress = ''
    const artifactAdapter = new SolCompilerArtifactAdapter(
      path.resolve('dist/artifacts'),
      path.resolve('contracts'),
    )
    const revertTraceSubprovider = new RevertTraceSubprovider(
      artifactAdapter,
      defaultFromAddress,
      true,
    )
    providerEngine.addProvider(revertTraceSubprovider)
  }

  providerEngine.addProvider(new GanacheSubprovider({ gasLimit: 8_000_000 }))
  providerEngine.start()

  return new ethers.providers.Web3Provider(providerEngine)
}

/**
 * This helper function allows us to make use of ganache snapshots,
 * which allows us to snapshot one state instance and revert back to it.
 *
 * This is used to memoize expensive setup calls typically found in beforeEach hooks when we
 * need to setup our state with contract deployments before running assertions.
 *
 * @param provider The provider that's used within the tests
 * @param cb The callback to execute that generates the state we want to snapshot
 */
export function snapshot(
  provider: ethers.providers.JsonRpcProvider,
  cb: () => Promise<void>,
) {
  if (process.env.DEBUG) {
    debug('Debugging enabled, snapshot mode disabled...')

    return cb
  }

  const d = debug.extend('memoizeDeploy')
  let hasDeployed = false
  let snapshotId = ''

  return async () => {
    if (!hasDeployed) {
      d('executing deployment..')
      await cb()

      d('snapshotting...')
      /* eslint-disable-next-line require-atomic-updates */
      snapshotId = await provider.send('evm_snapshot', undefined)
      d('snapshot id:%s', snapshotId)

      /* eslint-disable-next-line require-atomic-updates */
      hasDeployed = true
    } else {
      d('reverting to snapshot: %s', snapshotId)
      await provider.send('evm_revert', snapshotId)

      d('re-creating snapshot..')
      /* eslint-disable-next-line require-atomic-updates */
      snapshotId = await provider.send('evm_snapshot', undefined)
      d('recreated snapshot id:%s', snapshotId)
    }
  }
}

export interface Roles {
  defaultAccount: ethers.Wallet
  oracleNode: ethers.Wallet
  oracleNode1: ethers.Wallet
  oracleNode2: ethers.Wallet
  oracleNode3: ethers.Wallet
  oracleNode4: ethers.Wallet
  stranger: ethers.Wallet
  consumer: ethers.Wallet
}

export interface Personas {
  Default: ethers.Wallet
  Carol: ethers.Wallet
  Eddy: ethers.Wallet
  Nancy: ethers.Wallet
  Ned: ethers.Wallet
  Neil: ethers.Wallet
  Nelly: ethers.Wallet
  Norbert: ethers.Wallet
}

interface Users {
  roles: Roles
  personas: Personas
}

/**
 * Generate roles and personas for tests along with their correlated account addresses
 */
export async function users(
  provider: ethers.providers.JsonRpcProvider,
): Promise<Users> {
  const accounts = await Promise.all(
    Array(8)
      .fill(null)
      .map(async (_, i) =>
        createFundedWallet(provider, i).then((w) => w.wallet),
      ),
  )

  const personas: Personas = {
    Default: accounts[0],
    Neil: accounts[1],
    Ned: accounts[2],
    Nelly: accounts[3],
    Nancy: accounts[4],
    Norbert: accounts[5],
    Carol: accounts[6],
    Eddy: accounts[7],
  }

  const roles: Roles = {
    defaultAccount: accounts[0],
    oracleNode: accounts[1],
    oracleNode1: accounts[2],
    oracleNode2: accounts[3],
    oracleNode3: accounts[4],
    oracleNode4: accounts[5],
    stranger: accounts[6],
    consumer: accounts[7],
  }

  return { personas, roles }
}
