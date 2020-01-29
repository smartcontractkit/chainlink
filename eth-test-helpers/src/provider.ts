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
const debug = makeDebug('helpers')

/**
 * Create a test provider which uses an in-memory, in-process chain
 */
export function makeTestProvider(): ethers.providers.JsonRpcProvider {
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

  providerEngine.addProvider(new GanacheSubprovider({}))
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
export function useSnapshot(
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
