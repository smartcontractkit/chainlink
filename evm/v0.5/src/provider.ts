import { ethers } from 'ethers'
import {
  SolCompilerArtifactAdapter,
  Web3ProviderEngine,
  RevertTraceSubprovider,
} from '@0x/sol-trace'
import {
  FakeGasEstimateSubprovider,
  GanacheSubprovider,
} from '@0x/subproviders'
import * as path from 'path'

/**
 * Create a test provider which uses an in-memory, in-process chain
 */
export function makeTestProvider(): ethers.providers.JsonRpcProvider {
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

  const providerEngine = new Web3ProviderEngine()
  providerEngine.addProvider(new FakeGasEstimateSubprovider(4 * 10 ** 6)) // Ganache does a poor job of estimating gas, so just crank it up for testing.
  providerEngine.addProvider(revertTraceSubprovider)
  providerEngine.addProvider(new GanacheSubprovider({}))
  providerEngine.start()

  return new ethers.providers.Web3Provider(providerEngine)
}
