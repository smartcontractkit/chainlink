import { CoordinatorFactory } from './generated/CoordinatorFactory'
import { MeanAggregatorFactory } from './generated/MeanAggregatorFactory'
import {
  registerPromiseHandler,
  DEVNET_ADDRESS,
  createProvider,
  deployContract,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'

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

  return {
    linkToken,
    coordinator,
    meanAggregator,
  }
}
