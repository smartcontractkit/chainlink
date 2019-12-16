import { ethers } from 'ethers'
import ganache from 'ganache-core'

/**
 * Create a test provider which uses an in-memory, in-process chain
 */
export function makeTestProvider(): ethers.providers.JsonRpcProvider {
  return new ethers.providers.Web3Provider(ganache.provider() as any)
}
