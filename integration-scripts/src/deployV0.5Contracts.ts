import { Coordinator__factory } from '@chainlink/contracts/ethers/v0.5/factories/Coordinator__factory'
import { MeanAggregator__factory } from '@chainlink/contracts/ethers/v0.5/factories/MeanAggregator__factory'
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
    { Factory: Coordinator__factory, signer, name: 'Coordinator' },
    linkToken.address,
  )

  const meanAggregator = await deployContract({
    Factory: MeanAggregator__factory,
    signer,
    name: 'MeanAggregator',
  })

  return {
    linkToken,
    coordinator,
    meanAggregator,
  }
}
