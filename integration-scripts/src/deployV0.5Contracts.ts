import { CoordinatorFactory } from '@chainlink/contracts/ethers/v0.5/CoordinatorFactory'
import { MeanAggregatorFactory } from '@chainlink/contracts/ethers/v0.5/MeanAggregatorFactory'
import {
  createProvider,
  deployContract,
  DEVNET_ADDRESS,
  registerPromiseHandler,
} from './common'
import { deployLinkTokenContract } from './deployLinkTokenContract'

async function main() {
  registerPromiseHandler()
  await deployContracts()
}
main()

export async function deployContracts() {
  const provider = createProvider()
  const signer = provider.getSigner(DEVNET_ADDRESS)

  const linkToken = await deployLinkTokenContract()

  const coordinator = await deployContract(
    { Factory: CoordinatorFactory, signer, name: 'Coordinator' },
    linkToken.address,
  )

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
