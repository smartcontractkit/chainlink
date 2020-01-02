import { ethers } from 'ethers'
import { generated as chainlink } from 'chainlinkv0.5'
import {
  registerPromiseHandler,
  DEVNET_ADDRESS,
  createProvider,
  deployContract,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'
const {
  CoordinatorFactory,
  MeanAggregatorFactory,
  PrepaidAggregatorFactory,
} = chainlink

async function main() {
  registerPromiseHandler()

  await deployContracts()
}
main()

// export async function deployContracts(provider: ethers.providers.JsonRpcProvider, DEVNET_ADDRESS: string) {
export async function deployContracts() {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  // deploy LINK token
  const linkToken = await deployLinkTokenContract()

  // deploy Coordinator
  const coordinator = await deployContract(
    { Factory: CoordinatorFactory, signer, name: 'Coordinator' },
    linkToken.address,
  )

  // deploy MeanAggregator
  const meanAggregator = await deployContract({
    Factory: MeanAggregatorFactory,
    signer,
    name: 'MeanAggregator',
  })

  // deploy PrepaidAggregator
  const prepaidAggregator = await deployContract(
    { Factory: PrepaidAggregatorFactory, signer, name: 'PrepaidAggregator' },
    linkToken.address,
    1,
    60,
    1,
    ethers.utils.formatBytes32String('USD/ETH'),
  )

  return {
    linkToken,
    coordinator,
    meanAggregator,
    prepaidAggregator,
  }
}
