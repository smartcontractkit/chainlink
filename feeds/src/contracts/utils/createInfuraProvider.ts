import { ethers } from 'ethers'
import { JsonRpcProvider } from 'ethers/providers'
import { Config } from '../../config'
import { networkName, Networks } from '../../utils'

/**
 * Initialize the infura provider for the given network
 *
 * @param networkId The network id from the whitelist
 */
export function createInfuraProvider(
  networkId: Networks = Networks.MAINNET,
): JsonRpcProvider {
  const provider = new ethers.providers.JsonRpcProvider(
    Config.devProvider() ??
      `https://${networkName(networkId)}.infura.io/v3/${Config.infuraKey()}`,
  )
  provider.pollingInterval = 8000

  return provider
}
